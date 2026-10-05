# Phase 1 security backlog

Status: findings recorded 2026-10-05; **not implemented in Phase 0**. Priorities are engineering risk assessments from current source, not proof of exploitation. The generic auth endpoints still exist even though their frontend screens were removed. Owner roles below do not invent named approvers. No actual JWT values, passwords from real accounts, or database contents are disclosed.

## Phase 1A implementation status — verified 2026-10-06

The finding tables below preserve the Phase 0 audit evidence. Current implementation and validation are in [PHASE1_FOUNDATION.md](PHASE1_FOUNDATION.md).

- **SEC-001 configuration safeguard implemented:** production rejects missing/blank/default/placeholder/short/low-diversity signing material before infrastructure. Claim verification and session design are still SEC-002–006 work.
- **SEC-009 configuration safeguard implemented:** safe URL encoding/parsing, explicit production verify-full TLS, minimum TLS 1.2, and no fallback. Unit tests inspect effective pool config; actual deployment CA/host/TLS connectivity is unrun.
- **SEC-011 development seed guard implemented:** opt-in CLI plus seeder guard deny production/test before IO. Final production bootstrap-admin procedure remains unresolved and has no seeded product-account policy.
- **SEC-008 partially implemented:** automatic Docker/startup migration removed; distinct CLI up/down/status and explicit local Compose workflow exist. SQL/version bookkeeping still use separate commits, no concurrent-run lock, and Down still conflates lookup failures with no rows. Atomicity, exclusivity, least-privilege release/runtime credentials and isolated PostgreSQL recovery tests remain open. Do not claim SEC-008 complete.
- **SEC-013 configuration portion implemented:** explicit production HTTPS frontend/origin allowlist, no wildcard/credential/path defaults. Cookie/CSRF policy and deployment/preflight verification remain open.
- **SEC-002–007, SEC-010, SEC-012, SEC-014–019 remain open.** No authentication/session redesign, product decisions, deployment or CI suite were implemented in Phase 1A.

## Phase 1B implementation status — 2026-10-06

- **SEC-006 current-account authorization implemented:** every retained protected route resolves current account role/is_active from PostgreSQL; token/body role cannot authorize it. Missing accounts get generic 401, inactive/unknown-role accounts get generic 403. In-flight snapshot windows, access-token/logout/session-family revocation and final institutional permissions remain open.
- **SEC-007 production containment implemented:** register route exists only for explicit development/test. Account creation remains reusable and local-only; no new provisioning API/UI or final eligibility/role policy is accepted.
- **SEC-005 identity locator portion implemented:** middleware rejects a zero verified uid before database access. Issuer/algorithm/audience/type/required-time/subject validation remains Phase 1C work; JWT adapter semantics are unchanged.
- **Additional discovered self-update risk contained:** old self-profile service loaded a full user and repository.Update rewrote password/role/is_active. Self update now accepts display name only and uses a narrow SQL update with current access predicate; security fields cannot be restored from a stale snapshot. No account-security migration was required.
- **SEC-017 auth-boundary coverage expanded:** real-router HTTP tests, real JWT adapter boundary test, safe repository projection/write-scope tests and frontend trusted-account adapter tests pass without a database. Live PostgreSQL concurrency/integration and CI remain unrun/open.
- SEC-002–004/010 and full JWT/session design are explicitly deferred to Phase 1C/1D. Other Phase 1A partial/open findings remain outstanding. No Phase 1C/1D implementation is claimed.

## Phase 1C implementation status — 2026-10-06

The Phase 0 finding tables and Phase 1A/1B notes are historical evidence; this section describes current source. See PHASE1_FOUNDATION.md for validation and limits.

