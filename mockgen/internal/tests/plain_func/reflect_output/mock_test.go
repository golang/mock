package mock_plain_func

import (
	"testing"

	"github.com/golang/mock/gomock"
)

// TestPlainFuncMock contains a simple example of how to use mocks that
// were constructed from plain function types.
func TestPlainFuncMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := NewMockFunc(ctrl)
	s.EXPECT().Call(5).Return(7)

	if r := s.Call(5); r != 7 {
		t.Errorf("s.Call(5) == %d; wanted 7", r)
	}
}
