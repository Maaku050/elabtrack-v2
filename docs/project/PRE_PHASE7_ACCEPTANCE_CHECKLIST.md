# Pre-Phase7 owner acceptance checklist

Prepared 2026-10-09. **HUMAN ACCEPTANCE: PENDING.** Engineering verification is recorded separately in [the readiness report](PRE_PHASE7_READINESS_REPORT.md). No item below has been confirmed by the owner. Phase7 remains unstarted and unauthorized.

Reviewer: ____________________  Date: ____________________  Environment: ____________________

For each row, enter **PASS / FAIL / NOT TESTED / NOT APPLICABLE** in the result space and record observations. Every blank result means NOT TESTED. A screenshot or automated test does not establish human acceptance.

## Startup and review data

From `/home/marvin/projects/eLabTrack_V2`:

```sh
docker compose up -d postgres
make migrate-status
make dev
```

Open **http://localhost:5173/login**. Check **http://localhost:5173/status** and **http://localhost:8080/api/v1/ready**. Ctrl+C stops the API/frontend; PostgreSQL and its volume stay running. The normal database is already migrated through000007; startup performs no migration or seeding. Do not run migration-down, reset, or seed to obtain a login.

| Source | Available review / limits |
|---|---|
| A — normal development data | Login, public status, anonymous guards, missing activation link, development visual previews. There are currently **zero accounts and zero equipment**. Authenticated review requires separately authorized initial Admin provisioning; there is no default password, public signup, or product Admin-creation endpoint. Do not enter real Student records during this review. |
| B — isolated development fixtures | Real authenticated Admin/Staff/Student/Faculty directory, XLSX, stock and image workflows; [Batch1 reproduction](../../integration/BATCH1.md). The previous private environment, random credentials, fixtures and listeners were cleaned up. Recreate a dedicated test environment before interactive review; never substitute the normal database for the harness target. |
| C — mock email delivery | The test-only API harness captures random activation links privately in the isolated environment. This supports password establishment/single-use/error review without sending mail. It is not a normal-runtime email option and does not prove Brevo delivery. Never paste captured tokens/passwords in this checklist. |
| D — approved visuals / prior engineering captures | [Binding PNG manifest](../ux/approved/MANIFEST.json), [Phase5 captures](../ux/verification/phase5/acceptance.json), [Phase6 captures](../ux/verification/phase6/acceptance.json). Actual development previews are listed below. They contain illustration fixtures and cannot validate account permissions, persistence or inventory mutations. |

Authenticated routes below require the indicated current account role. Replace `{id}` with a UUID obtained from the actual directory/catalog link, not an invented record. Admin inherits Staff operations; **the Admin borrower directory uses `/staff/borrowers`**.

| Flow | Actual local URL / control |
|---|---|
| Login / activation | `http://localhost:5173/login` / `http://localhost:5173/activate`; real activation requires a privately delivered fragment token |
| Admin/Staff borrower directory | `http://localhost:5173/staff/borrowers` |
| Individual Student / Faculty creation (Admin) | `http://localhost:5173/staff/borrowers/new`; select STUDENT or FACULTY on the same form |
| Student Excel bulk creation / deactivation (Admin) | `http://localhost:5173/admin/borrowers/bulk`; choose **Bulk Create** or **Bulk Deactivate** on the same page |
| Borrower account detail | `http://localhost:5173/staff/borrowers/{id}` |
| Restricted administrative directory | `http://localhost:5173/admin/administration/accounts`; initial privileged Admin creation is not available |
| Borrower equipment catalog / detail | `http://localhost:5173/borrower/equipment` / `http://localhost:5173/borrower/equipment/{id}`; existing current-terms gate applies |
| Operational inventory / categories / new equipment | `http://localhost:5173/staff/inventory` / `http://localhost:5173/staff/inventory/categories` / `http://localhost:5173/staff/inventory/new` |
| Operational equipment detail / edit / stock adjustment | `http://localhost:5173/staff/inventory/{id}` / `http://localhost:5173/staff/inventory/{id}/edit` / `http://localhost:5173/staff/inventory/{id}/adjust` |
| Admin quantity reconciliation | `http://localhost:5173/admin/inventory/{id}/reconcile` |
| Borrower first-use terms | `http://localhost:5173/borrower/terms`; no official terms have been published in normal data |
| Development visual previews | `http://localhost:5173/__preview/borrower/home`, `/__preview/borrower/equipment`, `/__preview/staff/dashboard`, `/__preview/staff/requests`, `/__preview/admin/dashboard` on the same origin; excluded from production |

## Authentication

