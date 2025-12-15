#!/bin/bash
set -e

echo "Starting Temporal dev server..."
temporal server start-dev \
  --namespace default \
  --db-filename /tmp/ziggy-temporal.db \
  --ui-port 8233
