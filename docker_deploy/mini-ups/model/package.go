package model

import (
	"gorm.io/datatypes"
)

type Package struct {
	ID      string `gorm:"primaryKey" json:"package_id"`
	UserID  uint   `gorm:"not null" json:"user_id"`
	User    User   `gorm:"foreignKey:UserID"`
	TruckID *uint  `json:"truck_id"` // nullable
	//Truck       *Truck         `gorm:"foreignKey:TruckID"`
	Items       datatypes.JSON `json:"items"`
	Destination Coordinate     `gorm:"not null;embedded" json:"coord"`
	WarehouseID int            `gorm:"not null" json:"warehouse_id"`
	Status      PackageStatus  `gorm:"type:varchar(20)" json:"status"`
}

type PackageStatus string

const (
	StatusCreated          PackageStatus = "created"
	StatusWaitingForPickup PackageStatus = "waiting_for_pickup"
	StatusPickupAssigned   PackageStatus = "pickup_assigned"
	StatusReadyForPickup   PackageStatus = "ready_for_pickup"
	StatusPickupComplete   PackageStatus = "pickup_complete"
	StatusLoading          PackageStatus = "loading"
	StatusLoaded           PackageStatus = "loaded"
	StatusOutForDelivery   PackageStatus = "out_for_delivery"
	StatusDelivered        PackageStatus = "delivered"
)

type Item struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Quantity    int    `json:"quantity"`
}

func (p *Package) SetCoord(x int, y int) {
	p.Destination.X = x
	p.Destination.Y = y
}
