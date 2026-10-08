# Foundation API contracts and operational observability

**Batch1 update,2026-10-09:** Phase5 account management and Phase6 catalog/inventory now have implemented contracts below. Other future resource drafts remain unimplemented; no borrowing/return/replacement/fine settlement is included.


Phase 1F technical contracts with Phase 4A authentication and Phase 4B terms updates, 2026-10-08 (Asia/Shanghai). Current source and tests govern implementation; this document does not authorize later phases. Product roles are BORROWER, STAFF and ADMIN; terms infrastructure is implemented; Phase5 account provisioning is implemented under Batch1; borrowing remains deferred.

## Routes and success convention

All API routes use `/api/v1`. JSON successes retain `{success:true,message,data,meta}`. `meta` is null except the existing generic user list, whose existing `{page,per_page,total,last_page}` contract is retained. No pagination design is created for future modules. Logout and accepted CORS preflight return bodyless 204. Success request IDs are in `X-Request-ID`, without adding observability metadata to domain entities.

| Route | Success / visibility | Important failures |
|---|---|---|
| POST /auth/register | Unavailable in every environment | 404 after the existing perimeter; no account/session/cookie creation. |
| POST /auth/login | 200, safe browser session | 400 input, 401 generic invalid credentials, 403 Origin, 500 issuance/persistence failure |
| POST /auth/refresh | 200, rotated safe browser session; cookie only | 401 same invalid-session response for missing/malformed/unknown/expired/revoked/replayed/disallowed account; 403 Origin; 500 infrastructure |
| POST /auth/logout | 204, idempotent cookie revoke/clear | 403 Origin, 500 persistence failure; cookie still cleared |
| GET /auth/me | 200, safe current account | 401 missing/invalid bearer/account; 403 inactive/unknown role; 500 lookup failure |
| GET /users/me | 200, safe profile | Same current-account boundary |
| PATCH /users/me | 200, display name only | Same boundary; 400 strict input/fields; 403 repository access predicate; 500 storage |
| GET /users/ | 200, existing bounded list + existing meta | Same boundary plus current ADMIN permission accounts.read (BORROWER/STAFF denied); 400 invalid supplied query values; 500 storage |
| GET /health | 200, process liveness | General perimeter/errors; does not query DB |
| GET /ready | 200 when DB checker succeeds | 503 dependency absent/failing, safe standard error |
| Unknown route / wrong method | No successful API fallback | Standard 404 / Fiber's practical 405, no raw framework message |

Fiber retains automatic HEAD for GET. Cookie mutations remain POST-only and exact trusted-Origin guarded. Perimeter 403/413/415/429/431 may precede routing, including an absent registration route; this does not expose registration. Unknown API routes are distinct from frontend SPA 404. nginx errors generated before the API remain ingress responses, not an API-generated envelope; the frontend handles them generically. Phase 1G verifies the local nginx serving path; production edge behavior requires its own verification.

Safe role values are exactly `BORROWER`, `STAFF`, `ADMIN`; Student/Faculty are implemented BORROWER category attributes, not authorities. Old JWT role hints do not override current PostgreSQL state. Browser session data contains `access_token,expires_at,token_type,user`, with user `id,email,name,role,is_active`. Raw refresh travels only via HttpOnly cookie. No new credential transport, retry, family-revocation or cross-tab policy is implied.

## Failure envelope and mapping

```json
{
  "success": false,
  "message": "Validation failed.",
  "data": null,
  "meta": null,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed.",
    "requestId": "12345678-1234-4123-8123-123456789abc",
    "fields": {"email": "is required"}
  }
}
```

The top-level message mirrors error.message for retained envelope compatibility. Client code reads the typed nested error; no message parsing. Optional fields map JSON field names to one safe message per field; no submitted values or rejected payloads. No stack, SQL, connection string, pgx detail or JWT parser internals are public.

