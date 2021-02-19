//go:generate mockgen -package source -destination=../output/source_mock.go -source=source.go

// Package source makes sure output imports its. See #505.
package source

type Foo struct{}

type Bar interface {
	Baz(Foo)
}
