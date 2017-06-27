#!/bin/bash
# This script is used by the CI to check if the code is gofmt formatted.

set -euo pipefail
cd "$( dirname "$0" )"

GOFMT_DIFF=$( gofmt -d $( IFS='\n'; find . -type f -name '*.go' ) )
if [[ -n "${GOFMT_DIFF}" ]]; then
    echo "${GOFMT_DIFF}"
    echo
    echo "The go source files aren't gofmt formatted."
    exit 1
fi
