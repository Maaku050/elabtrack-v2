#!/usr/bin/env python3
"""Back up the stopped, marked presentation and verify a full isolated restore.

Never restores over an existing database. Retains backup and new restore volume.
"""
import hashlib
import importlib.util
import json
import os
from pathlib import Path
import re
import secrets
import time

spec = importlib.util.spec_from_file_location('demo', Path(__file__).resolve().parents[1] / 'presentation-demo.py')
m = importlib.util.module_from_spec(spec)
spec.loader.exec_module(m)


def query(container, database, statement):
    if not re.fullmatch(r'elabtrack_v2_(presentation_demo|presentation_restore_[a-z0-9]+)', database):
        raise RuntimeError('Refusing foreign recovery database')
    command = 'PGPASSWORD="$MIGRATION_DB_PASSWORD" exec psql -X -h127.0.0.1 -Uelabtrack_migrator -d' + database + ' -vON_ERROR_STOP=1 -t -A'
    out = m.run(['docker', 'exec', '-i', container, 'sh', '-c', command], input=statement.encode()).decode()
    return [line for line in out.splitlines() if line and line not in ('BEGIN', 'COMMIT')]


def snapshot(container, database):
    tables = query(container, database, "BEGIN READ ONLY;SELECT tablename FROM pg_tables WHERE schemaname='public' ORDER BY tablename;COMMIT;")
    if not tables or any(not re.fullmatch('[a-z_][a-z_0-9]*', t) for t in tables):
        raise RuntimeError('Unexpected recovery schema')
    expressions = []
    for table in tables:
        expressions.append("'" + table + "',(SELECT json_build_array(count(*),md5(COALESCE(string_agg(to_jsonb(t)::text,E'\\n' ORDER BY to_jsonb(t)::text),''))) FROM public.\"" + table + '\" t)')
    result = query(container, database, 'BEGIN ISOLATION LEVEL REPEATABLE READ READ ONLY;SELECT json_build_object(' + ','.join(expressions) + ');COMMIT;')
    return json.loads(result[0])


def main():
    values, identity = m.target()
    if m.sql("SELECT identity||':'||seeded FROM presentation_demo_identity WHERE id=1;").strip() != 'elabtrack-presentation-demo:true':
        raise RuntimeError('Refusing unmarked/unseeded source')
    process_file = m.PRIVATE / 'processes.json'
    if process_file.exists():
        for record in m.read_private(process_file).values():
            try:
                os.kill(int(record["pid"]), 0)
            except (ProcessLookupError, TypeError, ValueError):
                continue
            raise RuntimeError('Stop the presentation before a consistent backup/restore rehearsal')
    stamp = time.strftime('%Y%m%dT%H%M%SZ', time.gmtime()).lower() + secrets.token_hex(2)
    private = m.PRIVATE / ('recovery-' + stamp)
    private.mkdir(mode=0o700)
    before = snapshot(m.CONTAINER, m.DB)
    command = 'PGPASSWORD="$MIGRATION_DB_PASSWORD" exec pg_dump --format=custom -h127.0.0.1 -Uelabtrack_migrator -d' + m.DB
    data = m.run(['docker', 'exec', m.CONTAINER, 'sh', '-c', command])
    if not data.startswith(b'PGDMP'):
        raise RuntimeError('Invalid PostgreSQL archive')
    backup = private / 'presentation.dump'
    fd = os.open(backup, os.O_WRONLY | os.O_CREAT | os.O_EXCL, 0o600)
    with os.fdopen(fd, 'wb') as f:
        f.write(data)
    database = 'elabtrack_v2_presentation_restore_' + stamp
    container = database + '_postgres'
    volume = database + '_pgdata'
    config = dict(values, DB_NAME=database, DB_PORT='54837', DB_PASSWORD=secrets.token_urlsafe(32), MIGRATION_DB_PASSWORD=secrets.token_urlsafe(32), BOOTSTRAP_DB_PASSWORD=secrets.token_urlsafe(32))
    m.protected(private / 'compose.env', '\n'.join(k + '=' + v for k, v in config.items()) + '\n')
    m.protected(private / 'compose.yml', 'services:\n  postgres:\n    container_name: ' + container + '\n    ports: !override\n      - "127.0.0.1:54837:5432"\nvolumes:\n  elabtrack_v2_pgdata:\n    name: ' + volume + '\n')
    args = ['docker', 'compose', '--env-file', private / 'compose.env', '-p', 'elabtrack-recovery-' + stamp, '-f', m.ROOT / 'docker-compose.yml', '-f', private / 'compose.yml']
    m.run(args + ['up', '-d', '--wait', '--pull', 'never', 'postgres'])
    info = json.loads(m.run(['docker', 'inspect', container]))[0]
    if [x['Name'] for x in info['Mounts'] if x['Type'] == 'volume'] != [volume] or info['HostConfig']['PortBindings']['5432/tcp'] != [{'HostIp': '127.0.0.1', 'HostPort': '54837'}]:
        raise RuntimeError('Restore target identity mismatch')
    command = 'PGPASSWORD="$POSTGRES_PASSWORD" exec pg_restore --exit-on-error --single-transaction -h127.0.0.1 -Upostgres -d "$POSTGRES_DB"'
    m.run(['docker', 'exec', '-i', container, 'sh', '-c', command], input=data)
    after = snapshot(container, database)
    if before != after:
        raise RuntimeError('Restored application row digests differ')
    sequence_query = "BEGIN READ ONLY;SELECT coalesce(jsonb_agg(to_jsonb(s) ORDER BY sequencename),'[]') FROM pg_sequences s WHERE schemaname='public';COMMIT;"
    if query(m.CONTAINER, m.DB, sequence_query) != query(container, database, sequence_query):
        raise RuntimeError('Restored sequence state differs')
    applied = query(container, database, 'BEGIN READ ONLY;SELECT count(*) FROM schema_migrations;COMMIT;')[0]
    if applied != '11':
        raise RuntimeError('Restore migration tracking differs')
    result = {'result': 'PASS', 'full_restore': True, 'tables_verified': len(before), 'migration_rows': 11, 'sequence_state_equal': True, 'backup_sha256': hashlib.sha256(data).hexdigest(), 'container': container, 'database': database, 'port': 54837, 'volume_retained': volume, 'normal_data_untouched': True}
    m.protected(private / 'result.json', json.dumps(result, indent=2) + '\n')
    m.protected(m.PRIVATE / 'latest-recovery.json', json.dumps(result, indent=2) + '\n')
    # Release the rehearsal port without deleting recovered data or its volume.
    m.run(['docker', 'stop', container])
    print(json.dumps(result))


if __name__ == '__main__':
    main()
