# Phase 3B — Approved Visual Design System & Application Shell

2026-10-08. **PHASE 3B — COMPLETE AND OWNER VISUALLY APPROVED.** The owner reviewed real implementation screenshots and confirmed faithful reproduction of the approved mockups. All four previews, navy/violet FSMO identity, compact branding, shadcn foundation, light/dark themes, mobile responsiveness and desktop operational layouts are accepted, including the disclosed minor adaptations. The 32 remaining production screens are not implemented; their approved references and accessibility requirements remain binding. DEC-064 records approval; Phase 4 is not started or authorized. Current environment/checkpoint evidence is in [LOCAL_ENVIRONMENT_READINESS_REPORT](LOCAL_ENVIRONMENT_READINESS_REPORT.md).

## 1. Approved visual references verified

The approved README and manifest were read before code changes. The standalone seal, B01, B02, both S01 frames and every one of the 32 individual remaining PNGs were opened and visually inspected. All required files were readable. [Reference package](../ux/approved/README.md), [verification records](../ux/verification/phase3b/ASSETS.json).

## 2. Visual asset manifest status

All 36 manifest files match their supplied SHA-256 fingerprints and dimensions. There are 35 screen PNGs and one brand PNG; the combined S01 image contains two screens. Consequently there are **36 screen targets: four baselines and 32 future screens**. The approved package, including README, manifest and overview images, remains unchanged.

## 3. Superseded evergreen design documented

The owner-approved navy/indigo/violet package supersedes the earlier evergreen proposal. VISUAL_SYSTEM, HIGH_FIDELITY_MOCKUPS, COMPONENT_COMPOSITION, MOCKUP_REVIEW_CHECKLIST, DEC-063 and ROADMAP record that authority. Original SVGs and the Phase 3A.2 report remain historical evidence, including their superseded approval wording. No historical business decisions were rewritten.

## 4. Final visual token system

The semantic system is defined in `frontend/src/styles/index.css`, with scoped composition rules in `application.css`. Light primary is `#2B10BB`, page `#F5F7FD`, text `#10132E`, surface white. Staff navigation uses `#061E52` to `#061741`, with `#384CE6` active navigation. Gold `#FFC529` is controlled; success, warning and danger have separate foreground/backing pairs. Essential input boundaries and focus rings are stronger than decorative card borders. Existing Tailwind v4/shadcn aliases remain compatible. See the complete [token table and sampled reference coordinates](../ux/VISUAL_SYSTEM.md).

## 5. Light/dark theme implementation

Both themes reuse the existing Zustand UI store. Dark page is `#0D1428`, surface `#17213A`, text `#EEF1FC`, primary `#B1A0FF` with dark text. The institutional navy sidebar remains. Actual theme-button switching passes at every recorded viewport. Dark is an authorized engineering translation of light references; no new approved dark PNG was supplied. The owner subsequently approved the real dark implementation screenshots. No replacement theme library or session persistence was added.

## 6. Typography system

The existing system sans font is retained. Borrower titles are 24px, Staff titles 27px, sections 18px, operational panels 16px, body/control text 14px, normal metadata 12–13px, counts 22–27px with tabular numerals. The catalog has subordinate 10px tags and a 7.5px institutional caption. These small subordinate labels and generated-font differences were included in the owner-reviewed, accepted baseline; no exact typeface or accessibility certification is claimed.

## 7. Spacing/layout system

Spacing tokens are 4/8/12/16/20/24/32px. Borrower phone edges are 14px, wide edges 24px; header 64px; bottom navigation 72px plus safe area. Card radius is 12px, control radius 8px, borders one pixel and shadows restrained. Staff rail is 238px at 1440, 210px at 1024 and 76px when collapsed; workspace margins are 24/20px. Controls use 44px targets. Readable text and interaction targets increase vertical content height relative to compressed source artwork; the fidelity matrices disclose those differences.

## 8. shadcn primitives reused

Retained Button, Card, Badge, Input, Checkbox, Dialog, Sheet, Sidebar, Table, Alert, Skeleton and Empty primitives support the actual compositions. Existing common empty/error/loading components remain available. All **62 primitive files**, `components.json` and dependency manifests are unchanged. Lucide provides icons; no starter, replacement component library, font or chart dependency was installed.

## 9. New shared components

Five files under `frontend/src/components/application/` provide AppBrand, AppButton, ThemeControl, SurfaceCard, PageHeading, SectionHeading, MetricCard, semantic/accountability/fine badges, equipment thumbnails/cards/quantity selection, table/search/pagination, loading/conflict feedback, read-only PreviewDialog and Borrower/Staff shell compositions. Only actual foundation/preview needs were implemented; future workflow components remain deferred. [Composition inventory](../ux/COMPONENT_COMPOSITION.md).

