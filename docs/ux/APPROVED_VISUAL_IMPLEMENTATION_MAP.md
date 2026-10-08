# Approved visual implementation map

2026-10-08. Every manifest PNG was opened and inspected; hashes/dimensions are verified. [Approved README](approved/README.md), [manifest](approved/MANIFEST.json), [fidelity contract](VISUAL_FIDELITY_CONTRACT.md), [Phase3B report](../project/PHASE3B_REPORT.md).

**36 files =35 screen PNGs +1 brand PNG. S01 contains two baseline screens, so36 screen targets =4 baselines +32 future screens.** Role below is the confirmed capability boundary, not the manifest's device label. Named components refer to current shared compositions or clearly deferred retained primitives, not evidence that a full screen exists.

Expected phase is sequencing/dependency guidance, never authorization. **Phase 4A update:** B06-02 sign-in/session recovery is implemented and runtime-reviewed; B04-02 has a limited safe read-only account surface. The other30 targets remain NOT IMPLEMENTED, including S06-01 whose shared feedback components do not implement its complete feature states. Visual inspection of a reference is not runtime review of a future feature. Existing59 Phase3A.1 behavior surfaces still apply where PNGs combine variants.

## Four baseline previews

### B01 — Borrower Home

- **Role:** Borrower. **Reference file:** [core/B01-borrower-home.png](approved/core/B01-borrower-home.png).
- **Layout structure:** Compact header; greeting;2×2 metrics; Browse;3 activity rows;3 quick actions; amber notice; bottom nav.
- **Shared components:** BorrowerShell, MetricCard, EquipmentThumbnail, BorrowingStatusBadge, headings.
- **Unique components/content:** Synthetic activity and local deferred-feature disclosures.
- **Expected responsive behavior:** 320/390 single column;768/1280 wider metrics/activity, same four destinations.
- **Expected implementation phase:**3B representative presentation; production Home/Dashboard metrics10, Catalog6/7/11, Pending7.
- **Current implementation status:** Development preview implemented; no functional production feature or business API.
- **Visual review status:** Chromium viewport/theme comparison completed. Owner reviewed real implementation screenshots and accepted both themes, responsiveness and disclosed minor adaptations, 2026-10-08 (DEC-064). This approves the baseline preview, not future production functionality.

### B02 — Equipment Catalog

- **Role:** Borrower. **Reference file:** [core/B02-borrower-equipment-catalog.png](approved/core/B02-borrower-equipment-catalog.png).
- **Layout structure:** Header/cart; title/search/filter/sort; chips; horizontal item rows; selected summary; bottom nav.
- **Shared components:** BorrowerShell, SearchField, EquipmentCard, QuantitySelector, AvailabilityBadge, Sheet.
- **Unique components/content:** Local query/category/availability/sort, unique selection, read-only summary.
- **Expected responsive behavior:** 320 controls below facts,390 side control;768 two columns;1280 three; no page scroll horizontally.
- **Expected implementation phase:**3B representative presentation; production Home/Dashboard metrics10, Catalog6/7/11, Pending7.
- **Current implementation status:** Development preview implemented; no functional production feature or business API.
- **Visual review status:** Chromium viewport/theme comparison completed. Owner reviewed real implementation screenshots and accepted both themes, responsiveness and disclosed minor adaptations, 2026-10-08 (DEC-064). This approves the baseline preview, not future production functionality.

### S01-upper — Staff Dashboard

- **Role:** Staff/Admin. **Reference file:** [core/S01-staff-dashboard-and-pending-requests.png](approved/core/S01-staff-dashboard-and-pending-requests.png).
- **Layout structure:** Navy sidebar; toolbar; title;4 metrics; bar-trend/donut/activity panels.
- **Shared components:** StaffShell, PageHeading, MetricCard, SurfaceCard, ThemeControl.
- **Unique components/content:** Static chart dataset, accessible text table and physical-stock legend.
- **Expected responsive behavior:** 1440 three panels;1024 two plus spanning activity; below1024 nav Sheet and stacked panels.
- **Expected implementation phase:**3B representative presentation; production Home/Dashboard metrics10, Catalog6/7/11, Pending7.
- **Current implementation status:** Development preview implemented; no functional production feature or business API.
- **Visual review status:** Chromium viewport/theme comparison completed. Owner reviewed real implementation screenshots and accepted both themes, responsiveness and disclosed minor adaptations, 2026-10-08 (DEC-064). This approves the baseline preview, not future production functionality.

