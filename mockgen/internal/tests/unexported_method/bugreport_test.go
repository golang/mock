package bugreport

import (
	"testing"

	"go.uber.org/mock/gomock"
)

func TestCallExample(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := NewMockExample(ctrl)
	e.EXPECT().someMethod(gomock.Any()).Return("it works!")

	CallExample(e)
}
