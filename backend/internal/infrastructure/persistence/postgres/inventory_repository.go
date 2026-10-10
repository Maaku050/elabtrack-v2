package postgres

import (
	"context"
	"encoding/json"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/inventory"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
)

type InventoryRepository struct{ pool *pgxpool.Pool }

func NewInventoryRepository(p *pgxpool.Pool) *InventoryRepository { return &InventoryRepository{p} }
func (r *InventoryRepository) executor(c context.Context) executor {
	if tx, ok := database.TxFromContext(c); ok {
		return tx
	}
	return r.pool
}

const equipmentColumns = `e.id,e.name,e.description,e.category_id,COALESCE(c.name,''),e.status,e.available,e.reserved,e.checked_out,e.damaged_held,e.total_tracked,e.stock_sequence,e.metadata_version,e.image_id,e.created_at,e.updated_at`

func scanEquipment(s scanner) (v d.Equipment, e error) {
	e = s.Scan(&v.ID, &v.Name, &v.Description, &v.CategoryID, &v.CategoryName, &v.Status, &v.Stock.Available, &v.Stock.Reserved, &v.Stock.CheckedOut, &v.Stock.DamagedHeld, &v.Stock.Total, &v.Sequence, &v.Version, &v.ImageID, &v.CreatedAt, &v.UpdatedAt)
	return v, accountError(e)
}
func (r *InventoryRepository) LockEquipment(c context.Context, id uuid.UUID) error {
	var v uuid.UUID
	return accountError(r.executor(c).QueryRow(c, `SELECT id FROM equipment WHERE id=$1 FOR UPDATE`, id).Scan(&v))
}
func (r *InventoryRepository) Get(c context.Context, id uuid.UUID) (d.Equipment, error) {
	return scanEquipment(r.executor(c).QueryRow(c, `SELECT `+equipmentColumns+` FROM equipment e LEFT JOIN equipment_categories c ON c.id=e.category_id WHERE e.id=$1`, id))
}
func (r *InventoryRepository) List(c context.Context, f d.Filter) (p d.Page, err error) {
	p.Items = []d.Equipment{}
	p.Page = f.Page
	p.PerPage = f.PerPage
	search := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(f.Search)
	args := []any{f.Status, f.CategoryID, "%" + search + "%", f.AvailableOnly}
	where := ` FROM equipment e LEFT JOIN equipment_categories c ON c.id=e.category_id WHERE ($1='' OR e.status=$1) AND ($2::uuid IS NULL OR e.category_id=$2) AND (e.name ILIKE $3 OR e.description ILIKE $3) AND (NOT $4::bool OR (e.status='ACTIVE' AND e.available>0))`
	e := r.executor(c).QueryRow(c, `SELECT count(*),COALESCE(sum(available),0)::bigint,COALESCE(sum(reserved),0)::bigint,COALESCE(sum(checked_out),0)::bigint,COALESCE(sum(damaged_held),0)::bigint,COALESCE(sum(total_tracked),0)::bigint`+where, args...).Scan(&p.Total, &p.Totals.Available, &p.Totals.Reserved, &p.Totals.CheckedOut, &p.Totals.DamagedHeld, &p.Totals.Total)
	if e != nil {
		return p, accountError(e)
	}
	order := `lower(e.name),e.id`
	if f.Sort == "available" {
		order = `e.available DESC,lower(e.name),e.id`
	}
	args = append(args, f.PerPage, (f.Page-1)*f.PerPage)
	rows, e := r.executor(c).Query(c, `SELECT `+equipmentColumns+where+` ORDER BY `+order+` LIMIT $5 OFFSET $6`, args...)
	if e != nil {
		return p, accountError(e)
	}
	defer rows.Close()
	for rows.Next() {
		v, e := scanEquipment(rows)
		if e != nil {
			return p, e
		}
		p.Items = append(p.Items, v)
	}
	return p, accountError(rows.Err())
}
func (r *InventoryRepository) Insert(c context.Context, v d.Equipment) error {
	_, e := r.executor(c).Exec(c, `INSERT INTO equipment(id,name,description,category_id,status,created_at,updated_at) VALUES($1,$2,$3,$4,$5,$6,$7)`, v.ID, v.Name, v.Description, v.CategoryID, v.Status, v.CreatedAt, v.UpdatedAt)
	return accountError(e)
}
func (r *InventoryRepository) UpdateMetadata(c context.Context, v d.Equipment) error {
	_, e := r.executor(c).Exec(c, `UPDATE equipment SET name=$2,description=$3,category_id=$4,status=$5,metadata_version=$6,image_id=$7,updated_at=clock_timestamp() WHERE id=$1`, v.ID, v.Name, v.Description, v.CategoryID, v.Status, v.Version, v.ImageID)
	return accountError(e)
}
func (r *InventoryRepository) UpdateStock(c context.Context, v d.Equipment) error {
	_, e := r.executor(c).Exec(c, `UPDATE equipment SET available=$2,total_tracked=$3,stock_sequence=$4,updated_at=clock_timestamp() WHERE id=$1`, v.ID, v.Stock.Available, v.Stock.Total, v.Sequence)
	return accountError(e)
}
func (r *InventoryRepository) Categories(c context.Context, active bool, page int) (out []d.Category, err error) {
	out = []d.Category{}
	rows, e := r.executor(c).Query(c, `SELECT id,name,is_active,version FROM equipment_categories WHERE (NOT $1::bool OR is_active) ORDER BY lower(name),id LIMIT 100 OFFSET $2`, active, (page-1)*100)
	if e != nil {
		return nil, accountError(e)
	}
	defer rows.Close()
	for rows.Next() {
		var v d.Category
		if e = rows.Scan(&v.ID, &v.Name, &v.Active, &v.Version); e != nil {
			return nil, accountError(e)
		}
		out = append(out, v)
	}
	return out, accountError(rows.Err())
}
func (r *InventoryRepository) GetCategory(c context.Context, id uuid.UUID, lock bool) (v d.Category, err error) {
	suffix := ""
	if lock {
		suffix = " FOR UPDATE"
	}
	e := r.executor(c).QueryRow(c, `SELECT id,name,is_active,version FROM equipment_categories WHERE id=$1`+suffix, id).Scan(&v.ID, &v.Name, &v.Active, &v.Version)
	return v, accountError(e)
}
func (r *InventoryRepository) PutCategory(c context.Context, v d.Category, create bool) error {
	var e error
	if create {
		_, e = r.executor(c).Exec(c, `INSERT INTO equipment_categories(id,name,is_active) VALUES($1,$2,$3)`, v.ID, v.Name, v.Active)
	} else {
		_, e = r.executor(c).Exec(c, `UPDATE equipment_categories SET name=$2,is_active=$3,version=$4,updated_at=clock_timestamp() WHERE id=$1`, v.ID, v.Name, v.Active, v.Version)
	}
	return accountError(e)
}
func (r *InventoryRepository) Movement(c context.Context, v d.Movement) error {
	_, e := r.executor(c).Exec(c, `INSERT INTO inventory_movements(id,equipment_id,actor_id,sequence,kind,delta_available,delta_reserved,delta_checked_out,delta_damaged_held,delta_total,after_available,after_reserved,after_checked_out,after_damaged_held,after_total,reason) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`, v.ID, v.EquipmentID, v.ActorID, v.Sequence, v.Kind, v.Delta.Available, v.Delta.Reserved, v.Delta.CheckedOut, v.Delta.DamagedHeld, v.Delta.Total, v.After.Available, v.After.Reserved, v.After.CheckedOut, v.After.DamagedHeld, v.After.Total, v.Reason)
	return accountError(e)
}
func (r *InventoryRepository) Movements(c context.Context, id uuid.UUID, page int) (out []d.Movement, err error) {
	out = []d.Movement{}
	rows, e := r.executor(c).Query(c, `SELECT id,equipment_id,actor_id,sequence,kind,delta_available,delta_reserved,delta_checked_out,delta_damaged_held,delta_total,after_available,after_reserved,after_checked_out,after_damaged_held,after_total,reason,created_at FROM inventory_movements WHERE equipment_id=$1 ORDER BY sequence DESC LIMIT 25 OFFSET $2`, id, (page-1)*25)
	if e != nil {
		return nil, accountError(e)
	}
	defer rows.Close()
	for rows.Next() {
		var v d.Movement
		e = rows.Scan(&v.ID, &v.EquipmentID, &v.ActorID, &v.Sequence, &v.Kind, &v.Delta.Available, &v.Delta.Reserved, &v.Delta.CheckedOut, &v.Delta.DamagedHeld, &v.Delta.Total, &v.After.Available, &v.After.Reserved, &v.After.CheckedOut, &v.After.DamagedHeld, &v.After.Total, &v.Reason, &v.At)
		if e != nil {
			return nil, accountError(e)
		}
		out = append(out, v)
	}
	return out, accountError(rows.Err())
}
func (r *InventoryRepository) Audit(c context.Context, actor uuid.UUID, equip, cat *uuid.UUID, action string, before, after any) error {
	b, e := json.Marshal(before)
	if e != nil {
		return e
	}
	a, e := json.Marshal(after)
	if e != nil {
		return e
	}
	_, e = r.executor(c).Exec(c, `INSERT INTO inventory_audit_events(id,actor_id,equipment_id,category_id,action,before_data,after_data) VALUES($1,$2,$3,$4,$5,$6,$7)`, uuid.New(), actor, equip, cat, action, b, a)
	return accountError(e)
}
func (r *InventoryRepository) ReadReceipt(c context.Context, actor uuid.UUID, op, key, hash string) ([]byte, error) {
	var h string
	var data []byte
	e := r.executor(c).QueryRow(c, `SELECT payload_hash,result FROM inventory_operation_receipts WHERE actor_id=$1 AND operation=$2 AND key=$3`, actor, op, key).Scan(&h, &data)
	if e != nil {
		return nil, accountError(e)
	}
	if h != hash {
		return nil, shared.ErrConflict
	}
	return data, nil
}
func (r *InventoryRepository) WriteReceipt(c context.Context, actor uuid.UUID, op, key, hash string, data []byte) error {
	_, e := r.executor(c).Exec(c, `INSERT INTO inventory_operation_receipts(actor_id,operation,key,payload_hash,result) VALUES($1,$2,$3,$4,$5)`, actor, op, key, hash, data)
	return accountError(e)
}