### S01-lower — Pending Requests

- **Role:** Staff/Admin. **Reference file:** [core/S01-staff-dashboard-and-pending-requests.png](approved/core/S01-staff-dashboard-and-pending-requests.png).
- **Layout structure:** Same shell; heading/search/filter; lifecycle pills; identity-first table; footer paging.
- **Shared components:** StaffShell, SearchField, DataTableShell, FineStatusBadge, Table, PaginationControls.
- **Unique components/content:** Eight synthetic rows, local checkboxes/search/fine filter; read-only Review.
- **Expected responsive behavior:** 1440 full row actions;1024 named horizontal table scrolling; smaller shell Sheet.
- **Expected implementation phase:**3B representative presentation; production Home/Dashboard metrics10, Catalog6/7/11, Pending7.
- **Current implementation status:** Development preview implemented; no functional production feature or business API.
- **Visual review status:** Chromium viewport/theme comparison completed. Owner reviewed real implementation screenshots and accepted both themes, responsiveness and disclosed minor adaptations, 2026-10-08 (DEC-064). This approves the baseline preview, not future production functionality.

## 32 future approved targets

### B02-01 — Borrower - Request Cart

- **Role:** Borrower. **Reference file:** [remaining/B02-01.png](approved/remaining/B02-01.png).
- **Layout structure:** Line-item cart above totals and informational banner; bottom Continue dock.
- **Shared components:** BorrowerShell, EquipmentThumbnail, QuantitySelector, SurfaceCard; any workflow-specific composition is deferred.
- **Unique components/content:** Unique selections, remove control, type/unit totals; draft not reserved.
- **Expected responsive behavior:** Mobile item stack, tablet list plus summary, wide side summary only if list remains readable.
- **Expected implementation phase:** Phase 7; end-to-end validation11.
- **Current implementation status:** NOT IMPLEMENTED. Shared shell/component availability is not screen completion.
- **Visual review status:** Actual PNG inspected; future runtime fidelity is NOT VERIFIED. Individual content/dependency findings are in the fidelity contract's32-reference ledger.

### B02-02 — Borrower - Review and Submit

- **Role:** Borrower. **Reference file:** [remaining/B02-02.png](approved/remaining/B02-02.png).
- **Layout structure:** Date/time review card above selected-item card and visit notice; bottom Submit.
- **Shared components:** BorrowerShell, PageHeading, SurfaceCard, Input, Button; any workflow-specific composition is deferred.
- **Unique components/content:** Requested date and time, current-terms evidence, review/submission boundary.
- **Expected responsive behavior:** Single logical mobile flow; wider date/time pair and item/summary panels.
- **Expected implementation phase:** Phase 7; terms foundation4; validation11.
- **Current implementation status:** NOT IMPLEMENTED. Shared shell/component availability is not screen completion.
- **Visual review status:** Actual PNG inspected; future runtime fidelity is NOT VERIFIED. Individual content/dependency findings are in the fidelity contract's32-reference ledger.

### B02-03 — Borrower - Pending Confirmation

- **Role:** Borrower. **Reference file:** [remaining/B02-03.png](approved/remaining/B02-03.png).
- **Layout structure:** Success/status header; request/expiry facts; visit instructions; reserved items.
- **Shared components:** BorrowerShell, BorrowingStatusBadge, SurfaceCard, Button; any workflow-specific composition is deferred.
- **Unique components/content:** Committed Pending reference, absolute submission/expiry, physical visit next step.
- **Expected responsive behavior:** Stack status then quantities; wider fact pairs without hiding visit/expiry.
- **Expected implementation phase:** Phase 7.
- **Current implementation status:** NOT IMPLEMENTED. Shared shell/component availability is not screen completion.
- **Visual review status:** Actual PNG inspected; future runtime fidelity is NOT VERIFIED. Individual content/dependency findings are in the fidelity contract's32-reference ledger.

### B02-04 — Borrower - Denied Expired Cancelled

- **Role:** Borrower. **Reference file:** [remaining/B02-04.png](approved/remaining/B02-04.png).
- **Layout structure:** Terminal request card with Denied reason; adjacent Expired/Cancelled examples.
- **Shared components:** BorrowerShell, BorrowingStatusBadge, SurfaceCard; any workflow-specific composition is deferred.
- **Unique components/content:** Separate retained terminal outcomes and visible denial reason.
- **Expected responsive behavior:** One readable mobile card; wider context groups; no historical reactivation.
- **Expected implementation phase:** Phase 7.
- **Current implementation status:** NOT IMPLEMENTED. Shared shell/component availability is not screen completion.
- **Visual review status:** Actual PNG inspected; future runtime fidelity is NOT VERIFIED. Individual content/dependency findings are in the fidelity contract's32-reference ledger.

