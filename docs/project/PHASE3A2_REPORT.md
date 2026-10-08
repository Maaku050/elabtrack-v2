# Phase 3A.2 — High-Fidelity Mockups & Visual Design Specification

**2026-10-08 · CURRENT · DESIGN PACKAGE COMPLETE · OWNER VISUAL APPROVAL PENDING.** Phase 3B is NOT STARTED. Current explicit owner direction confirms the Phase 3A.1 architecture as approved and authorizes this static visual design package only. Creation of these artifacts does not approve their visual choices or authorize production implementation.

## 1. Sources reviewed

Current phase attachment and AGENTS.md; PROJECT_CHARTER, SOURCE_OF_TRUTH, DECISIONS, OPEN_DECISIONS, ROADMAP, STACK, FOUNDATION_AUDIT and Phase 2/2.5/3A.1 reports; all eight current domain documents and seven existing UX documents. Current manifests, components.json, styles/font/theme, router/providers/query boundaries and UI primitive inventory inspected read-only. Current rules supersede historical alternatives and earlier awaiting-review/not-started wording. No personal capstone content or V1/private borrower photographs copied. V1 audits and boss audit are already-recorded technical references, not visual or operational authority; no new V1 Firebase access or boss work.

The palette's normal text and essential control contrast targets are grounded in [W3C contrast minimum](https://www.w3.org/WAI/WCAG22/Understanding/contrast-minimum.html) and [W3C non-text contrast](https://www.w3.org/WAI/WCAG22/Understanding/non-text-contrast.html). No private notes/plugin/external design service is required. A temporary local SVG generator and browser capture scripts live in `/tmp`, outside the delivered product source; they only render static synthetic designs.

## 2. Visual direction

Calm, precise, welcoming and accountable FSMO identity: deep evergreen, warm neutrals, modest monoline catalog placeholders and task-oriented hierarchy. Borrower obligations/discovery and Staff counter work take precedence over decorative heroes, vanity metrics and generic SaaS analytics. Visual choices are proposed, not locked product policy. [Visual system](../ux/VISUAL_SYSTEM.md).

## 3. Color/token strategy

Both themes share semantic background/foreground/card/muted/muted-foreground/border/control-border/primary/primary-foreground/secondary/accent/destructive/warning/success/information/focus. Light primary `#12645C`; dark primary `#80D7BC` with dark foreground. Essential control boundaries are stronger than decorative dividers. Lifecycle and derived conditions retain explicit labels/icons; money and physical/replacement quantities have separate cards.44 calculated contrast pairs meet proposed targets; detailed values and limits are in the visual system.

## 4. Typography

Preserved existing system font; no font dependency installed. Mobile task title24px, wide28px, section/card18–20px, body16px, labels/metadata/tables14px minimum. Tabular numerals for quantities/money/timestamps; no microscopic table text or oversized desktop hero.

## 5. Spacing/layout

4px baseline and8px rhythm; borrower padding16/24/32 at390/768/1280, centered max1120 on desktop. Staff216px rail at1440, compact176px at1024,24–32px working padding. Card16–24px, section20–24px,44px controls;8/12/16px control/card/dialog radii. Terms/sign-in max640. Fixed navigation/action/keyboard/safe-area behavior explicitly specified; full-page SVG exports do not simulate scroll or keyboards.

## 6. Component styling

Primary/secondary/outline/ghost/destructive/disabled/loading/focus visual rules; inputs/selects/cards/badges/tabs/tables/sheets/dialogs/AlertDialogs/search/filter/quantity/stat/empty/skeleton/alerts/toasts specified. VS01 shows control/status/sign-in examples. Persistent summaries provide consequential results; Sonner only supplements them. Future HTML focus/keyboard/semantic behavior remains an implementation gate.

## 7. Borrower mobile shell

390px source composition: compact brand/notifications bar, four visible destinations Home/Equipment/My Borrowings/Account, clear active state and action separated above navigation. “My Borrowings” wraps between words, not inside the label. Catalog selected-equipment access appears near the top and in the action dock. No borrower sidebar base.

## 8. Borrower desktop adaptation

768 keeps approved bottom destinations;1280 translates the same labels into top navigation and adds readable paired panels/three catalog columns. Mobile complete tasks remain the source; no mandatory right summary or desktop-only action. Desktop is an enhancement, not a different product.

## 9. Borrower high-fidelity composition count

