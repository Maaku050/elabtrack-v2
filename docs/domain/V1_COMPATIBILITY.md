# Phase 2.5 V1 compatibility and migration inputs

2026-10-08. **ACTUAL V1 BEHAVIOR** means the supplied [static V1 audit](../reference/v1-audit/README.md), not fresh source/Firebase/deployed inspection. Audit15's containment patch is undeployed. **CONFIRMED V2 DECISION** means the owner's current Phase 2.5 working direction, informed by implementing V1 and knowing FSMO's manual workflow. Intended operational policy and V1's actual code may differ; preserve both. No live export, importer, conversion, migration or cutover occurs here.

## Preserve / modernize / supersede

| Capability | Audited factual V1 behavior | Current V2 disposition / evidence |
|---|---|---|
| Borrower population | Audit03/07/13 student-oriented UI/role; no distinct faculty role | Owner D1/2: generic active Borrower; Student/Faculty category, not roles; supersedes student-only assumption |
| Provisioning / bulk creation | Audit05: Staff/Admin accounts, bulk spreadsheet including password, no signup | DEC-070 supersedes D4 Staff grant: Admin alone creates accounts; Student requires unique textual official ID/SKSU email, Faculty individual-only valid accessible email/no required ID; standard Excel creation/deactivation Student-only with complete preview, no passwords/role assignment; secure borrower-owned activation remains unimplemented |
| Staff/Admin distinction | Audit03/15 broad shared operations including privileged user management/fines | Supersede with three roles D3/26; named multiple Admins; Admin-only fine clearing/deactivation/privileged administration |
| Active account | Audit05/09 client status check; current-session weaknesses | D2: active registered account gate, no independent suspension/eligibility subsystem. Current Phase 1 server freshness remains; old custody retained |
| Request reservation | Audit07/09 available decreases/borrowed increases on submit; race-prone writes | Preserve intent D6: explicit reserved pool, server transaction/locks, available not locally authoritative |
| Approval as release | Audit07 Request→Ongoing resets borrowedDate, no stock recheck | Preserve operational combined edge D5; recheck target/stock and immutable issue evidence; no waiting-approved stage |
| Request expiry | No audited24h reservation expiry | New working default D7: EXPIRED after24h, release/history; versioned configurable candidate |
| Borrower cancellation | Audit07 absent; manual delete restores inaccurately and loses history | New D8 own PENDING cancel; no issued delete/cancel restoration |
| Denial | Audit07/08 restores stock/deletes request, no reason record | D9 retained DENIED plus required borrower-visible reason; transactional release |
| Direct checkout | Audit07 active student selected, Ongoing immediately; no same maximum | Preserve D10 with active existing Borrower, due/time/terms/availability and actor/target checks |
| Duration/time | Audit07/09 borrower UI seven-day maximum; midnight Manila/calendar variants | Supersede fixed maximum D11; required due_at date+time, Asia/Manila D12; absolute timestamps |
| Terms | Audit04/07 first-request flag and time, no version/hash | D13 current-version first-use acceptance user/version/time, not per-loan; material update gates next request |
| Partial/mixed return | Audit07 Staff/Admin disposition helpers, cumulative good/damage/loss, only good restocks; intermediate state overwritten. Return-photo workflow is not established by the audit | Preserve partial capability D17/18/23; immutable structured condition/quantity events, duplicate protection, one authoritative map. Clarified D17: a borrower may only physically show their own phone photo in person; Staff/Admin alone records the return. No software photo/evidence submission or return-photo feature |
| Damage/loss | Audit07/09 actual code prices snapshots and closes when good+damage+loss=issued | Owner D19–24 supersede priced liability and premature closure: separate replacement obligation, physical custody0 can stay open |
| Original damaged stock | Audit06/07 good restocks, damaged/lost disappear from borrowed but not total reconciliation | Selected hybrid four physical counts, loss removes total, damage nonusable held, replacements acquire stock; disposal policy later |
| Overdue fine | Audit07/09 PHP 10/calendar day; daily/current calculations differ | D14/15 confirmed PHP 10/day while physical OR replacement unresolved; selected ceil elapsed24h; freeze at completion |
| Existing fine/new request | Audit07 no source-visible auto-block | D16 request/reserve permitted; staff may deny after human review |
| Fine clearing | Audit04/08 zeroes record fine and sets finePaidAt without updating separate fine truth | D25–27 Admin full clear, immutable assessed/clear history/method; no partial/payment allocation; no erased assessment |
| History | Audit04/07/08 creates terminal copies/fines/notices before batch deleting active | Retain one canonical borrowing; atomic terminal events, fine freeze, audit/outbox; no copy/delete |
| Notifications | Audit08/12/13 producers/two shapes, consumer absent | Preserve intent; outbox/attempts extended to expiry/partial/replacement; provider/cadence later |
| Reporting | Audit08/10 unbounded client aggregation; assessed fines called revenue | Bounded scoped projections; recorded PAID distinct from waived/other and assessed/outstanding |
| Inventory edit/delete | Audit05/06 client guards/hard delete, image replacement risks | Future metadata version, stock ledger, no hard consequential deletion; guarded archive/least privilege |
| Interactive Kiosk | Audit01/13 shared responsive borrower routes; no device-auth lifecycle evidenced | Owner D28/29 corrects meaning: catalog/search/filter/cart/review/request in normal mobile-first flow; reject dedicated hardware assumptions |