`response.Map` is the centralized HTTP adapter. It uses `errors.Is/As` for wrapped known errors and structured field errors. Domain/application remain independent of Fiber/pgx. Internal classification takes precedence over joined causes; unsupported framework statuses and arbitrary failures become generic 500. `response.Error` records safe diagnostic metadata even when a handler writes JSON directly. The global error handler and completion middleware cover returned errors.

| Status | Code | Current meaning |
|---|---|---|
| 400 | BAD_REQUEST | Malformed, unknown-field or multiple JSON bodies / framework syntax error |
| 400 | VALIDATION_ERROR | Missing/invalid fields, domain input validation or current list query values |
| 401 | UNAUTHORIZED | Authentication required; invalid bearer or missing current account |
| 401 | INVALID_CREDENTIALS | Generic login rejection, no account existence/status detail |
| 401 | TOKEN_INVALID | Generic invalid/expired refresh session, no reason differentiation |
| 403 | FORBIDDEN | Current-account authorization, role, Origin/CORS denial |
| 404 | NOT_FOUND | Unknown resource/route or absent registration |
| 405 | METHOD_NOT_ALLOWED | Fiber rejects a method for a retained route |
| 409 | CONFLICT | Real duplicate/current-state conflict; currently local duplicate email |
| 413 | PAYLOAD_TOO_LARGE | Existing global/auth body bounds |
| 415 | UNSUPPORTED_MEDIA_TYPE | Current JSON mutations require application/json |
| 429 | RATE_LIMITED | Existing effective-IP budget; Retry-After retained |
| 431 | BAD_REQUEST | Socket parser rejects oversized request headers; safe generic text |
| 500 | INTERNAL_ERROR | Unexpected/internal failure, generic safe message |
| 503 | SERVICE_UNAVAILABLE | Readiness checker unavailable/failing, no dependency error text |

Codes are uppercase underscore-separated, intentionally small and retained where meaningful. They are public semantics, never Go type names. New modules extend the catalog and frontend normalization together. 409 does not replace generic validation and creates no loan/inventory conflict behavior yet.

One input policy: **400 for malformed and semantic request validation**, including wrapped FieldErrors; 422 is normalized to 400. 413/415 describe transport bounds/media, not field validation. Current user-list page/per_page must be positive bounded integers; omitted values retain defaults, valid large per_page clamps to 100. Offset overflow, unsupported supplied sort/order, and search over 256 bytes produce structured fields. There are no retained dynamic-ID routes needing path binding; an unmatched path stays 404.

Expected login failures remain uniform. Phase 1E also hid credential-lookup infrastructure failure behind the same 401; Phase 1F preserves that containment while recording an ERROR diagnostic and completion with status 401. This narrowly marked lookup exception is deliberate; token issuance/session persistence and other unexpected failures remain 500. Timing equalization is not proven.

InternalFailure is a transport-independent safe carrier: code-owned operation and root cause type, generic Error(), and ErrInternal classification. Sensitive infrastructure text is discarded at auth/transaction/session storage boundaries; raw wrapped user/register errors can reach the mapper but never its public text or log value. Logs use type/operation metadata rather than err.Error(), zap.Error or arbitrary object formatting. Operators investigate using request ID, operation/class and authorized database/network diagnostics. Logs intentionally omit detailed database exception text, even internally.

## Correlation and logging

All incoming X-Request-ID values are ignored, including UUID-looking, oversized and malformed values. The server generates UUIDv4 per request, sets its response header, and stores it in a private-key standard context via `observability.WithRequestID`. It is correlation only, not identity, authorization or an idempotency key. Error JSON and logs share that ID. Application/repository ports receive context.Context; transaction context additions preserve correlation. Pre-handler error rendering also ensures an ID. CORS already exposes X-Request-ID/Retry-After; no inbound-ID trust or proxy correlation protocol is added.

Production defaults to existing Zap JSON; local console remains configurable, with the same structured fields. There is one logging framework and an injectable Zap core, tested with observer records. No Fiber plaintext request logger remains. Startup's separate Fiber banner is disabled. Final request duration is numeric milliseconds measured around middleware/handler/error rendering, not network delivery or distributed tracing.

