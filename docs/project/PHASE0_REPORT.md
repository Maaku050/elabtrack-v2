# Phase 0 report — eLabTrack V2 rebaseline and template adaptation

Date: 2026-10-05, Asia/Shanghai. **Requested adaptation is implemented; the unconditional Phase 0 exit gate is not satisfied because backend quality gates and real Git verification are blocked.** No Phase 1 work, business features, product tables, V1 Firebase access, database changes, commit, push or deployment occurred.

## 1. Repository identity changes

App/config/env examples and HTML title now identify eLabTrack V2. Frontend package and lockfile root package are `elabtrack-v2-frontend`; browser storage keys are `elabtrack_v2.*`. Development DB is `elabtrack_v2`; JWT issuer is `elabtrack-v2`. Compose names are `elabtrack_v2_postgres`, `elabtrack_v2_backend`, `elabtrack_v2_frontend`, and the explicitly named `elabtrack_v2_pgdata` volume. Compose backend DB_USER/DB_NAME now follow root credentials; full-profile frontend defaults to nginx `/api/v1`. These changes do not move data or alter existing volumes. Seed display names are Development Admin/User; seed credentials/behavior were not redesigned.

The module/import path stays `github.com/fullstacktemplate/backend`: no authoritative remote can be established in this session. See ELAB-V2-DEC-020. No GitHub repository URL was invented.

## 2. Template files retained

Go module/sum/dependencies, backend domain/application/infrastructure/HTTP/bootstrap structure, auth/user handlers/services/repositories, JWT/bcrypt adapters, middleware, migrations, development seeder, Go tests, both Dockerfiles/dockerignore files, Compose, Make wrappers and optional Windows dev script remain. Frontend strict TypeScript/Vite/Vitest/lint configuration, dependency versions, shadcn configuration and all 62 UI primitives, CSS tokens, hooks, shared API/storage/types, Query and Zustand support, and provisional auth/user API/schema/hook/store modules remain. Dormant auth/user adapters are documented as inherited infrastructure, not implemented V2 account policy.

## 3. Template files removed

Removed `.specify/` (26 substantive files), `.devin/` (100 substantive files including bundled unrelated examples/specialists), `docs/SPECKIT-PLANNER.md`, showcase feature (1 substantive file), dashboard feature (2), demo layout (1), generic login/register/profile screens (3), and all 62 root component examples including `example.tsx`. Removed 397 Windows `*:Zone.Identifier` sidecars across these and retained areas; new sidecars are ignored. In total 196 substantive files plus 397 sidecars were removed. A temporary pre-edit snapshot and exact deletion manifest were kept under `/tmp` for review; no source-controlled reference DOCX was created.

## 4. AGENTS.md changes

Replaced private template-owner/project/planner/Obsidian instructions with FSMO scope, source hierarchy/evidence labels, phase boundaries, accepted stack, dependency direction, frontend state ownership, server authority, database integrity, security, accessibility, testing and Git/reference safety. No unavailable private tool or plugin is a V2 requirement. Useful rules are also retained in `docs/ARCHITECTURE.md`; actual limitations are explicit.

## 5. .specify decision

Removed entirely with dependent `.devin` workflow material after runtime/build/test reference inspection. It was optional agent planning infrastructure; manifests, Makefiles, Dockerfiles and runtime source do not require it. No workflow subset, initializer or historical example spec remains to masquerade as implemented product behavior. Project documents and ordinary repository tooling govern future work. This is a rebaseline request, not a full speculative numbered feature proposal; no proposal ZIP was generated.

## 6. shadcn/ui status

Official foundation preserved: `frontend/components.json`, base-nova style, Base UI, Lucide, zinc CSS variables, aliases, and **all 62 original primitive files byte-for-byte unchanged**. Placeholder uses existing Card and Button and retained theme styles. No alternative UI framework, registry regeneration or full design system was introduced.

## 7. Backend architecture status

