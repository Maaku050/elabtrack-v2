# Open decisions and source conflicts

Recorded 2026-10-05. These entries are investigation results, not accepted policy. V1 evidence is the supplied static audit; its emergency patch is recorded as undeployed. Capstone locators follow `SOURCE_OF_TRUTH.md`; no personal/sample borrower data is reproduced. Decisions must be confirmed by an accountable stakeholder before dependent implementation. “Not established” means the engineering passages do not settle the question.

## ELAB-V2-OPEN-001

- **Question:** Public registration or administrative provisioning?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 01/05/13: no public signup; staff/admin create active accounts. Audit 15 adds undeployed server authorization around provisioning.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Scope/account management P421 requires approved/manual verification; conclusion P1527 says borrowers can register; Add User description P1316 describes admin registration. This is an internal documentation conflict as well as a V1 conflict.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Determines onboarding, approval/verification, abuse controls and routes in Phases 4–5. Template /auth/register remains inherited only.
- **Required decision-maker:** FSMO accountable administrator with school project sponsor (responsibility, not an invented named approver).

## ELAB-V2-OPEN-002

- **Question:** Which borrower types are eligible: students only, faculty, staff?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 01/03/13: borrower UI is student-oriented; staff/admin share operational UI; no distinct faculty borrower role.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Stakeholder benefits P401–402 and Sprint 2 P725 mention students and faculty; objectives emphasize students/borrowers.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Eligibility model, profile fields, permissions and test matrix in Phases 2/4/5.
- **Required decision-maker:** FSMO policy owner with faculty/student representatives (responsibility, not an invented named approver).

## ELAB-V2-OPEN-003

- **Question:** How do staff and administrator privileges differ?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 03 and 15: shared operational UI and broad equivalent privileges, including creating/deleting privileged users and clearing fines; containment patch retains that equality.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Scope P421 names administrators; SOP P364 and browsing P425 name staff approval; distinct borrower/admin access described P717. No complete operation matrix is established.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Least privilege and server operation matrix must be agreed before Phase 4 roles or Phase 5 provisioning.
- **Required decision-maker:** FSMO head/accountable administrator with school IT/security owner (responsibility, not an invented named approver).

## ELAB-V2-OPEN-004

- **Question:** What responsibilities would a Super Administrator have?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 03/13: only student/staff/admin; no Super Admin role or supervisor workflow.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Recommendations P1537 propose administrator supervision, configurations, access/privilege management, monitoring and integrity across laboratories/departments.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Potential high-privilege recovery/configuration authority; must not simply rename admin or grant unrestricted powers.
- **Required decision-maker:** Project sponsor and FSMO/IT governance owners (responsibility, not an invented named approver).

## ELAB-V2-OPEN-005

- **Question:** Does a Super Admin belong in current FSMO V2 or only future expansion?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 01/10/13: globally shared FSMO inventory, no organizational model and no Super Admin.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** P1537 is a recommendation tied to multiple laboratories/departments; it is not an accepted V2 requirement. Current stakeholder direction explicitly excludes campus-wide scope.
- **Current status:** Deferred; no confirmed V2 product-policy resolution.
- **Impact:** No Phase 0 implementation. Any narrowly scoped FSMO need requires an explicit decision; multi-department expansion remains a separate proposal.
- **Required decision-maker:** Project sponsor and FSMO governance owner (responsibility, not an invented named approver).

## ELAB-V2-OPEN-006

- **Question:** What fine amount, currency, grace period, cap and damage/loss rules apply?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 07/09: PHP 10/day overdue, unit-price damage/loss charges, client/timezone rounding; no evidenced cap or grace period.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Scope P430–431 and equipment definition P441 intend fines/collection/accountability; engineering descriptions do not establish a complete approved amount/rounding/waiver policy.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Phase 2 model and Phase 8 calculations, day boundaries and tests. V1 amounts are evidence, not confirmed V2 policy.
- **Required decision-maker:** FSMO policy owner with finance/accountability authority (responsibility, not an invented named approver).

## ELAB-V2-OPEN-007

