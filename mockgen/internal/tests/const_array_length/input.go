package const_length

//go:generate mockgen -package const_length -destination mock.go -source input.go

const C = 2

type I interface {
	Foo() [C]int
	Bar() [2]int
}
