# Visual fidelity contract

**Current provisioning policy overlay, 2026-10-08:** [DEC-070 / approved rules and remaining gates](../project/ACCOUNT_PROVISIONING_POLICY.md) governs future behavior: only Admin creates accounts; Student requires unique textual official Student ID/approved SKSU email; Faculty is individual-only with valid unique accessible email and no required Student ID. Standard bulk creation/deactivation is Student-only. Both use borrower-owned separate passwords and current officially published Phase4B terms. This documentation overlay changes no implementation or approved mockup asset; original phase checkpoints remain historical.

2026-10-08. The owner-approved PNG is the visual target. [Manifest](approved/MANIFEST.json), [approved README](approved/README.md), [coverage map](APPROVED_VISUAL_IMPLEMENTATION_MAP.md) and [visual system](VISUAL_SYSTEM.md) bind later implementation. Confirmed Phase2.5 business rules override inaccurate image copy; preserve its visual treatment. Phase3A.1 governs behavior where not superseded by approved visual composition. Evergreen primary colors are superseded.

## Binding future feature gate

Before implementing any screen, open its individual PNG (not just an overview or filename), identify its layout and role, verify the current domain/API contract, then reuse the approved component/token system. Compare actual browser captures at the relevant phone/tablet/desktop widths and both themes. Track branding, color, type, spacing, containers, cards, controls, navigation, icons, borders/radii/shadows, responsive behavior, density, and empty/loading/error/conflict states. Keep screenshots and observed geometry with the feature report; unit tests alone are insufficient.

Correct policy errors without aesthetic substitution. Do not create separate Approved/Released enums, return photos/evidence uploads, price liability, Staff fine clearance, campus scope or hardware kiosk to imitate an image. Record ambiguous content as a dependency rather than inventing persistence/permissions. Functional feature acceptance still requires backend authorization, transactions/concurrency/idempotency and error/session boundaries.

**Match** below is an agent comparison finding about the composition family, not pixel identity or owner approval of every detail. **Deviation** retains the observed difference and its reason. The owner explicitly accepted the disclosed minor differences in the four current baseline previews after reviewing real implementation screenshots (DEC-064, 2026-10-08); future deviations still require owner review. **Blocked** means unavailable evidence/required unresolved contract; no current baseline row is blocked. No pixel-perfect claim. The generated font/antialiasing, static synthetic copy and sizes preclude meaningful whole-image pixel equality; compare at normalized width and preserve independent layout relationships.

## Four baseline comparison matrices

Borrower source941×1672 normalizes to390×693; actual captures retain900px viewport and full scroll content. S01 source1536×1024 normalizes to1440×960, upper0–526 and lower538–1024 reviewed separately. Its two frames are not one production page. Compare equal widths, not a distorted fixed-height image. [Comparison viewer](verification/phase3b/comparison.html), [browser data](verification/phase3b/RESULTS.json).

### B01 — Borrower Home

