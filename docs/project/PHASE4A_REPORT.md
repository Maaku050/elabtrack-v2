# Phase 4A — Authentication, Roles & Protected Navigation

Verified 2026-10-08. **COMPLETE within the authorized local implementation and verification scope.** Phase 4B and Phase 5 are NOT STARTED. No commit, push or deployment.

## 1. Executive summary

Real login now connects the approved FSMO visual system to the existing hardened authentication architecture. BORROWER, STAFF and ADMIN are current PostgreSQL account roles; protected navigation confirms current account authority before displaying a workspace. Logout clears memory/private queries immediately, coordinates other tabs and preserves honest pending/failure/retry feedback. Product workspaces contain explicit feature placeholders and safe own-account metadata, never preview borrowing records or fabricated operational reports.

Paired migration 000004 safely maps legacy roles and preserves identity/session data. Public registration is unavailable in every environment. Actual PostgreSQL migration, privilege and authentication regressions, real Chromium acceptance, original preview regression and quality gates pass. [Safe machine-readable summary](../ux/verification/phase4a/VALIDATION.json), [rendered evidence](../ux/verification/phase4a/README.md), [reproduction](../../integration/PHASE4A.md).

## 2. Sources reviewed

The explicit Phase 4A owner request governs this implementation. Repository instructions, README, Makefile, PROJECT_CHARTER, SOURCE_OF_TRUTH, DECISIONS, OPEN_DECISIONS, ROADMAP, STACK and current Go/frontend manifests establish scope and architecture. PHASE1_FOUNDATION, PHASE2_REPORT, PHASE2_5_REPORT, PHASE3B_REPORT and LOCAL_ENVIRONMENT_READINESS_REPORT supply historical engineering/product/visual evidence. Current domain documents establish confirmed product policy and future invariants; UX architecture, flows, wireframes, visual/status/responsive/composition rules, checklists, implementation map, fidelity contract and approved manifest establish the visual baseline.

Actual backend bootstrap/configuration, handlers/routes/middleware, account resolver/repository, JWT and refresh services, migrations/runner/privileges and their tests were inspected. Frontend router, provider/bootstrap, feature API/hooks, centralized transport, coordinator, state/query clearing, shells/tokens and tests were inspected. Approved B01/S01 core and B04-02/B06-02 images were opened for comparison. Current source and real verification establish implementation; historical V1/capstone/boss observations remain referenced evidence, not newly inspected live systems. No boss-repository or V1 Firebase access occurred.

## 3. Existing authentication architecture

The Go API remains Fiber v3 → application auth use cases/account ports → PostgreSQL repositories and existing JWT/session adapters. Domain/application do not acquire Fiber or pgx dependencies. The React boundary remains page → feature hook → TanStack Query → feature API → the single centralized transport. Zustand holds memory session/UI state, React local password-visibility state, and React Hook Form/Zod form validation.

The existing transport still establishes sessions from its login/refresh replies, attaches access tokens from memory, coordinates cookie mutations, fences stale responses and bounds refresh recovery. Lazy feature imports reference that same singleton; they do not create another client/provider. Query owns server account reads; UI role metadata never authorizes an API operation.

## 4. Final Borrower/Staff/Admin role model

| Current database role | Implemented workspace intent | Existing privileged API access |
|---|---|---|
| BORROWER | Own borrower shell | Own account only; directory denied |
| STAFF | FSMO operational shell | Own account only; directory denied |
| ADMIN | Staff operations plus Reports/Administration destinations | Own account and retained Admin account directory |

Student/Faculty remain borrower categories, not permissions. No category schema, tenancy, shared privileged login or Super Admin is introduced. Future inventory/borrowing/provisioning permissions remain the confirmed domain contracts for their authorized phases.

## 5. Foundation-role migration approach

Legacy `user` maps only to BORROWER; legacy `admin` maps to ADMIN. The former ordinary account receives no Staff/Admin privilege; legitimate Admin authority remains represented. Existing JWT role claims do not authorize requests, so a valid pre-migration token continues resolving identity against the newly mapped database role. Role validity now accepts exactly the three uppercase product values and fails closed for legacy/category/unknown strings.

