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

package models

type Car interface {
	Automatic() bool
	Color() string
	Seats() []Seat
}

func NewCar(automatic bool, color string, seats []Seat) Car {
	return &car{
		automatic: automatic,
		color:     color,
		seats:     seats,
	}
}

type car struct {
	automatic bool
	color     string
	seats     []Seat
}

func (c *car) Automatic() bool {
	return c.automatic
}

func (c *car) Color() string {
	return c.color
}

func (c *car) Seats() []Seat {
	return c.seats
}

type Seat struct {
	Material string
	Color    string
}

var (
	ClothSeat = Seat{
		Material: "cloth",
		Color:    "brown",
	}

	LeatherSeat = Seat{
		Material: "leather",
		Color:    "black",
	}
)
