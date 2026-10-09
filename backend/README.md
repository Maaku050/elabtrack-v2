# eLabTrack V2 API foundation

Go Fiber v3 with domain/application/infrastructure/HTTP boundaries, composed in `internal/bootstrap/`. Read [architecture](../docs/ARCHITECTURE.md) and [foundation audit](../docs/project/FOUNDATION_AUDIT.md). Phase 4A applies confirmed BORROWER/STAFF/ADMIN policy to the preserved Phase 1 authentication infrastructure. Business/account-management modules remain deferred.

Requires Go 1.27.1 or newer and a dedicated local PostgreSQL database. Copy `.env.example` to `.env`, configure matching local credentials. The published signing key is permitted only for local development; production requires explicit strong signing material. The authoritative module/import path is `github.com/Maaku050/elabtrack-v2/backend`, verified against `origin` during Phase 0 closure.

```bash
go run ./cmd/api                    # database required at bootstrap
go run ./cmd/api --migrate-up       # explicit, authorized schema change
go run ./cmd/api --migrate-down     # destructive rollback; use deliberately
go run ./cmd/api --migrate-status   # locked, read-only checksum/status inspection
go run ./cmd/api --migrate-adopt-legacy # known local history only; read docs first
go run ./cmd/api --seed             # explicit; development only
go run ./cmd/api --sessions-cleanup # explicit; at most 1000 terminal rows older than seven days
go fmt ./...
go vet ./...
go test ./...
```

Retained endpoints: `GET /api/v1/health` (liveness); `GET /api/v1/ready` (dependency readiness); `POST /api/v1/auth/{login,refresh,logout}`; `GET /api/v1/auth/me`; `GET/PATCH /api/v1/users/me`; and admin-gated `GET /api/v1/users/`. Public registration is unavailable in every environment. Protected requests resolve current active accounts from PostgreSQL, and roles are exactly `BORROWER`, `STAFF`, `ADMIN`. Self PATCH accepts display name only; generic user listing requires the current database ADMIN permission. See the [security backlog](../docs/project/PHASE1_SECURITY_BACKLOG.md). No business-domain tables or workflows exist.

Access JWT verification pins HS256, configured issuer, fixed `elabtrack-v2-api` audience, explicit `purpose=access`, required exp/iat/nbf and consistent nonzero uid/subject. Database role/status remain authoritative. Refresh tokens retain crypto/rand 32-byte opaque hexadecimal credentials; only SHA-256 digests reach the repository. Rotation uses one PostgreSQL transaction with a session row lock, safe current-account lock and conditional consumption. Logout revokes the cookie's session idempotently and clears that cookie; it does not immediately revoke access JWTs or an independently rotated successor.