| Review row | Approved reference | Implemented result | Match / deviation / blocked | Observations | Action required |
|---|---|---|---|---|---|
| Header | B01 compact seal/divider/wordmark/bell | AppBrand + BorrowerTopBar64px | Match | Separate supplied seal, secondary notification access, discreet organization caption | Keep original master replacement dependency |
| Navigation | Four labeled destinations; top active marker | Four destinations72px plus safe area | Match | Same icons/task order; actual focus targets44px+ | None |
| Page composition | Greeting →2×2 metrics →Browse →Activity →Quick Actions →notice | Same ordered sections; four wide metrics | Match | Synthetic Alex and fixture snapshot replace illustrative personal/date data | None |
| Typography | Bold compact title, strong count and subdued metadata | 24px title/22px counts/system font/12–14px body | Deviation | Exact generated font unknown; tiny source labels enlarged; brand caption remains subordinate | Accepted by owner; retain rendered type, no font dependency invented |
| Spacing | ≈14px card edges, modest section gaps | 14px edges;10px card gaps;16–18px sections | Deviation | 44px interaction targets and readable rows increase scroll height over693px source | Retain accessibility; review density |
| Cards | White rounded bordered metrics/activity; quiet three quick cards | 12px surface cards, icon bubbles, quiet quick cards | Match | Original columns and picture/facts/status relationship retained | None |
| Controls | Violet Browse, secondary View All, three quick actions | Real shadcn Button/Link,44px target; read-only disclosures | Deviation | Preview-only deferred actions cannot imply recorded workflows | Remove disclosures only under later authorized features |
| Status styling | Active green, Pending amber, Completed quiet | Icon + semantic text/backing; PHP fine | Match | No color-only status; fine is PHP rather than source dollar symbol | Keep domain correction |
| Images/icons | Circular FSMO seal and generic equipment photos | Replaceable reconstructed seal, isolated catalog samples, Lucide | Match | No full-screen bitmap, private data or return photograph | Replace catalog samples in authorized feature |
| Responsive behavior | Approved mobile composition; wider behavior not imaged | 320/390 mobile;768/1280 enhancement; both themes | Deviation | Same four-destination nav at wide widths; dark is authorized engineering translation | Accepted by owner; retain inferred wide/dark composition |

### B02 — Equipment Catalog

| Review row | Approved reference | Implemented result | Match / deviation / blocked | Observations | Action required |
|---|---|---|---|---|---|
| Header | B02 compact brand and selected count/cart | Same shell; local selected-types counter | Match | Count is marked preview and does not imply a reservation | None |
| Navigation | Same four destinations, Equipment active | Labeled navigation, active top marker | Match | Summary clears72px navigation and safe area | None |
| Page composition | Title/search/filter/sort/chips/item rows/summary | Same anatomy; category wrapping; six local item cards | Match | No production inventory endpoint or payment surface | None |
| Typography | Strong item name, subordinate category/tags | 14px mobile name/12px category/10px tags;16px wide names | Deviation | Source compressed tag/quantity copy enlarged; system-font metrics differ | Accepted by owner; retain small subordinate labels |
| Spacing | Approx14px edges; compact horizontal rows | 14px edges;10px grid gaps;44px controls | Deviation | Rows/category hits taller;320 moves controls under facts | Keep44px targets; review full-page capture |
| Cards | Picture left, facts middle, availability/selection right | Same390 relationship; wide cards aligned action row | Match | Six individually extracted equipment pictures, not one UI bitmap | None |
| Controls | Selected −/quantity/+; Add; filter/sort; Review | Local quantity/search/category/sort/availability; Sheet and read-only summary | Deviation | No submit;0 removes draft; invalid integers show linked errors; type/unit totals explicit | Functional submit remains Phase7 |
| Status styling | Available green, Low stock amber, Unavailable red | Text/icons/backing, disabled Out of stock | Match | Attention is an explicit fixture view state, no threshold policy inferred | Future availability from authoritative contract |
| Images/icons | Six generic illustrative equipment photos and cart | Equipment-only PNG regions and Lucide | Match | No private photography or return-evidence feature | Replace fixtures only with approved catalog content |
| Responsive behavior | One mobile list; approved mobile hierarchy | 320/390 one column;768 two;1280 three, both themes | Deviation | Wide grid and dark translate supplied mobile target; selected dock is sticky and bounded | Accepted by owner; retain adaptation |

### S01 upper — Staff Dashboard