### S02-01 — Staff - Request Review

- **Role:** Staff/Admin. **Reference file:** [remaining/S02-01.png](approved/remaining/S02-01.png).
- **Layout structure:** Identity and request facts across top; held-equipment table; handover banner/actions.
- **Shared components:** StaffShell, SurfaceCard, Table, Badge, Dialog; any workflow-specific composition is deferred.
- **Unique components/content:** Bound terms, current target/hold, due confirmation; approve WITH physical release.
- **Expected responsive behavior:** Two desktop panels; stack early on tablet; readable mobile review and full-height confirmation.
- **Expected implementation phase:** Phase 7.
- **Current implementation status:** NOT IMPLEMENTED. Shared shell/component availability is not screen completion.
- **Visual review status:** Actual PNG inspected; future runtime fidelity is NOT VERIFIED. Individual content/dependency findings are in the fidelity contract's32-reference ledger.

### S02-02 — Staff - Direct Checkout

- **Role:** Staff/Admin. **Reference file:** [remaining/S02-02.png](approved/remaining/S02-02.png).
- **Layout structure:** Four-stage direct-checkout stepper; borrower search table left, selected borrower right.
- **Shared components:** StaffShell, SearchField, DataTableShell, SurfaceCard; any workflow-specific composition is deferred.
- **Unique components/content:** Active borrower selection, current terms, equipment/due/review steps; immediate issue.
- **Expected responsive behavior:** Desktop main/summary; tablet condensed summary; mobile one step at a time.
- **Expected implementation phase:** Phase 7.
- **Current implementation status:** NOT IMPLEMENTED. Shared shell/component availability is not screen completion.
- **Visual review status:** Actual PNG inspected; future runtime fidelity is NOT VERIFIED. Individual content/dependency findings are in the fidelity contract's32-reference ledger.

### B03-01 — Borrower - Active Borrowing

- **Role:** Borrower. **Reference file:** [remaining/B03-01.png](approved/remaining/B03-01.png).
- **Layout structure:** Active identity, due/progress card, issued/physical quantities, in-person return notice.
- **Shared components:** BorrowerShell, BorrowingStatusBadge, SurfaceCard; any workflow-specific composition is deferred.
- **Unique components/content:** Read-only original due, remaining units and Staff return instructions.
- **Expected responsive behavior:** Stack quantities and timing; pair panels wide, no borrower-return action.
- **Expected implementation phase:** Phase 7 detail;8 accountability.
- **Current implementation status:** NOT IMPLEMENTED. Shared shell/component availability is not screen completion.
- **Visual review status:** Actual PNG inspected; future runtime fidelity is NOT VERIFIED. Individual content/dependency findings are in the fidelity contract's32-reference ledger.

### B03-02 — Borrower - Partial and Replacement

- **Role:** Borrower. **Reference file:** [remaining/B03-02.png](approved/remaining/B03-02.png).
- **Layout structure:** Active/partial header; incident quantities; replacement required/accepted/remaining cards.
- **Shared components:** BorrowerShell, AccountabilityBadge, SurfaceCard; any workflow-specific composition is deferred.
- **Unique components/content:** Physical0/replacement3 remains open; original incident and deadline retained.
- **Expected responsive behavior:** Mobile labeled counts, tablet paired accountability, wide events plus summary.
- **Expected implementation phase:** Phase 8.
- **Current implementation status:** NOT IMPLEMENTED. Shared shell/component availability is not screen completion.
- **Visual review status:** Actual PNG inspected; future runtime fidelity is NOT VERIFIED. Individual content/dependency findings are in the fidelity contract's32-reference ledger.

### B03-03 — Borrower - Overdue Fine

- **Role:** Borrower. **Reference file:** [remaining/B03-03.png](approved/remaining/B03-03.png).
- **Layout structure:** Overdue notice; original due; live PHP fine card; unresolved replacements.
- **Shared components:** BorrowerShell, FineStatusBadge, SurfaceCard, Alert; any workflow-specific composition is deferred.
- **Unique components/content:** PHP10 per started24h, live clock until P0/U0; no payment gateway.
- **Expected responsive behavior:** Money/timing before item detail; widen fact pairs without replacing explanatory copy.
- **Expected implementation phase:** Phase 8.
- **Current implementation status:** NOT IMPLEMENTED. Shared shell/component availability is not screen completion.
- **Visual review status:** Actual PNG inspected; future runtime fidelity is NOT VERIFIED. Individual content/dependency findings are in the fidelity contract's32-reference ledger.

