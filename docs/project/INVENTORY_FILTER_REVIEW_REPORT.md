# Pre-Phase 7 inventory filter review

Final checkpoint status, 2026-10-09: **engineering complete; Stage A owner approved; Stage B owner visually approved; Inventory filter owner manual verification PASSED** (explicit current owner reply “Verified and passed”). Final reconstruction visual/functional acceptance is confirmed. [Checkpoint audit and Phase7 handoff](FRONTEND_RECONSTRUCTION_CHECKPOINT.md) record the local-only checkpoint scope and remaining dependencies; Phase7 is not started or authorized. Historical pending statements and test counts below retain their original context.


2026-10-09 — **CORRECTED AND VERIFIED; ready for the final Stage B checkpoint and owner review.** The owner approved Stage B's overall visual implementation. Phase 7 remains unauthorized.

## Root cause and prior behavior

“Available only” meant **A: physical available quantity greater than zero**, regardless of administrative equipment status. The frontend sent the existing `available_only` flag independently of Status/Category. PostgreSQL applied `NOT $4::bool OR e.available>0`. Consequently Staff/Admin's all-status directory included inactive Spoon with 20 available physical units. This was a listing-semantics defect, not incorrect stock arithmetic.

Borrower catalog reads already enforced ACTIVE equipment in the application service, so they did not expose inactive equipment. Authenticated reads, current-account checks, actor-scoped TanStack Query keys and the persistent router shell were correct and are unchanged.

## Final contract and isolated correction

Both Inventory and Borrower Catalog now say **Available for borrowing**. The existing API flag `available_only=true` requires equipment status ACTIVE **and** available physical quantity > 0. No additional equipment eligibility restriction was invented: inactive categories still do not hide existing active equipment. Account authorization and the separate officially published terms/acceptance borrowing gate retain their existing behavior; catalog availability does not grant permission to submit a borrowing request.

| Equipment status | Available stock | Checked filter |
| --- | ---: | --- |
| ACTIVE | Positive | Included |
| ACTIVE | 0 | Excluded |
| INACTIVE | Positive or 0 | Excluded |
| ARCHIVED | Any | Excluded |

Status, Category and Search remain independent and intersect with availability. Explicit INACTIVE or ARCHIVED + checked returns no matching records, total 0 and zero matching-stock aggregates. Unchecking preserves the selected Status, including its actual physical quantities. Unchecked/omitted availability retains existing operational archived visibility and Borrower ACTIVE-only visibility. Sorting and filters persist in URLs; switching filters resets page 1; Clear filters removes the independent inventory filters.

A frontend-only override of Status to ACTIVE would silently discard an explicit INACTIVE/ARCHIVED selection. The existing API could not express the desired independent intersection. The necessary backend change is therefore one read predicate in `InventoryRepository.List`: `NOT $4::bool OR (e.status='ACTIVE' AND e.available>0)`. The same predicate feeds rows, total count and all stock aggregates **before LIMIT/OFFSET**. There is no client filtering of a returned page, new parameter, schema change or inventory command change.

Physical stock, reserved/checked-out/damaged-held custody, ledger history, reconciliation, activation/status transitions, authorization, idempotency, locks and borrowing policy are unchanged. Empty filtered aggregates describe the matching set; they do not change the inactive Spoon's physical stock.

## Verification

| Check | Actual result / evidence |
| --- | --- |
| Focused frontend regressions | 4 passed: label and independent status/category integration; switch/clear/reset; search intersection; remount/URL restoration and authoritative page/count rendering |
| Full frontend tests | **268 passed, 21 files** — [log](verification/inventory-filter/frontend-tests.txt) |
| Frontend lint | Exit 0, no errors, **19 existing primitive warnings** — [log](verification/inventory-filter/frontend-lint.txt) |
| Production build | Passed with installed Node 24.19.0 — [log](verification/inventory-filter/frontend-build.txt) |
| Go formatting / vet / all tests | Passed with installed Go 1.27.1; integration cases that need opt-in environment are separately exercised below — [vet](verification/inventory-filter/go-vet.txt), [tests](verification/inventory-filter/go-tests.txt) |
| Real PostgreSQL / API filter and authority regressions | New `TestBatch1InventoryBorrowabilityFilter` with 15 subtests plus existing `TestBatch1InventoryHTTP` passed — [log](verification/inventory-filter/pg-api.txt) |
| Existing real inventory regressions | `TestRealInventory` with 8 subtests passed, including rollback, arithmetic, custody, ledger, replay, archive and images — [log](verification/inventory-filter/pg-inventory.txt) |
| Existing concurrency/correction regressions | `TestRealInventoryContention` with 3 subtests and `TestRealInventoryCorrectionSafety` passed — [log](verification/inventory-filter/pg-stock-safety.txt) |
| Real Chromium acceptance | **9 checks passed, 10 screenshots**, no runtime exceptions — [log](verification/inventory-filter/chromium.txt), [visual evidence](../ux/verification/inventory-filter/README.md), [sanitized manifest](../ux/verification/inventory-filter/acceptance.json) |
| Git whitespace checks | `git diff --check` passed — [log](verification/inventory-filter/diff-check.txt) |
| Disposable environment cleanup | Owned API/Vite stopped; labeled test container/volume/network and generated private fixtures removed; test ports closed; normal local API/frontend healthy before and after — [log](verification/inventory-filter/cleanup.txt) |

