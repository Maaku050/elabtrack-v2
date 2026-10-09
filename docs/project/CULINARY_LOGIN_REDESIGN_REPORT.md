# Culinary login redesign and visual verification

2026-10-09. **OWNER ACCEPTANCE: PENDING.** The owner accepted the preceding Dashboard and Staff/Admin navigation as the current visual direction and requested this subsequent login refinement. Their acceptance does not authorize Phase 7 or approve the new login. This report supersedes the earlier login screenshots only; the preceding [UI refinement report](PRE_PHASE7_UI_REFINEMENT_REPORT.md) and its captures remain historical evidence.

## Selected visual and rights

The login now uses [culinary-flatlay.webp](../../frontend/src/assets/login/culinary-flatlay.webp), an original AI-generated photo-style asset created with the built-in image generation tool. A light blue overhead work surface, edge-arranged whisk, wooden spoon, flour, eggs and bowls convey cooking/equipment while preserving central space for branding. No stock preview, Google Images image or owner screenshot was copied. The existing FSMO seal is rendered separately and unchanged.

The optimized, locally bundled WebP is **1122×1402, 266,698 bytes**, versus a 2,813,128-byte source PNG, preserving source dimensions at native Chromium encoder quality 0.88. No image-service hotlink or new package was added. [Asset provenance, applicable output-rights references, SHA-256 and exact final prompt](../../frontend/src/assets/login/README.md) document its source and use. It has no third-party stock license/attribution requirement; this is not a claim of exclusive copyright, public-domain status or institutional approval.

## Layout and themes

The desktop two-column arrangement remains balanced: a bright photographic branding panel and the existing authentication form. A small pale backdrop makes the original seal, eLabTrack wordmark and full Food Service Management Organization identity readable. Navy typography and the concise “Equipment, managed with care.” tagline sit over the blue surface. The photograph has no dark overlay or filter, including in dark mode.

The form uses existing shadcn/ui primitives, navy/violet colors, semantic labels, 50px inputs, full-width primary action and the retained assistance copy. Light mode has a white form surface; dark mode uses the existing navy surface, contrasting fields and violet action. Form spacing is adjusted locally without editing shared application styling.

At 768px the two columns use reduced padding/type proportions. Below 768px the photograph becomes a 144–152px compact banner, prioritizing the form. The cropped edges still show utensils/ingredients. All requested widths (320, 390, 768, 1024, 1366, 1440) were captured in **both** actual themes. No horizontal overflow or cropped form controls appeared. At ordinary 900px-high desktop sizes (1024px and above), the page fits vertically; short-height/error/mobile pages scroll normally.

## Session-notice correction and security boundary

The broad existing `sessionEnded` state also covers deliberate sign-out and account changes, so it did not reliably describe an expired session. The login now uses separate, in-memory presentation metadata in the existing session-action store. Only terminal denial/invalidation of a live authenticated session causes the involuntary-end notice. Intentional logout, fresh anonymous bootstrap with no refresh cookie, and network/server failures stay neutral. Successful refresh/login clears the notice. A genuine notice retains `role="status"`.

The API client records this metadata beside its existing guarded denial/refresh/logout paths. Existing session clearing, token validation, revocation, generation checks, retry limits, cross-tab coordination and authority remain intact. No backend, endpoint, cookie, token, authentication store, coordinator, account policy or terms logic changed. Invalid credentials still clear/focus the password; normal loading, disabled controls, password visibility, username/current-password autocomplete and safe return navigation are retained. No registration or recovery link was introduced.

## Verification results

| Check | Actual result |
|---|---|
| Frontend tests | **PASS: 14 files, 223 tests**, including 12 new notice-classification/regression cases; previous baseline 211 |
| Frontend lint | **PASS: 19 inherited warnings**, no new warnings/errors |
| Production build | **PASS**, TypeScript and Vite; local WebP bundled/fingerprinted |
| Real Chromium + isolated HTTP/PostgreSQL | **PASS: 10 assertion groups**, 12 viewport/theme cases, 16 captures, zero runtime exceptions |
| Keyboard/focus | PASS: Tab order, Enter visibility activation, submit access at 1440×600, first invalid field at 390×480, password retry focus after actual rejection |
| Authentication regression | PASS: real Admin login/logout, later SPA login after deliberate logout, server-revoked refresh→401 notice, fresh anonymous entry, subsequent login/logout |
| Backend formatting/vet/tests | PASS: `go fmt ./...`, `go vet ./...`, `go test ./...`; backend files unchanged |
| Whitespace and preservation | `git diff --check` PASS; starting file hash comparison confirms the restricted source changes listed below |

