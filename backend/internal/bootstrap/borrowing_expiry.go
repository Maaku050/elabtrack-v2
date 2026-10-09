package bootstrap

import (
	"context"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"go.uber.org/zap"
	"time"
)

// Deadlines are durable; the timer only initiates bounded reconciliation.
func (a *App) runExpiry(ctx context.Context) {
	defer close(a.expiryDone)
	timer := time.NewTicker(30 * time.Second)
	defer timer.Stop()
	for {
		cycle, cancel := context.WithTimeout(ctx, 10*time.Second)
		count, e := a.borrowing.Sweep(cycle, 100)
		cancel()
		if e != nil && ctx.Err() == nil {
			class, op := shared.FailureDetails(e)
			a.log.Error("pending expiration sweep failed", zap.String("event", "borrowing.expiry_failed"), zap.String("error_class", class), zap.String("operation", op))
		} else if count > 0 {
			a.log.Info("pending reservations expired", zap.String("event", "borrowing.expiry_completed"), zap.Int("count", count))
		}
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
		}
	}
}
