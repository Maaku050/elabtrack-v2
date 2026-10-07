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
- Decision: React + TypeScript + Vite, Tailwind CSS v4, shadcn/ui and Lucide remain the eLabTrack V2 UI foundation. Borrower/customer experience is mobile-first: mobile composition first, tablet/desktop progressively enhanced. Kiosk is touch-first with large targets, privacy/session-reset awareness and a dedicated shell. Staff/admin is desktop/tablet-first, responsive on smaller screens where practical; complex administration need not be forced into tiny mobile layouts.
- Authority/date: explicit stakeholder Phase 1F UX direction, 2026-10-07.
- Impact: These are accepted design principles, not approved product screens/permissions/kiosk workflows. Existing 62 primitives and components.json remain unchanged. No mockups or feature UI created.

## ELAB-V2-DEC-039 — Mockups before feature UI completion

- Status: Accepted (future delivery process)
- Decision: Phase 3A owns UX architecture, wireframes and high-fidelity mockups: borrower mobile views first, tablet/desktop adaptations, dedicated kiosk layouts and staff/admin layouts. Phase 3B implements the approved mockups through the shadcn-based eLabTrack design system, application shell and reusable domain UI components. Later feature UI is complete only after comparison against approved mockups and responsive behavior.
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
