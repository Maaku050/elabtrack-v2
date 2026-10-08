# eLabTrack V2

eLabTrack V2 is a ground-up modernization of the existing FSMO laboratory equipment borrowing system. **eLabTrack V2 is NOT the campus-wide system.** Phase 0 is complete; Phases 1A–1I add configuration safeguards, current-account authorization, strict access JWTs, hash-only transactional refresh sessions, browser sessions, HTTP perimeter controls and API contracts/structured observability, real local verification and atomic/checksummed/locked migrations with database role separation and cross-tab session coordination to the Go Fiber/React/PostgreSQL foundation. Phase 3B is complete and owner visually approved; the four implemented previews and both themes are the accepted visual baseline. Phase 4 has not started. See the [local readiness report](docs/project/LOCAL_ENVIRONMENT_READINESS_REPORT.md) for the verified full-stack setup and checkpoint plan.

Start with the [project charter](docs/project/PROJECT_CHARTER.md), [source policy](docs/project/SOURCE_OF_TRUTH.md), [decision register](docs/project/DECISIONS.md), [open decisions](docs/project/OPEN_DECISIONS.md), and [roadmap](docs/project/ROADMAP.md). The [foundation audit](docs/project/FOUNDATION_AUDIT.md) distinguishes inherited code from approved product requirements. See the [Phase 1 security backlog](docs/project/PHASE1_SECURITY_BACKLOG.md), [Phase 0 report](docs/project/PHASE0_REPORT.md), and [Phase 1 foundation report](docs/project/PHASE1_FOUNDATION.md) before treating this foundation as production ready.

## Architecture and requirements

React SPA → REST `/api/v1` → Go Fiber v3 → PostgreSQL.

- Go **1.27.1 or newer**, as declared in `backend/go.mod`. The backend selects Go 1.27.1 with `GOTOOLCHAIN=auto`; run Go commands from `backend/` and preserve the directive.
- Node satisfying `^22.12.0 || ^24.0.0 || >=26.0.0`; npm with `npm ci` and the committed lockfile. Current Node 24.19.0 satisfies this range.
- PostgreSQL; the development Compose image is `postgres:18.6-alpine`. Docker Engine and Compose are optional for host development.
- Linux shell; Make is optional. See [stack evidence](docs/STACK.md) and [architecture rules](docs/ARCHITECTURE.md).

## Local development

Start Docker Desktop on Windows with its Linux engine. In **Settings → Resources → WSL Integration**, enable this distribution if needed. From this WSL repository, verify `docker version`, `docker info`, and `docker compose version`. If the engine cannot be reached, start Desktop/enable integration and retry; do not install a competing daemon.

The closeout restored all three ignored local files with independent generated secrets and mode `0600`. Preserve them. For a new checkout, create only missing files:

```bash
test -e .env || cp .env.example .env
test -e backend/.env || cp backend/.env.example backend/.env
test -e frontend/.env || cp frontend/.env.example frontend/.env
chmod 600 .env backend/.env frontend/.env
```

Replace example passwords and JWT signing key with independent development-only random secrets (at least 32 random bytes each), stored only in these ignored files. Root `.env` configures Compose; backend `.env` configures host runs. Match root `DB_PASSWORD` with backend runtime `DB_PASSWORD`; match root `MIGRATION_DB_PASSWORD` with the password in backend `MIGRATION_DATABASE_URL`. Keep the bootstrap password separate. Never put connection credentials into documentation or shell command arguments.

Current verified non-secret settings:

| Configuration | Local setting |
|---|---|
| Root Compose database | `DB_NAME=elabtrack_v2`, `DB_PORT=5434`, `BOOTSTRAP_DB_USER=postgres`, `DB_USER=elabtrack_runtime` |
| Backend environment/listener | `APP_ENV=development`, `APP_PORT=8080` |
| Backend runtime connection | `DB_HOST=127.0.0.1`, `DB_PORT=5434`, `DB_NAME=elabtrack_v2`, `DB_USER=elabtrack_runtime`, `DB_SSLMODE=disable` |
| Migration URL components | `elabtrack_migrator`, host `127.0.0.1`, port `5434`, database `elabtrack_v2`, `sslmode=disable`; password stays private |
| Browser Origin | `FRONTEND_URL=http://localhost:5173`, `ALLOWED_ORIGINS=http://localhost:5173,http://localhost:4173`, `TRUSTED_PROXIES` empty |
| Frontend host configuration | `VITE_API_URL=http://localhost:8080/api/v1` |

