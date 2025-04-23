package vnetcontroller

import (
	"encoding/binary"
	"fmt"
	"mini-ups/dao"
	"mini-ups/model"
	"mini-ups/protocol"
	"mini-ups/protocol/worldupspb"
	"mini-ups/util"
	"net"

	"google.golang.org/protobuf/proto"
)

type Sender struct {
	conn       net.Conn
	sendWindow *SendWindow
}

func NewSender(conn net.Conn, sw *SendWindow) *Sender {
	return &Sender{
		conn:       conn,
		sendWindow: sw,
	}
}

func (s *Sender) Send(msg proto.Message, seqnum int64) error {
	err := util.SendMsg(s.conn, msg)
	if err != nil {
		return err
	}
	s.sendWindow.Add(seqnum, msg)
	return nil
}

func (s *Sender) SendWorldRequestToGoPickUp(truckID model.TruckID, warehouseID uint, seqnum int64) error {
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
func (s *Sender) SendWorldDeliveryRequest(packageID string, seqnum int64) error {

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
func (s *Sender) SendWorldTruckQuery(truckID model.TruckID, seqnum int64) error {
	query := protocol.MakeTruckQuery(int32(truckID), seqnum)
	cmd := protocol.CreateUPSCommands(nil, nil, 0, false, []*worldupspb.UQuery{query}, nil)

	if err := protocol.SendUPSCommands(cmd); err != nil {
		return fmt.Errorf("error sending delivery request: %w", err)
	}

	return nil
}

// reserved for future use
func (s *Sender) SendMsg(conn net.Conn, msg proto.Message) error {
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
