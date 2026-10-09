# Frontend reconstruction — Stage A

Final checkpoint status, 2026-10-09: **engineering complete; Stage A owner approved; Stage B owner visually approved; Inventory filter owner manual verification PASSED** (explicit current owner reply “Verified and passed”). Final reconstruction visual/functional acceptance is confirmed. [Checkpoint audit and Phase7 handoff](FRONTEND_RECONSTRUCTION_CHECKPOINT.md) record the local-only checkpoint scope and remaining dependencies; Phase7 is not started or authorized. Historical pending statements and test counts below retain their original context.


Date: 2026-10-09. **OWNER VISUAL ACCEPTANCE: PENDING.** Stage A ends here; Stage B and Phase 7 are not authorized. Existing uncommitted work was preserved. No commits, pushes or deployment.

## Delivered scope

**IMPLEMENTED:** persistent operational shell, Borrower Directory, operational Inventory Directory, shared directory composition, server pagination and URL state. These are production routes using the existing APIs, not mock-data prototypes.

**PROTOTYPED:** none. Synthetic records exist only in the explicitly isolated browser-test database; they were never added to the owner's local database or application source.

**NOT YET MIGRATED:** Student/Faculty creation and account details/editing, restricted Staff management, equipment create/edit/details and uploader, category-management screen, bulk validation results, stock adjustment/correction screens, borrower-facing catalog/shell and historical previews. Existing functionality remains available through existing routes.

**FAILED/BLOCKED:** owner visual approval remains pending. Verification failures and remaining limitations are listed below; none authorizes expanding scope.

## Findings and component decisions

The inherited operational sidebar bypassed installed SidebarProvider behavior with `collapsible="none"`, hand-sized rails and a separate drawer breakpoint. Directory composition mixed custom cards/controls, inconsistent pagination and typography; accumulated CSS serves multiple historical and current consumers. Installed components already include 62 UI modules and Base UI, so reinstallation or a second library was unnecessary. See the [component audit](../ux/SHADCN_COMPONENT_AUDIT.md) for the complete inventory, consumer mapping and migration risks.

The operational frame now uses installed SidebarProvider, Sidebar/Header/Content/Group/Menu/MenuItem/MenuButton/Footer/Trigger, Breadcrumb, DropdownMenu and Sidebar's existing Tooltip/Sheet behavior. Table, Card, Input, NativeSelect, Skeleton, Empty, Button, Badge and Pagination primitives support the representative directories. No primitive or dependency manifest was changed. TanStack Table is not installed; these simple server tables use semantic Table and existing TanStack Query instead of adding a dependency to reproduce documentation examples.

New application components are `OperationalFrame`, meaningful directory header/filter/panel/table-region/stat compositions, `ServerPagination`, a URL-state hook and bounded page-range helper. Existing AppButton, ToneBadge and protected EquipmentImage remain useful adapters. Legacy forms, historical shell previews, image security hooks and borrower shell remain custom until their separately authorized migration. The sidebar replacement is presentation-only: routes, guards, current-account revalidation, auth/session transport and query ownership stay unchanged.

## Visual and navigation results

The expanded shell uses a 260px navy rail, centered FSMO seal/identity, role-aware links and an actual-account footer. Its 72px icon rail preserves Zustand preferences, centered seal, keyboard tooltips and accessible navigation controls. Theme and Sign Out remain available; the footer menu offers existing account navigation and logout, without demonstration features. Mobile uses installed Sheet behavior and its focus trap/Escape return; navigation closes the drawer. The persistent router parent keeps the shell while its Outlet changes.

Borrowers use a readable identity/type/program/status/activation table, grouped search/type/status filters and Admin-only existing create/bulk links. Student IDs remain strings. Missing accountability data is explicitly unavailable; older unclassified accounts are shown honestly rather than fabricated as Students. No account provisioning, activation or identity validation changed.

Inventory groups thumbnails/names, retains category/status/available filters and existing name/availability sort, and displays all four authoritative custody buckets. Metrics explicitly describe all matching equipment, not only the visible page. Quantity sizing also accommodates large real API values from disposable boundary fixtures. Existing image hooks retain authenticated access and generation fencing; Stage A did not modify upload behavior.

Both directories use actual server totals and page responses. Numbered pages/ellipsis are bounded; unavailable Previous/Next have no destination and cannot execute; pending navigation disables another page change. URLs preserve filters, browser Back restores state, and modified links retain browser behavior. The category picker shares pagination while explicitly stating that total pages are unavailable because its existing endpoint returns no total. A full last batch may require one subsequent empty page; no final-page count is invented.

