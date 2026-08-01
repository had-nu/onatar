#!/usr/bin/env bash
# Stops Onatar backend + frontend dev servers.
# Kills by listening port (ss) so `go run` orphaned binaries are caught too.
set -euo pipefail

stop_port() {
  local port="$1" name="$2"
  local pids
  pids="$(ss -tlnp 2>/dev/null | grep ":$port" | grep -oE 'pid=[0-9]+' | cut -d= -f2 | sort -u || true)"
  if [[ -n "$pids" ]]; then
    echo "stopping $name (pid(s): $(echo "$pids" | tr '\n' ' '))"
    kill $pids 2>/dev/null
    echo "$name stopped"
  else
    echo "$name not running"
  fi
}

pkill -f "go run ./cmd/server" 2>/dev/null || true

stop_port 8090 "backend"
stop_port 5173 "frontend"
exit 0
