# Phase 4A rendered authentication evidence

Verified 2026-10-08 with Chrome 140.0.7339.16 against the real Go API and isolated PostgreSQL 18.6 database. [Reproduction and cleanup](../../../../integration/PHASE4A.md), [phase report](../../../project/PHASE4A_REPORT.md).

[RESULTS.json](RESULTS.json) records **24 rendered screenshots and 29 passing integrated checks**, zero application console errors/warnings, no horizontal overflow, and no interactive target below 44px in the captured layouts. Authentication replies were never fabricated. The harness holds actual refresh requests, toggles Chromium offline or blocks the API URL for adverse cases. Four random-password, synthetic accounts were created only in the isolated database and removed afterward; no password, access token or cookie value appears in these artifacts.

| Surface | Widths | Rendered evidence |
|---|---|---|
| Login | 320, 390, 1280 | [320 light](login-320-light.png), [320 dark](login-320-dark.png), [390 light](login-390-light.png), [390 dark](login-390-dark.png), [1280 light](login-1280-light.png), [1280 dark](login-1280-dark.png) |
| Borrower protected catalog placeholder | 320, 390, 1280 | [320 light](borrower-320-light.png), [320 dark](borrower-320-dark.png), [390 light](borrower-390-light.png), [390 dark](borrower-390-dark.png), [1280 light](borrower-1280-light.png), [1280 dark](borrower-1280-dark.png) |
| Own account | 390 | [Read-only safe account metadata](account-390-light.png) |
| Borrower forbidden | 390 | [Light](forbidden-390-light.png), [dark](forbidden-390-dark.png) |
| Staff workspace | 1024, 1440 | [1024 light](staff-1024-light.png), [1024 dark](staff-1024-dark.png), [1440 light](staff-1440-light.png), [1440 dark](staff-1440-dark.png) |
| Staff forbidden from Admin | 1024 | [Dark](staff-forbidden-1024-dark.png) |
| Admin workspace | 1440 | [Light](admin-1440-light.png), [dark](admin-1440-dark.png) |
| Logout network failure | 390 | [Dark, persistent honest failure and retry](logout-network-error-390-dark.png) |
| Session restoration network failure | 390 | [Light, bounded manual retry](session-network-error-390-light.png) |

The integrated run covers real form login, generic invalid/inactive feedback, safe deep-link restoration, role/navigation denial, current-account demotion and disable after issuance, stale Admin JWT denial, private-cache removal, refresh/reload without privileged flash, revoked-session restoration denial, native cross-tab logout, non-secret lifecycle messages, memory-only access state and HttpOnly refresh-cookie invisibility. Keyboard checks exercise password focus after failure, show/hide naming/state and named focus targets. This is targeted accessibility evidence, not a full WCAG or cross-browser audit.

Representative final renders were opened and inspected, including narrow/light Borrower, mobile/dark Borrower, mobile login in both themes, wide login, Staff light, Admin dark, account, forbidden and session/logout failure. Comparisons opened the established B01/S01 core references plus B04-02 and B06-02. No pixel-perfect or new owner visual approval is claimed. The login keeps the compact FSMO brand, heading/form hierarchy, violet action and semantic theme tokens; anonymous login has no protected bottom navigation. Larger accessible labels, responsive wide/dark translations and conditional real notices differ from compressed generated references. Authenticated workspaces intentionally have honest placeholders rather than preview metrics and fabricated borrowing records. Staff sign-out contrast was corrected after screenshot inspection; collapsed navigation retains its accessible icon button.

[PREVIEW_REGRESSION.json](PREVIEW_REGRESSION.json) separately records **24 original development-preview viewport/theme cases and eight interaction checks**, with zero console errors/warnings. Its anonymous refresh interception is labeled preview-only and is not live authentication evidence. [PRODUCTION_REGRESSION.json](PRODUCTION_REGRESSION.json) records four quick preview cases and ten checks, including production login at root, preview-route 404 and absence of synthetic routes/fixtures in production JavaScript. Historical Phase 3B and local-environment evidence is unchanged.

The local HTTP refresh cookie is HttpOnly, host-only, SameSite=Lax and scoped to `/api/v1/auth`; Secure=false is the existing explicit development/test policy. Production Secure=true remains covered by existing backend tests; production HTTPS deployment was not performed. [VALIDATION.json](VALIDATION.json) contains the safe quality/database/cleanup summary.