| Review row | Approved reference | Implemented result | Match / deviation / blocked | Observations | Action required |
|---|---|---|---|---|---|
| Header | Compact operational search, notification/account | 64px header with search-shaped disclosure, avatar/account control | Deviation | Source≈48px after normalization;44px controls and readable account text require64px | Accepted by owner; retain header/density |
| Navigation | Navy sidebar, indigo active, Admin additions in illustrated Admin frame | 238px navy/210px tablet rail; Staff operational destinations; Admin variant adds privileged entries | Deviation | Role correctness removes Admin entries from Staff; source normalized rail≈230px | Keep permission distinction; review small width difference |
| Page composition | Title/date →4 counters →trend/stock/activity panels | Same sequence and three-panel desktop relationship | Match | Static explicit fixture metrics, not unavailable backend readings | None |
| Typography | Compact heading, count-led cards, dense panels | 27px title/27px counts/12–16px panel labels | Deviation | Microscopic raster labels enlarged; exact generated font unspecified | Accepted by owner; retain density |
| Spacing | Compact edge/panel separation within upper frame | 24px workspace,16px gaps, readable operational panels | Deviation | Larger text/control targets extend content beyond cropped≈493px frame; no added business sections | Keep accessible density; inspect comparison |
| Cards | Four metric cards, bar chart, donut, event list | Retained cards and same three panel types | Match | Charts use corrected blue/lavender spectrum and textual alternatives | None |
| Controls | Date/range/View/account/theme controls | Existing store theme, scoped Button, read-only disclosures | Deviation | No fake search/query/account mutation; original toggle rendered as named button | Later functional controls need feature contracts |
| Status styling | Colored metric icons; illustrative maintenance/retired labels | Semantic metric icons; Available/Checked out/Damaged held/Reserved | Deviation | Four physical buckets reconcile124; inaccurate stock labels corrected | Keep Phase2.5 arithmetic |
| Images/icons | FSMO seal, operational icons, blue chart geometry | Supplied seal, actual Lucide, SVG chart from static local data | Match | Source charts are UI vectors; no raster dashboard used | None |
| Responsive behavior | Desktop framed dashboard | 1440 three panels;1024 two plus spanning activities; narrower Sheet; both themes | Deviation | Tablet/dark/small-width behavior inferred from approved hierarchy and UX rules | Accepted by owner; retain translation |

### S01 lower — Pending Requests

| Review row | Approved reference | Implemented result | Match / deviation / blocked | Observations | Action required |
|---|---|---|---|---|---|
| Header | Same operational header | Same StaffShell/search/account/theme | Deviation | Accessible64px header differs from cropped compact source | Accepted by owner alongside Dashboard |
| Navigation | Requests active in shared navy rail | Same selected destination; Admin entries omitted for Staff | Deviation | Server-role distinction wins over illustrated Admin identity | Keep restricted additions |
| Page composition | Title/filter/search, lifecycle pills, five-row table, footer paging | Same hierarchy and five-row first page | Match | Eight local rows; actual local search/filter/pagination; no issuing action | None |
| Typography | Compact identities, secondary program/times | 13px table/12px metadata; identity/secondary lines retained | Deviation | Readable type increases row heights, exact generated font unspecified | Accepted by owner; retain table density |
| Spacing | Tight rows with thumbnail/quantity/action alignment | 44px header, roughly68–78px identity rows; bounded scroll | Deviation | 44px Review targets; narrower labels make desktop Review fully visible | Keep accessible controls |
| Cards | White rounded table container, subdued headers/dividers | DataTableShell + unchanged Table primitives | Match | No whole-page overflow; keyboard named region at1024 | None |
| Controls | Selection/filter/search/Review/pagination | Local checkboxes/filter/search/paging; read-only Review dialog | Deviation | No bulk approval or row handover; Review cannot commit | Full review/commands Phase7 only |
| Status styling | Pending/Approved/Released etc; illustrative fines/replacement context | Canonical Pending/Active/Completed/Denied/Cancelled/Expired; PHP fine context | Deviation | Approval/release collapsed into single Active lifecycle; sample old fine cannot automatically block | Keep domain labels; add real derived filters later |
| Images/icons | Borrower avatars, equipment thumbnails, action arrows | Synthetic initials, isolated samples, Lucide | Match | No capstone people/contact data copied | None |
| Responsive behavior | Dense desktop queue | 1440 actions fit;1024 named table scroll; small-width nav Sheet; dark | Deviation | Core operations remain desktop/tablet optimized; full production phone cards are Phase7 work | Future narrow production feature still needs review |

## 32 future-reference inspection and dependency ledger

All below PNGs were visually inspected. These are findings for later authorized implementation, **not approved deviations or functional screens**. Their visual structures remain binding; only copy/permissions/contracts require the following care. Shared foundation coverage cannot close their individual runtime gate.

