# Phase 7 local integration and backup-restore recovery report

**Formal owner acceptance, 2026-10-10 — DEC-081:** Phase7 Borrowing & Reservations is OWNER ACCEPTED following manual core-workflow review and automated safeguards. Student/Faculty mobile UI is accepted as functional Phase7, not final presentation; DEC-080 governs the future Phase11 browse/cart. Only a safe LOCAL checkpoint of verified Phase7/demo tooling/UX documents is authorized now. Do not begin Phases8–14, push or deploy. Official FSMO terms/current consent, approved production Student domains and required activation readiness remain live-use gates. Earlier pending-acceptance/correction statements below are historical. See [the checkpoint report](PHASE7_OWNER_ACCEPTANCE_CHECKPOINT.md).

2026-10-09 — **PHASE 7 LOCAL INTEGRATION VERIFIED / OWNER ACCEPTANCE PENDING.** The owner expressly authorized safe recovery of the failed disposable restore, followed by guarded migration000008. Required restore/rehearsal/preservation/API/workflow/regression/browser gates were executed. No product-code or migration-file changes were needed. No Phase8, commit, push, deployment, normal reset, fake consent, account activation change or live mail.

## 1. Original application error and cause

Requests & Borrowings originally returned500 INTERNAL_ERROR. Actual authenticated Chromium generated `GET /api/v1/borrowings/`500 on the normal8080 API; runtime read-only SQL returned42P01, absent `borrowings`. Normal migration history was000007, with000008 pending. The expiry worker hit the same schema gap. This was a missing normal schema, not an empty-directory defect. After migration, actual directory GET succeeds200 with total0/items[]; the UI renders its empty state.

## 2. Normal and disposable database identities

| Target | Database / loopback port | Container / volume |
|---|---|---|
| Normal owner data | elabtrack_v2 /5434 | elabtrack_v2_postgres / elabtrack_v2_pgdata |
| Backup restore and real-data rehearsal | elabtrack_v2_phase7_restore_20261009t150053z /54833 | elabtrack_v2_phase7_restore_20261009t150053z_postgres / elabtrack_v2_phase7_restore_20261009t150053z_pgdata |
| Fresh synthetic regressions/browser | elabtrack_v2_batch1_test /54832 | elabtrack_v2_phase7_recovery_regression_postgres / elabtrack_v2_phase7_recovery_20261009t151930z_pgdata |

All PostgreSQL18.6 Alpine. Normal target confirmed from existing private configs, actual Compose/project-root labels, container health, port/mounts and SQL identity. Disposable databases/volumes are distinct, and the normal volume is absent from them. No frontend, expiry worker or mail adapter was connected to the restored copy; its only mutations were full restore and migration rehearsal. Synthetic API18085/frontend15175 are separate from normal8080/5173. Normal configuration files were not overwritten.

## 3. Initial baseline and migration status

Before recovery, all17 normal public-table counts/digests matched the paused original backup baseline. Existing000001–000007 checksums verified; only000008 pending. All fourteen historical paired SQL files match HEAD byte-for-byte.

| Original data | Actual count/balance |
|---|---:|
| Admin / Staff / Borrower accounts |1 /0 /2|
| Borrower profiles |2|
| Equipment / categories / catalog images |1 /1 /1|
| Account audits / operation receipts / roster batches / activation tokens |6 /2 /5 /2|
| Inventory audits / movements / operation receipts |8 /3 /8|
| Refresh-session rows at original paused baseline |46|
| Terms versions / acceptances |0 /0|
| Available / reserved / checked-out / damaged-held / total-tracked |20 /0 /0 /0 /20|
| Conservation violations |0|

There is one equipment row; its actual balance equals these totals. Private comparisons include all original fields without publishing real identities, password hashes, token values or binary catalog images. Aggregate digests remain in ignored0700 private helper storage.

## 4. Migration000008 review and safeguards

Entire paired migration and existing runner inspected. Four new tables: borrowings/items/events/operation receipts. Five explicit borrowing indexes plus unique inventory source-effect index; PK/unique backing indexes also created. Restrictive FKs retain accounts, owned original terms acceptance, equipment and aggregate/item/equipment sources. Checks enforce entry/state/timestamps,24-hour expiry, due-after-handover, positive bounded quantities, denial reason, exact movement source/vector and automated actor. Triggers protect transition/history/receipts.

