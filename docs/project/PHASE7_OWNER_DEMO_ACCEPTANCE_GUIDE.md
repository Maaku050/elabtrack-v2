# Phase7 owner demonstration and acceptance guide

**Formal owner acceptance, 2026-10-10 — DEC-081:** Phase7 Borrowing & Reservations is OWNER ACCEPTED following manual core-workflow review and automated safeguards. Student/Faculty mobile UI is accepted as functional Phase7, not final presentation; DEC-080 governs the future Phase11 browse/cart. Only a safe LOCAL checkpoint of verified Phase7/demo tooling/UX documents is authorized now. Do not begin Phases8–14, push or deploy. Official FSMO terms/current consent, approved production Student domains and required activation readiness remain live-use gates. Earlier pending-acceptance/correction statements below are historical. See [the checkpoint report](PHASE7_OWNER_ACCEPTANCE_CHECKPOINT.md).

2026-10-10. This is the isolated Phase7 demo at committed baseline `73f3fe4`. Owner acceptance remains pending. The normal development database is not a scenario/reset target. Returns, replacements and fines are outside this demo; Phases8–14 are paused by the latest owner correction.

## Start and retrieve private credentials

Run in your own WSL terminal from `/home/marvin/projects/eLabTrack_V2`:

```bash
python3 integration/phase7-demo.py setup
python3 integration/phase7-demo.py start
python3 integration/phase7-demo.py credentials
```

Open **http://127.0.0.1:15176/login**. The footer explicitly identifies **ISOLATED PHASE7 DEMO · FICTIONAL DATA**. Normal development remains on5173/8080/5434; demo frontend/API/PostgreSQL are15176/18086/54835. Do not use the normal frontend for these scenarios. Open the demo in a private browser window or a separate browser profile. The demo uses `127.0.0.1` while normal development uses `localhost`, because port numbers alone do not isolate the existing host-only refresh cookie. If you access normal development using `127.0.0.1`, the separate browser profile/private window is required; do not mix their sessions in one profile.

Setup is explicitly invoked and refuses duplicate seeding. Start reuses the prepared local API binary and production frontend assets, checks health and does not migrate/seed/reset normal data. Credentials are generated separately for five accounts, stored only in ignored0700 directories/0600 files and displayed only by the credentials command in an interactive terminal. Do not paste them into chat, documents or logs. Reset generates new passwords; retrieve them again afterward.

| Login key | Fictional email | Role / category |
|---|---|---|
|student|student.one@students.example.invalid|BORROWER / STUDENT; textual ID DEMO-001|
|student2|student.two@students.example.invalid|BORROWER / STUDENT; textual ID DEMO-002|
|faculty|faculty.one@example.invalid|BORROWER / FACULTY; no Student ID|
|staff|staff.one@example.invalid|STAFF|
|admin|admin.one@example.invalid|ADMIN|

All addresses are non-deliverable `.invalid` examples. Students use the demo-only `students.example.invalid` domain. Each login uses email plus its separate eLabTrack password. No SSO, role impersonation or existing normal password is involved. Local trusted initialization makes these fictional accounts ready without pretending email was delivered. Normal activation enforcement remains unchanged, and future newly provisioned demo accounts still follow the existing activation workflow.

## First-use demonstration terms

A fresh setup/reset publishes version **DEMO-1**, titled **DEMONSTRATION TERMS — NOT OFFICIAL FSMO POLICY**. No acceptance is seeded. Log in separately as each Borrower, read the displayed synthetic policy, select the initially unchecked consent checkbox and press **Accept and continue**. The real API stores the account/version/server-time receipt. Staff/Admin do not accept Borrower terms for another account. No unofficial policy or consent is inserted into normal data.

The automated version-change test publishes the next DEMO version (DEMO-2 from a fresh dataset) only in its disposable test scenario, confirms stale/missing-current acceptance rejection, then records actual UI consent. The owner-ready reset returns to DEMO-1 with zero acceptances. Official institutional approval remains external.

## Scenario A — Student request and immediate reservation

1. Log in as **student**, accept DEMO-1 through the actual form, and browse Equipment.
2. Find **DEMO Culinary Ladle**, open its details and select **Request Equipment**.
3. Set quantity2, review the request and confirm once.
4. Verify **Pending**, the SUBMITTED event and quantity2 reserved in the borrowing detail.
5. Log out and log in as Staff; inspect Inventory and the same request.

On a fresh demo, Ladle changes from **A20/R0/C0/D0/T20** to **A18/R2/C0/D0/T20**. A=available, R=reserved, C=checked out, D=damaged held, T=total tracked. The other equipment pools are independent. The request expires24 elapsed hours after actual submission. Repeat tests change the starting values; compare actual before/after, not just these first-run examples.

Faculty uses the same request workflow after its own first-use consent. Student/Faculty are both BORROWER; neither is Staff. The catalog includes multiple categories, six active positive-stock pools, an active zero-stock Rolling Pin and an inactive Grater. The inactive Grater's physical19 units remain tracked but cannot be newly borrowed.

## Scenario B — cancellation

As the owning Student, open the pending request, select **Cancel Pending Request**, review and confirm. Verify **Cancelled**, retained SUBMITTED/CANCELLED history and reserved quantity0. The preceding Ladle request returns to **A20/R0/C0/D0/T20**. Refresh/reopen the row: stock must not release twice. An issued loan has no cancellation action.

## Scenario C — Staff denial

