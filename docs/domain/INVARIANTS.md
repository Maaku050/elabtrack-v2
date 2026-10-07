# Phase 2 invariant specification

Design review draft, 2026-10-08. **57 test-ready invariant specifications in nine categories. None has been implemented or tested as a business workflow in this phase.** These are ENGINEERING RECOMMENDATION safeguards within CURRENT V2 DECISION backend authority/history/transaction boundaries (DEC-013/017/037). Policy-dependent predicates mean the later approved FSMO rule, not a guessed numerical default. See [policy gate](BUSINESS_RULES.md#policy-decision-matrix), [state scenarios](STATE_MACHINES.md#fifteen-required-walkthroughs) and [database enforcement](DATA_MODEL.md).

DB = structural constraint/grant/trigger; APP = authoritative application/domain command; BOTH = both layers with a shared transaction. “PG” means a future disposable **real PostgreSQL** test, not an in-memory substitute. Browser tests prove UX/privacy only; they cannot prove SQL integrity. Enforcement recommendations include restricted runtime grants and immutable-evidence triggers; operator/migration repair is an separately reviewed exceptional procedure.

## Inventory — 12

| ID | Exact rule / why | Enforcement point | Affected operations | Future test |
|---|---|---|---|---|
| INV-01 | Every stock bucket and total is a nonnegative integer; negative/fractional stock cannot represent custody | BOTH: typed domain quantities; NOT NULL/integer/CHECK; overflow-safe arithmetic | Initial/add/remove/reserve/issue/return/disposition/correction | Unit boundaries + PG invalid writes |
| INV-02 | `total = available + reserved + checked_out + damaged + lost + retired`; every accounted unit belongs to exactly one bucket | BOTH: row arithmetic CHECK using wide arithmetic | All inventory movements | Unit vectors + PG constraint |
| INV-03 | For equipment E, `E.reserved = sum(item.reserved_qty for E)` and held lines belong only to permitted preissue states; row sums alone cannot reconcile holds | APP transaction under equipment/borrowing locks; DB row guards/FKs; reconciled query | Submit/approve/deny/cancel/release/import later | PG workflow + last-unit concurrency |
| INV-04 | `E.checked_out = sum(issued-good-damaged-lost for items of E)`; no ghost custody | APP atomic normalized return/issue; DB per-item arithmetic | Issue/partial/final return/import | PG reconciliation after each command |
| INV-05 | Each bucket current count and total equal sum of immutable ledger deltas from explicit baseline; stock sequence equals last movement sequence; no unexplained current-row edits | BOTH: transaction, `(equipment_id,sequence)` unique, append-only/grants; APP sequence/continuous reconciliation | Every quantity mutation | PG direct-write denial and ledger/current comparison |
| INV-06 | Damaged bucket equals damage-return and authorized classification inflows less approved repair/retire/remove outflows; do not equate it forever with historical cumulative damaged loans | APP ledger reconciliation + DB movement structure | Damage/repair/reclassification | PG mixed return and repair |
| INV-07 | Lost/retired buckets equal their signed movement histories after recovery/removal; history still records original disposition | APP reconciliation; DB immutable ledger | Loss/recovery/retire/write-off | PG disposition chains |
| INV-08 | Damage/loss custody disposition alone changes no accounted total; only authorized external add/remove/correction changes total | BOTH: operation-specific vector validation and `delta_total=sum(delta_buckets)` CHECK | Returns/stock add/remove/write-off | Unit vector property cases + PG |
| INV-09 | Under approved reservation policy, no hold/issue consumes more locked available/reserved than exists; competing last-unit claims cannot both consume it | BOTH: ordered FOR UPDATE, fresh predicates/conditional writes; nonnegative CHECK | Submit/approve/direct/release | PG two connections/rollback |
| INV-10 | Proposed archive requires zero reserved and outstanding custody reconciled across rows; archived equipment cannot newly issue | APP guard under same equipment lock; DB lifecycle constraints | Archive versus submit/approval/release | PG concurrent archive/issue |
| INV-11 | Metadata edits cannot set stock or alter old snapshots; stale expected metadata version conflicts | BOTH: narrow repository columns/version predicate; role grants | Rename/category/value/status/image edits | PG/HTTP stale-edit and snapshot regression |
| INV-12 | Ordinary manual adjustment cannot mutate reserved/checked_out, fabricate a return or hide custody discrepancy; all adjustments require approved actor/reason | APP typed movement allowlist, atomic audit/ledger | Stocktake/add/remove/repair/retire | Unit rejection + PG rollback |

## Borrowing — 8

| ID | Exact rule / why | Enforcement point | Affected operations | Future test |
|---|---|---|---|---|
| BOR-01 | Every committed borrowing has ≥1 line; requested quantities>0; `(borrowing,equipment)` unique; no duplicate pool interpretation | BOTH: prevalidated bounded input and commit aggregate validation; DB positive/UNIQUE/FK | Submit/direct/import | Unit duplicate/empty + PG |
| BOR-02 | `reserved_qty` changes only at the confirmed boundary; pending demand under B is not held inventory | APP approved policy; DB quantities/FKs | Submit/approve/release/deny/cancel | Unit A/B traces + PG reconciliation |
| BOR-03 | Only an approved legal state edge may commit; two decisions/releases cannot both transition the same source | APP domain machine + borrowing lock/conditional state write; DB stable-state CHECK | Approve/deny/cancel/release/return | Unit all edges + PG competing staff |
| BOR-04 | Proposed initial issue is all-or-nothing: each `issued_qty=requested_qty`, positive at checkout; reservation consumed exactly once | BOTH: transaction, quantity CHECK and conditional legal edge | Release/direct | PG multiequipment success/rollback |
| BOR-05 | Pending/approved/denied/cancelled have issued=0 and dispositions=0; checked_out has some outstanding; completed has none | APP cross-row aggregate checks; DB local arithmetic | Every edge | Unit transition matrix + PG invalid aggregate |
| BOR-06 | Denial/cancellation/completion retain canonical ID/items/decision evidence; no destructive restoration of already issued custody | BOTH: no delete command, restrict grants/triggers/FKs | Deny/cancel/complete | PG denied/history and double action |
| BOR-07 | Submission, decision, release and completion times are distinct; agreed due/policy/terms/issue evidence is frozen after issue | BOTH: narrow updates/immutability guard; APP approved due calculation | Submit/approve/issue/return/history | Unit calendar + PG immutable issue |
| BOR-08 | Borrower/reference identity is stable, reference unique; resource access never trusts submitted borrower ID for ownership | BOTH: FK/unique reference; APP actor/target access | All reads/mutations | PG duplicate reference + HTTP object access |

## Returns — 8

| ID | Exact rule / why | Enforcement point | Affected operations | Future test |
|---|---|---|---|---|
| RET-01 | Repeated borrowing-item ID in a return is rejected, even when totals fit; UNIQUE(event,item) independently prevents duplicate stored lines | BOTH: input uniqueness before mutation; database unique pair | Every return path/direct HTTP | Unit duplicate + PG/HTTP boss two-good-lines counterexample |
| RET-02 | Each line processed `good+damage+loss` is positive, nonnegative integer components and ≤locked outstanding; event nonempty | BOTH: integer/CHECK; domain/borrowing lock | Partial/final/mixed/concurrent return | Unit boundary + PG competing returns |
| RET-03 | `good_returned=sum(events.good)`, `damaged=sum(events.damaged)`, `lost=sum(events.lost)` per borrowing item; totals never rely on a checkbox | APP transaction/reconciliation; DB immutable event rows | Every return/import | PG multiple events and counters |
| RET-04 | The exact normalized line drives custody, event and ledger vector `checked_out:-processed, available:+good, damaged:+damage, lost:+loss` | APP single validated map; DB movement arithmetic | Mixed good/damage/loss | Unit vector + PG stock/item/event equality |
| RET-05 | Event item belongs to the event's borrowing and referenced borrowing item; unrelated IDs cannot move stock | BOTH: composite ownership FKs + APP authorization | Return | PG cross-borrowing line rejection + HTTP |
| RET-06 | `issued=good+damage+loss+outstanding` for every line; completed iff every outstanding=0 with reconciled movement; loss counts as accounted custody, not good return | BOTH: per-item CHECK/generated outstanding + APP aggregate check | Partial/final | Unit closure + PG all eight return examples |
| RET-07 | Return event/header/lines are immutable after commit; later notes/disposition corrections cannot silently rewrite them | BOTH: INSERT/SELECT only, immutable triggers; no edit/delete command | Return/history | PG denied UPDATE/DELETE/TRUNCATE; HTTP no arbitrary editing |
| RET-08 | Failure anywhere in a multiequipment return leaves stock, cumulative items, event, charges, ledger/audit/outbox/receipt unchanged | APP one transaction; inward tx port; DB atomic commit | Return and final assessment | PG fault injection after first line/required audit/outbox |

## Accountability — 8

| ID | Exact rule / why | Enforcement point | Affected operations | Future test |
|---|---|---|---|---|
| ACC-01 | One immutable base assessment and append-only signed adjustments define balance; records/fines/paid flag cannot define another truth | BOTH: immutable charges/adjustment protection; derived view | Assess/adjust/history/report | PG immutable base + view reconciliation |
| ACC-02 | Assessment/adjustment minor-unit amounts have bounded safe arithmetic and explicit same currency; no fractional floating-point money | BOTH: BIGINT/CHECK/currency refs + APP checked multiplication/sums | All financial calculation | Unit overflow/currency mismatch + PG |
| ACC-03 | `outstanding=base+increases-reductions-waivers-settlements+approved_reversals≥0`; concurrent adjustments cannot over-discharge | APP locks charge row then derives balance; DB row constraints/immutable rows | Reduce/waive/settle/increase/reverse | PG concurrent reductions/settlements |
| ACC-04 | Each assessment has unique source/window/component, immutable basis and applicable policy/value snapshot; estimates never double-post cumulative damage/loss | BOTH: source-key UNIQUE; APP policy/source/ownership checks | Assess/final return/scheduler later | Unit dates/rates + PG duplicate posting/partial events |
| ACC-05 | Custody completion independent of financial balance; zero balance independent of good stock return | APP separate aggregates; no payment condition on borrowing close | Final return/settlement | PG completed-with-open-charge |
| ACC-06 | Assessments, waivers, administrative settlements and actual payments/collections are distinct report measures; settlement alone is not revenue | APP read semantics and approved payment evidence | Reports/settle/payment branch | PG report reconciliation + HTTP projections |
| ACC-07 | Every adjustment has approved actor, reason, amount, kind and timestamp; it never erases the base or prior evidence | BOTH: required fields/FKs/immutable triggers + APP authority | Adjust/waive/settle | PG constraints/rollback + HTTP permissions |
| ACC-08 | A proposed positive reversal references a same-charge reduction/waiver/settlement and cumulative reversals≤original amount; repeated reversal cannot create credit | BOTH: self/composite FK, positive amount; APP charge lock/approved reversal policy | Correction/reversal if approved | Unit invalid kind + PG concurrent reversal |

## Historical integrity — 5

| ID | Exact rule / why | Enforcement point | Affected operations | Future test |
|---|---|---|---|---|
| HIS-01 | Canonical borrowing, returns, charge/adjustment, ledger and audit persist with noncascade evidence refs; no active→copied history deletion | BOTH: restrictive FKs/grants; APP commands | Deny/complete/archive/account lifecycle | PG retained join chain |
| HIS-02 | Metadata/profile changes cannot rewrite issued name/category/value/currency/borrower snapshots or effective policy/terms references | BOTH: column guards; APP write boundaries | Metadata/issue/history | PG historical price/name regression |
| HIS-03 | Referenced policy/terms versions/content/hash and linked acceptance are immutable; historical reports use them, not latest settings | BOTH: immutable rows/hash validation/FKs | Policy publication/issue/assessment | PG new version with old loan |
| HIS-04 | High-value mutation records durable who/what/when/entity/meaningful before-after/context atomically; operational log is not substitute | APP tx audit port; DB audit required fields/immutability | Decisions/issue/return/stock/archive/account role/status/charges | PG audit failure rollback |
| HIS-05 | Legacy missing dates/events/payment facts remain unknown with provenance; privacy processing cannot fabricate or destroy custody/financial evidence | APP future mapping/reconciliation/privacy procedure | Migration/retention later | Sanitized mapping tests + PG reconciliation |

## Authorization — 4

| ID | Exact rule / why | Enforcement point | Affected operations | Future test |
|---|---|---|---|---|
| AUTH-01 | Current server authority for approved operation, resource and actor/target owns every command/read; UI/JWT/body role cannot grant permission | APP current-account/capability ports | All commands/projections | HTTP approved operation matrix |
| AUTH-02 | Command authorization remains valid through its transaction linearization point; status/permission writers honor same account locks | BOTH: ordered user/profile locks and narrow writers; APP recheck | Approval/issue/return/adjust/account edit | PG demotion/inactivation race |
| AUTH-03 | Borrower gets own private records; narrow catalog/image does not grant report/evidence/download permission; replay reauthorizes | APP scoped read/replay/download ports | Lists/detail/export/image/replay | HTTP cross-user/object/file-purpose tests |
| AUTH-04 | Account/role/eligibility changes preserve custody/history and record actor/reason; role matrix remains institutionally approved | APP transactional future identity commands; DB FKs/audit | User management/eligibility | PG rollback/old loans + HTTP escalation denial |

## Eligibility — 4

| ID | Exact rule / why | Enforcement point | Affected operations | Future test |
|---|---|---|---|---|
| ELG-01 | Auth account existence/active role is not borrowing eligibility; category/manual institutional/email gates are distinct approved predicates | APP BorrowerProfile eligibility port | Submit/approve/release/direct | Unit eligible/ineligible categories; HTTP |
| ELG-02 | Target borrower exists and passes current approved eligibility before issue even when actor is staff | BOTH: user/profile locks/FKs; APP recheck | Direct/approve/release | PG inactive/suspended target counterexample |
| ELG-03 | Reservation/old approval cannot bypass later account/catalog/terms/due restrictions at physical release | APP transaction-time validation with locked rows/version refs | Release/combined issue | PG inactive-between-request-and-issue |
| ELG-04 | Account/catalog inactivation preserves existing custody and permits authorized staff disposition under approved return policy; never grants new issue | APP separate issue versus return predicates | Inactivate/return/history | PG inactive target return + new-issue denial |

## Command idempotency — 4

| ID | Exact rule / why | Enforcement point | Affected operations | Future test |
|---|---|---|---|---|
| CMD-01 | Same actor/operation/concrete resource/channel/context/key + same payload returns same committed result identity, once; key reuse with changed hash conflicts | BOTH: unique command receipt claim in business tx + APP normalization | Submit/direct/returns/selected irreversible commands | PG same-key concurrent/retry + HTTP |
| CMD-02 | Retry/replay authenticates and reauthorizes before exposing stored private outcome; key is not bearer authority | APP current-account/resource checks before replay | Every receipt replay | HTTP revoked/different actor |
| CMD-03 | Fingerprint covers concrete resource, target borrower, normalized quantities, selected terms/due inputs and channel; no missing source-context alias | APP canonical request encoder; DB explicit context columns | Web/staff/conditional kiosk commands | Unit changed resource/channel + HTTP conflict |
| CMD-04 | Lost commit acknowledgment is treated as unknown; retry original key inspects receipt, not blind second operation; failed transaction leaves no successful receipt | APP tx/replay handling; DB atomic receipt | Every idempotent command | PG connection interruption/replay/fault injection |

## Notifications — 4

| ID | Exact rule / why | Enforcement point | Affected operations | Future test |
|---|---|---|---|---|
| NTF-01 | Event key includes distinct business event identity; repeated partial returns have different ReturnEvent IDs/keys; scheduler keys include approved period/policy | BOTH: unique event key + APP event construction | Return/decisions/reminders | PG two partial notices + duplicate scheduler |
| NTF-02 | Required outbox insertion commits with business change; later provider failure never rolls back business | APP tx outbox port; DB atomic commit | Business command/delivery | PG enqueue failure rollback + provider failure after commit |
| NTF-03 | Attempts preserve outcome/time/provider acceptance separately from delivery proof; finite retry/dead-letter handling, no swallowed persistence failure | APP worker port, DB append-only attempts | Future delivery/retry | PG/provider crash/error/duplicate-acceptance cases |
| NTF-04 | Lease completion fenced by unique claim/attempt token; expired owner cannot overwrite newer attempt; bounded claims use database locking | BOTH: claim update predicate/unique attempt; APP finite lease | Worker claim/lease/retry | PG two workers + stale owner completion |

## Enforcement and review limits

Simple structural arithmetic belongs to database CHECK/NOT NULL/FK/UNIQUE. Cross-row sums, eligibility and state transitions require the application transaction protocol and future reconciliation tests; PostgreSQL row CHECK cannot query other rows to enforce them ([PostgreSQL constraints](https://www.postgresql.org/docs/current/ddl-constraints.html)). The design does not claim protection from arbitrary schema-owner changes. Future migrations must restrict runtime update/delete rights on evidence and verify trigger/grant behavior with the actual Phase 1H role split.

Counts: inventory12, borrowing8, returns8, accountability8, history5, authorization4, eligibility4, idempotency4, notifications4 = **57**. Acceptance requires real PostgreSQL concurrency/fault tests plus authoritative HTTP tests, with pure domain/calendar/money unit tests and later browser accessibility/session checks. Phase 1 tests verify foundation only; passing them is not proof of these future business invariants.
