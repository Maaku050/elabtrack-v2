package inventory

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/Maaku050/elabtrack-v2/backend/internal/application"
	a "github.com/Maaku050/elabtrack-v2/backend/internal/domain/accounts"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/inventory"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/google/uuid"
	"strings"
	"time"
)

type Service struct {
	repo     d.Repository
	accounts a.Repository
	tx       application.Transactions
	images   d.ImageValidator
}

func NewService(r d.Repository, a a.Repository, t application.Transactions, v d.ImageValidator) *Service {
	return &Service{r, a, t, v}
}
func (s *Service) authority(ctx context.Context, id uuid.UUID, write, admin bool) (bool, error) {
	u, e := s.accounts.Get(ctx, id)
	if e != nil {
		return false, e
	}
	if !u.IsActive || u.ActivationRequired {
		return false, shared.ErrUnauthorized
	}
	if admin && u.Role != user.RoleAdmin {
		return false, shared.ErrForbidden
	}
	if write && u.Role != user.RoleStaff && u.Role != user.RoleAdmin {
		return false, shared.ErrForbidden
	}
	return u.Role == user.RoleBorrower, nil
}
func (s *Service) List(ctx context.Context, actor uuid.UUID, f d.Filter) (d.Page, error) {
	b, e := s.authority(ctx, actor, false, false)
	if e != nil {
		return d.Page{}, e
	}
	if !f.Valid() {
		return d.Page{}, shared.ErrInvalidInput
	}
	f.Borrower = b
	if b {
		if f.Status != "" && f.Status != "ACTIVE" {
			return d.Page{}, shared.ErrForbidden
		}
		f.Status = "ACTIVE"
	}
	return s.repo.List(ctx, f)
}
func (s *Service) Detail(ctx context.Context, actor, id uuid.UUID) (d.Equipment, error) {
	b, e := s.authority(ctx, actor, false, false)
	if e != nil {
		return d.Equipment{}, e
	}
	v, e := s.repo.Get(ctx, id)
	if e != nil {
		return v, e
	}
	if b && v.Status != "ACTIVE" {
		return d.Equipment{}, shared.ErrNotFound
	}
	if !b {
		v.ArchiveSafety, e = s.repo.ArchiveSafety(ctx, id)
	}
	return v, e
}
func (s *Service) Categories(ctx context.Context, actor uuid.UUID, page int) ([]d.Category, error) {
	b, e := s.authority(ctx, actor, false, false)
	if e != nil {
		return nil, e
	}
	if page < 1 || page > 10000 {
		return nil, shared.ErrInvalidInput
	}
	return s.repo.Categories(ctx, b, page)
}
func run[T any](s *Service, ctx context.Context, actor uuid.UUID, op, key string, input any, admin bool, fn func(context.Context) (T, error)) (out T, err error) {
	k, e := uuid.Parse(key)
	if e != nil || k == uuid.Nil {
		return out, shared.ErrInvalidInput
	}
	payload, e := json.Marshal(input)
	if e != nil {
		return out, shared.ErrInvalidInput
	}
	h := sha256.Sum256(payload)
	hash := hex.EncodeToString(h[:])
	err = s.tx.Within(ctx, func(c context.Context) error {
		if e := s.accounts.LockAccounts(c, []uuid.UUID{actor}); e != nil {
			return e
		}
		if _, e := s.authority(c, actor, true, admin); e != nil {
			return e
		}
		data, e := s.repo.ReadReceipt(c, actor, op, k.String(), hash)
		if e == nil {
			return json.Unmarshal(data, &out)
		}
		if !errors.Is(e, shared.ErrNotFound) {
			return e
		}
		out, e = fn(c)
		if e != nil {
			return e
		}
		data, e = json.Marshal(out)
		if e != nil {
			return e
		}
		return s.repo.WriteReceipt(c, actor, op, k.String(), hash, data)
	})
	return
}
func (s *Service) category(ctx context.Context, id *uuid.UUID) error {
	if id == nil {
		return nil
	}
	v, e := s.repo.GetCategory(ctx, *id, true)
	if e != nil {
		return e
	}
	if !v.Active {
		return shared.ErrConflict
	}
	return nil
}
func (s *Service) Create(ctx context.Context, actor uuid.UUID, key string, in d.Create) (d.Equipment, error) {
	in.Name = strings.TrimSpace(in.Name)
	in.Reason = strings.TrimSpace(in.Reason)
	if !in.Metadata.Valid() || in.ExpectedVersion != 0 || in.Opening < 0 || in.Opening > d.MaxQuantity || !d.Text(in.Reason, 1000) || in.Opening > 0 && in.Reason == "" {
		return d.Equipment{}, shared.ErrInvalidInput
	}
	return run(s, ctx, actor, "equipment.create", key, in, false, func(c context.Context) (d.Equipment, error) {
		if e := s.category(c, in.CategoryID); e != nil {
			return d.Equipment{}, e
		}
		v := d.Equipment{ID: uuid.New(), Name: in.Name, Description: in.Description, CategoryID: in.CategoryID, Status: "ACTIVE", Version: 1, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()}
		if e := s.repo.Insert(c, v); e != nil {
			return v, e
		}
		if in.Opening > 0 {
			var e error
			v.Stock, e = v.Stock.ChangeAvailable(in.Opening)
			if e != nil {
				return v, e
			}
			v.Sequence = 1
			if e = s.repo.UpdateStock(c, v); e != nil {
				return v, e
			}
			if e = s.repo.Movement(c, d.Movement{ID: uuid.New(), EquipmentID: v.ID, ActorID: &actor, Sequence: 1, Kind: "OPENING", Delta: d.Stock{Available: in.Opening, Total: in.Opening}, After: v.Stock, Reason: in.Reason}); e != nil {
				return v, e
			}
		}
		if e := s.repo.Audit(c, actor, &v.ID, nil, "EQUIPMENT_CREATED", nil, v); e != nil {
			return v, e
		}
		return s.repo.Get(c, v.ID)
	})
}
func (s *Service) Edit(ctx context.Context, actor, id uuid.UUID, key string, in d.Metadata) (d.Equipment, error) {
	in.Name = strings.TrimSpace(in.Name)
	if !in.Valid() || in.ExpectedVersion < 1 {
		return d.Equipment{}, shared.ErrInvalidInput
	}
	return run(s, ctx, actor, "equipment.edit:"+id.String(), key, in, false, func(c context.Context) (d.Equipment, error) {
		v, e := s.lock(c, id)
		if e != nil {
			return v, e
		}
		if v.Version != in.ExpectedVersion || v.Status == "ARCHIVED" {
			return v, shared.ErrConflict
		}
		sameCategory := (v.CategoryID == nil && in.CategoryID == nil) || (v.CategoryID != nil && in.CategoryID != nil && *v.CategoryID == *in.CategoryID)
		if !sameCategory {
			if e = s.category(c, in.CategoryID); e != nil {
				return v, e
			}
		}
		old := v
		v.Name = in.Name
		v.Description = in.Description
		v.CategoryID = in.CategoryID
		v.Version++
		if e = s.repo.UpdateMetadata(c, v); e != nil {
			return v, e
		}
		if e = s.repo.Audit(c, actor, &id, nil, "EQUIPMENT_EDITED", old, v); e != nil {
			return v, e
		}
		return s.repo.Get(c, id)
	})
}
func (s *Service) lock(ctx context.Context, id uuid.UUID) (d.Equipment, error) {
	if id == uuid.Nil {
		return d.Equipment{}, shared.ErrInvalidInput
	}
	if e := s.repo.LockEquipment(ctx, id); e != nil {
		return d.Equipment{}, e
	}
	return s.repo.Get(ctx, id)
}
func (s *Service) Status(ctx context.Context, actor, id uuid.UUID, key string, in d.StatusInput) (d.Equipment, error) {
	if (in.Status != "ACTIVE" && in.Status != "INACTIVE" && in.Status != "ARCHIVED") || !in.Confirm || in.ExpectedVersion < 1 {
		return d.Equipment{}, shared.ErrInvalidInput
	}
	return run(s, ctx, actor, "equipment.status:"+id.String(), key, in, false, func(c context.Context) (d.Equipment, error) {
		v, e := s.lock(c, id)
		if e != nil {
			return v, e
		}
		if v.Version != in.ExpectedVersion || v.Status == "ARCHIVED" {
			return v, shared.ErrConflict
		}
		if in.Status == "ARCHIVED" {
			safe, e := s.repo.ArchiveSafety(c, id)
			if e != nil {
				return v, e
			}
			if (safe != "NOT_INSTALLED" && safe != "CLEAR") || v.Stock.Reserved != 0 || v.Stock.CheckedOut != 0 {
				return v, shared.ErrConflict
			}
		}
		old := v
		v.Status = in.Status
		v.Version++
		if e = s.repo.UpdateMetadata(c, v); e != nil {
			return v, e
		}
		if e = s.repo.Audit(c, actor, &id, nil, "EQUIPMENT_STATUS", old, v); e != nil {
			return v, e
		}
		return s.repo.Get(c, id)
	})
}
func (s *Service) Adjust(ctx context.Context, actor, id uuid.UUID, key string, in d.Adjustment) (d.Equipment, error) {
	in.Reason = strings.TrimSpace(in.Reason)
	if !in.Confirm || in.Reason == "" || !d.Text(in.Reason, 1000) || in.Quantity < 0 || in.Quantity > d.MaxQuantity || (in.Kind != "ADD" && in.Kind != "REMOVE" && in.Kind != "RECONCILE") || (in.Kind != "RECONCILE" && in.Quantity == 0) || (in.Kind == "RECONCILE" && (in.ExpectedSequence == nil || *in.ExpectedSequence < 0)) {
		return d.Equipment{}, shared.ErrInvalidInput
	}
	return run(s, ctx, actor, "equipment.stock:"+id.String(), key, in, in.Kind == "RECONCILE", func(c context.Context) (d.Equipment, error) {
		v, e := s.lock(c, id)
		if e != nil {
			return v, e
		}
		if v.Status == "ARCHIVED" || in.ExpectedSequence != nil && v.Sequence != *in.ExpectedSequence {
			return v, shared.ErrConflict
		}
		delta := in.Quantity
		if in.Kind == "REMOVE" {
			delta = -delta
		}
		if in.Kind == "RECONCILE" {
			delta = in.Quantity - v.Stock.Available
		}
		if delta == 0 {
			return v, shared.ErrConflict
		}
		old := v
		v.Stock, e = v.Stock.ChangeAvailable(delta)
		if e != nil {
			return v, e
		}
		v.Sequence++
		if e = s.repo.UpdateStock(c, v); e != nil {
			return v, e
		}
		if e = s.repo.Movement(c, d.Movement{ID: uuid.New(), EquipmentID: id, ActorID: &actor, Sequence: v.Sequence, Kind: in.Kind, Delta: d.Stock{Available: delta, Total: delta}, After: v.Stock, Reason: in.Reason}); e != nil {
			return v, e
		}
		if e = s.repo.Audit(c, actor, &id, nil, "STOCK_"+in.Kind, old, v); e != nil {
			return v, e
		}
		return s.repo.Get(c, id)
	})
}
func (s *Service) Movements(ctx context.Context, actor, id uuid.UUID, page int) ([]d.Movement, error) {
	if _, e := s.authority(ctx, actor, true, false); e != nil {
		return nil, e
	}
	if page < 1 || page > 100000 {
		return nil, shared.ErrInvalidInput
	}
	if _, e := s.repo.Get(ctx, id); e != nil {
		return nil, e
	}
	return s.repo.Movements(ctx, id, page)
}
func (s *Service) Category(ctx context.Context, actor, id uuid.UUID, key string, in d.CategoryInput) (d.Category, error) {
	in.Name = strings.TrimSpace(in.Name)
	if in.Name == "" || !d.Text(in.Name, 100) || in.Active == nil || (id != uuid.Nil && in.ExpectedVersion < 1) || (id == uuid.Nil && in.ExpectedVersion != 0) {
		return d.Category{}, shared.ErrInvalidInput
	}
	return run(s, ctx, actor, "category:"+id.String(), key, in, false, func(c context.Context) (d.Category, error) {
		v := d.Category{ID: uuid.New(), Version: 1}
		var old any
		if id != uuid.Nil {
			var e error
			v, e = s.repo.GetCategory(c, id, true)
			if e != nil {
				return v, e
			}
			if v.Version != in.ExpectedVersion {
				return v, shared.ErrConflict
			}
			old = v
			v.Version++
		}
		v.Name = in.Name
		v.Active = *in.Active
		if e := s.repo.PutCategory(c, v, id == uuid.Nil); e != nil {
			return v, e
		}
		return v, s.repo.Audit(c, actor, nil, &v.ID, "CATEGORY_SAVED", old, v)
	})
}
func (s *Service) Image(ctx context.Context, actor, id, imageID uuid.UUID) (d.Image, error) {
	v, e := s.Detail(ctx, actor, id)
	if e != nil {
		return d.Image{}, e
	}
	if v.ImageID == nil || *v.ImageID != imageID {
		return d.Image{}, shared.ErrNotFound
	}
	return s.repo.GetImage(ctx, id, imageID)
}
func (s *Service) SaveImage(ctx context.Context, actor, id uuid.UUID, key string, expected int64, raw []byte) (d.Equipment, error) {
	if _, e := s.authority(ctx, actor, true, false); e != nil {
		return d.Equipment{}, e
	}
	im, e := s.images.Validate(raw)
	if e != nil {
		return d.Equipment{}, e
	}
	if expected < 1 {
		return d.Equipment{}, shared.ErrInvalidInput
	}
	payload := struct {
		Hash    string
		Version int64
	}{im.Hash, expected}
	return run(s, ctx, actor, "equipment.image:"+id.String(), key, payload, false, func(c context.Context) (d.Equipment, error) {
		v, e := s.lock(c, id)
		if e != nil {
			return v, e
		}
		if v.Version != expected || v.Status == "ARCHIVED" {
			return v, shared.ErrConflict
		}
		old := v
		im.ID = uuid.New()
		im.EquipmentID = id
		im.ActorID = actor
		if e = s.repo.PutImage(c, im); e != nil {
			return v, e
		}
		v.ImageID = &im.ID
		v.Version++
		if e = s.repo.UpdateMetadata(c, v); e != nil {
			return v, e
		}
		if e = s.repo.Audit(c, actor, &id, nil, "CATALOG_IMAGE_UPDATED", old, v); e != nil {
			return v, e
		}
		return s.repo.Get(c, id)
	})
}
