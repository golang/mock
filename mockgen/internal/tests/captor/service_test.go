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
	"github.com/golang/mock/mockgen/internal/tests/captor/models"
)

// TestAddIDs is an example of how to use an ArgumentCaptor with a slice of int values
func TestAddIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedIDs := []int{1, 4, 253}

	mockDao := NewMockDao(ctrl)
	idCaptor := gomock.AnyCaptor()
	mockDao.EXPECT().InsertIDs(idCaptor)

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

// TestAddIDPointerWithMutation is an example of how to use an ArgumentCaptor with an Eq Matcher.
// In this case the Eq Matcher alone is not enough to test the pointer value's mutation.
// The ArgumentCaptor is used to ensure the mutation happens as expected.
func TestAddIDPointerWithMutation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedID := 100

	mockDao := NewMockDao(ctrl)
	idCaptor := gomock.Captor(gomock.Eq(&expectedID))
	mockDao.EXPECT().InsertIDPointer(idCaptor)

	AddIDPointerWithMutation(mockDao, &expectedID)

	actualID := idCaptor.Value().(*int)
	if *actualID != 101 {
		t.Errorf("expected actualID value to be %d, but got %d", 101, *actualID)
	}
}

// TestAddSportsCarAndSUV is an example of how to use an ArgumentCaptor for a more complex use case.
// In this case an AnyCaptor is used to capture multiple values being passed to the same method, dao.InsertCar.
// AllValues is used to verify that InsertCar was called with both Car values, in the expected order.
func TestAddCars(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDao := NewMockDao(ctrl)
	carCaptor := gomock.AnyCaptor()

	mockDao.EXPECT().InsertCar(carCaptor).Times(2)

	AddCars(mockDao)

	if len(carCaptor.AllValues()) != 2 {
		t.Errorf("expected values length to be %d, but got %d", 2, len(carCaptor.AllValues()))
	}

	for i, val := range carCaptor.AllValues() {
		actualCar := val.(models.Car)
		if i == 0 {
			verifyCar(t, actualCar, false, "red", []models.Seat{models.LeatherSeat, models.LeatherSeat})
		} else {
			verifyCar(t, actualCar, true, "blue", []models.Seat{
				models.ClothSeat, models.ClothSeat, models.ClothSeat, models.ClothSeat, models.ClothSeat})
		}
	}
}

func verifyCar(t *testing.T, actual models.Car, expectedAutomatic bool, expectedColor string, expectedSeats []models.Seat) {
	if expectedAutomatic != actual.Automatic() {
		t.Errorf("expected automatic to be %v, but got %v", expectedAutomatic, actual.Automatic())
	}
	if expectedColor != actual.Color() {
		t.Errorf("expected color to be %s, but got %s", expectedColor, actual.Color())
	}
	if len(expectedSeats) != len(actual.Seats()) {
		t.Errorf("expected %d seats but got %d", len(expectedSeats), len(actual.Seats()))
		return
	}
	for i, seat := range actual.Seats() {
		if seat != expectedSeats[i] {
			t.Errorf("expected seat material to be %v but got %v", seat, expectedSeats[i])
		}
	}
}
