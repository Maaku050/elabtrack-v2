# FSMO frontend reconstruction contract

**Phase 7 execution update, 2026-10-09:** The owner authorized complete milestones 7A–7D from baseline `366b6ef`. Borrowing/reservation/physical-checkout implementation and final verification are recorded in [the completion report](../project/PHASE7_COMPLETION_REPORT.md). Current source and the Phase 7 sections govern these surfaces; earlier unstarted/unauthorized statements are historical checkpoints. New-screen owner acceptance remains pending. Phase 8 is not started. Official FSMO publication/current acceptance, approved production Student domains and verified live activation delivery remain external gates.

Stage B, 2026-10-09. The owner explicitly approved Stage A and authorized Stage B in the current reconstruction brief. **STAGE B OWNER VISUAL ACCEPTANCE: APPROVED** in the subsequent pre-Phase 7 inventory-filter review; the final checkpoint follows the focused correction documented in [the inventory filter report](../project/INVENTORY_FILTER_REVIEW_REPORT.md). Phase 7 remains unauthorized. Product content/permissions follow the charter, confirmed decisions, API contracts and current backend. Owner screenshots guide layout quality, not demonstration features.

## Status and architecture

**IMPLEMENTED:** the approved operational frame/directories plus all existing Phase 5–6 borrower, Excel, category, equipment, stock/correction and Administration surfaces listed in [the Stage B report](../project/FRONTEND_RECONSTRUCTION_STAGE_B_REPORT.md). Screens use real existing APIs, not mock-data prototypes. Stage B reconstruction retained backend behavior, auth, API clients, query hooks, schemas, permissions and manifests. The subsequent owner-authorized filter correction changes only the inventory listing read predicate to require ACTIVE equipment and positive available stock; command/business logic remains unchanged. Final verification and limitations are recorded in the report; the later explicit owner approval supersedes the historical pending visual-acceptance state.

Preserve page → feature hook → TanStack Query → feature API → centralized transport. Query owns returned accounts, stock and images. Zustand owns only existing UI preferences/session coordination. URL search parameters own representative directory filters/page; local form state stays in React Hook Form/Zod. Sidebar presentation reads the existing account store; it never fetches or authorizes an account independently. ProtectedRoute/current-account hooks and server permission checks remain authoritative.

## Shared composition and tokens

| Concern | Convention / existing component |
| --- | --- |
| Workspace | SidebarProvider controlled by existing sidebarOpen/setSidebar; installed icon collapse and mobile Sheet; persistent pathless router parent |
| Identity/navigation | FSMO seal, navy sidebar, violet/indigo active state; SidebarHeader/Content/Group/Menu/Footer; Breadcrumb uses Router Links |
| Account menu | Base UI DropdownMenu, actual account identity, Your account and existing logout; no tenant/team switcher |
| Content | Main max-width 1600px; padding 24px desktop, 16px below 768px; page gaps 24px / 20px mobile |
| Spacing | Existing 4/8/12/16/20/24/32px tokens; sections/cards generally 20–24px; related controls 8–16px |
| Text | Headings 30px / 26px mobile; section headings 17px; body/controls 14–15px; secondary metadata 13px; small labels 12px only |
| Surface | Existing Card with scoped 12px radius, border and subtle shadow; no globally overridden Card internals |
| Actions | Existing AppButton/Button; minimum 44px target, primary action first visually prominent, descriptive action names; wrap on small screens |
| Filters | DirectorySearch + installed Input; DirectorySelect + NativeSelect, semantic associated label; search expands, other filters wrap predictably |
| Table | Installed Table semantic headers/cells; named focusable horizontal scroll region; rows minimum 76px, identity grouping, full data accessible |
| Status | Existing ToneBadge/Badge semantic labels plus icon; color never the only status signal |
| Loading/empty/error | Skeleton with busy/status semantics; Empty composition; scoped error message and Retry; never fabricate zero stock during failure |
| Statistics | DirectoryStat uses backend quantities; label aggregate scope “Across all matching equipment”; preserve A/R/C/D values |
| Focus/motion | Existing Base UI behavior; visible rings, keyboard menus/tooltips, mobile focus trap/Escape/return, inherited reduced-motion support |

Use `reconstruction.css` scoped to the operational frame/directory classes. Legacy styles remain for historical previews and unmigrated workflows; Stage B uses the same stylesheet for shared management/forms/uploader composition. Legacy styles remain where historical previews or borrower-facing cards still consume them; do not blindly delete those consumers. Do not add another override layer for each page. Login assets/styles and borrower-facing shell remain independent.

