package sample

import "github.com/golang/mock/mockgen/tests/golang_x_context/pass"

// PassInterface parses correctly
type PassInterface interface {
	pass.Context
}
