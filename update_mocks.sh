#! /bin/bash -e

mockgen -source gomock/matchers.go \
  > gomock/mock_matcher/mock_matcher.go

echo >&2 "OK"
