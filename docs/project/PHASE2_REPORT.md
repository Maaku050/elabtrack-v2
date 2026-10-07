# Phase 2 — Domain & Database Design report

> **Historical Phase2 handoff snapshot, superseded by Phase2.5 (2026-10-08).** The BLOCKED status,57 invariants,15 traces,28 table entries and policy alternatives below describe the original delivery, not the current baseline. Current owner decisions and readiness are in [PHASE2_5_REPORT](PHASE2_5_REPORT.md); current domain documents contain71 active invariants plus1 retired ID,18 walkthroughs and26 table entries. Core UX policy blockers are resolved; Phase3A is ready for separate authorization and remains unstarted. Original phase evidence is retained, not rewritten as if answers existed earlier.

Date: **2026-10-08, Asia/Shanghai**. **Design documentation delivered for review; Phase 2 BLOCKED FOR IMPLEMENTATION.** No institutional policy was silently selected. State/schema-sensitive decisions remain open, so the unconditional COMPLETE/exit gate is not claimed. Phase 0 and Phase 1A–1I remain complete; their source and evidence are preserved.

The package defines the FSMO domain contract, candidate PostgreSQL structure, reconciliation/concurrency rules, seven evidence labels, **57 future invariant specifications** and **15 conditional scenario traces**. They are design specifications, not executed business tests or implemented tables. No business source, configuration, migration, mockup, UI, boss repository, V1 Firebase, commit, push or deployment was changed. No Phase 3A began.

## 1. Files changed

Created eight domain documents:

- [DOMAIN_MODEL.md](../domain/DOMAIN_MODEL.md): scope, evidence, modules, entities, ownership, vocabulary and snapshots.
- [BUSINESS_RULES.md](../domain/BUSINESS_RULES.md): sourced rules, policy alternatives and 20-row decision gate.
- [STATE_MACHINES.md](../domain/STATE_MACHINES.md): formal transitions, diagrams, return examples and 15 scenarios.
- [DATA_MODEL.md](../domain/DATA_MODEL.md): 28-entry existing/candidate inventory, fields/keys/FKs/checks/indexes, ER, movements/locks/retention.
- [INVARIANTS.md](../domain/INVARIANTS.md): 57 identified rules with reason/enforcement/operation/future test type.
- [V1_COMPATIBILITY.md](../domain/V1_COMPATIBILITY.md): preserve/change map, source shapes, unknowns, provenance/reconciliation.
- [BOSS_REBUILD_REFERENCE.md](../domain/BOSS_REBUILD_REFERENCE.md): read-only provenance, salvage/rejection map and critical-defect guardrails.
- [API_RESOURCE_DRAFT.md](../domain/API_RESOURCE_DRAFT.md): resource/command/mobile pagination/error draft, no handlers.

Updated [OPEN_DECISIONS.md](OPEN_DECISIONS.md), [DECISIONS.md](DECISIONS.md), [ROADMAP.md](ROADMAP.md), and [SOURCE_OF_TRUTH.md](SOURCE_OF_TRUTH.md). Created this **PHASE2_REPORT.md**. **13 Markdown files total: nine new, four modified.** SOURCE_OF_TRUTH is the required documentation correction to reflect the expressly updated authority order and seven labels; stale capstone locators in the open register were also corrected. No Phase 1 verification report was rewritten.

## 2. Sources reviewed

Current request and AGENTS; PROJECT_CHARTER, SOURCE_OF_TRUTH, DECISIONS, OPEN_DECISIONS, ROADMAP, FOUNDATION_AUDIT, PHASE1_SECURITY_BACKLOG, PHASE1_FOUNDATION and PHASE0_REPORT; docs/ARCHITECTURE and STACK; integration/README and the existing Phase 1I evidence/closure context; current backend/frontend manifests; current users/refresh SQL and account/transaction port implementation. The six historical SQL files remain unchanged.

