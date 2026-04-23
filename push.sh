#!/bin/bash

# -----------------------------
# Configuration
# -----------------------------
VPS_HOST="203.161.48.179"
REMOTE_PATH="/home/samiul/apps/bin/noor-api"
SERVICE_NAME="noorapi.service"
PING_URL="https://noor-api.pssoft.xyz/api/v1/ping"
REMOTE_SCRIPT_PATH="/home/samiul/apps/bin/update-and-reload.sh"
BINARY_PATH="./bin/app"

# -----------------------------
# Step 1: Proficient Change Detection
# -----------------------------

SHOULD_BUILD=false

if [[ ! -f "$BINARY_PATH" ]]; then
    echo "Binary missing. Forcing build..."
    SHOULD_BUILD=true
else
    # 1. Check for 'Dirty' files (Modified, Staged, or Untracked)
    # We ignore the 'bin/' directory and documentation files using grep -v
    # The --porcelain flag ensures the output is stable for scripts
    CHANGES=$(git status --porcelain | grep -vE ' bin/|\.md$|docs/|^\?\? \.env' | wc -l)

    if [[ "$CHANGES" -gt 0 ]]; then
        echo "Detected $CHANGES modified/new files. Rebuilding..."
        SHOULD_BUILD=true
    else
        # 2. Workspace is clean, but is the binary older than the last commit?
        # This handles cases where you just pulled changes but haven't built them yet.
        LAST_COMMIT_DATE=$(git log -1 --format=%ct)
        BINARY_MOD_DATE=$(stat -c %Y "$BINARY_PATH" 2>/dev/null || stat -f %m "$BINARY_PATH")

        if [[ "$LAST_COMMIT_DATE" -gt "$BINARY_MOD_DATE" ]]; then
            echo "Binary is out of date compared to the latest commit. Rebuilding..."
            SHOULD_BUILD=true
        else
            echo "Environment is clean and binary is up to date. Skipping compilation."
        fi
    fi
fi

if [[ "$SHOULD_BUILD" == true ]]; then
    echo "Removing old local binary..."
    rm -f "$BINARY_PATH"

    echo "Cross-compiling Go binary for Linux..."
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-s -w" -o "$BINARY_PATH" ./cmd/api || { echo "Build failed"; exit 1; }
fi

# -----------------------------
# Step 2: Upload to /tmp on VPS
# -----------------------------
echo "Uploading binary to /tmp on VPS..."
scp -C -o ServerAliveInterval=60 "$BINARY_PATH" samiul@"$VPS_HOST":/tmp/app || { echo "SCP failed"; exit 1; }

# -----------------------------
# Step 3: Execute Remote Script
# -----------------------------
echo "Triggering remote deployment script..."
ssh -t -o ServerAliveInterval=60 samiul@"$VPS_HOST" "sudo $REMOTE_SCRIPT_PATH $REMOTE_PATH $SERVICE_NAME" || { echo "Remote script failed"; exit 1; }

# -----------------------------
# Step 4: Ping the API
# -----------------------------
echo "Pinging API..."
sleep 2 
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$PING_URL")
echo "API response code: $HTTP_CODE"

if [[ "$HTTP_CODE" -ge 200 && "$HTTP_CODE" -lt 300 ]]; then
    echo "Deployment successful!"
else
    echo "Warning: API returned code $HTTP_CODE"
fi