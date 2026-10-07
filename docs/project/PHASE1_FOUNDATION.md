# Phase 1A — Production configuration and runtime safety

Task started 2026-10-05; closure verified 2026-10-06 (Asia/Shanghai). Phase 0 remains COMPLETE. Phase 1A is infrastructure foundation work; authentication/session design is **not complete** and this is not deployment approval. Phase 1B has not begun. No product workflow, table, policy resolution, dependency upgrade, commit, push or deployment is included.

## 1. Files changed

The exact file list and working-tree state are recorded in section 16. Changes cover centralized configuration/tests, startup/CLI wiring, guarded seeding, read-only migration status, Docker/Compose, Make commands, environment examples/ignore rules and current engineering documentation. Deleted fallback-only `config/env.go` and `bootstrap/helpers.go` are replaced by strict parsing. The existing HTTP contract fixture changes only its now-typed body limit. Frontend implementation and shadcn primitives are untouched; only the public environment example comment changes.

## 2. Configuration architecture before and after

Before: `Load()` returned config without validation; helpers read the process environment repeatedly and silently replaced malformed values with defaults. Dotenv mutated global environment. Environments were arbitrary strings. Body limits silently fell back during server wiring. Database DSN concatenation partly escaped passwords and unconditionally disabled TLS. CLI commands bootstrapped the HTTP application, and infrastructure failure called a fatal process exit.

After: `config.Load() (*Config, error)` snapshots process environment once, merges missing values from a local `.env` only for development, and calls `Parse(map[string]string)`. Explicit production/test/unknown environments never read `.env`. Actual environment, including explicitly blank values, takes precedence. Missing optional settings receive documented defaults; invalid/blank critical values return errors rather than a usable config. No global environment is mutated. Dotenv supports simple `KEY=value`, `export`, matching quotes and full-line comments; no interpolation, multiline values or inline comment expansion. Malformed lines/read failures are errors without echoing contents.

`Environment`, integer body/pool/rate limits and durations are typed. Bootstrap validates before logger, pool, service wiring or listener. Each server/one-shot command loads once. Validated config is passed to infrastructure. One-shot commands do not create HTTP handlers/listeners; migration scaffold creation does not connect to PostgreSQL. The app no longer constructs/exposes a migrator or seeder. Connection/config errors are safe descriptions with setting names and no supplied secret/URL values; infrastructure returns errors instead of exiting internally. No heavy dependency was added. Existing pgx initializes pool internals; effective credentials, runtime parameters and TLS are pinned to the validated snapshot.

## 3. Supported environments and complete variable inventory

Only exact `development`, `test`, `production` values are accepted. Missing `APP_ENV` selects local development for compatibility; blank/unknown names fail. Production must explicitly set `APP_ENV=production`. Test requires an explicit isolated signing key, ignores dotenv and rejects development seeds; use a dedicated test database for integration tests. Configuration unit tests use synthetic maps and do not connect.

“Optional” means omission may use the documented fallback; supplying an invalid value is an error. “URL alternative” means DATABASE_URL replaces all five discrete connection identity settings, which must then be absent, even when blank. Production password-based TCP authentication is the supported baseline. Example hosts/values below are illustrative only.

| Backend setting | Purpose; required/optional | Development default | Production requirement | Sensitivity; example format |
|---|---|---|---|---|
| APP_ENV | Typed runtime environment; optional locally | development | Explicit production | Public; `development`, `test`, `production` |
| APP_PORT | API TCP port; optional | PORT, then 8080 | Decimal 1–65535 | Public; `8080` |
| PORT | Platform alias; optional | 8080 if both absent | Same as APP_PORT if both supplied | Public; `8080` |
| APP_NAME | API identity; optional | eLabTrack V2 | Nonblank; same default | Public; `eLabTrack V2` |
| APP_READ_TIMEOUT | HTTP read timeout; optional | 10s | Positive Go duration | Operational; `10s` |
| APP_WRITE_TIMEOUT | HTTP write timeout; optional | 15s | Positive Go duration | Operational; `15s` |
| APP_BODY_LIMIT | HTTP payload ceiling; optional | 1MB | Positive bytes/KB/MB/GB, maximum 1GB | Operational; `512KB` |
| DATABASE_URL | PostgreSQL URL alternative; optional | Absent | Explicit user/password/host/database and verify-full strategy when used | **Secret-bearing**; `postgresql://USER:ENCODED_PASSWORD@HOST:5432/DATABASE?sslmode=verify-full` |
| DB_HOST | Single TCP host; optional locally, URL alternative | 127.0.0.1 | Explicit nonblank host | Infrastructure; `db.example.invalid` |
| DB_PORT | TCP database port; optional, URL alternative | 5432 | Decimal 1–65535; same default | Infrastructure; `5432` |
| DB_NAME | Database identity; optional locally, URL alternative | elabtrack_v2 | Explicit nonblank database | Infrastructure; `elabtrack_v2` |
| DB_USER | Database login; optional locally, URL alternative | postgres | Explicit nonblank user | Credential metadata; `runtime_user` |
| DB_PASSWORD | Database password; optional locally, URL alternative | Empty in code; example uses published local-only value | Explicit nonblank password | **Secret**; inject externally, no production value supplied |
| DB_SSLMODE | Transport policy; optional locally | disable | Explicit verify-full here or URL query | Operational; `disable` locally / `verify-full` |
| DB_SSLROOTCERT | Trust source; optional | System trust for verify-full | Omitted/`system` or readable, valid existing PEM CA bundle; incompatible with disable | Infrastructure path; `system` or `/path/to/existing-ca.pem` |
| DB_MAX_CONNS | Pool maximum; optional | 20 | Integer 1–2147483647; capacity must fit deployment | Operational; `20` |
| DB_MIN_CONNS | Pool minimum; optional | 2 | Integer 0–DB_MAX_CONNS | Operational; `2` |
| DB_MAX_CONN_LIFETIME | Pool connection lifetime; optional | 1h | Positive Go duration | Operational; `1h` |
| DB_MAX_CONN_IDLE_TIME | Idle connection lifetime; optional | 15m | Positive Go duration | Operational; `15m` |
| JWT_SECRET | Existing HMAC signing material | Published local-only key from code/example | Explicit key satisfying section 4; no fallback | **Secret**; externally generated random material, never a documentation value |
| JWT_ISSUER | Existing token issuance identifier; optional | elabtrack-v2 | Nonblank; same default | Public; `elabtrack-v2` |
| JWT_ACCESS_TTL | Existing access TTL; optional | 15m | Positive Go duration | Operational; `15m` |
| JWT_REFRESH_TTL | Existing refresh TTL; optional | 168h | Positive and greater than access TTL | Operational; `168h` |
| FRONTEND_URL | Intended frontend origin; optional locally | http://localhost:5173 | Explicit exact HTTPS origin | Public; `https://app.example.invalid` |
| ALLOWED_ORIGINS | Credentialed CORS allowlist; optional locally | http://localhost:5173,http://localhost:4173 | Explicit comma-separated HTTPS origins; no wildcard/user info/path/query/fragment/empty item | Public; `https://app.example.invalid` |
| RATE_LIMIT_MAX | Existing global rate ceiling; optional | 120 | Integer 1–2147483647 | Operational; `120` |
| RATE_LIMIT_WINDOW | Existing global rate window; optional | 1m | Positive Go duration | Operational; `1m` |
| LOG_LEVEL | App logger verbosity; optional | info | debug/info/warn/error | Operational; `info` |
| LOG_FORMAT | App logger encoder; optional | console | json default; json/console accepted | Operational; `json` |

No seed or migration environment toggle is supported. `SEED_ENABLED`, `AUTO_MIGRATE`, `MIGRATE_ON_START` have no effect; use the guarded CLI actions below. Nonempty `PG*` environment settings are rejected by name to prevent pgx's implicit credential/service/TLS configuration from becoming a second configuration path. Password/service files, keyword DSNs, multi-host URLs, Unix sockets and client-certificate/passwordless authentication are outside this baseline. URL query allowlist is only one `sslmode` and/or one `sslrootcert`; duplicate/unknown keys and conflicting DB_SSL* settings fail. Pool sizing remains in DB_* even with a URL. General unrelated environment variables are ignored.

| Frontend/Compose setting | Purpose; required/optional | Local default | Production requirement | Sensitivity; example |
|---|---|---|---|---|
| VITE_API_URL | Public browser API base; optional; bundled at build time | http://localhost:8080/api/v1 in source/standalone Docker; `/api/v1` in Compose | Supply intended public API path/origin at build; no backend secrets | Public; `/api/v1` |
| Root DB_USER, DB_PASSWORD, DB_NAME | Local Compose PostgreSQL initialization/backend credentials; optional | postgres / published local-only password / elabtrack_v2 | This Compose file is **development only**, not production configuration | DB_PASSWORD sensitive outside local examples; match backend host-run credentials |
| Root DB_PORT | Local Compose host published PostgreSQL port; optional | 5432 | Development only | Public; `5432`; container always uses internal 5432 |

Vite's built-in MODE/DEV/PROD/BASE_URL are build metadata, not server secrets. No JWT/database/mail credential has a VITE_* counterpart. FRONTEND_URL is retained configuration metadata; current server CORS reads ALLOWED_ORIGINS, and no new redirect or auth contract is introduced.

## 4. Production JWT secret policy

Production has no fallback: missing/blank/whitespace, the former template secret, the current published development key, recognized placeholders, fewer than 32 bytes, fewer than 12 distinct characters, whitespace/control characters, common alphabet/digit/keyboard runs and repeated patterns of up to 32 bytes fail startup. Externally generate at least 32 **random bytes**, encoded as hex or base64, and inject them using the deployment's chosen secret mechanism. Application startup never generates/persists a key and never logs it. The screening is a heuristic for obvious weakness, not proof of cryptographic entropy or key ownership. A passing hand-written value is not a recommendation.

This changes signing configuration only. Existing JWT/session issuance, claims verification, refresh storage/rotation and browser storage remain unchanged and provisional.

## 5. PostgreSQL and TLS policy

Production must explicitly choose **verify-full** in DB_SSLMODE or URL query. The effective pgx connection requires TLS 1.2+, validates certificate chain and hostname, has InsecureSkipVerify=false and no plaintext/fallback connection entries. System trust is the default; an existing readable PEM CA file can be supplied. Invalid/missing supplied bundles fail validation without printing paths/content. No certificate, hostname/provider, proxy exception or alternative transport security has been invented. Disable is allowed only in development/test; allow/prefer/require/verify-ca and missing production strategy are rejected.

No actual production database or certificate handshake was exercised. Unit tests inspect effective TLS settings, CA error paths and credential round trips; deployment CA trust, hostname matching, routing and authentication need later authorized environment checks.

## 6. Credential and URL handling

Discrete user/password/database fields are encoded with net/url.UserPassword and URL path/query encoding; host/port uses net.JoinHostPort. Passwords/usernames containing `%`, `:`, `@`, `/`, `#`, `&`, `?`, quotes, spaces, plus signs or Unicode survive round trips. DATABASE_URL is parsed as a PostgreSQL URL and validated with a narrow TLS query allowlist. It replaces discrete identity fields; ambiguity, malformed escapes, NUL and unsafe transport fail without echoing URLs/credentials. pgx parses the encoded connection settings, and the resulting pool config explicitly pins credentials/runtime parameters and validated TLS. The full config/pool/connection string must never be logged. Connection failures return generic connectivity/credential/TLS guidance rather than raw pgx errors containing server/connection details.

## 7. Development seeding policy

Normal startup has no seed path. Only explicit `--seed`/`make seed` in development is permitted. CLI guard runs before connecting; Seeder.Run independently checks environment and explicit request before reading files/using its pool. Production, test, unknown environment and a non-requested call fail safely. Tests use a nil pool and nonexistent directory to prove rejection before IO. Synthetic SQL/accounts remain unchanged and local-only; no final eLabTrack administrator/user/role policy is accepted. Logs identify applied SQL filenames, never passwords or SQL/server error contents.

## 8. Migration policy and commands

| Action | Host command from backend/ | Behavior |
|---|---|---|
| API | `go run ./cmd/api` | Validate, connect, serve; no migration or seed |
| Apply | `go run ./cmd/api --migrate-up` | Explicit pending paired SQL migrations |
| Roll back | `go run ./cmd/api --migrate-down` | Explicit destructive latest rollback; operator responsibility |
| Inspect | `go run ./cmd/api --migrate-status` | Read-only pending/applied state; never creates tracking table |
| Scaffold | `go run ./cmd/api --migrate-create NAME` | Create paired SQL files, no database connection |
| Seed | `go run ./cmd/api --seed` | Explicit development-only seed command |

Multiple actions, invalid/empty scaffold names and positional arguments fail. Root/backend Makefiles expose the corresponding commands. Application startup and container restart never imply schema mutation. Production migrations must be a separate authorized release step/job, with no cloud-specific tooling introduced. All existing migration SQL is preserved. Status reflects a snapshot of current bookkeeping and packaged files, not schema-integrity proof.

**SEC-008 remains partial:** no lock or atomic SQL/version bookkeeping redesign was attempted. Down's inherited lookup-error handling also remains open. Operators must serialize releases until those mechanisms and recovery tests are hardened; this policy is not a claim of concurrency safety or a least-privilege deployment.

## 9. Docker and Compose

Backend Docker uses exec ENTRYPOINT `["./api"]`; removed the shell command that ran migrations first. An explicit container invocation can pass a one-shot flag. Compose retains frontend, backend, PostgreSQL, health dependency and the same named volume. Backend env_file is optional, with local development defaults, port 8080/internal DB port 5432 and local disabled TLS. APP_ENV from the backend file is preserved, so unknown values fail and production settings cannot silently turn into development; production is rejected with this local disabled-TLS configuration. It shares the PostgreSQL initialization credentials; the development-only examples match. `make compose-migrate-up` runs `docker compose --profile full run --rm --build backend --migrate-up`, then the full profile can start normally. Local .env overrides do not turn this file into a production manifest.

Existing volumes keep their initialized credentials; editing env files does not rotate passwords or migrate their data. No volume was deleted, recreated or accessed. Container images/services were not built/started in this verification.

## 10. Environment file hygiene

`.env` and all `.env.*` are ignored at every repository depth; only `.env.example`/`.env.*.example` are exceptions. Backend Docker context also excludes real env files. Root example now contains only local Compose inputs; backend example covers local API settings; backend/.env.production.example intentionally leaves production-required credentials/hosts/origins blank and fails until supplied. Production/test read injected environment, not local dotenv. Frontend example documents public build values only. There were no real env files in the inspected workspace and no real secret was introduced/staged/committed. Only published local values and synthetic unit fixtures appear in changes. Production-required values are inventoried above.

## 11. Tests added

- config: valid development/test/production; invalid/blank environment; strict ports/body/durations/pool/rate/log settings; PORT precedence/conflicts; test signing isolation; missing/default/blank/placeholder/short/repeated/sequential/low-diversity production secrets; strong fixture accepted.
- config: missing/unsafe production TLS/password/origins; local disabled TLS accepted; missing CA bundle; safe discrete/URL credential round trips and effective host-verified TLS/no fallback; malformed/duplicate/unsupported/conflicting URL settings and redacted errors; ambient PostgreSQL rejection/pinning; dotenv precedence, malformed input redaction and test/production isolation.
- seeder: production/test/unknown/non-requested guard before any filesystem/database IO, explicit development guard accepted without executing seeds.
- CLI: exclusive command selection and safe invalid-argument errors; startup validation before connection; production seed rejection before connection using canceled context/unreachable synthetic host.

No production database is needed. Existing backend domain/security/HTTP and frontend tests remain enabled. Actual migration/status success and schema immutability on a live database were not integration-tested; the startup/status command paths were inspected for absence of DDL calls.

## 12. Backend validation

`go fmt ./...`, `go vet ./...`, `go test ./...` PASS. Final structured test run: 20 top-level tests plus 55 subtests (75 passing test nodes), across 7 tested packages. Phase 1A adds 11 top-level tests; 9 existing top-level tests remain. Go 1.27.1 is selected inside backend/; the launcher outside the module remains 1.26.5. Dependencies/go.mod/go.sum are unchanged. Gates used `GOCACHE=/tmp/elabtrack-phase1a-go-cache` because the default home cache is not writable in this sandbox. No database or Docker socket is required for these tests.


## 13. Frontend validation

`npm run lint` PASS with the same 19 inherited shadcn/hook warnings and no errors; frontend implementation/lint rules are unchanged. `npm run test:run` PASS: 21 tests, 3 files. `npm run build` PASS: TypeScript and Vite, JS 475.69 kB (149.97 kB gzip), CSS 179.91 kB. No lint/test/type-check gate was weakened.

## 14. Compose validation

`docker compose config --quiet` and `docker compose --profile full config --quiet` PASS with all services/named volume retained and no real backend .env required. This verifies configuration resolution only, not image builds, startup, health or database execution.

## 15. Diff, security review and scope

`git diff --check` PASS. Review confirmed:

| Boundary | Evidence / result |
|---|---|
| Published production JWT key | Missing/default/development/blank/weak material fails before infrastructure; CLI startup test uses canceled context |
| Redacted validation/startup errors | Tests use synthetic secret sentinels, malformed URL input and invalid command arguments; no supplied values emitted; DB connection errors are generic |
| Production transport | Effective pool TLS verifies host/chain, minimum TLS 1.2, no plaintext fallback; unsafe or undefined modes fail |
| Credential encoding | Reserved characters and Unicode round-trip in discrete/URL credentials; no manual escape helper remains |
| Seeder isolation | CLI and seeder require development and explicit request before IO; nil-pool/missing-directory guard tests pass |
| Startup migration policy | HTTP infrastructure contains no migrator/seeder; Docker ENTRYPOINT serves only API; migration actions are explicit |
| Browser config | Only public VITE_API_URL is supported; all 102 frontend source files and UI primitives are unchanged |
| Environment/secret hygiene | Real .env variants ignored, example exceptions retained, production example blank, no real env file/secret added and nothing staged/committed |
| Phase boundary | All 4 auth migration SQL files, 1 seed SQL file, 10 domain files, 10 application files, 3 security adapters, 2 repositories and 16 V1 audit files match baseline hashes |

No equipment, inventory, borrowing, approval, return, fine, notification, reporting, kiosk or campus-wide feature/table was created. Existing users/refresh_tokens SQL and schema_migrations bookkeeping are the only schema mechanisms. PHASE0_REPORT.md, OPEN_DECISIONS.md and dependency manifests/lockfiles are unchanged. Technical accepted decisions are limited to DEC-021–023; no stakeholder product policy was resolved.


## 16. Exact git status

All 32 changes are unstaged: 24 modified tracked files, 2 deleted tracked files and 6 new files. Exact `git status --short` at closure:

```text
 M .env.example
 M .gitignore
 M Makefile
 M README.md
 M backend/.dockerignore
 M backend/.env.example
 M backend/Dockerfile
 M backend/Makefile
 M backend/README.md
 M backend/cmd/api/main.go
 M backend/internal/bootstrap/app.go
 D backend/internal/bootstrap/helpers.go
 M backend/internal/bootstrap/infrastructure.go
 M backend/internal/bootstrap/migration_test.go
 M backend/internal/bootstrap/server.go
 M backend/internal/config/config.go
 M backend/internal/config/database.go
 M backend/internal/config/dotenv.go
 D backend/internal/config/env.go
 M backend/internal/infrastructure/database/migrator.go
 M backend/internal/infrastructure/database/postgres.go
 M backend/internal/infrastructure/database/seeder.go
 M docker-compose.yml
 M docs/project/DECISIONS.md
 M docs/project/PHASE1_SECURITY_BACKLOG.md
 M frontend/.env.example
?? backend/.env.production.example
?? backend/cmd/api/main_test.go
?? backend/internal/config/config_test.go
?? backend/internal/config/seed.go
?? backend/internal/infrastructure/database/seeder_test.go
?? docs/project/PHASE1_FOUNDATION.md
```

No staging/commit/push is authorized or performed.

## 17. Remaining Phase 1 work and evidence limits

SEC-001 and SEC-009 configuration safeguards and SEC-011 seed guard are implemented. SEC-008 (startup policy only) and SEC-013 (configuration validation only) remain partial. SEC-002–007, SEC-010, SEC-012, SEC-014–019 remain open: browser/refresh-token storage, transactional single-use rotation, exact JWT claim/algorithm/issuer verification, current-account authorization/revocation, public onboarding containment, bounded refresh retries, trusted proxy/auth abuse limits, deployment headers/TLS/CSP, safe structured observability/envelopes, isolated repository/concurrency/CI coverage, retention, dependency/image review and inherited lint warnings. See the preserved finding tables and current status in PHASE1_SECURITY_BACKLOG.md.

Production CA/connectivity, image builds, container runtime, live development API success, migration up/down/status execution, schema permissions/locks/recovery and deployment checks were not run. No V1 services/data were accessed. No configured Obsidian MCP is exposed in this session, so private-note reads/version-checked persistence could not be performed. The current repository AGENTS.md does not require that service; removed `.specify/memory/constitution.md` and `docs/SPECKIT-PLANNER.md` remain historical Phase 0 removals, not recreated dependencies.

## 18. Phase 1A exit gate

**SATISFIED** for the authorized Phase 1A implementation and requested local gates. Backend fmt/vet/tests, frontend lint/tests/build, both Compose configuration checks and diff hygiene pass. Production-secret, explicit TLS, credential encoding, seed isolation and explicit migration policies have automated/source evidence. Unrun live database/container/deployment checks are listed in section 17. Phase 1A exit concerns implemented configuration/runtime safeguards and requested local validation, not authentication completion or production readiness. Open product decisions and Phase 0 historical reports are preserved. Do not begin Phase 1B without a separate instruction.

# Phase 1B — Authentication & Authorization Foundation

Verified 2026-10-06 (Asia/Shanghai). Phase 0 and Phase 1A remain COMPLETE; their reports above are historical closure evidence and are preserved. Phase 1B began from a clean working tree. This section records the current technical boundary, **not a final institutional account/role policy or completed session-security design**. No Phase 1C/1D, business feature, production infrastructure access, commit, push or deployment is included.

### B1. Files changed

The exact Phase 1B working-tree list is in B18. Changes are limited to safe account projection/read port, application resolver/DTOs, narrow self-profile write, middleware/context helpers, route/bootstrap wiring, focused backend/frontend tests, frontend auth/profile adapters, and current engineering docs. Removed legacy claims/string-local accessors and the unused context-ID path; downstream handlers use the typed current-principal helper. The old CurrentUser method was replaced by the already-resolved principal response. refresh.go only loses an unused Claims import/placeholder after that removal; **refresh behavior is unchanged**.

Unchanged: all existing SQL migrations/seeds, JWT/bcrypt/opaque-token adapters, raw refresh repository, browser token transport/storage/store implementation, placeholder router/providers/pages, 62 shadcn primitives, Phase 1A config/migrator/seeder safeguards, dependencies/lockfiles and V1 references. No migration is necessary: existing `users.role` and `users.is_active` already represent the generic boundary.

### B2. Previous auth architecture and route audit

Before Phase 1B: bearer middleware verified a token and copied the entire application Claims plus user ID to string-key Fiber locals. RequireRole read token role only. Auth/me and users/me independently looked up a full User, without an access-boundary active/known-role check. Profile updates bound name/email but loaded a full aggregate and rewrote password/role/is_active; an in-flight self update could restore security fields changed since its read. Registration was public in all environments. Login checked inactive state before password verification. Refresh/logout used raw opaque refresh credentials in request bodies and remain separate from access-token authentication.

| Current route | Authentication required before / retained transport | Previous authorization | Previous trust source | Risk | Phase 1B action |
|---|---|---|---|---|---|
| GET /api/v1/health | None | Public diagnostic | Health adapter | Known raw-envelope exception | Keep unchanged; no account lookup |
| POST /api/v1/auth/register | None | Unrestricted starter signup | DTO + account-creation service | Unresolved onboarding/eligibility; unreviewed production surface | Production absent; explicit development/test only; strict allowed fields |
| POST /api/v1/auth/login | Email/password, no access bearer required | Existing active-account login | PostgreSQL user + bcrypt | Inactive-account existence revealed before password check | Verify password before status; preserve hashing/token mechanics |
| POST /api/v1/auth/refresh | Opaque refresh credential, no access bearer required | Existing token validity + active account | Raw refresh row + PostgreSQL user | Raw storage/non-atomic rotation/issuer validation | Keep flow unchanged; defer to 1C/1D |
| POST /api/v1/auth/logout | Opaque refresh credential, no access bearer required | Holder can revoke supplied token | Raw refresh row | Access token remains usable; revocation redesign pending | Keep existing revocation flow; defer redesign |
| GET /api/v1/auth/me | Bearer access token | Token verification; subsequent user lookup | Token uid + full User lookup | Missing current active/known-role boundary | Auth + safe current-account snapshot; five safe fields |
| GET /api/v1/users/me | Bearer access token | Token verification only | String-local token uid + full User lookup | Inactive/unknown account could still read | Same boundary; safe profile DTO from principal |
| PATCH /api/v1/users/me | Bearer access token | Token-identified self; name/email DTO | Token uid + previously loaded full User | Identity editing policy unresolved; stale full-row security overwrite | Same boundary; name-only strict DTO and narrow SQL write |
| GET /api/v1/users/ (also /users) | Bearer access token | Cached JWT admin role | Token role | Demoted/inactive token holder could list users | Current active PostgreSQL admin only; bounded read-only list |

There is no arbitrary-user GET/PATCH/DELETE, account-create administration route, role/status mutation endpoint or alternate production public-signup route. Reusable domain/repository account capabilities are not exposed APIs. Existing generic list DTOs never return password hashes.

### B3. New trusted authentication boundary

Bootstrap composes one `middleware.Auth(issuer, accounts)` implementation for all retained protected routes. It reads the existing Authorization Bearer transport, invokes the existing TokenIssuer verifier, rejects invalid/absent credentials and zero account ID, constructs typed `auth.Identity{UserID}`, then calls the application CurrentAccountService. Token email/role are discarded for authorization. No request body/path/frontend identity is used as proof.

CurrentAccountService depends on the narrow domain AccountRepository port. PostgreSQL implements a parameterized `SELECT id,email,name,role,is_active,created_at,updated_at FROM users WHERE id=$1`. It does not select passwords/refresh/security records. The safe domain Account projection becomes a value Principal; it cannot carry password/session fields. Identity/principal are stored using private typed Fiber keys only after resolution succeeds. `PrincipalFromContext` and `IdentityFromContext` are the only HTTP context accessors; handlers no longer cast arbitrary locals or read cached token Claims. Application/domain have no Fiber/pgx/infrastructure dependency.

### B4. Registration containment