All non-role account fields, uniqueness, UUIDs, password hashes, status/timestamps and refresh records are retained. Real migration tests compare before/after account and session snapshots excluding only the expected role change. Down maps BORROWER/ADMIN back to user/admin but refuses atomically if any STAFF exists; an operator must explicitly reconcile that unrepresentable role before rollback. No silent demotion or promotion occurs.

## 6. Backend authorization changes

The central domain permission matrix defines borrower workspace, Staff workspace, Administration and account-directory access. Reusable `RequirePermission` consumes the principal already resolved by the existing `Auth` middleware. That middleware verifies JWT identity then queries current account status/role. The retained `GET /api/v1/users/` directory now explicitly requires the Admin directory permission. Self-account endpoints retain their verified-identity boundary; profile updates accept the new valid role vocabulary without allowing role/status editing.

Real HTTP/PG checks deny Borrower and Staff directory access, allow Admin, deny an old Admin JWT after database demotion, deny disabled/missing accounts and deny inactive restoration. Route menus are presentation. There are no new business commands/API resources behind the placeholder workspaces; future commands must apply their own reviewed server permission and integrity rules.

## 7. Database migration changes

New paired `000004_product_roles.up.sql` / `.down.sql` changes only `users.role` values, its default and CHECK constraint. It uses the existing transactional/checksummed/advisory-locked runner; no table, extension, grant or migration-framework replacement is added. Historical 000001–000003 SQL pairs remain byte-identical. The development seed now creates no accounts or published passwords.

The isolated PostgreSQL 18.6 database verified legacy up/down/up, exact identity/session preservation, invalid-role rejection, Staff rollback refusal and runtime privilege denials. Existing runner regression was adapted for a fourth suffix while preserving the three-file explicit legacy-attestation boundary; 15 real fault/concurrency/privilege subtests pass under `-race`.

The normal preserved local database received 000004 through explicit `make migrate-up`; `make migrate-status` verified all four pairs. It had zero accounts/sessions before and after. This normal application does not substitute for the synthetic legacy-data preservation test. Normal startup still performs no migrations/seeding. Coordinate schema/application releases and rollback; do not run incompatible legacy application code against new role strings.

## 8. Login UI

`/login` has the compact FSMO seal/wordmark, navy/violet identity, existing theme switch, Welcome back heading and restrained single-column sign-in card. It uses existing application/shadcn compositions and Lucide icons. Email/password are labeled, with browser autocomplete, accessible show/hide state, inline associated validation, pending submit prevention and safe feedback. Password is cleared and focus restored after failed authentication.

Email trimming/format/255-character limits and password 8–72 limits match the existing server login contract. Both themes and 320/390/1280 widths pass real rendering checks. Provisioning guidance says contact FSMO; no signup, role picker or nonfunctional password-recovery link exists. Recovery is an explicit later dependency.

## 9. Login API integration

The feature adapter reuses actual `POST /api/v1/auth/login`, `/refresh`, `/logout` and `GET /api/v1/auth/me`. Existing self-profile endpoints are retained but no profile-edit module is added. The login adapter exposes safe returned account metadata; the transport alone establishes memory access/cookie-based refresh state. Successful role navigation is followed by a protected current-account read, so a concurrent role/status change cannot rely solely on the earlier login response.

`POST /api/v1/auth/register` is unmounted for development/test/production and returns the normal safe 404. Internal registration service/handler definitions remain unreachable rather than being repurposed as account management. Tests prove registration cannot create accounts/sessions or set cookies, including privilege-injection input. Older registration-based browser harnesses are historical; the current isolated fixture guide replaces them without re-enabling signup.

## 10. Protected route model

| Role | Exact permitted application paths |
|---|---|
| BORROWER | `/borrower/home`, `/borrower/equipment`, `/borrower/borrowings`, `/borrower/account`, `/borrower/notifications` |
| STAFF, ADMIN | `/staff/dashboard`, `/staff/requests`, `/staff/inventory`, `/staff/borrowers`, `/staff/account`, `/staff/notifications` |
| ADMIN only | `/admin/reports`, `/admin/administration` |