Typography, spacing, 44px controls, cards, table rows, focus treatments and responsive filter/action layouts live in scoped `reconstruction.css` using existing FSMO colors/tokens. New duplicate rules were consolidated. Legacy CSS was retained because unmigrated screens and historical previews still consume it. Mobile actions wrap without clipped labels; wide tables scroll inside named keyboard-focusable regions rather than widening the document. Full legacy CSS removal belongs to Stage B after each consumer migrates.

## Verification

| Gate | Actual result | Evidence |
| --- | --- | --- |
| Frontend lint | PASS; 0 errors, 19 inherited warnings | [Log](verification/reconstruction-stage-a/lint-final.txt) |
| Frontend tests, serial | 263 passed, 20 files; baseline 252 + 11 new pagination/URL-state tests | [Log](verification/reconstruction-stage-a/serial-final.txt) |
| Frontend tests, default parallel | 263 passed, 20 files; independently executed | [Log](verification/reconstruction-stage-a/parallel.txt) |
| Production TypeScript/Vite build | PASS | [Log](verification/reconstruction-stage-a/build.txt) |
| Backend go fmt / go vet | PASS; no formatting changes | [fmt](verification/reconstruction-stage-a/go-fmt.txt), [vet](verification/reconstruction-stage-a/go-vet.txt) |
| Backend go test ./... | PASS; 25 test-bearing packages; opt-in DB checks skipped without their environment | [Log](verification/reconstruction-stage-a/go-test.txt) |
| Fresh real PostgreSQL regressions | 8 top-level tests, 34 subtests passed | [Terms](verification/reconstruction-stage-a/pg-terms.txt), [Accounts](verification/reconstruction-stage-a/pg-accounts.txt), [Inventory](verification/reconstruction-stage-a/pg-inventory.txt), [Contention/correction/delivery](verification/reconstruction-stage-a/pg-contention.txt), [HTTP](verification/reconstruction-stage-a/pg-http.txt) |
| Real Chromium | 17 acceptance checks passed, 25 screenshots, 0 uncaught runtime errors | [Log](verification/reconstruction-stage-a/browser.txt), [Acceptance manifest](../ux/verification/reconstruction-stage-a/acceptance.json) |
| git diff --check | PASS | Final verification recorded below |

Frontend commands: `npm run lint`, `npm run test:run -- --no-file-parallelism`, independently `npm run test:run`, and `npm run build`. Toolchains: Node 24.19.0 and the installed Go 1.27.1 Linux toolchain; no package/toolchain upgrades.

Backend commands: `go fmt ./...`, `go vet ./...`, `go test ./...`; opt-in PostgreSQL commands ran `TestRealTerms`, then `TestRealAccounts`, then `TestRealInventory`, then `TestRealInventoryContention|TestRealInventoryCorrectionSafety|TestRealAccountDeliveryRace`, and `TestBatch1HTTP|TestBatch1InventoryHTTP`. Existing migrations and grants were applied only to the owned disposable `elabtrack_v2_batch1_test` database on 54832. No migration files changed. Its container/volume, test-only API, Vite port 15175 and private fixture credentials were removed after verification. Normal owner PostgreSQL/API/Vite services were preserved.

Chromium exercised real Admin login; server pagination and browser Back; Student/Faculty/inactive/empty filters; stock metrics, sorting, available-only filtering, protected thumbnail; unknown-total categories; collapsed seal centering and keyboard tooltip; Ctrl+B; footer Escape focus return; 390/768/1024/1366px light/dark layouts; mobile drawer Escape focus return and close-on-navigation; existing detail routes; controlled loading/error/Retry; Staff controls/direct-route/backend permission denials; logout; and borrower denial from operational Inventory. One initial document navigation occurred for login; route changes caused none. Sidebar/header/main node identity stayed intact through desktop route navigation; viewport breakpoint changes legitimately switch desktop Sidebar and mobile Sheet.

Mobile action labels were checked for internal clipping, document width was checked at every capture, and tables use bounded horizontal scrolling. A screenshot review corrected mobile button wrapping before the final pass. Native keyboard/menu/dialog behavior was retained; no incompatible DOM event mechanics were added.

No backend/API/schema/authentication implementation was changed by Stage A. Existing frontend tests continue to cover current-account revalidation, role/account invalidation, session recovery, refresh coordination and cross-tab logout. Browser role checks supplement those boundaries; test success is not institutional or visual approval.

## Failures encountered and remaining limits

