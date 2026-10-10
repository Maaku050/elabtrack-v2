package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/accounts"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountsRepository struct{ pool *pgxpool.Pool }

func NewAccountsRepository(pool *pgxpool.Pool) *AccountsRepository {
	return &AccountsRepository{pool: pool}
}
func (r *AccountsRepository) executor(ctx context.Context) executor {
	if tx, ok := database.TxFromContext(ctx); ok {
		return tx
	}
	return r.pool
}
func accountError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return shared.ErrNotFound
	}
	var p *pgconn.PgError
	if errors.As(err, &p) {
		if p.Code == "23505" {
			if p.ConstraintName == "borrower_profiles_student_id_key" {
				return d.ErrStudentIDExists
			}
			if p.ConstraintName == "users_email_key" {
				return user.ErrEmailAlreadyExists
			}
			return shared.ErrConflict
		}
		if p.Code == "23514" || p.Code == "23503" {
			return shared.ErrConflict
		}
	}
	return shared.Internal("accounts.store", err)
}

const managementColumns = `u.id,u.name,u.email,u.role,u.is_active,COALESCE(p.borrower_type,''),COALESCE(p.student_id,''),COALESCE(p.course,''),COALESCE(p.contact_number,''),u.activation_required,u.activation_delivery_status,u.created_at,u.updated_at`
const managementFrom = ` FROM users u LEFT JOIN borrower_profiles p ON p.user_id=u.id `

