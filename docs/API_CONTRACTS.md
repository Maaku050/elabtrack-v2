# Foundation API contracts and operational observability

Phase 1F technical contracts with Phase 4A authentication updates, 2026-10-08 (Asia/Shanghai). Current source and tests govern implementation; this document does not authorize later phases. Product roles are BORROWER, STAFF and ADMIN; feature provisioning, terms and business resources remain deferred.

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

Safe role values are exactly `BORROWER`, `STAFF`, `ADMIN`; Student/Faculty are future category attributes, not authorities. Old JWT role hints do not override current PostgreSQL state. Browser session data contains `access_token,expires_at,token_type,user`, with user `id,email,name,role,is_active`. Raw refresh travels only via HttpOnly cookie. No new credential transport, retry, family-revocation or cross-tab policy is implied.

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

Phase 4B will introduce versioned terms/acceptance evidence at the reviewed domain ports and future request/direct-issue gate. Those candidate routes/tables in docs/domain/API_RESOURCE_DRAFT.md are not implemented endpoints. Phase 5 owns secure provisioning/import/status/privileged account controls. Authentication does not imply accepted terms or an activated onboarding process.