Middleware: outer recovery fallback → request ID → existing client-IP resolver → completion logger → inner recovery → existing security headers → CORS → rate limits → request safety → compression → routes. Inner recovery allows the completion logger to observe panic 500; outer recovery contains rare correlation/IP/completion failures. Normal returned errors are rendered before completion is recorded. Direct error writers and returned errors receive one safe diagnostic each; panic receives its own diagnostic, without a duplicate request-error diagnostic.

| Field / event | Policy |
|---|---|
| request_id | Server-generated UUID, context/header/error/log correlation |
| method, route | Method and registered route template; unmatched fallback, never raw path/query/URL |
| client_ip | Phase 1E shared trusted-proxy/effective-IP signal; not identity |
| status, duration_ms | Final routed HTTP status (parser lifecycle limitation below); numeric elapsed milliseconds |
| error_code | Stable public category when available |
| error_class, operation | Safe cause type and code-owned operation, never raw error text |
| http.request_completed | Completion record; no body/header/account dump |
| http.request_failed | Unexpected failure diagnostic, including masked login lookup |
| http.panic_recovered | Panic type plus server source stack; no panic value |
| security_event | auth.login_succeeded/failed, auth.refresh_succeeded/failed, auth.logout_completed/failed when an auth route executed |
| server.starting | Environment, listen address, explicit migration policy, startup seed=false, trusted-proxy enabled flag, origin count |
| server.shutdown_started/completed/failed, server.listen_failed, database.init_failed | Existing lifecycle hooks, safe event/class metadata |

Account IDs/emails/names are omitted from logs because current diagnostics do not need them. An idempotent successful logout is logged as completed, not proof a particular session row existed. Early perimeter denial still has its public category/status; it does not fabricate an authenticated account or successful session event.

DEBUG: routine 4xx diagnostics, especially unknown routes/validation/method failures. INFO: successful requests, normal lifecycle, expected generic login/refresh 401 events. WARN: 403 access/Origin denial and 429 rate limiting. ERROR: unexpected/internal/dependency failure and panic; masked lookup failure keeps public 401 but logs ERROR. Routine wrong passwords are not ERROR. Levels are configurable; production default info suppresses routine 404 noise. Existing explicit debug can be used for controlled investigation; no production SLO/slow threshold is invented. High-volume health success sampling/retention and operational budget tuning require measured deployment needs.

Panic responses always use generic INTERNAL_ERROR/500—even panic(fiber.ErrUnauthorized). The diagnostic logs type/source stack and request ID, with no raw panic value. Recovery does not dump requests. Every infrastructure failure path remains safe across configured log formats.

Forbidden log values: passwords/hashes, JWT signing material, access/refresh credentials, refresh digests, Authorization/Cookie/Set-Cookie headers, DB passwords/full credential-bearing URLs, SMTP passwords/API keys, reset tokens and future kiosk secrets. Also avoid emails, account metadata, raw URLs/query values, arbitrary headers/bodies, config structs and unreviewed error/object formatting. Enforcement is field allowlisting and dropping raw detail, not a regex secret scanner. Tests use synthetic sentinels in request bodies/headers/query/response cookies/panic/error/startup settings. Future callers must retain these rules.

API parser failures precede normal handler execution. Fiber's installed server-error path traverses global middleware before applying its final parser status; App.Test reports global oversized-body parse errors rather than an HTTP response. The Phase 1G listener callback wrapper preserves that traversal but defers completion logging until the final parser response, retaining its already-generated ID after Fiber context release. Actual body/header socket limits return safe 413/431 and exactly one correctly correlated final-status completion. Slow-client timeout runtime injection was not performed; configured timeout evidence remains source/offline. Production ingress-generated errors still need chosen-edge verification.

## Health, client use and extension points

