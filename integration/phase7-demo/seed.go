// Trusted local fixture tool, built ONLY inside the private committed Phase7 snapshot.
// This is not an API route, production account bootstrap or activation bypass.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	a "github.com/Maaku050/elabtrack-v2/backend/internal/application"
	appborrowing "github.com/Maaku050/elabtrack-v2/backend/internal/application/borrowing"
	appinventory "github.com/Maaku050/elabtrack-v2/backend/internal/application/inventory"
	appterms "github.com/Maaku050/elabtrack-v2/backend/internal/application/terms"
	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/borrowing"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/inventory"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/catalogimage"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/database"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/persistence/postgres"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/security"
	"github.com/google/uuid"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"os"
	"strings"
	"time"
)

const demoName = "elabtrack_v2_phase7_demo"

type joinedTx struct{ base a.Transactions }

func (t joinedTx) Within(c context.Context, fn func(context.Context) error) error {
	if _, ok := database.TxFromContext(c); ok {
		return fn(c)
	}
	return t.base.Within(c, fn)
}

type inputAccount struct {
	Key, Email, Name, Password, Category, StudentID string
	Role                                            user.Role
}
type input struct{ Accounts []inputAccount }
type pastClock struct {
	*postgres.BorrowingRepository
	at time.Time
}

func (r pastClock) Clock(context.Context) (time.Time, error) { return r.at, nil }
func main() {
	if e := run(); e != nil {
		fmt.Fprintln(os.Stderr, "Demo command failed; verify the isolated identity, empty seed target or accepted demo terms. No credentials or database row details logged.")
		os.Exit(1)
	}
}
func run() error {
	c := context.Background()
	cfg, e := config.Load()
	if e != nil {
		return errors.New("configuration invalid")
	}
	if cfg.DB.Name != demoName || cfg.DB.Port != "54835" || cfg.DB.Host != "127.0.0.1" || cfg.App.Env != "development" || os.Getenv("ELABTRACK_PHASE7_DEMO") != "1" {
		return errors.New("refusing non-demo target")
	}
	ownerCfg, e := cfg.MigrationConnection()
	if e != nil {
		return errors.New("migration connection invalid")
	}
	db, e := database.New(c, ownerCfg)
	if e != nil {
		return errors.New("database unavailable")
	}
	defer db.Close()
	var safe bool
	e = db.Pool.QueryRow(c, `SELECT current_database()='elabtrack_v2_phase7_demo' AND current_user='elabtrack_migrator' AND NOT rolsuper AND NOT rolcreatedb AND NOT rolcreaterole AND (SELECT identity='elabtrack-phase7-owner-demo' FROM phase7_demo_identity WHERE id=1) FROM pg_roles WHERE rolname=current_user`).Scan(&safe)
	if e != nil || !safe {
		return errors.New("database identity guard failed")
	}
	tx := joinedTx{database.NewTxManager(db.Pool)}
	ar := postgres.NewAccountsRepository(db.Pool)
	ur := postgres.NewUserRepository(db.Pool)
	ir := postgres.NewInventoryRepository(db.Pool)
	br := postgres.NewBorrowingRepository(db.Pool)
	ts := appterms.NewService(postgres.NewTermsRepository(db.Pool), ur, tx)
	is := appinventory.NewService(ir, ar, tx, catalogimage.Validator{})
	if len(os.Args) > 1 && os.Args[1] == "expire" {
		var borrower, equipment uuid.UUID
		if e = db.Pool.QueryRow(c, `SELECT id FROM users WHERE email='student.one@students.example.invalid'`).Scan(&borrower); e != nil {
			return errors.New("demo Student absent")
		}
		if e = db.Pool.QueryRow(c, `SELECT id FROM equipment WHERE name='DEMO Culinary Ladle'`).Scan(&equipment); e != nil {
			return errors.New("demo equipment absent")
		}
		past := pastClock{br, time.Now().UTC().Add(-25 * time.Hour)}
		s := appborrowing.NewService(past, ar, ir, ts, tx)
		v, e := s.Submit(c, borrower, uuid.NewString(), d.Input{Confirm: true, Items: []d.Line{{EquipmentID: equipment, Quantity: 1}}})
		if e != nil {
			return errors.New("expiry fixture requires active Student, accepted current demo terms and available stock")
		}
		live := appborrowing.NewService(br, ar, ir, ts, tx)
		_, e = live.Sweep(c, 100)
		if e != nil {
			return e
		}
		_, e = live.Sweep(c, 100)
		if e != nil {
			return e
		}
		got, e := live.Read(c, borrower, v.ID)
		if e != nil || got.Status != "EXPIRED" || got.Items[0].Reserved != 0 {
			return errors.New("expiry did not release exactly once")
		}
		return json.NewEncoder(os.Stdout).Encode(map[string]any{"id": got.ID, "status": got.Status, "reserved": 0, "expires_at": got.ExpiresAt, "events": got.Events})
	}
	var in input
	dec := json.NewDecoder(io.LimitReader(os.Stdin, 32768))
	dec.DisallowUnknownFields()
	if e = dec.Decode(&in); e != nil || len(in.Accounts) != 5 {
		return errors.New("requires exactly five private fixture accounts")
	}
	roles := map[user.Role]int{}
	keys := map[string]bool{}
	for _, x := range in.Accounts {
		roles[x.Role]++
		if !strings.HasSuffix(x.Email, ".invalid") || !strings.HasPrefix(x.Name, "DEMO ") || !x.Role.Valid() || len(x.Password) < 16 || len(x.Password) > 72 || keys[x.Key] {
			return errors.New("invalid private fixture input")
		}
		keys[x.Key] = true
	}
	if roles[user.RoleAdmin] != 1 || roles[user.RoleStaff] != 1 || roles[user.RoleBorrower] != 3 {
		return errors.New("invalid fixture role counts")
	}
	hasher := security.NewBcryptHasher(12)
	users := map[string]*user.User{}
	for _, x := range in.Accounts {
		hash, e := hasher.Hash(x.Password)
		if e != nil {
			return errors.New("credential hashing failed")
		}
		u := user.NewUser(x.Email, x.Name, hash)
		u.Role = x.Role
		users[x.Key] = u
	}
	output := map[string]any{}
	e = tx.Within(c, func(c context.Context) error {
		q, _ := database.TxFromContext(c)
		var n int
		if e = q.QueryRow(c, `SELECT count(*) FROM users`).Scan(&n); e != nil || n != 0 {
			return errors.New("refusing to seed existing accounts; use guarded reset")
		}
		if e = q.QueryRow(c, `SELECT count(*) FROM schema_migrations`).Scan(&n); e != nil || n != 8 {
			return errors.New("requires exactly migrations000001–000008")
		}
		admin := users["admin"]
		if admin == nil || admin.Role != user.RoleAdmin {
			return errors.New("Admin fixture absent")
		}
		if e = ur.Create(c, admin); e != nil {
			return e
		}
		for _, x := range in.Accounts {
			u := users[x.Key]
			if x.Key != "admin" {
				if e = ur.Create(c, u); e != nil {
					return e
				}
			}
			if x.Role == user.RoleBorrower {
				if x.Category != "STUDENT" && x.Category != "FACULTY" {
					return errors.New("fixture category invalid")
				}
				var sid *string
				if x.Category == "STUDENT" {
					if x.StudentID == "" {
						return errors.New("Student ID required")
					}
					sid = &x.StudentID
				}
				_, e = q.Exec(c, `INSERT INTO borrower_profiles(user_id,borrower_type,student_id,course,contact_number) VALUES($1,$2,$3,$4,'')`, u.ID, x.Category, sid, "Fictional Culinary Program")
				if e != nil {
					return e
				}
			}
			if e = ar.AddAudit(c, admin.ID, u.ID, "DEMO_INITIALIZED_NO_EMAIL"); e != nil {
				return e
			}
			output[x.Key] = map[string]any{"id": u.ID, "email": u.Email, "name": u.Name, "role": u.Role, "borrower_type": x.Category}
		}
		version, e := ts.Publish(c, appterms.PublishCommand{ActorID: admin.ID, Version: "DEMO-1", Title: "DEMONSTRATION TERMS — NOT OFFICIAL FSMO POLICY", Body: "Synthetic demonstration policy for isolated Phase7 owner review only. This is not approved FSMO institutional policy and is never published to normal or production data. By accepting you authorize this fictional demo account to test request reservation, cancellation, denial, physical checkout and direct issuance. No real borrowing or email delivery is represented."})
		if e != nil {
			return e
		}
		output["terms"] = map[string]any{"id": version.ID, "version": version.Version}
		categories := []uuid.UUID{}
		active := true
		for _, name := range []string{"DEMO Utensils", "DEMO Preparation", "DEMO Cookware"} {
			cat, e := is.Category(c, admin.ID, uuid.Nil, uuid.NewString(), inventory.CategoryInput{Name: name, Active: &active})
			if e != nil {
				return e
			}
			categories = append(categories, cat.ID)
		}
		equipment := []inventory.Equipment{}
		names := []string{"DEMO Culinary Ladle", "DEMO Mixing Bowl", "DEMO Chef Knife", "DEMO Measuring Cup", "DEMO Saucepan", "DEMO Whisk", "DEMO Rolling Pin", "DEMO Inactive Grater"}
		for n, name := range names {
			opening := int64(12 + n)
			if n == 0 {
				opening = 20
			}
			if n == 6 {
				opening = 0
			}
			cat := categories[n%3]
			v, e := is.Create(c, admin.ID, uuid.NewString(), inventory.Create{Metadata: inventory.Metadata{Name: name, Description: "Fictional owner-demonstration equipment; never copied from normal inventory.", CategoryID: &cat}, Opening: opening, Reason: "Explicit isolated demo opening acquisition"})
			if e != nil {
				return e
			}
			v, e = is.SaveImage(c, admin.ID, v.ID, uuid.NewString(), v.Version, illustration(n))
			if e != nil {
				return e
			}
			if n == 7 {
				v, e = is.Status(c, admin.ID, v.ID, uuid.NewString(), inventory.StatusInput{Status: "INACTIVE", ExpectedVersion: v.Version, Confirm: true})
				if e != nil {
					return e
				}
			}
			equipment = append(equipment, v)
			if n == 0 {
				output["equipment"] = v
			}
		}
		output["equipment_records"] = equipment
		output["terms_acceptances_seeded"] = 0
		_, e = q.Exec(c, `UPDATE phase7_demo_identity SET seeded=true WHERE id=1`)
		return e
	})
	if e != nil {
		return e
	}
	return json.NewEncoder(os.Stdout).Encode(output)
}

