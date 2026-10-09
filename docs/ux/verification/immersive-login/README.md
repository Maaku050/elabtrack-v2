# Immersive login visual verification

2026-10-09. **OWNER ACCEPTANCE: PENDING.** The owner's Sarcita Sign In screenshot is the primary composition reference for this full-viewport login. These are new captures; earlier culinary/card attempts and the accepted Dashboard/navigation direction remain preserved.

[Report](../../../project/IMMERSIVE_LOGIN_REDESIGN_REPORT.md), [asset, rights and exact prompt](../../../../frontend/src/assets/login/CULINARY_LABORATORY.md), [local hero WebP](../../../../frontend/src/assets/login/culinary-laboratory.webp), [browser evidence](acceptance.json), [PNG dimensions and hashes](MANIFEST.json).

## Full-page responsive screenshots

| Width | Light | Dark |
|---|---|---|
| 320px | [Light](login-320-light.png) | [Dark](login-320-dark.png) |
| 390px | [Light](login-390-light.png) | [Dark](login-390-dark.png) |
| 768px | [Light](login-768-light.png) | [Dark](login-768-dark.png) |
| 1024px | [Light](login-1024-light.png) | [Dark](login-1024-dark.png) |
| 1366px | [Light](login-1366-light.png) | [Dark](login-1366-dark.png) |
| 1440px | [Light](login-1440-light.png) | [Dark](login-1440-dark.png) |

320/390 viewports are 844px high, 1366 is 768px high (matching the supplied reference), and the others are 900px high. Every PNG is a **full-page** capture, so the entire form/help/footer remains reviewable on mobile and short screens. Actual PNG dimensions are recorded in the manifest; a full page can be taller than its viewport. Ordinary desktops fit their viewport without vertical scrolling. Desktop ratio is 59/41; tablet 52/48; mobile uses a compact 200–204px culinary hero above the form.

## Interaction and short-height evidence

- [1440×600 dark keyboard focus](login-keyboard-1440x600-dark.png): Tab order and native Enter activation of password visibility, followed by submit focus.
- [390×480 dark invalid fields](login-validation-390x480-dark.png): first-invalid-field focus and validation messages; full page shows reachable content below the initial viewport.
- [Actual invalid credentials, 1440 light](login-authentication-error-1440-light.png): real server rejection, password cleared and focused. The private synthetic email is then cleared before capture, which additionally produces client email validation.
- [Actual revoked session, 390 light](login-session-invalidated-390-light.png): real API logout revokes the refresh cookie while client access state remains live; the next refresh gets 401 and redirects to the genuine session notice. A subsequent fresh anonymous page entry is neutral.

The browser checks loaded local imagery, actual light/dark themes, full viewport bounds, zero card border/radius, desktop hero ratio/height, tablet ratio, compact mobile hero, no horizontal overflow, unclipped minimum-44px form controls, autocomplete, focus and runtime exceptions. Hero-text contrast sampling uses actual text-node rectangles over a temporary text-hidden capture of the rendered photo/gradient, with minimum 4.5:1 for normal text and 3:1 for text at least 24px. It does not certify the entire application's accessibility or replace assistive-technology testing.

## Reproduction and boundaries

The dedicated `frontend/scripts/immersive-login-qa.mjs` uses the retained dependency-free Chromium CDP helper and the existing [guarded Batch 1 harness](../../../../integration/BATCH1.md). It expects private synthetic fixtures at `/tmp/elabtrack-batch1-fixtures.json`, test API localhost:18085 and Vite localhost:15175, backed by disposable `elabtrack_v2_batch1_test` on PostgreSQL port 54832. Never use owner credentials or point this fixture script at the normal database. Set `CHROMIUM_EXECUTABLE` (and any local library path needed), then run `node scripts/immersive-login-qa.mjs` from `frontend` after starting those guarded test services. No browser dependency was installed.

The screenshots contain no passwords or session tokens. Synthetic terms and fake email are confined to the disposable harness; no institutional terms or real Brevo delivery are implied. Browser rendering used Vite development mode; production build passed separately. Physical devices, Safari, Firefox and screen readers remain untested. Owner visual review is pending.
