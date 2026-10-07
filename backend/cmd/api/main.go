package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/Maaku050/elabtrack-v2/backend/internal/bootstrap"
	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/database"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/logger"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/persistence/postgres"
)

type command struct {
	up, down, status, seed, cleanup, adopt bool
	create                                 string
}

func main() {
	cmd, err := parseCommand(os.Args[1:])
	if errors.Is(err, flag.ErrHelp) {
		return
	}
	if err == nil {
		err = run(context.Background(), cmd)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
func parseCommand(args []string) (command, error) {
	var cmd command
	flags := flag.NewFlagSet("api", flag.ContinueOnError)
	// Flag errors can contain supplied values; print only our safe errors.
	flags.SetOutput(io.Discard)
	flags.BoolVar(&cmd.adopt, "migrate-adopt-legacy", false, "attest and adopt verified local foundation history (development only)")
	flags.BoolVar(&cmd.up, "migrate-up", false, "apply pending migrations")
	flags.BoolVar(&cmd.down, "migrate-down", false, "roll back the latest migration")
	flags.BoolVar(&cmd.status, "migrate-status", false, "show migration status without schema changes")
	flags.StringVar(&cmd.create, "migrate-create", "", "create a paired SQL scaffold")
	flags.BoolVar(&cmd.seed, "seed", false, "explicitly run local development seeds")
	flags.BoolVar(&cmd.cleanup, "sessions-cleanup", false, "delete at most 1000 sessions terminal for over seven days")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			fmt.Println("api [--migrate-up | --migrate-down | --migrate-status | --migrate-adopt-legacy | --migrate-create NAME | --seed | --sessions-cleanup]")
			return cmd, flag.ErrHelp
		}
		return cmd, errors.New("command: invalid arguments; use --help")
	}
	count := 0
	for _, active := range []bool{cmd.up, cmd.down, cmd.status, cmd.adopt, cmd.create != "", cmd.seed, cmd.cleanup} {
		if active {
			count++
		}
	}
	if count > 1 || flags.NArg() > 0 {
		return cmd, errors.New("command: choose exactly one action, without positional arguments")
	}
	if cmd.create != "" && strings.Trim(sanitizeName(cmd.create), "_") == "" {
		return cmd, errors.New("migrate-create: name must contain letters or digits")
	}
	// An explicitly empty create flag must not accidentally launch a server.
	flags.Visit(func(f *flag.Flag) {
		if f.Name == "migrate-create" && cmd.create == "" {
			count = -1
		}
	})
	if count < 0 {
		return cmd, errors.New("migrate-create: name is required")
	}
	return cmd, nil
}
func run(ctx context.Context, cmd command) error {
	if !cmd.up && !cmd.down && !cmd.status && !cmd.adopt && cmd.create == "" && !cmd.seed && !cmd.cleanup {
		app, err := bootstrap.New(ctx)
		if err != nil {
			return err
		}
		return app.Run()
	}
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if cmd.seed {
		if err := config.CheckDevelopmentSeed(cfg.App.Env, true); err != nil {
			return err
		}
	}
	if cmd.create != "" {
		_, err := database.NewMigrator(nil, "migrations").Create(sanitizeName(cmd.create))
		return err
	}
	if cmd.adopt && cfg.App.Env != config.Development {
		return errors.New("migration: legacy adoption permitted only in development")
	}
	dbCfg := cfg.DB
	if !cmd.cleanup {
		dbCfg, err = cfg.MigrationConnection()
		if err != nil {
			return err
		}
	}
	db, err := database.New(ctx, dbCfg)
	if err != nil {
		return err
	}
	defer db.Close()
	if cmd.cleanup {
		count, err := postgres.NewAuthRepository(db.Pool).Cleanup(ctx, time.Now().UTC().Add(-7*24*time.Hour), 1000)
		if err != nil {
			return errors.New("sessions: cleanup failed")
		}
		fmt.Printf("sessions: deleted %d terminal records\n", count)
		return nil
	}
	log := logger.Must(cfg.Log.Level, cfg.Log.Format)
	defer log.Sync()
	migrator := database.NewMigrator(db.Pool, "migrations", log.Logger)
	switch {
	case cmd.adopt:
		err = migrator.AdoptLegacy(ctx)
	case cmd.up:
		err = migrator.Up(ctx)
	case cmd.down:
		err = migrator.Down(ctx)
	case cmd.status:
		err = migrator.Status(ctx)
	case cmd.seed:
		return database.NewSeeder(db.Pool, "seeds", cfg.App.Env).Run(ctx, true)
	}
	// SQL errors may include server details. Do not send them to startup logs.
	if err != nil {
		var safe *database.MigrationError
		if errors.As(err, &safe) {
			return safe
		}
		return errors.New("migration: command failed")
	}
	return nil
}
func sanitizeName(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '-' || r == ' ' || r == '_':
			b.WriteRune('_')
		}
	}
	return b.String()
}
