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
