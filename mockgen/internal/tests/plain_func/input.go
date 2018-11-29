//go:generate mockgen -destination reflect_output/mock.go github.com/golang/mock/mockgen/internal/tests/plain_func Func
//go:generate mockgen -source input.go -package plain_func -destination source_output/mock.go
package plain_func

type Func func(int) int
