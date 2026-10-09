# Sarcita-inspired immersive login redesign

2026-10-09. **OWNER ACCEPTANCE: PENDING.** The owner's attached Sarcita Sign In screenshot is the primary structural reference. It supersedes the two previous centered-card login attempts; the earlier [culinary report](CULINARY_LOGIN_REDESIGN_REPORT.md) and captures remain historical. Dashboard and Staff/Admin navigation remain the already accepted direction. This request authorizes only sign-in presentation, not Phase 7 or any product-policy change.

## Full-viewport composition and typography

The centered rounded authentication card is removed. The login now spans the viewport edge to edge with **59% photographic hero / 41% integrated form panel** at desktop sizes, closely following the supplied reference's 806/560 split at 1366×768. The hero is full height and the form has no outer card, border, shadow, radius or enclosing gutter. At 768px the ratio adapts to 52/48; below 768px a compact 200–204px photographic branding header precedes the form.

The unchanged FSMO seal and eLabTrack lockup sit at the upper left, without the previous pale rectangular logo container. The full Food Service Management Organization name remains readable in HTML. The large two-line headline is “Every tool in its place. / Every session on track.” Supporting copy describes organizing equipment and borrowing for FSMO practical sessions. Desktop headline scale, left inset, vertical placement and form hierarchy follow the primary reference while preserving the application's navy/violet identity. Mobile retains a compact two-line headline and omits the extra supporting paragraph/eyebrow for space.

## Culinary asset and photographic treatment

The new [culinary-laboratory.webp](../../frontend/src/assets/login/culinary-laboratory.webp) is an original **AI-generated** photo-style visual created with the built-in image tool: a spacious sunlit culinary laboratory with stainless stations, cookware, whisk, cutting board, vegetables and practical utensils. No identifiable people, borrowed stock preview or owner reference photo was supplied to generation. This is illustrative imagery, not a claim to photograph the actual FSMO facility.

The locally imported/fingerprinted WebP is **1122×1402, 204,424 bytes**, reduced from the 2,171,662-byte PNG with the native Chromium WebP encoder at quality 0.88; original dimensions preserved, no new dependency or external hotlink. [Origin, usage-rights references, SHA-256 and exact final generation prompt](../../frontend/src/assets/login/CULINARY_LABORATORY.md) are documented. No third-party stock was reused; no stock-site license/attribution requirement applies. Output rights remain governed by the applicable OpenAI agreement, without claims of exclusive copyright, guaranteed uniqueness, CC or public-domain status.

A controlled navy gradient supplies text contrast while preserving workstation detail and natural lighting. The desktop image/overlay stay the same across both themes. Tablet/mobile crops have stronger navy treatment because bright steel/windows occupy more of the narrow text region. This adjustment follows actual sampled hero-text contrast checks. There is no flat opaque replacement, teal palette or artificial purple wash. The previous blue flat-lay and provenance remain intact but are no longer imported by this page.

## Integrated authentication panel

The right panel uses the existing off-white background in light mode and deep navy in dark mode. A contextual “FSMO · Account access” header occupies the reference's top navigation spacing, with the existing theme switch aligned right. The welcome eyebrow, “Sign in to your account” heading, concise login instructions, 50px email input, 58px password input and full-width 50px violet action follow the reference's hierarchy and proportions. Fields have subtle surfaces, clear borders and visible keyboard focus; the existing password button has a 48px-wide target.

Email/password labels, autocomplete, password visibility, React Hook Form/Zod validation, disabled/loading states, safe server errors, cleared-password retry focus, assistance copy and intentional logout feedback are preserved. The existing session-notice classification is used unchanged: ordinary anonymous entry and deliberate logout are neutral; actual involuntary invalidation warrants the notice.

No login API, JWT/refresh protocol, token storage, browser coordinator, authentication guard, role route, account-status enforcement, backend, migration, terms infrastructure or policy was altered. The only test edit updates the expected semantic heading to the new heading; the anonymous-access/theme/protected-fetch assertions remain intact.

## Actual verification

| Gate | Result |
|---|---|
| Frontend tests | PASS: **14 files, 223 tests**, unchanged count; existing authentication, errors, loading, role, session and safe-destination regressions retained |
| Lint | PASS: **19 inherited warnings**, zero new warnings/errors |
| Production build | PASS: TypeScript and Vite; optimized local hero bundled |
| Backend formatting / vet / tests | PASS: `go fmt ./...`, `go vet ./...`, `go test ./...`, using a writable `/tmp` Go cache; backend files unchanged |
| Chromium responsive/theme checks | PASS: all six requested widths in both actual themes; full viewport geometry, desktop/tablet ratios, compact mobile banner, no horizontal overflow, minimum 44px controls and autocomplete |
| Real authentication / session checks | PASS: invalid credentials, Admin login/logout, later SPA login after deliberate logout, actual server-revoked refresh→401 notice, fresh anonymous entry and subsequent login/logout |
| Keyboard / focus | PASS: Tab to password/visibility/submit, native Enter visibility activation, first invalid-field focus at short mobile height, retry focus after real rejection |
| Hero contrast sampling | PASS: at least 4.5:1 for normal text and 3:1 for text at least 24px over the sampled rendered photographic background, using actual text-node bounds; minimum observed 4.64:1 normal and 4.53:1 large |
| Browser evidence | **10 assertion groups, 16 distinct full-page PNGs, zero runtime exceptions** |
| Whitespace / preservation | `git diff --check` PASS; starting-workspace hashes confirm only the three existing files below changed |

