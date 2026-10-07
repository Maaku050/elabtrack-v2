# Source-of-truth policy

Status: accepted Phase 0 policy, 2026-10-05; refined by explicit Phase2 direction (DEC-048) and authoritative Phase2.5 owner working decisions (DEC-051–061), 2026-10-08.

| Priority | Source | Can establish | Cannot establish |
|---|---|---|---|
| 1 | Current explicit stakeholder / project owner / Dean direction, with stated authority and date | Current FSMO scope, authorized phase, express changes to prior direction | That a feature is implemented or deployed |
| 2 | Explicitly confirmed current V2 decisions | Locked engineering choices and approved requirements with authority/date | Approval of still-proposed institutional policy |
| 3 | Latest local capstone/product documentation, `.project-reference/ELABTRACK-SEMIFINAL.docx` | Product purpose, intended workflows, stakeholders, kiosk intent, reports, recommendations | Exact V1 implementation/security/schema, automatically approved expansion/policy |
| 4 | V1 technical audit and actual audited behavior; source inspection only when separately authorized/available | Audited rules/data paths/security/debt at stated revision; concrete behavior verification | Live Firebase contents/rules/deployment of emergency patch, approved V2 policy |
| 5 | Boss rebuild audit/source as read-only engineering reference | Useful design mechanisms and concrete defects at stated snapshot | FSMO product authority, permission to port/code/migrate/expand |
| 6 | Engineering recommendation | Reviewable proposed safeguards/alternatives and their tradeoffs | Stakeholder policy acceptance or implementation proof |

Current V2 source/manifests establish what exists, independently of this product-policy hierarchy. The completed Phase 1 foundation is not replaced by boss choices. The earlier Phase2 scope authorized read-only `/home/marvin/projects/elabtrack` inspection; its `docs/audit/BOSS_REBUILD_AUDIT.md` remains historical technical evidence, not approval to modify it or copy its migrations.

The stakeholder brief originally expected `docs/reference/v1-audit/`; the supplied audit was at `docs/v1-audit/` and Phase 0 relocates it without rewriting historical evidence. V1 source itself is not present here. Source-level V1 claims therefore cite the supplied audit, not a fresh inspection of that repository.

## Evidence labels and conflicts

Every policy discussion distinguishes **FACTUAL V1 BEHAVIOR**, **INTENDED LEGACY REQUIREMENT**, **NEW V2 PROPOSAL**, and **CONFIRMED V2 DECISION**. An implementation detail in the starter is labeled inherited template behavior, not a V2 decision.

Phase 2 uses the more explicit labels **CONFIRMED REQUIREMENT**, **CURRENT V2 DECISION**, **ACTUAL V1 BEHAVIOR**, **CAPSTONE INTENT**, **BOSS DESIGN REFERENCE**, **UNRESOLVED PRODUCT POLICY**, and **ENGINEERING RECOMMENDATION**. Earlier labels remain historical equivalents; no recommendation becomes accepted merely by appearing in a diagram or SQL-shaped table.

When sources disagree, retain both with document/section/path references in `OPEN_DECISIONS.md`, record impact and the responsible decision-maker, and keep the affected policy unresolved until explicit higher-priority direction resolves it. Do not silently prefer whichever source is easiest to implement. Higher-priority explicit direction governs current scope; lower-priority recommendations cannot override it. Ask the accountable stakeholder before dependent product implementation. Accepted technical choices need not await unresolved product policy.

When resolved, add/update an `ELAB-V2-DEC-*` entry with status, authority/date, rationale, superseded evidence and affected documents; link the open question to it. Use statuses Accepted, Proposed, Deferred, Needs Stakeholder Input, and Rejected. Do not mark policy Accepted from code alone.

## Limitations and temporal evidence

The capstone includes relational/PHP/MySQL descriptions while the actual audited V1 uses Expo/Firebase (capstone Sprint 1; audit 01/04/12). Use that discrepancy as an evidence conflict, not a stack migration instruction. Go/Fiber/React/PostgreSQL is already confirmed for V2.

