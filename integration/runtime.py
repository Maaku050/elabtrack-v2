"""Local runtime evidence. Invoke through env-run.py with the isolated stack up."""
import concurrent.futures
import http.client
import json
import os
from pathlib import Path
import subprocess
import tempfile
import time
import uuid

root = Path(__file__).resolve().parent.parent
docker = os.environ.get("DOCKER_BIN", "docker")
api = "/tmp/elabtrack-phase1g-api"
evidence = {"requests": [], "checks": {}}

def sql(query):
    result = subprocess.run(["psql", "-h", "127.0.0.1", "-p", "15432", "-U", "postgres", "-d", "elabtrack_v2_integration", "-At", "-v", "ON_ERROR_STOP=1"],
                            input=query, env={**os.environ, "PGPASSWORD": os.environ["DB_PASSWORD"]}, capture_output=True, text=True)
    if result.returncode:
        raise RuntimeError("isolated fixture SQL failed (detail intentionally withheld)")
    return result.stdout.strip()

def request(path, method="GET", headers=None, body=None, port=18080, expected=None):
    conn = http.client.HTTPConnection("127.0.0.1", port, timeout=10)
    conn.request(method, "/api/v1"+path, body=body, headers=headers or {})
    r = conn.getresponse()
    data = r.read()
    hdr = dict((k.lower(),v) for k,v in r.getheaders())
    conn.close()
    if expected is not None:
        assert r.status == expected, (path,r.status,expected)
    ident = hdr.get("x-request-id")
    if ident:
        uuid.UUID(ident)
    if r.status >= 400:
        parsed = json.loads(data)
        assert parsed["error"]["requestId"] == ident
        assert parsed["data"] is None and parsed["meta"] is None and not parsed["success"]
    evidence["requests"].append({"path":path,"status":r.status,"id":ident,"port":port})
    return r.status, hdr, data

def check(name, result):
    evidence["checks"][name] = result
    print("PASS",name,flush=True)

schema = sql("SELECT column_name FROM information_schema.columns WHERE table_name='refresh_tokens' ORDER BY ordinal_position")
assert "token_hash" in schema.splitlines() and "token" not in schema.splitlines()
check("final_schema", {"columns":schema.splitlines(),
    "constraints":sql("SELECT conname FROM pg_constraint WHERE conrelid='refresh_tokens'::regclass ORDER BY conname").splitlines(),
    "indexes":sql("SELECT indexname FROM pg_indexes WHERE tablename='refresh_tokens' ORDER BY indexname").splitlines(),
    "versions":sql("SELECT version FROM schema_migrations ORDER BY version").splitlines(),
    "roles":sql("SELECT current_user, rolsuper FROM pg_roles WHERE rolname=current_user")})

# A synthetic migration probe is confined to this disposable database and an
# external temporary directory. Existing migration files remain untouched.
with tempfile.TemporaryDirectory(prefix="elabtrack-phase1g-migrator-") as directory:
    migrations = Path(directory)/"migrations"
    migrations.mkdir()
    version="000099_phase1g_probe"
    (migrations/(version+".up.sql")).write_text("SELECT pg_sleep(1); CREATE TABLE phase1g_runner_probe(id integer);")
    (migrations/(version+".down.sql")).write_text("DROP TABLE phase1g_runner_probe;")
    def run_up():
        return subprocess.run([api,"--migrate-up"],cwd=directory,capture_output=True,text=True)
    try:
        sql("ALTER TABLE schema_migrations ADD CONSTRAINT phase1g_reject_probe CHECK(version <> '000099_phase1g_probe')")
        failed=run_up()
        assert failed.returncode==1
        assert sql("SELECT to_regclass('phase1g_runner_probe') IS NOT NULL") == "t"
        assert sql("SELECT count(*) FROM schema_migrations WHERE version='000099_phase1g_probe'") == "0"
        check("DDL_committed_before_failed_bookkeeping",True)
    finally:
        sql("DROP TABLE IF EXISTS phase1g_runner_probe; ALTER TABLE schema_migrations DROP CONSTRAINT IF EXISTS phase1g_reject_probe")
    try:
        with concurrent.futures.ThreadPoolExecutor(max_workers=2) as executor:
            runs=list(executor.map(lambda _:run_up(),range(2)))
        assert sorted(r.returncode for r in runs)==[0,1]
        assert sql("SELECT count(*) FROM schema_migrations WHERE version='000099_phase1g_probe'")=="1"
        check("concurrent_runner", {"exitCodes":[r.returncode for r in runs],"oneDDLFailure":True,"advisoryLock":False})
        (migrations/(version+".up.sql")).write_text("SELECT 'modified migration content';")
        assert run_up().returncode==0
        check("modified_applied_migration_not_detected",True)
    finally:
        sql("DROP TABLE IF EXISTS phase1g_runner_probe; DELETE FROM schema_migrations WHERE version='000099_phase1g_probe'")