- **Question:** How are fine assessment, payment, clearing and waivers distinguished?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 04/07/08/13: records/fines/finePaid representations diverge; UI clearing changes records without reconciling fines; no coherent payment ledger.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Dashboard/report scope P430–431 calls for payment records and fines collection; reporting intent alone does not define settlement or collection procedures.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Phase 8 balances/ledger and Phase 10 revenue accuracy; payment assessment must not be reported as collected cash.
- **Required decision-maker:** FSMO finance/accountability owner and administrator (responsibility, not an invented named approver).

## ELAB-V2-OPEN-008

- **Question:** What happens to inactive/suspended accounts and existing sessions/loans?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 03/05/09/13: login blocks non-active profile; current sessions are not centrally ejected; status types differ. Audit 15 denies inactive callers at its HTTP boundary.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Scope P421 requires approved users; Sprint 1 P717/719 describes deactivation, without session revocation/active-loan policy.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Phase 4 session checks/revocation and Phase 5 eligibility; preserve outstanding accountability.
- **Required decision-maker:** FSMO policy owner and IT/security owner (responsibility, not an invented named approver).

## ELAB-V2-OPEN-009

- **Question:** Is email verification required, and at what point?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 05: accounts created emailVerified=false; no send/check verification flow. Administrative verification is not verified email ownership.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Scope P424 describes manual credentials verification; P427–428 emphasize accurate reachable email. No explicit email-ownership verification gate is established.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Phase 4 onboarding/recovery and Phase 9 deliverability; do not confuse email verification with borrowing eligibility.
- **Required decision-maker:** FSMO administrator and IT/security owner (responsibility, not an invented named approver).

## ELAB-V2-OPEN-010

- **Question:** How does kiosk authentication, idle timeout and session cleanup work?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 01/03/13: shared responsive web/student flows; no separately hardened kiosk route/session lifecycle evidenced.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Scope P429 and definition P440 describe self-service browsing/requests on physical kiosk and personal devices; authentication, idle expiry and handoff cleanup are unspecified.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Phase 11 shared-device privacy, request attribution, session isolation and operational acceptance. No anonymous borrowing is assumed.
- **Required decision-maker:** FSMO operational owner with IT/security and kiosk operators (responsibility, not an invented named approver).

## ELAB-V2-OPEN-011

- **Question:** Are equipment categories required, and who maintains them?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 06/13: name/search/status discovery; no category field or category management.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Browsing P425 and Sprint 2 P725 describe category/availability filtering.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Phase 2 schema and Phase 6 filters; requires taxonomy and lifecycle ownership before category implementation.
- **Required decision-maker:** FSMO inventory custodian and administrator (responsibility, not an invented named approver).

## ELAB-V2-OPEN-012

- **Question:** What are equipment deletion, archival and retention rules?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 06/07/09/13: hard deletion with client-side active-transaction guard and concurrency/status gaps; archive not found; history depends on snapshots.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Scope P423 requires inventory updates and monitoring; Add Equipment description demonstrates entry, but no archive/retention policy is established.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Phase 2 references/history and Phase 6 lifecycle. Confirmed V2 history-preservation principle requires safe design; exact archival policy remains open.
- **Required decision-maker:** FSMO inventory/accountability owner (responsibility, not an invented named approver).

## ELAB-V2-OPEN-013

- **Question:** May a borrower or staff cancel a request/active transaction?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 07/13: borrower cancellation not found; staff can hard-delete transactions with stock restoration edge cases and no history.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** SOP P364 describes submission/approval/release/archive, not a defined cancellation lifecycle. Form Cancel buttons are not cancellation-policy evidence.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Phases 2/7 lifecycle, actor permissions, stock release, audit trail and retries.
- **Required decision-maker:** FSMO operations/policy owner (responsibility, not an invented named approver).

## ELAB-V2-OPEN-014

- **Question:** Must denial include a reason and retained history?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 07/08/09: denial restores stock, queues notification and deletes request; no durable denial record/reason.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Denial email figure/description P1342 communicates denial/contact guidance, but does not specify retained reasons or structured decision history.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Phase 7 evidence retention, borrower feedback, staff accountability and notification payload.
- **Required decision-maker:** FSMO operations/policy owner (responsibility, not an invented named approver).