Port **5432 already has an unrelated PostgreSQL listener** on this machine. It was preserved; the project database uses **5434**. On another machine, choose a free local port and keep root, backend, and migration URL consistent. Local HTTP and database TLS disabling are development settings; preserve the existing production validation. JWT issuer, limits and timeouts retain the environment-example defaults.

For a host PostgreSQL server, have its administrator create the dedicated local `elabtrack_v2` database and the same schema-owner/runtime role boundary shown in [the bootstrap script](backend/database/bootstrap-roles.sh); apply named runtime grants after migration. Never point this environment at V1.

For an available Docker installation, start the development database:

```bash
docker compose up -d postgres
```

Fresh Compose databases bootstrap separate `elabtrack_migrator` schema-owner and `elabtrack_runtime` DML logins. The examples contain obvious local-only passwords; match custom runtime/migration values between root and backend configuration. Bootstrap administrator credentials stay in the PostgreSQL container. Existing PostgreSQL volumes retain their original password; changing an environment value does not rotate it. Use the existing credentials or change them explicitly through your local database administrator; do not delete a volume to solve a password mismatch. Renaming Compose resources does not migrate old volumes; preserve and handle any existing data separately.

Install dependencies with a compatible toolchain and access to their registries:

```bash
cd backend
go mod download
cd ../frontend
npm ci
cd ..
```

Inspect the dedicated development database before applying anything:

```bash
make migrate-status
```

The restored database already has all three **generic auth** migrations with verified checksums and runtime grants. No migration command is needed to start it. For a fresh database or an explicitly reviewed pending suffix, use the existing runner:

```bash
make migrate-up
make compose-runtime-grants
make migrate-status
```

Preserve initialized history; investigate checksum/history errors instead of rewriting SQL, adopting history automatically or resetting the database. The runner retains atomic bookkeeping, paired checksums and advisory exclusion. No seed is required for startup or anonymous previews.

There are no equipment/borrowing/fine/report/notification domain tables. Optional `make seed` loads synthetic starter accounts with published development passwords; never use it on real data or production.

Migration 000003 explicitly invalidates existing refresh sessions when applied or rolled back; users must log in again. Coordinate schema and application versions. Historical migrations are preserved. `make sessions-cleanup` explicitly deletes at most 1000 expired/revoked records terminal for over seven days; operators must run enough batches regularly. It never runs at API startup. Phase 1G verified migrations, refresh concurrency/rollback and cleanup against disposable PostgreSQL. See the [integration run instructions](integration/README.md) and [Phase 1 foundation report](docs/project/PHASE1_FOUNDATION.md).

Browser access tokens live only in memory; refresh credentials use the host-only `elabtrack_v2_refresh` HttpOnly cookie at `/api/v1/auth`, with explicit SameSite=Lax and production Secure. Cookie-changing auth POSTs require an exact configured frontend Origin, including login and local registration. Configure explicit local origins (the example includes `http://localhost:5173`); local HTTP cookies derive Secure=false only from development/test APP_ENV. Production requires HTTPS and a same-site SPA/API arrangement, normally the same-origin nginx `/api/v1` proxy. Arbitrary cross-site deployments are unsupported. Login/refresh JSON carries access credentials and safe current account fields only; raw refresh JSON input has been removed.

The HTTP perimeter uses the socket peer by default. `TRUSTED_PROXIES` accepts deliberately chosen literal proxy IPs/CIDRs; forwarded client IP is honored only from those peers, using the first untrusted hop from the right. Proxy operators must sanitize forwarding headers and restrict bypass; no hosting provider or production network is selected. Limits are per process and effective client IP: login 10, refresh 60, local registration 5, general API/logout 120 per one-minute window, with safe 429/Retry-After. Shared/NAT clients share budgets; configure operational values deliberately before deployment. CORS permits only explicit canonical origins and current GET/POST/PATCH/OPTIONS methods. Global body ceiling defaults to 1 MiB, auth POSTs to 16 KiB; read/write/idle defaults are 10/15/60 seconds.

