#!/bin/bash
# This script is used by CI to check if the code passes golint.

set -euo pipefail

GOLINT_OUTPUT=$(IFS=$'\n'; golint ./... | grep -v "mockgen/internal/.*\|sample/.*")
if [[ -n "${GOLINT_OUTPUT}" ]]; then
    echo "${GOLINT_OUTPUT}"
    echo
    echo "The go source files aren't passing golint."
    exit 1
fi

