package mock_embed_test

import (
	reflect "reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/golang/mock/mockgen/internal/tests/embed"
	mock_embed "github.com/golang/mock/mockgen/internal/tests/embed/mock"
)

func TestEmbed(t *testing.T) {
	hoge := mock_embed.NewMockHoge(gomock.NewController(t))
	et := reflect.TypeOf((*embed.Hoge)(nil)).Elem()
	ht := reflect.TypeOf(hoge)
	if !ht.Implements(et) {
		t.Errorf("source interface has been not implemented")
	}
}
