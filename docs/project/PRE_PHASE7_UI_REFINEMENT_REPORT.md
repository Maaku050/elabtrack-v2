# Pre-Phase 7 UI refinement report

Verified 2026-10-09, Asia/Shanghai. **ENGINEERING VERIFICATION: PASS. OWNER ACCEPTANCE: PENDING.** The owner explicitly authorized this UI refinement after reviewing the local application. Phase 7 remains unstarted; this report does not authorize another phase.

Review the [before/after image table](../ux/verification/pre-phase7-ui/README.md) or open the [standalone comparison](../ux/verification/pre-phase7-ui/index.html) in a browser. Review the running application at `http://localhost:5173/login` with the existing local Admin account. No owner password was requested or used in this sprint.

## 1. Navigation flicker root cause

The application already uses React Router links, and the session provider already lives above the router. Ordinary navigation did not reload the document. Two composition problems caused the visible blink: the route guard replaces its outlet with a loading state during each route-specific `/auth/me` request, and operational pages independently instantiate `StaffShell`. The shared chrome therefore disappeared during revalidation and was recreated with the next page. Staff and Admin route groups also had separate guarded parents with no common operational layout.

The [before observations](../ux/verification/pre-phase7-ui/before/acceptance.json) show different sidebar/header DOM objects across all 18 repeated navigations, 39 observer batches in which the original chrome was disconnected, zero document requests and 18 account-check requests. The observer count is evidence of disconnection, not a count of separate reloads or frames.

## 2. Exact navigation fix

`OperationalLayout` is one pathless parent of both existing operational role-guard groups. It renders the sole operational `StaffShell` and a nested outlet. Account and inventory pages, including generic workspace pages, now render content under that parent. The sidebar, top bar and main container remain mounted as the nested guard/page changes.

The existing guard still calls the existing route-keyed `useCurrentUser(location.pathname)`. Its role, active-account, identity, pending/fetching and error checks are preserved. Loading, denial and retry states render inside the content region when the operational shell is present. Borrower/anonymous experiences retain their appropriate standalone boundary. The shell is presentation metadata; it cannot authorize a page or an API request. Cached page content is withheld during access revalidation. No fade, animation delay or cached-authority shortcut conceals the issue.

The centralized transport, bootstrap provider, authentication endpoints, refresh rotation, cross-tab coordination, private query cleanup, safe return destinations and focus/reconnect revalidation are unchanged. Sign-out continues through the existing hook, including immediate local clearing and retryable server acknowledgement failure.

## 3. Sidebar overflow root cause

The shared operational footer rendered `Sign Out` as a bare text node, while collapsed CSS only hid a child `span`. Its text therefore remained beside the icon and expanded beyond the 76px rail. Independently, the sidebar content used an unbounded minimum height; navigation and the footer could exceed short viewports.

## 4. Sidebar alignment and footer fix

The sole footer now wraps its text in a span and retains `aria-label` and a title in both modes. Expanded controls span the available rail width. Collapsed controls center the icon, hide the label and constrain width/minimum width. Theme, sign-out and avatar share the rail alignment. The seal remains square, undistorted and centered with the navigation icons; no page-specific offsets were introduced.

The sidebar content has a bounded viewport height. Navigation owns independent vertical scrolling, and the footer cannot shrink or overlap it. Brand, navigation groups, active/hover/focus states and footer spacing use one shared composition. Native destination/control titles and accessible names remain available when labels are hidden. The mobile sheet uses the same footer and an expanded navigation layout.

## 5. Top-bar changes

A small `FSMO › section` context replaces the input-shaped workspace label. The top bar retains the navigation toggle, truthful future-notifications destination and account identity. Account actions now use the retained shadcn dropdown primitive with Your account and the existing Sign Out action. Keyboard menu navigation, Escape and trigger focus restoration were verified. Mobile context and account copy adapt to available width; large page headings remain in the content region. No notification counts or global-search feature were invented.

## 6. Inventory design improvements

The heading/actions, real stock totals and filter/results panel have separate spacing and hierarchy. Category management and Add Equipment retain their original routes and authority. Search has a visible operational label; category/status/sort controls use an adaptive grid, and Available only occupies a clear row. Metric cards display the original API totals with restrained styling. The empty state, populated table and pagination share the results panel. The wide table scrolls internally and has a labeled, keyboard-focusable region.

State setters, filter values, category paging, query arguments, page size, equipment actions and stock arithmetic are unchanged. Supplemental populated screenshots show actual synthetic quantities persisted through the isolated API; the desktop comparison uses the empty state captured before test equipment was created. No fabricated stock values were added to the application.

## 7. Dashboard empty-state improvements

