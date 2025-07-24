#!/usr/bin/env bash

# This script is run as a pre-push hook to avoid pushing tags that does not
# match the version defined in the code.
#
# This is the pre-push script in .git/hooks/pre-push:
#
#   #!/bin/bash
#   set -euo pipefail
#
#   # Capture pushed refs from stdin
#   while read -r _ _ remote_ref _; do
#     if [[ "$remote_ref" =~ refs/tags/(v[0-9]+\.[0-9]+\.[0-9]+) ]]; then
#       tag="${BASH_REMATCH[1]}"
#
#       echo "🔍 Checking version for tag: $tag"
#       .github/scripts/check-version.sh "$tag"
#     fi
#   done

set -euo pipefail

TAG=${1:-}

if [[ -z "$TAG" ]]; then
    echo "Error: No tag provided."
    exit 1
fi

# Look for the tag string in the file that defines the version
if ! grep -q "wsl version $TAG" ./*.go; then
    echo "❌ Version constant does not match tag: $TAG"
    exit 1
fi

echo "✅ Version constant matches tag: $TAG"
