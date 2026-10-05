# eLabTrack V2 engineering boundaries

The accepted architecture is one React SPA, one Go Fiber REST API, and one PostgreSQL database. Current runtime code is generic auth infrastructure plus the Phase 0 frontend placeholder; product workflows are not implemented.

## Backend

```text
HTTP interface → Application use cases → Domain entities and ports
                        ↑                        ↑
                 Infrastructure implements ports
Bootstrap composes adapters, services, handlers, and routes.
```

Domain and application may use the Go standard library and domain/application ports, but must not import Fiber, pgx, or infrastructure. HTTP adapts requests; business policies belong to application/domain code. Parameterized repositories and one pgxpool are retained. Integrity-sensitive use cases must coordinate transactions through inward-facing ports rather than importing database adapters.

Current paths: `backend/internal/domain/`, `application/`, `infrastructure/`, `interface/http/`, `bootstrap/`, `config/`, and `shared/`. UUID is used as a shared identifier dependency. The health handler currently directly depends on the infrastructure health checker; this is a documented boundary adaptation for Phase 1. Generic role/session behavior is provisional.

REST is versioned `/api/v1`. Normal handlers use `{success,message,data,meta,error?}`. Health currently returns `{status,services}` and 204 responses have no body. Phase 1 must resolve the health exception and align clients/contracts. The placeholder explicitly adapts that endpoint without adding authentication.

## Frontend

Page → feature hook → TanStack Query → feature API → shared transport → API. Each real feature owns its API, hooks, components, schemas, types, and pages as needed. Avoid empty scaffolding for unapproved features.

- TanStack Query owns fetched server data.
- Zustand owns global client/UI state and suitable session metadata; the inherited token persistence requires Phase 1 review.
- React owns local UI state.
- React Hook Form and Zod support input UX; the backend remains authoritative.
- shadcn/ui primitives remain in `frontend/src/components/ui/`, with Base UI behavior, Tailwind v4 tokens, and Lucide icons. Product design tokens and shell come in Phase 3.

Phase 0 routes are `/`, `/status`, and a wildcard 404. Generic dashboard/showcase/auth/profile screens were removed. Unmounted auth/user API, schema, hook, and store modules remain as provisional starter adapters; their presence does not accept public registration or any V2 role policy.

## Cross-cutting constraints

Server-enforced permissions and bounded queries replace V1 client-owned business rules. Preserve material history and use database constraints/transactions for stock, returns, accountability, and later audit requirements. Define lifecycle vocabulary, retries and idempotency before mutating product workflows. Do not optimistically confirm stock deductions, permissions, fines, or irreversible actions.

No campus tenancy/hierarchy, distributed infrastructure, physical tracking, offline platform, or native mobile app is in current scope. See `project/DECISIONS.md`, `OPEN_DECISIONS.md`, and `ROADMAP.md` before changing a boundary or executing another phase.
