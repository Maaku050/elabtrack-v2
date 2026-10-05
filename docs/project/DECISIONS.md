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

## ELAB-V2-DEC-024 — Current-account request authorization

- Status: Accepted
- Decision: Every retained access-token-protected route verifies bearer identity through one middleware implementation, then resolves a safe current PostgreSQL account snapshot. A nonzero verified account ID locates the row; current is_active and known temporary role drive access. Token/body email, role and status are never authorization evidence. Missing/deleted accounts return generic 401; inactive/unknown-role accounts return generic 403; lookup failures fail closed with safe 500.
- Authority/date: stakeholder Phase 1B request, 2026-10-06; implemented application account resolver, narrow read port, private typed request keys/helpers and real-router HTTP tests.
- Rationale/impact: demotion/deactivation takes effect on subsequent protected requests without waiting for token expiry. A request already authorized can finish from its snapshot; this does not accept final institutional roles, session revocation or JWT validation policy. Existing schema needs no migration.

## ELAB-V2-DEC-025 — Reversible production registration containment

- Status: Accepted (technical containment only)
- Decision: POST /api/v1/auth/register is mounted only for explicit development/test route configuration. Production/unknown/unset routing environments leave it absent (404). Account-creation service/domain code and local-only adapter remain reusable; no alternative provisioning endpoint is added.
- Authority/date: stakeholder Phase 1B request, 2026-10-06; production/local/unknown environment route tests and creation-side-effect assertions.
- Rationale/impact: prevents unreviewed production signup while OPEN-001/002/003 remain unresolved. Local account creation assigns only the inherited generic user role and hashes passwords; it establishes no product signup/eligibility/approval policy.

## ELAB-V2-DEC-026 — Narrow self-service and temporary privileged reads

- Status: Accepted (technical containment only)
- Decision: Retained self profile updates accept only display name. Identity/email, role/status, passwords/session fields, IDs and timestamps cannot be bound from client mutation payloads. A parameterized name-only UPDATE rechecks is_active and known role without rewriting a security snapshot. Existing generic admin user listing remains read-only, bounded and authorized from current database role. No arbitrary-account CRUD/provisioning/status/role endpoint is exposed.
- Authority/date: stakeholder Phase 1B request, 2026-10-06; strict DTO decoding, repository projection/write-scope tests and HTTP mass-assignment/self-vs-other tests.
- Rationale/impact: removes the inherited full-row self-update that could overwrite concurrent demotion/deactivation/password changes. Email/identity changes and final user administration require later reviewed flows; display-name containment does not settle institutional staff/admin permissions or identity verification.


## ELAB-V2-DEC-027 — Strict access-token contract

- Status: Accepted (technical security contract)
- Decision: HS256 alone; configured Phase 1A issuer; fixed `elabtrack-v2-api` audience; signature-protected `purpose=access`; required exp/iat/nbf; zero clock leeway; nonzero uid and matching canonical UUID subject. Time ordering must satisfy iat <= nbf < exp, and iat/nbf cannot be in the future. JWT role/email remain non-authoritative hints; Phase 1B current PostgreSQL principal still authorizes protected requests.
- Authority/date: stakeholder Phase 1C request, 2026-10-06. Issuance and verification change together, retaining the existing jwt/v5 dependency and typed configuration. The fixed audience avoids an unnecessary environment option for a single API.
- Impact: pre-1C access tokens lack the required audience/purpose and require re-login. Opaque refresh credentials are a separate class; they never enter JWT verification. No institutional role/session or browser transport policy is accepted.

## ELAB-V2-DEC-028 — Hash-only single-use refresh sessions

