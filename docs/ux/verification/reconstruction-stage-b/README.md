# Stage B Chromium evidence

Current owner acceptance, 2026-10-09: **Stage A approved; Stage B visually approved; final Inventory filter manual verification PASSED**. [Checkpoint acceptance record](../../../project/verification/frontend-checkpoint/ACCEPTANCE.json) supersedes historical pending labels in the original capture manifests below; those manifests retain their original test-time context.


**OWNER VISUAL ACCEPTANCE: PENDING.** All images are real application viewport captures using synthetic records on isolated PostgreSQL; no owner data, passwords, activation links or live mail.

[Complete accepted index](INDEX.json) records 31 browser checks and 447 distinct screenshots. Group manifests record actual HTTP request paths, checks and screenshots. The seven requested widths (320/390/440/768/1024/1366/1440) were checked in both themes; short 560px heights are included. Additional projection/uploader captures scroll to the relevant region. These are viewport screenshots, not synthetic mockups or full-page composites.

| Surface | Desktop / primary | Mobile / supporting |
| --- | --- | --- |
| Student creation | [Capture](borrower-create-student-1366-light.png) | [Capture](borrower-create-student-390-dark.png) |
| Faculty creation | [Capture](borrower-create-faculty-1366-light.png) | [Capture](borrower-create-faculty-390-dark.png) |
| Borrower details/profile | [Capture](borrower-details-1366-light.png) | [Capture](borrower-details-390-dark.png) |
| Student workbook upload | [Capture](student-bulk-1366-light.png) | [Capture](student-bulk-390-dark.png) |
| Complete validation | [Capture](bulk-validation-1366-light.png) | [Capture](bulk-validation-390-dark.png) |
| Creation confirmation | [Capture](bulk-creation-confirmation-1366-light.png) | [Capture](bulk-creation-confirmation-390-dark.png) |
| Deactivation confirmation | [Capture](bulk-deactivation-confirmation-1366-light.png) | [Capture](bulk-deactivation-confirmation-390-dark.png) |
| Categories / dialog | [Capture](equipment-categories-1366-light.png) | [Capture](category-create-dialog-390-dark.png) |
| Equipment image create/edit | [Capture](equipment-add-image-area-light.png) | [Capture](equipment-edit-image-area-light.png) |
| Equipment details | [Capture](equipment-details-final-1366-light.png) | [Capture](equipment-details-final-390-dark.png) |
| Stock adjustment / review | [Capture](stock-adjustment-1366-light.png) | [Capture](stock-review-projection-320-dark.png) |
| Inventory correction / review | [Capture](inventory-correction-1366-light.png) | [Capture](correction-review-projection-320-dark.png) |
| Administration | [Capture](administration-overview-1366-light.png) | [Capture](administration-overview-390-dark.png) |
| Staff/Admin directory | [Capture](staff-admin-directory-1366-light.png) | [Capture](staff-admin-directory-390-dark.png) |
| Staff creation | [Capture](staff-create-1366-light.png) | [Capture](staff-create-390-dark.png) |
| Staff details | [Capture](staff-details-1366-light.png) | [Capture](staff-details-390-dark.png) |
| Borrower catalog | [Capture](borrower-equipment-catalog-1366-light.png) | [Capture](borrower-equipment-catalog-390-dark.png) |
| Borrower equipment details | [Capture](borrower-equipment-details-1366-light.png) | [Capture](borrower-equipment-details-390-dark.png) |
| Maximum quantities / archived status | [Capture](inventory-maximum-stock-1366-light.png) | [Capture](inventory-maximum-stock-1366-dark.png) |

Category Dialog and final inventory captures were refreshed after the last touch-target/column fixes. A file ending in `-failure.png` is diagnostic evidence from an intermediate harness/assertion run and is excluded from the accepted index. Published TEST ONLY terms/acceptance existed only in the disposable database, which was removed after verification. Normal local API/frontend remained healthy.