## 10. Borrower mobile shell

B01 governs the compact seal/divider/wordmark header, secondary notification access, content hierarchy and Home/Equipment/My Borrowings/Account bottom navigation. Safe-area variables, skip link, named controls, active destination and navigation clearance exist. Unsupported destinations open a clearly marked preview disclosure. The seal is a separately replaceable supplied AI reconstruction, not an authenticated official master. [Brand provenance](../../frontend/src/assets/brand/README.md).

## 11. Borrower tablet/desktop adaptation

The same four destinations and task order remain. Home metrics become four columns; the catalog uses one column at 320/390, two at 768 and three at 1280. Selected-equipment actions stay accessible. At 320, equipment controls move below the facts to preserve touch targets. The owner accepted these inferred wide compositions in the real implementation review; no unrelated desktop borrower navigation was introduced.

## 12. Staff/Admin shell

The shared shell has a navy sidebar, indigo active navigation, compact operational toolbar, search-shaped preview disclosure, account/notification controls and theme control. Staff sees Dashboard, Requests & Borrowings, Inventory and Borrowers. An explicit presentation-only Admin variant adds Reports and Administration. Rendering that variant does not authenticate anyone or grant a capability. Actual authority remains server-owned. Existing providers, transport, Query and auth/session stores are preserved.

## 13. Staff/Admin responsiveness

At 1440 the dashboard has four metrics and three operational panels; at 1024 the activity panel spans a second row. Below 1024 navigation uses the retained Sheet. Desktop collapse preserves named navigation and expands the workspace. Pending table scrolling is bounded, focusable and keyboard reachable at 1024; all columns and Review actions fit at 1440. Whole-page horizontal overflow was absent in every required capture.

## 14. Borrower Home fidelity

