# Phase 4B — First-use terms acceptance and account onboarding

2026-10-08, Asia/Shanghai. **IMPLEMENTED FOUNDATION — VERIFIED. BLOCKED FUNCTIONALITY: official FSMO terms publication/institutional consent and credential delivery/activation.** Phase5 is not started or authorized. No commit, push or deployment.

## 1. Executive summary

Implemented immutable versioned terms, restricted immediate Admin publication, authenticated Borrower acceptance with original idempotent receipts, own acceptance status, first-use/updated/missing/error UI and a reusable future-command transactional policy. Existing Phase1/4A authentication/session architecture remains authoritative. Official text is not approved; no invented document or acceptance was seeded in normal application storage. Activation is an explicit owner-directed decision gate.

## 2. Sources reviewed

AGENTS.md, README, Makefile, PROJECT_CHARTER, SOURCE_OF_TRUTH, DECISIONS, OPEN_DECISIONS, ROADMAP, STACK/manifests; foundation/current authentication, JWT, account/password, migration, handlers/security/authorization/logging, session coordination, central transport/Query, protected routing, shells/tokens, tests and Phase1/2/2.5/3B/4A/local-readiness reports; relevant domain/UX documentation and approved B04-03 PNG. Read-only V1 audit04/05/09/13 identifies static unversioned terms and spreadsheet passwords; local capstone paragraph search provides product intent, not current official terms. No live Firebase or boss-repository changes. Current explicit owner Phase4B brief/clarifications outrank historical wording.

## 3. Official terms availability

Owner explicitly confirms: no current official V2 FSMO Terms and Conditions document is approved. V1 display text, capstone examples and approved visual mockup copy cannot establish governing wording. Synthetic TEST documents are labeled NOT OFFICIAL and confined to the disposable verification database. Exact final text/version and responsible publication approval remain OPEN-016/DEC-067. No institutional acceptance is claimed.

## 4. Existing onboarding architecture

Current users store UUID, email/name, bcrypt password hash, BORROWER/STAFF/ADMIN, is_active and timestamps. Internal creation accepts a supplied password and hashes it; no provisioning endpoint is mounted. Bootstrap seeds remain explicit development-only tools, not product provisioning. Existing login requires valid credentials/current active recognized role. There is no activation/invitation, password-change-required, email-verification or incomplete-activation state. Phase4A supplies current-account/session/route hooks, not credential delivery. Borrower category is future profile data, not an authorization role.

## 5. Terms database schema

New paired migration000005 only: `terms_versions` UUID/version unique/title/body/SHA-256/publisher FK/server creation and publication times; `terms_publication` constant singleton nullable current-version FK; `terms_acceptances` UUID/user FK/version FK/server accepted_at, UNIQUE(user,version) and UNIQUE(id,user). Publisher/version FK indexes supplement unique/PK indexes. Consequential FKs RESTRICT account/version deletion. UPDATE/DELETE/TRUNCATE triggers reject version/receipt history changes, including owner DML. Runtime gets SELECT/INSERT history and only current-pointer UPDATE; no history edits/deletes/DDL/migration tracking. No borrowing/inventory/business-policy/category tables.

## 6. Terms publication/version model

Admin-only POST publishes one immediate mandatory immutable version atomically. Exact UTF-8 plain text is preserved/hashed; body max65,536 bytes, title max200 bytes, version1–64 validated ASCII characters. Actor comes from current authenticated DB account; timestamps from PostgreSQL. Required expected-current UUID/null serializes competing publishers with one winner and an explicit conflict. No generic legal CMS, drafts, scheduled effective dates, media or withdrawal endpoint. A text change requires a new version. Software mechanics do not authorize unapproved official content.

## 7. Acceptance evidence

Identity is authenticated Borrower UUID, version is exactly the reviewed current UUID, time is server-generated. Payload is `{}` only; user_id, accepted_at, accepted flags, null/array/multiple JSON and unknown fields are rejected. INSERT ON CONFLICT plus unique user/version returns the same original UUID/time after repeats or parallel submissions. Old evidence remains; accepting an old version after replacement conflicts before replay. No IP/device/extra personal data is collected; DTO omits user identity. Opening/reading/checking UI alone never records consent.

