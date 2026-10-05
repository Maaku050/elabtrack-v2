# eLabTrack V2 project instructions

Read `docs/project/PROJECT_CHARTER.md`, `SOURCE_OF_TRUTH.md`, `DECISIONS.md`, `OPEN_DECISIONS.md`, `ROADMAP.md`, `docs/STACK.md`, current manifests, and relevant source before substantial work. `docs/project/FOUNDATION_AUDIT.md` records actual inherited behavior. The roadmap is sequencing guidance, not authorization to execute another phase.

## Scope and evidence

- eLabTrack V2 modernizes the existing FSMO operation. It is NOT the campus-wide system.
- Follow current explicit stakeholder direction first. Use the local capstone for product intent, `docs/reference/v1-audit/` for audited V1 behavior, and current source for implementation evidence.
- Label factual V1 behavior, intended legacy requirements, new V2 proposals, and confirmed V2 decisions separately. Record conflicts and policy questions with evidence; never resolve them by guessing.
- Do not introduce tenancy, campus hierarchy, cross-department borrowing, physical tracking, native mobile apps, or unapproved product features.
- Keep `.project-reference/` ignored. Never copy capstone biographies, contact details, or sample users into engineering documents.

## Architecture

- React + TypeScript + Vite + Tailwind CSS v4; preserve shadcn/ui primitives in `frontend/src/components/ui/` and `components.json`. Lucide is the icon library.
- Keep frontend boundaries feature based: page → feature hook → TanStack Query → feature API → centralized transport. Query owns server state, Zustand owns appropriate global client/UI/session state, React owns local state. React Hook Form + Zod support forms.
- Go + Fiber v3 + PostgreSQL via pgx/pgxpool. Domain owns business concepts and ports, application coordinates use cases, infrastructure implements ports, HTTP handlers adapt requests, bootstrap wires dependencies. Domain and application must not depend on Fiber or pgx. Standard library and domain/application ports are allowed inward dependencies.
- REST uses `/api/v1` and a consistent response envelope. The retained health endpoint currently has an envelope exception; follow the Phase 1 backlog before declaring this fully compliant.
- Backend validation, authorization, stock arithmetic, lifecycle transitions, and integrity rules are authoritative. Client route guards are UX only.
- Parameterize SQL. Use explicit constraints, paired SQL migrations, and transactions for integrity-sensitive changes. Preserve consequential history and audit evidence. Define concurrency/idempotency when planning mutating workflows.
- Prefer focused changes and existing dependencies. Keep a single SPA/API/database; no Redis, queues, microservices, Kubernetes, GraphQL, CQRS, or event sourcing without later explicit justification.

## Quality and safety

- Apply least privilege, safe error responses, secret hygiene, bounded queries, logging without credentials, and explicit environment configuration. Inherited auth code is provisional, not production-ready policy.
- Preserve accessibility: semantic HTML, keyboard access, accessible names, focus states, responsive behavior, and reduced motion.
- Test business invariants and critical boundaries when implemented. Never weaken tests or type checking to pass a gate.
- Run backend `go fmt ./...`, `go vet ./...`, `go test ./...`; frontend `npm run lint`, `npm run test:run`, `npm run build`; and `git diff --check`. Report actual results, toolchain mismatches, and unrun database/browser/Docker/deployment checks.
- Inspect status and remotes before Git operations; redact URL credentials. Never guess a module/repository path. Do not commit, push, change remotes, or access V1 Firebase unless explicitly authorized.
- Future approved feature plans should include scope, evidence, decisions, contracts, data integrity, acceptance criteria, tests, and handoff. No private notes service or optional agent plugin is required.