Inventory movements gain borrowing_item_id; only historical kind/no-custody constraints are replaced by custody-aware checks. Existing conservation remains. Expiry alone permits null actor. All3 old normal movements have nonnull actors and zero R/C/D deltas, and survived unchanged. Narrow authorized feature grants permit R/C and legal transition fields; no runtime damaged-held update, tracking read, history deletion or superuser grant. Up performs no normal row reset/credential change/history erase. DDL can take exclusive table locks; normal API/frontend stayed paused through restore/rehearsal/migration/preservation. Existing advisory-lock runner applies SQL and paired-checksum tracking atomically. Only000008 was pending at normal application.

Down refuses consequential borrowing/receipt/source-ledger history; it is not an automatic recovery procedure for loans. Historical empty down/up and history refusal were re-executed in synthetic regressions. No down/reset was applied to normal data or the restored real-data clone.

## 5. Original restore failure and ownership diagnosis

The first attempt used restricted elabtrack_migrator for pg_restore and failed `must be owner of extension citext` at `COMMENT ON EXTENSION citext`. Its single transaction rolled back; normal migration was not attempted. The stop rule was observed. The new owner request explicitly authorized recovery and resumption.

Actual normal database owner is postgres, public schema owner elabtrack_migrator, citext1.8 and plpgsql owned by postgres. Migrator/runtime are nonsuperusers without CREATE ROLE/DB. Fresh disposable bootstrap uses the same ownership arrangement. Restoring that extension comment as migrator is insufficient; the existing bootstrap owner is appropriate for this isolated operation. Normal roles, database ownership and citext ownership were not changed.

## 6. Private archive and integrity

`backend/tmp/pre-phase7-backups/elabtrack_v2-pre-phase7-20261009T150053Z-000007.dump`;145915 bytes, file0600, directory0700, Git ignored. Custom PGDMP archive from supported pg_dump; original command/listing succeeded. SHA256 `4e008c0d40944344ab737adf71636aae2cf32d73a77d17ddf3e7605e87a7e8fa` rechecked before recovery and after final cleanup; archive is unchanged. All17 expected table-data entries and citext extension/comment metadata were inspected. No backup was uploaded/staged/copied into tracked source.

## 7. Successful full restoration and record comparison

Installed pg_restore18.6 version/help checked for --exit-on-error/--single-transaction. Positively identified disposable target resumed with its separate volume and empty public application schema. Full restore authenticated as its existing postgres bootstrap owner through container credential environment; password never entered argv/logs. pg_restore completed exit0 in one transaction, retaining archive object ownership/ACLs and extension metadata. No --no-owner, --no-acl, --no-comments, ignored failure or normal privilege workaround.

All17 source/restored table counts and canonical row digests match exactly. Restored migration history000007. Required schemas/extensions/owners,219 validated public constraints,38 indexes,11 noninternal triggers and all effective table permissions verified. Conservation violations0, no missing critical tables/constraints/indexes/audits. See [restore evidence](verification/phase7-local-integration/restore-recovery.json).

The raw metadata comparator paused on three **representation** differences, investigated read-only: title/body checks flatten associative AND parentheses after dump/reparse; expression comparisons on empty/blank/one-byte/limit/over-limit/tab/null boundaries show0 mismatches. schema_migrations changes an explicit all-owner ACL to null/default-owner representation; expanded effective ACLs match exactly across every table. Owners, schema ACLs, remaining constraints, indexes, triggers and extension versions/ownership match. This is full functionally equivalent restoration, not stripped metadata or suppressed missing checks.

## 8. Migration rehearsal against restored owner records

Existing CLI invoked from backend with a private guarded wrapper targeting only distinctive restore name/54833, APP_ENV=test and its own credentials. Rehearsal status confirmed000007/all7 checksums/only000008 pending. Up completed successfully; status all8 verified/current000008. All16 preexisting nontracking tables—including sessions—retain original digests/counts; schema_migrations only appends000008. Four borrowing tables empty. Validated290 public constraints,53 indexes and all6 explicit new indexes. Conservation0; runtime/migrator remain nonsuperuser, D/history/tracking privileges remain forbidden. Normal17-table baseline was reconfirmed unchanged after rehearsal. An initial status invocation from repo root reported missing migration directory without mutation; corrected to documented backend cwd before rehearsal.

## 9. Normal migration execution and version

