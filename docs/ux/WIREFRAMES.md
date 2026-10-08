# Low-fidelity wireframes — Phase 3A.1

**Current provisioning policy overlay, 2026-10-08:** [DEC-070 / approved rules and remaining gates](../project/ACCOUNT_PROVISIONING_POLICY.md) governs future behavior: only Admin creates accounts; Student requires unique textual official Student ID/approved SKSU email; Faculty is individual-only with valid unique accessible email and no required Student ID. Standard bulk creation/deactivation is Student-only. Both use borrower-owned separate passwords and current officially published Phase4B terms. This documentation overlay changes no implementation or approved mockup asset; original phase checkpoints remain historical.

**Text layouts only, 2026-10-08;38 numbered wireframe groups covering all 59 inventoried UX surfaces.** Multi-step/variant groups share layouts and are not38 production routes. Mobile borrower base and desktop/tablet operations are shown schematically; widths of plain text drawings are not CSS measurements. This package contains no high-fidelity colors, branded imagery, fonts, polished screens, React/HTML components or runnable prototype.

[Inventory/navigation](UX_ARCHITECTURE.md), [Borrower journeys](BORROWER_FLOWS.md), [operational journeys](STAFF_ADMIN_FLOWS.md), [responsive matrix](RESPONSIVE_RULES.md), [status system](STATUS_SYSTEM.md) and [future checklist](UX_ACCEPTANCE_CHECKLIST.md) define behavior beyond each drawing. Generic references/counts are synthetic examples; no real borrower identity/contact data. Square brackets represent controls/status/placeholder text, never photographs for returns. Catalog image placeholders are independently in scope.

## Shared interaction and state conventions

Borrower drawings start with a narrow single-column flow; bottom labels use short schematic text where necessary (final label My Borrowings). Sticky actions never overlap navigation, safe area, focus/errors or keyboard. Operational drawings inherit one-level sidebar/app bar in WF-15; do not repeat the sidebar in every detail. Variant screenshots are text states, not new backend lifecycle enums. Buttons visible only under current role/ownership; server enforces all commands.

Every group inherits Skeleton for region loading, truthful Empty with next action, Alert+Retry for reads, associated field errors and persistent command result; WF-37 provides examples. Domain-sensitive success/reservation/expiry/completion/fine clearance appears only after server-confirmed outcome. Unknown write checks original result before reissue; staff conflicts refresh quantities and require review. Toast is supplementary. All controls labeled, keyboard/focus accessible; semantic statuses independent of color/light-dark; no return-media controls for any role.

## Review index

| Group | Layout | Surface IDs | Base |
|---|---|---|---|
| WF-01 | Sign in | SH-01 | Borrower mobile |
| WF-02 | First-use and updated terms | B-01 | Borrower mobile |
| WF-03 | Borrower home — mobile base | B-02 | Borrower mobile |
| WF-04 | Catalog and search/filter sheet | B-03, B-04 | Borrower mobile |
| WF-05 | Equipment detail | B-05 | Borrower mobile |
| WF-06 | Selected equipment / cart | B-06 | Borrower mobile |
| WF-07 | Requested date/time and request review | B-07, B-08 | Borrower mobile |
| WF-08 | Submitted and pending request | B-09, B-10 | Borrower mobile |
| WF-09 | Denied, expired and cancelled detail variants | B-11, B-12, B-13 | Borrower mobile |
| WF-10 | Active borrowing and overdue fine | B-14, B-17 | Borrower mobile |
| WF-11 | Partial, replacement and completed status | B-15, B-16, B-18 | Borrower mobile |
| WF-12 | My Borrowings and history — mobile lists | B-19, B-20 | Borrower mobile |
| WF-13 | Profile / account | B-22 | Borrower mobile |
| WF-14 | Notifications / activity — scoped list | B-21, O-23 | Borrower mobile |
| WF-15 | Operational dashboard — desktop | O-01 | Operations desktop/tablet |
| WF-16 | Pending request queue | O-02 | Operations desktop/tablet |
| WF-17 | Request review and physical release | O-03 | Operations desktop/tablet |
| WF-18 | Reasoned denial | O-04 | Operations desktop/tablet |
| WF-19 | Direct checkout — four steps | O-07, O-08, O-09, O-10 | Operations desktop/tablet |
| WF-20 | Return processing — quantities now | O-11 | Operations desktop/tablet |
| WF-21 | Return review and partial/final result | O-12 | Operations desktop/tablet |
| WF-22 | Replacement selection, acceptance and result | O-13, O-14 | Operations desktop/tablet |
| WF-23 | Borrower list | O-15 | Operations desktop/tablet |
| WF-24 | Borrower detail | O-16 | Operations desktop/tablet |
| WF-25 | Admin individual Student/Faculty creation | O-17 | Admin-only desktop/tablet |
| WF-26 | Inventory list — physical counts | O-18 | Operations desktop/tablet |
| WF-27 | Inventory detail and movement history | O-19 | Operations desktop/tablet |
| WF-28 | Equipment create/edit — catalog metadata | O-20 | Operations desktop/tablet |
| WF-29 | Stock movement and inactive/archive impact | O-21, O-22 | Operations desktop/tablet |
| WF-30 | Admin full fine clearance and history | A-01, A-02 | Operations desktop/tablet |
| WF-31 | Admin Student bulk import — template, preview and result | A-03, A-04, A-05 | Operations desktop/tablet |
| WF-32 | Admin account deactivation/reactivation | A-06 | Operations desktop/tablet |
| WF-33 | Admin named Staff/Admin accounts | A-07, A-08 | Operations desktop/tablet |
| WF-34 | Admin current policy and versioned terms | A-09, A-10 | Operations desktop/tablet |
| WF-35 | Admin reports with filter/export concept | A-11 | Operations desktop/tablet |
| WF-36 | Admin audit | A-12 | Operations desktop/tablet |
| WF-37 | Access, loading, empty and error examples | SH-02 | Shared |
| WF-38 | Operational borrowings list and detail | O-05, O-06 | Operations desktop/tablet |

