#!/bin/bash
set -e

cd "$(dirname "$0")/../../web"

echo "Starting web dev server..."
npm run dev
