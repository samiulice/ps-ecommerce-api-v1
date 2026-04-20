#!/usr/bin/env bash

# Usage:
#   ./killport.sh        (prompt for port)
#   ./killport.sh 8080   (pass port directly)

PORT=$1

# Ask for port if not provided
if [ -z "$PORT" ]; then
  read -p "Enter port: " PORT
fi

# Detect OS
OS_RAW="$(uname)"

case "$OS_RAW" in
  Linux)
    OS="Linux"
    ;;
  Darwin)
    OS="macOS"
    ;;
  MINGW*|MSYS*|CYGWIN*)
    OS="Windows"
    ;;
  *)
    OS="Unknown"
    ;;
esac

echo "🖥️ OS: $OS"
echo "🔍 Checking port: $PORT"

# Function to kill PID
kill_pid() {
  PID=$1
  if [ -n "$PID" ]; then
    echo "⚠️ Killing PID: $PID"
    kill -9 "$PID" 2>/dev/null && echo "✅ Killed $PID" || echo "❌ Failed to kill $PID"
  fi
}

case "$OS" in
  Linux|macOS)
    PIDS=$(lsof -t -i :"$PORT")

    if [ -z "$PIDS" ]; then
      echo "✅ No process found on port $PORT"
    else
      for PID in $PIDS; do
        kill_pid "$PID"
      done
    fi
    ;;

  Windows)
    PIDS=$(netstat -ano | grep ":$PORT" | awk '{print $5}' | sort -u)

    if [ -z "$PIDS" ]; then
      echo "✅ No process found on port $PORT"
    else
      for PID in $PIDS; do
        echo "⚠️ Killing PID: $PID"
        taskkill //PID "$PID" //F >/dev/null 2>&1 && echo "✅ Killed $PID" || echo "❌ Failed to kill $PID"
      done
    fi
    ;;

  *)
    echo "❌ Unsupported OS"
    exit 1
    ;;
esac