Clean Architecture with pragmatic domain/user/auth boundaries remains. Backend validation and parameterized repositories exist; user list paging is bounded. Generic users/auth migrations remain byte-for-byte unchanged. Documented adaptation needs include the health handler infrastructure dependency, non-atomic auth/migration behavior, raw refresh tokens, unsafe production defaults and stale authorization. No eLabTrack domain tables were created. Backend is not claimed production-ready.

## 8. Frontend architecture status

Feature API/hooks/state patterns, Query server state, Zustand appropriate global state and React local state remain. Providers now mount only Query infrastructure needed by the placeholder; no token hydration or auth gate occurs at Phase 0 startup. Unmounted auth/user adapters retain starter behavior pending Phase 1/4. Deleted page exports and obsolete navigation targets were corrected. Dependency versions and strict type checking were preserved.

## 9. Generic demo/showcase cleanup

Routes are `/` (eLabTrack V2 placeholder/theme toggle), `/status` (explicit service check) and wildcard 404. Removed demo routes, dashboard, component showcase and account screens. The manual health feature reads the raw existing `/api/v1/health` response without credentials/tokens, schema-checks it, and handles failure/manual retry. Tests verify React, routing/theme, health boundary behavior and removed routes. Mocked success is not proof of live API/database connectivity; a running API/database is required.

## 10. Reference-document handling

The capstone remains local at `.project-reference/ELABTRACK-SEMIFINAL.docx`; `.project-reference/` was already ignored and its rule was preserved/documented. Its engineering paragraphs and relevant tables were inspected without copying biographies, acknowledgements, personal contacts or screenshot/sample borrower records. Ignore-pattern behavior was checked in an isolated temporary Git repository; **actual tracked status cannot be inspected in this workspace**.

Relocated the supplied `docs/v1-audit/` to `docs/reference/v1-audit/`, preserving all 16 Markdown files byte-for-byte (README plus 15 numbered audits). Preserved the distinction between the original static audit, the later undeployed emergency hardening patch and unknown deployed behavior. No V1 source/Firebase inspection was performed.

## 11. Files created under docs/project/

- [PROJECT_CHARTER.md](PROJECT_CHARTER.md)
- [SOURCE_OF_TRUTH.md](SOURCE_OF_TRUTH.md)
- [DECISIONS.md](DECISIONS.md)
- [OPEN_DECISIONS.md](OPEN_DECISIONS.md)
- [ROADMAP.md](ROADMAP.md)
- [FOUNDATION_AUDIT.md](FOUNDATION_AUDIT.md)
- [PHASE1_SECURITY_BACKLOG.md](PHASE1_SECURITY_BACKLOG.md)
- [PHASE0_REPORT.md](PHASE0_REPORT.md) (this report)

Also adapted root/backend/frontend READMEs, `docs/STACK.md`, `docs/ARCHITECTURE.md`, env/Compose/Make/browser identity and the minimal frontend. Added five foundation feature files: two pages, health API, hook and boundary tests. No additional speculative backlog or product implementation was added.

## 12. Accepted technical decisions

ELAB-V2-DEC-001–018 record the stakeholder-confirmed FSMO focus, React/TypeScript/Vite, Go/Fiber v3/pgx, PostgreSQL/migrations/constraints/transactions, REST/version/envelope, shadcn, Tailwind v4/Lucide, Query server state, appropriate Zustand/React client state, Hook Form/Zod, Clean Architecture/pragmatic DDD, single SPA/API/database, server authority, Docker/Linux/Git development, no premature distributed infrastructure, Phase 0 boundaries, historical integrity and authorized optional-workflow cleanup. Acceptance records direction, not completion.

Super Administrator expansion is Deferred (DEC-019); authoritative module path Needs Stakeholder Input (DEC-020). No unresolved product policy was marked Accepted.

## 13. Open product decisions

