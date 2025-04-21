package protocol

import (
	"fmt"
	"net"

	"mini-ups/protocol/worldupspb"

	"google.golang.org/protobuf/proto"
)

// createTruck constructs a UInitTruck given ID, x, and y coordinates
func createTruck(id, x, y int32) *worldupspb.UInitTruck {
	return &worldupspb.UInitTruck{
		Id: proto.Int32(id),
		X:  proto.Int32(x),
		Y:  proto.Int32(y),
	}
}

// connectUPS connects to UPS with a given worldID and a list of initial trucks
// This assumes that Amazon connects first. If we connect first, may consider not providing worldID
func connectUPS(worldID int64, trucks []*worldupspb.UInitTruck) net.Conn {
	conn, err := net.Dial("tcp", "vcm-47478.vm.duke.edu:12345")
	if err != nil {
		panic(err)
	}
	fmt.Println("Sending UConnect...")

	uconnect := &worldupspb.UConnect{
		IsAmazon: proto.Bool(false),
		Worldid:  proto.Int64(worldID),
		Trucks:   trucks,
	}
	if err := sendMsg(conn, uconnect); err != nil {
		panic(err)
	}

	resp := &worldupspb.UConnected{}
	if err := recvMsg(conn, resp); err != nil {
		panic(err)
	}
	fmt.Println("UConnected:", resp)
	return conn
}

func createGoPickupCommand(truckID, whID int32, seqnum int64) *worldupspb.UGoPickup {
	return &worldupspb.UGoPickup{
		Truckid: proto.Int32(truckID),
		Whid:    proto.Int32(whID),
		Seqnum:  proto.Int64(seqnum),
	}
}

func createDeliveryLocation(packageID int64, x, y int32) *worldupspb.UDeliveryLocation {
	return &worldupspb.UDeliveryLocation{
		Packageid: proto.Int64(packageID),
		X:         proto.Int32(x),
		Y:         proto.Int32(y),
	}
}

func makeDelivery(truckID int32, seqnum int64, deliveries []*worldupspb.UDeliveryLocation) *worldupspb.UGoDeliver {
	return &worldupspb.UGoDeliver{
		Truckid:  proto.Int32(truckID),
		Seqnum:   proto.Int64(seqnum),
		Packages: deliveries,
	}
}

func createUPSCommands(
	pickups []*worldupspb.UGoPickup,
	deliveries []*worldupspb.UGoDeliver,
	simspeed uint32,
	disconnect bool,
	queries []*worldupspb.UQuery,
	acks []int64,
) *worldupspb.UCommands {
	return &worldupspb.UCommands{
		Pickups:    pickups,
		Deliveries: deliveries,
		Simspeed:   proto.Uint32(simspeed),
		Disconnect: proto.Bool(disconnect),
		Queries:    queries,
		Acks:       acks,
	}
}

func sendUPSCommands(conn net.Conn, cmd *worldupspb.UCommands) error {
	if err := sendMsg(conn, cmd); err != nil {
		return fmt.Errorf("failed to send UCommands to UPS: %w", err)
	}
	fmt.Println("Sent UCommands to UPS:", cmd)
	return nil
}