## Server table/pagination contract

TanStack Query is installed; TanStack Table is not. These straightforward server tables use the installed Table without another dependency or fake client sorting. A future advanced table needs a separate justified decision; no library upgrade is required by this contract.

Known-total directories use the actual response page/per_page/total/items, numbered first/final pages, adjacent pages and ellipsis for large ranges. Previous/Next unavailable actions have no href, aria-disabled and tabIndex=-1; event handlers suppress commands. Busy pagination cannot submit another page change. Mobile keeps Previous/Next plus visible current/total page text. Invalid URL page inputs fall back to page 1; out-of-range results show an explicit explanation and a route back to valid pages.

Anchors expose real filter-preserving URLs and link semantics; ordinary clicks use router query navigation. Modified clicks retain expected browser link behavior. Search replaces the current history entry to avoid a Back entry per character; discrete filters/page changes create history entries. Changing a filter resets the server page. Actor-scoped query keys and session invalidation stay unchanged. Never sort/slice only the currently returned 25 records and present it as the whole dataset.

Categories currently return up to 100 items with **no total**. The Inventory category picker shares previous/next presentation, displays “Total pages unavailable”, and enables Next only for a full batch. A full last batch can produce one subsequent empty page because the API has no has-more flag; retain honest unknown-total behavior. Category management, equipment form options and borrower catalog categories now use this pattern without inventing counts. Account audit and inventory movement history likewise use their actual 25-item batches with unknown totals. Staff/Admin directories use actual known totals; embedded Administration keeps staff_page/staff_search/staff_status separate from other view state. Bulk validation is bounded by its existing roster contract; it must not be confused with a server directory.

Borrowers retain only the existing search/borrower_type/status/page/per_page parameters; no unapproved sort is added. Student/Faculty type is separate from authorization role. Admin-only Add Borrower/Student Bulk links are absent for Staff; guards/backend still enforce authority. Activation-not-required does not claim email ownership or delivery; obligations UNAVAILABLE never means clear. Inventory keeps category/status/available_only/sort/name-or-available parameters, all four stock buckets and actor-scoped thumbnails/detail links.

Inventory and Borrower Catalog label `available_only` **Available for borrowing**. The server intersects ACTIVE status + available>0 with Status, Category and Search before counts/totals/paging. Explicit INACTIVE or ARCHIVED + checked returns zero matching equipment; unchecking retains the selected Status. Archived visibility otherwise stays unchanged. Inactive categories do not hide existing active equipment. Physical bucket values and backend authorization/terms enforcement are unchanged.

## Stage B form contract — implemented

Apply existing Field, FieldLabel, FieldDescription, FieldError, Input, NativeSelect/Select and Textarea to React Hook Form/Zod without changing validation or command payloads. Use a short heading/description, logical field groups, consistent 16–20px gaps, inline field feedback associated with aria-describedby and invalid semantics, and an error summary for submission failures. Desktop forms should use the available width with paired fields where useful, not narrow empty columns; single-column reading order on mobile.

A shared FormSection earns its abstraction through field grouping/labels; FormActions provides consistent alignment, cancel/back, primary submit, disabled pending state and feedback. Do not create wrappers that merely rename a div. Consequential actions retain the existing Dialog/AlertDialog or separate review step, focus behavior, command idempotency key lifecycle, stale-version errors and retry discipline. Toast may supplement persistent success/error text; it must not be the only confirmation.

Student creation retains institutional domain + string Student ID; Faculty individual creation permits its existing valid email requirements and no Student ID. Borrowers choose their own activation password. Staff/Admin administration remains separate. No new role controls, password assignment, Excel policy, terms publication or live email action is part of Stage B.

Equipment/category editing preserves metadata versions and current API commands. Stock adjustment and Admin-only correction retain explanations, projected quantities, separate review/confirmation, stale sequence protection and submission locks. Correction affects verified available only; it never overrides reserved/checked-out/damaged-held custody or resolves returns/fines/replacements. No backend arithmetic, ledger, transaction or audit change is part of the rollout.

## Stage B image uploader contract — implemented

Use the existing FileDropzone/EquipmentImage security behavior, not a new storage model. One editor section owns the existing or selected image **inside the upload area**, format/size guidance, browse control and validation feedback. Support drag-and-drop, click/keyboard browse, visible drag-over, replace selection and remove **unsaved** selection. Revoke object URLs on replacement/unmount; preserve actor/generation fencing and protected image access.

