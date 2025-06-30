#!/usr/bin/env bash

# Watch for Go file changes, rebuild, and kill the running process
while true; do
    echo "Building..."
    # Build the application
    go build -o sift && pkill -f './sift'
    
    # Watch for changes to any Go file
    fswatch -1 -l 0.3 -e ".*" -i "\.go$" . || exit
done
