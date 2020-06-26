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
