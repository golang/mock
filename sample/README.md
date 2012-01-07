This directory contains an example of a package containing a non-trivial
interface that can be mocked with GoMock. The interesting files are:

 *  `user.go`: Source code for the sample package, containing interfaces to be
    mocked. This file depends on the packages named imp[1-4] for various things.

 *  `user_test.go`: A pretend test for the sample package, in which mocks of the
    interfaces from `user.go` are used. This demonstrates how to create mock
    objects, set up expectations, and so on.

If you want to build the sample and run the test, it's best to clone the GoMock
git repository into your `$GOPATH` directory so that the `go` command can deal
with its dependencies:

    git clone https://github.com/dsymonds/gomock.git $GOPATH/src/github.com/dsymonds/gomock

You can build the sample package as follows:

    go build github.com/dsymonds/gomock/sample

To run the test, you'll need to first use MockGen to generate the `mock_user`
package used by the test:

    cd $GOPATH/src/github.com/dsymonds/gomock/sample
    mkdir -p mock_user
    mockgen --source=user.go --aux_files=imp1=imp1/imp1.go --imports=.=github.com/dsymonds/gomock/sample/imp3,imp_four=github.com/dsymonds/gomock/sample/imp4 > mock_user/mock_user.go

You can now verify that the mock package builds:

    go build github.com/dsymonds/gomock/sample/mock_user

You can invoke the following command to run the tests in `user_test`.go:

    go test
