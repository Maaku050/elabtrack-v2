# Responsive rules — Phase 3A.1

Low-fidelity behavior contract, 2026-10-08; no CSS/components implemented or rendered layout verified. Applies to [59 surfaces](UX_ARCHITECTURE.md#screen-inventory-and-counting-rule) and [wireframes](WIREFRAMES.md). Screen behavior depends on available width, not detected role/device brand. Exact pixels are design proposals for Phase 3A.2 review, not locked domain policy.

## Width ranges and shell

Mobile **320–767 CSS px** (390 base wireframe,320 smallest target); tablet **768–1023** (768/834 review); desktop **>=1024** (1280/1440 review). Existing use-mobile boundary768 is retained as evidence;1024 is proposed desktop enhancement. At200% browser zoom, layout follows resulting CSS width and must not require whole-page horizontal scroll. Review320/390/768/1024/1440 and landscape/keyboard states later; no browser testing performed now.

Borrower mobile compact app bar and four bottom destinations; account/notifications secondary. Content bottom padding clears nav/sticky action/safe-area inset; keyboard-open input form can hide bottom navigation while preserving Back/Submit above keyboard. No doubled overlapping fixed bars: one action dock above nav where useful, or on focused form with navigation out of keyboard path. Tablet keeps bottom nav, desktop same labels in side navigation; no different feature workflow. Page headings and Back visible even in modal/sheet variants.

Staff/Admin desktop one-level Sidebar; tablet collapsible icon rail with labeled toggle/tooltips or Sheet when working area insufficient; mobile menu Sheet plus app bar. Desktop action bars and side summaries become stacked cards at narrow width. Inventory/docs/core control actions stay accessible, dense matrices optimized desktop/tablet with a mobile summary and scoped scroll alternative.

## Borrower behavior by major surface

| Surface / WF | Mobile320–767 | Tablet768–1023 | Desktop>=1024 |
|---|---|---|---|
| Sign in/terms WF-01/02 | One labeled form/scrolling terms column; inline acceptance/action, no nested scroll trap | Centered readable column; terms action after content | Same centered readable max-width; no unnecessary two-column login |
| Home WF-03 | Next-action card first; stacked pending/active summaries and Browse CTA | Summary cards2 columns, obligations still first |2–3 summary columns plus activity; content hierarchy unchanged |
| Catalog/filters WF-04 | **1 card column**; search full width; filter Sheet with Apply/Clear; selected count app-bar link | **2 columns**; search+filter row, Sheet/Popover; selected-request summary collapsible | **3 columns** at1024+, **4** only when card minimum width permits around1440; filters lateral/Popover; side selected summary only if grid still fits |
| Equipment detail WF-05 | Catalog image then name/description/available/quantity; bottom Add dock above nav with unobscured content | Image+facts side-by-side, quantity within facts | Readable2-column detail, action in summary; no enlarged meaningless image |
| Cart WF-06 | Unique item cards with labels, qty/remove; sticky Continue below totals (item/unit counts only) | Item list and due/help summary if sufficient width, otherwise stack | List plus right review summary; editable lines still ordered same |
| Due/time+review WF-07 | Date+time stack, timezone always visible; Calendar Popover can become Sheet; typed-input alternative; primary Submit clear of keyboard | Inline date/time fields alongside summary; review remains single logical flow | Items left, due+physical-visit summary right; keyboard order remains main→summary→submit |
| Submitted/Pending WF-08 | Status/Visit FSMO first, absolute expiry+advisory timer, item cards, Cancel secondary | Info pairs2 columns; visit/expiry remain prominent | Items/evidence alongside status sidebar; no extra release/approval controls |
| Denied/Expired/Cancelled WF-09 | Reason or expiry/cancel explanation before items/history; Start new request inline | Readable card with 2-column context | Same detail composition wider, no historical state conflation |
| Active/Overdue WF-10 | Due/Overdue/live fine, physical/replacement remaining cards; return instructions visible; no data table base | Aligned2-column obligation/fine summaries | Item table with labeled quantities + accountability side panel; read-only borrower action set |
| Partial/Replacement/Completed WF-11 | Per-item issued/good/damage/loss/P/U labels; replacement remaining highlighted; collapsible event history |2-column physical vs replacement panels | Labeled quantity table plus event timeline/final-fine panel; no ghost C |
| My Borrowings/History WF-12 | Tabs accessible by keyboard; cards, Load more; statuses wrap | Two-list-column layout only if chronological reading order unambiguous; default single list | Wider rows or bounded table; filters toolbar, no additional lifecycle |
| Notifications WF-14 | Stacked own event rows, meaningful link/time/status, no unread-count dependency | Wider single readable list | Single readable feed, not enormous grid; operational version different permitted scope |
| Account WF-13 | Label/value stack, supported edit fields one column; terms/theme/logout reachable | Two-column display groups; edit form readable | Readable identity/preferences panel; no role control introduced |

Available quantity text/counter remains readable at320px. Filter chips wrap; no horizontal chip-only control necessary. Input controls target at least44×44 CSS px for primary touch interactions; text/input layout responds to text enlargement. Image has informative alt or decorative handling; absence/error preserves card hierarchy/name/available/action.

## Staff/Admin behavior by major surface

| Surface / WF | Desktop>=1024 | Tablet768–1023 | Mobile320–767 fallback |
|---|---|---|---|
| Dashboard WF-15 | Sidebar+action queues,3-column summaries |2-column summaries, rail/collapsible sidebar | Stacked action summaries; menu Sheet; counts labeled |
| Pending WF-16 | Bounded table, borrower/expiry/accountability key columns, Open only | Keep key columns, secondary fields in expanded row | Request cards with expiry/account/fine; Open review, no row approve |
| Request review/deny WF-17/18 | Borrower/request main, accountability/actions side; confirmation Dialog | Two areas if fit, otherwise stack; denial Dialog | Identity→time→quantities→due→context→actions; full-height denial/confirm Sheet, focus managed |
| Direct issue WF-19 | Stepper + selection table + side summary | Stepper compact, summary collapsible | One step at a time, labeled cards, sticky Next/Issue with keyboard clearance; no four-column form |
| Borrowings WF-38 | Filtered table, detail item/event panels | Essential columns with detail links, panels stack | Loan cards; detail stack; all return/replacement actions reachable |
| Return WF-20/21 | One row/item; prior counts separated from now-input columns; preview to right or below | Essential P/U+now inputs, prior counts disclosure; preview below | One item card, Good now/Damaged now/Lost now labeled full-width input groups; prior/P/U text; Review action at bottom; no squeezed table |
| Replacement WF-22 | Obligation table+quantity/type fields, effect review | Compact row/cards, split panels only if fit | Obligation cards, accepted quantity then type/rationale; stacked before/after stock; full-height review |
| Borrowers WF-23/24/25 | Table/list+detail; current/history tabs; create form single logical column | Essential borrower/category/status, action panel below | Cards; detail/form fields stack; privileged actions only Admin |
| Inventory WF-26/27/28 | Table A/R/C/D/T; detail counts+secondary ledger; metadata form/image | Essential A/R/C then count detail; ledger scoped scroll if needed | Cards show availability/lifecycle, expand exact physical counts; ledger labeled scroll with sticky context; edit single column |
| Stock/archive WF-29 | Typed Dialog and before/after exact quantity preview | Dialog / Sheet if space limited | Full-height Sheet, clear effect/reason, no hidden confirmation |
| Fine clear/history WF-30 | Dialog amount/basis/method/note/history; Admin only | Same readable width, history below | Full-height Sheet, fixed amount not editable, stacked history, confirmation clear of keyboard |
| Bulk import WF-31 | Multi-step template/file/validation table/results; row error detail | Preview reduces nonessential columns; detail for errors | Step-by-step file input and valid/invalid summary; row cards for errors; optional labeled scroll for full preview; never force whole-page overflow |
| Deactivate/privileged accounts WF-32/33 | Named-account table and impact Dialog | Condensed table, full context review | Account cards and full-height impact form; no unexpected role dropdown for Staff |
| Settings/terms WF-34 | Current policy panel + terms version list/editor review | Panels stack, readable terms width | Display current policy; text editor/review stacked; publish confirmation; no provider controls |
| Reports/audit WF-35/36 | Filters and bounded dense table with key row context | Filter Sheet, essential columns/scroll region | Summary cards + labeled horizontally scrollable result where unavoidable; recommend larger screen for complex comparison, retain filter/back/action access |

## Dialogs, scrolling and action placement

Use modal Dialog for short desktop actions, AlertDialog for consequential confirmations, Sheet/Drawer for mobile editing/filtering. In full-height mobile forms use one content scroll region, titled header, Close/Back and footer without hiding last field/error behind actions. Consent/time selection must work without drag gestures. Dismissible overlays trap focus while open, Escape safely cancels, restore focus to trigger; no forced focus trap on a normal full page. Closing critical unsaved form asks Keep editing / Discard when needed.

Desktop split-pane is optional layout only, not a new route or changed processing order. Tablet stack early if labels/controls become cramped. Any horizontal scroll is contained and labeled (ledger/report/import), reachable by keyboard, with sticky row identifier where useful; never a whole-page scrollbar for borrower forms. Load/error/empty regions retain approximate reading structure at every width. A sticky action must never hide error summary, totals, keyboard focus or nav.

## Review and themes

Both light/dark modes use identical semantic labels, error messages, focus hierarchy and actions. Lifecycle/overdue/fine states cannot depend on color; icon accompanies text, never replaces it. No color palette, type scale, final spacing or polished imagery chosen in this phase. Reduced motion has no automatic scrolling/carousel/confetti; loading uses static or reduced-motion skeleton behavior later. Future checks belong to [UX_ACCEPTANCE_CHECKLIST](UX_ACCEPTANCE_CHECKLIST.md), all pending until implementation exists.
