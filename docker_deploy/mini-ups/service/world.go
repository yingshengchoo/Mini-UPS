package service

import (
	"fmt"
	"mini-ups/dao"
	"mini-ups/model"
	"mini-ups/protocol"
	"mini-ups/protocol/worldupspb"

	"google.golang.org/protobuf/proto"
)

//This service file handles oeprations related to sending commands to the world simulation.

// SendRequestToGoPickUp sends a pickup request to the UPS world process
func SendWorldRequestToGoPickUp(truckID model.TruckID, warehouseID uint, seqnum int64) error {
	cmd := &worldupspb.UCommands{
		Pickups: []*worldupspb.UGoPickup{{
			Truckid: proto.Int32(int32(truckID)),
			Whid:    proto.Int32(int32(warehouseID)),
			Seqnum:  proto.Int64(seqnum),
		}},
	}

	fmt.Printf("Sending UGoPickup command: TruckID=%d, WarehouseID=%d, Seqnum=%d\n", truckID, warehouseID, seqnum)

	if err := protocol.SendUPSCommands(cmd); err != nil {
		return fmt.Errorf("failed to send pickup command: %w", err)
	}

	return nil
}

// Sends a UGoDelivery command to world that tells the truck to delivery the package
func SendWorldDeliveryRequest(packageID string, seqnum int64) error {

	pack, err := dao.GetPackagesByPackageID(packageID)
	if err != nil {
		return err
	}
	delivery := protocol.CreateDeliveryLocation(int64(*pack.TruckID), int32(pack.Destination.X), int32(pack.Destination.Y))
	goDeliver := protocol.MakeDelivery(int32(*pack.TruckID), seqnum, []*worldupspb.UDeliveryLocation{delivery})
	cmd := protocol.CreateUPSCommands(nil, []*worldupspb.UGoDeliver{goDeliver}, 0, false, nil, nil)

	if err := protocol.SendUPSCommands(cmd); err != nil {
		return fmt.Errorf("error sending delivery request: %w", err)
	}

	return nil
}

// Sends a Truck Query to World of truckID
func SendWorldTruckQuery(truckID model.TruckID, seqnum int64) error {
	query := protocol.MakeTruckQuery(int32(truckID), seqnum)
	cmd := protocol.CreateUPSCommands(nil, nil, 0, false, []*worldupspb.UQuery{query}, nil)

	if err := protocol.SendUPSCommands(cmd); err != nil {
		return fmt.Errorf("error sending delivery request: %w", err)
	}

	return nil
}