As Student, create another pending request for3 Ladles. Log out, log in as Staff, open it under Requests & Borrowings and choose **Deny Request**. An explanation is mandatory; record a clearly fictional operational reason, review and confirm. Verify Denied and stock release. Log back in as that Student and verify the same visible reason and retained history. A fresh starting20 pool goes A17/R3 → A20/R0.

## Scenario D — physical checkout

Create a pending request for2 Ladles as Student. As Staff, open it, set a future date **and time** in Asia/Manila, confirm actual fictional physical handover, review and choose **Confirm Physical Checkout**. Verify **Checked out**, issued2, reserved0, recorded handover/due timestamps and history. There is no intermediate Approved state. A fresh starting20 pool becomes **A18/R0/C2/D0/T20**. Choose a future due time; there is no hardcoded seven-day maximum.

## Scenario E — Direct Issuance

As Staff, open Requests & Borrowings → **Direct Issuance**. Search for the fictional Student or Faculty email, select that Borrower, select active stocked equipment, enter quantity and a future Manila due timestamp, confirm physical handover, review and confirm. The target must already have accepted current demo terms. Verify DIRECT entry, Checked out, issued quantity, no pending expiry and exact available-to-custody movement.

After Scenario D with two Ladles issued, directly issue one additional Ladle: **A17/R0/C3/D0/T20**. Inspect Inventory detail and retained CHECKOUT/DIRECT movements. Stock remains authoritative after refresh. Separate issued loans remain open because Phase8 returns are deliberately not included.

## Scenario F — negative eligibility and authority

Use only these fictional demo records. Student2 can demonstrate missing consent before its first acceptance. Direct issue must reject that target. As demo Admin, Student2 can be deactivated through the normal account UI; a new login/request/direct issue must then fail until reactivated. An existing signed-in token does not override current account status.

Try inactive Grater, zero-stock Rolling Pin, a quantity above available, missing handover confirmation and a nonfuture due timestamp. The API must reject them without changing inventory. As a Borrower, navigate to `/staff/requests/direct`: **Access denied**. Borrowers cannot invoke Staff actions; Staff cannot create accounts. Automated checks also verify pending activation using temporary flags solely in the disposable demo fixture, restoring them afterward. Never toggle the normal pending accounts.

## Scenario G — deterministic expiry without waiting a day

First accept current demo terms as **student** and ensure one Ladle is available. In the same WSL terminal:

```bash
python3 integration/phase7-demo.py expire
```

This guarded local tool uses a test-only clock for a new fictional request created25 hours in the past, through the existing submit/reservation/consent services. It then invokes the real persisted-deadline sweep twice. It neither edits an existing loan's dates nor adds a backend clock override/API. The output gives a demo detail URL; open it as Student One, or filter My Borrowings by Expired. Verify EXPIRED, reserved0, automatic actor and retained history. Net stock must be unchanged after the one reserve/release. Insufficient stock or unaccepted terms makes the command fail safely.

## Reset and shutdown

For a new clean rehearsal, explicitly confirm this exact disposable target:

```bash
python3 integration/phase7-demo.py reset --confirm elabtrack_v2_phase7_demo
python3 integration/phase7-demo.py start
python3 integration/phase7-demo.py credentials
python3 integration/phase7-demo.py status
```

Reset verifies configuration, actual container/loopback port/volume/database and a database identity marker. It stops only owned demo processes, retains the old volume/private recovery configuration, creates a new distinct demo volume, applies only000001–000008 and atomically initializes fresh fixtures. It does not erase/reuse a normal or prior volume. Wrong confirmation, normal database/volume or wrong marker refuses before resetting. Accounts/equipment are not duplicated by repeating setup. Reset removes current rehearsal loans/consent from the active demo by switching to the fresh dataset; prior scenario data remains in its retained volume. Retrieve the newly generated credentials and reload the browser after reset; old sessions/passwords belong to the retained dataset.

Stop the application without deleting data:

```bash
python3 integration/phase7-demo.py stop
```

PostgreSQL remains available and its volumes remain. Restart with start. Setup/build needs preinstalled Go1.27.1-compatible tooling/modules, Node supported by the manifest, installed frontend dependencies, Docker Compose with override support and cached postgres:18.6-alpine. The initial tool deliberately snapshots73f3fe4 so paused future migrations cannot leak into this demo. Once prepared, the API/production assets/catalog illustrations/login image/fonts run locally without external authentication, CDN or email. No Docker/npm/Go download is required for ordinary start/stop. The reset rebuild was also tested with Go module downloads disabled.

## Troubleshooting and acceptance

If login fails, confirm15176 and use the latest credentials command; normal Admin credentials do not apply. After reset, refresh the page or sign out before signing in again. If borrowing redirects to terms, explicitly accept current DEMO version as that Borrower. Newly created pending accounts are not silently activated. If a port is occupied, stop only the process you identify as the demo; the controller will not stop unexpected processes. Private logs are under `backend/tmp/phase7-owner-demo/`; do not share credentials/config files. A failed setup does not authorize resets of other targets.

The existing normal development records remain unchanged; read the [engineering report](PHASE7_OWNER_DEMO_IMPLEMENTATION_REPORT.md) for executed checks and evidence. Review Student/Faculty login and terms, request/reservation, cancellation, denial, physical checkout, direct issue, history/expiry, stock and role boundaries in this demo. **The owner has now formally granted Phase7 acceptance (DEC-081). Phases8–14 still require separate implementation authorization; this acceptance permits only the safe local checkpoint.**
