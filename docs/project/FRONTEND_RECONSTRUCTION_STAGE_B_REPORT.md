# Frontend reconstruction — Stage B

Final checkpoint status, 2026-10-09: **engineering complete; Stage A owner approved; Stage B owner visually approved; Inventory filter owner manual verification PASSED** (explicit current owner reply “Verified and passed”). Final reconstruction visual/functional acceptance is confirmed. [Checkpoint audit and Phase7 handoff](FRONTEND_RECONSTRUCTION_CHECKPOINT.md) record the local-only checkpoint scope and remaining dependencies; Phase7 is not started or authorized. Historical pending statements and test counts below retain their original context.


Current owner-review update: **overall Stage B visuals approved**; the subsequent focused [Inventory Filter Review](INVENTORY_FILTER_REVIEW_REPORT.md) records the last correction and final checkpoint readiness. Original verification and pending-approval statements below are historical.


Owner explicitly approved Stage A and authorized Stage B in the current brief. Phase 7 remains unauthorized. **OWNER VISUAL ACCEPTANCE: PENDING** for Stage B. No commits, pushes or deployment.

## Pre-change migration inventory

Recorded before application changes. Current source, the Stage A audit/contract, approved visual package and router establish this inventory. API clients, hooks, schemas, auth architecture and backend behavior remain the functional baseline.

| Surface / route | Existing component | Planned migration / preserved behavior |
| --- | --- | --- |
| Student / Faculty creation `/staff/borrowers/new` | AccountCreate | Shared Field sections; conditional Student ID/domain rules; no Admin-assigned password |
| Borrower details/profile/status `/staff/borrowers/:id` | AccountDetail / ProfileForm | Identity/contact/status, compact truthful accountability, existing edit and confirm actions, audit paging |
| Student Excel `/admin/borrowers/bulk` | StudentBulk | Shared uploader, full validation Table, step/results/confirmation; preserve every row/selection/parser/idempotency |
| Categories `/staff/inventory/categories` | CategoryManagement / CategoryForm | Shared list/Table, unknown-total pagination, Field form; retain version/activation relationships |
| Equipment creation `/staff/inventory/new` | EquipmentForm / CatalogImagePicker | Shared metadata/opening sections; image preview inside uploader; separate create/upload retries |
| Equipment editing `/staff/inventory/:id/edit` | EquipmentForm / CatalogImagePicker | Same editor; existing image inside uploader, replacement and unsaved removal; no stored deletion |
| Equipment details `/staff/inventory/:id` | EquipmentDetail | Display image; action/physical-stock grouping; movement Table/paging; remove duplicate editor only after verified Edit replacement |
| Routine stock `/staff/inventory/:id/adjust` | StockAdjustment | Shared Fields/cards/review; frozen quantity/stock sequence and pending/idempotency locks unchanged |
| Admin correction `/admin/inventory/:id/reconcile` | InventoryReconciliation | Distinct warning/review; available-only projection, custody and audit unchanged |
| Administration `/admin/administration` | AdministrationPage / EmbeddedStaffDirectory | Compact overview and shared Staff/Admin directory; existing terms status only; no new publication surface |
| Staff/Admin directory `/admin/administration/accounts` | StaffDirectory / LegacyAccountDirectory | Shared Table/filter/pagination; name/email search; role and honest activation labels |
| Staff creation `/admin/administration/accounts/new` | StaffCreate | Reuse shared form; fixed existing STAFF provisioning, no Admin role selector |
| Staff/Admin details `/admin/administration/accounts/:id` | StaffDetail | Reuse details/audit; existing Admin remains read-only; no new Staff profile edit |
| Activation `/activate` | ActivationPage | Shared Field composition only; preserve fragment removal/token/password/single-use flow and missing-link behavior |
| Borrower equipment catalog `/borrower/equipment` and details | BorrowerCatalog / BorrowerEquipmentDetail | Existing card catalog/borrower shell; shared filters/paging/detail composition; preserve terms guard and read-only availability |
| Approved operational Inventory | InventoryDirectory | Investigate desktop Table minimum width/column proportions; preserve all buckets and accessible Open |

Approved Stage A shell, Borrower Directory and login remain intact. Administration has no existing privileged Admin provisioning or terms publishing UI; these are not invented. Borrower home and Phase 7/8 workspace placeholders are not reconstructed into unimplemented workflows. Phase 4B borrower terms behavior is retained.

