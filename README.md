# eLabTrack V2

eLabTrack V2 is a ground-up modernization of the existing FSMO laboratory equipment borrowing system. **eLabTrack V2 is NOT the campus-wide system.** The current repository is Phase 0: rebaseline and adaptation of a Go Fiber/React/PostgreSQL starter. Product workflows are planned for later approved phases.

Start with the [project charter](docs/project/PROJECT_CHARTER.md), [source policy](docs/project/SOURCE_OF_TRUTH.md), [decision register](docs/project/DECISIONS.md), [open decisions](docs/project/OPEN_DECISIONS.md), and [roadmap](docs/project/ROADMAP.md). The [foundation audit](docs/project/FOUNDATION_AUDIT.md) distinguishes inherited code from approved product requirements. See the [Phase 1 security backlog](docs/project/PHASE1_SECURITY_BACKLOG.md) and [Phase 0 report](docs/project/PHASE0_REPORT.md) before treating this foundation as production ready.

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

Set database credentials and replace the development JWT secret in `backend/.env`. Root `.env` configures Compose; backend `.env` configures host runs. For a host PostgreSQL server, create a dedicated local `elabtrack_v2` database using your PostgreSQL administration tools. Never point this environment at V1.

For an available Docker installation, start the development database:

```bash
docker compose up -d postgres
```

Compose defaults to development credentials `postgres`/`postgres`. Set `DB_PASSWORD=postgres` in `backend/.env` when connecting from the host, or use matching custom root/backend credentials. Renaming Compose resources does not migrate old volumes; preserve and handle any existing data separately.

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

The optional full development profile is `docker compose --profile full up --build`; its frontend uses `/api/v1` through nginx by default. Backend currently applies auth migrations on container startup. That inherited behavior is a Phase 1 hardening item, not a production deployment procedure. Docker execution is unverified in this workspace.

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
