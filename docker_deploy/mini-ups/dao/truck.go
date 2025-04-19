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
