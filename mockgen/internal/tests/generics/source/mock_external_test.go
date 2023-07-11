package source

import (
	"context"
	"testing"

	"go.uber.org/mock/gomock"
	"go.uber.org/mock/mockgen/internal/tests/generics"
)

func TestMockEmbeddingIface_One(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockEmbeddingIface[int, float64](ctrl)
	m.EXPECT().One("foo").Return("bar")
	if v := m.One("foo"); v != "bar" {
		t.Errorf("One() = %v, want %v", v, "bar")
	}
}

func TestMockUniverse_Water(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockUniverse[int](ctrl)
	m.EXPECT().Water(1024)
	m.Water(1024)
}

func TestNewMockGroup_Join(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockGroup[generics.Generator[any]](ctrl)
	ctx := context.TODO()
	m.EXPECT().Join(ctx).Return(nil)
	if v := m.Join(ctx); v != nil {
		t.Errorf("Join() = %v, want %v", v, nil)
	}
}
