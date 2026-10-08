# eLabTrack V2 stack evidence

Inspected 2026-10-05 (Asia/Shanghai). This is a source inventory, not a claim that registry versions or prior template validation were independently verified today. Manifests and lockfiles govern exact versions.

| Concern | Current source evidence | Accepted direction |
|---|---|---|
| Backend | `backend/go.mod`: Go 1.27.1; Fiber v3.5.0; pgx/v5 v5.11.0 | Go + Fiber v3 + pgx/pgxpool |
| Security/utilities | jwt/v5 v5.3.1, validator/v10 v10.30.5, uuid v1.6.0, zap v1.28.0, x/crypto v0.57.0 | Retain useful adapters; harden before product authentication |
| Frontend | React/react-dom ^19.3.0, Vite ^8.3.2, TypeScript ~6.0.3 | React + TypeScript + Vite |
| UI | Tailwind ^4.3.3, Base UI 1.8.0, shadcn `base-nova`, Lucide ^1.51.0 | Tailwind v4 + official shadcn primitive foundation + Lucide |
| State/forms | TanStack Query ^5.104.1, Zustand ^5.0.15, React Hook Form ^7.89.0, Zod ^4.6.5 | Query: server state; Zustand: appropriate global client state; React: local state |
| Routing/transport | React Router ^7.18.4, Axios ^1.20.0 | SPA routing, centralized transport; Phase 1F also normalizes health/error contracts |
| Testing | Vitest ^5.0.3, Testing Library, jsdom; Go testing | Retain tests; expand around approved invariants later |
| Lint/build | oxlint ^1.86.0; `tsc -b && vite build` | Keep lint and strict type checking |
| Database | SQL migrations, pgxpool, `postgres:18.6-alpine` in Compose | PostgreSQL with explicit constraints and transactions |
| Containers | Go 1.27 Alpine build; Node 22 Alpine build; nginx 1.27 Alpine frontend; Alpine 3.23 API runtime | Docker/Compose development; deployment choices await Phase 14 |

Closure environment evidence, 2026-10-05: the backend selects Go 1.27.1 linux/amd64 with `GOTOOLCHAIN=auto`. The global launcher outside the Go module remains Go 1.26.5; backend commands use the compatible automatically selected toolchain. Node 24.19.0 and npm 11.17.0 satisfy the frontend engine. Docker Engine 29.8.0 and Compose v5.5.1 are accessible, and `docker compose config` passes. No Docker services were started. Earlier blocked environment evidence is preserved in `project/PHASE0_REPORT.md`.

The authoritative Go module is now `github.com/Maaku050/elabtrack-v2/backend`. The closure request supplied the repository identity, and read-only Git inspection verified `origin` as `https://github.com/Maaku050/elabtrack-v2.git` on `main`. Module declaration and internal imports were renamed consistently; dependency versions and `go.sum` were preserved. The original identity blocker remains recorded as historical evidence in `project/PHASE0_REPORT.md`.

Historical Windows/Laragon validation, private template-owner preferences, and unrelated project specs are not current eLabTrack build evidence. No dependencies were upgraded or required versions lowered during Phase 0.

Batch1 addition,2026-10-09: `backend/go.mod` adds Excelize v2.11.0 for the bounded Student-only XLSX parser/template; maintained library use is wrapped with ZIP/XML size/formula/external-file preflight. Brevo uses a backend standard-library HTTP adapter, with no frontend SDK or credentials. Inventory canonical PNG validation/storage uses Go image/png/image/jpeg and PostgreSQL bytea, with strict encoded/pixel limits; no paid storage dependency. Existing frontend stack and shadcn primitives are retained. Batch1 selected Go1.27.1 and local Node24.19/npm11.17 for newly run verification; prior closure/container evidence remains historical.
