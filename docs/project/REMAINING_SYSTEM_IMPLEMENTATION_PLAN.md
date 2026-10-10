# Remaining system implementation plan — Phases 8–14

2026-10-10. Owner-authorized sequential implementation from accepted checkpoint `c2741c5`, after DEC-081. DEC-082 records this renewed authority. No commit, push, deployment or normal database modification is authorized. Final engineering and owner acceptance are separate. Existing Phase7 demo and the eleven privately shelved Phase8 files remain preserved.

## Authority and boundaries

Approved policy: BORROWER Student/Faculty, STAFF, ADMIN; current activation/terms enforcement; immediate reservations,24h expiry, physical checkout; operational returns/replacement, PHP10 per loan per started24h overdue, Admin full manual clear with PAID/WAIVED/OTHER_RESOLUTION. Deactivation never resolves obligations. Damaged originals stay held; no new repair/disposal workflow, return photos, payments, due extensions or V1 imports.

Existing contracts: Phase7 services/transactions, current-account locks, protected catalog images, consistent REST envelope, server pagination, query/client-state boundaries, shadcn/Base UI and persistent shells. DEC-080 and both original V1 screenshots govern selection interaction, not their prices/wizard or old styling.

Safe technical choices: additive paired migrations; borrowing aggregate lock before sorted equipment locks; immutable disposition/acceptance/clearance evidence and exact ledger sources; durable database notifications with recipient/event uniqueness; bounded filters/CSV limits and formula escaping; protected raster profile images through the existing validator; explicit isolated presentation targets and private generated credentials.

External decisions: official terms/production domains/activation readiness remain live-use gates. Live Brevo is deferred. Detailed damaged-original disposition and institutional production reminder cadence remain outside this implementation. In-app reminder scheduling must be configurable; any demo cadence is explicitly a demo fixture, not silently approved production policy.

## Sequential deliverables and gates

| Phase | Deliverable | Required exit evidence |
|---|---|---|
|Preflight|Review/reproduce shelved000009 failure, inventory every target, preserve original files, validate corrected candidate in rollback-only isolated transactions|Root cause, rollback proof, normal read-only baseline, no applied000009 on current targets; no unresolved destructive/privilege issue|
|8|Partial mixed returns, damage/loss obligations, physical replacement, automatic all-obligation completion/frozen fine, Admin full clear, visible history and real account/archive projections|Domain/fine boundaries, PostgreSQL/concurrency/idempotency/rollback/authorization/ledger, fresh+representative restore migration rehearsal, frontend and Chromium|
|9|Durable scoped notification center/count/read state, transaction event recipients, configured due/overdue reminders and restart/concurrent catch-up|Recipient isolation, event dedup/read persistence, actual restart/concurrency and browser checks; no mail|
|10|Truthful role-scoped dashboards, bounded operational/admin reports and filtered CSV|Metrics reconcile to records; filters/counts/authorization; deterministic CSV/formula safety; browser|
|11|Images/facts/adjacent quantity+Add, cart editing/removal, mobile review, desktop browse/cart, reference/layout refinements|DEC-080 image/reference comparison and mobile/desktop/light/dark Chromium; real multi-item API/stock/conflicts/duplicate guards|
|12|Fresh-start protected profile images and dedicated isolated presentation:40 equipment,20 Students,5 Faculty,2 Staff,1 Admin, avatars/local catalog assets and actual synthetic lifecycle histories|Actual authentication/consent, dataset/ledger/metric reconciliation, wrong-target reset rejection, repeatable setup/reset with retained volumes|
|13|Full-system regressions/security/technical-debt verification|Serial+parallel frontend, Go/race, SQL/constraints/rollback/concurrency, restore preservation, real Chromium role/workflow/light/dark/responsive, no secret artifacts|
|14|Fresh setup/offline start/stop/health/reset/recovery and presentation walkthrough|Clean isolated installation, no hidden credential dependency, blocked external requests/download-free runtime, actual presentation scenarios; normal preservation|

Each phase is planned, implemented, tested, integrated, verified and documented before the next phase. A failed gate is diagnosed/corrected/rerun; no success by assertion. Genuine policy/security/destructive/permission blockers stop dependent work. Track actual files, SQL/API/UI changes and results in [durable progress](PHASE8_14_PROGRESS.md) and phase reports. Phase14 is the stopping point; final owner acceptance remains pending.

## Database and recovery protocol

Never rewrite000001–000008, grant normal runtime privileges, reset normal volumes or publish demo terms in normal storage. Capture normal state read-only without publishing row data or password hashes. New migrations use verified disposable targets only. Full private backup restores use the authorized bootstrap owner only on disposable databases (citext ownership), then verify actual rows/schema before forward migration. Test empty rollback/reapply and history-bearing rollback refusal; never delete history to force down migration. Preserve all prior volumes and private paused artifacts. Prepare a separate reviewed normal rollout procedure; do not execute it.

## Presentation and handoff

Dedicated presentation data must be fictional, illustrated, local and clearly labeled. Reuse real versioned terms and authentication; never seed real-user consent or claim email verification. Credentials are generated privately, passed without argv and retrieved in the owner's terminal. Preserve accepted five-account/eight-equipment Phase7 demo. All40 catalog images and28 profile avatars must load offline. Seed lifecycles through services with controlled clocks where possible, with exact stock/fine/report/notification reconciliation.

Final report must identify phase statuses, migration preflight, actual gates/results and failures, dataset counts/images/avatars, normal preservation, limitations, exact Git status and safe next actions. Claim engineering completion only once every required gate is executed and passed; owner acceptance is not automatic.
