package integration

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	app "github.com/Maaku050/elabtrack-v2/backend/internal/application/accounts"
	appterms "github.com/Maaku050/elabtrack-v2/backend/internal/application/terms"
	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/accounts"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/terms"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/database"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/persistence/postgres"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/security"
	"github.com/google/uuid"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func batchDatabase(t *testing.T) (context.Context, *database.Postgres, *database.Postgres, *config.Config) {
	t.Helper()
	if os.Getenv("ELABTRACK_BATCH1") != "1" {
		t.Skip("requires strictly owned disposable Batch 1 database")
	}
	cfg, e := config.Load()
	require(t, e == nil && cfg.DB.Name == "elabtrack_v2_batch1_test" && cfg.DB.Port == "54832", "isolated config guard")
	ctx := context.Background()
	runtime, e := database.New(ctx, cfg.DB)
	require(t, e == nil, "runtime connect")
	t.Cleanup(runtime.Close)
	ownerCfg, e := cfg.MigrationConnection()
	require(t, e == nil, "owner config")
	owner, e := database.New(ctx, ownerCfg)
	require(t, e == nil, "owner connect")
	t.Cleanup(owner.Close)
	for _, pair := range []struct {
		db   *database.Postgres
		role string
	}{{runtime, "elabtrack_runtime"}, {owner, "elabtrack_migrator"}} {
		var safe bool
		e = pair.db.Pool.QueryRow(ctx, `SELECT current_database()='elabtrack_v2_batch1_test' AND current_user=$1 AND NOT rolsuper AND NOT rolcreatedb AND NOT rolcreaterole FROM pg_roles WHERE rolname=current_user`, pair.role).Scan(&safe)
		require(t, e == nil && safe, "least privileged named target")
	}
	return ctx, runtime, owner, cfg
}

type capturedMail struct {
	mu       sync.Mutex
	messages []d.Mail
	err      error
}

func (m *capturedMail) SendActivation(_ context.Context, v d.Mail) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messages = append(m.messages, v)
	return "synthetic-provider-acceptance", m.err
}
func (m *capturedMail) lastToken() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.messages) == 0 {
		return ""
	}
	u, _ := url.Parse(m.messages[len(m.messages)-1].ActivationURL)
	return strings.TrimPrefix(u.Fragment, "token=")
}
func (m *capturedMail) count() int { m.mu.Lock(); defer m.mu.Unlock(); return len(m.messages) }

type fixtureObligations struct{}