### S03-01 — Staff - Record Return

- **Role:** Staff/Admin. **Reference file:** [remaining/S03-01.png](approved/remaining/S03-01.png).
- **Layout structure:** Borrowing/remaining facts left; good/damage/loss entry right; stock-effect review below.
- **Shared components:** StaffShell, SurfaceCard, Input, Table, ConflictAlert; any workflow-specific composition is deferred.
- **Unique components/content:** Prior vs now quantities, required incident notes, unique return lines and P/U prediction.
- **Expected responsive behavior:** Desktop split panels; tablet review below; mobile item input cards, no squeezed input table.
- **Expected implementation phase:** Phase 8.
- **Current implementation status:** NOT IMPLEMENTED. Shared shell/component availability is not screen completion.
- **Visual review status:** Actual PNG inspected; future runtime fidelity is NOT VERIFIED. Individual content/dependency findings are in the fidelity contract's32-reference ledger.

### S03-02 — Staff - Return Result

- **Role:** Staff/Admin. **Reference file:** [remaining/S03-02.png](approved/remaining/S03-02.png).
- **Layout structure:** Three disposition metrics; immutable movement/result table; resulting state panel.
- **Shared components:** StaffShell, MetricCard, DataTableShell, BorrowingStatusBadge; any workflow-specific composition is deferred.
- **Unique components/content:** Server-confirmed return result, good/damage/loss effects, automatic closure criterion.
- **Expected responsive behavior:** Metrics wrap; scoped effect table or labeled cards; result stays before optional next action.
- **Expected implementation phase:** Phase 8.
- **Current implementation status:** NOT IMPLEMENTED. Shared shell/component availability is not screen completion.
- **Visual review status:** Actual PNG inspected; future runtime fidelity is NOT VERIFIED. Individual content/dependency findings are in the fidelity contract's32-reference ledger.

### S03-03 — Staff - Accept Replacement

- **Role:** Staff/Admin. **Reference file:** [remaining/S03-03.png](approved/remaining/S03-03.png).
- **Layout structure:** Replacement obligation context left; accepted quantity/type right; before/after table.
- **Shared components:** StaffShell, SurfaceCard, QuantitySelector, Table, Dialog; any workflow-specific composition is deferred.
- **Unique components/content:** Appropriate equivalent acceptance, new-stock acquisition, preserved original incident.
- **Expected responsive behavior:** Two panels wide, stacked narrow; named scroll only for genuine before/after comparison.
- **Expected implementation phase:** Phase 8.
- **Current implementation status:** NOT IMPLEMENTED. Shared shell/component availability is not screen completion.
- **Visual review status:** Actual PNG inspected; future runtime fidelity is NOT VERIFIED. Individual content/dependency findings are in the fidelity contract's32-reference ledger.

### B04-01 — Borrower - My Borrowings History

- **Role:** Borrower. **Reference file:** [remaining/B04-01.png](approved/remaining/B04-01.png).
- **Layout structure:** My Borrowings filter pills, open transaction cards, terminal history cards.
- **Shared components:** BorrowerShell, BorrowingStatusBadge, SurfaceCard, PaginationControls; any workflow-specific composition is deferred.
- **Unique components/content:** Own bounded Pending/Active/history views with distinct terminal meanings.
- **Expected responsive behavior:** Mobile single chronology; wider rows/optional table retain state/due context.
- **Expected implementation phase:** Phase 7; accountability8.
- **Current implementation status:** NOT IMPLEMENTED. Shared shell/component availability is not screen completion.
- **Visual review status:** Actual PNG inspected; future runtime fidelity is NOT VERIFIED. Individual content/dependency findings are in the fidelity contract's32-reference ledger.

### B04-02 — Borrower - Account

