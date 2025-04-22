package service

import (
	"fmt"
	"mini-ups/dao"
	"mini-ups/model"
	"mini-ups/util"
)

// register a truck
func RegisterTruck(truckID model.TruckID, x int, y int) error {
	return dao.CreateTruck(model.NewTruck(
		truckID,
		x, y,
		model.TruckStatus.IDLE,
	))
}

// get truck by id
func GetTruckByID(truckID model.TruckID) (*model.Truck, error) {
	return dao.GetTruckByID(truckID)
}

func GetIDByTruck(truck *model.Truck) (model.TruckID, error) {
	if truck == nil {
		return -1, fmt.Errorf("truck is nil")
	}
	return truck.ID, nil
}

// communicate with world
func GetUpdatedTruckInfo(truckID model.TruckID) error {
	seqnum := util.GenerateSeqNum()
	return SendWorldTruckQuery(truckID, seqnum)
}

func GetFirstIdleTruck() (*model.Truck, error) {
	return dao.GetFirstIdleTruck()
}

func ChangeTruckStatus(truckID int, newStatus model.Status) error {
	// Add business rules here if needed
	return dao.UpdateTruckStatus(truckID, newStatus)
}
