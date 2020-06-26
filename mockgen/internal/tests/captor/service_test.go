// Copyright 2020 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

	actualIDs := idCaptor.Value().([]int)
	if len(expectedIDs) != len(actualIDs) {
		t.Errorf("expected ids length to be %d, but got %d", len(expectedIDs), len(actualIDs))
	}
	for i, expectedID := range expectedIDs {
		if expectedID != actualIDs[i] {
			t.Errorf("expected id to be %d, but got %d", expectedID, actualIDs[i])
		}
	}
}
