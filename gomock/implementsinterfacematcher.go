package gomock

import (
	"io"
	"reflect"
)

// ImplementsInterface is a Matcher that matches if the string parameter to this
// function is registered (see RegisterType) and the parameter to the mock
// function implements this interface
//
// Examples:
//
// 		writeFooMock.EXPECT().WriteFoo(gomock.ImplementsInterface("io.Writer"))
//
func ImplementsInterface(x interface{}) Matcher {
	switch v := x.(type) {
	case string:
		return implementsInterfaceMatcher{reflectType: stringToTypeMap[v]}
	case reflect.Type:
		return implementsInterfaceMatcher{reflectType: v}
	default:
		return nil
	}
}

func RegisterType(s string, rt reflect.Type) {
	stringToTypeMap[s] = rt
}

var stringToTypeMap = map[string]reflect.Type{
	"error":         reflect.ValueOf(new(error)).Type().Elem(),
	"io.Reader":     reflect.ValueOf(new(io.Reader)).Type().Elem(),
	"io.Writer":     reflect.ValueOf(new(io.Writer)).Type().Elem(),
	"io.ReadWriter": reflect.ValueOf(new(io.ReadWriter)).Type().Elem(),
}

type implementsInterfaceMatcher struct {
	reflectType reflect.Type
}

func (m implementsInterfaceMatcher) Matches(x interface{}) bool {
	return m.reflectType != nil && x != nil && reflect.TypeOf(x).Implements(m.reflectType)
}

func (m implementsInterfaceMatcher) String() string {
	return "implements interface " + m.reflectType.Name()
}
