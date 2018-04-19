package gomock

import (
	"reflect"
)

// IsOfType is a Matcher that matches if the string parameter to this function
// accurately describes the type of the parameter to the mock function.
// This is probably best explained with an example.
//
// Examples:
//
// 		addMock.EXPECT().
// 			Insert(gomock.IsOfType("int"), gomock.IsOfType("int"), gomock.IsOfType("*int"))
//
// 		walkerMock.EXPECT().Walk(gomock.IsOfType("Dog"))
//
func IsOfType(typeStr string) Matcher {
	return isOfTypeMatcher{targetTypeStr: typeStr}
}

type isOfTypeMatcher struct {
	targetTypeStr string
}

func (m isOfTypeMatcher) Matches(x interface{}) bool {
	return m.TypeStrOf(x) == m.targetTypeStr
}

func (m isOfTypeMatcher) TypeStrOf(x interface{}) string {
	return m.TypeName(reflect.TypeOf(x))
}

func (m isOfTypeMatcher) TypeName(t reflect.Type) string {
	if t == nil {
		return "nil"
	}
	switch t.Kind() {
	case reflect.Chan:
		return "chan " + m.TypeName(t.Elem())
	case reflect.Ptr:
		return "*" + m.TypeName(t.Elem())
	default:
		return t.Name()
	}
}

func (m isOfTypeMatcher) String() string {
	return "is of type " + m.targetTypeStr
}
