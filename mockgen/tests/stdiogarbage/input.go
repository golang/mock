//go:generate mockgen -destination output.go -package test github.com/golang/mock/mockgen/tests/stdiogarbage I
package test

import "fmt"

func init() {
	fmt.Print("garbage")
}

type I interface{}
