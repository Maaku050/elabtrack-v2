# Phase 4B isolated verification

Use only disposable `elabtrack_v2_phase4b_test` / container `elabtrack_v2_phase4b_postgres` / project `elabtrack_v2_phase4b` / volume `elabtrack_v2_phase4b_pgdata`, PostgreSQL port35432, API18084, Vite15174, compiled preview14174. Do not redirect these tests to normal `elabtrack_v2`, Phase4A or production accounts. The source helper refuses other identities and keeps privileged migration settings out of API startup. No official terms are approved; TEST documents are not institutional wording.

Run from repository root unless a block says otherwise. Requires configured Go toolchain, npm dependencies, Docker, Chromium and Linux browser libraries. The recorded run used Go1.27.1, Node24.19.0/npm11.17.0, PostgreSQL18.6 and Chromium140. Browser paths below are this workstation's cached tools; use installed equivalents elsewhere. Existing user dev processes on8080/5173 are preserved.

Create ignored random infrastructure credentials only when the file does not already exist:

```bash
python3 - <<'PY'
import os, secrets
from pathlib import Path
p=Path('backend/.env.phase4b')
values={'DB_NAME':'elabtrack_v2_phase4b_test','DB_PORT':'35432','BOOTSTRAP_DB_USER':'postgres',
        'BOOTSTRAP_DB_PASSWORD':secrets.token_urlsafe(48),'DB_PASSWORD':secrets.token_urlsafe(48),
        'MIGRATION_DB_PASSWORD':secrets.token_urlsafe(48),'JWT_SECRET':secrets.token_urlsafe(64)}
fd=os.open(p,os.O_CREAT|os.O_EXCL|os.O_WRONLY,0o600)
with os.fdopen(fd,'w') as f: f.write(''.join(k+'='+v+'\n' for k,v in values.items()))
PY
docker compose --env-file backend/.env.phase4b -p elabtrack_v2_phase4b -f docker-compose.yml -f integration/compose.phase4b.yml config --quiet
docker compose --env-file backend/.env.phase4b -p elabtrack_v2_phase4b -f docker-compose.yml -f integration/compose.phase4b.yml up -d postgres
```

Do not use the base `full` profile: its API/frontend container names belong to normal development. Wait for isolated PostgreSQL health. In `backend`, apply the existing hardened runner, then explicitly grant foundation DML (000005 itself grants narrow terms permissions):

```bash
python3 ../integration/phase4b-env-run.py go run ./cmd/api --migrate-up
python3 ../integration/phase4b-env-run.py go run ./cmd/api --migrate-status
```

From root, grant without putting passwords into argv/logs:

```bash
python3 - <<'PY'
import os, subprocess
from pathlib import Path
v=dict(l.split('=',1) for l in Path('backend/.env.phase4b').read_text().splitlines() if l and not l.startswith('#'))
assert v['DB_NAME']=='elabtrack_v2_phase4b_test' and v['DB_PORT']=='35432'
env=os.environ.copy();env['PGPASSWORD']=v['MIGRATION_DB_PASSWORD']
subprocess.run(['docker','exec','-i','--env','PGPASSWORD','elabtrack_v2_phase4b_postgres','psql','-X','-q','-v','ON_ERROR_STOP=1','-U','elabtrack_migrator','-d','elabtrack_v2_phase4b_test','-f','/opt/elabtrack/runtime-grants.sql'],env=env,check=True)
PY
```

In `backend`, run migration fault/race tests **before publishing synthetic history**. These suites intentionally refuse an existing terms history; do not delete immutable records to make a rerun pass. A whole fresh owned disposable database is required for a complete rerun.

```bash
python3 ../integration/phase4b-env-run.py go test -race -v -count=1 ./internal/infrastructure/database -run '^TestRealMigrator$'
python3 ../integration/phase4b-env-run.py go test -race -v -count=1 ./tests/integration -run '^TestRealTerms$'
go build -o /tmp/elabtrack-phase4b-api ./cmd/api
python3 ../integration/phase4b-env-run.py /tmp/elabtrack-phase4b-api
```