All **16 V1 audit documents** (README and 01–15) were reviewed: architecture/roles/data/auth/equipment/borrowing/notifications/history/rules/performance/debt/deployment/features/planning/undeployed hardening. Their live-deployment limits remain explicit.

The local ignored capstone DOCX was inspected for engineering/business text: SOP P364; borrower benefits P401/402; scope P419–432; definitions P436–443; Sprint 1 P604–608; borrowing Sprint 2 P613–617; returns P620–624; email/reminders P628–636; kiosk P639/640; relational description P702; registration/denial interface descriptions P827/P851; conclusion P912; expansion recommendation P922. Only rule intent/locators were retained, never personal/sample content. Several original paragraph citations were stale and are now corrected, without changing the unresolved policy evidence.

Boss audit `/home/marvin/projects/elabtrack/docs/audit/BOSS_REBUILD_AUDIT.md` (2026-10-07, boss snapshot `42b8c7c5d8b10ba13de0c7c53dfd007cbcd50cfc`) plus selected read-only inventory entity/repository, borrowing entity/service/requests/returns/charges, return/charge migration structure, notification repository, reporting repository, file access and kiosk service. Referenced source paths and inspected-versus-audit-only limits appear in BOSS_REBUILD_REFERENCE. No boss write/port occurred. Its historical current-V2 migrator/session findings predate completed 1H/1I.

