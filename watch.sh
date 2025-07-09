#!/usr/bin/env bash

# Watch for Go file changes, rebuild, and kill the running process
build_failed=false
while true; do
    echo "Building..."
    # Build the application
    if go build -o sift; then
        if [ "$build_failed" = true ]; then
            echo -e "\033[1;30;42m  BUILD SUCCEEDED  \033[0m"
            afplay /System/Library/Sounds/Glass.aiff
            build_failed=false
        else
            echo -e "\033[1;40m  BUILD SUCCEEDED  \033[0m"
        fi
        pkill -f './sift'
    else
        echo -e "\033[1;30;41m  BUILD FAILED  \033[0m"
        afplay /System/Library/Sounds/Bottle.aiff
        build_failed=true
    fi
    
    # Watch for changes to any Go file
    fswatch -1 -l 0.3 -e ".*" -i "\.go$" . || exit
done