Root selects the authenticated role home or login. Anonymous direct links go to login with an internal pathname; login restores only an exact path permitted by that current role. External/protocol/traversal/preview/query destinations are not accepted. Borrower→Staff and Staff→Admin display accessible forbidden feedback. Admin inherits Staff access; no Borrower role inheritance is invented.

Before rendering protected children, `ProtectedRoute` requests current account per pathname, with no automatic account-read retry, fresh mount/focus/reconnect checks and no protected content while pending/refetching. Network failure gives a safe manual retry. Identity, active status and current role must agree with the current session. Router lazy-load failures use a generic reload experience instead of default stack output.

## 11. Borrower shell integration

The approved compact brand/header and Home/Equipment/My Borrowings/Account bottom navigation remain. Notifications are secondary; theme and visible Sign Out are accessible. Current own account name/email/role are read from the real session/account service. Borrower Home, Catalog, My Borrowings and Notifications display honest feature placeholders pending their services; no preview metric, equipment inventory or synthetic borrowing is represented as live data.

Account is read-only safe metadata. It does not fabricate category, program/contact, terms acceptance, preferences or profile editing. The illustrative B04-02 surface is consequently only partially implemented.

## 12. Staff/Admin shell integration

The navy sidebar and approved header/workspace structure are reused. Staff has Dashboard, Requests & Borrowings, Inventory and Borrowers. Admin additionally sees Reports and Administration. Current name/role replaces preview identity; navigation uses real guarded destinations. The placeholder workspace does not pretend to offer live dashboard metrics/search/activity or management actions.

Sign Out remains visible and readable in light/dark themes. Collapsed sidebar preserves a named 44px sign-out icon; only its text label is visually hidden. These production additions use scoped styles and optional shell props; default preview rendering remains intact.

## 13. Inactive account handling

Inactive login uses the same generic credential/account-unavailable feedback as invalid or missing credentials, avoiding enumeration. Current-account denial after issuance invalidates memory/private queries and denies protected routes. Refresh cannot restore an inactive account. Missing-account and revoked-session cases also deny restoration/access. No client workaround, eligibility suspension policy, deactivation UI or hard delete is added.

Current account authority is checked on API requests and route mount/focus/reconnect. This is not a new push mechanism for immediate cross-device revocation; institutional session/global-revocation policy remains OPEN-027.

## 14. Session bootstrap

The existing idle/bootstrapping/authenticated/unauthenticated/error semantics and coordinated single-flight restore remain. Unknown/restoring state displays loading without protected children. Expected anonymous refresh denial becomes ordinary login without a misleading persistent error. Session expiration gives understandable sign-in feedback. Transient network failure is distinct from bad credentials and offers bounded deliberate Retry session; no automatic refresh storm or recursive retry is introduced.

Signing-in UI state keeps the login form mounted during its own pending request. Current-account loading independently hides workspace content. Real held-refresh/reload checks prove protected content does not flash before authority resolves.

## 15. Logout behavior

Sign Out synchronously clears local access/user/private queries and removes protected UI, then invokes the existing transport logout and navigates to login. The transport preserves its peer signal, Web Lock sequencing and late-response fences. An in-memory UI action store carries pending/failure feedback across route unmount. Duplicate sign-out is fenced; a network failure honestly states local sign-out succeeded while server revocation could not be confirmed, with Retry sign out.

Successful logout revokes the presented refresh session and clears the scoped cookie with existing matching flags. Back navigation and reload cannot expose the previous workspace. Failed network logout cannot guarantee server-cookie revocation until retry; the UI explicitly communicates that distinction. Public query data survives clearing; all private query/mutation state is removed. Same-account role changes remove old private data and advance generation while preserving the resolving `auth.me` query to avoid canceling its own authority check.

## 16. Cross-tab regression results

Real native Web Locks/BroadcastChannel checks create a peer from the actual shared refresh cookie, obtain that peer's own access response and verify logout clears both tabs' memory/private UI. Observed lifecycle messages contain only the existing non-secret version/tab/attempt/epoch/time/type fields. Access/refresh credentials and account metadata are not broadcast or persisted.

