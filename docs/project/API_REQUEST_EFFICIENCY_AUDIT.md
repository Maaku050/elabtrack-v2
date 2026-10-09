# API request efficiency & navigation audit

Date: 2026-10-09. Owner review pending. No Phase 7 work, policy change, migration change, commit, push or deployment.

The confirmed redundant requests were repeated administrative lookups of unpublished terms. Protected navigation already makes exactly one authoritative current-account request per route and reuses cached equipment/borrower data. The original unmatched GET 500 could not be identified from the evidence available; it remains open. All 236 frontend tests pass with one worker, but the default parallel run retains an initial lazy-login loading timeout.

## Necessary requests and their exact triggers

- **`GET /auth/me`:** `ProtectedRoute` is the only production caller of `useCurrentUser`. Its query key includes `location.pathname`, with `staleTime: 0`, `retry: false`, `refetchOnMount: 'always'`, and focus/reconnect revalidation enabled. Each protected pathname change needs fresh server role/active-account verification. The guard withholds even previously cached privileged content while verification is pending or failing. Returning to an already visited pathname still revalidates it. Search/filter changes retain the pathname and do not themselves recreate this authority query. These checks were preserved.
- **Equipment, categories and borrowers:** feature hooks use actor/filter-scoped TanStack Query keys, a 30-second freshness window and authenticated-query metadata. Initial visits, changed filters, stale remounts, and successful feature mutations can legitimately fetch/refetch data. The measured three-cycle sequence fetched each list only once. Dashboard currently has no operational data endpoint to fetch in this sequence.
- **Staff directory and current terms:** Administration needs these informational reads on its initial visit. Published and unpublished administrative publication information now use the same bounded freshness window. This cached information cannot authorize borrowing or represent consent.
- **`OPTIONS 204`:** the local SPA and API are different origins. Authenticated requests carrying `Authorization`, and applicable JSON headers, trigger browser-managed CORS preflight. The existing allowlist and credential requirements remain intact; preflight responses already advertise a 300-second max age. Subsequent calls to warmed paths reused browser preflight permission during the measurement. An OPTIONS log with `route: unmatched` is expected because CORS handles it before endpoint matching; its status is 204.
- **Refresh:** document startup restores the HttpOnly-cookie session through the centralized transport. Actual 401 recovery, explicit refresh, deliberate document reload and manual recovery can also refresh. Ordinary SPA navigation performed no refresh.

Inspection covered `app/router.tsx`, `providers.tsx`, `session-bootstrap.tsx`, the persistent `OperationalLayout`, access boundaries, current-account/profile hooks, auth/session stores, transport/coordinator, query configuration, and account/inventory/terms hooks. `useProfile` is a separate `/users/me` adapter but has no production page callers. There is no second component independently fetching `/auth/me` in this navigation sequence.

`SessionBootstrap` runs its effect at provider mount; the transport shares its bootstrap promise across StrictMode effects. Router links navigate inside the SPA. The shell's header and main container retained their DOM identities through all measured route changes. The sidebar is also inside that persistent layout in the source. No document requests occurred in either measured sequence.

No navigation-triggered invalidation or effect loop was found. Profile updates intentionally invalidate current-account metadata; account/inventory mutations invalidate their feature roots. Role changes intentionally clear private queries while allowing the in-flight account verification to settle. Focus and reconnect each produced one current-account check in Chromium, with no continuing request loop.

## Confirmed redundancy and correction

Previously, `useCurrentTerms` rejected `503 TERMS_NOT_PUBLISHED` and left the query without successful data. `retry: false` disabled retry attempts within that fetch, but did not disable TanStack Query's default failed-query loading on remount. Leaving and reopening Administration therefore repeated the same expected unavailable-policy request.

The administrative query now maps **only** the centralized `ApiRequestError` with status 503 and code `TERMS_NOT_PUBLISHED` to successful `null` data. It caches that informational state for 30 seconds. Remount after expiry or explicit terms invalidation rechecks publication; no timer polls while the page stays mounted. Network failures, other 5xx responses, 401 and 403 remain failures, with automatic retries disabled. Actor-scoped private state clears on logout/account authority changes. The page retains its existing institutional-approval explanation, rather than displaying an internal-failure message.

