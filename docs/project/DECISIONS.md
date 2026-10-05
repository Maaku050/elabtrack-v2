# eLabTrack V2 decision register

Recorded 2026-10-05. Status vocabulary: **Accepted**, **Proposed**, **Deferred**, **Needs Stakeholder Input**, **Rejected**. Acceptance records direction, not completed implementation. All accepted entries below derive from the current stakeholder Phase 0 request; workflow cleanup follows the authorized inspection and dependency check. Unresolved product policies are in `OPEN_DECISIONS.md`.

## ELAB-V2-DEC-001 — FSMO scope

- Status: Accepted
- Decision: V2 modernizes the existing FSMO operation; it is NOT the campus-wide system.
- Authority/date: current stakeholder Phase 0 direction, 2026-10-05.
- Rationale/impact: No organizational tenancy, campus hierarchy or cross-department borrowing.

## ELAB-V2-DEC-002 — Frontend platform

- Status: Accepted
- Decision: React + TypeScript + Vite.
- Authority/date: current stakeholder Phase 0 direction, 2026-10-05.
- Rationale/impact: Responsive SPA; no initial native mobile app.

## ELAB-V2-DEC-003 — Backend platform

- Status: Accepted
- Decision: Go + Fiber v3 with pgx/pgxpool.
- Authority/date: current stakeholder Phase 0 direction, 2026-10-05.
- Rationale/impact: Retain manifest versions; no framework migration.

## ELAB-V2-DEC-004 — Persistence

- Status: Accepted
- Decision: PostgreSQL with SQL migrations, explicit constraints and transactions.
- Authority/date: current stakeholder Phase 0 direction, 2026-10-05.
- Rationale/impact: No product-domain tables until approved Phase 2 design.

## ELAB-V2-DEC-005 — API

- Status: Accepted
- Decision: REST, versioned /api/v1, consistent response envelope.
- Authority/date: current stakeholder Phase 0 direction, 2026-10-05.
- Rationale/impact: Inherited health exception is an adaptation item, not an exception accepted permanently.

## ELAB-V2-DEC-006 — UI foundation

- Status: Accepted
- Decision: Official shadcn/ui primitives are retained.
- Authority/date: current stakeholder Phase 0 direction, 2026-10-05.
- Rationale/impact: Use project design tokens later; do not replace the component framework or ship a default demo.

## ELAB-V2-DEC-007 — Styling and icons

- Status: Accepted
- Decision: Tailwind CSS v4 and Lucide.
- Authority/date: current stakeholder Phase 0 direction, 2026-10-05.
- Rationale/impact: Accessible responsive behavior and reduced motion.

## ELAB-V2-DEC-008 — Server state

- Status: Accepted
- Decision: TanStack Query owns server state.
- Authority/date: current stakeholder Phase 0 direction, 2026-10-05.
- Rationale/impact: Do not duplicate server collections in Zustand.

## ELAB-V2-DEC-009 — Client state

- Status: Accepted
- Decision: Zustand only for appropriate global client/UI/session state; React local state for local UI.
- Authority/date: current stakeholder Phase 0 direction, 2026-10-05.
- Rationale/impact: Existing localStorage token transport is not approved secure session policy.

## ELAB-V2-DEC-010 — Forms

- Status: Accepted
- Decision: React Hook Form + Zod.
- Authority/date: current stakeholder Phase 0 direction, 2026-10-05.
- Rationale/impact: Frontend validation supports UX; backend validation remains authoritative.

## ELAB-V2-DEC-011 — Architecture

- Status: Accepted
- Decision: Clean Architecture + pragmatic DDD.
- Authority/date: current stakeholder Phase 0 direction, 2026-10-05.
- Rationale/impact: Keep inward dependency direction; focused modules/use cases; no speculative abstractions.

## ELAB-V2-DEC-012 — Deployment shape

- Status: Accepted
- Decision: Single React SPA → Go Fiber REST API → PostgreSQL.
- Authority/date: current stakeholder Phase 0 direction, 2026-10-05.
- Rationale/impact: Separate distributed services need later justification.

## ELAB-V2-DEC-013 — Business authority

- Status: Accepted
- Decision: Server-authoritative rules and authorization.
- Authority/date: current stakeholder Phase 0 direction, 2026-10-05.
- Rationale/impact: UI guards and frontend calculations cannot enforce stock, eligibility or privileges.

## ELAB-V2-DEC-014 — Development

- Status: Accepted
- Decision: Docker/Compose, Linux-compatible commands, Git/GitHub workflow.
- Authority/date: current stakeholder Phase 0 direction, 2026-10-05.
- Rationale/impact: Do not invent a GitHub remote or assume Docker/deployment checks passed.

