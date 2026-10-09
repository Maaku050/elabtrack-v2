# Pre-Phase 7 UI review evidence

**OWNER ACCEPTANCE: PENDING.** Actual application screenshots; synthetic isolated identities only.

[Open the side-by-side comparison](index.html) (open this standalone HTML file in a browser). All images below are also directly accessible.

| State | Before | After |
|---|---|---|
| Login · light | [PNG](before/login-1440-light.png) | [PNG](after/login-1440-light.png) |
| Login · dark | [PNG](before/login-1440-dark.png) | [PNG](after/login-1440-dark.png) |
| Dashboard · expanded · light | [PNG](before/dashboard-expanded-1440-light.png) | [PNG](after/dashboard-expanded-1440-light.png) |
| Dashboard · expanded · dark | [PNG](before/dashboard-expanded-1440-dark.png) | [PNG](after/dashboard-expanded-1440-dark.png) |
| Dashboard · collapsed · light | [PNG](before/dashboard-collapsed-1440-light.png) | [PNG](after/dashboard-collapsed-1440-light.png) |
| Dashboard · collapsed · dark | [PNG](before/dashboard-collapsed-1440-dark.png) | [PNG](after/dashboard-collapsed-1440-dark.png) |
| Inventory · expanded · light | [PNG](before/inventory-expanded-1440-light.png) | [PNG](after/inventory-expanded-1440-light.png) |
| Inventory · expanded · dark | [PNG](before/inventory-expanded-1440-dark.png) | [PNG](after/inventory-expanded-1440-dark.png) |
| Inventory · collapsed · light | [PNG](before/inventory-collapsed-1440-light.png) | [PNG](after/inventory-collapsed-1440-light.png) |
| Inventory · collapsed · dark | [PNG](before/inventory-collapsed-1440-dark.png) | [PNG](after/inventory-collapsed-1440-dark.png) |

## Additional verification

The before set contains 15 unique PNGs; the after set contains 43 unique PNGs (48 captures, including five repeated filenames). The core desktop comparison retains the first refined empty-inventory captures. Subsequent responsive/table captures use real synthetic equipment stored through the isolated API; the normal development database remains unchanged. Some mobile images deliberately retain keyboard focus outlines.

[Before navigation observations](before/acceptance.json) · [After assertions and viewport matrix](after/acceptance.json) · [Closeout report](../../../project/PRE_PHASE7_UI_REFINEMENT_REPORT.md)

Widths: 320, 390, 768, 1024, 1366, 1440. Both actual themes. Desktop rail modes at 1024+. All six sections at 1366×480 and 1366×600 in both themes/modes. Mobile drawer at 390×600; login validation at 390×480.

## Reproduction

Use the existing guarded disposable environment in [integration/BATCH1.md](../../../../integration/BATCH1.md), start `TestBatch1BrowserServer` and Vite at localhost:15175 with the test API at localhost:18085. The fixture file must remain mode 0600. From `frontend/`, run `node scripts/pre-phase7-ui-qa.mjs` with `CHROMIUM_EXECUTABLE` set to your Chromium executable. The WSL library workaround is environment-specific; it is not an application dependency. The harness signs in with private random fixture credentials, delays a real `/auth/me` request, exercises the UI and creates synthetic equipment only in that disposable database. It never invokes Brevo.

`--before` captures the checked-out implementation; the stored before set was captured before source edits. Do not overwrite that baseline using the refined source. The default run preserves existing core after comparisons, so remove/regenerate those after images when intentionally evaluating a subsequent visual revision. Owner approval is separate from this engineering verification.
