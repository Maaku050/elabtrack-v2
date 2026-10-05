# Source-of-truth policy

Status: accepted Phase 0 policy, 2026-10-05.

| Priority | Source | Can establish | Cannot establish |
|---|---|---|---|
| 1 | Current stakeholder direction and explicitly confirmed V2 decisions | Current FSMO scope, locked technical choices, authorized phase, approved changes to legacy policy | That a proposed or documented feature is already implemented |
| 2 | Latest local capstone/product documentation, `.project-reference/ELABTRACK-SEMIFINAL.docx` | Product purpose, intended workflows, stakeholders, kiosk intent, reports, recommendations | Exact V1 implementation/security/data schema, approval of future expansion, production delivery configuration |
| 3 | V1 technical audit, `docs/reference/v1-audit/` | Audited architecture, actual rules/data paths, security evidence, technical debt and scalability observations at the stated revision | Live deployed rules, current Firebase contents, deployment of the emergency patch, automatically approved V2 policy |
| 4 | V1 source behavior, when its revision is available for authorized inspection | Reproduction/verification of concrete code paths and clarification of audit evidence | Stakeholder approval, unstated deployed configuration, approved V2 requirements |
| 5 | Current template architecture/source | Bootstrap code, manifests, routes, infrastructure and defects that actually exist in this repository | eLabTrack product requirements, borrower eligibility, role authority, fine policy, deployment readiness |

The stakeholder brief originally expected `docs/reference/v1-audit/`; the supplied audit was at `docs/v1-audit/` and Phase 0 relocates it without rewriting historical evidence. V1 source itself is not present here. Source-level V1 claims therefore cite the supplied audit, not a fresh inspection of that repository.

## Evidence labels and conflicts

Every policy discussion distinguishes **FACTUAL V1 BEHAVIOR**, **INTENDED LEGACY REQUIREMENT**, **NEW V2 PROPOSAL**, and **CONFIRMED V2 DECISION**. An implementation detail in the starter is labeled inherited template behavior, not a V2 decision.

When sources disagree, retain both with document/section/path references in `OPEN_DECISIONS.md`, record impact and the responsible decision-maker, and keep the affected policy unresolved. Do not silently prefer whichever source is easiest to implement. Higher-priority explicit direction governs current scope; lower-priority recommendations cannot override it. Ask the accountable stakeholder before dependent product implementation. Accepted technical choices need not await unresolved product policy.

When resolved, add/update an `ELAB-V2-DEC-*` entry with status, authority/date, rationale, superseded evidence and affected documents; link the open question to it. Use statuses Accepted, Proposed, Deferred, Needs Stakeholder Input, and Rejected. Do not mark policy Accepted from code alone.

## Limitations and temporal evidence

The capstone includes relational/PHP/MySQL descriptions while the actual audited V1 uses Expo/Firebase (capstone Sprint 1; audit 01/04/12). Use that discrepancy as an evidence conflict, not a stack migration instruction. Go/Fiber/React/PostgreSQL is already confirmed for V2.

Audit `15-v1-emergency-hardening.md` describes a subsequent user-management containment patch and explicitly says it is **not deployed**. Audit 01/05 describe pre-patch unauthenticated behavior. Neither establishes current production behavior. Preserve both temporal facts; V2 must independently enforce server authorization.

Read capstone engineering sections and relevant tables directly from the local DOCX; do not copy biographies, acknowledgements, contact details, or screenshot/sample personal records. Paragraph references in `OPEN_DECISIONS.md` use zero-based paragraphs directly under `word/document.xml` (excluding nested table paragraphs), alongside visible section names. They are locators for this supplied file, not stable IDs across future editions.

Inherited template instructions referring to private Obsidian notes or another owner's projects do not govern V2. No Obsidian tool is exposed in this session; no notes were read or persisted. Git metadata and dependency/toolchain availability also constrain verification; record actual outcomes in `PHASE0_REPORT.md`. Never treat an unavailable source as verified.
