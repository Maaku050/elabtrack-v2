# Phase 2.5 — Product Decision Integration & Domain Rebaseline report

Date: **2026-10-08, Asia/Shanghai**. **COMPLETE — documentation/domain-design reconciliation only.** All 29 owner authoritative working decisions integrated. Core workflow policy blockers resolved. **Phase 3A readiness: YES, pending separate authorization; NOT STARTED.** No business implementation, migrations, UI, mockups, Phase 1 session changes, boss writes, commit, push or deployment.

Owner direction governs working policy; this report does not imply independent institutional signoff, deployed behavior or passing business tests. Exact hybrid stock, lock/schema choices and live fine-clearance checkpoints are selected engineering recommendations. [Source hierarchy](SOURCE_OF_TRUTH.md), [decision register](DECISIONS.md), [remaining details](OPEN_DECISIONS.md) and [domain baseline](../domain/DOMAIN_MODEL.md) distinguish them. The original [Phase 2 report](PHASE2_REPORT.md) is preserved as a superseded historical handoff.

## 1. Files changed

**14 Markdown files:13 tracked updates and1 new report.** Updated all 8 docs/domain files: DOMAIN_MODEL, BUSINESS_RULES, STATE_MACHINES, DATA_MODEL, INVARIANTS, V1_COMPATIBILITY, BOSS_REBUILD_REFERENCE and API_RESOURCE_DRAFT. Updated5 docs/project files: DECISIONS, OPEN_DECISIONS, ROADMAP, SOURCE_OF_TRUTH and PHASE2_REPORT. Created this PHASE2_5_REPORT. No separate stakeholder packet.

Sources reviewed: root AGENTS, charter, source/decision/open/roadmap/Phase 2 documents, all domain documents, stack/manifests and relevant supplied V1 audit04/06/07/08/09/13 identity/inventory/borrowing/history/requirements evidence. Current Phase 1 foundation/auth roles and migration contracts remain implementation evidence; boss audit/source findings are prior read-only Phase 2 technical reference, not product authority. No private capstone identities/contact examples copied, no V1 Firebase access.

## 2. Decisions integrated

