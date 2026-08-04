#!/bin/bash
# Onatar development server runner
# Uses tmux to keep the server alive in the background

set -e

SESSION="onatar-api"
DIR="/home/hadnu/workspace/homelab/onatar"

# Kill existing session if exists
tmux kill-session -t "$SESSION" 2>/dev/null || true

# Create new tmux session and run server with signals disabled
tmux new-session -d -s "$SESSION" -c "$DIR" "DEV_NO_SIGNALS=1 go run ./cmd/server"

echo "Server started in tmux session: $SESSION"
echo "Attach with: tmux attach -t $SESSION"
echo "Stop with: tmux kill-session -t $SESSION"
echo ""
echo "Logs available in tmux buffer. Server should stay alive."