OPEN-001–020 record all requested questions with V1 evidence, capstone evidence, status, impact and responsible decision-maker: public registration/provisioning; student/faculty/staff borrowers; staff/admin separation; Super Admin responsibilities/current scope; fine policy; assessment/payment/waivers; suspension; email verification; kiosk sessions; categories; archival/deletion; cancellation; denial/history; overdue vocabulary/time policy; terms versioning; provider/delivery; migration scope. Additional conflicts cover legacy technology/schema descriptions and borrowing limits/reservation/direct-checkout differences.

Actual conflicts are recorded, including capstone public-registration language versus administrative provisioning, category filtering versus absent V1 category support, fine collection versus inconsistent payment representations, and Super Admin expansion versus current FSMO scope. Capstone PHP/MySQL descriptions do not override audited Firebase V1 behavior or accepted V2 Go/PostgreSQL direction.

## 14. Phase 1 security findings

Recorded **1 CRITICAL, 10 HIGH, 6 MEDIUM and 2 LOW** findings with paths, owner roles and acceptance evidence. Critical: production can retain a published default JWT signing secret. High: both tokens in localStorage, raw DB refresh tokens, race-prone/non-atomic rotation, incomplete JWT validation, stale role/status authorization, unresolved public-registration surface, automatic/unsafe migration behavior, forced DB TLS disablement/URL encoding, recursive client refresh retries, and unguarded development seeding. Medium: proxy/IP abuse controls, production CORS defaults, API/SPA header coverage, logging/error visibility, response contracts and missing integration/CI evidence. Low: refresh cleanup/retention and retained warning/dependency review.

Hardening is documented, not executed. Proposed cookie/session details remain proposals; stakeholder account/privilege policy is still unresolved. No new business-table migration was created.

## 15. Git remote status

`git remote -v` was attempted read-only in normal and elevated execution; exit **128** in both. Exact error:

```text
fatal: not a git repository (or any of the parent directories): .git
```

The mounted `.git` directory is empty/read-only in this session. Remote presence, URL, branch, prior dirty state and credentials cannot be determined. This is not evidence that the real project has no remote. No Git metadata, remotes or remote repositories were changed.

## 16. Validation results

| Check | Actual outcome |
|---|---|
| `go fmt ./...` | BLOCKED, exit 1: go.mod requires Go >=1.27.1; local Go 1.26.5 |
| `go vet ./...` | BLOCKED, same version error |
| `go test ./...` | BLOCKED, same version error |
| `npm run lint` | PASS, exit 0; 19 warnings in retained primitives/use-mobile; no lint rules weakened |
| `npm run test:run` | PASS, 3 files / 21 tests; includes 8 new Phase 0 boundary cases and 13 retained tests |
| `npm run build` | PASS, strict `tsc -b` plus Vite 8.3.2 build; JS 475.69 kB (149.97 kB gzip), no chunk-size warning in final build |
| `git diff --check` | BLOCKED, exit 129: no accessible Git repository |
| Isolated snapshot whitespace checks | PASS across 37 modified/new files using `git diff --no-index --check`; not a substitute for actual Git/index verification |
| Primitive/migration/audit preservation | PASS, hashes match the pre-edit snapshot |
| Ignore-pattern check | PASS in temporary repo: DOCX directory, node_modules/dist, .env and sidecars ignored; actual index status unknown |
| Authored Markdown links/evidence checks | PASS: all authored links resolve, eight project docs, 20 open questions, 15 roadmap phases and 19 report sections; no business-domain migrations |

Go commands were run with `GOTOOLCHAIN=local` to avoid an automatic toolchain download/upgrade; requirements were not downgraded. Standalone gofmt inspection of the changed Go config file showed no formatting diff. Backend gates must be rerun using a compatible toolchain.

Frontend dependencies were installed with `npm ci` from the unchanged dependency graph. Initial offline/no-cache and sandbox DNS attempts failed; authorized elevated npm installation succeeded (373 packages). No package version upgrade or dependency audit result is claimed. An unnecessary toast/tooltip provider was removed from the minimal entrypoint; final lint/tests/build passed afterward. Reusable source primitives were retained.