The real component preview follows greeting, 2×2 metrics, Browse Equipment, three activities, three quick actions, reservation notice and bottom navigation. Synthetic copy replaces illustrative names/dates/currency. The [B01 ten-row matrix](../ux/VISUAL_FIDELITY_CONTRACT.md#b01--borrower-home) records composition matches and differences in font metrics, density, preview disclosures and wide/dark adaptation. No pixel-perfect claim.

## 15. Equipment Catalog fidelity

The preview reproduces title/search/filter/sort, category chips, equipment pictures/facts/availability, quantity controls, Add actions and selected summary. Search, category, sort, availability, selection and filters affect only local fixtures. Review opens a read-only dialog. Six equipment-only sample crops were extracted from the approved catalog; no whole mockup renders as UI. [Catalog asset provenance](../../frontend/src/features/visual-preview/assets/README.md), [B02 matrix](../ux/VISUAL_FIDELITY_CONTRACT.md#b02--equipment-catalog).

## 16. Staff Dashboard fidelity

The S01 upper frame is reproduced as title/date, four counters, request trend, physical stock summary and activities. SVG charts use static fixture data and blue/lavender tokens with textual alternatives. Physical counts reconcile: 82 available + 28 checked out + 8 damaged held + 6 reserved = 124. Dates align with the explicit fixture snapshot. No backend metrics are fabricated. [Dashboard matrix](../ux/VISUAL_FIDELITY_CONTRACT.md#s01-upper--staff-dashboard).

## 17. Pending Requests fidelity

The S01 lower frame supplies the same shell, heading, status pills, search/filter toolbar, identity-first table, accountability context, Review and pagination. Eight synthetic rows paginate five per page; search, filter and checkboxes are local. Canonical lifecycle labels replace Approved/Released illustration errors. Review cannot approve, deny, issue or record a return. Initial clipped desktop Review actions were fixed through scoped cell spacing and compact fine labels. [Pending matrix](../ux/VISUAL_FIDELITY_CONTRACT.md#s01-lower--pending-requests).

## 18. Visual coverage map completeness

The [implementation map](../ux/APPROVED_VISUAL_IMPLEMENTATION_MAP.md) contains four baseline entries and all 32 individually named future screen entries, each with role, reference, structure, shared/unique components, responsive expectation, phase, implementation status and review status. A separate brand record covers the remaining manifest asset. All 32 future screens explicitly remain **NOT IMPLEMENTED**; shared shell/feedback coverage does not imply production screen completion. The fidelity contract also contains a 32-reference inspection/dependency ledger.

## 19. Screenshot QA results

Real Chromium **140.0.7339.16** rendered the previews through Vite. **24 viewport/theme cases pass**, plus **10 keyboard/boundary/production checks**. Evidence includes 24 viewport PNGs, two full-scroll mobile captures, a selection dialog and a navigation Sheet: **28 PNGs**. Actual captures were opened and compared with approved PNGs at equal normalized widths; S01 frames were reviewed separately. [Comparison viewer](../ux/verification/phase3b/comparison.html), [raw measurements/results](../ux/verification/phase3b/RESULTS.json), [reproduction instructions](../ux/verification/phase3b/README.md).

Application JS exceptions, unexpected console errors and browser warnings: **zero**. Business API requests: **zero**. The runner intercepts only existing session-bootstrap refresh and returns an anonymous 401; **31 expected 401 network diagnostics** are separately recorded. No token or fake role session is supplied. Production JS excludes fixture data and preview routing, production preview URLs show the existing 404, and the original foundation root works.

## 20. Responsive QA results

| Screen family | Widths | Themes | Result |
|---|---|---|---|
| Borrower Home and Catalog | 320, 390, 768, 1280 | Light and dark | 16/16 pass |
| Staff Dashboard and Pending | 1024, 1440 | Light and dark | 8/8 pass |

All captures use a 900px viewport height. Checks cover render/title, whole-page overflow, selected text clipping, control sizes, landmark geometry, actual theme switching and bottom-navigation clearance. The catalog dock clears navigation; scrolling makes final content reachable. Table content is 1152px inside a 1152px desktop region; tablet content scrolls inside its named region. Full-page Chromium PNGs retain fixed elements at their viewport position and therefore require the separate scroll-clearance measurements for obstruction assessment.

## 21. Accessibility checks

Semantic header/main/nav, heading order, skip links, visible focus, accessible control names, icon/text statuses, linked quantity errors, disabled bounds and textual chart alternatives exist. Browser checks pass for Sheet Escape/focus restoration, Dialog Tab trapping, smaller-width Staff navigation, table keyboard scrolling, sidebar collapse and reduced motion. **28 selected opaque semantic contrast pairs pass** 4.5:1 text or 3:1 essential boundary/focus targets. [Contrast calculations](../ux/verification/phase3b/CONTRAST.json).

These checks do not certify every primitive/overlay or WCAG conformance. Physical devices/notches, virtual keyboards, full browser zoom and complete future loading/error/mutation states were not exercised. The 320px capture establishes reflow evidence, not a claim of a 200% browser zoom run.

## 22. Component tests

Six new tests verify integer quantity bounds and linked errors, theme changes without session changes, canonical lifecycle labels, Staff/Admin navigation without authority grants, local catalog filtering/selection without submit/API calls, and local pending pagination/search with read-only Review. Tests assert critical boundaries rather than fixture snapshots. All pre-existing tests remain passing.

## 23. Frontend lint/tests/build

Node **24.19.0**, npm **11.17.0**, Vite **8.3.2**; no toolchain mismatch. Actual final commands:

```sh
npm --prefix frontend run lint
npm --prefix frontend run test:run
npm --prefix frontend run build
```

Lint exits 0; **145 tests in nine files pass**; TypeScript and production build pass. Final build reports 2157 transformed modules with no build warning. The 24-case Chromium run and 10 extra checks also pass. The final run includes the named loading fallback for lazy preview routes and explicit browser warning collection.

Backend code is unchanged. Full Go, PostgreSQL integration, Docker and deployment checks were deliberately not run in this visual-only phase, as requested. No real production domain workflow was tested or claimed.

## 24. Existing warnings

Frontend lint retains **19 baseline warnings**: 17 `only-export-components` warnings in retained primitive files and two `set-state-in-effect` warnings in the existing carousel/mobile hook. Their files remain unchanged. No rule, type check or test was weakened to pass the gate.

## 25. New warnings/errors

Final added application code has **zero new lint warnings or errors**; final build has zero warnings. Initial browser findings involving Card direction, Link/native-button semantics, quantity sizing, table clipping, navigation padding, chart palette and the lazy-route hydration loading warning were corrected before the recorded final run. Browser warnings are collected explicitly alongside errors. Expected anonymous refresh network diagnostics are disclosed in section 19 rather than misreported as application errors.

## 26. Domain conflicts corrected in sample UI

Generated copy was corrected to respect one approval/physical-release transition; Pending reserves rather than checks out; Staff/Admin alone records returns; physical stock and replacement obligations stay distinct; fines use PHP and the existing PHP10 per started overdue 24-hour rule. Historical fine context is not a pending-request fine or automatic block. Staff navigation omits Admin-only destinations. Borrower return-photo upload/evidence/attachments are absent; catalog imagery is independent. No domain document, backend contract, state machine or arithmetic implementation changed.

## 27. Unresolved visual deviations

The four ten-row matrices disclose system-font differences, readable text/control sizes, increased scroll height, 64px Staff header versus the compressed reference, small sidebar width differences, role-correct Staff navigation, local preview disclosures and inferred wide/dark layouts. The owner reviewed and accepted the subordinate catalog/brand labels and operational density; accessibility remains an active requirement. Exact generated typeface and an authenticated original seal are unavailable. The owner has now **accepted these disclosed minor adaptations** after inspecting real implementation screenshots. This acceptance preserves accessibility, readable type/control sizes and future individual PNG gates; no pixel-perfect claim is made. [Owner review checklist](../ux/MOCKUP_REVIEW_CHECKLIST.md).

## 28. Files changed

Tracked files changed by this phase: `docs/project/DECISIONS.md`, `docs/project/ROADMAP.md`, `frontend/src/app/router.tsx`, `frontend/src/styles/index.css`.

New implementation/evidence groups: five shared application composition files; six preview modules plus component tests, six catalog-only PNGs and their README; separate brand PNG/README; `application.css`; two standard-library Chromium QA scripts; verification artifacts and this report. Four already-untracked visual documents were rebaselined as explicitly requested; the implementation map and fidelity contract were added.

At the original Phase 3B delivery, the Phase 3A.2 report, approved package, old mockups and handoff ZIP were pre-existing untracked files, not newly authored Phase 3B changes. The closeout preserves these files and excludes the redundant ZIP from Git; its current status is recorded in the local readiness report. A before/after audit checks all 320 baseline tracked files and all 249 protected files: only the four allowed tracked files changed, **zero protected files changed**. Protected scope includes backend, domain documents, auth/session, providers/Query, API_CONTRACTS, UI primitives, manifests/configuration and approved assets. [Scope evidence](../ux/verification/phase3b/SCOPE.json).

## 29. git diff --check

`git diff --check` passes. New/rebaselined text files were also checked for trailing whitespace and excess EOF blank lines because untracked files are outside that Git check. Documentation links and 36-target/32-future/matrix coverage were checked. No commit, push, remote change or deployment was performed.

## 30. Exact git status

Historical Phase 3B delivery `git status --short` snapshot, including preserved pre-existing work. The current closeout status is in the local readiness report:

```text
 M docs/project/DECISIONS.md
 M docs/project/ROADMAP.md
 M frontend/src/app/router.tsx
 M frontend/src/styles/index.css
?? docs/project/PHASE3A2_REPORT.md
?? docs/project/PHASE3B_REPORT.md
?? docs/ux/APPROVED_VISUAL_IMPLEMENTATION_MAP.md
?? docs/ux/COMPONENT_COMPOSITION.md
?? docs/ux/HIGH_FIDELITY_MOCKUPS.md
?? docs/ux/MOCKUP_REVIEW_CHECKLIST.md
?? docs/ux/VISUAL_FIDELITY_CONTRACT.md
?? docs/ux/VISUAL_SYSTEM.md
?? docs/ux/approved/
?? docs/ux/mockups/
?? docs/ux/verification/
?? eLabTrack_V2_APPROVED_VISUAL_HANDOFF.zip
?? frontend/scripts/
?? frontend/src/assets/
?? frontend/src/components/application/
?? frontend/src/features/visual-preview/
?? frontend/src/styles/application.css
```

## 31. Whether the Phase 3B exit gate is satisfied

**Yes — COMPLETE within the authorized scope.** Reference verification, rebaseline, both themes, preserved shadcn/security foundation, shared shells, four rendered/compared previews, responsive/accessibility/component/frontend gates, complete future map, disclosure and whitespace requirements are met. The owner has approved the real implementation and accepted the disclosed minor adaptations (DEC-064). Completion is not an assertion that 36 production screens work, nor authorization to advance the roadmap.

## 32. Whether owner review is needed before remaining production screens

**The current baseline review is complete.** Owner review of future material visual changes and relevant feature content/contracts remains required; the 32 future production screens are not covered by implementation approval. Each later screen must open its matching approved PNG, follow the map/fidelity contract, implement reviewed server-owned behavior and record fresh screenshot/integrity tests under explicit phase authorization. **Phase 4 has not begun.** No migrations, real borrowing/return/fine-clearance service, boss port, commit, push or deployment was performed.