- **Role:** Borrower. **Reference file:** [remaining/B04-02.png](approved/remaining/B04-02.png).
- **Layout structure:** Profile identity card, accepted-terms/preferences card, supported account actions.
- **Shared components:** BorrowerShell, SurfaceCard, ThemeControl, Input; any workflow-specific composition is deferred.
- **Unique components/content:** Own account display, supported name edit, terms access and existing sign-out.
- **Expected responsive behavior:** Readable single mobile column; paired display groups with one logical edit form.
- **Expected implementation phase:** Phase 4 account/access;5 profile.
- **Current implementation status:** PARTIAL, Phase 4A: safe real read-only name/email/role and shared theme/sign-out. No category/contact/terms acceptance or profile editor is fabricated.
- **Visual review status:** PNG opened and limited account capture inspected; full profile/preferences/terms composition remains deferred to its authorized feature. See Phase 4A deviations.

### B04-03 — Borrower - Terms Acceptance

- **Role:** Borrower. **Reference file:** [remaining/B04-03.png](approved/remaining/B04-03.png).
- **Layout structure:** Terms version card with readable rules and acceptance control/action.
- **Shared components:** BorrowerShell, SurfaceCard, Checkbox, Button; any workflow-specific composition is deferred.
- **Unique components/content:** Once-per-current-version consent; existing obligations remain accessible.
- **Expected responsive behavior:** Readable constrained column at all widths; acceptance reachable without nested scrolling.
- **Expected implementation phase:** Phase 4.
- **Current implementation status:** NOT IMPLEMENTED. Shared shell/component availability is not screen completion.
- **Visual review status:** Actual PNG inspected; future runtime fidelity is NOT VERIFIED. Individual content/dependency findings are in the fidelity contract's32-reference ledger.

### S04-01 — Staff - Borrower Directory

- **Role:** Staff/Admin; privileged actions Admin only. **Reference file:** [remaining/S04-01.png](approved/remaining/S04-01.png).
- **Layout structure:** Borrower directory toolbar, Add/Bulk actions, identity/accountability table.
- **Shared components:** StaffShell, SearchField, DataTableShell, Badge, PaginationControls; any workflow-specific composition is deferred.
- **Unique components/content:** Borrower-category/status search, Staff single provision; Admin-only bulk/status actions.
- **Expected responsive behavior:** Essential columns at tablet; mobile cards; privileged actions retain explicit scope.
- **Expected implementation phase:** Phase 5.
- **Current implementation status:** NOT IMPLEMENTED. Shared shell/component availability is not screen completion.
- **Visual review status:** Actual PNG inspected; future runtime fidelity is NOT VERIFIED. Individual content/dependency findings are in the fidelity contract's32-reference ledger.

### S04-02 — Staff - Borrower Detail

- **Role:** Staff/Admin; deactivation/clearance Admin only. **Reference file:** [remaining/S04-02.png](approved/remaining/S04-02.png).
- **Layout structure:** Borrower identity/context cards then fine/accountability and borrowing-history table.
- **Shared components:** StaffShell, SurfaceCard, FineStatusBadge, AccountabilityBadge, Table; any workflow-specific composition is deferred.
- **Unique components/content:** Operational identity/history; inactive loans remain resolvable; privileged contextual actions.
- **Expected responsive behavior:** Panels stack and histories scope-scroll; mobile preserves record identity and current obligations.
- **Expected implementation phase:** Phase 5; fine clearance8.
- **Current implementation status:** NOT IMPLEMENTED. Shared shell/component availability is not screen completion.
- **Visual review status:** Actual PNG inspected; future runtime fidelity is NOT VERIFIED. Individual content/dependency findings are in the fidelity contract's32-reference ledger.

### S04-03 — Admin - Bulk Import

- **Role:** Admin only. **Reference file:** [remaining/S04-03.png](approved/remaining/S04-03.png).
- **Layout structure:** Four-stage bulk-import path; template/file left, validation summary/table right.
- **Shared components:** StaffShell, SurfaceCard, Input, Table, AlertDialog; any workflow-specific composition is deferred.
- **Unique components/content:** Row validity/duplicates, explicit valid subset, durable import result; no password spreadsheet.
- **Expected responsive behavior:** Step-by-step mobile; compact tablet preview and named full-row scroll; desktop split.
- **Expected implementation phase:** Phase 5.
- **Current implementation status:** NOT IMPLEMENTED. Shared shell/component availability is not screen completion.
- **Visual review status:** Actual PNG inspected; future runtime fidelity is NOT VERIFIED. Individual content/dependency findings are in the fidelity contract's32-reference ledger.

### S05-01 — Staff - Inventory Overview

