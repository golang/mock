package const_array_length

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/golang/mock/mockgen/internal/tests/const_array_length/consts"
)

func TestMockConstArrayLength_AuxLength(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	auxLength := [consts.C]int{1}
	constArrayLength := NewMockConstArrayLength(ctrl)
	constArrayLength.EXPECT().
		AuxLength().
		Return(auxLength)

	actual := constArrayLength.AuxLength()
	if actual != auxLength {
		t.Fatalf("Expected AuxLenght to be %v, but got %v", auxLength, actual)
	}
}

func TestMockConstArrayLength_PackagePrefixConstLength(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	packagePrefixConstLength := [consts.C]int{1}
	constArrayLength := NewMockConstArrayLength(ctrl)
	constArrayLength.EXPECT().
		PackagePrefixConstLength().
		Return(packagePrefixConstLength)

	actual := constArrayLength.PackagePrefixConstLength()
	if actual != packagePrefixConstLength {
		t.Fatalf("Expected PackagePrefixConstLength to be %v, but got %v", packagePrefixConstLength, actual)
	}
}

func TestMockConstArrayLength_ConstLength(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	constLength := [C]int{1, 2}
	constArrayLength := NewMockConstArrayLength(ctrl)
	constArrayLength.EXPECT().
		ConstLength().
		Return(constLength)

	actual := constArrayLength.ConstLength()
	if actual != constLength {
		t.Fatalf("Expected ConstLenght to be %v, but got %v", constLength, actual)
	}
}

func TestMockConstArrayLength_LiteralLength(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	literalLength := [3]int{1, 2, 3}
	constArrayLength := NewMockConstArrayLength(ctrl)
	constArrayLength.EXPECT().
		LiteralLength().
		Return(literalLength)

	actual := constArrayLength.LiteralLength()
	if actual != literalLength {
		t.Fatalf("Expected LiteralLength to be %v, but go %v", literalLength, actual)
	}
}

func TestMockConstArrayLengthMockRecorder_SliceLength(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sliceLength := []int{1, 2, 3, 4}
	constArrayLength := NewMockConstArrayLength(ctrl)
	constArrayLength.EXPECT().
		SliceLength().
		Return(sliceLength)

	actual := constArrayLength.SliceLength()
	if !reflect.DeepEqual(actual, sliceLength) {
		t.Fatalf("Expected SliceLength to be %v, but got %v", sliceLength, actual)
	}
}