## Execution and evidence

Implementation groups A–I completed; final verification results follow. Group order follows the owner brief: shared form/uploader; borrower screens; bulk; categories; equipment forms/images; details; stock/correction; Administration; consistency/responsive pass. Record final VERIFIED, IMPLEMENTED BUT NOT VERIFIED, DEFERRED and OWNER REVIEW PENDING status at closeout.

## Completed migration groups

- **A — VERIFIED:** shared ManagementPage/Card, semantic Field sections/actions/feedback and workflow steps. Existing FileDropzone now holds the image preview/current image inside the uploader, with unchanged validators and unsaved removal. Shared browser file-input check and 14 focused tests passed.
- **B — VERIFIED:** Student/Faculty creation, borrower details/profile/status/audit and activation presentation. Chromium verified textual leading-zero Student ID, configured test domain, external Faculty email/no Student ID, no Admin password fields, pending activation, profile persistence and confirmed deactivation/reactivation. Four browser checks, 65 screenshots; 36 focused frontend tests and build passed.
- **C — VERIFIED:** Student-only Excel upload/complete validation/selection/confirmation/results. All four test rows remain visible; three invalid rows cannot be selected. Double confirmation sends one request; roster reuse deactivates the Student and preserves Faculty. Two browser checks, 81 screenshots; three uploader tests and build passed.
- **D — VERIFIED:** category table, real unknown-total server paging and add/edit Dialog. Inactivation preserves existing equipment relationships. One browser check, 33 screenshots; 11 pagination tests and build passed.
- **E — VERIFIED:** Add/Edit equipment use the same Field sections and integrated image picker. Chromium forced an image-upload failure after successful metadata creation: retry uploaded only the image and did not create duplicate equipment. Existing image was retained after removing an unsaved replacement; Edit persisted replacement and metadata. Two browser checks, 35 screenshots; five focused tests and build passed.
- **F — VERIFIED:** Equipment Details displays the protected image and links to Edit; duplicate uploader removed only after E passed. Real status changes preserve physical stock and ledger; movement history uses shared unknown-total paging. One browser check, 17 screenshots; 14 focused tests and build passed.

All browser fixtures are disposable records in the separate PostgreSQL test database. Test terms are explicitly synthetic and confined to that database. No live email, owner credentials or normal development records were used.
- **G — VERIFIED:** routine Stock Adjustment and separate Admin Inventory Correction use the shared management composition. Chromium verified seven cases covering duplicate confirmation, signed add/remove and correction effects, mandatory explanation, invalid/insufficient quantities, concurrent stale review, persisted decrease/ledger and unauthorized correction rejection. 64 screenshots; 14 projection/schema tests and build passed. Existing PostgreSQL tests separately verify rollback and nonzero unavailable-custody preservation.
- **H — VERIFIED:** Administration now includes the shared Staff/Admin directory directly, with compact linked cards and actual terms status. Standalone search labels name/email; both directories use actual totals, disabled unavailable controls and independent URL page keys. Chromium verified fixed-STAFF creation, pending activation, status persistence and read-only existing Admin. Four checks, 65 screenshots; 51 auth/pagination tests and build passed. No Admin provisioning, terms publishing UI or role selection was invented.
- **I — VERIFIED (core pass):** Borrower catalog/details retain their existing shell and terms guard, read-only equipment/availability and protected image access. Shared filters/paging replace legacy controls. Full inventory stock columns fit 1366px expanded desktop; narrower screens retain controlled scroll. Chromium verifies unchanged shell node identity/document count across routes, centered collapsed seal, account-menu focus return and short 320/390px mobile Sheet keyboard containment. TERMS_NOT_PUBLISHED interception follows the real envelope and shows expected unavailable policy once across immediate navigation. Five checks, 75 screenshots; final tests/edge checks recorded below.

## Shared components and preserved behavior

Installed Base UI/shadcn Card, Field/FieldSet/Legend/Group/Label/Description/Error, Input, NativeSelect, Textarea, Alert, Table, Pagination, Dialog and AlertDialog compose the pages. Existing Sidebar, Sheet, DropdownMenu, Breadcrumb, Tooltip, Skeleton, Empty, Badge, Button, EquipmentImage and role-aware shells remain in use. No new package, incompatible Radix primitive or standalone page stylesheet was added. New management abstractions provide meaningful field associations, responsive sections, persistent feedback and action placement. Existing DirectoryHeader gained an optional domain label; defaults remain unchanged.