The PostgreSQL matrix has 31 equipment records: active/inactive × positive/zero stock, archived positive stock, and 26 additional active positive records in another category. Availability matches 27 records with available-stock aggregate 46. Actual server pages return 25 then 2 unique records with the same total/aggregate; page 3 is empty with the complete total retained. Both Admin and Staff run every filter case. Borrower ACTIVE-only visibility, anonymous rejection, invalid flag rejection and forbidden Borrower status filters are also checked. Explicit search/category/status intersections and inactive-category relationship behavior pass. Read operations preserve every equipment stock vector, status, sequence, metadata version, update timestamp and movement/audit/receipt counts.

Chromium uses real HTTP/PostgreSQL, exercises Admin/Staff controls and Borrower catalog, checks filter clearing/switching, search/category/sort, unchanged inactive physical quantity 20, inactive/archived empty intersections, 25/2 paging and page-two full refresh. Combined Status/Category/Search/Sort/Availability survive document refresh and authenticated session recovery. Ordinary filter navigation causes no document reload. The Borrower fixture accepts only clearly synthetic test-only terms in the disposable database under the existing terms gate. No official terms were published or accepted in the normal application. Light/dark screenshots at 320, 390 and 1366px show the label fitting the existing layout without document overflow; there is no CSS or visual redesign.

## Exact changed files for this correction

Pre-existing uncommitted Stage A/B and refinement work was preserved. This list is relative to the state at the start of this filter review, rather than the entire dirty worktree.

| File | Change |
| --- | --- |
| [backend/internal/infrastructure/persistence/postgres/inventory_repository.go](../../backend/internal/infrastructure/persistence/postgres/inventory_repository.go) | One listing predicate: ACTIVE + positive available stock |
| [backend/internal/bootstrap/inventory_filter_integration_test.go](../../backend/internal/bootstrap/inventory_filter_integration_test.go) | New real PostgreSQL/HTTP filter matrix, pagination/count and read-safety regression |
| [frontend/src/features/inventory/pages/inventory-directory.tsx](../../frontend/src/features/inventory/pages/inventory-directory.tsx) | Rename operational checkbox only |
| [frontend/src/features/inventory/pages/equipment-list.tsx](../../frontend/src/features/inventory/pages/equipment-list.tsx) | Rename Borrower catalog checkbox only |
| [frontend/src/features/inventory/pages/inventory-directory.test.tsx](../../frontend/src/features/inventory/pages/inventory-directory.test.tsx) | Four page/hook/API-boundary integration regressions |
| [frontend/scripts/inventory-filter-qa.mjs](../../frontend/scripts/inventory-filter-qa.mjs) | Reproducible isolated Chromium/API acceptance runner |
| [docs/API_CONTRACTS.md](../API_CONTRACTS.md) | Explicit availability semantics and independent filter intersections |
| [docs/ux/FRONTEND_RECONSTRUCTION_CONTRACT.md](../ux/FRONTEND_RECONSTRUCTION_CONTRACT.md) | Final terminology, owner visual approval and narrow read-query correction |
| [docs/project/FRONTEND_RECONSTRUCTION_STAGE_B_REPORT.md](FRONTEND_RECONSTRUCTION_STAGE_B_REPORT.md) | Subsequent owner approval and focused correction overlay; historical evidence retained |
| [docs/project/INVENTORY_FILTER_REVIEW_REPORT.md](INVENTORY_FILTER_REVIEW_REPORT.md) | This report |
| [docs/ux/verification/inventory-filter/README.md](../ux/verification/inventory-filter/README.md) | New screenshot/evidence index |

Generated evidence files are enumerated individually in [changed-files.json](verification/inventory-filter/changed-files.json), including the quality logs, sanitized acceptance JSON and ten PNGs. No migrations, manifests, authentication/session source, inventory command/business logic or visual styles changed.

## Readiness and remaining issues

No confirmed inventory-filter defect or security regression remains. The focused correction is ready for the final Stage B checkpoint; stop for owner review. Existing lint warnings, physical assistive-device/non-Chromium verification and external live-email/institutional approvals remain separate documented limitations, not new filter blockers. No Phase 7 work, commit, push or deployment occurred.