GET /health is liveness only and never checks DB: data `{status:"ok",service:"elabtrack-v2"}`. GET /ready uses the existing two-second PostgreSQL checker: data `{status:"ready",service:"elabtrack-v2"}` on success, or standard 503 on absent/failing dependency. Both no-store, both public general-budget reads. Bootstrap still requires initial DB connection before a listener; liveness availability during dependency loss refers to an already-running process. Readiness proves ping only, not migrations/schema correctness, business capability, privilege sufficiency or uptime SLO. Phase 1G verified real local dependency loss/reconnect: health stays 200, ready becomes safe 503 then 200.

Frontend ApiResponse is a typed success/failure union. ApiRequestError carries status/code/message/requestId/optional field errors, without retaining Axios config/request/response. Central normalization rejects nonstandard/legacy/error-type payloads, malformed IDs and unsafe fields, preserves a validated header reference, and always substitutes generic 5xx text. Network errors use status 0/NETWORK_ERROR. Future codes must be added deliberately. Support text can add a reference for unexpected failures via apiErrorMessage; routine validation does not display UUIDs.

Current hooks consume normalized errors. Query does not retry any 4xx; server/network queries retain at most one ordinary retry, mutations none. Phase 1D automatic protected 401 refresh still shares one promise and retries at most once, with auth-endpoint exclusions, generation checks and serialized cookie operations. No generic global toast handler is added. Future forms map error.fields to Hook Form fields, session handling owns authorization/expiry, ordinary failures use safe messages, and unexpected failures may expose a support reference. Health uses the centralized dynamically loaded client, with no bearer, no credentials or automatic auth refresh, and a safe schema adapter; its existing minimal status view now says liveness explicitly.

Operational logs are disposable operations/debug/security/performance records, subject to later chosen retention/access/rotation. They are not durable business evidence. Later business audit events capture who changed what/when under the approved domain model, transaction/history/ledger policy and retention. The stakeholder's boss-rebuild audit emphasizes durable evidence; logs here do not implement that schema or establish its policy.

Zap's core/output and standard request context are future vendor-neutral extension points. No Sentry, telemetry SDK, collector, metrics server, dashboard, cloud provider or business audit schema is added. Operators still need chosen log collection/storage/access/rotation and incident procedures before deployment. Phase 1G verifies local proxy/browser/PostgreSQL/Docker, parser and captured log redaction; its report records exact scope and remaining deployment checks. Runner atomic bookkeeping/checksums/advisory lock/lookup error semantics belong to separately authorized Phase 1H, before relying heavily on product migrations.

## Phase 4A browser integration and deferred boundaries

The product login form uses the existing centralized transport and safe account projection. Protected frontend routes revalidate GET /auth/me; no request-body role or JWT decoding authorizes navigation. Existing generic GET /users/ remains ADMIN-only through the central permission guard; there are no provisioning/status/role-management endpoints. Self PATCH retains display-name-only strict binding with a current active product-role repository predicate.

Login invalid/missing/inactive accounts share INVALID_CREDENTIALS/401. Deactivation after issuance denies protected reads (403) and refresh restoration (401). Ordinary resource 403 does not invalidate an otherwise valid session; current-account denial and persistent 401 follow existing fenced invalidation. Network failures remain recoverable. Logout preserves bodyless 204, serialized cookie revocation/clearing, immediate client/private-cache clearing and non-secret peer signals.

Migration 000004 maps `user`→`BORROWER`, `admin`→`ADMIN`, adds `STAFF` to the role constraint and changes only the role default/values/constraint. Identity, password hashes, status, timestamps and refresh records are preserved. Down reverses representable roles and refuses while any Staff exists; migrate application and schema together. Historical 000001–000003 SQL/checksums are unchanged.

Phase 4B implements versioned terms and acceptance at the contracts below. The borrowing request/direct-issue enforcement point remains a future transactional integration requirement; no borrowing endpoint exists. Phase 5 owns secure provisioning/import/status/privileged account controls. Authentication does not imply accepted terms or an activated onboarding process.

## Phase 4B implemented terms contracts

All routes use `/api/v1`, the existing envelope/request ID/current-account middleware and `Cache-Control: no-store`. GET requires bearer authentication and current active BORROWER/STAFF/ADMIN. Borrower-specific status/acceptance requires current BORROWER; publication requires current ADMIN. Staff/Admin never accepts on a Borrower's behalf. POST also requires an exact trusted Origin and JSON; ordinary API rate/security rules remain. No cookies are issued by terms endpoints.

