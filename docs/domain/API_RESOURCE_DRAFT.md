# Phase 2 API resource draft

**ENGINEERING RECOMMENDATION**, design-only, 2026-10-08. No routes, handlers, DTO packages or OpenAPI artifact are implemented. Final actor permissions and policy-dependent commands require the [decision gate](BUSINESS_RULES.md#policy-decision-matrix). Existing [Phase 1 API envelope](../API_CONTRACTS.md) and security/session architecture remain unchanged.

## Resource and command shape

All paths below are beneath `/api/v1`. “Capability” means approved future server authority, not a role implicitly granted by a menu/JWT/body field. Existing auth/user endpoints remain as implemented; they are not replaced by this draft.

| Resource / command | Intended contract and authority | Gate / integrity |
|---|---|---|
| GET /equipment; GET /equipment/{id} | Bounded approved catalog projection/search/filter/detail; never borrower directory/private loan data | Active catalog/availability visibility, optional categories/images; server paging |
| POST /equipment; PATCH /equipment/{id} | Inventory metadata creation/narrow edit with expected metadata version; fields exclude reserved/issued counters | Approved inventory capability; stale version conflict; new initial stock via one audited command |
| POST /equipment/{id}/stock-adjustments | Typed authorized add/remove/repair/retire/correction; positive quantities, reason, expected/current guard | OPEN-022; current lock/ledger/audit; no generic set-all-counts |
| POST /equipment/{id}/archive; optional reactivate/inactivate commands | Approved catalog lifecycle operation; retain history | OPEN-012; no archive with holds/outstanding under recommendation |
| GET /equipment-categories | Conditional approved taxonomy for catalog filter | OPEN-011; optional admin category commands require policy |
| GET /equipment/{id}/image; conditional POST image attachment | Narrow image purpose, bounded verified upload, catalog access permission | No GET /files/{id} granting private reports by image permission |
| GET /borrowers; GET /borrowers/{id} | Approved staff borrower management projection; no public directory | OPEN-001/002/003/008/009; eligibility/provisioning commands belong later approved identity phase |
| GET /me/borrowing-eligibility | Own minimal current eligibility/guidance, not arbitrary policy impersonation | Distinct from /auth/me; exact disclosure policy after confirmation |
| GET /borrowing-policy; GET /terms/current | Approved effective policy/terms text/version and governing timezone/deadline inputs | OPEN-006/015/016/023; no fabricated terms/rate |
| POST /terms/acceptances | Attributable acceptance if selected account-level flow; per-borrowing acceptance may be inside submit instead | OPEN-016; immutable version/content evidence and authenticated borrower |
| POST /borrowings | Submit own request: unique equipment/positive quantity lines, requested due input, agreed terms evidence/version; channel bound by server | Idempotency-Key recommended; OPEN-020 reservation; server chooses owner/state/value/rate |
| GET /borrowings; GET /borrowings/{id} | Own borrowings unless current approved staff capability; includes requested/reserved/issued/disposition/outstanding, separate calendar/financial summaries | Query ownership applied before paging; no client-only UID filter |
| POST /borrowings/{id}/approve | Authoritative legal decision, expected version/reason if required | OPEN-021: separate approve or confirmed combined release; not a guessed endpoint meaning |
| POST /borrowings/{id}/deny | Retained denial decision/reason and actual hold release | OPEN-014; lock pending state; no DELETE |
| POST /borrowings/{id}/cancel | Conditional approved owner/staff preissue command; omitted if not allowed | OPEN-013; no active custody erase |
| POST /borrowings/{id}/checkout | Actual release only if separate release adopted; current actor/target/catalog/terms/due checks | OPEN-020/021/023; issue snapshots/fresh locks; state-safe retry |
| POST /borrowings/direct-checkout | Conditional staff direct issue specifying target existing borrower and item intent | OPEN-020; shared issue rules; Idempotency-Key; no inactive target bypass |
| POST /borrowings/{id}/returns | One immutable event: unique borrowingItemId lines, good/damaged/lost quantities and notes; return operator permission | Idempotency-Key; DB ownership/uniqueness and locked outstanding; duplicate IDs fail whole request |
| GET /borrowings/{id}/returns | Authorized chronological immutable events and disposition evidence | Same parent authorization; no writable returned checkbox or arbitrary event edit |
| GET /charges; GET /charges/{id} | Own liability or approved staff financial projection: original assessment, adjustment components and derived balance separately | No assessed total labeled cash/revenue; approved source/currency rules |
| POST /charges/{id}/adjustments | Typed reasoned approved waiver/reduction/increase/administrative settlement/reversal | OPEN-007/024; charge lock/nonnegative balance/source key; no base amount PATCH |
| GET /charges/{id}/adjustments | Authorized immutable charge history | Resource ownership/current capability |
| GET /reports/inventory, /reports/borrowings, /reports/returns, /reports/accountability, /reports/usage | Approved staff bounded read metrics/filters and as-of calendar semantics | OPEN-003/006/007/015; no default public/global export; formats deferred |
| GET /audit-events | Conditional restricted staff audit review, bounded entity/action/date filters | Approved audit access/retention; no credential/provider raw data |
| /kiosk/catalog, /kiosk/handoffs, /kiosk/claims (conditional concepts only) | Narrow shared-device catalog/attributed intent and approved handoff, if selected | OPEN-010; no device auth protocol accepted; normal canonical borrowing use cases; no private session hydration by default |

No standalone fine truth, `/records` copy, arbitrary borrowing status PATCH, ordinary transactional DELETE, anonymous checkout, tenant/lab routes, serialized units, online payment gateway or generic public/private file endpoint is introduced. Optional cash/payment routes require OPEN-007 confirmation and reviewed Payment/Allocation/refund contracts; they are deliberately not defined as existing scope.

## Mobile-first contract check

| Borrower need | Server response/interaction |
|---|---|
| Browse/search/filter | Paginated concise catalog cards, optional confirmed categories/availability; detail only on demand; server-authoritative current availability |
| Equipment details | One requested pool and allowed image/description/status; no global collection download or private borrower names |
| Request submission | Bounded unique item intent; explicit selected due/terms; definitive receipt/entity ID after commit; availability conflict returns safe authorized item details |
| Request status | One canonical borrowing with submitted/decision/issue evidence, derived due/custody and explicit hold semantics; not multiple copied status documents |
| Active/due borrowings | Owner-scoped paginated `/borrowings?view=active` plus approved due filters; one calendar contract; no client-generated fine authority |
| History/accountability | Owner-scoped stable paging, immutable issue snapshots/return events/charge components; estimates visibly unposted, no private staff audit by default |

Recommendation for list contract is reuse current numeric `page`/`per_page` envelope conventions where suitable: default 20 / maximum 100 for ordinary lists, validated server bounds, stable `(createdAt,id)` or `(name,id)` ordering, total metadata only when justified. These are **engineering transport bounds**, not institutional borrowing limits. Cursor paging may replace offsets for growing chronological lists if measured need warrants it; no new generic pagination subsystem now. Query text length and allowed filters/sorts are bounded/validated, SQL parameters bound and sort columns allowlisted. Invalid dates/enums are errors, not silently ignored. Owner filters precede count and page. Sensitive authenticated Query data must use Phase 1I private-cache classification/fences in later frontend code.

Proposed submission returns canonical ID/reference, current state, quantities, immutable accepted versions, explicit reservation effect and agreed/requested due evidence. Proposed return result identifies ReturnEvent, new outstanding/completion and any posted charge IDs. HTTP success occurs only after commit. A same-key replay identifies the original result, with fresh current correlation ID and current access check; a stock estimate never constitutes confirmed issuance. Client-supplied actors/prices/state/fee totals are ignored/rejected as authoritative inputs.

## Future domain error catalog

Draft codes are intentionally small and unimplemented. Use existing `{success,message,data,meta,error}`/safe requestId boundary; foundation validation remains 400, authentication401, authorization403, absent authorized resource404, integrity conflict409. Error details must not disclose another borrower's profile or catalog/charge evidence outside current authority. Body/parser/rate/server statuses remain Phase 1 behavior.

| Proposed code | HTTP | Meaning |
|---|---|---|
| EQUIPMENT_NOT_AVAILABLE | 409 | Requested pool has insufficient eligible available/held stock at approved boundary |
| BORROWER_NOT_ELIGIBLE | 403 | Caller/selected target fails approved borrowing gate; disclose only authorized guidance |
| BORROWING_STATE_CONFLICT | 409 | Command is incompatible with current legal state/expected version |
| RETURN_EXCEEDS_OUTSTANDING | 409 | Submitted disposition exceeds locked remaining custody |
| DUPLICATE_RETURN_ITEM | 400 | Input repeats a borrowing item ID; reject entire event |
| STOCK_CONFLICT | 409 | Stock/count/archive/version condition changed or fails reviewed integrity rule |
| CHARGE_ALREADY_SETTLED | 409 | No positive balance for requested discharge; waiver/reduction can also create zero balance |
| ACCOUNTABILITY_CONFLICT | 409 | Adjustment would over-discharge/reverse or source already assessed |
| IDEMPOTENCY_CONFLICT | 409 | Same scoped command key with different normalized payload |

Invalid quantity/notes/due/terms/version fields use existing VALIDATION_ERROR with safe field details; ownership mismatches may return404 to avoid disclosure. No handler mapping, automatic retry on every409 or guessed business policy is added now. [INVARIANTS](INVARIANTS.md) ties the later HTTP contract to domain/real-database tests.
