# Pre-Phase7 Git, database and local application readiness

Verified **2026-10-09, Asia/Shanghai** against the actual repository and normal local development database. **GIT CHECKPOINT: PASS. DATABASE MIGRATIONS: PASS. LOCAL FULL STACK: PASS. HUMAN ACCEPTANCE: PENDING.** Phase7 has not started and this checkpoint does not authorize it.

These are technical readiness results. The normal database has no accounts, so authenticated manual owner review still needs authorized initial Admin access or recreated isolated test fixtures. No default login, real Student, official terms or live email was created during this checkpoint.

## Git state, staging and checkpoint

| Requested evidence | Actual result |
|---|---|
| Branch / starting HEAD | `main`, `dc04db7fec3cf57fe991bca3092df58ef9faddc5` (`Terms & Onboarding`) |
| Existing Phase3B/4A/4B history | Phase3B `451d5cb`, Phase4A `effb5f2`, Phase4B `dc04db7` already committed; no duplicate staging of their historical snapshots |
| Remote | `origin` fetch/push `https://github.com/Maaku050/elabtrack-v2.git`; unchanged; no push/fetch/deployment performed |
| Preflight status | Empty index;38 modified tracked paths and203 untracked individual paths. Exact match with [BATCH1_FILES.txt](BATCH1_FILES.txt); no mixed unrelated changes |
| Selective staging | Explicit NUL-delimited path list; original241 attributable Batch1 paths plus acceptance checklist and fresh verification JSON;243 files total:38 modified,205 added, no deleted/renamed/type-changed paths |
| Reviewed content | Backend/frontend feature source, security/integrity tests, paired000006/000007 migrations, isolated development scripts, current policy/domain/API/UX documentation and135 relevant synthetic engineering PNG captures. README reconciled with actual schema/missing initial login |
| Staging integrity | `git diff --cached --check` PASS; exact staged blob/worktree agreement PASS; summary243 files,10171 insertions,99 deletions. Removed lines are reviewed replacements/refactoring, not deleted files or unexplained destructive changes |
| Secrets/privacy inspection | All108 staged text files scanned; no private-key/provider-key/JWT/activation-link literals found. Five connection-URL candidates reviewed: published local example, commented production placeholders and synthetic tests only. Email literals use reserved example domains. Test captures come from isolated synthetic accounts; representative bulk/activation captures inspected;135 PNG signatures verified with no text/EXIF metadata. No real roster, fixture password or captured activation token staged |
| Secret log inspection | Configured nontrivial password/secret/connection values checked privately against newly generated logs: no exact matches. No secret values printed or included in reports |
| Preserved files | All10 original000001–000005 migration files,62 shadcn primitives and47 approved visual-package files retain initial hashes |
| Local feature checkpoint | **PASS — `29b49ecd917c8d72c64db5450b9269270bf9dda1`**, `feat: implement eLabTrack account and inventory foundations` |
| Current HEAD/status at report capture | Feature checkpoint above; `main`; clean index and working tree. This report is then added in a separate documentation closeout commit, without changing feature code. Final HEAD is the commit containing this report; obtain it with `git log -1 --format=%H -- docs/project/PRE_PHASE7_READINESS_REPORT.md`. The final execution response supplies that hash and actual post-commit status |

Excluded from both checkpoints: root/backend/frontend real `.env`, old ignored `.env.phase4a/.env.phase4b`, `.project-reference/`, redundant owner handoff ZIP, frontend `node_modules/` and `dist/`, private local logs/scripts/fixtures under `/tmp`, and the ignored backup under `backend/tmp/`. Generated `.env.batch1` and original private Batch1 credentials/mail/rosters had already been removed. No database dumps, institutional reference documents or runtime artifacts were staged. Exclusion did not discard or overwrite user material.

## Normal database preflight and backup

| Requested evidence | Actual result |
|---|---|
| Confirmed target | Database `elabtrack_v2`, development host `127.0.0.1:5434`; Docker container `elabtrack_v2_postgres`, project `elabtrack_v2`, persistent volume `elabtrack_v2_pgdata` |
| Identity checks | Root Compose configuration, backend runtime configuration, actual published port, container project/config/working-directory labels and SQL `current_database()` agree. Runtime is `elabtrack_runtime`; CLI owner connection is `elabtrack_migrator`. No production connection |
| Docker/WSL | Docker Desktop Linux engine accessible; Engine29.8.0; PostgreSQL18.6 Alpine healthy before/after; `docker compose up -d postgres` succeeds and preserves the running container/volume |
| Pre-migration schema | `000005_terms_acceptance`; runner verifies all five applied paired checksums and reports only000006/000007 pending |
| Existing data | Before: users0,refresh_tokens0,terms_versions0,terms_acceptances0; one terms_publication row with NULL current-version pointer. Aggregate inspection only; no private records exported to docs |
| Roles | Migrator/runtime logins are not superuser and have no create-role/create-database privileges. Runtime has no membership in owner role, no public-schema CREATE, and no schema_migrations SELECT |
| Safe backup | **PASS**, custom-format `pg_dump` using configured migrator password internally; no credentials in filename/argv/output. Private directory0700 and archive0600; verified ignored by Git |

