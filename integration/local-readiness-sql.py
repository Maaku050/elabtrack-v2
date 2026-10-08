#!/usr/bin/env python3
"""Expire/clean only a synthetic user created by the local readiness smoke test."""
import json
import os
from pathlib import Path
import subprocess
import sys
import uuid

root = Path(__file__).resolve().parent.parent
values = dict(line.split("=", 1) for line in (root / "backend/.env").read_text().splitlines()
              if line and not line.startswith("#"))
if not (values.get("APP_ENV") == "development" and values.get("DB_NAME") == "elabtrack_v2"
        and values.get("DB_USER") == "elabtrack_runtime" and values.get("DB_HOST") == "127.0.0.1"):
    raise SystemExit("Refusing a non-development runtime target")
request = json.load(sys.stdin)
fixture_id = str(uuid.UUID(request["id"]))
guard = "id=:'fixture_id' AND role='user' AND email LIKE 'local-readiness-%@example.invalid'"
if request["action"] == "expire":
    sql = """UPDATE refresh_tokens SET created_at=NOW()-INTERVAL '2 days',
expires_at=NOW()-INTERVAL '1 day' WHERE revoked_at IS NULL
AND user_id IN (SELECT id FROM users WHERE """ + guard + ") RETURNING id;"
elif request["action"] == "cleanup":
    sql = "DELETE FROM users WHERE " + guard + " RETURNING id;"
else:
    raise SystemExit("Unsupported smoke action")
result = subprocess.run([
    "psql", "-X", "-h", values["DB_HOST"], "-p", values["DB_PORT"],
    "-U", values["DB_USER"], "-d", values["DB_NAME"], "-At",
    "-v", "ON_ERROR_STOP=1", "-v", "fixture_id=" + fixture_id,
], input=sql, env={**os.environ, "PGPASSWORD": values["DB_PASSWORD"], "PGCONNECT_TIMEOUT": "5"},
    capture_output=True, text=True)
if result.returncode:
    raise SystemExit("Restricted smoke SQL failed; credentials not displayed")
matched = fixture_id in result.stdout if request["action"] == "cleanup" else "UPDATE 1" in result.stdout
print(json.dumps({"action": request["action"], "matched": matched, "ok": matched}))
if not matched:
    raise SystemExit(1)
