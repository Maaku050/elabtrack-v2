# Phase 3A.1 — UX Architecture & Low-Fidelity Wireframes report

**2026-10-08, Asia/Shanghai. COMPLETE — requested UX/text-design deliverables and documentation exit gate.** Package is a low-fidelity proposal for owner review, not approved high-fidelity mockups or implemented functionality. Phase 3A.2 is ready for separate authorization and NOT STARTED; Phase 3B remains NOT STARTED. No locked Phase 2.5 policy changed.

## 1. Sources reviewed

Owner's current Phase 3A.1 brief; root AGENTS; [charter](PROJECT_CHARTER.md), [source hierarchy](SOURCE_OF_TRUTH.md), [decisions](DECISIONS.md), [open questions](OPEN_DECISIONS.md), [roadmap](ROADMAP.md), historical [Phase 2 report](PHASE2_REPORT.md), current [Phase 2.5 report](PHASE2_5_REPORT.md). All eight domain documents: DOMAIN_MODEL, BUSINESS_RULES, STATE_MACHINES, DATA_MODEL, INVARIANTS, API_RESOURCE_DRAFT, V1_COMPATIBILITY, BOSS_REBUILD_REFERENCE. Working source already reviewed in prior phases remains evidence; owner correction explicitly excludes software return photos. No capstone personal/sample data is copied.

Read-only frontend inspection: App/router/providers/Query/session bootstrap and foundation pages, UI theme store, use-mobile, components.json and62 ui TSX-file inventory, with representative Sidebar/form/dialog/sheet/calendar/combobox/empty primitive availability. Current router is only foundation `/`, `/status`, wildcard; no product shell implemented. Stack/manifests remain React/TypeScript/Vite/Tailwind v4/shadcn Base UI/Lucide and preserved Go/PostgreSQL foundation; no package/library change or runtime verification.

## 2. UX roles

BORROWER, STAFF, ADMIN only; Student/Faculty categories share Borrower flow. Admin inherits Staff operations plus privileged administration. Current Phase 1 generic user/admin source is unchanged; future approved implementation reconciles roles through server authorization. No student-only navigation, separate eligibility or new Super Admin.

## 3. Borrower navigation model

Mobile compact app bar and four labeled bottom destinations: **Home, Equipment, My Borrowings, Account**. Notifications secondary in app bar/account; Selected equipment/cart contextual within catalog. Rationale: direct access to finding equipment and following obligations without crowded tabs/sidebar. Tablet keeps primary model; desktop adds persistent same-label navigation and wider composition. See [architecture](../ux/UX_ARCHITECTURE.md).

## 4. Staff/Admin navigation model

Desktop/tablet one-level sidebar **Dashboard, Requests & Borrowings, Inventory, Borrowers**; operational loan tabs Pending/Active/Overdue/Replacements/History. Admin adds Reports and Administration (Accounts/Terms & Policy/Audit). Fine clear contextual on account/loan panels, bulk import on Borrowers; Staff has no privileged controls. Operational activity and account/theme/sign-out secondary. Two separate readable Mermaid navigation maps created.

## 5. Total screen inventory

