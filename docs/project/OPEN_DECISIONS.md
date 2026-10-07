# Open decisions and source conflicts

Recorded initially 2026-10-05; consolidated for Phase 2 on 2026-10-08. These entries are investigation results, **not accepted institutional policy**. Existing OPEN-001–020 IDs/evidence are retained; OPEN-021–027 identify newly separated decision questions. V1 evidence is the supplied static audit, whose emergency patch is undeployed. No personal/sample borrower data is reproduced.

Phase 2 is **BLOCKED FOR IMPLEMENTATION** where a decision changes the state machine/schema or authoritative policy. This is a documentation handoff, not an instruction to stop independent design or begin another phase. The [20-row policy decision matrix](../domain/BUSINESS_RULES.md#policy-decision-matrix) records WHY, V1, capstone, boss reference, engineering recommendation, database impact and whether dependent implementation can proceed. NON-BLOCKING means core design can proceed; later dependent work still needs the stated stakeholder decision. No Phase 1 engineering question is reopened by this register.

**Citation correction:** direct body-paragraph checks of the current local capstone found stale locators in the original register. The policy evidence remains unchanged; current locators now identify registration conclusion P912, administrative registration description P827, Super Admin recommendation P922, Sprint 1 access/deactivation P605/607 and PHP/MySQL P606, Sprint 2 faculty/category P613, denial description P851, and email Sprint 4 P628–631. Visible section names/text govern. Audit 04/07's stored status vocabulary is Ondue/Overdue and combined Incomplete variants; earlier Due/Missing shorthand did not establish actual stored labels.

## IDENTITY / ELIGIBILITY

## ELAB-V2-OPEN-001

- **Question:** Public registration or administrative provisioning?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 01/05/13: no public signup; staff/admin create active accounts. Audit 15 adds undeployed server authorization around provisioning.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Scope/account management P421 requires approved/manual verification; conclusion P912 says borrowers can register; Add User description P827 describes admin registration. This is an internal documentation conflict as well as a V1 conflict.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Determines onboarding, approval/verification, abuse controls and routes in Phases 4–5. Template /auth/register remains inherited only.
- **Phase 1B engineering evidence (2026-10-06):** Production registration is now absent (404); only explicit development/test retains local signup. DEC-025 accepts reversible containment, not public-signup or provisioning policy. No alternative creation API/UI is implemented.
- **Required decision-maker:** FSMO accountable administrator with school project sponsor (responsibility, not an invented named approver).
- **Phase 2 gate (2026-10-08):** **BLOCKING DOMAIN IMPLEMENTATION; STAKEHOLDER CONFIRMATION REQUIRED** for dependent feature/contract.

## ELAB-V2-OPEN-002

- **Question:** Which borrower types are eligible: students only, faculty, staff?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 01/03/13: borrower UI is student-oriented; staff/admin share operational UI; no distinct faculty borrower role.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Stakeholder benefits P401–402 and Sprint 2 P613 mention students and faculty; objectives emphasize students/borrowers.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Eligibility model, profile fields, permissions and test matrix in Phases 2/4/5.
- **Required decision-maker:** FSMO policy owner with faculty/student representatives (responsibility, not an invented named approver).
- **Phase 2 gate (2026-10-08):** **BLOCKING DOMAIN IMPLEMENTATION; STAKEHOLDER CONFIRMATION REQUIRED** for dependent feature/contract.

## ELAB-V2-OPEN-003

- **Question:** How do staff and administrator privileges differ?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 03 and 15: shared operational UI and broad equivalent privileges, including creating/deleting privileged users and clearing fines; containment patch retains that equality.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Scope P421 names administrators; SOP P364 and browsing P425 name staff approval; distinct borrower/admin access described P605. No complete operation matrix is established.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Least privilege and server operation matrix must be agreed before Phase 4 roles or Phase 5 provisioning.
- **Required decision-maker:** FSMO head/accountable administrator with school IT/security owner (responsibility, not an invented named approver).
- **Phase 2 gate (2026-10-08):** **BLOCKING DOMAIN IMPLEMENTATION; STAKEHOLDER CONFIRMATION REQUIRED** for dependent feature/contract.

## ELAB-V2-OPEN-004

- **Question:** What responsibilities would a Super Administrator have?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 03/13: only student/staff/admin; no Super Admin role or supervisor workflow.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Recommendations P922 propose administrator supervision, configurations, access/privilege management, monitoring and integrity across laboratories/departments.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Potential high-privilege recovery/configuration authority; must not simply rename admin or grant unrestricted powers.
- **Required decision-maker:** Project sponsor and FSMO/IT governance owners (responsibility, not an invented named approver).
- **Phase 2 gate (2026-10-08):** **NON-BLOCKING / DEFERRED**; no extra role/platform/stack is authorized. Stakeholder confirmation is required for any proposed scope change.

## ELAB-V2-OPEN-005

- **Question:** Does a Super Admin belong in current FSMO V2 or only future expansion?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 01/10/13: globally shared FSMO inventory, no organizational model and no Super Admin.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** P922 is a recommendation tied to multiple laboratories/departments; it is not an accepted V2 requirement. Current stakeholder direction explicitly excludes campus-wide scope.
- **Current status:** Deferred; no confirmed V2 product-policy resolution.
- **Impact:** No Phase 0 implementation. Any narrowly scoped FSMO need requires an explicit decision; multi-department expansion remains a separate proposal.
- **Required decision-maker:** Project sponsor and FSMO governance owner (responsibility, not an invented named approver).
- **Phase 2 gate (2026-10-08):** **NON-BLOCKING / DEFERRED**; no extra role/platform/stack is authorized. Stakeholder confirmation is required for any proposed scope change.

## ELAB-V2-OPEN-008

- **Question:** What happens to inactive/suspended accounts and existing sessions/loans?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 03/05/09/13: login blocks non-active profile; current sessions are not centrally ejected; status types differ. Audit 15 denies inactive callers at its HTTP boundary.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Scope P421 requires approved users; Sprint 1 P605/719 describes deactivation, without session revocation/active-loan policy.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Phase 4 session checks/revocation and Phase 5 eligibility; preserve outstanding accountability.
- **Required decision-maker:** FSMO policy owner and IT/security owner (responsibility, not an invented named approver).

- **Phase 1B engineering evidence (2026-10-06):** Existing generic is_active=false denies subsequent protected API requests with 403; current role is looked up each request. No suspended-state taxonomy, session-family/logout revocation or loan/accountability policy was added. Those questions remain unresolved.
- **Phase 2 gate (2026-10-08):** **BLOCKING DOMAIN IMPLEMENTATION; STAKEHOLDER CONFIRMATION REQUIRED** for dependent feature/contract.

## ELAB-V2-OPEN-009

- **Question:** Is email verification required, and at what point?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 05: accounts created emailVerified=false; no send/check verification flow. Administrative verification is not verified email ownership.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Scope P424 describes manual credentials verification; P427–428 emphasize accurate reachable email. No explicit email-ownership verification gate is established.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Phase 4 onboarding/recovery and Phase 9 deliverability; do not confuse email verification with borrowing eligibility.
- **Required decision-maker:** FSMO administrator and IT/security owner (responsibility, not an invented named approver).

- **Phase 1B engineering evidence (2026-10-06):** Self-service email/identity changes are temporarily unavailable; only display-name editing remains. This containment avoids inventing verification/provisioning behavior and does not settle email ownership/eligibility policy.
- **Phase 2 gate (2026-10-08):** **BLOCKING DOMAIN IMPLEMENTATION; STAKEHOLDER CONFIRMATION REQUIRED** for dependent feature/contract.

## ELAB-V2-OPEN-027

- **Question:** What institutional concurrent-session, global/family logout and immediate access-revocation behavior is required?
- **ACTUAL V1 BEHAVIOR:** Audit 05/15: client login/status and undeployed current-caller checks; no family/global policy.
- **CAPSTONE INTENT:** Scope P421 and deactivation Sprint 1 P605/607 do not establish global session/device rules.
- **BOSS DESIGN REFERENCE / ENGINEERING RECOMMENDATION:** Current Phase 1A–1I has verified single-use rotation/current account/local-browser coordination; no family/cross-device logout guarantee. Boss persisted-session fencing is a reference, not an approved replacement.
- **Current status:** Needs Stakeholder Input; no institutional policy accepted.
- **Phase 2 gate:** **NON-BLOCKING core domain design; BLOCKING dependent notification/retention/identity feature; STAKEHOLDER CONFIRMATION REQUIRED.**
- **Impact:** Later approved identity policy; no Phase 1 redesign here. Transaction-time account freshness is separately recommended for business mutations.
- **Required decision-maker:** FSMO administrator and school IT/security owner (responsibility, not an invented named approver).


## BORROWING

## ELAB-V2-OPEN-020

- **Question:** Which borrowing limits, reservation point and direct-checkout rules are approved?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 07/09: requests reserve stock immediately, student due date within seven days; staff direct checkout differs; no evidenced concurrent-loan limits or fine-based block.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** SOP P364 says release after approval; scope P425 requires staff approval; reservation timing, limits and direct-checkout exceptions are not fully specified.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Phase 2 stock/lifecycle invariants and Phase 7 concurrency and request/checkout rules.
- **Required decision-maker:** FSMO policy owner and operational staff (responsibility, not an invented named approver).
- **Phase 2 gate (2026-10-08):** **BLOCKING DOMAIN IMPLEMENTATION; STAKEHOLDER CONFIRMATION REQUIRED** for dependent feature/contract.
- **Phase 2 clarification:** Select request-time hold (A), approval-time hold (B), or justify a hybrid (C); decide all-quantity issuance, direct-checkout permission/exception, duration and any borrower limits. Recommendations are compared in [BUSINESS_RULES](../domain/BUSINESS_RULES.md#reservation-options). No numerical limit/policy is accepted.

## ELAB-V2-OPEN-021

- **Question:** Does staff approval record a permission decision, actual physical release, or two separate actions?
- **ACTUAL V1 BEHAVIOR:** Audit 07: Request→Ongoing, borrowedDate reset, no separately recorded pickup/release.
- **CAPSTONE INTENT:** SOP P364 and scope P425 release after staff approval; whether there is a waiting interval is not specified.
- **BOSS DESIGN REFERENCE / ENGINEERING RECOMMENDATION:** Boss audit §5 merges pending→checked_out. Recommend separate approved state and issue timestamps; combined command only if explicitly accepted.
- **Current status:** Needs Stakeholder Input; no institutional policy accepted.
- **Phase 2 gate:** **BLOCKING DOMAIN IMPLEMENTATION; STAKEHOLDER CONFIRMATION REQUIRED.**
- **Impact:** State machine, API actions, inventory timing and actor/issue snapshot fields fundamentally differ.
- **Required decision-maker:** FSMO operations/policy owner (responsibility, not an invented named approver).

## ELAB-V2-OPEN-013

- **Question:** May a borrower or staff cancel a request/active transaction?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 07/13: borrower cancellation not found; staff can hard-delete transactions with stock restoration edge cases and no history.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** SOP P364 describes submission/approval/release/archive, not a defined cancellation lifecycle. Form Cancel buttons are not cancellation-policy evidence.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Phases 2/7 lifecycle, actor permissions, stock release, audit trail and retries.
- **Required decision-maker:** FSMO operations/policy owner (responsibility, not an invented named approver).
- **Phase 2 gate (2026-10-08):** **BLOCKING DOMAIN IMPLEMENTATION; STAKEHOLDER CONFIRMATION REQUIRED** for dependent feature/contract.

## ELAB-V2-OPEN-014

- **Question:** Must denial include a reason and retained history?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 07/08/09: denial restores stock, queues notification and deletes request; no durable denial record/reason.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Denial email figure/description P851 communicates denial/contact guidance, but does not specify retained reasons or structured decision history.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Phase 7 evidence retention, borrower feedback, staff accountability and notification payload.
- **Required decision-maker:** FSMO operations/policy owner (responsibility, not an invented named approver).
- **Phase 2 gate (2026-10-08):** **NON-BLOCKING core domain design; STAKEHOLDER CONFIRMATION REQUIRED** before the dependent feature or operational procedure.

## ELAB-V2-OPEN-015

- **Question:** What is the canonical due/overdue vocabulary and day-boundary policy?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 04/07/09: the stored audit vocabulary is Request/Ongoing/Ondue/Overdue/Incomplete and combined variants (audit 04/07), rather than the earlier Due/Missing shorthand; Asia/Manila scheduler/calendar rules and mixed timestamp arithmetic.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Scope P427/430 and reminder figures use due/overdue; they do not establish exact canonical states or transition times.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Phases 2/7/8/9 shared lifecycle vocabulary and date calculation. Client timezone must not implicitly define policy.
- **Required decision-maker:** FSMO policy owner, with engineering for implementation semantics (responsibility, not an invented named approver).
- **Phase 2 gate (2026-10-08):** **BLOCKING DOMAIN IMPLEMENTATION; STAKEHOLDER CONFIRMATION REQUIRED** for dependent feature/contract.
- **Phase 2 clarification:** Choose date-only local deadline versus exact instant, FSMO business timezone/cutoff and due/overdue/grace meanings. Asia/Manila is a V1-grounded recommendation only; developer Asia/Shanghai is not product policy.

## ELAB-V2-OPEN-016

- **Question:** How are terms versions, acceptance evidence and reacceptance handled?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 04/09/13: static terms UI; acceptance boolean and timestamp only, no version/content hash.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** SOP and borrowing intent require accountability, but engineering passages reviewed do not define versioned acceptance or policy-change handling.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Phases 2/5/7 consent evidence, retention and eligibility. No legal terms text is invented.
- **Required decision-maker:** FSMO policy owner and school policy/document custodian (responsibility, not an invented named approver).
- **Phase 2 gate (2026-10-08):** **BLOCKING DOMAIN IMPLEMENTATION; STAKEHOLDER CONFIRMATION REQUIRED** for dependent feature/contract.
- **Phase 2 clarification:** Select once-per-account/version versus per-borrowing acceptance and effective version at request/issue, approved text/publication authority and assisted acceptance evidence. Final uniqueness/consent predicates cannot be migrated before the decision.

## ELAB-V2-OPEN-011

- **Question:** Are equipment categories required, and who maintains them?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 06/13: name/search/status discovery; no category field or category management.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Browsing P425 and Sprint 2 P613 describe category/availability filtering.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Phase 2 schema and Phase 6 filters; requires taxonomy and lifecycle ownership before category implementation.
- **Required decision-maker:** FSMO inventory custodian and administrator (responsibility, not an invented named approver).
- **Phase 2 gate (2026-10-08):** **NON-BLOCKING core domain design; STAKEHOLDER CONFIRMATION REQUIRED** before the dependent feature or operational procedure.


## RETURNS

## ELAB-V2-OPEN-022

- **Question:** Which return/disposition evidence, repair/recovery/retirement/write-off and correction procedures are authorized?
- **ACTUAL V1 BEHAVIOR:** Audit 06/07/09: partial good/damaged/lost quantities, notes required in UI; damage/loss not explicit stock buckets; no evidenced repair/write-off or return-correction workflow.
- **CAPSTONE INTENT:** Scope P424/426 and definitions P441/443 require manual condition inspection/accountability, without detailed disposition procedures.
- **BOSS DESIGN REFERENCE / ENGINEERING RECOMMENDATION:** Boss audit §§6/14: buckets/reclassification useful, notes optional; duplicate-return defect must be rejected. Recommend immutable unique-line return events, required damage/loss reasons and explicit movement corrections; no invented backdating or generic edit.
- **Current status:** Needs Stakeholder Input; no institutional policy accepted.
- **Phase 2 gate:** **BLOCKING DOMAIN IMPLEMENTATION; STAKEHOLDER CONFIRMATION REQUIRED.**
- **Impact:** Bucket meaning/movement authority, event validation, assessment/correction policy. Core defect prevention is required regardless of the selected policy.
- **Required decision-maker:** FSMO inventory/operations and accountability owners (responsibility, not an invented named approver).


## ACCOUNTABILITY

## ELAB-V2-OPEN-006

- **Question:** What fine amount, currency, grace period, cap and damage/loss rules apply?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 07/09: PHP 10/day overdue, unit-price damage/loss charges, client/timezone rounding; no evidenced cap or grace period.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Scope P430–431 and equipment definition P441 intend fines/collection/accountability; engineering descriptions do not establish a complete approved amount/rounding/waiver policy.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Phase 2 model and Phase 8 calculations, day boundaries and tests. V1 amounts are evidence, not confirmed V2 policy.
- **Required decision-maker:** FSMO policy owner with finance/accountability authority (responsibility, not an invented named approver).
- **Phase 2 gate (2026-10-08):** **BLOCKING DOMAIN IMPLEMENTATION; STAKEHOLDER CONFIRMATION REQUIRED** for dependent feature/contract.
- **Phase 2 clarification:** Confirm currency, day/cutoff/grace/cap/rounding and damage/loss value formula. V1 PHP 10/day is factual evidence only. Posting and binding timing are OPEN-023/024; no duplicate financial truth is proposed.

## ELAB-V2-OPEN-007

- **Question:** How are fine assessment, payment, clearing and waivers distinguished?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 04/07/08/13: records/fines/finePaid representations diverge; UI clearing changes records without reconciling fines; no coherent payment ledger.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Dashboard/report scope P430–431 calls for payment records and fines collection; reporting intent alone does not define settlement or collection procedures.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Phase 8 balances/ledger and Phase 10 revenue accuracy; payment assessment must not be reported as collected cash.
- **Required decision-maker:** FSMO finance/accountability owner and administrator (responsibility, not an invented named approver).
- **Phase 2 gate (2026-10-08):** **BLOCKING DOMAIN IMPLEMENTATION; STAKEHOLDER CONFIRMATION REQUIRED** for dependent feature/contract.
- **Phase 2 clarification:** Decide administrative discharge versus actual cash/receipt tracking; identify assessment/waiver/reduction/increase/settlement/reversal actors and reason requirements. Conditional payments/allocations are evaluated, not approved. Assessments/settlements cannot establish revenue.

## ELAB-V2-OPEN-023

- **Question:** When do equipment value, borrower presentation, due policy and governing terms become binding for a borrowing?
- **ACTUAL V1 BEHAVIOR:** Audit 04/07: price/name snapshots captured at creation; approval resets borrowedDate; totalPrice semantics uncertain.
- **CAPSTONE INTENT:** P364/421/425/441 accountability/release intent, no price-lock or policy-change rule.
- **BOSS DESIGN REFERENCE / ENGINEERING RECOMMENDATION:** Boss audit §7 captures request price while sometimes described as checkout snapshot. Recommend explicit request quotation and immutable issue evidence; effective version/value at issue unless stakeholder selects request binding.
- **Current status:** Needs Stakeholder Input; no institutional policy accepted.
- **Phase 2 gate:** **BLOCKING DOMAIN IMPLEMENTATION; STAKEHOLDER CONFIRMATION REQUIRED.**
- **Impact:** Snapshot fields/timing, policy/terms references and liability basis; approval-time binding may require model revision.
- **Required decision-maker:** FSMO policy/inventory/accountability and document custodians (responsibility, not an invented named approver).

## ELAB-V2-OPEN-024

- **Question:** At what event/window are overdue, damage and loss charges posted, and how are partial returns, zero balance and reversals handled?
- **ACTUAL V1 BEHAVIOR:** Audit 07/08/09: mutable active fine and cumulative final assessment; divergent clearing/payment representations.
- **CAPSTONE INTENT:** P430/431 fines/collection reports, no detailed posting/day/partial/correction policy.
- **BOSS DESIGN REFERENCE / ENGINEERING RECOMMENDATION:** Boss audit §14 final-only charges, active estimates and zero-status inconsistency. Recommend one immutable assessment per source/window, append-only changes, separate unposted estimate and derived balance; no duplicate cumulative posting.
- **Current status:** Needs Stakeholder Input; no institutional policy accepted.
- **Phase 2 gate:** **BLOCKING DOMAIN IMPLEMENTATION; STAKEHOLDER CONFIRMATION REQUIRED.**
- **Impact:** Charge source uniqueness/window, return transaction assessment, estimate projections and permitted reversal/reopening commands.
- **Required decision-maker:** FSMO finance/accountability and policy owners (responsibility, not an invented named approver).


## NOTIFICATIONS

## ELAB-V2-OPEN-017

- **Question:** Which email delivery provider/consumer and retry/receipt mechanism are approved?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 08/12/13: two notification shapes; queue producers present, consumer/extension absent; deployed delivery cannot be inferred.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Scope P427–428 and Sprint 4 P628–631 intend confirmations, reminders, delivery logs and failure handling; screenshots imply received mail but do not identify the operational consumer/provider.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Phase 9 bounded scheduled work, delivery failures/retries/deduplication and operations; no message broker is implied.
- **Required decision-maker:** School IT/email custodian and FSMO operational owner (responsibility, not an invented named approver).
- **Phase 2 gate (2026-10-08):** **NON-BLOCKING core domain design; STAKEHOLDER CONFIRMATION REQUIRED** before the dependent feature or operational procedure.

## ELAB-V2-OPEN-025

- **Question:** Which notifications, recipients and reminder cadence are institutionally required, and does grace change reminder wording/timing?
- **ACTUAL V1 BEHAVIOR:** Audit 08: approval/denial/receipt and once-per-transaction due/tomorrow/overdue producer flags; delivery consumer absent.
- **CAPSTONE INTENT:** P427/428 and Sprint 4–5 P628–636 require confirmations/reminders, logged delivery/failure intent, not a provider/cadence guarantee.
- **BOSS DESIGN REFERENCE / ENGINEERING RECOMMENDATION:** Boss audit §16 repeats partial-return keys and uses daily overdue/grace-inconsistent scheduler. Recommend per-event keys, approved schedule windows, bounded fenced worker and attempts; provider/operations remain OPEN-017.
- **Current status:** Needs Stakeholder Input; no institutional policy accepted.
- **Phase 2 gate:** **NON-BLOCKING core domain design; BLOCKING dependent notification/retention/identity feature; STAKEHOLDER CONFIRMATION REQUIRED.**
- **Impact:** Core outbox is independent; triggers/templates/calendar scheduling/recipient handling and delivery operations gate Phase 9.
- **Required decision-maker:** FSMO operational owner and IT/email custodian (responsibility, not an invented named approver).


## KIOSK

## ELAB-V2-OPEN-010

- **Question:** How does kiosk authentication, idle timeout and session cleanup work?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 01/03/13: shared responsive web/student flows; no separately hardened kiosk route/session lifecycle evidenced.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Scope P429 and definition P440 describe self-service browsing/requests on physical kiosk and personal devices; authentication, idle expiry and handoff cleanup are unspecified.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Phase 11 shared-device privacy, request attribution, session isolation and operational acceptance. No anonymous borrowing is assumed.
- **Required decision-maker:** FSMO operational owner with IT/security and kiosk operators (responsibility, not an invented named approver).
- **Phase 2 gate (2026-10-08):** **NON-BLOCKING core domain design; STAKEHOLDER CONFIRMATION REQUIRED** before the dependent feature or operational procedure.


## DATA RETENTION / MIGRATION

## ELAB-V2-OPEN-012

- **Question:** What are equipment deletion, archival and retention rules?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 06/07/09/13: hard deletion with client-side active-transaction guard and concurrency/status gaps; archive not found; history depends on snapshots.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Scope P423 requires inventory updates and monitoring; Add Equipment description demonstrates entry, but no archive/retention policy is established.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Phase 2 references/history and Phase 6 lifecycle. Confirmed V2 history-preservation principle requires safe design; exact archival policy remains open.
- **Required decision-maker:** FSMO inventory/accountability owner (responsibility, not an invented named approver).
- **Phase 2 gate (2026-10-08):** **NON-BLOCKING core domain design; STAKEHOLDER CONFIRMATION REQUIRED** before the dependent feature or operational procedure.

## ELAB-V2-OPEN-018

- **Question:** What V1 data should migrate, reconcile or be excluded?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 04/07/14: users/Auth IDs, aggregate equipment, active transactions, snapshots, conflicting fines, notifications and images; live volumes/rules/state not inspected.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Product history/reporting intent P422/431 does not define cutover, historical retention or migration scope.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Phase 12 read-only discovery authorization, mapping/reconciliation, dry run, rollback and signoff. Phase 0 accesses no Firebase.
- **Required decision-maker:** FSMO data owner and IT/Firebase custodian (responsibility, not an invented named approver).
- **Phase 2 gate (2026-10-08):** **NON-BLOCKING core domain design; STAKEHOLDER CONFIRMATION REQUIRED** before the dependent feature or operational procedure.
- **Phase 2 clarification:** Unknown past returns/actors/dates/consent/payment are not reconstructable facts. Any explicitly labeled cumulative baseline or nullable historical-time accommodation needs a reviewed migration-specific schema/invariant amendment; otherwise quarantine. No importer or live discovery is authorized here.

## ELAB-V2-OPEN-026

- **Question:** What privacy categories, retention, anonymization, evidence access and legal/institutional data procedures apply?
- **ACTUAL V1 BEHAVIOR:** Audit 04/05/08/14: PII snapshots, public profile images, retained history and unbounded exposure; no retention schedule evidenced.
- **CAPSTONE INTENT:** P422/431 history/reporting intent, no legal retention duration or anonymization procedure established.
- **BOSS DESIGN REFERENCE / ENGINEERING RECOMMENDATION:** Boss audit §§15/20: audit/provenance concepts; retention unresolved. Recommend restrictive refs, private owner projections, explicit approved evidence-preserving procedure; no made-up period.
- **Current status:** Needs Stakeholder Input; no institutional policy accepted.
- **Phase 2 gate:** **NON-BLOCKING core domain design; BLOCKING dependent notification/retention/identity feature; STAKEHOLDER CONFIRMATION REQUIRED.**
- **Impact:** Retention/export/access/privacy handling and cutover. Do not anonymize away custody/balances/source evidence; normal transactional deletion remains excluded.
- **Required decision-maker:** FSMO data/document/accountability custodian and school IT/privacy authority (responsibility, not an invented named approver).

## ELAB-V2-OPEN-019

- **Question:** How should contradictory legacy technology/schema descriptions be treated?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 01/04/12: Expo/React Native + Firebase, Firestore collections and Functions; source is not supplied in this V2 workspace.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Sprint 1 P606 describes PHP/MySQL; database tables use relational types. These contradict actual audited V1 architecture.
- **Current status:** Deferred; no confirmed V2 product-policy resolution.
- **Impact:** No V2 stack question: current Go/Fiber/React/PostgreSQL direction is accepted. Verify audit revision/source during later authorized migration discovery; do not copy document tables literally.
- **Required decision-maker:** Engineering lead and V1 technical custodian (responsibility, not an invented named approver).
- **Phase 2 gate (2026-10-08):** **NON-BLOCKING / DEFERRED**; no extra role/platform/stack is authorized. Stakeholder confirmation is required for any proposed scope change.

## Proposal handling

Any suggested policy remains ENGINEERING RECOMMENDATION / NEW V2 PROPOSAL until linked to an Accepted decision with stakeholder authority/date. Confirmed FSMO scope overrides expansion recommendations. Documented alternatives are not a generic runtime configuration requirement. Confirm the operation matrix, eligibility, approval/release, reservation, due/terms/value binding, dispositions and accountability before their dependent implementation. Category naming, provider choice, kiosk protocol and retention/cutover details gate their later activities. Phase 1A–1I remains complete; current source/verification is in PHASE1_FOUNDATION.md.
