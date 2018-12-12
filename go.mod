module github.com/golang/mock

require (
	golang.org/x/tools v0.0.0-20190221204921-83362c3779f5

	test.src/a v0.0.0 // Fake - Only used for tests
	test.src/b v0.0.0 // Fake - Only used for tests
)

replace test.src/a => ./mockgen/internal/tests/vendor_dep/vendor/test.src/a

replace test.src/b => ./mockgen/internal/tests/vendor_pkg/vendor/test.src/b
