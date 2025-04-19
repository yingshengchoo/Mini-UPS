package dao

import (
	"mini-ups/db"
	"mini-ups/model"
)

func CreateTruck(truck *model.Truck) error {
	return db.DB.Create(truck).Error
}
