#!/usr/bin/env python3
"""Explicit disposable-role psql helper; passwords use environment, never argv."""
import os
from pathlib import Path
import subprocess
import sys
if os.environ.get('DB_NAME') != 'elabtrack_v2_integration' or os.environ.get('DB_PORT') != '15432':
    raise SystemExit('Refusing non-integration database')
role = sys.argv[1] if len(sys.argv)>1 else 'runtime'
settings = {'runtime':('elabtrack_runtime','DB_PASSWORD'),
            'migration':('elabtrack_migrator','MIGRATION_DB_PASSWORD'),
            'bootstrap':('postgres','BOOTSTRAP_DB_PASSWORD')}
if role not in settings:
    raise SystemExit('Choose runtime, migration or bootstrap')
user, key = settings[role]
query = Path(sys.argv[2]).read_text() if len(sys.argv)>2 else sys.stdin.read()
result = subprocess.run(['psql','-X','-h','127.0.0.1','-p','15432','-U',user,'-d','elabtrack_v2_integration','-At','-v','ON_ERROR_STOP=1'],
                        input=query,env={**os.environ,'PGPASSWORD':os.environ[key]},capture_output=True,text=True)
if result.returncode:
    raise SystemExit('Disposable SQL command failed (driver detail withheld)')
print(result.stdout.strip())