- An initial PostgreSQL run reused the populated browser fixture database. Empty-history migration rollback checks and the unpublished-terms precondition correctly failed. See the [original failed run](verification/reconstruction-stage-a/pg.txt). These failures were not suppressed or fixed by weakening tests; the final run uses a fresh owned disposable database in the required order.
- The original 45-minute browser harness expired during diagnostic work. It was restarted only against the owned test database.
- Full-page DevTools capture transiently resized the viewport to 1px and switched the installed sidebar into its mobile branch. Ordinary navigation diagnostics preserved sidebar/header/main nodes and made no document request. Capture now uses the stable viewport; the pagination screenshot scrolls to the footer. This was a test-capture artifact, not an application remount.
- A malformed correction authorization probe omitted `expected_sequence`, correctly returning 400; the final probe sends a valid payload and verifies permission denial.
- A role-switch probe forced a document navigation before the existing logout completed, aborting the cookie-clearing request. The corrected harness stays on the existing SPA login page until pending logout clears. Authentication code was not changed.
- The earlier reported parallel lazy-login timeout did not reproduce in the separately executed default parallel suite. This run does not establish that the intermittent baseline issue is permanently resolved.
- Automated Chromium coverage and screenshot inspection are not a full screen-reader audit. Real hardware, all browsers and institutional acceptance remain unrun.
- Brevo delivery configuration/testing and officially approved FSMO terms remain existing external dependencies; Stage A did not send mail, publish terms in the owner database or claim deployment readiness. Synthetic terms were confined to the disposable test harness.
- Prior unidentified unmatched-GET logging concerns are outside this frontend stage; no new claim of resolving them is made.

## Stage B handoff and owner gate

Review expanded/collapsed navigation, mobile drawer, both directories, tables/pagination and light/dark themes in the [screenshot index](../ux/verification/reconstruction-stage-a/README.md). The [reconstruction contract](../ux/FRONTEND_RECONSTRUCTION_CONTRACT.md) specifies shared tokens, forms and the future integrated image uploader.

After explicit owner visual approval, separately authorize migration of account forms/details and restricted Staff management; equipment create/edit/details with uploader previews inside forms; categories; stock adjustment/correction review composition; then bulk results and legacy CSS consolidation. Preserve RHF/Zod, backend image limits, unsupported stored-image deletion, idempotency/concurrency, custody and audit semantics. Return photographs remain outside eLabTrack. No Stage B form/uploader implementation or Phase 7 borrowing behavior was added.

## Exact files and Git status

The workspace already contained substantial modified/untracked Phase 5–6, login and prior refinement work. The following Stage A list is measured against the starting SHA256 snapshot, not against HEAD. [Exact machine-readable changed files](verification/reconstruction-stage-a/changed-files.json) includes every new source/document/log/screenshot path. [Preservation check](verification/reconstruction-stage-a/preservation.json) confirms all other starting Git-visible files and historical evidence remain byte-identical; no starting file was removed.

Modified existing files (Stage A only):

- `frontend/src/components/application/operational-shell.tsx`
- `frontend/src/features/accounts/pages/account-directory.tsx`
- `frontend/src/features/auth/product-auth.test.tsx`
- `frontend/src/features/inventory/pages/equipment-list.tsx`
- `frontend/src/styles/index.css`

The auth test change only updates the shell selector for the same collapse assertion. No test was removed or relaxed. Directory entry wrappers retain the existing Staff account directory and borrower equipment catalog while selecting the new operational screens.

New Stage A files (exact paths):