All safety gates passed before applying normal schema. Normal target/unchanged baseline reconfirmed; historical SQL intact; existing make migrate-status reported only pending000008. Existing **make migrate-up** executed once with normal private migrator config on2026-10-09 at23:18 Asia/Shanghai;000008 transaction succeeded. No unrelated migration/reset/manual schema mechanism. make migrate-status then confirmed000001–000008 applied/checksum verified/current000008_borrowing. Exact output: [migration-status.txt](verification/phase7-local-integration/migration-status.txt). No normal restore/credential change/seed occurred.

## 10. Normal preservation after migration and browser checks

Immediately post-migration all16 original nontracking table counts/digests match, including accounts/profiles/activation tokens, sessions, catalog images, audit/ledger/receipts/terms. Tracking appends one legitimate version; four new tables empty. Schema checks match rehearsal:290 validated constraints,53 indexes, six explicit new indexes,0 conservation violations and reviewed narrow grants.

After real authorized login/refresh/logout browser checks, all15 original non-session/nontracking tables remain identical; new borrowing tables still empty. Session rows alone legitimately reflect authenticated testing after their unchanged migration baseline. No existing password, category, account activation, terms acceptance, equipment quantity or history was changed by integration. Final normal A20/R0/C0/D0/T20; postgres still owns citext and runtime/migrator remain nonsuperusers.

## 11. Normal borrowing API and original error resolution

Actual real Admin session directoryGET200, total0/items[]/page1/per_page25. Search/status combinations and explicit page2 agree with server total0; no page-local client filtering. The absent-table500 is resolved without product-code changes. Anonymous borrowing requests401. Actual readiness and normal frontend200. Final browser sessions show no unexpected500/runtime exceptions; expected unpublished terms503 is separately verified.

## 12. Normal local UI and session verification

Owner entered credentials only at hidden terminal prompt in private secure-login.py. Real Phase4A login validated ADMIN; only a temporary0600 server-issued HttpOnly refresh cookie saved. Chromium restored via normal refresh/current-account resolution, never injected roles/fabricated tokens. All8 final normal checks pass: directory/empty/results, desktop/tablet/mobile themes, URL search/status+refresh, out-of-range pagination recovery, bounded failure+Retry, persistent SPA/direct form, expected missing terms, keyboard focus and real UI logout. Direct opens actual account/equipment/due/handover form, review disabled without eligible selection. No normal checkout or artificial eligibility/consent.

Retry check intercepts only TEST-RETRY-ONLY browser GET responses with simulated503; one automatic Query retry then explicit Retry returns actual backend200. Normal server is not altered. Missing current terms actual503 TERMS_NOT_PUBLISHED requested once as GET; OPTIONS preflight is counted separately. Temporary owner cookie removed after UI logout.

## 13. Isolated workflow results

Fresh guarded synthetic database, TEST-only policy, random synthetic users and equipment. Borrower creation/review/one submission reserves2 immediately; cancellation confirmed restores2 exactly, history retained. Staff denial requires reason, persists borrower-visible reason and releases once. Physical checkout requires future Manila due and explicit handover, reserves-to-custody2 once. Direct eligible issuance has no pending hold and issues1. Fixture final A17/R0/C3/D0/T20. Real sweep creates EXPIRED history/system actorNULL with0 residual reservation. No corresponding normal fixture or borrowing row.

## 14. Concurrency, rollback, idempotency and crash recovery

Real -race PostgreSQL tests assert last-unit competition, multiple reservations, multi-item rollback, same-key concurrent replay, changed/duplicate requests, cancel/checkout, expiry/checkout, denial/cancel and direct/reservation; exact nonnegative/conserved vectors and no double release/checkout. Required audit failure rolls back every consequence. Read snapshots cannot mix transition header/items. Current account deactivation/role changes and terms revision serialize with commands; deactivation preserves liabilities/stock/due/history. Archive boundaries and inventory correction/stock contention remain guarded.

Actual compiled product API was killed while pending, restarted after controlled persisted deadline and started again. Exactly one EXPIRED release verified. Product DB clock/startup worker authority unchanged; no24-hour real wait, scheduler bypass or normal fixture clocks.

## 15. Authorization, activation and institutional terms

Real suites cover anonymous401, Borrower staff-operation denial, owner-scoped404, stale role/status, inactive target, activation-required Borrower, absent current policy/acceptance, inactive equipment and current authority before replay. Foundation refresh rotation/concurrency/logout/session invalidation and product permission matrix passed. Terms publication/acceptance race/ownership/rollback/immutability tests passed. Browser real logout/role switching and normal session recovery retain Phase4A mechanisms.

