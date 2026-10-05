# eLabTrack V2

eLabTrack V2 is a ground-up modernization of the existing FSMO laboratory equipment borrowing system. **eLabTrack V2 is NOT the campus-wide system.** Phase 0 is complete; Phase 1A adds configuration and runtime safeguards to the Go Fiber/React/PostgreSQL foundation. Product workflows are planned for later approved phases.

Start with the [project charter](docs/project/PROJECT_CHARTER.md), [source policy](docs/project/SOURCE_OF_TRUTH.md), [decision register](docs/project/DECISIONS.md), [open decisions](docs/project/OPEN_DECISIONS.md), and [roadmap](docs/project/ROADMAP.md). The [foundation audit](docs/project/FOUNDATION_AUDIT.md) distinguishes inherited code from approved product requirements. See the [Phase 1 security backlog](docs/project/PHASE1_SECURITY_BACKLOG.md), [Phase 0 report](docs/project/PHASE0_REPORT.md), and [Phase 1A foundation report](docs/project/PHASE1_FOUNDATION.md) before treating this foundation as production ready.

## Architecture and requirements

React SPA → REST `/api/v1` → Go Fiber v3 → PostgreSQL.

- Go **1.27.1 or newer**, as declared in `backend/go.mod`. The backend selects Go 1.27.1 with `GOTOOLCHAIN=auto`; run Go commands from `backend/` and preserve the directive.
- Node satisfying `^22.12.0 || ^24.0.0 || >=26.0.0`; npm with `npm ci` and the committed lockfile. Current Node 24.19.0 satisfies this range.
- PostgreSQL; the development Compose image is `postgres:18.6-alpine`. Docker Engine and Compose are optional for host development.
- Linux shell; Make is optional. See [stack evidence](docs/STACK.md) and [architecture rules](docs/ARCHITECTURE.md).

## Local development

```bash
cp .env.example .env
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env
```

Set matching local database credentials in `backend/.env`. Its published signing key is development-only and production rejects it. Root `.env` configures Compose; backend `.env` configures host runs. For a host PostgreSQL server, create a dedicated local `elabtrack_v2` database using your PostgreSQL administration tools. Never point this environment at V1.

For an available Docker installation, start the development database:

```bash
docker compose up -d postgres
```

Compose and the examples use the published local credentials `postgres`/`postgres`. Use matching custom root/backend credentials if desired. Existing PostgreSQL volumes retain their original password; changing an environment value does not rotate it. Use the existing credentials or change them explicitly through your local database administrator; do not delete a volume to solve a password mismatch. Renaming Compose resources does not migrate old volumes; preserve and handle any existing data separately.

Install dependencies with a compatible toolchain and access to their registries:

```bash
cd backend
go mod download
cd ../frontend
npm ci
```

Only against your dedicated development database, apply the **generic auth** migrations:

```bash
make migrate-up
```

There are no equipment/borrowing/fine/report/notification domain tables. Optional `make seed` loads synthetic starter accounts with published development passwords; never use it on real data or production.

Start each service in its own terminal:

```bash
make backend
make frontend
```

The frontend at `http://localhost:5173/` displays a minimal eLabTrack V2 placeholder, a theme toggle, and a link to `/status`. The status check explicitly reads `GET http://localhost:8080/api/v1/health` when clicked. It does not authenticate or mutate data. The API must be running and connected to PostgreSQL for a real success result.

The full development profile requires an explicit migration first:

```bash
make compose-migrate-up
# Equivalent: docker compose --profile full run --rm --build backend --migrate-up
docker compose --profile full up --build
```

The frontend uses `/api/v1` through nginx. Normal host/container API startup performs no migrations or seeding. `make migrate-status` reads migration state without creating bookkeeping tables. `make seed` requires development and an explicit request; production and test reject it. Production must inject validated settings and run migrations as a distinct, authorized release step. See the [environment inventory and requirements](docs/project/PHASE1_FOUNDATION.md). Image builds, service startup and actual PostgreSQL migration/TLS checks were not run during Phase 1A.

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

Report blocked and unrun checks accurately. See [Phase 0 report](docs/project/PHASE0_REPORT.md) for this run. No commit, push, deployment, or V1 Firebase access is authorized by this setup.

## Reference hygiene

The capstone is local at `.project-reference/ELABTRACK-SEMIFINAL.docx`, ignored by Git and excluded from any engineering deliverable. Source-controlled engineering evidence is in [the V1 audit](docs/reference/v1-audit/README.md), including an undeployed emergency hardening supplement. Product intent and implementation evidence may disagree; record conflicts rather than silently choosing a policy. Personal biographies and sample borrower/contact data must not be copied into project documentation.
# elabtrack-v2
