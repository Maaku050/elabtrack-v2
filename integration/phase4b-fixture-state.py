#!/usr/bin/env python3
"""Mutate only owned Phase 4B fixtures; never delete immutable terms history."""
import json
import os
from pathlib import Path
import subprocess
import sys
import uuid

root = Path(__file__).resolve().parent.parent
env_path = root / 'backend/.env.phase4b'
if env_path.stat().st_mode & 0o077:
    raise SystemExit('Unsafe isolated environment permissions')
values = dict(line.split('=', 1) for line in env_path.read_text().splitlines() if line and not line.startswith('#'))
if values.get('DB_NAME') != 'elabtrack_v2_phase4b_test' or values.get('DB_PORT') != '35432':
    raise SystemExit('Refusing a non-isolated target')
path = Path('/tmp/elabtrack-phase4b-fixtures.json')
if path.stat().st_mode & 0o077:
    raise SystemExit('Unsafe fixture permissions')
fixtures = json.loads(path.read_text())
action = sys.argv[1] if len(sys.argv) > 1 else ''
kind = sys.argv[2] if len(sys.argv) > 2 else ''
guard = "DO $$ BEGIN IF current_database()<>'elabtrack_v2_phase4b_test' OR current_user<>'elabtrack_runtime' THEN RAISE EXCEPTION 'Unsafe test identity'; END IF; END $$;"
if action == 'reset-current' and len(sys.argv) == 2:
    # Test-only reset simulates an empty publication pointer while preserving
    # all existing synthetic versions and receipts. No product withdrawal API.
    statement = """DO $$ BEGIN
    IF EXISTS(SELECT 1 FROM terms_versions v JOIN users u ON u.id=v.published_by
      WHERE v.version NOT LIKE 'TEST-%' OR v.title NOT LIKE 'SYNTHETIC%'
      OR u.email::text NOT LIKE 'phase4b-%@example.invalid') THEN
      RAISE EXCEPTION 'Non-synthetic history'; END IF;
    END $$;
    UPDATE terms_publication SET current_version_id=NULL WHERE id=1;"""
elif action in ('disable', 'enable', 'revoke', 'role') and kind in fixtures:
    fixture = fixtures[kind]
    ident = str(uuid.UUID(fixture['id']))
    email = fixture['email']
    if not email.startswith('phase4b-browser-') or not email.endswith('@example.invalid') or any(c in email for c in "'\\\n\r"):
        raise SystemExit('Refusing an unowned fixture')
    where = f"id='{ident}'::uuid AND email='{email}'::citext"
    if action == 'disable':
        statement = f'UPDATE users SET is_active=false WHERE {where};'
    elif action == 'enable':
        statement = f'UPDATE users SET is_active=true WHERE {where};'
    elif action == 'revoke':
        statement = f'UPDATE refresh_tokens SET revoked_at=now() WHERE user_id IN(SELECT id FROM users WHERE {where}) AND revoked_at IS NULL;'
    elif len(sys.argv) == 4 and sys.argv[3] in ('BORROWER', 'STAFF', 'ADMIN'):
        statement = f"UPDATE users SET role='{sys.argv[3]}' WHERE {where};"
    else:
        raise SystemExit('Invalid fixture role')
else:
    raise SystemExit('Unsupported isolated fixture action')
env = os.environ.copy()
env['PGPASSWORD'] = values['DB_PASSWORD']
result = subprocess.run(['docker', 'exec', '-i', '--env', 'PGPASSWORD', 'elabtrack_v2_phase4b_postgres', 'psql', '-X', '-qAt', '-v', 'ON_ERROR_STOP=1', '-U', 'elabtrack_runtime', '-d', 'elabtrack_v2_phase4b_test'], input=guard+'\nBEGIN;\n'+statement+'\nCOMMIT;\n', text=True, env=env, capture_output=True)
if result.returncode:
    raise SystemExit('Isolated fixture action failed (details suppressed)')
print('Owned isolated fixture action completed; history preserved.')
