# eLabTrack V2 — final human-review UI polish

**Published evidence:** The release retains 42 representative screenshots across current light/dark desktop/mobile, before/after, real transaction, PDF and offline checks. All 600 excluded screenshots are preserved locally. Original full-run capture totals below remain historical verification facts, not the number of files shipped. Galleries and capture lists are curated; [the evidence manifest](../ux/verification/RELEASE_EVIDENCE_MANIFEST.json) records every included/excluded screenshot. Raw owner reference crops are excluded.

**Current owner acceptance and publication authority, 2026-10-11 — DEC-085:** Final Staff/Admin desktop and Student/Faculty mobile visual review passed; the owner reported no remaining UI issues and authorized auditing, committing and pushing the completed Phases 8–14 application and shared UI to `release/elabtrack-v2-presentation-ready`. Main and production deployment remain untouched. Normal schema 000008 and existing records are preserved. Official FSMO terms/current acceptance, production Student domains, verified activation/Brevo and operational deployment readiness remain separate gates. Earlier pending-review/no-commit statements below are historical task snapshots. [Release audit and publication handoff](RELEASE_PUBLICATION_REPORT.md).

Owner review/task start: 2026-10-10. Final verification: 2026-10-11 (Asia/Shanghai). Authority: explicit final human-review master request and DEC-084. This is a focused shared-source refinement of the previously engineering-verified reconstruction. **ELABTRACK V2 — HUMAN-REVIEW UI POLISH: ENGINEERING VERIFIED; OWNER VISUAL ACCEPTANCE PENDING.** Owner visual acceptance of this final pass remains pending; no business-policy changes, normal database migration/reset, commit, push or deployment.

## Staff/Admin completion matrix

| Item | Implemented result | Verification |
|---|---|---|
| HR-01 | Thin native themed scrollbars on page, activity, table, issuance/cart and category scrollers; hover tokens, Firefox progressive enhancement, forced-color fallback. Accepted Dashboard composition retained. | COMPLETE — presentation/operational Chromium PASS |
| HR-02 | Compact operational toolbar and real supported-state count strip; separate reference/borrower identity/equipment/units/time/status/action columns; avatar or initials; strong Pending Review and existing detail Open; compact 52px rows, sticky header, independent table scrolling, visible server pager. Original shell retained. | COMPLETE — unit/presentation/operational PASS |
| HR-03 | Original three-stage wizard and frozen review retained. Stage 2 contains search, content-sized catalog rows, an independently scrollable side cart and accessible Back/Next footer. Adjacent compact controls; card quantities reflect the actual selected cart quantity. | COMPLETE — wizard, presentation/operational and real issuance/return PASS |

## Student/Faculty completion matrix

| Item | Implemented result | Verification |
|---|---|---|
| MR-01 | Accepted Home source/composition and existing navigation retained. | ACCEPTED/PRESERVED — unchanged source; both browser environments PASS |
| MR-02A | Existing Base UI Popover provides an anchored, portalled, responsive panel above chips/cards; collision handling, autofocus, Escape/outside close and returned trigger focus. Existing server search/category/sort/availability preserved. | COMPLETE — unit/presentation/operational PASS |
| MR-02B | `[−][quantity][+][Add]` is a contiguous visual cluster with transparent minimal quantity controls and emphasized Add; mobile targets remain at least 44×44px. Stock bounds, cart state/edit/remove and validation retained; explicit visible keyboard focus and invalid-input border. | COMPLETE — unit/presentation/operational PASS |
| MR-03A | Exactly one red Sign Out in the Account heading, matching the Terms heading position; bottom action removed; existing logout hook unchanged. | COMPLETE — unit and real Student/Faculty logout in both environments PASS |
| MR-03B | 24px spacing between personal information and terms; Account-specific bottom safe-area allowance preserves access above fixed navigation. Profile/avatar/readonly details retained. | COMPLETE — presentation/operational PASS |

## Direct comparison with references