The nginx container serves same-origin production assets/API with a baseline document CSP and security headers; its Docker build defaults to `/api/v1`. External API build overrides require a reviewed matching connect-src policy. Host Vite development retains the explicit separate localhost API origin. The supplied nginx template is HTTP development only and does not assert HSTS. Production's authoritative HTTPS edge must own TLS/HSTS; API HSTS is emitted only in production on actual TLS or exact `https` from an explicitly trusted peer. Phase 1G verifies the local nginx/browser/proxy/rate/logging behavior; production HTTPS and hosting topology remain deployment checks. See the [foundation report](docs/project/PHASE1_FOUNDATION.md).

After PostgreSQL and the existing explicit migration/grant setup above are ready, start both host services in one terminal (Linux/WSL, GNU Make, Bash 5.1+ and util-linux `setsid`):

```bash
make dev
```

The supervisor runs the existing `make backend` and `make frontend` targets concurrently, shows both logs and stops both process groups (including `go run`/npm children) on **Ctrl+C** or when either target exits. Startup errors remain visible and return a nonzero status. It allows up to twelve seconds for shutdown before killing a remaining group. It never starts, resets, seeds or migrates PostgreSQL. The database stays running when you stop the servers.

Individual targets still work independently, in separate terminals if desired:

```bash
make backend   # Go API only; requires PostgreSQL and backend configuration
make frontend  # Vite only; fails if port 5173 is occupied
make help
```

Default frontend: **http://localhost:5173/**. Default backend: **http://localhost:8080/**; liveness **http://localhost:8080/api/v1/health**, readiness **http://localhost:8080/api/v1/ready**. The root frontend route remains the foundation page, with a theme toggle and `/status` link.

Phase 3B development-only visual previews:

| Preview | Exact URL |
|---|---|
| Borrower Home | http://localhost:5173/__preview/borrower/home |
| Equipment Catalog | http://localhost:5173/__preview/borrower/equipment |
| Staff Dashboard | http://localhost:5173/__preview/staff/dashboard |
| Pending Requests | http://localhost:5173/__preview/staff/requests |

These use synthetic local fixtures; they do not submit, approve, issue, return or clear anything. They remain absent from production routing. Final business screens are not implemented.

Host Vite has **no API proxy**. `frontend/.env.example` and the transport default set `VITE_API_URL=http://localhost:8080/api/v1`, so requests go directly to the backend. Root `.env` controls Compose only; its `/api/v1` setting is for the nginx container and does not configure host Vite. Application loading makes one credentialed refresh-cookie restoration attempt through the unchanged transport/cross-tab coordinator. A missing/invalid session returns the expected 401 and becomes unauthenticated; network/server failures show a recoverable session banner and allow public content. Restoration can rotate an existing refresh session. The `/status` action reads backend health without cookie/bearer credentials. Memory-only access credentials and HttpOnly refresh cookies remain unchanged.

Connection troubleshooting:

- **Backend database startup failure:** check `backend/.env`, credentials and database availability. For this project's Compose workflow, run `docker compose up -d postgres`; inspect `docker compose ps postgres` and `docker compose logs postgres`. If Docker Desktop is stopped or unavailable in WSL, enable its Linux-engine/WSL integration, or use your explicitly configured host PostgreSQL. Do not delete existing volumes to resolve credentials. Initial migrations and grants remain separate commands above.
- **Port already occupied:** `make frontend` uses `--strictPort` so it cannot silently switch to a different, unapproved CORS origin. Stop the conflicting server; a target failure stops its peer under `make dev`. If intentionally changing ports, keep the displayed Vite URL, backend `APP_PORT`, frontend API URL, `FRONTEND_URL` and exact `ALLOWED_ORIGINS` consistent.
- **CORS/Origin errors:** open `http://localhost:5173`, as allowed by the backend example. `127.0.0.1` and other ports are different origins. For an intentional alternate origin, configure it explicitly in `FRONTEND_URL`/`ALLOWED_ORIGINS`; retain exact allowlists and the auth Origin checks. Restart the API after backend configuration changes and Vite after `VITE_*` changes.
- **Wrong API URL or connection refused:** check `frontend/.env` and the backend listener. Do not set host Vite to `/api/v1` expecting a nonexistent proxy. `curl http://localhost:8080/api/v1/health` checks process liveness; `curl http://localhost:8080/api/v1/ready` checks database connectivity. Bootstrap requires PostgreSQL; a running process may remain live during later database loss while readiness returns 503.
- **Anonymous 401 versus server failure:** an empty refresh session is expected to return 401. Connection refused, readiness 503 or a restoration error banner points to connectivity/configuration, not permission to bypass authentication. Use the same `localhost` host consistently for the SPA and API cookies.

