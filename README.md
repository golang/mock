gomock [![Build Status](https://travis-ci.org/golang/mock.svg?branch=master)](https://travis-ci.org/golang/mock)
======

GoMock is a mocking framework for the [Go programming language][golang]. It
integrates well with Go's built-in `testing` package, but can be used in other
contexts too.


Installation
------------

Once you have [installed Go][golang-install], run these commands
to install the `gomock` package and the `mockgen` tool:

    go get github.com/golang/mock/gomock
    go get github.com/golang/mock/mockgen


Documentation
-------------

After installing, you can use `go doc` to get documentation:

    go doc github.com/golang/mock/gomock

Alternatively, there is an online reference for the package hosted on GoPkgDoc
[here][gomock-ref].


Running mockgen
---------------

`mockgen` has two modes of operation: source and reflect.
Source mode generates mock interfaces from a source file.
It is enabled by using the -source flag. Other flags that
may be useful in this mode are -imports and -aux_files.

Example:

	mockgen -source=foo.go [other options]

Reflect mode generates mock interfaces by building a program
that uses reflection to understand interfaces. It is enabled
by passing two non-flag arguments: an import path, and a
comma-separated list of symbols.

Example:

	mockgen database/sql/driver Conn,Driver

The `mockgen` command is used to generate source code for a mock
class given a Go source file containing interfaces to be mocked.
It supports the following flags:

 *  `-source`: A file containing interfaces to be mocked.

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

*  `-build_flags`: (reflect mode only) Flags passed verbatim to `go build`.

* `-mock_names`: A list of custom names for generated mocks. This is specified 
	as a comma-separated list of elements of the form
	`Repository=MockSensorRepository,Endpoint=MockSensorEndpoint`, where 
	`Repository` is the interface name and `MockSensorRepository` is the desired
	mock name (mock factory method and mock recorder will be named after the mock).
	If one of the interfaces has no custom name specified, then default naming
	convention will be used.

For an example of the use of `mockgen`, see the `sample/` directory. In simple
cases, you will need only the `-source` flag.


TODO: Brief overview of how to create mock objects and set up expectations, and
an example.

[golang]: http://golang.org/
[golang-install]: http://golang.org/doc/install.html#releases
[gomock-ref]: http://godoc.org/github.com/golang/mock/gomock

## grpc

See the [./sample/grpc](./sample/grpc) directory. **NB:** we use the reflect instead of source mode.

```
$ cd sample/grpc && make
protoc -I . --go_out=plugins=grpc:. ./route_guide.proto
mockgen github.com/banaio/mock/sample/grpc RouteGuideClient,RouteGuide_ListFeaturesClient,RouteGuide_RecordRouteClient,RouteGuide_RouteChatClient > mock/route_guide.go

go test ./...
ok  	github.com/banaio/mock/sample/grpc	(cached)
?   	github.com/banaio/mock/sample/grpc/mock	[no test files]
```

If you see errors like below, please try the reflect mode:

```sh
mockgen -source=./route_guide.pb.go
2018/04/01 19:44:38 imported package collision: "context" imported twice
make: *** [all] Error 1
```