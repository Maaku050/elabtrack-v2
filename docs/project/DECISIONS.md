# eLabTrack V2 decision register

Recorded initially 2026-10-05; subsequent entries cite their own phase/date/authority. Status vocabulary: **Accepted**, **Proposed**, **Deferred**, **Needs Stakeholder Input**, **Rejected**. Acceptance records direction, not completed implementation. Unresolved product policies are in `OPEN_DECISIONS.md`; Phase 2 proposals do not accept institutional policy.

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

## ELAB-V2-DEC-032 — Explicit proxy trust and process-local abuse budgets

- Status: Accepted (technical HTTP perimeter contract)
- Decision: Empty TRUSTED_PROXIES uses the socket peer. Configured literal IPs/CIDRs enable a single shared client-IP resolver; only a trusted immediate source can supply XFF, validated as a complete bounded chain and traversed right-to-left to the first untrusted hop. Malformed/ambiguous/all-trusted chains fall back to the peer. Alternate IP/host/scheme headers and Fiber proxy parsing remain disabled. A single exact trusted XFP=https may establish transport authority for production API HSTS. Rate limits use this effective IP, with independent process-local fixed-window login 10, refresh 60, local register 5 and general/logout 120 defaults per minute; values are configuration-backed and 429 uses safe Retry-After.
- Authority/date: stakeholder Phase 1E request, 2026-10-06; config/proxy/limiter/concurrency/normal-refresh tests.
- Impact: IP is an abuse signal, never identity. Operators must choose actual narrow proxy sources, sanitize headers and prevent bypass; no provider/topology chosen. NAT/proxied peers share budgets; process restarts reset them and multiple instances do not coordinate. No Redis, account lockout or institutional policy introduced. Runtime proxy/rate verification remains Phase 1G.

## ELAB-V2-DEC-033 — Exact credentialed CORS and separate serving policies

- Status: Accepted (technical browser/perimeter contract)
- Decision: Normalize configured origins to canonical HTTP(S) host/default-port form, trim configuration separators, reject paths/credentials/wildcards/malformed hosts/ports and require frontend membership. Incoming Origin matches exactly for both CORS and cookie CSRF checks; methods/headers are limited to current needs. API uses explicit nosniff/frame/referrer/permissions headers and production authoritative-HTTPS HSTS; SPA nginx owns document CSP and common headers, with same-origin scripts/API, no external/inline/eval scripts, inline styles for retained UI primitives, no framing/objects/base override. HTTP nginx template emits no HSTS; authoritative production HTTPS edge owns it.
- Authority/date: stakeholder Phase 1E request, 2026-10-06; installed Fiber inspection, header/CORS/Origin tests and existing nginx/Docker serving model.
- Impact: same-origin container API default aligns with its connect-src policy. Other serving/API origins require an explicit reviewed policy and browser verification. nginx overwrites XFF/XFP from its actual transport and strips alternate forwarding fields; an additional ingress requires separate trust review. No hosting/TLS termination provider, CSP reporting platform or document UI redesign accepted.

## ELAB-V2-DEC-034 — Narrow parser and leakage safeguards

- Status: Accepted (technical HTTP safety contract)
- Decision: Retain 1 MiB default global ceiling, add 16 KiB auth POST ceiling before credential work, require JSON on current JSON mutation routes and strict-bind login. Explicit 10/15/60-second read/write/idle defaults and 8 KiB header buffer. Move recovery outermost; log safe panic type/source stack and error type rather than arbitrary values. Framework/validation HTTP text is generic; login account/status/lookup failures share one generic 401 envelope while service password-before-status ordering remains.
- Authority/date: stakeholder Phase 1E request, 2026-10-06; HTTP parser/body/panic/redaction/health/login tests.
- Impact: no full envelope/error taxonomy/logging/correlation redesign. Health shape stays unchanged; network timing equalization is not claimed. Future upload/report routes require separately reviewed bounds. Live HTTP/proxy/browser enforcement remains Phase 1G.


## ELAB-V2-DEC-035 — Foundation API contract and validation

- Status: Accepted (technical contract)
- Decision: Retain success/message/data/meta successes and bodyless 204; all API failures add required error.code/message/requestId with optional JSON-field-to-message details. Central wrapped-error mapping; safe generic infrastructure/panic responses, stable uppercase foundation taxonomy and current 409 duplicate semantics. One 400 syntax/semantic input policy, with transport 413/415; no speculative product pagination/errors.
- Authority/date: stakeholder Phase 1F request, 2026-10-07.
- Impact: Health adopts the envelope and becomes liveness; /ready checks existing dependency with standard 503. Client error normalization changes with the API. Phase 1E generic login lookup containment remains, now internally observable. See [contracts](../API_CONTRACTS.md); no product permission policy accepted.

## ELAB-V2-DEC-036 — Server-owned correlation and structured operational logging

- Status: Accepted (technical contract)
- Decision: Ignore all incoming request IDs; generate UUIDv4, return X-Request-ID and error.requestId, propagate a private-key context.Context value to ports. Use existing Zap for structured completion, unexpected failure, panic and lifecycle events. Route templates/effective IP/status/numeric milliseconds/code/type/operation only; no raw paths/query, headers/bodies/config/error values, credentials/digests or unnecessary PII. Panic source stack stays server-side. Production retains JSON default; injectable core enables assertions.
- Authority/date: stakeholder Phase 1F request, 2026-10-07.
- Impact: Correlation is not authentication/idempotency; no vendor chosen. Expected password failures are not ERROR. Live parser/ingress/log propagation still requires Phase 1G verification. This supersedes plaintext logging and arbitrary incoming-ID preservation.

## ELAB-V2-DEC-037 — Operational logs and durable business evidence

- Status: Accepted (architectural distinction)
- Decision: Operational logs serve debugging, operations, security and performance, with later chosen retention/rotation/access. They do not substitute for future durable business audit events/history/ledger evidence recording who changed what and when under approved domain integrity and retention policy.
- Authority/date: stakeholder Phase 1F direction and supplied boss-rebuild audit finding, 2026-10-07.
- Impact: No business audit schema/workflow or retention policy is implemented in this phase; no observability vendor integration.

## ELAB-V2-DEC-038 — Audience-specific responsive UX foundation

- Status: Accepted (confirmed future design direction)
- Decision: React + TypeScript + Vite, Tailwind CSS v4, shadcn/ui and Lucide remain the eLabTrack V2 UI foundation. Borrower/customer experience is mobile-first: mobile composition first, tablet/desktop progressively enhanced. The original dedicated touch-first kiosk-shell assumption is **superseded by DEC-061 (2026-10-08)**: Interactive Kiosk is the normal borrower catalog/cart/request experience, mobile-first and responsive; no device-authenticated hardware shell. Staff/admin is desktop/tablet-first, responsive on smaller screens where practical; complex administration need not be forced into tiny mobile layouts.
- Authority/date: explicit stakeholder Phase 1F UX direction, 2026-10-07.
- Impact: These are accepted design principles, not approved product screens/permissions/kiosk workflows. Existing 62 primitives and components.json remain unchanged. No mockups or feature UI created.

## ELAB-V2-DEC-039 — Mockups before feature UI completion

- Status: Accepted (future delivery process)
- Decision: Phase 3A owns UX architecture, wireframes and high-fidelity mockups: borrower mobile views first, tablet/desktop adaptations and staff/admin layouts. The earlier dedicated kiosk layouts are **superseded by DEC-061**: one coherent borrower catalog/cart flow satisfies Interactive Kiosk. Phase 3B implements the approved mockups through the shadcn-based eLabTrack design system, application shell and reusable domain UI components. Later feature UI is complete only after comparison against approved mockups and responsive behavior.
- Authority/date: explicit stakeholder Phase 1F direction, 2026-10-07.
- Impact: ROADMAP records the split; Phase 3A/3B have not begun. Product decisions still gate dependent designs.