The full development profile requires an explicit migration first:

```bash
make compose-migrate-up
# Equivalent: docker compose --profile tools run --rm --build migrator
make compose-runtime-grants
docker compose --profile full up --build
```

The frontend uses `/api/v1` through nginx. Normal host/container API startup performs no migrations or seeding. `make migrate-status` reads migration state without creating bookkeeping tables. `make seed` requires development and an explicit request; production and test reject it. Production must inject validated settings and run migrations as a distinct, authorized release step. See the [environment inventory and requirements](docs/project/PHASE1_FOUNDATION.md). Image builds, service startup and actual PostgreSQL migration/TLS checks were not run during Phase 1A.

API JSON successes retain the existing envelope; errors contain code/message/requestId and optional field errors. Server-owned X-Request-ID and safe structured Zap events support correlation. See [API contracts and logging](docs/API_CONTRACTS.md) and the [Phase 1G report](docs/project/PHASE1_FOUNDATION.md#phase-1g--real-integration-verification). Phase 1H verifies atomic up/down bookkeeping, paired SHA-256 checksums, PostgreSQL advisory exclusion and API/auth under a DML-only database role. Applied files are immutable. CLI migration actions use MIGRATION_DATABASE_URL; production refuses missing credentials or the runtime user. API startup does not need that secret. Known checksum-free local history requires the explicit development-only adoption procedure, never automatic backfill. See the [Phase 1H report](docs/project/PHASE1_FOUNDATION.md#phase-1h--migration-runner--database-privilege-hardening). Phase 1I cross-tab session coordination is complete; browser access credentials remain memory-only and credential-free lifecycle messages coordinate cookie mutations and peer logout. See the [Phase 1I report](docs/project/PHASE1_FOUNDATION.md#phase-1i--cross-tab-session-coordination).

## Validation

```bash
cd backend
go fmt ./...
go vet ./...
go test ./...
cd ../frontend
npm run lint
npm run test:run
npm run build
cd ..
git diff --check
```

For a live browser/API/session regression, keep `make dev` running, provide an installed headless Chromium executable and local `psql`, then run from the root:

```bash
CHROMIUM_EXECUTABLE=/path/to/chromium node frontend/scripts/local-environment-smoke.mjs
```

This uses real API requests, creates one temporary unprivileged synthetic fixture through development registration, expires only its own session for the expiry check, and removes only that fixture in cleanup. It refuses a non-development runtime target. Results/captures go to `/tmp/elabtrack-local-environment-smoke`; no tokens, cookies or passwords are exported. Chromium may need system libraries or an environment-specific `LD_LIBRARY_PATH`. This is separate from the original visual QA, which uses an anonymous 401 fixture. See [closeout evidence](docs/ux/verification/local-environment/README.md).

Report blocked and unrun checks accurately. See the [local readiness report](docs/project/LOCAL_ENVIRONMENT_READINESS_REPORT.md) for this run. No commit, push, deployment, or V1 Firebase access is authorized by this setup.

## Reference hygiene

The capstone is local at `.project-reference/ELABTRACK-SEMIFINAL.docx`, ignored by Git and excluded from any engineering deliverable. Source-controlled engineering evidence is in [the V1 audit](docs/reference/v1-audit/README.md), including an undeployed emergency hardening supplement. Product intent and implementation evidence may disagree; record conflicts rather than silently choosing a policy. Personal biographies and sample borrower/contact data must not be copied into project documentation.
