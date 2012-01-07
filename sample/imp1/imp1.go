package imp1

import "bufio"

type Imp1 struct{}

type ForeignEmbedded interface {
	// The return value here also makes sure that
	// the generated mock picks up the "bufio" import.
	ForeignEmbeddedMethod() *bufio.Reader
}