| Method/path | Success | Data / input |
|---|---|---|
| GET `/terms/current` | 200 | Current immutable document; 503 `TERMS_NOT_PUBLISHED` if none |
| GET `/terms/status` | 200 | Own status: `unpublished`, `required`, `updated`, or `accepted`; never lists other accounts |
| POST `/terms/{versionID}/accept` | 200 | Exact reviewed UUID in path, body `{}` only; returns original receipt on repeat acceptance of still-current version |
| POST `/terms/versions` | 201 | Admin-only immediate mandatory publication: `version`, `title`, `body`, required `expected_current_version_id` (null for initial publication, otherwise last-read UUID) |

Document fields: `id`, `version`, `title`, `body`, `content_hash` (SHA-256 over the exact UTF-8 plain-text body), `published_at`. Publisher identity and creation time are retained in storage but not included in the content DTO. Version identifiers are 1–64 ASCII alphanumeric/dot/underscore/hyphen characters, beginning alphanumeric; title is nonblank and at most 200 UTF-8 bytes with no line breaks; body is nonblank plain text, at most 65,536 UTF-8 bytes. Body is preserved exactly and rendered as text, never executed HTML. No drafts, scheduled effective dates, withdrawal endpoint or generic CMS. Every new publication requires new acceptance before future new borrowing commands.

Status fields: `state`, `current_terms` (document or null), `acceptance` (current receipt or null), `has_previous_acceptance`, `acceptance_required`, `can_initiate_borrowing`. Unpublished means no current document, `acceptance_required:false` (there is no published version to accept), `can_initiate_borrowing:false`; the distinct unpublished state never implies consent. Historical acceptance makes a missing current receipt `updated`; no historical receipt makes it `required`. Only an actual matching current receipt makes the status `accepted`. These flags describe policy eligibility, not an implemented borrowing capability.

Receipt fields: `id`, `terms_version_id`, `accepted_at` (PostgreSQL server timestamp). Identity comes exclusively from the authenticated principal; arbitrary user IDs/times/consent flags and unknown payload fields are rejected. Unique `(user_id,terms_version_id)` makes duplicates/parallel submissions return the same original receipt/time. No separate Idempotency-Key is required for this intrinsically idempotent operation. Checking a superseded version occurs before receipt replay, so old acceptance cannot silently approve new content. Acceptance JSON maximum is 16 KiB; publication JSON maximum is 96 KiB under the global parser bound.

| Error code | HTTP | Meaning |
|---|---|---|
| TERMS_NOT_PUBLISHED | 503 | No applicable document; acceptance/new-command policy fails closed |
| TERMS_VERSION_NOT_FOUND | 404 | Submitted UUID has no version |
| TERMS_VERSION_CHANGED | 409 | Reviewed version is superseded; fetch/review the new document with unchecked consent |
| TERMS_ACCEPTANCE_REQUIRED | 409 | Reusable future-command policy lacks a matching current receipt |
| TERMS_VERSION_EXISTS | 409 | Immutable version identifier already exists |
| TERMS_PUBLICATION_CHANGED | 409 | Expected pointer no longer matches; publication cannot overwrite a concurrent change |

Invalid UUID/payload/content is 400; JSON media 415; body bound 413; missing/invalid session 401; wrong current role/inactive account/untrusted Origin 403; database failures 500 with safe reference only. Unknown/malformed client error payloads retain the existing safe fallback. No success is inferred from network failure.

Publication/acceptance lock the active account before the singleton publication row, using one transaction. Publication takes the pointer exclusively; status/acceptance/policy take it shared. An acceptance that commits first retains its reviewed version; a publication that commits first makes the old review conflict. No operation substitutes another version for the submitted UUID.