## WF-01 — Sign in

Surfaces: SH-01. Likely future shadcn primitives: **Card, Field, Input, Button, Alert**.

```text
+--------------------------------------+
| eLabTrack                     Theme  |
| Sign in to your provisioned account  |
| Email                                |
| [__________________________________] |
| Password                             |
| [________________________] [Show]    |
| [ Sign in ]                          |
| Need an account? Contact FSMO.       |
| Field error / service error region   |
+--------------------------------------+
```

No signup or role picker. Show password has an accessible name/state; invalid login safe inline message. SH-02 handles current-account/session denial. Desktop remains one readable column.

## WF-02 — First-use and updated terms

Surfaces: B-01. Likely future shadcn primitives: **Card, Field, Checkbox, Button, Alert**.

```text
+--------------------------------------+
| < Back          Terms of borrowing   |
| Version [current]  Effective [date]   |
| Review before creating a request     |
| ------------------------------------ |
| Readable approved terms content      |
| Full content available by scrolling  |
| No truncated acceptance-only text    |
| ------------------------------------ |
| [ ] I accept this terms version      |
| [ Accept terms ]                     |
| [ View existing borrowings ]         |
| Failure / Retry region               |
+--------------------------------------+
```

Acceptance once/version; label names the version. Existing loans remain reachable if updated terms not accepted. No forced per-request checkbox. Busy/error does not fabricate acceptance; approved copy is a later content dependency.

## WF-03 — Borrower home — mobile base

Surfaces: B-02. Likely future shadcn primitives: **Card, Badge, Button, Item, Skeleton, Alert**.

```text
+--------------------------------------+
| Home           Notifications [link] |
| Next action                          |
| Pending BR-004                       |
| Visit FSMO before [expiry, Manila]   |
| [ View request ]                     |
| ------------------------------------ |
| Active borrowing: due [date + time]  |
| Physical remaining 3                 |
| Replacements remaining 0  [Open]      |
| ------------------------------------ |
| [ Browse equipment ]                 |
| Recent permitted updates [View all] |
| ------------------------------------ |
| Home | Equipment | My loans | Account|
+--------------------------------------+
```

Short navigation text in drawings is schematic: final third label is My Borrowings. Next action remains above summaries; do not total overlapping status counts. Empty active region invites browse; read failures are not0 active.

## WF-04 — Catalog and search/filter sheet

Surfaces: B-03, B-04. Likely future shadcn primitives: **Input, Button, Card, Badge, Sheet, Select, Checkbox, Skeleton, Empty**.

```text
+--------------------------------------+
| Equipment              Selected [2] |
| Search equipment                     |
| [________________________] [Search]  |
| [ Filter ] [Available] [Clear]       |
| ------------------------------------ |
| [catalog image / placeholder]        |
| Spoons                               |
| Category [if confirmed]              |
| Available5   Short description       |
| [View detail] [Add to request]       |
| ------------------------------------ |
| Next equipment card;1-column mobile |
| [Load more]                          |
| Home | Equipment | My loans | Account|
+--------------------------------------+
Filter Sheet (same catalog):
+--------------------------------------+
| Filters                      [Close] |
| Category [All v]                      |
| [ ] Available equipment only         |
| [ Clear filters ] [ Apply filters ]  |
+--------------------------------------+
```

Quantity is selected in detail/add sheet, not automatic1 without review. Filters/search have labeled input and no-results Clear filters; no internal R/C/D buckets. Catalog image failure preserves name/availability. Cart count is selected items, label counts explicitly.

## WF-05 — Equipment detail

Surfaces: B-05. Likely future shadcn primitives: **Card, AspectRatio, Field, Input, Button, Badge, Alert**.

```text
+--------------------------------------+
| < Equipment                Selected |
| [ Equipment catalog image ]          |
| Spoons   [category if confirmed]      |
| Description / permitted use           |
| Available5   Checked [time]          |
| Quantity                             |
| [ - ] [ 2 ] [ + ]                    |
| Not reserved until request submitted |
| ------------------------------------ |
| [ Add2 to request ]                  |
| Home | Equipment | My loans | Account|
+--------------------------------------+
```

Add merges same equipment into one draft line; quantity<=shown availability is advisory, server rechecks on submit. Out of stock disables Add with text. Large touch controls, typed number alternative, sticky action clears bottom navigation.

## WF-06 — Selected equipment / cart

Surfaces: B-06. Likely future shadcn primitives: **Card, Field, Input, Button, Empty, Alert**.

```text
+--------------------------------------+
| < Equipment       Selected equipment |
| Spoons                               |
| Quantity [ - ] [ 2 ] [ + ] [Remove]   |
| Available now5; availability can vary|
| ------------------------------------ |
| Plates                               |
| Quantity [ - ] [ 1 ] [ + ] [Remove]   |
| ------------------------------------ |
| 2 equipment types /3 units           |
| Selection is not a reservation       |
| [ Add more equipment ]               |
| [ Continue to request review ]       |
| Empty: [Browse equipment]             |
+--------------------------------------+
```

No purchase subtotal or payment CTA. Remove is local with optional undo, not a destructive borrowing cancellation. Stock conflict annotates affected line and requires deliberate edit/review; preserve unaffected intent.

## WF-07 — Requested date/time and request review

Surfaces: B-07, B-08. Likely future shadcn primitives: **Field, Input, Calendar, Popover/Sheet, Card, Button, Alert**.

