# Culinary login visual verification

2026-10-09. **OWNER ACCEPTANCE: PENDING.** The accepted Dashboard and Staff/Admin navigation direction is preserved; these captures document only the subsequent login refinement.

[Redesign report](../../../project/CULINARY_LOGIN_REDESIGN_REPORT.md), [asset, rights and final prompt](../../../../frontend/src/assets/login/README.md), [original bundled visual](../../../../frontend/src/assets/login/culinary-flatlay.webp), [browser assertions](acceptance.json), [PNG dimensions and SHA-256](MANIFEST.json).

## Final responsive captures

| Width | Light | Dark |
|---|---|---|
| 320px | [Light](login-320-light.png) | [Dark](login-320-dark.png) |
| 390px | [Light](login-390-light.png) | [Dark](login-390-dark.png) |
| 768px | [Light](login-768-light.png) | [Dark](login-768-dark.png) |
| 1024px | [Light](login-1024-light.png) | [Dark](login-1024-dark.png) |
| 1366px | [Light](login-1366-light.png) | [Dark](login-1366-dark.png) |
| 1440px | [Light](login-1440-light.png) | [Dark](login-1440-dark.png) |

320/390 captures are 844px high; the other widths are 900px high. Tablet and desktop retain two columns; smaller widths use a compact photograph/banner. These 12 cases assert actual theme, loaded local image, no horizontal overflow, minimum 44px form controls, correct autocomplete and column count. Ordinary desktop cases at 1024px and above also assert no document vertical overflow.

## Interaction and short-height captures

- [1440×600 dark: keyboard focus](login-keyboard-1440x600-dark.png) — Tab reaches email, password, visibility and submit; Enter activates the native visibility button.
- [390×480 dark: invalid fields](login-validation-390x480-dark.png) — both validation messages appear and first invalid input receives focus; vertical scrolling is expected at this short height.
- [1440×900 light: actual authentication failure](login-authentication-error-1440-light.png) — server returns invalid credentials; password is cleared and focused for retry. The synthetic email was then cleared before capture, which also exposes its client validation message.
- [390×844 light: actual revoked session](login-session-invalidated-390-light.png) — a live test session's refresh cookie is revoked by the existing logout API; a subsequent real refresh receives 401 and the notice appears. A new anonymous page entry remains neutral. The extra notice adds vertical scrolling on mobile.

## Reproduction and limits

`frontend/scripts/culinary-login-qa.mjs` uses the existing dependency-free CDP helper and a real Chromium executable. It expects the existing guarded Batch 1 browser harness: Vite at localhost:15175, synthetic API at localhost:18085, PostgreSQL database `elabtrack_v2_batch1_test` at port 54832, and private generated fixtures at `/tmp/elabtrack-batch1-fixtures.json`. Follow the existing isolated-harness procedure; do not point the script at the normal database or supply owner credentials. Set `CHROMIUM_EXECUTABLE` to a local browser binary (and any local library path needed), then run `node scripts/culinary-login-qa.mjs` from `frontend`. No browser dependency was added to the manifest.

All 10 assertion groups pass; 16 distinct captures; zero runtime exceptions. Actual Admin login/logout and subsequent login after revocation are verified. The anonymous and login captures contain no password or access/refresh token. Only the disposable synthetic environment was used, with the existing fake mail adapter; no live email or normal database writes occurred. The browser rendered the development build; production compilation was checked separately. Physical devices, Safari, Firefox and a screen reader were not tested. Engineering verification does not constitute owner approval.
