# eLabTrack V2 — APPROVED VISUAL HANDOFF

## Status
The project owner APPROVED these UI references in ChatGPT on 2026-10-08. This pack is the approved visual reference set for implementation. Approval does NOT change the locked Phase 2/2.5 domain rules or authorize implementation beyond Phase 3B.

## Mandatory precedence
1. Current approved FSMO blue/violet **visuals here** for palette, typography, brand/header, layout, component shape, spacing, hierarchy, visual density, and role-specific navigation.
2. Current domain and Phase 2.5 decisions for meanings, actual state names, actions, permissions, quantity/fine formulas, and data consistency. If a generated image contains a conflicting word, status, or sample, DOMAIN RULES OVERRIDE the image's erroneous copy; preserve visual treatment.
3. Phase 3A.1 UX architecture for flows and responsive constraints.
4. Phase 3A.2 legacy evergreen docs/SVGs as supplementary structural inventory ONLY. Their evergreen tokens and screenshots are SUPERSEDED and MUST NOT set colors or visual direction.

## Files
- `approved/brand/FSMO-seal-reference.png` — high-resolution recreated seal used as visual reference. It is NOT independently verified official vector art.
- `approved/core/B01-borrower-home.png` — final approved borrower mobile Home (compact header).
- `approved/core/B02-borrower-equipment-catalog.png` — final approved borrower mobile Equipment Catalog (compact header).
- `approved/core/S01-staff-dashboard-and-pending-requests.png` — approved staff desktop Dashboard (upper) and Pending Requests (lower).
- `approved/remaining/B*.png` — individual borrower mockups from batches 2–6.
- `approved/remaining/S*.png` — individual staff/admin mockups from batches 2–6.
- `approved/overview/batch-*.png` — overview boards for convenience; individual PNGs are primary for details.
- `approved/MANIFEST.json` — all individual screen IDs, titles and SHA256 fingerprints.

## Safety/accuracy cautions
These are static illustrative designs, NOT screenshots of working software. Generated sample labels/data sometimes contradict confirmed requirements. Specifically NEVER implement Approved and Released as separate borrowing lifecycle states; approval and physical handover are the SAME action. Pending means reserved, not borrowed. Borrowers cannot record authoritative returns. Damage/loss creates replacement obligations, not automatic money charges. Admin-only full fine clearance preserves assessments/history. No return-photo upload/attachment feature. Do not expose fake functionality or fake network states in production.

## Fidelity and scoping rules
- Phase 3B is shared shadcn design system, browser shell, component catalog and representative UI previews; it does NOT authorize full backend business features or all production functional screens.
- Do not use the artwork as one full-screen flattened `<img>` or bitmap UI. Rebuild with accessible React, CSS/Tailwind, existing shadcn/Base UI and Lucide primitives.
- Preserve composition, type scale, hierarchy, spacing, border radii, icon family, table density, states, light/dark semantics, nav placement, and responsive behavior as closely as practicable.
- No speculative added interactions. If a visual is cropped, ambiguous or contains an AI artifact, document it in a fidelity/deviation ledger; do not invent a rule to match it.
- Use the same component system across later feature phases. Each actual screen must get viewport screenshot comparisons against its corresponding reference, with mismatches reviewed.
- Match **layout and design intent**, not hardcoded example names, values, greeting, request counts, or fixed screenshots. Data comes later from authoritative API/domain.

## Existing project UX docs
Read `docs/ux/UX_ARCHITECTURE.md`, `BORROWER_FLOWS.md`, `STAFF_ADMIN_FLOWS.md`, `RESPONSIVE_RULES.md`, `WIREFRAMES.md`, `STATUS_SYSTEM.md`, and current Phase 2.5 domain documents before building. The earlier Phase 3A.2 `VISUAL_SYSTEM.md` evergreen values need explicit supersession, not silent reuse.