```text
+--------------------------------------+
| < Selected equipment   Review request|
| Spoons2 / Plates1        [Edit]      |
| Requested return date                 |
| [YYYY-MM-DD] [Calendar]               |
| Requested return time [HH:MM]         |
| Timezone: Asia/Manila                 |
| Staff confirms final due at release  |
| ------------------------------------ |
| Terms version [accepted version]     |
| Updated? [Review current terms]       |
| Submit reserves equipment.           |
| Visit FSMO for approval and handover.|
| Pending reservation expires in 24h.   |
| Field error summary                  |
| [ Submit request ]                   |
+--------------------------------------+
Date picker: Calendar + typed-date alternative;
mobile full-height sheet, labeled time input.
```

Due date AND time required; no seven-day maximum or hidden default. Already accepted current terms shown as evidence, no per-loan checkbox. Read/write failure leaves draft; unknown submit reconciles original result before another command.

## WF-08 — Submitted and pending request

Surfaces: B-09, B-10. Likely future shadcn primitives: **Card, Badge, Alert, Button, AlertDialog, Item**.

```text
+--------------------------------------+
| < My Borrowings    BR-004 [Pending]   |
| Your equipment is reserved.          |
| Visit FSMO for approval and physical |
| release before [expiry date + time]. |
| This request is not yet approved.    |
| ------------------------------------ |
| Submitted [date/time, Asia/Manila]    |
| Expires   [date/time, Asia/Manila]    |
| Advisory time left [hours/minutes]   |
| Requested due [date/time, Manila]    |
| Reserved: Spoons2 / Plates1          |
| [ Browse equipment ]                 |
| [ Cancel request ]                   |
+--------------------------------------+
Cancel confirmation:
| Release this request's reservation? |
| History remains.                     |
| [Keep request] [Cancel request]      |
```

B-09 success context and B-10 retained detail share layout. Cancel own pending only. Timer0 shows checking status until server response; concurrent approval/expiry shows actual state, never local release. Success summary remains accessible beyond a toast.

## WF-09 — Denied, expired and cancelled detail variants

Surfaces: B-11, B-12, B-13. Likely future shadcn primitives: **Card, Badge, Alert, Item, Button**.

```text
+--------------------------------------+
| < History                 BR-004     |
| [Denied] / [Expired] / [Cancelled]   |
|                                      |
| Denied: borrower-visible reason      |
| [Reason text]  Decided [date/time]    |
| OR Expired: reservation time ended.  |
| OR Cancelled: request was cancelled. |
| ------------------------------------ |
| Reservation released; history kept.  |
| Requested items and original times  |
| [ Start a new request ]              |
| [ View history ]                     |
+--------------------------------------+
```

Three variants of one canonical detail, not three new statuses beyond domain. Reason only for Denied; no blame/denial label for expiry. New request creates fresh draft/checks, never reopens this row. No Delete or overdue fine for these terminal requests.

## WF-10 — Active borrowing and overdue fine

Surfaces: B-14, B-17. Likely future shadcn primitives: **Card, Badge, Item, Alert, Accordion, Button**.

```text
+--------------------------------------+
| < My Borrowings    BR-005 [Active]    |
| [Overdue if due crossed]              |
| Due [date + time, Asia/Manila]        |
| Physical remaining 3                 |
| Replacements remaining 0             |
| ------------------------------------ |
| Spoons issued5 / good returned2      |
| Recorded damaged0 / lost0           |
| Still physically outstanding3       |
| [View return / replacement history] |
| ------------------------------------ |
| Current overdue fine PHP 10           |
| As of [time] [Outstanding]           |
| Can grow until obligations resolved. |
| Bring equipment to FSMO; staff       |
| records the authoritative return.    |
+--------------------------------------+
```

Read-only borrower quantities. No Mark Returned, photo, payment or complete control. B-17 overdue variant shares active layout and requires server unresolved/due snapshot. Fine unavailable on read error is not PHP 0; no replacement-valued price fine.

## WF-11 — Partial, replacement and completed status

Surfaces: B-15, B-16, B-18. Likely future shadcn primitives: **Card, Badge, Item, Accordion, Alert**.

```text
+--------------------------------------+
| BR-006 [Active] [Replacement required]|
| [Overdue if original due crossed]    |
| Spoons: issued5 / good2 / lost3       |
| Physical remaining 0                 |
| Replacement required3 (Lost)        |
| Accepted1 / Remaining2              |
| Bring2 appropriate replacements     |
| to FSMO for staff acceptance.         |
| ------------------------------------ |
| Original due [date + time]            |
| Live fine / clear history [View]     |
| Event history (read-only)            |
| Lost3 recorded [time]                |
| Replacement1 accepted [time]        |
+--------------------------------------+
Completion variant (after accepted remaining 2):
| [Completed] Physical0 / Replacement0 |
| Completed [time]  Final fine PHP 20    |
| [Fine outstanding] or [Fine cleared] |
| Original loss3 remains in history.  |
```

B-15 partial-good variant shows physical remaining 3/replacement0; B-16 replacement-only remains Active with C0; B-18 completed uses final fine. Fine20 illustration assumes completion25h late. Staff acceptance is authoritative; no borrower editing or original-disposition workflow.

## WF-12 — My Borrowings and history — mobile lists

Surfaces: B-19, B-20. Likely future shadcn primitives: **Tabs, Card, Badge, Select, Button, Skeleton, Empty**.

```text
+--------------------------------------+
| My Borrowings                        |
| [Pending] [Active] [History]          |
| Filter history [All v]               |
| ------------------------------------ |
| BR-006 Active / Replacement required |
| Spoons: physical0 / replacement2    |
| Due [date/time]         [Open detail]|
| ------------------------------------ |
| BR-004 Denied                        |
| Reason summary [Open detail]         |
| ------------------------------------ |
| BR-003 Completed / Fine outstanding |
| Completed [date/time] [Open detail]  |
| [Load more]                          |
| Home | Equipment | My loans | Account|
+--------------------------------------+
```

