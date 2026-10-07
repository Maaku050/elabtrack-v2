# Phase 2 boss rebuild design reference

Read-only review, 2026-10-08. **BOSS DESIGN REFERENCE**: `/home/marvin/projects/elabtrack/docs/audit/BOSS_REBUILD_AUDIT.md`, audit dated 2026-10-07 at boss HEAD `42b8c7c5d8b10ba13de0c7c53dfd007cbcd50cfc`. Its current-V2 comparison is historical: it predates completed Phase 1H/1I and must not reopen resolved foundation defects. Selected source was inspected locally read-only. No file, migration or UI component was ported or modified there.

## Salvage map

Classification is an **ENGINEERING RECOMMENDATION for later approved work**, not an assertion of accepted policy or implemented reuse. ADOPT CONCEPT means use its idea in the design, not copy code. Paths below are relative to the boss repository, not current V2 source.

| Concept | Classification | Source evidence | FSMO adaptation / safeguard |
|---|---|---|---|
| Explicit stock buckets | ADAPT | `backend/internal/domain/inventory/entities.go`; `backend/migrations/000006_equipment_inventory.up.sql`; audit §6 | Six proposed buckets, accounted-total semantics; no lab FK/tenancy; OPEN-022 confirms loss/retirement treatment; cross-row reconciliation required |
| Append-only movements | ADOPT CONCEPT | Inventory domain/repository; audit §§6/15 | Current row authoritative, signed evidence/sequence reconciles; no event sourcing; reason/actor/source for every movement |
| Ordered row locks | ADOPT CONCEPT | `backend/internal/application/borrowing/service.go` lockEquipmentSorted; `.../persistence/postgres/inventory_repository.go` FindByIDForUpdate | Deduplicate equipment IDs, UUID ordering; extend order to users/profiles/borrowing/charges, status writers; all participating repos in one transaction |
| Optimistic metadata version | ADOPT CONCEPT | `.../inventory_repository.go` Update/version predicate | Narrow metadata fields, separate immutable snapshots and stock guard; stale edits conflict |
| Retained canonical terminal borrowing | ADOPT CONCEPT | `backend/internal/domain/borrowing/entities.go`; audit §§5/15 | Denied/cancelled/completed IDs stay; no copied/deleted active/history records |
| Separate approval/physical release | ADAPT | `backend/internal/application/borrowing/requests.go`; audit §5 merges pending→checked_out | Preserve separate decision/release evidence; prefer separate approved state, conditional combined command only with OPEN-021 confirmation |
| Historical snapshots | ADAPT | `.../borrowing/requests.go`, request capture; audit §7 | Define request versus issue binding OPEN-023; freeze borrower display/ID as well as item metadata/value; don't mislabel request price as checkout capture |
| ReturnEvent / ReturnEventItem | ADAPT | `.../application/borrowing/returns.go`; domain ApplyReturn; migration 000008; audit C01 | Reject input duplicate IDs; UNIQUE(event,item), composite same-borrowing FKs; single normalized map drives every effect |
| Assessment + append-only adjustments | ADAPT | `backend/migrations/000008_returns_charges.up.sql`; `.../application/borrowing/charges.go`; audit §14 | Immutable base protected at DB as well as APP; derived zero balance; same-currency/source uniqueness; settlements distinct from payments; no boss financial defaults |
| Policy/terms versions | ADAPT | Borrowing policy/terms in domain/migration 000007; audit §§7/14 | Single FSMO immutable typed versions; actual approved text; binding/acceptance/time/rates not inherited; no organization hierarchy/CMS |
| Business audit | ADOPT CONCEPT | service auditEvent and audit schema/repository; audit §15 | Meaningful before/after for custody/count/value/status changes; atomic required writes; protected minimal actor/entity context, not logs |
| SMTP transactional outbox | ADAPT | `.../application/borrowing/service.go` emit; `.../persistence/postgres/notification_repository.go`; audit §16 | Unique ReturnEvent keys, bounded claims, lease-token-fenced completion, preserved persistence failures/attempts; delivery separate from acceptance; future worker/deployment |
| SQL report pagination/filtering | ADAPT | `.../persistence/postgres/reporting_repository.go`; audit §17 | Single FSMO approved reader/owner scopes; unique borrower rows; validate all filters/search/date range; shared calendar/currency; no fabricated revenue |
| Small CSV sanitization concept | DEFER | `backend/internal/application/reporting/csv.go` / tests; audit §24 | Useful future export detail only after approved formats/permissions; no export feature now |
| Device/handoff/narrow kiosk DTO/idle reset | DEFER | `backend/internal/application/kiosk/service.go`; `frontend/src/features/kiosk/`; audit §13 | OPEN-010 first; same canonical borrowing; independent privacy/session shell; reject late-result/reset race, first-page-only browse and missing image auth |
| Transaction-time actor freshness | ADAPT | `backend/internal/application/authz/fencer.go`; audit §§8/10 | Inward current-account/eligibility ports and shared lock protocol; do not replace verified Phase 1 tokens/cookies/client |
| Provenance hashes/mapping/quarantine | ADAPT | Audit §20 and `backend/internal/migration/{export,importer}.go` as defect evidence | Future transactional mappings over real V1 multi-item shapes; immutable unknowns; safe non-writing/isolated dry run; no importer reuse |
| Hierarchy/global Super Admin/scoped role catalogs | REJECT | migration 000003 identity; admin organization UI; audit §3 | Unsupported FSMO scope; no organizations/units/labs/global-super-admin platform |
| Duplicate-return processing | REJECT | C01; domain ApplyReturn accumulates duplicates, ProcessReturn stock loop uses first match | Counterexample explicitly designed out and future regression required, not accepted ordinary-path correctness |
| Target borrower eligibility bypass | REJECT | `.../borrowing/requests.go` DirectCheckoutExt checks scopes, not account/eligibility; audit H03 | Current actor authority never substitutes for current target eligibility at release |
| Generic file/export authorization | REJECT | `.../application/inventory/files.go` Open; audit H02 | Equipment images only in narrow model; private export ownership/current privilege/expiry rechecked through its purpose-specific route |
| Unsafe V1 importer/migration numbering | REJECT | `backend/internal/migration/importer.go`, cmd/migrate-legacy; audit H08 | No flattened single-item mapping, fabricated NOW dates, partial claims, writing dry-run, copied SQL pairs/numbers |
| Repeated partial-return notification key | REJECT | borrowing service emit + BorrowingEventKey; audit §16 | Event-specific IDs, not one borrowing-returned key |
| Due/grace/report inconsistencies | REJECT | audit §§14/16/17; domain OverdueDays versus scheduler/report SQL | One approved cutoff/day/grace contract; date ranges half-open, grace distinct from custody lateness; typed currency |
| Whole auth/config/UI foundation replacement | REJECT | audit §§8/12/20/24 | Current tested Phase 1 stays; shared primitives already present; no raw startup migrations, unguarded seed, plaintext DB policy or copied showcase/binaries |
| Async export jobs/million-row materialization | DEFER | reporting service/repository; audit §17 | No requirement/measurement justifies jobs/storage/large memory path; later bounded formats first if approved |

