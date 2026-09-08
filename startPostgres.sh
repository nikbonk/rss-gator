#!/usr/bin/env bash

set -euo pipefail

CONTAINER_NAME="postgres"
VOLUME="postgres-data"

if command -v podman >/dev/null 2>&1; then
    RUNTIME="podman"
elif command -v docker >/dev/null 2>&1; then
    RUNTIME="docker"
else
    echo "Error: neither podman nor docker is installed."
    exit 1
fi

# Already running? Nothing to do.
if "$RUNTIME" ps --format '{{.Names}}' | grep -qx "$CONTAINER_NAME"; then
    echo "Postgres is already running."
    exit 0
fi

# Container exists but is stopped? Start it.
if "$RUNTIME" container exists "$CONTAINER_NAME" 2>/dev/null; then
    echo "Starting existing Postgres container..."
    "$RUNTIME" start "$CONTAINER_NAME"
    exit 0
fi

# Create volume if needed.
"$RUNTIME" volume inspect "$VOLUME" >/dev/null 2>&1 ||
    "$RUNTIME" volume create "$VOLUME" >/dev/null

echo "Creating Postgres container..."

"$RUNTIME" run -d \
    --name "$CONTAINER_NAME" \
    -v "$VOLUME":/var/lib/postgresql \
    -e POSTGRES_PASSWORD=postgres \
    -p 5432:5432 \
    postgres:latest