# A temporary restricted login is solely a fault fixture, not a new credential
# architecture. Denied tracking reads must not masquerade as an empty database.
lookup_password=uuid.uuid4().hex+uuid.uuid4().hex
sentinels=Path('/tmp/elabtrack-phase1g-secrets.json')
sentinels.write_text(json.dumps(json.loads(sentinels.read_text())+[lookup_password]))
assert sql("SELECT count(*) FROM pg_roles WHERE rolname='phase1g_lookup_probe'")=="0"
try:
    sql("CREATE ROLE phase1g_lookup_probe LOGIN PASSWORD '"+lookup_password+"'; GRANT USAGE, CREATE ON SCHEMA public TO phase1g_lookup_probe")
    p=subprocess.run([api,"--migrate-down"],env={**os.environ,"DB_USER":"phase1g_lookup_probe","DB_PASSWORD":lookup_password},cwd=root/"backend",capture_output=True,text=True)
    assert p.returncode==0 and "no migrations to roll back" in p.stdout
    check("denied_last_applied_lookup_swallowed",{"exitCode":0,"realState":"three versions applied","reported":"no migrations to roll back"})
finally:
    sql("REVOKE USAGE, CREATE ON SCHEMA public FROM phase1g_lookup_probe; DROP ROLE IF EXISTS phase1g_lookup_probe")

# Config validation must fail before the unreachable production endpoint is used.
production={**os.environ,"APP_ENV":"production","DB_HOST":"127.0.0.1","DB_SSLMODE":"disable","FRONTEND_URL":"https://app.example.invalid","ALLOWED_ORIGINS":"https://app.example.invalid"}
for args,label in [([],"production_unsafe_TLS_rejected"),(["--seed"],"production_seed_rejected")]:
    p=subprocess.run([api,*args],env=production,cwd=root/"backend",capture_output=True,text=True)
    assert p.returncode==1 and "DB_SSLMODE" in p.stderr
    check(label,"startup validation rejects disabled TLS before connection/seed IO")
# Use valid verified TLS configuration to independently reach the seed guard.
production["DB_SSLMODE"]="verify-full"
p=subprocess.run([api,"--seed"],env=production,cwd=root/"backend",capture_output=True,text=True)
assert p.returncode==1 and "development" in p.stderr
check("production_seed_guard", "valid production config refuses seed before database IO")

for incoming in ["8d2fd646-00d0-4c88-943f-13ba0c65a68a", "malformed-id", "x"*2048]:
    status,h,_=request("/health",headers={"X-Request-ID":incoming},expected=200)
    assert h["x-request-id"] != incoming
check("incoming_IDs_ignored",True)
for label,port,hdr in [("direct",18080,{}),("direct_spoof",18080,{"X-Forwarded-For":"198.51.100.42","X-Forwarded-Proto":"https"}),("proxied_spoof",15173,{"X-Forwarded-For":"198.51.100.42","X-Forwarded-Proto":"https"})]:
    _,h,_=request("/health",port=port,headers=hdr,expected=200)
    assert "strict-transport-security" not in h
    evidence["checks"][label]={"requestId":h["x-request-id"]}
check("local_HSTS_absent",True)

# Real socket parser boundary, deliberately beyond the global body/header limits.
_,h,_=request("/not-implemented","POST",{"Content-Type":"application/json"},"x"*((1<<20)+1),expected=413)
check("socket_global_body_limit",{"requestId":h.get("x-request-id")})
_,h,_=request("/health",headers={"X-Large-Header":"x"*10000},expected=431)
check("socket_header_limit",{"requestId":h.get("x-request-id")})
_,_,_=request("/auth/login","POST",{"Origin":"http://localhost:15173","Content-Type":"application/json"}," "*16384+"{}",expected=413)
check("socket_auth_body_limit",True)

before=sql("SELECT string_agg(version||':'||applied_at::text,',' ORDER BY version) FROM schema_migrations")
subprocess.run([docker,"stop","elabtrack_v2_phase1g_postgres"],check=True,capture_output=True)
try:
    request("/health",expected=200)
    request("/ready",expected=503)
    check("DB_unavailable_liveness_and_readiness",{"liveness":200,"readiness":503})
finally:
    subprocess.run([docker,"start","elabtrack_v2_phase1g_postgres"],check=True,capture_output=True)
for _ in range(20):
    status,_,_=request("/ready")
    if status==200:break
    time.sleep(.25)
assert status==200
check("API_pool_reconnect",True)

for path,limit,method in [("/auth/login",10,"POST"),("/auth/refresh",60,"POST"),("/health",120,"GET")]:
    time.sleep(5.2)
    for i in range(limit):
        status,h,_=request(path,method,{"Origin":"http://localhost:15173","Content-Type":"application/json"},"{}" if method=="POST" else None)
        assert status!=429,(path,i)
    status,h,_=request(path,method,{"Origin":"http://localhost:15173","Content-Type":"application/json"},"{}" if method=="POST" else None,expected=429)
    assert int(h["retry-after"])>=1
    check("rate_limit_"+path,{"allowed":limit,"next":429,"retryAfter":h["retry-after"],"requestId":h["x-request-id"],"window":"5s"})

assert sql("SELECT string_agg(version||':'||applied_at::text,',' ORDER BY version) FROM schema_migrations")==before
check("explicit_migrations_unchanged_by_reconnect",True)
Path('/tmp/elabtrack-phase1g-runtime.json').write_text(json.dumps(evidence,indent=2))
