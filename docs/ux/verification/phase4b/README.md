# Phase 4B real browser evidence

2026-10-08, Chromium140.0.7339.16, real isolated PostgreSQL/API and actual form interactions. [RESULTS.json](RESULTS.json) records21 passing checks and26 screenshots with zero application console errors/warnings. No official FSMO terms document is approved. Every visible document is explicitly SYNTHETIC TEST / NOT OFFICIAL and existed only in the disposable `elabtrack_v2_phase4b_test` database. No fake terms or account is added to normal application storage. Credentials/tokens are excluded from artifacts.

| State | Captured sizes/themes |
|---|---|
| Unpublished | 320/390, light/dark |
| Required, complete long document | 320/390/768/1280, light/dark |
| Reading end, unchecked checkbox/button | 390, light/dark |
| Accepted / updated / stale conflict | 390, light/dark each |
| Acceptance transport failure | 390, dark |
| Terms-read transport failure | 390, light |
| Staff / Admin unaffected | 1280, light/dark each |

All automated viewport captures report no page horizontal overflow and effective interactive targets at least44px. The shadcn checkbox's16px visual indicator has a full-width/minimum44px label target. Keyboard Space selects explicit consent. Long content is plain text with page scrolling; no hidden or tiny-modal document. Browser cases verify first-use routing, missing content, read-only access, pending single submission, original repeat receipt, three tabs/reload, changed mandatory version/stale review reset, network failure, cross-tab logout/private-cache clearing, disabled/generic inactive login, current-role separation, memory-only token/HttpOnly cookie and login after prior acceptance.

Compared with approved [B04-03](../../approved/remaining/B04-03.png): existing FSMO brand/header, navy/violet tokens, card border/spacing, version/status, checkbox/button hierarchy and bottom navigation are retained. Required runtime adaptations are disclosed: official policy copy is absent; synthetic long text verifies reading; unchecked consent replaces the static checked mockup; existing accepted shell brand is retained; theme/logout/read-only actions are present; the acceptance button follows the complete document rather than floating over long text. Dark/wide translations use the already approved Phase3B tokens/shells. These screenshots are implementation evidence, not a new owner visual approval or institutional wording approval.

[Preview regressions](preview-regression/RESULTS.json) independently record24 preserved approved-preview layouts,10 checks and production fixture exclusion, zero application errors/warnings. That inherited preview harness uses only an anonymous refresh fixture; authenticated terms checks use real API/PG without fake responses. Reproduce using [integration guide](../../../../integration/PHASE4B.md). See [phase report](../../../project/PHASE4B_REPORT.md) for remaining content/activation gates and explicit Phase7 integration requirements. No production deployment claim.
