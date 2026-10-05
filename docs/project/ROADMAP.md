# Initial delivery roadmap

Confirmed sequence from current stakeholder direction. **Only Phase 0 is authorized in this task.** Later deliverables describe phase-level outcomes, not a speculative task backlog or permission to execute them. Dependencies include the relevant accepted policy decisions, not merely a checked phase number.

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

## Phase 2 — Domain & Database Design

- **Goal:** Agree domain boundaries, lifecycle vocabulary and integrity model before product tables.
- **Primary deliverables:** Approved domain model, SQL schema/migration design, constraints/indexes, transaction/concurrency strategy, history/audit/retention design.
- **Dependencies:** Phase 1; stakeholder decisions affecting roles, eligibility, stock and accountability.
- **Explicit non-goals:** No speculative campus entities or implementation of later workflows before design approval.
- **Exit gate:** FSMO model and policy-dependent contracts approved; invariants, migration reversibility and concurrency expectations reviewable before business-domain tables are created.

## Phase 3 — UI Design System & Application Shell

- **Goal:** Create the eLabTrack visual language and responsive accessible shell.
- **Primary deliverables:** Project tokens, shadcn compositions, navigation/layout, loading/error/empty states, accessible responsive patterns.
- **Dependencies:** Phases 1–2 and role/navigation direction sufficient to design the shell.
- **Explicit non-goals:** No full business dashboard, inventory/loan implementation or replacement UI framework.
- **Exit gate:** Shell and tokens approved; keyboard/mobile/theme checks pass; primitives remain reusable.

## Phase 4 — Authentication & Authorization

- **Goal:** Implement confirmed product account and permission policy on hardened infrastructure.
- **Primary deliverables:** Approved role/operation matrix, access/recovery/onboarding flows, current-account authorization, verification/session behavior and tests.
- **Dependencies:** Phases 1–3; decisions on provisioning, staff/admin, Super Admin, suspension and verification.
- **Explicit non-goals:** No guessed public signup, campus roles or borrower/equipment operations.
- **Exit gate:** Server denies unauthorized/inactive access and stale privileges per approved policy; session/recovery and critical abuse tests pass.

## Phase 5 — Users & Borrower Management

- **Goal:** Manage approved borrower identities and account lifecycle.
- **Primary deliverables:** Authorized provisioning/profile/eligibility flows, bounded lists, lifecycle history and tests.
- **Dependencies:** Phase 4; borrower types and account/terms policy; Phase 2 data model.
- **Explicit non-goals:** No stock, loan/return behavior or campus hierarchy.
- **Exit gate:** Approved actor matrix and lifecycle scenarios pass; private data stays scoped and queries bounded.

## Phase 6 — Equipment & Inventory

- **Goal:** Implement authorized aggregate inventory and discovery with historical integrity.
- **Primary deliverables:** Inventory administration/discovery, approved categories/statuses, stock constraints, archival behavior and tests.
- **Dependencies:** Phases 2–5; taxonomy, deletion/archival and inventory policy decisions.
- **Explicit non-goals:** No loans/returns, serial tracking or RFID/barcode/QR scope.
- **Exit gate:** Inventory invariants, authorization, bounded search and lifecycle/history tests pass; no unsafe direct client writes.

## Phase 7 — Borrowing & Approval

- **Goal:** Implement confirmed request, decision and checkout lifecycle.
- **Primary deliverables:** Requests, approve/deny, direct checkout, active borrowing, stock reservation/release, cancellation/audit rules and tests.
- **Dependencies:** Phases 4–6; limits/due dates/reservation/denial/cancellation decisions and Phase 2 concurrency model.
- **Explicit non-goals:** No fine/payment or returns redesign beyond approved boundary contracts.
- **Exit gate:** Concurrent/duplicate requests cannot corrupt stock; decisions/history and authorization follow approved policy.

## Phase 8 — Returns & Accountability

- **Goal:** Handle partial/final returns and approved damage/loss/fine policy reliably.
- **Primary deliverables:** Return transitions, quantity reconciliation, accountability/fine assessment and approved settlement evidence/history.
- **Dependencies:** Phase 7; fine amount, payment/waiver, retention and history decisions.
- **Explicit non-goals:** No speculative online payment gateway or reporting module.
- **Exit gate:** Return/stock/fine/history changes are atomic or explicitly recoverable; duplicate/concurrency and calculation tests pass.

## Phase 9 — Notifications & Scheduled Work

- **Goal:** Deliver approved confirmations/reminders with observable failures.
- **Primary deliverables:** Provider adapter, consistent messages, bounded scheduled jobs, delivery/retry/deduplication records and tests.
- **Dependencies:** Phases 7–8; email provider, due/overdue timing and delivery responsibilities agreed.
- **Explicit non-goals:** No unapproved SMS/push provider or distributed broker.
- **Exit gate:** Agreed triggers are reliable, retries are safe, failures visible, and test delivery/schedule boundaries pass.

## Phase 10 — Dashboards & Reporting

- **Goal:** Provide authorized, source-defined operational insight.
- **Primary deliverables:** Approved metrics, bounded dashboard/report queries and required export/print formats.
- **Dependencies:** Phases 5–9; precise status/fine assessment/collection definitions and access policy.
- **Explicit non-goals:** No generic analytics platform, unbounded global reads or guessed revenue.
- **Exit gate:** Metrics reconcile with authoritative records; authorization, pagination/export limits and accuracy tests pass.

## Phase 11 — Interactive Kiosk

- **Goal:** Provide approved walk-in access with shared-device safety.
- **Primary deliverables:** Kiosk-oriented responsive flow, session/idle/reset behavior and operator guidance.
- **Dependencies:** Phases 4/6/7/9; kiosk authentication/privacy and hardware/operating expectations agreed.
- **Explicit non-goals:** No physical tracking integration, native app or assumed anonymous borrowing.
- **Exit gate:** Walk-in journeys and device handoff preserve identity/privacy; usability and unattended-session checks pass.

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
