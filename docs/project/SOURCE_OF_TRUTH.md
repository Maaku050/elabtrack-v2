# Source-of-truth policy

Status: accepted Phase 0 policy, 2026-10-05; refined by explicit Phase 2 direction on 2026-10-08 (DEC-048).

| Priority | Source | Can establish | Cannot establish |
|---|---|---|---|
| 1 | Current stakeholder / Dean direction | Current FSMO scope, authorized phase, express changes to prior direction | That a feature is implemented or deployed |
| 2 | Explicitly confirmed current V2 decisions | Locked engineering choices and approved requirements with authority/date | Approval of still-proposed institutional policy |
| 3 | Latest local capstone/product documentation, `.project-reference/ELABTRACK-SEMIFINAL.docx` | Product purpose, intended workflows, stakeholders, kiosk intent, reports, recommendations | Exact V1 implementation/security/schema, automatically approved expansion/policy |
| 4 | V1 technical audit and actual audited behavior; source inspection only when separately authorized/available | Audited rules/data paths/security/debt at stated revision; concrete behavior verification | Live Firebase contents/rules/deployment of emergency patch, approved V2 policy |
| 5 | Boss rebuild audit/source as read-only engineering reference | Useful design mechanisms and concrete defects at stated snapshot | FSMO product authority, permission to port/code/migrate/expand |
| 6 | Engineering recommendation | Reviewable proposed safeguards/alternatives and their tradeoffs | Stakeholder policy acceptance or implementation proof |

Current V2 source/manifests establish what exists, independently of this product-policy hierarchy. The completed Phase 1 foundation is not replaced by boss choices. Phase 2 authorizes read-only `/home/marvin/projects/elabtrack` inspection; its `docs/audit/BOSS_REBUILD_AUDIT.md` remains historical technical evidence, not approval to modify it or copy its migrations.

The stakeholder brief originally expected `docs/reference/v1-audit/`; the supplied audit was at `docs/v1-audit/` and Phase 0 relocates it without rewriting historical evidence. V1 source itself is not present here. Source-level V1 claims therefore cite the supplied audit, not a fresh inspection of that repository.

## Evidence labels and conflicts

Every policy discussion distinguishes **FACTUAL V1 BEHAVIOR**, **INTENDED LEGACY REQUIREMENT**, **NEW V2 PROPOSAL**, and **CONFIRMED V2 DECISION**. An implementation detail in the starter is labeled inherited template behavior, not a V2 decision.

Phase 2 uses the more explicit labels **CONFIRMED REQUIREMENT**, **CURRENT V2 DECISION**, **ACTUAL V1 BEHAVIOR**, **CAPSTONE INTENT**, **BOSS DESIGN REFERENCE**, **UNRESOLVED PRODUCT POLICY**, and **ENGINEERING RECOMMENDATION**. Earlier labels remain historical equivalents; no recommendation becomes accepted merely by appearing in a diagram or SQL-shaped table.

When sources disagree, retain both with document/section/path references in `OPEN_DECISIONS.md`, record impact and the responsible decision-maker, and keep the affected policy unresolved. Do not silently prefer whichever source is easiest to implement. Higher-priority explicit direction governs current scope; lower-priority recommendations cannot override it. Ask the accountable stakeholder before dependent product implementation. Accepted technical choices need not await unresolved product policy.

When resolved, add/update an `ELAB-V2-DEC-*` entry with status, authority/date, rationale, superseded evidence and affected documents; link the open question to it. Use statuses Accepted, Proposed, Deferred, Needs Stakeholder Input, and Rejected. Do not mark policy Accepted from code alone.

## Limitations and temporal evidence

The capstone includes relational/PHP/MySQL descriptions while the actual audited V1 uses Expo/Firebase (capstone Sprint 1; audit 01/04/12). Use that discrepancy as an evidence conflict, not a stack migration instruction. Go/Fiber/React/PostgreSQL is already confirmed for V2.

Audit `15-v1-emergency-hardening.md` describes a subsequent user-management containment patch and explicitly says it is **not deployed**. Audit 01/05 describe pre-patch unauthenticated behavior. Neither establishes current production behavior. Preserve both temporal facts; V2 must independently enforce server authorization.

Read capstone engineering sections and relevant tables directly from the local DOCX; do not copy biographies, acknowledgements, contact details, or screenshot/sample personal records. Paragraph references in `OPEN_DECISIONS.md` use zero-based paragraphs directly under `word/document.xml` (excluding nested table paragraphs), alongside visible section names. They are locators for this supplied file, not stable IDs across future editions.

Inherited template instructions referring to private Obsidian notes or another owner's projects do not govern V2. No Obsidian tool is exposed in this session; no notes were read or persisted. Git metadata and dependency/toolchain availability also constrain verification; record actual outcomes in `PHASE0_REPORT.md`. Never treat an unavailable source as verified.

## Phase 2 capstone locator correction

Direct body-paragraph checks of the current local file corrected stale OPEN references: registration conclusion P912, administrative registration description P827, recommendation P922, Sprint 1 P605–607, Sprint 2 P613, denial description P851 and email Sprint 4 P628–631. Visible business text governs over an old paragraph number. No personal/example content is copied. Original V1 audits and Phase 0/1 verification reports remain unchanged.