Login/refresh issue a host-only `elabtrack_v2_refresh` cookie: HttpOnly, Path=/api/v1/auth, SameSite=Lax, Secure in production, Expires matching the new session. Only development/test permit Secure=false for local HTTP. All cookie-changing auth POSTs require exact Origin membership in validated ALLOWED_ORIGINS, with no missing/null Origin exception. Cookie reads replace raw JSON body input for refresh/logout. Session JSON contains only access_token, expires_at, token_type and safe user fields; all session responses use Cache-Control: no-store. No failed refresh clears a potentially newer cookie; deliberate logout clears it. Logout clears it even on storage failure, but such a failure cannot establish server-side revocation. Protected APIs still require a memory-held bearer token; the cookie alone cannot authorize them. Production SPA/API must be same-site HTTPS. See the [Phase 1D report](../docs/project/PHASE1_FOUNDATION.md#phase-1d--browser-session-architecture) for retry, logout acknowledgement and cross-tab limitations.

HTTP perimeter configuration: empty `TRUSTED_PROXIES` uses the direct socket peer. Literal IPs/CIDRs enable only the shared ClientInfo/ClientIP boundary; it validates the complete bounded X-Forwarded-For chain and walks right-to-left through trusted hops. Invalid chains fall back to the peer. X-Real-IP/Forwarded/forwarded host and alternate scheme headers are ignored; Fiber proxy parsing remains disabled. Only trusted immediate peers with a single exact `X-Forwarded-Proto: https`, or actual TLS, establish HTTPS for production API HSTS. Operators must sanitize XFF/XFP, choose narrow proxy ranges and prevent untrusted bypass. IP is an abuse signal, never account identity. Without configuration, nginx/other proxied clients share the peer's budget.

Process-local fixed windows default to one minute: `LOGIN_RATE_LIMIT_MAX=10`, `REFRESH_RATE_LIMIT_MAX=60`, `REGISTER_RATE_LIMIT_MAX=5`, `RATE_LIMIT_MAX=120` for the rest including logout/current-user/admin reads and public health/readiness reads. They count successes/failures, use no email/token keys, expire locally and reset on restart; multiple instances do not share enforcement. OPTIONS is exempt. 429 keeps the current safe envelope with conservative Retry-After; ordinary cookie refresh/retry fits the refresh budget. CORS allows exact canonical configured origins, GET/POST/PATCH/OPTIONS, Content-Type/Accept/Authorization/X-Request-ID, credentials, and 300-second preflight caching. Frontend URL must be in the allowlist. Login failures hide account/status/lookup distinctions; constant-time network/account-enumeration resistance is not claimed.

Global body ceiling defaults to 1 MiB, auth POST ceiling is 16 KiB. Login/self PATCH require JSON; strict binding rejects unknown fields/multiple documents. Read/write/idle are 10/15/60 seconds with explicit 8 KiB header buffer. Future uploads/reports need separately reviewed limits/timeouts. API emits nosniff/frame/referrer/permissions headers; document CSP belongs to the SPA server. Panic diagnostics retain type/source stack without raw panic/error values; HTTP framework/wrapped-validation errors use generic text. Health/readiness use minimal enveloped service status, no-store. Phase 1F implements the error contract, structured logging and correlation. Phase 1G verifies actual socket 413/431 limits and fixes final parser completion logging; see the [integration report](../docs/project/PHASE1_FOUNDATION.md#phase-1g--real-integration-verification).

Migrations are paired SQL for `users` and `refresh_tokens` only, plus runner bookkeeping. New 000003 replaces the raw token column with a unique hash and adds revocation/replacement metadata and cleanup indexes. Both up/down invalidate all sessions; re-login and coordinated schema/application versions are required. No migrations were executed in Phase 1C/1D; Phase 1D adds none. Normal host/container startup runs neither migrations, seeds nor cleanup. Container commands use `./api` as their entrypoint; pass maintenance actions explicitly. Production migrations are a separate release job/step. Phase 1H implements atomic migration bookkeeping, paired checksums and PostgreSQL advisory exclusion, and verifies real API/auth/browser behavior with runtime-only DML credentials. Registration service code remains internal and unmounted; no public provisioning endpoint is introduced. See the [Phase 1 foundation report](../docs/project/PHASE1_FOUNDATION.md). Do not deploy this retained auth foundation as production-ready.

`--sessions-cleanup` removes one batch of at most 1000 records whose earlier expiry/revocation is over seven days old. An operator must arrange regular runs (initially daily, repeating batches to drain the eligible backlog), monitor table size and adjust cadence to measured rotation volume. This is technical session-record retention, not institutional borrowing/audit policy. No scheduler is installed or run. No token/hash/user detail is printed, only the deleted count.

Phase 1F uses the retained JSON success envelope and standard nested error code/message/requestId/optional fields. Input validation is 400; unknown/internal errors are generic. Server-owned UUIDv4 IDs propagate through standard context; existing Zap logs safe structured final status/duration/type/operation without credentials, raw URLs/query or headers/bodies. /health does not ping DB; /ready uses the existing checker and safe standard 503. See [current API/log contracts](../docs/API_CONTRACTS.md). Phase 1G records real local socket/proxy/browser/DB/Docker evidence; reproduce with the [isolated integration instructions](../integration/README.md). Completed Phase 1I Web Locks/non-secret BroadcastChannel coordination is preserved by Phase 4A. No phase is automatically authorized by these instructions.

Migration CLI actions and development seeds select typed `MIGRATION_DATABASE_URL`;
API startup and session cleanup use only `DATABASE_URL`/`DB_*`. Production requires
separate migration credentials targeting the same host/port/database and explicit
verify-full TLS; it never falls back. Development/test retain a documented local
single-connection fallback when the migration URL is absent. Fresh Compose uses
the separate roles by default; its tools profile runs an explicit migrator CLI,
not API startup. Apply foundation runtime grants separately, as documented in
[the integration guide](../integration/README.md). Existing volumes are not
silently re-owned or repaired. Immutable SQL 000001–000003 remains unchanged.

Phase 4A adds paired 000004_product_roles: legacy `user`→`BORROWER`, legacy `admin`→`ADMIN`, new valid `STAFF`, default BORROWER. It preserves all non-role identity/account/session data. Down refuses if any Staff would lose its unrepresentable role; coordinate application rollback and explicit account disposition before retry. No table or runtime grant is added. Public signup remains absent; development seeds now contain no accounts or passwords. Current-account status/role resolution and a central permission guard protect the existing directory. See [Phase 4A](../docs/project/PHASE4A_REPORT.md) and [isolated real-PG/browser commands](../integration/PHASE4A.md).

## Local Student-domain configuration and roster validation

For the host development API, set this in the **ignored** `backend/.env`, then restart the API (`make backend`, or restart `make dev`):

```dotenv
STUDENT_EMAIL_DOMAINS=sksu.edu.ph
```

This is the owner's identified local development domain, not a declaration of production institutional approval. Root `.env` controls Compose; it does not configure the host API. Multiple explicitly approved domains can be comma-separated, without `@`, wildcards or URLs. Matching is exact and case-insensitive; a subdomain requires its own entry. An empty allowlist safely blocks Student creation. Faculty uses a valid unique email without a Student ID or institutional-domain restriction. The frontend obtains this policy from the protected API; no domain is hardcoded in its forms. Restart after configuration changes; an already-running API retains its startup configuration. Do not enable live Brevo delivery for isolated testing.

Download the Student `.xlsx` template. It contains one worksheet and exactly `studentId`, `name`, `email`, `course`, `contactNumber`, with text-formatted Student ID/contact input cells for up to 500 rows. A cell already stored as numeric must be re-entered as text; changing Excel's display format alone cannot recover lost leading zeroes. Structure/header/formula failures return safe bounded locations where available. Numeric IDs, missing fields, domains, duplicate rows and existing identities are validation-preview issues; invalid rows cannot be confirmed. Faculty/Staff/Admin identities are never converted by Student import/deactivation. Preserve the rejected workbook locally if further diagnosis is needed; a screenshot cannot establish its internal structure.