- **SEC-003 implemented:** raw refresh tokens no longer enter persistence; SHA-256 hash-only records and hash lookup/revoke, database hash uniqueness/format constraint, and paired migration 000003. Both migration directions invalidate sessions; historical migrations remain unchanged. SQL execution on PostgreSQL is required Phase 1G coverage and was not run.
- **SEC-005 implemented:** exact HS256/signature, configured issuer, fixed API audience, required exp/iat/nbf with zero leeway, canonical nonzero uid/subject and explicit access purpose. Opaque refresh credentials cannot pass access verification. Existing current-account authorization/stale-role tests remain.
- **SEC-004 refresh rotation implemented:** one PostgreSQL READ COMMITTED transaction; FOR UPDATE session lock; FOR SHARE safe account lock; conditional once-only consumption; replacement insertion and revocation commit together. Replay fails with generic 401. Service rollback/concurrent-double, actual transaction-manager, repository statement and HTTP tests pass. **Live PostgreSQL concurrency/rollback proof remains REQUIRED Phase 1G.** Local registration's separate user/session writes remain a recovery/atomicity item for later account provisioning work. Stronger family-wide revocation/concurrent-session policy remains unresolved; SEC-004 as a whole is not declared closed.
- **SEC-006 refresh enforcement extended:** absent/inactive/unknown-role current accounts cannot refresh. Hash-based logout is idempotent, revoke-all remains internal. Access JWTs remain usable until expiry if the current account permits access; no immediate session/access or family-wide revocation policy is claimed. Revoke-all affects rows visible to its statement; concurrent/future issuance is not globally serialized.
- **SEC-018 cleanup capability implemented:** indexed terminal time, seven-day terminal-record retention and an explicit one-batch maintenance command (maximum 1000; SKIP LOCKED). Operators must schedule sufficient regular batches/monitor growth before deployment. No automatic scheduler/operational ownership is established; actual cleanup/index performance remains Phase 1G validation.
- **SEC-017 coverage expanded:** deterministic JWT, hashing, rotation/replay/rollback, test-double concurrency, actual transaction/context/isolation lifecycle, repository locking/conditional SQL/error redaction and actual-router refresh/logout contracts. No live database integration/CI claim.
- **SEC-002/010 and transport part of SEC-013 remain Phase 1D:** localStorage access/refresh storage, cookie/CSRF decision, bounded retry, hydration/cross-tab/session handling. No Phase 1D work performed. SEC-008 migration-runner integrity and the other Phase 1A/1B partial/open findings remain outstanding.

## CRITICAL

| ID | Finding / evidence | Required hardening and acceptance evidence | Owner |
|---|---|---|---|
| SEC-001 | `config/config.go` and env examples permit published development JWT secret in all environments; no production validation. A deployment retaining it would allow forged privileged claims. | Fail closed on missing/default/insufficient production signing material; validate secret configuration without logging secrets. Tests reject defaults in production and permit explicit development setup. | Backend/security |

## HIGH

| ID | Finding / evidence | Required hardening and acceptance evidence | Owner |
|---|---|---|---|
| SEC-002 | `lib/storage.ts`, auth-store and api-client persist/read **access and refresh tokens in localStorage**. XSS exposes long-lived credentials. | Establish secure transport/storage before product auth: evaluate short-lived access in memory and HttpOnly/Secure/SameSite refresh cookie with explicit CSRF/origin protection. This is a NEW V2 PROPOSAL, not a confirmed cookie contract. Test reload/expiry/logout/cross-tab behavior and absence of long-lived readable credentials. | Frontend/backend/security |
| SEC-003 | `refresh_tokens.token` and AuthRepository store raw bearer refresh tokens. | Store a one-way digest; keep raw value only for issuance/transport. Design paired auth-schema migration and test DB values cannot be used directly as bearer tokens. | Backend/data/security |
| SEC-004 | Refresh does read/check/revoke/create without a transaction; Revoke has no revoked=false predicate and can succeed repeatedly. | Atomically consume once and replace under transaction/conditional update; define family/reuse revocation and failure recovery. Concurrent requests using the same token must not both create successors. Registration user/token writes also need atomicity/recovery. | Backend/security |
| SEC-005 | JWT verification accepts HMAC family; does not require configured issuer/exact HS256/expiry presence, uid validity or subject consistency. | Pin accepted algorithm/issuer and required identity/time claims; agree audience if needed. Reject wrong issuer/algorithm/missing expiry/empty uid/inconsistent subject in tests. | Backend/security |
| SEC-006 | Auth middleware authorizes cached JWT role; no current active-status/version check on privileged routes. Logout only revokes refresh token. | Define current-account authorization, role change/suspension/revocation semantics before product roles. Test stale privilege and inactive access denial per agreed policy. Do not invent staff/admin/Super Admin authority (OPEN-003/004/005/008). | Backend/security + FSMO policy owner |
| SEC-007 | Public `/auth/register` is live starter behavior; generic user/admin roles conflict with unresolved provisioning/eligibility policy. | Contain unintended public onboarding while policy is unresolved, then implement only confirmed Phase 4 flows. Do not treat hidden frontend signup as server protection. Test unauthorized/privileged account creation and deny role escalation. Policy is OPEN-001/002/003, not accepted public registration. | Backend/security + FSMO administrator |
| SEC-008 | Docker CMD automatically runs migrations; runner commits DDL separately from its version record and has no concurrent-run lock. | Define explicit authorized release migrations, least-privilege runtime credentials, atomic SQL/bookkeeping and exclusive execution. Prove failed migration/duplicate-run recovery against isolated PostgreSQL. Do not run these checks on V1. | Backend/DevOps/data |
| SEC-009 | `config/database.go` forces sslmode=disable; URL assembly incompletely escapes credentials. | Use validated connection configuration/encoding and explicit environment-specific TLS; production must require approved transport security. Test special characters and production rejection of unsafe configuration. | Backend/DevOps |
| SEC-010 | `api-client.ts` retries a 401 after refresh without a retry marker; another 401 can reenter refresh repeatedly. | Bound one request retry and clear session safely on failure; synchronize actual credentials/session metadata and consider cross-tab rotation. Test persistent 401, concurrent refresh and revoked session paths. | Frontend/security |
| SEC-011 | Dev seeder has no enforced production guard; SQL creates a privileged synthetic account with known password. | Reject production seed execution, separate bootstrap admin procedure and enforce environment isolation. Test guard without executing real seeds. No seeded account policy is approved for V2. | Backend/security/DevOps |