Primary PostgreSQL documentation verifies structural constraint and lock boundaries: [constraints](https://www.postgresql.org/docs/current/ddl-constraints.html), [explicit locking](https://www.postgresql.org/docs/current/explicit-locking.html). These support technical enforcement choices, not FSMO policy.

## 3. FSMO scope conclusion

**CONFIRMED REQUIREMENT / CURRENT V2 DECISION:** existing FSMO aggregate equipment borrowing transaction and tracing modernization. No tenancy, campus/unit/lab hierarchy, cross-department borrowing, global role platform, individual assets/physical tracking, native mobile, offline platform or speculative distributed infrastructure. Audience-specific responsive direction and existing shadcn/Phase 1 boundaries remain intact.

## 4. Domain modules

Eleven conceptual modules: Identity & Access, Borrowers, Equipment & Inventory, Borrowing, Returns, Accountability/Charges, Notifications, Reporting, Audit/History, Settings/Policies, Interactive Kiosk. Each has purpose/entities/rules/dependencies/nonownership in DOMAIN_MODEL. Recommended implementation is pragmatic Go domain/application modules and ports in the one API; Returns can remain within Borrowing. No scaffolding or services were added.

## 5. Canonical terminology

Equipment stock pool/unit quantity; stable User versus eligible Borrower; borrowing transaction/request/reservation/approval/release/checkout; active/partial/completed custody; derived due_today/overdue; good/damaged/lost dispositions; immutable assessment/adjustment/waiver/settlement versus actual payment/collection; distinct history/audit/ledger. V1 actual Ondue/Overdue/Incomplete spellings retain raw provenance. Public wording remains OPEN-015 where institutional approval matters.

## 6. Proposed user/borrower model

Preserve current UUID users/session infrastructure. Proposed one-to-zero/one BorrowerProfile carries approved category/institutional identifier/eligibility/manual verification evidence/version. Active login, institutional eligibility and email ownership are distinct. Categories, onboarding, actor matrix, sanctions and verification gates remain OPEN-001–009/020/027; no generic role is presumed eligible. Current actor and target checked again at business mutation time.

## 7. Proposed equipment model

One aggregate pool per catalog entry; name/description/condition observation, optional category/narrow image, value/currency, catalog lifecycle, current buckets and versions. No serial/unit table. Proposed active/inactive/archived is separate from counts; metadata updates cannot set custody counters or alter old snapshots. Archive guard/condition/disposition policy remains OPEN-012/022.

## 8. Proposed inventory buckets

**ENGINEERING RECOMMENDATION:** available, reserved, checked_out, damaged, lost, retired; accounted total sums all six. Only active eligible available stock is borrowable. Accounted total is not physical on-hand stock or an assertion that lost/retired units remain owned. Loss/retirement recognition/write-off rules await OPEN-022.

## 9. Inventory invariants

Nonnegative integer buckets; exact row sum; reserved equals held line quantities; checked_out equals total outstanding issued custody; cumulative dispositions equal immutable event sums; every count equals signed ledger baseline/deltas; damaged/lost/retired reconcile after repair/recovery/removal. Row checks alone cannot enforce cross-entity truth. INV-01–12 define enforcement and future tests.

## 10. Inventory ledger approach

Current Equipment row is authoritative for commands, immutable movement ledger is supporting reconciled evidence. Explicit initial/add/remove/reserve/release/issue/good-damage-loss/repair/retire/recovery/correction vectors carry sequence/actor/reason/source/time. No event sourcing. Every quantity command writes current counts and ledger/audit together; correction preserves original evidence.

## 11. Borrowing state machine

Proposed pending→approved→checked_out→completed with conditional combined approve/release and direct checkout; pending→denied or permitted cancellation; conditional approved cancellation before issue. Partial/due/overdue are projections. All terminal canonical transactions persist. Approval/release/cancellation choices remain open; formal transition table gives actors, guards, stock/history/notification effects and failures.

## 12. Reservation recommendation/status

Compare A request-time hold, B approval-time hold, C hybrid/expiry for oversubscription, expectations, utilization, workload, complexity and V1 compatibility. Recommend A if staff confirm pending hold expectations; B is explicit alternative; defer C absent justified operational need. **No option selected**; OPEN-020 is blocking dependent implementation. Explicit per-line reserved quantity accommodates A/B; C needs another reviewed expiry design.

## 13. Approval versus physical-release recommendation/status

Recommend distinct decision and actual issue timestamps/state; allow a combined atomic command only if staff policy confirms simultaneous release. V1/boss merge them; capstone says release follows approval without defining the interval. **OPEN-021 is blocking** the final transition/API/schema contract.

## 14. Direct checkout design

Existing active eligible target borrower, current authorized operator, active equipment/capacity, approved due/duration/exception, borrower-attributable terms and frozen issue evidence are mandatory proposed checks. Same canonical issue transaction and idempotency as requested checkout. No staff bypass based on target scope/grant alone; OPEN-020/002/003/016 decide permission/exception specifics.

## 15. Snapshot strategy

Request evidence is frozen separately from decision and actual issue evidence. Minimal borrower display/ID and equipment name/category/value/currency snapshots, governing policy/terms/acceptance refs, distinct times and source channel. Historical joins use those snapshots, not mutable current profiles/prices. Binding at request/approval/checkout is **OPEN-023**; recommendation freezes liability at actual issue unless request binding is confirmed.

## 16. Terms/policy versioning

Immutable published typed FSMO policy and terms content/hash/version; no generic config engine/CMS. Acceptance evidence belongs to the user and exact version/context; borrowing references what it used. Once-per-account/version versus each-borrowing acceptance and assisted consent are **OPEN-016**, affecting final constraints. No legal text/defaults fabricated.

## 17. Due/timezone design

UTC exchange/timestamptz event instants; one explicitly selected FSMO IANA calendar, retained for old policy. Asia/Manila is V1 evidence/recommendation only, not accepted V2 timezone; developer Asia/Shanghai is not policy. Date-only cutoff versus exact instant, today/maximum/grace/cap/calendar day basis remain OPEN-015/020/006. One rule must govern DTOs/assessment/reminders/reports; inclusive date filters become half-open instant ranges.

## 18. Partial-return design

One immutable ReturnEvent header with unique owning BorrowingItem lines, good/damaged/lost quantities and notes. Cumulative item totals/outstanding updated from that one normalized input in the same transaction as movement, evidence and approved assessments. Same canonical borrowing remains until every issued unit is accounted. Examples cover full/partial/mixed/multiple events and failure rollback.

## 19. Duplicate-return protection

Reject duplicate borrowing-item IDs before mutation; DB UNIQUE(event,item); composite same-borrowing event/item FKs; one validated map for custody/events/stock/ledger/charge basis. Boss's issued2/two good1 lines must reject without moving stock. Same-key retry returns original event; different-key concurrent over-return conflicts under borrowing lock. RET-01–08/CMD-01 specify future real regression tests.

## 20. Damage/loss design

Explicit disposition accounts custody into damaged/lost buckets; never restocks good or silently reduces accounted total. Authorized later repair/recovery/retirement/disposal appends new movement without erasing old return/charge evidence. Notes/authority/write-off/correction policy OPEN-022, asset/value rule OPEN-023. No physical per-unit condition tracking introduced.

## 21. Charge/accountability design

One immutable Charge assessment, append-only ChargeAdjustments, derived nonnegative balance under charge-row lock. Source/window uniqueness and frozen basis/currency prevent duplicate/cumulative posting. No records/fines/paid/status parallel truth. Proposed positive reversals restore only prior discharge adjustments; original referenced amount caps reversal. Custody completion and financial closure are independent.

## 22. Payment/settlement and overdue policy status

V1 PHP10/day and snapshot-price formula are factual, **not accepted V2 rates**. Amount/currency/day/grace/cap/rounding/binding/posting/partial-return basis remain OPEN-006/015/023/024. Recommend visibly unposted active estimates, immutable posted assessments at approved event. Administrative settlement differs from actual receipt/payment/revenue; OPEN-007 determines procedure/authority. Payment/Allocation tables conditional only, no gateway or cash procedure assumed.

## 23. History model

Retain canonical request/decision/issued/completed transaction and items; immutable return events, charge adjustments, policy/terms evidence and snapshots. No active→records copy/delete, no denied/manual-deleted history reconstruction. Ledger/audit/notification evidence have distinct purposes, not substitute transaction truth.

## 24. Business audit model

Minimal durable who/action/when/entity/reason/meaningful before-after/context and server correlation. Required high-value decision/issue/return/disposition/stock/archive/account status/role/charge changes commit audit in same transaction. Restricted append-only evidence, no every-read audit, request body/credential or operational-log substitute. Legal retention/access remains OPEN-026.

## 25. Notification/outbox model

Atomically inserted business event intent, unique event identity (ReturnEvent ID per partial return), bounded future claims, fenced lease token, immutable attempts and safe outcome classes. No ignored persistence error. Provider acceptance is distinct from recipient delivery; duplicate delivery after provider acceptance/persistence failure remains possible. Provider/trigger/cadence/recipient/worker operations OPEN-017/025; no SMTP or worker implementation here.

## 26. Reporting design implications

Authorized bounded relational queries for stock/availability/active and overdue custody/unique borrowers/history/returns/damage-loss/usage/assessed-adjusted-waived-settled-outstanding accountability. Shared approved calendar/currency, historical versus current names explicit. No report tables by default, unbounded collections, campus aggregates or assessment-as-revenue. Export formats/storage remain later requirements, purpose authorization rechecked.

## 27. Kiosk design boundary

Additional touch-first/privacy-aware client of same borrowing domain. Narrow catalog, optional distinct device identity, expiring attributable handoff and idle reset are reference concepts. OPEN-010 must select authentication/assistance; no device protocol, anonymous borrowing or ordinary borrower cookie hydration is chosen. Late-result/reset fencing and all-page browsing are later acceptance requirements.

## 28. File/attachment model

Prefer conditional narrow EquipmentImage metadata/association, verified image content/size/type, safe replacement and bounded orphan cleanup. No unrestricted generic bucket shared with private reports/borrower/evidence documents. Other independently authorized file purposes remain deferred; public catalog access never grants private export permission. **Current D17 correction:** return photographs are wholly outside eLabTrack, only personally shown on the borrower's phone in person. Staff/Admin alone records authoritative condition/quantities. No return-photo/evidence submission, upload, transmission, storage, retention, attachments or evidence-upload workflow is a software requirement; catalog images are unaffected. No storage implementation now.

## 29. Idempotency strategy

Selected submit/direct/return/financial-disposition commands use unique receipt bound to current actor, operation, concrete resource/context, source channel and normalized payload hash. Same input replay returns original entity/event; different payload conflicts; reauthenticate/authorize before any private replay. Claim/result commit with business transaction. Retain compact key/hash/entity binding while effect exists; expiry cannot allow duplicate recreation. Metadata uses expected version. No credential or requestId as idempotency authority.

## 30. Concurrency/locking strategy

Receipt claim → ordered User/eligibility rows → existing Borrowing header → deduplicated Equipment UUID order → affected Charges (conditional payment parent first). Recheck after locks; all account/catalog/custody writers honor same order. Archive holds Equipment and reads guarded item sums without taking Borrowing locks afterward. READ COMMITTED/conditional updates/unique constraints, no application-only mutex. Atomically write every participating fact/evidence/outbox/receipt; failures roll back all. Unknown commit outcome uses original receipt key.

## 31. Proposed table count/list

**25 new candidates evaluated + three existing = 28 inventory entries.** Sixteen proposed core: borrower_profiles, equipment, inventory_ledger, borrowings, borrowing_items, return_events, return_event_items, borrowing_policy_versions, terms_versions, terms_acceptances, charges, charge_adjustments, notification_outbox, notification_attempts, audit_events, command_receipts. Two conditional catalog: equipment_categories/equipment_images. Three deferred provenance: legacy_import_batches/legacy_mappings/legacy_reconciliation_issues. Two conditional receipts: payments/payment_allocations. Two conditional kiosk: kiosk_devices/kiosk_handoffs. Existing users/refresh_tokens/schema_migrations unchanged. No blanket implementation of this count. DATA_MODEL describes purpose/PK/important columns/FKs/unique/check/index/mutability/history for every candidate and provides Mermaid ER.

## 32. Major constraints

Nonnegative/widened bucket arithmetic; unique canonical reference/equipment line; per-item issued/disposed arithmetic; duplicate event-line uniqueness; composite same-borrowing ownership FKs; required snapshot/version evidence at approved stages; unique charge source and command/event claims; restrictive historical refs; immutable-evidence grants/triggers. Cross-row custody/ledger/balance/eligibility checks belong to authoritative application transactions and real reconciliation tests, not impossible giant CHECK expressions.

## 33. Index strategy

Owner-before-page chronological history; active-state/deadline queue; catalog lifecycle/name/category; item equipment reconciliation; chronological return/movement/charge/audit lookups; pending and expired-lease outbox claims; unique source/command keys. Stable time/name plus ID tie-breaker. PK/unique indexes not duplicated, volatile now() partial predicate avoided; search/materialization justified only by later measured workload.

## 34. Delete/archive strategy

Stable User/account references preserved through deactivation; equipment inactivate/archive rather than consequential deletion; categories inactivate, images retire/stage cleanup; borrowing/returns/ledger/charges/adjustments/audit/versions/acceptance ordinarily retained. No blanket soft-delete column or cascade history loss. Outbox attempts/receipt identity retained under explicit later retention. Down migrations are not safe operational undo once history exists; later authorized backup/forward-repair review required.

## 35. V1 migration compatibility and privacy gaps

Map actual users/equipment/embedded multi-item transactions/terminal records/fines/notification shapes/settings/images through separate provenance. Reconcile reserved versus issued borrowedQuantity, cumulative disposition, snapshot prices/names/due dates, duplicate record/ref collisions and conflicting fine-clearing evidence. Missing original dates/return events/denials/actors/consent/payment cannot be invented. Imported cumulative baselines/unknown time accommodation require reviewed migration-specific amendments or quarantine. No importer or live V1 access.

Identity/profile, custody, finance, audit and notification data are private with owner/capability scope; no invented legal retention duration. Evidence-preserving anonymization/access/cutover is OPEN-026/018, blocking dependent procedures. No capstone biography/contact/sample user content or secrets copied to repository.

## 36. Boss concepts adopted

ADOPT CONCEPT as engineering recommendations: atomic ordered stock locks, optimistic metadata versions, movement evidence, retained canonical terminal loans, durable business audit. This means design ideas, **no code port**. Detailed provenance/conditions in BOSS_REBUILD_REFERENCE.

## 37. Boss concepts adapted

Adapt buckets/accounted-total semantics, ownership/locking for one FSMO, request-versus-issue snapshots, immutable return events with duplicate/ownership fixes, base assessment+adjustments, typed policy/terms, outbox event keys/leases/errors, bounded report SQL/calendar/currency and future provenance concepts. No boss defaults become approved policy.

## 38. Boss concepts rejected

Hierarchy/global-super-admin platform; duplicate-return interpretation; inactive-target checkout bypass; generic export-file access; unsafe flattened importer/fabricated dates/partial claims/writing dry-run; copied migrations/numbers; repeated-return key suppression; inconsistent due/grace/report definitions; wholesale replacement of verified auth/config/shadcn; startup migration/unsafe seed/plaintext DB or showcase/binaries.

## 39. Boss concepts deferred

Kiosk device/handoff implementation until OPEN-010, actual cash/receipt branch until OPEN-007, async exports/CSV utility until required approved formats/permissions/workload, import implementation until Phase 12 authorization/data decisions. No unsupported organization role platform is deferred into core scope; it is rejected for this project.

## 40. Invariant count

**57 unique IDs:** Inventory12, Borrowing8, Returns8, Accountability8, History5, Authorization4, Eligibility4, Command idempotency4, Notifications4. Each has precise rule/reason/enforcement/affected operations/future unit/real-PG/concurrency/HTTP/browser test type. They are future acceptance specifications, not tests added or reported passing now.

## 41. Scenario walkthrough results

**15 design traces documented, zero executed business scenarios.** Normal request/decision/release/full return; denied; conditional cancellation; staff direct; partial/final; good+damage; loss; last-stock competition; competing approval/checkout; duplicate/concurrent returns; inactive borrower; archive race; charge adjustment; historical value change; provider failure after commit. Every trace states expected quantities/evidence/invariants and remaining policy. Mixed return examples also cover duplicate payload and rollback after first of multiple pools. Conditional representability is not approval/readiness.

## 42. Blocking product-policy decisions

Dependent domain implementation blockers: OPEN-001/002/003 onboarding/category/actor authority; OPEN-008/009 eligibility/verification gates; OPEN-020 reservation/duration/direct/limits; OPEN-021 approval versus actual release; OPEN-015 due calendar/format; OPEN-016 acceptance; OPEN-013 cancellation if enabled; OPEN-006/023/024 currency/formula/value/policy-binding/posting; OPEN-007 settlement/payment authority/procedure; OPEN-022 disposition/repair/write-off/correction policy. The 20-row matrix records each database/domain impact and whether final dependent work can proceed without it.

The most direct state/schema choices are **approval/release, A/B/C reservation, date-only versus instant due, account-level versus per-borrowing acceptance, binding/posting basis, and administrative settlement versus receipt tracking**. Confirm these before freezing migrations/contracts. No numeric defaults or branch selection were silently accepted.

## 43. Non-blocking core decisions and later gates

Optional category taxonomy (011), denial-reason validation (014), image/archive details (012), provider/notification cadence (017/025), kiosk protocol (010), retention/anonymization/cutover (018/026), institution session-family/global behavior (027) and historical stack conflict (019) do not prevent independent core design discussion. They still gate their dependent UX/features/operations; archive/return/disposition edges need policy before implementation. Extra Super Admin 004/005 remains deferred/out of accepted scope, not a missing platform implementation.

## 44. Stakeholder confirmations required

FSMO operations/policy owner: eligibility and actual approval/release/cancellation/due/reservation/direct workflow. Accountable administrator/IT: provisioning, approved operation matrix, verification and later institutional session policy. Inventory custodian: categories/condition/archive/repair/loss/retirement/value evidence. Finance/accountability owner: rates/currency/assessment timing/waiver/correction/settlement/actual collection. Terms/document/data custodians: content/binding/acceptance/privacy/retention/migration. Kiosk/email operators with IT: assistance/auth/privacy/cadence/delivery ownership. No named approver or legal rule invented. Seven grouped OPEN sections preserve original IDs and add 021–027.

## 45. Recommended implementation order

Recommendation only, aligned with existing roadmap:

1. Confirm blocking policies; revise/approve Phase 2 state/schema/API/invariants and record actual decisions. Final schema/constraint/grant review before any business migration.
2. Separately authorize Phase 3A UX architecture/mockups from approved borrower mobile/staff/kiosk direction; Phase 3B shadcn design system/shell follows approved mockups, no business logic disguised as shell.
3. Phase 4 approved identity/operation/verification policy on preserved foundation; Phase 5 borrower profiles/provisioning and fresh eligibility/actor mutation protocol.
4. Phase 6 approved catalog/inventory schema, constraints, ledger/movement ports and locking; real PG stock/version/archive/reconciliation tests immediately.
5. Phase 7 canonical request/decision/issue, terms/snapshots/reservation/idempotency and atomic business audit/outbox insertion ports; real competing-stock/staff/eligibility tests. Audit/outbox schema must be available before commands depend on them; full delivery worker remains Phase 9.
6. Phase 8 immutable returns, partial/mixed closure, approved assessments/adjustments and all duplicate/rollback/concurrent reconciliation tests; payment receipt branch only if approved.
7. Phase 9 bounded notification/reminder/provider worker, leases/attempts/recovery/operations; Phase 10 approved scoped reporting/formats/calendar/financial metrics.
8. Phase 11 approved shared-device kiosk; Phase 12 separately authorized V1 adapter/discovery/reconciliation/cutover; Phases 13–14 QA/security/measured performance/operational deployment/pilot evidence.

No implementation step or next phase is started by this report.

## 46. Schema/domain implementation readiness and risks

**NOT IMPLEMENTATION-READY.** Reviewable structural proposals and stable integrity safeguards exist, but unresolved policy changes the final state/schema/constraints and command meaning. Policy alternatives are deliberately not turned into a generic runtime config system. Main risks: treating design as accepted policy, using stock CHECK alone, inconsistent normalized return interpretations, stale target authority, mutable liability/snapshots, unverified receipts mislabeled revenue, missing provider/retention operations and guessed V1 lineage/dates. Required mitigation is policy confirmation plus approved contract and later meaningful real PostgreSQL/HTTP tests.

## 47. Whether Phase 3A can begin

**No automatic/unconditional start; not begun and not authorized by this turn.** This package can inform a separately requested stakeholder review or UX planning of confirmed scope, but final role/workflow/reservation/terms/kiosk mockups require sufficient policy confirmation. OPEN-021/020/015/016 especially affect promise-bearing borrower and staff actions; OPEN-010 affects kiosk flow. Current Phase 2 instruction explicitly prohibits beginning Phase 3A.

## 48. Whether business coding can begin

**No.** Resolve the relevant blocking stakeholder policies, revise/approve model/API/constraints/test plan, follow separately authorized roadmap phases. No new domain schema, service, UI or production migration is created. Selected boss concepts are references, not folder-copy permission.

## 49. Validation / git diff --check

**PASS**, including new documents checked for whitespace without staging. Required docs are present; relative links/anchors, Markdown table shape/fences, seven labels, 57 unique invariant IDs, 15 scenario IDs, 20 gate rows, 27 unique OPEN entries and the 28-entry table count were checked. Diff/source-baseline inspection shows **only documentation changes**: all 298 other previously tracked files retain the captured SHA-256 baseline, including Phase 1 source/config/tests/reports, six historical SQL files, manifests/locks, shadcn and integration/V1 evidence. All new files are Markdown in docs/domain or this report. Exact preservation comparison is supplemental local evidence, not a production test.

No Go fmt/vet/test, frontend lint/test/build or Phase 1 integration rerun: current explicit request says not to rerun those for documentation-only changes. No database/browser/Docker/service/SMTP/import/deployment test is claimed. Prior Phase 1 results remain prior evidence; no gate or test was weakened. Read-only boss inspection created no changes there; no live Firebase accessed. Git remotes/status inspected before Git reads; origin remains `https://github.com/Maaku050/elabtrack-v2.git`, no remote/network Git operation/stage/commit/push.

## 50. Exact git status

Final `git status --short` (unstaged; Git collapses the eight untracked domain files into the directory entry):

```text
 M docs/project/DECISIONS.md
 M docs/project/OPEN_DECISIONS.md
 M docs/project/ROADMAP.md
 M docs/project/SOURCE_OF_TRUTH.md
?? docs/domain/
?? docs/project/PHASE2_REPORT.md
```

The full 13-file inventory is section1; no other tracked/untracked file change exists. No staging/commit/push/deployment.

## 51. Phase 2 exit gate

**Design deliverables and documentation validation delivered; unconditional Phase 2 exit gate NOT SATISFIED. Status: BLOCKED FOR IMPLEMENTATION.** Scope/modules/vocabulary/buckets/cross-invariants/returns/duplicate protection/one financial truth/distinct history-audit-ledger/candidate relational constraints-indexes/concurrency/V1 compatibility/boss map/test specifications and no implementation/diff boundaries are fulfilled. The additional policy-sensitive gate remains: state/schema choices are unresolved, so none can be claimed harmless to dependent implementation or final workflow mockups. Confirm the exact decisions in §42, revise the affected contract and obtain design approval before claiming COMPLETE or starting separately authorized work.

No foundation/toolchain/runtime failure is the blocker. The blocker is explicit institutional policy identified with evidence; documenting it completes the authorized independent design work without guessing the decision.

## Phase 2.5 supersession addendum

All29 authoritative working owner decisions are now integrated into the eight domain documents and DEC-051–061 (stock recommendation DEC-062). Borrower is generic with Student/Faculty categories; active provisioned account gate, three named roles and narrow Staff provisioning replace earlier open identity assumptions. Reserve on submit,24h EXPIRED, own pending cancel, visible denial and combined physical approval/release replace alternatives. Direct issue is immediate, date+time due uses Asia/Manila with no seven-day maximum; terms acceptance is once per current version.

Damage/loss is replacement liability separate from physical inventory. Current hybrid arithmetic A/R/C/damaged_held excludes lost/retired physical buckets; replacement acceptance is a new acquisition, not ghost custody or erased incident. Complete requires both physical and replacement zero. PHP10/day selected ceiling elapsed24h continues until completion and freezes. Admin full-clear PAID/WAIVED/OTHER_RESOLUTION retains assessed history; generalized finance/device-kiosk candidates are removed. Interactive Kiosk means one borrower mobile-first catalog/cart flow; Staff/Admin desktop/tablet responsive.

The original report's unresolved core-policy gates and separate-approved/device proposals are historical, superseded by current [OPEN_DECISIONS](OPEN_DECISIONS.md) and [Phase2.5 report](PHASE2_5_REPORT.md). Remaining disposal/equivalence/retention/provider/onboarding/migration/operations details are NON-BLOCKING for core UX but gate their affected later activities. Phase3A can begin after separate authorization; neither it nor any business implementation begins in Phase2.5. Updated design includes26 table entries (3 existing/18 new core/2 conditional catalog/3 deferred legacy),71 active invariants+1 retired ID and18 reconciled design walkthroughs, not executed business tests.

Phase2.5 validation is Markdown scope/hash/link/count review, `git diff --check` and exact status; original runtime validation references above remain historical. Phase1 source, migrations, frontend, mockups and boss files are unchanged. No commit/push/deploy.