func scanManagement(s scanner) (a d.Record, err error) {
	err = s.Scan(&a.ID, &a.Name, &a.Email, &a.Role, &a.IsActive, &a.BorrowerType, &a.StudentID, &a.Course, &a.ContactNumber, &a.ActivationRequired, &a.DeliveryStatus, &a.CreatedAt, &a.UpdatedAt)
	a.Obligations = d.Obligations{Availability: "UNAVAILABLE"}
	return a, accountError(err)
}
func (r *AccountsRepository) LockIdentity(ctx context.Context) error {
	tx, ok := database.TxFromContext(ctx)
	if !ok {
		return shared.ErrInternal
	}
	_, e := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtextextended('elabtrack.account.identity.v1',0))`)
	return accountError(e)
}
func (r *AccountsRepository) LockAccounts(ctx context.Context, ids []uuid.UUID) error {
	tx, ok := database.TxFromContext(ctx)
	if !ok {
		return shared.ErrInternal
	}
	rows, e := tx.Query(ctx, `SELECT id FROM users WHERE id=ANY($1::uuid[]) ORDER BY id FOR UPDATE`, ids)
	if e != nil {
		return accountError(e)
	}
	defer rows.Close()
	for rows.Next() {
		var id uuid.UUID
		if e = rows.Scan(&id); e != nil {
			return accountError(e)
		}
	}
	return accountError(rows.Err())
}
func (r *AccountsRepository) Get(ctx context.Context, id uuid.UUID) (d.Record, error) {
	return scanManagement(r.executor(ctx).QueryRow(ctx, `SELECT `+managementColumns+managementFrom+`WHERE u.id=$1`, id))
}
func (r *AccountsRepository) FindIdentity(ctx context.Context, email, studentID string) (out []d.Record, err error) {
	rows, e := r.executor(ctx).Query(ctx, `SELECT `+managementColumns+managementFrom+`WHERE u.email=$1 OR ($2<>'' AND p.student_id=$2) ORDER BY u.id LIMIT 3`, email, studentID)
	if e != nil {
		return nil, accountError(e)
	}
	defer rows.Close()
	out = []d.Record{}
	for rows.Next() {
		a, e := scanManagement(rows)
		if e != nil {
			return nil, e
		}
		out = append(out, a)
	}
	return out, accountError(rows.Err())
}
func (r *AccountsRepository) List(ctx context.Context, f d.Filter) (out d.Page, err error) {
	where := `WHERE u.role='BORROWER'`
	if f.EligibleForIssuance {
		where += ` AND u.is_active AND NOT u.activation_required AND p.borrower_type IN ('STUDENT','FACULTY') AND EXISTS(SELECT 1 FROM terms_acceptances ta JOIN terms_publication tp ON tp.current_version_id=ta.terms_version_id WHERE tp.id=1 AND ta.user_id=u.id)`
	}
	if f.Staff {
		where = `WHERE u.role IN ('STAFF','ADMIN')`
	}
	args := []any{}
	add := func(value any) string { args = append(args, value); return fmt.Sprintf("$%d", len(args)) }
	if f.Search != "" {
		search := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(f.Search)
		arg := add("%" + search + "%")
		where += ` AND (u.name ILIKE ` + arg + ` OR u.email::text ILIKE ` + arg + ` OR p.student_id ILIKE ` + arg + `)`
	}
	if f.BorrowerType != "" {
		where += ` AND p.borrower_type=` + add(f.BorrowerType)
	}
	if f.Status == "ACTIVE" {
		where += ` AND u.is_active=TRUE`
	} else if f.Status == "INACTIVE" {
		where += ` AND u.is_active=FALSE`
	} else if f.Status == "PENDING" {
		where += ` AND u.activation_required=TRUE`
	}
	out.Items = []d.Record{}
	out.Page = f.Page
	out.PerPage = f.PerPage
	if e := r.executor(ctx).QueryRow(ctx, `SELECT COUNT(*)`+managementFrom+where, args...).Scan(&out.Total); e != nil {
		return out, accountError(e)
	}
	limit, offset := add(f.PerPage), add((f.Page-1)*f.PerPage)
	rows, e := r.executor(ctx).Query(ctx, `SELECT `+managementColumns+managementFrom+where+` ORDER BY u.created_at DESC,u.id LIMIT `+limit+` OFFSET `+offset, args...)
	if e != nil {
		return out, accountError(e)
	}
	defer rows.Close()
	for rows.Next() {
		a, e := scanManagement(rows)
		if e != nil {
			return out, e
		}
		out.Items = append(out.Items, a)
	}
	return out, accountError(rows.Err())
}
func (r *AccountsRepository) Insert(ctx context.Context, a d.Record) error {
	_, e := r.executor(ctx).Exec(ctx, `INSERT INTO users(id,name,email,password,role,is_active,activation_required,activation_delivery_status,created_at,updated_at) VALUES($1,$2,$3,'!activation-pending',$4,TRUE,TRUE,'PENDING',$5,$5)`, a.ID, a.Name, a.Email, a.Role, a.CreatedAt)
	if e != nil {
		return accountError(e)
	}
	if a.Role == user.RoleBorrower {
		var sid any
		if a.BorrowerType == "STUDENT" {
			sid = a.StudentID
		}
		_, e = r.executor(ctx).Exec(ctx, `INSERT INTO borrower_profiles(user_id,borrower_type,student_id,course,contact_number) VALUES($1,$2,$3,$4,$5)`, a.ID, a.BorrowerType, sid, a.Course, a.ContactNumber)
	}
	return accountError(e)
}
func (r *AccountsRepository) UpdateProfile(ctx context.Context, a d.Record) error {
	_, e := r.executor(ctx).Exec(ctx, `UPDATE users SET name=$2,updated_at=statement_timestamp() WHERE id=$1`, a.ID, a.Name)
	if e != nil {
		return accountError(e)
	}
	_, e = r.executor(ctx).Exec(ctx, `UPDATE borrower_profiles SET course=$2,contact_number=$3 WHERE user_id=$1`, a.ID, a.Course, a.ContactNumber)
	return accountError(e)
}
func (r *AccountsRepository) SetActive(ctx context.Context, id uuid.UUID, active bool) error {
	_, e := r.executor(ctx).Exec(ctx, `UPDATE users SET is_active=$2,updated_at=statement_timestamp() WHERE id=$1`, id, active)
	if e != nil {
		return accountError(e)
	}
	if !active {
		_, e = r.executor(ctx).Exec(ctx, `UPDATE refresh_tokens SET revoked_at=COALESCE(revoked_at,statement_timestamp()) WHERE user_id=$1 AND revoked_at IS NULL`, id)
	}
	return accountError(e)
}
func (r *AccountsRepository) SetPassword(ctx context.Context, id uuid.UUID, hashed string) error {
	tag, e := r.executor(ctx).Exec(ctx, `UPDATE users SET password=$2,activation_required=FALSE,updated_at=statement_timestamp() WHERE id=$1 AND is_active AND activation_required`, id, hashed)
	if e != nil {
		return accountError(e)
	}
	if tag.RowsAffected() != 1 {
		return d.ErrActivationInvalid
	}
	return nil
}
func (r *AccountsRepository) PutToken(ctx context.Context, t d.Token) error {
	_, e := r.executor(ctx).Exec(ctx, `INSERT INTO account_activation_tokens(token_hash,user_id,created_at,expires_at) VALUES($1,$2,$3,$4)`, t.Hash, t.AccountID, t.CreatedAt, t.ExpiresAt)
	if e != nil {
		return accountError(e)
	}
	_, e = r.executor(ctx).Exec(ctx, `UPDATE users SET activation_delivery_status='PENDING',updated_at=statement_timestamp() WHERE id=$1`, t.AccountID)
	return accountError(e)
}
func (r *AccountsRepository) FindToken(ctx context.Context, hash string, lock bool) (t d.Token, err error) {
	q := `SELECT token_hash,user_id,created_at,expires_at,invalidated FROM account_activation_tokens WHERE token_hash=$1`
	arg := hash
	if strings.HasPrefix(hash, "account:") {
		q = `SELECT token_hash,user_id,created_at,expires_at,invalidated FROM account_activation_tokens WHERE user_id=$1 ORDER BY created_at DESC,token_hash LIMIT 1`
		arg = strings.TrimPrefix(hash, "account:")
	}
	if lock {
		q += ` FOR UPDATE`
	}
	err = r.executor(ctx).QueryRow(ctx, q, arg).Scan(&t.Hash, &t.AccountID, &t.CreatedAt, &t.ExpiresAt, &t.Invalid)
	return t, accountError(err)
}
func (r *AccountsRepository) InvalidateTokens(ctx context.Context, id uuid.UUID) error {
	_, e := r.executor(ctx).Exec(ctx, `UPDATE account_activation_tokens SET invalidated=TRUE WHERE user_id=$1 AND NOT invalidated`, id)
	return accountError(e)
}
func (r *AccountsRepository) SetDelivery(ctx context.Context, id uuid.UUID, hash, state, message string) error {
	if len(message) > 256 {
		message = ""
	}
	_, e := r.executor(ctx).Exec(ctx, `UPDATE account_activation_tokens SET delivery_status=$3,provider_message_id=$4 WHERE user_id=$1 AND token_hash=$2 AND NOT invalidated`, id, hash, state, message)
	if e != nil {
		return accountError(e)
	}
	_, e = r.executor(ctx).Exec(ctx, `UPDATE users SET activation_delivery_status=$2,updated_at=statement_timestamp() WHERE id=$1`, id, state)
	return accountError(e)
}
func (r *AccountsRepository) ReadReceipt(ctx context.Context, actor uuid.UUID, op, key, hash string) (result []byte, err error) {
	var stored string
	err = r.executor(ctx).QueryRow(ctx, `SELECT payload_hash,result FROM account_operation_receipts WHERE actor_id=$1 AND operation=$2 AND command_key=$3`, actor, op, key).Scan(&stored, &result)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, accountError(err)
	}
	if stored != hash {
		return nil, shared.ErrConflict
	}
	return result, nil
}
func (r *AccountsRepository) WriteReceipt(ctx context.Context, actor uuid.UUID, op, key, hash string, result []byte) error {
	_, e := r.executor(ctx).Exec(ctx, `INSERT INTO account_operation_receipts(actor_id,operation,command_key,payload_hash,result) VALUES($1,$2,$3,$4,$5)`, actor, op, key, hash, json.RawMessage(result))
	return accountError(e)
}
func (r *AccountsRepository) AddAudit(ctx context.Context, actor, id uuid.UUID, action string) error {
	_, e := r.executor(ctx).Exec(ctx, `INSERT INTO account_audit_events(id,actor_id,account_id,action) VALUES($1,$2,$3,$4)`, uuid.New(), actor, id, action)
	return accountError(e)
}
func (r *AccountsRepository) Audits(ctx context.Context, id uuid.UUID, page int) (out []d.Audit, err error) {
	rows, e := r.executor(ctx).Query(ctx, `SELECT id,actor_id,account_id,action,occurred_at FROM account_audit_events WHERE account_id=$1 ORDER BY occurred_at DESC,id LIMIT 25 OFFSET $2`, id, (page-1)*25)
	if e != nil {
		return nil, accountError(e)
	}
	defer rows.Close()
	out = []d.Audit{}
	for rows.Next() {
		var a d.Audit
		if e = rows.Scan(&a.ID, &a.ActorID, &a.AccountID, &a.Action, &a.At); e != nil {
			return nil, accountError(e)
		}
		out = append(out, a)
	}
	return out, accountError(rows.Err())
}
func (r *AccountsRepository) PutBatch(ctx context.Context, b d.Batch) error {
	raw, e := json.Marshal(b)
	if e != nil {
		return shared.ErrInternal
	}
	_, e = r.executor(ctx).Exec(ctx, `INSERT INTO account_roster_batches(id,actor_id,operation,preview,expires_at) VALUES($1,$2,$3,$4,$5)`, b.ID, b.ActorID, b.Operation, json.RawMessage(raw), b.ExpiresAt)
	return accountError(e)
}
func (r *AccountsRepository) GetBatch(ctx context.Context, id uuid.UUID, lock bool) (b d.Batch, err error) {
	q := `SELECT actor_id,COALESCE(result,preview) FROM account_roster_batches WHERE id=$1`
	if lock {
		q += ` FOR UPDATE`
	}
	var actor uuid.UUID
	var raw []byte
	if err = r.executor(ctx).QueryRow(ctx, q, id).Scan(&actor, &raw); err != nil {
		return b, accountError(err)
	}
	if json.Unmarshal(raw, &b) != nil {
		return b, shared.ErrInternal
	}
	b.ActorID = actor
	return b, nil
}
func (r *AccountsRepository) ConfirmBatch(ctx context.Context, b d.Batch) error {
	raw, e := json.Marshal(b)
	if e != nil {
		return shared.ErrInternal
	}
	tag, e := r.executor(ctx).Exec(ctx, `UPDATE account_roster_batches SET result=$2 WHERE id=$1 AND result IS NULL`, b.ID, json.RawMessage(raw))
	if e != nil {
		return accountError(e)
	}
	if tag.RowsAffected() != 1 {
		return shared.ErrConflict
	}
	return nil
}

var _ d.Repository = (*AccountsRepository)(nil)
