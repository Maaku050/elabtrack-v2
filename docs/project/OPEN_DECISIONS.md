# Open decisions and resolved source conflicts

**Current account-policy gates, 2026-10-08:** DEC-070 confirms Admin-only creation, required unique textual Student IDs and approved SKSU Student email, Faculty individual creation with any valid unique accessible email and no required Student ID, and Student-only bulk creation/deactivation. Both categories use separate borrower-chosen eLabTrack passwords and Phase4B official terms; no SSO/mailbox passwords. OPEN-001/009 retain activation/ownership/recovery implementation, DEC-071 selects Brevo; OPEN-017 retains API-key/sender/configuration and tested delivery, OPEN-028 retains exact Student domain/roster formatting and technical matching. DEC-073 resolves OPEN-029: warned/confirmed Student deactivation is allowed regardless of outstanding obligations, which remain unchanged. DEC-072 retains OPEN-016 official terms approval after presenting the application to FSMO, gating official publication/live borrowing but not independent account-management or inventory development. These dependencies do not reopen the approved product rules. Phase4B remains a verified foundation; Batch 1 now authorizes Phase5 and conditional Phase6; [implemented technical contracts](../API_CONTRACTS.md#phase-5-implemented-account-management-contracts) distinguish verified engineering from external gates. See [current policy](ACCOUNT_PROVISIONING_POLICY.md) and [historical Phase4B report](PHASE4B_REPORT.md).


Reconciled **2026-10-08, Phase 2.5**. All 29 owner working decisions are recorded in DEC-051–061; selected stock engineering baseline DEC-062. **No remaining policy materially blocks core borrower/review/return/replacement/fine/role-navigation UX.** NON-BLOCKING means core design/Phase 3A planning can proceed when separately authorized, not authorization to send mail, migrate, destroy records or deploy. Original 27 IDs and source evidence are retained; resolved questions no longer carry the historical Phase 2 BLOCKED gate. No Phase 1 architecture is changed.

## Current residual details

| Detail | IDs | Responsible review / dependent gate |
|---|---|---|
| Activation lifecycle, expiry/reissue, mailbox verification/recovery and existing-account treatment | OPEN-001/009 | FSMO/IT; login/borrower-owned password and Student/Faculty requirements resolved in DEC-068/070; hash-only activation/expiry/reissue are implemented with bounded technical defaults; institutional ownership/recovery/existing-account procedures remain |
| Approved SKSU Student domain configuration, official Student-ID format inputs and optional roster field/matching details | OPEN-028 | Institutional inputs; implemented exact text-ID/email matching and bounded template contract; Student ID is already required/unique/textual, Faculty has no institutional-domain or required Student-ID gate |
| Official terms approval after application presentation | OPEN-016 | Owner/FSMO; gates official publication/consent and live borrowing, not independent account-management/inventory development (DEC-072) |
| Taxonomy/content, ordinary inventory action details | OPEN-011/012 | FSMO catalog/inventory owner; seed/content and specialized inventory actions |
| Damaged-original disposition / equivalent acceptance guidance | OPEN-022 | FSMO inventory owner; disposal/repair feature. Core flow uses nonusable held originals and staff discretion |
| Selected Brevo integration/API key/verified sender/live delivery; later notification cadence | OPEN-017/025 | Provider selection resolved DEC-071; backend adapter and fake-submission tests implemented; real configuration/tested delivery and future recovery remain; later scheduled business messages stay separate |
| Legacy discovery/cutover and evidence gaps | OPEN-018/019 | Authorized V1 custodian/FSMO; real import/reconciliation |
| Institutional retention/privacy | OPEN-026 | Institutional/legal custodian; destructive cleanup and retention operations |
| Institutional session policy / final operations | OPEN-027 and Phase 14 roadmap | FSMO/IT; optional future session policy and deployment readiness, existing foundation preserved |

The ceiling24h interpretation is adopted as an engineering detail; PHP 10/day is not open. Original registration, population, role distinction, reservation/release, seven-day duration, terms cadence, damage-price liability, generic fine settlement and device-kiosk questions are resolved/superseded. Final deployment environment/backups/operational ownership remain Phase 14 work, not a core UX gate.

## Original register IDs and disposition

Historical evidence labels below retain factual V1 versus intended legacy/capstone distinctions. They do not override current owner direction. V1 audit15 patch is undeployed. Capstone paragraph references use corrected direct body-paragraph locators P912/P827/P922/P605–607/P613/P851/P628–631, not personal/sample records.

## ELAB-V2-OPEN-001

**Historical Phase 1 engineering evidence; later product resolutions below govern:**

- **Phase 1B engineering evidence (2026-10-06):** Production registration is now absent (404); only explicit development/test retains local signup. DEC-025 accepts reversible containment, not public-signup or provisioning policy. No alternative creation API/UI is implemented.

- **Original question:** Public registration or administrative provisioning?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 01/05/13: no public signup; staff/admin create active accounts. Audit 15 adds undeployed server authorization around provisioning.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Scope/account management P421 requires approved/manual verification; conclusion P912 says borrowers can register; Add User description P827 describes admin registration. This is an internal documentation conflict as well as a V1 conflict.
- **Current status:** **RESOLVED**, 2026-10-08.
- **Resolution / authority:** Provisioned Borrowers; no public signup. Only Admin creates any account; Staff assists operationally. Student-only bulk workflows; Faculty individual creation. Final owner update DEC-070 supersedes the Staff grant in DEC-053 in [DECISIONS](DECISIONS.md).
- **Remaining detail / scope:** DEC-070 confirms email/separate borrower-chosen password for both categories; Student uses approved institutional email, Faculty any valid unique accessible email. Single-use links to the respective email are recommended. Provider delivery (OPEN-017), SKSU Student domain/roster inputs (OPEN-028), activation lifecycle/expiry/reissue/existing-account treatment and general recovery remain dependent implementation details; no plaintext password import or public signup.
- **Gate:** NON-BLOCKING for core UX/domain planning; only stated later dependent work requires its review/authorization.

## ELAB-V2-OPEN-002

- **Original question:** Which borrower types are eligible: students only, faculty, staff?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 01/03/13: borrower UI is student-oriented; staff/admin share operational UI; no distinct faculty borrower role.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Stakeholder benefits P401–402 and Sprint 2 P613 mention students and faculty; objectives emphasize students/borrowers.
- **Current status:** **RESOLVED**, 2026-10-08.
- **Resolution / authority:** Active registered Borrower, at least Student/Faculty categories; category not role. Project owner Phase 2.5; DEC-051/052 in [DECISIONS](DECISIONS.md).
- **Remaining detail / scope:** No student-only or separate faculty authorization gate remains.
- **Gate:** NON-BLOCKING for core UX/domain planning; only stated later dependent work requires its review/authorization.

## ELAB-V2-OPEN-003

- **Original question:** How do staff and administrator privileges differ?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 03 and 15: shared operational UI and broad equivalent privileges, including creating/deleting privileged users and clearing fines; containment patch retains that equality.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Scope P421 names administrators; SOP P364 and browsing P425 name staff approval; distinct borrower/admin access described P605. No complete operation matrix is established.
- **Current status:** **RESOLVED**, 2026-10-08.
- **Resolution / authority:** Staff operational approval/direct/returns/replacements/history and user assistance, without account creation; Admin creates all accounts and controls deactivation/Student bulk/privileged accounts/fine clear/policy/audit. Final owner DEC-070 supersedes DEC-053 creation authority; DEC-052/060 in [DECISIONS](DECISIONS.md).
- **Remaining detail / scope:** Detailed ordinary inventory action definitions later; no ambiguity over Admin-only clear.
- **Gate:** NON-BLOCKING for core UX/domain planning; only stated later dependent work requires its review/authorization.

## ELAB-V2-OPEN-004

- **Original question:** What responsibilities would a Super Administrator have?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 03/13: only student/staff/admin; no Super Admin role or supervisor workflow.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Recommendations P922 propose administrator supervision, configurations, access/privilege management, monitoring and integrity across laboratories/departments.
- **Current status:** **REJECTED FOR CURRENT SCOPE**, 2026-10-08.
- **Resolution / authority:** Three named roles suffice; no Super Admin role. Project owner Phase 2.5; DEC-052 in [DECISIONS](DECISIONS.md).
- **Remaining detail / scope:** Any later governance expansion requires a new explicit proposal.
- **Gate:** NON-BLOCKING for core UX/domain planning; only stated later dependent work requires its review/authorization.

## ELAB-V2-OPEN-005

- **Original question:** Does a Super Admin belong in current FSMO V2 or only future expansion?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 01/10/13: globally shared FSMO inventory, no organizational model and no Super Admin.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** P922 is a recommendation tied to multiple laboratories/departments; it is not an accepted V2 requirement. Current stakeholder direction explicitly excludes campus-wide scope.
- **Current status:** **REJECTED FOR CURRENT SCOPE**, 2026-10-08.
- **Resolution / authority:** FSMO only, no campus hierarchy or global Super Admin. Project owner Phase 2.5; DEC-001/052 in [DECISIONS](DECISIONS.md).
- **Remaining detail / scope:** Separate future campus proposal remains outside V2.
- **Gate:** NON-BLOCKING for core UX/domain planning; only stated later dependent work requires its review/authorization.

## ELAB-V2-OPEN-006

- **Original question:** What fine amount, currency, grace period, cap and damage/loss rules apply?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 07/09: PHP 10/day overdue, unit-price damage/loss charges, client/timezone rounding; no evidenced cap or grace period.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Scope P430–431 and equipment definition P441 intend fines/collection/accountability; engineering descriptions do not establish a complete approved amount/rounding/waiver policy.
- **Current status:** **RESOLVED**, 2026-10-08.
- **Resolution / authority:** PHP 10/day overdue; ceil elapsed24h selected, replacement for damage/loss, no price charge. Project owner Phase 2.5; DEC-058/059/060 in [DECISIONS](DECISIONS.md).
- **Remaining detail / scope:** No rate/currency gate remains; do not invent grace/cap.
- **Gate:** NON-BLOCKING for core UX/domain planning; only stated later dependent work requires its review/authorization.

## ELAB-V2-OPEN-007

- **Original question:** How are fine assessment, payment, clearing and waivers distinguished?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 04/07/08/13: records/fines/finePaid representations diverge; UI clearing changes records without reconciling fines; no coherent payment ledger.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Dashboard/report scope P430–431 calls for payment records and fines collection; reporting intent alone does not define settlement or collection procedures.
- **Current status:** **RESOLVED**, 2026-10-08.
- **Resolution / authority:** Admin full-clear PAID/WAIVED/OTHER_RESOLUTION, no partial/payment allocation, history retained. Project owner Phase 2.5; DEC-052/060 in [DECISIONS](DECISIONS.md).
- **Remaining detail / scope:** Live full-clear checkpoints selected engineering recommendation; recorded paid distinguished from waiver/other.
- **Gate:** NON-BLOCKING for core UX/domain planning; only stated later dependent work requires its review/authorization.

## ELAB-V2-OPEN-008

**Historical Phase 1 engineering evidence; later product resolutions below govern:**

- **Phase 1B engineering evidence (2026-10-06):** Existing generic is_active=false denies subsequent protected API requests with 403; current role is looked up each request. No suspended-state taxonomy, session-family/logout revocation or loan/accountability policy was added. Those questions remain unresolved.

- **Original question:** What happens to inactive/suspended accounts and existing sessions/loans?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 03/05/09/13: login blocks non-active profile; current sessions are not centrally ejected; status types differ. Audit 15 denies inactive callers at its HTTP boundary.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Scope P421 requires approved users; Sprint 1 P605/719 describes deactivation, without session revocation/active-loan policy.
- **Current status:** **RESOLVED**, 2026-10-08.
- **Resolution / authority:** Admin deactivation blocks normal authentication/request actions; no suspension subsystem. Preserve old obligations, Staff/Admin resolves them. Project owner Phase 2.5; DEC-051/058 in [DECISIONS](DECISIONS.md).
- **Remaining detail / scope:** Institutional cross-device/session family details remain OPEN-027. DEC-073 resolves Student individual/bulk deactivation with obligations: Admin warns and confirms, but outstanding fines/active or overdue loans/unreturned equipment/replacements do not block. Preserve all obligations/history/due/overdue calculations and separate authoritative resolution; OPEN-029 is resolved, not a dependent execution gate.
- **Gate:** NON-BLOCKING for core UX/domain planning; only stated later dependent work requires its review/authorization.

## ELAB-V2-OPEN-009

**Historical Phase 1 engineering evidence; later product resolutions below govern:**

- **Phase 1B engineering evidence (2026-10-06):** Self-service email/identity changes are temporarily unavailable; only display-name editing remains. This containment avoids inventing verification/provisioning behavior and does not settle email ownership/eligibility policy.

- **Original question:** Is email verification required, and at what point?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 05: accounts created emailVerified=false; no send/check verification flow. Administrative verification is not verified email ownership.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Scope P424 describes manual credentials verification; P427–428 emphasize accurate reachable email. No explicit email-ownership verification gate is established.
- **Current status:** **PARTIALLY RESOLVED**, 2026-10-08; residual implementation review remains.
- **Resolution / authority:** DEC-070 confirms both categories use email and separate borrower-chosen passwords during secure activation. Students require approved SKSU email; Faculty may use any valid unique accessible email and require no Student ID. Single-use activation links to the respective mailbox remain recommended, not configured delivery. Syntax/domain/roster presence alone is not mailbox ownership. Product authority and category requirements are resolved; technical ownership evidence remains to be implemented.
- **Remaining detail / scope:** FSMO/IT must finalize activation lifecycle, link origin, token expiry/reissue limits, failed/lost delivery, mailbox verification evidence, existing-account treatment, email changes and password change/recovery/revocation rules. Provider is OPEN-017; exact approved SKSU Student configuration/roster formatting and technical matching are OPEN-028. No actual activation/recovery code or delivery is claimed.
- **Gate:** Planning and independently verifiable infrastructure may proceed within authorized scope; dependent activation/recovery and real delivery require the stated reviews. No Phase5 implementation authorization from this update.

## ELAB-V2-OPEN-010

- **Original question:** How does kiosk authentication, idle timeout and session cleanup work?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 01/03/13: shared responsive web/student flows; no separately hardened kiosk route/session lifecycle evidenced.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Scope P429 and definition P440 describe self-service browsing/requests on physical kiosk and personal devices; authentication, idle expiry and handoff cleanup are unspecified.
- **Current status:** **RESOLVED**, 2026-10-08.
- **Resolution / authority:** Interactive Kiosk means normal responsive catalog/cart/request; hardware/device/handoff authentication removed. Project owner Phase 2.5; DEC-061 in [DECISIONS](DECISIONS.md).
- **Remaining detail / scope:** No kiosk protocol/hardware approval is required for this flow.
- **Gate:** NON-BLOCKING for core UX/domain planning; only stated later dependent work requires its review/authorization.

## ELAB-V2-OPEN-011

- **Original question:** Are equipment categories required, and who maintains them?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 06/13: name/search/status discovery; no category field or category management.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Browsing P425 and Sprint 2 P613 describe category/availability filtering.
- **Current status:** **NON-BLOCKING**, 2026-10-08.
- **Resolution / authority:** Interactive browse/search/filter capability confirmed; category/type optional catalog structure. Project owner Phase 2.5; DEC-061 in [DECISIONS](DECISIONS.md).
- **Remaining detail / scope:** FSMO catalog owner confirms taxonomy/content and metadata details before final catalog seeding; placeholder labels sufficient for mockups.
- **Gate:** NON-BLOCKING for core UX/domain planning; only stated later dependent work requires its review/authorization.

## ELAB-V2-OPEN-012

- **Original question:** What are equipment deletion, archival and retention rules?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 06/07/09/13: hard deletion with client-side active-transaction guard and concurrency/status gaps; archive not found; history depends on snapshots.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Scope P423 requires inventory updates and monitoring; Add Equipment description demonstrates entry, but no archive/retention policy is established.
- **Current status:** **NON-BLOCKING**, 2026-10-08.
- **Resolution / authority:** History retained; archive guard physical/holds/replacement zero selected; ordinary metadata/stock operations versus exceptional Admin correction defined. Project owner Phase 2.5; DEC-049/050/062 in [DECISIONS](DECISIONS.md).
- **Remaining detail / scope:** FSMO inventory owner defines repair/discard/retire procedures; legal destructive retention under OPEN-026. Not a core loan-flow gate.
- **Gate:** NON-BLOCKING for core UX/domain planning; only stated later dependent work requires its review/authorization.

## ELAB-V2-OPEN-013

- **Original question:** May a borrower or staff cancel a request/active transaction?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 07/13: borrower cancellation not found; staff can hard-delete transactions with stock restoration edge cases and no history.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** SOP P364 describes submission/approval/release/archive, not a defined cancellation lifecycle. Form Cancel buttons are not cancellation-policy evidence.
- **Current status:** **RESOLVED**, 2026-10-08.
- **Resolution / authority:** Owner cancels own PENDING only; release/history. No issued cancellation/delete. Project owner Phase 2.5; DEC-055 in [DECISIONS](DECISIONS.md).
- **Remaining detail / scope:** Optional reasoned Staff/Admin pending cancellation is an engineering recommendation, select only if needed in a future feature plan.
- **Gate:** NON-BLOCKING for core UX/domain planning; only stated later dependent work requires its review/authorization.

## ELAB-V2-OPEN-014

- **Original question:** Must denial include a reason and retained history?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 07/08/09: denial restores stock, queues notification and deletes request; no durable denial record/reason.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Denial email figure/description P851 communicates denial/contact guidance, but does not specify retained reasons or structured decision history.
- **Current status:** **RESOLVED**, 2026-10-08.
- **Resolution / authority:** Staff/Admin denial requires borrower-visible reason; retained DENIED releases holds. Project owner Phase 2.5; DEC-055 in [DECISIONS](DECISIONS.md).
- **Remaining detail / scope:** V1 deletion behavior superseded; no denial-content gate remains.
- **Gate:** NON-BLOCKING for core UX/domain planning; only stated later dependent work requires its review/authorization.

## ELAB-V2-OPEN-015

- **Original question:** What is the canonical due/overdue vocabulary and day-boundary policy?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 04/07/09: the stored audit vocabulary is Request/Ongoing/Ondue/Overdue/Incomplete and combined variants (audit 04/07), rather than the earlier Due/Missing shorthand; Asia/Manila scheduler/calendar rules and mixed timestamp arithmetic.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Scope P427/430 and reminder figures use due/overdue; they do not establish exact canonical states or transition times.
- **Current status:** **RESOLVED**, 2026-10-08.
- **Resolution / authority:** Required due date+time, Asia/Manila, unresolved physical/replacement overdue; adopted ceil elapsed24h. Project owner Phase 2.5; DEC-056/058/059 in [DECISIONS](DECISIONS.md).
- **Remaining detail / scope:** No evidenced conflict requiring rounding confirmation; completed amount frozen, exact due zero.
- **Gate:** NON-BLOCKING for core UX/domain planning; only stated later dependent work requires its review/authorization.

## ELAB-V2-OPEN-016

- **Original question:** How are terms versions, acceptance evidence and reacceptance handled?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 04/09/13: static terms UI; acceptance boolean and timestamp only, no version/content hash.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** SOP and borrowing intent require accountability, but engineering passages reviewed do not define versioned acceptance or policy-change handling.
- **Current status:** **RESOLVED**, 2026-10-08.
- **Resolution / authority:** Once/current material version user/version/time, first activation/use; changed version gates next request, not per-loan. Project owner Phase 2.5; DEC-057 in [DECISIONS](DECISIONS.md).
- **Remaining detail / scope:** DEC-072 confirms official FSMO wording will be finalized after presenting the application to FSMO. Preserve implemented Phase4B infrastructure and isolated TEST fixtures; never invent/publish/automatically accept institutional text. Actual approved document/version/publication responsibility remains external; live borrowing must enforce official publication and documented acceptance. Pending wording does not block independent account-management or inventory development when separately authorized.
- **Gate:** NON-BLOCKING for independent account-management/inventory development under separate authorization; required before official publication/consent and live borrowing.

## ELAB-V2-OPEN-017

- **Original question:** Which email delivery provider/consumer and retry/receipt mechanism are approved?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 08/12/13: two notification shapes; queue producers present, consumer/extension absent; deployed delivery cannot be inferred.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Scope P427–428 and Sprint 4 P628–631 intend confirmations, reminders, delivery logs and failure handling; screenshots imply received mail but do not identify the operational consumer/provider.
- **Current status:** **PROVIDER SELECTION RESOLVED; deployment/technical delivery dependencies remain**, 2026-10-08.
- **Resolution / authority:** Brevo selected for account activation and future password recovery by the owner in DEC-071; vendor choice is no longer open. Transactional outbox concept/delivery observability retains its separately planned scope under DEC-049 in [DECISIONS](DECISIONS.md).
- **Remaining detail / scope:** Secure backend-only Brevo integration during authorized account-management development, configured backend API key, verified sender, trusted link origin and successful actual delivery/failure testing before live activation/recovery claims. No credentials in browser/commits/logs. Activation/recovery lifecycle/ownership/operations remain OPEN-001/009; later borrowing notification cadence remains OPEN-025, not vendor-selection work. No provider adapter is implemented by this policy update.
- **Gate:** NON-BLOCKING for independent authorized development/test-adapter verification; required before real activation/recovery delivery and dependent notification sending. A selected provider, mock or send intent is not verified delivery.

## ELAB-V2-OPEN-018

- **Original question:** What V1 data should migrate, reconcile or be excluded?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 04/07/14: users/Auth IDs, aggregate equipment, active transactions, snapshots, conflicting fines, notifications and images; live volumes/rules/state not inspected.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Product history/reporting intent P422/431 does not define cutover, historical retention or migration scope.
- **Current status:** **NON-BLOCKING**, 2026-10-08.
- **Resolution / authority:** Legacy migration is deferred and requires authorized access/export/cutover; no live Firebase. Project owner Phase 2.5; DEC-049/050 in [DECISIONS](DECISIONS.md).
- **Remaining detail / scope:** V1 custodian/FSMO approves dataset, source lineage, opening stock, historically completed damage/loss and fine reconciliation, backup/rollback. Gates actual Phase 12 migration, not core UX.
- **Gate:** NON-BLOCKING for core UX/domain planning; only stated later dependent work requires its review/authorization.

## ELAB-V2-OPEN-019

- **Original question:** How should contradictory legacy technology/schema descriptions be treated?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 01/04/12: Expo/React Native + Firebase, Firestore collections and Functions; source is not supplied in this V2 workspace.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** Sprint 1 P606 describes PHP/MySQL; database tables use relational types. These contradict actual audited V1 architecture.
- **Current status:** **NON-BLOCKING**, 2026-10-08.
- **Resolution / authority:** Current V2 Go/Fiber/React/PostgreSQL stack accepted; actual V1 audit uses Expo/Firebase despite capstone PHP/MySQL prose. Project owner Phase 2.5; DEC-002/003/004 in [DECISIONS](DECISIONS.md).
- **Remaining detail / scope:** Engineering/V1 custodian verifies source revision during authorized migration; not a V2 stack choice.
- **Gate:** NON-BLOCKING for core UX/domain planning; only stated later dependent work requires its review/authorization.

## ELAB-V2-OPEN-020

- **Original question:** Which borrowing limits, reservation point and direct-checkout rules are approved?
- **Evidence from V1 (FACTUAL V1 BEHAVIOR):** Audit 07/09: requests reserve stock immediately, student due date within seven days; staff direct checkout differs; no evidenced concurrent-loan limits or fine-based block.
- **Evidence from latest documentation (INTENDED LEGACY REQUIREMENT):** SOP P364 says release after approval; scope P425 requires staff approval; reservation timing, limits and direct-checkout exceptions are not fully specified.
- **Current status:** **RESOLVED**, 2026-10-08.
- **Resolution / authority:** Reserve-on-submit, active Borrower, atomic issue/direct checkout, no seven-day maximum/automatic fine block. Project owner Phase 2.5; DEC-051/054/056 in [DECISIONS](DECISIONS.md).
- **Remaining detail / scope:** No new concurrency quotas or duration maximum guessed.
- **Gate:** NON-BLOCKING for core UX/domain planning; only stated later dependent work requires its review/authorization.

## ELAB-V2-OPEN-021

- **Original question:** Does staff approval record a permission decision, actual physical release, or two separate actions?
- **ACTUAL V1 BEHAVIOR:** Audit 07: Request→Ongoing, borrowedDate reset, no separately recorded pickup/release.
- **CAPSTONE INTENT:** SOP P364 and scope P425 release after staff approval; whether there is a waiting interval is not specified.
- **Current status:** **RESOLVED**, 2026-10-08.
- **Resolution / authority:** Approval and physical handover are one PENDING→CHECKED_OUT edge; direct immediate issue. Project owner Phase 2.5; DEC-054 in [DECISIONS](DECISIONS.md).
- **Remaining detail / scope:** Separate approved-waiting proposal superseded.
- **Gate:** NON-BLOCKING for core UX/domain planning; only stated later dependent work requires its review/authorization.

## ELAB-V2-OPEN-022

- **Original question:** Which return/disposition evidence, repair/recovery/retirement/write-off and correction procedures are authorized?
- **ACTUAL V1 BEHAVIOR:** Audit 06/07/09: partial good/damaged/lost quantities, notes required in UI; damage/loss not explicit stock buckets; no evidenced repair/write-off or return-correction workflow.
- **CAPSTONE INTENT:** Scope P424/426 and definitions P441/443 require manual condition inspection/accountability, without detailed disposition procedures.
- **Current status:** **NON-BLOCKING**, 2026-10-08.
- **Resolution / authority:** Damage/loss replacement, physical custody separation, Staff/Admin acceptance and completion gating are RESOLVED. D17 clarification also resolves the return-photo boundary: a borrower may physically show a personally stored phone photo in person, entirely outside eLabTrack; Staff/Admin alone records condition/quantities. No borrower return/evidence submission or return-photo/attachment software feature, storage or upload step is required or deferred. Project owner Phase 2.5; DEC-058/062 in [DECISIONS](DECISIONS.md).
- **Remaining detail / scope:** FSMO inventory owner later specifies original repair/discard/retain/retire and detailed equivalence guidance; retain originals nonusable and staff-certify equivalence meanwhile. Gates a specific disposal/repair feature, not borrowing/replacement flow.
- **Gate:** NON-BLOCKING for core UX/domain planning; only stated later dependent work requires its review/authorization.

## ELAB-V2-OPEN-023

- **Original question:** When do equipment value, borrower presentation, due policy and governing terms become binding for a borrowing?
- **ACTUAL V1 BEHAVIOR:** Audit 04/07: price/name snapshots captured at creation; approval resets borrowedDate; totalPrice semantics uncertain.
- **CAPSTONE INTENT:** P364/421/425/441 accountability/release intent, no price-lock or policy-change rule.
- **Current status:** **RESOLVED**, 2026-10-08.
- **Resolution / authority:** Due/policy/terms and issue snapshots explicit; no equipment-price monetary assessment. Project owner Phase 2.5; DEC-056/057/058/060 in [DECISIONS](DECISIONS.md).
- **Remaining detail / scope:** Catalog price metadata informational if retained; pending terms binding survives new material version. No damage-value gate remains.
- **Gate:** NON-BLOCKING for core UX/domain planning; only stated later dependent work requires its review/authorization.

## ELAB-V2-OPEN-024

- **Original question:** At what event/window are overdue, damage and loss charges posted, and how are partial returns, zero balance and reversals handled?
- **ACTUAL V1 BEHAVIOR:** Audit 07/08/09: mutable active fine and cumulative final assessment; divergent clearing/payment representations.
- **CAPSTONE INTENT:** P430/431 fines/collection reports, no detailed posting/day/partial/correction policy.
- **Current status:** **RESOLVED**, 2026-10-08.
- **Resolution / authority:** Live deterministic overdue projection, immutable completion freeze and full-clear history; no generalized charge posting. Project owner Phase 2.5; DEC-059/060 in [DECISIONS](DECISIONS.md).
- **Remaining detail / scope:** No daily cumulative posting, price charges, partial payments, adjustment/reversal architecture.
- **Gate:** NON-BLOCKING for core UX/domain planning; only stated later dependent work requires its review/authorization.

## ELAB-V2-OPEN-025

- **Original question:** Which notifications, recipients and reminder cadence are institutionally required, and does grace change reminder wording/timing?
- **ACTUAL V1 BEHAVIOR:** Audit 08: approval/denial/receipt and once-per-transaction due/tomorrow/overdue producer flags; delivery consumer absent.
- **CAPSTONE INTENT:** P427/428 and Sprint 4–5 P628–636 require confirmations/reminders, logged delivery/failure intent, not a provider/cadence guarantee.
- **Current status:** **NON-BLOCKING**, 2026-10-08.
- **Resolution / authority:** Events confirmed: submit, deny, approve/checkout, expire, due reminder, overdue, partial return, replacement created/resolved, complete. Project owner Phase 2.5; DEC-054/055/058/059 in [DECISIONS](DECISIONS.md).
- **Remaining detail / scope:** FSMO operations selects exact reminder lead times/repeat cadence/templates with IT delivery constraints before scheduled notification feature; expiry is fixed24h independently.
- **Gate:** NON-BLOCKING for core UX/domain planning; only stated later dependent work requires its review/authorization.

## ELAB-V2-OPEN-026

- **Original question:** What privacy categories, retention, anonymization, evidence access and legal/institutional data procedures apply?
- **ACTUAL V1 BEHAVIOR:** Audit 04/05/08/14: PII snapshots, public profile images, retained history and unbounded exposure; no retention schedule evidenced.
- **CAPSTONE INTENT:** P422/431 history/reporting intent, no legal retention duration or anonymization procedure established.
- **Current status:** **NON-BLOCKING**, 2026-10-08.
- **Resolution / authority:** Canonical consequential history/evidence preserved; retention periods not invented. Project owner Phase 2.5; DEC-049 in [DECISIONS](DECISIONS.md).
- **Remaining detail / scope:** Institutional/legal data custodian confirms retention, anonymization, export and audit access obligations before destructive cleanup/production retention rollout; minimal PII throughout.
- **Gate:** NON-BLOCKING for core UX/domain planning; only stated later dependent work requires its review/authorization.

## ELAB-V2-OPEN-027

- **Original question:** What institutional concurrent-session, global/family logout and immediate access-revocation behavior is required?
- **ACTUAL V1 BEHAVIOR:** Audit 05/15: client login/status and undeployed current-caller checks; no family/global policy.
- **CAPSTONE INTENT:** Scope P421 and deactivation Sprint 1 P605/607 do not establish global session/device rules.
- **Current status:** **NON-BLOCKING**, 2026-10-08.
- **Resolution / authority:** Phase 1I browser coordination complete and unchanged; inactive-account current checks remain. Project owner Phase 2.5; DEC-045/046/047/051 in [DECISIONS](DECISIONS.md).
- **Remaining detail / scope:** FSMO/IT may later decide concurrent sessions, family/global logout or stronger immediate cross-device access revocation. No Phase 1 reopening and no new borrowing suspension policy.
- **Gate:** NON-BLOCKING for core UX/domain planning; only stated later dependent work requires its review/authorization.

## ELAB-V2-OPEN-028 — SKSU Student domain configuration and roster technical contract

- **Status:** Remaining institutional inputs and technical dependencies; core Student/Faculty requirements RESOLVED by DEC-070.
- **APPROVED PRODUCT RULES:** Admin-only creation; Students require unique textual Student ID matching official identification and approved SKSU institutional email. Faculty is individual-only, uses any valid unique accessible email and does not require Student ID. Standard bulk workflows are Student-only; no Faculty allowlist, optional Student-ID or Faculty-bulk eligibility decision remains.
- **REMAINING INSTITUTIONAL APPROVALS/INPUTS:** FSMO/registrar/IT supplies exact approved SKSU Student domains and official ID/roster format examples without copying personal records. Confirm actual course/contact availability and requiredness for the recommended studentId/name/email/course/contactNumber template; do not reopen required Student ID/name/email identity inputs.
- **REMAINING TECHNICAL DEPENDENCIES:** Reviewed ID normalization and spreadsheet string/leading-zero preservation, domain matching, changed-email reconciliation and stable-ID matching to internal UUID. Reject missing required Student ID and ambiguous/conflicting identities; no silent missing-ID fallback, numeric coercion or Faculty relabeling. Field/category persistence and enforcement remain unimplemented.
- **Gate:** Finalize actual configuration/format/matching contract before dependent implementation, within separately authorized Phase5. Do not invent a domain or ID regex. String representation and uniqueness are approved requirements, not proposals.

## ELAB-V2-OPEN-029 — Student deactivation with outstanding obligations (RESOLVED)

- **Status:** **RESOLVED**, 2026-10-08; removed from active open/deployment gates. Retained ID records supersession, not an outstanding question.
- **Authority:** Owner additional policy decisions; DEC-073. Previous obligation-based exclusion/override alternatives are superseded.
- **Resolution:** Admin may deactivate a Student individually or through the Student-only bulk workflow regardless of fines, active/overdue borrowing, unreturned equipment or replacement obligations, for graduation/withdrawal/transfer/other authorized deactivation. Display appropriate warnings and require confirmation; obligations must not block deactivation.
- **Preserved boundaries:** Current Admin authority and Student-only bulk identity/role/category checks, preview/confirmation/audit/idempotency/history remain required. No automatic fine clearance/waiver/payment, return, replacement resolution, borrowing closure, overdue change or history deletion. FSMO resolves separately; fine clearance remains Admin-only and auditable.
- **Remaining technical work:** Implement honest obligation projections/warnings, confirmation and status/audit commands under separately authorized Phase5; missing projections cannot claim zero. This is implementation work, not an unresolved policy or terms-content dependency.

## Proposal handling

Record any later changed stakeholder instruction with authority/date and affected decision IDs; preserve factual V1 history. Engineering design choices (hybrid inventory, live-clear checkpoints, lock order, endpoints) are explicitly recommendations, not invented institutional rules. No separate stakeholder packet is needed: the owner provided working answers directly. The current policy update changes documentation only. Existing Phase1/4A/4B foundations retain their recorded implementation status; Phase5/6/later work is not automatically authorized.

## Batch1 inventory dependency reconciliation

Phase6 has an authorized, locally supportable bounded database raster strategy (DEC-075), not a blocked paid-provider choice. Future object-storage scaling/retention and approved catalog taxonomy/content remain distinct technical/institutional dependencies, not reasons to disable independent operator-maintained catalog CRUD. No seeded institutional categories or stock reconciliation guesses. Future borrowing/replacement liability integration must precede archive when those tables exist; the current boundary fails closed. Detailed damaged-original repair/disposal/correction procedures remain OPEN-021 scope and are not implemented as ordinary inventory. Phase7 still requires separate authorization plus official terms/acceptance before live borrowing. Brevo live verification and actual Student-domain inputs remain unchanged external account gates.