The Dashboard now has a restrained introduction explaining that insights arrive in a later phase, followed by links to the implemented Inventory, Borrowers and Account workspaces. It contains no dashboard statistics, activity data, chart data or reporting implementation. The later-feature notice remains explicit.

## 8. Sign-in redesign

Before captures were made and inspected before editing source. The desktop composition combines a flat navy FSMO brand panel with a focused form; mobile uses a single column and compact brand introduction. The supplied seal, system typography and approved semantic navy/violet palette remain in use. Inputs have refined spacing and invalid states, and the primary action is full width. Both themes were rendered.

Existing React Hook Form/Zod validation, field IDs, labels, password-manager autocomplete, password visibility, busy-state disabling, error messages, password clearing/focus after failure, API integration, return navigation and logout feedback remain intact. No public registration, unimplemented password-recovery link, institutional SSO or demo credentials were introduced.

## 9–10. Responsive and theme verification

| Surface | Widths | Themes / modes | Executed checks |
|---|---|---|---|
| Sign-in | 320, 390, 768, 1024, 1366, 1440 | Actual light and dark | Document overflow, labels, password autocomplete, minimum 44px form controls |
| Dashboard and Inventory | All six widths | Light/dark; both desktop rail modes at 1024+ | Document/rail overflow, seal alignment/shape, footer separation, contained collapsed Sign Out |
| All six operational sections | 1366×480 and 1366×600 | Both themes and both modes; 48 cases | Footer remains visible; navigation scrolls independently; no rail/page overflow |
| Mobile navigation drawer | 390×600 | Dark, expanded sheet | Footer separation, Escape dismissal and focus return |
| Short sign-in | 390×480 | Dark | Invalid-field state, scrolling to usable controls, visibility toggle and keyboard focus |
| Populated inventory | 1440, 390, 768 | Representative light/dark | Real table, internal scrolling, actual stock, bounded search/no-match state |

Mobile themes are switched using the visible sheet control and explicitly asserted on the document. Some supplemental captures retain keyboard focus outlines. Filters/actions wrap in the inspected mobile/tablet screenshots. Login and long content scroll vertically on short screens. Chromium headless checks and visual inspection were performed; physical-device, Safari and Firefox verification was not performed.

## 11–13. Security, automated and browser regressions

| Check | Actual result |
|---|---|
| `npm run lint` | PASS; 19 inherited warnings, zero new warnings |
| `npm run test:run` | PASS; 14 files, 211 tests (208 existing + 3 new persistent-layout/access-boundary tests) |
| `npm run build` | PASS; TypeScript and Vite production build |
| `go -C backend fmt ./...` | PASS; no backend files changed |
| `go -C backend vet ./...` | PASS |
| `go -C backend test ./...` | PASS; separately gated real-database suites remain opt-in in this command |
| Real Chromium / HTTP / PostgreSQL | PASS; 13 assertion groups in [after acceptance evidence](../ux/verification/pre-phase7-ui/after/acceptance.json), zero runtime exceptions |
| `git diff --check` | PASS |
| Preservation hash comparison | PASS; all tracked backend files, all 62 retained UI primitive files and all 47 approved-package files unchanged |
| Private secret scan | PASS; no private fixture passwords or generated environment secrets in changed text/artifacts |

New unit tests verify stable chrome during pending/denied/recovered authority, withheld cached route content with changed authority, and transient access failure/retry. Existing authentication tests still cover invalid credentials, rate/network/server failures, duplicate submissions, visibility, safe return destinations, inactive accounts, role boundaries and immediate private-data clearing on logout.

Browser verification repeats Dashboard, Requests & Borrowings, Inventory, Borrowers, Reports and Administration three times through UI links. Sidebar, header, main container and document identities remain identical; theme and collapsed state remain stable; content changes; no chrome-removal batches, document reload requests or refresh requests occur during these 18 navigations. Existing account checks still run. A real `/auth/me` request is deliberately paused and then continued, proving that protected page content stays withheld while chrome remains visible. An intentional full reload restores the session with one existing refresh request.

Real Admin, Staff and Borrower logins/logouts, anonymous redirects, restricted Staff navigation/Admin-route denial and Borrower operational denial pass. The Borrower terms flow uses only the clearly synthetic version in the disposable test harness; no official terms were published in the normal database. The shared account menu and mobile sheet pass keyboard/Escape focus checks. Populated inventory/search verification uses the actual isolated equipment API.

The test environment used the existing guarded Batch 1 harness: database `elabtrack_v2_batch1_test`, PostgreSQL port 54832, API port 18085 and Vite port 15175. Credentials were random, private and never rendered in captures. Mail uses the harness adapter; no Brevo delivery occurred. Only owned disposable resources/private fixture files were cleaned up afterward. Normal PostgreSQL, the owner's Admin account/data and the existing API/Vite services were preserved. The screenshot browser ran the development build; production compilation passed separately.

