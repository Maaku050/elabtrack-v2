# Phase 9 — engineering exit gate

Engineering PASS,2026-10-10 under DEC-082. Owner final acceptance remains pending. Normal development and accepted Phase7 demo records were not modified. No live email.

Persistent notification dispatch derives from committed immutable borrowing events, return dispositions and replacement obligations. Unique source/recipient keys and a transaction advisory lock protect concurrent catch-up. Delivery and dispatch markers commit together; transaction rollback leaves neither. Startup catches up immediately and a bounded30s worker reconciles subsequently. Recipients and safe detail links are scoped to actual loan ownership/current operational role. List/count/read mutations independently recheck current active/activated identity; demotion hides retained operational notifications.

All16 required kinds are implemented: submitted, denied, cancelled, expired, checkout/direct, return, damage, loss, replacement-required/accepted, complete, fine-cleared, due-soon, overdue and fine-assessed. Financial notifications do not post charges or alter custody. Due lead is explicitly configured0–604800 seconds, default0 disabled; one reminder per kind/loan. Production reminder policy remains OPEN-025. The24h integration lead is demonstration-only.

Migration000010 passed real up/down/up on two isolated databases with original application values preserved; nonempty history makes down refuse safely. Runtime grants allow notification insert/read and read_at updates only. Immutable content/history enforcement was exercised. Historical migration1–8 remain unchanged.

Evidence:

- `TestRealNotifications`:4 PostgreSQL race-enabled groups, PASS; recipients, cross-user rejection, demotion/deactivation, rollback, four simultaneous workers, source-kind coverage, read/unread persistence, exact counts/nonoverlapping pages, grants and history preservation.
- `TestRealNotificationProcessRestart`: PASS; real product API starts with committed backlog, abrupt process kill/restart, exactly three expected notifications, no duplicate, read state retained. This is actual process recovery, beyond service reconstruction.
- `TestPhase9HTTP`: PASS; authentication, bounds/read allowlist, required boolean, unknown fields, allowed Origin, object ownership, read/unread and current inactive-token denial. Existing inactive-token middleware returns403; no authentication behavior changed.
- Frontend288 tests/25 files PASS with two-worker parallelism;4 new notification UI regressions. Lint PASS with19 inherited primitive warnings; production build PASS. Go fmt/vet/all tests PASS; ordinary suites skip guarded PostgreSQL tests, which were explicitly executed above. `git diff --check` PASS.
- Real Chromium5 checks PASS,4 fictional light/dark desktop/mobile screenshots in `docs/ux/verification/phase9/acceptance.json`; center, accurate badge, persisted read mutation/count invalidation, filter reload, owned deep link, Staff shell/action recipient. No page exceptions;390px no horizontal overflow.
- Read-only normal snapshot:21 tables unchanged against preimplementation values.

Initial harness failures are retained as evidence: fixture omitted explicit confirmation; HTTP assertion expected401 although existing inactive middleware correctly returns403; Base UI rendered navigation controls as button role, corrected to link role in new notification components; browser expected raw PENDING rather than displayed Pending, repeated test attempted to mark an already-read record, and reload sampled outgoing DOM. Fixture/assertion/repeatability/navigation semantics were corrected and gates rerun. No weaker authorization, timeout or business assertions were introduced.

Limitations:30s foreground UI polling and bounded worker eventual delivery; production repeat-reminder cadence is not approved; live Brevo is deferred. Historical recipient eligibility is captured when an undelivered source is reconciled, so a subsequently activated operator does not receive already-dispatched old alerts. No production notification retention/delete policy was invented. Phase10 may proceed.