Existing coordinator and transport tests still cover serialization, owner/waiter failure, stale events, late-response fencing, fallback, logout/invalidation, bounded retries and private-cache removal. The full historical Phase 1I three-tab pressure suite was not re-run here; its source/evidence is retained. Phase 4A supplies new live two-tab product-form acceptance plus the preserved unit and real PostgreSQL concurrency regressions.

## 17. Light/dark visual results

Login, Borrower/Staff/Admin navigation, forbidden and session/logout states use the existing semantic theme infrastructure. No evergreen identity or new palette is introduced. Captured layouts have zero horizontal overflow and no target under 44px. The real preview regression separately passes 24 theme/viewport cases, with eight interaction/accessibility checks.

Representative final screenshots were opened for visual comparison. Brand, typography hierarchy, card/focus styling, violet action and navy Staff navigation match the established system. Screenshot review found the new light-theme Staff Sign Out contrast issue; scoped styling corrected it and a regression checks readable expanded/collapsed states.

## 18. Responsive screenshots

[Evidence index](../ux/verification/phase4a/README.md) links all **24** PNGs. Login and protected Borrower catalog placeholder each have 320/390/1280 light/dark captures. Staff has 1024/1440 light/dark; Admin has 1440 light/dark. Additional captures show 390 own account, 390 forbidden in both themes, 1024 Staff forbidden, 390 dark failed logout and 390 light failed restoration.

Representative examples: [320 Borrower](../ux/verification/phase4a/borrower-320-light.png), [390 login](../ux/verification/phase4a/login-390-dark.png), [1280 login](../ux/verification/phase4a/login-1280-light.png), [1024 Staff](../ux/verification/phase4a/staff-1024-light.png), [1440 Admin](../ux/verification/phase4a/admin-1440-dark.png). These show real isolated accounts and honest placeholders, not the anonymous synthetic preview dashboards.

## 19. Authentication security regression

Existing JWT signature/algorithm/issuer/audience/purpose/time validation, memory-only tokens, HttpOnly host-only refresh cookies, hash-only storage, transactional rotation, replay denial and safe losing-response cookie behavior remain. Current DB status/role overrides old JWT metadata. Origin/CORS, limits, request IDs, safe errors/structured logging, runtime-only DML and migration controls remain unchanged in purpose and pass retained tests.

Real runtime refresh pressure yields one winner/eleven safe denials with no losing Set-Cookie, a coherent usable successor, replay denial, logout/idempotence, post-consumption transaction rollback and bounded cleanup. Local/browser access is never written into localStorage/sessionStorage/IndexedDB; the refresh cookie is invisible to JavaScript. Local HTTP Secure=false is explicit development/test policy; production Secure=true is covered by existing tests. Production HTTPS was not deployed/tested here.

## 20. PostgreSQL tests

`TestRealProductRoleMigration` passed against isolated PostgreSQL 18.6 with distinct non-superuser migrator/runtime roles. It verifies two synthetic legacy identities/sessions survive exact mapping/reversal/reapplication, invalid role constraints, Staff-safe rollback refusal, named-table runtime DML and denied runtime DDL/history/role creation.

`TestRealMigrator` passed under `-race` with **15 subtests**, covering atomic up/down SQL/bookkeeping failures, checksum tampering, missing/gapped/duplicate history, permission failure, unsupported/transaction-control SQL, advisory/concurrent connections, runtime/migrator privilege separation and explicit three-file legacy adoption.

`TestRealFoundation` passed against the real API/PG under `-race` with **eight subtests**: health/safe account, generic login, role matrix, stale authority, hash/rotation/concurrent HTTP/logout, transaction rollback, missing account after issuance and bounded retention. Borrower/Staff/Admin/inactive/disabled/missing/demoted/revoked behavior is exercised through actual repositories/HTTP, not mocked account responses.

Both normal and disposable databases finish at four verified migrations with **zero users and zero refresh rows**. Runtime is not superuser and cannot create in public or SELECT migration bookkeeping. Browser fixtures were four explicit random-password bcrypt accounts in the isolated database; only owned UUID/email records were mutated/deleted. The private mode-0600 fixture file was removed. No normal data reset occurred.

## 21. Backend test results

