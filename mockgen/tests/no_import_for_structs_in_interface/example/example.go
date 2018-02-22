//go:generate mockgen -destination ../bugreport_mock.go -package bugreport -source=example.go -source_package=github.com/golang/mock/mockgen/tests/no_import_for_structs_in_interface/example

package example

type ParamType struct {
	Value string
}

type ReturnType struct {
	Value string
}

type Example interface {
	Method(param ParamType) ReturnType
}
