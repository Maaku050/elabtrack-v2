# Phase 2.5 domain model

Rebaselined 2026-10-08. **CONFIRMED V2 DECISION** means the project owner's authoritative working Phase 2.5 decisions, not a claim of deployed behavior or independent institutional signoff. Relational structures, locks and arithmetic below are **ENGINEERING RECOMMENDATIONS** selected for the design baseline. [Decision register](../project/DECISIONS.md), [source hierarchy](../project/SOURCE_OF_TRUTH.md) and [Phase 2.5 report](../project/PHASE2_5_REPORT.md) govern. Phase 1A–1I remains implemented and unchanged; this document implements nothing.

## Scope and modules

One FSMO operation, one SPA/API/database, aggregate equipment pools. No tenancy, campus hierarchy, serialized assets, cross-department borrowing or native mobile application.

| Domain boundary | Concepts and responsibilities | Dependencies through ports |
|---|---|---|
| Identity / borrower management | Named User, active status, BORROWER/STAFF/ADMIN capabilities, borrower category, provisioning/bulk import | Existing identity/session infrastructure; future account-management use cases |
| Equipment / catalog | Equipment pool, metadata, optional category/image, lifecycle; interactive search/cart experience | Inventory availability projection; bounded catalog queries |
| Inventory | Current quantities and append-only movement ledger | Borrowing/return/replacement operations share transaction and lock protocol |
| Borrowing | Request, reservation, atomic approval/release, direct issue, due time, expiry and retained terminal history | Current account, equipment, terms/policy, transaction, clock, audit/outbox ports |
| Returns / replacements | Immutable ReturnEvent, damage/loss incident, ReplacementObligation and acceptance | Borrowing lock, stock acquisition/movement, final completion |
| Overdue fines | Dynamic assessment, immutable final assessment, Admin full-clear evidence | Borrowing due/completion basis and policy; no equipment-price liability |
| Notifications | Transactional intent, delivery attempts, due/overdue scheduling | Business event IDs and future provider adapter |
| Reporting | Scoped bounded projections over canonical history, stock, replacement and fine evidence | Same definitions as operational use cases |

Domain owns concepts/invariants/ports. Application coordinates use cases and transactions. Infrastructure implements pgx repositories; HTTP handlers adapt Fiber requests; bootstrap wires dependencies. Domain/application do not import Fiber or pgx. No event sourcing, generalized finance platform or additional services.

## Identity and vocabulary

**BORROWER** is the generic business actor. Expected categories include Student and Faculty; category is an account attribute, never a permission role. Active registered borrower accounts can initiate borrowing; no separate eligibility/suspension table or automatic outstanding-fine block exists. Account deactivation blocks login and normal authenticated borrower actions while retaining all obligations. Staff/Admin can still record returns/replacements for that account.