## Duplicate-return counterexample and correction

Source-derived audit C01: issued item quantity2, raw payload repeats that item ID with good1 twice. Boss domain accumulates good2 and completes; stock loop takes first matching line and restores1, leaving one ghost checked-out unit. Each pool CHECK still passes because totals remain locally consistent. Audit describes a trace, not an executed exploit, and Phase 2 likewise does not execute it.

V2 recommendation rejects the duplicate payload **before mutation**, even if summing would fit. DB UNIQUE(event,item) adds independent prevention; composite FKs prevent a line from another borrowing. Only one validated map feeds event quantities, cumulative custody, stock deltas, movement rows and assessment basis. The [invariants](INVARIANTS.md) require exact cross-reconciliation, concurrent same/different-key return tests and all-or-nothing rollback after the first equipment line. This corrects an interpretation defect; locking alone would not.

## Evidence limits and port guardrails

Actual selected source inspection covered inventory entity/repository locking/versioning, borrowing ApplyReturn/ProcessReturn/DirectCheckoutExt/idempotency/emit, return/charge migration structure, charge application, notification claiming/attempt persistence, report SQL, file Open and kiosk service DTO/device concepts. Other rows cite the full audit's path/section; they do not claim a fresh runtime proof or exhaustive source re-audit. Boss unit/build results and skipped DB tests remain the audit's historical results, never current V2 business acceptance.

Do not copy boss migrations into current numbering: the existing runner history/role contract is separately verified. Future approved schemas must be authored against the current users/session tables, constraints, checksums and named-object runtime grants. Each later selected code reuse needs provenance, a reviewed contract and targeted defect/real PostgreSQL tests. Audit 24's salvage classification is a recommendation, not phase authorization. No UI reuse before Phase 3A approved mockups and Phase 3B shell/design system.