Accepted existing backend input: PNG/JPEG, at most 512 KiB, at most 2048px per dimension and 4 million pixels. Backend re-encoding, ownership checks, Origin, idempotency and expected_version remain authoritative. Inline busy/progress and validation text must describe the actual upload phase; selection is not evidence of persistence.

- Create: image section belongs inside the creation form. Preserve existing create-then-upload outcome handling; do not imply a new atomic combined endpoint. Metadata creation succeeding while upload fails must leave an honest recoverable status.
- Edit: show the existing image inside the same section, replace selection there, and keep existing image if an unsaved replacement is removed. Do not expose unsupported stored-image deletion.
- Details: display the image and link to authorized Edit; Stage B removed duplicated full editor composition after Chromium verified Edit replacement, retained current image after unsaved removal, and image-only retry following a metadata-success/upload-failure outcome. Do not delete image data/history or invent an endpoint.

This contract is for equipment **catalog** images only. Return photographs/evidence remain completely outside eLabTrack under Decision 17; no return upload/attachment workflow is introduced.

## Stage B rollout and owner gate

The owner approved Stage A in the current brief; its original report and screenshots retain their historical pending status. Stage B executes A shared forms/uploader → B borrower surfaces → C Student bulk → D categories → E Add/Edit equipment/images → F details → G stock/correction → H Administration → I consistency/responsive, with frontend and Chromium gates between groups. The report records actual results and diagnostic reruns.

ManagementPage/Card reuse DirectoryHeader and installed Card. FormSection/Field/Actions preserve RHF events, refs, validation and payloads; Notice uses Alert with semantic status/error text. WorkflowSteps indicates real Excel progress. Category add/edit uses the installed Dialog with pending dismissal locks. Consequential account/equipment status and bulk operations retain AlertDialog; stock retains a separate frozen review. No new dependency or incompatible Radix component was added.

Inventory keeps every physical bucket and Open action. The fixed table minimum is 1020px, with 28% equipment, four 12% stock columns, 12% status and 12% actions; quantities use concise padding and headers can wrap. At 1366px expanded desktop it fits the content region; narrower screens retain named keyboard-focusable scrolling rather than hiding stock. Mobile form reading order is single-column; image preview remains contained and controls wrap.

Student bulk preview always shows every returned row without fake paging; row selection labels provide a 44px target. No default activation email or automatic consent is introduced. Administration offers fixed existing Staff creation and read-only Admin details, name/email search, current policy publication status and existing account audit links. TERMS_NOT_PUBLISHED remains expected null-policy state through the unchanged hook; no institutional text is invented.

**OWNER VISUAL APPROVAL RECEIVED:** the focused inventory filter correction is complete, and owner manual verification is explicitly PASSED at the final checkpoint (DEC-076). Verification covers the seven requested widths and short heights; assistive-device/hardware checks and external institutional/delivery gates remain distinct limitations. The final owner instruction authorizes a reviewed LOCAL Git checkpoint. Stop after closure; no Phase 7 implementation, push or deployment is authorized.

Reference APIs verified against local installed source and [Base Sidebar](https://ui.shadcn.com/docs/components/base/sidebar), [Base Pagination](https://ui.shadcn.com/docs/components/base/pagination), [Base Data Table](https://ui.shadcn.com/docs/components/base/data-table), [Base Field](https://ui.shadcn.com/docs/components/base/field). Use Base UI render composition, not incompatible Radix asChild examples.

## Phase 7 implementation reuse and acceptance boundary

Owner-approved StageA/B primitives, tokens and persistent shells are reused for catalog-to-request selection/review, own history/details/cancel, Staff request directory/details/denial/physical checkout and direct issuance. Existing ManagementCard/FormField/FormActions/WorkflowSteps/Notice, DirectoryTableRegion/search, ServerPagination, badges, AlertDialog, RHF/Zod and centralized TanStack Query transport remain the architecture. No approved visual assets or UI primitives were rewritten.

New screens have loading/empty/error states, stable reviewed payload/key, duplicate-click guards, conflict revalidation and session-generation fences. Direct issuance selects an existing Borrower; dates explicitly use Manila. Catalog-prefill mounts a form initialized from loaded eligible equipment, with no state-copy effect loop. Own history is outside the new-request terms gate. Pending detail revalidates every30s and mutations invalidate borrowing/inventory/accounts without permanently caching roles.

[Actual synthetic browser evidence](verification/phase7/) covers mobile/tablet/desktop, both themes and keyboard focus. Engineering evidence does not establish owner acceptance of Phase7 screens. StageA/B acceptance remains approved; Phase7 owner acceptance is pending. No kiosk expansion or Phase8 return interface is included.
