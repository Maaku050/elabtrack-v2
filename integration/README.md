# Disposable foundation integration verification

Phase 1H hardens migrations and verifies actual runtime/migrator separation on
PostgreSQL 18.6, the repository API/nginx images and Chromium. Historical Phase 1G
[evidence](evidence/2026-10-07.json) and its foundation report remain immutable;
old runner-defect expectations have been replaced by assertions of safe behavior.
No business schema or product UI is involved. Run from the repository root.

## Isolation and prerequisites

Use the declared Go toolchain, compatible Node/npm, `psql`, Docker and Compose
supporting `!override`/`!reset`. The integration target is ONLY
`elabtrack_v2_integration`, loopback PostgreSQL **15432**, API **18080**, nginx
**15173**. Inspect existing containers/volumes first. If any Phase 1H resource
already exists with unknown ownership, stop rather than reuse/delete its data.
The base manifest has fixed names, so `-p` alone is insufficient.

Generate ignored 0600 credentials; do not print or render the environment:

```bash
python3 - <<'PY'
from pathlib import Path
import secrets
p = Path('backend/.env.phase1h')
if p.exists():
    raise SystemExit('Review the existing integration file first')
p.write_text('BOOTSTRAP_DB_USER=postgres\nBOOTSTRAP_DB_PASSWORD='+secrets.token_hex(24)+
             '\nMIGRATION_DB_PASSWORD='+secrets.token_hex(24)+
             '\nDB_PASSWORD='+secrets.token_hex(24)+'\nJWT_SECRET='+secrets.token_hex(32)+
             '\nDB_NAME=elabtrack_v2_integration\nDB_USER=elabtrack_runtime\nDB_PORT=15432\n')
p.chmod(0o600)
PY
phase1h_compose() {
  docker compose -p elabtrack_v2_phase1h --env-file backend/.env.phase1h \
    -f docker-compose.yml -f integration/compose.phase1h.yml "$@"
}
phase1h_compose --profile full --profile tools config --quiet
phase1h_compose up -d postgres
```

Fresh-volume init runs `backend/database/bootstrap-roles.sh`: bootstrap admin
creates two non-superuser logins, owns the database/installs citext, transfers the
public schema to `elabtrack_migrator`, revokes public schema CREATE/database TEMP
and grants runtime schema USAGE/database CONNECT. Neither login gets cluster
administration or membership in the other. Current UUID keys need no sequences;
standard built-ins/citext retain their required ordinary type/function access.
No custom SECURITY DEFINER function is installed.

Wait for PostgreSQL health. Apply migrations explicitly as the owner, then grant
DML on ONLY the current application tables:

```bash
(cd backend && GOCACHE=/tmp/elabtrack-phase1b-go-cache go build -o /tmp/elabtrack-phase1h-api ./cmd/api)
(cd backend && python3 ../integration/env-run.py /tmp/elabtrack-phase1h-api --migrate-status)
(cd backend && python3 ../integration/env-run.py /tmp/elabtrack-phase1h-api --migrate-up)
# This session selects the migration role; the helper never puts passwords in argv.
python3 integration/env-run.py python3 integration/sql.py migration backend/database/runtime-grants.sql
(cd backend && python3 ../integration/env-run.py /tmp/elabtrack-phase1h-api --migrate-status)
```

For explicit zero/latest/zero/latest testing, run `--migrate-down` three times
before the final `--migrate-up`, then reapply runtime-grants.sql. Existing 000003
invalidates sessions in either direction; do this before starting the API or
creating authentication fixtures. No historical SQL is edited. Future approved
migrations grant DML on named new application tables inside their own transaction,
with only required sequence/function permissions. Never grant ALL TABLES/default
runtime table privileges: `schema_migrations` stays private to the migrator.

The helper constructs typed `MIGRATION_DATABASE_URL` from the ignored local
credentials and keeps `DB_USER=elabtrack_runtime` for API/repository tests. Runtime
API Compose configuration receives neither migration URL nor migration/bootstrap
passwords. Seeds select the explicit migration connection and retain the
before-IO development-only guard. Session cleanup selects runtime credentials.

