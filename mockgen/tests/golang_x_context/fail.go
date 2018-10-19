package sample

import "github.com/golang/mock/mockgen/tests/golang_x_context/fail"

// FailInterface creates a collision
type FailInterface interface {
	fail.Context
}
