# Approved high-fidelity visual package

**Phase 3A.2 references approved; Phase 3B COMPLETE AND OWNER VISUALLY APPROVED, 2026-10-08.** Owner reviewed real implemented screenshots, including light/dark and responsive adaptations (DEC-064). The [approved handoff](approved/README.md) and [manifest](approved/MANIFEST.json) are the visual target. Their navy/indigo/violet direction supersedes the original evergreen SVG proposal. Old `mockups/` and [Phase 3A.2 report](../project/PHASE3A2_REPORT.md) remain historical structure/coverage evidence only; their original awaiting-review wording describes that earlier handoff. They are not current palette, geometry or implementation authority.

## Current package and implementation scope

36 manifest files: three core PNGs, one separately replaceable brand PNG and32 individual remaining screen PNGs. The combined S01 core contains two screens, so35 screen image files represent **36 screen targets:4 baselines +32 future targets**. Overview boards supplement the individual images and are not additional targets. All manifest hashes/dimensions/files verified and all individual PNGs visually inspected before coding.

| Baseline | Approved reference | Current evidence |
|---|---|---|
| Borrower Home | [B01](approved/core/B01-borrower-home.png) | Real component development preview; [390 capture](verification/phase3b/B01-390-light-full.png) |
| Equipment Catalog | [B02](approved/core/B02-borrower-equipment-catalog.png) | Real component development preview; [390 capture](verification/phase3b/B02-390-light-full.png) |
| Staff Dashboard | [S01 upper](approved/core/S01-staff-dashboard-and-pending-requests.png) | Real component development preview; [1440 capture](verification/phase3b/S01-dashboard-1440-light.png) |
| Pending Requests | [S01 lower](approved/core/S01-staff-dashboard-and-pending-requests.png) | Real component development preview; [1440 capture](verification/phase3b/S01-pending-1440-light.png) |

The remaining32 individual references are binding future targets, mapped individually with structure, shared/unique composition, role, responsive expectation, feature phase and review/implementation status in [APPROVED_VISUAL_IMPLEMENTATION_MAP](APPROVED_VISUAL_IMPLEMENTATION_MAP.md). None is implemented merely because a shell or shared status component exists. This36-target counting differs from the historical59 UX surfaces /34 evergreen compositions; it does not delete those behavior/variant requirements.

## Review method

Open [comparison.html](verification/phase3b/comparison.html) beside [visual system](VISUAL_SYSTEM.md), [fidelity contract and four matrices](VISUAL_FIDELITY_CONTRACT.md), and [browser results](verification/phase3b/RESULTS.json). Borrower941px images normalize to390 CSS px; S01 normalizes1536→1440, with upper frame0–526 and lower538–1024 inspected separately. Compare landmarks, hierarchy, colors, density and real controls, rather than trusting a unit-test result or comparing fictional names byte-for-byte.

Actual Chromium renders cover Borrower320/390/768/1280 and Staff1024/1440 in both themes. Mobile content scrolls as needed; navigation and cart clear safe areas. Dark artwork was not supplied in this approved set, so dark renders are a documented semantic translation now accepted by the owner through review of real implementation screenshots. Exact generated font identity and original vector seal are unavailable. Accessibility sizing, domain-copy corrections, preview markings and responsive interpretations are disclosed; no pixel-perfect claim.

## Domain copy controls

Approval is handover in one CHECKED_OUT transition; Pending reserves and expires24h after submission. The preview shows canonical lifecycle filter labels rather than Approved/Released, PHP money, borrower fine context that does not automatically block requests, and A/R/C/damaged-held physical counts instead of invented maintenance/retired buckets. Chart “Issued” represents the same approval-with-handover event. The fixture snapshot is08 Oct2026,10:00 AM Asia/Manila; examples are not live metrics.

Borrowers cannot mark returned or submit photographic evidence. All return photos stay entirely outside the software. Lost/damaged dispositions require replacement quantities, not price charges. Admin alone fully clears fines; completion requires all physical/replacement quantities zero and freezes fine calculation. Correct inaccurate image words while retaining composition; dependency questions are recorded in the fidelity contract, not silently turned into policy.

Each later authorized feature must inspect its matching PNG and run fresh screenshots with domain/integrity/authorization acceptance. Phase3B success does not authorize Phase4 or claim functional borrowing/inventory/accounts/reporting.
