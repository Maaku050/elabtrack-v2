# Disposable foundation integration verification

Phase 1G uses real PostgreSQL, the repository API image, nginx and Chromium.
No product screens, test endpoints, production configuration or business tables
are added. Run from the repository root. These scripts deliberately target only
`elabtrack_v2_integration` on loopback port 15432 and the isolated API/frontend
ports 18080/15173. They do not consume an ambient production database URL.

## Prerequisites and isolation

Use the declared Go toolchain, compatible Node/npm, PostgreSQL `psql`, Docker and
Compose supporting `!override`/`!reset`. Supply an external Playwright installation
and its Chromium runtime. No browser dependency is added to the application.
The recorded run used cached Playwright 1.55.0/Chromium 140.0.7339.16. Missing
Ubuntu libraries were downloaded/extracted under `/tmp`, not installed globally.

The base manifest has fixed container/volume names: **`-p` alone is insufficient**.
Always combine it with `integration/compose.override.yml`. First inspect Docker
containers and volumes and ensure these integration names are unused. Stop if an
unknown existing resource uses them; do not reuse or delete its data.

Generate disposable credentials without printing them:

```bash
python3 - <<'PY'
from pathlib import Path
import secrets
p = Path('backend/.env.phase1g')
if p.exists():
    raise SystemExit('Review the existing local integration file first')
p.write_text('DB_PASSWORD='+secrets.token_hex(24)+'\nJWT_SECRET='+secrets.token_hex(32)+
             '\nDB_NAME=elabtrack_v2_integration\nDB_USER=postgres\nDB_PORT=15432\n')
p.chmod(0o600)
PY
```

The file is ignored and excluded from the backend Docker context. The override
resets backend `env_file`, explicitly selects local development, retains the
repository's `unless-stopped` policy, and changes only integration serving/abuse
settings. Limits remain 10/60/120 but use five-second windows for bounded tests.
`INTEGRATION_TRUSTED_PROXY` starts empty. Do not trust a whole Docker network.

For readability, create this shell function (or substitute equivalent explicit
arguments in every command):

```bash
phase1g_compose() {
  docker compose -p elabtrack_v2_phase1g --env-file backend/.env.phase1g \
    -f docker-compose.yml -f integration/compose.override.yml "$@"
}
phase1g_compose config --quiet
phase1g_compose up -d postgres
```

Wait for PostgreSQL health. Build the normal images, then apply migrations
**explicitly**, before starting the API/frontend:

```bash
phase1g_compose --profile full build
(cd backend && GOCACHE=/tmp/elabtrack-phase1b-go-cache go build -o /tmp/elabtrack-phase1g-api ./cmd/api)
(cd backend && python3 ../integration/env-run.py /tmp/elabtrack-phase1g-api --migrate-status)
(cd backend && python3 ../integration/env-run.py /tmp/elabtrack-phase1g-api --migrate-up)
(cd backend && python3 ../integration/env-run.py /tmp/elabtrack-phase1g-api --migrate-down)
(cd backend && python3 ../integration/env-run.py /tmp/elabtrack-phase1g-api --migrate-up)
(cd backend && python3 ../integration/env-run.py /tmp/elabtrack-phase1g-api --seed)
phase1g_compose --profile full up -d
phase1g_compose --profile full ps
docker exec elabtrack_v2_phase1g_frontend nginx -t
```

Wait for `/api/v1/ready` at `http://localhost:15173` to return 200. Compose restart
starts processes before dependency readiness; the existing restart policy may
retry the API's fail-closed connection startup. Never infer readiness from `Up`.
Migration and seed commands are never startup/restart side effects.

For the trusted-proxy experiment, inspect **only** nginx's network address:

```bash
docker inspect elabtrack_v2_phase1g_frontend --format '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}'
```

Append `INTEGRATION_TRUSTED_PROXY=<observed-IP>/32` to the ignored fixture, then
`phase1g_compose --profile full up -d backend`. Recheck the address after any nginx
container recreation. The recorded address was 172.18.0.4; it is not a deployment
recommendation or a portable value.

## Tests and browser seam

```bash
(cd backend && python3 ../integration/env-run.py go test -race -v -count=1 ./tests/integration)
node integration/build-browser.mjs
docker cp frontend/phase1g.local/dist/. elabtrack_v2_phase1g_frontend:/usr/share/nginx/html/
```

The normal production index/assets remain in nginx. The additional ignored build
serves `/phase1g.local/index.html`: it imports the existing app, client singleton
and store in one bundle. Its test bridge exists only in this disposable build.
The forced access denial changes test-page memory; bounded persistent-denial
cases mock only protected HTTP responses, while refresh always uses the real API
and PostgreSQL. The real normal SPA is also tested separately. The harness's
static import combines the lazy transport and emits an expected size warning;
the unchanged normal production build retains its existing lazy chunk.

Set `PLAYWRIGHT_MODULE` to the external Playwright `index.mjs`. If required, set
`LD_LIBRARY_PATH` to locally extracted browser libraries. Run these sequentially
because the runtime script stops PostgreSQL and fills rate budgets:

```bash
python3 integration/env-run.py node integration/browser.mjs
DOCKER_BIN=docker python3 integration/env-run.py python3 integration/runtime.py
COMPOSE_BIN=/path/to/docker-compose python3 integration/env-run.py node integration/restart.mjs
```

For restart, `COMPOSE_BIN` is a standalone Compose executable (not a shell string).
Docker Desktop's Windows `docker.exe`/`docker-compose.exe` also work through WSL
interop when native integration is unavailable; that was the recorded run's
fallback. Local socket/browser/Windows interop required sandbox escalation in
the verification environment. No requirements were downgraded.

`runtime.py` creates/drops a temporary runner probe table, constraint, tracking
row and restricted login exclusively in the disposable PostgreSQL instance; it
does not edit migrations or redesign the runner. It records real bookkeeping,
concurrency, edited-content and lookup-error failures. Do not run it against
development/production data. The restricted login is a fault fixture, not a
runtime/migration credential design.

Scripts write private test sentinels to `/tmp/elabtrack-phase1g-secrets.json` (0600)
and safe result JSONs under `/tmp`. At the beginning of a fresh verification run,
create an empty `/tmp/elabtrack-phase1g-final.log`. Capture the API log before any
API container recreation and again after all tests:

```bash
docker logs elabtrack_v2_phase1g_backend >> /tmp/elabtrack-phase1g-final.log 2>&1
python3 integration/logs.py
```

The checker deduplicates repeated captures of the same log line. Optional
`/tmp/elabtrack-phase1g-api.log` and `...-before-fix.log` preserve this run's older
segments for extra redaction checks; they are not prerequisites for a new run.
Do not display the sentinel file or environment values.
The checked-in [evidence summary](evidence/2026-10-07.json) contains no credentials.

Normal `go test ./...` skips the opt-in live suite. All standard gates remain
required; no mocks substitute for this suite or Chromium. Tests create synthetic
accounts and delete their own fixtures. Operational retention ownership,
production TLS/hosting and institutional account policy remain separate work.

## Cleanup

Only after confirming ownership of this integration project:

```bash
phase1g_compose --profile full down -v --rmi local
rm backend/.env.phase1g /tmp/elabtrack-phase1g-secrets.json
rm -r frontend/phase1g.local
```

This removes only the named integration containers/network/volume and generated
project images. Never run global Docker prune, remove the normal
`elabtrack_v2_pgdata` volume, or use these commands on another project. Shared
downloaded base-image caches can remain. No V1, production or school data is used.
