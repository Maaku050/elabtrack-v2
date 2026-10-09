# shadcn component audit — Stage A

Audit recorded before modifying application screens, 2026-10-09. Current source is implementation evidence; the owner’s Stage A brief authorizes layout changes while preserving product rules. Historical previews and verification remain evidence, not production data.

## Installed foundation

`components.json` selects **base-nova**, CSS variables and Lucide. React 19, TypeScript 6, Tailwind v4, Base UI 1.8.0, shadcn React 0.3.x and TanStack Query 5 are already installed. There are **62 UI modules**, **38 with direct Base UI imports**, and **no Radix imports**. TanStack Table is **not installed**; use the existing semantic Table with server queries. Do not add or upgrade dependencies to reproduce documentation demos.

Reviewed all TypeScript/TSX/CSS files under components/ui, components/application, features, app and styles; inspected composition and the affected render paths. The inventory below distinguishes direct Base UI imports from compositions/vendor helpers; absence of a direct import does not imply Radix.

## Installed module inventory

| Module | Implementation imports |
| --- | --- |
| accordion | @base-ui/react/accordion |
| alert-dialog | @base-ui/react/alert-dialog |
| alert | Native HTML / existing primitive composition |
| aspect-ratio | Native HTML / existing primitive composition |
| attachment | @base-ui/react/merge-props, @base-ui/react/use-render |
| avatar | @base-ui/react/avatar |
| badge | @base-ui/react/merge-props, @base-ui/react/use-render |
| breadcrumb | @base-ui/react/merge-props, @base-ui/react/use-render |
| bubble | @base-ui/react/merge-props, @base-ui/react/use-render |
| button-group | @base-ui/react/merge-props, @base-ui/react/use-render |
| button | @base-ui/react/button |
| calendar | react-day-picker |
| card | Native HTML / existing primitive composition |
| carousel | embla-carousel-react |
| chart | recharts, recharts |
| checkbox | @base-ui/react/checkbox |
| collapsible | @base-ui/react/collapsible |
| combobox | @base-ui/react |
| command | cmdk |
| context-menu | @base-ui/react/context-menu |
| dialog | @base-ui/react/dialog |
| direction | @base-ui/react/direction-provider |
| drawer | @base-ui/react/drawer |
| dropdown-menu | @base-ui/react/menu |
| empty | Native HTML / existing primitive composition |
| field | Native HTML / existing primitive composition |
| hover-card | @base-ui/react/preview-card |
| input-group | Native HTML / existing primitive composition |
| input-otp | Native HTML / existing primitive composition |
| input | @base-ui/react/input |
| item | @base-ui/react/merge-props, @base-ui/react/use-render |
| kbd | Native HTML / existing primitive composition |
| label | Native HTML / existing primitive composition |
| marker | @base-ui/react/merge-props, @base-ui/react/use-render |
| menubar | @base-ui/react/menu, @base-ui/react/menubar |
| message-scroller | Native HTML / existing primitive composition |
| message | Native HTML / existing primitive composition |
| native-select | Native HTML / existing primitive composition |
| navigation-menu | @base-ui/react/navigation-menu |
| pagination | Native HTML / existing primitive composition |
| popover | @base-ui/react/popover |
| progress | @base-ui/react/progress |
| questionnaire | Native HTML / existing primitive composition |
| radio-group | @base-ui/react/radio, @base-ui/react/radio-group |
| resizable | Native HTML / existing primitive composition |
| scroll-area | @base-ui/react/scroll-area |
| select | @base-ui/react/select |
| separator | @base-ui/react/separator |
| sheet | @base-ui/react/dialog |
| sidebar | @base-ui/react/merge-props, @base-ui/react/use-render |
| skeleton | Native HTML / existing primitive composition |
| slider | @base-ui/react/slider |
| sonner | next-themes, sonner |
| spinner | Native HTML / existing primitive composition |
| switch | @base-ui/react/switch |
| table | Native HTML / existing primitive composition |
| tabs | @base-ui/react/tabs |
| textarea | Native HTML / existing primitive composition |
| toast | @base-ui/react/toast |
| toggle-group | @base-ui/react/toggle, @base-ui/react/toggle-group |
| toggle | @base-ui/react/toggle |
| tooltip | @base-ui/react/tooltip |

## Migration mapping

