# eLabTrack V2 API foundation

Go Fiber v3 with domain/application/infrastructure/HTTP boundaries, composed in `internal/bootstrap/`. Read [architecture](../docs/ARCHITECTURE.md) and [foundation audit](../docs/project/FOUNDATION_AUDIT.md). Current auth/user code is retained starter infrastructure; it does not establish eLabTrack product policies.

Requires Go 1.27.1 or newer and a dedicated local PostgreSQL database. Copy `.env.example` to `.env`, configure matching local credentials. The published signing key is permitted only for local development; production requires explicit strong signing material. The authoritative module/import path is `github.com/Maaku050/elabtrack-v2/backend`, verified against `origin` during Phase 0 closure.

```bash
go run ./cmd/api                    # database required at bootstrap
go run ./cmd/api --migrate-up       # explicit, authorized schema change
go run ./cmd/api --migrate-down     # destructive rollback; use deliberately
go run ./cmd/api --migrate-status   # read-only status, no bookkeeping creation
go run ./cmd/api --seed             # explicit; development only
go fmt ./...
go vet ./...
go test ./...
```

Retained endpoints: `GET /api/v1/health`; `POST /api/v1/auth/{register,login,refresh,logout}`; `GET /api/v1/auth/me`; `GET/PATCH /api/v1/users/me`; and admin-gated `GET /api/v1/users/`. Public registration and `user/admin` roles are inherited, unapproved for V2, and tracked in the [security backlog](../docs/project/PHASE1_SECURITY_BACKLOG.md). No business-domain tables or workflows exist.

Migrations are paired SQL for `users` and `refresh_tokens` only, plus runner bookkeeping. Normal host/container startup runs neither migrations nor seeds. Container commands use `./api` as their entrypoint; pass `--migrate-up`, `--migrate-down`, or `--migrate-status` explicitly. Production migrations are a separate release job/step. Atomic bookkeeping and migration locks remain Phase 1 work, alongside token storage/rotation. See [Phase 1A configuration and environment documentation](../docs/project/PHASE1_FOUNDATION.md). Do not deploy this retained auth foundation as production-ready.