Shows variant examples, not a single tab containing incompatible live filters. History retains completed/denied/cancelled/expired. Empty pending/active/history each has own accurate copy/Browse action. Desktop can align rows; mobile is not a shrunken data table.

## WF-13 — Profile / account

Surfaces: B-22. Likely future shadcn primitives: **Card, Field, Input, Button, Switch, Separator**.

```text
+--------------------------------------+
| Account                              |
| Name [own registered name]           |
| Email [own account, read-only]        |
| Borrower category [Student/Faculty]  |
| Program / contact [if approved]      |
| Status Active (account)              |
| [Edit supported display name]        |
| ------------------------------------ |
| Terms [accepted version] [View]      |
| [Notifications / activity]           |
| Theme [Light / Dark]                 |
| [Sign out]                           |
| Account changes? Contact FSMO.        |
+--------------------------------------+
```

Role/status/email/category not self-edit controls. Category does not grant permissions. Sign out follows existing session behavior; no new password delivery/reset mechanics or photos/evidence. Theme means semantic content readable in either mode, no final colors chosen.

## WF-14 — Notifications / activity — scoped list

Surfaces: B-21, O-23. Likely future shadcn primitives: **Item, Card, Badge, Button, Skeleton, Empty, Alert**.

```text
+--------------------------------------+
| < Back       Notifications / activity|
| Latest updates                       |
| ------------------------------------ |
| Request pending BR-004 [time]         |
| Visit FSMO before [expiry] [Open]     |
| ------------------------------------ |
| Replacement accepted BR-006 [time]   |
| 2 replacements remain       [Open]   |
| ------------------------------------ |
| Completed BR-003 [time] [Open]       |
| [Load more]                          |
| Empty: No updates for this view.      |
| Error: Updates unavailable [Retry]    |
+--------------------------------------+
```

Borrower own events; Staff/Admin permitted operational activity uses same readable structure. No mark-read/unread badge, media, SMTP settings or fabricated delivery proof. Future owner-scoped projection/API is a dependency; canonical borrowing detail remains truth even if notices fail.

## WF-15 — Operational dashboard — desktop

Surfaces: O-01. Likely future shadcn primitives: **Sidebar, Card, Badge, Button, Table/Item, Skeleton, Empty**.

```text
+-------------------+----------------------------------------------+
| Dashboard         | Dashboard        Notifications  Account      |
| Requests & loans  | [Direct checkout]                            |
| Inventory         | Pending12 [Open] | Active8 [Open]             |
| Borrowers         | Due today3       | Overdue2 [Open]            |
| Admin: Reports    | Replacements4 [Open obligations]             |
| Administration    | -------------------------------------------- |
|                   | Next pending requests: borrower / expiry      |
|                   | BR-004 [Borrower] [expiry] [Open review]       |
|                   | Open obligations / relevant operational list |
+-------------------+----------------------------------------------+
```

Illustrative synthetic counts, overlapping measures not summed. Admin links absent for Staff; Staff action queues are not report permission. Pending empty state “No requests awaiting review”; filtered empty clears filter. Mobile cards retain actionable labels.

## WF-16 — Pending request queue

Surfaces: O-02. Likely future shadcn primitives: **Tabs, Input, Button, Table, Badge, Pagination, Empty**.

```text
+---------------------------------------------------------------+
| Requests & Borrowings [Pending][Active][Overdue][Replacements] |
| Search borrower/reference [__________]  [Expiry filter]        |
| Borrower | Ref | Submitted | Expires | Types/units | Account   |
| [name]   |004  |[time]     |[time]   |2 /3        |Active      |
| Fine outstanding PHP 20 / replacement context  [Open request]  |
| ------------------------------------------------------------- |
| Pagination / results count                                   |
| Empty queue: No pending requests.  [Direct checkout]           |
+---------------------------------------------------------------+
```

No Approve row action; Open request reveals verification context. Expires absolute timestamp plus advisory remaining, not a stock-release timer authority. Card fallback shows borrower/account/expiry/fine and Open. Server-filtered bounded list, no bulk approval.

## WF-17 — Request review and physical release

Surfaces: O-03. Likely future shadcn primitives: **Card, Badge, Table, Field, Input, AlertDialog, Button, Alert**.

```text
+------------------------------------+--------------------------+
| BR-004 Pending                     | Accountability context   |
| Borrower [name] / [category]        | Existing fine PHP 20      |
| Account Active                     | Replacements [remaining] |
| Submitted / expires [date/time]    | [View borrower history]  |
| Bound accepted terms [version/time]| Human review; not a block|
| ---------------------------------- |                          |
| Equipment | Requested | Held       | Final due date [date]    |
| Spoons    |2          |2            | Final due time [HH:MM]   |
| Plates    |1          |1            | Asia/Manila              |
| Current pool status [usable]       | [Approve & Release]      |
| Requested due [date/time]          | [Deny]                   |
+------------------------------------+--------------------------+
Confirm: all listed equipment to this borrower, final due [instant].
Confirm only while physically handing over now.
[Keep reviewing] [Confirm release]
```

One issue edge; no Approved state/secondary checkout. Current account/hold/due/expiry server checks. Old terms bound to pending acceptable; fine context allows staff discretion. If quantities need change, current full-request assumption forbids silent edit. Expired/concurrent decision swaps action panel for actual status.

## WF-18 — Reasoned denial

Surfaces: O-04. Likely future shadcn primitives: **Dialog/Sheet, Field, Textarea, Alert, Button**.

