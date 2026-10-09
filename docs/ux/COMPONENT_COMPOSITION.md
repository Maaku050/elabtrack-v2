# Component composition — approved Phase 3B foundation

**Future Phase11 composition,2026-10-10:** [DEC-080 binding selection](PHASE11_EQUIPMENT_SELECTION_UX.md) composes existing equipment-image/card, quantity and button primitives into a combined mobile quantity+Add row; cart/review uses a sheet or dedicated screen, with desktop browse/cart together. Reuse real Phase7 hooks/APIs and protected image handling. These are future requirements, not newly implemented components; the DEC-081 implementation pause remains.

2026-10-08. Owner-approved navy/violet screenshots supersede the former evergreen component proposal. [Visual system](VISUAL_SYSTEM.md), [coverage map](APPROVED_VISUAL_IMPLEMENTATION_MAP.md) and [fidelity contract](VISUAL_FIDELITY_CONTRACT.md) govern. Phase 3B is **COMPLETE AND OWNER VISUALLY APPROVED** after real implementation screenshot review (DEC-064). It implements shared presentation plus four **development previews**, not production business pages. All62 existing UI primitives and `components.json` are unchanged.

## Implemented compositions

| Actual component | Retained foundation | Reference / responsibility |
|---|---|---|
| AppBrand | Separate supplied PNG | B01/B02 compact seal/divider/wordmark; S01 FSMO identity; replaceable `sealSrc`, provenance recorded |
| AppButton / SurfaceCard | Button / Card |44px targets, scoped geometry and semantic tokens; Link rendering explicitly uses `nativeButton=false` |
| ThemeControl | Button + existing UI store | Same light/dark mechanism; no competing provider or persistence |
| BorrowerTopBar / BottomNavigation / Shell | Button, Router Link, semantic landmarks | Compact64px header, four destinations,72px safe-area navigation; no mobile sidebar |
| StaffSidebar / TopBar / Shell | SidebarProvider, Sidebar, Sheet, Button, Router Link | Navy238/210px sidebar, operational header/search/account, small-width Sheet; Admin additions are presentation only |
| PageHeading / SectionHeading | Semantic h1/h2 | Task hierarchy, description, optional action; no generic hero |
| MetricCard | Card, Lucide | B01 2×2 counters; S01 four operational metrics; explicit supplied values, no query/calculation |
| BorrowingStatusBadge | Badge, Lucide | Canonical lifecycle display; never Approved and Released as separate statuses |
| AccountabilityBadge / FineStatusBadge | Badge, Lucide | Physical/replacement and money separate; call only with known authoritative data in future features; unknown reads need feedback |
| AvailabilityBadge | Badge, Lucide | Supplied available/attention/unavailable state; no low-stock threshold inferred |
| EquipmentThumbnail / Card | Card, Badge, Button | B02 imagery/facts/availability/action; image fallback preserves name and layout |
| QuantitySelector | Input, Button | Local integer selection, explicit bounds/errors and zero-removal; never authoritative inventory validation |
| SearchField / DataTableShell / PaginationControls | Input, Table, Button | Labeled search, named keyboard scroll region, local pagination in preview |
| PreviewDialog | Dialog, Button | Clearly read-only development disclosure; no business mutation confirmation |
| ConflictAlert / LoadingRegion | Alert / Skeleton / Card | Persistent review conflict and busy geometry; existing reduced motion retained |
| Existing EmptyState / ErrorState / LoadingState | Empty / Alert / Spinner | Retained feedback from `components/feedback/states.tsx`; no duplicate primitive library |

Components live in `frontend/src/components/application/`. Lucide owns icons. Presentation props do not read/create role sessions or replace server authority. Sheets/Dialog use Base UI keyboard/focus behavior. Staff's built-in Sidebar is composed with the existing UI store and a1024px workspace Sheet boundary; the primitive's original768px hook and behavior remain unchanged. No sidebar-cookie/session persistence is added by these compositions.

## Representative previews

`frontend/src/features/visual-preview/` owns all synthetic fixtures, local filters/selection, read-only disclosures, charts and four pages. Router lazy imports them only inside `import.meta.env.DEV`. The Admin route renders the same Dashboard with a visual navigation audience; it is not a fifth business implementation or role switch.

| Target | Reused composition | Unique preview content |
|---|---|---|
| B01 Home | BorrowerShell, metrics, badges, thumbnail, headings | Greeting, Browse, activity rows, three quick actions, reservation notice |
| B02 Catalog | BorrowerShell, search, equipment cards, quantity, Sheet | Category filters, available-only toggle, alphabetical sort, selected-types/unit summary, read-only selection dialog |
| S01 upper Dashboard | StaffShell, metrics, cards, headings | Static stacked blue-spectrum bar chart with text table; physical-stock donut and text legend; synthetic activity |
| S01 lower Pending | StaffShell, search, table, badges, pagination | Eight local fixture requests, five-row pages, visible-row checkboxes, fine-context filter, read-only Review |

Selection is local intent, not a reservation. Review does not approve/deny/release. Filters against fixtures do not imply backend endpoints. Nonimplemented navigation explains that its feature is deferred. Existing Providers, QueryClient, SessionBootstrap, auth store/transport, cross-tab coordination and `/`, `/status`,404 remain intact.

## Later compositions — not implemented

Request cart/review/due-time, submission/Pending terminal detail, Staff issue/denial, return/replacement entry and results, fine clearance/history, account/terms/provisioning/import, inventory/movements/corrections, notifications and reports must follow their individual PNGs in the [map](APPROVED_VISUAL_IMPLEMENTATION_MAP.md). Shell or badge reuse does not mean those32 screens exist. ConfirmationDialog for consequential writes belongs to those feature contracts; PreviewDialog cannot serve as proof of a recorded operation.

Later pages follow page → feature hook → TanStack Query → feature API → centralized transport. React owns local draft state; Zustand appropriate global UI/session state only. Private responses obey existing generation fences. Backend owns permissions, lifecycle, stock/fine arithmetic, transaction/idempotency and integrity. No terms/activation/import/export/stock-correction policy is invented here.

No return photograph, camera, upload, evidence attachment/transmission/storage/retention or return-evidence step exists for any role. Equipment catalog images and Admin account-import file controls are independent later requirements.