PASS `go fmt ./...`, `go vet ./...`, `go test ./...` using selected Go 1.27.1. The global outside-module launcher remains 1.26.5; no manifest/toolchain requirement was lowered. Normal opt-in live tests skip without explicit isolation; they were separately run as described above rather than counted as offline evidence.

PASS targeted `go test -race` across nine packages: application/auth, infrastructure/security, database, persistence/postgres, HTTP routes, middleware, bootstrap, config and response. Real migrator/foundation race runs are additional. Existing authority/CSRF/JWT/cookie/assertion coverage is preserved; role fixtures and migration counts were updated to the actual product schema.

## 22. Frontend test results

PASS **172 tests in ten files**. New meaningful product tests cover login validation/no premature submission, pending duplicate prevention, safe auth/network/rate/server-ID feedback and password focus, role destinations/safe deep links, anonymous/role guard matrix, current-account loading/no flash/demotion/inactive/error recovery, logout immediate clearing/pending/failure retry and themes. Store tests verify role-change cache removal/generation fencing. Existing coordinator/transport/session/bootstrap/visual tests remain passing. Registration-adapter absence is explicitly tested.

Mocked UI tests are supplemented by the real PG/API/browser run; they are not the only acceptance evidence. No test/typechecking gate was disabled or relaxed to accept an insecure outcome.

## 23. Browser checks

PASS Chrome **140.0.7339.16**, actual localhost Vite/API plus PostgreSQL. [RESULTS.json](../ux/verification/phase4a/RESULTS.json): **29 integrated checks**, **24 screenshots**, zero application errors/warnings. Real forms cover role login/destinations, safe deep links, forbidden, own metadata, reload/restore/loading, inactive login, disable/demotion after issuance, old Admin JWT denial, revoked restoration, cross-tab logout, history/reload denial, storage/cookie hygiene, keyboard/focus and offline recovery.

Fault cases hold an actual refresh, use Chromium offline or block only API traffic; auth replies are never fabricated. The harness was corrected to wait for disabled login inputs to become enabled after pending logout retry; the production submission fence remained enforced. Earlier unsuccessful harness attempts are not counted as passing evidence.

PASS original development preview matrix: **24 cases/eight checks**. PASS quick production regression: **four cases/ten checks**, including root product login, preview URL 404 and production JavaScript exclusion of development fixture routes. Preview-only anonymous-refresh interception is clearly separated from real authentication evidence.

The normal `make dev` workflow was also started against preserved local PostgreSQL: actual `/api/v1/health` and `/ready` 200, anonymous refresh 401, registration 404, frontend 200. Ctrl+C terminated both groups and the API logged pool closure. Final ports 8080/5173/4173 are not listening; normal PG remains healthy on 5434. Only isolated PG was stopped, with its volume retained; unrelated resources remain intact.

## 24. Lint/build results

PASS frontend lint: **zero errors, 19 existing primitive warnings, zero new warnings**. PASS strict TypeScript/Vite production build with no new build warnings; login/workspace and feature transport are lazy-loaded. Main JS is 478.43kB before gzip. Node 24.19.0/npm 11.17.0 remain compatible with the unchanged manifest. No dependency, lockfile or shadcn primitive changes.

PASS quiet normal Compose default/dev/production and isolated Phase 4A validation. Actual PostgreSQL containers were used; no fresh API/frontend Docker image build or deployment is claimed in this phase. Historical container verification remains evidence at its original date.

## 25. Remaining visual deviations

The B06-02 generated mobile reference includes protected bottom navigation and a permanent illustrative expiry notice. Runtime login removes that navigation and shows notices only for real session conditions, preserving the approved hierarchy. Labels/control density are larger for accessible narrow screens; wide/dark adaptations use established tokens rather than new generated PNGs. Production placeholder content changes B01/S01 content density because services are deferred. Own account omits unavailable B04-02 category/contact/terms/edit fields.

The supplied reconstructed seal and exact generated-font fidelity remain existing dependencies. The original four approved previews are preserved and independently verified. New authentication screenshots were inspected, but no pixel-perfect equivalence or new owner visual approval is asserted. Material future visual changes retain the existing review requirement.

## 26. Open Phase 4B dependencies

