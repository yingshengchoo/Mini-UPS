package model

// Truck represents a shipping truck
type Truck struct {
	ID     TruckID    `gorm:"primaryKey" json:"id"`
	Coord  Coordinate `gorm:"not null;embedded" json:"coord"`
	Status Status     `gorm:"type:varchar(20)" json:"status"`
}

// Coordiante represents the location of an object.
type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type TruckID int
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

// NewTruck creates a new a truck
func NewTruck(truckID TruckID, x int, y int, status Status) *Truck {
	return &Truck{
		ID:     truckID,
		Coord:  Coordinate{X: x, Y: y},
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