**Future Phase 7 requirement:** inside the borrowing command's existing transaction, after current actor/target authorization and sorted participating-account locks, call `terms.Service.RequireCurrentAcceptance(ctx, borrowerID)`. Retain its shared publication lock through stock/borrowing/history/receipt writes and commit. Persist returned receipt ID and borrower ID with a composite FK to `(terms_acceptances.id,user_id)`, and bind the reviewed terms version. Direct checkout separately authorizes Staff/Admin, then checks the target Borrower's own evidence; the actor cannot create acceptance. Pending submissions retain their original bound terms during later issue. Calling this policy in a separate transaction would leave a race and is unsupported. No live request/direct-checkout endpoint or end-to-end borrowing enforcement is claimed.

**Content/activation gates:** No official V2 FSMO document is approved; normal storage has no synthetic publication/acceptance, and isolated tests alone use TEST terms. Both Student/Faculty Borrowers require current officially published terms via Phase4B. DEC-070 is the current future provisioning contract: only Admin creates accounts; Student requires unique textual official Student ID and approved SKSU email; Faculty is individual-only with any valid unique accessible email and no required Student ID. Bulk creation/deactivation is Admin-only and Student-only, with current role/category checks, complete preview before confirmation and retained history. Both use separate borrower-chosen passwords via secure activation; links to the respective mailbox remain recommended. Exact SKSU configuration/roster formatting/matching (OPEN-028), activation/ownership/recovery (OPEN-001/009), Brevo backend integration/API key/verified sender/successful live delivery testing for activation and future recovery (DEC-071/OPEN-017) remain technical/deployment dependencies. DEC-073 resolves OPEN-029: individual/Student bulk deactivation allows outstanding fines/active or overdue loans/unreturned equipment/replacements with warnings/confirmation; it cannot resolve those obligations, change overdue calculations or delete history. Fine clearance remains separately Admin-only/auditable. DEC-072 defers official terms until after FSMO presentation (OPEN-016), gating official publication/live borrowing with documented acceptance, without blocking independent account-management/inventory development. No product provisioning, activation, domain enforcement, bulk deactivation, password-change/recovery or invitation endpoint exists. The [domain routes](domain/API_RESOURCE_DRAFT.md) are future proposals requiring separate Phase5 authorization; current Phase4A authentication is unchanged. See [policy](project/ACCOUNT_PROVISIONING_POLICY.md), [historical Phase4B report](project/PHASE4B_REPORT.md) and [reproduction guide](../integration/PHASE4B.md).

## Phase 5 implemented account-management contracts

Batch 1 implementation, 2026-10-09. Current active PostgreSQL authority is checked by middleware and again under application transaction locks. No public registration, Admin creation/promotion, role/category/identity rewriting or account deletion is exposed. Student/Faculty are BORROWER attributes. Existing unclassified accounts retain their identity and remain outside Student roster matching.

| Route under `/api/v1` | Authority / result |
|---|---|
| GET `/borrowers` | Staff/Admin; bounded `page` 1–10000, `per_page` 1–100, search ≤100 bytes, type STUDENT/FACULTY and status ACTIVE/INACTIVE/PENDING; `data={items,total,page,per_page}` |
| GET `/borrowers/:id` | Staff/Admin; safe identity/profile/status, activation submission state and explicit accountability availability |
| POST `/borrowers` | Admin; strict name/email/borrower_type/student_id/course/contact_number; fixed BORROWER role and pending activation; 201 |
| PATCH `/borrowers/:id` | Admin; name/course/contact_number only, expected_updated_at; no email/ID/category/role mutation |
| PATCH `/borrowers/:id/status` | Admin; required boolean active, confirm=true, expected_updated_at; revokes refresh/activation tokens on deactivation; obligations never veto |
| POST `/borrowers/:id/activation` | Admin; empty object; bounded cooldown, rotate/invalidate old token, submit a new link without account duplication |
| GET `/borrowers/:id/audit` | Admin; immutable events, 25/page, safe actor/target IDs/action/time |
| GET `/borrowers/policy` | Admin; configured Student domains plus technical activation TTL/cooldown; no mail secret |
| GET `/borrowers/student-template` | Admin; binary `.xlsx` download exception, no-store, exact five-column text-cell template |
| POST `/borrowers/rosters` | Admin; trusted-Origin multipart `roster` .xlsx plus `operation` CREATE/DEACTIVATE; full immutable owned preview |
| GET `/borrowers/rosters/:id` | Owning current Admin; reviewed preview/result |
| POST `/borrowers/rosters/:id/confirm` | Owning current Admin; rows integer array + confirm=true; all selected valid rows atomic, identical selection replays committed result |
| GET/POST `/staff-accounts`, GET `/staff-accounts/:id`, PATCH `/:id/status`, POST `/:id/activation`, GET `/:id/audit` | Admin-only separate directory/STAFF provisioning/status/activation/audit; existing ADMIN read-only |
| POST `/auth/activate` | Public single-use credential operation, trusted-Origin JSON token/password, login rate bucket; establishes bcrypt password atomically, no session or terms consent |

