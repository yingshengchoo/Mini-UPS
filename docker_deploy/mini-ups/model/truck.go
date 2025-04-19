package model

type Truck struct {
	ID     uint   `gorm:"primary_key" json:"id"`
	CoordX int    `gorm:"not null" json:"coordx"`
	CoordY int    `gorm:"not null" json:"coordy"`
	Status Status `gorm:"type:varchar(20)" json:"status"`
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

/**
 * transit status
 * use truck status as parameter
 */
func (t *Truck) Transit(status Status) {
	t.Status = status
}
