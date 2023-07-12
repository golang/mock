package core

//go:generate mockgen -package core -self_package go.uber.org/mock/mockgen/internal/tests/self_package -destination mock.go go.uber.org/mock/mockgen/internal/tests/self_package Methods

type Info struct{}

type Methods interface {
	getInfo() Info
}