Both documents were read completely and all **18** embedded screenshots inspected individually. [The reference manifest](../ux/verification/human-review/reference-manifest.json) records document/image hashes, figure IDs, dimensions and inspection status. Original DOCX/extracted media remain private and ignored. Three separately privacy-reviewed crops D-03/D-04/A-03 contain no real personal names/emails/credentials and accompany the safe fictional screenshots.

- **D-01–04:** Accepted four-KPI/three-panel/Stock & Accountability composition remains. Native scrollbars use transparent tracks, rounded subtle neutral/violet thumbs and brighter hover; no JavaScript fake-scroll library or hidden production scrollbars.
- **R-01/R-02:** The generic directory becomes an operational queue with a compact search/status toolbar, integrated count pills, dense aligned identity/equipment/unit/time rows and visible Review/Open actions and pagination. Navy/violet tokens and the accepted sidebar/header remain. Only PENDING, CHECKED_OUT (Active), COMPLETED, DENIED, CANCELLED, EXPIRED and All are used. Mockup Approved/Released/For Return states, unsupported bulk checkboxes, Purpose/Course and accountability fields are omitted rather than invented. Overdue remains the actual row-level flag; no unsupported overdue status-filter endpoint is introduced. Full references remain accessible via text/title and the detail link. Times use one Asia/Manila context, without repeating the timezone per row.
- **I-01–04:** The V1 distribution of space is adopted: bounded image-led list on the left, selected-equipment cart on the right, reachable Back/Next below. Stage progress and selected borrower remain visible. Minimal nearby minus/quantity/plus/Add replace separated boxes. Prices, currency totals and payment controls are omitted. At 600px height the catalog is deliberately smaller and scrolls internally; the footer stays available.
- **H-01:** Home is preserved. Actual full names/data can differ from reference values; no fictional values are hardcoded into application components.
- **E-01–03:** The complete filter panel now layers above catalog content intentionally rather than being painted behind the first card. V1 supplies adjacency/minimal controls only; V2 retains its B02 cards, actual protected imagery and 44px touch targets.
- **A-01–04:** Centered profile/avatar/readonly metadata remain; one red Sign Out sits alongside Account, as the Terms positioning reference requests. Personal Information and Borrowing Terms now have 24px separation, and the latter is reachable above bottom navigation.

[Before/after gallery](../ux/verification/human-review/index.html). The 24 captured BEFORE views are genuine pre-change fictional runtime screenshots. Their inherited harness hid native scrollbars; D-03/D-04 are the owner-provided native BEFORE evidence. Final AFTER automation explicitly enables native scrollbar rendering. No synthetic “before” was created by changing styles after implementation.

## Functional safeguards and query behavior

Frontend-only application changes; backend source, API contracts and all migrations remain byte-identical to the task-start hashes. No stock, reservation, checkout, return, replacement, fine, terms, activation, authorization, ledger, notification or report business logic changed.

Status counts are returned by **seven bounded existing list requests** with `per_page=1`, one per supported state plus All, intersecting the existing search and borrower-ID constraints. These are TanStack Query caches scoped to the signed-in actor, enabled for Staff/Admin only, fresh for 30 seconds and invalidated by existing borrowing commands. They are not page-slice counts or authorization data. Errors display unknown counts rather than zero. Each API request enforces current authorization. Separate count reads are not one atomic database snapshot: concurrent changes can briefly make counters differ until refetch/invalidation. This is a display limitation of the retained API, not relaxed stock or eligibility authority. No backend endpoint or contract was added.

Server pagination, URL restoration, search/status intersections and detail navigation remain. Local quantity/cart controls do not reserve/issue stock. Existing mandatory future due date, current borrower eligibility/terms, physical handover, fresh review, frozen payload/key, duplicate prevention and authoritative final mutation remain.

## Verification

Final engineering gates passed on 2026-10-11. [Machine evidence](../../integration/evidence/2026-10-10-human-review-ui-polish.json) and browser JSON provide check names/screens. Earlier reconstruction’s 121 database checks/83 browser groups remain historical, not newly rerun claims.

