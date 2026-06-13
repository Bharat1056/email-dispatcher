#!/usr/bin/env bash

# Exit early if there are no staged Go files
if ! git diff --cached --name-only | grep -q '\.go$'; then
  exit 0
fi

# Create a temporary patch file of staged Go files
git diff --cached -- "*.go" > staged.patch

# If patch is empty, clean up and exit
if [ ! -s staged.patch ]; then
  rm -f staged.patch
  exit 0
fi

# Run golangci-lint using the patch file to only show errors in staged changes
golangci-lint run --new-from-patch staged.patch
exit_code=$?

# Clean up
rm -f staged.patch

exit $exit_code
