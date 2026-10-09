# Phase 7 isolated owner demonstration — implementation report

**Formal owner acceptance, 2026-10-10 — DEC-081:** Phase7 Borrowing & Reservations is OWNER ACCEPTED following manual core-workflow review and automated safeguards. Student/Faculty mobile UI is accepted as functional Phase7, not final presentation; DEC-080 governs the future Phase11 browse/cart. Only a safe LOCAL checkpoint of verified Phase7/demo tooling/UX documents is authorized now. Do not begin Phases8–14, push or deploy. Official FSMO terms/current consent, approved production Student domains and required activation readiness remain live-use gates. Earlier pending-acceptance/correction statements below are historical. See [the checkpoint report](PHASE7_OWNER_ACCEPTANCE_CHECKPOINT.md).

2026-10-10. **PHASE 7 OWNER DEMONSTRATION READY / ENGINEERING VERIFIED / OWNER ACCEPTANCE PENDING.** This task enables manual acceptance of the existing Phase7 implementation; it does not implement Phases8–14. The latest owner correction withdraws the earlier Phase7 acceptance and pauses further advancement (DEC-079).

## Environment and delivered scope

| Surface | Isolated demonstration | Normal development |
|---|---|---|
|Frontend|http://127.0.0.1:15176/login|http://localhost:5173/login|
|API|127.0.0.1:18086|localhost:8080|
|PostgreSQL|127.0.0.1:54835 / `elabtrack_v2_phase7_demo`|localhost:5434 / `elabtrack_v2`|
|Container|`elabtrack_v2_phase7_owner_demo_postgres`|`elabtrack_v2_postgres`|
|Volume|Distinct `elabtrack_v2_phase7_owner_demo_*_pgdata` for each reset|`elabtrack_v2_pgdata`, never reset|
|Source/schema|Private committed snapshot `73f3fe4`, migrations000001–000008|Existing Phase7 source/schema retained|

The frontend is the real production-built Phase7 application, with a demo-only footer identifying fictional data. Requests use its real APIs, PostgreSQL, authorization, transaction services and inventory ledger. Normal frontend/backend product source and historical migration files have no diff. No mock authentication or frontend-only transaction flow was added.

Five legitimate fictional accounts comprise2 Students,1 Faculty,1 Staff and1 Admin. Students and Faculty remain BORROWER categories. Separate randomly generated passwords are bcrypt-hashed at cost12; trusted initialization is allowed only against the guarded empty disposable database. It records `DEMO_INITIALIZED_NO_EMAIL`, not email verification/delivery. `.invalid` addresses and the isolated Student domain are deliberately non-deliverable. Private configuration/credentials are ignored and restricted to0700 directories/0600 files; passwords travel to the seed process through stdin, never shell arguments. Owner retrieval is interactive-terminal-only.

Eight fictional equipment records span3 categories, with7 ACTIVE records (6 positive-stock and1 zero-stock) and1 INACTIVE positive-stock record. Eight distinct local PNG illustrations use the existing image validator/service/storage. The opening inventory is A114/R0/C0/D0/T114, including the Ladle A20/R0/C0/D0/T20. Seven nonzero opening acquisitions produce7 actual OPENING movements; the zero-stock pool correctly needs no quantity movement. Twenty inventory audit events and20 operation receipts retain category/create/image/status evidence. Final aggregate-to-ledger comparison reports zero mismatches.

Synthetic DEMO-1 is explicitly titled **DEMONSTRATION TERMS — NOT OFFICIAL FSMO POLICY**. The atomic seed records no consent. Browser tests exercise unchecked consent, real account/version/time receipts, DEMO-2 publication, rejected stale/current-missing acceptance and actual reacceptance. After scenario testing, a guarded reset leaves the owner dataset at DEMO-1, **zero accepted terms and zero borrowings**. Owner actions must create consent and borrowing history themselves.

## Safety and operational controls

`integration/phase7-demo.py` offers explicit setup/start/stop/status/credentials/expire/reset/build. Setup reuses an initialized dataset and refuses ambiguous partial seeding; initialization services join one database transaction. Guarded reset requires the exact database confirmation and validates configuration, disposable marker, actual database/container/Compose owner, loopback port and exact volume. It retains the previous volume/private recovery configuration and uses a new uniquely named volume. There is no `down -v`, normal reset, normal privilege grant or normal seed path.

