#!/bin/sh

# Start worker in background
./ziggy worker &

# Start API server in foreground
./ziggy serve
