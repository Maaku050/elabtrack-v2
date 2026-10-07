#!/usr/bin/env python3
"""Run a local integration command without putting generated secrets in argv."""
import os
from pathlib import Path
import sys

root = Path(__file__).resolve().parent.parent
values = dict(line.split("=", 1) for line in (root / "backend/.env.phase1g").read_text().splitlines() if line and not line.startswith("#"))
if values.get("DB_NAME") != "elabtrack_v2_integration" or values.get("DB_PORT") != "15432":
    raise SystemExit("Refusing a non-integration target")
env = os.environ.copy()
env.update(values)
env.update(APP_ENV="development", DB_HOST="127.0.0.1", DB_SSLMODE="disable", APP_PORT="18081", LOG_FORMAT="json", LOG_LEVEL="debug",
           FRONTEND_URL="http://localhost:15173", ALLOWED_ORIGINS="http://localhost:15173,http://localhost:14173",
           GOCACHE="/tmp/elabtrack-phase1b-go-cache", ELABTRACK_INTEGRATION="1",
           INTEGRATION_SECRET_FILE="/tmp/elabtrack-phase1g-secrets.json")
if len(sys.argv) < 2:
    raise SystemExit("Usage: integration/env-run.py COMMAND [ARGS]")
os.execvpe(sys.argv[1], sys.argv[1:], env)