Borrower terms status, consent submission, terms gates, backend acceptance/version checks and account authority were not relaxed or changed. An unpublished administrative result never means that terms were accepted or that borrowing is authorized.

The backend also previously classified every 5xx as unexpected, so an intentional unpublished-policy 503 generated both an ERROR completion and an ERROR failure record. This particular trusted mapping now emits one INFO completion retaining status 503, code and correlation ID. Genuine service failures, internal errors and internal errors joined with `ErrNotPublished` still produce ERROR records. HTTP envelopes and private `no-store` headers remain unchanged. Sanitized before/after log samples are in [checks.json](verification/request-efficiency/checks.json).

## Chromium before/after counts

Each measurement used a fresh real Chromium browser at 1440×900, the existing isolated HTTP server on port 18085, Vite on port 15175, and PostgreSQL `elabtrack_v2_batch1_test` on port 54832. Random synthetic accounts were used; the owner's Admin credentials were not used. Login/bootstrap traffic is excluded. All repeated navigation finished within the data freshness window.

Begin at Dashboard after login. Click Inventory → Borrowers → Dashboard three times. Then click Administration → Dashboard three times with no current published terms.

| Sequence / request | Before | After |
| --- | ---: | ---: |
| Nine Dashboard/Inventory/Borrowers transitions: `/auth/me` | 9 | 9 |
| Equipment list | 1 | 1 |
| Equipment categories | 1 | 1 |
| Borrower list | 1 | 1 |
| Preflights in that sequence | 3 | 3 |
| **That sequence total** | **15** | **15** |
| Six Administration/Dashboard transitions: `/auth/me` | 6 | 6 |
| Staff directory | 1 | 1 |
| `/terms/current` unavailable-policy responses | 3 | 1 |
| Preflights in that sequence | 2 | 2 |
| **That sequence total** | **12** | **10** |
| **Combined application calls excluding OPTIONS** | **22** | **20** |
| **Combined requests including OPTIONS** | **27** | **25** |
| Navigation document reloads / refresh calls | 0 / 0 | 0 / 0 |

Every measured preflight returned 204; data/current-account reads returned 200; the expected terms response remained 503. There were no 500 responses or browser runtime exceptions in the measured sequences. Counts describe this controlled sequence, not a claim that all real-world browsing should have identical counts: cold/warm preflight state, elapsed time, filters, focus and mutations affect legitimate traffic.

Evidence: [before.json](verification/request-efficiency/before.json), [after.json](verification/request-efficiency/after.json). The repeatable CDP script is [request-efficiency-qa.mjs](../../frontend/scripts/request-efficiency-qa.mjs). It reads the existing private synthetic fixture file from the isolated harness; it cannot be run against the normal local accounts. Its after mode changes only the guarded disposable fixture's role/active flag and restores it on exit. Follow the existing [isolated verification setup](../../integration/BATCH1.md) and `TestBatch1BrowserServer` harness when rerunning; measure a before baseline before applying the caching change. An unpublished fixture can be selected by clearing only that disposable database's singleton publication pointer, preserving immutable version/acceptance history.

## Unmatched GET 500: unresolved original incident

The owner supplied `unmatched`, GET, status 500 and `INTERNAL_ERROR`, but not the actual request path, request ID or correlated failure/recovery log. A clarification requested those details; none was available by closeout.

The completion logger deliberately stores route templates rather than raw URLs or query strings. `unmatched` is a logging fallback, not a request path, and does not establish that a normal 404 was incorrectly mapped to 500. A correlated failure or panic record and the browser's Network entry are needed to distinguish an actual unexpected error from unmatched routing.

Read-only probes against the existing normal API returned health 200 and **404 `NOT_FOUND`** for `/`, `/favicon.ico`, `/api/v1/request-audit-unmatched` and `/api/v1/auth/request-audit-unmatched`. Paths, statuses and request IDs are recorded in [unmatched-probes.json](verification/request-efficiency/unmatched-probes.json). The existing full-middleware test also passes its unmatched path, returned-error and panic cases: framework 404 maps to `NOT_FOUND`; real internal failures and panics map to `INTERNAL_ERROR`.

**The original failing path and root cause are not established. No speculative unmatched-routing change was made, and the incident is not declared fixed.** Preserve the failing browser Network path plus `X-Request-ID` and the corresponding backend completion/failure/panic records for follow-up. Do not include credentials, cookies or tokens.

