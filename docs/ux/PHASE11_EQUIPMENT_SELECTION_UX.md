# Phase 11 — binding equipment selection UX

2026-10-10. **Confirmed owner requirements, DEC-080; current implementation authority DEC-082 after Phase10 exit gate.** DEC-081 accepted functional Phase7 and its local checkpoint; the later master prompt renewed Phases8–14 authority. Both V1 references are re-inspected. Phase11 engineering gate PASS; [verification](../project/PHASE11_REPORT.md). Final owner acceptance remains pending. [Current plan](../project/PHASE11_IMPLEMENTATION_PLAN.md). The exact comparison below remains binding.

## Reference and current evidence

The owner's two V1 equipment-selection screenshots are the binding interaction references, including desktop browse-and-cart composition. Retain FSMO branding, current navy/violet design tokens, approved shadcn/Base UI components and the persistent application shell. Do not copy V1 prices/payment features or infer business rules from its presentation.

**Reference received and inspected, 2026-10-10:** The owner supplied `elabtrack_phase11_v1_references.zip` in the repository root. The previous missing-reference dependency is resolved. The table records the exact archive members, decoded PNG dimensions and SHA-256 hashes so the comparison is traceable without importing V1 assets into the application.

| Member under `elabtrack_phase11_v1_references/` | Dimensions | SHA-256 |
|---|---|---|
| `elabtrack_v1_quantity_add_reference.png` | 591 × 154 px | `d0dbd41d7748c9b579ef341ca931e63c334617c90f197647495098f950e443cc` |
| `elabtrack_v1_browse_cart_reference.png` | 1311 × 588 px | `5265a1918bf6d4f2df8a5f515f110c1100a4292cd02dcdb86558e19f40b73510` |

The archive's `REFERENCE_NOTES.txt` was read as supporting reference material. It is not separate execution authority: the earlier reference-preparation message authorized documentation only; the latest master prompt/DEC-082 supplies current implementation authority. Screenshots establish visible layout and control affordances; they do not establish click behavior, editable-input semantics, accessibility, backend arithmetic or mobile responsiveness. No V1 application or Firebase access was performed. Background personal/sample records are excluded from this specification.

**Phase7 source baseline before Phase11:** `frontend/src/features/borrowing/pages/request-compose.tsx` lists name/category/available stock with Select, initially adding one unit to a local form. Quantity editing appears separately under Your selection, then a dedicated review stage. These selection rows have no equipment images or inline minus/quantity/plus/Add controls. Phase7 is the functional baseline, not evidence that this Phase11 requirement is implemented.

