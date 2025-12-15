#!/bin/bash
set -e

cd "$(dirname "$0")/../.."

echo "Waiting for Temporal server..."
until nc -z localhost 7233 2>/dev/null; do
  sleep 1
done
echo "Temporal server is ready"

# Give worker a moment to start
sleep 2

echo "Starting Ziggy API server with air..."
cd worker
air -c .air.api.toml