| Gate | Executed result |
|---|---|
| Frontend full suite | **334/334 tests, 32 files, PASS**, two workers; final run 79.23s. Includes auth/role/session/terms, cart/wizard/stale/duplicate and existing feature regressions. |
| Focused human-review/wizard | **10/10 tests, 2 files, PASS**. Additional existing wizard assertion verifies selected card/cart quantities agree. |
| TypeScript / production build | `npm run build` (`tsc -b && vite build`), PASS. No dependencies added. |
| Lint | `npm run lint`, PASS; **19 inherited Oxlint warnings**. Explicit ESLint on changed TypeScript files, PASS/zero diagnostics. |
| Go | `go fmt ./...`, `go vet ./...`, `go test ./...`, PASS; **28 public packages** report tests passing, with existing cache reuse. No backend/migration source changed. |
| Presentation Chromium | **8 check groups / 146 screenshots**, PASS; native scrollbars visible, zero JavaScript errors. |
| Normal production-bundle Chromium | **8 check groups / 146 screenshots**, PASS against isolated compatible API; no presentation decoration. Automatically restored presentation afterward. |
| Real PostgreSQL-backed workflow browser | **6 check groups / 13 screenshots**, PASS; actual Staff checkout/future due/current terms/frozen review/one POST, inspected good return/exact full stock restoration/immutable history, Faculty read-only access, menu logout. |
| Responsive acceptance | 1366/1440 × 600/768, expanded/collapsed, both themes, independent keyboard/native scrolling and visible actions/pagers; Student and Faculty 320/375/390/440 × 844, both themes, popover hit-testing/Escape/outside/focus, 44px adjacent controls, real quantity/cart/remove and genuine logout. PASS. |
| Presentation integrity | **11 read-only categories PASS**: counts, stock, ledger, movement sums, chronology, ledger chronology, synthetic data, replacement overflow, notification uniqueness, distinct fictional avatars, asset provenance hashes. |
| Data preservation | Fresh baseline digests compare **all 21 normal development tables** and **all 22 accepted Phase 7 demo tables**, identical. No normal writes/migrations/reset. |
| Shared-source consistency | **265 frontend source/manifest files** match the presentation snapshot byte-for-byte. Home/Dashboard source hashes and all backend/migration hashes unchanged. |
| Secret/Git gates | `git diff --check` and final private-secret audit PASS; exact inventory linked below. Index remains empty. |

Generic Go database-dependent cases may skip without fixture configuration. Dedicated fresh PostgreSQL suites were **not rerun** for this frontend-only task. The actual isolated PostgreSQL-backed issuance/return test above is newly executed, and verifies unchanged stock boundaries/audit rather than substituting fixture-only simulation. Production deployment/normal migration remain unrun and unauthorized.

Earlier failures were resolved before declaring PASS: inherited footer/row spacing, grid track shrinking and Account bottom-nav clearance; a complete TermsStatus test fixture fixed a strict TypeScript error. Cart focus and server-page transitions required the browser harness to wait for the actual settled UI/data before clicking/asserting. No test assertions or type checking were weakened; exact quantity, pager, hit-test and auth outcomes pass.

## Files and environments

Application changes: `frontend/src/styles/{human-review,index}.css`; `features/borrowing/pages/{borrowing-directory,request-compose,catalog-cart}.tsx`; `features/borrowing/hooks/use-borrowing.ts`; `features/borrowing/components/equipment-selection-card.tsx`; `features/workspace/workspace-page.tsx`. Regression changes: `features/borrowing/{human-review,direct-wizard}.test.tsx`; `frontend/scripts/browser-cdp.mjs` adds an opt-in visible-scrollbar option while preserving previous harness defaults; new human-review BEFORE/browser/workflow scripts. Documentation: this report/progress, DECISIONS, SOURCE_OF_TRUTH, ROADMAP, exact Git inventory, public evidence/gallery/reference manifest and `.gitignore` exclusions for the two new private owner documents.

Normal production and presentation builds use the same source components/styles. The normal bundle is verified on the isolated compatible schema/API origin without presentation decoration; it is not tested by migrating normal schema8. The full presentation remains at `http://127.0.0.1:15177/login`, API18087/PostgreSQL54836/schema11. Normal development and accepted Phase 7 target remain untouched.

### Exact task document/source paths

