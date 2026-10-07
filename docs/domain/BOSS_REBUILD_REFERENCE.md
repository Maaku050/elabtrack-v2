# Phase 2.5 boss rebuild design reference

Read-only historical reference, reclassified 2026-10-08. `/home/marvin/projects/elabtrack/docs/audit/BOSS_REBUILD_AUDIT.md`, dated2026-10-07 at HEAD `42b8c7c5d8b10ba13de0c7c53dfd007cbcd50cfc`, is **BOSS DESIGN REFERENCE**, never FSMO product authority. Phase 2 inspected selected source; this reconciliation changes no boss file or current Phase 1 implementation. Historical foundation comparisons predate completed Phase 1H/I.

## Rebaselined salvage map

These are selected engineering concepts for later approved work; ADOPT CONCEPT is not code reuse or implementation acceptance.

| Concept | Classification | Historical source evidence (boss-relative) | Current FSMO design / safeguard |
|---|---|---|---|
| Current stock buckets | ADAPT | inventory entities; migration000006; audit§6 | Four physical counts A/R/C/damaged_held; lost is historical incident, replacement separate. Reject six-bucket accounted-total as current stock semantics |
| Append-only stock movements | ADOPT CONCEPT | inventory domain/repository; §§6/15 | Explicit baseline, exact signed vector, sequence/source/actor; current row authoritative, no event sourcing |
| Ordered equipment locks | ADOPT CONCEPT | borrowing/service.go lockEquipmentSorted; postgres/inventory_repository.go | Extend to account/policy/borrowing/obligation/equipment/fine shared order; deduplicate IDs; one transaction |
| Metadata version/snapshots | ADAPT | inventory_repository Update; requests.go; §7 | Narrow metadata edits; request/issue binding and immutable user/item snapshots; no price-based damage fine |
| Canonical terminal history | ADOPT CONCEPT | borrowing entities; §§5/15 | DENIED/CANCELLED/EXPIRED/COMPLETED remain; no copied record/delete |
| Combined approval/issue | ADAPT | borrowing/requests.go; §5 pending→checked_out | Owner D5 confirms single physical edge; no former separate-approved proposal; current target/hold/due checks |
| ReturnEvent/ReturnEventItem | ADAPT | borrowing/returns.go; ApplyReturn; migration000008; C01 | Unique input map and stored event/item, composite parent ownership, exact effects; add replacement obligations/events |
| Generic assessments/adjustments | REJECT FOR CURRENT SCOPE | migration000008, borrowing/charges.go; §14 | Replace with overdue_fines/fine_clearances; PHP 10 dynamic/final, full Admin clearance; no damage-price component/allocation/reversal default |
| Policy/terms versions | ADAPT | domain/migration000007; §§7/14 | Single FSMO typed immutable policy and user/version acceptance; no per-loan checkbox or organization hierarchy |
| Business audit | ADOPT CONCEPT | service auditEvent; §15 | Required actor/entity/before-after/source in same transaction; logs distinct |
| Transactional outbox/attempts | ADAPT | service emit; notification_repository; §16 | Distinct submission/issue/expiry/partial/replacement/completion IDs; token-fenced leases, observable failure; provider deferred |
| SQL report paging/filtering | ADAPT | reporting_repository; §17 | Owner/role scope, bounded validated filters, same clock math; PAID versus waived/other, no assessed-as-revenue |
| Device/handoff kiosk model | REJECT | kiosk/service.go; frontend/features/kiosk; §13 | Owner D28/29: ordinary borrower catalog/cart is Interactive Kiosk. No device registration/credentials/handoff/privacy shell infrastructure |
| Transaction-time actor freshness | ADAPT | authz/fencer.go; §§8/10 | Current account/target check under shared locks; no new eligibility table or change to Phase 1 sessions |
| Provenance/migration reconciliation | ADAPT CONCEPT, DEFER IMPLEMENTATION | migration/export/importer; §20 defect evidence | Future multiitem source mapping/quarantine; no importer reuse, V1 access or fabricated replacement/payment dates |
| Hierarchy/global Super Admin | REJECT | migration000003; organization/admin UI; §3 | FSMO only and three product roles; no tenancy/campus expansion |
| Duplicate return loop | REJECT | C01 ApplyReturn sums duplicates; ProcessReturn takes first match | Reject duplicate IDs before any effect; single map drives counters/evidence/stock/obligations |
| Target eligibility bypass | REJECT | DirectCheckoutExt; H03 | Active existing Borrower required independent of issuing Staff authority; no hidden suspension subsystem |
| Generic file/export auth | REJECT | inventory/files.go Open; H02 | Narrow catalog image purpose; private reports and structured return/audit records need current role/ownership; no generic storage auth or return-photo/attachment feature |
| Unsafe importer/SQL numbering | REJECT | migration/importer.go/cmd/migrate-legacy; H08 | No flattened data, fake NOW, writing dry run, partial claims or copied migrations/grants |
| Borrowing-wide partial notice key | REJECT | emit/BorrowingEventKey; §16 | Each ReturnEvent/ReplacementEvent has distinct key; no missed later partial notice |
| Due/grace/report discrepancies | REJECT | §§14/16/17 | One elapsed24h ceiling PHP 10 basis, Asia/Manila interpretation, no invented grace/amount |
| Foundation/auth/config/UI replacement | REJECT | §§8/12/20/24 | Preserve current tested Phase 1 architecture/62 primitives; no copied starter/seeds/binaries |
| CSV sanitization / async large exports | DEFER | reporting/csv.go/service; §§17/24 | Later approved bounded export formats; no jobs/large-memory platform needed now |

## Duplicate return counterexample retained

Boss audit C01: item issued2; input repeats same item with good1 twice. Domain accumulates2 and completes while stock loop uses first line, restores1 and leaves ghost C1. Row stock CHECK alone passes. This is a source-derived trace, not an executed exploit. V2 rejects duplicates before mutation, has UNIQUE(event,item), composite same-borrowing FKs and one normalized map. New replacement acceptance has the analogous duplicate obligation protection; same destination equipment may legitimately have distinct obligations, whose quantities aggregate once for stock and reconcile with every immutable source line.

## New safeguards driven by owner decisions

D17 clarification adds no file capability: borrowers may physically show a personally stored phone photo during an in-person return, entirely outside eLabTrack. Staff/Admin may inspect actual equipment and alone records condition/quantities. No return-photo upload, transmission, storage, retention, attachment, borrower evidence submission or evidence-step infrastructure is adopted from any reference. Equipment catalog images remain an independent purpose; structured return events/audit remain the software history.

Known loss exits C and total; received damaged original exits C into nonusable D. Neither priced charge nor ghost C keeps a borrowing open: replacement liability does. New accepted units enter A/total as acquisitions; originals/incident history remain. Final physical/replacement zero closes once and freezes the fine; Admin full clear leaves assessed history intact. These are V2 design choices based on current policy, not implied boss fixes.

Prior Phase 2 source inspection covered inventory locking/versioning, ApplyReturn/ProcessReturn/DirectCheckoutExt/idempotency/emit, return/charge migrations, notification claim/attempt behavior, report SQL, file Open and kiosk device concepts. Other source rows cite audit sections rather than fresh runtime proof. Historical boss unit/build/DB results are not V2 acceptance. No source/migration/UI was ported. Future approved reuse needs provenance, actual current foundation compatibility and targeted real PostgreSQL/HTTP tests; no Phase 3A or feature phase is authorized by this reference.
