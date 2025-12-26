#!/bin/bash

WATCH_DIR="."
GO_PATTERN=".*\.go$"

while inotifywait -r -e modify,create,delete --format '%f' "$WATCH_DIR" | grep -qE "$GO_PATTERN"; do
    pkill -f go
done
