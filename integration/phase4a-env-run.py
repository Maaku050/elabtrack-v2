#!/usr/bin/env python3
"""Execute verification against only the named disposable Phase 4A database."""
import os
from pathlib import Path
import sys
from urllib.parse import quote

root = Path(__file__).resolve().parent.parent
path = root / 'backend/.env.phase4a'
if path.stat().st_mode & 0o077:
    raise SystemExit('Integration environment must have mode 0600')
values = dict(line.split('=', 1) for line in path.read_text().splitlines() if line and not line.startswith('#'))
if values.get('DB_NAME') != 'elabtrack_v2_phase4a_test' or values.get('DB_PORT') != '25432':
    raise SystemExit('Refusing a non-Phase-4A target')
env = {key: value for key, value in os.environ.items() if not key.startswith('PG') and key not in ('DATABASE_URL', 'MIGRATION_DATABASE_URL', 'PORT')}
env.update(values)
env.update(APP_ENV='test', APP_PORT='8080', DB_HOST='127.0.0.1', DB_USER='elabtrack_runtime', DB_SSLMODE='disable',
           FRONTEND_URL='http://localhost:5173', ALLOWED_ORIGINS='http://localhost:5173,http://localhost:4173',
           LOG_FORMAT='json', LOG_LEVEL='info', RATE_LIMIT_WINDOW='5s',
           GOCACHE='/tmp/elabtrack-phase1b-go-cache', ELABTRACK_PHASE4A='1')
# Privileged settings are supplied only to explicit migration/test commands.
if '--migrate-up' in sys.argv or '--migrate-status' in sys.argv or 'test' in sys.argv:
    env['MIGRATION_DATABASE_URL'] = 'postgresql://elabtrack_migrator:' + quote(values['MIGRATION_DB_PASSWORD'], safe='') + '@127.0.0.1:25432/elabtrack_v2_phase4a_test?sslmode=disable'
else:
    env.pop('MIGRATION_DB_PASSWORD', None)
    env.pop('BOOTSTRAP_DB_PASSWORD', None)
if len(sys.argv) < 2:
    raise SystemExit('Usage: integration/phase4a-env-run.py COMMAND [ARGS]')
os.execvpe(sys.argv[1], sys.argv[1:], env)
