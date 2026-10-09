# Phase 7 isolated verification and owner-run rollout

**Phase 7 local integration recovery, 2026-10-09: VERIFIED / OWNER ACCEPTANCE PENDING.** Full private-backup restoration succeeded using the existing bootstrap owner only on the isolated target. All17 restored tables and effective schema permissions match; only equivalent PostgreSQL metadata representations differ. Migration000008 passed restored-data rehearsal and was then applied through the existing normal-local migrator. All8 checksums verify; prior business data and physical20/0/0/0/20 stock are preserved. Normal directory200/empty/search/filter/direct-form checks, isolated workflows/races/restart,279 frontend tests in both modes, Go gates and22 Chromium checks pass. Official FSMO terms/acceptance, production Student domains and verified activation delivery remain live-use dependencies. Phase8 is not started. See [the local integration report](../docs/project/PHASE7_LOCAL_INTEGRATION_REPORT.md).

Executed evidence is in [the completion report](../docs/project/PHASE7_COMPLETION_REPORT.md). Normal local data was not migrated during the original implementation run; the subsequent authorized recovery applied000008 after verified full restoration and rehearsal. All mutating regression/browser checks use `elabtrack_v2_batch1_test` at127.0.0.1:54832 and private0600 `backend/.env.batch1`; synthetic domain `students.example.invalid`, TEST-only policy and random account credentials. Do not point these suites at normal data or commit that environment/private fixture file.

## Isolated database setup

Use a fresh, separately named **test** volume when rerunning historical empty-database migration tests. Never reset/delete an existing volume. The `compose.phase7.yml` override owns `elabtrack_v2_phase7_pgdata`; `compose.phase7-regression.yml` accepts PHASE7_REGRESSION_VOLUME with an isolated default. Final full-regression run used `elabtrack_v2_phase7_regression_v3_pgdata`; previous volumes were retained. Only one isolated container may bind54832; normal5434 remains separate.

Create the private environment using random credentials, not example/default passwords, and confirm its DB_NAME/DB_PORT guard. Start only the overridden postgres service:

```bash
docker compose --env-file backend/.env.batch1 -p elabtrack-phase7-regression -f docker-compose.yml -f integration/compose.phase7-regression.yml up -d postgres
```

After readiness, from backend run guarded migration status/up and the existing explicit foundation grants **only in this isolated container**:

```bash
python3 ../integration/batch1-env-run.py go run ./cmd/api --migrate-up
python3 ../integration/batch1-env-run.py go run ./cmd/api --migrate-status
```

Foundation `backend/database/runtime-grants.sql` must be applied as the isolated schema owner; credentials come from the existing private container environment, never shell arguments. Feature grants belong to the paired migrations. Runtime must remain a nonsuperuser without migration tracking/schema ownership. Empty000008 down/reapply was verified before fixtures; history-bearing down is deliberately refused.

## Order-sensitive real database regressions

Run these **sequentially** on a fresh isolated database. Historical terms/accounts/inventory rollback tests assume no later consequential history; borrowing fixtures must not be created first. Paired rollback tests now remove/reapply empty000008 before exercising older pairs, while preserving meaningful history refusal and permission assertions.

```bash
cd backend
python3 ../integration/batch1-env-run.py go test -race -v ./tests/integration -run '^TestRealTerms$' -count=1
python3 ../integration/batch1-env-run.py go test -race -v ./tests/integration -run '^TestRealAccounts$' -count=1
python3 ../integration/batch1-env-run.py go test -race -v ./tests/integration -run '^TestRealInventory$' -count=1
python3 ../integration/batch1-env-run.py go test -race -v ./tests/integration -run '^(TestRealAccountDeliveryRace|TestRealBorrowing.*|TestRealInventoryContention|TestRealInventoryCorrectionSafety)$' -count=1
python3 ../integration/batch1-env-run.py go test -race -v ./internal/bootstrap -run '^(TestPhase7HTTP|TestBatch1HTTP|TestBatch1InventoryHTTP|TestBatch1InventoryBorrowabilityFilter)$' -count=1
```

Actual process persistence is also verified, with no external listener already occupying18085:

```bash
cd backend
go build -o /tmp/elabtrack-phase7-api ./cmd/api
ELABTRACK_PHASE7_API_BINARY=/tmp/elabtrack-phase7-api python3 ../integration/batch1-env-run.py go test -race -v ./tests/integration -run '^TestRealBorrowingProcessRestart$' -count=1
```

The test launches the normal binary, terminates it while pending, restarts after deadline and starts it again, asserting exactly one release. Only fixture submission time is shortened; product worker time/deadline authority is unchanged. All child processes are stopped by the test. Ordinary opt-in suites skip this process check unless the binary environment is explicitly supplied; the standalone command above was actually executed and passed.

Network suites TestRealFoundation/TestRealTermsHTTP additionally require an actual API listening18085 via the same guarded environment. Ordinary go test intentionally skips opt-in database suites without their isolated flags. Do not claim database coverage from ordinary unit tests alone.

## Chromium fixtures and evidence

Only the test build exposes `TestPhase7BrowserServer`; it is guarded by ELABTRACK_PHASE7_BROWSER plus the isolated database guard. It starts normal services/routes/expiry worker, creates synthetic random accounts/catalog and clearly labelled TEST terms/acceptance, and writes private0600 `/tmp/elabtrack-phase7-fixtures.json`. There is no production test endpoint, real-account activation bypass or live email. A test-only expiry sentinel creates a past-deadline synthetic request through normal services with an injected fixture clock, then runs the actual persisted sweep.

