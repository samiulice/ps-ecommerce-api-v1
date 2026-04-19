#!/bin/bash

# -----------------------------
# Configuration
# -----------------------------
VPS_HOST="203.161.48.179"
REMOTE_PATH="/home/samiul/apps/bin/urban-elegance-pse-api"
SERVICE_NAME="urbaneleganceapi.service"
PING_URL="https://urbanelegance-api.pssoft.xyz/api/v1/ping"
REMOTE_SCRIPT_PATH="/home/samiul/apps/bin//update-and-reload.sh"

# -----------------------------
# Step 1: Clean and Compile
# -----------------------------

echo "Removing old local binary..."
rm -f ./bin/app

echo "Cross-compiling Go binary for Linux..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-s -w" -o ./bin/app ./cmd/api || { echo "Build failed"; exit 1; }

# -----------------------------
# Step 2: Upload to /tmp on VPS
# -----------------------------
echo "Uploading binary to /tmp on VPS..."
scp -C -o ServerAliveInterval=60 ./bin/app samiul@"$VPS_HOST":/tmp/app || { echo "SCP failed"; exit 1; }

# -----------------------------
# Step 3: Execute Remote Script (WILL PROMPT FOR PASSWORD)
# -----------------------------
echo "Triggering remote deployment script..."
ssh -t -o ServerAliveInterval=60 samiul@"$VPS_HOST" "sudo $REMOTE_SCRIPT_PATH $REMOTE_PATH $SERVICE_NAME" || { echo "Remote script failed"; exit 1; }

# -----------------------------
# Step 4: Ping the API
# -----------------------------
echo "Pinging API..."
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$PING_URL")
echo "API response code: $HTTP_CODE"

if [[ "$HTTP_CODE" -ge 200 ]]; then
    echo "Deployment successful!"
else
    echo "Warning: API returned non-2xx response."
fi