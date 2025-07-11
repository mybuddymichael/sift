#!/usr/bin/env bash

# Run app continuously with restart on exit
while true; do
    ./sift
    
    # Add small delay to prevent death spin if app crashes immediately
    sleep 0.25
done
