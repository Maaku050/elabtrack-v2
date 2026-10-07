# Future UX implementation acceptance checklist

Phase 3A.1, 2026-10-08. **Specification only; every checkbox below is pending.** No screen/component/UI/browser/accessibility/theme implementation was tested. Apply later against approved Phase 3A.2 mockups and actual role-authorized implementation; current text wireframes are hierarchy/behavior drafts. [Inventory](UX_ARCHITECTURE.md), [wireframes](WIREFRAMES.md), [responsive contract](RESPONSIVE_RULES.md), [status system](STATUS_SYSTEM.md) govern.

## Per-screen review record

Use one record per inventoried surface ID; variants sharing a route still need their conditions reviewed.

| Evidence field | Future reviewer records |
|---|---|
| Surface / role / feature version | Inventory ID, approved mockup reference, current commit and supported actor |
| Domain contract | Applicable rule/invariant/source and expected event/stock/quantity/fine outcome |
| Viewports / input / theme |320,390,768,1024,1440 CSS px; touch/keyboard;200% zoom; light/dark |
| Scenario data | Sanitized fixtures, lifecycle/balance/partial/replacement and inactive/authorization cases |
| Results | Observed behavior, evidence/link, pass/fail, defect, owner and retest disposition |

No placeholder checklist tick proves compliance. Do not weaken backend invariants/type checking to pass visual review. Full WCAG certification and production delivery evidence are separate work, not claimed by this design.

## Every implemented screen

- [ ] Match approved hierarchy/layout/copy/actions at mobile, tablet and desktop; document approved deviations.
- [ ] Preserve four borrower destinations and operational/Admin visibility; deep-link/back/filter/scroll behavior works.
- [ ] Handle initial Skeleton, background refresh, genuine empty, filtered empty, network/server/validation/session/ordinary-forbidden states.
- [ ] Show authoritative counts/amounts and as-of/live/final state; loading/error never displayed as0 truth.
- [ ] Handle stale/conflict/unknown-write response without silent data overwrite, automatic clipping or duplicate submission.
- [ ] Verify permission and ownership at server boundary, including replay; hiding a control isn't the authorization test.
- [ ] Maintain private query/cache/draft/session-generation behavior inherited from Phase 1I; no cross-user late result.
- [ ] Verify light/dark semantic readability, text status, focus and error indicators without color dependence.
- [ ] Use existing shadcn/Lucide foundation without altering primitives just to fit a one-off screen.
- [ ] Check semantic landmarks/headings, labels, accessible names, descriptions, validation summary and field-error association.
- [ ] Keyboard order is logical, every action reachable, focus visible; overlay traps/restores focus and Escape cancels safely.
- [ ] Primary touch controls at least44×44 CSS px design target;320px/200%zoom has no whole-page overflow/obscured content.
- [ ] Sticky actions clear nav/safe area/keyboard; focus/error summary never hidden; scoped table scrolling labeled and keyboard reachable.
- [ ] Equipment images have appropriate alt/fallback; all icon-only controls have accessible names; charts if later added have text alternatives.
- [ ] Respect reduced motion, no essential drag gesture or pointer-only calendar/time selection; live regions don't spam countdown updates.

## Borrower feature checks

- [ ] SH-01 has provisioned-account sign-in and contact instruction, no public signup/role impersonation or unapproved reset flow.
- [ ] B-01 accepts current terms once/version; material update gates new submission while existing obligations remain viewable; no per-loan consent.
- [ ] B-03–06 show available catalog quantity and draft selection, not reserved stock, commerce payment or internal lost/damaged buckets.
- [ ] B-07/08 require date AND time with Asia/Manila interpretation, no hard seven-day maximum or autoextension.
- [ ] B-09/10 confirms actual committed Pending/reserved state, expiry and face-to-face approval/release next step.
- [ ] Own pending Cancel confirms release/history and handles expiry/other staff action races; no issued cancel/delete.
- [ ] B-11 visible required denial reason; B-12 expiry distinct from denial; B-13 cancellation retained; new request uses fresh checks.
- [ ] B-14–18 shows original due, physical vs replacement remaining, partial/overdue and final completion correctly; no Mark Returned/Complete.
- [ ] No role has a return-photo upload, camera, attachment, transmission/retention or borrower return-evidence interface; D17 photos stay outside software.
- [ ] Replacement-only loan P0/U>0 stays Active/Overdue as due passes; incident history survives accepted replacements.
- [ ] Fine live/final/outstanding/clear history is read-only for Borrower; no price-based damage charge, Pay Online or partial settlement.
- [ ] Old fine permits new request/reservation; no automatic fine block disguised as UI validation.
- [ ] History mobile cards and bounded lists include Completed/Denied/Cancelled/Expired; final money can remain after completion.
- [ ] Notifications link to own allowed event/detail, don't invent read/delivery proof or expose staff-only notes/audit.
- [ ] Account supported edits only; role/status/category/email changes follow separately approved controls; sign-out/theme behavior preserved.

