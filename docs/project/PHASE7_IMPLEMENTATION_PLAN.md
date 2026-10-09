# Phase 7 implementation plan

Owner authorization: COMPLETE Phase 7 one-shot implementation, 2026-10-09. Baseline 366b6ef. This plan precedes feature code; implementation proceeds through verified 7A–D without routine approval. Phase 8 is excluded.

## Approved scope and reuse

Students and Faculty are BORROWER. Reuse current account locks/activation, authentication/refresh/cross-tab coordination, Phase4B current acceptance within the command transaction, TxManager, equipment repository, physical stock model, immutable ledger and paired migration runner. Use one `borrowings` aggregate with item rows, REQUEST/DIRECT entry paths; do not duplicate requests and loans. Existing pending submissions retain their original acceptance on checkout. New submission/direct issuance require current acceptance. No automatic fine eligibility veto.

Reuse persistent OperationalLayout, BorrowerShell/catalog/details, ManagementPage/Card/FormField/WorkflowSteps, directory tables/search/select, ServerPagination, feedback, shadcn/Base UI, RHF/Zod, TanStack Query and centralized transport. Local request intent does not reserve stock before submission. No kiosk expansion or redesign.

## Domain/state and database

PENDING → CHECKED_OUT (approval and confirmed physical handover), DENIED (visible reason), CANCELLED (owner), EXPIRED (24 elapsed hours). Direct issuance begins CHECKED_OUT. No APPROVED, returns, completion mutation, fines, payments or replacements. Issued records retain exact issued units and due timestamp for Phase8. Persist timestamptz, render Asia/Manila; future due required, no seven-day maximum.

Paired 000008 adds aggregate/items, immutable transition events and actor/operation/key/hash receipts; owned acceptance composite FK; unique borrowing/equipment lines; positive bounded quantities, path/state/timestamp constraints, history and pending-expiry indexes. Extend existing inventory movement kinds with source item FK and unique source/effect; reserve/release/issue vectors preserve total and damaged stock. Narrow R/C update privileges, retain existing history immutability; down refuses consequential data. Fine policy basis is immutable issuance metadata for later work, not fine assessment. Durable transition events provide future notification source evidence; no notification delivery/outbox worker in Phase7.

## Transaction and retry protocol

Normalize/sort distinct equipment IDs; reject duplicate lines. Lock participating users sorted, publication shared when current acceptance required, existing borrowing, equipment sorted. Capture database clock after locks; conditional PENDING transition and vector updates, immutable history/movement and receipt commit together. Reauthorize actor and target on replay; identical normalized key replays original identity, different payload conflicts. Late cancel/deny/approve commits EXPIRED and release, then returns BORROWING_EXPIRED; never roll back expiry to signal conflict. Account deactivation never resolves borrowing or changes stock; expiry still releases inactive borrowers' pending holds. Stock arithmetic remains server owned.

Sweeper selects a bounded indexed set without locking aggregates first, then uses account→borrowing→equipment order per candidate. Immediate restart catch-up and periodic sweep, bounded context, concurrent-safe conditional state. No in-memory-only timers, Redis or broker. Log safe sweep counts/failures. Reads display stored stock until committed release. Archive remains fail closed with future unintegrated liabilities; no invented zero replacement/fine balances.

## REST contracts

`GET /borrowings` scoped own/operational directory: page/per_page/status/search, stable created_at/id, real scoped count. `GET /borrowings/:id` owned/operational details with chronological immutable history. `POST /borrowings` own unique equipment_id/quantity lines, confirm; server reserves. `POST /borrowings/:id/cancel` confirm. `POST /borrowings/:id/deny` reason/confirm. `POST /borrowings/:id/approve` due_at/physical_handover_confirmed/confirm. `POST /borrowings/direct-checkout` borrower_id/items/due_at/physical_handover_confirmed. All commands trusted Origin, strict JSON and UUID Idempotency-Key; no generic state PATCH or DELETE. Standard safe response envelope; 400 validation, 403 authority/eligibility, scoped404, 409 stock/state/key/expiry, existing503 unpublished terms. No browser automatic mutation retries.

## Screens and integration

Borrower catalog/details link into composition/review; own borrowing directory/details/cancel/history outside current terms read gate. Submission remains terms gated and server verified. Staff Requests & Borrowings directory/details, deny review, physical checkout review with explicit future due selection, direct issuance to an existing active/activated borrower. Dates explicitly Manila (+08:00 conversion). Display live availability and stale/conflict feedback; frozen reviewed payload and stable key prevent duplicate submissions. Query actor keys/session fences, invalidation of borrowing/inventory/account directories after commands; no permanent role caching.

## Milestones and gates

| Milestone | Delivery | Gate |
|---|---|---|
|7A|Domain,000008,repositories,atomic reservation foundation|Domain, empty up/down, isolated PG conservation/concurrent reservation/rollback/authority|
|7B|Submit/cancel/history/replay/expiry worker, Borrower UI|API, expiry/restart/retry, terms/activation, frontend and actual Borrower Chromium|
|7C|Staff denial/physical checkout/direct issue and UI|Real PostgreSQL transition races/direct contention/rollback, API/authority/frontend/Chromium|
|7D|Cross-feature wiring, regression, evidence and docs|fmt/vet/tests/race, isolated PG/migrations, serial+parallel frontend tests/lint/type/build, responsive/theme/keyboard Chromium, diff check|

## Test matrix and acceptance

Two borrowers last unit, many reservations, multi-item partial conflict, cancellation/checkout, expiry/checkout, denial/cancel, direct/reserve, replay/change-key conflicts, duplicate checkout, transaction interruption/audit rollback, restart/repeated expiry, boundary time, invalid quantity, inactive equipment/borrower, pending activation, missing terms/acceptance, role revocation and ownership. Assert exact A/R/C/D/T, ledger sums, one legal history outcome and no orphan holds. Existing Phase5/6/auth regression retained. Screenshots use only synthetic fixtures in isolated database, both themes and mobile/tablet/desktop; owner visual acceptance remains pending.

## Dependencies, risk and rollout

No unresolved borrowing policy blocks this authorized scope. Official FSMO terms await approval after presentation; production Student domain approval and verified Brevo sender/key/delivery remain external. Isolated synthetic terms/users only; normal borrowers may remain blocked. No live email. High complexity lies in lock ordering, receipt replay, lazy expiry commit and cross-feature stock authority; independent tests precede UI integration. StageB visual approval does not approve new screens. Existing account fine/replacement projection remains unavailable, accurately labelled.

Verify000008 only in disposable database first. Never automatically migrate normal local database. Document explicit backup/status/up/status procedure for owner-run normal rollout after review. No reset, volume deletion, production config, commit/push/deploy. Stop after Phase7 engineering report for owner acceptance.

## Verified implementation refinements

Final review adds aggregate locking for coherent detail header/item/history reads, with real concurrent-read coverage. Directory records return empty detail arrays rather than nullable arrays. Catalog-prefill initializes the form at mount after loading eligible equipment, avoiding a state-copy effect and new lint warnings. Operator-led inventory audits retain issue-time labels separately from original request snapshots. Actual borrowing obligations are PARTIAL; fine/replacement facts remain unavailable. Installed borrowing-aware archive checks permit only unrelated clear equipment, while future unintegrated liability modules remain closed. These are integrity/architecture refinements within the approved policies, not new borrowing behavior. See progress/completion for executed gates; owner acceptance remains separate.
