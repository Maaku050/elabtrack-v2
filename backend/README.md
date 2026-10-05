# eLabTrack V2 API foundation

Go Fiber v3 with domain/application/infrastructure/HTTP boundaries, composed in `internal/bootstrap/`. Read [architecture](../docs/ARCHITECTURE.md) and [foundation audit](../docs/project/FOUNDATION_AUDIT.md). Current auth/user code is retained starter infrastructure; it does not establish eLabTrack product policies.

Requires Go 1.27.1 or newer and a dedicated local PostgreSQL database. Copy `.env.example` to `.env`, configure credentials, and replace the development JWT secret. The authoritative module/import path is `github.com/Maaku050/elabtrack-v2/backend`, verified against `origin` during Phase 0 closure.

```bash
go run ./cmd/api                    # database required at bootstrap
go run ./cmd/api --migrate-up       # dedicated development database only
go run ./cmd/api --migrate-down     # destructive rollback; use deliberately
go run ./cmd/api --seed             # synthetic development accounts only
go fmt ./...
go vet ./...
go test ./...
```

Retained endpoints: `GET /api/v1/health`; `POST /api/v1/auth/{register,login,refresh,logout}`; `GET /api/v1/auth/me`; `GET/PATCH /api/v1/users/me`; and admin-gated `GET /api/v1/users/`. Public registration and `user/admin` roles are inherited, unapproved for V2, and tracked in the [security backlog](../docs/project/PHASE1_SECURITY_BACKLOG.md). No business-domain tables or workflows exist.

Migrations are paired SQL for `users` and `refresh_tokens` only, plus runner bookkeeping. The Dockerfile currently runs migrations automatically at startup; host API startup does not. Production migration safety and token storage/rotation are deferred to Phase 1. Do not deploy this retained auth foundation as production-ready.
