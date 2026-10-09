# Phase 7 formal owner acceptance and local checkpoint

2026-10-10. **PHASE 7 OWNER ACCEPTED — DEC-081.** The owner formally accepted Borrowing & Reservations after manual core-workflow review and automated engineering safeguards. The Student/Faculty mobile interface is accepted as functional Phase7, not final presentation. [DEC-080](../ux/PHASE11_EQUIPMENT_SELECTION_UX.md) governs future equipment cards, adjacent quantity/Add, desktop browse/cart and mobile cart/review. **Phases8–14 remain paused.** This task authorizes only the safe LOCAL checkpoint; no push or deployment.

## Checkpoint scope and review

Before this checkpoint, branch `main` was at `73f3fe4 — Phase 7: Borrowing & Reservations`. That commit already contains the verified application and migration000008. This checkpoint adds the isolated owner-demo controller, trusted fictional seed/image fixtures, safety/API/Chromium verification tools, sanitized evidence, owner guide/report, DEC-080 reference comparison and current acceptance/status records. Application source, dependencies, authentication, borrowing rules and historical migrations are unchanged.

All modified and visible untracked paths were inventoried and reviewed before explicit staging. Text/source review and exact-value scanning found no candidate containing any of72 distinct private credential/secret values read without displaying them. Additional checks covered JWTs, bcrypt hashes, private keys and API-key patterns. Apparent email syntax in a Go toolchain module path was identified as a path, not an account. Real email addresses were absent from the reviewed text; `.invalid` fictional demo identities remain intentional fixtures. All47 PNGs were visually inspected through private contact sheets and checked for embedded metadata:44 successful workflow images and3 retained, labeled failure diagnostics. These curated verification artifacts contain fictional identities and masked password input; private runtime output is excluded.

The sanitized [checkpoint evidence](verification/phase7-owner-checkpoint.json) records current test results, file inventory and normal-preservation status. The commit containing this report is the local checkpoint; obtain its exact hash with `git log -1 --format='%H %s'`. Earlier [demo evidence](PHASE7_OWNER_DEMO_IMPLEMENTATION_REPORT.md) and its pending-acceptance timestamps remain historical; the current owner decision is recorded separately rather than rewriting test outcomes.

## Fresh verification for this checkpoint

| Check | Actual result |
|---|---|
|Backend formatting/vet/unit/race|PASS: `go fmt ./...`, `go vet ./...`, `go test ./...`, `go test -race ./...` with Go1.27.1; database suites explicitly enabled separately|
|Demo illustration fixture|PASS:1 race-enabled test in the private committed Phase7 snapshot;8 canonical bounded/distinct images|
|Demo utility safeguards|PASS:6 Python tests; private permissions, normal-target/volume rejection, confirmation/marker guards and noninteractive credential refusal|
|Frontend tests|PASS:279 tests across23 files, serial unchanged suite,70.06s|
|Frontend lint/build|PASS:0 lint errors,19 inherited warnings; TypeScript/Vite production build in a private temporary output directory|
|Real PostgreSQL/API regressions|PASS:12 suites,70 named subtests across5 groups, race-enabled and uncached; terms/accounts/inventory/borrowing/HTTP on a fresh isolated54832 target|
|Real Chromium/API smoke|PASS: mobile390px and desktop1440px demo login, hidden password input, demo indicator, no overflow; public API readiness200. No external requests or runtime exceptions; no credential entry/consent/borrowing mutations|
|Normal database preservation|PASS:20 business/schema tables' row counts/digests unchanged in read-only before/after comparison; refresh-session rows68→68. No database row contents or digests are published|
|Repository/documents|PASS: historical SQL/source comparison, local document links, secret/privacy audit and unstaged/staged whitespace checks before commit|

The previous complete demo verification remains12 PostgreSQL suites/70 subtests,16 API eligibility/version checks and24 Chromium workflow checks with44 successful screenshots. Full mutating browser scenarios were not rerun for this checkpoint because source is unchanged and the owner's existing demo scenarios must be preserved. Fresh smoke checks do not replace or inflate those historical workflow counts.

The first database attempt safely refused the retained test database because it already had terms history. No test guard was bypassed and no history was deleted. Only that isolated test container was stopped; its container/volume were retained. A new distinct checkpoint test volume was initialized through the existing bootstrap/migrator, with only000001–000008 and existing runtime grants. All five regression groups then passed. Normal5434 and owner-demo54835 were untouched. Chromium initially lacked the default executable/shared-library lookup; the successful checks used the already cached headless shell and private library bundle, without installation/download or application changes.

## Exclusions and preservation

- `elabtrack_phase11_v1_references.zip`: retained locally and narrowly ignored; its legacy background includes personal/sample records. The safe specification retains exact member dimensions/hashes and visual facts; the archive is not an application asset or staged attachment.
- `.project-reference/`: private legacy/capstone references, including personal material.
- Root/backend `.env*` private configuration and local integration credentials; only already tracked example contracts remain in the repository.
- `backend/tmp/`: demo passwords/configuration/session material, backups/database dumps, logs, binaries/private baseline snapshots, retained recovery configurations and paused Phase8 source/manifest/patch. No excluded file is deleted or exposed.
- `frontend/node_modules/`, build/cache directories and `/tmp/elabtrack-phase7-checkpoint/`: dependencies, generated builds, raw logs, digest snapshots and audit contact sheets. Sanitized verification records and reviewed screenshots are intentional engineering evidence, not private runtime artifacts.
- Existing normal and owner-demo database records/volumes: no reset, migration, activation, password change, consent publication or borrowing/stock mutation in this task. Earlier retained disposable volumes remain available.

Staging uses an explicit reviewed path allowlist; no blanket `git add .` and no unrelated change is discarded. Git history is preserved and the existing credential-free origin is unchanged. The local checkpoint does not grant push/deployment or future-phase authority.

## Remaining limitations and handoff

Official FSMO terms wording/publication and documented current acceptance still gate live borrowing. Production Student-domain inputs and required activation/ownership readiness remain external dependencies; live Brevo delivery/production email verification are deferred. Demo synthetic terms and trusted fictional initialization are not institutional approval or verified email delivery.

Returns/replacements/fines/notifications/reports and the larger40-equipment/28-account presentation dataset remain future scope. The privately preserved partial Phase8 migration previously failed and rolled back in isolation; that failed gate is unresolved and must be diagnosed only after separate advancement authorization. DEC-080's final Phase11 presentation is unimplemented and has no rendered/runtime acceptance.

Inherited lint warnings, original unmatched GET500 lacking exact request evidence, worker-backlog limitations, non-Chromium/physical assistive-device testing, a fully disconnected physical-device presentation and production deployment remain disclosed. This checkpoint neither introduces nor resolves them.

**Next action:** review the local checkpoint and, when ready, separately authorize the next phase. Do not begin Phases8–14 under this checkpoint task.