- **Role:** Staff/Admin. **Reference file:** [remaining/S05-01.png](approved/remaining/S05-01.png).
- **Layout structure:** Inventory heading/search; four physical counts; equipment stock table.
- **Shared components:** StaffShell, MetricCard, DataTableShell, SearchField, Badge; any workflow-specific composition is deferred.
- **Unique components/content:** A/R/C/damaged-held reconciled stock, catalog lifecycle distinct from availability.
- **Expected responsive behavior:** Dense desktop counts/table; tablet essential columns; mobile pool cards and count disclosure.
- **Expected implementation phase:** Phase 6.
- **Current implementation status:** NOT IMPLEMENTED. Shared shell/component availability is not screen completion.
- **Visual review status:** Actual PNG inspected; future runtime fidelity is NOT VERIFIED. Individual content/dependency findings are in the fidelity contract's32-reference ledger.

### S05-02 — Staff - Equipment Detail and Movements

- **Role:** Staff/Admin. **Reference file:** [remaining/S05-02.png](approved/remaining/S05-02.png).
- **Layout structure:** Equipment identity/metadata; physical reconciliation; separate liability banner and ledger.
- **Shared components:** StaffShell, SurfaceCard, MetricCard, DataTableShell; any workflow-specific composition is deferred.
- **Unique components/content:** T=A+R+C+D, immutable movements, metadata vs typed stock operations.
- **Expected responsive behavior:** Count groups wrap, ledger named scroll; edit form stacks and does not hide custody.
- **Expected implementation phase:** Phase 6.
- **Current implementation status:** NOT IMPLEMENTED. Shared shell/component availability is not screen completion.
- **Visual review status:** Actual PNG inspected; future runtime fidelity is NOT VERIFIED. Individual content/dependency findings are in the fidelity contract's32-reference ledger.

### S05-03 — Admin - Full Fine Clearance

- **Role:** Admin only. **Reference file:** [remaining/S05-03.png](approved/remaining/S05-03.png).
- **Layout structure:** Fine basis/history left; noneditable full outstanding/method/note review right.
- **Shared components:** StaffShell, SurfaceCard, FineStatusBadge, Select, Textarea, AlertDialog; any workflow-specific composition is deferred.
- **Unique components/content:** Full PAID/WAIVED/OTHER clearance, preserved basis and active-accrual warning.
- **Expected responsive behavior:** Two desktop panels; narrow full-height readable confirmation with fixed amount.
- **Expected implementation phase:** Phase 8.
- **Current implementation status:** NOT IMPLEMENTED. Shared shell/component availability is not screen completion.
- **Visual review status:** Actual PNG inspected; future runtime fidelity is NOT VERIFIED. Individual content/dependency findings are in the fidelity contract's32-reference ledger.

### S05-04 — Admin - Reports

- **Role:** Admin only. **Reference file:** [remaining/S05-04.png](approved/remaining/S05-04.png).
- **Layout structure:** Report tabs/filter toolbar; result metrics; bounded physical/operational table/export.
- **Shared components:** StaffShell, MetricCard, DataTableShell, SearchField, Tabs; any workflow-specific composition is deferred.
- **Unique components/content:** Source-defined reports, filters and separately approved export/print format.
- **Expected responsive behavior:** Compact tablet filters; named comparison table scroll; mobile summaries and filter Sheet.
- **Expected implementation phase:** Phase 10.
- **Current implementation status:** NOT IMPLEMENTED. Shared shell/component availability is not screen completion.
- **Visual review status:** Actual PNG inspected; future runtime fidelity is NOT VERIFIED. Individual content/dependency findings are in the fidelity contract's32-reference ledger.

### S05-05 — Admin - Administration

- **Role:** Admin only. **Reference file:** [remaining/S05-05.png](approved/remaining/S05-05.png).
- **Layout structure:** Administration cards for accounts/terms/audit; named account list and current policy.
- **Shared components:** StaffShell, SurfaceCard, DataTableShell, Tabs; any workflow-specific composition is deferred.
- **Unique components/content:** Named identities, versioned terms, restricted audit; initial rate/TTL read-only.
- **Expected responsive behavior:** Three cards wide then stack; bounded account table; readable single terms column.
- **Expected implementation phase:** Phase 4 roles;5 accounts;10 audit/reporting.
- **Current implementation status:** NOT IMPLEMENTED. Shared shell/component availability is not screen completion.
- **Visual review status:** Actual PNG inspected; future runtime fidelity is NOT VERIFIED. Individual content/dependency findings are in the fidelity contract's32-reference ledger.

