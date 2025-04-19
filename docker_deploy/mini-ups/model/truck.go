package model

type Truck struct {
	ID     uint       `gorm:"primary_key" json:"id"`
	Coord  Coordinate `gorm:"not null" json:"coord"`
	Status Status     `gorm:"type:varchar(20)" json:"status"`
}

type Coordinate struct {
	X int
	Y int
}

type Status string

// truck statuses
var TruckStatus = struct {
	IDLE       Status
	PICKING    Status
	ARRIVED    Status
	LOADED     Status
	DELIVERING Status
	DELIVERED  Status
}{
	IDLE:       "idle",
	PICKING:    "picking",
	ARRIVED:    "arrived",
	LOADED:     "loaded",
	DELIVERING: "delivering",
	DELIVERED:  "delivered",
}

// new a truck
func NewTruck(coord Coordinate, status Status) *Truck {
	return &Truck{
		Coord:  coord,
		Status: status,
	}
}

/**
 * transit status
 * use truck status as parameter
 */
func (t *Truck) Transit(status Status) {
	t.Status = status
}

// set coord
func (t *Truck) SetCoord(x int, y int) {
	t.Coord.X = x
	t.Coord.Y = y
}