Real PostgreSQL connection/integration tests, migration up/down/seed execution, live API/browser end-to-end tests, Docker/Compose config/build/run and production deployment were **not run**. No database or live Firebase was touched. Docker's WSL shim says integration is unavailable; static Compose inspection does not replace Docker execution.

## 17. Environment/toolchain blockers

- Local Go 1.26.5 is below declared Go 1.27.1. No unsafe upgrade or manifest downgrade was attempted.
- Real `.git` metadata is unavailable despite read-only elevated retry, preventing exact status/remote/diff and authoritative module rename.
- Docker/Compose cannot execute through this distro's current Windows/WSL integration shim.
- No running configured PostgreSQL/live API was used for verification.
- Node 24.19.0 and npm 11.17.0 are compatible; frontend registry access was resolved for installation and final gates passed.
- No Obsidian MCP tool is exposed; no inherited private notes were read or persisted. V2 instructions remove that dependency.

## 18. Exact git status

`git status` was attempted; exit **128**. The complete status output is:

```text
fatal: not a git repository (or any of the parent directories): .git
```

There is no legitimate porcelain status to report; do not infer a clean or dirty branch. A filesystem baseline was saved before edits for supplemental preservation/whitespace comparisons, but it is not the Git index/HEAD. No commit or push was performed.

## 19. Phase 0 exit gate

**Not satisfied unconditionally.** Scope/identity cleanup, authoritative documentation, open-conflict register, minimal application and frontend quality gates are complete. Before marking Phase 0 closed, expose the actual repository metadata, verify remote/status/diff and the authoritative module path, and rerun required backend fmt/vet/tests with Go >=1.27.1. These are verification/foundation identity blockers, not authorization to implement Phase 1.

Open product questions may remain open until their dependent phases; hiding them is not an exit criterion. Database/Docker/browser/deployment checks remain explicitly unrun. Current V2 stays FSMO-focused. No later phase was started.

## Phase 0 Closure Verification

Date: **2026-10-05 (Asia/Shanghai)**. This section records the final closure verification and supersedes the earlier blocked exit assessment for current status. The complete earlier report above is preserved unchanged. Closure establishes the Phase 0 engineering baseline; it does not declare the inherited auth implementation production-ready or authorize Phase 1.

| Item | Verification | Result |
|---|---|---|
| 1 | Git repository valid | PASS: real repository inspected; starting HEAD `f950e5cef9c8b38e8784a502de9e88c1271e7d48` |
| 2 | Branch | `main`; local branch up to date with the locally recorded `origin/main`; no fetch/push performed |
| 3 | Authoritative remote | Fetch and push origin both `https://github.com/Maaku050/elabtrack-v2.git`; does not point to the template; remote unchanged |
| 4 | Working tree before task | Clean: `git status` reported nothing to commit; `git status --porcelain=v1` returned empty output |
| 5 | Previous Go module | `github.com/fullstacktemplate/backend` (historical identity) |
| 6 | Final Go module | `github.com/Maaku050/elabtrack-v2/backend`; module declaration and every active internal import renamed consistently in 37 files |
| 7 | Go toolchain | `go version go1.27.1 linux/amd64` from `backend/`; `go env GOTOOLCHAIN` → `auto`; Go directive stays 1.27.1 |
| 8 | Backend `go fmt ./...` | PASS, exit 0; no files changed by formatting after the module replacement |
| 9 | Backend `go vet ./...` | PASS, exit 0, no diagnostics |
| 10 | Backend `go test ./...` | PASS, exit 0; JSON evidence reports 9 top-level tests plus 7 subtests across 4 passing test packages; packages without tests reported normally |
| 11 | Frontend `npm run lint` | PASS, exit 0; the same 19 retained UI/use-mobile warnings, no configuration changes |
| 12 | Frontend `npm run test:run` | PASS: 3 files, 21 tests |
| 13 | Frontend `npm run build` | PASS: strict TypeScript plus Vite 8.3.2; JS 475.69 kB / gzip 149.97 kB; no chunk-size warning |
| 14 | shadcn/ui foundation | PASS: original base-nova/zinc/CSS-variable/Lucide configuration and all 62 original UI primitive files unchanged |
| 15 | Docker version | PASS: Docker Engine client/server 29.8.0; read-only verification only |
| 16 | Compose verification | PASS: Compose v5.5.1; `docker compose config` exit 0; rendered environment values withheld; no services started or volumes modified |
| 17 | Reference hygiene | PASS: `.project-reference/` ignored; capstone exists locally, untracked and unstaged; all 16 V1 audit files unchanged |
| 18 | Phase 0 business boundary | PASS: only inherited backend auth/user/shared and frontend auth/users/foundation modules; only four paired generic auth migration files; no equipment, inventory, borrowing, approval, return, damage/loss, fine, notification, report, kiosk or campus hierarchy implementation |
| 19 | `git diff --check` | PASS, exit 0, no output; checked again after recording final status |
| 20 | Exact final Git status | Recorded verbatim below; closure changes are unstaged; no commit or push |
| 21 | Remaining closure blockers | None. Earlier Git/toolchain/Docker blockers are resolved for the required Phase 0 gates |
| 22 | Recommended Phase 1 order | Recorded below from the existing security backlog; recommendations only, no implementation |
| 23 | Phase 0 status | **COMPLETE** |

