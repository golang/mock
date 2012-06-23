#! /bin/bash -e

mockgen code.google.com/p/gomock/gomock Matcher \
  > gomock/mock_matcher/mock_matcher.go
mockgen code.google.com/p/gomock/sample Index,Embed,Embedded \
  > sample/mock_user/mock_user.go
gofmt -w gomock/mock_matcher/mock_matcher.go sample/mock_user/mock_user.go

echo >&2 "OK"
