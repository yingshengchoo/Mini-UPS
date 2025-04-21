package protocol

import (
	"encoding/binary"
	"fmt"
	"net"

	"mini-ups/protocol/worldamazonpb"

	"google.golang.org/protobuf/proto"
)

// Serialize and send protobuf message
func sendMsg(conn net.Conn, msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	// Varint encode the length
	var lenBuf [binary.MaxVarintLen32]byte
	n := binary.PutUvarint(lenBuf[:], uint64(len(data)))

	// Send length followed by data
	if _, err := conn.Write(lenBuf[:n]); err != nil {
		return err
	}
	if _, err := conn.Write(data); err != nil {
		return err
	}

	fmt.Printf("Sent %d bytes header + %d bytes data\n", n, len(data))
	return nil
}

func recvMsg(conn net.Conn, msg proto.Message) error {
	// Read varint length prefix
	var size uint64
	var err error
	var sizeBuf [1]byte
	var shift uint

	for {
		_, err = conn.Read(sizeBuf[:])
		if err != nil {
			return err
		}
		b := sizeBuf[0]
		size |= uint64(b&0x7F) << shift
		if b < 0x80 {
			break
		}
		shift += 7
		if shift >= 64 {
			return fmt.Errorf("recvMsg: varint size too long")
		}
	}

	// Read full message of decoded size
	data := make([]byte, size)
	total := 0
	for total < int(size) {
		n, err := conn.Read(data[total:])
		if err != nil {
			return err
		}
		total += n
	}

	return proto.Unmarshal(data, msg)
}

// createWarehouse constructs an AInitWarehouse given ID, x, and y coordinates
func createWarehouse(id, x, y int32) *worldamazonpb.AInitWarehouse {
	return &worldamazonpb.AInitWarehouse{
		Id: proto.Int32(id),
		X:  proto.Int32(x),
		Y:  proto.Int32(y),
	}
}

func connectAmazon(warehouses []*worldamazonpb.AInitWarehouse) (net.Conn, int64) {
	conn, err := net.Dial("tcp", "vcm-47478.vm.duke.edu:23456")
	if err != nil {
		panic(err)
	}
	fmt.Println("Sending AConnect...")

	aconnect := &worldamazonpb.AConnect{
		IsAmazon: proto.Bool(true),
		Initwh:   warehouses,
	}
	fmt.Printf("AConnect before marshal: %+v\n", aconnect)
	raw, _ := proto.Marshal(aconnect)
	fmt.Println("Serialized AConnect length:", len(raw))

	if err := sendMsg(conn, aconnect); err != nil {
		panic(err)
	}

	resp := &worldamazonpb.AConnected{}
	if err := recvMsg(conn, resp); err != nil {
		panic(err)
	}
	fmt.Println("AConnected:", resp)
	return conn, resp.GetWorldid()
}

func createProduct(id int64, description string, count int32) *worldamazonpb.AProduct {
	return &worldamazonpb.AProduct{
		Id:          proto.Int64(id),
		Description: proto.String(description),
		Count:       proto.Int32(count),
	}
}

func createPurchaseMore(whnum int32, seqnum int64, products []*worldamazonpb.AProduct) *worldamazonpb.APurchaseMore {
	return &worldamazonpb.APurchaseMore{
		Whnum:  proto.Int32(whnum),
		Seqnum: proto.Int64(seqnum),
		Things: products,
	}
}

func createPack(whnum int32, shipid int64, seqnum int64, products []*worldamazonpb.AProduct) *worldamazonpb.APack {
	return &worldamazonpb.APack{
		Whnum:  proto.Int32(whnum),
		Shipid: proto.Int64(shipid),
		Seqnum: proto.Int64(seqnum),
		Things: products,
	}
}

func createPutOnTruckCommand(whnum, truckid int32, shipid, seqnum int64) *worldamazonpb.APutOnTruck {
	return &worldamazonpb.APutOnTruck{
		Whnum:   proto.Int32(whnum),
		Truckid: proto.Int32(truckid),
		Shipid:  proto.Int64(shipid),
		Seqnum:  proto.Int64(seqnum),
	}
}

func createAmazonDisconnectCommand() *worldamazonpb.ACommands {
	return &worldamazonpb.ACommands{
		Disconnect: proto.Bool(true),
	}
}

func createAmazonCommands(
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

func sendAmazonCommands(conn net.Conn, cmd *worldamazonpb.ACommands) error {
	if err := sendMsg(conn, cmd); err != nil {
		return fmt.Errorf("failed to send ACommands to Amazon: %w", err)
	}
	fmt.Println("Sent ACommands to Amazon:", cmd)
	return nil
}