Image management uses the original PNG/JPEG limits, signature/decode validation, protected image Query, object-URL cleanup, actor/generation fencing, expected metadata version and idempotency keys. Metadata and upload remain separate existing commands with honest partial-success recovery. Removing a selected replacement preserves the stored image; stored deletion is not exposed. Return photographs remain outside eLabTrack.

Known-total account/equipment directories use actual page/per_page/total. Categories/options and movement/audit histories use actual batches without invented totals. Bulk preview shows all returned rows and selections. Embedded Staff filters/pages use independent staff_* URL keys; ordinary router paging preserves filters/history. Student/Faculty rules, fixed STAFF provisioning and read-only Admin remain unchanged. Auth, current-account enforcement, session refresh/logout, cross-tab coordination, terms acceptance, API clients/hooks/schemas and backend business logic were not edited.

The operational shell, approved Borrower Directory and immersive login were retained. Inventory's minimum width and column proportions were adjusted to fix confirmed expanded-desktop scroll and Archived badge overflow while preserving every physical bucket and the Open action. Borrower card catalog remains a card catalog, now using shared filters and pagination; its existing shell and terms gate remain intact. Stock arithmetic, review snapshots/sequence, explanations, idempotency, custody boundaries and immutable movement/audit history are unchanged.

## Verification diagnostics and limits

Initial harness runs exposed test-only issues: an escaped selector, out-of-band category fixture creation requiring explicit cache invalidation, transient deferred-search rendering, an unreferenced awaited CDP promise, a missing requestId in the synthetic error envelope, and checking Base UI focus guards before focus restoration. The harness was corrected and gates rerun; product auth/query policy was not changed to accommodate them. Chromium found a real Archived status-column overflow; proportions were corrected and boundary captures rerun.

Running serial and default parallel frontend suites simultaneously reproduced the inherited lazy-login one-second findByLabelText timeout (both 263 passed, one failed). This is the same baseline symptom recorded by Stage A. No test or timeout was weakened; independently rerun results are reported below. Passing reruns do not prove that inherited timing sensitivity is permanently resolved.

Screenshots use only explicitly synthetic disposable accounts/equipment. Actual institution domains/official terms and verified live Brevo delivery remain external dependencies. Synthetic test terms and acceptance were confined to the isolated test database; no normal database policy or personal accounts were changed. No hardware/screen-reader session, live mail or deployment check was performed. DOM semantics, labels, focus containment, keyboard navigation, touch-target composition, clipping and viewport overflow were checked where testable.

## Final quality gates

| Gate | Actual result | Evidence |
| --- | --- | --- |
| Frontend serial, final independent run | **VERIFIED — 264 tests / 20 files passed** | [Serial](verification/reconstruction-stage-b/final-serial.txt) |
| Frontend default parallel, final independent run | **VERIFIED — 264 tests / 20 files passed** | [Parallel](verification/reconstruction-stage-b/final-parallel.txt) |
| TypeScript and production build | **VERIFIED — tsc -b and Vite build passed** | [Build](verification/reconstruction-stage-b/final-build.txt) |
| Frontend lint | **VERIFIED — zero errors; 19 inherited warnings; no new warnings** | [Lint](verification/reconstruction-stage-b/final-lint.txt) |
| Backend go fmt ./... | **VERIFIED — passed, no files rewritten** | [Format](verification/reconstruction-stage-b/go-fmt.txt) |
| Backend go vet ./... | **VERIFIED — passed** | [Vet](verification/reconstruction-stage-b/go-vet.txt) |
| Backend go test ./... | **VERIFIED — passed; general run includes cached results / env-gated integration skips** | [Go tests](verification/reconstruction-stage-b/go-test.txt) |
| Fresh real PostgreSQL regressions | **VERIFIED — eight top-level tests / 34 subtests passed** | [Terms](verification/reconstruction-stage-b/pg-terms.txt), [Accounts](verification/reconstruction-stage-b/pg-accounts.txt), [Inventory](verification/reconstruction-stage-b/pg-inventory.txt), [Contention/correction/delivery](verification/reconstruction-stage-b/pg-contention.txt), [HTTP](verification/reconstruction-stage-b/pg-http.txt) |
| Chromium acceptance | **VERIFIED — 31 checks, 447 distinct accepted viewport screenshots, 11 group manifests** | [Visual index](../ux/verification/reconstruction-stage-b/README.md), [Machine-readable index](../ux/verification/reconstruction-stage-b/INDEX.json) |
| git diff --check | **VERIFIED — passed** | [Diff check](verification/reconstruction-stage-b/diff-check.txt) |
| Disposable environment cleanup | **VERIFIED — owned test API/frontend stopped, labeled test container/volume/network removed; private credentials/fixture files removed** | [Cleanup](verification/reconstruction-stage-b/cleanup.txt) |

