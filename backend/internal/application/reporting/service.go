package reporting

import (
	"bytes"
	"context"
	"encoding/csv"
	"github.com/Maaku050/elabtrack-v2/backend/internal/application"
	a "github.com/Maaku050/elabtrack-v2/backend/internal/domain/accounts"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/reporting"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/google/uuid"
)

type Service struct {
	repo     d.Repository
	accounts a.Repository
	tx       application.Transactions
}

func NewService(r d.Repository, a a.Repository, t application.Transactions) *Service {
	return &Service{r, a, t}
}
func (s *Service) actor(c context.Context, id uuid.UUID) (user.Role, error) {
	if e := s.accounts.LockAccounts(c, []uuid.UUID{id}); e != nil {
		return "", e
	}
	u, e := s.accounts.Get(c, id)
	if e != nil {
		return "", e
	}
	if !u.IsActive || u.ActivationRequired {
		return "", shared.ErrUnauthorized
	}
	if !u.Role.Valid() {
		return "", shared.ErrForbidden
	}
	return u.Role, nil
}
func (s *Service) Dashboard(c context.Context, id uuid.UUID, ranges ...int) (v d.Dashboard, e error) {
	days := 7
	if len(ranges) > 1 {
		return v, shared.ErrInvalidInput
	}
	if len(ranges) == 1 {
		days = ranges[0]
	}
	if days != 7 && days != 30 {
		return v, shared.ErrInvalidInput
	}
	e = s.tx.Within(c, func(c context.Context) error {
		role, err := s.actor(c, id)
		if err != nil {
			return err
		}
		v, err = s.repo.Dashboard(c, id, role, days)
		return err
	})
	return
}
func (s *Service) report(c context.Context, id uuid.UUID, key string, f d.Filter, limit int) (v d.Page, e error) {
	def, ok := d.Find(key)
	if !ok || !f.Valid() || !def.Accepts(f) || f.DueToday && key != "active" {
		return v, shared.ErrInvalidInput
	}
	e = s.tx.Within(c, func(c context.Context) error {
		role, err := s.actor(c, id)
		if err != nil {
			return err
		}
		if role == user.RoleBorrower || def.AdminOnly && role != user.RoleAdmin {
			return shared.ErrForbidden
		}
		v, err = s.repo.Report(c, key, f, limit)
		return err
	})
	return
}
func (s *Service) Report(c context.Context, id uuid.UUID, key string, f d.Filter) (d.Page, error) {
	return s.report(c, id, key, f, f.PerPage)
}
func (s *Service) Export(c context.Context, id uuid.UUID, key string, f d.Filter) ([]byte, error) {
	f.Page = 1
	v, e := s.report(c, id, key, f, 5000)
	if e != nil {
		return nil, e
	}
	if v.Total > 5000 {
		return nil, d.ErrExportLimit
	}
	var b bytes.Buffer
	w := csv.NewWriter(&b)
	w.UseCRLF = true
	if e = w.Write(v.Columns); e != nil {
		return nil, e
	}
	for _, row := range v.Rows {
		out := make([]string, len(row))
		for i, cell := range row {
			out[i] = d.SafeCell(cell)
		}
		if e = w.Write(out); e != nil {
			return nil, e
		}
	}
	w.Flush()
	return b.Bytes(), w.Error()
}
func (s *Service) Definitions(c context.Context, id uuid.UUID) (out []d.Definition, e error) {
	out = []d.Definition{}
	e = s.tx.Within(c, func(c context.Context) error {
		role, err := s.actor(c, id)
		if err != nil {
			return err
		}
		if role == user.RoleBorrower {
			return shared.ErrForbidden
		}
		for _, def := range d.Definitions {
			if !def.AdminOnly || role == user.RoleAdmin {
				out = append(out, def)
			}
		}
		return nil
	})
	return
}

// ExportData uses the same authorized filter/order/snapshot and limit as CSV.
func (s *Service) ExportData(c context.Context, id uuid.UUID, key string, f d.Filter) (d.Page, error) {
	f.Page = 1
	v, e := s.report(c, id, key, f, 5000)
	if e != nil {
		return v, e
	}
	if v.Total > 5000 {
		return d.Page{}, d.ErrExportLimit
	}
	return v, nil
}