### S05-06 — Staff - Borrowings List

- **Role:** Staff/Admin. **Reference file:** [remaining/S05-06.png](approved/remaining/S05-06.png).
- **Layout structure:** Borrowings list tabs, physical/replacement/fine context table and explanatory banner.
- **Shared components:** StaffShell, DataTableShell, BorrowingStatusBadge, AccountabilityBadge; any workflow-specific composition is deferred.
- **Unique components/content:** Canonical lifecycle plus derived Overdue/Replacements; no enum conflation.
- **Expected responsive behavior:** Essential identity/due/P/U columns tablet; mobile record cards and stacked detail.
- **Expected implementation phase:** Phase 7;8 accountability.
- **Current implementation status:** NOT IMPLEMENTED. Shared shell/component availability is not screen completion.
- **Visual review status:** Actual PNG inspected; future runtime fidelity is NOT VERIFIED. Individual content/dependency findings are in the fidelity contract's32-reference ledger.

### B06-01 — Borrower - Activity

- **Role:** Borrower. **Reference file:** [remaining/B06-01.png](approved/remaining/B06-01.png).
- **Layout structure:** Today/earlier own event-feed sections with icon/status/time rows.
- **Shared components:** BorrowerShell, SurfaceCard, BorrowingStatusBadge; any workflow-specific composition is deferred.
- **Unique components/content:** Owner-scoped immutable activity projection; no delivery/read receipts implied.
- **Expected responsive behavior:** Single readable feed at all widths, no oversized wide grid.
- **Expected implementation phase:** Phase 9.
- **Current implementation status:** NOT IMPLEMENTED. Shared shell/component availability is not screen completion.
- **Visual review status:** Actual PNG inspected; future runtime fidelity is NOT VERIFIED. Individual content/dependency findings are in the fidelity contract's32-reference ledger.

### S06-01 — Shared - Empty Loading and Error States

- **Role:** Shared. **Reference file:** [remaining/S06-01.png](approved/remaining/S06-01.png).
- **Layout structure:** Three comparative empty/loading/conflict panels; error and unavailable-state panels below.
- **Shared components:** StaffShell where authenticated, EmptyState, LoadingRegion, ErrorState, ConflictAlert; any workflow-specific composition is deferred.
- **Unique components/content:** Context-specific empty/busy/read failure/409 recovery; safe request ID.
- **Expected responsive behavior:** Stack examples below tablet; feature states preserve their actual content geometry.
- **Expected implementation phase:** Phase 3B shared foundation only; concrete states each feature.
- **Current implementation status:** NOT IMPLEMENTED. Shared shell/component availability is not screen completion.
- **Visual review status:** Actual PNG inspected; future runtime fidelity is NOT VERIFIED. Individual content/dependency findings are in the fidelity contract's32-reference ledger.

### B02-05 — Borrower - Equipment Detail

- **Role:** Borrower. **Reference file:** [remaining/B02-05.png](approved/remaining/B02-05.png).
- **Layout structure:** Catalog image/detail, item facts, quantity and unreserved Add dock.
- **Shared components:** BorrowerShell, EquipmentThumbnail, AvailabilityBadge, QuantitySelector, SurfaceCard; any workflow-specific composition is deferred.
- **Unique components/content:** Equipment description and selection; catalog picture only.
- **Expected responsive behavior:** Image then facts mobile; paired image/facts wide; action clears navigation.
- **Expected implementation phase:** Phase 6 catalog;7 selection;11 validation.
- **Current implementation status:** NOT IMPLEMENTED. Shared shell/component availability is not screen completion.
- **Visual review status:** Actual PNG inspected; future runtime fidelity is NOT VERIFIED. Individual content/dependency findings are in the fidelity contract's32-reference ledger.

### S02-03 — Staff - Deny Request Confirmation

- **Role:** Staff/Admin. **Reference file:** [remaining/S02-03.png](approved/remaining/S02-03.png).
- **Layout structure:** Request identity/held items left; required denial reason and release warning right.
- **Shared components:** StaffShell, SurfaceCard, Textarea, AlertDialog; any workflow-specific composition is deferred.
- **Unique components/content:** Borrower-visible nonblank reason; confirmed denial releases hold and preserves history.
- **Expected responsive behavior:** Short desktop confirmation; readable mobile Sheet with safe cancellation.
- **Expected implementation phase:** Phase 7.
- **Current implementation status:** NOT IMPLEMENTED. Shared shell/component availability is not screen completion.
- **Visual review status:** Actual PNG inspected; future runtime fidelity is NOT VERIFIED. Individual content/dependency findings are in the fidelity contract's32-reference ledger.