func (fixtureObligations) Read(context.Context, uuid.UUID) (d.Obligations, error) {
	active, overdue, unreturned, replacements := 3, 2, 9, 4
	fine := int64(12500)
	return d.Obligations{Availability: "AVAILABLE", ActiveBorrowings: &active, OverdueBorrowings: &overdue, UnreturnedUnits: &unreturned, ReplacementUnits: &replacements, FineMinor: &fine}, nil
}
func TestRealAccounts(t *testing.T) {
	ctx, db, owner, cfg := batchDatabase(t)
	repo := postgres.NewAccountsRepository(db.Pool)
	tx := database.NewTxManager(db.Pool)
	hasher := security.NewBcryptHasher(4)
	mail := &capturedMail{}
	svc := app.NewService(repo, tx, hasher, mail, cfg.Accounts.Policy)
	t.Run("empty_down_up_preserves_existing_identities_and_sessions", func(t *testing.T) {
		var history int
		require(t, owner.Pool.QueryRow(ctx, `SELECT (SELECT count(*) FROM account_audit_events)+(SELECT count(*) FROM borrower_profiles)`).Scan(&history) == nil && history == 0, "migration rollback test must precede account history")
		snapshot := func() string {
			var data string
			require(t, owner.Pool.QueryRow(ctx, `SELECT jsonb_build_object('users',(SELECT jsonb_agg(jsonb_build_object('id',id,'name',name,'email',email,'password',password,'role',role,'active',is_active,'created',created_at,'updated',updated_at) ORDER BY id) FROM users),'sessions',(SELECT jsonb_agg(to_jsonb(r) ORDER BY id) FROM refresh_tokens r))::text`).Scan(&data) == nil, "snapshot")
			return data
		}
		before := snapshot()
		m := database.NewMigrator(owner.Pool, "../../migrations")
		require(t, m.Down(ctx) == nil, "empty inventory rollback before account pair")
		require(t, m.Down(ctx) == nil && m.Up(ctx) == nil, "empty paired migration")
		require(t, snapshot() == before, "existing identities and sessions unchanged")
	})
	admin := termsUser(t, ctx, db, user.RoleAdmin)
	staff := termsUser(t, ctx, db, user.RoleStaff)
	borrower := termsUser(t, ctx, db, user.RoleBorrower)
	input := func() d.Input {
		return d.Input{Name: "Synthetic Student", Email: uuid.NewString() + "@students.example.invalid", BorrowerType: "STUDENT", StudentID: "00" + uuid.NewString(), Course: "TEST program"}
	}
	create := func(t *testing.T, in d.Input) d.Record {
		t.Helper()
		r, e := svc.Create(ctx, admin.ID, uuid.NewString(), in, false)
		require(t, e == nil && r.Role == user.RoleBorrower && r.ActivationRequired, "pending borrower create")
		return r
	}
	var student, faculty d.Record
	var firstToken string
	t.Run("individual_student_faculty_text_identity_and_idempotency", func(t *testing.T) {
		in := input()
		key := uuid.NewString()
		var e error
		student, e = svc.Create(ctx, admin.ID, key, in, false)
		require(t, e == nil && student.StudentID == in.StudentID && student.DeliveryStatus == "ACCEPTED", "Student text creation")
		firstToken = mail.lastToken()
		require(t, len(firstToken) == 43, "unpredictable activation length")
		before := mail.count()
		again, e := svc.Create(ctx, admin.ID, key, in, false)
		require(t, e == nil && again.ID == student.ID && mail.count() == before, "retry no duplicate account or email")
		in.Name = "Changed"
		_, e = svc.Create(ctx, admin.ID, key, in, false)
		require(t, errors.Is(e, shared.ErrConflict), "changed payload key rejected")
		faculty = create(t, d.Input{Name: "Synthetic Faculty", Email: uuid.NewString() + "@external.example.invalid", BorrowerType: "FACULTY"})
		require(t, faculty.StudentID == "" && faculty.BorrowerType == "FACULTY", "Faculty noninstitutional no ID")
		sum := sha256.Sum256([]byte(firstToken))
		var hash string
		require(t, owner.Pool.QueryRow(ctx, `SELECT token_hash FROM account_activation_tokens WHERE user_id=$1 AND NOT invalidated`, student.ID).Scan(&hash) == nil && hash == hex.EncodeToString(sum[:]), "only token hash persisted")
	})
	t.Run("domain_duplicates_and_authority", func(t *testing.T) {
		for _, actor := range []uuid.UUID{staff.ID, borrower.ID} {
			_, e := svc.Create(ctx, actor, uuid.NewString(), input(), false)
			require(t, errors.Is(e, shared.ErrForbidden), "non-Admin creation denied")
		}
		in := input()
		in.Email = student.Email
		_, e := svc.Create(ctx, admin.ID, uuid.NewString(), in, false)
		require(t, errors.Is(e, user.ErrEmailAlreadyExists), "duplicate email")
		in = input()
		in.StudentID = student.StudentID
		_, e = svc.Create(ctx, admin.ID, uuid.NewString(), in, false)
		require(t, errors.Is(e, d.ErrStudentIDExists), "duplicate ID")
		in = input()
		in.Email = uuid.NewString() + "@wrong.example.invalid"
		_, e = svc.Create(ctx, admin.ID, uuid.NewString(), in, false)
		require(t, errors.Is(e, d.ErrStudentDomain), "institutional domain")
		empty := app.NewService(repo, tx, hasher, mail, d.Policy{})
		_, e = empty.Create(ctx, admin.ID, uuid.NewString(), input(), false)
		require(t, errors.Is(e, d.ErrDomainsMissing), "missing domains closed")
		in = input()
		in.BorrowerType = "ADMIN"
		_, e = svc.Create(ctx, admin.ID, uuid.NewString(), in, false)
		require(t, e != nil, "privileged escalation denied")
		_, e = svc.List(ctx, borrower.ID, d.Filter{Page: 1, PerPage: 25})
		require(t, errors.Is(e, shared.ErrForbidden), "Borrower directory denied")
		_, e = svc.List(ctx, staff.ID, d.Filter{Page: 1, PerPage: 25})
		require(t, e == nil, "Staff operational read")
		_, e = svc.List(ctx, staff.ID, d.Filter{Page: 1, PerPage: 25, Staff: true})
		require(t, errors.Is(e, shared.ErrForbidden), "Staff privileged directory denied")
	})
	t.Run("bounded_search_filters_and_missing_mail_state", func(t *testing.T) {
		page, e := svc.List(ctx, admin.ID, d.Filter{Page: 1, PerPage: 1, Search: faculty.Email, BorrowerType: "FACULTY", Status: "ACTIVE"})
		require(t, e == nil && page.Total == 1 && len(page.Items) == 1 && page.Items[0].ID == faculty.ID, "bounded Faculty search/filter")
		page, e = svc.List(ctx, admin.ID, d.Filter{Page: 2, PerPage: 1, Search: faculty.Email})
		require(t, e == nil && page.Total == 1 && len(page.Items) == 0, "real empty second page")
		_, e = svc.List(ctx, admin.ID, d.Filter{Page: 1, PerPage: 101})
		require(t, errors.Is(e, shared.ErrInvalidInput), "page ceiling")
		unconfigured := app.NewService(repo, tx, hasher, nil, cfg.Accounts.Policy)
		r, e := unconfigured.Create(ctx, admin.ID, uuid.NewString(), input(), false)
		require(t, e == nil && r.ActivationRequired && r.DeliveryStatus == "UNCONFIGURED", "missing mail does not undo committed pending account or claim delivery")
		failure := &capturedMail{err: d.ErrMailUnknown}
		failedSvc := app.NewService(repo, tx, hasher, failure, cfg.Accounts.Policy)
		r, e = failedSvc.Create(ctx, admin.ID, uuid.NewString(), input(), false)
		require(t, e == nil && r.ActivationRequired && r.DeliveryStatus == "UNKNOWN", "ambiguous provider outcome stays pending and truthful")
	})
	t.Run("single_use_password_and_current_terms", func(t *testing.T) {
		password := uuid.NewString() + "-borrower-password"
		require(t, svc.Activate(ctx, firstToken, password) == nil, "password established")
		require(t, errors.Is(svc.Activate(ctx, firstToken, password), d.ErrActivationInvalid), "activation replay")
		var hash string
		require(t, db.Pool.QueryRow(ctx, `SELECT password FROM users WHERE id=$1`, student.ID).Scan(&hash) == nil && hasher.Compare(hash, password) == nil, "bcrypt password")
		termsSvc := appterms.NewService(postgres.NewTermsRepository(db.Pool), postgres.NewUserRepository(db.Pool), tx)
		status, e := termsSvc.Status(ctx, student.ID)
		require(t, e == nil && status.Acceptance == nil && !status.CanInitiateBorrowing, "no automatic terms acceptance")
		e = tx.Within(ctx, func(ctx context.Context) error { _, e := termsSvc.RequireCurrentAcceptance(ctx, student.ID); return e })
		require(t, errors.Is(e, terms.ErrNotPublished) || errors.Is(e, terms.ErrAcceptanceRequired), "terms still authoritative")
	})
	t.Run("reissue_expiry_inactive_and_no_session_resurrection", func(t *testing.T) {
		pending := create(t, input())
		old := mail.lastToken()
		_, e := svc.Resend(ctx, admin.ID, pending.ID, uuid.NewString())
		require(t, errors.Is(e, d.ErrCooldown), "reissue cooldown")
		now := time.Now().UTC().Add(2 * time.Minute)
		svc.Now = func() time.Time { return now }
		_, e = svc.Resend(ctx, admin.ID, pending.ID, uuid.NewString())
		require(t, e == nil, "safe reissue")
		newToken := mail.lastToken()
		require(t, newToken != old && errors.Is(svc.Activate(ctx, old, "testing-password"), d.ErrActivationInvalid), "old token revoked")
		svc.Now = func() time.Time { return now.Add(25 * time.Hour) }
		require(t, errors.Is(svc.Activate(ctx, newToken, "testing-password"), d.ErrActivationInvalid), "expiry")
		svc.Now = func() time.Time { return time.Now().UTC() }
		pending = create(t, input())
		raw := mail.lastToken()
		sum := sha256.Sum256([]byte(uuid.NewString()))
		_, e = db.Pool.Exec(ctx, `INSERT INTO refresh_tokens(id,user_id,token_hash,expires_at)VALUES($1,$2,$3,now()+interval '1 day')`, uuid.New(), pending.ID, hex.EncodeToString(sum[:]))
		require(t, e == nil, "owned refresh fixture")
		inactive, e := svc.Status(ctx, admin.ID, pending.ID, uuid.NewString(), app.StatusInput{Confirm: true, ExpectedUpdatedAt: pending.UpdatedAt}, false)
		require(t, e == nil && !inactive.IsActive, "inactive")
		require(t, errors.Is(svc.Activate(ctx, raw, "testing-password"), d.ErrActivationInvalid), "inactive activation")
		_, e = svc.Status(ctx, admin.ID, pending.ID, uuid.NewString(), app.StatusInput{Active: true, Confirm: true, ExpectedUpdatedAt: inactive.UpdatedAt}, false)
		require(t, e == nil, "reactivate")
		var live int
		require(t, db.Pool.QueryRow(ctx, `SELECT count(*) FROM refresh_tokens WHERE user_id=$1 AND revoked_at IS NULL`, pending.ID).Scan(&live) == nil && live == 0, "old refresh cannot resurrect after reactivation")
		require(t, errors.Is(svc.Activate(ctx, raw, "testing-password"), d.ErrActivationInvalid), "old activation stays revoked")
	})
	t.Run("bulk_preview_conflicts_selection_idempotency_and_rollback", func(t *testing.T) {
		a, b := input(), input()
		rows := []d.Row{{Number: 2, Input: a}, {Number: 3, Input: b}, {Number: 4, Input: a}}
		preview, e := svc.Prepare(ctx, admin.ID, "CREATE", rows)
		require(t, e == nil && preview.Rows[0].Error != "" && preview.Rows[2].Error != "" && preview.Rows[1].Error == "", "both duplicate rows invalid")
		before := mail.count()
		done, e := svc.Confirm(ctx, admin.ID, preview.ID, []int{3}, true)
		require(t, e == nil && done.Confirmed && done.Rows[1].AccountID != nil && mail.count() == before, "selected pending creation no unintended email")
		again, e := svc.Confirm(ctx, admin.ID, preview.ID, []int{3}, true)
		require(t, e == nil && again.Rows[1].AccountID != nil && *again.Rows[1].AccountID == *done.Rows[1].AccountID, "batch idempotency")
		_, e = svc.Confirm(ctx, admin.ID, preview.ID, []int{2}, true)
		require(t, errors.Is(e, shared.ErrConflict), "changed selection denied")
		a, b = input(), input()
		preview, e = svc.Prepare(ctx, admin.ID, "CREATE", []d.Row{{Number: 2, Input: a}, {Number: 3, Input: b}})
		require(t, e == nil, "preview")
		create(t, b)
		_, e = svc.Confirm(ctx, admin.ID, preview.ID, []int{2, 3}, true)
		require(t, e != nil, "TOCTOU conflict abort")
		matches, e := repo.FindIdentity(ctx, a.Email, a.StudentID)
		require(t, e == nil && len(matches) == 0, "first selected row rolled back")
		state, e := svc.Preview(ctx, admin.ID, preview.ID)
		require(t, e == nil && !state.Confirmed, "failed batch remains uncommitted")
		failSvc := app.NewService(repo, controlledTx{base: tx, after: func() error { return shared.ErrInternal }}, hasher, mail, cfg.Accounts.Policy)
		in := input()
		before = mail.count()
		_, e = failSvc.Create(ctx, admin.ID, uuid.NewString(), in, false)
		require(t, e != nil && mail.count() == before, "forced transaction failure sends no email")
		matches, e = repo.FindIdentity(ctx, in.Email, in.StudentID)
		require(t, e == nil && len(matches) == 0, "forced creation rollback")
	})
	t.Run("student_only_deactivation_ignores_obligations_and_retains_history", func(t *testing.T) {
		svc.SetObligationReader(fixtureObligations{})
		m := create(t, input())
		preview, e := svc.Prepare(ctx, admin.ID, "DEACTIVATE", []d.Row{{Number: 2, Input: d.Input{Name: m.Name, Email: m.Email, BorrowerType: "STUDENT", StudentID: m.StudentID}}, {Number: 3, Input: d.Input{Name: faculty.Name, Email: faculty.Email, BorrowerType: "STUDENT", StudentID: "fake-faculty-id"}}})
		require(t, e == nil && preview.Rows[0].Error == "" && preview.Rows[1].Error != "" && *preview.Rows[0].Obligations.FineMinor == 12500, "Student match warns despite obligations")
		done, e := svc.Confirm(ctx, admin.ID, preview.ID, []int{2}, true)
		require(t, e == nil && done.Confirmed, "obligations do not veto deactivation")
		current, e := repo.Get(ctx, m.ID)
		require(t, e == nil && !current.IsActive && current.StudentID == m.StudentID, "identity retained")
		f, e := repo.Get(ctx, faculty.ID)
		require(t, e == nil && f.IsActive, "Faculty retained")
		audit, e := svc.Audits(ctx, admin.ID, m.ID, 1)
		require(t, e == nil && len(audit) >= 3, "audit history preserved")
		preview, e = svc.Prepare(ctx, admin.ID, "DEACTIVATE", []d.Row{{Number: 2, Input: d.Input{Name: m.Name, Email: m.Email, BorrowerType: "STUDENT", StudentID: m.StudentID}}})
		require(t, e == nil && preview.Rows[0].Error == "Already inactive", "already inactive preview")
		svc.SetObligationReader(d.UnavailableObligations{})
	})
	t.Run("concurrent_identity_and_activation", func(t *testing.T) {
		in := input()
		var winners atomic.Int32
		var wg sync.WaitGroup
		for range 8 {
			wg.Go(func() {
				_, e := svc.Create(ctx, admin.ID, uuid.NewString(), in, false)
				if e == nil {
					winners.Add(1)
				} else if !errors.Is(e, user.ErrEmailAlreadyExists) {
					t.Error("unexpected concurrent creation result")
				}
			})
		}
		wg.Wait()
		require(t, winners.Load() == 1, "one concurrent creation winner")
		matches, e := repo.FindIdentity(ctx, in.Email, in.StudentID)
		require(t, e == nil && len(matches) == 1, "one identity persisted")
		raw := mail.lastToken()
		winners.Store(0)
		for range 6 {
			wg.Go(func() {
				e := svc.Activate(ctx, raw, "concurrent-password")
				if e == nil {
					winners.Add(1)
				} else if !errors.Is(e, d.ErrActivationInvalid) {
					t.Error("unexpected activation race result")
				}
			})
		}
		wg.Wait()
		require(t, winners.Load() == 1, "atomic single activation winner")
	})
	t.Run("database_constraints_history_and_nonempty_rollback", func(t *testing.T) {
		_, e := db.Pool.Exec(ctx, `UPDATE borrower_profiles SET student_id='bad' WHERE user_id=$1`, student.ID)
		require(t, e != nil, "runtime cannot rewrite Student identity")
		_, e = db.Pool.Exec(ctx, `UPDATE account_audit_events SET action='tampered'`)
		require(t, e != nil, "runtime audit immutable")
		_, e = owner.Pool.Exec(ctx, `UPDATE account_audit_events SET action='tampered'`)
		require(t, e != nil, "owner audit immutable trigger")
		_, e = db.Pool.Exec(ctx, `DELETE FROM account_operation_receipts`)
		require(t, e != nil, "runtime cannot erase receipt")
		_, e = owner.Pool.Exec(ctx, `UPDATE users SET role='STAFF' WHERE id=$1`, student.ID)
		require(t, e != nil, "classified borrower role integrity")
		m := database.NewMigrator(owner.Pool, "../../migrations")
		require(t, m.Down(ctx) == nil, "empty inventory rollback before account pair")
		require(t, m.Down(ctx) != nil, "history-bearing rollback denied")
		var count int
		require(t, owner.Pool.QueryRow(ctx, `SELECT count(*) FROM schema_migrations WHERE version='000006_account_management'`).Scan(&count) == nil && count == 1, "migration tracking preserved")
		require(t, m.Up(ctx) == nil, "restore empty inventory after account history gate")
	})
}