## ELAB-V2-OPEN-015

- **Question:** What is the canonical due/overdue vocabulary and day-boundary policy?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 07/09: Request, Ongoing, Due, Missing and terminal snapshots; V1 Missing functions as overdue; Asia/Manila scheduler/calendar rules and mixed timestamp arithmetic.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Scope P427/430 and reminder figures use due/overdue; they do not establish exact canonical states or transition times.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Phases 2/7/8/9 shared lifecycle vocabulary and date calculation. Client timezone must not implicitly define policy.
- **Required decision-maker:** FSMO policy owner, with engineering for implementation semantics (responsibility, not an invented named approver).

## ELAB-V2-OPEN-016

- **Question:** How are terms versions, acceptance evidence and reacceptance handled?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 04/09/13: static terms UI; acceptance boolean and timestamp only, no version/content hash.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** SOP and borrowing intent require accountability, but engineering passages reviewed do not define versioned acceptance or policy-change handling.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Phases 2/5/7 consent evidence, retention and eligibility. No legal terms text is invented.
- **Required decision-maker:** FSMO policy owner and school policy/document custodian (responsibility, not an invented named approver).

## ELAB-V2-OPEN-017

- **Question:** Which email delivery provider/consumer and retry/receipt mechanism are approved?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 08/12/13: two notification shapes; queue producers present, consumer/extension absent; deployed delivery cannot be inferred.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Scope P427–428 and Sprint 4 P740–743 intend confirmations, reminders, delivery logs and failure handling; screenshots imply received mail but do not identify the operational consumer/provider.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Phase 9 bounded scheduled work, delivery failures/retries/deduplication and operations; no message broker is implied.
- **Required decision-maker:** School IT/email custodian and FSMO operational owner (responsibility, not an invented named approver).

## ELAB-V2-OPEN-018

- **Question:** What V1 data should migrate, reconcile or be excluded?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 04/07/14: users/Auth IDs, aggregate equipment, active transactions, snapshots, conflicting fines, notifications and images; live volumes/rules/state not inspected.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Product history/reporting intent P422/431 does not define cutover, historical retention or migration scope.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Phase 12 read-only discovery authorization, mapping/reconciliation, dry run, rollback and signoff. Phase 0 accesses no Firebase.
- **Required decision-maker:** FSMO data owner and IT/Firebase custodian (responsibility, not an invented named approver).

## ELAB-V2-OPEN-019

- **Question:** How should contradictory legacy technology/schema descriptions be treated?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 01/04/12: Expo/React Native + Firebase, Firestore collections and Functions; source is not supplied in this V2 workspace.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Sprint 1 P718 describes PHP/MySQL; database tables use relational types. These contradict actual audited V1 architecture.
- **Current status:** Deferred; no confirmed V2 product-policy resolution.
- **Impact:** No V2 stack question: current Go/Fiber/React/PostgreSQL direction is accepted. Verify audit revision/source during later authorized migration discovery; do not copy document tables literally.
- **Required decision-maker:** Engineering lead and V1 technical custodian (responsibility, not an invented named approver).

## ELAB-V2-OPEN-020

- **Question:** Which borrowing limits, reservation point and direct-checkout rules are approved?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 07/09: requests reserve stock immediately, student due date within seven days; staff direct checkout differs; no evidenced concurrent-loan limits or fine-based block.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** SOP P364 says release after approval; scope P425 requires staff approval; reservation timing, limits and direct-checkout exceptions are not fully specified.
- **Current status:** Needs Stakeholder Input; no confirmed V2 product-policy resolution.
- **Impact:** Phase 2 stock/lifecycle invariants and Phase 7 concurrency and request/checkout rules.
- **Required decision-maker:** FSMO policy owner and operational staff (responsibility, not an invented named approver).

## Proposal handling

Any solution suggested during later planning is a **NEW V2 PROPOSAL** until linked to an Accepted decision with stakeholder authority/date. The confirmed FSMO scope overrides expansion recommendations. Questions can remain open at Phase 0; only their dependent later work is gated. Authoritative module identity is a separate foundation blocker in ELAB-V2-DEC-020.
