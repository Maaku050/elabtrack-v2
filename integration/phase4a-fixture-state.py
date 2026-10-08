#!/usr/bin/env python3
"""Narrow, parameter-validated mutations of owned disposable browser fixtures."""
import json, os, subprocess, sys, uuid
from pathlib import Path

root = Path(__file__).resolve().parent.parent
values = dict(line.split('=',1) for line in (root/'backend/.env.phase4a').read_text().splitlines() if line and not line.startswith('#'))
if values.get('DB_NAME') != 'elabtrack_v2_phase4a_test' or values.get('DB_PORT') != '25432':
    raise SystemExit('Refusing non-isolated database')
path=Path('/tmp/elabtrack-phase4a-fixtures.json')
if path.stat().st_mode & 0o077: raise SystemExit('Refusing unsafe fixture-file permissions')
fixtures=json.loads(path.read_text())
action=sys.argv[1] if len(sys.argv)>1 else ''
kind=sys.argv[2] if len(sys.argv)>2 else ''
selected=list(fixtures.values()) if action=='cleanup' else [fixtures[kind]]
statements=[]
for fixture in selected:
    ident=str(uuid.UUID(fixture['id']))
    email=fixture['email']
    if not email.startswith('phase4a-browser-') or not email.endswith('@example.invalid') or any(c in email for c in "'\\\n\r"):
        raise SystemExit('Refusing unowned account')
    where=f"id='{ident}'::uuid AND email='{email}'::citext"
    if action=='cleanup': statement=f'DELETE FROM users WHERE {where}'
    elif action=='disable': statement=f'UPDATE users SET is_active=false WHERE {where}'
    elif action=='enable': statement=f'UPDATE users SET is_active=true WHERE {where}'
    elif action=='revoke': statement=f"UPDATE refresh_tokens SET revoked_at=now() WHERE user_id IN(SELECT id FROM users WHERE {where}) AND revoked_at IS NULL"
    elif action=='role' and len(sys.argv)==4 and sys.argv[3] in ('BORROWER','STAFF','ADMIN'): statement=f"UPDATE users SET role='{sys.argv[3]}' WHERE {where}"
    else: raise SystemExit('Unsupported fixture action')
    statements.append(statement+';')
sql="SELECT current_database()='elabtrack_v2_phase4a_test' AND current_user='elabtrack_runtime';\nBEGIN;\n"+'\n'.join(statements)+'\nCOMMIT;\n'
env=os.environ.copy();env['PGPASSWORD']=values['DB_PASSWORD']
result=subprocess.run(['docker','exec','-i','--env','PGPASSWORD','elabtrack_v2_phase4a_postgres','psql','-X','-qAt','-v','ON_ERROR_STOP=1','-U','elabtrack_runtime','-d','elabtrack_v2_phase4a_test'],input=sql,text=True,env=env,capture_output=True)
if result.returncode or result.stdout.strip()!='t': raise SystemExit('Isolated fixture action failed (details suppressed)')
if action=='cleanup': path.unlink()
print('Owned isolated fixture action completed.')
