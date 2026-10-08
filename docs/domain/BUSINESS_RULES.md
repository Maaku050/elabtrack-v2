# Phase 2.5 business rules and authorization draft

**Current Batch1 implementation,2026-10-09:** Phases5–6 now implement the scoped account/catalog/inventory boundaries described in the overlays below and current API_CONTRACTS.md. Earlier unimplemented/candidate statements are historical Phase2/4 design evidence, not current status. See PHASE5_REPORT.md (external-gate partial) and PHASE6_REPORT.md (reviewed inventory scope complete). Phase7 remains unauthorized.


2026-10-08. Owner decisions D1–D29 and the final Student/Faculty update DEC-070 and additional owner decisions DEC-071–073 are **CONFIRMED V2 DECISIONS / authoritative working product policy**. Engineering choices are called out separately. They override conflicting capstone/V1/boss assumptions; none is implementation evidence. See [source policy](../project/SOURCE_OF_TRUTH.md), [accepted decisions](../project/DECISIONS.md), [remaining questions](../project/OPEN_DECISIONS.md).

## Decision integration matrix

D numbers refer to the owner's Phase 2.5 instruction, with permanent register IDs below.

| Owner decision | Current working rule | Register / principal design effect |
|---|---|---|
| D1 | Active registered BORROWER, including at least Student/Faculty categories; category is not role | DEC-051; user attributes, no StudentRole/FacultyRole |
| D2 | Active account is borrowing gate; Admin deactivation blocks normal auth actions; no eligibility subsystem or fine block | DEC-051; reuse identity status, preserve existing obligations |
| D3 | BORROWER/STAFF/ADMIN; named multiple Admins; Admin inherits staff operations | DEC-052; future server permission matrix |
| D4, updated by DEC-070 | No public signup; only Admin creates any account; individual Student/Faculty; Student-only bulk workflows | DEC-053 Staff grant/generic import partially superseded; secure activation remains future implementation |
| D5 | Approval is physical handover, no approved-waiting state | DEC-054; one issue edge/transaction |
| D6 | Reserve on submission, atomic under stock locks | DEC-054; available→reserved |
| D7 | Pending expires 24h after submission, retained EXPIRED history | DEC-055; TTL policy snapshot and expiry index |
| D8 | Borrower cancel own PENDING only | DEC-055; release hold, retained history |
| D9 | Staff/Admin denial of PENDING; required borrower-visible reason | DEC-055; release and reason commit together |
| D10 | Staff/Admin direct checkout, active existing borrower, availability/due/server checks | DEC-054; immediate CHECKED_OUT |
| D11 | Required issued due date+time; no fixed seven-day maximum | DEC-056; due_at and nullable future maximum |
| D12 | Asia/Manila local interpretation, absolute timestamptz persistence | DEC-056; shared calendar boundary |
| D13 | First-use/current-version terms acceptance; new material version gates new request | DEC-057; user/version/time evidence, no per-loan acceptance |
| D14 | Overdue means now>due_at and unresolved physical OR replacement obligation | DEC-058; not checked_out count alone |
| D15 | PHP 10/day; adopt recommended 24h ceiling arithmetic as engineering detail | DEC-059; dynamic/frozen final amount |
| D16 | Existing fines permit submission/reservation; human review may deny at handover | DEC-051/054; no auto-block |
| D17 | Staff/Admin alone records return condition/quantities; personally stored photos may only be physically shown in person, entirely outside eLabTrack | DEC-058 clarified; no borrower return/evidence submission or return-photo feature |
| D18 | Partial returns retain original deadline and leave loan open | DEC-058; fine continues |
| D19 | Damage/loss requires replacement, not automatic peso price charge | DEC-058; immutable incident plus obligation |
| D20 | Disposed damaged/lost units leave checked_out; liability tracked separately | DEC-058/062; no ghost custody |
| D21 | Staff/Admin accepts appropriate type/equivalent, restocks and reduces liability | DEC-058/062; acquisition plus acceptance event |
| D22 | Original damaged disposition later, NON-BLOCKING for borrowing UX | DEC-058/062; nonusable held default; no mandatory repair flow |
| D23 | Item reconciliation and completion include both physical and replacement outstanding | DEC-058; exact normalized maps/duplicate protection |
| D24 | Primary money liability overdue fine only | DEC-060; focused fine schema |
| D25 | Full Admin clear with assessed history, actor/time/method/note; no online payment | DEC-060; fine_clearances |
| D26 | Fine clearance ADMIN only | DEC-052/060; Staff denied by backend |
| D27 | OUTSTANDING/CLEARED, no partial payment/allocation | DEC-060; full current-balance clear |
| D28 | Kiosk is catalog/cart/request, no hardware/device/handoff infrastructure | DEC-061; remove device candidates |
| D29 | One borrower mobile-first responsive flow; staff desktop/tablet responsive | DEC-061; supersede dedicated shell assumption |