Phase 4B is NOT STARTED and requires explicit authorization. Institutional terms content/publication responsibility remains OPEN-016; secure initial activation/password delivery/email ownership/recovery design remains OPEN-001/009. Confirmed policy is acceptance once per material terms version, evidence linked to account/version/time, current terms required for new borrowing/direct issue and continued access to existing obligations. Authentication success is not terms acceptance.

Integration belongs at immutable terms/version and unique acceptance persistence, own-account content/acceptance API plus transactional future request/direct-issue guards, with stale-version/concurrency/idempotency and ownership tests. No existing implemented terms enforcement was found to bypass; no terms table, acceptance claim, text, route or gate is added here. Engineering authentication is ready to support the separately authorized phase; content/security-dependent details must be settled before implementing their workflows.

## 27. Phase 5 dependencies

Account provisioning remains authorized FSMO personnel policy. Later Phase 5 must implement bounded borrower directory, narrow Staff/Admin Borrower creation, Admin bulk import/status/privileged-account management and justified category/profile fields. It needs secure onboarding/recovery, reviewed import/duplicate/idempotency behavior, named actors/history, inactive-obligation preservation and privileged recovery/last-admin safeguards. None is implemented by placeholder routes or internal test fixture scripts.

There is no default product account/password and the development seed creates none. The preserved normal DB is empty; actual operational accounts require the later authorized provisioning path. Random isolated fixtures are verification-only and were cleaned.

## 28. Risks/blockers

**No remaining blocker to the authorized Phase 4A exit gate.** Deferred terms/activation/recovery and management/business services are dependencies of later phases. Failed network logout clears local authority but server revocation needs retry; the UI states this clearly. Staff rollback requires explicit account disposition because the legacy schema cannot represent it. Deployment must coordinate application/schema versions.

Production HTTPS/TLS, release/backup/restore, CI/operations, full cross-browser/WCAG validation and new owner visual acceptance are not claimed complete. The existing 19 warnings and reconstructed seal/font limitations remain disclosed. Historical signup-based integration/browser scripts are no longer current executable acceptance paths; current reproduction uses the isolated guide. No production deployment readiness, global session revocation or broader product completion is inferred.

All 171 selected protected baseline files are unchanged, including historical SQL, approved artwork/manifest, previous verification evidence, primitive components, dependency manifests/locks, Makefile and ignore rules. Existing three local environment files remain ignored mode 0600; isolated credentials also remain ignored 0600. Credential-value scanning of changed/new artifacts found zero matches. Screenshots expose only synthetic display identity, never passwords/tokens/cookie values.

## 29. Exact files changed

46 modified tracked files and 48 new files (expanded directories; 24 rendered PNGs). All changes are unstaged.

