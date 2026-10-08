# Local full-stack closeout evidence

2026-10-08. Phase 3B closeout; [readiness report](../../../project/LOCAL_ENVIRONMENT_READINESS_REPORT.md). The owner reviewed real implementation screenshots and accepted the current baseline in DEC-064. Original [Phase 3B evidence](../phase3b/README.md) remains historical and unchanged.

- `LIVE_RESULTS.json`: **40 checks and 24 viewport/theme cases** using the real API on localhost:8080 and Vite on localhost:5173, without API interception. Anonymous, valid, revoked and expired session handling; simultaneous two-tab restoration and peer logout; direct credentialed readiness/request IDs; memory-only access and HttpOnly refresh-cookie policy; exact Origin rejection; mobile navigation and keyboard table scrolling. No unexpected JavaScript exceptions or browser warnings.
- `OUTAGE.json`: real API stopped, actual connection refusal and recovery banner. The live run then confirms the banner disappears after connectivity returns.
- `HEALTH.json`: actual liveness/readiness 200 and exact localhost credentialed CORS/request-ID headers.
- `MIGRATIONS.json`: existing runner verifies all three paired checksums; actual held advisory lock excludes a second runner, then status succeeds after release.
- `PRIVILEGES.json` / `DATABASE.json`: runtime DDL and migration-history access denied; named DML grants without grant option; migrator/runtime non-superuser attributes, citext, final fixture counts zero and healthy PostgreSQL on host port 5434.
- `SHUTDOWN.json`: actual make-dev descendant groups before Ctrl+C; no remaining group members, no listeners on 8080/5173/4173 and PostgreSQL still healthy afterward. API logged graceful pool closure. GNU Make emitted interruption/termination diagnostics, including `wait: No child processes`, during the intentional stop; no orphan or shutdown timeout resulted.
- `REFERENCE_QA_RESULTS.json`: **4 light baseline cases and 10 keyboard/boundary/production checks** from the existing visual runner. This separate run uses its documented anonymous-refresh 401 interception; it is not the live API/session evidence. Production JS excludes previews/fixtures, preview URL shows the existing 404, and foundation root remains available.
- Four selected live screenshots: Borrower Home 390 light, Catalog 390 dark, Staff Dashboard 1440 light, Pending Requests 1024 dark. Each was opened and visually inspected. Other live viewport captures remain in `/tmp/elabtrack-closeout-live`; their results are retained in JSON rather than duplicating the original 24-case reference matrix.
- `VALIDATION.json`: actual command outcomes and tools; `WORKTREE.json`: exact file inventory, protected scope, ignored-secret checks and redundant ZIP comparison.

Reproduce from the repository root with Docker/PostgreSQL ready, ignored local configuration restored, local `psql`, Python 3, Node and Chromium available:

```sh
make migrate-status
make dev
# In a second terminal:
CHROMIUM_EXECUTABLE=/path/to/chromium node frontend/scripts/local-environment-smoke.mjs
```

The runner creates one synthetic unprivileged development-registration fixture, expires only that fixture's active session and cleans only that fixture in `finally`. The SQL helper verifies development/runtime target, UUID, role and synthetic email prefix; it never targets existing accounts or business data. Credentials/cookies/access tokens remain in memory or ignored environment files and are not exported. The test refuses other configured database names/roles. Default output is `/tmp/elabtrack-local-environment-smoke`; `--output` selects a directory.

Actual browser: Chromium 140.0.7339.16, using the existing cached executable. Missing shared libraries were provided from `/tmp/elabtrack-browser-libs/root/usr/lib/x86_64-linux-gnu` via `LD_LIBRARY_PATH`; no system daemon/package or application dependency was installed. The final runner waits for a new document before reading bootstrap state; an intermediate stale-document timing failure was corrected in the verification utility, with application authentication code unchanged.

To reproduce the separate production check, build and serve `npm run preview -- --host 127.0.0.1 --port 4173 --strictPort` from frontend, then run its existing QA with `--url http://localhost:5173 --production-url http://127.0.0.1:4173 --quick --output /tmp/elabtrack-reference-qa`. Stop both development targets with Ctrl+C; the database stays running.

This is local readiness evidence. Production TLS/deployment, the separate isolated destructive integration harness, physical devices and all future feature states were not exercised here.
