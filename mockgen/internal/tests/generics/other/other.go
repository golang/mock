package other

type One[T any] struct{}

type Two[T any, R any] struct{}

type Three struct{}

type Four struct{}

type Five interface{}

type Either[T, R, K, V any] interface {
	First() T
	Second() R
	Third() K
	Fourth() V
}