| Manual check | Expected result | Result / notes |
|---|---|---|
| Open login; use Check connection on `/status` (A) | Email/separate-password form; actual API connected. No default account/password advertised. | ____________________ |
| Sign in with a private activated fixture; try incorrect credentials (B) | Correct role workspace; safe generic denial for invalid credentials; no institutional SSO/mailbox-password request. | ____________________ |
| Reload, logout, and check another tab (B) | Session restores or honestly fails; logout clears protected data across tabs; no access token persisted or exposed. | ____________________ |
| Deactivate a signed-in fixture and reopen protected content (B) | Backend refuses current inactive account; no stale-token bypass. | ____________________ |

## Roles/permissions

| Manual check | Expected result | Result / notes |
|---|---|---|
| Open management deep links while anonymous (A) | Redirect to login; protected APIs return401. | ____________________ |
| Attempt account creation/bulk/activation resend as Staff (B) | No privileged control; current backend denies direct mutation. Only Admin creates accounts. | ____________________ |
| Compare Student/Faculty, Staff and Admin accounts (B) | Student/Faculty both BORROWER; Staff/Admin remain separate authorization roles. | ____________________ |
| Inspect restricted administrative directory (B) | Staff provisioned separately by Admin; existing Admin rows read-only; no Admin promotion/creation/status writer. | ____________________ |

## Student provisioning

| Manual check | Expected result | Result / notes |
|---|---|---|
| Select STUDENT (B; fail-closed normal configuration after authorized Admin access) | Student ID/institutional email required. Absent approved domain configuration displays unavailable onboarding. | ____________________ |
| Create with a textual ID including leading zeroes; repeat ID/email (B) | Exact text preserved; unique identity enforced; conflicts do not overwrite another account. | ____________________ |
| Try a disallowed domain or blank ID (B) | Server rejects; no account is created from invalid inputs. | ____________________ |

## Faculty provisioning

| Manual check | Expected result | Result / notes |
|---|---|---|
| Select FACULTY on the same form (B) | No Student ID/institutional-domain requirement; valid unique external email accepted individually. | ____________________ |
| Inspect password and authorization controls (B/D) | No Admin-assigned password or Staff/Admin role choice on borrower form; account awaits secure activation. | ____________________ |

## Bulk Student creation

| Manual check | Expected result | Result / notes |
|---|---|---|
| Download template and upload an isolated synthetic `.xlsx` (B) | Exact `studentId,name,email,course,contactNumber`; first three required; Student IDs are text. | ____________________ |
| Include duplicate IDs/emails, missing IDs, disallowed domains, formulas or forbidden columns (B) | Complete validation preview or safe workbook rejection; no writes before confirmation; no Faculty/password/role columns. | ____________________ |
| Select valid rows, cancel, then confirm and retry (B) | Cancel does not create accounts; confirmed selection creates BORROWER/STUDENT only; conflict prevents partial commit; retry returns original result. | ____________________ |
| Inspect new accounts and mail behavior (B/C) | Activation pending; no automatic bulk email campaign. Open each account to send its activation link. | ____________________ |

## Bulk Student deactivation

| Manual check | Expected result | Result / notes |
|---|---|---|
| Reuse isolated Student roster; include unmatched/conflicting/already-inactive identities (B) | Every result shown; exact stable ID/email matching; no roster-absent deactivation. | ____________________ |
| Attempt Faculty/Staff/Admin targeting (B) | Excluded/rejected; workflow affects Students only. | ____________________ |
| Review warnings and cancel/confirm (B) | Explicit confirmation; obligations do not veto deactivation. Current obligations show **UNAVAILABLE**, never fabricated zero balances. | ____________________ |
| Inspect retained history and unresolved obligations (B; future integration where applicable) | No deletion, fine clearance/payment, returns, loan closure, replacement resolution, overdue or inventory change. Live borrowing/fine projections are not installed yet. | ____________________ |

## Activation states

| Manual check | Expected result | Result / notes |
|---|---|---|
| Open `/activate` without a link (A) | Honest missing-link help; no account mutation. | ____________________ |
| Deliver/capture a fixture link, set a separate password and reuse it (B/C) | Fragment removed from visible URL; borrower chooses password; single use; activation does not automatically authenticate or accept terms. | ____________________ |
| Review unavailable/failed/ambiguous delivery, expiry, resend and deactivation (B/C) | Truthful states; bounded resend; invalidated/expired token denied; reactivation does not revive old credentials. | ____________________ |

## Brevo external dependency

