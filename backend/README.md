# eLabTrack V2 API foundation

Go Fiber v3 with domain/application/infrastructure/HTTP boundaries, composed in `internal/bootstrap/`. Read [architecture](../docs/ARCHITECTURE.md) and [foundation audit](../docs/project/FOUNDATION_AUDIT.md). Current auth/user code is retained starter infrastructure; it does not establish eLabTrack product policies.

Requires Go 1.27.1 or newer and a dedicated local PostgreSQL database. Copy `.env.example` to `.env`, configure matching local credentials. The published signing key is permitted only for local development; production requires explicit strong signing material. The authoritative module/import path is `github.com/Maaku050/elabtrack-v2/backend`, verified against `origin` during Phase 0 closure.

```bash
go run ./cmd/api                    # database required at bootstrap
go run ./cmd/api --migrate-up       # explicit, authorized schema change
go run ./cmd/api --migrate-down     # destructive rollback; use deliberately
go run ./cmd/api --migrate-status   # read-only status, no bookkeeping creation
go run ./cmd/api --seed             # explicit; development only
go run ./cmd/api --sessions-cleanup # explicit; at most 1000 terminal rows older than seven days
go fmt ./...
go vet ./...
go test ./...
```

Retained endpoints: `GET /api/v1/health`; `POST /api/v1/auth/{register,login,refresh,logout}`; `GET /api/v1/auth/me`; `GET/PATCH /api/v1/users/me`; and admin-gated `GET /api/v1/users/`. Production does not mount public registration; development/test retain local convenience only. Protected requests resolve current active accounts from PostgreSQL, and the temporary `user/admin` roles are unapproved as final V2 policy. Self PATCH accepts display name only; generic user listing requires the current database admin role. See the [security backlog](../docs/project/PHASE1_SECURITY_BACKLOG.md). No business-domain tables or workflows exist.

Access JWT verification pins HS256, configured issuer, fixed `elabtrack-v2-api` audience, explicit `purpose=access`, required exp/iat/nbf and consistent nonzero uid/subject. Database role/status remain authoritative. Refresh tokens retain crypto/rand 32-byte opaque hexadecimal credentials; only SHA-256 digests reach the repository. Rotation uses one PostgreSQL transaction with a session row lock, safe current-account lock and conditional consumption. Logout revokes the supplied session idempotently; it does not immediately revoke access JWTs or a rotated successor.

Migrations are paired SQL for `users` and `refresh_tokens` only, plus runner bookkeeping. New 000003 replaces the raw token column with a unique hash and adds revocation/replacement metadata and cleanup indexes. Both up/down invalidate all sessions; re-login and coordinated schema/application versions are required. No migrations were executed in Phase 1C. Normal host/container startup runs neither migrations, seeds nor cleanup. Container commands use `./api` as their entrypoint; pass maintenance actions explicitly. Production migrations are a separate release job/step. Atomic migration bookkeeping/locks, local registration user/session write atomicity, browser transport and live PostgreSQL tests remain outstanding. See the [Phase 1 foundation report](../docs/project/PHASE1_FOUNDATION.md). Do not deploy this retained auth foundation as production-ready.

`--sessions-cleanup` removes one batch of at most 1000 records whose earlier expiry/revocation is over seven days old. An operator must arrange regular runs (initially daily, repeating batches to drain the eligible backlog), monitor table size and adjust cadence to measured rotation volume. This is technical session-record retention, not institutional borrowing/audit policy. No scheduler is installed or run. No token/hash/user detail is printed, only the deleted count.