Commands used the installed Node 24.19.0 and Go 1.27.1 toolchains. Frontend: `npm run test:run -- --no-file-parallelism`, separately `npm run test:run`, `npm run lint`, `npm run build`. Backend: `go fmt ./...`, `go vet ./...`, `go test ./...`; existing real integration cases ran with `-count=1` against the private Batch1 test environment on PostgreSQL port 54832, before browser fixtures. No dependency, test timeout, test assertion or type-check gate was weakened.

Database verification independently covers textual Student identity/domain/uniqueness, Faculty individual creation, Student-only selected roster matching, confirmation/replay/concurrency, activation single use/reissue/deactivation, terms version/acceptance rollback, physical stock conservation, reason/ledger/audit/receipt atomicity, invalid/insufficient/overflow quantities, simultaneous writers, expected-sequence correction, nonzero reserved/checked-out/damaged-held preservation, image validation/ownership and archive boundary. No backend test/source or migration was changed in Stage B.

The final selection-label bulk rerun passed both creation and Student-only deactivation. Edge verification confirmed maximum supported quantities, long equipment names, Archived badge and Open fit at 1366px expanded desktop; category Dialog focus restoration and 44px close target; actual Chromium DataTransfer drop/rejection; and final stock/correction review projections at 320/1366 in both themes. Normal local API `8080/api/v1/health` and frontend `5173/login` returned 200 before and after cleanup. Owned ports 54832/18085/15175 were closed. Production data was not accessed.

## Stage B acceptance checklist

| Requirement | Status |
| --- | --- |
| Shared form/uploader composition and installed Base UI primitives | VERIFIED |
| Student/Faculty creation, details/profile/status/activation presentation/audit | VERIFIED |
| Complete Student Excel preview, selection, confirmations and actual results | VERIFIED |
| Category table, add/edit, active status and relationships | VERIFIED |
| Add/Edit equipment, image preview/current image, replacement/removal and failure recovery | VERIFIED |
| Equipment Details, Edit action, physical quantities and movement history | VERIFIED |
| Routine stock and distinct Admin correction, frozen review, signed projections, persistence/audit | VERIFIED |
| Administration overview, embedded/standalone Staff/Admin directories, existing restricted controls | VERIFIED |
| Actual known/unknown-total paging, filter preservation and URL/history | VERIFIED |
| Seven widths, both themes, short heights, named table scrolling and unclipped fields | VERIFIED |
| Persistent shell, collapsed brand, menu/Sheet/Dialog keyboard focus and accessible names | VERIFIED |
| Current role/status, backend authority, refresh/logout/cross-tab and terms behavior | VERIFIED through unchanged source, auth/session tests and Chromium/backend permission checks |
| No code surface implemented but lacking the required local gates | IMPLEMENTED BUT NOT VERIFIED: none |
| Physical screen-reader/device and non-Chromium browser sessions | DEFERRED; not claimed as tested |
| Live Brevo delivery / final institutional terms and Student-domain approvals | DEFERRED external dependencies; not simulated in the normal application |
| Stage B visual acceptance | OWNER REVIEW PENDING |
| Phase 7 / borrowing, returns, replacements and other future features | NOT AUTHORIZED; no implementation begun |

No confirmed functional defect remains from this reconstruction. The inherited lazy-login test timing sensitivity remains a verification limitation, despite both final independent suites passing. Unknown-total APIs can legitimately offer one empty trailing page for a full final batch; the UI does not fabricate totals. Native assistive-device evaluation and owner visual acceptance remain outstanding.

