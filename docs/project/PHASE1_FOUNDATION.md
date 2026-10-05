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
