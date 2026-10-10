# eLabTrack V2

**Complete shared UI reconstruction, 2026-10-10:** Staff/Admin desktop and Student/Faculty mobile engineering verification is complete; **owner visual acceptance remains pending**. See the [implementation report and 14-refinement matrix](docs/project/COMPLETE_FRONTEND_UI_RECONSTRUCTION_REPORT.md) and [before/after evidence gallery](docs/ux/verification/complete-ui/index.html). Final gates: 328 frontend tests in serial and parallel, 121 named PostgreSQL/API checks, 83 Chromium check groups, lint/ESLint/build/Go checks and data preservation. The isolated presentation remains at http://127.0.0.1:15177/login; normal data/schema are retained and require separately authorized compatible integration. No commit, push or deployment.

eLabTrack V2 modernizes FSMO laboratory equipment management and borrowing. It is not a campus-wide system. Phase7 is owner accepted at `c2741c5`; the authorized Phases8–14 implementation adds returns, replacements, loan-level fines, persistent notifications, operational dashboards/reports/CSV, DEC-080 equipment selection/cart, protected profile images and an isolated offline presentation. Phases8–14 engineering gates passed; final owner acceptance remains pending. See the [full implementation report](docs/project/SYSTEM_IMPLEMENTATION_REPORT.md). Current status and executed evidence are in the [progress register](docs/project/PHASE8_14_PROGRESS.md) and phase reports.

For the full-system presentation, follow [setup and secure credential retrieval](docs/project/PRESENTATION_SETUP.md) and the [FSMO walkthrough](docs/project/PRESENTATION_WALKTHROUGH.md). Use **http://127.0.0.1:15177/login**, API18087 and PostgreSQL54836. This contains fictional data and demonstration-only terms. Normal development5434/8080/5173 and the accepted Phase7 demo54835/18086/15176 are preserved separately. The normal database remains migration000008; do not run the new API against it before a separately reviewed and authorized000009–000011 rollout.

Only Admin provisions accounts. Students and Faculty are BORROWER categories; Student textual IDs/institutional domains and Student-only Excel workflows remain enforced. Faculty is individual-only with a valid unique email and no required Student ID. Live borrowing still requires verified activation and current officially published FSMO terms with documented acceptance. Live Brevo delivery and institutional approvals are deferred, not bypassed. See [account policy](docs/project/ACCOUNT_PROVISIONING_POLICY.md) and [production readiness gaps](docs/project/PRODUCTION_READINESS_GAPS.md).
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

The normal development database was backed up, migrated through **000008_borrowing** and owner accepted in the Phase7 checkpoint. It contains existing records and must remain unchanged during this presentation work. All historical000001–000008 files remain immutable. The full-system source requires the paired000009 accountability,000010 notification and000011 profile-image migrations; these were applied and verified only on isolated targets. Applying them to normal data requires a separately authorized backup, restore rehearsal, checksum/status review, migration and named runtime grants. `make migrate-status` is read-only; generic migration commands below are tools, not authorization to change normal data.

Account, equipment, borrowing, return/replacement/fine, notification and reporting workflows are implemented. There is no published default product login. An initial real Admin requires the approved secure bootstrap procedure; demonstration credentials are generated privately by the isolated presentation tool. `make seed` does not provide product accounts or passwords.
Migration 000003 explicitly invalidates existing refresh sessions when applied or rolled back; users must log in again. Coordinate schema and application versions. Historical migrations are preserved. `make sessions-cleanup` explicitly deletes at most 1000 expired/revoked records terminal for over seven days; operators must run enough batches regularly. It never runs at API startup. Phase 1G verified migrations, refresh concurrency/rollback and cleanup against disposable PostgreSQL. See the [integration run instructions](integration/README.md) and [Phase 1 foundation report](docs/project/PHASE1_FOUNDATION.md).

Browser access tokens live only in memory; refresh credentials use the host-only `elabtrack_v2_refresh` HttpOnly cookie at `/api/v1/auth`, with explicit SameSite=Lax and production Secure. Cookie-changing auth POSTs require an exact configured frontend Origin, including login. Configure explicit local origins (the example includes `http://localhost:5173`); local HTTP cookies derive Secure=false only from development/test APP_ENV. Production requires HTTPS and a same-site SPA/API arrangement, normally the same-origin nginx `/api/v1` proxy. Arbitrary cross-site deployments are unsupported. Login/refresh JSON carries access credentials and safe current account fields only; raw refresh JSON input has been removed.

