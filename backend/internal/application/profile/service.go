package profile

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/Maaku050/elabtrack-v2/backend/internal/application"
	accounts "github.com/Maaku050/elabtrack-v2/backend/internal/domain/accounts"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/profile"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/google/uuid"
)

type Service struct {
	repo      d.Repository
	accounts  accounts.Repository
	tx        application.Transactions
	validator d.Validator
}

func NewService(r d.Repository, a accounts.Repository, t application.Transactions, v d.Validator) *Service {
	return &Service{r, a, t, v}
}
func (s *Service) authorize(c context.Context, actor, target uuid.UUID, write bool) error {
	if actor == uuid.Nil || target == uuid.Nil {
		return shared.ErrInvalidInput
	}
	if e := s.accounts.LockAccounts(c, []uuid.UUID{actor, target}); e != nil {
		return e
	}
	a, e := s.accounts.Get(c, actor)
	if e != nil {
		return e
	}
	if !a.IsActive || a.ActivationRequired {
		return shared.ErrUnauthorized
	}
	if !a.Role.Valid() {
		return shared.ErrForbidden
	}
	if actor == target {
		return nil
	}
	if write || a.Role == user.RoleBorrower {
		return shared.ErrForbidden
	}
	t, e := s.accounts.Get(c, target)
	if e != nil {
		return e
	}
	if a.Role == user.RoleAdmin || (a.Role == user.RoleStaff && t.Role == user.RoleBorrower) {
		return nil
	}
	return shared.ErrForbidden
}
func (s *Service) Metadata(c context.Context, actor, target uuid.UUID) (v d.Metadata, e error) {
	e = s.tx.Within(c, func(c context.Context) error {
		if e := s.authorize(c, actor, target, false); e != nil {
			return e
		}
		var e error
		v, e = s.repo.Metadata(c, target)
		return e
	})
	return
}
func (s *Service) Image(c context.Context, actor, target, id uuid.UUID) (v d.Image, e error) {
	if id == uuid.Nil {
		return v, shared.ErrInvalidInput
	}
	e = s.tx.Within(c, func(c context.Context) error {
		if e := s.authorize(c, actor, target, false); e != nil {
			return e
		}
		var e error
		v, e = s.repo.Image(c, target, id)
		return e
	})
	return
}
func (s *Service) Save(c context.Context, actor, target uuid.UUID, key string, expected int64, raw []byte, remove bool) (out d.Metadata, e error) {
	if id, e := uuid.Parse(key); e != nil || id == uuid.Nil || expected < 0 || expected >= 9223372036854775807 {
		return out, shared.ErrInvalidInput
	}
	var image d.Image
	if !remove {
		image, e = s.validator.Validate(raw)
		if e != nil {
			return out, e
		}
	} else if len(raw) > 0 {
		return out, shared.ErrInvalidInput
	}
	payload, _ := json.Marshal(struct {
		Target  uuid.UUID
		Version int64
		Hash    string
		Remove  bool
	}{target, expected, image.Hash, remove})
	sum := sha256.Sum256(payload)
	hash := hex.EncodeToString(sum[:])
	operation := "PROFILE_IMAGE:" + target.String()
	e = s.tx.Within(c, func(c context.Context) error {
		if e := s.authorize(c, actor, target, true); e != nil {
			return e
		}
		receipt, e := s.accounts.ReadReceipt(c, actor, operation, key, hash)
		if e != nil {
			return e
		}
		if receipt != nil {
			return json.Unmarshal(receipt, &out)
		}
		current, e := s.repo.Metadata(c, target)
		if e != nil {
			return e
		}
		if current.Version != expected {
			return shared.ErrConflict
		}
		image.Metadata = d.Metadata{AccountID: target, Version: expected + 1}
		action := "PROFILE_IMAGE_REMOVED"
		if !remove {
			id := uuid.New()
			image.ImageID = &id
			action = "PROFILE_IMAGE_UPDATED"
		}
		if e = s.repo.Save(c, image); e != nil {
			return e
		}
		if e = s.accounts.AddAudit(c, actor, target, action); e != nil {
			return e
		}
		out = image.Metadata
		receipt, e = json.Marshal(out)
		if e != nil {
			return shared.ErrInternal
		}
		return s.accounts.WriteReceipt(c, actor, operation, key, hash, receipt)
	})
	return
}

type OwnAccount struct {
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Role          user.Role `json:"role"`
	BorrowerType  string    `json:"borrower_type"`
	StudentID     string    `json:"student_id"`
	Course        string    `json:"course"`
	ContactNumber string    `json:"contact_number"`
}

func (s *Service) Own(c context.Context, actor uuid.UUID) (out OwnAccount, e error) {
	e = s.tx.Within(c, func(c context.Context) error {
		if e := s.authorize(c, actor, actor, false); e != nil {
			return e
		}
		u, e := s.accounts.Get(c, actor)
		if e != nil {
			return e
		}
		out = OwnAccount{u.Name, u.Email, u.Role, u.BorrowerType, u.StudentID, u.Course, u.ContactNumber}
		return nil
	})
	return
}
