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

Current paths: `backend/internal/domain/`, `application/`, `infrastructure/`, `interface/http/`, `bootstrap/`, `config/`, and `shared/`. UUID is used as a shared identifier dependency. The health handler accepts a small context-based check interface; bootstrap supplies the database adapter. Generic role/session behavior is provisional.

REST is versioned `/api/v1`. All JSON handlers use `{success,message,data,meta,error?}`; errors require nested code/message/requestId and optional field details. 204 has no body. Health is enveloped process liveness; readiness checks the dependency. See [API contracts and logging](API_CONTRACTS.md).

## Frontend

Page → feature hook → TanStack Query → feature API → shared transport → API. Each real feature owns its API, hooks, components, schemas, types, and pages as needed. Avoid empty scaffolding for unapproved features.

- TanStack Query owns fetched server data.
- Zustand owns global client/UI state and suitable session metadata; access credentials are non-persisted memory and refresh credentials are HttpOnly cookies under Phase 1D.
- React owns local UI state.
- React Hook Form and Zod support input UX; the backend remains authoritative.
- shadcn/ui primitives remain in `frontend/src/components/ui/`, with Base UI behavior, Tailwind v4 tokens, and Lucide icons. Approved mockups come in Phase 3A; the shadcn-based design system and shell follow in Phase 3B.

Phase 0 routes are `/`, `/status`, and a wildcard 404. Generic dashboard/showcase/auth/profile screens were removed. Unmounted auth/user API, schema, hook, and store modules remain as provisional starter adapters; their presence does not accept public registration or any V2 role policy.

## Cross-cutting constraints

Server-enforced permissions and bounded queries replace V1 client-owned business rules. Preserve material history and use database constraints/transactions for stock, returns, accountability, and later audit requirements. Define lifecycle vocabulary, retries and idempotency before mutating product workflows. Do not optimistically confirm stock deductions, permissions, fines, or irreversible actions.

No campus tenancy/hierarchy, distributed infrastructure, physical tracking, offline platform, or native mobile app is in current scope. See `project/DECISIONS.md`, `OPEN_DECISIONS.md`, and `ROADMAP.md` before changing a boundary or executing another phase.