Individual mutations require UUID `Idempotency-Key`, normalized payload binding, current-actor reauthorization on replay and immutable transaction receipt. The CORS allowlist explicitly permits this header. Profile/status PATCH accepts trusted Origins; cookie-changing auth handlers remain POST-only. Private responses are no-store. Unknown fields, null/invalid object bodies, role/password additions to provisioning, oversized JSON and unsupported media are rejected.

Workbook limits: 768 KiB compressed, 64 ZIP entries, 4 MiB per entry/XML, 16 MiB total inflated, one worksheet, at most 500 data rows. Headers are exactly `studentId,name,email,course,contactNumber`; text IDs retain zeroes, numeric ID cells are invalid, formulas/macros/external links/embedded files are rejected without evaluating formulas. Required fields are first three; optional fields remain optional. Duplicate rows are all invalid. Preview lifetime is a technical 30-minute default. Confirmation revalidates identity/role/version and aborts all selected changes on a conflict; no silent overwrite/partial selected success. Unselected rows are NOT_SELECTED. Bulk creation leaves activation pending and sends no automatic mass mail; open each created account to send its link. Deactivation matches exact existing Student ID + normalized email even if the current onboarding domain configuration changes; it never guesses by name or treats roster absence as a command.

Activation stores SHA-256 hashes of 32 random bytes (43-character URL-safe token); the raw link appears only in the backend mail call and the user's URL fragment, which the activation page removes immediately. Technical defaults: TTL24h (configurable15m–72h), resend cooldown1m (1m–1h), trusted frontend `/activate` URL, provider timeout5s, bounded delivery transaction8s. Reissue invalidates prior tokens, and inactive/Admin accounts cannot activate. Database/password/audit/consumption commit together. Existing accounts are not silently forced into activation. Email states PENDING/UNCONFIGURED/ACCEPTED/FAILED/UNKNOWN describe submission only. Brevo HTTP201 + messageId means provider acceptance, never recipient delivery; no provider response/raw token is exposed or logged. Live key/sender/recipient testing and approved Student domains remain external readiness gates.

Current obligation reader returns UNAVAILABLE with absent numeric fields; future borrowing/fine/replacement modules must supply authoritative projections. No zero balances, settlement, return, overdue freeze or transaction closure is inferred. Official terms remain independently required before live borrowing through the existing Phase4B guard.

## Phase 6 implemented equipment and inventory contracts

All routes use current-account authentication, `/api/v1`, standard envelope, and `Cache-Control: no-store`. Authenticated Borrowers read only ACTIVE equipment; hidden detail/image IDs return404 and non-ACTIVE filters return403. Staff/Admin read full inventory/history and maintain catalog/categories/ordinary available stock. Reconciliation is Admin-only in the application as well as its protected UI. Mutations require trusted Origin, strict bounded JSON, explicit confirmation where consequential, and UUID `Idempotency-Key`; same actor/resource/operation/key/payload replays one result, changed payload conflicts. Authority is rechecked under sorted account locks before equipment/category locks and receipt writes.

