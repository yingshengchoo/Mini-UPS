package dao

import (
	"mini-ups/db"
	"mini-ups/model"
)

// create
func CreateTruck(truck *model.Truck) error {
	return db.DB.Create(truck).Error
}

// update
func UpdateTruck(truck *model.Truck) error {
	return db.DB.Save(truck).Error
}

// get truck by id
func GetTruckByID(truckID model.TruckID) (*model.Truck, error) {
	var truck model.Truck
	if err := db.DB.Where("ID = ?", truckID).First(&truck).Error; err != nil {
		return nil, err
	}
	return &truck, nil
}

func GetFirstIdleTruck() (*model.Truck, error) {
	var truck model.Truck
	if err := db.DB.Where("status = ?", model.TruckStatus.IDLE).First(&truck).Error; err != nil {
		return nil, err
	}
	return &truck, nil
}

func UpdateTruckStatus(truckID int, newStatus model.Status) error {
	var truck model.Truck
	if err := db.DB.First(&truck, truckID).Error; err != nil {
		return err
	}

	truck.Transit(newStatus)

	return db.DB.Save(&truck).Error
}