- `docs/project/FRONTEND_RECONSTRUCTION_STAGE_A_REPORT.md`
- `docs/project/verification/reconstruction-stage-a/browser.txt`
- `docs/project/verification/reconstruction-stage-a/build.txt`
- `docs/project/verification/reconstruction-stage-a/changed-files.json`
- `docs/project/verification/reconstruction-stage-a/diff-check.txt`
- `docs/project/verification/reconstruction-stage-a/git-status.txt`
- `docs/project/verification/reconstruction-stage-a/go-fmt.txt`
- `docs/project/verification/reconstruction-stage-a/go-test.txt`
- `docs/project/verification/reconstruction-stage-a/go-vet.txt`
- `docs/project/verification/reconstruction-stage-a/lint-final.txt`
- `docs/project/verification/reconstruction-stage-a/parallel.txt`
- `docs/project/verification/reconstruction-stage-a/pg-accounts.txt`
- `docs/project/verification/reconstruction-stage-a/pg-contention.txt`
- `docs/project/verification/reconstruction-stage-a/pg-http.txt`
- `docs/project/verification/reconstruction-stage-a/pg-inventory.txt`
- `docs/project/verification/reconstruction-stage-a/pg-terms.txt`
- `docs/project/verification/reconstruction-stage-a/pg.txt`
- `docs/project/verification/reconstruction-stage-a/preservation.json`
- `docs/project/verification/reconstruction-stage-a/serial-final.txt`
- `docs/ux/FRONTEND_RECONSTRUCTION_CONTRACT.md`
- `docs/ux/SHADCN_COMPONENT_AUDIT.md`
- `docs/ux/verification/reconstruction-stage-a/README.md`
- `docs/ux/verification/reconstruction-stage-a/acceptance.json`
- `docs/ux/verification/reconstruction-stage-a/account-footer-menu-light.png`
- `docs/ux/verification/reconstruction-stage-a/borrowers-empty-light.png`
- `docs/ux/verification/reconstruction-stage-a/borrowers-error-light.png`
- `docs/ux/verification/reconstruction-stage-a/borrowers-expanded-dark.png`
- `docs/ux/verification/reconstruction-stage-a/borrowers-expanded-light.png`
- `docs/ux/verification/reconstruction-stage-a/borrowers-loading-light.png`
- `docs/ux/verification/reconstruction-stage-a/borrowers-pagination-page2-light.png`
- `docs/ux/verification/reconstruction-stage-a/borrowers-responsive-1024-dark.png`
- `docs/ux/verification/reconstruction-stage-a/borrowers-responsive-1024-light.png`
- `docs/ux/verification/reconstruction-stage-a/borrowers-responsive-1366-dark.png`
- `docs/ux/verification/reconstruction-stage-a/borrowers-responsive-1366-light.png`
- `docs/ux/verification/reconstruction-stage-a/borrowers-responsive-390-dark.png`
- `docs/ux/verification/reconstruction-stage-a/borrowers-responsive-390-light.png`
- `docs/ux/verification/reconstruction-stage-a/borrowers-responsive-768-dark.png`
- `docs/ux/verification/reconstruction-stage-a/borrowers-responsive-768-light.png`
- `docs/ux/verification/reconstruction-stage-a/collapsed-tooltip-light.png`
- `docs/ux/verification/reconstruction-stage-a/failure.png`
- `docs/ux/verification/reconstruction-stage-a/inventory-catalog-image-light.png`
- `docs/ux/verification/reconstruction-stage-a/inventory-collapsed-dark.png`
- `docs/ux/verification/reconstruction-stage-a/inventory-collapsed-light.png`
- `docs/ux/verification/reconstruction-stage-a/inventory-expanded-dark.png`
- `docs/ux/verification/reconstruction-stage-a/inventory-expanded-light.png`
- `docs/ux/verification/reconstruction-stage-a/inventory-mobile-dark.png`
- `docs/ux/verification/reconstruction-stage-a/inventory-mobile-light.png`
- `docs/ux/verification/reconstruction-stage-a/mobile-navigation-dark.png`
- `docs/ux/verification/reconstruction-stage-a/mobile-navigation-light.png`
- `frontend/scripts/reconstruction-stage-a-qa.mjs`
- `frontend/src/components/application/directory.tsx`
- `frontend/src/components/application/operational-frame.tsx`
- `frontend/src/components/application/server-pagination.test.tsx`
- `frontend/src/components/application/server-pagination.tsx`
- `frontend/src/features/accounts/pages/borrower-directory.tsx`
- `frontend/src/features/inventory/pages/inventory-directory.tsx`
- `frontend/src/hooks/use-directory-state.ts`
- `frontend/src/lib/page-range.ts`
- `frontend/src/styles/reconstruction.css`

Exact `git status --short` at handoff, including preexisting work ([plain text](verification/reconstruction-stage-a/git-status.txt)):

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
?? frontend/scripts/culinary-login-qa.mjs
?? frontend/scripts/immersive-login-qa.mjs
?? frontend/scripts/phase56-refinement-qa.mjs
?? frontend/scripts/pre-phase7-ui-qa.mjs
?? frontend/scripts/reconstruction-stage-a-qa.mjs
?? frontend/scripts/request-efficiency-qa.mjs
?? frontend/src/assets/login/
?? frontend/src/components/application/directory.tsx
?? frontend/src/components/application/file-dropzone.test.tsx
?? frontend/src/components/application/file-dropzone.tsx
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

Final scope check: five existing files changed by this stage; every other starting Git-visible file was preserved. Backend sources/migrations, manifests/lockfile, auth transport/providers/guards, router, account/equipment form logic, login assets/styles and earlier screenshots remain unchanged from the start of this task. Normal frontend/API readiness endpoints returned HTTP 200 after isolated-service cleanup. Live email, deployment, institutional terms approval, broader screen-reader/hardware testing and owner visual approval remain unverified.
