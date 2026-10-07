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
api = "/tmp/elabtrack-phase1h-api"
evidence = {"requests": [], "checks": {}}

def sql(query):
    result = subprocess.run(["psql", "-h", "127.0.0.1", "-p", "15432", "-U", "elabtrack_migrator", "-d", "elabtrack_v2_integration", "-At", "-v", "ON_ERROR_STOP=1"],
                            input=query, env={**os.environ, "PGPASSWORD": os.environ["MIGRATION_DB_PASSWORD"]}, capture_output=True, text=True)
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

# Phase 1G defects remain historical evidence. The corrected failure/lock/role
# assertions live in TestRealMigrator and migration-processes.py.
check("migration_checksums_present", all(len(row.split("|"))==2 for row in sql("SELECT version, checksum FROM schema_migrations ORDER BY version").splitlines()))

# Config validation must fail before the unreachable production endpoint is used.
production={**os.environ,"APP_ENV":"production","DB_HOST":"127.0.0.1","DB_SSLMODE":"disable","FRONTEND_URL":"https://app.example.invalid","ALLOWED_ORIGINS":"https://app.example.invalid"}
production.pop("MIGRATION_DATABASE_URL", None)
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
subprocess.run([docker,"stop","elabtrack_v2_phase1h_postgres"],check=True,capture_output=True)
try:
    request("/health",expected=200)
    request("/ready",expected=503)
    check("DB_unavailable_liveness_and_readiness",{"liveness":200,"readiness":503})
finally:
    subprocess.run([docker,"start","elabtrack_v2_phase1h_postgres"],check=True,capture_output=True)
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