```bash
cd backend
ELABTRACK_PHASE7_BROWSER=1 python3 ../integration/batch1-env-run.py go test -v ./internal/bootstrap -run '^TestPhase7BrowserServer$' -count=1 -timeout 45m
```

In another terminal, use frontend Vite15175 configured to `http://localhost:18085/api/v1`. Run the two browser groups in order on fresh fixtures; set CHROMIUM_EXECUTABLE and required local shared-library path to the installed Chromium environment:

```bash
cd frontend
node scripts/phase7-browser-qa.mjs borrower
node scripts/phase7-browser-qa.mjs staff
```

The script uses real UI sign-out before switching roles. It writes sanitized checks/request paths/screenshots to `docs/ux/verification/phase7/`; no credentials/tokens/private fixture content are copied. Stop the owned test harness by creating its `/tmp/elabtrack-phase7-browser-stop` sentinel; remove only those task-owned private fixture/sentinel files after it exits. Stop owned test listeners/container if desired, retaining all data volumes. Never stop/reset normal services as cleanup.

## Full quality checks

Backend: `go fmt ./...`, `go vet ./...`, `go test ./...`, and relevant `-race` packages plus the explicit isolated suites above. Frontend: `npm run lint`, `npm run test:run -- --no-file-parallelism`, `npm run test:run -- --maxWorkers=2`, `npm run build` (includes TypeScript). Both full test modes preserve all279 assertions. The initial unrestricted23-worker run under concurrent load hit inherited lazy-login timing failures; no timeout/test weakening was used. Finish with `git diff --check`. No new lint warnings; inherited19 remain.

## Historical normal local rollout plan — subsequently executed after recovery

1. Review Phase7 code/migration and owner checklist. Stop the normal API while applying schema changes; do not reset storage. Confirm backend's existing private runtime/migrator configuration targets the intended normal local DB, not the isolated test database or production. Keep credentials in existing environment/config, never argv or Git.
2. Make a **new** private backup of current normal accounts/history; earlier checkpoint backups predate later owner activity. Set umask077, use the existing container credential environment to run pg_dump in custom format, and retain the file under ignored backend/tmp. Validate its archive listing and SHA256. Do not overwrite an earlier backup. A backup listing is not a tested restoration.
3. From repository root run `make migrate-status`; verify existing000001–000007 checksums and only the expected pending000008 suffix. Resolve any drift/target mismatch before proceeding.
4. Explicitly run `make migrate-up` with the configured migrator credential, then `make migrate-status`. Existing transactional/advisory-lock runner applies000008 and its narrow grants. Runtime/API never migrates automatically.
5. Restart existing `make dev`; verify readiness/login/protected reads and safe missing-terms behavior. Do not seed/publish TEST policy or artificially activate real accounts. Only approved official publication and each real borrower's documented acceptance can authorize real requests; required Student-domain/activation delivery gates remain.
6. Do not use migrate-down after borrowing history exists:000008 intentionally refuses destructive rollback. Recovery/restore requires a separately reviewed plan and authorization. Do not alter checksums, force migrations, hard-delete evidence or remove volumes.

The original implementation run performed no normal rollout. The later authorized local integration is now verified in the linked report; production configuration, live Brevo, commit, push, deployment and Phase8 remain outside scope.

## Final local recovery verification

Current normal target: elabtrack_v2/127.0.0.1:5434/container elabtrack_v2_postgres/volume elabtrack_v2_pgdata,000008 applied with all8 checksums verified. No reset, seed or normal restore. Private archive: backend/tmp/pre-phase7-backups/elabtrack_v2-pre-phase7-20261009T150053Z-000007.dump,0600. Recovery used PostgreSQL18.6 bootstrap owner **only** on elabtrack_v2_phase7_restore_20261009t150053z at54833, preserving migrator/runtime restrictions and normal citext ownership. Supported full pg_restore --exit-on-error --single-transaction succeeded; all17 tables, owners/effective ACLs,219 constraints,38 indexes and11 triggers verified. Equivalent two AND-check parenthesis shapes and default owner ACL representation are explained in the report. Existing metadata was not stripped or ignored.

Actual normal migration used make migrate-up after only000008 was confirmed pending. Do not blindly repeat recovery helpers: the restored target is now populated and at000008. For a new restore drill create a fresh separately named disposable target/volume with a distinctive database name and verify identities/empty application schema before restoring. Never restore over normal data. Preserve owners and ACLs; use bootstrap credentials from the disposable container environment, not argv or extra normal privileges. Use the existing CLI from backend for isolated migration commands; normal make migrate-status is read-only. Revalidate a matching backup/baseline if normal records have changed.

Final synthetic regressions used private backend/.env.batch1 and a fresh retained volume recorded in the report.16 suites/82 named subtests, including actual process crash/restart, passed. The existing guarded browser harness plus private evidence adapters ran against15175/18085:5 Borrower and9 Staff/Admin checks; normal5173/8080 contributed8 checks.38 published screenshots (5 normal,33 synthetic) are under docs/ux/verification/phase7-local-integration; normal-form captures with Student-ID metadata were excluded. Private adapters live under ignored backend/tmp/phase7-local and make no product-code change. Inherited lazy-login unit sensitivity did not recur:279 pass serially and with2 workers.

Owner review can launch the normal app via make dev if it is not already running, inspect make migrate-status, sign in at http://localhost:5173/login, and open Requests & Borrowings. Actual final normal services are already running; disposable containers/listeners are stopped, all volumes/backup retained, and temporary session/test credentials removed. Official terms and activation gates remain; do not create normal demo checkouts. No Phase8 authority is implied.