```text
+----------------------------------------------------+
| Deny request BR-004                         [Close] |
| Borrower will see this reason.                      |
| Reason (required)                                  |
| [_______________________________________________]  |
| [_______________________________________________]  |
| Required-field message / server-state error         |
| Reserved equipment released; request history stays. |
| Review: borrower, quantities, visible reason        |
| [Keep reviewing] [Confirm denial]                   |
+----------------------------------------------------+
```

Required trimmed reason; no internal-only reason substituting borrower-visible content. Confirmation can be the same form after review, not a chain of redundant modals. Concurrent expiry/decision cannot be overridden; successful denial recorded before released copy appears.

## WF-19 — Direct checkout — four steps

Surfaces: O-07, O-08, O-09, O-10. Likely future shadcn primitives: **Combobox/Command, Input, Field, Table/Card, Button, AlertDialog, Alert**.

```text
+--------------------------------------------------------------+
| Direct checkout   1 Borrower >2 Equipment >3 Due >4 Review     |
| STEP1: Search [____] Active Borrower results                  |
| [Choose Borrower] Category/account/terms/fine context         |
| Missing current terms: Borrower must accept; staff cannot.    |
| STEP2: Equipment search, available, selected unique qty       |
| Spoons available5    Selected2  [-][2][+]                    |
| STEP3: Due date [YYYY-MM-DD] Time [HH:MM] Asia/Manila          |
| STEP4: Borrower / equipment / quantities / final due          |
| Creates active borrowing immediately, no pending stage.      |
| [Back] [Next] or final [Issue equipment]                      |
+--------------------------------------------------------------+
Final confirmation: physically handing equipment over now.
[Keep reviewing] [Confirm issue]
```

One step at a time on mobile; progress labels clickable only for visited steps. Existing fine displayed, not automatic prohibit. Borrower active/current terms server revalidated; direct stock conflict requires review. Draft progress not reservation; no expiry shown after direct issue.

## WF-20 — Return processing — quantities now

Surfaces: O-11. Likely future shadcn primitives: **Table/Card, Field, Input, Textarea, Button, Alert**.

```text
+---------------------------------------------------------------------+
| Record return BR-005 Active  Borrower [name]  Due [instant]          |
| [Fill all physically outstanding as good]                           |
| Item  |Issued|Prior good/damage/loss|Physical left|Replacement left  |
| Spoons|5     |0 /0 /0              |5            |0                 |
| ------------------------------------------------------------------- |
| Good returned NOW [2]   Damaged NOW [0]   Lost NOW [3]                |
| Process now5 <= physical left5; after P0 / new replacement3         |
| Incident rationale (required for damage/loss in current design)      |
| [Lost units confirmed___________________________________________]   |
| Repeat one unique row per borrowing item                            |
| [Cancel draft] [Review return]                                      |
+---------------------------------------------------------------------+
Mobile: same labels in one item card; each NOW input readable/full width.
```

Now counts default0; prior counts read-only. Full-good shortcut sets only P, explicitly reviews overwrite of existing draft, cannot resolve U. Actor Staff/Admin only, including inactive borrower loan. No photograph/camera/attachment/evidence input; photo viewing outside software.

## WF-21 — Return review and partial/final result

Surfaces: O-12. Likely future shadcn primitives: **Card, Alert, Badge, Table/Item, Button**.

```text
+------------------------------------------------------------+
| Review return BR-005                                       |
| Spoons: good2 / damaged0 / lost3 NOW                       |
| Before physical5 / replacement0                           |
| After physical0 / replacement3                            |
| Stock effect: available+2 / checked out-5 / total tracked-3 |
| Borrowing remains Active; original due unchanged.           |
| Lost3 incident remains;3 replacements required.             |
| Fine continues if overdue.                                 |
| [Back to quantities] [Record return]                       |
+------------------------------------------------------------+
After committed result:
| Return recorded [time] / quantities / remaining P0 U3       |
| [Open borrowing] [Record replacement]                       |
Final-good alternative with P0/U0:
| Completed [time]; final fine frozen; history retained.      |
```

Prediction labeled review until commit; never auto-complete from processed original counts if U remains. Field/server failure keeps draft unrecorded. Unknown result check original command; no toast-only critical result. No generic Complete action.

## WF-22 — Replacement selection, acceptance and result

Surfaces: O-13, O-14. Likely future shadcn primitives: **Card, Table, Field, Input, Textarea, Button, AlertDialog, Alert**.

```text
+------------------------------------------------------------+
| Record Replacement BR-005                                  |
| Obligation Spoons / Lost                                   |
| Required3 / Already accepted0 / Remaining3                 |
| Accept NOW [3]                                             |
| Type: appropriate Spoons or staff-accepted equivalent       |
| Equivalent rationale [if applicable____________________]    |
| [Review acceptance]                                        |
| ---------------------------------------------------------- |
| Review:3 new usable units received now                     |
| Available2 ->5  Total tracked2 ->5                        |
| Physical remains0 / Replacements3 ->0                     |
| Completes only if all borrowing items P0/U0.               |
| [Back] [Accept replacement]                                 |
+------------------------------------------------------------+
Result: Accepted3; remaining 0; Completed [time]; final fine.
Damage-original case: damaged_held unchanged; new unit adds total.
```

Unique obligation selection, no accepted>remaining or automatic clipping on race. Staff certifies appropriate/equivalent for original pool. Acquired new units, not original good-return or implicit disposal. Final event/stock/fine freeze atomic; incident remains historical after acceptance.

## WF-23 — Borrower list

Surfaces: O-15. Likely future shadcn primitives: **Input, Select, Table/Card, Badge, Button, Pagination, Empty**.

```text
+------------------------------------------------------------+
| Borrowers    [Admin:Create] [Admin:Student bulk]        |
| Search [____] Category [All v] Account status [All v]        |
| Borrower | Category | Program | Account | Current obligations|
| [name]   | Faculty  | [if set]| Active  | P3/U0/Fine20      |
| [Open borrower]                                             |
| Pagination                                                 |
| Empty: No Borrowers [Admin:Create Borrower]                     |
+------------------------------------------------------------+
```

