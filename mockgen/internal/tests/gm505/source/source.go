//go:generate mockgen -package source -destination=../output/source_mock.go -source=source.go

package source

type Foo struct{}

type Bar interface {
	Baz(Foo)
}
