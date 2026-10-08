package integration

import (
	"context"
	"encoding/json"
	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/database"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/persistence/postgres"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/security"
	"github.com/google/uuid"
	"os"
	"testing"
)

// Explicit isolated browser fixtures only; never a product provisioning endpoint.
// Immutable terms history requires disposal of the whole owned test database;
// the fixture helper never deletes acceptance or publication evidence.
func TestPhase4BBrowserFixtures(t *testing.T) {
	if os.Getenv("ELABTRACK_PHASE4B") != "1" || os.Getenv("PHASE4B_BROWSER_FIXTURES") != "1" {
		t.Skip("requires explicit isolated browser fixture generation")
	}
	const path = "/tmp/elabtrack-phase4b-fixtures.json"
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	require(t, err == nil, "private fixture file must not exist")
	defer file.Close()
	cfg, err := config.Parse(map[string]string{"APP_ENV": "test", "JWT_SECRET": os.Getenv("JWT_SECRET"), "DB_HOST": "127.0.0.1", "DB_PORT": "35432", "DB_NAME": "elabtrack_v2_phase4b_test", "DB_USER": "elabtrack_runtime", "DB_PASSWORD": os.Getenv("DB_PASSWORD")})
	require(t, err == nil, "isolated config")
	ctx := context.Background()
	db, err := database.New(ctx, cfg.DB)
	require(t, err == nil, "isolated connection")
	defer db.Close()
	var safe bool
	err = db.Pool.QueryRow(ctx, `SELECT current_database()='elabtrack_v2_phase4b_test' AND current_user='elabtrack_runtime' AND NOT rolsuper FROM pg_roles WHERE rolname=current_user`).Scan(&safe)
	require(t, err == nil && safe, "isolated runtime identity")
	fixtures := map[string]map[string]string{}
	ids := []uuid.UUID{}
	complete := false
	defer func() {
		if !complete {
			_, _ = db.Pool.Exec(ctx, `DELETE FROM users WHERE id=ANY($1::uuid[])`, ids)
			_ = os.Remove(path)
		}
	}()
	repo := postgres.NewUserRepository(db.Pool)
	for _, kind := range []string{"BORROWER", "STAFF", "ADMIN", "INACTIVE"} {
		password := uuid.NewString() + uuid.NewString()[:16]
		hash, err := security.NewBcryptHasher(0).Hash(password)
		require(t, err == nil, "random fixture hash")
		u := user.NewUser("phase4b-browser-"+uuid.NewString()+"@example.invalid", "Synthetic "+kind, hash)
		u.Role = user.Role(kind)
		if kind == "INACTIVE" {
			u.Role = user.RoleBorrower
			u.IsActive = false
		}
		require(t, repo.Create(ctx, u) == nil, "owned fixture creation")
		ids = append(ids, u.ID)
		fixtures[kind] = map[string]string{"id": u.ID.String(), "email": u.Email, "password": password}
	}
	require(t, json.NewEncoder(file).Encode(fixtures) == nil, "private fixture write")
	complete = true
	t.Log("four isolated browser fixtures created; generated passwords retained only in mode-0600 /tmp file; cleanup required")
}
