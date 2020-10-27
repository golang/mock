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

import "github.com/golang/mock/mockgen/internal/tests/captor/models"

func AddIDs(dao Dao, ids []int) {
	dao.InsertIDs(ids)
}

func AddIDPointerWithMutation(dao Dao, id *int) {
	*id += 1
	dao.InsertIDPointer(id)
}

func AddCars(dao Dao) {
	sportsCar := models.NewCar(
		false,
		"red",
		[]models.Seat{models.LeatherSeat, models.LeatherSeat})

	suv := models.NewCar(
		true,
		"blue",
		[]models.Seat{models.ClothSeat, models.ClothSeat, models.ClothSeat, models.ClothSeat, models.ClothSeat})

	cars := []models.Car{sportsCar, suv}

	for _, car := range cars {
		dao.InsertCar(car)
	}
}