Normal terms versions/acceptances remain0; live borrowing cannot bypass publication/current consent or activation. Official FSMO wording after presentation, approved production Student domains/roster inputs and configured Brevo key/verified sender/tested delivery/operational ownership remain external. No live mail, unofficial normal policy, Admin-auth backdoor or normal Student/Faculty sample creation.

## 16. Current executed quality results

| Check | This recovery result |
|---|---|
| Go fmt ./... | PASS;0 formatting changes |
| Go vet ./... / go test ./... | PASS; ordinary opt-in DB cases skip without explicit guards |
| Relevant Go race packages | PASS domain/application/bootstrap/response/database/security |
| Real PostgreSQL/API regressions |16 unique suites,82 named subtests PASS with -race |
| Actual API crash/restart | PASS as explicit standalone guarded suite |
| Frontend lint | PASS,19 inherited warnings,0 new |
| Frontend serial |23 files,279 tests PASS;69.04s |
| Frontend2-worker parallel |23 files,279 tests PASS;42.89s |
| TypeScript + production Vite build | PASS via npm run build |
| Chromium |22 checks PASS;38 published PNGs |
| git diff --check | PASS after final documentation |

Toolchains: Go1.27.1, Node24.19.0, installed Chromium1187 with existing local shared-library path. No mismatch/installation workaround. Inherited lazy-login timing failures did not recur in either unit mode; no timeout/assertion/typecheck weakening. Sanitized log paths/hashes and suite/subtest names are in [verification.json](verification/phase7-local-integration/verification.json).

## 17. Browser execution and earlier incomplete attempts

Final normal8/Borrower5/Staff-Admin9 checks completed on actual Chromium. Additional populated isolated page2 matched API items/counts; status reset page, search intersection yielded0 and clear restored original total. Responsive390/768/1024/1440, light/dark, basic focus/dialog/feedback/no-overflow checks passed. This is not a full assistive-device/WCAG/non-Chromium certification.

Earlier incomplete private-harness attempts are disclosed: immediate status input before deferred search request completed lost its search in automation; waiting for actual search completion/native selection preserved both. Page.reload initially inspected the old document; final requires fresh-document marker before evaluating filters. A request counter initially included OPTIONS and failed its expected one-GET assertion; corrected to GET-only. The first isolated Borrower attempt had no15175 Vite server and could not find login; explicitly started its isolated frontend before final run. None of these attempts is counted as a pass; final original functional assertions remained and no product changes were required.

## 18. Screenshots and privacy

[Evidence folder](../ux/verification/phase7-local-integration/README.md):5 normal directory/filter/retry screenshots with account identity hidden,13 synthetic Borrower and20 synthetic Staff/Admin screenshots, plus per-group acceptance JSON. Four representative light/dark/direct/pending captures visually inspected. Final acceptance generated41 captures but **3 normal Direct form captures were removed** because real Borrower metadata could expose an official Student ID. Normal form verification still passed; synthetic Direct screenshots document presentation without publishing real identities. No passwords/tokens/dumps appear in published evidence. Original31 Phase7 implementation screenshots and eight unrelated diagnostics remain untouched.

## 19. Limitations, cleanup and remaining blockers

No remaining engineering blocker for this local integration. Owner acceptance remains pending. Official terms/activation/live-domain external gates above still block real borrowing when unmet. Future return/fine/replacement resolution, large-backlog expiry lag/monitoring, prior unrelated unmatched-GET investigation, production rollout, non-Chromium and physical assistive-device checks remain outside this task. Expiry outages can delay release until persisted sweep catch-up; bounded worker behavior unchanged.

Normal API8080/frontend5173 remain running for owner review. Only owned synthetic harness/API/Vite and the restore/regression containers were stopped; all volumes retained. Private synthetic credential/sentinel files and owner cookie removed. Original custom backup unchanged0600; ignored private restore/test environment/helpers protected. No normal reset/table cleanup/password change/citext ownership or added normal superuser privileges. No commit/push/deploy/Phase8–14.

## 20. Safe repeat instructions and owner checklist

Normal inspection (read-only migration status) and launch if servers are not already running:

```bash
cd /home/marvin/projects/eLabTrack_V2
make migrate-status
# Only if the existing local servers are stopped:
make dev
```

Sign in normally at `http://localhost:5173/login`; Requests & Borrowings `/staff/requests`, Direct Issuance `/staff/requests/direct`. Ctrl+C on an owner-started make dev stops only API/frontend, not database. No migrate-up is needed again: normal is already000008. Never use restore/down/reset/remove volumes as a routine launch step.

