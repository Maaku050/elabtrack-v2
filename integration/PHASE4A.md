# Phase 4A isolated verification

This runs the actual migration runner, restricted-runtime PostgreSQL repositories, HTTP API and Chromium login forms. No registration endpoint, published password or permanent local account is needed. The normal `elabtrack_v2` database and its volume are preserved. The older registration-based browser/local-readiness runners are historical and are superseded for current product authentication by this guide.

Requirements: declared Go/Node tools, Docker/Compose supporting `!override`, Python 3 and headless Chromium with its shared libraries. Run from the repository root. Stop other servers on 8080/5173 before the isolated host run. Inspect containers/volumes/ports first; refuse unknown existing Phase 4A resources. The test identity is exactly `elabtrack_v2_phase4a_test`, PostgreSQL 25432, container `elabtrack_v2_phase4a_postgres`, volume `elabtrack_v2_phase4a_pgdata`.

Generate only a missing ignored mode-0600 environment; existing credentials are never overwritten:

```sh
python3 - <<'PY'
from pathlib import Path
import os, secrets
p = Path('backend/.env.phase4a')
values = dict(DB_NAME='elabtrack_v2_phase4a_test', DB_PORT='25432', BOOTSTRAP_DB_USER='postgres',
              BOOTSTRAP_DB_PASSWORD=secrets.token_hex(32), DB_PASSWORD=secrets.token_hex(32),
              MIGRATION_DB_PASSWORD=secrets.token_hex(32), JWT_SECRET=secrets.token_hex(32))
fd = os.open(p, os.O_WRONLY | os.O_CREAT | os.O_EXCL, 0o600)
with os.fdopen(fd, 'w') as f:
    f.write(''.join(f'{key}={value}\n' for key, value in values.items()))
PY
phase4a_compose() {
  docker compose --env-file backend/.env.phase4a -p elabtrack_v2_phase4a \
    -f docker-compose.yml -f integration/compose.phase4a.yml "$@"
}
phase4a_compose config --quiet
phase4a_compose up -d postgres
phase4a_compose ps postgres
```

Wait for healthy PostgreSQL. The fresh volume creates separate non-superuser migrator/runtime roles using the existing bootstrap script. Run the role migration test **before any browser accounts or API process**; it refuses to replay rollback on a nonempty disposable schema. It copies the unchanged three foundation pairs to a temporary fixture directory, applies them with the real runner, grants only the named application-table DML, inserts two synthetic legacy accounts/session records, verifies exact mapping and all non-role data, reverses/reapplies, checks invalid-role/runtime-DDL denials and Staff rollback refusal, then removes only its own fixtures. The final schema is version 000004.

```sh
(cd backend && python3 ../integration/phase4a-env-run.py go test -v -count=1 \
  ./internal/infrastructure/database -run '^TestRealProductRoleMigration$')
(cd backend && go build -o /tmp/elabtrack-phase4a-api ./cmd/api)
# Terminal 1: restricted API; no migrator/bootstrap settings are inherited.
(cd backend && python3 ../integration/phase4a-env-run.py /tmp/elabtrack-phase4a-api)
# Terminal 2: normal Vite frontend, with localhost API URL from frontend/.env.
(cd frontend && npm run dev -- --strictPort)
```

The environment helper pins the test database/port, removes ambient PG*/URL overrides, uses explicit APP_ENV=test (no local dotenv fallback), and exposes migration credentials only to explicit migration/test commands. API and fixture mutations use `elabtrack_runtime`. The isolated HTTP rate window is five seconds; the normal login/refresh/general maxima and enforcement are retained. Local production defaults are unchanged.

Run security tests, then create browser fixtures explicitly. `TestRealFoundation` uses actual HTTP on 8080 and PostgreSQL 25432, checks all three roles/inactive/missing/stale JWT/demotion/disabled restoration, single-use rotation/replay, one concurrent winner/eleven safe denials, logout and transaction rollback. The browser fixture generator uses cryptographically random passwords and bcrypt; it refuses to overwrite its private file. Never display/copy that file.

```sh
(cd backend && python3 ../integration/phase4a-env-run.py go test -race -v -count=1 \
  ./tests/integration -run '^TestRealFoundation$')
(cd backend && PHASE4A_BROWSER_FIXTURES=1 python3 ../integration/phase4a-env-run.py \
  go test -v -count=1 ./tests/integration -run '^TestPhase4ABrowserFixtures$')
# Chromium executable/library paths depend on the host installation.
(cd frontend && CHROMIUM_EXECUTABLE=/path/to/chromium node scripts/phase4a-browser-qa.mjs)
# Always run this after browser success OR failure. It validates the isolated
# target, UUID and synthetic email, and deletes exactly the private file's IDs.
python3 integration/phase4a-fixture-state.py cleanup
# Existing hardened-runner regression, including atomic failures, checksum
# mismatch, advisory contention, privilege denials and explicit legacy adoption:
(cd backend && python3 ../integration/phase4a-env-run.py go test -race -v -count=1 \
  ./internal/infrastructure/database -run '^TestRealMigrator$')
```

The browser uses real user input, the real API and unchanged replies. Fault cases hold an actual refresh request, toggle Chromium offline or block only the API URL; they never fabricate successful auth responses. It validates login/logout/restore/deep links, all role shells, forbidden access, SQL changes confined to owned synthetic IDs, post-issuance disable/demotion, native peer logout, storage hygiene, keyboard access and 320/390/1280 Borrower plus 1024/1440 Staff/Admin light/dark. It exports screenshots and safe status/geometry/lifecycle metadata only. Passwords/access/cookies never enter artifacts.

The captured browser result is in `docs/ux/verification/phase4a/RESULTS.json`. Tests fail on unexpected runtime console errors, new warnings, overflow, undersized targets or missing authority. Network 401/403/404/offline responses are expected adverse cases, separate from application exceptions.

For approved preview regression (anonymous refresh interception, separately labeled), build the normal production bundle and optionally serve it on 4173 in another terminal:

```sh
(cd frontend && npm run build)
(cd frontend && npm run preview -- --host 127.0.0.1 --strictPort)
(cd frontend && CHROMIUM_EXECUTABLE=/path/to/chromium node scripts/phase3b-browser-qa.mjs \
  --production-url http://127.0.0.1:4173 --output /tmp/elabtrack-phase4a-preview-regression)
```

Stop host services with Ctrl+C. Stop only the isolated database with `phase4a_compose stop postgres`; retain the isolated volume unless its disposal is explicitly desired. Never run `down -v` against the normal development project. Re-running the role-migration test requires an empty owned disposable schema; if a browser run failed, clean exactly its fixtures first. Normal local development resumes with `make dev` and the preserved local `.env` files. Existing local role migration is an explicit `make migrate-up`, followed by `make migrate-status`; 000004 requires no new grants and preserves sessions.