## Account, terms and issue rules

Active registered Borrower with current terms accepted is the borrowing baseline; no course, suspension, fine-balance or active-loan-count restriction is added. DEC-070 defines Student/Faculty as BORROWER categories. Only Admin creates accounts; Staff assists operationally. Student requires approved SKSU institutional email and a unique official Student ID stored as a string; Faculty is individual-only with any valid unique accessible email and no required Student ID. Both choose separate eLabTrack passwords in secure activation, never mailbox passwords or SSO. No public registration/shared admin credential. Deactivation retains obligations/history for Staff/Admin resolution; DEC-073 allows Admin Student individual/bulk deactivation despite fines, active/overdue borrowing, unreturned equipment and replacements, with appropriate warnings and confirmation; outstanding obligations never veto deactivation. No automatic clearance/waiver/payment, return, replacement resolution, unresolved-loan closure, overdue calculation change or history deletion. OPEN-029 is resolved.

Terms acceptance is unique per user/version and immutable. Current material version must be accepted before a new submission. Recommended direct checkout requires the borrower to have accepted current terms; Staff cannot synthesize that acceptance. Submission binds its accepted version; publishing new terms does not force reacceptance for an already pending request. Official terms will be finalized after presenting the application to FSMO (DEC-072). Pending wording does not block independent account-management/inventory development under separate authorization; live borrowing still requires official publication/documented acceptance and secure onboarding, never placeholder or automatic consent.

Request input is a nonempty, bounded unique equipment/positive integer quantity map. Recommended requested due date/time conveys intent; staff confirms/edits and audits actual future `due_at` at issue. Initial fulfillment is all-or-nothing; changing quantities requires a new request rather than silent partial issue. No hard duration maximum. Issue rechecks active actor/target, role, pool usability, held quantities, captured policy and due time under locks. Staff sees fine/replacement accountability, may approve or deny based on face-to-face review; no automatic prohibition or forced payment UI.

Borrower can cancel only own PENDING request, and only before its expiry linearization time. Staff/Admin may retain reasoned administrative pending cancellation as an **engineering recommendation**, not a required new permission: future feature plan must explicitly decide whether it is needed. No issued-loan cancellation/deletion/restoration route.

## Future permission matrix

All cells are future server-authoritative capabilities, unchanged Phase 1 auth code. Active identity and resource ownership are checked for every command/read/replay. BORROWER row means provisioned borrower capability; staff roles do not automatically become borrower accounts or obtain other users' borrower impersonation rights.

| Operation | BORROWER | STAFF | ADMIN | Server boundary |
|---|---|---|---|---|
| Browse/search/filter equipment | Yes | Yes | Yes | Bounded catalog; no private account data |
| Build own cart, submit own request | Yes | No automatic borrower grant | No automatic borrower grant | Current terms, active borrower, stock lock; old fine allowed |
| Cancel own pending request | Yes | No impersonation | No impersonation | Owner + PENDING + before expiry |
| View own history/status/accountability | Yes | Own borrower flow only if separately authorized | Same | Owner-scoped reads |
| View borrowers' operational history/accountability | No | Yes | Yes | Operational purpose and bounded query |
| Review/approve and physically issue | No | Yes | Yes | One atomic PENDING→CHECKED_OUT |
| Deny pending request | No | Yes | Yes | Required visible reason; release hold |
| Direct checkout | No | Yes | Yes | Active existing borrower; no bypass |
| Process good partial/final return | No | Yes | Yes | Locked physical outstanding |
| Record damage/loss | No | Yes | Yes | Immutable incident + exact replacement obligation |
| Accept replacement | No | Yes | Yes | Locked remaining; acquisition + history |
| Equipment metadata CRUD / ordinary stock operations | No | Yes, normal operations | Yes | Create/read/update/archive; guard active holds/custody/liability; ledger every movement |
| Exceptional stock correction | No | No | Yes | Explicit reason/audit; cannot hide active custody or fabricate return |
| Provision individual Student/Faculty Borrower | No | No | Yes | Current Admin backend check; BORROWER fixed; conditional Student ID/SKSU email versus Faculty accessible email/no required ID; separate privileged workflow |
| Deactivate/reactivate Borrower | No | No | Yes | Student deactivation allowed despite all obligations with warnings/confirmation (DEC-073); no automatic resolution or changed overdue/history; reactivation remains separate engineering lifecycle counterpart |
| Excel bulk Student creation | No | No | Yes | BORROWER/STUDENT fixed; required unique textual Student ID/SKSU email; complete validation preview before explicit confirmation, results/retries/audit; no password or role/category assignment columns |
| Excel Bulk Deactivate Students | No | No | Yes | Stable Student ID/email matching; current BORROWER AND STUDENT; exclude Faculty/Staff/Admin; matched/unmatched/conflicts/already-inactive/obligation warnings; explicit selected-account confirmation, idempotency/history/audit; DEC-073 allows deactivation despite all obligations without an obligation-based veto |
| Create/manage STAFF/ADMIN | No | No | Yes | Named actors, anti-escalation/last-admin recovery design later |
| Clear fine (PAID/WAIVED/OTHER_RESOLUTION) | No | No | Yes | Entire current balance only; actor/time/basis retained |
| Manage policy/terms/settings | No | No | Yes | Immutable published versions; audit |
| View administrative audit/reports | No | No | Yes | Bounded purpose-specific access; operational summaries for Staff do not imply admin exports |
| Administrative pending cancellation | No | Not yet selected | Not yet selected | Optional recommendation, needs future scope decision |
| Automatic pending expiry | No client command | No client command | No client command | System use case, deadline and stock locks, explicit system audit actor |