The HTTP perimeter uses the socket peer by default. `TRUSTED_PROXIES` accepts deliberately chosen literal proxy IPs/CIDRs; forwarded client IP is honored only from those peers, using the first untrusted hop from the right. Proxy operators must sanitize forwarding headers and restrict bypass; no hosting provider or production network is selected. Limits are per process and effective client IP: login 10, refresh 60, general API/logout 120 per one-minute window, with safe 429/Retry-After. Shared/NAT clients share budgets; configure operational values deliberately before deployment. CORS permits only explicit canonical origins and current GET/POST/PATCH/OPTIONS methods. Global body ceiling defaults to 1 MiB, auth POSTs to 16 KiB; read/write/idle defaults are 10/15/60 seconds.

The nginx container serves same-origin production assets/API with a baseline document CSP and security headers; its Docker build defaults to `/api/v1`. External API build overrides require a reviewed matching connect-src policy. Host Vite development retains the explicit separate localhost API origin. The supplied nginx template is HTTP development only and does not assert HSTS. Production's authoritative HTTPS edge must own TLS/HSTS; API HSTS is emitted only in production on actual TLS or exact `https` from an explicitly trusted peer. Phase 1G verifies the local nginx/browser/proxy/rate/logging behavior; production HTTPS and hosting topology remain deployment checks. See the [foundation report](docs/project/PHASE1_FOUNDATION.md).

After a separately authorized schema/application upgrade and existing explicit migration/grant setup are ready, start both host services in one terminal (Linux/WSL, GNU Make, Bash 5.1+ and util-linux `setsid`):

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

Default frontend: **http://localhost:5173/**. Default backend: **http://localhost:8080/**; liveness **http://localhost:8080/api/v1/health**, readiness **http://localhost:8080/api/v1/ready**. The root frontend route redirects to `/login` or the current role workspace. `/status` remains the explicit public connection check.

Phase 3B development-only visual previews:

| Preview | Exact URL |
|---|---|
| Borrower Home | http://localhost:5173/__preview/borrower/home |
| Equipment Catalog | http://localhost:5173/__preview/borrower/equipment |
| Staff Dashboard | http://localhost:5173/__preview/staff/dashboard |
| Pending Requests | http://localhost:5173/__preview/staff/requests |

These use synthetic local fixtures; they do not submit, approve, issue, return or clear anything. They remain absent from production routing. These are historical previews; current authenticated business screens use the role routes and real API.

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

For a live browser/API/session regression, provide headless Chromium and follow the isolated database/fixture instructions.

Current terms/authentication verification uses the [Phase 4B isolated reproduction guide](integration/PHASE4B.md) and `frontend/scripts/phase4b-browser-qa.mjs`. It generates random-password accounts and explicitly synthetic TEST documents only in a separate disposable database. Immutable acceptance/publication history is retained until that whole owned test database is discarded. The [Phase 4A guide](integration/PHASE4A.md) is historical authentication evidence; its exact four-migration replay refuses newer migration history and must use its historical isolated schema. The previous `local-environment-smoke.mjs` and its [closeout evidence](docs/ux/verification/local-environment/README.md) are historical: their public-registration fixture step is unavailable after Phase 4A. Do not re-enable registration to run an old harness.

## Reference hygiene

The capstone is local at `.project-reference/ELABTRACK-SEMIFINAL.docx`, ignored by Git and excluded from any engineering deliverable. Source-controlled engineering evidence is in [the V1 audit](docs/reference/v1-audit/README.md), including an undeployed emergency hardening supplement. Product intent and implementation evidence may disagree; record conflicts rather than silently choosing a policy. Personal biographies and sample borrower/contact data must not be copied into project documentation.

Batch 1 account-management implementation adds Admin-only Student/Faculty provisioning, Student Excel preview/confirmation, secure activation infrastructure and fixed-STAFF administration. See [the implementation plan](docs/project/BATCH1_IMPLEMENTATION_PLAN.md) and [account API contracts](docs/API_CONTRACTS.md#phase-5-implemented-account-management-contracts). Approved Student domains and backend Brevo settings are deliberately blank in example configuration; missing domains block Student onboarding, and no live email delivery or official terms approval is inferred. Disposable tests use [the isolated Batch 1 setup](integration/compose.batch1.yml); normal `.env` files are not overwritten.

Batch1 Phase6 implements the protected Borrower catalog/details and Staff/Admin catalog/category/inventory/history/available-adjustment UI, with Admin reviewed available-count reconciliation and bounded authenticated PNG catalog images. New paired000007 enforces physical conservation and immutable evidence. Use the existing explicit local migration/runtime-grant workflow for another database. Normal local storage was migrated and verified separately at the pre-Phase7 checkpoint; the original isolated Batch1 evidence remains historical. [Phase5 report](docs/project/PHASE5_REPORT.md), [Phase6 report](docs/project/PHASE6_REPORT.md) and [isolated verification](integration/BATCH1.md) distinguish core verification from live Brevo/approved-domain/official-term dependencies. Phase7 is now owner accepted; Phases8–14 use the separate current DEC-082 authority and presentation environment. Historical Batch1 reports retain their original scope.
