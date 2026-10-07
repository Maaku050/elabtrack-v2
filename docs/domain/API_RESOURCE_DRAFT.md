# Phase 2.5 future REST resource draft

2026-10-08. **ENGINEERING RECOMMENDATION**, unimplemented: `/api/v1`, existing consistent envelope `{success,message,data,meta,error}`, safe server-generated request ID and Phase 1 HTTP/session protections. Product capabilities follow [permission matrix](BUSINESS_RULES.md#future-permission-matrix), not route guards. No handler/source or auth role has changed. Resource naming/casing are future implementation contracts, not permission to start a phase.

## Proposed resources and commands

All private resources authenticate current account, validate actor capability and scope to borrower/operational purpose before paging/detail/replay. Actor/role/state/fee totals are server-derived. Commands marked key require Idempotency-Key and canonical normalized payload; authorization is repeated on replay.

| Draft route | Actor / purpose | Required boundary |
|---|---|---|
| GET /equipment; GET /equipment/{id} | Borrower/Staff/Admin catalog/search/filter | Bounded visible pool metadata/image/availability; no borrower/private evidence |
| POST /equipment; PATCH /equipment/{id} | Staff/Admin ordinary metadata | Narrow allowed fields, expected metadata version; no raw stock fields |
| POST /equipment/{id}/archive | Staff/Admin | No hold/physical/replacement outstanding; history retained; no ordinary hard delete |
| GET /equipment/{id}/movements | Staff/Admin operational ledger | Bounded filters; actor/source quantities; no secrets |
| POST /equipment/{id}/stock-movements | Staff/Admin ordinary acquisition/removal | Key; typed vectors/reason; cannot alter hold/custody. Exceptional correction Admin only; disposal/repair commands await later inventory policy |
| GET /users/me | Existing authenticated self | Current Phase 1 contract unchanged; future role/category projection only under approved feature |
| GET /borrowers; GET /borrowers/{id} | Staff/Admin operational selection/history | Active selection for new issue; old inactive obligations still resolvable; bounded minimal profile |
| POST /borrowers | Staff/Admin Borrower provisioning | Key; role forced BORROWER, secure onboarding later; no public signup |
| POST /borrowers/{id}/deactivate; POST /borrowers/{id}/reactivate | Admin | Key; retained obligations/audit; no identity delete |
| POST /borrowers/import | Admin bulk provision | Future validated secure batch/preview/reconciliation contract; no plaintext-password mandate |
| POST /staff-accounts; PATCH /staff-accounts/{id} | Admin Staff/Admin management | Key for irreversible changes; current actor authority, named users and privilege recovery design; no shared admin |
| GET /borrowing-policy; GET /terms/current | Registered account policy/content | Current published version/zone/TTL/fine basis; no generic tenant settings |
| POST /terms/{versionId}/accept | Borrower own acceptance | Key; current published terms, unique user/version; no staff accepting on behalf |
| POST /borrowing-policy/versions; POST /terms/versions | Admin | Immutable published typed versions/content; audited effective boundary |
| POST /borrowings | Borrower own cart submission | Key; unique equipment/positive quantity, requested due intent, existing current terms; reserve atomically; old fine allowed |
| GET /borrowings; GET /borrowings/{id} | Own Borrower; Staff/Admin operational view | Stable scoped history, physical/replacement/fine projections; lifecycle/active/due filters validated |
| POST /borrowings/{id}/approve | Staff/Admin approval WHILE handover | Key; PENDING/unexpired, active target, confirmed due_at, held stock; one CHECKED_OUT edge, no separate checkout route |
| POST /borrowings/{id}/deny | Staff/Admin | Key; PENDING/unexpired, required borrower-visible reason; release/history |
| POST /borrowings/{id}/cancel | Borrower owner only | Key; PENDING/unexpired; release/history; administrative cancellation optional unselected |
| POST /borrowings/direct-checkout | Staff/Admin immediate issue | Key; active existing borrower/current terms, unique quantities, future due_at, current availability |
| POST /borrowings/{id}/returns | Staff/Admin authoritative good/damage/loss | Key; unique item IDs, integer quantities<=physical remaining; immutable evidence/obligations/stock; closure if both0 |
| GET /borrowings/{id}/returns | Owner Borrower; Staff/Admin | Immutable chronological events; no writable borrower returned checkbox |
| GET /borrowings/{id}/replacement-obligations | Owner Borrower; Staff/Admin | Required/accepted/remaining by kind/source/equipment, historical incident retained |
| POST /borrowings/{id}/replacements | Staff/Admin acceptance | Key; unique obligation IDs, integer accepted<=remaining, appropriate type/equivalence evidence; new stock acquisition, possible completion |
| GET /borrowings/{id}/replacements | Owner Borrower; Staff/Admin | Immutable acceptance history, no edit/delete |
| GET /fines; GET /fines/{id} | Own Borrower; Staff/Admin operational review | Live/final assessed, original basis, clear history, outstanding, as_of; no automatic issue block |
| POST /fines/{id}/clear | ADMIN ONLY | Key; method PAID/WAIVED/OTHER_RESOLUTION, optional note, expected outstanding; server calculates entire current balance; no amount chosen by caller |
| GET /fines/{id}/clearances | Owner Borrower; Staff/Admin operational context | Preserved basis/amount/time/Admin/method; no inferred bank/gateway evidence |
| GET /reports/inventory, /reports/borrowings, /reports/replacements, /reports/fines, /reports/usage | Admin administrative reports | Bounded approved measures; assessed/PAID/waived/other distinct; Staff ordinary summaries require narrow future operation grant |
| GET /audit-events | Admin restricted audit | Bounded entity/action/date/actor; safe minimal context, no credentials/provider raw payload |

Expiry is an internal system command/sweep, not borrower-controlled endpoint; it releases holds at locked expiry and writes history/outbox. Due/overdue reminders are internal future scheduled uses of the same formula. No `/charges`, adjustment/payment-allocation/device/handoff routes, separate approved-release route, mutable `/records` copy, arbitrary state PATCH or consequential DELETE. Return photographs are entirely outside eLabTrack. No borrower return-evidence submission, return-photo upload/transmission/storage/retention, attachment API or evidence-upload step exists in this draft; Staff/Admin return input has only structured condition/quantity/operational-note fields and rejects return-photo/media fields. Kiosk uses ordinary /equipment and /borrowings resources, no special auth namespace.

## Mobile-first request and response projection

| Borrower interaction | Draft response semantics |
|---|---|
| Browse/search/filter/cart | Concise bounded catalog; client cart is local intent, not a hold until submission commits |
| Submit | Canonical ID/ref, PENDING, reserved quantities, submitted_at/expires_at and bound terms/policy; no claim of checkout |
| Status/history | Stored lifecycle plus derived partial/overdue/replacement conditions; timestamped denial reason visible |
| Active issued borrowing | issued/good/damaged/lost/physical outstanding plus replacement required/accepted/remaining per item; immutable due_at |
| Accountability | Physical/replacement obligations separate from live/final PHP fine and clearance history; show as_of and fine-final marker |
| Completion | Server final event/time, physical0/replacement0, frozen final fine; money may remain outstanding |

Staff approval view shows outstanding accountability for human review; server does not auto-deny old fine. At the face-to-face return the borrower may physically show a personally stored phone photo; Staff/Admin may inspect actual equipment if desired and alone records authoritative condition/quantities in eLabTrack. The photo never enters the software, request payload, attachment, storage or notification. Borrowers cannot mark returned or submit return evidence through eLabTrack. At fine clear, API returns original cumulative assessment, full cleared amount, remaining0 at captured time, method/event/time/Admin, live/final marker and warning copy that an unresolved loan can accrue later (product wording, not implementation internals).

## Validation, pagination and errors

Recommend current `page`/`per_page` conventions default20/max100, stable `(created_at,id)` or `(name,id)`, explicit owner filters before count/page; cursor later only if measured need. These are engineering bounds, not borrowing quantity/time quotas. Bound search/body/item counts, validate enums/date offsets/Asia-Manila interpretation, allowlist sort columns and parameterize SQL. Scoped counts must not leak others' private records. Authenticated TanStack Query data must use existing Phase 1I cache classification/generation fences in later frontend implementation.

HTTP success follows commit. A late pending command can **commit expiry** and then return409 BORROWING_EXPIRED with the authorized new state; handler must not accidentally roll back the release. Same-key valid receipt replays committed result identity with new current request ID after reauthorization. Different key/state command fails safely; no automatic retry of every 409.

| Proposed code | HTTP | Meaning |
|---|---|---|
| VALIDATION_ERROR | 400 | Invalid quantities/dates/required denial reason/type/clear method |
| DUPLICATE_RETURN_ITEM / DUPLICATE_REPLACEMENT_OBLIGATION | 400 | Duplicate line ID; whole command rejected before effects |
| ACCOUNT_INACTIVE / BORROWER_NOT_ACTIVE | 403 | Current caller/authorized selected borrower inactive; no invented general eligibility state |
| TERMS_ACCEPTANCE_REQUIRED | 409 | Current version must be accepted before new submission/direct first-use |
| EQUIPMENT_NOT_AVAILABLE | 409 | Locked available/held quantity or usable status cannot fulfill intent |
| BORROWING_STATE_CONFLICT | 409 | Legal lifecycle/ownership/version precondition changed |
| BORROWING_EXPIRED | 409 | Expiry reached; committed release/state reported safely |
| RETURN_EXCEEDS_OUTSTANDING | 409 | Disposition exceeds locked physical remaining |
| REPLACEMENT_EXCEEDS_OUTSTANDING | 409 | Acceptance exceeds locked unresolved replacement |
| STOCK_CONFLICT | 409 | Archive/count/metadata version/integrity conflict |
| FINE_ALREADY_CLEARED / STALE_FINE_BALANCE | 409 | No current positive outstanding / expected full balance changed |
| IDEMPOTENCY_CONFLICT | 409 | Existing scoped key with different normalized payload |

Authentication401, forbidden operation403, unauthorized resource404 where appropriate, existing parser/body/rate/server responses retain Phase 1 contracts. Do not disclose private borrower history/amounts to unauthorized callers. No business HTTP implementation/tests run in this task.