## 8. Terms API contracts

[Implemented API specification](../API_CONTRACTS.md#phase-4b-implemented-terms-contracts): GET `/api/v1/terms/current` (active account), GET `/terms/status` (own Borrower), POST `/terms/{versionID}/accept` (Borrower), POST `/terms/versions` (Admin). Existing envelopes/request IDs/private no-store, trusted-Origin POST, JSON/body ceilings/rates apply. Errors: unpublished503; unknown404; changed/current-acceptance/duplicate-version/expected-publication conflicts409; input400/media415/body413/auth401/role-or-inactive403; safe DB500. Terms endpoints never issue cookies. Status is unpublished/required/updated/accepted; `can_initiate_borrowing` is eligibility only, no live borrowing feature.

## 9. Backend authorization/enforcement

Reuse the central role permission matrix and current-account middleware. Service repeats active/current-role checks under account read locks inside each transaction. Domain/application have no Fiber/pgx dependency. PostgreSQL repositories parameterize SQL/use the shared transaction context. `RequireCurrentAcceptance` requires its caller's transaction, target active Borrower and actual current receipt, retaining account/publication locks until caller commit. Missing/stale evidence fails closed. There is no fake borrowing endpoint and no claim of implemented Phase7 request enforcement.

## 10. First-use UI

Real `/borrower/terms` feature follows page → hook → TanStack Query → feature API → central transport, with server-state Query and local checkbox state. Approved Borrower shell/branding/tokens, shadcn checkbox, full plain-text article, version/status, readable page scrolling, initially unchecked consent, disabled pending submit and safe error/retry states. Theme/logout and account/existing-obligation links remain. Successful acceptance re-fetches server status and navigates only if the same authenticated generation/account remains.

## 11. Updated-terms UI

Historical evidence plus no current receipt displays Updated terms. Each new document UUID resets the checkbox. A409 from publishing during review fetches the new version, explains the conflict and requires a new explicit checkbox/action; never substitutes unreviewed content. Already accepted review displays actual server receipt time with a continue link. Missing/network/server failures never fabricate successful consent.

## 12. Protected route integration

Existing AuthProvider/ProtectedRoute/current-account reauthorization remain. Borrower home/equipment initiation paths are gated by current server status and matching receipt; initial/default/deep-link access redirects to terms without feature flashes. Terms is a Borrower-only route. Account, existing borrowings and notifications remain reachable without acceptance, preserving safe read-only/existing-obligation access. Restoring an explicitly requested read-only deep link can access that view; it does not unlock initiation. Staff/Admin workspace and Admin inheritance remain unchanged. No new server business mutation is exposed.

## 13. Logout/session handling

Memory-only access JWT/HttpOnly refresh cookie, single provider, central interceptor, database-authoritative status/roles, strict rotation/replay and Web Locks/BroadcastChannel are preserved. Terms Query keys are account-scoped/private and cleared by existing logout/role invalidation. Late acceptance cannot navigate or restore authority after logout/account generation change. Logout remains available on loading/unpublished/error/required states; no second session system/storage token is added.

## 14. Account activation design

Owner directs a decision gate rather than a silently selected delivery method. Concrete engineering proposals: expiring one-time activation link/code handed over in person by authorized FSMO personnel (recommended for this scoped operation), or an expiring one-time link delivered to an institutionally verified address. Either requires approved identity checks, channel ownership, expiry/reissue/lost-delivery handling, first-password semantics and accountable operational ownership. No endpoint/token model is implemented before those choices.

## 15. Credential provisioning decision status

OPEN-001/009, DEC-067: **UNRESOLVED**. Existing password hashing is not secure credential distribution. Required owner/FSMO/IT decision: choose the channel/one-time mechanism, authorize verification and lifecycle rules, and determine who provisions/delivers/reissues access. No predictable default/shared password, plaintext spreadsheet password column, invitation email or production delivery choice. Phase5 provisioning/import must use the approved method and server-forced permissions.

## 16. Password recovery status

Current V2 has no authenticated password-change/reset/recovery endpoint or link; V1 Firebase reset is historical, not V2 policy. Approve responsibility/identity verification/change and session-revocation rules before dependent implementation. No unrestricted reset, invented recovery link or public self-registration is added. First-password change is not falsely claimed complete.

## 17. Cross-tab behavior

Real tabsA/B accept the same version and retain one original receipt; tabC reloads current accepted status. Focus/reconnect/mount re-fetch uses server state; no sensitive consent/account/token broadcasts or new channel. A new mandatory version requires new review. Existing cross-tab logout clears peer memory/account/private terms cache and protected views; observed messages contain only existing allowlisted lifecycle fields. Tests do not claim every background tab immediately pushes a new terms publication.

## 18. Version-race handling

Lock order: current participating account(s), singleton publication row, future borrowing/equipment/history writes. Publication exclusively locks the pointer; acceptance/status/policy hold shared locks. Controlled real-PG tests observe lock waiters: publication-first rejects old acceptance; acceptance-first commits reviewed evidence then publication makes current status updated. Two publishers sharing one expected pointer have exactly one success. Twelve parallel acceptances return the same original evidence. Rollback removes uncommitted versions/pointer changes/receipts.

## 19. Missing-terms behavior

Own status returns unpublished, no current document/receipt, `can_initiate_borrowing:false`; acceptance_required:false means no published version exists to accept, not consent. Current read/acceptance/policy returns TERMS_NOT_PUBLISHED503. No checkbox/action pretending official agreement. Home/catalog initiation paths remain unavailable; own account/existing obligations/notifications/logout remain permitted. Service/transport errors are separate safe manual-retry screens. Normal DB is deliberately empty of terms until exact wording is approved and an authorized Admin publishes it.

## 20. Light/dark visual fidelity

Actual screenshots inspected against approved B04-03 and accepted Phase3B shell/tokens. Retained FSMO brand, navy/violet, typography/card/version/status/button hierarchy and bottom navigation. Runtime adaptations disclosed: long synthetic plain text instead of unapproved institutional copy, initially unchecked checkbox instead of static checked mockup, existing accepted branding, visible logout/theme/read-only links, acceptance action after the entire document rather than floating over it. Dark/wide views use established tokens. This is verified implementation fidelity, not a new owner approval or approved legal copy.

## 21. Responsive results

Real320/390/768/1280 required screens in both themes pass no horizontal page overflow and effective controls≥44px. Checkbox indicator16px uses a full-width≥44px label. Missing320/390, reading end390, accepted/updated/stale/error390 and Staff/Admin1280 screenshots also pass. Long article remains complete/page-scrollable; keyboard Space selects explicit consent. No excessive animation/new palette, token or shell redesign.

## 22. Backend tests

PASS `go fmt ./...`, `go vet ./...`, `go test ./...`. New domain/application validation/status/authorization/exact-version/idempotence/missing/policy-transaction/publication/conflict/rollback coverage and six safe error-contract cases. Extended role permission matrix retains existing cases. PASS targeted `go test -race ./internal/application/... ./internal/domain/... ./internal/infrastructure/persistence/postgres ./internal/interface/http/...`. Default test command deliberately skips external opt-in fixtures; actual external results are separate below. Selected Go1.27.1 works; global module-external launcher remains1.26.5, unchanged.

## 23. Frontend tests

PASS199 tests/11 files, all172 existing cases preserved plus21 feature cases and6 transport-code cases. Existing Borrower auth-test setup supplies accepted TEST status for the newly required gate; old assertions remain. Coverage includes first/returning/updated/current/malformed/missing states, exact `{}` API submission, pending duplicate fencing, network/server/stale/version reset, expired/inactive, logout/late result/private cache, read-only paths, Staff/Admin, theme/long escaped text. PASS lint: zero errors, same19 inherited shadcn Fast Refresh warnings, zero new warnings. PASS strict TypeScript/production build, main479.41kB, no chunk-size warning. No manifest/lock changes or weakened checks.

## 24. PostgreSQL integration

PASS12 `TestRealTerms` subtests under race against isolated PostgreSQL18.6: fresh empty terms down/up preserves all accounts/sessions; missing fails closed; roles/publication attribution/hash; first/repeat/account scope; new/old retention;12 simultaneous acceptance requests; acceptance/publication rollback; unique/FK/NOT NULL/history immutability/runtime DML-DDL/RESTRICT checks; both lock-race orders; competing publishers; populated down refusal and outer-transaction enforcement. Failed initial expectations/setup were corrected: actual ON DELETE RESTRICT SQLSTATE23001 and explicit foundation grants. No assertion/permission was weakened. Final suite passed on a freshly owned isolated database before HTTP/browser publication; historical data was never deleted to rerun it.

## 25. Browser checks

PASS real Chromium140.0.7339.16:21 checks/26 screenshots, zero application errors/warnings. Actual random-password login/forms/API/PG, missing → first review → pending acceptance → reload/logout/login → updated → publication during review409/new unchecked review → transport error/retry; same original receipt across tabs; third tab reload; cross-tab logout/cache; inactive denial; Staff/Admin unaffected; memory-only token/HttpOnly cookie. The harness's interrupted runs were corrected for initial document readiness/explicit post-logout deep-link selection/fresh first-use identities and headless tab activation; only the final full run is counted. [Artifacts](../ux/verification/phase4b/README.md), [reproduction](../../integration/PHASE4B.md). Separate inherited preview-only harness PASS24 layouts/10 checks/production exclusions with anonymous refresh fixture, not authenticated acceptance. No actual production user or fake authenticated API response.

## 26. Security regressions

PASS real existing Foundation suite8 subtests plus terms HTTP4 under race: safe own account/health, generic login, role matrix/stale authority, rotation/replay/logout/rollback/missing account/retention; terms auth/role denial, cross-account/time/flag/null/array/multiple JSON rejection, trusted Origin, current version/repeat/conflicts, stale Admin JWT and inactive Borrower. Concurrent refresh yields1 winner/11 safe denials and no denial cookie. Browser confirms JavaScript cannot read refresh; local HTTP cookie HttpOnly/Lax/auth path (Secure:false intentionally local only), memory-only access/no auth storage and nonsecret messages. Captured structured logs contain safe metadata, no credentials. Secret scan matches zero generated environment/fixture secret values in changed artifacts. CORS/rates/strict JWT/HTTP/logging existing tests remain passing; no production HTTPS/TLS/proxy/deployment claim.

## 27. Migration tests

PASS15 real hardened-runner subtests under race: locks, atomic failure, ordering/checksum/pairs/tracking/fault recovery/least privilege/legacy attestation. Historical product-role test now explicitly uses its original four-pair prefix and refuses newer tracking; its old assertions remain. All8 old SQL files byte-identical; new005 grants no history mutation and refuses lossy down once any version/receipt exists. Normal local runner applied005 successfully without seeds; verified0 users/refresh/terms/acceptances, null pointer and five correct tracked versions. Base/isolated Compose configs pass. This task did not stop user normal8080/5173 processes. Final inspection finds both normal app ports not listening; run `make dev` to start the updated full stack. Normal PostgreSQL remains healthy. Disposable Phase4B resources are separately disposed after verification; normal/Phase4A/unrelated volumes/containers are preserved.

## 28. Remaining policy dependencies

**Blocked:** official exact V2 FSMO text, first version identifier, responsible Admin publication/material-update review (OPEN-016); chosen secure first-access delivery/verification/lifetime/reissue/first-password/recovery ownership (OPEN-001/009). Owner explicitly authorizes independent foundation only and defers these decisions. No policy was inferred from V1/code/mockups. Recommended alternatives are in DEC-067. No other Phase2.5 operational decision is changed, including entirely external return photographs.

## 29. Phase5/7 integration requirements

Phase5 needs explicit authorization and approved provisioning/delivery policy before dependent account activation/import. Preserve current roles/status/hash-only credentials, narrow Staff creation and Admin privileged operations, no public signup or password spreadsheet, retained publisher/acceptance identity. Do not hard-delete users referenced by history. Phase7 submission/direct issue: authorize actor/target, sorted participating-account locks, call policy in the **same transaction**, retain publication lock through stock/borrowing/history/receipt commit, bind returned `(receipt_id,borrower_id)` composite FK plus exact terms version. Direct issue checks target Borrower evidence; Staff cannot accept for them. Existing pending submissions keep bound terms. No business endpoints/tables/UI, live inventory arithmetic, fine/payment or directory CRUD/import is implemented here.

## 30. Exact files changed

The following snapshot includes modified tracked files and new source/document/evidence files; generated secrets/configs/node_modules/dist/cache are ignored and excluded. No staged files.

```text
README.md
backend/internal/application/terms/service.go
backend/internal/application/terms/service_test.go
backend/internal/bootstrap/dependencies.go
backend/internal/domain/terms/terms.go
backend/internal/domain/terms/terms_test.go
backend/internal/domain/user/permissions.go
backend/internal/domain/user/permissions_test.go
backend/internal/infrastructure/database/migrator_integration_test.go
backend/internal/infrastructure/database/product_roles_integration_test.go
backend/internal/infrastructure/persistence/postgres/terms_repository.go
backend/internal/interface/http/handlers/terms_handler.go
backend/internal/interface/http/middleware/request_safety.go
backend/internal/interface/http/response/errors.go
backend/internal/interface/http/response/errors_test.go
backend/internal/interface/http/routes/routes.go
backend/internal/interface/http/routes/terms_routes.go
backend/migrations/000005_terms_acceptance.down.sql
backend/migrations/000005_terms_acceptance.up.sql
backend/tests/integration/foundation_test.go
backend/tests/integration/terms_fixtures_test.go
backend/tests/integration/terms_http_test.go
backend/tests/integration/terms_test.go
docs/API_CONTRACTS.md
docs/domain/API_RESOURCE_DRAFT.md
docs/domain/BUSINESS_RULES.md
docs/domain/DATA_MODEL.md
docs/domain/DOMAIN_MODEL.md
docs/domain/INVARIANTS.md
docs/domain/STATE_MACHINES.md
docs/project/DECISIONS.md
docs/project/OPEN_DECISIONS.md
docs/project/PHASE4B_REPORT.md
docs/project/ROADMAP.md
docs/project/SOURCE_OF_TRUTH.md
docs/ux/verification/phase4b/README.md
docs/ux/verification/phase4b/RESULTS.json
docs/ux/verification/phase4b/VALIDATION.json
docs/ux/verification/phase4b/admin-1280-dark.png
docs/ux/verification/phase4b/admin-1280-light.png
docs/ux/verification/phase4b/preview-regression/B01-1280-dark.png
docs/ux/verification/phase4b/preview-regression/B01-1280-light.png
docs/ux/verification/phase4b/preview-regression/B01-320-dark.png
docs/ux/verification/phase4b/preview-regression/B01-320-light.png
docs/ux/verification/phase4b/preview-regression/B01-390-dark.png
docs/ux/verification/phase4b/preview-regression/B01-390-light-full.png
docs/ux/verification/phase4b/preview-regression/B01-390-light.png
docs/ux/verification/phase4b/preview-regression/B01-768-dark.png
docs/ux/verification/phase4b/preview-regression/B01-768-light.png
docs/ux/verification/phase4b/preview-regression/B02-1280-dark.png
docs/ux/verification/phase4b/preview-regression/B02-1280-light.png
docs/ux/verification/phase4b/preview-regression/B02-320-dark.png
docs/ux/verification/phase4b/preview-regression/B02-320-light.png
docs/ux/verification/phase4b/preview-regression/B02-390-dark.png
docs/ux/verification/phase4b/preview-regression/B02-390-light-full.png
docs/ux/verification/phase4b/preview-regression/B02-390-light.png
docs/ux/verification/phase4b/preview-regression/B02-768-dark.png
docs/ux/verification/phase4b/preview-regression/B02-768-light.png
docs/ux/verification/phase4b/preview-regression/B02-selection-dialog.png
docs/ux/verification/phase4b/preview-regression/RESULTS.json
docs/ux/verification/phase4b/preview-regression/S01-dashboard-1024-dark.png
docs/ux/verification/phase4b/preview-regression/S01-dashboard-1024-light.png
docs/ux/verification/phase4b/preview-regression/S01-dashboard-1440-dark.png
docs/ux/verification/phase4b/preview-regression/S01-dashboard-1440-light.png
docs/ux/verification/phase4b/preview-regression/S01-navigation-sheet.png
docs/ux/verification/phase4b/preview-regression/S01-pending-1024-dark.png
docs/ux/verification/phase4b/preview-regression/S01-pending-1024-light.png
docs/ux/verification/phase4b/preview-regression/S01-pending-1440-dark.png
docs/ux/verification/phase4b/preview-regression/S01-pending-1440-light.png
docs/ux/verification/phase4b/staff-1280-dark.png
docs/ux/verification/phase4b/staff-1280-light.png
docs/ux/verification/phase4b/terms-accept-network-error-390-dark.png
docs/ux/verification/phase4b/terms-accepted-390-dark.png
docs/ux/verification/phase4b/terms-accepted-390-light.png
docs/ux/verification/phase4b/terms-read-network-error-390-light.png
docs/ux/verification/phase4b/terms-reading-end-390-dark.png
docs/ux/verification/phase4b/terms-reading-end-390-light.png
docs/ux/verification/phase4b/terms-required-1280-dark.png
docs/ux/verification/phase4b/terms-required-1280-light.png
docs/ux/verification/phase4b/terms-required-320-dark.png
docs/ux/verification/phase4b/terms-required-320-light.png
docs/ux/verification/phase4b/terms-required-390-dark.png
docs/ux/verification/phase4b/terms-required-390-light.png
docs/ux/verification/phase4b/terms-required-768-dark.png
docs/ux/verification/phase4b/terms-required-768-light.png
docs/ux/verification/phase4b/terms-stale-390-dark.png
docs/ux/verification/phase4b/terms-stale-390-light.png
docs/ux/verification/phase4b/terms-unpublished-320-dark.png
docs/ux/verification/phase4b/terms-unpublished-320-light.png
docs/ux/verification/phase4b/terms-unpublished-390-dark.png
docs/ux/verification/phase4b/terms-unpublished-390-light.png
docs/ux/verification/phase4b/terms-updated-390-dark.png
docs/ux/verification/phase4b/terms-updated-390-light.png
frontend/scripts/phase4b-browser-qa.mjs
frontend/src/app/query-client.ts
frontend/src/app/router.tsx
frontend/src/components/application/borrower-navigation.ts
frontend/src/features/auth/navigation.ts
frontend/src/features/auth/product-auth.test.tsx
frontend/src/features/terms/api/terms.api.ts
frontend/src/features/terms/components/account-terms.tsx
frontend/src/features/terms/components/terms-frame.tsx
frontend/src/features/terms/components/terms-gate.tsx
frontend/src/features/terms/hooks/use-terms.ts
frontend/src/features/terms/pages/terms-page.tsx
frontend/src/features/terms/terms.test.tsx
frontend/src/features/terms/types.ts
frontend/src/features/workspace/workspace-page.tsx
frontend/src/lib/api-error.test.ts
frontend/src/lib/api-error.ts
frontend/src/lib/query-keys.ts
frontend/src/styles/index.css
frontend/src/styles/terms.css
integration/PHASE4B.md
integration/compose.phase4b.yml
integration/phase4b-env-run.py
integration/phase4b-fixture-state.py
```

## 31. git diff --check

PASS exit0, final whitespace check. No commit/stage/push/remote change/deployment. Source hashes verify8 old SQL,47 approved-package files,62 ui primitives,5 manifest/lock/config files and85 historical report/evidence files unchanged. Baseline609 tracked files. [Machine validation](../ux/verification/phase4b/VALIDATION.json) records exact counts, normal DB state and secret-scan result.

## 32. Exact git status

`git status --short`, captured after this report existed (including collapsed new directories exactly as Git reports):

```text
 M README.md
 M backend/internal/bootstrap/dependencies.go
 M backend/internal/domain/user/permissions.go
 M backend/internal/domain/user/permissions_test.go
 M backend/internal/infrastructure/database/migrator_integration_test.go
 M backend/internal/infrastructure/database/product_roles_integration_test.go
 M backend/internal/interface/http/middleware/request_safety.go
 M backend/internal/interface/http/response/errors.go
 M backend/internal/interface/http/response/errors_test.go
 M backend/internal/interface/http/routes/routes.go
 M backend/tests/integration/foundation_test.go
 M docs/API_CONTRACTS.md
 M docs/domain/API_RESOURCE_DRAFT.md
 M docs/domain/BUSINESS_RULES.md
 M docs/domain/DATA_MODEL.md
 M docs/domain/DOMAIN_MODEL.md
 M docs/domain/INVARIANTS.md
 M docs/domain/STATE_MACHINES.md
 M docs/project/DECISIONS.md
 M docs/project/OPEN_DECISIONS.md
 M docs/project/ROADMAP.md
 M docs/project/SOURCE_OF_TRUTH.md
 M frontend/src/app/query-client.ts
 M frontend/src/app/router.tsx
 M frontend/src/features/auth/navigation.ts
 M frontend/src/features/auth/product-auth.test.tsx
 M frontend/src/features/workspace/workspace-page.tsx
 M frontend/src/lib/api-error.test.ts
 M frontend/src/lib/api-error.ts
 M frontend/src/lib/query-keys.ts
 M frontend/src/styles/index.css
?? backend/internal/application/terms/
?? backend/internal/domain/terms/
?? backend/internal/infrastructure/persistence/postgres/terms_repository.go
?? backend/internal/interface/http/handlers/terms_handler.go
?? backend/internal/interface/http/routes/terms_routes.go
?? backend/migrations/000005_terms_acceptance.down.sql
?? backend/migrations/000005_terms_acceptance.up.sql
?? backend/tests/integration/terms_fixtures_test.go
?? backend/tests/integration/terms_http_test.go
?? backend/tests/integration/terms_test.go
?? docs/project/PHASE4B_REPORT.md
?? docs/ux/verification/phase4b/
?? frontend/scripts/phase4b-browser-qa.mjs
?? frontend/src/components/application/borrower-navigation.ts
?? frontend/src/features/terms/
?? frontend/src/styles/terms.css
?? integration/PHASE4B.md
?? integration/compose.phase4b.yml
?? integration/phase4b-env-run.py
?? integration/phase4b-fixture-state.py
```

## 33. Phase4B completion status

**IMPLEMENTED FOUNDATION: VERIFIED.** Versioning/publication mechanics, durable user-bound acceptance, first-use/updated/missing UI, transactional policy port and security/PG/HTTP/browser/visual gates pass. **BLOCKED FUNCTIONALITY:** approved official document publication/institutional consent and secure initial-access activation/credential delivery/recovery. **EXACT REQUIRED DECISION:** approve text/version/publication authority and select delivery/verification/expiry/reissue/first-password/recovery rules. **SAFE CURRENT BEHAVIOR:** no normal publication/acceptance; missing terms never unlock borrowing eligibility; permitted read-only/account/existing obligations/logout continue. Do not declare full Phase4B COMPLETE or blocked activation/content complete.

## 34. Phase5 readiness

NOT STARTED; no automatic next-phase authorization. Technical terms/role/session foundation is available. Account-management planning may use it when separately authorized, but dependent provisioning/activation requires the owner/FSMO/IT credential decision; official content remains required before institutional consent/new borrowing. No production launch or full onboarding readiness is claimed.
