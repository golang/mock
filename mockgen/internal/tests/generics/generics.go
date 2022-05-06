package generics

import (
	"github.com/golang/mock/mockgen/internal/tests/generics/other"
)

//go:generate mockgen --source=generics.go --destination=source/mock_generics_test.go --package source
////go:generate mockgen --destination=reflect/mock_test.go --package reflect . Bar,Bar2

type Bar[T any, R any] interface {
	One(string) string
	Two(T) string
	Three(T) R
	Four(T) Foo[T, R]
	Five(T) Baz[T]
	Six(T) *Baz[T]
	Seven(T) other.One[T]
	Eight(T) other.Two[T, R]
	Nine(Iface[T])
	Ten(*T)
	Eleven() (*other.One[T], error)
	Twelve() (*other.Two[T, R], error)
}

type Foo[T any, R any] struct{}

type Baz[T any] struct{}

type Iface[T any] interface{}