Backup location:

```text
/home/marvin/projects/eLabTrack_V2/backend/tmp/pre-phase7-backups/elabtrack_v2-pre-000006-20261008T181715Z.dump
```

The filename uses UTC; the local date is2026-10-09. Archive size **20209 bytes**, SHA-256 `d5342d31199dcf7b506eab56344788dfba2b47795428b13caf1b3444558086bf`. Validated custom archive signature,64 table-of-contents lines and full SQL decode via `pg_restore --file=/dev/null`. This verifies readable content without restoring over normal data; restoration to another database was not run. An initial private `/tmp` copy also remains; the retained ignored `backend/tmp` copy is the backup of record. No database/volume was dropped, reset or removed.

## Migration and PostgreSQL verification

Executed existing **`make migrate-up`**, then **`make migrate-status`**, using the configured matching local migrator URL. The original advisory-lock, paired-checksum, ordered pending-suffix and transaction runner was unchanged. No down migration, forced reapplication, adoption or manual SQL workaround was used.

| Check | Result |
|---|---|
| `000006_account_management` | PASS, applied once; account profiles, hash-only activation, roster batches, immutable audit and receipts available |
| `000007_inventory` | PASS, applied once; categories, equipment, canonical images, movements, audit and receipts available |
| Final version/checksums | `000007_inventory`; all seven applied pairs **checksum=verified** |
| Table preservation | All six original public tables remain. Eleven approved new tables added; all17 public tables owned by migrator; no unexpected removal |
| Data preservation | Users/sessions/terms/acceptances stay0; singleton publication pointer staysNULL. New feature data remains empty; final normal count check after browser/API smoke unchanged |
| Runtime connection | Actual runtime login succeeds; SELECT/count on account/inventory schema succeeds; backend readiness200 through the actual configured host runtime connection |
| Foundation permissions | Existing named users/refresh DML retained; no new blanket grants run. Foundation runtime-grants script did not need reapplication |
| New permissions | Named SELECT/INSERT and specified column UPDATE only; no account/history/equipment DELETE, no history table UPDATE; equipment A/T updates allowed, R/C/damaged-held UPDATE denied |
| Privilege separation | Runtime owner membership/schema CREATE/migration-history SELECT still false; owner/runtime remain separate non-superuser roles |
| Guards | Account audit/receipt immutability, one-time roster confirmation, classified-Borrower role guards, inventory movement/audit/receipt/image immutability and existing terms triggers are enabled |
| Smoke strategy | Real PostgreSQL read-only transactions with owner and runtime roles; no disposable-fixture mutation suite was pointed at normal data |

This checkpoint's database checks are newly executed normal-data smoke checks. The extensive real PostgreSQL rollback, contention, activation and authenticated feature tests in [Phase5](PHASE5_REPORT.md), [Phase6](PHASE6_REPORT.md) and [BATCH1_VERIFICATION.json](BATCH1_VERIFICATION.json) remain prior isolated evidence; they are not claimed as rerun here. Ordinary `go test` skips opt-in real-database suites without their explicit isolated environment.

## Startup, API and browser verification

Started the unchanged normal **`make dev`** workflow after confirming8080/5173 had no listeners. Backend selects Go1.27.1, compiles, connects with runtime credentials and listens on8080; Vite8.3.2 starts on5173 with strictPort. Host frontend configuration actually targets `http://localhost:8080/api/v1`; root Compose `/api/v1` is not incorrectly used for host Vite.

New smoke evidence in [PRE_PHASE7_VERIFICATION.json](PRE_PHASE7_VERIFICATION.json):

- **12 HTTP checks PASS**: health200,ready200; login/activate malformed-body400 and anonymous refresh401; borrower directory/policy/template, restricted staff directory, equipment/categories and current-terms routes all registered and return authenticated-denial401. Response envelopes, server request IDs and exact localhost CORS origin checked. A401 proves route protection/reachability, not an authenticated business workflow.
- **10 real Chromium checks PASS**, zero runtime exceptions: real anonymous refresh401; login renders; `/status` Check connection reaches the actual API through the frontend transport; five management/reconciliation deep links redirect to login;390/1440 login has no horizontal overflow; dark-theme control works. No browser network mocking or credentials used in this checkpoint.
- No migration-related startup error, startup migration/seed, institutional terms publication, account creation or provider delivery. Official terms remain unpublished; Brevo is not needed for process/health readiness.
- **Graceful API shutdown PASS**: Ctrl+C invokes supervisor cleanup; API records `server.shutdown_started` then `server.shutdown_completed` and closes its runtime pool. Both listeners released; no forced-kill grace message. Chromium exits. PostgreSQL remains healthy and its volume stays intact.

The interrupted make command returned nonzero and nested make printed `No child processes` during signal handling. This is disclosed supervisor shutdown output; backend graceful completion and released listeners were independently verified. No workflow code was changed to suppress it.