| Current implementation | Recommended component | Reason | Migration risk | Affected screens |
| --- | --- | --- | --- | --- |
| StaffShell uses Sidebar collapsible=none plus separate Sheet and manual width/margins | Controlled SidebarProvider, Sidebar icon, Header/Content/Groups/Menu/Footer/Trigger | Built-in collapse, tooltip, keyboard and mobile state must agree with Zustand | Medium: persistent identity, viewport breakpoint, portalled menus | Operational shell; legacy preview retained |
| Custom topbar context | Breadcrumb + Router Link | Semantic navigation without document reload | Low: existing role destinations preserved | Operational header |
| Title/toolbar spacing differs by management page | Shared DirectoryHeader and filter section | Consistent title/actions/44px controls and readable text | Low; scoped styles only | Borrowers, Inventory |
| Borrower and inventory raw HTML tables, bespoke loading/error layouts | Table + Card + Skeleton + Empty | Semantic headers, bounded scrolling and common state presentation | Medium: actual API values and row links | Two representative directories |
| Previous/Next only; preview pagination renders every page | Shared shadcn Pagination composition | Known totals, bounded numbered ranges/ellipsis, disabled links, URL state | Medium: server paging, no client slicing | Two directories; category picker unknown-total navigation |
| Raw select labels | Installed NativeSelect with associated label | Browser keyboard support and consistent control presentation | Low | Representative filters |
| SurfaceCard/AppButton/ToneBadge | Retain wrappers around Card/Button/Badge | Existing FSMO semantic tones and control sizing remain useful | Low; avoid global overrides | New shared directory patterns |
| Multiple stock metric compositions | Compact shared stock cards | All four authoritative buckets remain visible, scope labeled | Low: no arithmetic changes | Inventory |
| Custom FileDropzone and EquipmentImage | Retain security-aware implementations; Stage B uploader contract | File selection and actor-scoped image caching are project behavior | Medium: object URL lifetime, upload validation | Equipment forms/details deferred |
| React Hook Form/Zod forms, stock review panel and confirmation | Retain now; Field/Input/Textarea/Dialog contract for Stage B | Working validation, transaction/idempotency safety must survive migration | High for stock/activation workflows | Account, equipment, category, stock and correction forms deferred |
| BorrowerShell / immersive login / visual preview routes | Retain unchanged | Outside representative operational redesign | High for auth; unauthorized redesign | Borrower-facing and login routes |

## Confirmed composition and accessibility issues

- Shell collapse is currently separate from SidebarProvider state. Its keyboard shortcut can change provider state without matching the manual CSS width. The custom drawer uses 1024px while installed Sidebar uses 768px. Seal/nav centering and footer affordances depend on overrides rather than coherent primitive state.
- CSS repeats `.staff-sidebar`, `.staff-workspace`, `.staff-top-bar`, metrics and toolbar rules across application/refinement/inventory/management. SurfaceCard globally overrides Card layout. Some former sign-in selectors remain alongside the approved immersive login styles; they are not safe to delete indiscriminately because historical previews remain.
- Directories use locally held filters/page; browser Back cannot restore a directory query. Empty/loading/error layouts and pagination presentation differ. Labels/status text can be small; empty panels take disproportionate space.
- Borrower activation text “Established” overstates what activation-not-required proves. Use the existing account-status presentation helper. Obligations unavailable must never be represented as clear/no obligations.
- Existing focus mechanics for Base UI menus/sheets/tooltips must remain. New native filters require associated labels; horizontal tables need a named focusable scroll region; collapsed links require accessible names and real tooltips. Disabled pagination anchors need click suppression and tab exclusion.

## Stage A decision and boundaries

Implement a new operational frame using existing Sidebar primitives, retaining the legacy shell only for visual previews. Preserve persistent route parent, auth hooks, guards, query keys and current-account revalidation. Share directory composition and pagination only for Borrowers and operational Inventory; keep Staff directory and borrower catalog on their existing render paths.

Add one scoped reconstruction stylesheet; avoid further global management/login overrides. No primitives reinstalled, no dependency changes, no backend/API/schema/auth changes. Full CSS consolidation and remaining form/category/roster migrations belong to Stage B after owner approval. Unused installed UI modules remain available; removal is not justified by this audit.

## Primary references

Reviewed [Base Sidebar](https://ui.shadcn.com/docs/components/base/sidebar), [Base Pagination](https://ui.shadcn.com/docs/components/base/pagination), [Base Data Table](https://ui.shadcn.com/docs/components/base/data-table) and [Base Field](https://ui.shadcn.com/docs/components/base/field). Local source defines compatible APIs (Base UI render props, not Radix asChild). The data-table guide uses TanStack Table; the existing installation supplies Query only, so no Table dependency is introduced. Owner screenshots guide proportions/density and expanded/collapsed/footer behavior, not demo content, tenancy or charts.