Operational directory includes Student/Faculty; only the bulk workflow is Student-only. Staff assists without account creation; Admin has all creation/status/privileged grants. Duplicate/status/error states distinguish empty collection. Private identity fields minimal; no sample contact/biography.

## WF-24 — Borrower detail

Surfaces: O-16. Likely future shadcn primitives: **Card, Tabs, Table/Item, Badge, Button**.

```text
+-----------------------------------+--------------------------+
| Borrower [name] [category]         | Account Active           |
| Email/program/contact where needed| Admin: [Deactivate]      |
| [Current][History][Accountability] | Staff: no status button  |
| Current BR-005: P0 / U3           |                          |
| [Open loan][Record replacement]   | Fine outstanding PHP 20   |
| --------------------------------- | Admin: [Clear Fine]       |
| History retained including denied | Staff: read-only fine    |
| / cancelled / expired / completed |                          |
+-----------------------------------+--------------------------+
```

Inactive account retains history/operations; no issue to inactive target. Admin controls visible only when current permission granted. Operational account history not unrestricted audit export; borrower’s own view has no staff controls.

## WF-25 — Provision individual Borrower

Surfaces: O-17, ADMIN only under DEC-070. Likely future shadcn primitives: **Card, Field, Input, Select, Button, Alert**.

```text
+------------------------------------------------------------+
| Provision Borrower                                         |
| Name [________________] Email [________________]           |
| Borrower type [STUDENT / FACULTY]               |
| Program [if applicable________] Contact [if needed________] |
| Role: BORROWER (fixed; privileged accounts separate) |
| Secure activation: Borrower chooses own password.       |
| Student ID [required for STUDENT; string, keep zeroes]  |
| [Cancel draft] [Review and create Borrower]                 |
| Result / validation / duplicate / uncertain-create region   |
+------------------------------------------------------------+
```

Exact activation/initial-password instructions depend on approved security design; no fake sent invitation or display of passwords. Actor Admin only, never Staff or public signup. Student requires unique textual official Student ID and approved SKSU email; Faculty may use any valid unique accessible email and does not require Student ID. No Admin-chosen borrower password; privileged STAFF/ADMIN creation stays separate. Mobile single-column fields despite wide schematic drawing.

## WF-26 — Inventory list — physical counts

Surfaces: O-18. Likely future shadcn primitives: **Input, Select, Table/Card, Badge, Button, Pagination, Skeleton, Empty**.

```text
+---------------------------------------------------------------------+
| Inventory [Add equipment] Search [____] Category/status filters      |
| Equipment|Available|Reserved|Checked out|Damaged held|Total tracked   |
| Spoons   |2        |0       |0          |0           |2               |
| ------------------------------------------------------------------- |
| Separate: Historical lost3 / Replacements remaining 3                |
| [Open equipment]                                                    |
| Counts reconcile T=A+R+C+D; incident/liability not extra buckets.    |
| Pagination / no equipment / filtered empty / read error             |
+---------------------------------------------------------------------+
```

Pool example after 5 issued/good2/loss3. Borrower catalog does not see these buckets. Physical counts labeled model exactly; total includes held-out units. Mobile card prioritizes available/state, expands the exact count labels.

## WF-27 — Inventory detail and movement history

Surfaces: O-19. Likely future shadcn primitives: **Card, Tabs, Table, Badge, Button, Pagination**.

```text
+------------------------------------------------------------+
| Spoons   Active   [Edit metadata] [Record stock movement]   |
| Available2 Reserved0 Checked out0 Damaged held0 Total2    |
| Separate obligations: Lost awaiting replacement3          |
| [Overview] [Movements]                                      |
| Overview: historical loss3 / unresolved3 / related BR-005 |
| Movements: Time | Operation | Quantity delta | Reason | Actor|
| [time] | Mixed return | available+2,checked out-5,total-3 | good2/lost3 | [name]|
| Related borrowing [BR-005]  [Load more]                      |
| [Set inactive / archive review]                             |
+------------------------------------------------------------+
```

Count vector and obligation/history are separate. Ledger secondary, bounded; human reference primary, not UUID. No repair/disposal workflow guessed; damaged originals nonusable pending later policy. Archived rows/history preserved; active-liability archive disallowed.

## WF-28 — Equipment create/edit — catalog metadata

Surfaces: O-20. Likely future shadcn primitives: **Field, Input, Textarea, Select, Card, Button, Alert**.

```text
+------------------------------------------------------------+
| Add equipment / Edit metadata                              |
| Name [________] Description [___________________________]   |
| Category [approved category v]                              |
| Catalog image [choose catalog image / keep current]         |
| Catalog-image purpose only, no return media                 |
| Create only: initial usable stock [quantity]                |
| Existing pool: change stock via typed movement, not this form|
| Another edit? Review the latest changes before saving.   |
| [Cancel draft] [Review and save]                            |
+------------------------------------------------------------+
```

Catalog images independently authorized as candidates; approved file/storage/security details later. Opening usable stock produces explicit baseline/acquisition, not metadata stock overwrite. No R/C fields, outdated version conflicts and requires review rather than silent overwrite.

## WF-29 — Stock movement and inactive/archive impact

Surfaces: O-21, O-22. Likely future shadcn primitives: **Dialog/Sheet, Select, Field, Input, Textarea, AlertDialog, Button**.

