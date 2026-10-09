# Phase 7 completion report

**Formal owner acceptance, 2026-10-10 — DEC-081:** Phase7 Borrowing & Reservations is OWNER ACCEPTED following manual core-workflow review and automated safeguards. Student/Faculty mobile UI is accepted as functional Phase7, not final presentation; DEC-080 governs the future Phase11 browse/cart. Only a safe LOCAL checkpoint of verified Phase7/demo tooling/UX documents is authorized now. Do not begin Phases8–14, push or deploy. Official FSMO terms/current consent, approved production Student domains and required activation readiness remain live-use gates. Earlier pending-acceptance/correction statements below are historical. See [the checkpoint report](PHASE7_OWNER_ACCEPTANCE_CHECKPOINT.md).

**Phase 7 local integration recovery, 2026-10-09: VERIFIED / OWNER ACCEPTANCE PENDING.** Full private-backup restoration succeeded using the existing bootstrap owner only on the isolated target. All17 restored tables and effective schema permissions match; only equivalent PostgreSQL metadata representations differ. Migration000008 passed restored-data rehearsal and was then applied through the existing normal-local migrator. All8 checksums verify; prior business data and physical20/0/0/0/20 stock are preserved. Normal directory200/empty/search/filter/direct-form checks, isolated workflows/races/restart,279 frontend tests in both modes, Go gates and22 Chromium checks pass. Official FSMO terms/acceptance, production Student domains and verified activation delivery remain live-use dependencies. Phase8 is not started. See [the local integration report](PHASE7_LOCAL_INTEGRATION_REPORT.md).

**Original implementation checkpoint; normal-local rollout is superseded by the recovery verification above.**

**2026-10-09 — PHASE 7 ENGINEERING COMPLETE / OWNER ACCEPTANCE PENDING.** Baseline `366b6ef4457378a2cdcf5d74aae9d1103e443637`, branch `main`. Explicit complete one-shot7A–D authority; engineering gates executed in isolated synthetic environments. Normal local PostgreSQL data was not migrated or reset during that original implementation run. No Phase8 work, commit, push, deployment or live email.

## 1. Phase 7A status

PASS. Paired000008, domain ports/entities/vectors, repository and transaction service implement one durable borrowing aggregate/items. Real PostgreSQL verifies conservation, last-unit/many reservations, multi-item failure rollback, current eligibility, duplicate/replay and reservation release. Empty up/down/reapply succeeds; retained-history down refuses. Historical14 SQL files000001–000007 match HEAD byte-for-byte.

## 2. Phase 7B status

PASS. Borrower own submit/reserve/review/pending/detail/history/cancel, strict API and replay,24h durable expiration/worker/lazy deadline handling, current terms and activation enforcement. Final Borrower Chromium:5 checks,13 screenshots. Actual submission reserves2; cancellation releases exactly2 with retained history. Request intent is local until authoritative submit.

## 3. Phase 7C status

PASS. Staff/Admin operational directory/details/search/status/borrower filtering, visible mandatory denial, explicit physical handover with future Manila due time, direct issuance and issued history. No separate APPROVED state. Real transition/concurrency/rollback/authority tests pass; final Staff/Admin Chromium:8 checks,18 screenshots. Confirmed physical checkout2 plus direct1 produces exactly A17/R0/C3/D0/T20 from the synthetic20-unit fixture.

## 4. Phase 7D status

PASS. Features use actual APIs and approved persistent shells/primitives. Real borrowing obligations integrate with account warnings as PARTIAL; actual installed borrowing liabilities integrate with inventory archive safety. Deactivation never resolves stock/loans/due/history. Frontend/backend/API/real PostgreSQL/migration/race/Chromium regressions pass. Domain, API, invariants, state machine, roadmap/source/decision and rollout docs updated. New-workflow owner acceptance remains pending.

## 5. Migrations and normal-data safety

New `backend/migrations/000008_borrowing.up.sql` and `.down.sql`. Adds aggregate/items/immutable events/receipts, indexes, original owned terms binding and reviewed custody-ledger source/vector extensions. Narrow runtime UPDATE grants for R/C and legal transition columns; D still forbidden. No historical checksum rewrite. All eight checksums verified on the isolated final target. Empty rollback/reapply and consequential-history refusal were actually executed.

Only named isolated DB `elabtrack_v2_batch1_test`, loopback54832, was used for mutations. Test containers own `elabtrack_v2_phase7_pgdata` and separate regression volumes (final `elabtrack_v2_phase7_regression_v3_pgdata`); older volumes retained. Normal `elabtrack_v2_postgres`/5434 was never a SQL/migration target during the original implementation run. Fresh separate volumes were needed because inherited historical rollback suites require empty later history; no volume was reset/deleted.

