// Trusted fresh, isolated presentation initializer. Never installed in the API.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	a "github.com/Maaku050/elabtrack-v2/backend/internal/application"
	ab "github.com/Maaku050/elabtrack-v2/backend/internal/application/borrowing"
	ai "github.com/Maaku050/elabtrack-v2/backend/internal/application/inventory"
	an "github.com/Maaku050/elabtrack-v2/backend/internal/application/notifications"
	ap "github.com/Maaku050/elabtrack-v2/backend/internal/application/profile"
	at "github.com/Maaku050/elabtrack-v2/backend/internal/application/terms"
	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	b "github.com/Maaku050/elabtrack-v2/backend/internal/domain/borrowing"
	i "github.com/Maaku050/elabtrack-v2/backend/internal/domain/inventory"
	t "github.com/Maaku050/elabtrack-v2/backend/internal/domain/terms"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/catalogimage"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/database"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/persistence/postgres"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/profileimage"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/security"
	"github.com/google/uuid"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const demoName = "elabtrack_v2_presentation_demo"

type joinedTx struct{ base a.Transactions }

func (x joinedTx) Within(c context.Context, f func(context.Context) error) error {
	if _, ok := database.TxFromContext(c); ok {
		return f(c)
	}
	return x.base.Within(c, f)
}

type account struct {
	Key, Email, Name, Password, Category, StudentID string
	Role                                            user.Role
}
type input struct {
	Accounts []account
	Assets   string
}
type equipment struct {
	Name, Category, Description, Asset, Status string
	Opening                                    int64
}

// Clock adapters affect only initial fictional history, never existing records.
type inventoryClock struct {
	*postgres.InventoryRepository
	at time.Time
}

func (r inventoryClock) Insert(c context.Context, v i.Equipment) error {
	v.CreatedAt = r.at
	v.UpdatedAt = r.at
	return r.InventoryRepository.Insert(c, v)
}
func (r inventoryClock) Movement(c context.Context, v i.Movement) error {
	q, ok := database.TxFromContext(c)
	if !ok {
		return errors.New("fixture transaction required")
	}
	_, e := q.Exec(c, `INSERT INTO inventory_movements(id,equipment_id,actor_id,sequence,kind,delta_available,delta_reserved,delta_checked_out,delta_damaged_held,delta_total,after_available,after_reserved,after_checked_out,after_damaged_held,after_total,reason,created_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`, v.ID, v.EquipmentID, v.ActorID, v.Sequence, v.Kind, v.Delta.Available, v.Delta.Reserved, v.Delta.CheckedOut, v.Delta.DamagedHeld, v.Delta.Total, v.After.Available, v.After.Reserved, v.After.CheckedOut, v.After.DamagedHeld, v.After.Total, v.Reason, r.at)
	return e
}

type termsClock struct {
	*postgres.TermsRepository
	at time.Time
}

func (r termsClock) Publish(c context.Context, v *t.Version) error {
	q, ok := database.TxFromContext(c)
	if !ok {
		return errors.New("fixture transaction required")
	}
	e := q.QueryRow(c, `INSERT INTO terms_versions(id,version,title,body,content_hash,published_by,created_at,published_at) VALUES($1,$2,$3,$4,$5,$6,$7,$7) RETURNING created_at,published_at`, v.ID, v.Version, v.Title, v.Body, v.ContentHash, v.PublishedBy, r.at).Scan(&v.CreatedAt, &v.PublishedAt)
	if e != nil {
		return e
	}
	_, e = q.Exec(c, `UPDATE terms_publication SET current_version_id=$1 WHERE id=1`, v.ID)
	return e
}
func (r termsClock) InsertAcceptance(c context.Context, u, v uuid.UUID) (*t.Acceptance, error) {
	q, ok := database.TxFromContext(c)
	if !ok {
		return nil, errors.New("fixture transaction required")
	}
	_, e := q.Exec(c, `INSERT INTO terms_acceptances(id,user_id,terms_version_id,accepted_at) VALUES($1,$2,$3,$4) ON CONFLICT(user_id,terms_version_id) DO NOTHING`, uuid.New(), u, v, r.at.Add(time.Hour))
	if e != nil {
		return nil, e
	}
	return r.FindAcceptance(c, u, v)
}

