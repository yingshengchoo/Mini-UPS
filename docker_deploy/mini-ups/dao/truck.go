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
