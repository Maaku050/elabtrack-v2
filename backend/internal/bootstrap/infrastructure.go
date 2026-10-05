package bootstrap

import (
	"context"

	"github.com/fullstacktemplate/backend/internal/config"
	"github.com/fullstacktemplate/backend/internal/infrastructure/database"
	"github.com/fullstacktemplate/backend/internal/infrastructure/logger"
	"github.com/fullstacktemplate/backend/internal/infrastructure/persistence/postgres"
	"github.com/fullstacktemplate/backend/internal/infrastructure/security"
)

// Infrastructure bundles the infrastructure-layer singletons created
// during bootstrap. They are passed to the dependency container which
// wires application services and handlers.
type Infrastructure struct {
	Config   *config.Config
	Logger   *logger.Logger
	DB       *database.Postgres
	Tx       *database.TxManager
	Migrator *database.Migrator
	Seeder   *database.Seeder
	Health   *database.HealthChecker
}

// initInfrastructure constructs all infrastructure singletons.
func initInfrastructure(ctx context.Context, cfg *config.Config) (*Infrastructure, error) {
	log := logger.Must(cfg.Log.Level, cfg.Log.Format)

	db, err := database.New(ctx, cfg.DB)
	if err != nil {
		log.Fatal("failed to init database")
		return nil, err
	}

	return &Infrastructure{
		Config:   cfg,
		Logger:   log,
		DB:       db,
		Tx:       database.NewTxManager(db.Pool),
		Migrator: database.NewMigrator(db.Pool, "migrations"),
		Seeder:   database.NewSeeder(db.Pool, "seeds"),
		Health:   database.NewHealthChecker(db.Pool),
	}, nil
}

// newSecurityAdapters constructs the password hasher and token issuer.
func newSecurityAdapters(cfg *config.Config) (*security.BcryptHasher, *security.JWTIssuer) {
	return security.NewBcryptHasher(0), security.NewJWTIssuer(cfg.JWT)
}

// newRepositories constructs the PostgreSQL repository implementations.
func newRepositories(db *database.Postgres) (*postgres.UserRepository, *postgres.AuthRepository) {
	return postgres.NewUserRepository(db.Pool), postgres.NewAuthRepository(db.Pool)
}