POST /api/v1/auth/register is mounted only when the route environment is development or test. Production, unknown or unset route environments return 404 with the existing error envelope; creation/refresh-record persistence is not invoked. Runtime routing receives the validated Phase 1A config: omitted APP_ENV still selects development locally, while production must explicitly set APP_ENV=production. No registration-enable flag or production bypass is introduced.

Development/test retain the existing local convenience, including immediate generic-user token issuance. Strict JSON accepts only email/name/password and rejects role/status/security-field injection with generic 400. Bcrypt persists the hash, not plaintext. The reusable account-creation use case/domain/repositories remain. No administrative provisioning UI/API, approval/verification flow or final borrower role is invented. OPEN-001/002/003 remain unresolved; this is reversible technical containment, not an accepted public-signup decision.

### B5. Current-account lookup

Every protected request resolves current PostgreSQL account state once before its handler/role check. It is not cached in JWTs, sessions or process memory. Me handlers reuse the safe snapshot rather than issuing a second aggregate/password read. Missing rows (including wrapped not-found errors), nil results or mismatched IDs fail with generic 401; database failure fails closed with safe 500 and no SQL/detail disclosure. A missing account does not become a resource-detail response that confirms account existence.

The profile update separately rechecks account permission in its SQL statement, and the privileged list performs its existing bounded read after authorization. Lookup/DTO/SQL behavior is tested through narrow fakes and the actual repository transaction seam; a real PostgreSQL integration run was not performed.

### B6. Account status and error handling

Only the inherited boolean `is_active` is used. Existing schema defines NOT NULL DEFAULT TRUE; no disabled/suspended vocabulary, disciplinary rule, loan restriction, migration or new school policy is added.

| Condition | Result |
|---|---|
| Missing/malformed bearer, failed verifier or zero identity | 401; no account query for invalid identity |
| Verified identity has missing/deleted/nil/mismatched account | Generic 401 |
| Current account has is_active=false or unknown role | Generic 403 |
| Active known-role account lacks required current role | Generic 403 |
| Account lookup unavailable | Safe 500; handler does not run |
| Arbitrary-account route genuinely absent | 404, irrespective of claimed privilege |
| Self mutation includes identity/security/unknown fields or malformed JSON | Generic 400; no mutation |
| Missing/invalid required display-name input | Existing 422 DTO validation / 400 domain validation |
| Account disappears/becomes inactive/unknown before name UPDATE matches | Generic 403; no update |

Login with bad credentials returns 401 without exposing inactive status; a verified password for an inactive account retains the existing 403. Protected-request error envelopes are consistent with existing API helpers. Health's raw envelope, richer global error/observability work and login timing/abuse design are deferred; this is not Phase 1E/1F error redesign.

### B7. Authorization source of truth

Privileged checks read Principal.Role from the **current PostgreSQL row**, after validating current activity/known role. JWT role may still be issued for transport compatibility but cannot independently authorize an API operation. Stale admin token + current user role yields 403 on list; old user token + current admin role can list. Request-body role/identity cannot override either. Subsequent protected requests observe demotion/deactivation without waiting for token expiry.

This is one account snapshot per request, not instantaneous global revocation. If role/status changes after a request has already passed its check, that in-flight read may complete from the earlier authorized snapshot. The display-name mutation writes no security fields and rechecks access in its UPDATE. Future consequential privileged mutations need their own transaction/concurrency authorization design. No session-family/version mechanism, access-token revocation, locking or concurrent-session policy is accepted here.

### B8. Temporary role model and technical route matrix

Only inherited generic `user` and `admin` are recognized. Unknown current or required role values cannot grant access; a role guard without the authentication boundary returns 401. These are temporary infrastructure values matching the existing CHECK constraint, **not final borrower/staff/administrator permissions**. No Super Administrator, organization scope, department or campus role is added.

| Route | Category | Authentication | Authorization | Trusted identity/state source |
|---|---|---|---|---|
| GET /health | Public | None | Public diagnostics | Health adapter |
| POST /auth/login | Credential entry | Email/password | Existing active-account credential check | PostgreSQL + bcrypt |
| POST /auth/register | Local-only account capability | None | Route mounted only development/test | Validated environment + fixed generic-user creation |
| POST /auth/refresh | Refresh credential | Opaque token body | Existing validity/active-account check | Raw refresh repository + PostgreSQL; deferred redesign |
| POST /auth/logout | Revoke supplied credential | Opaque token body | Existing token possession/revocation | Raw refresh repository; deferred redesign |
| GET /auth/me | Authenticated self | Shared bearer/current-account boundary | Active known-role self only | Current PostgreSQL Principal |
| GET /users/me | Authenticated self | Same boundary | Active known-role self only | Current PostgreSQL Principal |
| PATCH /users/me | Authenticated self | Same boundary | Display name of principal ID only; SQL recheck | Current principal + narrow PostgreSQL UPDATE |
| GET /users[/] | Temporary generic privileged read | Same boundary | Current active admin only | Current PostgreSQL Principal.Role |

All routes above are under /api/v1. This matrix is technical containment for the retained foundation, not a final business permission matrix.

### B9. Safe current-account endpoint and frontend consumer

GET /api/v1/auth/me returns the existing normal success envelope with only `id`, `email`, `name`, current `role`, and `is_active`. No password hash, refresh/access token, security version or credential appears. GET /users/me retains the safe profile shape with created/updated timestamps, also from the resolved snapshot.

Frontend AuthUser now includes is_active; useCurrentUser reads `/auth/me` with cancellation and its own `auth.me` query key so compact auth metadata cannot overwrite the full profile cache. Existing login/local-registration adapters already fetch /auth/me after token storage and continue doing so. Profile updates invalidate both profile and auth metadata. Backend remains authoritative; cached store/query role/status is UX metadata and may lag server changes. No new session UI/role layout is built, and the placeholder/router/providers still make no automatic registration/login/profile call.

### B10. Self-service restrictions

PATCH /users/me accepts a required display name only. ID is always the current principal's ID. Request email, role, active flag, passwords/hashes, refresh/security fields, IDs and timestamps are rejected. Email/identity editing remains contained until verification/provisioning contracts are reviewed; no password/reset flow is added. Frontend input/schema/transport likewise send name only, including if a runtime caller supplies extra fields.

Application validates/normalizes the name, then calls the repository's dedicated profile method. SQL updates only name and updated_at using parameters, with current active/known-role predicates and a safe RETURNING projection. It never invokes the broad aggregate Update or writes previously read password/role/status/email. Repository fake/statement tests protect this scope, and a stateful HTTP fake covers mid-request deactivation; actual PostgreSQL concurrent-write behavior remains an integration gate.

### B11. Privileged-route restrictions

The retained GET user list requires current active generic admin. It remains read-only, paginated, normalized to at most 100 per page, and returns existing safe UserDTOs. Generic user, inactive/deleted/unknown-role account, stale admin JWT or body role injection cannot reach the list handler. Arbitrary-user retrieval/update/deletion and all provisioning/role/status operations remain unavailable, including to generic admin. No final eLabTrack account administration is implemented.

### B12. Mass-assignment and password handling

Strict JSON decoding applies to registration and self update; unknown fields, unsupported content type, malformed input and additional JSON bodies fail generically. Explicit DTOs/commands contain only permitted input fields and never bind persistence entities. Display-name SQL cannot change ID/email/role/status/password/security fields. Local registration constructor fixes generic user/default active state and uses the retained bcrypt adapter; no new password library is introduced. Profile/current-account/list DTOs do not return hashes, and request logging records method/path/status/latency/request ID rather than bodies/tokens. No plaintext/hash logging was added.

Moving inactive login handling after password comparison fixes the discovered credential-enumeration detail without adding password reset, provisioning, dummy-hash timing, token mechanics or session-policy redesign.

### B13. Tests added

Backend adds 13 top-level tests with 63 subtests:

- Actual route authentication matrix covers all four protected route/methods against absent/malformed/invalid bearer; no account/write/list calls occur.
- Current-state matrix covers active user/admin, demoted stale-admin token, promoted old-user token, inactive account, unknown role and deleted account on current-account/self/list routes.
- Typed identity/principal availability, request isolation and exact safe /auth/me response field allowlist; no secret/hash/stale token email appears.
- Zero identity, wrapped missing account, unavailable database, nil/mismatched lookup and role-guard-without-auth/unknown required role denial.
- Body role injection, self-vs-arbitrary-user access, ten blocked profile fields, safe name-only update/security preservation and an in-flight deactivation fake.
- Production/unknown/unset route registration absence with no creation side effects; development/test local creation with real bcrypt and fixed generic-user role; local privilege injection rejection.
- Real JWT adapter through the actual router verifies current-account resolution and stale role denial without adding deeper issuer/algorithm/claim policy.
- Inactive login checks occur after bcrypt verification; password/hash do not appear in responses.
- Actual PostgreSQL repository through recording transaction/row fakes protects parameterized safe account projection, scan/error mapping, restricted SQL SET columns and mutation access predicate. No server is needed; no live race is claimed.

Frontend adds four tests: auth/me cancellation and server role/status, current-user metadata/cache separation with token values unchanged, no unauthenticated metadata request, and name-only transport despite runtime field injection. Existing placeholder/auth-store tests remain enabled; only the AuthUser fixture gains status.

### B14. Backend fmt/vet/test validation

PASS: `go fmt ./...`, `go vet ./...`, `go test ./...`. Structured final test count: **33 top-level tests + 118 subtests = 151 passing nodes**, across 10 tested packages; Phase 1A's 20 top-level/55 subtests remain. Go 1.27.1 is selected inside backend/, with `GOCACHE=/tmp/elabtrack-phase1b-go-cache` because the default home cache is not writable in this sandbox. No dependency or toolchain directive was changed.

### B15. Frontend lint/test/build validation

PASS: lint (0 errors, same 19 unchanged shadcn/hook warnings); **25 tests across 4 files**; TypeScript/Vite build. Bundle remains JS 475.69 kB / 149.97 kB gzip, CSS 179.91 kB. The production bundle and 62 primitives are unchanged because these adapters remain unmounted. No lint/test/type-check requirement was weakened.

### B16. Compose and migration validation

PASS: `docker compose config --quiet` and `docker compose --profile full config --quiet`. Compose/Dockerfiles/env config and the named volume are unchanged. No SQL migration was added or rewritten; existing role CHECK and active NOT NULL/default are sufficient, and all four migration files match baseline hashes. No services/images/volumes or migration were executed. Actual PostgreSQL account lookup/SQL/concurrency, browser end-to-end authentication, container runtime, production TLS and deployment remain unrun.

### B17. Diff and security review

`git diff --check` PASS. Source/test review verifies each retained protected route passes the same authentication/current-account boundary; current database role/status authorizes access; token/body role cannot elevate; inactive/unknown/deleted accounts fail closed; safe principal excludes password/session data; me response allowlist contains no secrets; self writes cannot mass-assign or overwrite security fields; production signup is absent with no alternate creation API; reusable password handling still hashes, and no body/password logging was added.

Baseline hashes preserve all JWT/bcrypt/opaque-token adapters, raw refresh repository, browser api-client/storage/store behavior, placeholder routes/providers, 62 UI primitives, four migration files, seed SQL, Phase 1A config/migrator/seeder, dependency manifests and 16 V1 audit files. The Phase 0 report and Phase 1A text above are preserved. No equipment/inventory/borrowing/approval/return/fine/notification/report/kiosk/campus feature or business-domain table was introduced.

### B18. Exact git status

All 40 changes are unstaged: 33 modified tracked files and 7 new files. Exact `git status --short` at Phase 1B closure:

```text
 M README.md
 M backend/README.md
 M backend/internal/application/auth/dto.go
 M backend/internal/application/auth/login.go
 M backend/internal/application/auth/refresh.go
 M backend/internal/application/auth/service.go
 M backend/internal/application/user/commands.go
 M backend/internal/application/user/dto.go
 M backend/internal/application/user/service.go
 M backend/internal/bootstrap/dependencies.go
 M backend/internal/bootstrap/migration_test.go
 M backend/internal/domain/user/repository.go
 M backend/internal/infrastructure/persistence/postgres/user_repository.go
 M backend/internal/interface/http/handlers/auth_handler.go
 M backend/internal/interface/http/handlers/user_handler.go
 M backend/internal/interface/http/middleware/auth.go
 M backend/internal/interface/http/routes/auth_routes.go
 M backend/internal/interface/http/routes/routes.go
 M backend/internal/interface/http/routes/user_routes.go
 M backend/internal/shared/constants/constants.go
 M docs/project/DECISIONS.md
 M docs/project/OPEN_DECISIONS.md
 M docs/project/PHASE1_FOUNDATION.md
 M docs/project/PHASE1_SECURITY_BACKLOG.md
 M frontend/src/features/auth/api/auth.api.ts
 M frontend/src/features/auth/hooks/use-auth.ts
 M frontend/src/features/users/api/users.api.ts
 M frontend/src/features/users/hooks/use-users.ts
 M frontend/src/features/users/schemas/update-profile.schema.ts
 M frontend/src/features/users/types/index.ts
 M frontend/src/lib/query-keys.ts
 M frontend/src/stores/auth-store.test.ts
 M frontend/src/types/common.ts
?? backend/internal/application/auth/principal.go
?? backend/internal/application/auth/principal_test.go
?? backend/internal/domain/user/account.go
?? backend/internal/infrastructure/persistence/postgres/user_account_test.go
?? backend/internal/interface/http/handlers/bind.go
?? backend/internal/interface/http/routes/auth_boundary_test.go
?? frontend/src/features/auth/auth-foundation.test.tsx
```

No staging, commit, push, remote modification or deployment is authorized/performed.

### B19. Remaining Phase 1C/1D risks and evidence limits

Deferred explicitly: exact JWT issuer/algorithm/audience/type/required-time/subject validation; refresh-token hashing at rest; atomic single-use rotation/replay/family handling; HttpOnly cookie/CSRF transport; removal of localStorage refresh/access tokens and final browser-memory strategy; bounded refresh retries; logout/access-token/family revocation redesign; concurrent-session policy and cross-tab coordination. Current valid stolen tokens still identify their account; account checks do not establish per-session logout revocation. In-flight reads use their authorized snapshot; frontend session metadata is cached and no 403 session UX redesign was performed.

SEC-006 current-account checks and SEC-007 production registration containment are implemented, with final role/provisioning/session policy unresolved. SEC-005 only gains the essential zero-ID rejection at the middleware boundary. SEC-008 locks/atomic migration bookkeeping, SEC-013 cookie/CSRF coordination and SEC-012/014–019 proxy/abuse/headers/observability/envelope/integration/CI/retention/dependency/lint work remain open or partial. Tests establish local source/unit/HTTP contracts, not live PostgreSQL concurrency or production readiness. No private Obsidian tool is exposed; no external note reads/persistence are claimed. No V1 service/data was accessed.

### B20. Product-policy decisions remain unresolved

OPEN-001 public signup versus provisioning; OPEN-002 borrower eligibility/types; OPEN-003 staff/admin authority; OPEN-004/005 Super Administrator scope; OPEN-008 deactivation/session/ongoing-loan consequences; OPEN-009 verification; and all other existing institutional/domain questions keep their prior Needs Stakeholder Input/Deferred statuses. Added engineering evidence only clarifies current containment. DEC-024–026 accept the technical boundary/containment, not any final signup, role taxonomy, eligibility or business permission matrix.

### B21. Phase 1B exit gate

**SATISFIED / COMPLETE** for the authorized Phase 1B scope: production registration contained; protected routes share authentication and PostgreSQL current-account resolution; current role/activity drive authorization; token/body role cannot elevate; generic routes guarded; safe /auth/me; focused backend/frontend tests and all requested local gates pass; no business feature added. Phase 0/1A historical evidence is intact. Phase 1C/1D has not begun. Completion is not authentication/session/deployment readiness; remaining and unrun gates are explicit above.

## Phase 1C — JWT & Refresh Session Security

Authorized 2026-10-06. Phase 1D, cookies, browser storage changes, business workflows, commit, push and deployment are outside this work.

### Pre-change audit (before implementation)

| Actual current behavior | Risk | Authorized action |
|---|---|---|
| JWTIssuer issues HS256 but accepts any HMAC method; issuer/expiry presence/subject consistency/purpose are not required. | Signed tokens with inappropriate algorithms, identity or purpose may pass. | Pin HS256, configured issuer, fixed API audience, required identity/time claims and access purpose; change issuance and verification together. |
| Access JWT embeds role/email; Phase 1B resolves a safe current PostgreSQL account on every protected request. | Returning to JWT-role authorization would restore stale privilege. | Preserve current-account boundary and its stale-role HTTP tests. |
| Refresh generator uses crypto/rand, 32 bytes, 64 lowercase hexadecimal characters. | Format is already sufficiently strong; changing it adds unnecessary incompatibility. | Retain opaque format/entropy. |
| Domain RefreshToken and PostgreSQL token column persist the raw bearer secret; errors can wrap database detail. | Database disclosure exposes usable sessions; detail may leak credentials. | Hash before persistence/lookup/revoke, hash-only repository interface, safe errors. |
| Refresh reads/checks/revokes/creates through unrelated calls. | Concurrent requests can both rotate; failed replacement can strand the old credential. | One application transaction through an inward port; PostgreSQL row lock and conditional consumption; commit before returning. |
| Refresh checks account only after revoking; checks active but not known role. | Failed/ineligible refresh can consume credentials; invalid current roles can authenticate. | Read/lock safe current account inside the transaction, reject absent/inactive/unknown role generically. |
| Logout updates raw token without an active-state predicate; revoke-all exists. | Repeated operations rewrite state; token details can leak. | Hash-based idempotent revoke; preserve internal revoke-all capability. |
| Historical migrations 000001/000002 have raw-token schema; no deployed V2 production sessions are evidenced. | Rewriting history or pretending production backfill exists would be misleading. | New paired 000003 migration, explicit session invalidation, constraints/indexes and documented down invalidation. |
| No expired/revoked session cleanup; no live PostgreSQL integration harness. | Unbounded table growth; unit fakes cannot prove database concurrency. | Bounded cleanup capability and explicit maintenance command; required live migration/concurrency coverage in Phase 1G. |
| JSON token-pair/body contract and localStorage adapters remain. | Browser-readable long-lived secrets and retry/cross-tab risks remain. | Keep compatibility; explicitly defer transport/storage/retry redesign to Phase 1D. |

No private Obsidian MCP is exposed; no private notes were read or persisted. The removed `.specify/memory/constitution.md` and `docs/SPECKIT-PLANNER.md` are absent. Current project AGENTS/source-of-truth policy and manifests govern; no unrelated template/private-note workflow is substituted.


### C1. Files changed

Auth domain/repository now represent hash-only sessions. Application adds hashing, transaction and safe locked-account ports, atomic refresh coordination, safe pair issuance and hash-based logout. Infrastructure implements the SHA-256 adapter, strict JWT policy, explicit transaction isolation/cancellation-safe rollback, session locking/consumption/revocation/cleanup and safe locked account lookup. Bootstrap wires these ports. Login also denies unknown current roles after password verification. CLI/Make wrappers expose explicit bounded cleanup. Paired 000003 is the only schema addition. Backend tests cover these boundaries; existing Phase 1B fixtures and domain session tests are adapted. Root/backend READMEs, accepted technical decisions, security backlog and this report are updated. C24 lists every changed file. No frontend source/manifest/lockfile/config change is necessary.

### C2. Previous JWT behavior

Issuer produced HS256 with uid/email/role, issuer, UUID subject and exp/iat/nbf. Verification accepted the whole HMAC family and validated optional time claims through library defaults; it did not require issuer/audience/expiry/iat/subject consistency or access purpose. Phase 1B already rejected zero identity and resolved current accounts independently of token roles. Refresh was a separate opaque credential, persisted raw and rotated through non-transactional calls.

### C3. Final access JWT validation policy

| Check | Final contract |
|---|---|
| Algorithm/signature | Exact HS256 allowlist and matching method; configured signing secret; no none/HS384/HS512 acceptance |
| Issuer | Required exact configured JWT_ISSUER (existing default elabtrack-v2); Phase 1A loader still rejects explicit blank values |
| Audience | Required membership of fixed elabtrack-v2-api; single REST API, no new environment option or institutional scope |
| Purpose | Required signature-protected purpose=access |
| exp | Required; current UTC time must be strictly before expiry, including rejection at the exact boundary |
| iat | Required and not in the future; issuance must precede expiry |
| nbf | Required and not in the future; iat <= nbf < exp |
| Clock | Zero leeway; real UTC clock in production, instance time seam in deterministic JWT/service tests; deployment hosts/database need synchronized clocks |
| Identity | Required nonzero UUID uid; required canonical UUID subject exactly equal to uid.String() |
| Authorization | uid locates the current account; current PostgreSQL role/is_active authorizes requests; token role/email cannot elevate access |
| Errors | One safe invalid-access-token error, no parser/header/claim/token/signature/secret detail |

Issuance and verification change together. Existing typed AccessTTL and RefreshTTL remain; configuration/defaults/dependencies are unchanged. Missing audience/purpose on older access tokens forces re-login. JSON token-pair response remains compatible. Claims do not establish final institutional role policy.

### C4. Token purpose/type design

JWT payload purpose=access is the authoritative token-class discriminator. JWT library's ordinary typ=JWT header is retained; no redundant header-type taxonomy is introduced. Refresh credentials are 64-character opaque random hex strings, never JWTs. They are accepted only by hashing/session lookup at refresh/logout, while access middleware requires a correctly signed JWT with the complete access contract. A JWT claiming purpose=refresh fails access verification; an opaque refresh string fails JWT parsing. No JWT refresh class exists.

### C5. Refresh token format

Retained crypto/rand.Read over 32 bytes, encoded as 64 lowercase hexadecimal characters: 256 bits of secret entropy. No timestamp, math/rand, UUID/user ID or deterministic hash generates the bearer secret. Malformed/noncanonical submitted refresh strings are rejected generically before lookup. Tests use explicitly synthetic deterministic strings in fakes; those fakes are not production security mechanisms.

### C6. Refresh hashing design

SHA-256 hashes the canonical encoded raw secret, yielding a separate 64-character lowercase hex digest. Password hashing is unnecessary for uniformly random 256-bit credentials. Raw tokens exist only in the application transport result/request; hash-only domain entities/repository parameters cannot accidentally copy a raw-token field. Login/local registration/new refresh issuance persist the digest, and presented raw refresh/logout credentials are hashed before repository access. Digests themselves are not accepted as bearer credentials: submitting one computes another digest and fails lookup. No salt/HMAC secret rotation complexity or new dependency is introduced. Repository errors discard PostgreSQL detail, including unique-constraint digests.

### C7. Database session schema changes

| Column/constraint/index | Purpose |
|---|---|
| id UUID primary key | Session identity |
| token_hash TEXT NOT NULL UNIQUE | Digest-only credential locator; unique constraint supplies lookup index |
| user_id UUID NOT NULL FK users ON DELETE CASCADE | Account ownership |
| created_at / updated_at TIMESTAMPTZ | Creation and latest state change |
| expires_at TIMESTAMPTZ NOT NULL | Authoritative persisted refresh expiry |
| revoked_at nullable TIMESTAMPTZ | NULL means unrevoked; timestamp means consumed/revoked |
| replaced_by nullable UUID FK refresh_tokens ON DELETE SET NULL | Consumed session points to its successor; no family/device model |
| Hash format CHECK | Exactly 64 lowercase hexadecimal characters |
| Lifetime / replacement CHECKs | expires_at > created_at; a replacement requires revocation; no self replacement |
| Active user index | Internal revoke-all lookup for unrevoked sessions |
| Terminal-time expression + id index | Bounded ordered cleanup by LEAST(expires_at, COALESCE(revoked_at, expires_at)) |
| Partial replaced_by index | Efficient optional replacement FK maintenance |

No raw token column remains after 000003 up. No business table, campus scope, fingerprint, geolocation or session dashboard is added.

### C8. Migration changes

New paired `000003_refresh_session_security.up.sql` / `.down.sql` only. Historical 000001/000002 up/down files match the pre-change hashes. Up explicitly DELETEs all existing raw sessions before changing the schema; no deployed V2 production sessions are evidenced and no speculative backfill is performed. Both schema directions require re-login. Down DELETEs hash-only sessions, drops new metadata/checks/indexes and restores the prior empty raw-column schema and indexes; hashes cannot reconstruct secrets. Coordinate application/schema rollback and do not run the hash-only API on the old schema.

Fresh sequence and down dependencies were reviewed from SQL/source, including historical constraint names and the runner's transactional file execution. **SQL was not executed against PostgreSQL**; fresh up/down/reapply, constraints, uniqueness, FK/link cleanup and failure recovery are REQUIRED Phase 1G integration checks. The existing runner still records migration versions separately and lacks a concurrent-run lock (SEC-008); this migration does not claim to fix that defect. No database/container/volume was modified.

### C9. Atomic rotation implementation

Application depends on Transactions.Within(ctx, callback), RefreshTokenHasher and RefreshAccounts ports; no pgx/Fiber/infrastructure import crosses inward. Bootstrap injects the existing PostgreSQL transaction manager and repositories. Within begins one pgx transaction with explicit READ COMMITTED isolation and binds that same transaction into the callback context. Both repositories use the bound transaction; session/account locking methods and Consume reject missing transaction context rather than releasing locks in autocommit.

1. Hash presented raw credential; begin transaction.
2. SELECT the hash-indexed session FOR UPDATE. Require existing, unrevoked, unexpired session and nonzero owner.
3. SELECT safe current account fields FOR SHARE. Require matching existing active account and a known temporary role. This lock conflicts with non-key role/status updates, unlike KEY SHARE. Recheck expiry after waiting.
4. Generate/sign a candidate pair; persist only the replacement digest in the same transaction. Insert precedes old-row update solely to satisfy the immediate successor FK.
5. Conditionally UPDATE the old hash only if revoked_at IS NULL and expiry exceeds both application UTC time and PostgreSQL clock_timestamp(); record revoked_at/updated_at and replaced_by. Exactly one affected row is required.
6. Commit both writes together. Only after commit succeeds does the application return the candidate pair.

Any callback/generation/lookup/write/conditional failure rolls back. Rollback gets a bounded five-second context detached from request cancellation. Begin/commit/session/account storage failures are sanitized and return no token pair. A lost commit acknowledgement or lost HTTP response has the ordinary distributed-system ambiguity: the database may already have committed while the client received no pair. No plaintext recovery cache or second successful retry is created; the client may need re-login. Live failure/connection tests remain Phase 1G.

### C10. Concurrent refresh behavior

PostgreSQL row locks, unique hashes and conditional consumption enforce at most one committed successor for a presented session across processes/API instances. A second READ COMMITTED locking reader waits, then sees the committed revoked row and returns generic 401; if the first transaction fails before commit, rollback permits a later attempt. No production mutex, Redis, grace replay or retry-success cache is used. Deterministic 16-caller service-double test yields one success and 15 denials; actual repository tests assert transaction-only FOR UPDATE and conditional mutation. **These are not proof of live PostgreSQL concurrency**; multi-connection/multi-instance verification is REQUIRED Phase 1G.