| Reference | Inspection finding / future action |
|---|---|
| [B02-01](approved/remaining/B02-01.png) | Cart layout/totals/Continue preserved; selection remains unreserved until actual submit, no monetary purchase total. |
| [B02-02](approved/remaining/B02-02.png) | Keep date/time review and visit banner; require Asia/Manila and existing current-version terms rather than a per-loan consent. |
| [B02-03](approved/remaining/B02-03.png) | Success/reference/expiry layout retained; Pending is reserved, not approved; confirm server commit before displaying success. |
| [B02-04](approved/remaining/B02-04.png) | Distinct denied/expired/cancelled history cards; denial reason mandatory and no reactivation of terminal request. |
| [S02-01](approved/remaining/S02-01.png) | Handover warning/actions retained; approval and physical issue commit once; no Approved-waiting screen. |
| [S02-02](approved/remaining/S02-02.png) | Four-stage direct issue retained; current borrower/terms/stock/date-time required; not a pending request. |
| [B03-01](approved/remaining/B03-01.png) | Active timing/quantity layout retained; any progress indicator is presentation, not a new duration/completion rule. |
| [B03-02](approved/remaining/B03-02.png) | Keep physical0/replacement3 active scenario and original due; incident history does not disappear after acceptance. |
| [B03-03](approved/remaining/B03-03.png) | Keep live fine panel; PHP10 ceil started24h, original due and unresolved replacements; no payment gateway. |
| [S03-01](approved/remaining/S03-01.png) | Good/damage/loss input and effect review retained; required incident notes, no photo/evidence field or acknowledgement. |
| [S03-02](approved/remaining/S03-02.png) | Result stock/history cards retained; automatic completion iff all physical/replacement0; no generic Mark complete. |
| [S03-03](approved/remaining/S03-03.png) | Acceptance adds new available/total; does not erase damage/loss or implicitly dispose original; partial remains open. |
| [B04-01](approved/remaining/B04-01.png) | Same cards/filter hierarchy; historical completed records may still have fine outstanding; no deletion. |
| [B04-02](approved/remaining/B04-02.png) | Keep profile/preferences; only supported approved edits (currently name) are actionable; email/category/contact edit unapproved. |
| [B04-03](approved/remaining/B04-03.png) | Correct illustrative staff-marks-complete wording to automatic physical0/replacement0; final institutional terms text needs review. |
| [S04-01](approved/remaining/S04-01.png) | DEC-070: account creation is Admin-only even if a generated Staff frame shows Add. Staff retains operational assistance. Admin individual form offers STUDENT/FACULTY with required textual Student ID/SKSU email for Student versus accessible email/no required Student ID for Faculty; Student-only bulk creation/deactivation. Preserve approved asset; adapt future behavior without Admin-assigned borrower passwords. |
| [S04-02](approved/remaining/S04-02.png) | Keep operational history/fine context; Deactivate and full-clear controls Admin-only, inactive obligations remain resolvable. |
| [S04-03](approved/remaining/S04-03.png) | Keep import stages/validation/subset review; illustrative .xlsx does not approve format/schema/activation or plaintext password columns. |
| [S05-01](approved/remaining/S05-01.png) | Separate catalog Active/Inactive/Archived from availability; A/R/C/damaged held only, liabilities never stock buckets. |
| [S05-02](approved/remaining/S05-02.png) | Keep reconciliation/movement history; Total tracked includes checked-out custody; typed stock authority, not arbitrary counters. |
| [S05-03](approved/remaining/S05-03.png) | Full fixed Admin clearance only; PAID/WAIVED/OTHER separate, history retained, unresolved loan can accrue later. |
| [S05-04](approved/remaining/S05-04.png) | Keep report tabs/filter/results; export format/limits remain reviewed dependency, not an image-created backend contract. |
| [S05-05](approved/remaining/S05-05.png) | Named accounts/terms/audit layout retained; current rate/TTL read-only until configuration separately approved. |
| [S05-06](approved/remaining/S05-06.png) | Keep filtered borrowings table; Overdue/Replacements are derived views rather than stored lifecycle enums. |
| [B06-01](approved/remaining/B06-01.png) | Keep own chronological event feed; no unread count/read receipt/provider delivery state inferred from decorative marks. |
| [S06-01](approved/remaining/S06-01.png) | Keep contextual empty/skeleton/conflict/error panels; failed reads never show fabricated zero/empty stock; request ID safe. |
| [B02-05](approved/remaining/B02-05.png) | Image/detail/quantity dock retained; catalog image purpose independent of returns; image placeholder not an upload contract. |
| [S02-03](approved/remaining/S02-03.png) | Required visible denial reason, release/history warning and safe cancellation; server state races still authoritative. |
| [B04-04](approved/remaining/B04-04.png) | Completed final fine/history retained; frozen final assessment can remain Outstanding; no current Overdue label. |
| [S04-04](approved/remaining/S04-04.png) | Named privileged accounts only; illustrative Pending secure activation is not a new persisted account enum; onboarding/recovery deferred. |
| [S05-07](approved/remaining/S05-07.png) | Exceptional correction Admin-only; exact approved procedure/locks/audit contract still needed and cannot hide custody. |
| [B06-02](approved/remaining/B06-02.png) | Preserve compact brand/form/recovery notice but remove authenticated borrower bottom navigation from sign-in; no signup or fake role session. |

