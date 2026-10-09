# Stage A visual review

Current owner acceptance, 2026-10-09: **Stage A approved; Stage B visually approved; final Inventory filter manual verification PASSED**. [Checkpoint acceptance record](../../../project/verification/frontend-checkpoint/ACCEPTANCE.json) supersedes historical pending labels in the original capture manifests below; those manifests retain their original test-time context.


**OWNER VISUAL ACCEPTANCE: PENDING.** Real Chromium screenshots of the Stage A implementation; synthetic data from isolated PostgreSQL/API fixtures only. No owner database records or credentials are included. Counts include previously executed disposable integration fixtures, including extreme quantity and nonzero custody buckets. They are not production seed data.

17 browser checks passed. [Machine-readable acceptance evidence](acceptance.json) lists exact checks and request metadata without credentials. These tests do not constitute owner visual approval.

| Review | Light | Dark |
| --- | --- | --- |
| Expanded Borrowers | [View](borrowers-expanded-light.png) | [View](borrowers-expanded-dark.png) |
| Expanded Inventory | [View](inventory-expanded-light.png) | [View](inventory-expanded-dark.png) |
| Collapsed sidebar | [View](inventory-collapsed-light.png) | [View](inventory-collapsed-dark.png) |
| Mobile navigation | [View](mobile-navigation-light.png) | [View](mobile-navigation-dark.png) |
| Mobile Inventory | [View](inventory-mobile-light.png) | [View](inventory-mobile-dark.png) |

Pagination: [server page 2 and footer](borrowers-pagination-page2-light.png). Keyboard tooltip: [collapsed navigation](collapsed-tooltip-light.png). Footer: [actual account menu](account-footer-menu-light.png). Protected image: [catalog thumbnail](inventory-catalog-image-light.png).

## Responsive Borrowers

- 390px: [Light](borrowers-responsive-390-light.png), [Dark](borrowers-responsive-390-dark.png).
- 768px: [Light](borrowers-responsive-768-light.png), [Dark](borrowers-responsive-768-dark.png).
- 1024px: [Light](borrowers-responsive-1024-light.png), [Dark](borrowers-responsive-1024-dark.png).
- 1366px: [Light](borrowers-responsive-1366-light.png), [Dark](borrowers-responsive-1366-dark.png).

## Feedback states

[Empty results](borrowers-empty-light.png), [loading](borrowers-loading-light.png), [error and Retry](borrowers-error-light.png). Loading/error screenshots intentionally use controlled network interception; ordinary navigation/filter/count/image/permission checks use the real existing API and PostgreSQL.

## Preserved diagnostic evidence

`failure.png` is a retained screenshot from a failed harness attempt. It is excluded from the 25 accepted captures and must not be treated as approved output. Historical evidence outside this folder was preserved unchanged. The report describes the corrected harness assumptions and the failed populated-database preconditions.

## All accepted captures

- [borrowers-expanded-light.png](borrowers-expanded-light.png)
- [borrowers-expanded-dark.png](borrowers-expanded-dark.png)
- [borrowers-pagination-page2-light.png](borrowers-pagination-page2-light.png)
- [borrowers-empty-light.png](borrowers-empty-light.png)
- [inventory-expanded-light.png](inventory-expanded-light.png)
- [inventory-expanded-dark.png](inventory-expanded-dark.png)
- [inventory-catalog-image-light.png](inventory-catalog-image-light.png)
- [inventory-collapsed-light.png](inventory-collapsed-light.png)
- [inventory-collapsed-dark.png](inventory-collapsed-dark.png)
- [collapsed-tooltip-light.png](collapsed-tooltip-light.png)
- [account-footer-menu-light.png](account-footer-menu-light.png)
- [borrowers-responsive-390-light.png](borrowers-responsive-390-light.png)
- [borrowers-responsive-390-dark.png](borrowers-responsive-390-dark.png)
- [borrowers-responsive-768-light.png](borrowers-responsive-768-light.png)
- [borrowers-responsive-768-dark.png](borrowers-responsive-768-dark.png)
- [borrowers-responsive-1024-light.png](borrowers-responsive-1024-light.png)
- [borrowers-responsive-1024-dark.png](borrowers-responsive-1024-dark.png)
- [borrowers-responsive-1366-light.png](borrowers-responsive-1366-light.png)
- [borrowers-responsive-1366-dark.png](borrowers-responsive-1366-dark.png)
- [mobile-navigation-light.png](mobile-navigation-light.png)
- [mobile-navigation-dark.png](mobile-navigation-dark.png)
- [inventory-mobile-light.png](inventory-mobile-light.png)
- [inventory-mobile-dark.png](inventory-mobile-dark.png)
- [borrowers-loading-light.png](borrowers-loading-light.png)
- [borrowers-error-light.png](borrowers-error-light.png)