## 14. Screenshot evidence

[Review table](../ux/verification/pre-phase7-ui/README.md), [standalone side-by-side comparison](../ux/verification/pre-phase7-ui/index.html) and [image hashes/dimensions](../ux/verification/pre-phase7-ui/MANIFEST.json) provide the exact artifacts. There are **15 before PNGs and 43 after PNGs**. The final browser script records 48 after captures, including five repeated filenames; this is not 48 distinct images. All ten requested desktop comparisons exist at 1440×900: Login light/dark; Dashboard expanded/collapsed light/dark; Inventory expanded/collapsed light/dark. Additional evidence covers mobile/tablet, short heights, validation, navigation and populated/no-match inventory.

The stored before set predates source edits. The core after comparison retains the refined empty-inventory captures; later fixture-backed responsive/table captures represent their actual datasets. `frontend/scripts/pre-phase7-ui-qa.mjs` documents the executable assertions. Reproduction/environment notes are in the evidence README. These screenshots are engineering evidence and do not amend approved mockup PNGs.

## 15. Remaining deviations and dependencies

Owner visual acceptance is pending. Dashboard metrics, Requests & Borrowings, Reports and functional Notifications remain future work under the existing roadmap. Official FSMO terms approval and verified live Brevo configuration/delivery remain the existing independent institutional/deployment dependencies. No borrowing features or Phase 7 work were performed. The approved package remains unchanged; the implementation-map overlay identifies this owner-authorized refinement without representing the resulting design as accepted.

## 16. Exact files changed

- `docs/ux/APPROVED_VISUAL_IMPLEMENTATION_MAP.md`
- `frontend/src/app/router.tsx`
- `frontend/src/components/application/operational-shell.tsx`
- `frontend/src/components/application/shells.tsx`
- `frontend/src/components/application/visual.tsx`
- `frontend/src/features/accounts/pages/account-create.tsx`
- `frontend/src/features/accounts/pages/account-detail.tsx`
- `frontend/src/features/accounts/pages/account-directory.tsx`
- `frontend/src/features/accounts/pages/administration-page.tsx`
- `frontend/src/features/accounts/pages/student-bulk.tsx`
- `frontend/src/features/auth/components/access-boundary.tsx`
- `frontend/src/features/auth/pages/login-page.tsx`
- `frontend/src/features/auth/product-auth.test.tsx`
- `frontend/src/features/inventory/pages/category-management.tsx`
- `frontend/src/features/inventory/pages/equipment-detail.tsx`
- `frontend/src/features/inventory/pages/equipment-form.tsx`
- `frontend/src/features/inventory/pages/equipment-list.tsx`
- `frontend/src/features/inventory/pages/stock-adjustment.tsx`
- `frontend/src/features/workspace/workspace-page.tsx`
- `frontend/src/styles/index.css`
- `frontend/src/styles/refinement.css`
- `frontend/scripts/pre-phase7-ui-qa.mjs`
- `docs/project/PRE_PHASE7_UI_REFINEMENT_REPORT.md`
- `docs/ux/verification/pre-phase7-ui/README.md`, `index.html`, `MANIFEST.json`, `before/acceptance.json`, `after/acceptance.json` and the exact 58 PNG paths listed in that manifest.

Account/inventory page edits remove per-page shell wrappers; the remaining feature logic is preserved. No backend, migration, API contract, schema, domain, package-manifest, normal environment or approved mockup changes were made.

## 17. Exact Git status

Branch `main`; starting and ending HEAD `ed89f9eb5a7d50c69a4a76deed4564d57c142e10`. The index remains empty. No staging, commit, push or deployment occurred.

```text
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
 M frontend/src/features/inventory/pages/category-management.tsx
 M frontend/src/features/inventory/pages/equipment-detail.tsx
 M frontend/src/features/inventory/pages/equipment-form.tsx
 M frontend/src/features/inventory/pages/equipment-list.tsx
 M frontend/src/features/inventory/pages/stock-adjustment.tsx
 M frontend/src/features/workspace/workspace-page.tsx
 M frontend/src/styles/index.css
?? docs/project/PRE_PHASE7_UI_REFINEMENT_REPORT.md
?? docs/ux/verification/pre-phase7-ui/
?? frontend/scripts/pre-phase7-ui-qa.mjs
?? frontend/src/styles/refinement.css
```

The untracked evidence directory contains only the documented review artifacts; the image manifest enumerates each PNG. Private test environment/fixture files are ignored or outside the repository.

## 18. Owner acceptance

**OWNER ACCEPTANCE: PENDING.** The owner must review the local interface and screenshot comparisons. Engineering checks do not grant visual approval or permission to begin Phase 7.
