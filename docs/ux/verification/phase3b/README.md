# Phase3B browser evidence

2026-10-08. Actual Chromium140.0.7339.16, Vite8.3.2, Node24.19.0/npm11.17.0. This verifies presentation and the preserved foundation route behavior, not production inventory/auth business features. [Phase report](../../../project/PHASE3B_REPORT.md), [fidelity matrices](../../VISUAL_FIDELITY_CONTRACT.md), [comparison viewer](comparison.html).

## Recorded artifacts

- `RESULTS.json`:24 viewport/theme cases plus10 boundary/keyboard/production checks; exact page/control/table/landmark measurements, application errors, anonymous refresh diagnostics and API request inventory.
- 24 viewport PNGs: B01/B02 at 320, 390, 768, 1280; S01-dashboard/S01-pending at 1024, 1440; each light/dark, 900px tall.
- 2 full-scroll 390px light PNGs; selection Dialog and small-width navigation Sheet captures: 28 PNGs total. Chromium full-page captures leave fixed navigation at its viewport position; use viewport captures and recorded scroll-clearance checks to assess obstruction.
- `ASSETS.json`: all36 approved manifest files verified against exact dimensions/SHA-256. Reference package retained unchanged.
- `CONTRAST.json`:28 opaque semantic foreground/backing pairs calculated using sRGB luminance; all meet selected4.5:1 text /3:1 essential boundary/focus targets. This does not certify rasterized text, every composite overlay, or all possible primitive variants. Charts also have textual data; pale chart segments are not the sole information channel.
- `SCOPE.json`: changed tracked files and protected-file invariance evidence.

Every viewport capture and the full-scroll/overlay captures were opened and visually inspected. Core comparisons use equal normalized widths: Borrower941→390, Staff1536→1440 with two S01 frames separately. Vertical dimensions are not stretched to conceal accessible row/control differences. Dark screenshots are an engineering translation; no approved dark PNG exists in this handoff. See disclosed differences/owner-review status in the fidelity contract.

## Reproduce

Use an installed Chromium executable with its shared libraries available. Start `npm run dev -- --host 127.0.0.1` in `frontend`; build with `npm run build`, and serve `npm run preview -- --host 127.0.0.1 --port 4173`. Then, from `frontend`:

```sh
CHROMIUM_EXECUTABLE=/absolute/path/to/chromium \
  node scripts/phase3b-browser-qa.mjs --production-url http://127.0.0.1:4173
```

The runner uses Node standard library/CDP, no added npm dependencies. Default capture destination is this directory; `--output /tmp/your-evidence` preserves existing evidence. `--quick` limits captures to four light baseline cases; it is not the full closure matrix. `--url` selects a local development server. For this environment, missing libnspr4/libnss3/libasound2t64 were downloaded/extracted under `/tmp`, with library search path supplied to Chromium; no system package or repository dependency change. Loopback servers/Chromium needed sandbox execution escalation and automatic review permitted them.

The QA runner intercepts only `/api/v1/auth/refresh` and returns the normal **anonymous401** response with origin-specific CORS. Existing SessionBootstrap executes unchanged and remains unauthenticated; no access token, fake role session or business success response is introduced. All other business API calls fail the run. Expected anonymous401 browser network diagnostics are counted separately; application JS exceptions/console errors and browser warnings must be zero. Lazy preview routes provide a named loading fallback. A manual preview uses the existing foundation backend/session behavior and can show its real session-restoration error if unavailable.

Checks include title/render, page overflow, selected text clipping, visible control targets, actual theme-button switching, content and cart-dock clearance above navigation, named table scrolling, Filter Sheet Escape/focus restoration, Dialog Tab trapping, small-width Staff Sheet, Admin navigation without authenticating, keyboard table scroll, desktop sidebar collapse, reduced motion, production fixture/route exclusion, production404 and original foundation root.

Physical touch devices, virtual keyboards/notches, full browser zoom and all future production loading/error/mutation states remain later feature/pilot checks.320px provides reflow evidence but is not a claim that200% browser zoom was exercised. No backend/PG/Docker/deployment regression is claimed for this visual-only phase.
