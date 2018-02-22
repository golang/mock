Types that appear in the signature of a mocked interface are not imported
properly for the mock interfaces.

Example:

```go
package example

type ParamType struct {
	value string
}

type ReturnType struct {
	value string
}

type Example interface {
	Method(param ParamType) ReturnType
}
```

```go
package mock_example

// Method mocks base method
func (m *MockExample) Method(param ParamType) ReturnType {
	ret := m.ctrl.Call(m, "Method", param)
	ret0, _ := ret[0].(ReturnType)
	return ret0
}
```

In the above example, `ParamType` and `ReturnType` are exported structs in the package
of the interface being mocked, but are not imported into the generated mock file,
and are not reference as `example.ParamType` and `example.ReturnType`.