For a repeat restore drill, use a **new, distinctly named empty disposable container/database/volume** and existing bootstrap-role script, PostgreSQL18.6 and random private0600 credentials. Positively verify project labels/name/port/volume/SQL identity and absence of application tables. Inspect installed pg_restore --version/--help and original --list; verify archive SHA256. The executed restore shape was:

```bash
# Example only after positively verifying a NEW empty disposable target.
# Credentials expand inside its container, never in host arguments.
docker exec -i "$verified_restore_container" sh -c '
  case "$POSTGRES_DB" in elabtrack_v2_phase7_restore_*) ;; *) exit 1 ;; esac
  PGPASSWORD="$POSTGRES_PASSWORD" exec pg_restore     --exit-on-error --single-transaction -h 127.0.0.1     -U "$POSTGRES_USER" -d "$POSTGRES_DB"
' < backend/tmp/pre-phase7-backups/elabtrack_v2-pre-phase7-20261009T150053Z-000007.dump
```

Do not point this at normal data or the already populated retained clone. Compare all restored tables, owners/effective ACLs/extensions/constraints/indexes/triggers/conservation, not only archive listing. Then run existing CLI from backend against explicitly guarded disposable credentials: --migrate-status (expect7/only8pending), --migrate-up, --migrate-status. Keep passwords in private config/container env. Repeat guarded synthetic regressions in the documented order on a **fresh test volume**, not historical rollback suites on a history-bearing retained target; [integration/PHASE7.md](../../integration/PHASE7.md) provides exact commands. Private recovery adapters are ignored under backend/tmp/phase7-local, including actual guarded comparisons, migration wrapper, secure session bridge and browser evidence scripts.

Owner manual acceptance remains unchecked:

- Existing Admin sign-in and existing2 Borrower accounts remain usable/present under their real eligibility.
- Equipment/category/catalog image and20/0/0/0/20 stock preserved.
- Requests directory loads200/empty; search/status/clear/refresh/pagination work.
- Direct form opens without bypassing activation/official acceptance.
- Review the isolated Borrower/Staff request/cancel/deny/checkout/direct/expiry history/screens and populated pagination; use only isolated fixtures for mutating manual scenarios.
- Verify themes/navigation/feedback/focus and absence of unexplained errors.
- Confirm official publication/activation dependencies stay enforced before actual normal borrowing.

## 21. Exact changes in this recovery

Project/status documents: this report; PHASE7_COMPLETION_REPORT.md; PHASE7_PROGRESS.md; ROADMAP.md; SOURCE_OF_TRUTH.md; OPEN_DECISIONS.md; DECISIONS.md evidence boundary; integration/PHASE7.md. Machine evidence: verification/phase7-local-integration/verification.json, restore-recovery.json, migration-status.txt, git-status.txt, changed-files.json. Visual evidence: phase7-local-integration/README.md, normal-acceptance.json, borrower-acceptance.json, staff-acceptance.json and38 exact PNG paths listed in [changed-files.json](verification/phase7-local-integration/changed-files.json).

No frontend/backend product/test source, migration SQL, approved mockup or shadcn primitive changed in this recovery. Required fmt made0 source changes. Existing intended Phase7 code is preserved uncommitted. Git-ignored additions/updates only for recovery: backend/tmp/phase7-local/recover.py, check-restored.py, finish-restore-check.py, verify-migration.py, restore-run.py, regression-setup.py, run-regressions.py, normal-browser.mjs, isolated-browser.mjs, existing secure-login.py and private target/schema/digest/result/env/compose artifacts. Compiled API and execution logs are under /tmp. Backup archive is retained, not a Git addition. Temporary credential handoff/fixture files removed.

## 22. Exact Git state and final status

Branch main, unchanged HEAD366b6ef4457378a2cdcf5d74aae9d1103e443637, index untouched. Final per-file status: 32 modified, 126 untracked, 0 staged. Exact final status and every path are [git-status.txt](verification/phase7-local-integration/git-status.txt) / [changed-files.json](verification/phase7-local-integration/changed-files.json); counts are recorded in machine verification. Eight unrelated diagnostic PNGs remain untracked/untouched. All private config/backup/helper paths remain ignored. No commit, push, deploy or history rewrite.

**PHASE 7 LOCAL INTEGRATION VERIFIED / OWNER ACCEPTANCE PENDING. Stop for owner review; Phase8 is not started.**
