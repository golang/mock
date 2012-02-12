#! /bin/bash -e

mockgen -source gomock/matchers.go \
  > gomock/mock_matcher/mock_matcher.go
mockgen -source sample/user.go \
  -aux_files=imp1=sample/imp1/imp1.go \
  -imports=.=github.com/dsymonds/gomock/sample/imp3,imp_four=github.com/dsymonds/gomock/sample/imp4 \
  > sample/mock_user/mock_user.go

echo >&2 "OK"
