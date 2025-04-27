package model

import (
	"time"

	"gorm.io/datatypes"
)

// Package is a model representing a package of a user.
type Package struct {
	ID            PackageID      `gorm:"primaryKey" json:"package_id"`
	Username      string         `gorm:"not null" json:"username"`
	User          User           `gorm:"foreignKey:Username;references:Username"`
	TruckID       *TruckID       `json:"truck_id"` // nullable
	Truck         *Truck         `gorm:"foreignKey:TruckID"`
	Items         datatypes.JSON `json:"items"`
	Destination   Coordinate     `gorm:"not null;embedded" json:"coord"`
	WarehouseID   uint           `gorm:"not null" json:"warehouse_id"`
	Status        PackageStatus  `gorm:"type:varchar(20)" json:"status"`
	UpdatedAt     time.Time      `json:"updated_at"`
	IsPrioritized bool           `gorm:"default:false" json:"is_prioritized"`
}

type PackageID string
type PackageStatus string

// These are the different status which a package can have.
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

// SetCoords sets the destination coordinates of a package.
func (p *Package) SetCoord(x int, y int) {
	p.Destination.X = x
	p.Destination.Y = y
}

// NewPackage creates a new package given the fields. By default they are not prioritized.
func NewPackage(packageID PackageID, username string, items datatypes.JSON, x int, y int, warehouseID uint, status PackageStatus) *Package {
	return &Package{
		ID:          packageID,
		Username:    username,
		TruckID:     nil, // or pointer to uint if assigned
		Items:       items,
		Destination: Coordinate{X: x, Y: y},
		WarehouseID: warehouseID,
		Status:      status, // use constant if defined
	}
}
