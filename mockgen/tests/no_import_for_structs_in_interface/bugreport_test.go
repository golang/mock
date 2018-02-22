package bugreport

import (
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/golang/mock/mockgen/tests/no_import_for_structs_in_interface/example"
)

func TestExample_Method(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := NewMockExample(ctrl)

	p := example.ParamType{"a"}
	r := example.ReturnType{"a"}

	m.EXPECT().Method(p).Return(r)
	if m.Method(p).Value != r.Value {
		t.Fail()
	}
}
