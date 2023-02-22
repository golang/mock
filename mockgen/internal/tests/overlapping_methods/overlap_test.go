package overlap

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

// TestValidInterface assesses whether or not the generated mock is valid
func TestValidInterface(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := NewMockReadWriteCloser(ctrl)
	s.EXPECT().Close().Return(errors.New("test"))

	s.Close()
}
