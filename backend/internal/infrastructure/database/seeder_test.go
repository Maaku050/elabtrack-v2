package database

import (
	"context"
	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	"strings"
	"testing"
)

func TestSeederGuardBeforeIO(t *testing.T) {
	for _, env := range []config.Environment{config.Production, config.Test, "unknown", config.Development} {
		for _, requested := range []bool{false, true} {
			if env == config.Development && requested {
				continue
			}
			seed := NewSeeder(nil, "nonexistent-seeds-directory", env)
			err := seed.Run(context.Background(), requested)
			if err == nil || !strings.HasPrefix(err.Error(), "seed:") {
				t.Fatal("guard must reject before filesystem/database access")
			}
		}
	}
	if err := config.CheckDevelopmentSeed(config.Development, true); err != nil {
		t.Fatal(err)
	}
}
