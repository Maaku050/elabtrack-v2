# Phase7 owner demo browser evidence

2026-10-10. Real Chromium, actual Phase7 frontend/API/PostgreSQL and isolated fictional accounts. These are historical engineering runs made before owner acceptance. Formal Phase7 acceptance is now confirmed in DEC-081; see [the checkpoint report](../../../project/PHASE7_OWNER_ACCEPTANCE_CHECKPOINT.md). All successful final runs block non-demo HTTP origins; no external requests or runtime exceptions are recorded. Browser credentials/session headers are never included in JSON or screenshots.

| Group | Checks | Successful screenshots | Scope |
|---|---:|---:|---|
|[borrower](borrower-acceptance.json)|6|14|Unchecked consent, catalog/detail, quantity/review, real reservation, cancellation/history;390/768/1024/1440 widths, light/dark|
|[staff](staff-acceptance.json)|8|18|Denial/reason/release, checkout/due/handover, direct issuance, A17/R0/C3/T20 ledger, Admin/Borrower history and actual expiry|
|[identities](identities-acceptance.json)|5|7|Student2/Faculty actual consent, Staff/Admin login, Student reacceptance and protected navigation|
|[owner-ready](owner-ready-acceptance.json)|5|5|Final reset's five fresh logins, decoded local images, visible demo identity and zero borrower consent|
|Total|24|44|Engineering evidence, not owner acceptance|

Representative views: [mobile dark catalog](borrower-catalog-390-dark.png), [mobile request review](borrower-review-390-light.png), [desktop issued quantities/history](staff-issued-loan-1440-light.png), [dark direct review](staff-direct-review-1024-dark.png), [mobile expiry](staff-borrower-expired-history-390-dark.png), [fresh Faculty terms](owner-ready-faculty-fresh-390-dark.png).

Three `*-failure-1440-light.png` files retain earlier harness diagnostics (borrower checkbox/role-switch/asset assertions). They are deliberately not counted as successful evidence. Failed runs and corrections are disclosed in the [implementation report](../../../project/PHASE7_OWNER_DEMO_IMPLEMENTATION_REPORT.md). Successful JSON lists exactly which screenshots were captured by each final run. Screenshots represent synthetic test histories; those histories remain in retained demo volumes after the final clean owner reset.

From the repository root, the repeatable acceptance sequence on a freshly reset demo is:

```bash
python3 integration/phase7-demo.py reset --confirm elabtrack_v2_phase7_demo
python3 integration/phase7-demo.py start
cd frontend
node scripts/phase7-owner-demo-qa.mjs borrower
node scripts/phase7-owner-demo-qa.mjs staff
cd ..
python3 integration/phase7-demo/verify_api.py
cd frontend
node scripts/phase7-owner-demo-qa.mjs identities
cd ..
python3 integration/phase7-demo.py reset --confirm elabtrack_v2_phase7_demo
python3 integration/phase7-demo.py start
cd frontend
node scripts/phase7-owner-demo-qa.mjs owner-ready
```

These tests mutate only the guarded disposable demonstration: do not run them against a manually reviewed demo you want to preserve without first retaining/resetting it. Private fixture data must exist; do not copy passwords into commands. The browser helper requires an installed Chromium (`CHROMIUM_EXECUTABLE`) and platform libraries; see existing browser tooling. The final group deliberately does not accept terms or submit requests.
