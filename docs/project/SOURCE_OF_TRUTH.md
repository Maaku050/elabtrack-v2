# Source-of-truth policy

**Batch 1 execution overlay, 2026-10-09:** Explicit owner authorization covers Phase5 and conditional Phase6 only. Current implementation is governed by [the scoped plan](BATCH1_IMPLEMENTATION_PLAN.md), [actual account API contracts](../API_CONTRACTS.md#phase-5-implemented-account-management-contracts) and phase verification reports; older unstarted checkpoints remain historical. Product policy DEC-070–073 is unchanged. Official terms, approved Student domains and verified real Brevo delivery are independent external gates.
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


## Phase 4B content and activation clarification

**Historical handoff clarification:** the subsequent DEC-068/069 update, now refined by DEC-070 below, resolves the login/borrower-owned password model and replaces the earlier delivery alternatives with the institutional-email link recommendation. Official content and unimplemented provider/lifecycle/recovery gates remain.

2026-10-08, explicit owner direction: **no current official FSMO Terms and Conditions document is confirmed for V2**. V1 display/capstone/draft/mockup wording remains reference content and cannot authorize institutional publication or acceptance. Implement versioning, acceptance records, first-use UI/server validation and fail-closed missing-content behavior independently; synthetic TEST terms are confined to isolated verification. The owner directs an explicit credential-delivery decision gate, with terms implemented independently. DEC-066 records accepted engineering mechanisms; DEC-067/OPEN-001/009/016 retain the unresolved approval/delivery/recovery dependencies. No content approval, activation policy, full Phase4B completion, Phase5 authorization or deployment is inferred from tests.


## Confirmed account activation and bulk provisioning update

**Historical policy checkpoint:** DEC-070 below supersedes this update’s universal institutional-email/optional-ID assumptions, earlier Staff creation authority and generic bulk population. Its secure borrower-owned password, preview/confirmation/history/audit requirements remain approved.

2026-10-08, subsequent explicit project-owner direction supersedes unresolved institutional-login/initial-password alternatives: institutional email plus a separate eLabTrack password is approved; Borrowers create their own password in secure activation, no mailbox-password collection/SSO/public registration. Recommend single-use activation links delivered to institutional email, with hash-only expiring tokens and secure invalidation/reissue; an unconfigured provider is not verified delivery. Stable internal UUIDs/normalized unique email and roster-supplied institutional student IDs apply; approved configurable domains must cover actual Student/Faculty needs. The example domain is not an approved global allowlist.

Admin Excel bulk creation and separate Excel Bulk Deactivate Borrowers are required, with validation/matching/duplicates/preview/explicit selected-account confirmation/per-record results/safe retries/durable audit. Graduation uses history-preserving deactivation, never roster-absent deactivation, Staff/Admin targeting or normal hard deletion. Outstanding-obligation deactivation procedure remains explicitly unresolved. DEC-068/069 governed this historical checkpoint; current DEC-070 and [policy](ACCOUNT_PROVISIONING_POLICY.md) supersede the conflicting identity/category/authority/bulk assumptions; OPEN-017/028/029 and residual OPEN-001/009 preserve delivery/data/activation/recovery/obligation dependencies. OPEN-016 still requires approved official FSMO wording; Phase4B remains a verified foundation, not fully complete. This policy update changes documentation only and does not authorize Phase5/6/later implementation, mail, commit/push/deployment. Historical Phase4B report/DEC-067 alternatives remain evidence of their earlier handoff, not current unresolved login policy.

## Final Student/Faculty provisioning authority

2026-10-08, latest explicit owner direction is **APPROVED PRODUCT RULES** in DEC-070: only Admin creates any account; Staff provides operational assistance. Student/Faculty are BORROWER categories, never authorization roles. Students require unique textual official Student IDs and approved SKSU institutional email and may be created individually or via Admin Student Excel import. Faculty is Admin individual-only with any valid unique accessible email and no required Student ID. Standard Excel creation/deactivation is Student-only, with studentId/name/email/course/contactNumber, no password/privilege/category assignment, complete preview before confirmation, conflict rejection, current target role/category checks and retained history. Both use borrower-chosen separate eLabTrack passwords in secure activation and must accept current officially published FSMO terms.

**REMAINING TECHNICAL DEPENDENCIES:** category/Student-ID storage, reviewed conditional validation/API/import/matching/transaction/idempotency/audit contracts and activation/ownership/recovery/configured tested delivery. **REMAINING INSTITUTIONAL APPROVALS:** actual SKSU Student domain/official roster formatting inputs (OPEN-028), activation/provider/recovery responsibility (OPEN-001/009/017), Student bulk deactivation with obligations (OPEN-029) and official FSMO terms (OPEN-016). These do not reopen approved Admin authority, Student ID or Faculty requirements. [Current policy/handoff](ACCOUNT_PROVISIONING_POLICY.md) governs affected domain/API/UX documents. Historical phase reports/V1 evidence and approved visual assets are retained unchanged. No Phase5 authorization, code/authentication/migration/UI implementation, mail, commit/push/deployment is provided by this documentation update.

## Additional owner-approved decisions and current readiness

2026-10-08, latest owner direction leaves DEC-070 Student/Faculty provisioning unchanged and adds **APPROVED PRODUCT RULES** DEC-071–073: Brevo is selected for backend-only activation and future password recovery; live delivery requires configured API key, verified sender and successful delivery testing, with no committed credentials or premature delivery claims. Official FSMO terms will be finalized after presenting the application; keep Phase4B infrastructure and do not invent/publish/automatically accept institutional content. That external approval does not block independent account-management/inventory development under separate authorization, but live borrowing still requires official publication and documented acceptance.

Admin Student deactivation, individually or Student-only bulk, is allowed for graduation/withdrawal/transfer/other authorized reasons regardless of fines, active/overdue borrowing, unreturned equipment or replacement obligations, with warnings and confirmation. No automatic clearance/waiver/payment, return, replacement resolution, unresolved-loan closure, overdue calculation change or history deletion. FSMO resolves separately; fine clearance stays Admin-only/auditable. OPEN-029 is resolved and removed from active gates; earlier paragraphs/register dependency bullets describe superseded checkpoints on these three topics, not current alternatives.

**Remaining technical dependencies:** reviewed Phase5 storage/import/matching/activation/recovery/Brevo adapter/confirmation/transaction/audit/test contracts. **Remaining external dependencies:** actual SKSU domain/official roster input (OPEN-028), backend Brevo API key/verified sender/verified live delivery plus operational ownership/recovery responsibility (OPEN-001/009/017), official terms approval after presentation (OPEN-016). Provider choice and deactivation-with-obligations policy are resolved. Product policy is ready for a scoped Phase5 plan; Phase5 remains unstarted and needs separate implementation authorization. This update changes documentation only; code, migrations, sessions, official content and approved visual assets remain unchanged.

Batch1 Phase6 implements only the reviewed catalog/inventory boundary in000007: four physical buckets, server-authoritative available adjustments, Admin reconciliation, immutable history, guarded lifecycle and bounded catalog images. See current API contracts/DEC-075 and PHASE6_REPORT.md. Older future/candidate descriptions are dated design evidence. No Phase7/return/replacement/fine-settlement workflow is authorized or implemented.