### C11. Replay behavior

Already consumed/revoked, expired, unknown or malformed presented credentials produce the same generic 401 envelope. A rotated row retains its successor link until terminal cleanup, allowing safe internal inspection without disclosing token existence. Replay never issues a second pair and does not revoke the legitimate successor. Stronger family-wide replay response would require reviewed concurrent-session/recovery policy and remains unresolved. No token/hash/raw body or replay-specific client detail is logged. Cleanup eventually removes old replay evidence; deleted hashes still fail generically.

### C12. Expiration handling

RefreshTTL from Phase 1A typed config determines persisted expiry at issuance; refresh validates the stored timestamp independently of browser state. Exact expiry is invalid. Expiry is rechecked after acquiring account locks, then the conditional consumption checks application time and the live database clock immediately before mutation. Revoked records cannot become valid through an expiry extension. JWT expiry remains independent and enforced by strict verification. No institution-specific maximum session lifetime/sliding-session policy is invented.

### C13. Current-account refresh behavior

Safe locked account projection excludes passwords/session material. Missing/deleted/nil/mismatched/inactive/unknown-role accounts cannot rotate; the client receives generic 401 without account/token existence detail. Unavailable account storage fails closed with safe 500. Replacement access claims use current account ID/email/role, not historical token/session hints. A status/role write occurring after rotation commits can complete; subsequent protected requests still resolve current database state through Phase 1B. No speculative suspension/ongoing-loan restriction is added.

### C14. Revocation/logout behavior

Logout hashes the presented raw credential and updates only an unrevoked matching session. Unknown, already revoked and malformed credentials are idempotent 204 with no body; database failures return safe 500. Missing/invalid request bodies retain existing 400/422 validation semantics. Revoked sessions cannot refresh. Logout of an old rotated credential does not revoke the successor; no family traversal is invented.

Internal RevokeAllForUser remains parameterized and now sets revocation timestamps only on unrevoked rows, with safe errors. It affects rows visible to that SQL statement and does not serialize concurrent/future login/rotation; final global account-session policy remains open. No logout-all endpoint/UI is introduced. A previously issued access JWT still works until expiry when current account role/activity permits it; logout does not establish immediate access-session revocation. Disabled/missing accounts remain blocked by Phase 1B/current refresh checks.

### C15. Cleanup/retention policy

Technical refresh-record retention: delete sessions more than seven days after their earlier expiry/revocation timestamp, in one batch of at most 1000. Active/unexpired and recently terminal rows remain. Cleanup uses the terminal expression index, ordered candidates, FOR UPDATE SKIP LOCKED and one DELETE statement. Limits outside 1..1000 fail. Deleted successor targets clear optional replaced_by links; no bearer material/history reconstruction is possible.

Explicit `go run ./cmd/api --sessions-cleanup` from backend/ or root/backend `make sessions-cleanup`; mutually exclusive with migration/seed actions. Normal startup does not clean up or schedule jobs. CLI prints only deleted count and a fixed safe failure. Operators must arrange regular runs (initially daily), repeat sufficient batches to drain eligible backlog, and monitor rotation volume/table growth; cadence/ownership must be established before deployment. No scheduler/maintenance action was executed. This retention is not institutional borrowing/audit policy. Index behavior and cleanup/rotation contention remain REQUIRED Phase 1G checks.

### C16. Logging safety

Changed JWT/hash/session/application paths introduce no token/body/header/password logging. JWT verification drops parser errors; refresh/pair/logout/account-lock/session repository errors return safe sentinel errors without wrapping SQL values. Begin/commit errors are sanitized. Synthetic raw-token/hash/detail sentinels test adapter/service/repository and HTTP failures. Existing request logger records time/method/path/status/latency/request ID only; cleanup prints count only. No raw JWT, refresh string, digest, Authorization header, password or signing secret is emitted by these paths. Full structured recovery/error observability remains separate SEC-015/Phase 1F work.

### C17. Frontend compatibility changes

None. All 121 baseline frontend source/config/manifest/lockfile files match baseline hashes. Token JSON remains access_token / refresh_token / expires_at / token_type=Bearer; refresh/logout still accept refresh_token bodies. Existing localStorage, auth store, Axios bearer injection/single-flight behavior and account adapters remain compatible. Expired/revoked/unknown/inactive refresh now uses the documented generic 401. No auth screen/session hydration/cookie/storage redesign was performed; the placeholder and 62 primitives remain.

### C18. JWT tests added

New security tests use a fixed clock and cover valid issuer-to-verifier roundtrip; invalid signature; HS384/HS512/none rejection; expired/exact-boundary/missing exp; wrong/missing issuer and audience; wrong/missing purpose; missing/malformed/mismatched subject; missing/zero uid; missing/future iat and nbf; time ordering; opaque refresh rejected as access. Error output is constant and contains no credential. Opaque generator/digest tests retain 32-byte lowercase-hex format, verify SHA-256 semantics and reject malformed credential inputs. Phase 1B's real-JWT/current-lower-database-role HTTP denial remains enabled and passing.

### C19. Refresh/rotation tests added

Service tests cover digest-only initial issuance, raw-to-hash lookup, stored digest rejected as bearer, current account/current role and typed TTLs, successful rotation/old replay denial/replacement reuse, exact/past expiry, revoked/unknown/malformed credentials, missing/inactive/unknown-role/mismatched account and expiry while waiting. Failure injection at lookup/account/access issuance/random generation/replacement creation/consumption/commit returns no pair and rolls back the test double; old credential remains usable after clearing the injected failure. A 16-caller transaction-double test produces at most one successor. Logout tests cover repeated/unknown/malformed revocation, hash-only lookup, refresh denial and safe errors.

Actual transaction-manager tests verify explicit isolation, shared transaction context, commit/rollback/begin/commit failures and cancellation-independent rollback. Actual repository statements are exercised via recording pgx transaction/row doubles: hash parameters, safe scans, transaction-only session/account locks, once-only expiry/revocation predicates, zero-row loss, idempotent session/revoke-all, bounded cleanup and SQL-detail redaction. Actual-router tests verify identical 401 envelopes for invalid session/current-account cases, idempotent 204 logout and safe storage 500s. Domain test enforces exact expiry/revocation with a fixed clock. No live database is represented by these doubles.

### C20. Backend fmt/vet/test results

PASS: go fmt ./..., go vet ./..., go test ./... (structured JSON run). **51 top-level tests + 173 subtests = 224 passing nodes across 11 tested packages**, zero failures; 18 additional top-level tests and 55 additional subtests relative to Phase 1B. Also PASS: go test -race for application/auth, infrastructure/database, infrastructure/persistence/postgres, infrastructure/security and interface/http/routes. No race detected in the tested local code/doubles; this does not establish database race correctness.

Commands run from backend/ with Go 1.27.1 and GOCACHE=/tmp/elabtrack-phase1b-go-cache, reusing the writable prior-phase cache. No toolchain directive/dependency downgrade or test-gate weakening. Live PostgreSQL and container runtime checks are unrun.

### C21. Frontend lint/test/build results

PASS: npm run lint (zero errors, same 19 shadcn/hook warnings), npm run test:run (**25 tests / 4 files**), npm run build (TypeScript + Vite). Unchanged assets: JS 475.69 kB / 149.97 kB gzip; CSS 179.91 kB / 27.63 kB gzip. No frontend source or gate was weakened. Browser end-to-end/session transport checks remain unrun.

### C22. Compose validation

PASS: docker compose config --quiet and docker compose --profile full config --quiet. Compose/Dockerfiles/environment examples/config are unchanged. No image built, container/service started, database connection/migration/cleanup/seeding performed, named volume removed or deployment attempted. Production TLS/CA connectivity and release/runtime checks remain unrun.

### C23. git diff --check and preservation review

PASS. Baseline hashes preserve 121 frontend files, all four historical migration files, 16 V1 audit files, dependency manifests/lockfiles, Phase 0 report, OPEN_DECISIONS and Phase 1A typed configuration. This report starts with the exact pre-1C Phase 1A/1B text; prior audit/gate evidence is preserved. Production code has no process mutex/session store/Redis; domain/application import no Fiber/pgx/infrastructure; SQL is parameterized; only hash fields enter persistence; protected-route current-account authorization remains. No equipment/inventory/borrowing/approval/return/fine/email/report/kiosk/campus capability or business schema was introduced.

### C24. Exact git status

All 30 changes are unstaged: 22 modified tracked files and 8 new files. Exact git status --short at Phase 1C closure:

```text
 M Makefile
 M README.md
 M backend/Makefile
 M backend/README.md
 M backend/cmd/api/main.go
 M backend/cmd/api/main_test.go
 M backend/internal/application/auth/login.go
 M backend/internal/application/auth/refresh.go
 M backend/internal/application/auth/service.go
 M backend/internal/application/ports.go
 M backend/internal/bootstrap/dependencies.go
 M backend/internal/domain/auth/entity.go
 M backend/internal/domain/auth/repository.go
 M backend/internal/infrastructure/database/transaction.go
 M backend/internal/infrastructure/persistence/postgres/auth_repository.go
 M backend/internal/infrastructure/persistence/postgres/user_repository.go
 M backend/internal/infrastructure/security/jwt.go
 M backend/internal/interface/http/routes/auth_boundary_test.go
 M backend/tests/unit/domain/auth/auth_test.go
 M docs/project/DECISIONS.md
 M docs/project/PHASE1_FOUNDATION.md
 M docs/project/PHASE1_SECURITY_BACKLOG.md
?? backend/internal/application/auth/refresh_test.go
?? backend/internal/infrastructure/database/transaction_test.go
?? backend/internal/infrastructure/persistence/postgres/auth_repository_test.go
?? backend/internal/infrastructure/security/jwt_test.go
?? backend/internal/infrastructure/security/refresh_hash.go
?? backend/internal/interface/http/routes/refresh_session_test.go
?? backend/migrations/000003_refresh_session_security.down.sql
?? backend/migrations/000003_refresh_session_security.up.sql
```

No staging, commit, push, remote modification or deployment was performed.

### C25. Items remaining for Phase 1D

Browser transport/storage design, HttpOnly/Secure/SameSite refresh-cookie decision and CSRF/origin protection, access-token memory strategy, session hydration, bounded 401 retry, cross-tab rotation, store/query synchronization and logout UX. Access and refresh credentials remain readable from localStorage; stolen credentials remain usable subject to server session/current-account checks. Lost successful refresh response and simultaneous tabs can require re-login under strict once-only rotation. No grace policy, cookie, final session UI or Phase 1D implementation has begun.

### C26. Items remaining for Phase 1G

REQUIRED isolated real PostgreSQL tests: fresh 000001→000003 up, constraints/unique hashes, down/reapply/invalidation, transaction rollback on replacement/consumption failure, raw-token absence and presented-hash rejection; simultaneous refresh through separate connections/instances (at most one committed successor, loser 401, replacement usable); account deactivation/role update versus rotation; concurrent logout/revoke-all/rotation and cleanup/FK behavior; actual expiry clock predicates; bounded cleanup retention/index performance; canceled request/lost connection/commit acknowledgement behavior. CI reproducibility, migration-runner bookkeeping/lock recovery, live HTTP/database auth, browser/container/production transport remain unverified and require their authorized gates. Local registration user/session failure recovery remains open for later account provisioning design. No live concurrency claim is made here.

### C27. Product/security policies still unresolved

Institutional concurrent-session limits, family-wide replay revocation, logout-all/global serialization, immediate access-session revocation, device/session management and stronger recovery/grace policy are not accepted. OPEN-001 signup/provisioning, OPEN-002 eligibility/types, OPEN-003 staff/admin authority, OPEN-004/005 Super Administrator scope, OPEN-008 deactivation/ongoing-loan consequences, OPEN-009 verification and every other existing institutional question retain prior Needs Stakeholder Input/Deferred statuses. DEC-027–029 record technical token/transaction/cleanup choices only. No private-note access/persistence, V1 access or stakeholder policy resolution is claimed.

### C28. Phase 1C exit gate

**SATISFIED / COMPLETE** for the authorized Phase 1C scope: strict access JWT contract and purpose separation; retained cryptographically strong opaque refresh format; hash-only persistence/lookup/revoke; real PostgreSQL transaction boundary and once-only lock/conditional design; safe replay/expiry/revocation/current-account denial; safe error/log paths; focused JWT/session/HTTP/repository/transaction tests; backend fmt/vet/tests and frontend lint/tests/build; Compose configuration and diff checks; Phase 1A/1B evidence preserved. Live PostgreSQL concurrency/migration/rollback verification is explicitly required in Phase 1G as permitted by the request. Completion does not establish production/session transport readiness. No Phase 1D or business feature work, commit, push or deployment occurred.


## Phase 1D — Browser Session Architecture

Authorized 2026-10-06. Phase 1E, business features, migrations against live/unknown databases, commit, push and deployment are outside this work.

### Pre-change browser/session audit (before implementation)

| Source | Existing reads/writes/flow | Risk and selected action |
|---|---|---|
| frontend/src/stores/auth-store.ts | setSession writes both tokens to storage/state; hydrate reads storage; clear removes keys and Query cache. | Tokens survive reload and XSS can read long-lived credentials. Replace with a non-persisted access/user/status store; bootstrap through cookie only. |
| frontend/src/features/auth/api/auth.api.ts | Login/local register persist both JSON tokens before /auth/me; refresh/logout send raw refresh_token bodies. | Refresh is JavaScript-readable and has a second transport path. Remove refresh DTO/state, centrally coordinate cookie login/refresh/logout. |
| frontend/src/lib/api-client.ts | Reads access from storage for Authorization; reads refresh for JSON refresh; writes replacement tokens; in-instance refresh promise; retry lacks a request marker and includes login/logout. | Recursive retries and competing flows can consume a single-use credential. Use memory, one shared refresh promise, one retry marker, explicit auth-endpoint exclusions and lifecycle generation checks. |
| frontend/src/features/auth/hooks/use-auth.ts | Login/register setSession again; logout reads persisted refresh and clears locally. | Duplicate/later writes can restore stale state after logout. Let the central client own session mutations; clear immediately and serialize logout behind pending cookie responses. |
| frontend/src/lib/storage.ts | Defines readable access/refresh keys; generic preferences use localStorage. | Remove auth keys and purge known legacy auth entries without reading their values; retain non-sensitive preferences. |
| frontend/src/app/providers.tsx, main.tsx, query-client.ts | No session bootstrap is mounted; React StrictMode; Query retries once. | Reload restoration absent; repeated mount/query retries could rotate twice. Add one controlled bootstrap, recoverable error state and no Query retry on auth failures. |
| backend auth handlers | Login/local register/refresh return TokenPairDTO including raw refresh; refresh/logout bind raw JSON body; no cookie helper. | Raw token exposed to JavaScript. Issue/clear one host-only HttpOnly refresh cookie, return access + safe current account only, remove raw body input. |
| backend current-account/JWT/session code | Strict HS256 access; safe PostgreSQL authority; hash-only transactional once-only rotation and idempotent revoke. | Preserve all Phase 1B/1C server guarantees; no new migration or family/replay policy. |
| config/CORS/routes | Typed environment and explicit origins; credentialed CORS already enabled; auth mutations POST-only; no request Origin check. | CORS alone does not prevent cookie CSRF. Add exact trusted-Origin enforcement to cookie-changing auth endpoints, including login/local register; retain explicit allowlist CORS. |

Chosen model: memory-only access; opaque refresh cookie elabtrack_v2_refresh, host-only, Path=/api/v1/auth, HttpOnly, SameSite=Lax, Secure in production (also fail-safe for unknown environment), expiry taken directly from the newly persisted refresh session. HTTP development/test are intentional exceptions derived only from typed APP_ENV. Cross-site SPA/API cookie deployments are unsupported by this SameSite policy; no configurable None/insecure production bypass. Session bootstrap/login/refresh use current safe server account metadata, not JWT decoding. Known-dead refresh failures clear the cookie; transient storage failures retain it for deliberate recovery.

No Obsidian MCP is exposed; no private notes were read or persisted. Removed planner/constitution files remain absent; current project instructions and source govern.

### D1. Files changed

29 unstaged files: 24 modified tracked files and five new files; exact inventory in D30. Backend changes cover safe session DTOs, cookie/Origin handler adaptation, wiring and HTTP tests. Frontend changes cover the centralized transport, memory store, bootstrap, auth adapters/hooks, query retry and storage cleanup, with focused tests. Documentation updates this report, technical decisions, security backlog and root/backend READMEs. No dependencies, environment/configuration files, infrastructure, migrations, business schema or reusable UI primitives changed.

### D2. Previous browser token design

The audit above was recorded before implementation. Both access and raw refresh credentials were returned in JSON, written to localStorage and hydrated into Zustand. Axios read persistent credentials, submitted refresh JSON and stored rotated tokens. Refresh retry lacked a per-request bound; bootstrap was not mounted. Auth hooks duplicated session writes and could restore late responses after logout. Generic non-sensitive preferences also use storage and remain supported.

### D3. Final access-token storage strategy

Non-persisted Zustand holds accessToken, safe current user, status and a lifecycle generation. Access JWT travels only as Authorization: Bearer from this memory. No persistence middleware, localStorage/sessionStorage/IndexedDB write, access cookie or JWT decoding. Reload loses that access token and creates a fresh application store. Query owns server data, never credential state. Clear/account changes clear Query cache; generation checks reject stale asynchronous session writes/retries. Existing short access TTL and server JWT validation remain unchanged.

### D4. Final refresh-token browser transport

The browser receives raw refresh solely through Set-Cookie and sends it automatically through credentialed requests. Frontend state/types/requests contain no raw refresh credential or readable refresh-cookie parser. Refresh/logout ignore raw body input and read only the cookie; there is no alternate production transport. Application TokenPairDTO retains raw material only internally for the HTTP cookie boundary, with json:"-" as an additional guard. PostgreSQL still receives SHA-256 digests only; crypto/rand generation, transaction locks, current-account validation and once-only rotation are retained unchanged.

### D5. Cookie name and attributes

| Attribute | Selected contract |
|---|---|
| Name | elabtrack_v2_refresh, owned by this project |
| HttpOnly | true, frontend JavaScript cannot read this cookie |
| Secure | true in production; also the fail-safe attribute for unknown environments |
| SameSite | Explicit Lax; production SPA/API must be same-site HTTPS |
| Path | /api/v1/auth; delivery covers refresh/logout without ordinary /users routes |
| Domain | Omitted; host-only API cookie |
| Expires | Newly persisted refresh-session ExpiresAt; serialized to HTTP-date seconds |
| Clearing | Same name/path/host-only/HttpOnly/Secure/Lax, empty value, past expiry, MaxAge=-1 |

One helper issues/clears cookies. There is no independent cookie TTL or security environment override. Rotation replaces the cookie with expiry from the replacement session. Path is delivery scope, not a security boundary; HttpOnly does not prevent injected JavaScript from making authenticated requests. Lax does not support cross-site fetch-based SPA authentication. These semantics follow [MDN Set-Cookie](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Set-Cookie). Actual browser enforcement remains unrun.

### D6. Development cookie behavior

Explicit APP_ENV=development/test derives Secure=false while keeping HttpOnly, Lax, path and host-only scope. Existing explicit localhost origins permit local HTTP SPA/API ports on the same site. No implicit missing-Origin bypass is allowed for tools/tests; they send a trusted Origin deliberately. Production registration containment remains unchanged. HTTP localhost behavior is covered by handler tests, not a real browser run.

### D7. Production cookie behavior

Typed production requires Secure=true and explicit HTTPS trusted origins; no arbitrary flag can disable it. Unknown environment, empty/wildcard/invalid cookie origins and production HTTP origins deny cookie mutations before service IO. Production supports same-site HTTPS, normally the existing same-origin nginx /api/v1 proxy. Different registrable sites/schemes are unsupported by this Lax contract; no SameSite=None option was added. TLS/proxy/browser compatibility must be verified in Phase 1G after the broader Phase 1E review.

### D8. Login response changes

Login/local registration persist the existing hash-only refresh session, issue its raw value only as cookie, and return the normal envelope containing exactly access_token, expires_at, token_type and user. User contains id/email/name/role/is_active only. No refresh secret/hash, password hash, session identifier/replacement metadata or refresh expiry in normal JSON. Safe user metadata is projected from the server account used for issuance. Session endpoints set Cache-Control: no-store, including errors. Local registration retains the existing separate user/session-write recovery limitation; production register remains 404.

### D9. Refresh response changes

Bodyless POST /auth/refresh reads the cookie, validates current account and rotates through the unchanged Phase 1C transactional mechanism. Only committed success issues the replacement cookie and safe access/user JSON. Old/replayed/expired/revoked/unknown/malformed credentials fail generically; the HTTP adapter does not weaken single-use constraints. A presented raw JSON credential without a cookie cannot refresh. No database schema/migration change is needed.

### D10. Logout behavior

Frontend clears access/user/cache immediately and increments generation before any network acknowledgement. It submits bodyless credentialed logout, serialized behind pending refresh/login responses so their latest cookie can be revoked. Late promises cannot restore the cleared state. Backend revokes the cookie's session idempotently, clears matching cookie attributes and normally returns 204 even for missing/invalid/already-revoked credentials. Storage failure returns safe 500 but still clears the cookie.

Network failure leaves this document locally unauthenticated; it cannot prove server revocation or delete an HttpOnly cookie from JavaScript. A later reload may restore a still-valid cookie if logout never reached the server. Storage failure can likewise leave a backend record until expiry. Logout does not globally revoke concurrent sessions or immediately invalidate access JWTs; current account checks and expiry still apply. No stronger policy is implied.

### D11. Session bootstrap flow

Providers mounts SessionBootstrap before application content. Initial idle/bootstrapping shows an accessible restoring-session status, then one cached bootstrap promise invokes the shared refresh operation using browser credentials. Repeated StrictMode effects share the same client/module/promise. Success accepts safe server access/user and marks authenticated; 401/403 marks unauthenticated; unexpected network/server failure clears credentials and shows a recoverable error banner with an explicit Retry session action. Public placeholder/status content is then available. No automatic restoration loop or final login/dashboard UI was added.

Transport is dynamically imported to preserve the existing bundle warning gate; module loading failure is recoverable. An explicit retry starts one new bootstrap attempt only from error. Existing public health requests remain manual; bootstrap itself can rotate an existing session and is now part of application loading.

### D12. Current-account handling

Login/refresh returns a safe server-account projection. The client validates the response shape and known temporary role/active flag, then retains only safe fields. It never decodes JWT claims as authorization truth. /auth/me remains a current-account read, with hook generation/identity checks before store updates. Protected backend routes still verify bearer identity then resolve current PostgreSQL role/status. Frontend role is presentation data only; cookie alone cannot authenticate protected APIs. No institutional role or provisioning decision changed.

### D13. Single-flight refresh implementation

One runtime Axios singleton owns refreshPromise, reused by bootstrap, manual refresh and eligible concurrent 401s. Session mutations share a promise queue so refresh/login/logout cannot compete within this client. No sleeps, secondary runtime transport or server grace relaxation. A late 401 for the previous access token uses an already renewed memory token rather than rotating again. Lifecycle generation checks reject work belonging to a cleared/replaced session.

Coordination is per application document. Independent tabs are not coordinated. Simultaneous tabs may consume the same single-use cookie; a losing invalid-refresh response can also clear a winner's cookie. Cross-tab/lost-response behavior is explicitly required Phase 1G verification and may need later approved client coordination or policy; no cross-tab safety claim is made.

### D14. Retry limit

Only a protected 401 sent with this session's memory token is eligible. Each request receives a sessionRetry marker before refresh and is retried at most once using the replacement token. A second 401 clears the session and stops. Login/register/refresh/logout, unauthenticated requests, disabled auth/retry options, stale generations and 403/validation/server/business failures cannot start automatic refresh. Query also excludes 401/403 from its ordinary retry so it cannot restart that auth loop. Existing one-retry behavior for other Query failures remains unchanged.

### D15. Refresh failure behavior

Any refresh failure clears memory/user/cache for the current generation and rejects queued protected requests without recursive refresh. Auth 401/403 becomes unauthenticated; transient failures become an error/recoverable state. Backend clears known-dead refresh cookies on generic 401, but retains cookies on transient storage 500 to allow a deliberate recovery attempt. The client never inspects or deletes the HttpOnly credential. Stale responses cannot clear a newer session.

### D16. CORS changes

