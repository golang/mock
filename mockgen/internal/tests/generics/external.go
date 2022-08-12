package generics

import (
	"context"

	"github.com/golang/mock/mockgen/internal/tests/generics/other"
	"golang.org/x/exp/constraints"
)

//go:generate mockgen --source=external.go --destination=source/mock_external_test.go --package source

type ExternalConstraint[I constraints.Integer, F constraints.Float] interface {
	One(string) string
	Two(I) string
	Three(I) F
	Four(I) Foo[I, F]
	Five(I) Baz[F]
	Six(I) *Baz[F]
	Seven(I) other.One[I]
	Eight(F) other.Two[I, F]
	Nine(Iface[I])
	Ten(*I)
}

type EmbeddingIface[T constraints.Integer, R constraints.Float] interface {
	other.Twenty[T, StructType, R, other.Five]
	TwentyTwo[StructType]
	other.TwentyThree[TwentyTwo[R], TwentyTwo[T]]
	TwentyFour[other.StructType]
	Foo() error
	ExternalConstraint[T, R]
}

type TwentyOne[T any] interface {
	TwentyOne() T
}

type TwentyFour[T other.StructType] interface {
	TwentyFour() T
}

type Clonable[T any] interface {
	Clone() T
}

type Finder[T Clonable[T]] interface {
	Find(ctx context.Context) ([]T, error)
}

type UpdateNotifier[T any] interface {
	NotifyC(ctx context.Context) <-chan []T

	Refresh(ctx context.Context)
}

type EmbeddedW[W StructType] interface {
	EmbeddedY[W]
}

type EmbeddedX[X StructType] interface {
	EmbeddedY[X]
}

type EmbeddedY[Y StructType] interface {
	EmbeddedZ[Y]
}

type EmbeddedZ[Z any] interface {
	EmbeddedZ(Z)
}