Audit `15-v1-emergency-hardening.md` describes a subsequent user-management containment patch and explicitly says it is **not deployed**. Audit 01/05 describe pre-patch unauthenticated behavior. Neither establishes current production behavior. Preserve both temporal facts; V2 must independently enforce server authorization.

Read capstone engineering sections and relevant tables directly from the local DOCX; do not copy biographies, acknowledgements, contact details, or screenshot/sample personal records. Paragraph references in `OPEN_DECISIONS.md` use zero-based paragraphs directly under `word/document.xml` (excluding nested table paragraphs), alongside visible section names. They are locators for this supplied file, not stable IDs across future editions.

Inherited template instructions referring to private Obsidian notes or another owner's projects do not govern V2. No Obsidian tool is exposed in this session; no notes were read or persisted. Git metadata and dependency/toolchain availability also constrain verification; record actual outcomes in `PHASE0_REPORT.md`. Never treat an unavailable source as verified.

## Phase 2 capstone locator correction

Direct body-paragraph checks of the current local file corrected stale OPEN references: registration conclusion P912, administrative registration description P827, recommendation P922, Sprint 1 P605–607, Sprint 2 P613, denial description P851 and email Sprint 4 P628–631. Visible business text governs over an old paragraph number. No personal/example content is copied. Original V1 audits and Phase 0/1 verification reports remain unchanged.

## Phase 2.5 working policy authority and supersession

The project owner explicitly supplied29 authoritative working product decisions in **“PHASE 2.5 — PRODUCT DECISION INTEGRATION & DOMAIN REBASELINE”**, 2026-10-08. The owner implemented V1 and knows FSMO's earlier manual workflow. Those statements establish current working V2 policy at priority1; this is not a claim that every detail was independently signed off by the Dean or is already implemented. DEC-051–061 and [BUSINESS_RULES decision matrix](../domain/BUSINESS_RULES.md#decision-integration-matrix) trace each statement.

They resolve earlier population/eligibility, provisioning, role distinction, reservation/release, due/terms, replacement, fine-clearance and kiosk assumptions. In particular, “Interactive Kiosk” is the normal responsive catalog/cart/request experience, not dedicated hardware/shared-device authentication. Owner correction supersedes the kiosk-shell parts of DEC-038/039, the original Phase11 wording and any older project/charter interpretation. FSMO-only scope and borrower mobile-first remain.

Audited V1 price-based damage/loss fines and closure after quantities accounted remain **FACTUAL V1 BEHAVIOR**; current V2 instead requires replacements and completion after both physical and replacement quantities resolve. PHP10/day is confirmed; recommended ceiling elapsed24h is selected engineering detail. Hybrid stock (DEC-062), exact schema/locking and live full-clear checkpoints are explicitly engineering recommendations selected for this baseline, not newly invented institutional policy.

**D17 corrected by explicit owner clarification, 2026-10-08:** a borrower may physically show a photo stored on their own phone during a face-to-face return, entirely outside eLabTrack. Staff/Admin may inspect actual equipment and alone records authoritative return condition/quantities. Earlier “optional supporting evidence” language must not imply software media: no return-photo acceptance, upload, transmission, storage, retention, attachments, borrower evidence submission, file model or workflow step. Structured return records/audit and independent equipment catalog images retain their scope. Other Phase2.5 working decisions are unchanged.

The original [Phase2 report](PHASE2_REPORT.md) is a historical handoff snapshot; its BLOCKED readiness/counts/alternatives are superseded by [Phase2.5 report](PHASE2_5_REPORT.md) and current domain documents. Original V1 audits, charter's Phase0-era evidence, and Phase0/1 verification reports are preserved rather than rewritten as if these working answers existed earlier. OPEN_DECISIONS now distinguishes resolved IDs from specific NON-BLOCKING later details. The roadmap is sequencing, not authorization: Phase2.5 changes Markdown design only; Phase3A/business implementation remain unstarted.