## Staff/Admin processing checks

- [ ] Queue row opens review; no context-free row approval or bulk handover; expiry/account/accountability visible.
- [ ] Approve & Release rechecks current target/hold/due while handover, commits one CHECKED_OUT edge; no approved-waiting/extra release.
- [ ] Denial reason required/borrower-visible; concurrency cannot overwrite another terminal state or release twice.
- [ ] Direct issue selects active existing generic Borrower/current terms/available stock/date+time; immediate checkout/no pending.
- [ ] Return shows issued/prior good/damage/loss/physical remaining separately from now entries and replacement outstanding.
- [ ] Full-good shortcut fills only physical remaining and cannot resolve existing replacement obligations; draft overwrite is explicit.
- [ ] Duplicate return/obligation inputs/cross-parent refs and overreturn/overaccept rejected; UI preview and server normalized map reconcile.
- [ ] Damage/loss decreases physical custody, creates replacement quantity, retains history, no automatic price assessment.
- [ ] Replacement acceptance adds actual new usable stock/total once, leaves damaged original held/history; partial acceptance remains open.
- [ ] Completion only allP0/allU0, final overdue freeze atomic; no generic completion or custody-hiding correction control.
- [ ] Staff can view fine and provision Borrower only; cannot clear/deactivate/bulk/assign Admin or see administrative audit/export.
- [ ] Admin full-clear method/actor/time/history and whole locked balance; no amount input/partial payments; stale balance forces review.
- [ ] Active fine cleared now can later accrue delta; completed final amount cannot change; payment distinct from waiver/other.
- [ ] Admin account lifecycle confirms impact and preserves history; privileged accounts named; last-admin/recovery behavior separately approved.
- [ ] Bulk template/validation/preview/explicit subset/result handle duplicates and unknown batch outcome; no password spreadsheet/delivery assumptions.
- [ ] Inventory current A/R/C/D/T reconcile, liability/incidents separate; metadata cannot overwrite held/custody; archive guard enforced.
- [ ] Terms published as new approved version; policy defaults read-only until configurable rules approved; no generic engine/provider UI.
- [ ] Admin reports/audit scoped/bounded; appropriate filters/ranges/status/method totals; exports only after approved format and access.

## Required journey and adverse-case evidence

| Journey / risk | Later expected proof |
|---|---|
| A request→FSMO issue | Mobile selection/current terms/date-time/reservation/expiry, staff physical review, checked-out result |
| B expiry | Exact24h and open-view conflict, released hold only from committed server response, retained history |
| C direct | Active target/current terms/last stock; immediate issue; denial-safe inactive target |
| D full good return | P0/U0 automatic completion/fine freeze and correct restock |
| E partial then final | Unchanged due, multiple immutable returns, correct intermediate/final state |
| F lost/damaged then replacement | P0/U>0 open, correct acquisition vs damaged-held total, partial/last acceptance |
| G overdue→complete→clear |1min/24h/25h, replacement-only overdue, frozen final, Admin-only whole clear/history |
| H bulk provision | Validation/duplicates/explicit selected subset/results, secure activation dependency and retry semantics |
| Last-stock / concurrent staff | One legitimate stock claim/state edge, loser guided refresh/review, no double effect |
| Duplicate/unknown writes | Double taps/input duplicates/original-key replay/acknowledgment loss and fault rollback |
| Authorization/session | Cross-user/resource/Staff privilege attacks, account disable/demotion/current role, existing cross-tab lifecycle |
| Accessible responsive | All surface IDs exercised by keyboard/touch/light/dark/zoom/small-width/error+focus paths |

Acceptance evidence is recorded only after separately approved implementation. Phase 3A.1 exit is document coverage/consistency/scope/whitespace, not passing these future checks.