## 6. New API endpoints

Under `/api/v1`: GET `/borrowings`, GET `/borrowings/:id`, POST `/borrowings`, POST `/borrowings/:id/cancel`, POST `/borrowings/:id/deny`, POST `/borrowings/:id/approve`, POST `/borrowings/direct-checkout`. Borrower-owned reads, real scoped count/page; Staff/Admin operational reads/decisions/direct issue. All commands strict JSON, trusted Origin, UUID Idempotency-Key and current database authority. Exact DTOs/bounds/errors/replay outcomes: [API contracts](../API_CONTRACTS.md#phase-7-implemented-borrowing-contracts).

## 7. Frontend routes and reusable architecture

Borrower: `/borrower/borrowings`, `/borrower/borrowings/new`, `/borrower/borrowings/:id`; catalog/details link to new composition with optional equipment_id. Own history/details remain outside new-request terms gate. Staff/Admin: `/staff/requests`, `/staff/requests/direct`, `/staff/requests/:id`; account borrowing-history link supports operational borrower_id.

Feature pages → feature Query hooks → API adapter → existing centralized transport. RHF/Zod positive unique quantities/future Manila due/mandatory reason; local field state, frozen review/key, duplicate click guard, no automatic mutation retry and session-generation fence. Mutation/conflict invalidation covers borrowing/inventory/accounts; pending detail polls30s. Existing OperationalLayout/BorrowerShell/ManagementCard/forms/directories/ServerPagination/badges/AlertDialog/feedback retained. No UI primitive/approved asset changes or new dependencies.

## 8. Implemented state machine and exact vectors

| Operation | State/effect on physical `(A,R,C,D;T)` |
|---|---|
| Borrower submit | New PENDING; `(-q,+q,0,0;0)` |
| Owner pending cancel | CANCELLED; `(+q,-q,0,0;0)` |
| Staff denial + reason | DENIED; `(+q,-q,0,0;0)` |
| Deadline expiration | EXPIRED; `(+q,-q,0,0;0)` |
| Staff physical handover + due | PENDING→CHECKED_OUT; `(0,-q,+q,0;0)` |
| Staff direct handover + due | New DIRECT CHECKED_OUT; `(-q,0,+q,0;0)` |

A/R/C/D remain nonnegative and sum to T. Issued state/units/due persist for future returns; no checkout cancellation, arbitrary state PATCH, terminal deletion, APPROVED or completion edge. Current active equipment and target eligibility are checked at issue. No seven-day maximum, new quota or automatic fine prohibition.

## 9. Durable expiration

Database deadlines are24 elapsed hours after submission. Worker immediately catches up on API process start and runs every30s, at most100 candidates/cycle with10s budget, separate protected transaction per candidate. The ticker initiates reconciliation; persisted pending state/deadline is authority. Multiple workers/restarts use the same locks and conditional consumption; inactive borrowers' holds still expire. An actual compiled API was abruptly terminated while a synthetic request remained pending, restarted after its deadline, and restarted again: exact one-unit release and one EXPIRED movement were verified. This is process-level evidence in addition to fresh service/repository reconstruction tests.

Clock is captured after equipment locks; equality expires. Late approve/deny/cancel commits expiry/event/ledger/receipt, then returns409 BORROWING_EXPIRED; retries return that outcome without another release. Reads never free a persisted unprocessed hold. Worker outages/backlogs may delay release beyond the24h deadline until processing; startup catches up. Success/failure safe logs exist, but a dedicated lag dashboard/independent scheduler is not implemented.

## 10. Concurrency and history safeguards

Sorted participating user locks → current publication lock when needed → borrowing → sorted equipment; account authorization and target eligibility serialize against deactivation. Exact ledger vectors and conditional state/sequence changes, header/items/events/operator audit/receipt commit together. One normalized actor/op/resource/key/hash receipt prevents duplicate effects, detects changed payload and returns original committed evidence. Current actor authority precedes replay; submit/direct also recheck target eligibility. Historical receipt data may predate later transitions; detail is current.

Aggregate lock makes detail header/items/history coherent under competing checkout; list count/page use one SQL snapshot. Immutable events/item source effects, original borrower/equipment request snapshots and operator issue-time label audits preserve evidence. System expiry actorNULL accurately identifies automated action. Runtime cannot edit/delete audit/receipt/event history or D; down refuses consequential data.

## 11. Authorization, terms and account safety

BORROWER is authorization; STUDENT/FACULTY are categories. No Staff account creation privilege, SSO, borrower return authority or borrower evidence upload. Current roles/status, activated classified target, positive bounded items, active equipment, actual physical availability, mandatory due and handover are server enforced. Wrong Borrower detail is scoped404; unauthorized operators cannot issue. Existing auth/refresh/cross-tab/session mechanisms preserved.

Submission/direct call Phase4B RequireCurrentAcceptance inside the command transaction and bind its owned composite FK through commit. Direct never invents borrower consent. Existing pending checkout keeps its original accepted version; terms revision applies to new commitments. Missing terms503 and missing acceptance409 fail closed without creating holds. Normal unpublished official FSMO policy is not replaced by TEST content. Isolated activation tests use actual token lifecycle with a captured test adapter, not live delivery or real-account bypass.

Admin Student deactivation remains allowed regardless of obligations after warning/confirmation. Real tests prove unchanged issued quantities/due/stock/history; missing fine/replacement facts remain null. Borrowing-aware archive permits only clear unrelated equipment and blocks pending/issued custody; unintegrated replacement tables fail closed; future return processing must replace the issued-liability adapter.

## 12. Exact executed test results

| Gate | Result |
|---|---|
| Go fmt ./... | PASS; final formatting applied |
| Go vet ./... | PASS |
| Go test ./... | PASS; ordinary opt-in DB tests skip without isolated flags |
| Relevant Go -race packages | PASS domain/application/bootstrap/response/database/security |
| Guarded real DB/API tests |16 unique top-level suites,82 unique named subtests PASS across recorded runs |
| Final new borrowing PostgreSQL suites |3 suites/17 subtests PASS (foundation4, operations8, cross-feature5) with -race; additional actual API process-restart suite PASS |
| Phase7 HTTP final |1 suite/4 subtests PASS with -race |
| Existing real terms/accounts/inventory |12/10/8 named subtests PASS, fresh ordered target |
| Existing auth/terms real network |2 suites/12 subtests PASS against listening actual product API |
| Existing + Phase7 HTTP regressions |4 suites/20 subtests PASS including inventory borrowing filter |
| Frontend serial |23 files/279 tests PASS;268 baseline +11 new |
| Frontend bounded parallel --maxWorkers=2 |23 files/279 tests PASS |
| Frontend lint | PASS, exactly19 inherited warnings;0 new |
| TypeScript / production build | PASS via npm run build |
| Migration tests/status | Isolated fresh up, empty8 down/up, retained-history refusal,8 verified checksums PASS |
| Chromium |13 checks/31 screenshots PASS;0 runtime exceptions; no document horizontal overflow at tested widths |
| git diff --check | PASS |

The16-suite/82-subtest figure deduplicates suite/subtest names across reruns; it is not the sum of repeatedly executed passes. Exact names/log SHA256/evidence are recorded in [verification summary](verification/phase7/verification.json). Toolchains:Go1.27.1, Node24.19.0, installed Chromium1187 package. No dependency/toolchain downgrade.

Initial failures were investigated and corrected/rerun: fresh DB needed existing explicit foundation grants; later tables required empty migration8 removal before older rollback tests; historical terms TRUNCATE now also correctly encounters FK0A000; custody permission expectations changed to narrow approved R/C while preserving forbidden D. The first unrestricted23-worker frontend run under concurrent Go/DB load had277 passing and2 inherited lazy-login timeout failures. An additional serial rerun during concurrent work hit one of the same inherited lazy-login timeouts. Final quiet serial/two-worker full runs preserve assertions/timeouts and pass279; unrestricted23-worker stability is not claimed. New prefill-effect lint warnings were removed through mount-time initialization, without suppression. Browser test UTC-string and signed-in-role-switch assertions were corrected to instant comparison and actual UI logout. Read coherence found during final audit was fixed and independently race-tested.

## 13. Browser verification

Actual Chromium exercised real login, catalog/detail selection, quantity/review/single submit, pending hold/cancel/history; Staff denial/review/reason, physical checkout/due/review, existing eligible target direct issue, exact API stock/ledger readback, Admin navigation/keyboard, Borrower denied/issued/expired histories and real persisted sweep. Browser role switches perform normal UI logout; no auth store role injection. Responsive widths390/768/1024/1440, both themes and Tab focus. Synthetic names/emails only; no secrets in evidence.

[Borrower checks](../ux/verification/phase7/borrower-acceptance.json) and [Staff/Admin checks](../ux/verification/phase7/staff-acceptance.json) retain sanitized request paths and exact screenshots. Screens are viewport captures so fixed header/navigation positions accurately represent the browser; original full-page capture artifacts were corrected in the harness. Physical assistive technology and non-Chromium engines remain unrun.

## 14. Screenshot evidence

All31 representative captures are under [docs/ux/verification/phase7](../ux/verification/phase7/). Borrower13: light/dark catalog, equipment detail, mobile/tablet quantity selection, mobile/desktop review, pending, cancellation confirmation/released history and responsive directory. Staff/Admin18: directory, denial form/review, checkout form/review, issued details, direct form/review, stock/ledger after issuance, Admin directory and Borrower denied/issued/expired mobile histories. Both acceptance JSON files list exact names. PNGs contain only synthetic fixtures; new owner visual approval is not inferred.

## 15. Remaining external dependencies

Approved official FSMO terms after presentation, assigned publication/version responsibility and each borrower's current documented consent before live borrowing. Authoritative production Student domain/roster inputs. Required activation/ownership operational readiness, backend Brevo key, verified sender and successful actual delivery testing. Backend provider submission is not delivery proof. Normal Borrowers may remain blocked by these gates; isolated TEST consent and activation do not approve production readiness.

## 16. Limitations and deferred work

Owner manual/visual Phase7 acceptance pending. Normal local schema rollout deliberately unexecuted. No return/damage/loss/replacement, completion, fine assessment/clearance/payment, recovery endpoint, notifications/live delivery, full reporting or expanded kiosk. Fine basis metadata and durable transition evidence prepare future work only. PARTIAL account projection cannot report future fine/replacement balances; future modules must replace adapters and extend tested migrations/ledger/authorization.

Bounded expiry throughput/outage delay, no dedicated lag dashboard; real physical handover is explicitly confirmed by Staff/Admin rather than hardware verified. Directory search uses retained request identities, while immutable operator audit preserves issue-time labels.19 inherited lint warnings and unrestricted parallel lazy-login sensitivity remain. Original unmatched GET500 remains without exact path/request-ID evidence; no new unknown500 defect was established here. No production/edge deployment, non-Chromium, physical assistive-device or normal-data restore verification.

## 17. Exact changed files and preserved exclusions

[Changed-files manifest](verification/phase7/changed-files.json) lists every intended modified/new file relative to HEAD, including new directories' individual files and screenshots. Related historical integration tests are narrowly adapted for installed000008/allowed custody privileges; no assertions were disabled. No UI primitives, approved assets, historical migration, dependency manifest or production config changed.

Task-owned verification API/Vite/browser harness and isolated database containers were stopped; private browser fixture/sentinel files were removed, while all volumes were retained. Eight pre-existing diagnostic PNGs remain untouched/untracked: reconstruction-stage-a/failure.png and reconstruction-stage-b/{administration,borrowers,categories,consistency,edges,final-visual,stock}-failure.png. Ignored `.env`/`.env.batch1`, bootstrap utilities, runtime build/dist/cache, private test fixtures and retained DB volumes are excluded. No normal account or credential is overwritten. No staging/commit occurred.

## 18. Exact Git status

[Git status snapshot](verification/phase7/git-status.txt) records exact final `git status --short`, branch and HEAD; [manifest](verification/phase7/changed-files.json) expands grouped untracked directories. All Phase7 work remains reviewable and uncommitted, plus the eight preserved unrelated diagnostics. Index unchanged/empty diff; no Git history rewrite, push or remote changes.

## 19. Local migration instructions

[Isolated verification and safe owner-run rollout](../../integration/PHASE7.md#historical-normal-local-rollout-plan--subsequently-executed-after-recovery) contains the full procedure. Review first; stop normal API; confirm private local target; make/validate a new private backup; owner runs `make migrate-status`, checks existing checksums/pending000008, explicitly runs `make migrate-up`, then `make migrate-status` and restarts `make dev`. No automatic migration or reset. Never seed TEST terms/acceptance or artificially activate real users. History-bearing rollback refuses; recovery needs a separately reviewed restore plan. This report provides instructions, not evidence that normal rollout occurred.

## 20. Owner acceptance checklist and handoff

- [ ] Review new light/dark Borrower and Staff/Admin screenshots and representative responsive/keyboard interactions.
- [ ] In an authorized isolated acceptance environment, verify multi-item selection, review, one submission and exact reservation.
- [ ] Verify own-only history, cancellation release, mandatory denial reason and durable expired history.
- [ ] Verify physical handover confirmation, future Manila due time, no seven-day maximum and one checkout.
- [ ] Verify direct issuance to an eligible existing Borrower and exact available/reserved/checked-out quantities/ledger.
- [ ] Verify current role/status changes, deactivation warnings/unchanged obligations and safe conflicts/retries.
- [ ] Verify missing official terms/acceptance/activation blocks new commitment, with existing history access retained.
- [ ] Separately approve/review normal local rollout and new backup; confirm no synthetic policy/accounts enter normal data.
- [ ] Record official terms, production domains and verified delivery readiness before any live borrowing claim.

**Next action:** Owner review/acceptance of Phase7, then separately reviewed normal local migration if desired. Engineering does not approve institutional content or deployment. Stop after Phase7; Phase8 is NOT STARTED.
