#!/usr/bin/env bash
# Starts Onatar backend (Go, :8090) + frontend (Vite, :5173) in background.
# Shell-agnostic: call with `bash scripts/start.sh` from any shell (bash/fish/zsh).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

if [[ ! -f .env ]]; then
  echo "missing .env — cp .env.example .env" >&2
  exit 1
fi

set -a
# shellcheck disable=SC1091
source .env
set +a

if ss -tln | grep -q ':8090 '; then
  echo "port 8090 already in use — backend not started" >&2
else
  nohup go run ./cmd/server > /tmp/onatar-backend.log 2>&1 &
  disown
  echo "backend  -> http://localhost:8090/health  (log: /tmp/onatar-backend.log)"
fi

if ss -tln | grep -q ':5173 '; then
  echo "port 5173 already in use — frontend not started" >&2
else
  nohup npm run dev --prefix frontend > /tmp/onatar-frontend.log 2>&1 &
  disown
  echo "frontend -> http://localhost:5173/        (log: /tmp/onatar-frontend.log)"
fi

sleep 2
echo
echo "check: curl http://localhost:8090/health"
