#!/usr/bin/env bash

set -euo pipefail

# Check if we have uncommitted changes
STATUS_OUTPUT=$(jj status)
if ! echo "$STATUS_OUTPUT" | grep -q "The working copy has no changes."; then
    echo "Error: Working copy has uncommitted changes. Please commit or stash them first."
    exit 1
fi

# Get the current version from git-cliff
echo "Generating changelog..."
git-cliff --bump -o CHANGELOG.md

# Check if CHANGELOG.md was actually modified
STATUS_OUTPUT=$(jj status)
if echo "$STATUS_OUTPUT" | grep -q "CHANGELOG.md"; then
    # Stage and commit the changelog
    echo "Committing changelog..."
    jj commit -m "chore: Update changelog"
    
    # Get the version from the changelog (extract from the first ## line)
    VERSION=$(grep -m 1 "^## \[" CHANGELOG.md | sed 's/^## \[\(.*\)\] - .*/\1/')
    
    if [ -z "$VERSION" ]; then
        echo "Error: Could not extract version from changelog"
        exit 1
    fi
    
    echo "Creating tag v$VERSION..."
    git tag "v$VERSION"
    
    echo "Setting main bookmark to current revision..."
    jj bookmark set main
    
    echo "Pushing main and tag to GitHub..."
    jj git push --bookmark main
    git push origin "v$VERSION"
    
    echo "Release v$VERSION completed successfully!"
else
    echo "No changes to changelog detected. Release not needed."
fi
