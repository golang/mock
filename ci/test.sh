#!/bin/bash
# This script is used to ensure that the go.mod file is up to date.

set -euo pipefail

if [[ $(go version) != *"go1.18"* ]]; then
  exit 0
fi

for i in $(find $PWD -name go.mod); do
  pushd $(dirname $i)
  go mod tidy
  popd
done

if [ ! -z "$(git status --porcelain)" ]; then
  git status
  git diff
  echo
  echo "The go.mod is not up to date."
  exit 1
fi

BASE_DIR="$PWD"
TEMP_DIR=$(mktemp -d)
function cleanup() {
  rm -rf "${TEMP_DIR}"
}
trap cleanup EXIT

cp -r . "${TEMP_DIR}/"
cd $TEMP_DIR

for i in $(find $PWD -name go.mod); do
  pushd $(dirname $i)
  go generate ./...
  popd
done

if ! diff -r . "${BASE_DIR}"; then
  echo
  echo "The generated files aren't up to date."
  echo "Update them with the 'go generate ./...' command."
  exit 1
fi
