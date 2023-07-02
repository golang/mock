package gomock_test

import (
	"testing"

	"go.uber.org/mock/gomock"
)

func TestEcho_NoOverride(t *testing.T) {
	ctrl := gomock.NewController(t, gomock.WithOverridableExpectations())
	mockIndex := NewMockFoo(ctrl)

	mockIndex.EXPECT().Bar(gomock.Any()).Return("foo")
	res := mockIndex.Bar("input")

	if res != "foo" {
		t.Fatalf("expected response to equal 'foo', got %s", res)
	}
}

func TestEcho_WithOverride_BaseCase(t *testing.T) {
	ctrl := gomock.NewController(t, gomock.WithOverridableExpectations())
	mockIndex := NewMockFoo(ctrl)

	// initial expectation set
	mockIndex.EXPECT().Bar(gomock.Any()).Return("foo")
	// override
	mockIndex.EXPECT().Bar(gomock.Any()).Return("bar")
	res := mockIndex.Bar("input")

	if res != "bar" {
		t.Fatalf("expected response to equal 'bar', got %s", res)
	}
}
