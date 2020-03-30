package captor

import (
	"testing"

	"github.com/golang/mock/gomock"
)

func TestAddIdsWithAnyCaptor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedIDs := []int{1, 4, 253}

	mockDao := NewMockDao(ctrl)
	idCaptor := gomock.AnyCaptor()
	mockDao.EXPECT().InsertIDs(idCaptor).Times(1)

	AddIDs(mockDao, expectedIDs)

	actualIDs := idCaptor.LastValue().([]int)
	if len(expectedIDs) != len(actualIDs) {
		t.Errorf("expected ids length to be %d, but got %d", len(expectedIDs), len(actualIDs))
	}
	for i, expectedID := range expectedIDs {
		if expectedID != actualIDs[i] {
			t.Errorf("expected id to be %d, but got %d", expectedID, actualIDs[i])
		}
	}
}