The product roles are BORROWER, STAFF and ADMIN; multiple named Admins are required. Staff can provision Borrower accounts (Decision 4's narrow grant). Admin controls deactivation, bulk borrower import, Staff/Admin management, full fine clearance, exceptional inventory corrections, policies/terms and administrative audit/reports. Neither borrower category nor a client route grants authority. [The permission matrix](BUSINESS_RULES.md#future-permission-matrix) is the future server contract. Phase 4A now implements this role vocabulary: verified migration 000004 maps foundation `user`→BORROWER and `admin`→ADMIN, permits STAFF and preserves non-role account/session data. Feature capabilities beyond the existing authenticated/account-directory boundaries remain future implementation contracts; see [Phase 4A report](../project/PHASE4A_REPORT.md).

No public self-registration. Bulk onboarding is required, but V1's plaintext password spreadsheet is not a V2 contract. Activation/initial-password mechanics remain later security design work. Optional program/contact data must have a purpose and minimal retention; no capstone sample identities are engineering fixtures.

## Canonical borrowing

A Borrowing has stable ID/reference, borrower ID, entry path REQUEST or DIRECT, policy version and terms-acceptance reference, timestamps, required issued `due_at`, and at least one unique equipment line. User/account display and item descriptions are snapshotted for historical interpretation; catalog edits never rewrite history. Price, if retained as informative catalog metadata, does not assess damage/loss fines.

Stored lifecycle: **PENDING, CHECKED_OUT, DENIED, CANCELLED, EXPIRED, COMPLETED**. Partial return, overdue and replacement awaiting resolution are projections, not lifecycle enum combinations. Request submission reserves stock immediately; pending expires 24 hours after server submission. Approval is physical handover: one transaction moves PENDING to CHECKED_OUT. Direct checkout starts CHECKED_OUT without a pending hold. Denied/cancelled/expired rows remain canonical history.

Each BorrowingItem records requested/reserved/issued quantity and cumulative good/damaged/lost dispositions reconciled to immutable ReturnEventItems. At issue, quantity, borrower/equipment snapshots and due/policy evidence freeze. Recommended initial issue is the entire request or none; staff cannot quietly edit quantities into a different loan. Staff confirms the due date/time at handover; every issued loan has `due_at`, interpreted locally in Asia/Manila and stored as an absolute timestamp. No seven-day maximum or new quota is imposed.

For each issued item:

```text
physical_outstanding = issued - good_returned - damaged - lost
replacement_required = sum(damage/loss obligation.required_qty)
replacement_accepted = sum(immutable acceptance lines.accepted_qty)
replacement_outstanding = replacement_required - replacement_accepted
complete item iff physical_outstanding = 0 and replacement_outstanding = 0
complete borrowing iff every item is complete
```

A loss disposition records confirmed loss without pretending a physical item arrived. A damage disposition means staff received a damaged original; an item still with the borrower remains physically outstanding until staff records an authoritative disposition. Damage/loss creates replacement obligations exactly for those quantities. Fine clearance does not complete the borrowing; an outstanding monetary fine does not prevent operational completion.

## Selected inventory model and alternatives

Re-evaluated options: (A) six buckets including lost/retired conflates historical incidents with current physical stock; (B) only usable stock plus incidents hides reserved custody and nonusable originals; (C) **selected hybrid** uses four physical quantities plus separate incident/obligation evidence. C makes current stock explicit without treating replacement liability as physical custody.

```text
total_tracked = available + reserved + checked_out + damaged_held
usable_pool = available + reserved + checked_out
on_premises = available + reserved + damaged_held
```

`total_tracked` includes units still in borrower custody; it is not a count of units presently on the FSMO shelf. `damaged_held` is nonusable originals in FSMO custody. Lost units cease to be physically tracked; their immutable loss incidents survive. Retired/disposed units leave tracked stock through a later authorized ledger movement. No current lost/retired bucket or device table is required. UI “damaged awaiting replacement” and “lost awaiting replacement” derive unresolved obligations by incident kind, not `checked_out` or `damaged_held` counts.

Let a movement be `(available, reserved, checked_out, damaged_held; total_tracked)`:

| Operation | Delta vector | Obligation / history effect |
|---|---|---|
| Opening/add usable stock q | `(+q,0,0,0;+q)` | Explicit baseline/acquisition evidence |
| Submit q | `(-q,+q,0,0;0)` | PENDING hold |
| Cancel/deny/expire q | `(+q,-q,0,0;0)` | Terminal request retained |
| Approve and hand over q | `(0,-q,+q,0;0)` | Issue and decision same edge |
| Direct checkout q | `(-q,0,+q,0;0)` | Immediate issued borrowing |
| Good return g | `(+g,0,-g,0;0)` | Original physical custody decreases |
| Damaged original received d | `(0,0,-d,+d;0)` | Create d damage replacement units |
| Lost units confirmed l | `(0,0,-l,0;-l)` | Create l loss replacement units |
| Mixed disposition g,d,l | `(+g,0,-(g+d+l),+d;-l)` | One normalized map for every effect |
| Accept new replacements q | `(+q,0,0,0;+q)` | Reduce liability q; acquisition, not another original return |
| Later authorized discard q damaged units | `(0,0,0,-q;-q)` | Incident/acceptance history remains; later inventory policy |
| Later authorized repair q damaged units | `(+q,0,0,-q;0)` | Nonmandatory later inventory policy; never resolves liability automatically |

All stock vectors obey `delta_total = sum(delta_buckets)`. Loss followed by replacement restores original total. Damage followed by replacement while retaining the original yields two actual physical units: one damaged and one usable. This increase is correct acquisition arithmetic, not double counting. Acceptance never decrements the damaged original implicitly, changes historical damage/loss quantities, or adds good-return counts.

Examples from initial `(5,0,0,0;5)`: issue all gives `(0,0,5,0;5)`; return 2 good and report 3 lost gives `(2,0,0,0;2)`, liability3; accept3 gives `(5,0,0,0;5)`, liability0 and completion. For damage2, return3 good gives `(3,0,0,2;5)`, liability2; accept2 gives `(5,0,0,2;7)`, liability0. Later discarding2 originals gives `(5,0,0,0;5)` without changing incidents or completion.

## Replacement evidence

One ReplacementObligation per immutable return-line and incident kind DAMAGE or LOSS; required quantity is frozen. Remaining quantity is required minus immutable acceptance quantities. ReplacementEvent has named Staff/Admin actor, time, borrowing and optional equivalence rationale; each unique ReplacementEventItem references one obligation belonging to that same borrowing. The staff operator certifies appropriate equipment type or operationally accepted interchangeable equivalent for the original pool; no speculative cross-SKU conversion model is introduced.

Reject duplicate obligation IDs in an acceptance command, overacceptance and unrelated-parent references. Required and remaining quantities cannot be negative. An acceptance adds new usable stock exactly once and leaves the incident intact. Last physical disposition/acceptance may complete atomically; physical zero with replacement positive remains CHECKED_OUT. Detailed original disposition and equivalence guidance are NON-BLOCKING for core borrowing UX, with default originals held outside usable stock.

## Overdue fine and clearance

A CHECKED_OUT borrowing is operationally unresolved while either physical or replacement quantity remains. It is overdue iff unresolved and `now > due_at`. PENDING/DENIED/CANCELLED/EXPIRED never accrue. A completed borrowing retains historical late duration and a frozen final fine, without appearing currently overdue.

Selected engineering interpretation of confirmed PHP 10/day:

```text
effective_end = server now while CHECKED_OUT; completed_at when COMPLETED
elapsed = max(0, effective_end - due_at)
overdue_days = ceil(elapsed / 86,400 seconds)
assessed_minor = overdue_days * 1,000  # PHP centavos
```

Use bounded integer duration/ceiling arithmetic, no float or daily fine mutation. Exact due time:0; one minute late:PHP 10; exactly24h:PHP 10; 24h+one minute:PHP 20. Partial returns and replacement-only open loans retain original due time. Completion freezes final amount and days in the same transaction. No grace period, cap, compound rate, or equipment-value assessment is introduced.

**Selected engineering recommendation for live clearance:** one overdue_fines row per issued borrowing holds immutable policy/due basis and once-only final assessment; fine_clearances records each Admin full-clear of the outstanding balance at captured server time. Before completion, clearing does not stop further accrual. A later increase is new outstanding balance, requiring a later full clear; this is not installment/payment allocation. Preserve cumulative assessment at every clearance and final assessment at completion.

```text
outstanding = current assessed (or frozen final) - sum(clearances.cleared_amount_minor)
clear command amount = entire locked positive outstanding at command time
balance label = OUTSTANDING if positive, otherwise CLEARED
```

A zero-fine loan may display “no fine” rather than implying a payment. Clearances record PAID, WAIVED or OTHER_RESOLUTION, named Admin, time and optional note. Only PAID is a recorded payment; all clearances are not revenue. No generic charges/adjustments, payment gateway, installment, allocation or refund system is retained. See [fine contract](BUSINESS_RULES.md#overdue-fines-and-full-clearance).

## Terms, history, notifications and files

Versioned immutable TermsAcceptance belongs to user/version/time; once per current version, not once per loan. Submission references existing current-version acceptance. A material version change gates new requests; it does not retroactively invalidate pending submission evidence. Recommended direct issue requires current borrower acceptance evidence; staff cannot accept on their behalf. Published borrowing policy captures typed 24h TTL, Asia/Manila, PHP 10/24h ceiling basis and nullable future maximum duration; no generic rules engine.

Keep canonical borrowing history, immutable return/replacement events, stock ledger, fine clearance/final basis and business audit distinct from operational logs. Required audit/outbox/command receipt writes share each business transaction. Notification intent supports submission, denial, approval/checkout, expiry, due, overdue, partial return, replacement creation/resolution and completion. No provider is implemented. Every partial return/acceptance uses its own event ID.

Optional equipment images remain catalog-purpose attachments, not a generic file/export authorization model. **D17 clarification: return photographs are entirely outside eLabTrack.** A borrower may take a photo on their own phone and physically show that personally stored photo to Staff/Admin during a face-to-face return. Staff/Admin may inspect the actual equipment if desired, then alone records authoritative condition/quantities in eLabTrack; the system updates borrowing and inventory. The borrower cannot mark a transaction returned or submit return evidence through the software.

eLabTrack does not accept, upload, transmit, store or retain return photos, create return-photo attachments, or provide a return-evidence interface, file/storage model or upload step. This informal physical convenience creates no software feature. Structured ReturnEvent quantities/condition/notes, ledger and audit records remain required business history; references to return evidence mean those records, never photographic media.

## Interactive Kiosk and responsive boundary

**Interactive Kiosk** means the normal borrower catalog/search/filter/quantity/cart/review/request experience. Borrower mobile-first; Staff/Admin desktop/tablet-first responsive. No dedicated terminal shell, device credentials, registration, shared-device auth or handoff sessions. This supersedes the earlier kiosk portion of DEC-038/039; academic terminology may remain with this functional meaning.

[State machines](STATE_MACHINES.md), [schema candidates](DATA_MODEL.md), [invariants](INVARIANTS.md), [V1 compatibility](V1_COMPATIBILITY.md), [boss salvage map](BOSS_REBUILD_REFERENCE.md) and [future REST draft](API_RESOURCE_DRAFT.md) are one design baseline. Readiness for future approved planning is not permission to implement any phase.