type borrowingClock struct {
	*postgres.BorrowingRepository
	at time.Time
}

func (r *borrowingClock) Clock(context.Context) (time.Time, error) { return r.at, nil }
func main() {
	if e := run(); e != nil {
		fmt.Fprintln(os.Stderr, e.Error())
		os.Exit(1)
	}
}
func run() error {
	c := context.Background()
	cfg, e := config.Load()
	if e != nil {
		return errors.New("configuration gate failed")
	}
	if cfg.DB.Name != demoName || cfg.DB.Host != "127.0.0.1" || cfg.DB.Port != "54836" || cfg.App.Env != "development" || os.Getenv("ELABTRACK_PRESENTATION_DEMO") != "1" {
		return errors.New("refusing non-presentation target")
	}
	conn, e := cfg.MigrationConnection()
	if e != nil {
		return errors.New("migration connection invalid")
	}
	db, e := database.New(c, conn)
	if e != nil {
		return errors.New("isolated database unavailable")
	}
	defer db.Close()
	var safe bool
	e = db.Pool.QueryRow(c, `SELECT current_database()='elabtrack_v2_presentation_demo' AND current_user='elabtrack_migrator' AND NOT rolsuper AND NOT rolcreatedb AND NOT rolcreaterole AND (SELECT identity='elabtrack-presentation-demo' AND NOT seeded FROM presentation_demo_identity WHERE id=1) FROM pg_roles WHERE rolname=current_user`).Scan(&safe)
	if e != nil || !safe {
		return errors.New("unseeded database identity guard failed")
	}
	var in input
	dec := json.NewDecoder(io.LimitReader(os.Stdin, 65536))
	dec.DisallowUnknownFields()
	if e = dec.Decode(&in); e != nil || len(in.Accounts) != 28 {
		return errors.New("requires exactly28 private fictional accounts")
	}
	roles := map[user.Role]int{}
	categories := map[string]int{}
	keys := map[string]bool{}
	emails := map[string]bool{}
	users := map[string]*user.User{}
	hasher := security.NewBcryptHasher(12)
	now := time.Now().UTC().Truncate(time.Second)
	for _, v := range in.Accounts {
		roles[v.Role]++
		categories[v.Category]++
		if keys[v.Key] || emails[v.Email] || !strings.HasSuffix(v.Email, ".invalid") || !strings.HasPrefix(v.Name, "DEMO ") || !v.Role.Valid() || len(v.Password) < 16 || len(v.Password) > 72 {
			return errors.New("fictional identity gate failed")
		}
		keys[v.Key] = true
		emails[v.Email] = true
		if v.Category == "STUDENT" && (!strings.HasSuffix(v.Email, "@students.example.invalid") || v.StudentID == "") {
			return errors.New("Student identity gate failed")
		}
		hash, e := hasher.Hash(v.Password)
		if e != nil {
			return errors.New("password hashing failed")
		}
		u := user.NewUser(v.Email, v.Name, hash)
		u.Role = v.Role
		u.CreatedAt = now.Add(-32 * 24 * time.Hour)
		u.UpdatedAt = u.CreatedAt
		users[v.Key] = u
	}
	if roles[user.RoleAdmin] != 1 || roles[user.RoleStaff] != 2 || roles[user.RoleBorrower] != 25 || categories["STUDENT"] != 20 || categories["FACULTY"] != 5 {
		return errors.New("role/category gate failed")
	}
	raw, e := os.ReadFile(filepath.Join(in.Assets, "equipment.json"))
	if e != nil {
		return errors.New("equipment assets missing")
	}
	var specs []equipment
	if json.Unmarshal(raw, &specs) != nil || len(specs) != 40 {
		return errors.New("requires40 equipment assets")
	}
	tx := joinedTx{database.NewTxManager(db.Pool)}
	ar := postgres.NewAccountsRepository(db.Pool)
	ur := postgres.NewUserRepository(db.Pool)
	ir := inventoryClock{postgres.NewInventoryRepository(db.Pool), now.Add(-31 * 24 * time.Hour)}
	tr := termsClock{postgres.NewTermsRepository(db.Pool), now.Add(-30 * 24 * time.Hour)}
	ts := at.NewService(tr, ur, tx)
	is := ai.NewService(ir, ar, tx, catalogimage.Validator{})
	pr := ap.NewService(postgres.NewProfileRepository(db.Pool), ar, tx, profileimage.Validator{})
	br := &borrowingClock{postgres.NewBorrowingRepository(db.Pool), now}
	bs := ab.NewService(br, ar, ir, ts, tx)
	out := map[string]any{"scenario_clock": now, "accounts": map[string]any{}, "scenarios": map[string]any{}}
	stage := "initialization"
	e = tx.Within(c, func(c context.Context) error {
		q, _ := database.TxFromContext(c)
		var n int
		if e := q.QueryRow(c, `SELECT count(*) FROM users`).Scan(&n); e != nil || n != 0 {
			return errors.New("target must be empty")
		}
		if e := q.QueryRow(c, `SELECT count(*) FROM schema_migrations`).Scan(&n); e != nil || n != 11 {
			return errors.New("requires migrations1–11")
		}
		admin := users["admin"]
		staff := users["staff1"]
		if admin == nil || staff == nil {
			return errors.New("required actors missing")
		}
		stage = "accounts and avatars"
		for index, v := range in.Accounts {
			u := users[v.Key]
			if e := ur.Create(c, u); e != nil {
				return e
			}
			if u.Role == user.RoleBorrower {
				var sid *string
				if v.Category == "STUDENT" {
					sid = &v.StudentID
				}
				_, e := q.Exec(c, `INSERT INTO borrower_profiles(user_id,borrower_type,student_id,course,contact_number) VALUES($1,$2,$3,$4,$5)`, u.ID, v.Category, sid, "Fictional Food Service Laboratory", "demo-contact-"+v.Key)
				if e != nil {
					return e
				}
			}
			if e := ar.AddAudit(c, admin.ID, u.ID, "DEMO_INITIALIZED_NO_EMAIL"); e != nil {
				return e
			}
			img, e := os.ReadFile(filepath.Join(in.Assets, "avatars", fmt.Sprintf("avatar-%02d.png", index+1)))
			if e != nil {
				return e
			}
			if _, e = pr.Save(c, u.ID, u.ID, uuid.NewString(), 0, img, false); e != nil {
				return e
			}
			out["accounts"].(map[string]any)[v.Key] = map[string]any{"id": u.ID, "email": u.Email, "name": u.Name, "role": u.Role, "borrower_type": v.Category}
		}
		stage = "demonstration terms"
		v, e := ts.Publish(c, at.PublishCommand{ActorID: admin.ID, Version: "DEMO-8-14-1", Title: "DEMONSTRATION TERMS — NOT OFFICIAL FSMO POLICY", Body: "Fictional isolated presentation policy only. These are not official FSMO institutional terms. This demo account can rehearse requests, physical checkout, returns and equivalent replacement. Historical consent records belong exclusively to generated fictional actors. No normal or production policy, personal records or email delivery are represented."})
		if e != nil {
			return e
		}
		out["terms"] = v
		for _, x := range in.Accounts {
			if x.Role == user.RoleBorrower && x.Key != "student20" {
				if _, e = ts.Accept(c, users[x.Key].ID, v.ID); e != nil {
					return e
				}
			}
		}
		out["seeded_fictional_acceptances"] = 24
		stage = "equipment and catalog assets"
		cats := map[string]uuid.UUID{}
		equipment := []i.Equipment{}
		for _, spec := range specs {
			cat, ok := cats[spec.Category]
			if !ok {
				active := true
				v, e := is.Category(c, admin.ID, uuid.Nil, uuid.NewString(), i.CategoryInput{Name: spec.Category, Active: &active})
				if e != nil {
					return e
				}
				cat = v.ID
				cats[spec.Category] = cat
			}
			eq, e := is.Create(c, admin.ID, uuid.NewString(), i.Create{Metadata: i.Metadata{Name: spec.Name, Description: spec.Description + " Fictional demonstration inventory.", CategoryID: &cat}, Opening: spec.Opening, Reason: "Isolated fictional opening acquisition"})
			if e != nil {
				return e
			}
			img, e := os.ReadFile(filepath.Join(in.Assets, spec.Asset))
			if e != nil {
				return e
			}
			eq, e = is.SaveImage(c, admin.ID, eq.ID, uuid.NewString(), eq.Version, img)
			if e != nil {
				return e
			}
			if spec.Status == "INACTIVE" {
				eq, e = is.Status(c, admin.ID, eq.ID, uuid.NewString(), i.StatusInput{Status: "INACTIVE", ExpectedVersion: eq.Version, Confirm: true})
				if e != nil {
					return e
				}
			}
			equipment = append(equipment, eq)
		}
		out["equipment"] = equipment
		stage = "historical completed loans"
		for index, x := range in.Accounts {
			if x.Role != user.RoleBorrower || x.Key == "student20" {
				continue
			}
			br.at = now.Add(-28*24*time.Hour + time.Duration(index)*3*time.Hour)
			due := br.at.Add(2 * time.Hour)
			loan, e := bs.Direct(c, staff.ID, uuid.NewString(), b.Input{BorrowerID: users[x.Key].ID, Items: []b.Line{{EquipmentID: equipment[index].ID, Quantity: 1}}, DueAt: &due, Confirm: true, Handover: true})
			if e != nil {
				return e
			}
			br.at = br.at.Add(time.Hour)
			loan, e = bs.Return(c, staff.ID, loan.ID, uuid.NewString(), b.ReturnInput{ExpectedEvents: len(loan.Events), Lines: []b.ReturnLine{{ItemID: loan.Items[0].ID, Good: 1}}, Reason: "Fictional laboratory session complete; physically inspected usable", Confirm: true})
			if e != nil || loan.Status != "COMPLETED" {
				return errors.New("historical return failed")
			}
		}
		scenarios := out["scenarios"].(map[string]any)
		direct := func(key string, borrower string, eq int, qty int64, issued, due time.Time) (b.Record, error) {
			br.at = issued
			v, e := bs.Direct(c, staff.ID, uuid.NewString(), b.Input{BorrowerID: users[borrower].ID, Items: []b.Line{{EquipmentID: equipment[eq].ID, Quantity: qty}}, DueAt: &due, Confirm: true, Handover: true})
			if e == nil {
				scenarios[key] = v
			}
			return v, e
		}
		ret := func(key string, v b.Record, at time.Time, good, damaged, lost int64) (b.Record, error) {
			br.at = at
			v, e := bs.Return(c, staff.ID, v.ID, uuid.NewString(), b.ReturnInput{ExpectedEvents: len(v.Events), Lines: []b.ReturnLine{{ItemID: v.Items[0].ID, Good: good, Damaged: damaged, Lost: lost}}, Reason: "Fictional in-person verified return; condition recorded by Staff", Confirm: true})
			if e == nil {
				scenarios[key] = v
			}
			return v, e
		}
		stage = "damage loss replacement fine histories"
		// Oldest first; independent equipment keeps movement timestamps chronological.
		v1, e := direct("resolved_replacement", "student03", 1, 3, now.Add(-6*24*time.Hour), now.Add(-5*24*time.Hour))
		if e != nil {
			return e
		}
		v1, e = ret("resolved_replacement", v1, now.Add(-4*24*time.Hour), 1, 1, 1)
		if e != nil {
			return e
		}
		br.at = now.Add(-3 * 24 * time.Hour)
		for _, o := range v1.Obligations {
			v1, e = bs.Replace(c, staff.ID, v1.ID, uuid.NewString(), b.ReplacementInput{ExpectedEvents: len(v1.Events), ObligationID: o.ID, Quantity: o.Required, Reason: "Fictional equivalent replacement verified in person", Equivalent: true, Confirm: true})
			if e != nil {
				return e
			}
		}
		br.at = now.Add(-2 * 24 * time.Hour)
		v1, e = bs.Read(c, admin.ID, v1.ID)
		if e != nil {
			return e
		}
		v1, e = bs.ClearFine(c, admin.ID, v1.ID, uuid.NewString(), b.FineInput{Method: "PAID", Note: "Fictional manual settlement demonstration; no actual payment", ExpectedMinor: v1.Fine.Outstanding, Confirm: true})
		if e != nil {
			return e
		}
		scenarios["resolved_replacement"] = v1
		v2, e := direct("outstanding_replacements", "student04", 2, 4, now.Add(-5*24*time.Hour), now.Add(-4*24*time.Hour))
		if e != nil {
			return e
		}
		if _, e = ret("outstanding_replacements", v2, now.Add(-3*24*time.Hour), 1, 1, 1); e != nil {
			return e
		}
		v3, e := direct("overdue", "student05", 3, 2, now.Add(-3*24*time.Hour), now.Add(-2*24*time.Hour))
		if e != nil {
			return e
		}
		scenarios["overdue"] = v3
		v4, e := direct("waived_completed", "faculty01", 4, 1, now.Add(-3*24*time.Hour), now.Add(-2*24*time.Hour))
		if e != nil {
			return e
		}
		v4, e = ret("waived_completed", v4, now.Add(-24*time.Hour), 1, 0, 0)
		if e != nil {
			return e
		}
		br.at = now.Add(-23 * time.Hour)
		v4, e = bs.Read(c, admin.ID, v4.ID)
		if e != nil {
			return e
		}
		v4, e = bs.ClearFine(c, admin.ID, v4.ID, uuid.NewString(), b.FineInput{Method: "WAIVED", Note: "Fictional Admin waiver example with preserved assessment", ExpectedMinor: v4.Fine.Outstanding, Confirm: true})
		if e != nil {
			return e
		}
		scenarios["waived_completed"] = v4
		v5, e := direct("partial_return", "faculty02", 5, 4, now.Add(-24*time.Hour), now.Add(24*time.Hour))
		if e != nil {
			return e
		}
		if _, e = ret("partial_return", v5, now.Add(-12*time.Hour), 2, 0, 0); e != nil {
			return e
		}
		if _, e = direct("due_soon", "student06", 6, 2, now.Add(-6*time.Hour), now.Add(12*time.Hour)); e != nil {
			return e
		}
		stage = "request lifecycle histories"
		for index, kind := range []string{"expired", "cancelled", "denied", "request_checkout", "pending"} {
			br.at = now.Add(-time.Duration(5-index) * time.Hour)
			if kind == "expired" {
				br.at = now.Add(-26 * time.Hour)
			}
			borrower := users[fmt.Sprintf("student%02d", index+7)]
			v, e := bs.Submit(c, borrower.ID, uuid.NewString(), b.Input{Items: []b.Line{{EquipmentID: equipment[index+7].ID, Quantity: int64(index%3 + 1)}}, Confirm: true})
			if e != nil {
				return e
			}
			br.at = br.at.Add(10 * time.Minute)
			action := ""
			actor := staff.ID
			decision := b.Decision{Confirm: true}
			switch kind {
			case "cancelled":
				action = "CANCELLED"
				actor = borrower.ID
			case "denied":
				action = "DENIED"
				decision.Reason = "Fictional supervised session rescheduled; reservation released"
			case "request_checkout":
				action = "CHECKOUT"
				due := now.Add(36 * time.Hour)
				decision.DueAt = &due
				decision.Handover = true
			}
			if action != "" {
				v, e = bs.Decide(c, actor, v.ID, uuid.NewString(), action, decision)
				if e != nil {
					return e
				}
			}
			scenarios[kind] = v
		}
		br.at = now
		if _, e = bs.Sweep(c, 100); e != nil {
			return e
		}
		for key, old := range scenarios {
			v, e := bs.Read(c, admin.ID, old.(b.Record).ID)
			if e != nil {
				return e
			}
			scenarios[key] = v
		}
		stage = "notifications"
		ns := an.NewService(postgres.NewNotificationsRepository(db.Pool), ar, tx, 86400)
		for pass := 0; pass < 20; pass++ {
			n, e := ns.Sweep(c, 100)
			if e != nil {
				return e
			}
			if n < 100 {
				break
			}
		}
		var invalid int
		if e = q.QueryRow(c, `SELECT count(*) FROM equipment WHERE total_tracked<>available+reserved+checked_out+damaged_held OR LEAST(available,reserved,checked_out,damaged_held)<0`).Scan(&invalid); e != nil || invalid != 0 {
			return errors.New("stock invariant failed")
		}
		_, e = q.Exec(c, `UPDATE presentation_demo_identity SET seeded=true WHERE id=1 AND identity='elabtrack-presentation-demo'`)
		return e
	})
	if e != nil {
		return fmt.Errorf("presentation seed failed at %s; transaction rolled back", stage)
	}
	return json.NewEncoder(os.Stdout).Encode(out)
}