type gatedMail struct{ entered, release chan struct{} }

func (m gatedMail) SendActivation(context.Context, d.Mail) (string, error) {
	close(m.entered)
	<-m.release
	return "TEST-ACCEPTED", nil
}
func TestRealAccountDeliveryRace(t *testing.T) {
	ctx, db, _, cfg := batchDatabase(t)
	users := postgres.NewUserRepository(db.Pool)
	operator := user.NewUser(uuid.NewString()+"@example.invalid", "TEST delivery operator", "not-a-login")
	operator.Role = user.RoleAdmin
	for i := 0; i < 4; i++ {
		operator.ID[i] = 0
	}
	require(t, users.Create(ctx, operator) == nil, "owned low-sorting actor")
	repo := postgres.NewAccountsRepository(db.Pool)
	mail := gatedMail{make(chan struct{}), make(chan struct{})}
	svc := app.NewService(repo, database.NewTxManager(db.Pool), security.NewBcryptHasher(4), mail, cfg.Accounts.Policy)
	input := d.Input{Name: "TEST delivery race", Email: uuid.NewString() + "@external.example.invalid", BorrowerType: "FACULTY"}
	type created struct {
		record d.Record
		err    error
	}
	creation := make(chan created, 1)
	go func() {
		r, e := svc.Create(ctx, operator.ID, uuid.NewString(), input, false)
		creation <- created{r, e}
	}()
	select {
	case <-mail.entered:
	case <-time.After(5 * time.Second):
		t.Fatal("bounded delivery barrier")
	}
	matches, e := repo.FindIdentity(ctx, input.Email, "")
	require(t, e == nil && len(matches) == 1, "durable pending account before mail")
	target := matches[0]
	require(t, operator.ID.String() < target.ID.String(), "fixture actor sorts before target")
	status := make(chan error, 1)
	go func() {
		_, e := svc.Status(ctx, operator.ID, target.ID, uuid.NewString(), app.StatusInput{Confirm: true, ExpectedUpdatedAt: target.UpdatedAt}, false)
		status <- e
	}()
	waitForDatabaseLock(t, ctx, db)
	close(mail.release)
	select {
	case result := <-creation:
		require(t, result.err == nil && result.record.DeliveryStatus == "ACCEPTED", "mail audit must commit without actor-target FK deadlock")
	case <-time.After(5 * time.Second):
		t.Fatal("delivery deadlock")
	}
	select {
	case e := <-status:
		require(t, e == nil || errors.Is(e, shared.ErrConflict), "status observes serialization/version conflict, never lock-order internal error")
	case <-time.After(5 * time.Second):
		t.Fatal("status deadlock")
	}
	audit, e := repo.Audits(ctx, target.ID, 1)
	require(t, e == nil && len(audit) >= 2, "creation and mail audit preserved")
}
