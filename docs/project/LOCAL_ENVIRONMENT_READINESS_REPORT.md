# Local environment readiness report

2026-10-08. **PHASE 3B — COMPLETE AND OWNER VISUALLY APPROVED. Local full-stack foundation ready: YES. Phase 4: NOT STARTED.** Scope: approved visual closeout, safe checkpoint preparation and restoration/verification of the existing development stack. [Evidence index](../ux/verification/local-environment/README.md), [Phase 3B report](PHASE3B_REPORT.md), [development instructions](../../README.md#local-development).

## 1. Phase 3B approval recorded

[DEC-064](DECISIONS.md#elab-v2-dec-064--phase-3b-implementation-visually-approved) records explicit owner review of **real implementation screenshots**: Borrower Home, Equipment Catalog, Staff Dashboard, Pending Requests, FSMO navy/violet institutional palette, compact branding, shadcn foundation, mobile-first Borrower shell, desktop/tablet-first Staff/Admin shell, light/dark themes and responsive layouts. Disclosed minor adaptations are accepted within the current baseline. Roadmap, visual system, fidelity contract, composition, implementation map and review checklist reflect this disposition.

Accessibility, readable typography/control sizes and future individual PNG comparisons remain active. The reconstructed seal remains replaceable by an authenticated official vector. The other 32 screen targets remain unimplemented. Business rules and Phase 1 security were preserved.

## 2. Exact tracked/untracked changes

The initial snapshot covered **544 tracked/untracked nonignored files**, including the root ZIP. Final status has **7 modified tracked files and 241 untracked files**. Nothing is staged; no tracked file was deleted. The full, exact expanded untracked inventory is `untrackedFiles` in [WORKTREE.json](../ux/verification/local-environment/WORKTREE.json); section 29 records the exact directory-collapsed Git status.

Tracked changes:

- `.gitignore`
- `Makefile`
- `README.md`
- `docs/project/DECISIONS.md`
- `docs/project/ROADMAP.md`
- `frontend/src/app/router.tsx`
- `frontend/src/styles/index.css`

The Makefile, router and CSS implementation changes predate this closeout. Existing Phase 3A.2 references, Phase 3B components/captures and development supervisor were preserved. Closeout adds approval, restoration documentation, safe evidence and live verification utilities; section 27 distinguishes these changes.

## 3. Git checkpoint staging plan

Prepare one coherent visual-foundation/development checkpoint from the paths below. **Plan only: no staging, commit or push was performed.** An explicit dry-run verified **248 changed/new entries**, using a temporary copy of the index because repository Git metadata is read-only in the sandbox. The original index remains empty and unchanged.

| Group | Retain and rationale |
|---|---|
| Production frontend foundation | Router, semantic/application styles, 5 shared application modules, replaceable brand asset; development previews stay lazy and excluded in production |
| Development preview feature | 13 files: pages, fixtures/context, route module, meaningful tests and 6 cropped equipment assets with provenance README |
| Project/UX documentation | Phase 3A.2 historical report, Phase 3B approved report, this readiness report, decision/roadmap and 6 visual specifications |
| Binding approved references | All 47 extracted files, including 36 manifest targets, overview aids, manifest and source README; future individual images remain required |
| Historical mockup evidence | 110 files, including 99 SVGs, 10 PNGs and manifest; retain because the historical Phase 3A.2 report links them; superseded as current visual authority |
| Verification evidence | Original Phase 3B 34 files and closeout 15 files; preserve original measurements/captures and retain four live captures plus sanitized results |
| Reproducible tooling | Three browser scripts, restricted fixture SQL helper, Makefile/dev supervisor and README |

Reviewed future staging command (shown for the owner; **not executed**):

```sh
git add -- \
  .gitignore Makefile README.md \
  docs/project/DECISIONS.md docs/project/ROADMAP.md \
  docs/project/PHASE3A2_REPORT.md docs/project/PHASE3B_REPORT.md \
  docs/project/LOCAL_ENVIRONMENT_READINESS_REPORT.md \
  docs/ux/APPROVED_VISUAL_IMPLEMENTATION_MAP.md \
  docs/ux/COMPONENT_COMPOSITION.md docs/ux/HIGH_FIDELITY_MOCKUPS.md \
  docs/ux/MOCKUP_REVIEW_CHECKLIST.md docs/ux/VISUAL_FIDELITY_CONTRACT.md \
  docs/ux/VISUAL_SYSTEM.md docs/ux/approved docs/ux/mockups \
  docs/ux/verification/phase3b docs/ux/verification/local-environment \
  frontend/src/app/router.tsx frontend/src/styles/index.css \
  frontend/src/styles/application.css frontend/src/assets/brand \
  frontend/src/components/application frontend/src/features/visual-preview \
  frontend/scripts/browser-cdp.mjs frontend/scripts/phase3b-browser-qa.mjs \
  frontend/scripts/local-environment-smoke.mjs \
  integration/local-readiness-sql.py scripts/dev.sh
```

Exclude local `.env` files, the root ZIP, `.project-reference/`, caches, node_modules, dist and temporary browser output. Exact generated-secret scanning found **zero matches** in Git candidates. The approved references and verification PNGs are intentional retained evidence, distinct from temporary output. No dependency manifest/lockfile changed.

## 4. Redundant files identified

`eLabTrack_V2_APPROVED_VISUAL_HANDOFF.zip` contains **47 of 47 byte-identical extracted counterparts** in `docs/ux/approved/`. It remains on disk, unchanged, and is now explicitly ignored. No user archive or historical mockup was deleted. Additional closeout viewport captures and operational scratch output remain under `/tmp`; only four selected live screenshots are retained. Original Phase 3B evidence is unchanged.

## 5. Docker status

Docker Desktop **4.92.0**, Linux Engine **29.8.0 / API 1.56**, Compose **5.5.1** are reachable. `docker version`, `docker info`, `docker compose version` and service discovery succeeded. Default Compose services: `postgres`. The unrelated stopped hello-world container was preserved. No competing daemon, Docker cleanup or full-stack container deployment was introduced.

## 6. WSL integration status

The current WSL Linux kernel is **6.18.40.1**. Native WSL Docker CLI and Windows Docker Desktop interoperability reach the Linux engine. No distribution reconfiguration or Windows intervention is currently required. README documents starting Desktop with the Linux engine and enabling this distribution under Resources → WSL Integration if connectivity is lost.

## 7. Local environment-file status

Missing `.env`, `backend/.env` and `frontend/.env` were created from repository examples, with **0600 permissions**. All three are ignored and absent from the index. Four independent 32-byte random secrets were generated for bootstrap, migrator, runtime and JWT signing; values were never copied into documentation/evidence or printed. Runtime/migrator settings were checked for consistency.

| Concern | Verified local configuration |
|---|---|
| PostgreSQL | Dedicated `elabtrack_v2`, host port **5434**, container port 5432; root bootstrap administrator remains separate |
| API runtime | `APP_ENV=development`, port 8080, DB host 127.0.0.1, runtime login `elabtrack_runtime` |
| Migration CLI | Private `MIGRATION_DATABASE_URL`: schema-owner `elabtrack_migrator`, same database/host/port; API runtime uses its own DB settings |
| DB TLS | Existing local `DB_SSLMODE=disable` / migration `sslmode=disable`; production verify-full requirement unchanged |
| Session | Strong private JWT secret; issuer `elabtrack-v2`, access TTL 15m, refresh TTL 168h; no memory/cookie architecture changes |
| Exact Origin/CORS | Frontend http://localhost:5173; allowed localhost 5173/4173; trusted proxies empty |
| Browser transport | Frontend `VITE_API_URL=http://localhost:8080/api/v1`; root `/api/v1` applies only to Compose/nginx |
| Bounds/logging | Existing body, connection, timeout, rate-limit and logging defaults retained |

Host port 5432 already had an unrelated PostgreSQL listener. It was preserved; root/backend/migration settings use the verified free port 5434. Production credentials were not used. No credential validation was disabled.

## 8. PostgreSQL container status

Initial Docker inspection found **no existing project PostgreSQL container or project volume**. The supported `docker compose up -d postgres` created `elabtrack_v2_postgres`, `elabtrack_v2_pgdata` and the project network. PostgreSQL **18.6-alpine** is running and healthy, exposing 5434 → 5432. The new persistent volume is retained. No reset, deletion, `down -v`, adoption or privileged-user seed was run. [Database evidence](../ux/verification/local-environment/DATABASE.json).

## 9. Database/migration status

Initial `make migrate-status` found a fresh database, absent migration tracking, and three pending foundation migrations. Existing bootstrap installed `citext` and role separation. Commands actually executed:

```sh
make migrate-up
make compose-runtime-grants
make migrate-status
```

All three versions are applied with paired SHA-256 checksums verified: `000001_create_users`, `000002_create_refresh_tokens`, `000003_refresh_session_security`. Current version is `000003_refresh_session_security`. Actual PostgreSQL advisory lock contention denied a second `make migrate-status`; status succeeded after the holding connection closed. Existing atomic runner/history validation source and tests were retained. No historical SQL or migration implementation changed; no new SQL migration was created. Failure/rollback injection was not repeated on this persistent development database. [Migration evidence](../ux/verification/local-environment/MIGRATIONS.json).

## 10. Migrator/runtime privilege verification

The migrator owns public schema and all three foundation tables, and successfully applied the existing migrations. Both logins are non-superuser, with no create-database/create-role/replication/bypass-RLS attributes. Runtime has database CONNECT and schema USAGE; no database CREATE/TEMP, schema CREATE, migrator membership, tracking SELECT or users TRUNCATE. Exactly SELECT/INSERT/UPDATE/DELETE on users and refresh_tokens are granted without grant option.

A real runtime connection reported `current_user=elabtrack_runtime`, and `SELECT 1` succeeded. Real runtime CREATE TABLE and migration tracking reads failed with permission denial. Backend auth/session DML succeeded using runtime credentials. [Privileges](../ux/verification/local-environment/PRIVILEGES.json), [ownership/attributes](../ux/verification/local-environment/DATABASE.json).

## 11. Backend startup result

Actual standalone **`make backend` passed**: development configuration loaded, PostgreSQL pool connected with restricted runtime credentials, API listened on 8080, startup logs confirmed explicit migration policy, no seed, no trusted proxy and two configured origins. It shut down gracefully before the genuine outage check. Actual combined startup also passed. No database bypass or migrator runtime credential substitution was needed.

## 12. Health/readiness result

Actual GET **`/api/v1/health`** and **`/api/v1/ready`** returned **200**, standard success envelopes and alive/ready messages. Origin `http://localhost:5173` received exact ACAO, credentials true, exposed X-Request-ID/Retry-After and a present server request ID. Ready also succeeded from the real browser and during the final shutdown verification run. [Health evidence](../ux/verification/local-environment/HEALTH.json).

## 13. Frontend startup result

Actual `make dev` invoked `make frontend`; Vite **8.3.2** listened on **http://localhost:5173** with strictPort. Root and four previews were accessible. The existing project Vite process was identified by command/cwd and stopped solely to transfer supervision to `make dev`; unrelated terminals/listeners were preserved. Dependency installation was already present; no manifests or packages were changed.

## 14. API connectivity

Direct cross-origin localhost API requests succeeded through the existing centralized transport. Browser credentialed readiness exposed its request ID, development registration returned 201, login/refresh/current-account requests returned 200 and logout returned 204. Unconfigured Origin 5174 was denied 403 without ACAO; a cookie-changing POST without Origin was denied 403. No Vite proxy, CORS wildcard or Origin relaxation was added.

## 15. Session restoration checks

| Actual scenario | Result |
|---|---|
| Fresh anonymous visitor | Expected refresh 401; unauthenticated, no restoration-error banner |
| Valid session after real document reload | Refresh 200; own memory session restored; `/auth/me` 200 |
| Two tabs reloading simultaneously | Both restored independent memory sessions using existing cookie coordination |
| Logout in one tab | Both tabs cleared memory credentials; real server logout 204 |
| Revoked cookie replay | 401; safely unauthenticated; no connection-error banner; denied refresh did not clear the replayed cookie |
| Expired server-side session | 401; safely unauthenticated; no connection-error banner |
| Backend genuinely stopped | Actual ERR_CONNECTION_REFUSED; auth state error and “Unable to restore your session.” visible |
| Connectivity restored | Anonymous state recovered and banner disappeared on subsequent page load |

Refresh remained host-only, HttpOnly, SameSite=Lax, path `/api/v1/auth`, Secure=false for development HTTP. Access remained memory-only; no credentials were placed in browser storage or lifecycle broadcasts. The real runner created only a temporary generic **user** fixture through existing development registration, expired only its own session, then removed only that fixture and its session rows. Final dedicated database counts: **0 users / 0 refresh tokens**. No privileged starter accounts or product onboarding UI were introduced. No authentication UX bug remained; auth code was not changed to suppress errors. [Live results](../ux/verification/local-environment/LIVE_RESULTS.json), [genuine outage](../ux/verification/local-environment/OUTAGE.json).

## 16. make dev full-stack result

**PASS — actually run.** Backend and Vite remained concurrently operational during live session and all viewport checks. A second short run recorded the exact descendant groups and confirmed both HTTP listeners responded before interruption. `make dev` performed no implicit migration/seed/DB startup. The database remained separately managed.

## 17. Ctrl+C cleanup result

**PASS — actual PTY Ctrl+C.** The supervisor stopped both target groups, API logged `server.shutdown_started` and `server.shutdown_completed` with pool closure, and all recorded descendants/groups disappeared. Listeners 8080/5173 closed; the temporary production preview on 4173 was stopped too. PostgreSQL remained healthy on 5434. No grace timeout or orphan was observed.

GNU Make emitted termination/interrupt diagnostics, including `wait: No child processes`, during the intentional interruption. PTY completion reported nonzero 1; this is recorded rather than treated as a normal successful server exit. Both cleanup outcomes were verified independently. [Process evidence](../ux/verification/local-environment/SHUTDOWN.json).

## 18. Four preview route results

| Exact development URL | Live widths, both themes | Result |
|---|---|---|
| http://localhost:5173/__preview/borrower/home | 320 / 390 / 768 / 1280 | PASS |
| http://localhost:5173/__preview/borrower/equipment | 320 / 390 / 768 / 1280 | PASS |
| http://localhost:5173/__preview/staff/dashboard | 1024 / 1440 | PASS |
| http://localhost:5173/__preview/staff/requests | 1024 / 1440 | PASS |

All 24 viewport/theme cases loaded through real Vite with the live backend, correct title/theme, bounded page width, navigation, synthetic-data notice and no session-error banner. Borrower bottom navigation actually opened Catalog; the Staff tablet table responded to keyboard horizontal scrolling within its labeled region. No production business action was enabled.

## 19. Backend tests

From backend, **`go fmt ./...`, `go vet ./...`, `go test ./...` all exited 0**, using the module-selected Go **1.27.1**. Formatting produced no backend source changes. Database-dependent suites gated by `ELABTRACK_INTEGRATION` were not enabled by this standard run and are **not claimed passed**. The separate isolated destructive integration harness was not redirected onto the persistent development database. Actual PostgreSQL migration/checksum/lock/role/API/auth checks described above were run here.

## 20. Frontend tests/lint/build

From frontend: **`npm run lint` exited 0 with the existing 19 warnings and zero new warnings; `npm run test:run` passed 145 tests in 9 files; `npm run build` passed TypeScript and Vite production compilation.** Final lint rerun after verification-script edits still has 19 warnings. Node **24.19.0**, npm **11.17.0** satisfy the current manifest. Browser utilities passed `node --check`; restricted Python helper passed AST syntax parsing. No tests, type checks or security gates were weakened.

## 21. Browser verification

Actual Chromium **140.0.7339.16**: **40 live checks / 24 viewport-theme cases passed**, without API interception. A separate existing visual QA quick run passed **4 light cases / 10 boundary checks**, with its documented anonymous-refresh 401 fixture. That separate run verified Sheet Escape/focus, Dialog focus trap, smaller Staff navigation, Admin shell labels without authentication, keyboard table scrolling, sidebar collapse, reduced motion, production preview/fixture JS exclusion, production preview URL 404 and foundation root.

Four selected new live captures were opened and visually inspected; their paths and all case results are in the evidence index. Original 28 Phase 3B captures and measured results were preserved. An intermediate new verification utility read the old document too early after reload; a document-transition wait fixed the test race, and the final 40-check run passed without application/auth changes. Chromium shared libraries came from an existing `/tmp` extraction; no system or application dependency change. Physical devices, full zoom/virtual keyboards and future feature mutation/error states remain unrun.

## 22. Docker Compose checks

**PASS:** `docker compose config --services` lists postgres; `docker compose config --quiet` and `docker compose --profile full --profile tools config --quiet` both exited 0. Actual development postgres startup/health/bootstrap and explicit runtime grants passed. Full API/nginx/migrator image builds and production/container deployment were not run for this host-development restoration.

## 23. New warnings/errors

**Zero new frontend lint warnings; zero unexpected application exceptions/browser warnings in final verification; no unexpected startup/CORS/auth failure.** Existing 19 primitive/hook lint warnings remain. Expected 401 anonymous/revoked/expired responses, deliberately tested 403 Origin denials, advisory/privilege denials and the genuine outage are successful boundary evidence. Intentional GNU Make shutdown diagnostics and the corrected verification timing failure are disclosed above. No authentication banner was hidden or security check bypassed.

## 24. Remaining blockers

**None for the requested local foundation/readiness gate.** Application servers are deliberately stopped after Ctrl+C verification; restart with `make dev`. PostgreSQL remains healthy. Production TLS/hosting/deployment and isolated destructive harness checks are outside this local run, and physical-device/all-future-feature verification remains future work. Secure onboarding/recovery, product contracts and role reconciliation retain their recorded later-phase design work; readiness does not claim they are implemented or decided beyond current approved policy.

## 25. Exact commands required from the user

The restored local files and migrated/granted database are already ready. From this repository with Docker Desktop available:

```sh
docker version
docker compose up -d postgres
make migrate-status
make dev
```

Open http://localhost:5173/ and the four section-18 preview URLs. In another terminal:

```sh
curl http://localhost:8080/api/v1/health
curl http://localhost:8080/api/v1/ready
```

Stop both host servers with **Ctrl+C** in the `make dev` terminal; the PostgreSQL volume/service remain. Do not overwrite restored `.env` files or run seed/reset commands. No migration reapplication is needed now. For a future new checkout or reviewed pending suffix, follow README's status-first environment/bootstrap/migration/grant sequence. If connection fails, check Desktop WSL integration, container health, matching port 5434 and private credentials; if API/CORS fails, retain exact localhost Origin and direct frontend API URL. Live verification can be repeated with the section-21 tools and evidence README command.

## 26. Phase 4 readiness YES/NO

**YES — ready for a separately authorized Phase 4 task.** PostgreSQL, validated migrations, restricted runtime backend, actual frontend/API connection, session/cross-tab safety, accepted visual foundation and reproducible one-terminal dev workflow are verified. **Phase 4 remains NOT STARTED.** No business module, new migration, product auth UI, borrowing/inventory policy change or redesign was implemented.

## 27. Files changed

Relative to the task's initial 544-file snapshot, **12 existing nonignored files changed and 18 nonignored files were added**:

Changed existing files:

- `docs/project/PHASE3B_REPORT.md`
- `docs/ux/APPROVED_VISUAL_IMPLEMENTATION_MAP.md`
- `docs/ux/COMPONENT_COMPOSITION.md`
- `docs/ux/HIGH_FIDELITY_MOCKUPS.md`
- `docs/ux/MOCKUP_REVIEW_CHECKLIST.md`
- `docs/ux/VISUAL_FIDELITY_CONTRACT.md`
- `docs/ux/VISUAL_SYSTEM.md`
- `frontend/scripts/browser-cdp.mjs`
- `.gitignore`
- `README.md`
- `docs/project/DECISIONS.md`
- `docs/project/ROADMAP.md`

Added files:

- `docs/project/LOCAL_ENVIRONMENT_READINESS_REPORT.md`
- `docs/ux/verification/local-environment/B01-390-light.png`
- `docs/ux/verification/local-environment/B02-390-dark.png`
- `docs/ux/verification/local-environment/DATABASE.json`
- `docs/ux/verification/local-environment/HEALTH.json`
- `docs/ux/verification/local-environment/LIVE_RESULTS.json`
- `docs/ux/verification/local-environment/MIGRATIONS.json`
- `docs/ux/verification/local-environment/OUTAGE.json`
- `docs/ux/verification/local-environment/PRIVILEGES.json`
- `docs/ux/verification/local-environment/README.md`
- `docs/ux/verification/local-environment/REFERENCE_QA_RESULTS.json`
- `docs/ux/verification/local-environment/S01-dashboard-1440-light.png`
- `docs/ux/verification/local-environment/S01-pending-1024-dark.png`
- `docs/ux/verification/local-environment/SHUTDOWN.json`
- `docs/ux/verification/local-environment/VALIDATION.json`
- `docs/ux/verification/local-environment/WORKTREE.json`
- `frontend/scripts/local-environment-smoke.mjs`
- `integration/local-readiness-sql.py`

Additionally the three missing private local environment files were created outside Git. Existing backend, domain, Phase 1 auth/session/bootstrap/provider/Query transport, UI primitives, dependency manifests, approved references, original visual evidence and dev supervisor were preserved. All **249 protected file hashes** and **36 approved manifest hashes** still match. No baseline file was removed. The small CDP helper change adds an optional target-session argument for real two-tab verification; new code is verification tooling only. [Exact inventory/scope](../ux/verification/local-environment/WORKTREE.json).

## 28. git diff --check

**PASS:** tracked `git diff --check` returned 0. All 241 untracked files were also checked against `/dev/null` for whitespace issues; no diagnostics. Git's no-index “files differ” exit is distinct from a whitespace error. Final ignored-secret checks and exact-value candidate scanning passed; original index remained empty.

## 29. Exact final git status

`git status --short` after closeout documentation/evidence:

```text
 M .gitignore
 M Makefile
 M README.md
 M docs/project/DECISIONS.md
 M docs/project/ROADMAP.md
 M frontend/src/app/router.tsx
 M frontend/src/styles/index.css
?? docs/project/LOCAL_ENVIRONMENT_READINESS_REPORT.md
?? docs/project/PHASE3A2_REPORT.md
?? docs/project/PHASE3B_REPORT.md
?? docs/ux/APPROVED_VISUAL_IMPLEMENTATION_MAP.md
?? docs/ux/COMPONENT_COMPOSITION.md
?? docs/ux/HIGH_FIDELITY_MOCKUPS.md
?? docs/ux/MOCKUP_REVIEW_CHECKLIST.md
?? docs/ux/VISUAL_FIDELITY_CONTRACT.md
?? docs/ux/VISUAL_SYSTEM.md
?? docs/ux/approved/
?? docs/ux/mockups/
?? docs/ux/verification/
?? frontend/scripts/
?? frontend/src/assets/
?? frontend/src/components/application/
?? frontend/src/features/visual-preview/
?? frontend/src/styles/application.css
?? integration/local-readiness-sql.py
?? scripts/dev.sh
```

The ignored archive and local environment files are intentionally absent. No commit, push, remote change, deployment, V1 access or boss-repository write occurred.