All owner D1–D29 mapped explicitly in [BUSINESS_RULES](../domain/BUSINESS_RULES.md#decision-integration-matrix). DEC-051 covers population/account/fine eligibility;052 roles;053 provisioning;054 reservation/atomic physical issue/direct;055 expiry/cancel/deny;056 due/timezone;057 terms;058 partial/physical/replacement/completion;059 PHP 10/day with adopted ceiling engineering detail;060 money/full Admin clearing;061 kiosk/responsive correction. DEC-062 separately records selected hybrid inventory as **Proposed engineering baseline**, not invented institutional policy. Original 27 OPEN IDs retain source conflicts with current resolutions/precise NON-BLOCKING residuals. The subsequent explicit D17 owner correction (2026-10-08) is integrated throughout this package: personally stored photos may be physically shown in person, entirely outside eLabTrack; Staff/Admin alone records returns and no photographic evidence feature exists. Other working decisions, inventory arithmetic, table/invariant counts and readiness are unchanged.

## 3. Borrower/account model

Generic BORROWER, at least Student/Faculty categories, not distinct permissions. Active registered account sufficient; no borrower-profile eligibility/suspension subsystem or automatic outstanding-fine gate. Admin deactivation blocks login/request/normal authenticated actions; old history/obligations survive and Staff/Admin can resolve them. Current Phase 1 `user`/`admin` source is unchanged; future approved identity work will reconcile product roles.

## 4. Final Staff/Admin distinction

Staff handles request review/physical approval/denial/direct checkout, partial/final returns, damage/loss, accepted replacements, borrower operational history and normal inventory. Staff can individually provision Borrower accounts via D4's narrow grant. Admin inherits those operations plus deactivation/bulk import/Staff-Admin management, exceptional stock correction, policies/terms, administrative audit/reports and **exclusive fine clearance**. Multiple named Admins; no shared credential.

## 5. Account provisioning model

No public signup. Staff/Admin provision Borrowers; Admin bulk creation/import required. Conceptual category/program/contact fields need purpose; secure activation/initial-password delivery remains implementation design. Legacy plaintext-password spreadsheets are not a V2 requirement; no insecure password/import mechanics selected here.

## 6. Final request lifecycle

Stored **PENDING, CHECKED_OUT, DENIED, CANCELLED, EXPIRED, COMPLETED**. Partial/overdue/replacement conditions derived. Canonical transaction/lines remain through terminal outcomes. No APPROVED-waiting or issued cancellation/delete state. [Machine and18 traces](../domain/STATE_MACHINES.md).

## 7. Reservation behavior

Submit immediately moves available−q/reserved+q in one server transaction. Locks validate available, unique item IDs/current account/terms, and atomic history/audit/outbox/receipt. Cart edits are local intent until committed submission. Old fine permits reservation; staff may later deny in person.

## 8. Pending-expiry behavior

expires_at=submitted_at+24h captured policy; at now>=expiry one retained EXPIRED edge releases every hold. Configurable future TTL candidate, no operating-hours calculation. Bounded sweep plus lazy pending-command checks proposed; approval/expiry serialize on borrowing. Lazy expiry must commit before reporting expired409, never roll back its necessary release. Until processed, a persisted hold remains real reserved stock; worker lag observable. No scheduler built.

## 9. Approval/release model

Staff/Admin approval WHILE handover is one PENDING→CHECKED_OUT transaction, reserved−q/checked_out+q. Recheck active target/catalog/held stock/due, capture decision/issue snapshots and audit/outbox/receipt. No separate release endpoint. Full-request issue selected engineering assumption; any changed quantities require a new request.

## 10. Direct-checkout model

Staff/Admin selects active existing borrower, appropriate items/quantities and future due date/time; server availability/account/terms/actor checks. Immediate CHECKED_OUT, available−q/checked_out+q, no pending hold artifact. Recommended current borrower terms evidence required, not staff impersonation.

## 11. Due-date/time model

Required issued due_at with date+time; Asia/Manila interpretation, absolute PostgreSQL timestamptz. No fixed seven-day maximum or invented grace/quota. Future optional maximum nullable in typed policy until expressly approved. Staff confirms requested due intent at physical issue; issued deadline freezes through partial/replacement resolution.

## 12. Terms acceptance model

One immutable user/current-version/accepted_at evidence at first activation/use; applies to later loans while version current. Material new version gates next request; no per-loan checkbox. Pending requests keep bound submission acceptance. Approved final copy and secure onboarding remain later review details, not navigation blockers.

## 13. Final overdue definition

Currently overdue iff CHECKED_OUT operationally unresolved and now>due_at. Unresolved means any physical OR replacement outstanding. C0 with loss/damage awaiting replacement remains overdue after deadline. PENDING/DENIED/CANCELLED/EXPIRED never accrue. Completed borrowing retains historical late/fine evidence but is not currently operationally overdue.

## 14. Fine calculation

Adopt owner's recommended ceiling24h interpretation as deterministic engineering detail:

```text
end = now while unresolved; completed_at after completion
elapsed = max(0, end-due_at)
days = ceil(elapsed / 86,400 seconds)
assessed_minor = days * 1,000 PHP centavos
```

Exactdue0;1min PHP 10;24h PHP 10;24h+1min PHP 20. Bounded integer arithmetic, no floats/daily mutation. Completion freezes final days/amount/time. No evidenced conflict requires reopening rounding; legacy closed amounts keep historical formula provenance rather than recalculation.

## 15. Partial-return behavior

Staff/Admin authoritative return; unique immutable lines. Good quantity restocks, damage/loss leaves physical custody and creates replacement quantities. Any remaining physical/replacement keeps CHECKED_OUT and original deadline/fine clock. Corrected D17: Borrower approaches Staff/Admin in person → may physically show a photo personally stored on their own phone → Staff/Admin may inspect actual equipment if desired → Staff/Admin alone records authoritative condition/quantities in eLabTrack → the system updates borrowing and inventory. The photograph is completely outside the software. Borrowers cannot mark returned or submit return evidence. No return-photo acceptance, upload, transmission, storage, retention, attachments, file model, borrower evidence interface or evidence-upload step. Independent equipment catalog images and structured Staff/Admin return history remain in scope.

## 16. Damage/loss replacement model

Damage/loss incident creates corresponding replacement obligation, never automatic equipment-price peso charge. Confirmed loss is not returned stock; damage disposition receives a nonusable original. Staff/Admin accepts proper type/operational equivalent, reduces liability and acquires usable replacement stock. Incident survives. Detailed damaged-original repair/discard/retain/retire is NON-BLOCKING for borrowing UX and not a mandatory workflow.

## 17. Final inventory arithmetic model

Selected justified hybrid over six historical buckets or usable-only counts:

```text
T = A + R + C + D       # tracked physical units, D nonusable originals held
usable_pool = A + R + C
on_premises = A + R + D
mixed return(g,d,l): delta(A,R,C,D;T) = (+g,0,-(g+d+l),+d;-l)
accept replacement(q): (+q,0,0,0;+q)
```

[Full vectors and examples](../domain/DOMAIN_MODEL.md#selected-inventory-model-and-alternatives) cover submit/release/direct/good/partial/loss/damage/acceptance and later original disposition. Loss+replacement restores physical total. Retained damage original+new replacement legitimately increases total by the new unit; only replacement is usable. No lost/retired current bucket, ghost custody or double original restock. UI awaiting-replacement counts derive liability by kind, not C/D.

## 18. Replacement-obligation model

Immutable basis per source return-line/kind, required_qty equals corresponding damage/loss. Remaining=required−sum(immutable accepted lines), nonnegative. Immutable acceptance header/unique obligation lines, appropriate pool/type/equivalence and named operator; reject duplicate IDs/cross-borrowing refs/overacceptance. Multiple obligations to same pool aggregate exact acquisition quantities. Borrowing-first locks serialize final return/acceptance/closure.

## 19. Completion rule

Every item physical=issued−good−damage−loss and replacement=required−accepted. Complete only when **all physical0 AND replacement0**. Final resolution and inventory/event/fine freeze/audit/outbox/receipt commit atomically. Fine balance does not prevent completion; clearing fine does not complete borrowing. No separate complete button or historical quantity rewrite.

## 20. Fine-clearing model

Admin full current outstanding only, PAID/WAIVED/OTHER_RESOLUTION, original/final assessed basis and clear amount/time/named actor/method/optional note preserved. No partial chosen amount. Staff denied. Selected engineering checkpoint recommendation permits clearing live fine without stopping unresolved accrual: paid10 at due+1h; assessed20 at25h leaves new outstanding10; final20 freezes on completion, later full clear10 recorded separately. Each clear covers entire outstanding at its time, not installments. Expected-balance guard, same-key receipt and borrowing/fine locks prevent stale/double discharge.

## 21. Final monetary-accountability model

Focused overdue_fines+fine_clearances; no independent records/fines/paid flag truth, generic charges/adjustments/payment allocations or automatic damage valuation. Outstanding=live/final assessed minus clearances. PAID is recorded offline FSMO payment; waiver/other are not cash. Reports separate assessment, paid, waived/other and outstanding; no online gateway or refund/reversal architecture.

## 22. Final kiosk interpretation

**Interactive Kiosk = borrower browse/search/filter/quantity/cart/review/request**, normal web/mobile responsive flow. Borrower mobile-first, Staff/Admin desktop/tablet-first responsive. DEC-061 explicitly supersedes kiosk-shell parts of DEC-038/039. Academic name retained with functional meaning; Phase 11 reframed to Interactive Equipment Catalog / Kiosk Experience.

## 23. Dedicated kiosk entities removed/deferred

Removed device registration/credentials/tokens, device-authenticated kiosk sessions and handoff candidates from active design; they are **rejected for current scope**, not necessary deferred borrower architecture. No special kiosk API/auth namespace or hardware shell remains.

## 24. Updated role permission matrix

[BUSINESS_RULES future matrix](../domain/BUSINESS_RULES.md#future-permission-matrix) covers browse/own submit/cancel/history, operational review/issue/deny/direct/return/incidents/replacements/inventory, Borrower provisioning/deactivation/import, privileged users, Admin-only fine/policy/audit/report operations. Staff does not automatically inherit Admin or borrow on behalf of users. Optional administrative pending cancellation remains an unselected scoped recommendation. Authorization server-authoritative, current actor/target checks and ownership/replay included.

## 25. Updated table candidates

[DATA_MODEL](../domain/DATA_MODEL.md#candidate-inventory--26-table-entries): **26 entries =3 existing +18 new core +2 conditional catalog +3 deferred legacy;23 new candidates total.** Core: equipment, inventory_ledger, borrowings, borrowing_items, return_events, return_event_items, replacement_obligations, replacement_events, replacement_event_items, borrowing_policy_versions, terms_versions, terms_acceptances, overdue_fines, fine_clearances, notification_outbox, notification_attempts, audit_events, command_receipts. Existing users/refresh_tokens/schema_migrations unchanged; future users attributes/roles only proposed.

## 26. Tables removed/deferred

Remove borrower_profiles/separate eligibility; replace charges/charge_adjustments with focused fines/clearances; remove payment/allocation and kiosk device/handoff candidates. Categories/images remain conditional; legacy_import_batches/mappings/reconciliation_issues deferred to authorized Phase 12. Return-photo/evidence attachments and storage are excluded entirely from software scope, not optional or deferred features. Damaged-original disposal/repair remain nonmandatory later inventory work. Data model defines column/key/FK/CHECK/index/immutability/lock requirements; no SQL exists for these candidates.

## 27. Updated invariant count

**71 active +1 retired =72 catalog IDs.** Original 57→56 retained active +15 additions. ACC-08 generalized reversal rule explicitly retired. Added INV-13/14, BOR-09/10, ACC-09, REP-01–08, AUTH-05, ELG-05. [Catalog](../domain/INVARIANTS.md) includes expiry/history/atomic issue/unresolved overdue/no ghost custody/replacement bounds/completion/history/final freeze/Admin clear/duplicates/terms/inactive rules with enforcement and future test types.

## 28. Scenario walkthrough results

**18/18 required design traces reconcile** in [STATE_MACHINES](../domain/STATE_MACHINES.md#eighteen-required-walkthroughs): full/partial/late good returns; lost replacements before/after due; damaged replacements; mixed good/damage/loss; expiry; owner cancel; reasoned deny; direct checkout; old-fine new request; human fine-based denial; Admin clear; material terms change; last-stock race; duplicate return; concurrent replacement acceptance. Each identifies stock/physical/liability/clock/history and protections. No schema/state contradiction remains for those flows; traces are not executed tests. Future real DB concurrency/fault/HTTP tests remain required when implemented.

## 29. Remaining unresolved policies

All are **NON-BLOCKING for core UX**: exact original damaged disposition/equivalence guidance; taxonomy/content; secure activation/password/recovery/email ownership mechanics; institutional/legal retention; notification provider/sender/reminder cadence; authorized legacy dataset/reconciliation/cutover; optional institutional session policy; final deployment/backups/monitoring/ownership. [OPEN_DECISIONS](OPEN_DECISIONS.md) names responsible review and affected later activity. No stale reserve/approved/seven-day/price-fine/device-auth policy gate.

## 30. Remaining blocking policies

**None for core borrower, review/physical issue, return/replacement, fine and role-navigation UX or their approved feature planning.** Specific dependent activities remain gated: damaged-original repair/disposal needs its inventory procedure; production notifications need provider/cadence; actual migration needs authorized export/reconciliation/cutover; destructive cleanup needs legal retention; deployment needs Phase 14 controls. They do not block borrowing mockups. No business phase authorized by this report.

## 31. Phase 3A readiness

**YES — ready, not started; separate authorization required.** Working flows and role/navigation stable. Explicit mockup assumptions: full-request issue; current terms acceptance also for direct first-use; nonusable held originals until later disposition; Staff/Admin certification of equivalent replacement into original pool; Admin live full-clear checkpoints and later accrual; draft catalog/category/terms copy and final Clear Fine wording; no return-photo/evidence software feature or upload step; any personally shown phone photograph stays outside eLabTrack. Optional admin pending cancellation is excluded unless specifically needed. No dedicated kiosk hardware shell.

## 32. Business implementation readiness by domain

“Ready” below means core policy/design supports a later approved feature plan, not implemented acceptance or authorization. Appropriate implementation design/review/testing and the sequenced approved UI/foundation phases still apply.

| Domain | Readiness | Specific remaining boundary |
|---|---|---|
| Identity / Borrower management | Core policy ready; onboarding mechanics design still required | Secure activation/recovery/import before affected implementation; future role migration must preserve Phase 1; no signup/eligibility policy blocker |
| Equipment / Catalog | Core ready for approved planning | Final taxonomy/content and image storage before respective delivery; normal catalog/cart does not need hardware |
| Inventory | Core hybrid/locks/reconciliation ready; specialized original disposition not ready | Repair/discard/retire and exceptional correction procedure for those features; baseline retains nonusable originals |
| Borrowing | Core ready | Approve future feature plan with due/input/TTL/full-issue implementation detail; prove expiry/last-unit/actor races |
| Returns / replacements | Core ready; separate disposal feature gated | Appropriate-type discretion supported; no mandatory repair; prove duplicates/remaining/atomic acquisitions/closure |
| Overdue fines | Core ready with selected ceiling/live-clear engineering recommendations | Validate integer boundaries/checkpoint expectations and Admin-only atomic clear; no rate/payment policy blocker |
| Notifications | Outbox/events ready; production delivery/scheduling not ready | Provider/sender/cadence and operational retry ownership needed before affected implementation/production sending |
| Reporting | Core definitions/access ready; final deliverable formats pending | Approved report measures/filters/export formats and bounded queries; PAID distinct from waived/other; no generic analytics platform |
| Interactive catalog/kiosk | Core ready | Approved mobile-first Phase 3A mockups, implementation/usability/accessibility proof; no device/handoff dependency |

## 33. Exact git status

Starting working tree was clean. Final `git status --short`:

```text
 M docs/domain/API_RESOURCE_DRAFT.md
 M docs/domain/BOSS_REBUILD_REFERENCE.md
 M docs/domain/BUSINESS_RULES.md
 M docs/domain/DATA_MODEL.md
 M docs/domain/DOMAIN_MODEL.md
 M docs/domain/INVARIANTS.md
 M docs/domain/STATE_MACHINES.md
 M docs/domain/V1_COMPATIBILITY.md
 M docs/project/DECISIONS.md
 M docs/project/OPEN_DECISIONS.md
 M docs/project/PHASE2_REPORT.md
 M docs/project/ROADMAP.md
 M docs/project/SOURCE_OF_TRUTH.md
?? docs/project/PHASE2_5_REPORT.md
```

## 34. git diff --check

**PASS — exit0, no whitespace diagnostics.** Final `git diff --check` checks13 tracked updates; the new report was also explicitly checked for trailing whitespace and final newline. Local Markdown paths/anchors, unique invariant IDs/counts,29 decision rows,26 table rows,18 scenario rows and35 report sections checked.

Before/after SHA-256 fingerprints of all311 originally tracked V2 files confirm only the13 listed Markdown updates; the only new untracked file is this report. No Go, React/TypeScript, SQL/migration, mockup, configuration, Phase1 source or existing UI primitive changed. Boss tracked-file fingerprints also unchanged; no boss writes performed. No private capstone data copied; `.project-reference/` remains ignored.

D17 correction validation: all14 existing Phase2.5/domain Markdown files were updated relative to the correction-start fingerprint; no other workspace file changed or was added. Searched `photo`, `photograph`, `image`, `evidence`, `attachment`, `return evidence`, `supporting evidence` and `upload`; earlier optional return-photo software wording is replaced by the explicit outside-system boundary. Catalog image candidates and requirements are preserved. Local links/whitespace pass, and29 decision rows,26 table entries,18 walkthroughs,71 active+1 retired invariants and35 report sections are unchanged.

Runtime Go/frontend suites, real database/browser/Docker/deployment checks **not rerun**, as explicitly requested for Markdown-only reconciliation. The18 scenarios are design reviews, not executed tests. Prior Phase1 results remain historical; no business runtime verification claimed. No Git commit/push or deployment.

## 35. Whether Phase 2.5 exit gate is satisfied

**YES.** All 29 working decisions integrated; generic borrower/three roles/provisioning, immediate holds,24h expiry, retained pending outcomes, physical approval/direct issue, required date+time/no seven-day cap, versioned first-use terms, PHP 10 deterministic unresolved clock, replacement-based damage/loss, joint completion, Admin full-clear history, corrected responsive kiosk meaning, reconciled inventory/schema,71 active invariants/18 traces, resolved stale open gates and Markdown-only validation satisfy the requested gate. Phase 3A/mockups/business code/migrations and all Git publication/deployment actions remain unstarted.
