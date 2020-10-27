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
	"fmt"

	"github.com/golang/mock/mockgen/internal/tests/captor/models"
)

//go:generate mockgen -destination mock_dao.go -package captor -source dao.go

type Dao interface {
	InsertIDs(ids []int)
	InsertIDPointer(id *int)
	InsertCar(car models.Car)
}

type realDao struct{}

func (realDao) InsertIDs(ids []int) {
	fmt.Printf("inserting ids %d", ids)
}

func (realDao) InsertIDPointer(id *int) {
	fmt.Printf("inserting ids %v", id)
}

func (realDao) InsertCar(car models.Car) {
	fmt.Printf("inserting car %v", car)
}

func NewDao() Dao {
	return &realDao{}
}