```text
M README.md
M backend/README.md
M backend/internal/application/auth/refresh_test.go
M backend/internal/domain/user/account.go
M backend/internal/domain/user/entity.go
M backend/internal/infrastructure/database/migrator_integration_test.go
M backend/internal/infrastructure/database/migrator_test.go
M backend/internal/infrastructure/persistence/postgres/user_account_test.go
M backend/internal/infrastructure/persistence/postgres/user_repository.go
M backend/internal/infrastructure/security/jwt_test.go
M backend/internal/interface/http/middleware/auth.go
M backend/internal/interface/http/routes/auth_boundary_test.go
M backend/internal/interface/http/routes/auth_routes.go
M backend/internal/interface/http/routes/browser_session_test.go
M backend/internal/interface/http/routes/user_routes.go
M backend/seeds/development.sql
M backend/tests/integration/foundation_test.go
M backend/tests/unit/domain/user/user_test.go
M docs/API_CONTRACTS.md
M docs/domain/API_RESOURCE_DRAFT.md
M docs/domain/DATA_MODEL.md
M docs/domain/DOMAIN_MODEL.md
M docs/project/DECISIONS.md
M docs/project/ROADMAP.md
M docs/ux/APPROVED_VISUAL_IMPLEMENTATION_MAP.md
M docs/ux/VISUAL_FIDELITY_CONTRACT.md
M frontend/README.md
M frontend/scripts/phase3b-browser-qa.mjs
M frontend/src/app/query-client.ts
M frontend/src/app/router.tsx
M frontend/src/app/session-bootstrap.test.tsx
M frontend/src/app/session-bootstrap.tsx
M frontend/src/components/application/shells.tsx
M frontend/src/features/auth/api/auth.api.ts
M frontend/src/features/auth/auth-foundation.test.tsx
M frontend/src/features/auth/hooks/use-auth.ts
M frontend/src/features/auth/index.ts
M frontend/src/features/auth/schemas/login.schema.ts
M frontend/src/features/foundation/foundation.test.tsx
M frontend/src/lib/api-client.test.ts
M frontend/src/lib/api-client.ts
M frontend/src/stores/auth-store.test.ts
M frontend/src/stores/auth-store.ts
M frontend/src/styles/index.css
M frontend/src/types/common.ts
M integration/README.md
A backend/internal/domain/user/permissions.go
A backend/internal/domain/user/permissions_test.go
A backend/internal/infrastructure/database/product_roles_integration_test.go
A backend/migrations/000004_product_roles.down.sql
A backend/migrations/000004_product_roles.up.sql
A backend/tests/integration/product_fixtures_test.go
A docs/project/PHASE4A_REPORT.md
A docs/ux/verification/phase4a/PREVIEW_REGRESSION.json
A docs/ux/verification/phase4a/PRODUCTION_REGRESSION.json
A docs/ux/verification/phase4a/README.md
A docs/ux/verification/phase4a/RESULTS.json
A docs/ux/verification/phase4a/VALIDATION.json
A docs/ux/verification/phase4a/account-390-light.png
A docs/ux/verification/phase4a/admin-1440-dark.png
A docs/ux/verification/phase4a/admin-1440-light.png
A docs/ux/verification/phase4a/borrower-1280-dark.png
A docs/ux/verification/phase4a/borrower-1280-light.png
A docs/ux/verification/phase4a/borrower-320-dark.png
A docs/ux/verification/phase4a/borrower-320-light.png
A docs/ux/verification/phase4a/borrower-390-dark.png
A docs/ux/verification/phase4a/borrower-390-light.png
A docs/ux/verification/phase4a/forbidden-390-dark.png
A docs/ux/verification/phase4a/forbidden-390-light.png
A docs/ux/verification/phase4a/login-1280-dark.png
A docs/ux/verification/phase4a/login-1280-light.png
A docs/ux/verification/phase4a/login-320-dark.png
A docs/ux/verification/phase4a/login-320-light.png
A docs/ux/verification/phase4a/login-390-dark.png
A docs/ux/verification/phase4a/login-390-light.png
A docs/ux/verification/phase4a/logout-network-error-390-dark.png
A docs/ux/verification/phase4a/session-network-error-390-light.png
A docs/ux/verification/phase4a/staff-1024-dark.png
A docs/ux/verification/phase4a/staff-1024-light.png
A docs/ux/verification/phase4a/staff-1440-dark.png
A docs/ux/verification/phase4a/staff-1440-light.png
A docs/ux/verification/phase4a/staff-forbidden-1024-dark.png
A frontend/scripts/phase4a-browser-qa.mjs
A frontend/src/features/auth/components/access-boundary.tsx
A frontend/src/features/auth/navigation.ts
A frontend/src/features/auth/pages/login-page.tsx
A frontend/src/features/auth/product-auth.test.tsx
A frontend/src/features/workspace/workspace-page.tsx
A frontend/src/stores/session-action-store.ts
A frontend/src/styles/authentication.css
A integration/PHASE4A.md
A integration/compose.phase4a.yml
A integration/phase4a-env-run.py
A integration/phase4a-fixture-state.py
```

## 30. git diff --check

**PASS**, exit 0, including final documentation. No whitespace errors. Checks and evidence were completed without staging, committing, pushing, changing remotes or deploying.

## 31. Exact Git status

