# Phase 2.5 business rules and authorization draft

2026-10-08. Owner decisions D1–D29 are **CONFIRMED V2 DECISIONS / authoritative working product policy**. Engineering choices are called out separately. They override conflicting capstone/V1/boss assumptions; none is implementation evidence. See [source policy](../project/SOURCE_OF_TRUTH.md), [accepted decisions](../project/DECISIONS.md), [remaining questions](../project/OPEN_DECISIONS.md).

## Decision integration matrix

D numbers refer to the owner's Phase 2.5 instruction, with permanent register IDs below.

| Owner decision | Current working rule | Register / principal design effect |
|---|---|---|
| D1 | Active registered BORROWER, including at least Student/Faculty categories; category is not role | DEC-051; user attributes, no StudentRole/FacultyRole |
| D2 | Active account is borrowing gate; Admin deactivation blocks normal auth actions; no eligibility subsystem or fine block | DEC-051; reuse identity status, preserve existing obligations |
| D3 | BORROWER/STAFF/ADMIN; named multiple Admins; Admin inherits staff operations | DEC-052; future server permission matrix |
| D4 | No public signup; Staff/Admin provision Borrowers; Admin bulk import and privileged accounts | DEC-053; secure activation mechanics later, not plaintext spreadsheet contract |
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

Active registered Borrower is sufficient to request; no email-verification, student-only, course, suspension, fine balance or active-loan-count restriction is silently added. Borrower category does not grant staff authority. Staff's account grant is Borrower provisioning only; Admin controls deactivation, bulk imports and Staff/Admin administration. No shared admin credential. Accounts with loans are deactivated rather than deleted; staff processes their remaining obligations. Activation/password delivery and optional email ownership verification are later design questions, not new eligibility policy.

Terms acceptance is unique per user/version and immutable. Current material version must be accepted before a new submission. Recommended direct checkout requires the borrower to have accepted current terms; Staff cannot synthesize that acceptance. Submission binds its accepted version; publishing new terms does not force reacceptance for an already pending request. Final approved text and secure onboarding must be reviewed before production use, without blocking navigation mockups.

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
| Provision individual Borrower | No | Yes, narrow D4 grant | Yes | Cannot assign STAFF/ADMIN or clear fines |
| Deactivate/reactivate Borrower | No | No | Yes | Existing history/obligations preserved; reactivation engineering lifecycle counterpart |
| Bulk Borrower import | No | No | Yes | Future validated secure importer; preview/reconciliation |
| Create/manage STAFF/ADMIN | No | No | Yes | Named actors, anti-escalation/last-admin recovery design later |
| Clear fine (PAID/WAIVED/OTHER_RESOLUTION) | No | No | Yes | Entire current balance only; actor/time/basis retained |
| Manage policy/terms/settings | No | No | Yes | Immutable published versions; audit |
| View administrative audit/reports | No | No | Yes | Bounded purpose-specific access; operational summaries for Staff do not imply admin exports |
| Administrative pending cancellation | No | Not yet selected | Not yet selected | Optional recommendation, needs future scope decision |
| Automatic pending expiry | No client command | No client command | No client command | System use case, deadline and stock locks, explicit system audit actor |

Normal inventory means metadata maintenance, supported usable stock acquisition/removal and archive under guards; detailed damage repair/disposal/correction procedures remain later inventory policy. Hard deletion of consequential history is not inventory CRUD. React guards are UX only. Permission changes and deactivation writers must participate in the transaction-time account-lock protocol; existing session behavior is not modified here.

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

Core request, review, return, replacement, fine and role/navigation policy gates are resolved. [OPEN_DECISIONS](../project/OPEN_DECISIONS.md) retains precise NON-BLOCKING details: damaged-original disposition/equivalence guidance; taxonomy/content; secure activation/email ownership mechanics; retention/cutover authorization; provider/cadence and deployment operations. Disposal feature, production delivery and real migration remain dependent on their own approvals/details. These do not justify reopening reserve-on-submit, merged approval, PHP 10/day, replacement liability, full-clear or kiosk meaning.

Phase 3A is ready to begin **only after separate authorization**. Future mockups may assume held damaged originals, staff-certified equivalence, full request issue, current-version acceptance and live full-clear checkpoints, visibly marking draft terms/catalog copy. No UI, schema, business code or product tests have been implemented here.