// Original native line illustrations rendered locally; no borrowed photos/CDN.
func illustration(kind int) []byte {
	im := image.NewRGBA(image.Rect(0, 0, 320, 240))
	draw.Draw(im, im.Bounds(), &image.Uniform{color.RGBA{237, 237, 251, 255}}, image.Point{}, draw.Src)
	ink := &image.Uniform{color.RGBA{66, 57, 129, 255}}
	rect := func(x, y, w, h int) { draw.Draw(im, image.Rect(x, y, x+w, y+h), ink, image.Point{}, draw.Src) }
	switch kind {
	case 0:
		rect(150, 35, 18, 140)
		for y := 150; y < 205; y++ {
			w := 70 - (y-175)*(y-175)/12
			if w > 0 {
				rect(160-w/2, y, w, 1)
			}
		}
	case 1:
		rect(75, 85, 170, 8)
		for y := 93; y < 175; y++ {
			w := 170 - (y - 93)
			rect(160-w/2, y, 4, 1)
			rect(160+w/2-4, y, 4, 1)
		}
		rect(118, 175, 84, 6)
	case 2:
		rect(75, 62, 28, 135)
		for y := 55; y < 160; y++ {
			rect(107, y, 85-(y-55)/2, 1)
		}
	case 3:
		rect(100, 70, 105, 120)
		rect(201, 90, 40, 7)
		rect(234, 90, 7, 65)
		rect(202, 148, 39, 7)
		for y := 85; y < 175; y += 22 {
			draw.Draw(im, image.Rect(110, y, 139, y+4), &image.Uniform{color.RGBA{237, 237, 251, 255}}, image.Point{}, draw.Src)
		}
	case 4:
		rect(80, 105, 140, 70)
		rect(65, 95, 170, 10)
		rect(220, 112, 60, 12)
	case 5:
		rect(150, 145, 16, 65)
		for x := 120; x < 198; x += 14 {
			rect(x, 55, 4, 90)
		}
		rect(120, 53, 78, 5)
		rect(121, 139, 77, 7)
	case 6:
		rect(90, 87, 140, 70)
		rect(55, 108, 35, 24)
		rect(230, 108, 35, 24)
	default:
		rect(110, 70, 100, 125)
		rect(135, 45, 50, 6)
		rect(135, 45, 6, 25)
		rect(179, 45, 6, 25)
		for y := 88; y < 180; y += 22 {
			for x := 126; x < 204; x += 22 {
				draw.Draw(im, image.Rect(x, y, x+9, y+8), &image.Uniform{color.RGBA{237, 237, 251, 255}}, image.Point{}, draw.Src)
			}
		}
	}
	var b bytes.Buffer
	_ = png.Encode(&b, im)
	return b.Bytes()
}