Reuse the protected `EquipmentImage`, `useEquipment` server queries/pagination, `useBorrowingCommand`, `borrowingApi.submit`, RHF/Zod validation and reviewed idempotency-key/busy guard. [Phase7 API contracts](../API_CONTRACTS.md#phase-7-implemented-borrowing-contracts) govern authority, payloads and transactions.

## Exact screenshot comparison and V2 adaptation

**FACTUAL V1 VISUAL EVIDENCE — quantity/Add crop:** A single white, lightly outlined rounded card occupies almost the full 591px crop width. An isolated equipment thumbnail sits on the left; a bold name, green availability line and per-unit price are stacked to its right. At the lower right, minus, a centered quantity of `1`, plus and a filled green Add button form one uninterrupted horizontal row. Add is immediately beside plus, rather than below or in another panel. The visible Add button is approximately 66 × 35px; the minus/plus glyphs have no visible button outline. These screenshot dimensions do not prove actual hit areas and must not replace V2's accessible touch targets.

**FACTUAL V1 VISUAL EVIDENCE — browse/cart overview:** The 1311 × 588px image shows a large white New Transaction modal over a dimmed dashboard. Below its title/close control is a three-stage strip: Select Student → Browse Equipment → Set Due Date. The content has a **single vertical column of equipment cards on the left**, with search above it, and a cart on the right. It is not an equipment-card grid. The visible left column is approximately 568px wide, the cart 584px wide, with a 28px gap: roughly equal browse/cart space. Two complete cards and part of the next card are visible; this does not prove how scrolling or pagination works.

The overview repeats the crop's card structure for two equipment types, with quantity `1` and adjacent Add controls at each card's lower right. The cart has a heading with an item count, one white item row containing a name and monetary calculation, and right-aligned minus, quantity `3`, plus and trash controls. A monetary Total sits at the lower right below a divider. Wide Back and Next actions occupy the bottom row beneath both panes. The depicted catalog quantity and cart quantity are distinct; the screenshots do not show how repeated Add reconciles them.

| Concern | Visible V1 reference | Current V2 source baseline | Binding V2 treatment |
|---|---|---|---|
| Equipment identity | Left thumbnail, bold name, availability underneath; category is not visible in either image. | Composer rows show name, category and available stock, with no thumbnail. Existing catalog has protected image handling. | Each selection card shows the actual protected catalog image, name, category and available physical quantity. Category comes from the owner requirement, not an inferred V1 label. |
| Availability meaning | A green `Available: numerator/denominator` label. One visible pair is `9/5`, so neither denominator meaning nor stock consistency can be established from these images. | Composer uses `stock.available`; its query requires ACTIVE equipment and positive available stock. | Use the authoritative available physical quantity. Do not copy the V1 fraction, example values or infer total-tracked stock from it. Preserve existing eligibility and visibility. |
| Quantity beside Add | Both images show minus → centered quantity → plus → Add together at the card's lower right. | Select adds one unit; editing is in a separate Your selection section. | Move quantity choice into the same card action row as Add: `[ − ] [ quantity ] [ + ] [ Add ]`. Keep compact spacing and sufficient targets, including on mobile. |
| Desktop composition | Search and vertically stacked cards alongside a cart, approximately equal pane widths. | Choose equipment and Your selection are separate cards in the existing request composer; no V1-style paired image-card/cart composition. | Show browsing and cart concurrently where width permits. Preserve this functional adjacency and hierarchy using the existing application shell; the screenshot's pixel widths and modal are not new V2 design tokens. |
| Cart editing/removal | Minus, displayed quantity, plus and trash are visible in the item row; typing into the displayed number is not demonstrated. | Labeled numeric inputs and Remove are already present in Your selection. | Preserve editable whole quantities and item removal in the cart itself. Use accessible minus/plus, labeled input and Remove control; no return to browsing is required to edit a selected line. |
| Cart summary and review | Item count, monetary line/Total, Back and Next are visible. An actual Next transition is not shown. | Local selection → Review Request → explicit Confirm Request, with stable reviewed idempotency key and busy guard. | Show equipment-type and selected-unit counts, then a separate borrowing review/confirmation. No prices, monetary total, payment or purchase checkout. |
| Mobile experience | No mobile viewport or cart sheet is supplied. The small image is a card crop, not evidence of mobile rendering. | Existing selection/edit/review flow is responsive but lacks the binding combined card controls. | Mobile-friendly cards retain the entire quantity/Add row; cart editing/removal and review use a sheet or dedicated screen. This is confirmed owner direction, with implementation testing still pending. |
| Styling and flow | White modal, pale borders, blue stage strip, green Add/Next, red/green cart icons; Select Student and Set Due Date are visible. | Approved FSMO navy/violet shadcn/Base UI styling and persistent React Router shell; Borrower and Staff direct flows are distinct. | Reuse approved tokens, primitives, typography, status colors and shell. Do not copy the old palette, administrative student-picker wizard or borrower due-date step. Phase7 retains role-specific workflow authority. |

**CONFIRMED V2 INTERACTION REQUIREMENTS:** Card quantity controls prepare local intent; Add places the selected amount into a cart containing one line per equipment ID. The cart permits quantity editing and removal, with explicit feedback and whole-number validation. Search, filters and server pagination must not discard selections within browsing. A separate review precedes the existing confirmed submission. These requirements come from the owner's directions and Phase7 contracts; screenshot affordances alone do not prove V1 implemented them safely. Repeated Add must never create duplicate IDs or silently alter a cart quantity; present the existing selection/edit state clearly rather than guessing V1's unseen behavior.

**Remaining visual acceptance dependency:** Compare the eventual authorized Phase11 implementation against both references for identity hierarchy, the combined action row, desktop browse/cart adjacency and cart controls. Mobile layout, light/dark styling, keyboard behavior and touch targets require new V2 runtime evidence; neither supplied screenshot proves these gates passed.

## Equipment cards and mobile action row

Every equipment card in the selection interface must include:

- The actual equipment catalog image via existing authenticated image handling. Missing/failed images use an honest accessible fallback, not an unrelated photograph.
- Name and category, with the existing Uncategorized fallback where applicable.
- Latest available **physical** quantity from the API, distinct from the local selected quantity.
- Compact quantity controls and explicit Add-to-cart action.

Mobile keeps all selection controls in the **same action row**:

```text
[ − ] [ quantity ] [ + ] [ Add ]
```

Compact spacing must preserve the existing minimum44×44 CSS-pixel primary touch targets, visible keyboard focus, equipment-specific accessible control names and labeled editable quantity. Support keyboard operation and whole-number validation at320px/390px and enlarged text. Do not detach Add into another section/footer while leaving quantity elsewhere. A compelling responsive exception must document the actual width/text constraint and be reviewed; the combined row is the default.

Add changes only the local draft. Use positive whole quantities and one cart line per equipment ID, consistent with the existing payload. Give clear item/quantity feedback and cart editing/removal. Displayed stock guides controls but does not promise reservation or authorize borrowing. Never silently reduce a selected quantity after a stock conflict.

Where existing visibility rules show an unavailable equipment card, keep its action row with disabled controls/Add and a clear reason. Do not make zero-stock equipment selectable or expose inactive/archived equipment that the current API hides from that account.

## Mobile review and desktop composition

Mobile uses an accessible cart/review sheet **or a dedicated review screen**. Show selected equipment, editable quantities/removal, equipment-type and selected-unit counts, and clear review/confirmation actions. A sheet supports focus management, dismissal/return and safe-area/keyboard clearance. A dedicated screen remains an approved alternative.

Desktop shows equipment browsing and the cart together when width permits, inspired by the owner's V1 reference. Preserve readable cards/controls and browse → cart → review keyboard order. Tablet/narrow desktop may stack these areas while keeping each quantity/Add action row together. Search/filter/server-page changes retain selections within browsing. Preserve the persistent SPA shell and private-state cleanup on logout/account invalidation.

Show equipment-type/unit counts only. No equipment prices, monetary cart total, payments or purchase checkout are introduced.

## Backend authority and integration

Preserve page → feature hook → TanStack Query → feature API → centralized transport. Query owns catalog/server state; selection is client draft state. Reuse protected image caching/object-URL cleanup and server pagination/counts. Do not filter only the current page or download the entire inventory for cart selection.

Review/confirmation uses existing `POST /api/v1/borrowings`, `{items:[{equipment_id,quantity}],confirm:true}` and UUID Idempotency-Key. Add/edit/review must not create a borrowing or reserve stock. Only committed Phase7 submission creates PENDING and immediately reserves in the existing transaction. Keep separate review, duplicate-submit protection and a stable reviewed key for safe retry.

Backend checks remain authoritative: current role/account/activation/category, applicable current terms, ACTIVE equipment and sufficient available stock. Preserve locks, atomic multi-item reservation, inventory conservation, immutable ledger/audit/history,24-hour expiry and one-time cancellation/denial/expiry release. Existing Staff direct-issuance/physical-handover authority and due-time ownership remain unchanged. No schema, payload, stock arithmetic, activation, authentication or borrowing-policy change is needed for this presentation requirement.

Stock/terms conflicts refresh affected server data and retain the draft for explicit correction/review. A retained draft grants no privileges. Official FSMO publication/current consent and required activation readiness continue gating live borrowing; isolated synthetic demo terms are not institutional approval.

## Future acceptance gates — not executed results

1. Compare both inspected V1 references with future rendered selection views for the elements in the comparison table; record responsive exceptions and preserve FSMO visual direction. Receipt/inspection of the references does not satisfy this implementation gate.
2. Verify actual images/fallbacks, facts and combined controls in light/dark at320/390/768/1024/1440 widths, enlarged text and keyboard focus, without whole-page overflow or undersized targets.
3. Verify invalid/zero/excess quantities, unique cart lines, editing/removal/counts, search/filter/pagination preservation and empty cart. Add/review produce no borrowing/stock write.
4. Verify mobile sheet or dedicated review, desktop browse/cart, focus/dismissal/navigation, explicit confirmation and duplicate-submit protection.
5. Exercise real Phase7 APIs/PostgreSQL: atomic multi-item reserve, stale/insufficient stock, eligibility/terms/role/account changes, idempotent retry, cancellation/expiry conservation and preserved custody/history.
6. Run required quality gates and real Chromium acceptance with responsive light/dark evidence. Do not mark completion or owner acceptance from this specification alone.

**Current handoff:** Phase10 prerequisites passed and DEC-082 authorizes implementation under this specification. Phase11 engineering/visual/runtime acceptance is still pending. No migration, V1 import, commit, push or deployment is needed for selection UX. Historical documentation-only checks below do not count as implementation acceptance.

## Initial binding-document verification — before reference receipt

The new specification/DEC-080 cross-links and eight affected documents were checked for consistency. Go formatting, vet and tests passed (existing cached tests; database integration was not enabled). Frontend lint passed with19 inherited warnings and production TypeScript/Vite build passed. The first two-worker frontend run had277 passes and two lazy-login wait timeouts while the build ran concurrently; rerunning the unchanged suite alone passed all279 tests across23 files in43.35s. No test, timeout or application code was changed to obtain that result. Both attempts are disclosed; the concurrency explanation is an inference from the observed runs.

`git diff --check` and the new-file whitespace check passed. No PostgreSQL, Chromium, Docker or deployment verification was rerun for this documentation-only update; the future Phase11 acceptance gates above remain unexecuted. Existing application source/migrations have no diff, and prior demo work and unrelated files are preserved.

## Reference-comparison documentation verification

Both archive images were decoded and visually inspected at original resolution; member dimensions and hashes above identify the exact inputs. Archive notes were treated as reference material only. The comparison separates screenshot facts, current V2 source evidence and confirmed future requirements, including the absence of mobile/category evidence and the ambiguous V1 stock fractions. Related DEC-080/status references are reconciled to show the references are received while implementation remains paused. Document links and whitespace are checked for this update; application, database and browser gates are not rerun for this documentation-only comparison. The initial binding-document results above are historical results, not a claim of a new test run.
