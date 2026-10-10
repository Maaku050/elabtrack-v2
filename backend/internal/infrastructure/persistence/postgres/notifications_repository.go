package postgres

import (
	"context"
	"encoding/json"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/notifications"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationsRepository struct{ pool *pgxpool.Pool }

func NewNotificationsRepository(p *pgxpool.Pool) *NotificationsRepository {
	return &NotificationsRepository{p}
}
func (r *NotificationsRepository) executor(c context.Context) executor {
	if tx, ok := database.TxFromContext(c); ok {
		return tx
	}
	return r.pool
}
func (r *NotificationsRepository) TryWorkerLock(c context.Context) (v bool, e error) {
	e = r.executor(c).QueryRow(c, `SELECT pg_try_advisory_xact_lock(1853188,9)`).Scan(&v)
	return v, accountError(e)
}
func (r *NotificationsRepository) Pending(c context.Context, limit int, lead int64) (events []d.Event, e error) {
	events = []d.Event{}
	rows, e := r.executor(c).Query(c, `WITH tick AS MATERIALIZED(SELECT clock_timestamp() at), sources AS (
 SELECT 'EVENT:'||e.id AS event_key,e.borrowing_id,b.borrower_id,e.actor_id,e.kind,e.occurred_at AS at FROM borrowing_events e JOIN borrowings b ON b.id=e.borrowing_id
 UNION ALL SELECT 'DAMAGED:'||l.id,l.borrowing_id,b.borrower_id,l.actor_id,'DAMAGED',l.occurred_at FROM return_lines l JOIN borrowings b ON b.id=l.borrowing_id WHERE l.damaged>0
 UNION ALL SELECT 'LOST:'||l.id,l.borrowing_id,b.borrower_id,l.actor_id,'LOST',l.occurred_at FROM return_lines l JOIN borrowings b ON b.id=l.borrowing_id WHERE l.lost>0
 UNION ALL SELECT 'REPLACEMENT_REQUIRED:'||o.id,o.borrowing_id,b.borrower_id,l.actor_id,'REPLACEMENT_REQUIRED',l.occurred_at FROM replacement_obligations o JOIN return_lines l ON l.id=o.return_line_id JOIN borrowings b ON b.id=o.borrowing_id
 UNION ALL SELECT 'DUE_SOON:'||b.id,b.id,b.borrower_id,NULL::uuid,'DUE_SOON',b.due_at FROM borrowings b,tick WHERE $2::bigint>0 AND b.status='CHECKED_OUT' AND b.due_at>=tick.at AND b.due_at<=tick.at+make_interval(secs=>$2::double precision)
 UNION ALL SELECT 'OVERDUE:'||b.id,b.id,b.borrower_id,NULL::uuid,'OVERDUE',b.due_at FROM borrowings b,tick WHERE b.status='CHECKED_OUT' AND b.due_at<tick.at
 UNION ALL SELECT 'FINE_ASSESSED:'||b.id,b.id,b.borrower_id,NULL::uuid,'FINE_ASSESSED',b.due_at FROM borrowings b,tick WHERE b.status IN ('CHECKED_OUT','COMPLETED') AND COALESCE(b.fine_final_minor,borrowing_fine_amount(b.due_at,tick.at))>0
 ) SELECT s.event_key,s.borrowing_id,s.borrower_id,s.actor_id,s.kind,s.at FROM sources s LEFT JOIN notification_dispatch n ON n.event_key=s.event_key WHERE n.event_key IS NULL ORDER BY s.at,s.event_key LIMIT $1`, limit, lead)
	if e != nil {
		return events, accountError(e)
	}
	defer rows.Close()
	for rows.Next() {
		var event d.Event
		if e = rows.Scan(&event.Key, &event.BorrowingID, &event.BorrowerID, &event.ActorID, &event.Kind, &event.At); e != nil {
			return nil, accountError(e)
		}
		events = append(events, event)
	}
	return events, accountError(rows.Err())
}
func (r *NotificationsRepository) Emit(c context.Context, event d.Event, message d.Message) (bool, error) {
	tag, e := r.executor(c).Exec(c, `INSERT INTO notification_dispatch(event_key,borrowing_id,kind)VALUES($1,$2,$3) ON CONFLICT(event_key)DO NOTHING`, event.Key, event.BorrowingID, event.Kind)
	if e != nil {
		return false, accountError(e)
	}
	if tag.RowsAffected() == 0 {
		return false, nil
	}
	// Recipient identity is derived from committed loan/event relationships, never caller input.
	_, e = r.executor(c).Exec(c, `INSERT INTO notifications(id,recipient_id,event_key,borrowing_id,scope,kind,title,body)
 SELECT gen_random_uuid(),u.id,$1,$2,CASE WHEN u.id=$3 THEN 'BORROWER' ELSE 'OPERATIONS' END,$4,$5,$6 FROM users u
 WHERE u.id=$3 OR (u.is_active AND NOT u.activation_required AND u.role IN ('STAFF','ADMIN') AND ($7 OR u.id=$8::uuid)) ON CONFLICT(event_key,recipient_id)DO NOTHING`, event.Key, event.BorrowingID, event.BorrowerID, event.Kind, message.Title, message.Body, message.Operations, event.ActorID)
	return e == nil, accountError(e)
}

const notificationScope = `recipient_id=$1 AND ((scope='BORROWER' AND $2='BORROWER' AND borrowing_id IN(SELECT id FROM borrowings WHERE borrower_id=$1)) OR(scope='OPERATIONS' AND $2 IN ('STAFF','ADMIN')))`

func (r *NotificationsRepository) List(c context.Context, id uuid.UUID, role user.Role, f d.Filter) (p d.Page, e error) {
	p.Items = []d.Notification{}
	p.Page = f.Page
	p.PerPage = f.PerPage
	var raw []byte
	where := ` FROM notifications WHERE ` + notificationScope + ` AND ($3='' OR ($3='read')=(read_at IS NOT NULL))`
	e = r.executor(c).QueryRow(c, `SELECT (SELECT count(*)`+where+`),(SELECT count(*) FROM notifications WHERE `+notificationScope+` AND read_at IS NULL),COALESCE((SELECT jsonb_agg(to_jsonb(n) ORDER BY n.created_at DESC,n.id DESC) FROM(SELECT id,borrowing_id,scope,kind,title,body,created_at,read_at`+where+` ORDER BY created_at DESC,id DESC LIMIT $4 OFFSET $5)n),'[]'::jsonb)`, id, string(role), f.Read, f.PerPage, (f.Page-1)*f.PerPage).Scan(&p.Total, &p.Unread, &raw)
	if e != nil {
		return p, accountError(e)
	}
	e = json.Unmarshal(raw, &p.Items)
	return p, e
}
func (r *NotificationsRepository) Count(c context.Context, id uuid.UUID, role user.Role) (n int, e error) {
	e = r.executor(c).QueryRow(c, `SELECT count(*) FROM notifications WHERE `+notificationScope+` AND read_at IS NULL`, id, string(role)).Scan(&n)
	return n, accountError(e)
}
func (r *NotificationsRepository) Mark(c context.Context, id uuid.UUID, role user.Role, notification uuid.UUID, read bool) (v d.Notification, e error) {
	e = r.executor(c).QueryRow(c, `UPDATE notifications SET read_at=CASE WHEN $4 THEN COALESCE(read_at,clock_timestamp()) ELSE NULL END WHERE `+notificationScope+` AND id=$3 RETURNING id,borrowing_id,scope,kind,title,body,created_at,read_at`, id, string(role), notification, read).Scan(&v.ID, &v.BorrowingID, &v.Scope, &v.Kind, &v.Title, &v.Body, &v.CreatedAt, &v.ReadAt)
	return v, accountError(e)
}

var _ d.Repository = (*NotificationsRepository)(nil)

func (r *NotificationsRepository) MarkAll(c context.Context, id uuid.UUID, role user.Role) (int64, error) {
	tag, e := r.executor(c).Exec(c, `UPDATE notifications SET read_at=statement_timestamp() WHERE `+notificationScope+` AND read_at IS NULL`, id, string(role))
	return tag.RowsAffected(), accountError(e)
}
