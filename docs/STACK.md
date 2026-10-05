# eLabTrack V2 stack evidence

Inspected 2026-10-05 (Asia/Shanghai). This is a source inventory, not a claim that registry versions or prior template validation were independently verified today. Manifests and lockfiles govern exact versions.

| Concern | Current source evidence | Accepted direction |
|---|---|---|
| Backend | `backend/go.mod`: Go 1.27.1; Fiber v3.5.0; pgx/v5 v5.11.0 | Go + Fiber v3 + pgx/pgxpool |
| Security/utilities | jwt/v5 v5.3.1, validator/v10 v10.30.5, uuid v1.6.0, zap v1.28.0, x/crypto v0.57.0 | Retain useful adapters; harden before product authentication |
| Frontend | React/react-dom ^19.3.0, Vite ^8.3.2, TypeScript ~6.0.3 | React + TypeScript + Vite |
| UI | Tailwind ^4.3.3, Base UI 1.8.0, shadcn `base-nova`, Lucide ^1.51.0 | Tailwind v4 + official shadcn primitive foundation + Lucide |
| State/forms | TanStack Query ^5.104.1, Zustand ^5.0.15, React Hook Form ^7.89.0, Zod ^4.6.5 | Query: server state; Zustand: appropriate global client state; React: local state |
| Routing/transport | React Router ^7.18.4, Axios ^1.20.0 | SPA routing, centralized transport; health temporarily adapts its legacy response |
| Testing | Vitest ^5.0.3, Testing Library, jsdom; Go testing | Retain tests; expand around approved invariants later |
| Lint/build | oxlint ^1.86.0; `tsc -b && vite build` | Keep lint and strict type checking |
| Database | SQL migrations, pgxpool, `postgres:18.6-alpine` in Compose | PostgreSQL with explicit constraints and transactions |
| Containers | Go 1.27 Alpine build; Node 22 Alpine build; nginx 1.27 Alpine frontend; Alpine 3.23 API runtime | Docker/Compose development; deployment choices await Phase 14 |

Local executable evidence: Go 1.26.5, Node 24.19.0, npm 11.17.0. Go is below the required directive; Node satisfies the package engine. Docker resolves to a Windows/WSL integration shim and cannot execute Compose in this session. Dependency installation and checks are recorded in `project/PHASE0_REPORT.md`.

The Go module remains `github.com/fullstacktemplate/backend` because `.git` metadata is inaccessible here and no authoritative new remote is available. Keep every import consistent until the repository owner supplies/verifies the new path. Do not derive it from a folder name or invent a GitHub URL.

Historical Windows/Laragon validation, private template-owner preferences, and unrelated project specs are not current eLabTrack build evidence. No dependencies were upgraded or required versions lowered during Phase 0.