The first backend command encountered a read-only default Go cache in the sandbox. The complete checks passed with a writable task-specific `GOCACHE` under `/tmp`; no toolchain change was needed. Versions: Go 1.26.5 and Node 24.19.0. Standard Go test invocation leaves separately gated database suites opt-in; this task separately executed real login checks through the existing Batch 1 browser server, not a full database-suite rerun.

The 12 new frontend cases distinguish anonymous/bootstrap failure, live terminal denial, network/server errors, peer sign-out/invalidation/account-change events, intentional logout failure and actual login notice rendering. Existing authentication/session tests continue to pass. The browser harness sends complete native Enter key event data and waits for React updates before inspecting password visibility.

## Screenshot evidence and inspection

[Review links](../ux/verification/culinary-login/README.md), [assertion evidence](../ux/verification/culinary-login/acceptance.json) and [image hashes/dimensions](../ux/verification/culinary-login/MANIFEST.json) identify all **16 distinct PNGs**. The width/theme matrix supplies the requested desktop light/dark, tablet, 390px and 320px evidence. Supplemental captures show a keyboard focus outline, invalid fields, actual invalid credentials and actual involuntary invalidation.

Visual inspection confirmed bright recognizable cooking objects, central readable navy copy, cohesive light/dark forms, the compact mobile banner and no horizontal cropping. No password is rendered in captures. On short mobile screens the lower form/help can require vertical scrolling, especially when a notice/error is present. Physical-device, Safari, Firefox and screen-reader testing were not performed. Browser rendering used Vite development mode; production compilation passed separately. Owner aesthetic acceptance is still pending.

## Local environment and preservation

Real browser verification used the existing guarded disposable Batch 1 environment: `elabtrack_v2_batch1_test`, PostgreSQL port 54832, API 18085 and Vite 15175. All credentials were generated private synthetic fixtures; no owner password was used, exposed or reset. Synthetic harness terms and fake mail remained confined to the test environment. No live Brevo delivery or official terms publication occurred. After verification, the owned test API/Vite stopped, and the disposable Compose resources and private fixture/environment files were removed after project-label/target checks. Normal API `/api/v1/health` and login returned HTTP 200; the normal PostgreSQL container remained healthy. A scan of changed source/artifacts found none of the seven generated private fixture/environment secret values.

The starting workspace already contained the earlier UI stabilization edits. Hash comparison against that full 968-file starting snapshot confirms all other retained files are unchanged, including Dashboard, navigation/sidebar, inventory, borrower pages, backend source, paired migrations, 62 UI primitive files, the approved visual package and the preceding screenshot/report artifacts. The implementation map records the owner's accepted Dashboard/navigation direction and keeps this new login pending.

## Exact changes in this task

Existing files changed relative to the starting workspace:

- `docs/ux/APPROVED_VISUAL_IMPLEMENTATION_MAP.md`
- `frontend/src/features/auth/pages/login-page.tsx`
- `frontend/src/features/auth/product-auth.test.tsx`
- `frontend/src/lib/api-client.ts`
- `frontend/src/lib/api-client.test.ts`
- `frontend/src/stores/session-action-store.ts`

New source/documentation files:

- `frontend/src/assets/login/culinary-flatlay.webp`
- `frontend/src/assets/login/README.md`
- `frontend/src/styles/culinary-login.css`
- `frontend/scripts/culinary-login-qa.mjs`
- `docs/project/CULINARY_LOGIN_REDESIGN_REPORT.md`
- `docs/ux/verification/culinary-login/README.md`
- `docs/ux/verification/culinary-login/acceptance.json`
- `docs/ux/verification/culinary-login/MANIFEST.json`

The 16 new PNG paths are enumerated exactly by `MANIFEST.json`; no previous capture was overwritten. No dependency manifest, shared stylesheet, backend, migration, product-policy or approved mockup change was made.

## Review status

**OWNER ACCEPTANCE: PENDING.** Review the final login at the local application and the linked screenshot matrix. There are no known failing login checks; remaining verification limits are identified above. Official FSMO terms approval and verified live Brevo configuration/delivery remain the existing external dependencies, independent of this presentation change. No Phase 7 work, staging, commit, push or deployment occurred. Branch `main` and HEAD `ed89f9eb5a7d50c69a4a76deed4564d57c142e10` remain unchanged.