## ELAB-V2-DEC-040 — Migration-runner hardening phase boundary

- Status: Accepted (sequencing only)
- Decision: Phase 1H explicitly owns DDL/bookkeeping atomicity, migration checksum validation, advisory lock/concurrent-runner protection and correct last-applied lookup failure semantics. It follows Phase 1G and precedes heavy reliance on business-domain migrations.
- Authority/date: stakeholder Phase 1F request and stated boss-rebuild comparison, 2026-10-07; current runner source independently exhibits separate commits and missing checksum/lock handling.
- Impact: No runner/SQL change or migration execution during Phase 1F. Prior Phase 1A partial SEC-008 status remains; no automatic authorization for 1G/1H.

## ELAB-V2-DEC-041 — Immutable paired transactional migrations

- Status: Accepted (technical migration contract)
- Decision: Applied up/down files are immutable. SHA-256 covers a versioned marker and length-framed raw up/down bytes; the digest is stored with canonical version/applied timestamp and verified on inspection/execution. Six-digit numeric ordering, unique version mapping, paired files and applied-prefix consistency are mandatory; numbering gaps are allowed. SQL and version INSERT/DELETE share one PostgreSQL transaction per file. Transaction/session-control SQL is rejected; exceptional nontransactional support and dirty recovery are deferred until needed.
- Authority/date: stakeholder Phase 1H direction, 2026-10-07; current three historical pairs remain unchanged.
- Impact: No silent checksum replacement or history repair. Explicit development-only adoption attests known local foundation provenance, accepts only baseline-verified foundation prefixes and preserves timestamps; unknown external histories require separate reconciliation. Commit acknowledgement loss requires status inspection, not blind assumption of rollback.

## ELAB-V2-DEC-042 — PostgreSQL migration exclusion and explicit execution

- Status: Accepted (technical operator contract)
- Decision: Up/down/status/legacy adoption acquire project key `0x454c41424d494752` with pg_try_advisory_lock on one dedicated connection before history access. Contention fails immediately; connection/lock-query acquisition is bounded to five seconds. Hold through the entire command and close the connection on exit. Migration and seed CLI actions remain explicit and separate from API/Docker startup.
- Authority/date: stakeholder Phase 1H direction, 2026-10-07; verified across real concurrent connections and separate CLI processes.
- Impact: Status refuses concurrent mutation instead of reporting a mixed snapshot. No in-process mutex substitutes for PostgreSQL. No automatic release migration, repair or startup adoption.

## ELAB-V2-DEC-043 — Schema-owner migration and runtime DML credentials

- Status: Accepted (infrastructure privilege contract)
- Decision: Bootstrap administration creates the database/extension and two non-superuser logins. Stable `elabtrack_migrator` owns public schema/application tables/tracking; `elabtrack_runtime` gets CONNECT, schema USAGE and SELECT/INSERT/UPDATE/DELETE on named application tables only. No schema CREATE, role membership/cluster administration/grant option/tracking access. Future approved migrations grant explicit per-object DML and only needed sequences/functions in their transaction; blanket table/default runtime grants are excluded.
- Authority/date: stakeholder Phase 1H direction, 2026-10-07; real API/auth/browser/privilege verification in disposable PostgreSQL.
- Impact: Typed MIGRATION_DATABASE_URL selects operator credentials and shares production verify-full/no-fallback policy. Production CLI requires it, a different user and the same target; API needs no migration secret. Development/test permit a documented absent-URL local fallback; fresh Compose separates roles by default, with an explicit tools-profile CLI job. Seeds use the operator path, retain development-only guards; cleanup/readiness use runtime. Existing volumes/ownership are never silently converted. Hosted role/TLS rollout remains deployment verification, not product-role policy.

## ELAB-V2-DEC-044 — Required Phase 1I cross-tab session coordination

- Status: Accepted (sequencing and engineering requirement only)
- Decision: Phase 1I — Cross-Tab Session Coordination precedes Phase 2/product UI. Preserve strict single-use rotation while addressing shared-cookie concurrent refresh and peer logout/session lifecycle; do not broadcast/persist usable credentials. The detailed coordination/fallback design belongs to separately authorized 1I.
- Authority/date: stakeholder Phase 1H direction, 2026-10-07, accepting the Phase 1G adverse finding as required follow-up.
- Impact: One tab's losing refresh 401 can clear the winner's new cookie; another tab can retain authenticated memory after peer logout. Both recur in 1H browser regression. No 1I implementation, replay grace, institutional session policy or later-phase authorization is included.

## ELAB-V2-DEC-045 — Cooperative browser ownership without credential sharing