Normal inventory means metadata maintenance, supported usable stock acquisition/removal and archive under guards; detailed damage repair/disposal/correction procedures remain later inventory policy. Hard deletion of consequential history is not inventory CRUD. React guards are UX only. Permission changes and deactivation writers must participate in the transaction-time account-lock protocol; existing session behavior is not modified here.

## Confirmed account activation and roster operations

DEC-070 is the approved current rule: Admin-only creation, Student/Faculty as BORROWER attributes, Student-required unique textual official Student ID and approved SKSU email, Faculty individual-only with valid unique accessible email and no required Student ID. Preserve UUID/history and normalized unique email. String IDs retain leading zeroes; syntax/domain/roster membership is not ownership. Exact SKSU configuration and official format inputs plus normalization/changed-email reconciliation/course-contact requiredness remain OPEN-028; Student ID is no longer optional and Faculty has no institutional-domain gate.

Standard Excel creation is Admin-only and STUDENT-only. Recommended studentId/name/email/course/contactNumber; no password/admin role/staff role/faculty role or category assignment. Server fixes role BORROWER and borrower_type STUDENT. Validate complete file, required identity/domain/duplicates/ambiguity/conflicts; show complete preview before explicit selected-row confirmation; retain per-record results/safe retries/audit. Reject conflicting Faculty/privileged identities rather than relabeling. File selection never creates accounts.

Separate Admin Student bulk deactivation may reuse the roster. Match stable Student IDs/email, reject ambiguity/conflicts, show matched/unmatched/conflicting/already-inactive accounts and outstanding borrowing/fine/replacement obligations, then explicitly confirm selected Students. Recheck current role AND category under transaction locks, excluding Faculty/Staff/Admin even after preview changes. Preserve history and leave unselected/roster-absent accounts unchanged. DEC-073 resolves OPEN-029: outstanding fines/active or overdue loans/unreturned equipment/replacements do not block deactivation. Required warnings/confirmation preserve every balance, due time, unresolved state and overdue calculation; no payment, clearance, return, replacement resolution or history deletion. FSMO resolves separately; fine clearance remains Admin-only/auditable.

Both categories choose their own separate passwords in secure activation; single-use links target Student institutional email or Faculty accessible email, including non-institutional email. Hash-only purpose-bound expiring tokens/atomic consumption/rates/safe invalidation-reissue remain technical requirements. Brevo is selected for secure backend-only activation/future recovery (DEC-071); a configured API key, verified sender and successful real delivery testing (OPEN-017), plus activation/ownership/recovery (OPEN-001/009), remain unimplemented dependencies. Both require current officially published FSMO terms through Phase4B; official text remains OPEN-016. [Approved rules, technical dependencies and institutional approvals](../project/ACCOUNT_PROVISIONING_POLICY.md) are separate. Phase5 remains unstarted pending authorization; no implemented enforcement is claimed.

## Pending and return rules

**D17 authoritative interaction:** Borrower approaches Staff/Admin in person → may physically show a personally stored phone photo → Staff/Admin may inspect the actual equipment if desired → Staff/Admin alone records condition/quantities in eLabTrack → eLabTrack updates borrowing and inventory. The photo stays on the borrower's phone; this is an informal real-world convenience, not software evidence. The borrower cannot mark returned or submit return evidence. No return-photo acceptance, upload, transmission, storage, retention, attachment, file model or interface/workflow step is in scope. Equipment catalog images remain independently supported.