The required module-selected backend toolchain is **Go 1.27.1**; the global launcher outside the module remains Go 1.26.5. Node **24.19.0** satisfies the frontend engine. The initial foundation test expected the previous “Welcome back” heading; its single copy assertion was updated to “Sign in to your account,” then the entire 223-test suite passed. No assertions were removed or type checking weakened. Production compilation was rerun after the final CSS adjustment.

The browser assertions use actual Chromium, HTTP and PostgreSQL through the existing guarded Batch 1 test harness, not mocked login. The contrast probe temporarily hides text, samples a capture of the rendered image/gradient under actual text rectangles, then restores the page before review captures. This verifies sampled hero contrast, not whole-application WCAG certification or screen-reader behavior.

## Screenshots and remaining visual deviations

[Full-page screenshot index](../ux/verification/immersive-login/README.md), [assertion evidence](../ux/verification/immersive-login/acceptance.json) and [exact PNG dimensions/hashes](../ux/verification/immersive-login/MANIFEST.json) provide desktop light/dark, tablet, 390px and 320px captures in both themes. Additional captures show keyboard focus at 1440×600, invalid fields at 390×480, real authentication failure and actual session invalidation.

At ordinary desktop heights (1366×768 and 1024/1440×900), there is no document vertical overflow. Short-height and narrow mobile/error states can scroll vertically; full-page PNGs include the reachable fields, action, help and footer. No fixed viewport height clips the form. The actual screenshots were inspected for full-screen geometry, readable text, logo, image detail, borders, CTA and light/dark cohesion.

Intentional deviations from Sarcita are culinary photography, FSMO branding and navy/violet accents; contextual account-access text instead of a public-home link; Admin-provisioned assistance instead of registration, guest, social or recovery actions; and omission of sports/location badges or unapproved feature claims. The compact tablet/mobile adaptation preserves usable forms rather than squeezing desktop typography. Owner aesthetic acceptance remains pending. Physical devices, Safari, Firefox and screen readers were not tested. Browser captures use the development build; production compilation is verified separately. No deployment verification was performed.

## Local environment and scope preservation

Only the guarded disposable database `elabtrack_v2_batch1_test` on port 54832, test API 18085 and Vite 15175 were used for actual login checks. Random private fixture accounts and clearly synthetic test terms are confined to that environment, with the existing fake email adapter. No owner credentials, normal database writes, password resets, official terms publication or live Brevo delivery were used. A private-secret scan found none of the seven generated fixture/environment secret values in changed source/documents/captures. Existing migration pairs were applied only to the empty disposable database; migration files were not changed. After verification, the test API/Vite were stopped and the label-verified disposable Compose resources/private fixtures were removed. Normal API health and login returned HTTP 200; the normal PostgreSQL container remained healthy.

The starting workspace already contained prior UI refinements. Hash comparison against **992 starting files** confirms **989 unchanged files**, including **163 backend files, 62 UI primitives, 47 approved-package files**, Dashboard/navigation, borrower/inventory pages, auth transport/store/coordinator/guards, manifests, policies, previous assets and all previous reports/screenshots. Only the following existing files changed in this task:

- `frontend/src/features/auth/pages/login-page.tsx`
- `frontend/src/features/foundation/foundation.test.tsx`
- `docs/ux/APPROVED_VISUAL_IMPLEMENTATION_MAP.md`

New files:

- `frontend/src/styles/immersive-login.css`
- `frontend/src/assets/login/culinary-laboratory.webp`
- `frontend/src/assets/login/CULINARY_LABORATORY.md`
- `frontend/scripts/immersive-login-qa.mjs`
- `docs/project/IMMERSIVE_LOGIN_REDESIGN_REPORT.md`
- `docs/ux/verification/immersive-login/README.md`, `acceptance.json`, `MANIFEST.json`
- The 16 exact PNG paths enumerated in `MANIFEST.json`.

## Exact Git status

Branch `main`, HEAD `ed89f9eb5a7d50c69a4a76deed4564d57c142e10`; empty index. This status includes inherited unstaged/untracked edits from the preceding tasks; the hash comparison above distinguishes this task's changes. No staging, commit, push or deployment occurred.

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
 M frontend/src/features/foundation/foundation.test.tsx
 M frontend/src/features/inventory/pages/category-management.tsx
 M frontend/src/features/inventory/pages/equipment-detail.tsx
 M frontend/src/features/inventory/pages/equipment-form.tsx
 M frontend/src/features/inventory/pages/equipment-list.tsx
 M frontend/src/features/inventory/pages/stock-adjustment.tsx
 M frontend/src/features/workspace/workspace-page.tsx
 M frontend/src/lib/api-client.test.ts
 M frontend/src/lib/api-client.ts
 M frontend/src/stores/session-action-store.ts
 M frontend/src/styles/index.css
?? docs/project/CULINARY_LOGIN_REDESIGN_REPORT.md
?? docs/project/IMMERSIVE_LOGIN_REDESIGN_REPORT.md
?? docs/project/PRE_PHASE7_UI_REFINEMENT_REPORT.md
?? docs/ux/verification/culinary-login/
?? docs/ux/verification/immersive-login/
?? docs/ux/verification/pre-phase7-ui/
?? frontend/scripts/culinary-login-qa.mjs
?? frontend/scripts/immersive-login-qa.mjs
?? frontend/scripts/pre-phase7-ui-qa.mjs
?? frontend/src/assets/login/
?? frontend/src/styles/culinary-login.css
?? frontend/src/styles/immersive-login.css
?? frontend/src/styles/refinement.css
```

## Owner review

**OWNER ACCEPTANCE: PENDING.** Engineering checks do not grant visual approval. Review the current local login and final screenshot matrix. No Phase 7 work or subsequent phase is authorized by this report. Existing official FSMO terms approval and configured/tested live Brevo dependencies remain independent and unchanged.