## Open visual review, dependencies and disposition

The owner accepted the current engineered dark/wide adaptations, accessible density and disclosed small reference differences after reviewing real implementation screenshots. Future material visual changes still require review. The supplied seal remains replaceable pending an authentic original; generated fonts and equipment artwork are illustrative. These current baseline adaptations are disclosed and now owner-approved; no separate generated dark PNG or authenticated master seal is implied. Core color/navigation/hierarchy/material defects found during QA were corrected (metric flex direction, link semantics,44px controls, chart palette, toolbar/brand composition, visible desktop Review).

Future onboarding/activation, final terms/taxonomy, activity projection, import/export formats, original-equivalence/correction procedures and production data contracts follow existing OPEN_DECISIONS/feature gates. They do not authorize new schema or workflows in Phase3B. Remaining production screens require explicit feature-phase authorization, owner disposition of material visual changes and fresh screenshots against these same references. No return-photo feature is a dependency.

## Phase 4A authentication comparison addendum

[Rendered evidence](verification/phase4a/README.md) and the [Phase 4A report](../project/PHASE4A_REPORT.md) record new production-purpose login/session/forbidden and protected-shell verification. B06-02 was opened before final composition: compact seal/wordmark/theme, FSMO access eyebrow, Welcome back heading, labeled single-column card and violet action are retained. Protected bottom navigation is removed from anonymous login as required. Mobile body/control labels are larger than the compressed generated reference; desktop uses the existing responsive centered column, and dark uses approved semantic translations. The amber session notice appears only for an actual failure/expiry, not an always-present illustrative message. No pixel-perfect or new owner approval is asserted.

B01/S01 authenticated workspaces deliberately replace synthetic metrics/activity/table content with explicit feature placeholders. This changes content density because real services are deferred; approved shells, spacing/colors/navigation remain. The four original previews remain visually intact, verified independently over24 viewport/theme cases. B04-02 is a partial safe read-only identity surface, not completion of terms/contact/category/profile-edit/preferences functionality. Supplied reconstructed seal and exact generated font remain existing dependencies. Staff sign-out contrast was corrected after actual screenshot inspection; scoped styles do not alter preview shells.

The S04-03 import reference is bound to the Student-only roster and complete validation preview in [DEC-070 handoff](../project/ACCOUNT_PROVISIONING_POLICY.md); any generic category/role or password implication in static artwork cannot authorize Faculty import or credential assignment. Both categories require existing Phase4B officially published terms, with approved text still pending. This policy correction changes documentation only, not approved PNG/SVG/manifest files or implemented previews.
