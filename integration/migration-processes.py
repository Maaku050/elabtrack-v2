#!/usr/bin/env python3
"""Real CLI process exclusion and secret-safe failure logs on disposable PG."""
import json
import os
from pathlib import Path
import shutil
import subprocess
import tempfile
import time

root = Path(__file__).resolve().parent.parent
api = '/tmp/elabtrack-phase1h-api'
if os.environ.get('DB_NAME')!='elabtrack_v2_integration' or os.environ.get('DB_PORT')!='15432':
    raise SystemExit('Refusing non-integration target')

def sql(query):
    run = subprocess.run(['psql','-X','-h','127.0.0.1','-p','15432','-U','elabtrack_migrator','-d','elabtrack_v2_integration','-At','-v','ON_ERROR_STOP=1'],input=query,
        env={**os.environ,'PGPASSWORD':os.environ['MIGRATION_DB_PASSWORD']},capture_output=True,text=True)
    if run.returncode:
        raise RuntimeError('Disposable process fixture query failed; detail withheld')
    return run.stdout.strip()

def command(directory, flag):
    return subprocess.run([api,flag],cwd=directory,capture_output=True,text=True,timeout=15)

logs=[]
with tempfile.TemporaryDirectory(prefix='elabtrack-phase1h-process-') as directory:
    migrations=Path(directory)/'migrations';shutil.copytree(root/'backend/migrations',migrations)
    up=migrations/'000099_process_probe.up.sql';down=migrations/'000099_process_probe.down.sql'
    up.write_text('CREATE TABLE phase1h_process_probe(id int); SELECT pg_sleep(2);')
    down.write_text('DROP TABLE phase1h_process_probe;')
    first=subprocess.Popen([api,'--migrate-up'],cwd=directory,stdout=subprocess.PIPE,stderr=subprocess.PIPE,text=True)
    try:
        deadline=time.monotonic()+5
        while time.monotonic()<deadline:
            if sql("SELECT EXISTS(SELECT 1 FROM pg_locks WHERE locktype='advisory' AND classid=1162625346::oid AND objid=1296648018::oid AND granted)")=='t':break
            if first.poll() is not None:raise RuntimeError('First migrator finished before lock observation')
            time.sleep(.01)
        else:raise RuntimeError('Native migration lock not observed')
        second=command(directory,'--migrate-up');status=command(directory,'--migrate-status')
        assert second.returncode==1 and 'another migrator holds the lock' in second.stderr
        assert status.returncode==1 and 'another migrator holds the lock' in status.stderr
        out,err=first.communicate(timeout=10);assert first.returncode==0
        logs.extend([out,err,second.stdout,second.stderr,status.stdout,status.stderr])
        assert sql("SELECT count(*) FROM schema_migrations WHERE version='000099_process_probe'")=='1'
        assert sql("SELECT to_regclass('phase1h_process_probe') IS NOT NULL")=='t'
        verified=command(directory,'--migrate-status');assert verified.returncode==0 and 'checksum=verified' in verified.stdout
        reverted=command(directory,'--migrate-down');assert reverted.returncode==0
        up.write_text("CREATE TABLE phase1h_process_probe(id int); DO $$ BEGIN RAISE EXCEPTION 'sql-private-sentinel'; END $$;")
        failure=command(directory,'--migrate-up');assert failure.returncode==1 and '000099_process_probe' in failure.stderr
        logs.extend([failure.stdout,failure.stderr])
        assert sql("SELECT to_regclass('phase1h_process_probe') IS NULL AND NOT EXISTS(SELECT 1 FROM schema_migrations WHERE version='000099_process_probe')")=='t'
    finally:
        if first.poll() is None:first.kill();first.communicate()
        sql("DROP TABLE IF EXISTS phase1h_process_probe; DELETE FROM schema_migrations WHERE version='000099_process_probe'")
combined='\n'.join(logs)
sentinels=[os.environ[k] for k in ('DB_PASSWORD','MIGRATION_DB_PASSWORD','BOOTSTRAP_DB_PASSWORD','JWT_SECRET','MIGRATION_DATABASE_URL')]+['sql-private-sentinel']
assert all(value not in combined for value in sentinels)
structured=[json.loads(line) for line in combined.splitlines() if line.startswith('{')]
assert all(row.get('event') in ('migration.command','migration.apply') for row in structured)
result={'separateProcesses':True,'nativeAdvisoryLockObserved':True,'winnerExit':0,'loserExit':1,'statusDuringMigrationExit':1,'oneCoherentAppliedRecord':True,'lockReleasedAfterCommand':True,'SQLAndBookkeepingRollback':True,'structuredMigrationRecords':len(structured),'secretSentinels':len(sentinels),'secretMatches':0}
Path('/tmp/elabtrack-phase1h-processes.json').write_text(json.dumps(result,indent=2))
Path('/tmp/elabtrack-phase1h-migration.log').write_text(combined)
print('PASS real separate CLI processes: one lock owner, clean loser/status refusal, coherent history, rollback and redacted logs')