Setup and multiple actual reset/start/stop cycles passed. A reset rebuilt successfully with `GOPROXY=off GOSUMDB=off`; ordinary start uses prebuilt assets/API and needs no downloads. The Docker image, Go modules and npm dependencies must be installed/cached before an offline presentation. Live Brevo configuration is blank in the demo API environment; no emails were sent. All imagery/fonts/authentication and business operations are local. Chromium blocks non-demo HTTP origins and records zero external requests/runtime exceptions.

The existing host-only refresh-cookie name is shared across ports. Demo URLs therefore use `127.0.0.1`, while normal URLs use `localhost`. The guide requires a private window/separate profile when normal development also uses127.0.0.1, and recommends one for every rehearsal. Authentication/cookie/rotation code is unchanged. The automated checks run in dedicated headless Chromium processes; final checks allow only the demo frontend/API origins.

Expiry uses an existing service with a test-only clock for a newly created fictional request25 hours in the past, followed by two real sweeps of its persisted deadline. It verifies EXPIRED/reserved0, retained automatic history and one-time release without editing an existing request or adding a production clock override.

## Executed demonstration verification — before formal acceptance

| Gate | Actual result |
|---|---|
|Go formatting, vet, unit tests, race tests|PASS: `go fmt ./...`, `go vet ./...`, `go test ./...`, `go test -race ./...`; root Phase7 source and private committed baseline checked|
|New fixture image validator test|PASS: `go test -race ./cmd/phase7-demo`,1 test, validates8 bounded canonical/distinct images|
|New Python safety tests|PASS:6 tests covering private permissions, normal database/volume rejection, exact reset confirmation, marker refusal and noninteractive credential refusal|
|Frontend serial|PASS:279 tests /23 files, `npm run test:run -- --maxWorkers=1`,65.00s|
|Frontend parallel|PASS:279 tests /23 files, `npm run test:run -- --maxWorkers=2`,56.32s|
|Frontend lint|PASS,0 errors;19 inherited shadcn/hook warnings retained, no suppressed checks|
|Frontend production build|PASS: TypeScript and Vite; root build and isolated demo build; quality build uses a separate output directory|
|Real PostgreSQL + HTTP regression|PASS:12 suites /70 named subtests under race detection on isolated54832; account, terms, stock, image, borrowing, contention, correction and pagination/eligibility contracts|
|Demo API negative/version checks|PASS:16 checks, including pending activation, inactive actor/target, missing/stale consent, inactive/zero/insufficient stock, unauthorized Staff/Borrower actions, strict evidence-field rejection and unchanged rejected-command stock|
|Real Chromium acceptance|PASS:24 checks across4 groups;44 successful light/dark mobile/tablet/desktop screenshots; actual login/logout/terms/request/reserve/cancel/deny/checkout/direct/history/expiry and fresh owner-ready identities|
|Final clean demo ledger/schema|PASS:8 verified migration checksums,5 active ready accounts,8 equipment/images, no borrowings/acceptances, zero ledger mismatches|
|Normal preservation|PASS for all20 compared business/migration tables; actual stock20/0/0/0/20; details below distinguish refresh-session activity|
|Whitespace and secret hygiene|PASS; final audit evidence is recorded with the exact changed-file inventory|

[PostgreSQL named test evidence](verification/phase7-owner-demo/postgres-regressions.json), [API checks](verification/phase7-owner-demo/api-eligibility.json), [migration checksums](verification/phase7-owner-demo/migration-status.txt), [final demo state](verification/phase7-owner-demo/final-demo-state.json), [normal preservation](verification/phase7-owner-demo/normal-preservation.json), [combined verification](verification/phase7-owner-demo/verification.json), [screenshot index](../ux/verification/phase7-owner-demo/README.md).

The browser sequence proves Ladle20 → reserve2 (A18/R2) → cancellation (A20/R0), deny3 with retained reason/release, request2 physically checked out, direct1 (A17/R0/C3/T20), and deterministic expiry with unchanged final physical total. Regression tests independently cover Faculty borrowing, rollback, authorization/status changes, duplicate keys, last-unit competition, stale transactions, restart, overdue obligations, immutable movement/audit history and untouched other custody buckets.

## Failures encountered and corrected