The local capstone's faculty intent (P401–402/P613), conflicting registration wording (P421/P827/P912), kiosk wording (P429/P440) and Super Admin recommendation (P922) remain **CAPSTONE INTENT**, not implementations. Current owner direction resolves population, provisioning and functional kiosk interpretation; expansion is outside FSMO. No biography/contact/sample user is reproduced.

## Collection mappings for later authorized migration

| Source | Relational mapping candidate | Evidence gaps / required reconciliation |
|---|---|---|
| Firebase Auth/users | UID→stable user mapping; category from legacy role where supported; Staff/Admin grants separately approved | V1 staff/admin role does not grant V2 Admin automatically; missing identities quarantined; no password extraction/plaintext preservation; verification flags do not create eligibility rules |
| equipment | One record→aggregate pool plus explicit reconciled opening ledger | borrowedQuantity mixes pending holds/custody; split using linked actual transactions; damaged/lost total gaps require reconciliation, not invented stock buckets |
| transactions with embedded items | Canonical borrowing with multiple items, raw status/provenance, actual dates and issue evidence | Map Request/Ongoing/Incomplete/Ondue/Overdue through quantities/state; never invent timestamps or enforce24h retrospectively without an approved cutover policy |
| records | Retained terminal history linked to canonical legacy source | Old completion includes damaged/lost accounted, not proven replacement. Do not invent accepted replacements, reopen historical obligations automatically or price them under new rules. Approved historical-only classification/reconciliation may be needed |
| fines plus record.fineAmount/finePaid/finePaidAt | Source monetary evidence and reconciliation issues, then approved mapping to fine basis/clearance | Conflicting truths, zero overwritten amount, no payment proof. Never infer PAID from 0 or migrate price-based damage/loss as newly confirmed V2 overdue fine |
| notifications | Source intent/delivery unknown; preserve only approved required lineage | Do not replay all old email queues or fabricate provider acceptance/delivery; minimize recipient/body retention |
| Storage refs | Authorized catalog image mapping only when owner/rights/hash known | No assumption legacy URL grants private access; any profile assets need separately approved purpose. Return photographs/attachments are excluded from V2 software scope; no return-photo storage/import mapping |

Three deferred tables (legacy_import_batches, legacy_mappings, legacy_reconciliation_issues) preserve source system/collection/document, authorized export hash/revision, target mapping, unknowns/quarantine and signoff. No boilerplate legacy fields on every new entity, flattened single-item assumption, generated NOW due times, fake actor, inferred replacement acceptance or assessed-as-paid reporting.

## Cutover and privacy gate

OPEN-018/026 remain NON-BLOCKING for core UX but gate real migration/cleanup: authorized dataset and access, retention/legal responsibility, legacy balance classification, replacements with missing evidence, snapshot reconciliation, backup/rollback and signed cutover. Historical financial amounts must retain original basis; do not silently recalculate closed V1 calendar-day charges under the new elapsed-day interpretation. A future approved importer must prove multiitem ownership, atomic mappings, nonwriting/isolated dry runs and stock/fine/event reconciliation with sanitized fixtures. Nothing accesses V1 Firebase in Phase 2.5.
