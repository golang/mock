package other

type One[T any] struct{}

type Two[T any, R any] struct{}

type Three struct{}

type Four struct{}

type Five interface{}

type Twenty[T any] interface {
	Twenty() T
}
