package model

import (
	"time"

	"gorm.io/datatypes"
)

type Package struct {
	ID          string         `gorm:"primaryKey" json:"package_id"`
	UserID      uint           `gorm:"not null" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID"`
	TruckID     *uint          `json:"truck_id"` // nullable
	Truck       *Truck         `gorm:"foreignKey:TruckID"`
	Items       datatypes.JSON `json:"items"`
	WarehouseID int            `gorm:"not null" json:"warehouse_id"`
	Status      TruckStatus    `gorm:"type:varchar(20)" json:"status"` // Same enum type as Truck
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