The first seed setup attempted a separate control schema that the existing migrator could not create. The identity marker was moved into its already-owned public schema; no privilege increase was introduced. The first isolated regression target lacked foundational runtime grants; the existing grant script was applied only to that guarded test database before the successful rerun.

An ordinary quality build initially overwrote the private demo bundle with default API configuration. Quality artifacts now use a separate output directory; the explicit demo build restored the correct API and indicator. Browser harness checks were corrected to use the visible consent label, actual API status codes/routes, the terms page's text Sign Out control and decoded authenticated blob-image elements. An API rerun hit the existing login rate limit because the harness authenticated before each request; the harness now reuses one legitimate in-memory token per account within its run. The rate limit and current-status checks remain intact, and a full rerun passed. Repeated version tests use the next labeled DEMO version. The API suite now asserts that the submitting Student has accepted current terms and that stock failures carry `EQUIPMENT_NOT_AVAILABLE`, preventing a terms rejection from masquerading as a stock test; it was rerun from a fresh dataset. Two read-only ledger verification queries failed on PostgreSQL row/subquery type handling; a scalar aggregate comparison passed. Diagnostic screenshots are retained and excluded from successful counts. No failing test was suppressed and no business/security behavior was weakened. Full browser workflows were rerun after the cookie-host isolation refinement.

## Normal data preservation

Before/after read-only row counts and digests match for20 business/migration tables. Normal data remains1 Admin,2 Borrowers both pending activation,0 terms versions/acceptances,0 borrowings and1 equipment pool A20/R0/C0/D0/T20 at000008. No normal account activation, password reset, policy publication, test loan, stock operation, migration or reset was executed.

The separate operational `refresh_tokens` table increased from65 to66 rows while the existing normal application remained running. This is disclosed separately; its complete contents are not claimed unchanged and were not restored/deleted. No demo command targets the normal authentication API or session table. This audit does not attribute that new session to a particular browser/user. The unchanged business-table comparison includes users/password hashes/account status and audit history without exposing any row data or hashes.

## Git and paused future work — historical pre-checkpoint state

Branch `main`, HEAD `73f3fe4 — Phase 7: Borrowing & Reservations`. No commit, staging, push, deployment or history rewrite was performed. Four existing project documents are modified; the demo controller/fixture/verification tools, guide/report, paused-progress record and sanitized evidence are untracked intended changes. [Exact porcelain status and changed-file inventory](verification/phase7-owner-demo/git-status.json) lists every path, including diagnostic screenshots. Private configuration, credentials, session material, logs, binaries, dependency/build caches, backup archives and retained demo data are excluded by existing ignore rules. No unrelated file was removed.

Before the correction arrived, partial Phase8 work had started under the earlier authorization. All11 affected files and their checksum manifest/diff are preserved privately at `backend/tmp/system-implementation/.paused-phase8/`; active source is restored to73f3fe4. Already launched isolated migration processes applied through000008 and failed candidate000009 with rollback; no normal migration ran. Phase8 is unverified and paused; that failed gate must be diagnosed if the owner later authorizes resumption. Phases9–14 were not started. No future return/replacement/fine/notification/report/demo-avatar behavior is part of this deliverable.

## Owner handoff and remaining gates

The demo is prepared and running at **http://127.0.0.1:15176/login**. Follow the [owner acceptance guide](PHASE7_OWNER_DEMO_ACCEPTANCE_GUIDE.md), retrieve the private passwords in your terminal and exercise scenariosA–G using separate legitimate logins. Only the owner can grant Phase7 acceptance after that review; this report does not grant it automatically.

There is no remaining engineering blocker for the scoped Phase7 demonstration. Official FSMO terms/publication/current consent, approved production Student domains and verified activation/ownership readiness remain live-borrowing gates. Live Brevo/production email verification is deferred and was not tested. A fully disconnected physical-device presentation, other browsers/assistive technologies and deployment were not performed; Chromium blocked external requests, local asset decoding and download-disabled rebuilding provide the executed offline evidence. The larger40-equipment/28-account/avatars/history dataset remains future Phase12 scope. FSMO's approved fresh-start direction remains; no V1 importer is planned or implemented.

**Owner Phase7 review is now formally accepted in DEC-081. Phases8–14 remain paused; checkpoint results are recorded separately.**
