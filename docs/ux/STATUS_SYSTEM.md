# Status system — Phase 3A.1

Low-fidelity semantic contract, 2026-10-08. Exact final colors/icons/type styling deferred to Phase 3A.2. [Domain machine](../domain/STATE_MACHINES.md) is authoritative. No new lifecycle enums, flags, money policies or periodic mutation. Text conveys meaning in both themes and assistive technology.

## Lifecycle labels

| Stored lifecycle | User-facing label | Borrower meaning / permitted action | Staff/Admin meaning |
|---|---|---|---|
| PENDING | Pending | Reserved; visit FSMO before expiry; own confirmed Cancel only while pending/unexpired | Review then Approve & Release or Deny; no approve from queue row |
| CHECKED_OUT | Active borrowing (short: Active) | Issued, still physical or replacement obligations; view/visit FSMO | Process authoritative return/acceptance; no separate Approved state |
| DENIED | Denied | Required visible reason, hold released, history; new request | Retained decision, no destructive delete/reapprove |
| CANCELLED | Cancelled | Own pending cancellation, hold released/history | Retained terminal request |
| EXPIRED | Expired |24h reservation deadline reached and server released; new request | Retained terminal expiry, not staff denial |
| COMPLETED | Completed | All physical/replacement resolved; history/final fine may still be Outstanding | Automatic after final valid resolution; no generic Complete button |

“Approved/released” can appear as one past event description, not persistent status. Staff detail may use “Checked out” beside Active for operational clarity. Preserve one label vocabulary across card/detail/queue/notification/report. Client expiry countdown0 before confirmed terminal response uses “Reservation time elapsed; checking status”, not EXPIRED or “stock released” guessed locally.

## Derived labels and badge hierarchy

| Label | Derivation / priority | Scope |
|---|---|---|
| Overdue | Active unresolved physical OR replacement and now>issued due_at; highest operational urgency | Borrower own / Staff/Admin; never pending terminal or currently Completed |
| Partial return | Some original dispositions processed and still P or U positive | Borrower own / Staff/Admin; secondary to Overdue |
| Replacement required | Sum required−accepted>0; show quantity per type/kind | Borrower own / Staff/Admin; physical C may0 |
| Expiring soon | Advisory time remaining, no new stored state | Pending only; final warning threshold is a UX copy/timing review, not new TTL |
| Due today | Due local date in Asia/Manila, contextual filter | Does not redefine exact overdue cutoff or calendar-day fine |
| Returned good / Recorded damaged / Recorded lost | Immutable disposition quantities | Incident still shown after replacement, never restyled as “good” |
| Replacements accepted / Remaining replacements | Immutable acceptance sum / required−accepted | Not a current physical bucket |

Card badge order: one primary lifecycle label; then Overdue if true, Replacement required if relevant, Partial return lower priority. At narrow width show lifecycle+highest urgency plus a readable obligations line; detail shows all relevant labels. Completed with fine due reads **Completed · Fine outstanding**, never current Overdue. A loan with P0/U3 reads Active · Replacement required and becomes Overdue after due; no invented checked-out3. User-facing P=“Physically outstanding”, U=“Replacements remaining”.

## Fine labels and examples

| Balance/context | Label and data | Interpretation |
|---|---|---|
| No assessed amount and no clearance | No overdue fine | Not an implicit payment |
| Positive outstanding on open loan | Fine outstanding; current assessed/previously cleared/current outstanding/as-of | Fine can rise until operational resolution |
| Zero outstanding with clearance on open loan | Fine cleared as of [time]; previous cleared/current assessment, live marker | Additional accrual possible; does not complete loan |
| Completed with positive outstanding | Completed + Fine outstanding; final assessed/finalization/remaining | Final amount frozen, Admin offline resolution still due |
| Completed with cleared fine | Fine cleared; final assessed plus clearance method/date/history | Assessed history preserved;0 balance not0 original fine |

Methods: PAID→Paid, WAIVED→Waived, OTHER_RESOLUTION→Other resolution. Admin Clear Fine records full current balance only. Staff/Borrower sees permitted read-only totals/resolution; no Admin-only controls. No paid-as-waived conflation, blanket “revenue”, automatic damage-price liability or gateway. PHP 10 times ceiling positive elapsed24h; final at completion. Explain live vs final instead of fine mutation implementation details.

## Physical stock / account / form state

Operational count labels preserve exact model: Available (`available`), Reserved (`reserved`), Checked out (`checked_out`), Damaged held (`damaged_held`, nonusable originals), Total tracked (`total_tracked`). `T=A+R+C+D`. Display definition: total tracked includes units with borrowers, not just shelves. Historical damaged/lost incidents and unresolved replacements are separate sections; never add them to T. Borrower catalog needs Available and lifecycle usability, not internal buckets.

Account Active / Inactive distinct from borrowing Active. Borrower category Student/Faculty never becomes a role label. Inventory Active / Inactive / Archived distinct from overdue/loan status. Inactive account blocks new login/borrower actions without erasing loans; Staff/Admin can process old returns/replacements.

Loading/Refreshing/Error/No results are view states, not empty stock or a borrowing lifecycle. Critical command “Recording…” → server-confirmed result; timeout→“Checking whether this was recorded” until reconciliation. Color/spinner/toast alone never communicates state. Announce meaningful updates politely once; do not announce countdown every second. Busy control has visible label and accessible state, with safe navigation/error recovery.
