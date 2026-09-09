#!/usr/bin/env bash

set -euo pipefail

CONTAINER_NAME="postgres"
VOLUME="postgres-data"

# Prefer Podman, otherwise use Docker.
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
if ! "$RUNTIME" volume inspect "$VOLUME" >/dev/null 2>&1; then
    echo "Creating volume '$VOLUME'..."
    "$RUNTIME" volume create "$VOLUME" >/dev/null
fi

# Podman needs SELinux relabeling for the init bind mount.
INIT_MOUNT="ro"
if [ "$RUNTIME" = "podman" ]; then
    INIT_MOUNT="ro,Z"
fi

echo "Creating Postgres container using $RUNTIME..."

"$RUNTIME" run -d \
    --name "$CONTAINER_NAME" \
    -v "$VOLUME":/var/lib/postgresql \
    -v "./init:/docker-entrypoint-initdb.d:$INIT_MOUNT" \
    -e POSTGRES_PASSWORD=postgres \
    -p 5432:5432 \
    postgres:latest
