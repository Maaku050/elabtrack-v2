# Phase 2.5 invariant catalog

2026-10-08. **71 active invariants plus one explicitly retired ID (72 catalog entries).** Prior IDs remain; wording is rebased to current owner decisions. ACC-08's generalized financial reversal rule is retired, not an authorized feature. Added INV-13/14, BOR-09/10, ACC-09, REP-01–08, AUTH-05 and ELG-05. These are future acceptance criteria, not claims of implemented enforcement or executed tests.

APP means domain/application transaction or read boundary; DB means structural constraint/immutable grant or guard; BOTH means both are needed. PG tests mean **future real PostgreSQL** tests with actual migrator/runtime roles. Cross-row sums and authorization are not ordinary row CHECKs. See [rules](BUSINESS_RULES.md), [schema](DATA_MODEL.md), [18 design traces](STATE_MACHINES.md#eighteen-required-walkthroughs).

## Inventory — 14

| ID | Exact invariant | Enforcement | Operations | Future acceptance test |
|---|---|---|---|---|
| INV-01 | All A/R/C/D/total stock are nonnegative bounded integers. | BOTH typed quantities and DB NOT NULL/CHECK | All stock commands | Unit overflow and PG invalid writes |
| INV-02 | total_tracked=A+R+C+damaged_held; lost/retired history is not current physical stock. | BOTH vector rules and row arithmetic CHECK | Every movement | PG row constraint and vector properties |
| INV-03 | Equipment R=sum(item.reserved); only unprocessed PENDING holds own reservations. | APP locks/reconciliation plus DB line bounds | Submit/deny/cancel/expire/issue | PG cross-row and last-unit race |
| INV-04 | Equipment C=sum(issued-good-damaged-lost); replacement outstanding never inflates C. | APP exact normalized maps and DB counters | Issue/return/replacement | PG mixed/loss-only no-ghost reconciliation |
| INV-05 | Each current count equals explicit baseline plus immutable ledger deltas; sequence contiguous. | BOTH transaction, unique equipment/sequence, evidence guards | All quantity changes | PG ledger/current equality and write denial |
| INV-06 | D=received damaged originals less authorized disposition outflows; D need not equal unresolved damage obligations. | APP source-vector reconciliation | Damage/replacement/later disposition | PG retained-original versus replacement count |
| INV-07 | Historical damage/loss stays immutable after replacement/physical disposal; no lost stock bucket. | BOTH immutable incidents and ledger | Incident/acceptance/disposition | PG loss/replacement/history join |
| INV-08 | Damage keeps total; loss reduces total; new replacement increases total; delta_total=sum(bucket deltas). | BOTH typed vectors/DB CHECK | Mixed return/acquisition | Unit vectors and PG totals |
| INV-09 | Reserve/issue cannot exceed locked available/held stock; two last-unit claims cannot both succeed. | BOTH ordered locks and nonnegative predicates | Submit/approve/direct | PG two connections and rollback |
| INV-10 | Archive requires zero holds, physical custody and replacement liability; inactive/archived pools cannot newly issue. | APP equipment/borrowing lock protocol | Archive versus issue/acceptance | PG archive race |
| INV-11 | Metadata changes never set stock or rewrite snapshots; stale metadata version conflicts. | BOTH narrow writes/version predicate | Metadata edits | HTTP stale version and PG snapshots |
| INV-12 | Ordinary adjustment cannot alter R/C, fabricate return or hide custody; exceptional Admin correction must reconcile evidence. | APP role/typed operation guard and audit | Stocktake/correction | Unit rejection and PG audit rollback |
| INV-13 | Accepted replacement adds q available and q total exactly once, never original good_returned or implicit damaged removal. | APP same validated map/immutable acquisition source | Accept replacement | PG acquisition/obligation/ledger comparison |
| INV-14 | Expired pending hold is released once by a committed expiry edge, not merely subtracted in a read projection. | APP locked expiry and outbox/audit transaction | Expiry/lazy expiry | PG approval-expiry race and enqueue failure |

## Borrowing — 10

| ID | Exact invariant | Enforcement | Operations | Future acceptance test |
|---|---|---|---|---|
| BOR-01 | Every borrowing has a nonempty unique equipment map; requested quantities positive. | BOTH input uniqueness and DB UNIQUE/FK/CHECK | Submit/direct | Unit empty/duplicate and PG rejection |
| BOR-02 | Submission reserves each requested quantity immediately; all terminal pending edges release; issue consumes exactly the same hold. | APP transaction and DB local quantity bounds | Submit/approve/deny/cancel/expire | PG cross-row stock/hold equality |
| BOR-03 | Only legal lifecycle edges commit; approval and physical release are one atomic edge with no APPROVED state. | APP borrowing lock/domain machine and DB enum | Decisions/issue/return/acceptance | Unit edge matrix and PG competing decisions |
| BOR-04 | Initial issue is all-or-nothing issued=requested, positive; no duplicate consumption. | BOTH transaction/row bounds | Approval/direct | PG multiequipment rollback |
| BOR-05 | Preissue terminal and PENDING have issued/disposition0; CHECKED_OUT has physical or replacement outstanding; COMPLETED has neither. | APP aggregate state and DB local counts | All edges | PG invalid closure and replacement-only open |
| BOR-06 | Denied/cancelled/expired/completed keep canonical ID/lines; issued stock is never restored by deleting borrowing. | BOTH no delete command/restrict evidence FKs | Terminal edges | PG retained history and repeated decisions |
| BOR-07 | Submission/expiry/decision/issue/completion evidence is explicit; issue due/policy/terms/snapshots freeze. | BOTH narrow updates and immutable guards | Submit/issue/complete | PG timestamps and old-version history |
| BOR-08 | Stable borrower/reference, unique reference; caller IDs never establish ownership. | BOTH FK/reference UNIQUE and scoped authorization | All reads/commands | HTTP object access and PG duplicate ref |
| BOR-09 | Pending deadline=submitted_at+24h captured policy; at now>=expiry pending decision commits EXPIRED and releases, not rollback. | APP locked server-time decision and DB timestamps | Expiry/approve/cancel/deny | Unit exact equality and PG late command commit |
| BOR-10 | Every issue has future date+time due_at in absolute time; partial/incident/acceptance does not extend it; no seven-day limit. | APP due conversion/immutable issue and DB presence | Issue/return/accept | Unit Manila boundaries and PG due immutability |

## Returns — 8

D17 correction keeps all existing IDs/counts: a borrower may physically show a personally stored phone photograph in person, entirely outside eLabTrack. Staff/Admin may inspect actual equipment and alone records condition/quantities. Return records contain no uploaded/transmitted/stored/retained photos or return attachments; there is no borrower evidence interface, storage model or evidence-upload step. RET-07 makes that boundary explicit without adding another product feature. Catalog images remain separate.

| ID | Exact invariant | Enforcement | Operations | Future acceptance test |
|---|---|---|---|---|
| RET-01 | Duplicate borrowing-item IDs reject the whole input even if summed values fit; UNIQUE(event,item) independently protects storage. | BOTH input validation/DB UNIQUE | Return | Unit/HTTP boss duplicate counterexample and PG |
| RET-02 | Each nonnegative integer line processes positive good+damage+loss<=locked physical outstanding; event nonempty. | BOTH domain lock and DB CHECK | Partial/mixed/concurrent return | Unit boundaries and PG race |
| RET-03 | Item good/damage/loss totals equal immutable return-line sums. | APP counters/reconciliation and DB immutable rows | Return/history | PG multiple events and cumulative totals |
| RET-04 | One normalized line drives C−(g+d+l), A+g, D+d, total−l and required replacements d+l. | APP exact map and DB vector guards | Mixed return | PG all counters/evidence/stock equality |
| RET-05 | Return line belongs to same borrowing as its event/item; unrelated IDs cannot move stock. | BOTH composite FKs and authorization | Return | PG cross-parent rejection and HTTP scope |
| RET-06 | issued=good+damage+loss+physical; completion also requires replacement0; loss is not a good return. | BOTH row CHECK and APP aggregate closure | Partial/final/acceptance | Unit closure and PG zero-C/open-liability |
| RET-07 | Return event/header/lines are immutable Staff/Admin-recorded condition/quantity records; later resolution never rewrites incidents. Borrowers cannot mark returned or submit return evidence; no return-photo payload/attachment is accepted or retained. | BOTH INSERT/SELECT evidence grants/guards | Return/history | PG UPDATE/DELETE/TRUNCATE denial; future HTTP/schema contract rejects borrower return/evidence and return-photo/media fields |
| RET-08 | Return failure rolls back counters, stock, event, obligations, fine freeze, ledger/audit/outbox/receipt together. | APP one transaction | Multiitem/final return | PG failure after first line and required evidence |

## Overdue accountability — 8

| ID | Exact invariant | Enforcement | Operations | Future acceptance test |
|---|---|---|---|---|
| ACC-01 | One fine basis per issued borrowing; live formula/final snapshot minus full clearances is the only monetary truth. | BOTH unique borrowing fine and APP projection | Issue/complete/clear/report | PG basis/clearance reconciliation |
| ACC-02 | Money uses bounded integer centavos, PHP basis and safe multiplication; no float money or price liability. | BOTH BIGINT/CHECK and checked domain arithmetic | Fine/clear | Unit overflow/currency and PG |
| ACC-03 | Outstanding=assessed or frozen final minus all clearances>=0; each clear equals entire locked positive outstanding. | APP borrowing/fine locks and DB positive rows | Admin full clear | PG simultaneous clears/stale expectation |
| ACC-04 | PHP 10*ceil(max(0,effective_end-due_at)/24h); end=now while unresolved, completed_at after completion; no daily cumulative posting. | APP one shared typed formula/frozen policy | Projection/complete/reminder/report | Unit exact due/24h boundaries and replacement-only overdue |
| ACC-05 | Operational completion depends on physical/replacement only, independent of fine; clearing never resolves inventory. | APP distinct quantity/fine projections | Complete/clear | PG completed with open fine and paid with replacement due |
| ACC-06 | Assessed, PAID, WAIVED and OTHER_RESOLUTION are separate report measures; total cleared is not cash. | APP authoritative read semantics | Reports/clear history | PG method totals and bounded HTTP read |
| ACC-07 | Clearance preserves cumulative assessment, full cleared amount, named Admin, time, method and optional note; no erasure. | BOTH immutable rows/required fields and APP authority | Clear/history | PG evidence/write denial and HTTP |
| ACC-09 | Completion freezes final days/amount/time once in the same transaction; later clears cannot change final basis. | BOTH NULL-to-final guard and APP atomic close | Final return/replacement/clear | PG freeze rollback and post-clear immutable amount |

## Replacement obligations — 8

| ID | Exact invariant | Enforcement | Operations | Future acceptance test |
|---|---|---|---|---|
| REP-01 | Each damage/loss source creates exactly that positive required quantity, unique per return-line/kind. | BOTH UNIQUE/composite FKs and APP source equality | Incident | PG mixed source-kind reconciliation |
| REP-02 | Remaining=required-accepted>=0; accepted cannot exceed locked remaining. | APP borrowing/obligation locks and DB positive acceptance | Accept | Unit bounds and PG competing acceptance |
| REP-03 | Accepted equals sum immutable acceptance lines; no independent writable accepted flag or balance. | BOTH evidence guards and APP derived sums | Acceptance/history | PG sum reconciliation |
| REP-04 | Any unresolved obligation prevents COMPLETED even when physical C0. | APP aggregate completion after every resolution | Return/accept | PG zero-physical/open-liability |
| REP-05 | Resolution never erases original damage/loss kind/quantity/source/time. | BOTH immutable obligation basis and incident evidence | Accept/disposition/history | PG post-acceptance incident comparison |
| REP-06 | Appropriate type/equivalence is recorded by named Staff/Admin; accepted q is new usable acquisition, not original repair. | APP permission/equipment/type evidence and exact vector | Accept | HTTP actor/type and PG acquisition count |
| REP-07 | Duplicate obligation IDs reject input; UNIQUE(event,obligation) and same-borrowing composite FKs prevent duplicates/cross-parent stock. | BOTH validator/DB constraints | Accept | Unit duplicates and PG cross-parent rows |
| REP-08 | Final concurrent acceptance completes/freezes/emits once; failed acquisition/evidence rolls back acceptance and stock. | APP common borrowing-first lock and one transaction | Last acceptance versus return/acceptance | PG two connections/fault injection/same-key retry |

## Historical integrity — 5

| ID | Exact invariant | Enforcement | Operations | Future acceptance test |
|---|---|---|---|---|
| HIS-01 | Canonical borrowing, return/replacement, fine/clearance, ledger and audit retain restrictive evidence joins; no copied terminal truth. | BOTH RESTRICT/grants and APP commands | Terminal/history/archive | PG retained chain |
| HIS-02 | Metadata/account changes cannot rewrite issued item/borrower snapshots or accepted version evidence. | BOTH column guards/narrow repository | Edits/issue/history | PG old names/metadata unaffected |
| HIS-03 | Published policy/terms content/hash and linked acceptance immutable; history uses bound basis. | BOTH immutable rows/FKs and APP validation | Publish/issue/fine | PG new terms and old borrowing |
| HIS-04 | Consequential mutation has durable actor/action/time/entity/before-after/source atomically; logs are not substitute. | APP tx audit and DB immutable fields | Stock/account/decision/return/replacement/clear | PG audit failure rollback |
| HIS-05 | Missing legacy dates/events/replacement/payment facts stay unknown with provenance; privacy cannot fabricate evidence. | APP future mapping/retention procedure | Later migration/privacy | Sanitized mapping tests and reconciliation |

## Authorization — 5

| ID | Exact invariant | Enforcement | Operations | Future acceptance test |
|---|---|---|---|---|
| AUTH-01 | Current server account/capability/resource check owns reads/commands; UI/JWT/body role cannot grant authority. | APP current-account ports | All commands/reads | HTTP BORROWER/STAFF/ADMIN matrix |
| AUTH-02 | Actor/target authority valid at transaction linearization; role/status writers share sorted account-lock order. | BOTH account locks/narrow writers and APP recheck | Issue/return/accept/clear/account edit | PG deactivation/demotion race |
| AUTH-03 | Borrower private records own-only; catalog does not grant audit/export access; replay reauthorizes. | APP scoped projection/replay ports | List/detail/replay/image/report | HTTP cross-user and narrow file purpose |
| AUTH-04 | Role/status changes preserve history; Staff provisioning cannot assign Staff/Admin/deactivate or gain fine-clear access. | APP role matrix/transactional audit and DB FK | Provision/deactivate/bulk/roles | HTTP escalation and PG history |
| AUTH-05 | Only current named Admin can clear any fine; Staff capability explicitly excludes it. | APP current-account Admin check at locked clear | Clear PAID/WAIVED/OTHER | HTTP Staff403 and PG demotion race |

## Account and terms gates — 5

| ID | Exact invariant | Enforcement | Operations | Future acceptance test |
|---|---|---|---|---|
| ELG-01 | Active registered borrower is sufficient; categories Student/Faculty are not roles; no separate eligibility or automatic fine block. | APP current User predicate, no profile table | Submit/issue | HTTP categories and existing-fine allowed |
| ELG-02 | Target active existing borrower required at issue even if Staff/Admin actor active. | BOTH user locks/FK and APP target check | Direct/approve | PG inactive target counterexample |
| ELG-03 | Old hold cannot bypass current account/catalog/due issue guards; bound submission terms do not require per-loan reacceptance. | APP locked checks and acceptance FK | Approval | PG deactivation and version-change cases |
| ELG-04 | Deactivation/inactive catalog preserves obligations and Staff/Admin can resolve them; no new borrower action/issue. | APP separate issue versus resolution predicates | Deactivate/return/accept | HTTP inactive caller and PG staff return |
| ELG-05 | Before new request current material terms accepted by same borrower; acceptance unique user/version/time, no staff impersonation. | BOTH UNIQUE/composite FK and APP publication lock | Terms accept/submit/direct | PG publish-submit race and HTTP missing acceptance |

## Command idempotency — 4

| ID | Exact invariant | Enforcement | Operations | Future acceptance test |
|---|---|---|---|---|
| CMD-01 | Same actor/operation/resource/key and payload returns one committed identity; changed payload conflicts. | BOTH UNIQUE receipt and APP canonical hash | Submit/issue/return/accept/clear | PG concurrent/retry and HTTP conflict |
| CMD-02 | Replay authenticates/currently authorizes before returning private evidence; key not bearer authority. | APP replay authorization | All replays | HTTP revoked/different actor |
| CMD-03 | Hash covers target/concrete parent/unique quantities/due/terms/method/equivalence as relevant; kiosk has no device/channel identity. | APP canonical encoder and DB explicit scope | All idempotent commands | Unit changed input and HTTP conflict |
| CMD-04 | Lost commit acknowledgment is unknown until same-key receipt inspection; failed tx has no successful receipt. | APP transaction/replay and DB commit | All commands | PG connection interruption/faults |

## Notifications — 4

| ID | Exact invariant | Enforcement | Operations | Future acceptance test |
|---|---|---|---|---|
| NTF-01 | Event key identifies each distinct source; partial return and replacement events cannot collapse into one borrowing-wide key. | BOTH UNIQUE key and APP construction | Submit/deny/issue/expire/due/overdue/partial/replacement/complete | PG two partial/two replacement notices and reminder dedupe |
| NTF-02 | Required outbox intent commits with business effect; provider failure later cannot roll it back. | APP transaction/outbox port | Business/delivery | PG enqueue rollback and provider failure |
| NTF-03 | Attempts preserve sanitized outcome/time/provider acceptance distinct from delivery; finite retries/failure visibility. | APP future worker and DB attempts | Future delivery/retry | Provider/error/crash tests |
| NTF-04 | Unique lease token fences completion; expired worker cannot overwrite new owner; bounded claims. | BOTH claim predicate/attempt token and APP lease | Future worker | PG two workers/stale owner |

## Retired invariant — 1

| ID | Disposition | Reason |
|---|---|---|
| ACC-08 | RETIRED; historical positive adjustment-reversal invariant | Generic charge adjustments/reversals are outside confirmed full-clear scope. Do not implement this old rule or reuse its ID for a different feature. |

## Count and verification limits

Inventory14 + borrowing10 + returns8 + overdue accountability8 + replacements8 + history5 + authorization5 + account/terms5 + idempotency4 + notifications4 = **71 active**. Plus ACC-08 retired =72 entries; all IDs unique. Original 57 became56 retained active +15 added active. All 18 requested scenarios map to these rules, especially expiry INV-14/BOR-09, replacement REP-01–08, fine freeze ACC-09 and Admin-only AUTH-05.

Future enforcement must prove legal edges, exact reconciliation, boundary arithmetic, permissions, duplicate inputs, same/different-key retries and rollback/concurrency. Passing Phase 1 foundation tests is not evidence that these business invariants exist. No runtime suites run for this documentation-only change; Phase 1 implementations and security contracts remain unchanged.
