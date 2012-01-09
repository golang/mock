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


Running mockgen
---------------

The `mockgen` command is used to automatically generate source code for a mock
class given a Go source file containing definitions of interfaces to be mocked.
It supports the following flags:

 *  `--source=...`: Specify the file containing interfaces to be mocked. You
    must supply this flag.

 *  `--destionation=...`: A file to which to write the resulting source
    code. If you don't set this, the code is printed to stdout.

 *  `--packageOut=...`: The name of the package to use for the resulting mock
    class source code. If you don't set this, the package name is `mock_` plus
    the name of the package of the input file.

 *  `--imports=...`: A list of explicit imports that should be used in the
    resulting source code, specified as a comma-separated list of elements of
    the form `foo=bar/baz`, where `bar/baz` is the package being imported and
    `foo` is the identifier to use for the package in the generated source code.

 *  `--aux_files=...`: A list of additional files that should be consulted to
    resolve e.g. embedded interfaces defined in a different file. This is
    specified as a comma-separated list of elements of the form
    `foo=bar/baz.go`, where `bar/baz.go` is the source file and `foo` is the
    identifier to use for the package in the generated source code.

For an example of the use of `mockgen`, see the `sample/` directory. In simple
cases, you will need only the `--source` flag.


TODO: Brief overview of how to create mock objects and set up expectations,
links to documentation, and an example.

[golang]: http://golang.org/
[golang-install]: http://golang.org/doc/install.html#releases
[gomock-ref]: http://gopkgdoc.appspot.com/pkg/github.com/dsymonds/gomock/gomock
