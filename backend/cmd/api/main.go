package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Maaku050/elabtrack-v2/backend/internal/bootstrap"
)

// CLI flags. The server runs by default; passing any of these flags runs
// the corresponding one-shot command instead and exits.
func main() {
	migrateUp := flag.Bool("migrate-up", false, "apply pending database migrations and exit")
	migrateDown := flag.Bool("migrate-down", false, "roll back the latest migration and exit")
	migrateCreate := flag.String("migrate-create", "", "create a new migration pair (provide a name) and exit")
	seed := flag.Bool("seed", false, "run database seeders and exit")
	flag.Parse()

	ctx := context.Background()

	// One-shot commands that don't need the full HTTP server wiring.
	if *migrateUp || *migrateDown || *migrateCreate != "" || *seed {
		if err := runCommand(ctx, *migrateUp, *migrateDown, *migrateCreate, *seed); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	app, err := bootstrap.New(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to bootstrap application: %v\n", err)
		os.Exit(1)
	}
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "application error: %v\n", err)
		os.Exit(1)
	}
}

func runCommand(ctx context.Context, up, down bool, createName string, seed bool) error {
	app, err := bootstrap.New(ctx)
	if err != nil {
		return err
	}
	defer app.Logger().Sync()
	defer func() { _ = app.Shutdown(context.Background()) }()

	switch {
	case createName != "":
		if _, err := app.Migrator().Create(sanitizeName(createName)); err != nil {
			return err
		}
	case up:
		return app.Migrator().Up(ctx)
	case down:
		return app.Migrator().Down(ctx)
	case seed:
		return app.Seeder().Run(ctx)
	}
	return nil
}

// sanitizeName keeps migration file names safe.
func sanitizeName(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '-' || r == ' ':
			b.WriteRune('_')
		case r == '_':
			b.WriteRune('_')
		}
	}
	return b.String()
}
