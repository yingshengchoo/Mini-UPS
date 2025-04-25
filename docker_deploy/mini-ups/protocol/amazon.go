package protocol

import (
	"fmt"
	"net"

	"mini-ups/protocol/worldamazonpb"
	"mini-ups/util"

	"google.golang.org/protobuf/proto"
)

// createWarehouse constructs an AInitWarehouse given ID, x, and y coordinates
func CreateWarehouse(id, x, y int32) *worldamazonpb.AInitWarehouse {
	return &worldamazonpb.AInitWarehouse{
		Id: proto.Int32(id),
		X:  proto.Int32(x),
		Y:  proto.Int32(y),
	}
}

func ConnectAmazon(warehouses []*worldamazonpb.AInitWarehouse) (net.Conn, int64) {
	conn, err := net.Dial("tcp", "vcm-47478.vm.duke.edu:23456")
	if err != nil {
		panic(err)
	}
	fmt.Println("Sending AConnect...")

	aconnect := &worldamazonpb.AConnect{
		IsAmazon: proto.Bool(true),
		Initwh:   warehouses,
	}

	if err := util.SendMsg(conn, aconnect); err != nil {
		panic(err)
	}

	resp := &worldamazonpb.AConnected{}
	if err := util.RecvMsg(conn, resp); err != nil {
		panic(err)
	}
	fmt.Println("AConnected:", resp)
	return conn, resp.GetWorldid()
}

func ConnectAmazonWithWorldID(worldID int64, warehouses []*worldamazonpb.AInitWarehouse) net.Conn {
	conn, err := net.Dial("tcp", util.WORLD_HOST)
	if err != nil {
		panic(err)
	}
	fmt.Println("Sending AConnect...")

	uconnect := &worldamazonpb.AConnect{
		IsAmazon: proto.Bool(false),
		Worldid:  proto.Int64(worldID),
		Initwh:   warehouses,
	}
	if err := util.SendMsg(conn, uconnect); err != nil {
		panic(err)
	}

	resp := &worldamazonpb.AConnected{}
	if err := util.RecvMsg(conn, resp); err != nil {
		panic(err)
	}
	fmt.Println("AConnected:", resp)
	return conn
}

func CreateProduct(id int64, description string, count int32) *worldamazonpb.AProduct {
	return &worldamazonpb.AProduct{
		Id:          proto.Int64(id),
		Description: proto.String(description),
		Count:       proto.Int32(count),
	}
}

func CreatePurchaseMore(whnum int32, seqnum int64, products []*worldamazonpb.AProduct) *worldamazonpb.APurchaseMore {
	return &worldamazonpb.APurchaseMore{
		Whnum:  proto.Int32(whnum),
		Seqnum: proto.Int64(seqnum),
		Things: products,
	}
}

func CreatePack(whnum int32, shipid int64, seqnum int64, products []*worldamazonpb.AProduct) *worldamazonpb.APack {
	return &worldamazonpb.APack{
		Whnum:  proto.Int32(whnum),
		Shipid: proto.Int64(shipid),
		Seqnum: proto.Int64(seqnum),
		Things: products,
	}
}

func CreatePutOnTruckCommand(whnum, truckid int32, shipid, seqnum int64) *worldamazonpb.APutOnTruck {
	return &worldamazonpb.APutOnTruck{
		Whnum:   proto.Int32(whnum),
		Truckid: proto.Int32(truckid),
		Shipid:  proto.Int64(shipid),
		Seqnum:  proto.Int64(seqnum),
	}
}

func CreateAmazonDisconnectCommand() *worldamazonpb.ACommands {
	return &worldamazonpb.ACommands{
		Disconnect: proto.Bool(true),
	}
}

func CreateAmazonCommands(
	purchases []*worldamazonpb.APurchaseMore,
	packs []*worldamazonpb.APack,
	loads []*worldamazonpb.APutOnTruck,
	queries []*worldamazonpb.AQuery,
	acks []int64,
	simspeed uint32,
	disconnect bool,
) *worldamazonpb.ACommands {
	return &worldamazonpb.ACommands{
		Buy:        purchases,
		Topack:     packs,
		Load:       loads,
		Queries:    queries,
		Simspeed:   proto.Uint32(simspeed),
		Disconnect: proto.Bool(disconnect),
		Acks:       acks,
	}
}

func SendAmazonCommands(conn net.Conn, cmd *worldamazonpb.ACommands) error {
	if err := util.SendMsg(conn, cmd); err != nil {
		return fmt.Errorf("failed to send ACommands to Amazon: %w", err)
	}
	fmt.Println("Sent ACommands to Amazon:", cmd)
	return nil
}