| Manual check | Expected result | Result / notes |
|---|---|---|
| Inspect configured provider state without sending mail (A after authorized Admin access/B) | Missing key/sender means UNCONFIGURED; startup still works. No key in browser/logs. | ____________________ |
| Live delivery verification | **NOT TESTED here.** Separately configure key, verified sender and permitted activation origin; authorize delivery testing before claiming receipt/ownership. Provider acceptance is not verified delivery. | ____________________ |

## Equipment catalog

| Manual check | Expected result | Result / notes |
|---|---|---|
| Browse/search/filter/sort/page as accepted Borrower (B) | Actual ACTIVE equipment only; truthful empty states; detail quantities agree with backend. Normal official terms remain unpublished. | ____________________ |
| Try inactive/archived equipment as Borrower (B) | Not disclosed as active catalog entries; guarded detail access. | ____________________ |
| Compare catalog to binding visuals (D/B) | Navy/violet shell preserved; View Details works; no pretend cart/request action. | ____________________ |

## Inventory

| Manual check | Expected result | Result / notes |
|---|---|---|
| Create/edit category and equipment as Staff/Admin (B) | Real persisted data; no institutional taxonomy seeds; validation and stale-version failures visible. | ____________________ |
| Review opening/add/remove stock with reason and confirmation (B) | A/T move together; R/C/D unchanged; immutable movement/audit visible; over-removal denied. | ____________________ |
| Toggle inactive/archive and inspect history (B) | Role-aware visibility; archive guards custody/holds and future liability boundary; history retained; no archived edits/unarchive. | ____________________ |

## Images

| Manual check | Expected result | Result / notes |
|---|---|---|
| Upload a small synthetic PNG/JPEG to equipment (B) | Authorized canonical PNG persists/displays; current image bound to its equipment. | ____________________ |
| Try invalid/oversized/dimension-exceeding files (B) | Safe rejection; max512KiB encoded,2048 per dimension/4M pixels; no arbitrary URL/SVG upload. | ____________________ |
| Inspect scope (B/D) | Images belong to equipment catalog only. No return evidence/upload/attachment feature; phone photos stay outside eLabTrack. | ____________________ |

## Quantity reconciliation

| Manual check | Expected result | Result / notes |
|---|---|---|
| Reconcile observed available quantity as Admin, including explicit zero (B) | Review/reason/confirmation required; `T=A+R+C+D`; R/C/D unchanged; persisted RECONCILE movement. | ____________________ |
| Attempt reconciliation as Staff; retry stale stock sequence (B) | Backend denies Staff; stale sequence conflicts; retry does not duplicate movement. | ____________________ |

## Responsive layout

| Manual check | Expected result | Result / notes |
|---|---|---|
| Check390/768/1024/1440 widths and long values (A/B/D) | Readable forms/cards; bounded table scrolling; no viewport overflow or inaccessible action. | ____________________ |
| Use keyboard and confirmation dialogs (A/B) | Visible focus, labels, error announcements, cancel/confirm access and restored focus. | ____________________ |

## Light/dark appearance

| Manual check | Expected result | Result / notes |
|---|---|---|
| Toggle appearance on login, directories, detail, bulk and catalog (A/B/D) | Approved color hierarchy and readable contrast; no lost statuses/focus/content. | ____________________ |
| Compare real fixture screens with approved PNGs (B/D) | Disclosed scope adaptations reviewed; no approval inferred from engineering screenshots. | ____________________ |

## Known unavailable features

| Manual check | Expected result | Result / notes |
|---|---|---|
| First-use terms in normal environment | Unpublished official terms remain fail-closed; do not publish synthetic terms there. Test synthetic acceptance only in isolated fixtures. | ____________________ |
| Inspect later workflows/placeholders | No borrowing/cart/reservation/checkout/return/replacement/fine settlement/report/notification/password-recovery implementation. No live obligation totals. | ____________________ |
| Initial Admin access and Student policy inputs | Separately authorized bootstrap procedure/private credentials and actual approved SKSU domains required for normal authenticated/Student review. | ____________________ |

## Phase7 prerequisites

| Owner decision / review | Expected result | Result / notes |
|---|---|---|
| Review this checklist and failed/untested items | Owner records actual disposition; engineering PASS does not mark human items PASS. | ____________________ |
| Review authentic activation/domain/official-terms dependencies | External gates remain explicit; official terms gate live borrowing, not independent development. | ____________________ |
| Review Phase7 implementation handoff | Future scope must connect transactional current-terms guard, real obligations/archive adapters and approved inventory privileges/invariants; no work starts under this checkpoint. | ____________________ |
| Authorize Phase7 separately | Explicit owner instruction required after human acceptance; no implied authorization. | ____________________ |

Acceptance decision: ____________________  Follow-up / blocker references: ____________________
