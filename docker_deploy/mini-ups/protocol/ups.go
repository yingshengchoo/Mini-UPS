package protocol

import (
	"fmt"
	"log"
	"net"
	"time"

	"mini-ups/protocol/worldupspb"
	"mini-ups/util"

	"google.golang.org/protobuf/proto"
)

//These functions are help functions that constructs commands to communicate with the world simulator.

// createTruck constructs a UInitTruck given ID, x, and y coordinates
func CreateTruck(id, x, y int32) *worldupspb.UInitTruck {
	return &worldupspb.UInitTruck{
		Id: proto.Int32(id),
		X:  proto.Int32(x),
		Y:  proto.Int32(y),
	}
}

// connectUPS connects to UPS with a given worldID and a list of initial trucks
// This assumes that Amazon connects first. If we connect first, may consider not providing worldID
func ConnectUPSWithWorldID(worldID int64, trucks []*worldupspb.UInitTruck) net.Conn {
	conn, err := net.Dial("tcp", util.WORLD_HOST+":12345")
	if err != nil {
		panic(err)
	}
	fmt.Println("Sending UConnect...")

	uconnect := &worldupspb.UConnect{
		IsAmazon: proto.Bool(false),
		// Worldid:  proto.Int64(worldID),
		Trucks: trucks,
	}
	if err := util.SendMsg(conn, uconnect); err != nil {
		panic(err)
	}

	resp := &worldupspb.UConnected{}
	if err := util.RecvMsg(conn, resp); err != nil {
		panic(err)
	}
	fmt.Println("UConnected:", resp)
	return conn
}

func ConnectUPS(trucks []*worldupspb.UInitTruck) (net.Conn, int64) {

	ticker := time.NewTicker(1 * time.Second) // try connect once per second
	defer ticker.Stop()                       // stop ticker
	var conn net.Conn
	var err error

	// try to conenct every second until success
	for range ticker.C { // wait for ticker
		conn, err = net.Dial("tcp", util.WORLD_HOST+":12345")
		if err != nil {
			log.Println(err)
		} else {
			break
		}
	}
	fmt.Println("Sending UConnect...")

	aconnect := &worldupspb.UConnect{
		IsAmazon: proto.Bool(false),
		Trucks:   trucks,
	}

	if err := util.SendMsg(conn, aconnect); err != nil {
		panic(err)
	}

	resp := &worldupspb.UConnected{}
	if err := util.RecvMsg(conn, resp); err != nil {
		panic(err)
	}
	fmt.Println("UConnected:", resp)
	return conn, resp.GetWorldid()
}

func CreateGoPickupCommand(truckID, whID int32, seqnum int64) *worldupspb.UGoPickup {
	return &worldupspb.UGoPickup{
		Truckid: proto.Int32(truckID),
		Whid:    proto.Int32(whID),
		Seqnum:  proto.Int64(seqnum),
	}
}

func CreateDeliveryLocation(packageID int64, x, y int32) *worldupspb.UDeliveryLocation {
	return &worldupspb.UDeliveryLocation{
		Packageid: proto.Int64(packageID),
		X:         proto.Int32(x),
		Y:         proto.Int32(y),
	}
}

func MakeDelivery(truckID int32, seqnum int64, deliveries []*worldupspb.UDeliveryLocation) *worldupspb.UGoDeliver {
	return &worldupspb.UGoDeliver{
		Truckid:  proto.Int32(truckID),
		Seqnum:   proto.Int64(seqnum),
		Packages: deliveries, //好像老師 覺得可以一次做很多個delivery?
	}
}

func MakeTruckQuery(truckID int32, seqnum int64) *worldupspb.UQuery {
	return &worldupspb.UQuery{
		Truckid: proto.Int32(truckID),
		Seqnum:  proto.Int64(seqnum),
	}
}

func CreateUPSCommands(
	pickups []*worldupspb.UGoPickup, // Make UGoPickUP
	deliveries []*worldupspb.UGoDeliver,
	simspeed uint32, // Can include to make it go faster i think.
	disconnect bool, //Always use False unless closing connection
	queries []*worldupspb.UQuery,
	acks []int64,
) *worldupspb.UCommands {
	return &worldupspb.UCommands{
		Pickups:    pickups,
		Deliveries: deliveries,
		Simspeed:   proto.Uint32(100),
		Disconnect: proto.Bool(disconnect),
		Queries:    queries,
		Acks:       acks,
	}
}

func SendUPSCommands(cmd *worldupspb.UCommands) error {
	if err := util.SendMsg(util.UPSConn, cmd); err != nil {
		return fmt.Errorf("failed to send UCommands to UPS: %w", err)
	}
	fmt.Println("Sent UCommands to UPS:", cmd)
	return nil
}
