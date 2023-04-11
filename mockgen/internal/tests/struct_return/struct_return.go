// Package extra_import makes sure output does not import it. See #515.
package struct_return

//go:generate mockgen -destination mock.go -package struct_return . Foo

type Message struct {
	Text string
}

type Foo interface {
	Bar() struct {
		Value  string
		Nested struct {
			Value int
		}
	}
}
