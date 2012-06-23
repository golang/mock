GoMock is a mocking framework for the [Go programming language][golang]. It
integrates well with Go's built-in `testing` package, but can be used in other
contexts too.


Installation
------------

Once you have [installed Go][golang-install], run these commands
to install the `gomock` package and the `mockgen` tool:

    go install code.google.com/p/gomock/gomock
    go install code.google.com/p/gomock/mockgen


Documentation
-------------

After installing, you can use `go doc` to get documentation:

    go doc code.google.com/p/gomock/gomock

Alternatively, there is an online reference for the package hosted on GoPkgDoc
[here][gomock-ref].


Running mockgen
---------------

The `mockgen` command is used to generate source code for a mock
class given a Go source file containing interfaces to be mocked.
It supports the following flags:

 *  `-source`: The file containing interfaces to be mocked. You must
    supply this flag.

 *  `-destination`: A file to which to write the resulting source code. If you
    don't set this, the code is printed to standard output.

 *  `-package`: The package to use for the resulting mock class
    source code. If you don't set this, the package name is `mock_` concatenated
    with the package of the input file.

 *  `-imports`: A list of explicit imports that should be used in the resulting
    source code, specified as a comma-separated list of elements of the form
    `foo=bar/baz`, where `bar/baz` is the package being imported and `foo` is
    the identifier to use for the package in the generated source code.

 *  `-aux_files`: A list of additional files that should be consulted to
    resolve e.g. embedded interfaces defined in a different file. This is
    specified as a comma-separated list of elements of the form
    `foo=bar/baz.go`, where `bar/baz.go` is the source file and `foo` is the
    package name of that file used by the -source file.

For an example of the use of `mockgen`, see the `sample/` directory. In simple
cases, you will need only the `-source` flag.


TODO: Brief overview of how to create mock objects and set up expectations, and
an example.

[golang]: http://golang.org/
[golang-install]: http://golang.org/doc/install.html#releases
[gomock-ref]: http://gopkgdoc.appspot.com/pkg/code.google.com/p/gomock/gomock