## Runner contract and compatibility

Files are `NNNNNN_lowercase_name.up.sql` / `.down.sql`, nonzero six-digit numeric
versions, unique numeric mappings and required pairs. Gaps are allowed. Applied
history must be a prefix of the discovered files; missing historical files,
lower pending versions and mismatches fail. SHA-256 covers a versioned marker
plus big-endian uint64 byte length and raw bytes for up, then down. Line endings
and comments count. Applied files in both directions are immutable.

Each command discovers/loads its files, acquires a dedicated connection and
`pg_try_advisory_lock(0x454c41424d494752)` (decimal parameter passed by Go), then
reads tracking. Contention fails immediately; connection/lock-query acquisition
is bounded to five seconds. The connection closes on every exit, releasing the
session lock without returning a locked session to the pool. Status uses the
same lock for a coherent read, creates no tracking table and reports current
version/applied/pending/checksum verification. Query/scan/iteration failures fail
rather than becoming version zero.

Up: BEGIN → SQL → INSERT(version, checksum) → COMMIT. Down: BEGIN → SQL → DELETE
matching record → COMMIT. Failures roll back both. The lock covers the full
command; transactions commit per file, so earlier successful files remain
applied if a later file fails. A lost commit acknowledgement may be ambiguous;
inspect status before retry. No dirty flag is needed for supported transactional
SQL. Nontransactional statements fail inside PostgreSQL's transaction block;
transaction/session-control SQL is rejected lexically before execution. This
includes SET/RESET/prepared statements under the current convention. PostgreSQL
still parses SQL. Future exceptional support needs an explicit reviewed dirty/
recovery contract; none is implemented now.

Checksum-free history is NEVER silently upgraded. Prefer a fresh disposable
volume for integration. For a known, verified, pre-production local database,
review ownership, schema provenance and a backup first. The controlled
`APP_ENV=development --migrate-adopt-legacy` action attests that history: only a
prefix of the three exact Phase 1G foundation versions with baseline-matching
paired content is accepted. Unknown versions or altered baseline files fail.
One transaction adds checksum, backfills that verified mapping, and adds NOT
NULL/digest checks while preserving timestamps. Repeated adoption is refused.
This is an operator attestation, not proof of arbitrary historical DDL. Unknown
external databases require separate reconciliation, not fabricated checksums.
Existing volumes are not auto-re-owned or bootstrapped; any transfer of users,
refresh_tokens and tracking ownership to the migrator must be separately reviewed
by their administrator. The bootstrap admin remains responsible for extensions.

Production API requires only runtime connection settings. Production migration
CLI requires a separate MIGRATION_DATABASE_URL to the same host/port/database,
a different username, explicit verify-full TLS and no fallback. Put a private CA
path in that URL's sslrootcert if necessary. Do not inject this secret into the
API. Development/test explicitly allow absent migration URL to use the configured
local DB connection for convenience; fresh supported Compose uses separated roles.
Its `tools` profile is an explicit one-shot migrator command with the local owner
connection, and is excluded from ordinary/full API startup. Normal development
`make compose-migrate-up`, then `make compose-runtime-grants`, implements this flow.

## Real migration, privilege and auth tests

Run migration fault tests before API/auth/browser tests:

```bash
(cd backend && ELABTRACK_MIGRATION_INTEGRATION=1 python3 ../integration/env-run.py go test -race -v -count=1 ./internal/infrastructure/database)
python3 integration/env-run.py python3 integration/migration-processes.py
phase1h_compose --profile full build
phase1h_compose --profile full up -d
# Wait for http://localhost:15173/api/v1/ready to return 200.
DOCKER_BIN=docker python3 integration/env-run.py python3 integration/roles.py
(cd backend && python3 ../integration/env-run.py go test -race -v -count=1 ./tests/integration)
```

