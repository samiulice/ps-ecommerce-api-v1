#!/bin/bash

# -----------------------------
# Configuration
# -----------------------------
VPS_HOST="203.161.48.179"
REMOTE_PATH="/home/samiul/apps/bin/pse-api"
SERVICE_NAME="pseapi.service"
PING_URL="https://pse-api.pssoft.xyz/api/v1/ping"

# -----------------------------
# Step 1: Remove old local binary
# -----------------------------
echo "Removing old binary..."
rm -f ./bin/app

# -----------------------------
# Step 2: Cross-compile Go app for Linux
# -----------------------------
echo "Cross-compiling Go binary for Linux..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./bin/app ./cmd/api
if [[ $? -ne 0 ]]; then
    echo "Build failed. Exiting."
    exit 1
fi


# -----------------------------
# Step 3: Stop the service on VPS
# -----------------------------
echo "Stopping remote service..."
ssh samiul@"$VPS_HOST" "sudo systemctl stop $SERVICE_NAME"
if [[ $? -ne 0 ]]; then
    echo "Failed to stop service. Exiting."
    exit 1
fi

# -----------------------------
# Step 4: Upload new binary with fast cipher and compression
# -----------------------------
echo "Uploading new binary..."
scp -C -c chacha20-poly1305@openssh.com ./bin/app samiul@"$VPS_HOST":"$REMOTE_PATH"
if [[ $? -ne 0 ]]; then
    echo "SCP failed. Exiting."
    exit 1
fi

# -----------------------------
# Step 5: Restart the service
# -----------------------------
echo "Restarting remote service..."
ssh samiul@"$VPS_HOST" "sudo systemctl restart $SERVICE_NAME && systemctl status $SERVICE_NAME --no-pager"
if [[ $? -ne 0 ]]; then
    echo "Failed to restart service."
    exit 1
fi

# -----------------------------
# Step 6: Ping the API
# -----------------------------
echo "Pinging API..."
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$PING_URL")
echo "API response code: $HTTP_CODE"

if [[ "$HTTP_CODE" -ge 200 ]]; then
    echo "Deployment successful!"
else
    echo "Warning: API returned non-2xx response."
fi