Exact `git status --short` at handoff (untracked directories use Git's default abbreviation):

```text
 M README.md
 M backend/README.md
 M backend/internal/application/auth/refresh_test.go
 M backend/internal/domain/user/account.go
 M backend/internal/domain/user/entity.go
 M backend/internal/infrastructure/database/migrator_integration_test.go
 M backend/internal/infrastructure/database/migrator_test.go
 M backend/internal/infrastructure/persistence/postgres/user_account_test.go
 M backend/internal/infrastructure/persistence/postgres/user_repository.go
 M backend/internal/infrastructure/security/jwt_test.go
 M backend/internal/interface/http/middleware/auth.go
 M backend/internal/interface/http/routes/auth_boundary_test.go
 M backend/internal/interface/http/routes/auth_routes.go
 M backend/internal/interface/http/routes/browser_session_test.go
 M backend/internal/interface/http/routes/user_routes.go
 M backend/seeds/development.sql
 M backend/tests/integration/foundation_test.go
 M backend/tests/unit/domain/user/user_test.go
 M docs/API_CONTRACTS.md
 M docs/domain/API_RESOURCE_DRAFT.md
 M docs/domain/DATA_MODEL.md
 M docs/domain/DOMAIN_MODEL.md
 M docs/project/DECISIONS.md
 M docs/project/ROADMAP.md
 M docs/ux/APPROVED_VISUAL_IMPLEMENTATION_MAP.md
 M docs/ux/VISUAL_FIDELITY_CONTRACT.md
 M frontend/README.md
 M frontend/scripts/phase3b-browser-qa.mjs
 M frontend/src/app/query-client.ts
 M frontend/src/app/router.tsx
 M frontend/src/app/session-bootstrap.test.tsx
 M frontend/src/app/session-bootstrap.tsx
 M frontend/src/components/application/shells.tsx
 M frontend/src/features/auth/api/auth.api.ts
 M frontend/src/features/auth/auth-foundation.test.tsx
 M frontend/src/features/auth/hooks/use-auth.ts
 M frontend/src/features/auth/index.ts
 M frontend/src/features/auth/schemas/login.schema.ts
 M frontend/src/features/foundation/foundation.test.tsx
 M frontend/src/lib/api-client.test.ts
 M frontend/src/lib/api-client.ts
 M frontend/src/stores/auth-store.test.ts
 M frontend/src/stores/auth-store.ts
 M frontend/src/styles/index.css
 M frontend/src/types/common.ts
 M integration/README.md
?? backend/internal/domain/user/permissions.go
?? backend/internal/domain/user/permissions_test.go
?? backend/internal/infrastructure/database/product_roles_integration_test.go
?? backend/migrations/000004_product_roles.down.sql
?? backend/migrations/000004_product_roles.up.sql
?? backend/tests/integration/product_fixtures_test.go
?? docs/project/PHASE4A_REPORT.md
?? docs/ux/verification/phase4a/
?? frontend/scripts/phase4a-browser-qa.mjs
?? frontend/src/features/auth/components/
?? frontend/src/features/auth/navigation.ts
?? frontend/src/features/auth/pages/
?? frontend/src/features/auth/product-auth.test.tsx
?? frontend/src/features/workspace/
?? frontend/src/stores/session-action-store.ts
?? frontend/src/styles/authentication.css
?? integration/PHASE4A.md
?? integration/compose.phase4a.yml
?? integration/phase4a-env-run.py
?? integration/phase4a-fixture-state.py
```

`git diff --cached --name-only` is empty. Branch `main`; origin remains `https://github.com/Maaku050/elabtrack-v2.git`. No commit/push/remote change.

## 32. Phase 4A exit gate result

**SATISFIED / COMPLETE, 2026-10-08.** Real login and current database role/status authority; safe tested migration; absent public signup; protected/forbidden/deep-link loading/error behavior; inactive denial; immediate local logout and real peer invalidation; restore/rotation/replay/cookie security; preserved approved shells/development-only previews; themes/responsive/rendered evidence; actual PostgreSQL/race/frontend/backend/Compose/diff gates; safe fixtures/cleanup and no secrets are verified within the stated scope.

No Phase 4B, Phase 5 or inventory/borrowing/return/replacement/fine/report service was implemented. Authentication engineering is ready for separately authorized Phase 4B planning/implementation subject to its content/onboarding dependencies. Phase 4B remains **NOT STARTED**. Exact worktree changes above are all unstaged; owner review and any commit/push/deployment remain separate actions.
