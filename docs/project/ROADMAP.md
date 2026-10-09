# Initial delivery roadmap

Confirmed sequence from current stakeholder direction. **Phases 0, 1A–1I, 2, 2.5, 3A.1, 3A.2 and 3B are complete within their recorded scopes; Phase 3B is OWNER VISUALLY APPROVED.** The owner reviewed real implementation screenshots and accepted the disclosed minor adaptations. The navy/indigo/violet PNG package is binding for future screens. Phase 3B delivers shared presentation and four development previews. **Phase 4A is complete. Phase 4B terms foundation is implemented and verified; official terms after FSMO presentation and verified live activation delivery remain pending; secure activation and the backend Brevo adapter are implemented in Phase5; Admin-only Student/Faculty provisioning and email/separate-password rules are confirmed in DEC-070. Batch1 authorizes Phases5–6. Phase5 core gates pass with external activation/domain/official-term dependencies; see PHASE5_REPORT.md. Phase6 is complete within its reviewed scope. Frontend Reconstruction Stages A and B are engineering-complete and owner-approved; the corrected Inventory filter is engineering-verified and owner manually verified as passed at the final checkpoint. Phase7 remains NOT STARTED and unauthorized.** See [Phase 3B report](PHASE3B_REPORT.md), [DEC-064](DECISIONS.md#elab-v2-dec-064--phase-3b-implementation-visually-approved), [local readiness report](LOCAL_ENVIRONMENT_READINESS_REPORT.md) and [binding visual map](../ux/APPROVED_VISUAL_IMPLEMENTATION_MAP.md).

## Phase 0 — Rebaseline & Template Adaptation

- **Goal:** Establish the current FSMO project contract and clean starter identity.
- **Primary deliverables:** Charter, source policy, decisions/open questions, roadmap, foundation audit, security backlog, adapted agent rules, minimal placeholder, validation report.
- **Dependencies:** Stakeholder brief, capstone, V1 audit and source inventory.
- **Explicit non-goals:** No product features, domain tables, V1 Firebase access or Phase 1 hardening implementation.
- **Exit gate:** Source-grounded documents and cleanup reviewed; reference hygiene proven; quality gates pass or blockers explicitly prevent unconditional closure.

## Phase 1 — Engineering Foundation

- **Goal:** Make the retained bootstrap reproducible, testable and secure enough for later product work.
- **Primary deliverables:** Verified module identity/toolchains, environment validation, API contracts, session infrastructure hardening, migration safety, CI/quality baseline.
- **Dependencies:** Phase 0 report; resolve environment/identity blockers and prioritize security backlog.
- **Explicit non-goals:** No equipment/borrowing/report/kiosk business behavior; no premature infrastructure.
- **Exit gate:** Foundation quality gates and critical security acceptance tests pass; production defaults fail safely; remaining risks recorded. Product auth policy remains Phase 4 work.

## Phase 1A–1E — Completed source foundation

Production configuration/runtime safety; authentication/current-account authorization; JWT/hash-only transactional refresh; browser session architecture; HTTP/API security. Completion is scoped to each report, with all live integration requirements retained. No product policies were inferred.

## Phase 1F — API Contracts, Logging & Observability

- **Goal/deliverables:** Consistent API success/error/validation contracts, centralized safe mapping, request ID propagation, structured redacted operational/security/lifecycle logging, typed client failures and explicit health/readiness.
- **Dependencies:** Completed 1A–1E; current retained routes/session/perimeter preserved.
- **Non-goals:** No real PostgreSQL/migration/browser/Docker runtime, telemetry vendor, business audit schema, product UI or mockups.
- **Exit gate:** Foundation contracts/logging/security tests and required local quality/Compose checks pass; unverified runtime evidence carried to 1G. Current report records actual gate state.

## Phase 1G — Real Integration Verification

- **Goal/deliverables:** Authorized isolated PostgreSQL startup/migrations/rollback/TLS/credential restrictions, hash persistence/transaction/concurrent refresh, cleanup; Docker/nginx/browser cookie/bootstrap/logout/CORS/CSP/proxy/rate and health/readiness verification; real request-ID/log/redaction and cross-tab/lost-response checks.
- **Dependencies:** Phase 1F and an explicitly authorized isolated environment; no V1 access implied.
- **Non-goals:** No business features or automatic migration-runner redesign; record runner limitations for 1H.
- **Exit gate:** Actual runtime/integration evidence with recovery and concurrent cases; offline doubles alone do not establish completion.

## Phase 1H — Migration Runner & Database Privilege Hardening

- **Goal/deliverables:** Atomic DDL/bookkeeping, immutable paired checksums, native advisory exclusion, strict history/discovery/error handling, explicit CLI, migration/runtime credential and ownership separation, real failure/concurrency/privilege/API/auth tests.
- **Evidence:** Stakeholder's boss-rebuild audit comparison and current runner both identify DDL commit before separate bookkeeping, no checksums/lock and lookup failure ambiguity.
- **Dependencies:** Phase 1G evidence and separate authorization; must precede heavy reliance on business-domain migrations.
- **Non-goals:** No product tables, automatic startup migrations or infrastructure platform redesign.
- **Exit gate:** Satisfied locally: transaction/version integrity, corruption/concurrent-run denial, correct error propagation and separated-role API/auth operation demonstrated against isolated PostgreSQL; standard/race/Compose gates pass. Production role/TLS deployment remains unverified.

## Phase 1I — Cross-Tab Session Coordination

- **Goal/deliverables:** Coordinate same-origin cookie-changing operations and peer session/logout lifecycle before product UI, preserving memory-only access/HttpOnly refresh and strict single-use server rotation.
- **Evidence:** Real Phase 1G/1H browser race: one refresh succeeds; a late losing 401 clears the shared winner cookie. Peer logout leaves another tab apparently authenticated until refresh fails.
- **Dependencies:** Completed 1H and explicit Phase 1I authorization, received 2026-10-07; required before Phase 2/product UI.
- **Non-goals:** No usable-token broadcast/persistence, server replay relaxation, business schema/features or CI platform implementation inferred from this requirement.
- **Exit gate:** **Satisfied / COMPLETE, 2026-10-08.** Original defect reproduced before fixes; Web Locks serialize cookie mutations and BroadcastChannel carries strict non-secret lifecycle hints. Refresh failures never delete possibly newer cookies. Real two/three-tab pressure, simultaneous reload, logout/private-cache removal, delayed-success fencing, disable/revoke, owner close, stale events and unavailable-primitive fallback pass. PostgreSQL retains one winner/11 denials; standard/race/Compose/log/diff gates pass. See the [58-item report and whole-foundation assessment](PHASE1_FOUNDATION.md#phase-1i--cross-tab-session-coordination). CI remains open; no Phase 2 or 1J is authorized by closure.

## Phase 2 — Domain & Database Design

- **Historical result:** Domain/design delivered with unresolved product-policy blockers; [Phase 2 report](PHASE2_REPORT.md) preserves that snapshot.
- **Goal/deliverables:** FSMO domain contracts, relational candidates, stock/history/transaction invariants, evidence/source reconciliation and future API draft.
- **Current gate:** Superseded by Phase 2.5 owner decisions and current domain baseline; no business tables or code implemented.
- **Non-goals:** No feature implementation, migrations, mockups, boss writes or automatic next phase.

## Phase 2.5 — Product Decision Integration & Domain Rebaseline

- **Result:** COMPLETE — documentation reconciliation only, 2026-10-08; [report](PHASE2_5_REPORT.md).
- **Deliverables:** All 29 working decisions integrated, four-count physical stock plus replacement obligations, six-state lifecycle, deterministic PHP 10/day fine, Admin full clear, revised role matrix,71 active invariants and18 design walkthroughs.
- **Dependency/gate:** Core policy blockers resolved; specific NON-BLOCKING onboarding/disposal/retention/delivery/migration/operations details retained.
- **Non-goals:** No code, SQL, UI, mockups, Phase 1 changes, boss changes, commit/push/deployment; Phase 3A not begun.

## Phase 3A — UX Architecture & Mockups

- **Status:** COMPLETE — low-fidelity architecture and supplied high-fidelity visual targets owner-approved,2026-10-08.
- **Sequence:** Phase 3A.1 low-fidelity architecture → separately authorized Phase 3A.2 high-fidelity mockups → owner mockup approval → separately authorized Phase 3B implementation.
- **Locked audience model:** Borrower mobile-first; Staff/Admin desktop/tablet-first responsive; Interactive Kiosk is the normal catalog/cart/request flow, no dedicated hardware shell.
- **Shared dependency:** Current Phase 2.5 rules, working role matrix, physical/replacement reconciliation and outside-system return-photo clarification.

## Phase 3A.1 — UX Architecture & Low-Fidelity Wireframes

- **Result:** COMPLETE — architecture approved by the current Phase 3A.2 owner brief, 2026-10-08. The original [report](PHASE3A1_REPORT.md) remains the historical delivery snapshot; its awaiting-review wording is superseded by this authorization.
- **Deliverables:** Seven UX documents;59 inventoried surfaces (2 shared,22 Borrower,23 operational,12 Admin-only),38 text wireframe groups; two navigation maps; journeys A–H; concrete responsive/status/loading/empty/error/conflict/confirmation rules, future shadcn mapping and unchecked implementation acceptance checklist.
- **Readiness:** Phase3A.2 subsequently approved; Phase3B subsequently implemented and owner visually approved. Minor UX/content/import/projection/activation/export questions remain explicit visual review or later implementation dependencies; no new domain-policy blocker found.
- **Non-goals:** No high-fidelity styling/mockups, React/components/backend/SQL, Phase 1 session/security or Phase 2.5 policy change; no boss writes/commit/push/deployment.

## Phase 3A.2 — High-Fidelity Mockups

- **Status:** COMPLETE / OWNER VISUAL APPROVAL RECEIVED,2026-10-08. The original [report](PHASE3A2_REPORT.md) remains the historical evergreen proposal; its pending-review/color wording is superseded by the supplied [approved PNG package](../ux/approved/README.md).
- **Current deliverables:**36 verified manifest files (35 screen images + separate seal), representing4 baseline screens +32 remaining visual targets; navy/indigo/violet identity, compact borrower header and operational shell. Old99 SVG artboards remain structural history only.
- **Authority:** Locked Phase2.5 business rules win over inaccurate illustration copy. Approval does not invent a new lifecycle, stock model, permissions, return-photo feature or unresolved account/import/export policy.
- **Exit:** Owner visual approval and explicit Phase3B authorization received. No later business implementation implied.

## Phase 3B — Approved Visual Design System & Application Shell

- **Status:** COMPLETE within the authorized Phase 3B scope, 2026-10-08; final verification recorded in [Phase3B report](PHASE3B_REPORT.md). **OWNER VISUALLY APPROVED**: real implementation screenshots, both themes, responsive layouts and disclosed minor adaptations accepted in DEC-064. No later phase is authorized.
- **Deliverables:** Navy/violet semantic light/dark tokens; preserved shadcn foundation; shared presentation and Borrower/Staff/Admin shells; four dev-only real component previews; local fixture controls; Chromium viewport/theme and keyboard evidence;36-target map and fidelity contract.
- **Dependencies:** Approved package verified before code; Phase2.5 meanings/permissions retained; no Phase1 foundation replacement.
- **Explicit non-goals:** No production business dashboard/catalog/borrowing/return/fine/account/report screens, backend services or migrations; no fake role session, boss port, commit/push/deploy or Phase4.
- **Exit gate:** Actual screenshots compared and material defects corrected/disclosed; responsive/accessibility/component/full frontend/whitespace gates pass. The owner accepted the disclosed minor adaptations; the 32 future screens remain unimplemented and retain individual PNG/runtime gates.
- **Handoff:** Every later feature must open its matching approved PNG and record fresh runtime fidelity/contract/integrity evidence. Current baseline approval is recorded. Future material visual changes need owner disposition; explicit next-phase authorization is still required.

## Phase 4 — Authentication & Authorization

- **Status:** Phase 4A COMPLETE; Phase 4B IMPLEMENTED FOUNDATION with blocked official-content/activation functionality. This historical Phase4 gate does not itself authorize Phase5; the later explicit Batch1 request does.
- **Foundation:** Preserve completed Phase 1 security/session architecture and approved visual shells; confirmed product roles and provisioned-account policy apply.

### Phase 4A — Authentication, Roles & Protected Navigation

- **Result:** COMPLETE, 2026-10-08; [32-section report](PHASE4A_REPORT.md), DEC-065.
- **Deliverables:** Real responsive login, current database Borrower/Staff/Admin permission checks, paired reversible role mapping with safe Staff rollback refusal, protected exact routes/deep-link restoration, loading/error/recovery/logout/disabled handling, truthful feature placeholders, safe own-account display and existing development previews.
- **Verification:** Actual PostgreSQL migration and HTTP role/session tests, strict rotation/concurrent replay/rollback and runtime privilege regressions, real Chromium forms/cross-tab/role/status/theme/responsive checks, preserved visual preview checks and full quality gates. No public registration or published test passwords.
- **Non-goals:** No Phase 4B terms/activation workflow, Phase 5 directory/provisioning/import/status UI or business feature services.

### Phase 4B — First-Use Terms & Onboarding

- **Status:** Explicitly authorized and foundation implemented, 2026-10-08; [34-section report](PHASE4B_REPORT.md), DEC-066/067. Official publication/live borrowing remains pending owner-approved terms after FSMO presentation; Batch1 activation/Brevo implementation is verified with test adapters; live delivery remains pending, without blocking independent account-management/inventory development (DEC-071/072); DEC-068/070 confirm email/separate-password login for both Student/Faculty categories and recommend links to their respective mailboxes. Official content remains unapproved.
- **Delivered:** Migration 000005, immutable immediate mandatory Admin publication, user/version/server-time receipts, own status API, real first-use/updated/missing/error UI, reusable transactional future-command policy, existing session/cache integration and real PG/HTTP/Chromium verification. No official content or acceptance is seeded in the configured application database.
- **Integration points:** Immutable terms versions/content/hash, unique user/version/time acceptance; current material version before new request/direct issue, preserve access to existing obligations. Use existing account/session ports and future transactional request guards; authentication does not imply acceptance.
- **Dependencies:** Owner explicitly confirms no official V2 FSMO text is approved. Finalize approved exact document/publication responsibility after application presentation (OPEN-016); this is not an independent account-management/inventory development gate. Remaining technical dependencies include activation lifecycle/recovery details (OPEN-001/009), actual approved SKSU Student configuration/roster technical contract (OPEN-028), and selected Brevo backend adapter/API key/verified sender/link-origin plus tested real delivery (OPEN-017). Email/password and vendor selection are resolved. Batch1 adds the tested adapter and activation infrastructure; live delivery remains unverified. No fake official terms, reset or public signup is implemented; Batch1 adds independently verified provisioning/activation without proving live email delivery.
- **Exit gate:** Technical version/change/concurrency/ownership and accessible first-use verification passes. Official publication/institutional acceptance and secure activation are not complete. Read-only account/existing obligations and logout remain accessible without acceptance; new Borrower home/catalog paths fail closed until applicable terms are accepted. Phase 7 must call the policy inside its borrowing transaction.

## Phase 5 — Users & Borrower Management

- **Status:** PHASE5 PARTIALLY COMPLETE — EXTERNAL GATES. Batch1 authorized implementation; account/backend/PG/security/Chromium core gates PASS. See [report](PHASE5_REPORT.md). Brevo live delivery, approved domains and official FSMO terms remain external; independent Phase6 is authorized after this checkpoint.
- **Goal:** Admin provisions/manages Student/Faculty Borrowers and named Staff/Admin with authoritative backend permissions, category-specific identity validation and separate borrower-owned passwords; preserve Phase4A/4B.
- **Deliverables / APPROVED PRODUCT RULES:** Only Admin creates any account; Staff assists operationally. Individual borrower form offers STUDENT/FACULTY, role fixed BORROWER. Student requires unique textual Student ID and approved SKSU institutional email; Faculty is individual-only with valid unique accessible email and no required Student ID. Privileged STAFF/ADMIN management is a separate restricted Admin workflow. Standard Excel creation and separate bulk deactivation are Admin-only, Student-only; recommended studentId/name/email/course/contactNumber columns exclude passwords and role/category assignment. Complete validation preview precedes explicit confirmation; stable identity conflicts reject; results/retries/audit preserve history. Bulk deactivation excludes Faculty/Staff/Admin and shows unmatched/conflicts/already-inactive/obligations. DEC-073 allows Admin individual/Student bulk deactivation regardless of fines, active/overdue borrowing, unreturned equipment and replacements, with warnings/confirmation; no obligation veto. Preserve fine/payment/stock/return/replacement/borrowing/history evidence and overdue calculations; separate resolution/fine clearance remains authorized and auditable. Borrowers set their own passwords through secure activation and both categories accept current official FSMO terms. No hard deletion or roster-absent deactivation.
- **REMAINING TECHNICAL DEPENDENCIES:** Batch1 core API/storage/import/string-ID/matching/concurrency/idempotency/audit/activation contracts are implemented and tested; actual approved SKSU Student domain configuration/roster inputs (OPEN-028); operational ownership/recovery/existing-account procedures (OPEN-001/009); selected Brevo API-key/verified-sender/permitted link-origin configuration and tested live activation/future recovery delivery to Student institutional or Faculty accessible email (OPEN-017). Test adapters never prove real mail delivery. Approved mockup assets remain unchanged; later UI behavior must follow DEC-070 and the updated visual map.
- **REMAINING INSTITUTIONAL APPROVALS:** Exact approved SKSU Student domains/official roster format inputs (OPEN-028), operational activation/recovery responsibility (OPEN-001/009), official terms wording/version/publication after FSMO presentation (OPEN-016). Deployment-only Brevo API-key/verified-sender/successful delivery prerequisites remain OPEN-017 technical/external dependencies, not vendor approval. Brevo selection and Student deactivation regardless of obligations are resolved; OPEN-029 is closed. Official terms do not block independent account-management/inventory development; they remain mandatory before live borrowing/official consent. Creation authority, category/email/Student-ID requirements and Student-only bulk scope are resolved, not gates.
- **Non-goals:** No public signup, Google/Microsoft SSO, institutional mailbox passwords, password/role/category assignment in the standard Student spreadsheet, password exports, global single-domain/student-only eligibility, hard account deletion, roster-based Faculty/Staff/Admin deactivation, independent suspension table or campus hierarchy. Phase6 is explicitly authorized in Batch1 after the recorded core gate; no Phase7 work.
- **Exit gate:** Real role/current-account, domain/identity, preview/selection, duplicate/idempotency/race/rollback/audit/history and activation hash/expiry/single-use/reissue/invalidation tests; bounded queries/files/minimal private data; configured Brevo API key/verified sender and successful live delivery tests before delivery claims; Student deactivation tests cover all open-obligation cases with warnings/confirmation and unchanged loan/fine/stock/replacement/due/overdue/history evidence. Terms uses Phase4B evidence and cannot fabricate official wording. Unresolved gates are reported explicitly, not marked complete.

## Phase 6 — Equipment & Inventory

- **Status:** COMPLETE within reviewed Batch1 catalog/inventory scope. Core Go/frontend/fresh PostgreSQL/narrow privileges/history/concurrency/HTTP/Chromium gates PASS; see [report](PHASE6_REPORT.md). Future liability integration remains explicitly guarded, not implemented Phase7 borrowing. Phase7 is not authorized.

- **Goal:** Implement catalog pools and four-count physical inventory with append-only evidence.
- **Deliverables:** Bounded discovery, metadata/version, A/R/C/damaged_held constraints, stock ledger and guarded archive; categories/images as reviewed.
- **Dependencies:** Phases2.5–5 sequencing and explicit Batch1 authorization; catalog content/storage details and specific original-disposition/correction policy if included. Pending official terms approval does not block independent inventory development (DEC-072); live borrowing still requires official publication/acceptance.
- **Non-goals:** No serialized tracking, mandatory repair workflow, lost-history physical bucket or inventory edits concealing custody.
- **Exit gate:** Stock reconciliation, permissions, stale edits/archive races and history checks pass; unusable originals remain outside available.

## Frontend reconstruction — final local checkpoint

- **Stage A:** Engineering complete; owner approval confirmed by Stage B authorization and reconfirmed in the final checkpoint request.
- **Stage B:** Engineering complete; owner visual approval confirmed in the filter review and final checkpoint request.
- **Inventory filter:** “Available for borrowing” intersects ACTIVE + available>0 with independent Status/Category/Search before server counts/totals/paging. Owner manual verification **PASSED**, explicitly confirmed in the current checkpoint reply, 2026-10-09.
- **Acceptance:** Final visual/functional reconstruction milestone **ACCEPTED**. [Checkpoint audit/handoff](FRONTEND_RECONSTRUCTION_CHECKPOINT.md) separates the local Git snapshot from production readiness; existing Stage A/B/filter reports retain dated engineering evidence.
- **Remaining external gates:** Official FSMO terms after presentation, authoritative publication and borrower acceptance before live borrowing; approved production Student domains and required secure activation/ownership readiness; Brevo configured key/verified sender and verified live delivery. No placeholders or provider-submission-only claims substitute for these.
- **Technical limitations:** Original unmatched GET 500 lacks path/request-ID evidence; inherited lazy-login timing sensitivity and 19 lint warnings; unknown-total endpoints, physical assistive-device/non-Chromium and production deployment checks remain documented.
- **Boundary:** Authorized local checkpoint only. No push, deployment, new policy or Phase7 implementation.

## Phase 7 — Borrowing & Approval

- **Status:** NOT STARTED / NOT AUTHORIZED. Accepted frontend reconstruction provides reusable components and reviewed contracts; request separate Phase7 scope/implementation authorization before work.

- **Goal:** Implement reserve-on-submit and approval WITH physical release.
- **Deliverables:** PENDING holds,24h expiry/sweep/lazy check, owner pending cancel, required visible denial reason, atomic CHECKED_OUT issue/direct checkout, required due date+time/Asia-Manila/current terms, retained history and idempotency.
- **Dependencies:** Phases4–6 and current transaction contract; live borrowing requires officially published FSMO terms and documented acceptance via Phase4B (DEC-072/OPEN-016). No unresolved reservation/approved-waiting/seven-day gate.
- **Non-goals:** No separate approved-release stage, operating-hours TTL, quota, automatic old-fine prohibition or issued delete/cancel.
- **Exit gate:** Last-stock/decision-expiry/account races and duplicate/replay/rollback tests prove exact holds/custody/history. Expiry release is part of borrowing acceptance, independent of email delivery.

## Phase 8 — Returns & Accountability

- **Goal:** Implement partial returns, damage/loss replacements and full Admin fine clearance.
- **Deliverables:** Immutable unique return/acceptance lines, obligations and acquisition vectors, physical0+replacement0 completion, original due clock, PHP 10 ceiling24h live/final fine and immutable full clear methods/history.
- **Dependencies:** Phase 7, selected stock/fine design; original-disposition detail only gates an included disposal feature, not core replacements.
- **Non-goals:** No automatic damage price charge, partial payment/allocation, online gateway, return-photo/evidence upload, transmission, storage, retention or attachments; no borrower evidence step. Any personally shown phone photo remains outside eLabTrack; Staff/Admin alone records returns. No mandatory repair workflow; independent catalog images remain in Phase 6.
- **Exit gate:** All 18 domain traces proven by relevant unit/real-PG/HTTP tests, duplicate/cross-parent/concurrency/atomicity and Admin-only checks.

## Phase 9 — Notifications & Scheduled Work

- **Goal:** Deliver approved confirmations/reminders with observable failures.
- **Primary deliverables:** Provider adapter, consistent messages, bounded scheduled jobs, delivery/retry/deduplication records and tests.
- **Dependencies:** Phases 7–8; email provider, due/overdue timing and delivery responsibilities agreed.
- **Explicit non-goals:** No unapproved SMS/push provider or distributed broker.
- **Exit gate:** Agreed triggers are reliable, retries are safe, failures visible, and test delivery/schedule boundaries pass.

## Phase 10 — Dashboards & Reporting

- **Goal:** Provide authorized, source-defined operational insight.
- **Primary deliverables:** Approved metrics, bounded dashboard/report queries and required export/print formats.
- **Dependencies:** Phases 5–9; Phase 2.5 lifecycle/replacement/fine definitions and access matrix; final approved report formats.
- **Explicit non-goals:** No generic analytics platform, unbounded global reads or guessed revenue.
- **Exit gate:** Metrics reconcile with authoritative records; authorization, pagination/export limits and accuracy tests pass.

## Phase 11 — Interactive Equipment Catalog / Kiosk Experience

- **Goal:** Complete and validate the academic Interactive Kiosk requirement as normal responsive borrower browse/search/filter/quantity/cart/review/request UX.
- **Deliverables:** End-to-end catalog/cart usability/accessibility, responsive refinements and operator guidance; reuse earlier implemented canonical catalog/request flows, no duplicate app.
- **Dependencies:** Approved Phase 3A/3B flow and Phases4/6/7; Phase 9 notices where included. No device registration, shared-device identity, handoff or hardware policy dependency.
- **Non-goals:** No dedicated touchscreen shell, device tokens/authentication infrastructure, physical tracking or native application.
- **Exit gate:** Mobile-first cart/request and larger-screen adaptations preserve stock/ownership/terms semantics and pass accessibility/usability tests; functional kiosk intent documented.

## Phase 12 — V1 Data Migration

- **Goal:** Migrate only approved legacy data with reconciliation and recoverable cutover.
- **Primary deliverables:** Authorized discovery, mappings, reconciliation rules, dry-run tooling/results, backup/rollback/cutover plan.
- **Dependencies:** Approved scope/data access; stable Phase 2 model and implemented Phases 5–11.
- **Explicit non-goals:** No live V1 mutation without separate authorization or copying unverifiable data assumptions.
- **Exit gate:** Owners sign off reconciled dry run, exclusions and balances; backup/rollback proven before authorized cutover.

## Phase 13 — QA, Security & Optimization

- **Goal:** Validate complete workflows and address measured security/performance risks.
- **Primary deliverables:** End-to-end regression, security review, concurrency/load/accessibility checks, defect fixes and release evidence.
- **Dependencies:** Phases 1–12 with test data, expected workloads and acceptance targets agreed.
- **Explicit non-goals:** No speculative features or infrastructure scaling without measured need.
- **Exit gate:** Required quality gates and risk-based tests pass; critical/high findings resolved or explicitly accepted by accountable owners.

## Phase 14 — Deployment, Pilot & Evaluation

- **Goal:** Deploy through approved controls and evaluate real FSMO operation.
- **Primary deliverables:** Approved deployment/backups/monitoring, operator training, pilot, incident/rollback procedures and evaluation results.
- **Dependencies:** Phase 13 release approval, migration/cutover signoff and operational owners.
- **Explicit non-goals:** No campus rollout or automatic expansion proposal approval.
- **Exit gate:** Pilot meets agreed operational criteria; incidents/lessons and stakeholder evaluation recorded; future work remains separately approved.

## Future separate project — Campus-Wide Equipment System Proposal

- **Goal:** evaluate whether measured V2 outcomes justify a broader school proposal.
- **Primary deliverables:** separately requested scope/feasibility/governance proposal informed by the FSMO pilot.
- **Dependencies:** Phase 14 evidence and explicit school sponsorship.
- **Explicit non-goals:** no campus architecture, tenancy or cross-department implementation in current V2.
- **Exit gate:** independent stakeholder decision on whether to authorize that future project.