## MEDIUM

| ID | Finding / evidence | Required hardening and acceptance evidence | Owner |
|---|---|---|---|
| SEC-012 | Global limiter uses c.IP(); bootstrap has no trusted-proxy handling; nginx forwards X-Forwarded-For. Default 120/min shared across routes is not a login abuse policy. | Verify actual proxy/IP behavior, accept forwarding only from configured trusted proxies, and add focused auth abuse controls without introducing Redis. Test forged forwarded headers and client separation behind intended proxy. | Backend/DevOps/security |
| SEC-013 | CORS allows configured origins with credentials; production falls back to localhost allowlist, with no environment validation. | Require explicit production origins/credential policy and reject wildcard/invalid settings; test preflight denial and accepted origins. Coordinate with any proposed cookie/CSRF transport. | Backend/security |
| SEC-014 | API Helmet headers do not apply to nginx-served SPA; current defaults may not fit actual TLS/kiosk/embed assumptions. | Validate deployment headers/CSP/TLS separately for API and SPA; retain accessible app behavior and deny unwanted framing. Test actual serving path. | Frontend/DevOps/security |
| SEC-015 | Request logger is plaintext regardless of JSON app-log setting; recovery callback configured with stack trace disabled; domain errors sent by handlers may bypass central logging. | Establish consistent structured events/request IDs and redaction; verify error/panic observability without exposing payloads/tokens/PII. | Backend/operations |
| SEC-016 | Normal envelope differs from raw health and frontend error typings; production error paths need verified safe contracts. | Align API/frontend error types and health contract; test malformed inputs, unknown routes, panic paths and SQL failures for safe responses. Preserve strict type checking. | Backend/frontend/QA |
| SEC-017 | Unit/HTTP tests exist but no real auth repository/concurrency integration suite or CI workflow. | Add risk-based tests for the findings above using isolated PostgreSQL and synthetic identities; required vet/lint/build/tests run reproducibly in CI. Resolve Go toolchain mismatch first. | QA/backend/DevOps |

## LOW

| ID | Finding / evidence | Required hardening and acceptance evidence | Owner |
|---|---|---|---|
| SEC-018 | No refresh-token retention/cleanup schedule or operations policy in current bootstrap. | Define bounded retention and revocation cleanup with operational ownership; design before scheduled implementation and measure query/index needs. | Backend/operations |
| SEC-019 | Retained UI/hooks emit 19 lint warnings; dependency lockfile/image policies need repeatable review. | Triage warnings without deleting reusable primitives or weakening gates; review dependencies/images using current primary evidence when upgrades are proposed. No package/image security audit was claimed during Phase 0. | Frontend/DevOps |

## Phase boundary and completion

Phase 1 can harden generic foundation mechanisms without implementing borrower administration, equipment, fines, kiosk or notifications. Product onboarding/privilege policy decisions belong to stakeholders and Phase 4. Cookie transport and token family details are proposals until their contracts are reviewed. Business history/audit and stock integrity are confirmed design constraints for Phase 2 and later; no financial or inventory tables are created by this backlog.

Before declaring foundation ready, resolve critical findings and verify high-risk mechanisms with tests, track remaining items with owners, and record all unavailable environment/database/deployment evidence. `PHASE0_REPORT.md` records this session's actual validation; this backlog is not completion evidence.