Existing typed explicit origin validation and Fiber AllowCredentials=true were sufficient and remain unchanged. The centralized Axios instance uses withCredentials=true. Actual-router tests verify exact trusted preflight ACAO/credentials and no permissive origin response for an untrusted origin. No wildcard or arbitrary origin reflection is introduced. Credentialed CORS requires exact allowed-origin responses, consistent with the [Fetch standard](https://fetch.spec.whatwg.org/#origin-header). Full CORS deployment/proxy review remains Phase 1E; browser enforcement remains Phase 1G.

### D17. CSRF policy

Defense in depth consists of host-only HttpOnly cookie, explicit Lax, auth path, POST-only cookie mutations, exact trusted-Origin enforcement and explicit credentialed CORS. Guard applies to login/local register as well as refresh/logout, preventing login CSRF/cookie replacement from untrusted origins. Missing/null Origin is rejected; no GET mutation or JSON credential fallback exists. CORS alone is not the CSRF guard. Requests from trusted origins can use credentials; protected APIs require the separate bearer token and current account. No elaborate CSRF-token framework or API-wide HTTP redesign was needed. XSS/trusted-origin compromise remains a distinct risk.

### D18. Browser-storage cleanup

Complete source search shows no application write/read of persisted authentication credentials. Only removal code references known old keys: elabtrack_v2.access_token, elabtrack_v2.refresh_token, fst.access_token and fst.refresh_token. Bootstrap, accepted session and clear purge these from localStorage and sessionStorage without reading values; storage exceptions are handled. No prior auth IndexedDB use exists. Generic localStorage/theme helpers and sidebar preference cookie remain unchanged. No raw refresh field, request DTO or auth persistence middleware remains.

### D19. Backend cookie tests

New browser_session_test.go tests real mounted HTTP handlers with synthetic service/repository doubles: login development/test/production attributes and digest storage, production Secure, local HTTP/register, exact safe JSON, refresh replacement/expiry/replay, body-only rejection, logout body ignored/idempotency/clearing, invalid cookie clearing, transient refresh retention, logout storage failure clearing, exact Origin variants and invalid policy denial before IO, CORS preflight, GET non-mutation and cookie-only protected-route denial. Existing Phase 1B/C HTTP tests were adapted to the transport contract while preserving authorization/password/rotation assertions. These tests prove handler behavior, not PostgreSQL transaction or browser enforcement.

### D20. Frontend persistence tests

Store/API tests verify access/user only in memory, raw refresh absent, safe projection discards extra sensitive server fields, legacy key removal, no auth storage reads/writes and retained theme preference. Existing account/cache tests remain; login/register feature adapters return only safe user metadata to mutation hooks, and manual refresh returns no credentials to feature callers. Three adapter boundary tests verify those results; Query receives metadata rather than credentials. Login fixture, refreshed state, reload simulation and logout all leave readable credential stores empty. No production secrets appear in fixtures/snapshots.

### D21. Bootstrap tests

API adapter tests cover one shared bootstrap, successful safe current account, missing/invalid 401/403, simulated new-document restoration, bounded network/server failure and deliberate retry. A StrictMode Providers/jsdom test verifies loading gate, public recoverable banner, one initial transport call and one explicit recovery attempt. Test transport uses synthetic cookie-server response simulation; no browser cookie jar, actual navigation/reload or live API was exercised.

### D22. Concurrent/single-flight tests

Eight simultaneous protected 401s produce exactly one refresh and eight retries with the replacement Authorization header. Late old-token 401 uses the new token without another refresh. Six simultaneous callers sharing failed refresh 401/503 produce one refresh and clean failure without recursion. Deferred promises verify pending refresh/login and later protected denials cannot resurrect logged-out memory state. No sleeps are used to create concurrency; independent browser tabs are outside these doubles.

### D23. Retry-loop tests

Protected 401 → successful refresh → retried 401 stops after three transport calls, clears session/cache and does not refresh again. Query cache clearing intentionally cancels pending query work; separate Query tests prove 401/403 receive no library retry. Auth-endpoint denial cases do not recursively intercept; attachAuth=false and non-401 failures do not rotate credentials. Quality gates and previous security assertions remain enabled.

### D24. Logout tests

Frontend tests verify immediate memory/cache invalidation before response, credentialed bodyless backend call, empty readable storage after failed network logout and no state resurrection. Pending refresh/login responses are followed by queued logout. Backend HTTP cases verify repeated/missing/invalid cookie 204 clearing and safe storage 500 with clearing. No test claims an unreachable server has revoked its record or that JavaScript can delete HttpOnly material.

### D25. Token/log safety

Modified auth code contains no secret-bearing logging or token response console output. HTTP JSON uses explicit allowlists, internal refresh fields are json:"-", errors are generic and Axios errors are normalized without retaining request headers/config/raw transport errors. Existing request logging remains method/path/status rather than cookie/body/Authorization logging. Password, access, refresh/hash, signing secret and Cookie/Authorization values are not written to logs. Backend hash persistence/security sources remain unchanged. Broader structured logging/error contracts remain later-phase work.

### D26. Backend fmt/vet/test results

PASS: go fmt ./..., go vet ./..., go test ./... (structured JSON run). **59 top-level tests + 182 subtests = 241 passing nodes across 11 tested packages**, zero failures; eight additional top-level tests and nine subtests relative to Phase 1C. PASS targeted go test -race ./internal/application/auth ./internal/interface/http/routes ./internal/bootstrap; no detected races in tested code/doubles. Commands used Go 1.27.1 from backend/ and the writable GOCACHE=/tmp/elabtrack-phase1b-go-cache. No dependency/toolchain/gate downgrade. Live database checks remain unrun.

### D27. Frontend lint/test/build results

PASS: npm run lint (zero errors, same 19 unchanged shadcn/hook warnings), npm run test:run (**54 tests / six files**), npm run build (TypeScript + Vite). 29 additional tests/two new files relative to Phase 1C. Assets: entry JS 477.83 kB / 150.77 kB gzip; dynamically loaded API client 54.42 kB / 19.98 kB gzip; CSS 179.91 kB / 27.63 kB gzip. Transport splitting avoids the default 500 kB entry warning without changing warning limits. Manifests/locks/UI primitives/strict TypeScript settings are unchanged. No real browser end-to-end session test claimed.

### D28. Compose validation

PASS: docker compose config --quiet and docker compose --profile full config --quiet. Compose, Dockerfiles and environment examples are unchanged. No image build, service/container startup, database connection, migration, seed, cleanup, named-volume deletion or deployment occurred. Full Docker SPA/API/PostgreSQL auth remains Phase 1G.

### D29. git diff --check and preservation review

PASS. This report begins with the exact pre-1D Phase 1A/1B/1C content. Baseline hashes preserve all six migration files including unrun 000003, configuration, dependency manifests/locks, infrastructure, 62 reusable UI files, 16 V1 audit files, Phase 0 report and OPEN_DECISIONS. Phase 1C hashing/transaction/JWT infrastructure and Phase 1B current-account authorization remain. No equipment/borrowing/approval/return/fine/email/report/kiosk/campus code or schema was introduced. Final exact status is below.

### D30. Exact git status

All 29 changes are unstaged: 24 modified tracked files and five new files. Exact git status --short at Phase 1D closure:

```text
 M README.md
 M backend/README.md
 M backend/internal/application/auth/dto.go
 M backend/internal/application/auth/refresh.go
 M backend/internal/bootstrap/dependencies.go
 M backend/internal/bootstrap/migration_test.go
 M backend/internal/interface/http/handlers/auth_handler.go
 M backend/internal/interface/http/routes/auth_boundary_test.go
 M backend/internal/interface/http/routes/auth_routes.go
 M backend/internal/interface/http/routes/refresh_session_test.go
 M docs/project/DECISIONS.md
 M docs/project/PHASE1_FOUNDATION.md
 M docs/project/PHASE1_SECURITY_BACKLOG.md
 M frontend/src/app/providers.tsx
 M frontend/src/app/query-client.ts
 M frontend/src/features/auth/api/auth.api.ts
 M frontend/src/features/auth/auth-foundation.test.tsx
 M frontend/src/features/auth/hooks/use-auth.ts
 M frontend/src/features/auth/types/index.ts
 M frontend/src/lib/api-client.ts
 M frontend/src/lib/storage.ts
 M frontend/src/stores/auth-store.test.ts
 M frontend/src/stores/auth-store.ts
 M frontend/src/types/common.ts
?? backend/internal/interface/http/handlers/session_cookie.go
?? backend/internal/interface/http/routes/browser_session_test.go
?? frontend/src/app/session-bootstrap.test.tsx
?? frontend/src/app/session-bootstrap.tsx
?? frontend/src/lib/api-client.test.ts
```

No staging, commit, push, remote modification or deployment performed.

### D31. Remaining Phase 1E work

Complete CORS review, trusted proxy configuration/IP trust model, targeted auth abuse/rate controls, general API/SPA security headers, production reverse-proxy/TLS assumptions, structured logging and full API error-contract normalization. Only cookie transport, its no-store responses and minimum session Origin protections were implemented here. Phase 1E has not begun.

### D32. Remaining Phase 1G live verification

REQUIRED: isolated PostgreSQL migration 000001→000003, down/reapply invalidation, hash-only persisted values/constraints and presented-hash rejection; transactional rollback/commit/lost acknowledgement; simultaneous refresh on separate connections/instances; account changes versus rotation; logout/revoke-all races; cleanup/FK/index/retention behavior; migration runner lock/bookkeeping and reproducible CI. Migration 000003 remains unchanged and unrun; Phase 1D adds no migration.

REQUIRED real browser checks: Set-Cookie HttpOnly/Secure/Lax/host/path/expiry/clearing on intended HTTP development and HTTPS production origins; credentialed preflight and trusted-Origin denial; login, actual reload bootstrap, access expiry, concurrent protected 401s, replay/invalid cookie, recoverable failures and logout acknowledgement; network/lost response behavior; independent tabs competing for rotation and late cookie clearing; full Docker frontend/API/PostgreSQL flow and TLS/proxy behavior. Unit/HTTP/jsdom doubles do not prove these runtime properties. Cross-tab coordination is not implemented and must be evaluated before production; no grace/family policy was accepted.

### D33. Product/security policies still unresolved

All existing institutional OPEN_DECISIONS retain their statuses, including signup/provisioning, borrower eligibility, staff/admin/Super Administrator authority, verification and deactivation/ongoing-loan consequences. Device/session management, concurrent-session limits, family-wide replay revocation, logout-all/global issuance serialization, immediate access-session revocation and stronger lost-response/grace recovery are unresolved. Technical browser choices alone are accepted in DEC-030/031. Local registration atomicity/recovery, deployment/cleanup ownership and cross-tab coordination remain later work. No private notes were read/persisted, stakeholder policies resolved or V1 systems accessed.

### D34. Phase 1D exit gate

**SATISFIED / COMPLETE** for authorized Phase 1D source-level scope: cookie-only HttpOnly refresh transport, no frontend raw refresh, memory-only access, safe login/refresh JSON, controlled restoration, shared single-flight refresh, one bounded retry, immediate local logout plus idempotent backend revoke/clear, minimum explicit credentialed CORS/CSRF policy, focused backend/frontend tests and all required local gates. Prior Phase 1A/B/C evidence is preserved. Browser/PostgreSQL/deployment checks remain explicitly required Phase 1G; completion does not establish production readiness or cross-tab coordination. No Phase 1E/business feature, migration execution, commit, push or deployment occurred.

## Phase 1E — HTTP & API Security

Authorized 2026-10-06. Initial working tree was clean. No Phase 1F/1G execution, business feature, database/migration/maintenance action, commit, push or deployment is authorized here.

### Pre-change HTTP perimeter audit

| Order/source | Current purpose and trust assumptions | Risk / Phase 1E action |
|---|---|---|
| Fiber configuration, before middleware | Read 10s, write 15s; idle derived as their sum; global body 1MB; no Server banner; default proxy trust disabled. c.IP uses socket peer; forwarding scheme/host ignored by installed Fiber when proxy trust is off. | Proxied clients share a peer bucket; no deliberate operator proxy setting. Add narrow typed IP/CIDR configuration and one shared client-IP/HTTPS resolver; preserve socket-derived Fiber host/scheme. Explicit idle/header bounds. |
| 1 RequestID | Preserves supplied X-Request-ID or generates one; response/context correlation. | Full correlation/structured logging design stays Phase 1F. |
| 2 Helmet security headers | API nosniff/frame/referrer/permissions/CSP/cross-origin headers; HSTS depends on Fiber scheme. | Browser document is served separately; generic JSON/CORP/COEP policy can contradict cross-origin credentialed API access. Replace with explicit API baseline and production authoritative-HTTPS HSTS. |
| 3 CORS | Validated explicit origins, credentials; GET/POST/PUT/PATCH/DELETE/OPTIONS; Origin/Content-Type/Accept/Authorization/X-Request-ID. | Fiber normalizes request origin case, auth handler compares exactly; methods include absent routes. Align exact canonical origin matching and required methods/headers/preflight denial. |
| 4 Recovery | Recovers panic; stack handling disabled, custom callback therefore provides no intended diagnostic. | Preserve safe diagnostics without logging arbitrary panic/error values or leaking them to HTTP. Move outermost so perimeter middleware is covered. |
| 5 Request logger | Time/method/path/status/latency/response request ID, stdout; no body/auth/cookie headers. | Add shared effective-IP signal without parsing headers here; retain logging format architecture for Phase 1F. |
| 6 Global limiter | Fiber fixed window, process-local expiring memory, peer c.IP; 120/min across routes, health exempt; 429 standard envelope. | No focused credential-guessing or refresh budget; no distributed enforcement. Separate config-backed login/refresh/local-register/general limits using shared effective IP. |
| 7 Compression | Compresses normal responses. | Preserve behavior; no transport redesign. |
| Route auth/current account | Bearer JWT verification then current PostgreSQL principal; temporary admin predicate on read-only user list. | Preserve current-account authority; IP remains an abuse signal, never identity. |
| Handler Origin / parsing | Login/local register/refresh/logout POST exact trusted Origin; HttpOnly refresh cookie; strict register/self PATCH JSON, generic binder on login; no auth-specific body ceiling. | Preserve cookie/CSRF model; strict login JSON and small auth-body bound before credential work. |
| Health / missing route / errors | Health emits only status and API/database health; missing route uses Fiber error; generic infrastructure fallback, but Fiber message and wrapped validation err.Error can reach response. | Preserve health shape; safe framework/validation error text; no full envelope/catalog redesign. |
| frontend/nginx.conf | HTTP port 80 SPA fallback; same-origin /api proxy appends client-provided XFF, overwrites XFP; no SPA headers/CSP or timeout/body declarations. | Overwrite forwarding IP at this boundary; separate SPA baseline/header responsibility; no HSTS on this HTTP template or invented production topology. |

The installed Fiber v3.5.0 source, middleware internals, handlers, env examples, nginx/Docker/Compose and frontend entry/API source were inspected. Existing configuration validates production HTTPS origins but origin parsing needs canonical host/default-port and malformed-origin review. There are no password-reset routes, destructive admin CRUD or uploads to protect in this phase. No Obsidian MCP or removed planner/constitution is available; no private notes were read/persisted.

### E1. Files changed

33 unstaged files: 24 modified tracked files and nine new files, inventoried verbatim in E37. Changes cover typed HTTP configuration/examples, shared proxy/client-IP resolution, differentiated limiter, CORS/headers/recovery/parser safeguards, login/health HTTP adapters, focused backend tests, nginx/Docker serving snippets and Compose comments. Documentation updates root/backend READMEs, accepted technical decisions, backlog and this report. No dependencies, frontend application source, migrations, business schema or session storage/rotation mechanism changed.

### E2. Previous middleware order

Before routes: RequestID → Helmet security headers → CORS → Recovery → Logger → global peer-IP limiter → compression. Router then applied bearer verification/current-account lookup and temporary admin predicate, followed by handler Origin/DTO checks. Fiber parser/body limits precede this chain; the global error handler handled failures/not-found. The audit table above records purpose, assumptions and risks before implementation.

### E3. Final middleware order

Recovery → RequestID → ClientInfo → explicit API security headers → Logger → exact CORS → differentiated RateLimit → RequestSafety → compression → routes. Protected routes still verify bearer and current PostgreSQL principal before handlers/admin reads; cookie-changing handlers retain their trusted-Origin guard. Recovery is outermost to cover perimeter failures. ClientInfo precedes headers/logging/limiting so all use the same authority. CORS precedes rate/body rejection so allowed credentialed clients can read 429/413/415. Accepted preflight terminates without consuming an auth budget. The safe global error handler reapplies API headers, including parser/framework failures; no route side effect occurs on parser rejection.

### E4. Trusted proxy model

Default in every environment: no configured proxy trust, actual socket peer only. Operators may supply actual literal proxy IPs/CIDRs; there is no automatic Docker/private/loopback trust, DNS provider list or trust-all switch. Only an explicitly trusted immediate socket source permits XFF/XFP interpretation. Operators must sanitize/overwrite forwarding fields, choose narrow actual proxy sources and prevent bypass. A shared network containing untrusted clients is not a suitable blanket trusted range. No production provider or topology selected.

### E5. Client-IP resolution

ClientInfo is the sole forwarding parser; ClientIP supplies the canonical abuse key/log signal. It validates the complete XFF chain, including repeated field lines, at most 16 addresses/1024 characters. Ports/zones/empty/malformed/unspecified/multicast entries reject the chain. IPv4-mapped peers/addresses canonicalize to IPv4. Starting from a trusted peer, walk right-to-left through configured trusted hops and stop at the first untrusted address; never trust an attacker-controlled leftmost prefix beyond that hop. Missing/invalid/oversized/all-trusted chains fall back to the socket peer. X-Real-IP/Forwarded/forwarded host/alternate scheme fields are ignored. Fiber's own proxy parsing stays disabled, so Host/Scheme helpers remain socket/request-host based. IP is not account identity or authorization evidence.

### E6. Proxy-related configuration

| Setting | Default / validation |
|---|---|
| TRUSTED_PROXIES | Empty = direct peer; comma-separated literal IPs/CIDRs, masked canonical prefixes; invalid entries fail startup without echoing values. Reject wildcard/hostname/port/scoped/multicast/unspecified addresses, mapped CIDRs and /0 trust-all ranges. |
| LOGIN_RATE_LIMIT_MAX | 10, positive bounded integer |
| REFRESH_RATE_LIMIT_MAX | 60, positive bounded integer |
| REGISTER_RATE_LIMIT_MAX | 5, positive bounded integer |
| RATE_LIMIT_MAX | Retained 120 general budget |
| RATE_LIMIT_WINDOW | Retained 1m; positive duration between 1s and 1h; Fiber fixed-window accounting uses whole seconds |
| APP_IDLE_TIMEOUT | New explicit 60s, positive Go duration |

Development/production examples document these values and deliberate proxy selection. Production's existing validated HTTPS origins/secrets/database transport remain. FRONTEND_URL and ALLOWED_ORIGINS now share canonical validation and frontend membership. No networking secret or arbitrary environment variable enables forwarding trust.

### E7. Rate-limit architecture

Retain installed Fiber's concurrency-safe fixed-window limiter and default process-local expiring memory store; independent login/refresh/local-register/general buckets. Entries carry expiration and are garbage-collected by the library's one-second sweep; no custom scheduler or Redis added. Keys are only the shared effective IP, never email, body, token, cookie or frontend role. All successful/failed attempts count. Whole-second window accounting and conservative Retry-After are explicit; fixed windows may permit boundary bursts. Tests use long bounded windows and controlled socket contexts, with no wall-clock sleeps. No hard distinct-IP capacity cap or distributed denial-of-service protection is claimed.

### E8. Login rate limit

POST login: default 10 attempts per effective client IP per minute, independent of other API traffic; changing body email or spoofing forwarding fields from an untrusted peer cannot create another budget. Case/trailing-slash aliases match the protected route policy; encoded/extra-slash aliases cannot bypass it. HTTP failures for missing account/wrong password/inactive/unknown-role/lookup failure share the same generic 401 envelope and no cookie. Application password-before-status ordering remains separately asserted. Network timing equalization and exhaustive enumeration resistance are not claimed; production ingress/workload tuning remains necessary.

### E9. Refresh rate limit

POST refresh: independent default 60/minute/IP, comfortably above ordinary single-flight bootstrap/retry/reload usage. Nine sequential real-service/real-JWT adapter rotations through HTTP doubles fit the default budget, and old replay still fails without revoking the replacement. The limiter never reads/relaxes the credential or Phase 1C transaction semantics. A 429 retains cookie material for later deliberate recovery; it does not itself rotate/revoke. Phase 1D client clears memory into recoverable error rather than automatically looping.

### E10. Other rate limits

Local-only POST register: 5/minute/IP; production route remains absent. General API—including logout, protected/current-account/admin reads, unknown routes and public health's existing database probe—uses 120/minute/IP. Health no longer bypasses the budget, preventing unbounded public dependency probes; ordinary manual status/probe usage fits defaults. OPTIONS is exempt and trusted preflight terminates before the limiter. There are no password-reset or destructive admin routes; no speculative limits/routes added. An additional account-ID budget was evaluated but not added: current retained privileged routes are reads, and pre-auth IP protection avoids expensive lookup before rejection. Shared/NAT clients share budgets and operational settings must reflect measured traffic.

### E11. Rate-limit scaling limitation

Each process owns its counters; restart resets them and multiple instances do not coordinate. Fixed-window bursts, distributed source addresses and shared egress/NAT remain limits. Deploying additional instances requires a separately reviewed enforcement model; no distributed claim, Redis or infrastructure expansion. A wrongly broad proxy allowlist could make attacker input authoritative; operator trust/bypass verification is a required runtime gate, not proven by unit tests.

### E12. CORS policy

Configured origins are canonical HTTP(S) scheme://host[:port], with lowercase/canonical host and redundant default ports removed. Configuration list whitespace is trimmed; paths/trailing slash, credentials, wildcard/suffix policy, query/fragment, malformed hosts/ports, noncanonical numeric host spellings and production HTTP are rejected. Frontend must be in the explicit list. Request Origin matches exactly, with no case/path/default-port/whitespace alias, null or wildcard reflection. Disallowed supplied Origin returns safe 403 with no allow-origin/credentials headers. Originless non-browser reads may proceed; cookie mutations separately reject missing Origin.

Credentialed preflight permits GET/POST/PATCH/OPTIONS and Content-Type/Accept/Authorization/X-Request-ID only, checks requested method/header membership, sets exact ACAO plus credentials=true, Vary for Origin/method/headers and max-age=300. PUT/DELETE/TRACE/arbitrary headers are denied. Current read routes retain Fiber's automatic safe HEAD handling. Exposed response headers are X-Request-ID and Retry-After; no private-network opt-in. Credential rules follow the [Fetch standard](https://fetch.spec.whatwg.org/#cors-protocol); exact matching deliberately aligns CORS with the cookie guard.

### E13. Trusted-Origin/CSRF relationship

Phase 1D login/local register/refresh/logout remain POST-only, exact trusted-Origin guarded, host-only HttpOnly Lax-cookie endpoints, without JSON credential fallback. The guard and CORS use the same canonical configuration contract. Missing/null/untrusted Origin cannot mutate cookies; protected APIs still require bearer/current-account authority and cannot authenticate through refresh cookies. No new CSRF framework or relaxation of strict refresh single-use semantics.

### E14. Security headers

API: X-Content-Type-Options=nosniff; X-Frame-Options=DENY; Referrer-Policy=no-referrer; Permissions-Policy denying geolocation/microphone/camera. Production authoritative HTTPS adds HSTS max-age=31536000 without subdomain/preload commitments. API JSON no longer receives document CSP or blanket CORP/COEP/COOP policies that conflict with supported credentialed cross-origin API access. Headers apply to normal, denied and safe error responses. nginx document/static/error responses get the same common baseline via an included snippet; upstream common duplicates are hidden at the API proxy. Future kiosk/media permissions need explicit approval.

### E15. CSP ownership/policy

The SPA server owns document CSP; Go owns API responses. nginx baseline: default/script/font/connect self; style self plus unsafe-inline for existing React/Base UI/shadcn style attributes and chart style tags; images self/data; objects/base/frame-ancestors none; forms self. No inline/eval/external script allowance or external network dependency added. Production Vite modules/dynamic chunks and current CSS/assets use the same origin; /api/v1 proxy matches connect-src. Host Vite development does not receive this production-asset policy. External API build overrides require an explicitly reviewed matching connect-src policy; they are not silently admitted. Real UI/CSP enforcement remains Phase 1G. Inline-style allowance follows the actual retained primitives and [MDN style-src semantics](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/style-src).

### E16. HTTPS/HSTS behavior

Production browser traffic requires HTTPS; Secure refresh-cookie policy is unchanged. No TLS termination is implemented in Go and no provider is chosen. HSTS requires production plus actual socket TLS or a single exact XFP=https from an explicitly trusted immediate peer; untrusted/duplicated/ambiguous/alternate headers cannot assert it. Development/test never emits API HSTS, even with a trusted HTTPS signal. The supplied nginx listener is HTTP-only local Compose and emits no HSTS. Production's authoritative TLS edge must enforce HTTPS/HSTS and reviewed forwarding; a further ingress behind HTTP nginx requires explicit topology/real-IP review. No includeSubDomains/preload assumption. See [MDN HSTS](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Strict-Transport-Security); actual TLS/proxy enforcement remains unrun.

### E17. Request body limits

Retain explicit APP_BODY_LIMIT, default 1 MiB global safety ceiling; add 16 KiB auth POST ceiling before binding/credential work. nginx mirrors 1m global/16k auth scope. JSON mutation routes login/local register/self PATCH reject inappropriate media types with 415; login now uses the existing strict JSON binder, rejecting unknown fields and multiple documents. Bodyless refresh/logout remain bodyless cookie operations and do not bind credential bodies. The global parser rejects oversized bodies before route execution; app.Test returns ErrBodyTooLarge rather than exposing its HTTP response, tested separately from auth's real 413. Future product uploads require separately reviewed limits; no uploads implemented.

### E18. HTTP method review

Login/local register/refresh/logout are POST; self profile change is PATCH; health/current-user/list reads are GET with Fiber's safe implicit HEAD behavior. GET auth mutations return 404/405 without session creation; production register remains unmounted. CORS is a browser permission layer, not route registration or authorization. No state-changing GET/new endpoint introduced.

### E19. Server timeout configuration

Read 10s, write 15s, explicit idle 60s; all typed positive durations. Explicit read buffer 8192 bytes bounds headers; installed Fiber/fasthttp read timeout covers the full request including headers/body, with no separate public ReadHeaderTimeout in this API. nginx connect/read/send proxy timeouts are 5/20/15s. Defaults allow current small JSON flows; future uploads/large reports need reviewed budgets. Socket slow-client/timeouts and actual proxy correspondence remain Phase 1G. Installed source governs behavior; [Fiber configuration documentation](https://docs.gofiber.io/api/fiber/) provides the corresponding public controls.

### E20. Recovery/panic behavior

Outermost recovery prevents a handler/perimeter panic escaping the request. Its custom diagnostic now runs, recording panic type and source stack with safe path/method/IP/request-ID context; arbitrary panic values are excluded. HTTP receives generic 500 without stack/SQL/paths/secrets. Unknown error diagnostics record type/context, not raw err.Error. Existing logging sinks and request-ID architecture remain; no complete Phase 1F observability model introduced. Tests verify useful diagnostics and no sensitive synthetic payload in client/log fields.

### E21. Health endpoint exposure

Existing /api/v1/health shape stays exactly status plus API/database health labels, no-store; no credentials/versions/paths/raw errors. It remains a lightweight aggregate dependency check rather than a newly split liveness/readiness design. The general rate budget now covers it. Frontend's existing legacy health adapter remains unchanged; envelope/health normalization belongs to Phase 1F.

### E22. Server/framework identification

Fiber ServerHeader stays empty. nginx server_tokens off minimizes version disclosure; generic nginx identity may remain in its built-in banner/error pages. Existing FIBER_ERROR response code remains pending Phase 1F taxonomy. No cosmetic error-page redesign or server-image/dependency upgrade was introduced.

### E23. Error leakage changes

Global framework errors preserve a safe 4xx/5xx status but use generic HTTP status text rather than custom Fiber messages; invalid codes become 500. Wrapped validation errors use generic Invalid request text instead of err.Error. Unknown infrastructure/JWT/session errors retain safe current envelopes. Inactive/unknown-role login HTTP denial maps to the existing invalid-credentials envelope, while service ordering and protected-account 403 rules remain unchanged. No API-wide catalog/health/envelope rewrite.

### E24. Logging secret-safety review

No HTTP middleware/handler logs Authorization/Cookie/body/password/access/refresh/hash/signing/database secret fields. Request logger retains method/path/status/latency/request ID and adds shared effective IP through a custom tag, with no independent header parsing. Recovery/global error diagnostics exclude arbitrary values/messages that could contain credentials; server-only stack contains source locations, not HTTP response detail. Synthetic header/body/panic/error sentinel tests verify this boundary. Full structured events, correlation trust/field conventions and logging completeness/status consistency remain Phase 1F; no complete observability claim.

### E25. Proxy tests

Actual Fiber handler contexts use an explicit socket peer, not a spoofable request RemoteAddr or external service. Cases cover direct/untrusted spoofing, trusted single/multi-hop, untrusted hop stopping left-prefix spoofing, malformed/empty/port/zone/multicast/unspecified/all-trusted/bounded chains, IPv6/mapped canonicalization, repeated XFF/XFP and ignored alternate IP/host/scheme fields. Pure TLS authority is covered; literal proxy config and invalid startup values are tested. No internet/network service dependency or real production reverse proxy exercised.

### E26. Rate-limit tests

Below/over login threshold, distinct socket/forwarded keys, untrusted forwarding spoof resistance, independent refresh/register/general/health budgets, safe 429/Retry-After/no bucket counts, route alias resistance and 32 concurrent callers with exactly ten allowed/22 denied. Defaults support repeated real-service/JWT cookie rotations against HTTP repository/transaction doubles and retain replay denial. No sleeps/distributed store; expiry/GC mechanism is source-inspected, not a claim of live runtime cleanup/load measurement.

### E27. CORS tests

Explicit local HTTP and configured production HTTPS origins; exact credentialed preflight/method/header/max-age/Vary; arbitrary/null/wildcard/suffix/path/case/default-port/downgrade denial; absent-Origin non-browser read; unsupported method/header denial; canonical config/frontend membership. Existing Phase 1D Origin/cookie-only authorization tests remain. Browser enforcement stays unrun.

### E28. Header tests

Normal/error responses verify API baseline, absent document/CORP policy and production HTTPS differences. Trusted/untrusted/ambiguous scheme tests enforce HSTS only under documented authority; development/test remain usable. nginx header/CSP responsibility was source-reviewed with include inheritance/always handling; no nginx -t or real browser header test claimed. nginx's documented add_header inheritance requires repeating the shared include where CSP is added; see [nginx headers module](https://nginx.org/en/docs/http/ngx_http_headers_module.html).

### E29. HTTP safety tests

Actual auth over/exact 16 KiB, case alias ceiling, global parser rejection before handler, wrong form/text/missing media, JSON charset, malformed/multiple/unknown-field login input, bodyless refresh, panic/unknown/framework/wrapped-validation redaction, safe diagnostics, minimal health and wired timeouts/header/global bounds/banner/proxy settings. Tests preserve Phase 1B password-order/current-account and Phase 1C/D token/hash/rotation/transport assertions. No unrelated integration suite manufactured.

### E30. Frontend compatibility changes

All 106 frontend source files—including memory store, central Axios single-flight/bounded retry, bootstrap, storage cleanup, health adapter, auth tests and 62 UI primitives—are unchanged. Existing credentials/Authorization/JSON/Origin behavior fits the new CORS and body policy. Only frontend Docker/nginx serving configuration changes; its Docker API default is now /api/v1, matching the existing full Compose profile and CSP. No UI/business route/token persistence change.

### E31. Docker/nginx findings

Existing local HTTP SPA/API proxy retained. nginx overwrites XFF with its socket peer and XFP with its actual scheme; strips X-Real-IP/Forwarded/forwarded host, with explicit body/proxy timeouts, common headers and SPA CSP. Shared snippets are copied into the image and avoid unsupported newer nginx inheritance directives. Common API headers are not duplicated; CSP is scoped to document/static location. Compose only gains a proxy-trust comment: no automatic shared-network allowlist, port/network/volume change or invented production provider. An additional TLS ingress/client-IP model must be reviewed deliberately. No host nginx binary is available; syntax/runtime image/startup verification remains Phase 1G. No Docker image/container/volume action performed.

### E32. Backend fmt/vet/test results

PASS: go fmt ./..., go vet ./..., go test ./... (JSON run). **76 top-level tests + 259 subtests = 335 passing nodes across 12 tested packages**, zero failures; 17 additional top-level tests and 77 subtests relative to Phase 1D. Go 1.27.1, backend working directory, GOCACHE=/tmp/elabtrack-phase1b-go-cache. Initial default-cache attempt could not write the read-only cache; successful gates used the existing writable cache. No toolchain/dependency/gate downgrade. PostgreSQL checks unrun.

### E33. Targeted race results

PASS: go test -race ./internal/config ./internal/interface/http/middleware ./internal/interface/http/routes ./internal/bootstrap. No detected race in the exercised resolver/limiter/HTTP code and doubles, including concurrent fixed-window enforcement. This does not prove PostgreSQL multi-connection rotation, external proxies or multi-instance limiter behavior.

### E34. Frontend lint/test/build results

PASS: npm run lint (zero errors, same 19 unchanged shadcn/hook warnings), npm run test:run (**54 tests / six files**), npm run build (strict TypeScript + Vite). Measured entry JS 477.83 kB / 150.77 kB gzip; transport chunk 54.42 kB / 19.98 kB gzip; CSS 179.94 kB / 27.63 kB gzip. No warning limit/type/test rule changed, no dependencies upgraded. Application source and tests are hash-identical to Phase 1D. This host build does not prove nginx's served CSP/cookies.

### E35. Compose validation

PASS: docker compose config --quiet; docker compose --profile full config --quiet. Both repository-supported configurations validate. No Docker build/run, API/database bootstrap, live browser/proxy, migration/seed/cleanup or deployment attempted. No production manifest/provider assumed.

### E36. git diff --check and preservation review

PASS. Previous Phase 1A–1D report is an exact preserved prefix. Baseline hashes preserve frontend source, all six migrations including unrun 000003, 16 V1 audit files, dependency manifests/locks, auth application and hash/JWT/database transaction/repository infrastructure, Phase 0 report and OPEN_DECISIONS. Typed HTTP configuration changes retain prior production secret/TLS safety. No product table/workflow or Phase 1F/1G implementation introduced.

### E37. Exact git status

All 33 changes are unstaged: 24 modified tracked files and nine new files. Exact git status --short at Phase 1E closure:

```text
 M README.md
 M backend/.env.example
 M backend/.env.production.example
 M backend/README.md
 M backend/internal/bootstrap/migration_test.go
 M backend/internal/bootstrap/server.go
 M backend/internal/config/config.go
 M backend/internal/interface/http/handlers/auth_handler.go
 M backend/internal/interface/http/handlers/health_handler.go
 M backend/internal/interface/http/handlers/session_cookie.go
 M backend/internal/interface/http/middleware/cors.go
 M backend/internal/interface/http/middleware/logger.go
 M backend/internal/interface/http/middleware/rate_limit.go
 M backend/internal/interface/http/middleware/recovery.go
 M backend/internal/interface/http/middleware/security.go
 M backend/internal/interface/http/response/errors.go
 M backend/internal/interface/http/routes/auth_boundary_test.go
 M backend/internal/interface/http/routes/refresh_session_test.go
 M docker-compose.yml
 M docs/project/DECISIONS.md
 M docs/project/PHASE1_FOUNDATION.md
 M docs/project/PHASE1_SECURITY_BACKLOG.md
 M frontend/Dockerfile
 M frontend/nginx.conf
?? backend/internal/bootstrap/http_security_test.go
?? backend/internal/config/http.go
?? backend/internal/config/http_test.go
?? backend/internal/interface/http/middleware/client_ip.go
?? backend/internal/interface/http/middleware/perimeter_test.go
?? backend/internal/interface/http/middleware/request_safety.go
?? backend/internal/interface/http/routes/http_perimeter_test.go
?? frontend/nginx-api-proxy.conf
?? frontend/nginx-security-headers.conf
```

No staging, commit, push, remote modification or deployment performed.

### E38. Remaining Phase 1F work

Complete API envelope/health normalization, comprehensive error taxonomy and consistent handler mapping, structured logging architecture/field/redaction conventions, request/correlation-ID trust and propagation, logging completeness/accurate final-error status, audit/log semantics and production observability guidance. Minimum leakage corrections here do not claim those deliverables. Phase 1F has not begun.

### E39. Remaining Phase 1G live verification

Carry all Phase 1C/D requirements: isolated PostgreSQL migration 000003 up/down/reapply/invalidation, hash constraints/persistence, actual transaction rollback/commit/lost acknowledgement, separate-connection concurrent refresh, account/logout/revoke-all races, cleanup/FK/index/cadence, migration bookkeeping/exclusivity and CI reproducibility. Migration 000003 is unchanged/unrun; no new migration.

Required runtime checks: full Docker build/nginx -t/startup and SPA/API/database flow; actual HTTPS edge, source restrictions/bypass, XFF/XFP sanitization and effective-IP separation through chosen proxy chain; NAT/budget tuning and 429/preflight/retry behavior, parser/header/timeouts/health on real sockets; browser CSP/static/error headers/inline-style UI compatibility and external-origin policy; HttpOnly/Secure/Lax/path/expiry/clear cookies, login/reload/bootstrap/expiry/retry/logout/network failure and independent two-tab rotation/late cookie clearing. Unit/socket-context/jsdom tests do not prove these. No Phase 1G execution claimed.

### E40. Product/security policies still unresolved

OPEN_DECISIONS retain all statuses: provisioning/signup/eligibility, staff/admin/Super Administrator authority, verification/deactivation and ongoing-loan policy remain stakeholder work. Session families/replay-wide revocation, concurrent-session/device management, immediate access revocation, lost-response grace/cross-tab coordination remain unresolved. Local registration atomicity, migration-runner integrity, cleanup ownership and final deployment/TLS/proxy/budget tuning remain open. DEC-032–034 accept only technical perimeter contracts. No institutional IP/account lockout policy, hosting provider or private-note knowledge was invented/resolved.

### E41. Phase 1E exit gate

**SATISFIED / COMPLETE** for authorized source-level scope: explicit trusted-proxy/client-IP boundary, spoof resistance, focused auth/public/general rate controls with safe 429, exact credential-safe CORS and preserved Phase 1D CSRF/session architecture, separate API/SPA header/CSP baseline, coherent production/HTTP development HSTS responsibility, explicit body/header/timeouts and safe panic/error/log boundaries, focused tests, backend fmt/vet/tests/race, frontend lint/tests/build, all Compose configuration checks and diff check. Runtime reverse-proxy/browser/PostgreSQL/Docker evidence remains required Phase 1G as permitted by this request; this is not production-readiness proof. No Phase 1F/1G, business feature, migration execution, commit, push or deployment occurred.

## Phase 1F — API Contracts, Logging & Observability

Authorized and verified 2026-10-07 (Asia/Shanghai). Phase 0 and 1A–1E remain COMPLETE within their recorded scope. This phase is engineering foundation only. No Phase 1G/1H execution, business feature, product UI/mockup, migration/seed/cleanup execution, commit, push or deployment occurred. The preceding reports are preserved byte-for-byte as a prefix. Current permanent conventions and route inventory are in [API_CONTRACTS.md](../API_CONTRACTS.md).

### Pre-change source audit

Read project instructions/charter/source policy/decisions/open questions/roadmap, inherited foundation audit/backlog/prior reports, architecture/stack/manifests; inspected current HTTP/application/domain/infrastructure/bootstrap/config, actual installed Fiber, current auth/users/health, central client/types/hooks/Query and Docker/nginx sources. Working tree started clean. No optional plugin/private-note service or sub-agent was required.

| Actual pre-1F pattern | Retained route/path evidence | Resulting action |
|---|---|---|
| Success/message/data/meta normal envelope; logout 204 | Login/register/refresh/me/self/list HTTP handlers | Retain successes and existing list metadata; preserve bodyless logout |
| Failure outer success=false/message/data=null/meta=null, nested code/optional fields, no correlation | response/errors.go, middleware/handlers | Require nested code/message/requestId/optional fields; mirror safe outer message |
| Field validation 422, domain input 400, wrapped FieldErrors lose structure | Validator, response.Error | One 400 input policy; detect wrapped field errors and retain trusted field map |
| Raw status/services health with optional DB ping, 503 has no normal error | GET /api/v1/health | Enveloped minimal liveness; separate simple readiness using existing checker |
| Framework errors use FIBER_ERROR; default 404/405 text adapted generically | Global Fiber ErrorHandler | Stable current code/status mapping, no raw framework text |
| Direct response.Error handlers sanitize but bypass upstream diagnostics | Auth/user services/repositories | Safe outcome metadata and completion/failure logging for both direct and returned paths |
| Arbitrary incoming request IDs accepted; no standard application context propagation | Fiber request-id wrapper | Server-generated UUIDv4, private-key standard context, error/header/log propagation |
| Fiber plaintext method/raw path/status/latency/IP/ID plus separate Zap app logs | Logger/newServer | Existing Zap for structured safe route templates/status/numeric duration; render routed errors before completion |
| Recovery logs safe type/source but outer placement prevents request completion on panic | Recovery/Logger/bootstrap | Inner recovery for accurate panic completion plus retained outer fallback; any panic is generic 500 |
| Generic infrastructure errors discarded in session boundaries | Application auth, transaction/session/account adapters | Transport-independent generic InternalFailure with safe operation/type metadata, no raw credential-bearing detail |
| Frontend expects inconsistent nested error type and reads arbitrary outer messages; health uses separate fetch | API client/types/health/hooks/Query | One normalized typed failure, safe 5xx/proxy fallback, common health transport and typed Query errors |

There were no retained {error:string} or raw infrastructure-string handler successes to preserve. The relevant inconsistencies were raw health, missing IDs, FIBER_ERROR, mixed validation statuses, frontend type/message disagreement and logging completeness. No external production compatibility requirement prevents correction.

### F1. Files changed

54 unstaged entries: 45 modified tracked files and nine new files. F32 contains the exact inventory. Backend changes cover centralized mapping/envelope/correlation, safe structured request/recovery/lifecycle logging, safe cause metadata at existing auth/transaction/repository boundaries, query validation and liveness/readiness. Focused tests update only explicitly changed contracts and extend critical boundary evidence. Frontend changes normalize typed failures/Query/hooks and adapt the existing minimal status diagnostic to common transport/liveness. Documentation records contracts, accepted future UX/process direction, exact phase sequencing and remaining verification. No dependencies, migrations, product tables/workflows, mockups, UI primitives or serving/network configuration changed.

### F2. Previous API response/error patterns

The audit above is pre-change evidence. Every retained route is inventoried in API_CONTRACTS: local register, login, refresh/logout, auth/me, users/me GET/PATCH, generic admin users GET, health; not-found/method/perimeter/recovery paths. Normal successes had a consistent envelope except health. Generic failures were largely sanitized in 1E, but did not contain request IDs; field errors used 422, other validation used 400, framework errors had one broad code, and handler service failures bypassed diagnostics. The client had no requestId and incompatible nested message typing.

### F3. Final success convention

Retain `{success:true,message,data,meta}` for all JSON successes, including health/readiness. `meta:null` normally; only the already-existing user list retains its page/per_page/total/last_page metadata. No nonexistent business pagination is designed. 201 local registration remains; logout and accepted preflight are bodyless 204. Successful correlation stays in X-Request-ID rather than domain entities.

### F4. Final error envelope

`{success:false,message,data:null,meta:null,error:{code,message,requestId,fields?}}`. Nested and outer messages agree. Required safe code/message/requestId; optional map of field names to one safe message. Standard mapper/error handler/perimeter writers use it. No raw errors, SQL, pgx values, JWT internals, stack traces or credential values appear in client JSON. No raw framework 404/plain-text API error fallback in retained routed paths. Ingress-generated failures are outside API middleware and normalize to generic client failures.

### F5. Error taxonomy

Fourteen current uppercase codes: BAD_REQUEST, VALIDATION_ERROR, UNAUTHORIZED, INVALID_CREDENTIALS, TOKEN_INVALID, FORBIDDEN, NOT_FOUND, METHOD_NOT_ALLOWED, CONFLICT, PAYLOAD_TOO_LARGE, UNSUPPORTED_MEDIA_TYPE, RATE_LIMITED, INTERNAL_ERROR, SERVICE_UNAVAILABLE. Existing meaningful names retained; no Go type names or speculative business codes. API_CONTRACTS maps all current statuses/meanings. Future modules extend the backend/client catalog deliberately.

### F6. HTTP/application mapping

One response.Map uses errors.Is/As for wrapped auth/user/shared/field/framework errors. Domain/application retain transport-independent errors and standard context; no Fiber/pgx inward dependency introduced. Known absent resources 404, current duplicates/conflicts 409, login credentials 401, generic invalid session 401, authenticated authorization 403, validation 400. Internal classification wins over joined lower causes. Unknown or unsupported framework errors become generic 500. Handlers do not duplicate mapping switches. A narrow LoginError adapter preserves 1E uniform inactive/unknown-role/credential-lookup rejection while making unexpected lookup failure an ERROR event; successful credential verification followed by issuance/storage failure remains 500.

### F7. Validation behavior

400 for both malformed and semantic input, including wrapped FieldErrors. Wrong required JSON media remains 415 and body bounds remain 413. Strict JSON still rejects unknown/multiple documents; validator messages contain only code-owned rule text, no submitted values. Domain email/name/password errors get structured fields. Supplied invalid current list page/per_page/order/sort/search and unsafe offset now get 400 fields instead of silent fallback; omitted defaults and valid per_page cap 100 remain. Unknown paths are 404; no new dynamic-ID endpoint/path contract is fabricated. Frontend consumes field keys directly.

### F8. Request-ID design

One server-owned UUIDv4 per API request. All incoming X-Request-ID values are ignored, including syntactically valid IDs, malformed/oversized values. Header returned on normal/error/204 paths; error body carries same ID. No inbound trust/echo protocol, authentication or idempotency meaning. Error rendering ensures an ID even before regular middleware when possible. Existing CORS exposes the header; no CORS relaxation.

### F9. Request-ID propagation

Request middleware writes a private-key context.Context through Fiber SetContext; response mapper, completion/failure/panic logs and application/repository calls retrieve it. Existing transaction context wrapping preserves it, now asserted in actual transaction-manager double tests. No Fiber types or observability fields enter entities. Normal success context/header, error header/body/log and panic correlations are tested.

### F10. Structured logging design

One existing Zap logger/core for HTTP and application lifecycle, with production JSON default and configurable structured console locally. Removed Fiber plaintext logging. Completion middleware renders returned errors before recording routed final status; direct handler error outcomes are also diagnosed. Inner recovery permits panic 500 completion; outer recovery remains a fallback around correlation/IP/logger. Safe failure diagnostic once per direct/returned unexpected error; panic gets its own diagnostic without a duplicate request-error entry. No dependency graph rewrite or second logging framework.

### F11. Log fields

request_id, method, registered route template/unmatched, Phase 1E client_ip, final routed status, numeric duration_ms, optional error_code/security_event; diagnostics add error_class/operation, panic_type/source stack. No raw path/query/URL, headers/body/account/config. Current account IDs/emails are omitted as unnecessary. Lifecycle fields include environment/listen address, explicit migration policy, startup_seed_enabled=false, trusted_proxy_enabled and origin count. No invented build/version metadata.

### F12. Log levels

DEBUG routine input/404/405 and other ordinary 4xx diagnostics. INFO successful requests/lifecycle and expected generic login/refresh 401 events. WARN 403 and 429. ERROR internal/dependency/unexpected failure and panic; masked login lookup retains client 401 but logs ERROR. Wrong passwords are not ERROR. Production info suppresses routine 404 noise; explicit level configuration is retained. No production SLO/slow threshold invented.

### F13. Secret and PII policy

Metadata allowlists, not secret regex scanning. Never log passwords/hashes, signing material, access/refresh credentials/digests, Authorization/Cookie/Set-Cookie, DB password/full credential-bearing URL, mail/API/reset/future kiosk secrets, arbitrary body/header/config/error values. Avoid emails/names/IDs where unnecessary. Safe internal failures retain only type/operation while discarding raw detail; mapper/loggers never blindly call err.Error or zap.Error/Any. Synthetic captured-record tests cover body, query/unmatched path, bearer/cookie/response cookie, panic/error, successful access/raw refresh and startup signing/DB values. Public validation details are trusted rule messages.

### F14. Panic response/logging

Always generic 500 INTERNAL_ERROR with correlated ID, even panic carrying fiber.ErrUnauthorized. Server diagnostic has same ID, safe metadata/type and source stack without raw panic value. Stack stays server-side; completion logs actual routed 500. Panic and ordinary unexpected failure records remain testable through injectable Zap cores. Rare outer fallback contains foundational middleware failures; no live process fault/ingress proof claimed.

### F15. Auth/security logging

Safe security_event categories for executed login/refresh success/failure and logout completion/failure; generic reason via public error code/status. Successful idempotent logout means completed, not proof a particular row existed. 403/429 notable perimeter categories are visible without fabricating identity. Masked credential-lookup infrastructure failures are separately diagnosable by safe type/operation/ID. No email/token/body/user snapshot or business audit entry is emitted.

### F16. Operational logs vs business audit

Accepted DEC-037: operational debugging/security/performance/lifecycle records may be rotated/deleted under later operations policy. They are not durable product evidence. Future business audit/history/ledger design records who changed what/when under approved transactional/retention requirements, informed by the stakeholder's boss-rebuild audit emphasis. No business audit schema, retention duration or workflow is invented/implemented here.

### F17. Health/readiness

GET /api/v1/health: process liveness only, envelope data status=ok/service=elabtrack-v2, no DB access. GET /api/v1/ready: existing two-second DB checker, envelope data status=ready/service=elabtrack-v2 on success; nil/failing dependency standard 503 SERVICE_UNAVAILABLE with safe generic message. Both no-store and general-budget public reads. Bootstrap still requires initial DB connection; ping readiness does not prove schema/migration/privilege/business correctness. Handler accepts a tiny standard-context interface instead of infrastructure type. Offline doubles prove routing/data/context contracts; the existing two-second checker timeout is source-inspected. Live PG readiness remains 1G.

### F18. Frontend normalized error

ApiResponse is a typed success/failure union. ApiRequestError has numeric status, code, safe message, optional requestId/fields; no retained Axios config/body/header/raw error object. Central adapter validates known codes/UUID/rule field shapes, rejects arbitrary legacy/proxy bodies, prefers a valid response-header reference and substitutes generic 5xx text. Network status=0/NETWORK_ERROR; malformed/unknown error fallback is safe. No component string parsing.

### F19. Frontend handling changes

Auth/user hooks consume typed safe error messages; TanStack default error type is registered. Query never retries 4xx, including 429; server/network reads retain one bounded ordinary retry and mutations none. Existing Phase 1D refresh single flight, endpoint exclusions, one auth retry, generation guards, queued cookie operations and memory-only state remain. apiErrorMessage optionally adds support reference for unexpected failures, never routine field failures; no global toast spam or product UI. Health dynamically imports the same singleton transport with attachAuth=false/retryAuth=false/withCredentials=false and schema-checks minimal liveness. The existing status placeholder says liveness explicitly, without implementing a feature screen.

### F20. 404/405

Fiber unknown routes and practical method errors map to NOT_FOUND/METHOD_NOT_ALLOWED standard envelopes, generic text and ID. No custom unsupported-method matrix or catch-all overrides that bypass auth. Absent production registration remains contained; existing global perimeter rejection may occur first. GET auth mutations remain unhandled 404/405 and never mutate sessions. SPA 404 remains frontend routing, separate from API.

### F21. 429

Existing limiter still uses independent effective-IP budgets and safe conservative Retry-After/no-store. It automatically gains nested message/requestId via shared Fail. Client normalizes 429 and does not Query-retry it; Phase 1D session recovery remains deliberate. Proxy/counter/window/NAT/process scope is unchanged, with real ingress/IP/429 tuning still 1G.

### F22. Error-envelope tests

Actual bootstrap pipeline tests malformed JSON, structured validation, unauthorized/forbidden, unknown path, 405, rate limit, direct/returned unexpected errors and panics. Assert status/code/safe mirrored message, UUID header/body, fields and no internals. Real retained route tests exercise auth/me/self/admin, local validation/current query/name/domain failure, success/meta/204; existing account/refresh/Origin/storage tests remain. Map tests cover 34 known/framework/infrastructure/joined cases directly and wrapped, including pgx error detail redaction.

### F23. Request-ID tests

Absent ID generated; safe valid-looking incoming ID replaced by policy; malformed/oversized IDs ignored safely; distinct UUIDv4s, header/body/log agreement, application repository context and transaction context preservation. Normal success correlation, direct error, returned error, panic and readiness check all covered without listeners/external services.

### F24. Redaction tests

Captured Zap records assert obvious synthetic secret values absent from body/header/query/unmatched-route/Set-Cookie/error/panic/startup/auth-success paths, while effective IP and ID remain. Safe InternalFailure tests preserve classification/type/operation without invoking a raw credential Error method. Existing repository tests still reject raw token/hash exception detail; errors.Is assertions accommodate new safe metadata wrappers without dropping rollback/locking/secret assertions.

### F25. Panic tests

String panic and framework-error panic both return 500 INTERNAL_ERROR, generic text and same UUID in diagnostics/completion. Diagnostic has source stack/type, never value; client lacks stack/source/secret. Accurate numeric status/duration and event severity are asserted via observer fields, not terminal text. Direct vs returned failure diagnostic counts prevent silent bypass/duplicates.

### F26. Frontend tests

73 passing tests across seven files (19 added). New tests exercise real central transport for 400/401/403/404/405/409/429/500/503, structured fields, request-ID header/body fallback, malformed/legacy/unknown payload safety, generic 5xx/support reference, network/unexpected failure and Query policy. All previous memory/bootstrap/single-flight/bounded retry/logout/account tests remain passing. Existing status tests now exercise common transport options/liveness/schema/manual retry. No real browser cookie jar/PG/nginx behavior claimed.

### F27. Backend fmt/vet/tests

PASS: go fmt ./..., go vet ./..., go test ./... via JSON run. **83 top-level tests + 321 subtests = 404 passing nodes across 14 tested packages**, zero failures; seven new top-level tests and 62 subtests relative to Phase 1E. Commands from backend/ with Go 1.27.1 and GOCACHE=/tmp/elabtrack-phase1b-go-cache. No dependency/toolchain/gate changes. Initial root-directory fmt invocation was outside the module and failed; corrected backend command passes. Initial fixture assertions were adapted to intentional 400/UUID/metadata changes; no security/concurrency/type gates removed.

### F28. Race checks

PASS: go test -race for internal/application/auth, domain/shared, bootstrap, interface/http/middleware, interface/http/response, interface/http/routes, infrastructure/database and infrastructure/persistence/postgres. After the final logger level adjustment, affected bootstrap/middleware/routes race checks were repeated and passed. This is local code/test-double evidence, not real PostgreSQL concurrency or cross-tab proof.

### F29. Frontend lint/test/build

PASS: npm run lint, zero errors and same **19 unchanged** shadcn/hook warnings; npm run test:run, **73 tests/seven files**; npm run build, strict TypeScript + Vite. Node 24.19.0 / npm 11.17.0. Entry JS 459.35 kB / 146.55 kB gzip; lazy transport 54.07 kB / 19.87 kB gzip; CSS 179.94 kB / 27.63 kB gzip. Dynamic transport boundary retained, no chunk threshold/lint/type/test weakening. Initial root npm invocation lacked package.json; the correct frontend invocation passes. A new control-regex warning was eliminated in source, retaining character validation without rule suppression.

### F30. Compose validation

PASS for default and full profile using the installed Docker Desktop Windows Compose CLI: docker-compose.exe -f docker-compose.yml config --quiet; docker-compose.exe --profile full -f docker-compose.yml config --quiet. Current WSL docker/docker-compose wrappers report disabled/unavailable integration and sandbox Windows interop initially fails. Authorized read-only execution outside the sandbox successfully validated both configurations. No integration settings changed, alternate binary downloaded, service started, image built, container/volume touched or DB connected. Compose/nginx/Docker/env source is unchanged. This is configuration evidence only.

### F31. Diff and preservation

PASS git diff --check, including final documentation. Baseline hashes prove **111 protected files unchanged**: all six migrations, 16 V1 audit documents, all 62 UI primitives, manifests/locks/configuration, Phase 0/open-policy registers, JWT/hash/migrator/session-cookie/Origin/client-IP/CORS/rate/request-safety/header controls, memory/store/storage/bootstrap and Docker/nginx/Compose serving configuration. Prior 1A–1E foundation text is an exact preserved prefix. Repository SQL and transaction semantics remain unchanged; only safe diagnostic wrapping changed at those boundaries. Read-only remotes still point to the verified eLabTrack repository; no staging/remote mutation.

### F32. Exact Git status

All changes are unstaged; no commit/push/remote change. Exact git status --short at closure:

```text
 M README.md
 M backend/README.md
 M backend/internal/application/auth/login.go
 M backend/internal/application/auth/principal.go
 M backend/internal/application/auth/refresh.go
 M backend/internal/application/auth/refresh_test.go
 M backend/internal/bootstrap/app.go
 M backend/internal/bootstrap/http_security_test.go
 M backend/internal/bootstrap/infrastructure.go
 M backend/internal/bootstrap/migration_test.go
 M backend/internal/bootstrap/server.go
 M backend/internal/infrastructure/database/transaction.go
 M backend/internal/infrastructure/database/transaction_test.go
 M backend/internal/infrastructure/persistence/postgres/auth_repository.go
 M backend/internal/infrastructure/persistence/postgres/auth_repository_test.go
 M backend/internal/infrastructure/persistence/postgres/user_repository.go
 M backend/internal/interface/http/handlers/auth_handler.go
 M backend/internal/interface/http/handlers/health_handler.go
 M backend/internal/interface/http/handlers/user_handler.go
 M backend/internal/interface/http/middleware/logger.go
 M backend/internal/interface/http/middleware/perimeter_test.go
 M backend/internal/interface/http/middleware/recovery.go
 M backend/internal/interface/http/middleware/requestid.go
 M backend/internal/interface/http/response/errors.go
 M backend/internal/interface/http/response/response.go
 M backend/internal/interface/http/routes/http_perimeter_test.go
 M backend/internal/interface/http/routes/refresh_session_test.go
 M backend/internal/interface/http/routes/routes.go
 M backend/internal/shared/constants/constants.go
 M backend/internal/shared/pagination/pagination.go
 M backend/internal/shared/validator/validator.go
 M docs/ARCHITECTURE.md
 M docs/STACK.md
 M docs/project/DECISIONS.md
 M docs/project/PHASE1_FOUNDATION.md
 M docs/project/PHASE1_SECURITY_BACKLOG.md
 M docs/project/ROADMAP.md
 M frontend/src/app/query-client.ts
 M frontend/src/features/auth/hooks/use-auth.ts
 M frontend/src/features/foundation/api/health.api.ts
 M frontend/src/features/foundation/foundation.test.tsx
 M frontend/src/features/foundation/pages/status-page.tsx
 M frontend/src/features/users/hooks/use-users.ts
 M frontend/src/lib/api-client.ts
 M frontend/src/types/api.ts
?? backend/internal/bootstrap/contracts_test.go
?? backend/internal/domain/shared/failure.go
?? backend/internal/domain/shared/failure_test.go
?? backend/internal/interface/http/response/errors_test.go
?? backend/internal/interface/http/routes/contracts_test.go
?? backend/internal/shared/observability/
?? docs/API_CONTRACTS.md
?? frontend/src/lib/api-error.test.ts
?? frontend/src/lib/api-error.ts
```

### F33. Mobile-first decision documentation

Accepted future direction in DEC-038: borrower/customer mobile composition first, progressively enhanced tablet/desktop; kiosk touch-first, large targets, dedicated shell and privacy/session-reset awareness; staff/admin desktop/tablet-first with practical small-screen responsiveness. This does not settle borrower eligibility, staff/admin authority or kiosk authentication policy. No mockup/UI feature created.

### F34. shadcn decision documentation

DEC-038 reaffirms React/TypeScript/Vite, Tailwind v4, shadcn/ui and Lucide as the eLabTrack UI foundation. All 62 primitives/components.json unchanged. DEC-039 requires approved designs implemented using the shadcn-based system in 3B; no replacement framework or default demo.

### F35. Phase 3A roadmap

ROADMAP now explicitly sequences 1F contracts/observability → 1G real integration → 1H runner hardening → 2 domain/database → 3A UX architecture/wireframes/high-fidelity mockups → 3B shadcn design system/application shell/reusable domain components. 3A prioritizes mobile borrower views/adaptations, kiosk-specific and staff/admin layouts; subsequent feature UI must compare against approved mockups/responsive behavior before completion. No later phase authorized or begun.

### F36. Remaining Phase 1G live verification

Required isolated runtime evidence, carried forward without closure from unit tests:

- PostgreSQL bootstrap, actual migrations 000001–000003 up/down/reapply/invalidation, constraints/hash persistence/presented-digest rejection, rollback/error/commit/lost acknowledgement; actual DB verify-full TLS/CA/host behavior and restricted release/runtime credentials where supported.
- Separate-connection/process transactional refresh/concurrent use at most one successor, account status/role locks, logout/revoke-all/rotation races, expiry clock predicates, cleanup/FK/index/retention and operational batch ownership/cadence. Local registration atomicity/recovery remains unresolved account work.
- Docker image/service startup, nginx syntax/proxy/static/error behavior and SPA/API/database flow; chosen HTTPS edge/trust/source/bypass, effective-IP separation, forwarded-header sanitization/NAT budgets/preflight/429/Retry-After, body/header/timeouts on real sockets.
- Browser login/HttpOnly/Secure/Lax/host/path/expiry/clear cookies, reload/bootstrap/expiry/single-flight/bounded retry/logout acknowledgement, recoverable/lost network responses, independent two-tab rotation/loser clearing winner cookie; no cross-tab coordination/grace policy claimed.
- Actual CORS/CSP/styles/assets/API/error headers on chosen serving path; liveness/readiness transitions and ping-only limits; request-ID header/body/context/proxy/browser correlation, structured console/JSON/log output/redaction, server lifecycle and log collection/access/rotation.
- Fiber's pre-handler parser lifecycle can traverse middleware before final parser mapping; App.Test reports oversized global body parse errors instead of an HTTP413. Real socket parser final response/status/correlation/log ordering and ingress-generated error boundaries need explicit 1G evidence. Routed final-status tests do not prove this lifecycle; no framework parser redesign performed here.
- Reproducible integration/CI gates and recovery observations. Current runner defects must be recorded in 1G, with fixes explicitly owned by 1H.

No live PostgreSQL/browser/nginx/Docker/container/production check is claimed in Phase 1F.

### F37. Phase 1H runner issue

Stakeholder's boss-rebuild audit identifies the shared runner weakness; current source independently shows execTx commits DDL before separate schema_migrations INSERT/DELETE, no checksum validation or advisory/concurrent-runner lock, and lastAppliedVersion treats every Scan failure as no migrations. Record **Phase 1H—Migration Runner Hardening**, before heavy reliance on business-domain migrations. Runner and all SQL are hash-identical and unrun here. Phase 1G observes current integration limitations; it does not automatically authorize these fixes.

### F38. Product-policy decisions unresolved

OPEN_DECISIONS unchanged: signup/provisioning, borrower types/eligibility, staff/admin/Super Administrator scope, verification/deactivation/ongoing-loan consequences and all inventory/loan/fine/history/kiosk/notification policies retain their prior status. Session families/device management/concurrent-session/global issuance/logout-all/immediate access revocation/grace recovery/cross-tab policy remain unresolved. DEC-035–037 accept engineering contracts, DEC-038/039 explicit future UX direction, DEC-040 sequencing only; none resolve institutional rules.

### F39. Remaining Phase 1 security backlog

SEC-015/016 source-level contracts/logging implemented, runtime still 1G. SEC-017 offline coverage expanded, real integration/CI unverified. SEC-008 partial with hardening 1H. SEC-004/006 stronger session/registration policy limits; SEC-018 cleanup scheduling/ownership/load verification; SEC-014 serving/HTTPS/CSP runtime; SEC-019 known warnings/dependency/image review remain. Production deployment topology/proxy budgets/log access-retention/operational ownership still require explicit operational decisions. No Sentry/Datadog/New Relic/OpenTelemetry/Prometheus/Grafana/cloud SDK/provider integrated; existing Zap core and standard context remain vendor-neutral extension points.

### F40. Phase 1F exit gate and security review

**SATISFIED / COMPLETE** for the authorized source/offline Phase 1F scope: consistent retained JSON/error/field contract and centralized safe mapping, correlated server-owned UUIDs/context, structured safe metadata/log levels/lifecycle/auth/error/panic records, tested secret exclusions, explicit minimal health/readiness, typed normalized frontend failure and preserved 1D refresh lifecycle / 1E perimeter, accepted UX/mockup/shadcn documentation and exact sequencing, backend fmt/vet/tests/race, frontend lint/tests/build, both supported Compose configurations and diff check all pass. All 40 requested final-report items are recorded above. Live parser/PG/browser/proxy/Docker/logging verification remains1G; runner hardening remains1H. This is not production deployment approval. No later phase, business feature, UI/mockup, migration execution, staging/commit/push or deployment occurred.

## Phase 1G — Real Integration Verification

Authorized and verified **2026-10-07 (Asia/Shanghai)**. Phase 0/1A–1F reports above
are preserved as an exact prefix. This is real disposable local verification,
with two narrow HTTP foundation fixes. No Phase 1H, business workflow/table,
mockup/UI redesign, V1 Firebase, production data, commit, push or deployment.
Reproduction: [integration instructions](../../integration/README.md).
Sanitized machine results: [evidence](../../integration/evidence/2026-10-07.json).
The following G1–G60 cover every requested final-report item.

### G1. Files changed

Modified: root/backend READMEs; bootstrap/server.go; middleware/logger.go;
response/errors.go and errors_test.go; docs/API_CONTRACTS.md; this report;
PHASE1_SECURITY_BACKLOG.md; ROADMAP.md. Added: bootstrap/parser_logging_test.go;
middleware/parser_logging.go; tests/integration/foundation_test.go; integration
README, Compose override, env runner, runtime/browser/restart/log check scripts,
browser-only driver/build script and sanitized evidence JSON. G56 records exact
porcelain status. No frontend production source, dependency, migration, seed SQL,
domain/application/auth implementation or normal Docker/nginx manifest changed.

### G2. Environment/tool versions

Ubuntu 24.04.1 under WSL; Go **1.27.1 linux/amd64**, GOTOOLCHAIN=auto;
Node **24.19.0**, npm **11.17.0**; Docker Engine client/server **29.8.0**,
Desktop **4.92.0**, Compose **5.5.1**; PostgreSQL **18.6**;
external cached Playwright **1.55.0**, real Chromium **140.0.7339.16**.
Native WSL Docker wrappers reported integration unavailable; installed Windows
Docker/Compose executables worked through interop after starting the stopped
local Desktop engine. Real local sockets/Windows interop required execution
outside the sandbox. Chromium initially lacked libnspr4/libnss3/libasound2t64;
Ubuntu packages were downloaded/extracted under /tmp, without system installation
or application dependency changes. No toolchain requirement was weakened.

### G3. Disposable database setup

Fresh `elabtrack_v2_integration`, PostgreSQL service/container
`elabtrack_v2_phase1g_postgres`, localhost **15432**; inside Compose `postgres:5432`.
Project `elabtrack_v2_phase1g`, volume `elabtrack_v2_phase1g_pgdata`, dedicated
project network, API port **18080**, nginx **15173**. Generated credentials were
in ignored 0600 backend/.env.phase1g; normal backend env_file was reset. Resource
inspection found no pre-existing eLabTrack container/volume. Base project fixed
names require the override in addition to `-p`. Normal development data was never
reused. All owned integration resources and credentials were removed (G23).

### G4. PostgreSQL result

PASS: repository `postgres:18.6-alpine` pulled/started; pg_isready health passed;
server reported PostgreSQL 18.6 x86_64 Linux-musl. Host pgx/psql and container psql
authenticated to the isolated database; full API used that same instance.

### G5. Migration from-zero result

PASS: explicit compiled repository CLI `--migrate-status` showed all three pending;
`--migrate-up` applied all three, without API startup. Final live schema has users,
refresh_tokens and schema_migrations; digest-only session columns, unique hash,
ownership/replacement FKs, lifetime/hash/replacement/self-link checks and expected
indexes were inspected. PostgreSQL retains the historical NOT NULL constraint
name `refresh_tokens_token_not_null` after column rename; it refers to token_hash,
not a remaining plaintext token column.

### G6. Migration versions applied

Exactly `000001_create_users`, `000002_create_refresh_tokens`,
`000003_refresh_session_security`, in order. Live status and tracking rows agreed.
All six repository SQL files are unchanged.

### G7. Down/up result

PASS: fresh up → explicit latest `--migrate-down` removed version 000003 and
restored its prior empty raw-column schema → `--migrate-up` reapplied 000003.
Final status/schema are hash-only. Both directions intentionally invalidate
sessions; this disposable sequence preceded authentication fixtures. No API
was run against the old schema.

### G8. Migration-runner observations

Known defects reproduced safely; no runner fix:

- A temporary tracking CHECK rejected probe version 000099 after its SQL committed;
  probe table existed with **no tracking row**. Operator cleanup was required.
- Two concurrent probe runners: **one exit 0, one exit 1**, one applied record;
  both attempted the same pending DDL. No advisory lock exists.
- Editing an applied temporary SQL file was silently skipped; no checksum detected it.
- A temporary login allowed schema access/create but denied tracking SELECT;
  `--migrate-down` exited **0**, saying “no migrations to roll back” despite three
  applied versions. Source lastAppliedVersion swallows every Scan error.

Temporary probe table/constraint/version/login were removed. Source also confirms
Down commits SQL before deleting tracking. No dangerous production failure,
business migration or runner redesign was attempted.

### G9. Seeding safety result

PASS: explicit development `--seed` executed repository development.sql locally;
normal startup never seeds. Production disabled-TLS configuration failed before
IO; separately, valid production verify-full configuration plus `--seed` reached
and failed the development-only guard before database IO. Seed passwords were
never printed or copied into this report. Seed accounts are synthetic temporary
foundation conveniences, not approved institutional roles/provisioning.

### G10. API startup result

PASS: real repository API image built and started against migrated PostgreSQL.
Startup JSON records explicit-command migration policy, startup_seed_enabled=false,
development, listener and allowlist count without credentials. API restarts did
not change migration versions/timestamps or invoke seeds. No test endpoint added.

### G11. Health result

PASS: real GET /api/v1/health returns 200 standard envelope with minimal
status=ok/service=elabtrack-v2, no-store and server UUID. It remained 200 while the
isolated PostgreSQL container was stopped.

### G12. Readiness result

PASS: /ready returns 200 when connected; stopping only the owned database produced
safe enveloped **503 SERVICE_UNAVAILABLE** with matching request ID. Raw database
errors were absent. After PostgreSQL resumed, the running API pool reconnected and
readiness became 200. The probe proves connectivity, not full schema/operational
health or backup readiness.

### G13. Synthetic account setup

Repository user adapter fixtures used generated synthetic addresses/passwords,
real bcrypt, temporary user/admin and is_active=false. Browser fixtures used the
supported local-only register API; no new provisioning UI/route. Only disposable
fixture role/status updates used direct parameterized PostgreSQL. Successful
tests deleted their own accounts; earlier harness-run remnants were disposed with
the owned volume. No real personal/school data.

### G14. Real login result

PASS: actual HTTP/database login returns access_token/Bearer/expiry/safe user;
no raw refresh in JSON. Cookie is host-only elabtrack_v2_refresh, HttpOnly,
SameSite=Lax, Path=/api/v1/auth, Secure=false for documented development HTTP.
Persisted session lookup matches the digest; no raw credential is printed.

### G15. Invalid-login result

PASS: wrong password, unknown account and inactive account all return the same
generic **401 INVALID_CREDENTIALS** message and standard ID/envelope. Login rate
budget applies before credential work. No claim of timing equalization is made.

### G16. Current-account result

PASS: real /auth/me returns exactly five safe current-account fields, expected
synthetic ID/current role/status and standard envelope. Password/hash/session
material is absent. PostgreSQL, rather than historical JWT authority, governs.

### G17. Stale-role/status result

PASS, mandatory: temporary admin initially lists users; direct disposable
demotion then old access token → **403** on list. Disabling that account then old
token → **403** on both list/current-account reads. No token reissue was necessary
for current database authority to win. Final institutional roles remain unresolved.

### G18. Refresh hash persistence result

PASS: SHA-256 of the presented canonical opaque credential locates its row;
the raw credential does not match persisted values. Final schema has no token
column. Submitting the stored digest itself as a credential returns 401.

### G19. Real rotation result

PASS: cookie-driven refresh creates a different opaque credential/digest,
consumes the old row, records its successor FK and returns renewed access. Old
credential replay gives generic 401 and does not revoke the winner's successor;
successor rotates successfully. Access JWT serialization can be identical when
issued for the same account within one second; credential rotation/state and
issuance are verified without inventing an access-token uniqueness requirement.

### G20. Concurrent refresh result

PASS, mandatory: **12 separate concurrent real HTTP requests** submitted the same
cookie against PostgreSQL. **One 200, eleven 401s**; exactly one unrevoked/unexpired
replacement for that account, old row revoked with exactly one successor link.
Winner remained usable after old replay. Repeated under the Go race detector.

### G21. Rotation rollback/atomicity evidence

PASS: successful real HTTP rotation has both old consumption and replacement.
Integration-only decorator invokes the real repository Consume, then returns a
controlled error through the existing application port. Real TxManager rolls back
both actual writes; no pair returned, one original usable row remains and the old
credential subsequently rotates. No production hook/backdoor. Lost commit
acknowledgement remains the documented re-login ambiguity, not simulated here.

### G22. Logout/revocation result

PASS: real logout returns bodyless 204, revokes the presented session, expires/
clears the cookie; repeated logout is safe, later refresh is 401. Actual browser
memory clears immediately. Existing access JWT remains valid to expiry if current
account permits it; this phase does not invent immediate per-session access revocation.

### G23. Cleanup/retention result

PASS: synthetic stale expired rows deleted in bounded batches **2 then 1**;
active sessions preserved. Explicit --sessions-cleanup succeeded with zero further
eligible rows. No scheduler added. Final Compose down -v --rmi local removed only
the owned three containers, network, fresh volume and built project images;
filtered Docker inspection found none remaining. Ignored env, private sentinel
file and generated browser build removed. Shared downloaded base-image caches
and safe result/log artifacts under /tmp remain; unrelated hello-world resources
and normal development volume were not touched.

### G24. Full Docker stack result

PASS: normal API/frontend Dockerfiles built successfully, isolated full profile
ran PostgreSQL/API/nginx. Compose ps showed PG healthy and API/frontend Up;
actual API readiness and browser document/API responses established reachability.
API/frontend have no Docker HEALTHCHECK, so “Up” alone was not called healthy.
No worker/service added; explicit migrations preceded startup.

### G25. nginx/SPA result

PASS: nginx -t passed; normal / and /status load with SPA route fallback,
client navigation/status button calls real API, built assets use relative /api/v1
through nginx, cookies rotate through that same path. No browser-localhost API
URL or product UI change. The additional browser driver was copied only into
the disposable container and removed with it.

### G26. Browser login result

PASS in real Chromium: actual centralized ApiClient authenticate/login plus
current-account call populated the actual non-persisted store. Browser fixtures
are synthetic; the driver imports the real app/client/store, not a reimplementation.

### G27. Browser token-storage result

PASS: programmatically inspected localStorage/sessionStorage contents contain no
auth token/key or tested secret; document.cookie cannot read refresh cookie;
Chromium cookie jar reports HttpOnly/Lax/path/local Secure=false. Access is present
in application memory only. Preference persistence/theme remains functional.

### G28. Reload/bootstrap result

PASS: driver observed initially empty memory after reload, network counted exactly
**one** real bootstrap refresh, authenticated/current-account memory restored and
cookie changed. StrictMode produced no duplicate refresh or loop. Same behavior
passed after full Compose restart.

### G29. Invalid-session result

PASS: revoke only the synthetic account's active refresh rows, reload → one 401
bootstrap, unauthenticated memory and cleared dead cookie, no infinite loop or
sensitive UI error.

### G30. Bounded-retry result

PASS: protected HTTP boundary mock gives first 401, real refresh then one actual
successful retry. Persistent protected 401 mock gives exactly two protected
attempts/one actual refresh, then memory clears/stops. Only protected responses
were mocked; cookie/API/database refresh stayed real and server security unchanged.

### G31. Single-flight result

PASS: invalidated only test-page memory access, issued five concurrent protected
requests → real 401s → **one refresh per tab**, all five resumed successfully.
Network counts validate actual Axios/store behavior beyond unit adapters.

### G32. Cross-tab result

REQUIRED experiment completed: A logged in, B restored shared cookie into separate
memory. Held two actual browser requests at the network boundary until both sent
the same cookie: one refresh succeeded, the other received 401. One tab became
unauthenticated and the late denial **cleared the shared winner cookie** in the
observed runs. Winning tab initially retained authenticated memory. Logout in A
left B apparently authenticated with still-valid access; B's next forced refresh
failed and cleared its memory. No duplicate DB successor or privilege escalation.

### G33. Cross-tab recommendation

**COORDINATION REQUIRED — B. ADD CROSS-TAB SESSION COORDINATION BEFORE PRODUCT UI.**
Smallest later approach: serialize cookie-changing refresh/login/logout across
same-origin tabs, e.g. Web Locks with BroadcastChannel lifecycle notifications;
notify logout immediately, avoid losing-tab refresh races, and deliberately restore
each tab's own memory after a peer operation. Design fallback/late-operation and
logout generation handling before implementation. Do not broadcast or persist
usable access/refresh tokens in localStorage. BroadcastChannel notification alone
is not an atomic refresh lock. No coordination/replay-grace implementation or new
Accepted product/technical policy was introduced by this recommendation.

### G34. Real CORS result

PASS: locally controlled allowed http://localhost:14173 made credentialed real
browser health/login/refresh requests to nginx:15173; JSON-header preflight worked,
exact ACAO and credentials=true, never wildcard. Disallowed localhost:14174 fetch
was browser-blocked. Chromium CDP wire metadata proved real 403 even though
Playwright suppresses normal response events for CORS-blocked replies.

### G35. Trusted-Origin/CSRF result

PASS: disallowed local browser Origin on cookie-bearing refresh/logout → 403,
no ACAO/credentials headers and shared cookie unchanged. Subsequent allowed
refresh succeeded. Tests remained local; exact Origin policy was not relaxed.

### G36. CSP/security-header result

PASS: real normal nginx document has intended self-script CSP, DENY framing,
nosniff, no-referrer and denied camera/microphone/geolocation. Built JS/CSS,
retained shadcn card/button, theme, navigation and status all worked; **zero CSP
console violations**. API headers and real safe error/envelope paths also pass.
Parser final 413/431 now correlate correctly. Ingress-generated errors remain
nginx responses as documented, not claimed to be API-generated JSON.

### G37. HSTS/local HTTPS observation

PASS: local HTTP API/document emit no HSTS, including direct/proxied spoofed
X-Forwarded-Proto=https. Production HTTPS/HSTS ownership/configuration retains
prior tests; no fake local production TLS or hosted endpoint was invented.
Actual production TLS edge validation remains deployment-only.

### G38. Proxy/client-IP result

Only observed nginx **172.18.0.4/32** trusted. Direct API source = **172.18.0.1**;
direct spoofed XFF=198.51.100.42 still logs 172.18.0.1. nginx overwrites spoofed
forwarding; proxied effective client also =172.18.0.1 rather than nginx .4.
This Docker Desktop NAT collapses tested clients to one address; end-user
separation/cloud bypass policy is not proven. No blanket private-network trust.

### G39. Real rate-limit result

PASS: configured maxima **10 login / 60 refresh / 120 general**, five-second
integration windows. Each maximum passes its budget (invalid auth safely fails
normally); next request is **429**, standard envelope/UUID, Retry-After=5 and
no secret. General test uses live health. Origin OPTIONS is not a credential
bypass; canonical effective IP comes from the verified boundary. Defaults remain
per-minute in normal configuration; process-local/NAT/restart limits remain.

### G40. Request-ID correlation result

PASS: **243 real browser/runtime responses** matched exactly one final completion
status/ID in captured structured output. Real HTTP test checks error-body/header
ID equality. Absent and valid-looking/malformed/2048-character inbound IDs were
replaced by server UUIDs, per DEC-036; no caller ID acceptance is claimed.

### G41. Log-redaction result

PASS: **454 captured structured records**, **91 distinct** tested sentinel strings,
**zero matches**. Includes synthetic plaintext/bcrypt passwords, access credentials,
opaque refresh/digests, JWT key, DB/probe passwords and credential-header sentinels.
Forbidden payload/header/URL/secret fields absent. Real login success/failure,
refresh success/failure, logout, 401/403/429, startup, readiness and restart records
were captured. Source stack/type metadata remains server-side. This is tested
local output, not an assertion about future external sinks/retention/access.

### G42. Panic-runtime result or reason unrun

No runtime panic injected: production exposes no safe test-only panic seam.
Existing captured observer/actual-router panic/internal-error tests and race
checks pass. No production-accessible panic route was added for this phase.

### G43. Responsive foundation smoke result

PASS: actual normal SPA at **390×844**, **768×1024**, **1440×900**; no horizontal
overflow beyond one-pixel tolerance. Dark/light toggle, shadcn primitives and
service button/navigation function. Foundation smoke only; no mockups or Phase 3
business UI QA. Accepted mobile borrower/touch kiosk/staff responsive policy unchanged.

### G44. DB TLS separation result

PASS: documented local non-TLS development connection used; production disable
configuration rejected before connection. Prior effective verify-full/no-fallback
tests remain enabled. Real production certificate/hostname handshake unavailable
and deferred to legitimate hosting; no invented TLS endpoint.

### G45. DB role/privilege observations

GAP documented: migration and runtime currently share local postgres credentials;
current_user postgres, rolsuper=true. No accepted least-privilege release/runtime
roles are implemented. The temporary restricted fault-fixture login was deleted
and does not solve this gap. Carry separation/ownership into later hardening.

### G46. Clean-start reproducibility

PASS: from zero volume/config → repository PG → explicit compiled CLI migrations
→ normal API/nginx builds/start → supported synthetic registration/login. No
reclone needed. New integration README exposes every additional setup: fixed-name
override, ignored generated credentials, cached/external Playwright and missing
libraries, Windows CLI fallback, actual nginx /32 discovery, test-only browser
build/copy, readiness waiting, log capture and owned cleanup. No hidden manual
schema/bootstrap operation or migration/seed startup side effect.

### G47. Docker restart behavior

PASS with repository `unless-stopped` policy retained: full Compose restart,
PostgreSQL account/data and migration versions/timestamps unchanged, API/SPA
ready again, old access/current account works, browser cookie survives and reload
rotates/restores; logout works. Compose restart may start API before PostgreSQL
readiness, causing a safe failed connection startup that Docker retries. The
integration harness waits for actual readiness; it does not infer it from process
start or add an application retry policy.

### G48. Integration defects discovered

1. **INTEGRATION DEFECT FOUND:** actual body-parser 413/header-parser rejection
   completion logs incorrectly reported **200**, preceding final error rendering.
   **ROOT CAUSE:** Fiber 3.5 socket parser callback traverses USE middleware first,
   then renders the parser error after Logger returns.
2. **INTEGRATION DEFECT FOUND:** over-8192-byte request header returned **500**.
   **ROOT CAUSE:** central response mapper omitted Fiber's 431 error and used the
   internal-error fallback.

Known migrator defects are reproduced in G8, unchanged for Phase 1H. Harness-only
issues corrected: fixture pool closed before cleanup; missing browser libraries;
normal Playwright response event unavailable for CORS-blocked response; an initial
integration restart=no override suppressed the existing retry policy. These are
distinct from the two production HTTP defects and are not hidden exit failures.

### G49. Integration fixes made

1. **FIX:** listener hook wraps the existing Fiber/fasthttp parser callback;
   suppress provisional middleware completion, log once after final render.
   Retain existing mapping/perimeter traversal and restore the already-emitted
   server ID after Fiber clears its standard context on release.
   **REGRESSION TEST:** real callback body/header cases plus actual socket
   response/ID/final-status log assertions.
2. **FIX:** safe **431 BAD_REQUEST**, “Request headers are too large.”; no raw parser
   text. **REGRESSION TEST:** direct/wrapped map tests, real parser callback and
   actual socket 431/correlation checks. API catalog updated.

No migration/auth/session-policy/UI redesign or production test backdoor.

### G50. Integration tests added

Opt-in Go live suite: health/ready/current account, generic invalid login, stale
role/disable, digest persistence/presented-hash rejection, 12-way real HTTP refresh,
replay/successor/logout, actual post-consumption rollback and bounded retention.
Six subtests plus parent pass. External Playwright scripts exercise actual client
memory/cookies/retry/single-flight/tabs/CORS/CSP/responsive/restart. Runtime script
probes real migration errors, config guards, parser/socket limits, dependency
outage, rates and IDs; log script verifies captured evidence. Normal unit gate
skips live tests unless explicitly enabled; no production services are implicit.

### G51. Backend fmt/vet/test results

PASS final `go fmt ./...`, `go vet ./...`, `go test ./...` from backend with
GOCACHE=/tmp/elabtrack-phase1b-go-cache. **84 top-level + 324 subtests = 408 passing
offline nodes across 14 packages**; live opt-in suite explicitly skipped in this
ordinary run and passed separately. Early path/cache command errors were corrected;
final recorded gate exits are zero. Dependencies/types/test assertions unchanged.

### G52. Targeted race results

PASS: application/auth, bootstrap, HTTP middleware/routes/response,
infrastructure database/persistence. Also **real opt-in PostgreSQL/HTTP suite under
-race**, six subtests plus parent, including separate-connection refresh race and
real rollback. No race detected in these exercised paths.

### G53. Frontend lint/test/build results

PASS: lint zero errors with **same 19** retained primitive/hook warnings;
**73 tests / seven files**; strict TypeScript/Vite normal build. Entry JS
459.35 kB/146.55 gzip, lazy transport 54.07/19.87, CSS 179.94/27.63.
Production frontend source/dependencies/primitives unchanged. Actual critical
browser flows were rerun against the rebuilt API after the HTTP fixes.

### G54. Compose validation

PASS: normal default, normal full profile and isolated full override config
--quiet (three configurations), with installed Windows Compose executable.
Actual base Dockerfiles built; nginx -t, full start/restart and browser/API/DB
runtime also passed. The override preserves original restart behavior.

### G55. git diff --check

PASS after implementation/documentation/cleanup. Prior foundation is an exact
prefix; migrations, auth/session/client/primitive code, manifests/locks, V1 audit,
Phase 0, DECISIONS and OPEN_DECISIONS are unchanged. No tracked credential artifact.

### G56. Exact git status

Pre-test status/diff were clean on `main`; fetch/push origin unchanged at
https://github.com/Maaku050/elabtrack-v2.git. Normal Compose names before tests:
services postgres/backend/frontend, fixed elabtrack_v2_postgres/backend/frontend,
volume elabtrack_v2_pgdata. No staging/commit/push/fetch/remote mutation.
Final `git status --short`:

```text
 M README.md
 M backend/README.md
 M backend/internal/bootstrap/server.go
 M backend/internal/interface/http/middleware/logger.go
 M backend/internal/interface/http/response/errors.go
 M backend/internal/interface/http/response/errors_test.go
 M docs/API_CONTRACTS.md
 M docs/project/PHASE1_FOUNDATION.md
 M docs/project/PHASE1_SECURITY_BACKLOG.md
 M docs/project/ROADMAP.md
?? backend/internal/bootstrap/parser_logging_test.go
?? backend/internal/interface/http/middleware/parser_logging.go
?? backend/tests/integration/
?? integration/
```

### G57. Deployment-only checks remaining

Actual production HTTPS/TLS termination/HSTS, legitimate hosted PostgreSQL
verify-full certificate/hostname handshake, cloud proxy/bypass topology, production
backup/restore and external monitoring delivery remain unrun. No provider/secret/
production endpoint invented. Credential privilege separation, operational log/
cleanup ownership and CI are also unresolved foundation/operations work, not
disguised as completed deployment checks.

### G58. Phase 1H migration-runner inputs

Unchanged runner reads paired files, executes each file in one transaction, then
separately INSERTs/DELETEs version tracking. Fresh three-version up/latest-down/up
passed; controlled orphan DDL, concurrent runner failure, missing edited-content
detection and swallowed permission error demonstrate integrity/recovery gaps.
No lock/checksum tracking; local migration/runtime both superuser. Hardening needs
atomic SQL/bookkeeping, exclusive advisory locking, reviewed checksum compatibility,
correct no-row-versus-error handling, constrained release/runtime credentials and
recovery tests. **Phase 1H has not begun and requires separate authorization.**

### G59. Product-policy decisions still unresolved

OPEN-001–020 unchanged: provisioning/signup, borrower types/eligibility,
staff/admin/Super Administrator, verification/deactivation/ongoing loans,
inventory/lifecycle/fines/history/terms/kiosk/email/migration policies remain
Needs Stakeholder Input/Deferred. Families/global revocation/concurrent devices/
immediate access revocation/grace are unresolved; cross-tab engineering
recommendation is evidence for a later task, not institutional policy acceptance.

### G60. Phase 1G exit gate

**SATISFIED / COMPLETE** for the authorized local integration scope: real
PostgreSQL/migrations/up-down-up/API/health/readiness/login/hash-only persistence,
concurrent atomic refresh/rollback/current-account stale authority/logout/cleanup;
actual full Docker/nginx/Chromium cookie/memory/reload/retry/single-flight/CORS/
Origin/CSP/responsive/proxy/rate/ID/log/restart checks; explicit adverse cross-tab
evaluation/recommendation; narrow defects fixed with regression evidence; required
standard/race/Compose/diff gates pass and owned resources/secrets cleaned.
Cross-tab coordination must be addressed before product UI, and Phase 1H runner
hardening remains separately authorized. Completion does not mean production
readiness or resolution of product policy. No later phase, business feature,
mockup, commit, push or deployment occurred.

## Phase 1H — Migration Runner & Database Privilege Hardening

Authorized and locally verified **2026-10-07 (Asia/Shanghai)**. Earlier A–G reports
are preserved byte for byte as a prefix. This phase changes infrastructure only:
no business migration/feature, Phase 1I implementation, mockup, V1 access, commit,
push or deployment. [Reproduction](../../integration/README.md) and
[sanitized evidence](../../integration/evidence/2026-10-07-phase1h.json).
The following H1–H57 cover all requested closure-report items.

### H1. Files changed

Runner/legacy compatibility/SQL transaction guard and unit/live tests; typed migration configuration and tests; CLI/action tests and Make wrappers; local env examples, Compose/tools job and bootstrap/grant scripts; runtime-identity auth tests; integration environment/role/SQL/process/browser/runtime/restart/log helpers and README/evidence; root/backend READMEs and foundation/decision/roadmap/security-backlog docs. H54 gives exact Git status; no historical migration, production frontend, dependency/lockfile, domain/application/auth source or OPEN_DECISIONS change.

Exact changed-file inventory (37 files):

```text
.env.example
Makefile
README.md
backend/.env.example
backend/.env.production.example
backend/Makefile
backend/README.md
backend/cmd/api/main.go
backend/cmd/api/main_test.go
backend/database/bootstrap-roles.sh
backend/database/runtime-grants.sql
backend/internal/config/config.go
backend/internal/config/config_test.go
backend/internal/config/database.go
backend/internal/config/migration_test.go
backend/internal/infrastructure/database/migration_legacy.go
backend/internal/infrastructure/database/migration_sql.go
backend/internal/infrastructure/database/migrator.go
backend/internal/infrastructure/database/migrator_integration_test.go
backend/internal/infrastructure/database/migrator_test.go
backend/tests/integration/foundation_test.go
docker-compose.yml
docs/project/DECISIONS.md
docs/project/PHASE1_FOUNDATION.md
docs/project/PHASE1_SECURITY_BACKLOG.md
docs/project/ROADMAP.md
integration/README.md
integration/browser.mjs
integration/compose.phase1h.yml
integration/env-run.py
integration/evidence/2026-10-07-phase1h.json
integration/logs.py
integration/migration-processes.py
integration/restart.mjs
integration/roles.py
integration/runtime.py
integration/sql.py
```

### H2. Previous migration lifecycle

Source audit before implementation: Up → ensureTable (pool CREATE IF NOT EXISTS) → read applied version map → discover .up.sql suffixes/lexical full-name sort → for each pending read file → pool Begin → tx Exec SQL → tx Commit → separate pool INSERT tracking. Down → ensureTable → pool SELECT version ORDER BY version DESC LIMIT 1 → any Scan error treated as no applied version → read matching .down.sql → tx Exec/Commit → separate pool DELETE tracking. Applied-map Query/Scan/Rows.Err errors propagated, but lastAppliedVersion swallowed all failures. Status only tested table existence/read versions and printed pending/applied. Create chose highest parsed numeric prefix +1, fell back to 1 on discovery errors, and wrote pairs with overwrite-capable WriteFile. No filename/pair/duplicate/history/checksum/lock validation.

### H3. Previous atomicity defect

execTx committed migration DDL before applyOne INSERT or Down DELETE. A failed bookkeeping write could leave committed orphan DDL; a failed removal could retain an applied record after rollback SQL committed. Phase 1G reproduced orphan DDL and swallowed denied lookup, concurrent DDL attempts and undetected editing; those historical observations remain unchanged.

### H4. Final migration lifecycle

Discover/validate and load both files → acquire dedicated PostgreSQL connection → immediate native advisory try-lock → inspect/create tracking as appropriate → read all records and verify applied-prefix/checksums → per-file BEGIN → SQL → bookkeeping in same tx → COMMIT → connection close releases command lock. Up applies ordered pending suffix; Down reverses the highest applied file only. Earlier successful files remain applied if a later file fails; the whole batch is intentionally not one transaction. Only Up creates missing tracking, under the lock. Status/Down never create it.

### H5. schema_migrations schema changes

Fresh public.schema_migrations: version TEXT PRIMARY KEY with canonical-name CHECK; applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW(); checksum TEXT NOT NULL with 64 lowercase hex CHECK. Numeric file prefix is the canonical ordering key; full canonical basename remains the compatible stored identifier. Legacy adoption adds checksum/digest constraint while retaining existing version PK/timestamps. No name/dirty/progress payload or business migration.

### H6. Checksum algorithm/policy

SHA-256 over marker elabtrack-v2-migration-pair-v1 plus NUL, then big-endian uint64 up byte length/up raw bytes, then down byte length/down raw bytes. Both directions are protected because safe rollback is part of this paired model. Comments/line endings count. Files are loaded once, so applied SQL and stored digest describe the same bytes. Store on apply; verify before Up/Down/Status; never replace an applied digest.

### H7. Existing-record compatibility

Checksum-free records cause an actionable failure. Explicit development-only --migrate-adopt-legacy is an operator provenance/schema attestation, not historical proof: accept only a prefix of exact 000001–000003 foundation identifiers whose paired bytes match the embedded Phase 1G baseline digests. Unknown/duplicate versions, baseline edits, missing mapping, gaps, query/scan failures and already-checksummed tracking fail. One transaction adds/backfills/constrains checksum without changing applied_at. Real known-three-version adoption and unknown-history/repeated-adoption refusals pass. Fresh disposable volumes are preferred; no external database is auto-repaired or assigned fabricated checksums.

### H8. Up atomicity behavior

SQL and INSERT(version,checksum) execute on the same dedicated session/transaction. SQL or INSERT failure rolls back both; commit is the only success point. Fresh tracking may remain empty after a first migration failure, which is coherent unapplied state. Parameterized tracking writes.

### H9. Down atomicity behavior

Down SQL and DELETE matching version/checksum share one tx. Require exactly one removed row; SQL/delete/commit errors fail. Real partial DROP-then-error and delete-trigger failure restore the table and applied record. Existing 000003 keeps its session-invalidation behavior.

### H10. Nontransactional migration policy

Supported files are transaction-safe SQL only. Current six SQL files meet this policy. CREATE INDEX CONCURRENTLY in a fixture fails clearly under PostgreSQL SQLSTATE 25001 and rolls back prior fixture DDL. Lexical guard rejects explicit transaction/session control (including BEGIN/COMMIT/END/ABORT, SET/RESET and prepared statements), understanding comments/quotes/dollar bodies; PostgreSQL remains the parser. No exceptional mode now; a future real need requires reviewed dirty/recovery semantics.

### H11. Advisory lock implementation

Project key 0x454c41424d494752 (ELABMIGR), parameterized pg_try_advisory_lock on a pool-acquired connection hijacked into a dedicated session. Acquire before every history read/change; hold across all file transactions, Up/Down/Status/adoption. Close on every exit with an independent bounded cleanup context; never return a session lock to the pool. Real different connections and separate CLI processes verify exclusion.

### H12. Lock timeout/failure behavior

Immediate safe failure if another migrator holds the lock; no blocking advisory wait. Connection acquisition/lock-query context is five seconds. Closing the owner session releases the lock; a following status/down succeeds. Separate-process loser and concurrent status both exit 1 with a retry-after-finish reason, without credentials.

### H13. Migration discovery/order validation

Nonzero six-digit numeric versions with lowercase underscore-separated names and required regular up/down files. Reject invalid .sql filenames, non-regular SQL files, missing pairs and duplicate numeric mappings even with different names. Sort numerically; gaps allowed; applied records must be a repository prefix. Generator validates existing discovery, chooses highest +1, uses an exclusive local scaffold lock/O_EXCL and never overwrites pairs. Interrupted partial generation needs operator review; no business scaffold was generated.

### H14. Missing migration handling

Applied record absent from discovered pairs causes a hard failure before migration SQL. Missing one counterpart fails discovery. Actual temporary applied version 000099 with both files removed is refused. Unknown database identifiers are not echoed into logs/errors.

### H15. Checksum mismatch handling

Hard failure with safe canonical local version before any SQL/history change. Both up and down fixture tampering fail Status/Up/Down. Restoring original bytes passes status, proving stored checksum was not overwritten.

### H16. Database-error handling

Separate missing table from empty applied records; Query, Scan, iteration and metadata/connection errors become safe classified MigrationError. No arbitrary error becomes version zero. Denied runtime history read fails; pgx can surface a server query error during rows iteration and that path is also checked. Closed-pool Down/Status fail. Unknown driver/SQL details never escape CLI.

### H17. Dirty/partial-state policy

No dirty flag: supported PostgreSQL SQL and bookkeeping roll back together. Strict checksums/history reject inconsistent external state. Connection loss at COMMIT can make acknowledgement ambiguous while database state remains atomic; error explicitly requires status inspection before retry. Unsupported exceptional nontransactional execution would require a later recovery design.

### H18. Migration CLI changes

Retain mutually exclusive --migrate-up/down/status/create, --seed and --sessions-cleanup. Add explicit development-only --migrate-adopt-legacy. CLI errors preserve safe reason/version; stdout status/generation and structured events contain no SQL/config/URLs. Migration/seed select operator connection; cleanup selects runtime. Normal API and Docker ENTRYPOINT remain API-only. Local tools-profile job is an explicit action excluded from default/full startup.

### H19. Migration status behavior

Locked coherent read only; missing-table notice, canonical version, applied/pending, checksum=verified/not-applied and current version (none or latest). It creates no tracking table and refuses corrupt/legacy/concurrently changing state. No admin UI.

### H20. Migration logging

Existing Zap conventions: migration.command and migration.apply, canonical migration_version, direction, result, numeric duration_ms, and allowlisted command failure reason. No raw errors/SQL/config/password/URL. Final separate-process capture: six structured migration records, six SQL/config/credential sentinels, zero matches. API capture: 327 structured rows, 70 distinct sentinels with zero matches and 242 response/status/ID correlations; machine evidence contains actual counts.

### H21. Migration DB configuration

Optional typed *DBConfig MigrationDB from MIGRATION_DATABASE_URL parsed through the existing constrained URL/TLS validator. It must target the runtime host/port/name. No environment reads in migrator. CLI selection is explicit. Separate URL CA setting is supplied in sslrootcert; no runtime TLS downgrade or URL logging.

### H22. Runtime DB configuration

Retain validated DATABASE_URL or DB_* and pool settings. API bootstrap/readiness/repositories and session cleanup use cfg.DB only. API startup does not need MigrationDB. Pin standard_conforming_strings=on with the existing allowlisted pgx runtime parameters; the ambient-PG test now asserts both safe parameters and preserves override rejection.

### H23. Development credential behavior

Fresh Compose creates separate fixed infrastructure logins by default; examples contain clearly published local-only values. API gets runtime DML credentials, explicit tools job gets the local migrator connection, admin stays in PostgreSQL. Absent migration URL may use cfg.DB only in development/test as a documented convenience. Integration uses generated distinct ignored 0600 credentials. Existing volumes retain passwords/owners and are never silently adopted.

### H24. Production credential behavior

API requires only valid runtime settings. Production migration actions require MIGRATION_DATABASE_URL, different username and the same target; missing URL cannot fall back. Both connections require explicit verify-full TLS, TLS >=1.2, hostname/chain verification and nil fallback. Configuration cannot prove deployed role privileges; actual hosted roles/CA handshake still require deployment checks. Production/test legacy adoption is refused before IO.

### H25. PostgreSQL role structure

Bootstrap admin creates DB/citext and local logins; elabtrack_migrator and elabtrack_runtime are both NOSUPERUSER/NOCREATEDB/NOCREATEROLE/NOREPLICATION/NOBYPASSRLS. No role membership between them. Roles are infrastructure only, unrelated to provisional user/admin product fields. Real current_user/catalog checks verify both and actual API runtime identity.

### H26. Schema ownership strategy

Stable public schema/application-table/tracking owner elabtrack_migrator. Database/citext bootstrap administrator retains extension responsibility; migrator has no database CREATE or cluster administration. Historical CREATE EXTENSION IF NOT EXISTS works with preinstalled citext. No random ownership transfer; future migrations run as that owner. Existing volume transfers require separate administrator review.

### H27. Runtime grants

CONNECT/database plus USAGE/public, SELECT/INSERT/UPDATE/DELETE on named users/refresh_tokens only. No grant option, schema/database CREATE/TEMP or owner membership. UUID IDs require no sequences; standard PostgreSQL/citext type/function access is retained, no custom security-definer surface. runtime-grants.sql is an explicit owner step after current foundation migrations; future approved table migrations grant only their named required objects transactionally, never blanket defaults.

### H28. schema_migrations privileges

Owned by migrator, explicit REVOKE ALL FROM runtime/PUBLIC; no runtime read or write dependency. Real runtime SELECT/UPDATE refused. Bootstrap source creates only runtime pool/ping; actual API has no migration URL/password/bootstrap password, and startup/reconnect/restart leave version/checksum/timestamps unchanged.

### H29. API startup under runtime credentials

PASS: final repository image starts with DB_USER=elabtrack_runtime and only its generated password. Safe Docker inspection verifies migration/admin secrets absent without printing Env. /ready returns 200 via runtime ping. Runtime lacks schema CREATE/tracking access; actual auth and browser work. No startup schema/adoption/grant command or test endpoint.

### H30. Migration run under migrator credentials

PASS: explicit compiled CLI from zero using elabtrack_migrator, plus explicit Compose tools job on current state. Non-superuser role owns schema/tables/tracking; required DDL succeeds. API/auth tests use elabtrack_runtime, not owner/admin credentials.

### H31. Atomic up-failure test

PASS live/race: temp paired 000099 executes CREATE TABLE then division-by-zero. Probe table absent, applied row absent, prior three versions unchanged. No permanent broken SQL file.

### H32. Bookkeeping-failure test

PASS live/race: temporary tracking CHECK rejects probe version after valid CREATE within tx; probe DDL and row absent. Also a temporary tracking delete trigger rejects Down bookkeeping; table and record retained. Fixture constraints/triggers/functions are removed.

### H33. Down-failure test

PASS live/race: applied fixture Down drops probe table then errors; rollback restores table and keeps checksum/applied record. Probe cleaned using owner credentials afterward.

### H34. Checksum tampering test

PASS live/race: separately edit applied temporary up and down bytes. Status/Up/Down all fail; restoring bytes verifies original tracking. Repository SQL untouched.

### H35. Missing-applied-migration test

PASS live/race: apply probe, remove both paired fixture files, Up refuses absent applied history. No automatic delete/repair; fixture table/row cleaned explicitly.

### H36. Duplicate-version test

PASS offline and live: two names with identical numeric 000003 are rejected before SQL. Additional invalid filename/zero/missing counterpart and ordering-gap tests pass.

### H37. DB-error test

PASS live: runtime tracking permission failure causes Down/Status error instead of no migrations. PASS offline closed pool and safe failure/log tests. Query/scan/iteration source paths all preserve errors; no permissive no-row fallback remains.

### H38. Concurrent migrator test

PASS mandatory real PostgreSQL: live different connections plus separate compiled CLI processes. Observer sees native pg_locks key; winner applies one table/record, second Up and simultaneous Status exit 1 with lock contention. Winning connection closes; status/down subsequently succeeds and history is coherent. No simultaneous schema mutation or in-memory lock claim.

### H39. Runtime DDL denial test

PASS live/race: runtime INSERT/UPDATE/SELECT/DELETE in a rolled-back synthetic-account tx; CREATE TABLE, ALTER disposable owner probe, CREATE ROLE, tracking SELECT/UPDATE denied. GRANT without grant option may yield only warning; effective PUBLIC privileges remain absent and runtime lacks grant option. No destructive project-table modification.

### H40. Migrator DDL test

PASS live/race: actual non-superuser owner CREATE/ALTER/DROP disposable table succeeds in a rolled-back transaction. Catalog checks confirm no cluster privileges, and all foundation/tracking ownership is stable.

### H41. Existing migrations from-zero result

PASS: exactly 000001_create_users, 000002_create_refresh_tokens, 000003_refresh_session_security apply in numeric order from a fresh PostgreSQL 18.6 volume; all checksum statuses verified. No changes to six SQL files or 000003 session-hash constraints/security intent. No additional permanent migration.

### H42. Up/down/up result

PASS: all three Up → three explicit latest Down operations → current version none → all three Up, then named runtime grants. Final hash-only schema/three verified records. Both 000003 directions intentionally invalidate sessions; the round trip preceded auth fixtures.

### H43. API/auth regression result

PASS: real runtime-role Go HTTP/repository suite under -race, six subtests plus parent; current-account demotion/disable, hash persistence/replay, 12-way one-success/11-denial refresh, actual post-consumption rollback, logout and bounded cleanup. Chromium memory-only/HttpOnly/reload/single-flight/bounded retry/CORS/CSP/three viewports pass; restart preserves DB/history/browser restoration. Runtime parser 413/431, readiness 503→200, rates/proxy/IDs/logs pass. Known cross-tab defect deliberately remains.

### H44. TLS safety result

PASS offline migration effective TLS configuration/username/target/fallback tests and actual production disable rejection before IO. Test startup permits no migration secret; production CLI refuses fallback. No real hosted verify-full certificate/hostname handshake performed; explicitly remains deployment-only.

### H45. Seed safety result

PASS: explicit development --seed selects migrator and executes existing synthetic seed SQL; production valid-config guard refuses before IO, test guard retained. API startup never seeds. Explicit --sessions-cleanup succeeds with runtime credentials (zero further eligible rows). No new seed accounts/policy or credentials in reports.

### H46. integration/README changes

Complete safe fresh-volume/bootstrap roles, separate ignored credentials, owner migrate/status/down/up/adoption, named runtime grants, runtime API readiness/auth, live fault/process tests, browser/runtime/restart/log checks and owned cleanup. Existing browser driver/artifact names are retained only for harness reuse; target is Phase 1H. Commands are deterministic later-CI inputs, not a CI implementation.

### H47. Phase 1I cross-tab status

Required before Phase 2/product UI; ROADMAP/backlog/DEC-044 name Cross-Tab Session Coordination explicitly. Actual repeated adverse race: winner refresh then late losing 401 clears shared new cookie; peer logout leaves other memory authenticated until refresh failure. No 1I code, usable-token broadcast/persistence, replay relaxation or institutional session policy.

### H48. CI/deployment checks still open

No CI pipeline. Hosted PostgreSQL role/ownership setup and verify-full CA/hostname handshake, real HTTPS/HSTS/CSP edge/proxy bypass, authorized release/backup/restore/recovery, external logging/sink/access/retention, cleanup operational ownership/cadence and measured workloads remain open. Local evidence is not production readiness or permission to deploy.

### H49. Backend fmt/vet/test results

PASS final go fmt ./..., go vet ./..., go test ./... using declared Go 1.27.1 and writable GOCACHE=/tmp/elabtrack-phase1b-go-cache. 425 passing offline nodes: 92 top-level +333 subtests across 14 packages; live suites skip without opt-in. Real migration database package has 39 passing nodes including live parent/15 cases; real auth suite has seven. No dependency/test/type weakening. Early working-directory/helper syntax and two PostgreSQL-semantics harness assertions were corrected; final recorded gates all exit zero.

### H50. Targeted race results

PASS auth application, bootstrap, HTTP middleware/routes/response, database and PostgreSQL persistence; additional CLI/config race checks. Actual live migration database suite and runtime-role HTTP/auth suite pass under -race, including connection/process/session races and rollback. No detected race in exercised paths.

### H51. Frontend lint/test/build results

PASS: lint zero errors/same 19 existing warnings; 73 tests/seven files; strict TypeScript/Vite build. Unchanged normal entry 459.35 kB/146.55 gzip, lazy transport 54.07/19.87 and CSS 179.94/27.63. All frontend production source/primitives/manifests/lockfiles unchanged. Only disposable ignored browser harness build emits expected static-import/chunk warning.

### H52. Compose validation

PASS six configurations: base default/full/tools/full+tools, retained Phase 1G full override and Phase 1H full+tools override, config --quiet. API/frontend/final API and explicit CLI job images build; actual PG/API/nginx and nginx -t pass. Documented container grant command passes. Isolated fixed names/project/volume only; normal development volume untouched. Owned resources/secrets are removed at closure; shared base-image cache/unrelated hello-world resource retained.

### H53. git diff --check

PASS final. A–G report prefix, six historical SQL files, frontend production source/primitives, manifests/locks, V1 audit, Phase 0 and OPEN_DECISIONS preserved against pre-H hashes. No tracked secret/generated browser/Python artifact. No git staging/commit/push/fetch/remote mutation.

### H54. Exact git status

Pre-H tree was clean on main; origin fetch/push remains https://github.com/Maaku050/elabtrack-v2.git. No staging/commit/push/remote mutation. Final git status --short:

```text
 M .env.example
 M Makefile
 M README.md
 M backend/.env.example
 M backend/.env.production.example
 M backend/Makefile
 M backend/README.md
 M backend/cmd/api/main.go
 M backend/cmd/api/main_test.go
 M backend/internal/config/config.go
 M backend/internal/config/config_test.go
 M backend/internal/config/database.go
 M backend/internal/infrastructure/database/migrator.go
 M backend/tests/integration/foundation_test.go
 M docker-compose.yml
 M docs/project/DECISIONS.md
 M docs/project/PHASE1_FOUNDATION.md
 M docs/project/PHASE1_SECURITY_BACKLOG.md
 M docs/project/ROADMAP.md
 M integration/README.md
 M integration/browser.mjs
 M integration/env-run.py
 M integration/logs.py
 M integration/restart.mjs
 M integration/runtime.py
?? backend/database/
?? backend/internal/config/migration_test.go
?? backend/internal/infrastructure/database/migration_legacy.go
?? backend/internal/infrastructure/database/migration_sql.go
?? backend/internal/infrastructure/database/migrator_integration_test.go
?? backend/internal/infrastructure/database/migrator_test.go
?? integration/compose.phase1h.yml
?? integration/evidence/2026-10-07-phase1h.json
?? integration/migration-processes.py
?? integration/roles.py
?? integration/sql.py
```

### H55. Remaining product-policy decisions

OPEN-001–020 unchanged: provisioning/eligibility, staff/admin/Super Administrator, verification/inactive/ongoing loans, stock/lifecycle/fines/settlement/history/terms/kiosk/email/V1 migration remain stakeholder/deferred work. Infrastructure DB roles establish no product permissions. Session family/global/immediate-access revocation/concurrent device/grace policy remains unresolved.

### H56. Remaining Phase 1 security/foundation backlog

Required 1I cross-tab coordination; SEC-004/006 local registration atomicity and unresolved session/institutional policy; SEC-017 automated CI; SEC-014 actual hosted HTTPS/CSP/TLS/edge; SEC-015 external logging retention/access; SEC-018 cleanup schedule/owner/load; SEC-019 retained warnings/dependency/image review. SEC-008 runner and separated-role source/local findings are verified resolved; deployed operator controls still unverified.

### H57. Phase 1H exit gate

**SATISFIED / COMPLETE for the authorized source and disposable local scope.** Atomic SQL/bookkeeping, native advisory exclusion, immutable checksums/history validation, safe state errors, coherent all-file round trip, real process concurrency, local least-privilege ownership/API/auth, explicit CLI, live/standard/race/frontend/Compose/diff gates pass. No later phase, business schema/feature, mockup, commit/push or deployment. Completion does not close production/CI/product-policy checks or authorize Phase 1I.


## Phase 1I — Cross-Tab Session Coordination

Authorized 2026-10-07. No Phase 2, business features/migrations, mockups, CI pipeline, commit, push or deployment.

### Baseline recorded before implementation

The unchanged Phase 1G Chromium harness was executed on a fresh isolated Phase 1H PostgreSQL 18.6/API/nginx stack using its separated runtime/migrator roles. All tracked file hashes still matched the pre-work baseline when this evidence was saved. Two authenticated tabs sent the same cookie behind a real browser network barrier: one refresh returned 200 and one 401; the shared replacement cookie was absent afterward; one tab remained authenticated and the other became unauthenticated. After a fresh login and peer reload, logout cleared only the initiating tab; the peer remained authenticated and its existing access still authorized `/auth/me` until forced recovery failed. See [baseline evidence](../../integration/evidence/2026-10-07-phase1i-baseline.json).

Current root cause: per-document single-flight cannot serialize shared-cookie mutations across documents. `AuthHandler.Refresh` clears the cookie for every `ErrTokenInvalid`; that application error collapses missing/malformed/unknown/expired/revoked/consumed credentials and disabled/deleted/invalid-role accounts into the same generic 401. A late losing response consequently deletes a replacement established by the winner. Unexpected transaction/database failures return safe 500 without clearing; denied Origin returns 403 before credential work and without clearing. Logout explicitly clears its presented browser session even on storage failure. The HTTP response cannot know which cookie the browser holds when that response arrives. This baseline is recorded before any implementation change.

### Final verification — 2026-10-08

Confirmed implementation and recorded disposable local evidence. Source/application-domain/migrations and prior A–H report remain preserved. [Final machine evidence](../../integration/evidence/2026-10-08-phase1i.json) and [reproduction guide](../../integration/README.md#phase-1i-multi-tab-verification) accompany this complete 58-item report.

### I1. Files changed

23 unstaged files: 18 modified and five new. Backend: handler and three HTTP test files plus live integration test. Frontend: coordinator/code/tests, centralized client/tests, query client, store/tests. Integration: driver, browser runner, session-browser helper, log verifier, guide and two safe evidence files. Four project documents updated. Exact paths appear in I51.

### I2. Phase 1G defect reproduced

YES, before implementation: actual Chromium 140.0.7339.16 on fresh PostgreSQL 18.6/API/nginx, separated roles, identical presented cookie behind a browser barrier, one 200/one 401, winner cookie deleted, final tabs authenticated/unauthenticated. Peer logout left the other tab authenticated and its old access usable. The baseline above and checked-in baseline artifact were saved before code changes.

### I3. Root cause

Independent document promises cannot serialize a shared cookie; a late replay rejection blindly expires a successor installed by a winner. Logout also had no peer-memory/cache notification. The server single-use decision was correct.

### I4. Previous refresh-cookie clearing policy

Missing, malformed, unknown, expired, explicitly revoked, consumed/replayed and invalid current account all collapse to ErrTokenInvalid, generic 401 and unconditional cookie deletion. Missing/deleted/disabled/unrecognized account state is included. Unexpected storage/transaction failure returns safe 500 without deletion; denied Origin returns 403 before IO without deletion. Logout explicitly deletes even when storage revocation fails.

### I5. Final refresh-cookie clearing policy

NO failed refresh emits Set-Cookie, including definitive invalidation. A request cannot identify the cookie present when its response arrives; even expired/revoked/account-invalid input can race a later login. Invalid HttpOnly cookies remain until expiry, replacement or explicit trusted-Origin logout; server denial remains authoritative. Success issues a successor and explicit coordinated logout revokes/clears as before. No error taxonomy/replay relaxation or new cookie/schema.

### I6. Coordination technology

Dedicated non-React SessionCoordinator: BroadcastChannel for lifecycle/outcome metadata, native exclusive Web Locks for atomic origin/storage-partition ownership. Same lock protects refresh/login/local-register/logout. No new dependency. DEC-045–047 record confirmed engineering choices; their primary API contracts are linked there.

### I7. Tab identity model

crypto.randomUUID per document/application runtime, never persisted and never used as authentication/account/session identity. HMR creates a fresh runtime after disposing the old coordinator.

### I8. Attempt/epoch model

Fresh UUID per mutation/authority check. Epoch is a monotonic microsecond-scale high-resolution clock value advanced past observed epochs; timestamps are milliseconds. Equal epochs use deterministic attempt ordering and event rank; completed success rejects contradictory same-attempt failure. Local store generation and sent-token checks separately fence asynchronous HTTP.

### I9. Cross-tab message types

Version 1 strict union: refresh-started, refresh-succeeded, refresh-failed (allowlisted failure only), logout, session-invalidated, auth-state-changed. Fields only type/version/tabId/attemptId/epoch/timestamp and optional failure. Exact-key, UUID, numeric/age/type checks reject malformed/unknown/extra fields; there is no account object or authority.

### I10. Secret-sharing analysis

No access/refresh/hash/Authorization/Cookie/password/reset value is sent. Refresh replies are stored only in their own document. Unit extra-field/credential-sentinel tests and real observed channel payload checks pass. No message logs or arbitrary payload forwarding.

### I11. Same-tab single-flight relationship

Existing shared refreshPromise, bootstrapPromise and local mutation queue remain. Eight simultaneous unit protected 401s and five real browser requests share one refresh per document. Local and cross-document ownership are separate layers.

### I12. Cross-tab refresh ownership

Acquire named exclusive native lock, announce refresh-started, send one ordinary HttpOnly-cookie request, accept only current generation, store own access and publish outcome. The browser mutex, rather than channel timing or localStorage, establishes ownership.

### I13. Waiting-tab recovery

Wait for native ownership, independent of receiving a completion message. Every waiter rechecks its own generation before HTTP and accepts only its own trusted server result. Peer terminal events cancel stale queued recovery through the generation fence.

### I14. Serialized follow-up refresh

A completes and the browser installs its successor before releasing the lock; B then rotates that successor for its own memory token; C repeats. Actual two/three-tab tests require distinct consumed credentials and persisted replacement links.

### I15. Leader lease/timeout

Active Axios HTTP bound 15s, queued lock wait 20s, valid metadata age 60s, future tolerance 5s. COORDINATION_TIMEOUT is a safe typed recoverable status-0 error. Never steal a live lock or launch an overlapping mutation on timeout; deliberately retry once the owner resumes/releases/closes. These are implementation bounds, not product/session policy.

### I16. Owner crash behavior

Native closure/crash releases ownership without a completion broadcast. Unit owner disappearance and real closed Chromium owner recover the waiter. A permanently suspended live owner causes bounded error rather than endless waiting; it must resume/close before safe ownership is reacquired. Lost committed response/cookie delivery may still require re-login; no grace/recovery secret cache.

### I17. Stale-message fencing

Older epochs/attempts, duplicates/self events, expired/unknown input and contradictory post-success failure are ignored. Protected denial retains request-start epoch; late old token/generation failures cannot announce newer invalidation. Logout fences pending refresh/login replies. No durable distributed session-version claim.

### I18. Login propagation behavior

Login/local-register intent and successful completion both send fresh non-secret auth-state-changed hints; peers discard old memory/private presentation, including peers that missed intent. No auto-login/role transfer/automatic peer refresh; explicit reload or login obtains trusted state. Separate login-change counter lets newer login supersede old queued logout.

### I19. Logout propagation behavior

Initiator clears immediately and sends logout; peers promptly clear memory/private cache. The server request waits behind cookie-changing replies and clears their actual replacement. Completion hint handles delivery ordering. Peer logout/invalidation never cancels revocation: three simultaneous logouts each reach the server as 204; reload stays unauthenticated. A newer login is the only lifecycle change allowed to supersede old queued logout.

### I20. Session invalidation behavior

Participating exclusive refresh 401 publishes generic invalidation; without atomic ownership a denial is ambiguous and stays local. Current-account auth/me or users/me resolver 403 and persistent retried 401 invalidate with request/token fences. Ordinary resource 403, Origin denial, network/5xx and coordination timeout do not globally log out peers. PostgreSQL remains authority; internal reason is never sent.

### I21. TanStack Query cache handling

removeQueries synchronously destroys/cancels authenticated auth/users roots or meta.authenticated=true queries, including consumed AbortSignal; public health/status remains. Existing private test fixtures use real private roots, with preserved assertions and new public-cache checks. Current private mutation history is cleared to discard credentials/inputs; cache removal cannot undo server mutation, so future callbacks must retain generation guards.

### I22. Session-store behavior

Zustand remains non-persisted memory-only user/access/status/generation. No channel/lock is stored in React or general app state. The actual owners are the auth store and SessionBootstrap; there is no active separate token-manager or auth-provider module. No refresh-token field, JWT authority decoding or persistence middleware was added.

### I23. Storage usage

No coordination or credential writes to localStorage/sessionStorage/IndexedDB. Existing theme preference remains; known legacy token keys are removed without reading values. Browser IndexedDB listing is empty; no storage-event/token-lock fallback.

### I24. Browser fallback behavior

Chromium verified. Web Locks alone still serialize; missing BroadcastChannel removes reliable peer UI notifications. BroadcastChannel alone carries lifecycle hints without atomic refresh optimization; 401 becomes ambiguous/local. Neither primitive: per-document single-flight plus server no-clear safety, possible local denial and deliberate reload/re-login. Missing messages do not block native release; reload/stale-token recovery consults the server. No cross-browser matrix, durable logout history or cross-device/immediate access revocation claim.

### I25. Backend cookie-clearing tests

PASS HTTP missing/malformed/unknown/expired/revoked/consumed/disabled/deleted/invalid-role cases: same safe 401, no cookie. Unexpected database failure: safe 500, no cookie. Origin denial: 403 before IO, no cookie. Explicit logout: clears correct production/development flags, including storage failure. Replay/default-perimeter successor continuation and real denial header assertions pass.

### I26. Message-contract tests

PASS 54 coordinator cases overall, including all six event kinds, every failure enum, malformed/unknown/old/future values and extra secret/account/role fields. Exact typed allowlist excludes arbitrary payloads.

### I27. Two-tab unit tests

PASS independent coordinator runtimes with shared FIFO browser-lock port: one starts while the other waits; each obtains separate own memory access using successive synthetic server-cookie state; max overlap one, no token in messages.

### I28. Three-tab/multi-waiter tests

PASS three and five independent coordinators, all turns serialized, distinct consumed-cookie sequence and independent own memory results. Real three-tab proof is I37.

### I29. Owner-failure tests

PASS network and server failure release ownership and leave peers intact; queued peer can proceed. Exclusive authentication denial generically invalidates every peer. Uncoordinated denial is ambiguous and cannot globally clear peers.

### I30. Owner-timeout tests

PASS fake timers: hung active owner, waiter aborts at 20s without HTTP/steal; after simulated browser release, deliberate retry acquires and succeeds. Client bootstrap timeout settles error, then deliberate retry restores. No unit test uses real-time sleep.

### I31. Stale-event tests

PASS old attempt/invalidation after newer success, contradictory completed-attempt failure, old authority request-start denial after newer peer, duplicate/self/expired/malformed events and late superseded-token 401/403. Actual stale BroadcastChannel delivery also leaves all three tabs authenticated.

### I32. Logout tests

PASS immediate memory/private removal, public cache preservation, pending HTTP fencing, no unnecessary peer recovery and client subscription disposal. Three simultaneous peer logouts cannot cancel every server revocation. Actual delayed-success/three-tab logout cases pass.

### I33. Session-invalidation tests

PASS generic exclusive-auth denial propagation plus explicit invalidation to multiple peers, transport current-account 403 versus ordinary permission 403, private query cancellation and late-request fences. Real disable and session revoke pass.

### I34. Retry-bound tests

PASS 401 → one coordinated refresh → one retried request → 401 STOP, including Query cancellation/retry behavior. Session endpoints/public attachAuth=false do not recursively recover. Late old-token denial cannot clear a newer successful recovery.

### I35. Storage tests

PASS local/session Storage spies and IndexedDB open/delete spies; only preference availability probes occur. Login/refresh/logout and message sentinels show no usable credentials stored/transmitted. Real IndexedDB databases/local/session/document.cookie checks independently pass.

### I36. Real Chromium two-tab result

PASS [200, 200], max network refresh in flight 1, two distinct presented credentials, two consumed/linked rows, surviving usable cookie, both authenticated. Each access came through its own actual server reply. Separate deliberately uncoordinated same-cookie race remains [200, 401]; late loser has no Set-Cookie and cannot delete winner.

### I37. Real Chromium three-tab result

PASS [200, 200, 200], max in flight 1, three distinct consumed/linked credentials, surviving usable cookie and all authenticated. Native pending-lock introspection confirms two real waiting documents before the first request is released.

### I38. Simultaneous reload result

PASS three near-simultaneous reloads: exactly three bootstrap refreshes, each begins with empty memory and ends authenticated; no storm or cookie corruption.

### I39. Real logout-propagation result

PASS three peers promptly unauthenticated, private caches absent/public cache retained, zero peer refresh requests, cookie absent, subsequent protected requests denied. A held committed 200 cannot restore old memory; logout waits then revokes/clears replacement. Three simultaneous logout requests return 204 and subsequent reload stays logged out.

### I40. Real disable/revoke result

PASS synthetic runtime-role DML disable: current-account 403 causes generic invalidation across three tabs/private caches. Re-enable and fresh restoration, then revoke sessions: exclusive refresh 401 invalidates all three. Failed refresh retains invalid HttpOnly cookie by documented policy; explicit logout clears.

### I41. Real PostgreSQL refresh concurrency regression

PASS under -race, actual DML runtime role/HTTP/separate connections: 1 success / 11 safe denials, every denial no Set-Cookie. Single coherent successor, hash-only storage/digest-not-bearer, replay denial, winner usable, actual post-consumption rollback, current-account demotion/disable, idempotent logout and bounded cleanup. Seven passing test nodes (six subtests plus parent). Server application/domain/SQL/rotation infrastructure byte-identical.

### I42. Cookie security regression

PASS real development HttpOnly, host-only, Path=/api/v1/auth, SameSite=Lax, Secure=false on explicit local HTTP; document.cookie cannot read refresh. Production Secure=true and matching logout flags covered by HTTP tests. Access memory-only. Actual production HTTPS remains unrun.

### I43. CORS/Origin regression

PASS real trusted credentialed login/refresh/preflight and same-origin logout; untrusted localhost:14174 refresh/logout wire 403, no allowed-origin/credentials headers, no cookie change. Explicit server trusted-Origin policy/CORS source unchanged; no secret-readable CSRF cookie.

### I44. Request-ID/logging regression

PASS 282 response IDs/statuses each match exactly one final API completion; 585 structured rows checked against 377 distinct sentinel values, ZERO matches. Includes all local runtime/bootstrap/API credentials. Safe auth events remain, parser 413/431 and readiness/rates correct. No channel payload/token/header logging. Log artifacts include earlier attempts; current response correlations are from the final successful browser/runtime suites.

### I45. Listener/resource cleanup

PASS one singleton independent of React remounts; HMR closes channel/removes listener, aborts queued lock requests and disposes transport subscription. Pending old client cannot restore/recover auth. Test-owned coordinators disposed. Owned containers/network/volume/local image tags, private env/sentinels/CLI and ignored browser build removed; unrelated resource and shared caches preserved, verified by Docker names/labels.

### I46. Backend fmt/vet/test results

PASS go fmt ./..., go vet ./..., go test ./... with selected Go 1.27.1 and task GOCACHE: 426 passing test nodes (92 top-level/334 subtests), 15 passing packages; two opt-in live tests skipped offline, 11 packages without tests. Global outside-module launcher is still 1.26.5; no version/rule lowered. Unchanged Phase 1H live migration fault/process suite not rerun in 1I; prior real evidence and source hashes are preserved.

### I47. Targeted race results

PASS eight actual packages: application/auth, infrastructure/security, infrastructure/database, infrastructure/persistence/postgres, HTTP routes, HTTP middleware, bootstrap and config. Real foundation PostgreSQL suite independently passes -race. An initial command included a nonexistent repository directory; corrected to the discovered actual persistence path and completed successfully.

### I48. Frontend lint/test/build results

PASS npm run lint (zero errors, unchanged 19 warnings), npm run test:run (139 tests/eight files: 54 coordinator and 37 transport cases), npm run build (strict TypeScript/Vite). Node 24.19.0/npm 11.17.0; Docker Node 22 build also passes. Normal entry 459.49kB, transport 59.80kB; only ignored static-driver build has expected chunk/dynamic-import warnings. No dependencies/gate/typechecking changes.

### I49. Compose validation

PASS six quiet configurations: base default, full, tools, full+tools, retained Phase 1G override full and Phase 1H override full+tools. Actual API/frontend images build/run, nginx -t passes. Final nginx observed /32 still equals trusted configuration. Phase 1H disposable names intentionally retained; no production manifest/deployment.

### I50. git diff --check

PASS at closure. Exact A–H foundation prefix, Phase 0/open questions, charter/source/architecture/stack, all six historical SQL files, domain/application/database/security source, dependency manifests/locks, V1 audit and shadcn primitives preserved. Final safe evidence contains no credential values.

### I51. Exact git status

All 23 changes unstaged; 18 modified/five new. No staging/commit/push/remote modification.

```text
 M backend/internal/interface/http/handlers/auth_handler.go
 M backend/internal/interface/http/routes/browser_session_test.go
 M backend/internal/interface/http/routes/http_perimeter_test.go
 M backend/internal/interface/http/routes/refresh_session_test.go
 M backend/tests/integration/foundation_test.go
 M docs/project/DECISIONS.md
 M docs/project/PHASE1_FOUNDATION.md
 M docs/project/PHASE1_SECURITY_BACKLOG.md
 M docs/project/ROADMAP.md
 M frontend/src/app/query-client.ts
 M frontend/src/lib/api-client.test.ts
 M frontend/src/lib/api-client.ts
 M frontend/src/stores/auth-store.test.ts
 M frontend/src/stores/auth-store.ts
 M integration/README.md
 M integration/browser-driver.ts
 M integration/browser.mjs
 M integration/logs.py
?? frontend/src/lib/session-coordinator.test.ts
?? frontend/src/lib/session-coordinator.ts
?? integration/evidence/2026-10-07-phase1i-baseline.json
?? integration/evidence/2026-10-08-phase1i.json
?? integration/session-browser.mjs
```

### I52. Remaining CI work

CI/OPERATIONS: automated pipeline remains absent/open. Wire deterministic selected-toolchain unit/race/frontend/Compose plus opt-in disposable-role/live/browser gates later; do not add CI here. Guide serializes shared fixtures/rate budgets and includes a six-second reset between live HTTP and browser suites.

### I53. Remaining deployment-only verification

Hosted same-site HTTPS/Secure cookies/HSTS/CSP/assets/edge proxy topology; production verified PostgreSQL TLS, least-privilege role/ownership rollout and separately credentialed release migration controls; backups/restore/recovery and operator deployment evidence. Local success is not deployed production approval; none blocks domain/database design.

### I54. Remaining product-policy decisions

OPEN-001–020 unchanged: provisioning/eligibility/institutional roles, stock reservation/borrowing/return/accountability/fine/currency/day-boundary/history/kiosk/terms/notifications/data migration policies. Family/global/concurrent-device/immediate-access-revocation policy remains unresolved. These gate dependent product design/implementation, not the engineering foundation.

### I55. Remaining non-blocking foundation backlog

19 retained primitive/hook warnings; dependency/image review; local registration user/session atomicity before future approved product provisioning; additional browser matrix if later required. Operational log sink/access/retention/rotation and cleanup scheduling/cadence/load/ownership stay CI/OPERATIONS. No automatic 1J.

### I56. Any blocker before Phase 2

BLOCKING BEFORE PHASE 2: NONE in the engineering foundation. CI and deployed-production checks remain open in their categories. Policy-dependent design still needs stakeholder decisions; Phase 2 itself still needs separate authorization.

### I57. Phase 1I exit gate

SATISFIED / COMPLETE, 2026-10-08, for source and disposable local scope. Baseline/fixed real race, atomic non-secret coordination, wait/recovery/fences, peer lifecycle/private cache, local single-flight/one retry, actual two/three/reload/disable/revoke/fallback/late-response checks, unchanged PostgreSQL single-use semantics and all standard/race/Compose/log/diff gates pass. No Phase 2/business schema/feature/mockup/CI pipeline/commit/push/deployment.

### I58. Complete Phase 1 readiness for Phase 2

READY for separately authorized domain/database design: Phase 1A–1I engineering foundations have no remaining blocker. Production verification/CI are not claimed complete. Product policies remain open and dependent designs require decisions. Phase 2 has NOT begun.

### Phase 1A–1I remaining-work categories

| Category | Remaining work / boundary |
|---|---|
| **BLOCKING BEFORE PHASE 2** | **None in the engineering foundation.** Separate Phase 2 authorization still required. |
| **NON-BLOCKING FOUNDATION BACKLOG** | Retained 19 warnings; dependency/image review; local registration atomicity before product provisioning; optional future browser matrix. |
| **DEPLOYMENT-ONLY** | Hosted HTTPS/Secure cookies/HSTS/CSP/edge topology; production verified DB TLS/roles/ownership; explicit release/recovery/backup/restore verification. |
| **CI/OPERATIONS** | CI pipeline; deterministic toolchain/race/live/browser jobs; log sink/access/retention/rotation ownership; bounded cleanup schedule/cadence/load/owner. |
| **PRODUCT POLICY** | Preserved OPEN-001–020 plus institutional session/family/global/concurrent-device/immediate-access-revocation policy. Decide before dependent model/feature implementation. |

Phase 1A configuration fails closed; 1B current-account authority/production registration containment, 1C JWT/hash/atomic single-use sessions, 1D memory/HttpOnly/Origin transport, 1E perimeter, 1F safe contracts/observability, 1G real runtime and 1H migration/privilege evidence remain intact. Phase 1I closes the mandatory cross-tab engineering finding. No Phase 1J or later-phase implementation is inferred.