Tests prove partial SQL and bookkeeping failure rollback in both directions,
up/down checksum tampering, missing applied files, historical gaps, duplicates,
permission/query errors, unsupported nontransactional and transaction-control
SQL, simultaneous PG connections and separate CLI processes, effective GRANT
refusal, runtime DML/DDL/tracking denials, non-superuser owner DDL and explicit
legacy adoption. Temporary SQL/constraints/triggers/history/objects are cleaned;
no permanent broken migration or production failure injection exists.
Auth tests assert runtime identity and exercise real login/refresh/hash storage,
12-way single-use rotation, post-consumption rollback, stale account authority,
logout and bounded cleanup. Ordinary `go test ./...` skips live suites; it never
implicitly starts services or touches a configured production database.

## Browser, runtime and logs

Use an external Playwright installation/Chromium runtime, not an application
dependency. The recorded run reused Playwright 1.55.0/Chromium 140.0.7339.16;
missing Ubuntu libraries were previously extracted under /tmp. If needed set
PLAYWRIGHT_MODULE and LD_LIBRARY_PATH explicitly. Docker Desktop Windows
`docker.exe`/standalone `docker-compose.exe` through WSL interop are supported
fallbacks. Real local sockets/browser/interop required sandbox escalation in the
recorded environment; no toolchain requirement was weakened.

Inspect only nginx's network address, then append its observed `/32` as
INTEGRATION_TRUSTED_PROXY in the ignored file; recreate the backend and recheck
readiness. Never trust the whole private Docker network. Capture API logs before
recreation into `/tmp/elabtrack-phase1h-startup.log`.

```bash
docker inspect elabtrack_v2_phase1h_frontend --format '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}'
phase1h_compose --profile full up -d backend
node integration/build-browser.mjs
docker cp frontend/phase1g.local/dist/. elabtrack_v2_phase1h_frontend:/usr/share/nginx/html/
docker exec elabtrack_v2_phase1h_frontend nginx -t
python3 integration/env-run.py node integration/browser.mjs
DOCKER_BIN=docker python3 integration/env-run.py python3 integration/runtime.py
COMPOSE_BIN=/path/to/docker-compose python3 integration/env-run.py node integration/restart.mjs
docker logs elabtrack_v2_phase1h_backend > /tmp/elabtrack-phase1h-final.log 2>&1
python3 integration/logs.py
```

Run these sequentially: runtime.py stops PostgreSQL and fills rate budgets;
restart waits for actual API readiness rather than Docker's `Up` state. The
ignored `phase1g.local` driver and `/tmp/elabtrack-phase1g-*` browser/runtime/
restart/sentinel artifacts retain their original names to reuse the existing
harness; they now target the Phase 1H stack. Normal production HTML/source is
unchanged. The test build's static driver import causes an expected chunk-size
warning absent from the unchanged normal build. Logs are checked against current
run credential sentinels without displaying them; migration-processes.py checks
SQL/config/password/URL redaction independently.

**Phase 1I remains required and unimplemented:** racing same-cookie tabs can
produce one successful rotation and a late losing 401 which clears the winner's
new cookie. Peer logout also leaves another tab apparently authenticated until
refresh fails. No token broadcast, replay relaxation or client coordination is
added during 1H. The browser harness deliberately preserves this adverse test.
CI is still open. These deterministic commands plus standard fmt/vet/test,
targeted race, frontend lint/test/build, Compose config and diff checks are the
later CI input; local success is not a deployed production TLS/backup/operations
claim.

## Cleanup

After confirming this run owns these resources:

```bash
phase1h_compose --profile full --profile tools down -v --rmi local
rm backend/.env.phase1h /tmp/elabtrack-phase1g-secrets.json
rm -r frontend/phase1g.local
```

Remove only this project's containers/network/volume/built images and private
fixtures. Do not prune Docker globally or touch normal `elabtrack_v2_pgdata`,
unknown resources, V1 or production data. Shared base-image caches and sanitized
/tmp evidence can remain. [Phase 1H machine evidence](evidence/2026-10-07-phase1h.json)
and the foundation report record actual results and deployment-only limits.