## Security and verification

Real Chromium checks passed for login and intentional logout; live Admin→Staff denial on an already visited Admin route; private-policy cache clearance after role change; deactivation returning `/auth/me` 403 with protected content/private data removed; a deliberately rejected synthetic access credential causing exactly one 401, one successful refresh and one successful retry; explicit refresh; reload performing one successful refresh; a refresh network failure withholding protected content until manual recovery; logout clearing a second real tab; revoked refresh returning 401; and successful subsequent login/logout. Harmless disabled query placeholders can remain, but no private successful data remained after invalidation.

Existing real PostgreSQL/HTTP `TestRealFoundation` and `TestRealTermsHTTP` passed separately. These cover current role/status enforcement, stale authority, hashed refresh rotation/concurrent refresh/logout, absent accounts, transaction rollback, terms identity/origin boundaries, version conflict, and inactive-account consent denial. Synthetic publications existed only in the disposable test database; no institutional terms were invented or published to normal data. The full terms database suite requiring an empty publication history was not rerun against this populated browser fixture database.

| Check | Result |
| --- | --- |
| `go fmt ./...`, `go vet ./...`, `go test ./...` | Pass; gated database/browser cases skip in the ordinary full run |
| Focused full-middleware logging/unmatched-route tests | Pass |
| Separate existing real HTTP/PostgreSQL authentication and terms tests | Pass |
| `npm run lint` | Pass, 19 existing warnings |
| `npm run build` | Pass |
| `npm run test:run`, initial run while other checks ran | 234 pass / 2 fail: initial lazy-login loading timed out in unchanged auth/foundation tests |
| `npm run test:run`, default-worker rerun | 235 pass / 1 fail: initial lazy-login loading timeout in unchanged auth test |
| `npm run test:run -- --maxWorkers=1` | **236/236 pass**, 15 files, all assertions retained |
| Unchanged auth test file run independently | 31/31 pass |
| Real Chromium before/after and security checks | Pass |
| `git diff --check` | Pass |

The remaining parallel-test loading timeout is reported rather than hidden. No test assertions, timeouts, worker defaults or type checks were weakened. The new 13 publication-cache tests verify remount reuse, expiry/invalidation discovery, genuine failures, role restrictions and private-cache clearing; logging tests verify expected versus genuine failures.

The shell PATH selected Windows npm and omitted Linux Go. Checks used the already installed Linux Node 24.19.0/npm and Go 1.27.1 toolchain explicitly, with a writable `/tmp` Go cache. No dependency or toolchain policy was changed. Initial attempts using missing `/snap/bin/go`, Windows npm or the read-only default Go cache failed before the corrected commands; those failures were environment/tool selection issues.

No security regression was observed in the exercised paths. Refresh/session/route authority source was unchanged during this audit. Normal development data was not mutated. The owned test API/Vite listeners, label-verified disposable database/volume/network and private fixture/environment files were removed. The normal API and login page still returned 200 and the normal PostgreSQL container remained healthy. The owner's existing API process was left running; the logging change was exercised in the isolated rebuilt API and takes effect locally when the normal backend is rebuilt/restarted through the existing development workflow. Docker application image builds and production deployment were not performed.

## Exact audit changes and handoff

Changed existing files:

- `frontend/src/features/terms/hooks/use-terms.ts`
- `frontend/src/features/accounts/pages/administration-page.tsx`
- `backend/internal/interface/http/response/errors.go`
- `backend/internal/interface/http/middleware/logger.go`
- `backend/internal/bootstrap/contracts_test.go`

Added:

- `frontend/src/features/terms/current-terms.test.tsx`
- `frontend/scripts/request-efficiency-qa.mjs`
- This report and `docs/project/verification/request-efficiency/{before,after,checks,unmatched-probes}.json`.

SHA-256 comparison against the 1,016-file starting snapshot confirmed that only the five listed existing files changed. Earlier uncommitted UI/login work was preserved; the Git index stayed empty.

Owner review remains pending. Follow-up needs the original unmatched-500 path/correlation evidence and investigation of the default parallel test's initial lazy-route timing. No Phase 7, commit, push or deployment is authorized by this report.