Keep the API running in that terminal. In another `backend` terminal:

```bash
python3 ../integration/phase4b-env-run.py go test -race -v -count=1 ./tests/integration -run '^TestReal(Foundation|TermsHTTP)$'
PHASE4B_BROWSER_FIXTURES=1 python3 ../integration/phase4b-env-run.py go test -v -count=1 ./tests/integration -run '^TestPhase4BBrowserFixtures$'
```

Browser fixtures use random bcrypt credentials saved only in mode0600 `/tmp/elabtrack-phase4b-fixtures.json`, created with O_EXCL. No product provisioning endpoint or published passwords. The narrow helper may disable/re-enable/revoke/change only exact owned fixture IDs/emails. `reset-current` verifies isolated identity and synthetic-only history, clears just the current pointer, and preserves historical versions/receipts to simulate missing content; no product withdrawal feature exists.

From root:

```bash
python3 integration/phase4b-fixture-state.py reset-current
```

In `frontend`, start Vite:

```bash
VITE_API_URL=http://localhost:18084/api/v1 npm run dev -- --host 127.0.0.1 --port 15174 --strictPort
```

Then run actual login/forms/terms requests against the real API/PG:

```bash
CHROMIUM_EXECUTABLE=/home/marvin/.cache/ms-playwright/chromium-1187/chrome-linux/chrome LD_LIBRARY_PATH=/tmp/elabtrack-browser-libs/root/usr/lib/x86_64-linux-gnu node scripts/phase4b-browser-qa.mjs
```

The final run records21 checks/26 screenshots, including long-document bottom, both themes/four widths, missing/required/updated/accepted/stale/network states, three tabs, original receipt preservation, inactive accounts and Staff/Admin. It never fabricates authenticated API responses. Node's isolated Admin login holds its JWT only in process memory for synthetic publication; terms read/acceptance occurs through the app's real feature hook/Query/API/transport. Original AuthProvider/Web Locks/BroadcastChannel remain the only product session architecture. Fresh Borrower fixtures are required when rerunning first-use cases because acceptance history remains immutable.

Build and start a compiled preview in another `frontend` terminal, then run visual regressions:

```bash
npm run build
npm run preview -- --host 127.0.0.1 --port 14174 --strictPort
```

```bash
CHROMIUM_EXECUTABLE=/home/marvin/.cache/ms-playwright/chromium-1187/chrome-linux/chrome LD_LIBRARY_PATH=/tmp/elabtrack-browser-libs/root/usr/lib/x86_64-linux-gnu node scripts/phase3b-browser-qa.mjs --url http://localhost:15174 --production-url http://localhost:14174 --output ../docs/ux/verification/phase4b/preview-regression
```

That inherited preview-only harness stubs anonymous refresh401 to inspect24 fixture layouts and production exclusion; it supplies no authenticated state or terms acceptance. Its evidence is distinct from the real Phase4B browser suite. Full quality commands remain in AGENTS.md: Go fmt/vet/test, targeted race, frontend lint/test/build and diff check. Default Go tests deliberately skip opt-in external fixture suites; the explicit commands above are the real database/HTTP verification.

Stop only the API/Vite/preview processes started for this suite, using Ctrl+C in their terminals. After verifying the container and volume labels belong to project `elabtrack_v2_phase4b`, discard the **whole owned disposable fixture database**; never use this command for normal/Phase4A projects:

```bash
docker compose --env-file backend/.env.phase4b -p elabtrack_v2_phase4b -f docker-compose.yml -f integration/compose.phase4b.yml down -v
```

Remove only this suite's private `/tmp/elabtrack-phase4b-fixtures.json` and any explicitly generated retry fixture files. Do not DELETE acceptance/version history or reset normal application accounts. Keep ignored infrastructure config mode0600 for reproducibility, or explicitly remove it after disposing its volume. Normal 000005 contains empty terms tables/a null pointer and no fake content; rollback refuses any published history. Existing normal dev API processes may need `make dev` restarted to load source changes; they are not silently stopped by verification. No production HTTPS/TLS/proxy, deployment or institutional publication is tested here.