### Module identity and historical evidence

The stakeholder-supplied repository path was verified against the actual origin before renaming. Active Go source contains no `github.com/fullstacktemplate/backend` or other `fullstacktemplate` identifier. The dependency section of `go.mod` and all of `go.sum` are unchanged, so `go mod tidy` was unnecessary and no dependency upgrades occurred. All 37 backend changes were verified to be exactly the module/import substitution; formatting made no extra edits.

Updated current README/backend README/stack evidence and resolved ELAB-V2-DEC-020 to Accepted. The earlier report, historical audit references and V1 audit contents were not rewritten to erase the original identity or blocked history. Product-policy decisions remain unresolved where originally recorded.

### Execution details and unrun checks

The root Go launcher remains Go 1.26.5 outside the module; `backend/` automatically selects the compatible installed Go 1.27.1 toolchain. Initial sandbox attempts could not write the default Go build cache and emitted setup errors; the successful fmt/vet/test gates used `GOCACHE=/tmp/elabtrack-closure-go-cache`. This is a per-command execution setting, not a dependency, source, toolchain or test-policy change. No toolchain upgrade or directive downgrade was performed.

Docker socket access initially required elevated read-only execution; version/config checks then succeeded. Compose configuration validation is sufficient for this closure under the stakeholder's current instruction. Actual PostgreSQL/API integration, migration/seed execution, Docker service builds/runs, live browser integration and deployment remain unrun and belong to later separately authorized work. No running services, production environment, V1 Firebase or Docker volumes were changed.

The removed `.specify/`, `.devin/`, showcase/dashboard folders, generic login/register/profile pages and root example components remain absent. Retained Go Clean Architecture, pgx, SQL migrations, Docker configuration, React/Vite/TypeScript, Query, Zustand, Hook Form, Zod, Tailwind v4, shadcn and Lucide remain available. Source/configuration inspection and hash comparisons confirm the UI primitives, V1 audit, migrations, frontend manifests/configuration and Docker files are unchanged by closure.

### Recommended Phase 1 work order

This sequence is derived from `PHASE1_SECURITY_BACKLOG.md` and is **not executed**. Each future change should add its relevant tests immediately; the final step consolidates integration/CI evidence. Product policy questions in `OPEN_DECISIONS.md` require accountable stakeholder decisions before dependent product implementation.

