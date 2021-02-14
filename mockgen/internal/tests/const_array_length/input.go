//go:generate mockgen -package const_length -destination mock.go -source input.go
package const_length

const C = 2

type I interface {
	F() [C]int
}
