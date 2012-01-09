GoMock is a mocking framework for the [Go programming language][golang]. It
integrates well with Go's built-in `testing` package, but can be used in other
contexts too.


Installation
------------

Once you have [installed Go][golang-install], simply run the following commands
to install the `gomock` package and the `mockgen` tool:

    goinstall github.com/dsymonds/gomock/gomock
    goinstall github.com/dsymonds/gomock/mockgen


Documentation
-------------

After installing, you can use `go doc` to get documentation:

    go doc github.com/dsymonds/gomock/gomock

Alternatively, there is an online reference for the package hosted on GoPkgDoc
[here][gomock-ref].


TODO: How to run mockgen, brief overview of how to create mock objects and set
up expectations, links to documentation, and an example.

[golang]: http://golang.org/
[golang-install]: http://golang.org/doc/install.html#releases
[gomock-ref]: http://gopkgdoc.appspot.com/pkg/github.com/dsymonds/gomock/gomock