**15 (MB01–MB15)** covering all 22 borrower surfaces plus the scoped activity variant.99-artboard package has additional Home normal/pending/active variants, with attention as the base, explicit denied/expired/cancelled variants and completed/fine detail. [Composition index](../ux/HIGH_FIDELITY_MOCKUPS.md#composition-index).

## 10. Staff/Admin shell

Desktop/tablet working rail: Dashboard, Requests & Borrowings, Inventory, Borrowers; Admin additionally Reports and Administration. Theme/account/activity in the top bar.1024 compact rail plus identity-first fallback for dense tables; below that menu Sheet/stacked core fields are specified. Staff/Admin share visual identity, not identical authorization.

## 11. Staff/Admin composition count

**16 (MS01–MS16)** cover all 23 operational and12 Admin-only surfaces via detailed stage/dialog/tab/result variants. **Three shared boards (VS01–VS03)** cover sign-in/access and control/empty/loading/error examples. Total **34 unique compositions**, not99 unique pages.

## 12. Catalog design

Mobile1/tablet2/desktop3 columns. Each card has generic equipment illustration, name/category/available/short description, quantity controls and Add to request. Borrower catalog shows no reserved/damaged/lost/ledger counts. Search/category/available-only filters and mobile Sheet Apply/Clear are specified. Catalog images remain independent from return processing.

## 13. Request design

Selection cards → required date/time (Asia/Manila) → item/due/current-terms review → Submit request. No monetary checkout, seven-day cap or payment wording. Reminder explains submission reserves for24h and physical approval/release still requires FSMO visit. Availability/terms/account/date conflicts recover in the affected step.

## 14. Pending design

Equipment reserved/Pending/Visit FSMO headline; submitted/absolute expiry/advisory time/requested due/items, confirmed own cancellation. No Approved state or implied handover. Zero advisory timer checks authoritative status before showing Expired/released. Denied reason and terminal history remain visible.

## 15. Active borrowing design

Original issue/due, issued and physically outstanding quantities, in-person return instructions; borrower read-only. Staff/Admin alone records returns and condition. No Borrower Mark returned, completion, payment, camera/evidence/attachment control.

## 16. Partial/replacement design

Human labeled issued/good/damaged/lost/still-out plus required/accepted/remaining replacements. Example3issued:1good,1damaged,1physically out;1replacement owed. P0/Upositive variant remains Active and can become Overdue. Completion requires both zero; no ghost checked-out liability units or automatic damage/loss peso price charge.

## 17. Fine design

PHP amount visually separate from units, due/as-of/current or frozen final assessment, history and clearance method. PHP10 per started positive24h; exact due0,1min10,24h10,25h20. Completed detail freezes fine but may retain outstanding money. Clearing does not resolve quantities/close borrowing; an active fine can accrue later.

## 18. Staff request review

Borrower/category/active account/terms/fine/held quantities/submitted/expiry/final due before Approve & release. Confirmation says the equipment is physically handed over now. Deny is a separate destructive reason-required dialog, visible to Borrower and releases hold only after commit. No queue-row approval.

## 19. Direct checkout

Four visible design stages: Select Borrower, Select Equipment, Set Due Date & Time, Review & Issue. The static board shows stages for review; future flow presents one active step with edit/back. Active registered borrower/current terms/future due/available stock rechecked authoritatively. Immediate Active borrowing, no Pending.

## 20. Return processing

Per-item prior issued/good/damage/loss/physical remaining and new good/damage/loss quantities clearly separated; defaults0 and total cannot exceed remaining. Full-good shortcut fills remaining and clears draft incident quantities only after overwrite review; existing replacement liability persists. Damage/loss note follows current engineering rationale. No photos/evidence/file feature.

Mixed example: A12→13,R0→0,C3→1,D0→1,T15→15; P3→1,U0→1. Review exact effect before Record return; persistent result shows quantities, state, original due and fine. Lost units reduce total; only good restocks. Completion automatic only P0/U0; final fine freezes then. Temporary screenshots caught and corrected clipped footers and label wrapping before handoff.

## 21. Replacement acceptance

Borrower/original loan/equipment/reason, required/already accepted/accept-now/remaining; input max=current remaining. Accept1 adds A+1/T+1 as a new acquisition, retains damaged original D1, resolves U1→0. With P1 still Active; P0 variant completes automatically. Stale quantities reload and review again; no over-acceptance.

## 22. Inventory

Operational catalog image/name/category/A/R/C/D/status/incident-replacement indicator/Open. Detail stock, metadata, movements and incident/liability sections separate. T=A+R+C+D; loss history/U are not physical buckets. Metadata cannot override custody; ordinary movements typed/count/reason; archive guarded by zero R/C/U. Exceptional Admin correction procedure remains later dependency, not fabricated backend/UI contract.

## 23. Borrower management

Generic Borrower directory with Student/Faculty category, program where relevant, account status, current borrowing, replacement/fine/history. Staff individual-Borrower provisioning; Admin deactivate/reactivate impact confirmation preserves obligations/history. Account profile supported display-name edit only; email/category/program/status read-only in Borrower view.

## 24. Bulk import

Admin three-stage upload/validation/result board. Invalid and duplicate rows explained, explicit valid subset choice, created/rejected/unknown outcomes retained; unresolved command checks original result before retry. No plaintext password-column requirement. Exact template/activation mechanics remain later contracts, not guessed by visuals.

## 25. Admin fine clearance

Current/final assessed and full outstanding, noneditable amount, Paid/Waived/Other and optional note. Confirm full balance; retain original amount, Admin and time. Active example full-clearPHP10at1h can leave newPHP10at25h when assessment reaches20. No partial payment, gateway or waiver-as-revenue implication.

## 26. Reports/admin

Admin report types Inventory/Active/Overdue/History/Incidents-Replacement/Fines; meaningful date range, bounded results and disabled export placeholder (format undefined). Separate Paid/Waived/Other/assessed figures. Administration deliberate Accounts/Terms & policy/Audit tabs; named roles, reviewed versioned terms, current24h/PHP10 read-only, bounded actor/event/date history. No generic settings dumping ground.

## 27. Loading

VS02 has five explicit Skeleton contexts: catalog cards, borrowing cards, staff rows, detail panels, dashboard cards. Stable footprint, regional busy semantics and static reduced-motion option specified; no default full-screen spinner.

## 28. Empty states

VS02 renders nine required cases: empty catalog/cart/requests/active/history/overdue/replacements/borrowers/reports. Quiet icon/title/contextual action; distinguish empty filters from no data and failed reads. No decorative illustration dependency.

## 29. Error/conflict states

VS03 renders network/server/validation/403/session/availability/processed/expired/stale return-replacement/unknown write, with persistent human recovery copy and actions. Safe request-ID pattern for unexpected failure, no secrets/stack traces. Unknown write offers Check result with repeat mutation disabled. SH02 also documents inactive/absent-resource variants. No fake success or0balance after a read error.

## 30. Responsive coverage

**Nine critical borrower screens** each have390 light/dark,768 light,1280 light artboards: Home/Catalog/Detail/Cart/Review/Pending/Active/Partial-Replacement/History. **Seven critical operational screens** have1440 light/dark and1024 light: Dashboard/Queue/Review/Direct/Return/Inventory/Borrower detail. Remaining responsive behavior is explicit per composition, not every possible viewport. Dark at additional widths inherits the same documented geometry/tokens; those extra dark permutations are not separately rendered.

## 31. Dark/light coverage

All34 compositions have base light and dark SVG pairs (**68**). Additional18 borrower responsive +7 operational compact +6 lifecycle/home light variants = **99 SVGs**. Ten PNG previews: four borrower/Staff theme contact sheets and six readable selected full-page captures. Every contact/address and equipment graphic is synthetic/generic; no production image dependency.

## 32. Accessibility review

44 computed opaque token contrasts passed proposed text/control targets, including muted/status/disabled/selected/accent/danger/focus pairs. Strong control borders distinguish inputs in both themes; status has text+icon; touch target goal44×44; focus2px/offset2; persistent labels/errors and readable tables specified. Selected rendered light/dark/mobile/tablet/desktop examples and contact sheets visually inspected; label/footer fixes applied.

This is not WCAG certification. No interactive keyboard/screen-reader/focus-trap/zoom320/text enlargement/real mobile keyboard/safe-area/deployment check ran. Local browser used solely to rasterize static SVG review artifacts, not to run the SPA or runtime suites. Initial browser sandbox socket restriction was resolved by automatically approved isolated rendering; no requested work remains blocked.

## 33. shadcn composition mapping

**32 planned composition mappings** in [COMPONENT_COMPOSITION](../ux/COMPONENT_COMPOSITION.md); actual retained inventory **62 UI TSX files**. Button/Card/Badge/Field/Input/Table/Tabs/Dialog/AlertDialog/Sheet/Drawer/Calendar/Popover/Select/Combobox/Command/Skeleton/Empty/Alert/Progress/Sonner/Sidebar and Lucide suffice. Feature boundaries and current session architecture retained; no custom TimePicker/library replacement assumed. No component implementation added.

## 34. Sarcita desktop-first pitfalls avoided

390 borrower complete task is source. No persistent right summary required, five-column time grid, two-column base form, oversized hero or desktop-only navigation.1→2→3 catalog columns enhance discovery; desktop panels add no exclusive actions. This addresses the brief's specific lessons, without claiming a fresh audit of unrelated Sarcita source.

## 35. Remaining visual questions

Owner review: evergreen tone, information density, table/card fallback, tablet bottom nav, compact operational rail, label/confirmation wording, category taxonomy, final terms copy, generic-to-real catalog images, and Due today wording for the brief's due-soon summary. No additional due-soon threshold or policy chosen. Activation/template/import contract, notification projection/provider, export format and exceptional correction mechanics are later scoped dependencies, not new core-policy blockers.

## 36. Mockup review checklist status

[MOCKUP_REVIEW_CHECKLIST](../ux/MOCKUP_REVIEW_CHECKLIST.md) delivered with borrower/operational/visual/content review items, reviewer/date/revision/change log and separate approval/authorization record. **All owner approval checkboxes unchecked.** Designer artifact checks do not count as owner approval.

## 37. Phase 3B readiness

Reviewable visual specification and static artifacts ready for owner review and future implementation planning. Production implementation **not authorized / NOT STARTED**. Owner visual approval and separate explicit Phase 3B authorization required; later feature/domain contracts are not supplied by mockup creation.

## 38. Files changed

Added five Markdown documents: this report and `docs/ux/VISUAL_SYSTEM.md`, `HIGH_FIDELITY_MOCKUPS.md`, `COMPONENT_COMPOSITION.md`, `MOCKUP_REVIEW_CHECKLIST.md`. Added `docs/ux/mockups/`:99 SVGs,10 PNGs,1 JSON design manifest (**110 design assets**). Updated only tracked `docs/project/ROADMAP.md`. Total **115 new files +1 modified tracked file**. The generator/capture scripts remain temporary tooling in `/tmp`, not production code.

All320 previously tracked V2 file fingerprints compared against the starting clean baseline; only ROADMAP changed. All588 previously fingerprinted boss tracked files compared unchanged. No Go/React/TypeScript source, manifests/dependencies/security files, SQL/migrations or existing domain/UX documents changed. No commit/push/deploy/remote change. No private contacts/media copied; `.invalid` examples and explicit outside-software return-photo boundary retained.

## 39. git diff --check

Validation results recorded after final artifact review: `git diff --check` exit0. New Markdown/SVG/JSON text also checked for trailing whitespace/newline issues because unstaged untracked files are outside ordinary git diff. Browser text-bound audit found no out-of-artboard text across all99 SVGs. Manifest/99 SVG dimensions, light/dark pairing, required widths, all 59unique surface mappings, local links, safe SVG content and PNG signatures/dimensions checked. Only design/documentation paths changed.

Per current instruction, backend/frontend runtime suites, Go formatting/vet/tests, frontend lint/test/build, database/Docker/deployment checks **not run** for this design-only phase. Rasterization is static artifact verification, not application/browser-runtime acceptance.

## 40. Exact git status

Final `git status --short`:

```text
 M docs/project/ROADMAP.md
?? docs/project/PHASE3A2_REPORT.md
?? docs/ux/COMPONENT_COMPOSITION.md
?? docs/ux/HIGH_FIDELITY_MOCKUPS.md
?? docs/ux/MOCKUP_REVIEW_CHECKLIST.md
?? docs/ux/VISUAL_SYSTEM.md
?? docs/ux/mockups/
```

Starting worktree was clean; Git collapses the110 new assets into the untracked directory line. None staged.

## 41. Whether Phase 3A.2 design-production gate is satisfied

**YES — DESIGN PACKAGE COMPLETE.** Coherent identity, semantic light/dark tokens, mobile-first borrower set, operational set, required responsive geometry, correct return/replacement/fine/kiosk semantics,59-surface trace, future primitive map, shared state examples, owner review checklist and explicit desktop-first lessons exist. Scope verification and diff check pass. No production implementation occurred. Phase remains CURRENT while approval is pending.

## 42. Whether owner visual approval is still required before Phase 3B

**YES.** No Phase 3A.2 palette/mockup approval is inferred from artifact existence. Current owner direction mandates visual review before Phase 3B, which also needs separate authorization. No production code or another phase is started by this report.
