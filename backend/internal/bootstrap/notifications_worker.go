package bootstrap

import (
	"context"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"go.uber.org/zap"
	"time"
)

func (a *App) runNotifications(ctx context.Context) {
	defer close(a.notificationDone)
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		cycle, cancel := context.WithTimeout(ctx, 10*time.Second)
		for batch := 0; batch < 10; batch++ {
			n, e := a.notifications.Sweep(cycle, 100)
			if e != nil {
				if ctx.Err() == nil {
					class, op := shared.FailureDetails(e)
					a.log.Error("notification reconciliation failed", zap.String("event", "notifications.reconciliation_failed"), zap.String("error_class", class), zap.String("operation", op))
				}
				break
			}
			if n < 100 {
				break
			}
		}
		cancel()
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}
