# Phase 2.5 future REST resource draft

**Current Batch1 implementation,2026-10-09:** Phases5–6 now implement the scoped account/catalog/inventory boundaries described in the overlays below and current API_CONTRACTS.md. Earlier unimplemented/candidate statements are historical Phase2/4 design evidence, not current status. See PHASE5_REPORT.md (external-gate partial) and PHASE6_REPORT.md (reviewed inventory scope complete). Phase7 remains unauthorized.


2026-10-08. **ENGINEERING RECOMMENDATION**, unimplemented: `/api/v1`, existing consistent envelope `{success,message,data,meta,error}`, safe server-generated request ID and Phase 1 HTTP/session protections. Product capabilities follow [permission matrix](BUSINESS_RULES.md#future-permission-matrix), not route guards. Phase 4A implements only existing auth/current-account/privileged-directory boundaries with BORROWER/STAFF/ADMIN. Phase 4B now implements the four terms contracts documented in [API_CONTRACTS](../API_CONTRACTS.md#phase-4b-implemented-terms-contracts); borrowing and other business resources remain unimplemented. Phase5 account routes now have [implemented contracts](../API_CONTRACTS.md#phase-5-implemented-account-management-contracts); the remaining draft resources are future recommendations, not permission to start a phase.

## Proposed resources and commands

All private resources authenticate current account, validate actor capability and scope to borrower/operational purpose before paging/detail/replay. Actor/role/state/fee totals are server-derived. Commands marked key require Idempotency-Key and canonical normalized payload; authorization is repeated on replay.

| Draft route | Actor / purpose | Required boundary |
|---|---|---|
| GET /equipment; GET /equipment/{id} | Borrower/Staff/Admin catalog/search/filter | Bounded visible pool metadata/image/availability; no borrower/private evidence |
| POST /equipment; PATCH /equipment/{id} | Staff/Admin ordinary metadata | Narrow allowed fields, expected metadata version; no raw stock fields |
| POST /equipment/{id}/archive | Staff/Admin | No hold/physical/replacement outstanding; history retained; no ordinary hard delete |
| GET /equipment/{id}/movements | Staff/Admin operational ledger | Bounded filters; actor/source quantities; no secrets |
| POST /equipment/{id}/stock-movements | Staff/Admin ordinary acquisition/removal | Key; typed vectors/reason; cannot alter hold/custody. Exceptional correction Admin only; disposal/repair commands await later inventory policy |
| GET /users/me | Existing authenticated self | Current Phase 1 contract unchanged; product role projection implemented in Phase 4A; category projection remains deferred |
| GET /borrowers; GET /borrowers/{id} | Staff/Admin operational selection/history | Active selection for new issue; old inactive obligations still resolvable; bounded minimal profile |
| POST /borrowers | Admin-only individual Student/Faculty creation (future) | Key; current backend Admin check, role fixed BORROWER, borrower_type STUDENT/FACULTY; Student required unique string Student ID/approved SKSU email, Faculty valid unique accessible email/no required ID; no Admin-supplied borrower password, public signup or privileged role assignment |
| POST /borrowers/{id}/deactivate; POST /borrowers/{id}/reactivate | Admin | Key; Student deactivation requires warnings/confirmation but permits all outstanding obligations (DEC-073); status/audit only, no borrowing/stock/fine/overdue/resolution/history changes; reactivation separate; no identity delete |
| POST /borrowers/import; GET /borrowers/import/{batchId} | Admin-only Student Excel preparation/complete preview/results (future) | studentId/name/email/course/contactNumber; required unique textual official Student ID and configured SKSU email; duplicate/conflict checks; no password/role/Faculty/category assignment; complete preview has no account mutations; bounded owned batch/results |
| POST /borrowers/import/{batchId}/confirm | Admin-only Student creation confirmation (future) | Key; server fixes role BORROWER/borrower_type STUDENT; binds reviewed rows/payload, reauthorizes/revalidates current identities, rejects Faculty/privileged conflicts; safe retries/durable per-record results/audit; no mutation on file selection |
| POST /borrowers/bulk-deactivations; GET /borrowers/bulk-deactivations/{batchId} | Admin-only Student deactivation preparation/preview/results (future) | Reuse Student roster; stable Student ID/email matching; reject ambiguity/conflicts/Faculty/Staff/Admin; matched/unmatched/conflicts/already-inactive and authoritative borrowing/fine/replacement warnings |
| POST /borrowers/bulk-deactivations/{batchId}/confirm | Admin-only selected Student deactivation (future) | Key; current BORROWER AND STUDENT/identity/status/obligation checks under locks; exclude Faculty/Staff/Admin including category/role changes since preview; history/audit/results preserved; unselected/roster-absent unchanged; DEC-073: warnings/confirmation, no rejection based on positive outstanding obligations; no loan/fine/stock/replacement/due/overdue/history changes |
| POST /auth/activate | Valid single-use activation-token capability (future proposed name) | Existing provisioned Student/Faculty BORROWER only, no registration; verify Student institutional or Faculty accessible mailbox; token plus newly chosen eLabTrack password, never mailbox password; current lifecycle/account authority, hash-only token persistence/expiry/atomic consumption, rates/safe invalidation/reissue |
| POST /staff-accounts; PATCH /staff-accounts/{id} | Admin Staff/Admin management | Key for irreversible changes; current actor authority, named users and privilege recovery design; no shared admin |
| GET /borrowing-policy | Registered account policy/content (future) | Current published zone/TTL/fine basis; no generic tenant settings |
| GET /terms/current; GET /terms/status | IMPLEMENTED Phase 4B: active account content / Borrower own status | Exact current immutable document; fail closed missing; no institutional text approved |
| POST /terms/{versionId}/accept | IMPLEMENTED Phase 4B: Borrower own acceptance | Body `{}`; exact current UUID; unique user/version returns original receipt/time, no Idempotency-Key needed; no staff accepting on behalf |
| POST /borrowing-policy/versions | Admin (future) | Immutable published typed policy; reviewed effective boundary |
| POST /terms/versions | IMPLEMENTED Phase 4B: current Admin only | Immediate immutable mandatory content, publisher/server time/hash; required expected_current_version_id prevents concurrent overwrite; approved content still gated |
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


## Phase 4B implemented terms boundary

2026-10-08: immutable terms publication/acceptance and first-use UI are implemented independently of official content approval. Migration 000005 adds `terms_versions`, singleton `terms_publication`, and `terms_acceptances`; no borrowing, inventory, policy/outbox or category persistence is implemented. Immediate mandatory Admin publication uses an expected-current pointer and shared/exclusive row locks, superseding proposed terms scheduling/advisory selection. Historical versions/receipts are append-only; a material update is a new version. Exact public DTOs/errors are in [API contracts](../API_CONTRACTS.md#phase-4b-implemented-terms-contracts).

The owner confirms no official V2 terms text is approved. Normal configured storage has no publication/acceptance; synthetic documents appear only in disposable tests. Missing content or service failure cannot authorize new borrowing. Borrower home/catalog paths request current consent; account/existing obligations/notifications and logout remain accessible. Server evidence is separate from authentication. DEC-070 confirms email plus a separate borrower-chosen eLabTrack password for both categories: Student requires approved SKSU institutional email and unique textual Student ID; Faculty may use any valid unique accessible email and requires no Student ID. Only Admin creates accounts; standard bulk creation/deactivation is Student-only. Single-use activation links target each category’s respective mailbox. DEC-071 selects Brevo for backend-only activation and future password recovery; integration and API-key/verified-sender/successful live delivery verification remain unimplemented (OPEN-017), alongside activation/ownership/recovery and actual SKSU inputs (OPEN-001/009/028). DEC-072 schedules official FSMO terms after application presentation: pending content does not block independent account-management/inventory development under separate authorization, but live borrowing still requires official publication and documented acceptance. See [current account provisioning policy](../project/ACCOUNT_PROVISIONING_POLICY.md).

Phase 7 must authorize the acting account and target Borrower, lock participating accounts in sorted order, then call `terms.Service.RequireCurrentAcceptance(ctx, borrowerID)` **inside the same borrowing command transaction**. Keep its shared publication lock until stock/borrowing/history/receipt commit and bind the returned receipt/borrower through `(id,user_id)` composite FK. Direct checkout checks the target's evidence, never records Staff consent for them. Existing pending submissions retain their original terms binding. No live borrowing endpoint or full business enforcement has been implemented. See [report](../project/PHASE4B_REPORT.md), DEC-066/067 and [tests](../../integration/PHASE4B.md).


## Confirmed account-policy update and future API boundaries

DEC-070 confirms email plus a separate borrower-chosen password for Student/Faculty, Admin-only individual creation and Student-only Excel creation/deactivation. Student requires unique textual official Student ID and approved SKSU email; Faculty is individual-only with any valid unique accessible email and no required Student ID. All new activation/roster routes above are **unimplemented naming proposals** for an authorized Phase5 plan, not current API contracts. Existing auth/login/refresh/logout/me and Phase4B terms remain implemented; no SSO/public registration/general password-reset is implied.

The prepare/upload/preview stage must have zero account mutations; only explicitly reviewed Admin confirmation selects Students for execution. Validate configured SKSU Student domains/required textual Student ID/template identity rules under DEC-070 and OPEN-028; domain/email syntax is not ownership verification. ID/email disagreement rejects as ambiguous. Reauthorize current Admin actor and current target role AND borrower_type/identity/status on confirmation/replay, excluding Faculty/Staff/Admin; persist scoped command/audit/per-record results under the chosen transaction strategy. Missing obligation projections cannot report zero balances; DEC-073 resolves OPEN-029: warn/confirm and permit Student deactivation despite obligations, without loan/stock/fine/replacement/overdue changes or resolution writes. Never affect absent/unselected accounts or delete account/history through a roster.

A Phase5 plan must define allowed Excel input/bounds, exact required/optional fields, batch/row outcome and confirmation contracts, duplicate/stale-preview/conflict semantics, idempotency hashes/scopes, partial-failure/atomicity/retry behavior and minimal audit/file retention. Account results/ordinary errors/logs never return passwords or raw activation tokens. Email link delivery is recommended; Brevo is selected for backend-only activation and future recovery. Its API key/verified sender/link-origin configuration and successful live delivery verification remain OPEN-017; no raw provider credential is client-visible or committed. Token expiry/reissue/verification/existing-account lifecycle and general recovery remain OPEN-001/009. No invented numerical lifetimes/quotas, fake delivery success or official terms wording. [Detailed policy/integrity/acceptance handoff](../project/ACCOUNT_PROVISIONING_POLICY.md).

## Final provisioning contract boundary

DEC-070 governs these future proposals: all account creation is Admin-only, including direct calls/confirmation/replay. Staff operational assistance does not grant creation. Individual Student/Faculty borrower form and separate privileged account workflow have conditional identity rules; standard Excel has no Faculty or role selector. The exact API field serialization, file limits, string-cell handling, error codes and preview/transaction semantics require a reviewed Phase5 contract; recommended template columns are not a newly implemented HTTP schema. Preserve Phase4A authentication/Phase4B terms. Actual SKSU Student configuration/format inputs and technical matching (OPEN-028), activation/ownership/recovery (OPEN-001/009) and selected Brevo backend integration/API-key/verified-sender/tested delivery (OPEN-017) remain dependencies. Student deactivation with obligations is resolved in DEC-073. Official terms approval after FSMO presentation (DEC-072/OPEN-016) gates official publication/live borrowing, not independent account-management/inventory development. [Policy handoff](../project/ACCOUNT_PROVISIONING_POLICY.md) distinguishes technical work from institutional approvals; Phase5 remains unstarted.

## Batch1 Phase6 implemented inventory boundary

Paired000007 now implements operator-provided categories, catalog equipment metadata/lifecycle, four physical buckets, immutable movement/audit/command receipts and bounded equipment-only canonical PNG images. `T=A+R+C+D`; Staff/Admin usable acquisition/removal affects A/T only. Admin reviewed expected-sequence reconciliation sets observed A and preserves R/C/D; no custody or incident correction. Current per-pool quantity bound is2147483647; all arithmetic is checked server-side and by PostgreSQL. Loss history and replacement liability remain separate future concepts, with no invented current lost bucket or zero liability projection.

Equipment ACTIVE is Borrower-visible; INACTIVE/ARCHIVED is operational-only. Archive retains current stock/history, blocks R/C, and requires verified absence of future borrowing/replacement tables; their presence without integrated liability checks fails closed. No archival return/settlement, repair/disposal, borrowing, reservation, replacement or fine workflow exists. Metadata and stock versions are separate; locked transactions make movements/audit/receipts atomic, idempotent and stale-basis safe. Category inactivity preserves existing references. Local database image support is bounded PNG/JPEG re-encoding for catalog use only, not return photographs or arbitrary storage. Actual routes/limits/storage/permissions and future integration requirements are in [API contracts](../API_CONTRACTS.md#phase-6-implemented-equipment-and-inventory-contracts). Historical Phase2/4 candidate/unimplemented descriptions above are dated design evidence, superseded for these scoped implemented surfaces.