Submission captures server `submitted_at` and `expires_at=submitted_at+24h`; TTL is future configurable versioned policy candidate. No operating-hours-aware rule. At `now>=expires_at`, pending expiration wins over a newly attempted decision/cancellation and releases all holds exactly once. Background bounded expiry use case plus lazy locked checks at pending commands is the future engineering plan; no new scheduler code now.

Return lines distinguish good, damaged received and lost confirmed. Nonnegative integer components sum to positive processed quantity, at most physical outstanding. Repeated borrowing-item IDs reject the whole request even if totals fit. A single normalized map drives counters, ledger, incident, liability, audit and outbox. Damage/loss rationale is an engineering-required operational note recorded by Staff/Admin; return photographs are entirely outside eLabTrack. Partial/final returns cannot alter original due time. Replacement acceptance rejects duplicate obligation IDs, quantities beyond remaining and unrelated-borrowing obligations; each accepted unit adds new usable stock once. Historical incidents are never rewritten to “good”.

Completion occurs iff every item has physical0 and replacement0. It occurs atomically at final disposition/acceptance; no separate user-driven complete button is needed. It freezes the overdue fine and emits completion once. Accounts/equipment becoming inactive do not obstruct authorized resolution of old obligations. Archived equipment must not receive a new loan; recommendation: prevent archive while any hold/custody/replacement remains so the original pool remains a valid replacement destination.

## Overdue fines and full clearance

Authoritative engineering formula is shared across domain, HTTP projections, reminders and reports:

```text
unresolved = lifecycle CHECKED_OUT and (physical_outstanding>0 or replacement_outstanding>0)
currently_overdue = unresolved and server_now>due_at
effective_end = server_now if unresolved else completed_at for an issued completed borrowing
days = 0 if effective_end<=due_at else ceil((effective_end-due_at)/86,400 seconds)
assessed_minor = days*1,000 PHP centavos
outstanding_minor = assessed_minor - sum(immutable full-clear amounts)
```

Only issued borrowings have a fine basis. PENDING/DENIED/CANCELLED/EXPIRED assess zero. Completed late borrowing uses immutable final amount; it can be complete with OUTSTANDING money. At exact deadline0, one minute10PHP, exactly24h10PHP, 24h+one minute20PHP. No daily mutation, grace, cap, damage/loss valuation or partial-payment workflow. Original final assessment cannot be overwritten by a clear.

**Live-clear engineering recommendation:** Admin clears the entire positive outstanding at one locked server time; store the cumulative assessed amount at that time and the actual whole balance cleared, method, actor and optional note. If unresolved fine later increases, only that new delta becomes outstanding. E.g due+1h assessed10, Admin PAID clears10; due+25h assessed20, outstanding10; completion at25h freezes20; subsequent WAIVED clears remaining10. Recorded payment10 and waived10 remain distinct. Every clear is full at its own time, not partial against a chosen amount. Clearing never extends due or resolves replacement.

The command accepts no caller-selected payment amount. Recommended expected-balance token/amount prevents unknowingly clearing a changed balance; stale expectation conflicts. Serialize with borrowing/fine locks; same-key retry replays one event, a different-key attempt against a zero balance creates no second clear. If a 24h boundary created a new delta, a fresh explicit full-clear is legitimate. Queries return assessment, cleared total, outstanding, as-of time and live/final marker. Nonnegative bounds and integer overflow are authoritative errors; never fabricate negative balance or float money.

Fine clearance is **Admin only**, even when Staff can inspect it. PAID records FSMO's asserted offline payment; WAIVED and OTHER_RESOLUTION do not imply cash. No receipt/payment allocation schema or gateway. Formal financial receipt procedures and fine reversal/reopening are outside confirmed scope and cannot be inferred from the boss repository.

## Remaining policy questions and readiness

Core request/review/return/replacement/fine/role-navigation and Student deactivation-with-obligations policy are resolved. DEC-070 Student/Faculty rules remain unchanged; DEC-071 selects Brevo, DEC-073 resolves OPEN-029, and DEC-072 defers official terms until after presentation without blocking independent account-management/inventory development. Phase5 product policy is ready for a scoped plan; implementation remains unstarted and requires separate authorization. Genuine technical/external dependencies are actual SKSU domains/roster matching (OPEN-028), activation/ownership/recovery (OPEN-001/009), backend Brevo integration/API key/verified sender/successful live delivery testing (OPEN-017), and official terms publication/acceptance before live borrowing (OPEN-016). Taxonomy/content, original disposition/equivalence, retention/cutover/deployment and later notification cadence retain their separate scope. No borrowing/return/replacement/fine arithmetic or authorization is reopened.