| Method / path | Contract |
|---|---|
| GET /equipment | `page`1..100000, `per_page`1..100(default25), literal search<=100bytes, optional UUID category_id/status, available_only true/false, sort name/available; page items/count and filtered physical totals |
| GET /equipment/:id | Metadata, category, ACTIVE/INACTIVE/ARCHIVED, authoritative stock A/R/C/D/T, metadata_version, stock_sequence, current image_id/timestamps; operational archive_safety reflects future-domain installation |
| POST /equipment | name1..160bytes, description<=2000, optional category_id, opening_quantity0..2147483647, reason required for positive opening, expected_version0. 201 creation/replay; server ACTIVE; opening is one ledger movement |
| PATCH /equipment/:id | name/description/category_id/expected_version; no stock/status/identity fields. Stale metadata409; archived immutable; existing inactive category can remain but new association requires active category |
| PATCH /equipment/:id/status | status, positive expected_version, confirm true; archive irreversible in this phase, guards holds/custody and future liabilities; no delete |
| POST /equipment/:id/adjustments | kind ADD/REMOVE/RECONCILE, explicit integer quantity, bounded nonempty reason, confirm true, optional expected_sequence mandatory for RECONCILE. ADD/REMOVE quantity>0; reconciliation quantity is observed available count, including explicit0. Only A/T change. Unknown reserved/custody/damaged setters reject |
| GET /equipment/:id/movements | Staff/Admin only; page1..100000,25rows, descending sequence; actor/time/reason/kind/deltas/after-vector; immutable audit and receipts commit with stock |
| GET /equipment-categories |100rows/page, name order; Borrower active category choices; operational inactive choices retained |
| POST /equipment-categories; PATCH /equipment-categories/:id | Staff/Admin name<=100bytes, explicit is_active; create expected_version0 / edit current positive expected_version. 201 create/replay,200 update; case-insensitive unique name; no delete; inactivity does not recategorize or hide existing active equipment |
| POST /equipment/:id/image | Staff/Admin multipart exactly image + expected_version, UUID command key; PNG/JPEG input<=512KiB,2048/dimension and4M pixels. Re-encode canonical PNG, strip metadata; version-conflict safe, hash-based replay; no return photos |
| GET /equipment/:id/images/:image | Current image only, equipment-owned, authenticated/scoped visibility; PNG binary envelope exception, no-store/nosniff, no arbitrary locator or external URL |

Paired000007 uses bigint checked quantity arithmetic with per-pool2147483647 bound, conservation `T=A+R+C+D`, immutable stock movements/audit/receipts and immutable canonical PNG rows. Zero opening has no fictitious movement. Loss is historical and replacement liability is separate; neither gets a current-stock bucket. Runtime has named/narrow DML grants and no history updates/deletes/DDL. Down refuses inventory data/history. No synthetic catalog/category seeds or serialized tracking.

Archive currently allows R/C zero only when PostgreSQL proves the future borrowings/replacement_obligations tables are absent (`NOT_INSTALLED`). If either appears before a real liability adapter is connected, safety becomes UNAVAILABLE and archive fails closed. Phase7/8 must connect authoritative outstanding checks and migration/runtime grants under account→borrowing/obligation→equipment lock order. This installation check does not claim integrated replacement balances. Damaged-original repair/disposal and reconciliation of D are unsupported pending later policy; available reconciliation never hides those originals.

Database PNG storage is a bounded local strategy, not an object-storage-provider selection. Backups include image data; retained prior images preserve audit references. Large files, SVG/GIF, oversized dimensions, arbitrary external links, image removal/history deletion and return-evidence uploads are unsupported. Future object-storage scaling/retention requires a separately reviewed migration and provider decision.

Frontend routes: `/borrower/equipment[/UUID]`, `/staff/inventory[/new,/categories,/UUID,/UUID/edit,/UUID/adjust]`, `/admin/inventory/UUID/reconcile`. Existing Borrower terms gate stays; catalog reads do not initiate borrowing or imply consent. No cart/request submission, return, replacement, lost correction, fine settlement or Phase7 implementation is included.
