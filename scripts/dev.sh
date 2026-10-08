#!/usr/bin/env bash
# Supervise the existing Make targets, including go run/npm descendants.
set -uo pipefail
set +m # Background children must not be group leaders before setsid runs.

repo_root="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
make_command="${1:-make}"

# A recipe referencing $(MAKE) also runs under make -n. Respect that mode.
make_flags="${MAKEFLAGS-}"
short_flags="${make_flags%% *}"
if [[ "$short_flags" != -* && "$short_flags" == *n* ]]; then
  printf 'Would supervise: make backend and make frontend (no servers started).\n'
  exit 0
fi

if (( BASH_VERSINFO[0] < 5 || (BASH_VERSINFO[0] == 5 && BASH_VERSINFO[1] < 1) )); then
  printf '[dev] Bash 5.1+ is required for wait -n -p.\n' >&2
  exit 1
fi
for command in "$make_command" setsid; do
  if ! command -v "$command" >/dev/null 2>&1; then
    printf '[dev] Required command unavailable: %s\n' "$command" >&2
    exit 1
  fi
done

pids=()
cleanup() {
  trap '' INT TERM HUP
  # jobs also covers an interrupt between launching a child and saving its PID.
  local pid attempt running
  local -a groups=("${pids[@]}")
  while read -r pid; do groups+=("$pid"); done < <(jobs -pr)
  if ((${#groups[@]} == 0)); then return; fi
  printf '\n[dev] Stopping both development targets…\n'
  for pid in "${groups[@]}"; do
    kill -TERM -- "-$pid" 2>/dev/null || kill -TERM "$pid" 2>/dev/null || true
  done
  # The API has a ten-second graceful shutdown budget. Bound cleanup too.
  for ((attempt = 0; attempt < 120; attempt++)); do
    running=false
    for pid in "${groups[@]}"; do
      if kill -0 -- "-$pid" 2>/dev/null || kill -0 "$pid" 2>/dev/null; then
        running=true
      fi
    done
    if [[ "$running" == false ]]; then break; fi
    sleep 0.1
  done
  for pid in "${groups[@]}"; do
    if kill -0 -- "-$pid" 2>/dev/null || kill -0 "$pid" 2>/dev/null; then
      printf '[dev] Target group %s exceeded shutdown grace; terminating it.\n' "$pid" >&2
      kill -KILL -- "-$pid" 2>/dev/null || kill -KILL "$pid" 2>/dev/null || true
    fi
    wait "$pid" 2>/dev/null || true
  done
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
trap 'exit 129' HUP

printf '[dev] Starting make backend and make frontend. Logs share this terminal.\n'
printf '[dev] PostgreSQL must already be configured/running; no migrations or seeds run.\n'
printf '[dev] Press Ctrl+C to stop both servers.\n'
setsid "$make_command" --no-print-directory -C "$repo_root" backend &
backend_pid=$!
pids+=("$backend_pid")
setsid "$make_command" --no-print-directory -C "$repo_root" frontend &
frontend_pid=$!
pids+=("$frontend_pid")

finished_pid=''
wait -n -p finished_pid "$backend_pid" "$frontend_pid"
status=$?
target=frontend
if [[ "$finished_pid" == "$backend_pid" ]]; then target=backend; fi
printf '[dev] %s exited with status %s; stopping its peer.\n' "$target" "$status" >&2
if [[ "$target" == backend && "$status" != 0 ]]; then
  printf '[dev] For database errors, check backend/.env and start your configured PostgreSQL.\n' >&2
  printf '[dev] Local Compose workflow: docker compose up -d postgres (see README.md).\n' >&2
fi
# A development server ending successfully is still an unexpected shutdown.
if ((status == 0)); then status=1; fi
exit "$status"