```text
+------------------------------------------------------------+
| Record stock movement                                      |
| Type [Add usable / Remove available] Quantity [q]           |
| Reason [required_______________________________________]   |
| Before/after A/R/C/D/T exact vector                         |
| R/C cannot be changed through this adjustment.             |
| [Back] [Confirm movement]                                  |
| Admin only: [Exceptional correction] |
+------------------------------------------------------------+
Archive / inactive review (separate action within this group):
| Archive: requires R0/C0/replacements0; keep history.        |
| Blocked? show quantities and [Open relevant borrowing].     |
| Inactive: no new issue; old returns/replacements allowed.   |
| [Keep current state] [Confirm permitted change]             |
```

Stock removal/archive/correction consequential preview; ordinary add/removal exact permitted vector. Exceptional correction Admin only, no hiding custody, final procedure later. No damaged-original disposition controls created from unresolved policy.

## WF-30 — Admin full fine clearance and history

Surfaces: A-01, A-02. Likely future shadcn primitives: **Dialog/Sheet, Card, Field, Select, Textarea, Button, AlertDialog, Alert**.

```text
+------------------------------------------------------------+
| Clear Fine BR-007                          Admin only      |
| Borrower [name] Borrowing Completed [time]                  |
| Final assessed PHP 20 / Previously cleared PHP 0             |
| Entire outstanding PHP 20  (not editable)                    |
| Method [Choose: Paid / Waived / Other resolution]           |
| Note/reason [optional__________________________________]   |
| [Keep outstanding] [Confirm full clearance PHP 20]           |
| ---------------------------------------------------------- |
| History: Final assessed20 stays visible after clearance.   |
| [time] Cleared20 / Paid / Named Admin [name]               |
+------------------------------------------------------------+
Live-loan variant: Current fine/as-of, still P or U outstanding.
Warning: Fine can increase later until loan completed.
Staff/Borrower variant: read-only basis/balance/history; no Clear button.
```

No partial/editable amount or payment gateway. Clear full locked current balance; stale expected balance reloads/requires confirmation. Successful balance0 does not erase assessed20. PAID distinct from waived/other; private note not automatically borrower-visible.

## WF-31 — Admin Student bulk import — template, preview and result

Surfaces: A-03, A-04, A-05. Likely future shadcn primitives: **Card, Field, Input(file), Table/Card, Badge, Progress, Button, AlertDialog, Alert**.

```text
+------------------------------------------------------------------+
| Student Excel import   1 File >2 Validate/preview >3 Result        |
| STEP1 [Download approved template]                               |
| [Choose Student roster] [Validate complete file]                        |
| Columns: studentId/name/email/course/contactNumber        |
| STEP2 Valid18 / Invalid2 / Duplicate-email conflicts [details]    |
| Row | Student ID | Name | Email | Status / conflict error         |
| [Correct file and revalidate]                                    |
| [Import valid18 only] (explicit subset; others not created)       |
| Confirmation:18 STUDENT BORROWER accounts; no other roles           |
| STEP3 Created18 / Not created2 / Unresolved outcomes [Check]      |
| [View Borrowers] [Correct failed rows]                            |
+------------------------------------------------------------------+
```

Illustrative counts, no actual data import. Admin-only Student creation fixes BORROWER/STUDENT; require unique textual Student ID/SKSU email and preserve leading zeroes. Exclude password/admin/staff/faculty roles and category assignment; reject conflicting Faculty/privileged identities. Complete validation preview precedes confirmation. No fake activation delivery claim. Final file format/batch transaction/duplicate resolution still implementation design. Unknown batch result checks original identity; no blind upload/import again. This file input is account data, never return evidence.

Separate Admin-only **Bulk Deactivate Students** reuses the roster: match stable Student ID/email → preview unmatched/conflicting/already-inactive/obligation warnings → explicitly confirm selected Students → audited results. Recheck current BORROWER/STUDENT under locks; never target Faculty/Staff/Admin, hard-delete history or affect roster-absent accounts. DEC-073 resolves OPEN-029: outstanding fines/active or overdue loans/unreturned equipment/replacements cannot veto Student deactivation; warnings/confirmation preserve every obligation, borrowing state, due/overdue calculation and history record without payment/clearance/return/replacement/closure writes; no new surface count or approved mockup is introduced.

## WF-32 — Admin account deactivation/reactivation

Surfaces: A-06. Likely future shadcn primitives: **AlertDialog/Sheet, Card, Button, Alert**.

```text
+------------------------------------------------------------+
| Deactivate Borrower [name]?                                 |
| Cannot log in or make new requests/normal borrower actions. |
| Existing borrowing/history/replacements/fines retained.     |
| Staff/Admin can continue resolving existing obligations.    |
| Current open P/U and fine context [read-only]               |
| [Keep active] [Confirm deactivation]                         |
+------------------------------------------------------------+
Reactivation variant: restore normal account access, keep history
and existing obligations. [Keep inactive] [Confirm reactivation]
```

Admin only; DEC-073 permits Student individual/bulk deactivation for graduation/withdrawal/transfer/other authorized reasons despite all obligations. Show warnings and require confirmation; positive obligations never disable that action. No automatic fine/payment/return/replacement/closure/overdue/history changes; FSMO resolution remains separate and fine clearance Admin-only/auditable. No account hard delete or independent eligibility suspension. Actual current account state checked at command. On failure remain unchanged; preserve role/status history. No client workaround for inactive login.

## WF-33 — Admin named Staff/Admin accounts

Surfaces: A-07, A-08. Likely future shadcn primitives: **Table/Card, Field, Input, Select, Button, AlertDialog, Alert**.

```text
+------------------------------------------------------------+
| Administration > Accounts       [Create named account]      |
| Name | Role (Staff/Admin) | Account state | [Open]           |
| ---------------------------------------------------------- |
| Create/edit: Name [___] Email [___] Role [Staff/Admin v]     |
| Status [active/inactive]  Impact / privileged review        |
| Access setup [approved secure onboarding dependency]       |
| Review the effect on current Staff/Admin access.   |
| [Cancel] [Review privileged change]                         |
+------------------------------------------------------------+
```