## Newly executed quality checks and limits

| Gate | Actual result |
|---|---|
| `go fmt ./...` | PASS, no formatting changes |
| `go vet ./...` | PASS |
| `go test -count=1 ./...` | PASS,23 tested packages,12 with no test files; caching explicitly disabled after the initial cached invocation |
| Targeted `go test -count=1 -race` | PASS,8 packages: account domain, auth application, inventory domain, catalogimage, email, spreadsheet, config and HTTP middleware |
| `npm run lint` | PASS,19 inherited warnings,0 new; inherited primitive/mobile-hook warnings retained |
| `npm run test:run` | PASS,208 tests in14 files, matching Batch1 count; no count change |
| `npm run build` | PASS, TypeScript and Vite production build |
| Migration runner status | PASS, actual normal database all7 checksums verified |
| PostgreSQL/HTTP/browser smoke | PASS, fresh checks described above |
| `git diff --check` / staged equivalent | PASS |
| Source/preservation/privacy | PASS, audited exact selected paths and immutable original-file hashes |

Toolchains: backend automatically selects module-declared Go1.27.1; global Go launcher remains1.26.5. Local Node24.19.0/npm11.17.0 satisfy the existing manifest. Installed Chromium was used with its existing private dependency-library path. No dependency upgrades or weakened checks.

Fresh private log files: `/tmp/elabtrack-prephase7-{go,race,lint,tests,build,dev}.log`. Opt-in mutating isolated PostgreSQL/full authenticated Chromium suites, production/full-profile Docker image builds, TLS/deployment, live institutional mail and backup restore were **NOT RUN** here. Their prior results, where applicable, are explicitly separate.

## Owner startup, URLs and human acceptance

From repository root:

```sh
docker compose up -d postgres
make migrate-status
make dev
```

Frontend **http://localhost:5173/**; login **http://localhost:5173/login**; connection UI **http://localhost:5173/status**. Backend **http://localhost:8080/api/v1/health** and **http://localhost:8080/api/v1/ready**. Servers were cleaned up after verification, so run the commands before opening them. Do not rerun migrations just for ordinary startup; inspect status when uncertain.

Actual feature routes and controls, source categories A/B/C/D, expected results and blank human outcome spaces are in [PRE_PHASE7_ACCEPTANCE_CHECKLIST.md](PRE_PHASE7_ACCEPTANCE_CHECKLIST.md). Key routes are:

| Feature | Local route, prefixed by `http://localhost:5173` |
|---|---|
| Admin/Staff borrower directory | `/staff/borrowers` |
| Individual Student / Faculty | `/staff/borrowers/new`; select category on the same Admin-only form |
| Student bulk creation / deactivation | `/admin/borrowers/bulk`; select operation on the same page |
| Borrower detail | `/staff/borrowers/{id}` |
| Borrower catalog | `/borrower/equipment`; current published-terms gate applies |
| Inventory / categories / equipment creation | `/staff/inventory`, `/staff/inventory/categories`, `/staff/inventory/new` |
| Operational equipment detail / edit | `/staff/inventory/{id}`, `/staff/inventory/{id}/edit` |
| Reviewed opening/add/remove | `/staff/inventory/{id}/adjust` |
| Admin reconciliation | `/admin/inventory/{id}/reconcile` |

`{id}` must come from an actual account/equipment UUID. No `/admin/borrowers` directory, separate Faculty URL, or separate bulk-deactivation URL is invented. Normal anonymous/preview checks are currently available; authenticated scenarios need an authorized account or rebuilt isolated fixture environment. Captured test mail works only in that test harness. Approved development previews are illustrative and excluded from production.

Remaining **technical/operational review blockers**: authorized initial Admin access for this empty normal database; recreated private isolated fixtures if chosen for human review; operational activation/reissue/recovery ownership; real obligation/archive adapters and transactional borrowing/current-terms integration when Phase7 is separately authorized. Current obligation projections honestly remain UNAVAILABLE. No password-recovery endpoint, Admin creation/promotion/status endpoint, borrowing/cart/reservation/checkout/return/replacement/fine settlement/report/notification feature is available.

Remaining **external/institutional dependencies**: approved exact SKSU Student domains/roster inputs; Brevo API key, verified sender, permitted activation origin and separately authorized successful live delivery testing; official FSMO terms finalized/approved/published after presentation. Blank Student domains fail closed. Missing official terms do not block independent account/inventory engineering, but do gate live borrowing and institutional consent. No actual email was sent and live Brevo delivery is not claimed.

**Readiness for owner acceptance:** the Git/database/process prerequisites pass, with public connectivity/previews reviewable and authenticated-review data/access dependencies clearly documented. All human checklist outcomes remain unconfirmed. Phase4A remains complete; Phase4B foundation verified with external institutional gates; Phase5 partially complete with external gates; Phase6 complete within reviewed scope. Only a separate explicit owner acceptance and Phase7 instruction can advance the project.