## Exact changed files and Git status

Stage B changes are measured against the SHA-256 inventory of 1,223 Git-visible files captured at task start, preserving extensive pre-existing uncommitted work. The [complete change manifest](verification/reconstruction-stage-b/changed-files.json) enumerates every Stage B added/modified file, including all screenshot/manifests/logs. Backend source, migrations, manifests, API clients, feature hooks/schemas, auth/session/terms implementation, approved login/sidebar source and approved visual mockups match their starting hashes; no baseline file was removed.

Application source/test files modified in Stage B:

- `frontend/src/components/application/directory.tsx`
- `frontend/src/components/application/file-dropzone.tsx`
- `frontend/src/components/application/file-dropzone.test.tsx`
- `frontend/src/features/accounts/pages/account-create.tsx`
- `frontend/src/features/accounts/pages/account-detail.tsx`
- `frontend/src/features/accounts/pages/account-directory.tsx`
- `frontend/src/features/accounts/pages/activation-page.tsx`
- `frontend/src/features/accounts/pages/administration-page.tsx`
- `frontend/src/features/accounts/pages/student-bulk.tsx`
- `frontend/src/features/inventory/components/catalog-image-picker.tsx`
- `frontend/src/features/inventory/pages/category-management.tsx`
- `frontend/src/features/inventory/pages/equipment-form.tsx`
- `frontend/src/features/inventory/pages/equipment-detail.tsx`
- `frontend/src/features/inventory/pages/equipment-list.tsx`
- `frontend/src/features/inventory/pages/inventory-directory.tsx`
- `frontend/src/features/inventory/pages/stock-adjustment.tsx`
- `frontend/src/styles/reconstruction.css`

Added application/test-runner sources: `frontend/src/components/application/management.tsx` and `frontend/scripts/reconstruction-stage-b-qa.mjs`. Updated documentation: `docs/ux/FRONTEND_RECONSTRUCTION_CONTRACT.md`. Added this report and evidence under `docs/project/verification/reconstruction-stage-b/` and `docs/ux/verification/reconstruction-stage-b/`; the manifest lists exact individual paths.

The [full closeout git status](verification/reconstruction-stage-b/git-status.txt) includes inherited changes and is reproduced below. It is not a claim that all listed changes were made by Stage B. No commit, push, remote change or deployment occurred.