No shared credential, global Super Admin or Staff access to role management. Last-admin/recovery procedure not guessed here; final implementation must review it. Destructive-impact changes retain actor/history and current server authorization.

## WF-34 — Admin current policy and versioned terms

Surfaces: A-09, A-10. Likely future shadcn primitives: **Card, Table/Item, Field, Textarea, Button, AlertDialog, Alert**.

```text
+------------------------------------------------------------+
| Administration > Terms & Policy                            |
| Current policy (read-only initial baseline):                |
| Pending expiry24h / Fine PHP 10 per started24-hour overdue period      |
| Business timezone Asia/Manila / No fixed maximum duration   |
| Current policy values are read-only in this design.|
| ---------------------------------------------------------- |
| Terms current version [version/effective date] [Read]       |
| [Draft next terms version]                                 |
| Approved text [readable editor/review area]                 |
| Publish review: new requests require new acceptance.        |
| Existing bound loan/acceptance history remains unchanged.   |
| [Keep draft] [Confirm publish new version]                  |
+------------------------------------------------------------+
```

No generic policy engine, per-department settings or SMTP provider controls. Approved content required before publishing; final version/effective-date input design later. Read-only defaults are not sliders silently changing PHP 10/TTL. Terms editor is not photo/attachment capture.

## WF-35 — Admin reports with filter/export concept

Surfaces: A-11. Likely future shadcn primitives: **Tabs/Select, Field, Input, Calendar, Popover/Sheet, Table, Pagination, Button, Empty**.

```text
+------------------------------------------------------------------+
| Reports  [Inventory / Active / Overdue / History / Replacement / Fines]|
| Dates [from][to] Asia/Manila   Category/status/method [filters]    |
| [Apply] [Clear] [Export]  |
| Summary of selected approved measures (not campus totals)        |
| Ref | Borrower | Due/completed | Physical/U | Assessment/balance  |
| ... bounded rows ...                                             |
| Fine measures: assessed / recorded Paid / Waived / Other / owed   |
| Pagination / No results [Clear filters] / Error [Retry]            |
+------------------------------------------------------------------+
```

Admin administrative access; staff normal queues separate. Export entry is design dependency, not working file endpoint/job. No assessed-as-revenue, automatic damage-price total or all-campus report. Mobile summary+scoped table scroll with key identity, larger-screen recommendation for dense analysis.

## WF-36 — Admin audit

Surfaces: A-12. Likely future shadcn primitives: **Field, Input, Select, Table/Item, Dialog/Sheet, Pagination, Alert**.

```text
+------------------------------------------------------------+
| Administration > Audit                                     |
| Date range / actor / action / borrowing-equipment reference |
| [Apply] [Clear]                                            |
| Time | Named actor/system | Action | Human entity/reference  |
| [Open permitted details: meaningful before/after/reason]    |
| Bounded pagination                                         |
| No results / unavailable / not authorized state             |
+------------------------------------------------------------+
```

Restricted named Admin access, not borrower/staff event feed. Safe context/no credentials/provider raw messages/private bodies. Legal retention policy still later; this wireframe grants no delete/export privilege or photographic evidence retention.

## WF-37 — Access, loading, empty and error examples

Surfaces: SH-02. Likely future shadcn primitives: **Skeleton, Empty, Alert, Card, Button, Dialog/Sheet**.

```text
+--------------------------------------+
| Page title / retained permitted nav   |
| [Skeleton heading][Skeleton card]    |
| [Skeleton rows] (loading, not0 data) |
+--------------------------------------+
Empty: No results for these filters. [Clear filters]
Network: Unable to load this section. [Retry]
Write unknown: Checking whether this was recorded. [Check status]
Conflict: Quantities changed. [Refresh and review]
Session expired: Sign in again to continue. [Sign in]
Inactive account: Contact FSMO about account access. [Back to sign in]
Forbidden action: You cannot perform this action. [Back]
Unexpected error: Couldn't finish. [Support details: request ID]
```

Global region states supplement per-flow notes; not all errors share login. Ordinary403 does not invalidate all peers. Every empty state in architecture/checklist has role-appropriate action. No raw API/SQL/provider detail or private state after session change.

## WF-38 — Operational borrowings list and detail

Surfaces: O-05, O-06. Likely future shadcn primitives: **Tabs, Input, Select, Table/Card, Badge, Button, Pagination, Accordion**.

```text
+------------------------------------------------------------+
| Requests & Borrowings [Active][Overdue][Replacements][History]|
| Search borrower/reference [___] Filters [status/date]       |
| Ref | Borrower | Due | Physical left | Replacement left     |
| 005 | [name]   |[t]  |0              |3   [Open borrowing]  |
| Pagination                                                 |
| ---------------------------------------------------------- |
| Detail BR-005: Active, original due, P0/U3, live fine         |
| Item quantities / current obligations / immutable events   |
| [Record return: physical units only] [Record Replacement]   |
| Staff: fine read-only / Admin: full Clear Fine if outstanding|
+------------------------------------------------------------+
```

Return with P0 explains “No physical units outstanding” and replacement action remains. Completed detail has history/final fine, no return/complete shortcuts. Inactive target still resolvable; new-issue capability separately checked. Server scoped filters, no reports permission implied.

## Review boundary

These38 groups cover home/catalog/detail/cart/due-review/submitted-pending/active-partial-replacement/history/profile and dashboard/queue/review/direct/return/replacement/borrower/inventory/fine/bulk plus shared/Admin secondary screens. Journeys A–H have a corresponding reviewable sequence. Future implementation must compare with approved Phase 3A.2 mockups and prove every state/permission/domain invariant; no runtime acceptance is claimed now. High-fidelity palette/type/component styling belongs to the next separately authorized phase.