1. **Production configuration safety** — fail closed on default/missing JWT secrets (SEC-001); validate environment/TLS/connection encoding (SEC-009); guard development seeding (SEC-011).
2. **Registration and authorization boundaries** — contain the inherited public-registration surface while policy remains unresolved (SEC-007); define current-account status/privilege/revocation checks (SEC-006); obtain policy decisions without inventing borrower/staff/admin/Super Admin permissions.
3. **Safe auth migration infrastructure** — separate release migrations from runtime startup, make SQL/bookkeeping atomic and exclusive, and establish isolated PostgreSQL fixtures/least-privilege credentials (SEC-008). This precedes refresh-token storage migrations.
4. **JWT verification** — enforce exact algorithm, issuer, expiry and identity/subject claims with negative tests (SEC-005), coordinated with current-account authorization.
5. **Refresh-token storage and lifecycle** — one-way token digests (SEC-003), atomic single-use rotation, failure recovery, reuse handling and transactional registration where appropriate (SEC-004); define bounded cleanup/retention (SEC-018).
6. **Browser session/transport contract** — review the memory-access/HttpOnly-refresh-cookie proposal before adopting it (SEC-002); coordinate CSRF/origin/cookie protections with the API; bound recursive 401 retries and synchronize session/cross-tab behavior (SEC-010).
7. **HTTP abuse and serving controls** — explicit production CORS (SEC-013), trusted proxy/IP behavior and focused auth rate limits (SEC-012), and separate API/SPA security headers (SEC-014); verify them against the chosen session transport.
8. **API contracts and observability** — align error/envelope/health contracts (SEC-016); ensure structured logs, request correlation, panic/error visibility and secret redaction (SEC-015).
9. **Integration and CI closure** — PostgreSQL/concurrency/session/negative-path tests and reproducible gates (SEC-017); review unchanged warnings and dependency/image policy (SEC-019), without weakening lint/types or deleting reusable primitives.

### Exact final Git status

```text
 M README.md
 M backend/README.md
 M backend/cmd/api/main.go
 M backend/go.mod
 M backend/internal/application/auth/login.go
 M backend/internal/application/auth/refresh.go
 M backend/internal/application/auth/register.go
 M backend/internal/application/auth/service.go
 M backend/internal/application/ports.go
 M backend/internal/application/user/commands.go
 M backend/internal/application/user/service.go
 M backend/internal/bootstrap/app.go
 M backend/internal/bootstrap/dependencies.go
 M backend/internal/bootstrap/infrastructure.go
 M backend/internal/bootstrap/migration_test.go
 M backend/internal/bootstrap/server.go
 M backend/internal/domain/user/entity.go
 M backend/internal/domain/user/repository.go
 M backend/internal/infrastructure/database/postgres.go
 M backend/internal/infrastructure/persistence/postgres/auth_repository.go
 M backend/internal/infrastructure/persistence/postgres/user_repository.go
 M backend/internal/infrastructure/security/jwt.go
 M backend/internal/interface/http/handlers/auth_handler.go
 M backend/internal/interface/http/handlers/health_handler.go
 M backend/internal/interface/http/handlers/user_handler.go
 M backend/internal/interface/http/middleware/auth.go
 M backend/internal/interface/http/middleware/cors.go
 M backend/internal/interface/http/middleware/rate_limit.go
 M backend/internal/interface/http/middleware/recovery.go
 M backend/internal/interface/http/middleware/requestid.go
 M backend/internal/interface/http/response/errors.go
 M backend/internal/interface/http/response/response.go
 M backend/internal/interface/http/routes/auth_routes.go
 M backend/internal/interface/http/routes/routes.go
 M backend/internal/interface/http/routes/user_routes.go
 M backend/internal/shared/pagination/pagination.go
 M backend/tests/unit/domain/auth/auth_test.go
 M backend/tests/unit/domain/user/user_test.go
 M backend/tests/unit/security/bcrypt_test.go
 M docs/STACK.md
 M docs/project/DECISIONS.md
 M docs/project/PHASE0_REPORT.md
```

No changes were staged. No commit, push, remote alteration, deployment, business feature or business-domain table was created. Phase 1 has not begun.

**PHASE 0 STATUS: COMPLETE**
