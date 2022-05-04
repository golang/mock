package generics

//go:generate mockgen --source=generics.go --destination=source/mock_test.go --package source
////go:generate mockgen --destination=reflect/mock_test.go --package reflect . Bar,Bar2

type Bar[T any, R any] interface {
	One(string) string
	Two(T) string
	Three(T) R
	//Four(T) Foo[T, R]
	//Five(T) Baz[T]
	//Six(T) *Baz[T]
}

type Foo[T any, R any] struct{}

type Baz[T any] struct{}