// Equipment lock is held on archive writes. These liability reads acquire no
// borrowing locks; a competing reservation/issue must lock the equipment too.
func (r *InventoryRepository) ArchiveSafety(c context.Context, ids ...uuid.UUID) (string, error) {
	var borrowing, replacement bool
	e := r.executor(c).QueryRow(c, `SELECT to_regclass('public.borrowings') IS NOT NULL,to_regclass('public.replacement_obligations') IS NOT NULL`).Scan(&borrowing, &replacement)
	if e != nil {
		return "UNAVAILABLE", accountError(e)
	}
	if replacement {
		if len(ids) != 1 {
			return "UNAVAILABLE", nil
		}
		var complete bool
		e = r.executor(c).QueryRow(c, `SELECT to_regclass('public.return_lines') IS NOT NULL AND to_regclass('public.replacement_acceptances') IS NOT NULL`).Scan(&complete)
		if e != nil {
			return "UNAVAILABLE", accountError(e)
		}
		if !complete {
			return "UNAVAILABLE", nil
		}
		var blocked bool
		e = r.executor(c).QueryRow(c, `SELECT EXISTS(SELECT 1 FROM borrowing_items i JOIN borrowings b ON b.id=i.borrowing_id WHERE i.equipment_id=$1 AND (i.reserved_quantity>0 OR (b.status='CHECKED_OUT' AND i.issued_quantity>COALESCE((SELECT sum(good+damaged+lost) FROM return_lines l WHERE l.item_id=i.id),0)))) OR EXISTS(SELECT 1 FROM replacement_obligations o WHERE o.equipment_id=$1 AND o.required>COALESCE((SELECT sum(quantity) FROM replacement_acceptances a WHERE a.obligation_id=o.id),0))`, ids[0]).Scan(&blocked)
		if e != nil {
			return "UNAVAILABLE", accountError(e)
		}
		if blocked {
			return "BLOCKED", nil
		}
		return "CLEAR", nil
	}
	if !borrowing {
		return "NOT_INSTALLED", nil
	}
	if len(ids) != 1 {
		return "UNAVAILABLE", nil
	}
	var installed bool
	e = r.executor(c).QueryRow(c, `SELECT to_regclass('public.borrowing_items') IS NOT NULL AND to_regclass('public.borrowing_events') IS NOT NULL AND EXISTS(SELECT 1 FROM information_schema.columns WHERE table_schema='public' AND table_name='borrowings' AND column_name='entry_path')`).Scan(&installed)
	if e != nil {
		return "UNAVAILABLE", accountError(e)
	}
	if !installed {
		return "UNAVAILABLE", nil
	}
	var outstanding bool
	e = r.executor(c).QueryRow(c, `SELECT EXISTS(SELECT 1 FROM borrowing_items i JOIN borrowings b ON b.id=i.borrowing_id WHERE i.equipment_id=$1 AND b.status IN ('PENDING','CHECKED_OUT') AND (i.reserved_quantity>0 OR i.issued_quantity>0))`, ids[0]).Scan(&outstanding)
	if e != nil {
		return "UNAVAILABLE", accountError(e)
	}
	if outstanding {
		return "BLOCKED", nil
	}
	return "CLEAR", nil
}
func (r *InventoryRepository) PutImage(c context.Context, v d.Image) error {
	_, e := r.executor(c).Exec(c, `INSERT INTO equipment_images(id,equipment_id,png,sha256,width,height,actor_id) VALUES($1,$2,$3,$4,$5,$6,$7)`, v.ID, v.EquipmentID, v.PNG, v.Hash, v.Width, v.Height, v.ActorID)
	return accountError(e)
}
func (r *InventoryRepository) GetImage(c context.Context, id, im uuid.UUID) (v d.Image, err error) {
	e := r.executor(c).QueryRow(c, `SELECT id,equipment_id,png,sha256,width,height,actor_id FROM equipment_images WHERE id=$1 AND equipment_id=$2`, im, id).Scan(&v.ID, &v.EquipmentID, &v.PNG, &v.Hash, &v.Width, &v.Height, &v.ActorID)
	return v, accountError(e)
}

var _ d.Repository = (*InventoryRepository)(nil)