## ELAB-V2-DEC-015 — Infrastructure restraint

- Status: Accepted
- Decision: No premature Redis, Kafka, RabbitMQ, Kubernetes, microservices, GraphQL, CQRS or event sourcing.
- Authority/date: current stakeholder Phase 0 direction, 2026-10-05.
- Rationale/impact: A later explicit requirement must justify any addition.

## ELAB-V2-DEC-016 — Phase boundary

- Status: Accepted
- Decision: Phase 0 rebaseline/adaptation only.
- Authority/date: current stakeholder Phase 0 direction, 2026-10-05.
- Rationale/impact: No business features, V1 Firebase access, commit, push or deployment.

## ELAB-V2-DEC-017 — Historical integrity

- Status: Accepted
- Decision: Preserve consequential history; integrity-sensitive operations use transactions/constraints.
- Authority/date: current stakeholder Phase 0 direction, 2026-10-05.
- Rationale/impact: Detailed audit/retention and archival rules still require domain design and policy decisions.

## ELAB-V2-DEC-018 — Workflow cleanup

- Status: Accepted
- Decision: Remove optional inherited .specify/.devin tooling and unrelated examples.
- Authority/date: current stakeholder Phase 0 direction, 2026-10-05.
- Rationale/impact: Project docs and AGENTS.md govern; no private-note/plugin dependency.

## ELAB-V2-DEC-019 — Super Administrator expansion

- Status: Deferred
- Evidence: capstone recommendations propose supervision across laboratories/departments; V1 has no Super Admin.
- Decision: no implementation in Phase 0. Whether a narrowly scoped FSMO role is needed remains Needs Stakeholder Input under ELAB-V2-OPEN-004/005. Multi-department expansion remains outside V2.
- Authority/date: current scope boundary, 2026-10-05; this does not reject the capstone recommendation.

## ELAB-V2-DEC-020 — Authoritative Go module path

- Status: Accepted
- Decision: use `github.com/Maaku050/elabtrack-v2/backend` for the module declaration and every active internal Go import.
- Authority/date: stakeholder closure request and read-only verification of `origin` as `https://github.com/Maaku050/elabtrack-v2.git` on `main`, 2026-10-05.
- Earlier evidence: the initial Phase 0 session could not inspect Git metadata and temporarily retained `github.com/fullstacktemplate/backend`. That blocked validation history remains intact in `PHASE0_REPORT.md`.
- Rationale/impact: the verified remote resolves the identity blocker without guessing a URL or upgrading dependencies; it does not authorize Phase 1 work.

## ELAB-V2-DEC-021 — Validated runtime configuration

- Status: Accepted
- Decision: One backend environment snapshot, strict typed development/test/production configuration, and startup validation before infrastructure/HTTP. Explicit blank/invalid values fail instead of silently defaulting. Production signing material must pass minimum length/diversity and placeholder checks; test signing material is explicit.
- Authority/date: stakeholder Phase 1A request, 2026-10-05; implemented in `internal/config` and bootstrap with deterministic tests.
- Rationale/impact: prevents accidental production use of published starter keys. This accepts configuration safeguards only; JWT verification, transport, storage and session design remain unfinished.

## ELAB-V2-DEC-022 — Explicit verified PostgreSQL transport

- Status: Accepted
- Decision: TCP PostgreSQL uses either discrete credentials encoded by `net/url` or a constrained PostgreSQL URL. Production explicitly requires `verify-full`, TLS 1.2 or newer, certificate chain/host verification, and no plaintext fallback. System trust or an existing supplied PEM CA bundle is supported; local development/test can explicitly disable TLS. Conflicting credential/TLS sources and ambient PG* settings fail.
- Authority/date: stakeholder Phase 1A request, 2026-10-05; implemented configuration/pool tests verify encoding and TLS fields without a production database.
- Rationale/impact: no alternative secure-transport architecture or provider has been approved. Passwordless/client-certificate authentication and provider-specific configurations require a future reviewed change; no certificates or hosting assumptions are invented.

## ELAB-V2-DEC-023 — Explicit schema changes and development seeds

- Status: Accepted
- Decision: API/container startup never runs migrations or seeds. Separate mutually exclusive CLI actions apply/roll back/read migration status or create scaffolds. Development seeding requires an explicit CLI request and is denied in production/test before IO, with a second guard inside the seeder.
- Authority/date: stakeholder Phase 1A request, 2026-10-05; implemented API command separation and seed guard tests.
- Rationale/impact: production migrations are a distinct authorized release step. Migration locking/atomic bookkeeping, least-privilege deployment credentials, and final admin provisioning remain outstanding. Existing auth SQL and synthetic local accounts establish no product roles or eligibility policies.
