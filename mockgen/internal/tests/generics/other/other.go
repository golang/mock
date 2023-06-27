package other

type One[T any] struct{}

type Two[T any, R any] struct{}

type Three struct{}

type Four struct{}

type Five interface{}

type Twenty[R, S, T any, Z any] interface {
	Twenty(S, R) (T, Z)
}

type TwentyThree[U, V any] interface {
	TwentyThree(U, V) StructType
}

type StructType struct{}