- Status: **Accepted — engineering implementation**, 2026-10-08.
- Decision: One dedicated per-document session coordinator uses the origin/storage-partition Web Lock `elabtrack_v2.session-cookie.v1` for all retained cookie-changing endpoints. BroadcastChannel `elabtrack_v2.session.v1` carries only a strict versioned union of lifecycle/attempt metadata. Local refresh single-flight remains; every waiting document obtains its own memory access credential through a subsequent serialized normal refresh. No credential, digest, header, password, account object or authority is shared/persisted.
- Authority/evidence: Stakeholder's explicit Phase 1I scope/requirements, 2026-10-07; 2026-10-08 source/unit and actual Chromium two/three-tab, bootstrap, storage and owner-close verification. No role/provisioning/product policy is inferred.
- Impact: This supersedes DEC-031's within-document-only coordination limitation and implements DEC-044. Ephemeral UUID tab/attempt IDs, high-resolution timestamp/Lamport-style monotonic epoch, completion ordering and existing local generation/token fences reject stale operations. The browser's lock service is the mutex; channel messages never confer authorization. No server rotation change or bootstrap endpoint.
- Limits: Axios active request timeout 15s; lock acquisition wait 20s, aborting only the queued request. Never steal an active cookie mutation. Closure/crash releases native ownership without a completion broadcast; a suspended live owner causes a bounded recoverable error until it resumes/closes. A lost committed rotation response can still require re-login. No polling, token lock or coordination storage. Primary primitive contracts: [W3C Web Locks](https://www.w3.org/TR/web-locks/) and [WHATWG BroadcastChannel](https://html.spec.whatwg.org/multipage/web-messaging.html#broadcasting-to-other-browsing-contexts).

## ELAB-V2-DEC-046 — Rotation-safe refresh cookie rejection

- Status: **Accepted — engineering safety**, 2026-10-08.
- Decision: Any failed refresh leaves cookies untouched, including missing/malformed/unknown/expired/revoked/consumed credentials, invalid current accounts, denied Origin and unexpected storage/transaction/server failures. Success still installs a replacement. Trusted-Origin explicit logout revokes its presented session and clears the cookie, including safe storage-failure handling.
- Authority/evidence: Explicit Phase 1I cookie-race requirement; unchanged baseline produced one winner/one denial and no surviving cookie. Final actual late replay 401 has no Set-Cookie and cannot erase the winner; real 12-way PostgreSQL race remains one success/11 denials. HTTP tests cover all collapsed rejection classes.
- Impact: Supersedes unconditional failed-refresh clearing in DEC-030/Phase 1D. Even conclusive invalidity of a *presented* credential does not identify the browser cookie at response-delivery time; another login/rotation can have replaced it. Rejected HttpOnly credentials may remain until expiry, explicit logout or replacement; they remain unusable at the server. Do not disclose rejection classes or weaken replay denial. No new cookie/schema/family policy.

## ELAB-V2-DEC-047 — Fenced peer lifecycle and private-cache invalidation

- Status: **Accepted — engineering implementation**, 2026-10-08.
- Decision: Logout clears local memory immediately, broadcasts peer logout, and serializes server revocation behind pending cookie replies. A second completion hint handles delayed message/reply ordering. A separate login-change counter permits only a newer login intent to supersede queued logout; peer logout/invalidation cannot cancel revocation, including simultaneous logout by all tabs. Peer lifecycle clears memory and authenticated Query roots (`auth`, `users`) plus queries marked `meta.authenticated=true`; public health/status cache survives. Mutation history is cleared because current mutations hold private/credential inputs; removing cache cannot undo an HTTP mutation. New private features must declare classification and retain asynchronous generation guards.
- Authority/evidence: Phase 1I lifecycle/cache/failure requirements; unit cancellation/public-cache tests and actual three-tab logout, held committed refresh response, disable/current-account 403 and revoked-session refresh 401.
- Impact: Login/register intent and successful completion discard peers' old account presentation without transferring role/status/account objects or automatically signing them in. The completion hint covers peers that missed the intent. Reload/deliberate server restoration obtains trusted state. Exclusive refresh 401 infers generic invalidation among participating documents; without Web Locks it is an ambiguous denial and only clears the failing document. Network/5xx/Origin denial/coordination timeout do not globally log out peers. Current-account `/auth/me`/`/users/me` resolver 403 and persistent retried 401 invalidate with request-start fencing; ordinary resource 403 does not. Old token/generation/attempt failures cannot erase newer state.
- Browser limits: Chromium 140 verified. Web Locks without BroadcastChannel still serialize but cannot reliably propagate peer UI lifecycle; BroadcastChannel without Web Locks provides lifecycle hints without atomic refresh optimization. With neither, local single-flight and safe server cookie rejection remain, with possible local denial/re-login. Missed messages do not hold ownership; reload/stale-token request recovery consults the server. There is no durable client logout history or cross-device/access-token revocation guarantee. Institutional family/concurrent-session/immediate-access-revocation policy remains open.

## ELAB-V2-DEC-048 — Phase 2 design-only authorization and evidence order

- Status: **Accepted — scope and evidence contract**, 2026-10-08.
- Decision: Authorize FSMO domain/database documentation and read-only boss reference inspection. Authority is current stakeholder/Dean direction, current V2 decisions, capstone intent, audited V1 behavior, boss design reference, then engineering recommendation. Use the seven explicit Phase 2 evidence labels; current source establishes implementation facts, not institutional policy.
- Authority: Explicit current stakeholder Phase 2 request; expands source-policy detail without changing Phase 1 implementation.
- Impact: No business migrations/services/UI, Phase 3A/mockups, boss writes, commit/push/deployment. Phase 1A–1I is complete. At the original Phase2 handoff, policy-dependent state/schema was blocked pending confirmation. DEC-051–062 and [Phase2.5 report](PHASE2_5_REPORT.md) now supersede that gate; the Phase2 report remains historical evidence.

## ELAB-V2-DEC-049 — Retained canonical business history and distinct evidence

- Status: **Accepted — required design boundary**, 2026-10-08.
- Decision: Preserve canonical borrowing transactions through denial/cancellation/completion; distinguish borrowing history, immutable return evidence (Staff/Admin-recorded condition/quantity events, excluding photographs or return attachments), inventory movement ledger, financial assessment/adjustment history and durable business audit. Operational logs do not replace evidence. Do not replicate V1's terminal copy/delete or independent conflicting fine truths.
- Authority: Explicit Phase 2 Parts 18/21/22 and existing DEC-017/037. This accepts the required integrity boundary, not a fine amount, payment procedure, retention duration or role matrix.
- Impact: [Domain designs](../domain/DOMAIN_MODEL.md) recommend exact tables/immutability/concurrency for later review; none implemented. Assessment, administrative settlement and actual payment/revenue remain distinct; OPEN-007/024 were subsequently resolved by DEC-059/060; retention OPEN-026 remains nonblocking for core UX but gates destructive cleanup.

## ELAB-V2-DEC-050 — FSMO aggregate-stock design constraint

- Status: **Accepted — required scope constraint**, 2026-10-08.
- Decision: Phase 2 designs catalog stock pools and aggregate quantities, without per-physical-unit/serialized tracking unless a later requirement changes. Do not adopt boss organization/unit/lab/global-role scope, generic file authorization or unsafe importer/migrations wholesale.
- Authority: Explicit Phase 2 scope/Parts 4/26/41; charter/AGENTS; V1 audit 04/06 evidence.
- Impact: The original six-bucket/valuation design was provisional. Phase2.5 selects a four-count physical hybrid plus replacement obligations (DEC-058/062); OPEN-020/023 are resolved and OPEN-022 retains only nonblocking original-disposition/equivalence details. This entry does not accept boss's stock/role/financial defaults.

## Phase 2.5 authoritative working product decisions

Authority for DEC-051–061: project owner's explicit **“PHASE 2.5 — PRODUCT DECISION INTEGRATION & DOMAIN REBASELINE”** direction, 2026-10-08; owner implemented V1 and knows the former FSMO manual workflow. These are accepted working V2 decisions, not assertions of deployment or independent institutional approval. D numbers trace all29 owner statements in [BUSINESS_RULES](../domain/BUSINESS_RULES.md#decision-integration-matrix). Exact schema/locks and identified engineering details remain design selections for later review. Phase2.5 authorizes documentation only; no Phase3A, migrations, business code, boss changes or Git writes.

## ELAB-V2-DEC-051 — Borrower population, account gate and fine review

- Status: **Accepted — authoritative working product decision**.
- Owner statements: D1, D2, D16.
- Authority/date: explicit Phase2.5 owner direction, 2026-10-08.
- Decision: Active registered BORROWER may request; expected categories include Student and Faculty and do not grant roles. Admin can deactivate; inactive accounts cannot log in or perform normal authenticated borrower actions. No independent eligibility/suspension subsystem or automatic outstanding-fine submission block. Existing obligations/history survive; staff resolves them. Old fine is reviewed face-to-face and can inform a reasoned human denial.
- Supersession / impact: Resolves OPEN-002/008 and eligibility portion of OPEN-009/020; audited student UI/status behavior remains historical.
- Affected documents: [domain baseline](../domain/DOMAIN_MODEL.md), related rules/state/schema/invariants/API, OPEN_DECISIONS, SOURCE_OF_TRUTH, ROADMAP and [Phase2.5 report](PHASE2_5_REPORT.md).

## ELAB-V2-DEC-052 — Three product roles and named administrative authority

- Status: **Accepted — authoritative working product decision**.
- Owner statements: D3, D26.
- Authority/date: explicit Phase2.5 owner direction, 2026-10-08.
- Decision: BORROWER, STAFF, ADMIN are future product-facing roles. Admin inherits Staff operations; fine clearing, deactivation, bulk imports, Staff/Admin management, exceptional stock corrections, policy/terms and administrative audit/reports belong to Admin. Multiple named Admins, never one shared credential. Staff operational history/returns/replacements do not imply Admin fine authority. Current generic user/admin source remains unchanged until an approved implementation.
- Supersession / impact: Resolves OPEN-003; excludes Super Admin expansion OPEN-004/005. Future permission matrix in BUSINESS_RULES is authoritative design intent.
- Affected documents: [domain baseline](../domain/DOMAIN_MODEL.md), related rules/state/schema/invariants/API, OPEN_DECISIONS, SOURCE_OF_TRUTH, ROADMAP and [Phase2.5 report](PHASE2_5_REPORT.md).

## ELAB-V2-DEC-053 — Provisioned accounts and secure bulk onboarding

- Status: **Partially superseded by DEC-070**, 2026-10-08. The Staff creation grant and generic bulk population below are historical; only Admin now creates accounts, and standard bulk workflows are Student-only. No public signup and secure onboarding remain approved.
- Owner statements: D4.
- Authority/date: explicit Phase2.5 owner direction, 2026-10-08.
- Decision: No public self-registration in current intended workflow. Authorized Staff/Admin provision Borrower accounts; Staff cannot provision privileged roles. Admin bulk Borrower creation/import is required. Preserve conceptual identity/category/program/contact fields where purposeful; do not preserve insecure plaintext spreadsheet passwords as a requirement. Activation/initial-password delivery remains implementation security design.
- Supersession / impact: Resolves OPEN-001 and registration conflicts; Phase1 development/test-only signup containment remains source fact, not public product capability.
- Affected documents: [domain baseline](../domain/DOMAIN_MODEL.md), related rules/state/schema/invariants/API, OPEN_DECISIONS, SOURCE_OF_TRUTH, ROADMAP and [Phase2.5 report](PHASE2_5_REPORT.md).

## ELAB-V2-DEC-054 — Reservation requests and atomic physical issue

- Status: **Accepted — authoritative working product decision**.
- Owner statements: D5, D6, D10, D16.
- Authority/date: explicit Phase2.5 owner direction, 2026-10-08.
- Decision: Borrower submission immediately reserves quantities transactionally. Staff/Admin reviews request/accountability in person and approves WHILE handing over: one PENDING→CHECKED_OUT edge, no approved-waiting state. Staff/Admin direct checkout immediately issues against available stock for active existing Borrower with date/time due and current server validation. Old fine does not prohibit submission; staff may deny during human review.
- Supersession / impact: Resolves OPEN-020/021 and supersedes separate-approval engineering preference; no borrowing code is implemented.
- Affected documents: [domain baseline](../domain/DOMAIN_MODEL.md), related rules/state/schema/invariants/API, OPEN_DECISIONS, SOURCE_OF_TRUTH, ROADMAP and [Phase2.5 report](PHASE2_5_REPORT.md).

## ELAB-V2-DEC-055 — 24-hour expiry and retained pending terminal outcomes

- Status: **Accepted — authoritative working product decision**.
- Owner statements: D7, D8, D9.
- Authority/date: explicit Phase2.5 owner direction, 2026-10-08.
- Decision: Working pending TTL is24 hours after submission, candidate configurable policy later; no operating-hours awareness. Expiry becomes EXPIRED and releases holds/history. Borrower cancels only own PENDING; cancellation releases/retains. Staff/Admin denies PENDING with required borrower-visible reason and releases/retains. Administrative pending cancellation is optional scoped engineering recommendation only; no issued-loan cancellation/delete.
- Supersession / impact: Resolves OPEN-013/014; history boundary DEC-049 includes EXPIRED. Atomic expiry versus approval is required.
- Affected documents: [domain baseline](../domain/DOMAIN_MODEL.md), related rules/state/schema/invariants/API, OPEN_DECISIONS, SOURCE_OF_TRUTH, ROADMAP and [Phase2.5 report](PHASE2_5_REPORT.md).

## ELAB-V2-DEC-056 — Due date and time without fixed duration maximum

- Status: **Accepted — authoritative working product decision**.
- Owner statements: D11, D12.
- Authority/date: explicit Phase2.5 owner direction, 2026-10-08.
- Decision: Every checked-out borrowing has required due_at containing date and time. Interpret local values in Asia/Manila and persist absolute timestamptz. Remove the V1 seven-day maximum; no new maximum, grace or quota. Nullable configurable future maximum may be added only under later approved policy. Partial return/replacement retains original issued deadline.
- Supersession / impact: Resolves duration/time portions of OPEN-015/020; no organizational timezone hierarchy.
- Affected documents: [domain baseline](../domain/DOMAIN_MODEL.md), related rules/state/schema/invariants/API, OPEN_DECISIONS, SOURCE_OF_TRUTH, ROADMAP and [Phase2.5 report](PHASE2_5_REPORT.md).

## ELAB-V2-DEC-057 — First-use terms and material version changes

- Status: **Accepted — authoritative working product decision**.
- Owner statements: D13.
- Authority/date: explicit Phase2.5 owner direction, 2026-10-08.
- Decision: Borrower accepts current terms once at activation/first use; immutable user/version/accepted_at evidence applies to subsequent requests while current. Material new active version requires acceptance before another request. No per-loan checkbox. Existing pending request retains its submission acceptance. Engineering recommendation: direct issue also requires current borrower acceptance, never staff impersonation.
- Supersession / impact: Resolves OPEN-016 and terms-binding OPEN-023; final approved text/secure onboarding mechanics later.
- Affected documents: [domain baseline](../domain/DOMAIN_MODEL.md), related rules/state/schema/invariants/API, OPEN_DECISIONS, SOURCE_OF_TRUTH, ROADMAP and [Phase2.5 report](PHASE2_5_REPORT.md).

## ELAB-V2-DEC-058 — Replacement obligations and operational completion

- Status: **Accepted — authoritative working product decision**.
- Owner statements: D14, D17, D18, D19, D20, D21, D22, D23.
- Authority/date: explicit Phase2.5 owner direction, 2026-10-08.
- Decision: Staff/Admin alone records authoritative return condition/quantities. **D17 clarified by explicit owner correction, 2026-10-08:** a borrower may physically show a photo personally stored on their own phone during the face-to-face return; Staff/Admin may inspect actual equipment if desired. The photo is entirely outside eLabTrack. Borrowers cannot mark returned or submit return evidence; no photo acceptance, upload, transmission, storage, retention, attachment, return-evidence file model or interface/workflow step. Equipment catalog images remain separate. Partial returns keep loan/deadline open. Damage/loss creates replacement quantities rather than automatic equipment-price monetary liability; those disposed quantities leave physical checked_out. Appropriate type/accepted equivalent replacement reduces liability/restores usable stock without erasing incident. Complete iff ALL physical outstanding0 AND replacement outstanding0. Overdue depends on unresolved borrowing, including replacement-only open loan. Damaged-original disposition/equivalence details are nonblocking for borrowing UX; no mandatory repair workflow.
- Supersession / impact: Resolves core OPEN-022/023 and premature completion; detailed original disposition/equivalence remain OPEN-022.
- Affected documents: [domain baseline](../domain/DOMAIN_MODEL.md), related rules/state/schema/invariants/API, OPEN_DECISIONS, SOURCE_OF_TRUTH, ROADMAP and [Phase2.5 report](PHASE2_5_REPORT.md).

## ELAB-V2-DEC-059 — PHP10 overdue fine and selected deterministic day arithmetic

- Status: **Accepted — authoritative working product decision**.
- Owner statements: D15.
- Authority/date: explicit Phase2.5 owner direction, 2026-10-08.
- Decision: Confirmed working rate PHP10 per overdue day; continues through partial return/replacement until completed_at, then freezes. Selected ENGINEERING DETAIL is ceil(max(0,effective_end-due_at)/24h), effective_end=now while unresolved or completed_at when complete. At exact due0,1 minute10PHP, exactly24h10PHP,24h+1minute20PHP. Use bounded integer centavos/duration and one shared formula, no daily fine mutation. Pending/denied/cancelled/expired never accrue.
- Supersession / impact: Resolves OPEN-006/015/024. Ceiling interpretation is adopted for this baseline; no evidenced calendar-day conflict requiring a fresh gate. Legacy historical calculations retain provenance rather than silent recalculation.
- Affected documents: [domain baseline](../domain/DOMAIN_MODEL.md), related rules/state/schema/invariants/API, OPEN_DECISIONS, SOURCE_OF_TRUTH, ROADMAP and [Phase2.5 report](PHASE2_5_REPORT.md).

## ELAB-V2-DEC-060 — Focused fine history and full Admin clearance

- Status: **Accepted — authoritative working product decision**.
- Owner statements: D24, D25, D26, D27.
- Authority/date: explicit Phase2.5 owner direction, 2026-10-08.
- Decision: Confirmed primary money liability is overdue fine; damage/loss remains replacement-based. Admin alone clears entire current outstanding via PAID, WAIVED or OTHER_RESOLUTION, preserving original/final assessed amount, cleared amount/time/actor/method and optional note. No partial payments, online gateway or allocation architecture. Selected ENGINEERING RECOMMENDATION: overdue_fines plus immutable fine_clearances; live clearance covers full outstanding at its time, continued unresolved accrual can produce a later delta for another full clear. Fine clearance does not close/extend borrowing; assessment, recorded PAID and waiver/other are distinct reports.
- Supersession / impact: Resolves OPEN-007/024 and supersedes generic charge/adjustment/payment candidates; live checkpoint arithmetic is design detail, not a partial-payment product feature.
- Affected documents: [domain baseline](../domain/DOMAIN_MODEL.md), related rules/state/schema/invariants/API, OPEN_DECISIONS, SOURCE_OF_TRUTH, ROADMAP and [Phase2.5 report](PHASE2_5_REPORT.md).

## ELAB-V2-DEC-061 — Interactive Kiosk is the responsive borrower catalog/cart

- Status: **Accepted — authoritative working product decision**.
- Owner statements: D28, D29.
- Authority/date: explicit Phase2.5 owner direction, 2026-10-08.
- Decision: Interactive Kiosk is browse/search/filter/select quantities/cart/review/create request in normal borrower web/mobile flow. No dedicated hardware, terminal shell, shared-device identity, device registration/tokens, handoff codes or special authentication. Borrower mobile-first, Staff/Admin desktop/tablet-first responsive. Preserve academic name with functional meaning.
- Supersession / impact: Resolves OPEN-010. Supersedes kiosk-shell portions of DEC-038/039 and old Phase11 hardware assumption. Boss device/handoff candidates are rejected for current scope.
- Affected documents: [domain baseline](../domain/DOMAIN_MODEL.md), related rules/state/schema/invariants/API, OPEN_DECISIONS, SOURCE_OF_TRUTH, ROADMAP and [Phase2.5 report](PHASE2_5_REPORT.md).

## ELAB-V2-DEC-062 — Selected physical stock and replacement design

- Status: **Proposed — selected engineering baseline**, distinct from confirmed product policy.
- Authority/date: Phase2.5 explicitly authorizes re-evaluating the inventory model and schema; engineering recommendation, 2026-10-08.
- Decision: Choose hybrid `total_tracked=available+reserved+checked_out+damaged_held`; loss removes physical total, replacement acquisition adds available/total. Damage originals remain nonusable held until later approved disposition. Replacement liability derives required minus immutable acceptances and never inflates physical stock; completion requires physical0 and replacement0. Four physical counts plus obligation/events replace six accounted buckets.
- Rationale: Pure usable counts hide holds/custody/originals; six physical/history buckets conflate lost history with present stock. Acquisition of a replacement while damaged original retained legitimately adds a second physical unit. Ledger/counters/acceptance/reconciliation use one normalized map with sorted locks.
- Impact: [DOMAIN_MODEL arithmetic](../domain/DOMAIN_MODEL.md#selected-inventory-model-and-alternatives), DATA_MODEL, STATE_MACHINES, INVARIANTS, API_RESOURCE_DRAFT and boss salvage map; no repair/disposal policy guessed, no code/schema implemented. Engineering mechanics do not reopen confirmed core product flows.

## ELAB-V2-DEC-063 — Approved navy/violet visuals and scoped Phase3B implementation

- Status: **Accepted — owner visual target and explicit phase authorization**, not a new business-policy decision.
- Authority/date: Owner Phase3B brief and supplied approved handoff, 2026-10-08.
- Decision: Phase3A.2 is visually approved using `docs/ux/approved/` navy/indigo/violet PNGs. Those supersede the original evergreen visual proposal. Locked Phase2.5 rules override erroneous generated copy; Phase3A.1 remains behavioral structure where compatible. The36-file manifest contains35 screen images plus brand, representing36 screen targets because S01 has two frames.
- Authorized implementation: Existing React/TypeScript/Tailwind v4/shadcn/Lucide visual tokens, shared compositions, Borrower and Staff/Admin shells, four development-only fixture previews, light/dark responsiveness/accessibility and browser comparison. Remaining32 screens are binding later feature targets, not production implementations in this phase. At the original implementation handoff, dark/wide translations and disclosed density/font differences awaited owner disposition. **DEC-064 subsequently records their acceptance within the current baseline.**
- Safety: No Phase1 auth/session replacement, fake role session, business APIs/services/migrations, domain-policy changes, return-photo feature, boss port, commit/push/deployment or automatic Phase4. Supplied reconstructed seal stays separately replaceable and is not authenticated official master.
- Evidence: [Approved handoff](../ux/approved/README.md), [map](../ux/APPROVED_VISUAL_IMPLEMENTATION_MAP.md), [fidelity contract](../ux/VISUAL_FIDELITY_CONTRACT.md), [Phase3B report](PHASE3B_REPORT.md). The original Phase3A.2 report remains historical, including its superseded proposal and pending-review snapshot.

## ELAB-V2-DEC-064 — Phase 3B implementation visually approved

- Status: **Accepted — explicit project-owner approval**, 2026-10-08.
- Authority: Current Phase 3B closeout brief confirms the owner reviewed real implementation screenshots and found them faithful to the approved mockups.
- Decision: **PHASE 3B — COMPLETE AND OWNER VISUALLY APPROVED.** Borrower Home, Equipment Catalog, Staff Dashboard and Pending Requests; FSMO navy/violet identity, compact branding, retained shadcn foundation, light/dark themes, mobile responsiveness and desktop operational layouts are accepted. Minor implementation adaptations disclosed in the Phase 3B report are part of the accepted current baseline.
- Boundaries: Accessibility, readable typography and control sizes remain requirements. Every future screen must follow its individual approved PNG and receive fresh runtime verification; the 32 remaining screens are not implemented. The supplied reconstructed seal remains replaceable by an authenticated original vector. This is visual approval, not product-policy change or permission to start Phase 4.
- Evidence: [Phase 3B report](PHASE3B_REPORT.md), [fidelity contract](../ux/VISUAL_FIDELITY_CONTRACT.md), [local readiness report](LOCAL_ENVIRONMENT_READINESS_REPORT.md).

## ELAB-V2-DEC-065 — Product authentication and three-role reconciliation

- **Status:** Accepted; Phase 4A, 2026-10-08.
- **Authority:** Explicit owner Phase 4A request and confirmed Phase 2.5 DEC-051/052/053; implementation verified through real PostgreSQL/HTTP/Chromium and preserved Phase 1 regression tests.
- **Decision:** Product roles are exactly BORROWER, STAFF, ADMIN; Student/Faculty remain borrower categories. 000004 maps foundation `user`→BORROWER and `admin`→ADMIN without changing account identity/password/status/timestamps or session rows. STAFF is a new permitted role, never assigned by signup. Down reverses representable values and atomically refuses while any Staff would lose authority. Historical migrations are unchanged.
- **Authorization:** Current PostgreSQL account role/status is authoritative for every protected request. A central permission matrix/reusable guard protects existing privileged reads; a stale JWT hint cannot preserve privileges after demotion/deactivation. Public registration is absent in every environment, superseding DEC-025's temporary development/test exception. Internal creation service code remains unmounted; Phase 5 owns product provisioning.
- **Browser:** Reuse the memory-only JWT/HttpOnly cookie, existing SessionBootstrap, Web Locks, strict non-secret BroadcastChannel, bounded recovery and private-query fencing. Real login navigates to an exact role-permitted saved destination or Borrower Home/Staff Dashboard. Protected routes check current /auth/me before rendering. Logout clears local memory/private caches immediately and retains visible pending/failure/retry feedback; it preserves existing server revocation semantics.
- **Presentation:** Reuse approved Borrower/Staff/Admin shells and semantic light/dark tokens. Real authenticated routes expose honest feature placeholders and safe read-only own account data. Dev-only synthetic previews remain isolated and absent from production. New rendered comparisons/deviations are documented; Phase 3B owner approval is not an approval claim for all new Phase 4A details.
- **Deferred:** Phase 4B versioned terms/current-acceptance gate and activation details; Phase 5 secure provisioning/import/status/privileged management and reviewed recovery mechanism. No business services, category persistence or institutional terms wording is added.
- **Evidence:** [Phase 4A report](PHASE4A_REPORT.md), [API contracts](../API_CONTRACTS.md), [isolated reproduction](../../integration/PHASE4A.md), [browser evidence](../ux/verification/phase4a/README.md).


## ELAB-V2-DEC-066 — Versioned terms and first-use foundation

- **Status:** Accepted engineering implementation, 2026-10-08. Official institutional text remains Needs Stakeholder Input.
- **Authority:** Explicit owner Phase 4B request and clarification: no current official FSMO terms document is confirmed; implement infrastructure independently, never publish placeholders or record acceptance of unapproved official content. Confirmed product cadence is DEC-057.
- **Decision:** Immediate mandatory Admin-only immutable publication with exact UTF-8 plain-text hash, durable publisher/time and expected-pointer conflict check. Unique authenticated Borrower/version/server-time receipts, idempotent original evidence, immutable history/FK retention and transactional shared/exclusive publication locking. No draft CMS/effective scheduling. This narrows earlier proposed terms effective-date/advisory-selector mechanisms to the implemented singleton-row model; borrowing-policy scheduling remains future design.
- **UX:** First-use and updated-version review use the approved shells/palette and initially unchecked consent. Home/catalog initiation paths require actual current evidence; account/existing obligations/notifications and logout remain available. Missing content never means acceptance. No Staff/Admin Borrower consent mandate.
- **Integration:** Future Phase 7 submission/direct checkout must authorize actor/target and call `RequireCurrentAcceptance` inside its command transaction, holding the publication lock through commit and binding the receipt/borrower with a composite FK. Existing pending submission evidence is not retroactively invalidated. No live borrowing endpoint/enforcement is claimed.
- **Evidence:** New migration 000005 only; real database races/rollback/constraints/HTTP and Chromium verification; [report](PHASE4B_REPORT.md), [API contracts](../API_CONTRACTS.md), [reproduction](../../integration/PHASE4B.md). Historical migrations, approved assets, primitives and security architecture are preserved.

## ELAB-V2-DEC-067 — Official content and secure activation decision gates

- **Status:** Partially superseded by DEC-068/069/070, 2026-10-08. Official content, provider delivery and residual activation/recovery details remain pending. The original activation alternatives below are historical proposals; current category-specific email/separate-password policy and email-link recommendation are recorded in DEC-068/070.
- **Authority:** Owner clarification directs a credential-delivery decision gate and independent terms implementation. V1/capstone text remains reference material, not officially approved V2 content.
- **Required content decision (OPEN-016):** Provide/approve exact final FSMO terms text and initial version identifier, assign responsible Admin publication approval and subsequent material-version review. No draft is currently published in normal application storage.
- **Historical activation alternatives (superseded by DEC-068):** Choose (then-recommended engineering alternative) an expiring one-time activation link/code handed over in person by authorized FSMO personnel, or delivery of an expiring one-time link to an institutionally verified address. Approve identity verification, channel ownership, expiry, replacement/lost-delivery handling and whether activation sets the first password; define authenticated password changes/revocation and recovery responsibility before dependent code. These are proposals, not selected institutional policy. Neither alternative requires predictable/shared passwords or plaintext-password import columns.
- **Boundary:** Existing account creation hashes a supplied password; there is no activation, invitation, password-change-required, email-verification, change/reset/recovery workflow. Internal creation is unmounted, public registration remains absent, and Phase 5 provisioning/import remains unstarted. No fake links or new recovery endpoint is added.
- **Consequence:** Phase 4B foundation can be verified independently; official publication/consent and activation are not complete. Phase 5 requires its own authorization and the secure provisioning decision before dependent account delivery. See [report](PHASE4B_REPORT.md).


## ELAB-V2-DEC-068 — Institutional-email login and borrower-owned activation password

- **Status:** **Partially superseded by DEC-070**, 2026-10-08. Universal institutional email and optional Student-ID assumptions below are historical. Student requires approved SKSU email and a unique textual Student ID; Faculty needs a valid unique accessible email and no Student ID. Separate borrower-owned password/security requirements remain approved; email-link delivery remains recommended and unverified; provider selection is subsequently resolved as Brevo in DEC-071.
- **Authority/date:** Explicit owner “ACCOUNT ACTIVATION AND BULK PROVISIONING — CONFIRMED PRODUCT POLICY UPDATE”, 2026-10-08.
- **Decision:** Borrowers use institutional email and a separate eLabTrack password, never their institutional mailbox password. Preserve Phase4A authentication; no Google/Microsoft SSO or public self-registration. Borrowers set their own password during secure activation using existing hashing conventions.
- **Identity:** Stable internal account UUID; normalized unique email; separate institutional student ID where actually supplied by an approved roster. Use configurable approved institutional domains after Student/Faculty requirements are examined; the SKSU example does not select one global domain. An address-shaped string/domain/roster entry does not prove mailbox ownership.
- **Activation:** Recommend a secure single-use link delivered to institutional email. Require unpredictable tokens, hash-only storage, expiry, atomic single use, rate limits, safe invalidation/reissue, server account authority and no password exports/committed passwords/token logs. No real delivery claim before an approved provider is configured and tested. Exact operational limits/provider/domain/roster/lifecycle/recovery details remain explicit dependencies.
- **Terms:** After activation/login, current published FSMO terms must be accepted through Phase4B evidence; no invented wording. Official terms approval remains OPEN-016 and Phase4B is not fully complete.
- **Supersession:** Resolves the institutional-login/initial-password model alternatives in DEC-067/OPEN-001/009. In-person code/default-password/SSO alternatives are not the current recommendation. Does not resolve outbound provider OPEN-017, confirmed domains/roster OPEN-028 or general recovery design.
- **Implementation boundary:** Documentation update only; current code has email/password login/normalization/bcrypt but no activation/provider/domain allowlist/roster identity persistence. Reusable activation infrastructure can be independently verified only within separately authorized implementation; no Phase5/6/later authorization, migrations, mail, commit/push/deployment.
- **Handoff:** [Account provisioning policy](ACCOUNT_PROVISIONING_POLICY.md), current rules/domain/API draft/roadmap/open register. The Phase4B report remains its historical handoff snapshot.

## ELAB-V2-DEC-069 — Excel bulk borrower creation and separate bulk deactivation

- **Status:** **Partially superseded by DEC-070**, 2026-10-08. Staff individual creation, generic Borrower bulk scope and the old template below are historical. Only Admin creates accounts; bulk creation/deactivation is Student-only. Preview/confirmation/results/history/audit and safe retries remain approved. DEC-073 subsequently resolves the historical outstanding-obligation procedure below.
- **Authority/date:** Explicit owner account activation/bulk provisioning policy update, 2026-10-08; preserves DEC-052/053 role grants.
- **Creation:** Admin-only Excel upload, validation, duplicate detection, preview, explicit confirmation, durable per-record results, safe retries and audit history. File selection never imports automatically. Individual Borrower provisioning remains Staff/Admin within the narrow approved matrix. No required plaintext password column; preferred roster fields are institutional_email/institutional_id/full_name/borrower_type/program_or_course/contact_number, subject to actual school availability and required/optional review.
- **Deactivation:** Separate Admin-only Bulk Deactivate Borrowers operation may reuse the creation roster. Match stable institutional identity/normalized email to retained internal UUIDs; reject ambiguous/conflicting matches. Show matched/unmatched/already-inactive records and authoritative borrowing/fine/replacement warnings; require review and explicit selection/confirmation, then audit accurate results. Outstanding-obligation execution policy remains OPEN-029, rather than a guessed block/override.
- **Integrity:** Never hard-delete as the normal graduation process, deactivate accounts absent from the roster, select Staff/Admin through a borrower roster, silently reactivate or bypass current server role/status. Preserve all history. Idempotency/transactions protect repeated/concurrent confirmations and bind results/audit to actual selected accounts; exact batch strategy is future engineering design.
- **Supersession:** Refines existing provisioning/deactivation/import decisions; rejects V1 plaintext-password import/deletion practices as V2 graduation workflow. Borrowing/return/replacement/fine authority and external-return-photo policy remain unchanged.
- **Implementation boundary:** Phase5 remains unstarted pending separate authorization and its dependent data/provider/obligation gates. No business APIs/schema/UI are added here. [Policy and acceptance handoff](ACCOUNT_PROVISIONING_POLICY.md).

## ELAB-V2-DEC-070 — Final Student/Faculty provisioning policy

- **Status:** Accepted — APPROVED PRODUCT RULES; Student/Faculty requirements remain unchanged. DEC-071/072/073 subsequently resolve provider selection, terms-development sequencing and Student deactivation with obligations; those items in the historical dependency bullets below are superseded. Documentation only, not implemented provisioning behavior.
- **Authority/date:** Explicit owner “FINAL STUDENT/FACULTY PROVISIONING POLICY UPDATE”, 2026-10-08; supersedes conflicting portions of DEC-053/068/069 and associated UX/domain drafts.
- **Authority:** Only current active ADMIN creates any new account. STAFF assists operationally but cannot create Student, Faculty, Staff or Admin accounts. Enforce on the future backend, including direct requests and bulk confirmation/replay; client guards are UX only. Privileged accounts remain a separate restricted Admin workflow.
- **Student:** BORROWER role, borrower_type STUDENT; valid unique institutional email from approved SKSU domain configuration; required unique Student ID matching official identification. Treat IDs as strings, including leading zeroes; example `28366` is illustrative. Admin may create individually or via Student Excel import.
- **Faculty:** BORROWER role, borrower_type FACULTY; individual Admin creation only. A valid unique accessible email is sufficient, including non-institutional addresses, subject to normal activation/ownership verification. Student ID is not required. Faculty is neither STAFF nor ADMIN and is excluded from standard bulk workflows.
- **Credentials/terms:** Both categories use email plus their own separate eLabTrack password through secure activation; no SSO, mailbox passwords or Admin-assigned borrower passwords. Recommended single-use activation links target the Student institutional address or Faculty accessible address. Both must accept current officially published FSMO terms through Phase4B; no unapproved text may be published/accepted as official.
- **Student bulk creation:** Admin-only Excel; recommended columns `studentId`, `name`, `email`, `course`, `contactNumber`. Exclude password/admin role/staff role/faculty role. Server fixes BORROWER/STUDENT, validates Student ID/SKSU domain/uniqueness/conflicts, presents complete validation preview before explicit confirmation and retains durable per-record results/audit/safe retries. Reject conflicting identities rather than overwrite/coerce categories.
- **Student bulk deactivation:** Admin-only; may reuse Student roster. Stable identity matching, matched/unmatched/conflicting/already-inactive review and outstanding-obligation warnings before explicit selected-account confirmation. Recheck current role AND borrower_type to exclude Faculty/Staff/Admin, including changes after preview. Never hard-delete accounts/history or affect unselected/roster-absent accounts. Outstanding-obligation execution procedure remains OPEN-029.
- **Individual UI contract:** Admin borrower form offers STUDENT/FACULTY and conditional email/Student-ID requirements; no borrower-password assignment. Staff operational assistance retains its scope without account-creation controls. Adapt approved layout behavior later; approved mockup assets remain unchanged.
- **REMAINING TECHNICAL DEPENDENCIES:** Category/unique textual Student-ID persistence and conditional validation; reviewed API/import/matching/concurrency/idempotency/audit contracts; actual SKSU configuration and optional roster field handling (OPEN-028); activation/verification/recovery/existing-account implementation (OPEN-001/009); configured/tested provider delivery (OPEN-017). None is implemented by this update.
- **REMAINING INSTITUTIONAL APPROVALS:** Exact SKSU Student domains/official roster formatting inputs (OPEN-028), activation/provider/recovery responsibility (OPEN-001/009/017), outstanding Student deactivation procedure (OPEN-029), official terms text/version/publication approval (OPEN-016). Admin authority, Student-ID requiredness/uniqueness/string type, Faculty email/ID exemption and Student-only bulk scope are resolved, not open policy alternatives.
- **Implementation boundary/handoff:** [Current policy](ACCOUNT_PROVISIONING_POLICY.md), OPEN_DECISIONS, ROADMAP, SOURCE_OF_TRUTH, affected domain/API/UX guidance. Phase5 remains NOT STARTED pending separate authorization. No Phase4A change, migration, UI implementation, approved-mockup edit, commit, push or deployment. Historical reports/V1 evidence remain intact; unrelated operational decisions are unchanged.

## ELAB-V2-DEC-071 — Brevo transactional email provider

- **Status:** Accepted — APPROVED PRODUCT RULE; provider selected, integration/delivery unimplemented and unverified.
- **Authority/date:** Explicit owner additional policy decisions, 2026-10-08.
- **Decision:** Brevo is the selected transactional email provider for account activation and future password recovery. Implement a secure backend-only adapter during separately authorized account-management development; retain domain/application mail ports and provider-specific infrastructure boundaries.
- **Security/delivery:** API key stays in backend secret configuration, never browser code, commits, exports or logs. Live delivery requires a configured API key, verified sender and successful delivery testing. Selection, a mocked adapter or a send intent does not prove real delivery or ownership verification. Safe activation/recovery token lifecycle and truthful failure/retry handling remain technical dependencies.
- **Supersession:** Resolves email-vendor selection in OPEN-017 and earlier provider-approval alternatives. OPEN-017 retains configuration/sender/testing/operations dependencies; OPEN-001/009 retain activation/ownership/recovery lifecycle, not vendor choice. Future borrowing notification cadence remains separate OPEN-025 work.
- **Scope:** Documentation only; no integration, credential collection, API-key change, live mail, Phase5 implementation or deployment.

## ELAB-V2-DEC-072 — Official terms approval after presentation without blocking independent development

- **Status:** Accepted — APPROVED PRODUCT SEQUENCING; actual official wording remains an external institutional approval dependency.
- **Authority/date:** Explicit owner additional policy decisions, 2026-10-08.
- **Decision:** Official FSMO Terms and Conditions will be finalized after presenting the application to FSMO. Preserve existing Phase4B versioning/publication/acceptance infrastructure. Do not invent, publish or automatically accept institutional terms, or record acceptance of unapproved content.
- **Development boundary:** Missing official wording must not block independent account-management or inventory development when those phases are separately authorized. Their implementation/test readiness is distinct from official terms/publication and live borrowing readiness. This decision does not authorize Phase5/6 or change current authentication/terms UI behavior.
- **Live borrowing boundary:** Live borrowing requires officially published terms and documented acceptance, enforced by the future authoritative borrowing command and existing Phase4B evidence. Missing/unapproved terms remain fail closed for borrowing; retain read-only/account/existing-obligation access and logout under the existing rules.
- **Supersession/dependency:** OPEN-016 retains exact approved text/version/publication responsibility after presentation, not a blanket account-management/inventory development blocker. Historical Phase4B content gates remain evidence of their earlier checkpoint.

## ELAB-V2-DEC-073 — Student deactivation regardless of outstanding obligations

- **Status:** Accepted — APPROVED PRODUCT RULE; resolves OPEN-029. Account deactivation commands remain unimplemented.
- **Authority/date:** Explicit owner additional policy decisions, 2026-10-08.
- **Decision:** Admin may deactivate a Student regardless of outstanding fines, active/overdue borrowing, unreturned equipment or replacement obligations, for graduation, withdrawal, transfer or other authorized deactivation. Applies to individual Student deactivation and Admin Student-only bulk deactivation. Display appropriate warnings and require confirmation; outstanding obligations must not block deactivation or require an obligation-based override/exclusion policy.
- **Preserved authority:** Admin-only operation; bulk remains Student-only with current BORROWER/STUDENT identity/category checks, ambiguity rejection, explicit selected-account confirmation, history/audit/idempotency and no roster-absent effects. All DEC-070 Student/Faculty creation/email/ID/password/category rules remain unchanged. No new Faculty bulk or Staff deactivation authority.
- **No consequential shortcuts:** Deactivation never clears/waives fines, records payments, marks equipment returned, resolves replacements, closes unresolved borrowing, changes overdue calculations or deletes history. Preserve due times, physical/replacement balances, fine basis/clearances and all evidence; overdue computation continues under approved rules while obligations remain open.
- **Resolution after deactivation:** FSMO resolves obligations separately through authorized Staff/Admin return/replacement processing and Admin-only auditable fine clearance. Deactivated students cannot perform normal authenticated borrower actions; obligations/history remain available to authorized personnel. No auth/session policy change is implemented here.
- **Implementation handoff:** Warn/confirm/status mutation with locked current actor/target checks and atomic audit/receipt, without borrowing/stock/fine-resolution writes. Future tests must cover every outstanding-obligation case and combinations, preserved status/due/stock/overdue/financial/history evidence, same-command retries and concurrent resolution. Missing projections cannot be displayed as zero; define honest warning/error behavior in the technical contract without converting positive obligations into a veto.
- **Supersession/scope:** Removes OPEN-029 from active gates; retain its resolved audit trail. No Phase5 implementation, migration, UI/mockup edit, commit, push or deployment. [Current policy](ACCOUNT_PROVISIONING_POLICY.md) and ROADMAP separate development readiness from live delivery/borrowing prerequisites.

## ELAB-V2-DEC-074 — Authorized Batch 1 account engineering contract

- **Status:** Implemented technical contract, not new institutional policy; verification/exit is recorded in the Phase5 report.
- **Authority/date:** Explicit owner autonomous Batch1 authorization received2026-10-08; execution2026-10-09. Phase5 then independent Phase6 is authorized only after separate engineering/security/database/visual gates. Phase7 remains unauthorized.
- **Scope:** Admin-only account writes, Staff/Admin operational Borrower reads, normalized one-to-one Student/Faculty attributes on existing UUID accounts, restricted fixed-STAFF provisioning and read-only existing Admin. No creation/promotion of new Admin without a separate privileged security policy. Prior DEC-070–073 rules remain unchanged.
- **Technical safeguards:** Random32-byte hash-only, expiring, single-use account activation; atomic bcrypt/password/token/audit; trusted frontend fragment link, POST activation rate limit; cooldown/reissue/inactive invalidation. Defaults TTL24h/resend1m/preview30m are bounded engineering settings, not finalized institutional procedure. Backend Brevo adapter reports submission only; no fake provider can prove live delivery.
- **Roster contract:** Maintained bounded XLSX parser, exact five columns, text Student IDs, no formulas/privileged/credential fields; immutable owned full preview; selected valid rows all-or-nothing, current identity/version rechecks, replay-safe committed results and audit. Bulk creation leaves activation pending for explicit individual link submission. Exact existing Student ID/email matches govern deactivation even if onboarding domain configuration changes; no inferred identity or roster-absent effects.
- **History/access:** Deactivation changes account access and invalidates credentials only. Account audits/receipts are immutable; rollback refuses history. Current obligations are UNAVAILABLE, not zero, until later authoritative modules exist; positive synthetic projection tests do not prove future borrowing/fine/history integration. Existing authentication/terms/cross-tab architecture is preserved.
- **Evidence:** [API contracts](../API_CONTRACTS.md#phase-5-implemented-account-management-contracts), [plan](BATCH1_IMPLEMENTATION_PLAN.md), paired000006, real database/HTTP and browser tests. External gates remain approved Student domains, institutional ownership/recovery operations, live Brevo key/verified sender/tested delivery and official terms after presentation. No commit, push, deployment, official-content seed or V1 access is authorized.

## ELAB-V2-DEC-075 — Batch1 inventory implementation boundary

**Status:** CONFIRMED OWNER SCOPE / IMPLEMENTATION SAFEGUARDS,2026-10-09. Explicit Batch1 authorizes independent Phase6 after Phase5 core gates. DEC-049/050/062 selected stock/authority/archive policy is unchanged. Implemented000007 separates physical A/R/C/D/T, ledger/audit/idempotency and metadata version; bounded Staff/Admin catalog/category/ordinary available operations, Admin expected-sequence count reconciliation preserve custody/originals. No borrowing/return/replacement/fine settlement or Phase7 authority is inferred.

**Technical choice:** Up to512KiB canonical PNG in PostgreSQL, validated PNG/JPEG input with pixel limits and equipment-scoped authentication, provides local catalog images without a paid storage provider. Prior image/history is immutable; object storage/retention scaling remains later technical scope. No return-photo feature.

**Integration guard:** Archive blocks reservations/custody and fails closed when future borrowing/replacement tables appear without an authoritative liability adapter. Until those domains exist, verified table absence permits safe current catalog archive; this does not fabricate zero replacement liability. Later phases must wire real outstanding checks and compatible movement kinds/permissions. Repair/disposal/custody correction stays later policy. Categories are operator data, without institutional taxonomy seeds. See [API contracts](../API_CONTRACTS.md#phase-6-implemented-equipment-and-inventory-contracts) and [plan](BATCH1_IMPLEMENTATION_PLAN.md).