```text
 M backend/README.md
 M backend/internal/application/accounts/service.go
 M backend/internal/bootstrap/batch1_integration_test.go
 M backend/internal/bootstrap/contracts_test.go
 M backend/internal/infrastructure/spreadsheet/student.go
 M backend/internal/infrastructure/spreadsheet/student_test.go
 M backend/internal/interface/http/handlers/accounts_handler.go
 M backend/internal/interface/http/middleware/logger.go
 M backend/internal/interface/http/response/errors.go
 M backend/tests/integration/inventory_test.go
 M docs/ux/APPROVED_VISUAL_IMPLEMENTATION_MAP.md
 M frontend/src/app/router.tsx
 M frontend/src/components/application/operational-shell.tsx
 M frontend/src/components/application/shells.tsx
 M frontend/src/components/application/visual.tsx
 M frontend/src/features/accounts/pages/account-create.tsx
 M frontend/src/features/accounts/pages/account-detail.tsx
 M frontend/src/features/accounts/pages/account-directory.tsx
 M frontend/src/features/accounts/pages/activation-page.tsx
 M frontend/src/features/accounts/pages/administration-page.tsx
 M frontend/src/features/accounts/pages/student-bulk.tsx
 M frontend/src/features/auth/components/access-boundary.tsx
 M frontend/src/features/auth/pages/login-page.tsx
 M frontend/src/features/auth/product-auth.test.tsx
 M frontend/src/features/foundation/foundation.test.tsx
 M frontend/src/features/inventory/pages/category-management.tsx
 M frontend/src/features/inventory/pages/equipment-detail.tsx
 M frontend/src/features/inventory/pages/equipment-form.tsx
 M frontend/src/features/inventory/pages/equipment-image.tsx
 M frontend/src/features/inventory/pages/equipment-list.tsx
 M frontend/src/features/inventory/pages/stock-adjustment.tsx
 M frontend/src/features/terms/hooks/use-terms.ts
 M frontend/src/features/workspace/workspace-page.tsx
 M frontend/src/lib/api-client.test.ts
 M frontend/src/lib/api-client.ts
 M frontend/src/stores/session-action-store.ts
 M frontend/src/styles/index.css
?? backend/internal/domain/accounts/roster_error.go
?? backend/internal/interface/http/handlers/roster_error.go
?? backend/internal/interface/http/handlers/roster_error_test.go
?? docs/project/API_REQUEST_EFFICIENCY_AUDIT.md
?? docs/project/CULINARY_LOGIN_REDESIGN_REPORT.md
?? docs/project/FRONTEND_RECONSTRUCTION_STAGE_A_REPORT.md
?? docs/project/FRONTEND_RECONSTRUCTION_STAGE_B_REPORT.md
?? docs/project/IMMERSIVE_LOGIN_REDESIGN_REPORT.md
?? docs/project/PRE_PHASE7_PHASE5_6_REFINEMENT_REPORT.md
?? docs/project/PRE_PHASE7_UI_REFINEMENT_REPORT.md
?? docs/project/verification/
?? docs/ux/FRONTEND_RECONSTRUCTION_CONTRACT.md
?? docs/ux/SHADCN_COMPONENT_AUDIT.md
?? docs/ux/verification/culinary-login/
?? docs/ux/verification/immersive-login/
?? docs/ux/verification/phase56-refinement/
?? docs/ux/verification/pre-phase7-ui/
?? docs/ux/verification/reconstruction-stage-a/
?? docs/ux/verification/reconstruction-stage-b/
?? frontend/scripts/culinary-login-qa.mjs
?? frontend/scripts/immersive-login-qa.mjs
?? frontend/scripts/phase56-refinement-qa.mjs
?? frontend/scripts/pre-phase7-ui-qa.mjs
?? frontend/scripts/reconstruction-stage-a-qa.mjs
?? frontend/scripts/reconstruction-stage-b-qa.mjs
?? frontend/scripts/request-efficiency-qa.mjs
?? frontend/src/assets/login/
?? frontend/src/components/application/directory.tsx
?? frontend/src/components/application/file-dropzone.test.tsx
?? frontend/src/components/application/file-dropzone.tsx
?? frontend/src/components/application/management.tsx
?? frontend/src/components/application/operational-frame.tsx
?? frontend/src/components/application/server-pagination.test.tsx
?? frontend/src/components/application/server-pagination.tsx
?? frontend/src/features/accounts/lib/
?? frontend/src/features/accounts/pages/borrower-directory.tsx
?? frontend/src/features/inventory/components/
?? frontend/src/features/inventory/hooks/use-catalog-image.test.tsx
?? frontend/src/features/inventory/hooks/use-catalog-image.ts
?? frontend/src/features/inventory/lib/
?? frontend/src/features/inventory/pages/inventory-directory.tsx
?? frontend/src/features/terms/current-terms.test.tsx
?? frontend/src/hooks/use-directory-state.ts
?? frontend/src/lib/page-range.ts
?? frontend/src/styles/culinary-login.css
?? frontend/src/styles/immersive-login.css
?? frontend/src/styles/phase56-refinement.css
?? frontend/src/styles/reconstruction.css
?? frontend/src/styles/refinement.css
?? integration/phase56-workbook-fixtures.go
```

**Stage B implementation and local verification complete. OWNER VISUAL ACCEPTANCE: PENDING. Stop here for owner review; Phase 7 remains unauthorized.**

## Subsequent owner review — inventory filter checkpoint

The owner approved the overall Stage B visual implementation and identified one focused functional correction: inactive equipment with positive physical stock appeared under “Available only”. Historical pending-approval statements and original verification counts above describe the Stage B closeout before this review. Current visual acceptance is **OWNER APPROVED**.

The later, explicitly authorized [Inventory Filter Review](INVENTORY_FILTER_REVIEW_REPORT.md) changes the label to “Available for borrowing” and the server listing predicate to ACTIVE + available>0, preserving independent filters, server counts/totals/pagination and all stock/command/security behavior. Its own exact change list, regression evidence and final checkpoint readiness are recorded separately; it is outside the original frontend-only reconstruction verification totals. Phase 7 remains unauthorized.
