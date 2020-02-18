//go:generate mockgen -destination ../source_mock.go -source=source.go -source_package=github.com/golang/mock/mockgen/internal/tests/source_package/definition
package source

type X struct{}

type S interface {
	F(X)
}