Phase 3A is ready to begin **only after separate authorization**. Future mockups may assume held damaged originals, staff-certified equivalence, full request issue, current-version acceptance and live full-clear checkpoints, visibly marking draft terms/catalog copy. No UI, schema, business code or product tests have been implemented here.


## Phase 4B implemented terms boundary

2026-10-08: immutable terms publication/acceptance and first-use UI are implemented independently of official content approval. Migration 000005 adds `terms_versions`, singleton `terms_publication`, and `terms_acceptances`; no borrowing, inventory, policy/outbox or category persistence is implemented. Immediate mandatory Admin publication uses an expected-current pointer and shared/exclusive row locks, superseding proposed terms scheduling/advisory selection. Historical versions/receipts are append-only; a material update is a new version. Exact public DTOs/errors are in [API contracts](../API_CONTRACTS.md#phase-4b-implemented-terms-contracts).

The owner confirms no official V2 terms text is approved. Normal configured storage has no publication/acceptance; synthetic documents appear only in disposable tests. Missing content or service failure cannot authorize new borrowing. Borrower home/catalog paths request current consent; account/existing obligations/notifications and logout remain accessible. Server evidence is separate from authentication. DEC-068 confirms institutional-email login and a separate borrower-chosen eLabTrack password; single-use institutional-email activation links are recommended. Batch1 implements secure activation infrastructure and the backend Brevo adapter; live delivery, approved Student-domain inputs and institutional ownership/recovery procedures retain OPEN-001/009/017/028 gates. Password recovery remains outside this batch. See [current account provisioning policy](../project/ACCOUNT_PROVISIONING_POLICY.md).

Phase 7 must authorize the acting account and target Borrower, lock participating accounts in sorted order, then call `terms.Service.RequireCurrentAcceptance(ctx, borrowerID)` **inside the same borrowing command transaction**. Keep its shared publication lock until stock/borrowing/history/receipt commit and bind the returned receipt/borrower through `(id,user_id)` composite FK. Direct checkout checks the target's evidence, never records Staff consent for them. Existing pending submissions retain their original terms binding. No live borrowing endpoint or full business enforcement has been implemented. See [report](../project/PHASE4B_REPORT.md), DEC-066/067 and [tests](../../integration/PHASE4B.md).

## Phase 5 implementation boundary

Staff/Admin operational Borrower reads and Admin-only creation/profile/status/activation/roster/audit routes enforce current database authority. Borrower category is separate from role; no Admin creation/promotion is exposed. Account deactivation touches only account access, credential invalidation and immutable account audit/receipts. The accountability reader currently returns UNAVAILABLE, without fabricated zero balances. Tests inject positive active/overdue/unreturned/fine/replacement projections and verify they do not veto Student deactivation; Phase7/8 must supply real projections and their own historical conservation checks when those domains exist. Current Student-domain configuration gates new onboarding/activation, not safe deactivation of an exactly matched historical Student identity. See [implemented contracts](../API_CONTRACTS.md#phase-5-implemented-account-management-contracts).

## Batch1 Phase6 implemented inventory boundary

Paired000007 now implements operator-provided categories, catalog equipment metadata/lifecycle, four physical buckets, immutable movement/audit/command receipts and bounded equipment-only canonical PNG images. `T=A+R+C+D`; Staff/Admin usable acquisition/removal affects A/T only. Admin reviewed expected-sequence reconciliation sets observed A and preserves R/C/D; no custody or incident correction. Current per-pool quantity bound is2147483647; all arithmetic is checked server-side and by PostgreSQL. Loss history and replacement liability remain separate future concepts, with no invented current lost bucket or zero liability projection.

Equipment ACTIVE is Borrower-visible; INACTIVE/ARCHIVED is operational-only. Archive retains current stock/history, blocks R/C, and requires verified absence of future borrowing/replacement tables; their presence without integrated liability checks fails closed. No archival return/settlement, repair/disposal, borrowing, reservation, replacement or fine workflow exists. Metadata and stock versions are separate; locked transactions make movements/audit/receipts atomic, idempotent and stale-basis safe. Category inactivity preserves existing references. Local database image support is bounded PNG/JPEG re-encoding for catalog use only, not return photographs or arbitrary storage. Actual routes/limits/storage/permissions and future integration requirements are in [API contracts](../API_CONTRACTS.md#phase-6-implemented-equipment-and-inventory-contracts). Historical Phase2/4 candidate/unimplemented descriptions above are dated design evidence, superseded for these scoped implemented surfaces.
