# Phase 5 — user and borrower management

**PHASE 5 PARTIALLY COMPLETE — EXTERNAL GATES.** Batch1 core engineering, authorization, persistence and visual gates PASS on 2026-10-09. The owner-authorized independent Phase6 may proceed. Live email delivery, approved SKSU domain inputs and official FSMO terms remain unverified external dependencies; Phase4B is not fully complete. Phase7 is not authorized.

## Implemented boundary

Admin-only individual Student/Faculty Borrower creation, separate Staff provisioning, operational Staff/Admin borrower directory/search/filter/pagination/details, profile and warned confirmed status changes, activation/reissue, account audit, Student-only XLSX template/preview/selected atomic creation and deactivation. Faculty accepts a valid external email without Student ID; Student IDs stay strings, including leading zeroes. Existing unclassified Borrowers are preserved and excluded from Student bulk matching. Admin directory is read-only: no privileged Admin creation, promotion or status mutation.

New paired migration000006 adds normalized borrower profiles, pending activation metadata, hash-only activation tokens, immutable account audit/operation receipts and owned immutable roster previews with one committed result. No old migration changes, recategorization, historical deletion or business obligation settlement. Down refuses consequential history. Runtime grants exclude history rewriting/deletion. Normal local database was read only: zero users/terms/acceptances and migration000005; new migration was applied only to disposable Batch1 PostgreSQL. Use the existing operator migration workflow before running this feature against another database.

Actual endpoint/DTO/limits/retry contracts are in [API contracts](../API_CONTRACTS.md#phase-5-implemented-account-management-contracts). Routes are `/api/v1/borrowers`, `/borrowers/policy`, `/borrowers/student-template`, `/borrowers/rosters`, `/borrowers/rosters/:id[/confirm]`, `/borrowers/:id[/status,/activation,/audit]`, `/staff-accounts` and `/auth/activate`. Mutations require current backend authority, strict bounded payloads, trusted Origin and resource-bound UUID idempotency keys; public activation uses a purpose-bound single-use token instead of a management key.

New pending accounts have an unusable password sentinel. Secure random activation credentials are persisted only as SHA256, expire, rotate, and invalidate on deactivation. Borrowers choose their separate bcrypt password themselves. Redemption does not authenticate or accept terms. Deactivation revokes sessions/activation without changing obligations; reactivation does not revive credentials. Profile/status stale expectations and sorted account locks protect concurrent changes. A delivery/status lock-order issue found during final review was corrected and verified with a real PostgreSQL barrier race.

Student Excel processing uses maintained Excelize with bounded ZIP/XML preflight, string-ID formatting, formula/external/macro/embedded-file rejection and unsigned ZIP64 checks. All rows are previewed; explicitly selected valid rows are rechecked and commit atomically, with durable repeat results. Changed selection, stale identities and conflicts cannot produce partial writes. Bulk creation intentionally does not send a mail campaign: Admin opens each pending account to issue its activation link. This avoids storing raw tokens or adding an unapproved queue.

Brevo integration follows the [official transactional endpoint](https://developers.brevo.com/reference/send-transac-email): backend-only key, bounded requests/replies, no redirects, safe errors. A201 messageId means provider acceptance, not delivery. Adapter tests use synthetic transports; no live provider request was made. Missing configuration remains UNCONFIGURED; ambiguous submission remains UNKNOWN. Parser source is [Excelize release](https://github.com/qax-os/excelize/releases) and [bounded workbook options](https://xuri.me/excelize/en/workbook.html).

## Actual checks

| Gate | Result / evidence |
|---|---|
| Go fmt, vet, full tests | PASS, newly rerun after final lock-order change; `/tmp/elabtrack-batch1-phase5-go.log` |
| Targeted race detector | PASS auth/terms/security/config/account domain/mail/parser, including hostile ZIP64 fixture; `/tmp/elabtrack-batch1-phase5-race.log` |
| Frontend lint | PASS,19 inherited warnings, no new warnings |
| Frontend tests/build | PASS203 tests in12 files; production build PASS |
| Fresh real PostgreSQL | PASS migration integrity, empty paired rollback/reapply, privilege/constraint/history guards, terms regression, account identity/bulk/activation/rollback/concurrent commands |
| Delivery/status race | PASS `TestRealAccountDeliveryRace`, actual PG locks/barrier; `/tmp/elabtrack-batch1-delivery-race.log` |
| Real HTTP security | PASS `TestBatch1HTTP`, role/origin/CORS/strict body/current account checks including PATCH |
| Existing network regressions | PASS `TestRealFoundation` and `TestRealTermsHTTP`: real auth/refresh/replay/inactive/current-role/terms persistence |
| Integrated Chromium | PASS20 checks/27 screenshots; [acceptance](../ux/verification/phase5/acceptance.json), real API/PG with private test-only captured mail |
| Existing visual previews | PASS actual Chromium layouts, no runtime errors/warnings, production fixture exclusion; [results](../ux/verification/phase5/preview-regression/RESULTS.json). Anonymous refresh stub is isolated preview evidence only |
| Normal/isolated Compose, diff check | PASS |
| Preserved hashes | All10 original migration files,62 shadcn primitives,47 approved visual files unchanged |
| Live Brevo, actual SKSU compliance, official terms | BLOCKED external configuration/approval; synthetic fixtures do not prove these |
| Deployment/TLS/production provider tests | NOT RUN; no deployment authorized |

Frontend logs: `/tmp/elabtrack-batch1-phase5-{lint,test,build}.log`. PG logs: `/tmp/elabtrack-batch1-{migrator,terms,accounts,http,regressions}.log`. Tests use PostgreSQL18.6 on isolated54832, Go1.27.1, Node24.19/npm11.17 and actual installed Chromium. The global Go binary differs (1.26.5); backend selects its declared toolchain. No secret or private fixture credential is exported here.

## Screens and visual assessment

Actual approved S04-01/02/03/04 and S05-05 PNGs were inspected alongside implementation captures. Directory table/search/filter/action grouping, two-column detail/accountability cards, four-stage roster review and three-card Administration with restricted directory preserve the approved navy/violet shell, token palette, rounded surfaces and responsive structure. Light/dark and1440/1024/768/390 layouts, keyboard focus and dialogs were exercised. Inherited Phase3B shell adaptations remain; no approved PNG or primitive changed.

Deliberate scope adaptations: immutable ID/category/role; dedicated conditional Student/Faculty and restricted Staff forms; real activation status instead of invented last activity; account audit instead of fabricated borrowing history; accountability explicitly UNAVAILABLE instead of fictitious zero balances; no fine-clearance or privileged Admin controls. Missing official terms are never invented; synthetic terms shown in screenshots are disposable test records only. This is an engineering visual gate, not a new owner visual approval.

## Remaining dependencies and handoff

Configure approved Student domains, verified Brevo sender/key and permitted HTTPS activation origin, then test real receipt/ownership before declaring delivery working. FSMO must finalize/publish official terms and establish operational activation/recovery responsibility. Password recovery remains later scope. Phase7/8 must supply real obligation projections; positive projected obligations were injected into real-PG deactivation tests and did not veto deactivation. Current projection availability is not a claim of integrated loan/fine conservation.

The exact phase-boundary file/status snapshot is [PHASE5_FILES.txt](PHASE5_FILES.txt). Initial tree was clean at HEADdc04db7; index remains empty. New source and evidence are unstaged. No commit/push/deploy/V1 access or normal config changes occurred. Phase6 changes after this checkpoint are recorded separately.