- `.gitignore`
- `docs/project/COMPLETE_FRONTEND_UI_RECONSTRUCTION_REPORT.md`
- `docs/project/DECISIONS.md`
- `docs/project/HUMAN_REVIEW_UI_POLISH_GIT_STATUS.txt`
- `docs/project/HUMAN_REVIEW_UI_POLISH_PROGRESS.md`
- `docs/project/HUMAN_REVIEW_UI_POLISH_REPORT.md`
- `docs/project/ROADMAP.md`
- `docs/project/SOURCE_OF_TRUTH.md`
- `docs/ux/verification/human-review/index.html`
- `docs/ux/verification/human-review/reference-manifest.json`
- `frontend/scripts/browser-cdp.mjs`
- `frontend/scripts/human-review-before.mjs`
- `frontend/scripts/human-review-browser-qa.mjs`
- `frontend/scripts/human-review-workflow-qa.mjs`
- `frontend/src/features/borrowing/components/equipment-selection-card.tsx`
- `frontend/src/features/borrowing/direct-wizard.test.tsx`
- `frontend/src/features/borrowing/hooks/use-borrowing.ts`
- `frontend/src/features/borrowing/human-review.test.tsx`
- `frontend/src/features/borrowing/pages/borrowing-directory.tsx`
- `frontend/src/features/borrowing/pages/catalog-cart.tsx`
- `frontend/src/features/borrowing/pages/request-compose.tsx`
- `frontend/src/features/workspace/workspace-page.tsx`
- `frontend/src/styles/human-review.css`
- `frontend/src/styles/index.css`
- `integration/evidence/2026-10-10-human-review-ui-polish.json`

Evidence files under `docs/ux/verification/human-review/{before,after,operational,workflows}/` are additionally enumerated individually in the cumulative Git inventory and acceptance JSON. Private utilities/logs/snapshots/references remain excluded.

## Final repository and safety audit

Branch `main`; HEAD `c2741c54ba3aca95df705ce612ade1a277e0335d` unchanged. **78 modified tracked files, 954 untracked public paths, 1,032 cumulative changed/untracked paths; zero staged files.** These cumulative counts include preserved earlier Phases8–14/reconstruction work, not 1,032 new source changes from this polish. The exact task delta is the 25 paths above plus **332 PNG artifacts** (24 genuine runtime BEFORE + 3 privacy-reviewed owner crops, 146 presentation AFTER, 146 normal-bundle AFTER, 13 actual workflow captures) and acceptance JSON.

Known-private-secret audit: **2,917 public files**, **48 private configurations**, **378 distinct private values compared**, **zero matches**. No credentials printed. Raw owner DOCX documents/extracted media, private `.env`/configuration, generated credentials, sessions, database backups, logs/runtime/build outputs and incomplete evidence remain ignored and retained locally. Three public crops were separately reviewed for fictional/no identity content and PNG metadata; all other raw reference images remain excluded. Gallery links, all 332 PNGs and each acceptance JSON’s screenshots are verified.

`git diff --check` and `git diff --cached --check` pass. No business/API/auth/stock authority changed; no security regressions observed in the executed checks. No commit, push, deployment, normal migration or reset. Normal and accepted Phase 7 digests remain identical; the restored fictional presentation is available on15177.

## Remaining limitations and acceptance

Owner visual acceptance is pending. Native scrollbar appearance varies by browser/OS; Firefox and forced-color/reduced-motion branches are progressive enhancements, not physical assistive-device or Firefox runtime certification. Physical phone/non-Chromium checks and production deployment are unrun. Normal database schema8→11 requires separate migration authorization; this task does not supply it. Official FSMO terms/publication/current consent, approved Student domains and required activation ownership/Brevo live delivery readiness remain production gates. No terms were invented/published/accepted in normal data.

No commit, push or deployment. See [exact cumulative Git inventory](HUMAN_REVIEW_UI_POLISH_GIT_STATUS.txt) and [machine evidence](../../integration/evidence/2026-10-10-human-review-ui-polish.json). Existing unrelated uncommitted work is preserved. Stop for owner visual acceptance.
