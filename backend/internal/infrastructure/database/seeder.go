package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Seeder runs idempotent seed SQL files for development environments.
// It should never be used in production.
type Seeder struct {
	pool *pgxpool.Pool
	dir  string
	env  config.Environment
}

// NewSeeder constructs a Seeder pointing at the given seeds dir.
func NewSeeder(pool *pgxpool.Pool, dir string, env config.Environment) *Seeder {
	return &Seeder{pool: pool, dir: dir, env: env}
}

// Run executes every *.sql file in the seeds directory in alphabetical order.
func (s *Seeder) Run(ctx context.Context, requested bool) error {
	if err := config.CheckDevelopmentSeed(s.env, requested); err != nil {
		return err
	}
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return fmt.Errorf("read seeds dir: %w", err)
	}
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".sql" {
			continue
		}
		path := filepath.Join(s.dir, e.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read seed %s: %w", path, err)
		}
		if _, err := s.pool.Exec(ctx, string(content)); err != nil {
			return fmt.Errorf("apply development seed %s: database operation failed", e.Name())
		}
		fmt.Printf("applied seed %s\n", e.Name())
	}
	return nil
}