- Status: Accepted (technical security mechanism)
- Decision: Retain crypto/rand 32-byte opaque credentials encoded as 64 lowercase hex characters. SHA-256 of that canonical encoded secret is persisted/queried/revoked; raw values exist only for client transport. Unique token_hash and minimal ownership/time/replacement metadata replace raw persistence. Rotate through one explicit READ COMMITTED PostgreSQL transaction with a session FOR UPDATE lock, safe current account FOR SHARE lock, replacement insert and conditional once-only revocation. Return credentials after commit only; replay is generic 401 and does not revoke the successor.
- Authority/date: stakeholder Phase 1C request, 2026-10-06. This chooses local single-use/replacement semantics, not family-wide revocation or device/concurrent-session policy.
- Impact: paired migration 000003 preserves historical SQL and explicitly invalidates sessions on both up/down; hash values cannot recover bearer secrets. Live PostgreSQL concurrency, rollback and migration execution remain REQUIRED Phase 1G evidence. Lost commit acknowledgement/HTTP response cannot safely be retried as a second successful rotation; re-login may be required.

## ELAB-V2-DEC-029 — Bounded terminal-session cleanup capability

- Status: Accepted (technical record retention only)
- Decision: Keep expired/revoked refresh records for seven days after the earlier terminal timestamp, then permit one explicit maintenance batch of at most 1000 eligible rows with SKIP LOCKED. No startup cleanup or new scheduler. Operators must arrange sufficient regular batches, initially daily, and measure backlog/table size before deployment.
- Authority/date: stakeholder Phase 1C cleanup requirement, 2026-10-06; implemented repository/CLI capability. This does not establish institutional borrowing, disciplinary, audit, privacy or concurrent-session policy.
- Impact: active/recent sessions are retained; deleted replacement targets clear the optional link. Deployment ownership/cadence and PostgreSQL index/cleanup validation remain open operational work. No maintenance command was executed during implementation.

## ELAB-V2-DEC-030 — Browser credential transport and minimum CSRF policy

- Status: Accepted (technical browser contract)
- Decision: Access JWT stays in non-persisted application memory and travels as Authorization Bearer. Opaque refresh credentials travel only as the backend-issued host-only `elabtrack_v2_refresh` cookie: HttpOnly, Path=/api/v1/auth, explicit SameSite=Lax, Secure in production, expiry from the persisted new session. Development/test APP_ENV alone permits local HTTP Secure=false. Login/local register/refresh JSON exposes access and safe current-account fields only; refresh/logout no longer accept body credentials. Cookie-changing auth POSTs require an exact configured trusted Origin, including login/local register. Credentialed CORS retains explicit allowlists; no wildcard, missing/null-origin bypass or state-changing GET.
- Authority/date: stakeholder Phase 1D request, 2026-10-06; actual-router cookie/Origin/CORS tests and frontend transport/storage tests.
- Impact: production requires a same-site HTTPS SPA/API arrangement, normally the existing same-origin nginx proxy. Cross-site deployments are unsupported by this contract. Cookie Path limits delivery; it is not an authorization boundary. Bearer APIs still authorize current PostgreSQL accounts. No new schema or native-client compatibility endpoint is introduced. Broader HTTP review remains Phase 1E and real browser enforcement remains Phase 1G.

## ELAB-V2-DEC-031 — Memory session lifecycle and bounded refresh coordination

- Status: Accepted (technical client contract)
- Decision: One centralized Axios singleton owns credentialed transport, a shared refresh promise and a serialized queue for cookie-changing session operations. A non-persisted Zustand store owns access/user/status; no raw refresh field or JWT-derived authorization. Bootstrap runs once per application load, including repeated StrictMode effects; authentication failure becomes unauthenticated, unexpected failure becomes recoverable error with explicit retry. Eligible protected 401s share refresh, retry at most once and never recurse from session endpoints. Logout clears memory/cache immediately; generation checks prevent late responses restoring it. Known legacy auth storage keys are removed without reading their values; preferences remain.
- Authority/date: stakeholder Phase 1D request, 2026-10-06; persistence/bootstrap/concurrency/retry/logout tests using controlled transport adapters.
- Impact: coordination is within one document. Independent tabs are not coordinated; strict rotation/lost responses can require re-login. No family-wide replay policy, grace period or server single-use relaxation is accepted. Network failure leaves this document logged out but cannot prove deletion of an HttpOnly cookie or backend revocation; a later reload may restore a still-valid session. Live browser/multiple-tab verification remains Phase 1G.