**59 inventoried reviewable UX surfaces =2 shared +22 Borrower +23 Staff/Admin operational +12 Admin-only.** Pages, steps, sheets/dialogs, shared detail variants and sections are counted explicitly; this is not 59 implemented pages/routes. Admin inherits operational surfaces without duplicate count. All 59 mapped to38 wireframe groups in [inventory](../ux/UX_ARCHITECTURE.md#screen-inventory-and-counting-rule).

## 6. Borrower screen inventory

B-01–22: first-use/updated terms, Home, catalog, search/filters, equipment detail, selected-equipment cart, request review, required due date/time, submitted confirmation, Pending/Denied/Expired/Cancelled detail, Active borrowing, partial return, replacement status, overdue/fine, Completed detail, My Borrowings overview, history, notifications/activity, profile/account. SH-01 sign-in/SH-02 access/session recovery shared. No public signup or software return-evidence surface.

## 7. Staff/Admin screen inventory

O-01–23: dashboard/queue/review/denial, loan list/detail, four direct-issue steps, return form/review-result, replacement select/review-result, Borrower list/detail/provisioning, inventory list/detail/edit/ordinary movement/inactive-archive, operational activity. A-01–12: full fine clear/history, three bulk-import stages, deactivate/reactivate, Staff/Admin list/create-edit, policy/terms/version review, reports, audit. Specific exceptional correction section is Admin-only within stock review; no unapproved damaged-original disposal pipeline.

## 8. Mobile-first decisions

Borrower390px schematic base,320px minimum target; one-column content, labeled touch controls and explicit next action. Catalog1 column mobile,2 tablet,3 desktop (4 only when card width permits); mobile history cards rather than shrunken desktop table. Filter/full-height editing sheets, due date+time stack, one primary action region clear of bottom nav/safe area/keyboard. No borrower mobile sidebar or wide-first forms.

## 9. Desktop/tablet decisions

Staff/Admin operational sidebar, bounded tables, contextual accountability/summary side panels and readable prior-versus-now quantity fields; compact rail/drawer tablet. Mobile fallback stacks item cards and review panels, full-height sheets, keeps actions accessible. Dense ledger/report/import preview may have labeled contained scrolling; no promise every dense analysis is optimal on phone.

## 10. Catalog/kiosk interpretation

Interactive Kiosk is the same responsive borrower browse/search/filter/quantity/selected-equipment/review/request experience. No device registration/tokens/shared-device auth/handoff, hardware touchscreen shell, separate kiosk app or commercial Buy/Pay wording. Equipment catalog images independently supported; internal stock buckets hidden from Borrower.

## 11. Borrow request flow

Select unique equipment/quantities → required requested due date AND time Asia/Manila → review current-version acceptance/items/face-to-face instruction → Submit → actual committed Pending and immediate reservation → visit FSMO. Cart only local intent before commit; old fine doesn't automatically block. No 7-day maximum, per-loan checkbox, silently clipped quantity or claim of approval on submit. [Journey A](../ux/BORROWER_FLOWS.md#journey-a--browse-request-visit-fsmo-receive-equipment).

## 12. Pending/expiration flow

Pending shows reserved equipment/submission/expiry/requested due and advisory countdown; own Cancel confirms released hold/history retained.24h server expiry produces Expired/history/new-request action, distinct from denial. Timer0 shows checking until actual server response, never client-created release/state. Concurrent staff/expiry changes refresh before actions. Journey B and WF-08/09 cover these.

## 13. Approval/release flow

Queue opens review; no row Approve shortcut. Review borrower/category/active/terms bound to submission/held quantities/due/accountability. Staff human review may approve/deny; old fine not an automatic prohibition. Confirm physical handover now → one PENDING→CHECKED_OUT, user label Active. No Approved-waiting or second Release step; full-request issue selected domain baseline, no silent quantity editing. WF-16–18.

## 14. Direct checkout flow

Journey C: select active existing Borrower → select current equipment quantities → required due date/time → review → physically issue. Current terms evidence, availability/account/actor validated by server later. Immediate Active/CHECKED_OUT from available stock, no pending/24h reservation. WF-19; generic Borrower, not Select Student.

## 15. Return flow

Staff/Admin only: prior issued/good/damage/loss, physically outstanding, replacement outstanding; now-good/damage/loss positive totals<=remaining. Full-good shortcut fills physical only. Review exact effect before recording; immutable outcome shows partial/open or complete iff allP/U0. Borrower read-only, no generic Complete/Mark Returned. Personal phone photo viewing is entirely outside eLabTrack; no capture/upload/storage/transmission/attachment/evidence step for any role. WF-20/21, journeys D/E.

## 16. Replacement-obligation flow

Staff records damage/loss, physical quantity leaves checked-out, exact replacement owed shown; no price charge. Select unique obligation/accepted-now<=remaining/correct type or accepted equivalent → review usable stock acquisition and remaining → authoritative acceptance. Retained damaged originals remain nonusable; original incident survives. Partial acceptance open; last resolution closes only all physical/replacement0. WF-22, Journey F; original disposal later policy, not a required UI step.

## 17. Overdue/fine flow

Overdue when unresolved physical OR replacement and now>original due; PHP 10 per ceiling elapsed24h, live as-of amount until Completed freezes final.1min10PHP/24h10PHP/25h20PHP. Fine read-only Borrower/Staff; Admin full outstanding Paid/Waived/Other, optional note, named history and stale-balance confirmation. No partial/online payment, no assessment erasure, no clearing fine to close equipment liability. Open loan full-clear can accrue new delta later; final completed amount frozen. Journey G/WF-10/11/30.

## 18. Bulk import flow

Journey H/WF-31: download approved template → choose account-data file → validate rows → valid/invalid/duplicate preview → explicit valid-row subset choice and confirmation → authoritative created/not-created/unknown result. Admin only, role forced Borrower. Format, batch atomicity/duplicate rules and secure activation delivery are later implementation dependencies; no permanent plaintext-password column or fake delivery claim. Unknown result checks original batch before new import; this file input never accepts return evidence.

## 19. Loading states

Region-shaped Skeletons for catalog cards, loan summaries, queue/return data/identity lists/detail; navigation retained. Background Updating/as-of retains usable old view, not false0. Authority-sensitive actions depend on loaded required state; no full-page spinner as default, no optimistic physical stock/completion. Global WF-37 plus per-group notes.

## 20. Empty states

Explicit11-case matrix: no equipment, no search matches, empty cart, no pending, no active, no history, no overdue, no replacements, no Borrowers, no report results, no activity updates. Role-appropriate Browse/Clear filters/Provision/Direct issue/back actions; error distinguished from real emptiness. No inactive/forbidden creation control disguised as empty-state CTA.

## 21. Error states

Network read retry, associated field validation/errors, server-safe unexpected message/request-ID support details, inactive account and session reauthentication, ordinary403/404 return/back without new global logout. Input preserved within account/session boundaries; unknown write checks original outcome before resend. Raw backend/SQL/provider/credential messages never user copy. Phase 1 behavior unchanged.

## 22. Conflict states

Changed availability annotates affected equipment; no silent smaller request. Already processed/expired request shows actual state and removes obsolete actions; returns/replacements refresh remaining and require review, not auto-clipping/resubmit. Stale fine balance forces new whole-balance confirmation. Pending timer never decides backend state, and ordinary conflicts never cause fabricated success.

## 23. Status system

[STATUS_SYSTEM](../ux/STATUS_SYSTEM.md): Pending / Active borrowing / Denied / Cancelled / Expired / Completed map exact lifecycle. Overdue/Partial return/Replacement required and fine Outstanding/Cleared derived separately; hierarchy lifecycle→urgency→obligations. Completed may have fine outstanding; no current Overdue after operational completion. Labels/text independent of final colors, theme or icons. Physical stock labels A/R/C/damaged_held/T remain distinct from liability/incidents.

## 24. Responsive rules

[RESPONSIVE_RULES](../ux/RESPONSIVE_RULES.md): mobile320–767, tablet768–1023, desktop>=1024 proposals align existing768 mobile boundary. Concrete per-major-surface matrices specify stacking/nav/card grid/cart/filter/input/calendar/dialog/table-card/scroll/action behavior, keyboard/safe-area handling and wider layouts. Breakpoints are UX engineering assumptions, not code changes or runtime-proven behavior.

## 25. Accessibility baseline

Semantic headings/landmarks, proper visible labels, accessible icon names/alt/fallback, linked error summary and field errors, logical keyboard order/visible focus, modal trap/restore/Escape, touch target44×44 design goal,320px/zoom/reduced-motion behavior, semantic status text and nonspamming announcements. Both themes meaningful without color. No WCAG certification/accessibility runtime result claimed.

## 26. shadcn mapping

Each38 wireframe group names likely available future primitives: Sidebar, Button/Card/Badge/Item, Field/Input/Textarea/Select/Combobox/Command, Table/Tabs/Pagination, Calendar+typed date/time Input, Popover/Sheet/Drawer/Dialog/AlertDialog, Skeleton/Empty/Alert/Progress, supplementary Sonner. No missing TimePicker or bespoke production component is assumed;62 primitives/components.json untouched. Attachment/showcase inventory doesn't authorize a return feature.

## 27. Wireframes created

**38 numbered text wireframe groups**, all 59 surface IDs explicitly mapped in [WIREFRAMES](../ux/WIREFRAMES.md). Covers all required borrower mobile home/catalog/detail/cart/review/pending/active/partial-replacement/history/profile and desktop dashboard/queue/review/direct/return/replacement/borrower/inventory/fine/bulk plus terms/shared/Admin secondary surfaces. Multi-step and lifecycle variants intentionally share reviewable layouts. No polished mockups/image assets/HTML prototype.

## 28. User journeys created

A browse→reserve→visit→issue; B24h expiry; C direct issue; D full good return; E partial→final; F lost/damaged→replacement acceptance; G overdue growth→completion→Admin full clear; H Admin bulk provision. [Borrower flows](../ux/BORROWER_FLOWS.md) and [Staff/Admin flows](../ux/STAFF_ADMIN_FLOWS.md) document quantities/clock/history/permission and adverse results; journey names reused for cross-review, not duplicated independent implementations.

## 29. UX acceptance checklist

[UX_ACCEPTANCE_CHECKLIST](../ux/UX_ACCEPTANCE_CHECKLIST.md) creates pending per-screen mockup/role/viewport/loading/empty/error/conflict/accessibility/theme/session/status checks plus feature/journey/adverse cases. All boxes remain unchecked: future implementation QA only, not tests run here. Requires actual domain/DB/HTTP evidence when implemented; visual fit cannot substitute business integrity.

## 30. Open UX questions

UX-01 navigation/cart labels;02 selected full issue/live-clear/current-terms assumptions;03 owner-scoped notification/activity projection without unread/delivery inventions;04 required requested-date/time input (confirmed phase requirement), field copy;05 bulk format/explicit subset/batch/retry/activation mechanics;06 terms/taxonomy/image content and read-only initial policy settings;07 report formats/export/privileged recovery details. Each has working design assumption/review boundary in architecture, not new institutional rules. Minor usability/content questions need owner review, not a new stakeholder packet.

## 31. Whether any domain-policy blocker was discovered

**NO new core UX domain-policy blocker.** No locked policy altered. Remaining original-damaged disposal/equivalence guidance, institutional retention, provider/cadence, migration access/reconciliation and deployment gates still apply to their later activities. Secure activation/import/projection/export/recovery mechanics need implementation design but don't block hierarchy/low-fidelity or separately authorized high-fidelity styling. No disposed original or software photo feature invented.

Major UX risks explicitly designed around: reservation mistaken for issue; expiry countdown mistaken for server release; replacement-only loan mistaken for completion; fine clear mistaken for frozen/closed active loan; hidden Staff/Admin privilege differences; stale write/unknown result duplicated; inventory total double-counted; tiny mobile return matrices; notification update mistaken for delivery proof; bulk preview assumed silent partial success.

## 32. Phase 3A.2 readiness

**YES — ready for separate authorization; NOT STARTED.** Review navigation/copy/assumptions, then authorized high-fidelity layouts and both-theme visual states can follow without changing working rules. This report does not approve mockups or authorize Phase 3B; parent Phase 3A complete mockup gate still pending.

## 33. Files changed

**9 Markdown files:7 new UX documents +1 new phase report +1 tracked roadmap update.** No domain/decision/Phase 2.5 source policy edits.

- [UX_ARCHITECTURE](../ux/UX_ARCHITECTURE.md)
- [BORROWER_FLOWS](../ux/BORROWER_FLOWS.md)
- [STAFF_ADMIN_FLOWS](../ux/STAFF_ADMIN_FLOWS.md)
- [RESPONSIVE_RULES](../ux/RESPONSIVE_RULES.md)
- [WIREFRAMES](../ux/WIREFRAMES.md)
- [STATUS_SYSTEM](../ux/STATUS_SYSTEM.md)
- [UX_ACCEPTANCE_CHECKLIST](../ux/UX_ACCEPTANCE_CHECKLIST.md)
- [ROADMAP](ROADMAP.md)
- This [PHASE3A1_REPORT](PHASE3A1_REPORT.md).

## 34. git diff --check

**PASS — `git diff --check` exit0, no diagnostics.** The tracked roadmap diff passed; all eight untracked Markdown files were explicitly checked for final newline/trailing whitespace, local Markdown link targets/anchors, balanced code fences and consistent inventory/wireframe/report counts.

Before/after SHA-256 fingerprints of all312 original tracked V2 files show only ROADMAP.md changed. All Go/React/TypeScript/SQL/source/config/session/UI-primitive files, existing domain/Phase2.5 documents and manifests unchanged. Eight new files are Markdown only: seven UX documents and this report. Boss tracked-file fingerprints unchanged; no boss writes. `.project-reference/` remains ignored.

59 unique surface IDs (2 shared/22 Borrower/23 operational/12 Admin-only) map to38 numbered wireframe groups; two separate navigation maps and journeys A–H present;36 ordered report sections; future QA checkboxes unchecked. No high-fidelity assets, runnable prototype, React screens, components, migrations or implementation added.

Runtime Go/frontend suites, database/browser/Docker/deployment checks **not run**, as requested for documentation-only design. No rendered responsive/accessibility/theme or business execution results claimed; acceptance checklist is future work.

## 35. Exact git status

Starting working tree was clean. Exact final `git status --short` (default untracked-directory grouping):

```text
 M docs/project/ROADMAP.md
?? docs/project/PHASE3A1_REPORT.md
?? docs/ux/
```

The grouped `docs/ux/` contains exactly the seven new UX files listed in section33; no other untracked files. No commit/push/deployment performed.

## 36. Whether Phase 3A.1 exit gate is satisfied

**YES — documentation design gate.** Borrower mobile-first IA/flows, Staff/Admin desktop/tablet responsive IA, corrected Interactive Kiosk, complete59-surface inventory, request/physical issue/expiry/direct/returns/replacements/Admin fine/bulk wireframes, two navigation maps, journeys A–H, concrete width behaviors, loading11 empties/errors/conflicts/confirmations/statuses,38 mapped text layouts, future shadcn mapping and unchecked QA checklist satisfy the requested scope. No high-fidelity/production implementation and no Git publication/deployment. Owner review and separate authorization precede Phase 3A.2; no future phase marked complete.
