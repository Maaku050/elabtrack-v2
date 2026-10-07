#!/usr/bin/env python3
"""Inspect only safe role/config metadata; never print Docker environment."""
import json
import os
from pathlib import Path
import subprocess
import urllib.request

docker=os.environ.get('DOCKER_BIN','docker')
result=subprocess.run([docker,'inspect','elabtrack_v2_phase1h_backend'],capture_output=True,text=True,check=True)
container=json.loads(result.stdout)[0]
values=dict(value.split('=',1) for value in container['Config']['Env'])
assert values['DB_USER']=='elabtrack_runtime'
assert not any(key in values for key in ('MIGRATION_DATABASE_URL','MIGRATION_DB_PASSWORD','BOOTSTRAP_DB_PASSWORD'))
assert values['DB_PASSWORD']==os.environ['DB_PASSWORD']
assert values['DB_PASSWORD']!=os.environ['MIGRATION_DB_PASSWORD']
query="""SELECT rolname,rolsuper,rolcreatedb,rolcreaterole,rolreplication,rolbypassrls FROM pg_roles WHERE rolname IN ('elabtrack_runtime','elabtrack_migrator') ORDER BY rolname;
SELECT c.relname,r.rolname FROM pg_class c JOIN pg_roles r ON r.oid=c.relowner WHERE c.relnamespace='public'::regnamespace AND c.relname IN ('users','refresh_tokens','schema_migrations') ORDER BY c.relname;
SELECT has_schema_privilege('elabtrack_runtime','public','CREATE'),has_table_privilege('elabtrack_runtime','schema_migrations','SELECT'),has_table_privilege('elabtrack_runtime','schema_migrations','INSERT'),pg_has_role('elabtrack_runtime','elabtrack_migrator','MEMBER');
SELECT version,checksum FROM schema_migrations ORDER BY version;"""
p=subprocess.run(['psql','-X','-h','127.0.0.1','-p','15432','-U','elabtrack_migrator','-d','elabtrack_v2_integration','-At','-v','ON_ERROR_STOP=1'],input=query,env={**os.environ,'PGPASSWORD':os.environ['MIGRATION_DB_PASSWORD']},capture_output=True,text=True)
assert p.returncode==0
rows=p.stdout.strip().splitlines();assert rows[0]=='elabtrack_migrator|f|f|f|f|f' and rows[1]=='elabtrack_runtime|f|f|f|f|f'
assert all(row.endswith('|elabtrack_migrator') for row in rows[2:5])
assert rows[5]=='f|f|f|f'
assert len(rows[6:])==3
with urllib.request.urlopen('http://localhost:18080/api/v1/ready') as response:
    assert response.status==200
safe={'runtimeAPIOnlyCredentials':True,'migrationSecretAbsentFromAPI':True,'ready':200,'roles':[row.split('|')[0] for row in rows[:2]],'allRolesNonSuperuser':True,'allRolesNoClusterAdministration':True,'applicationAndTrackingOwner':'elabtrack_migrator','runtimeSchemaCreate':False,'runtimeTrackingAccess':False,'runtimeMigrationRoleMembership':False,'versions':[row.split('|')[0] for row in rows[6:]]}
Path('/tmp/elabtrack-phase1h-roles.json').write_text(json.dumps(safe,indent=2))
print('PASS API/runtime identity, no migration secret, readiness, non-superuser owner, private tracking and stable object ownership')
