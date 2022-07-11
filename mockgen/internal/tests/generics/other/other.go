package other

type One[T any] struct{}

type Two[T any, R any] struct{}

type Three struct{}

type Four struct{}

type Five interface{}

type Otherer[T any, R any] interface {
	DoT(T) error
	DoR(R) error
	MakeThem() (T, R, error)
}