### B04-04 — Borrower - Completed Detail

- **Role:** Borrower. **Reference file:** [remaining/B04-04.png](approved/remaining/B04-04.png).
- **Layout structure:** Completed reference, physical/replacement zero, frozen final fine and history.
- **Shared components:** BorrowerShell, BorrowingStatusBadge, FineStatusBadge, SurfaceCard; any workflow-specific composition is deferred.
- **Unique components/content:** Completed can retain money outstanding; clearance does not erase assessment.
- **Expected responsive behavior:** Status/final basis first mobile; paired facts and chronological events wide.
- **Expected implementation phase:** Phase 8.
- **Current implementation status:** NOT IMPLEMENTED. Shared shell/component availability is not screen completion.
- **Visual review status:** Actual PNG inspected; future runtime fidelity is NOT VERIFIED. Individual content/dependency findings are in the fidelity contract's32-reference ledger.

### S04-04 — Admin - Staff and Admin Accounts

- **Role:** Admin only. **Reference file:** [remaining/S04-04.png](approved/remaining/S04-04.png).
- **Layout structure:** Named Staff/Admin directory beside create/edit privileged account form.
- **Shared components:** StaffShell, DataTableShell, Input, Select, Dialog; any workflow-specific composition is deferred.
- **Unique components/content:** Role/status review, secure provisioning and last-Admin recovery dependency.
- **Expected responsive behavior:** Directory/form split wide; stacked tablet; mobile named account cards/form.
- **Expected implementation phase:** Phase 4 role controls;5 account management.
- **Current implementation status:** NOT IMPLEMENTED. Shared shell/component availability is not screen completion.
- **Visual review status:** Actual PNG inspected; future runtime fidelity is NOT VERIFIED. Individual content/dependency findings are in the fidelity contract's32-reference ledger.

### S05-07 — Admin - Exceptional Stock Correction

- **Role:** Admin only. **Reference file:** [remaining/S05-07.png](approved/remaining/S05-07.png).
- **Layout structure:** Stock basis and guarded exceptional-correction review side by side.
- **Shared components:** StaffShell, SurfaceCard, Input, Textarea, AlertDialog; any workflow-specific composition is deferred.
- **Unique components/content:** Typed direction/quantity/reason and immutable audit; procedure still requires approved contract.
- **Expected responsive behavior:** Split wide, stack narrow; review never hides holds/custody or original evidence.
- **Expected implementation phase:** Phase 6 only after correction policy resolved.
- **Current implementation status:** NOT IMPLEMENTED. Shared shell/component availability is not screen completion.
- **Visual review status:** Actual PNG inspected; future runtime fidelity is NOT VERIFIED. Individual content/dependency findings are in the fidelity contract's32-reference ledger.

### B06-02 — Shared - Sign in and Session Recovery

- **Role:** Shared access; no authenticated borrower nav. **Reference file:** [remaining/B06-02.png](approved/remaining/B06-02.png).
- **Layout structure:** Compact brand then provisioned-account sign-in form and session-recovery notice.
- **Shared components:** AppBrand, Input, Button, Alert; LoginPage, SessionBootstrap and shared recovery/forbidden compositions.
- **Unique components/content:** Existing secure session restoration/expired handling; no signup or invented activation state.
- **Expected responsive behavior:** Centered readable form at all widths; remove protected navigation from unauthenticated state.
- **Expected implementation phase:** Phase 4A.
- **Current implementation status:** IMPLEMENTED in Phase 4A: real provisioned-account login and existing secure restoration/logout/recovery.
- **Visual review status:** Approved PNG opened; new 320/390/1280 light/dark captures inspected against shared tokens/composition. Deviations and limits are recorded in the Phase 4A report; no new owner-approval claim.

## Separate brand manifest record

**BRAND — FSMO seal; all audiences.** [Reference](approved/brand/FSMO-seal-reference.png); standalone circular image used by AppBrand, compact borrower header and Staff identity. No screen or workflow of its own. Responsive sizes44px Borrower /56px Staff,44px collapsed rail. Phase3B asset integrated with replaceable `sealSrc`. Supplied AI reconstruction visually inspected; **not authenticated official master**. [Provenance](../../frontend/src/assets/brand/README.md). This is the36th manifest file and is not counted as an additional screen.
