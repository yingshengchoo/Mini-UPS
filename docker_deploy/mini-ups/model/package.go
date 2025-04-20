package model

import (
	"gorm.io/datatypes"
)

type Package struct {
	ID          PackageID      `gorm:"primaryKey" json:"package_id"`
	UserID      uint           `gorm:"not null" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID"`
	TruckID     *TruckID       `json:"truck_id"` // nullable
	Truck       *Truck         `gorm:"foreignKey:TruckID"`
	Items       datatypes.JSON `json:"items"`
	Destination Coordinate     `gorm:"not null;embedded" json:"coord"`
	WarehouseID uint           `gorm:"not null" json:"warehouse_id"`
	Status      PackageStatus  `gorm:"type:varchar(20)" json:"status"`
}

type PackageID string
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

func (p *Package) SetCoord(x int, y int) {
	p.Destination.X = x
	p.Destination.Y = y
}

func NewPackage(packageID PackageID, userID uint, items string, x int, y int, warehouseID uint, status PackageStatus) *Package {
	return &Package{
		ID:          packageID,
		UserID:      userID,
		TruckID:     nil, // or pointer to uint if assigned
		Items:       datatypes.JSON([]byte(items)),
		Destination: Coordinate{X: x, Y: y},
		WarehouseID: warehouseID,
		Status:      status, // use constant if defined
	}
}
