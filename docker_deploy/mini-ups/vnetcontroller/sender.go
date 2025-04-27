package vnetcontroller

import (
	"encoding/binary"
	"log"
	"mini-ups/dao"
	"mini-ups/model"
	"mini-ups/protocol"
	"mini-ups/protocol/worldupspb"
	"mini-ups/service"
	"mini-ups/util"
	"net"

	"google.golang.org/protobuf/proto"
)

//The sender object handles sending message to the world simulation.

type Sender struct {
	recvWindow *RecvWindow
	sendWindow *SendWindow
	conn       net.Conn
}

func NewSender(rw *RecvWindow, sw *SendWindow, conn net.Conn) *Sender {
	return &Sender{
		recvWindow: rw,
		sendWindow: sw,
		conn:       conn,
	}
}

// Send sends the command object to the world simulation.
func (s *Sender) Send(msg proto.Message) error {
	err := s.SendMsg(s.conn, msg)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Sends a command to the world simulation to havea truck go pick up a package at a warehouse
func (s *Sender) SendWorldRequestToGoPickUp(truckID model.TruckID, warehouseID uint) error {
	seqnum := util.GenerateSeqNum()

	goPickup := protocol.CreateGoPickupCommand(int32(truckID), int32(warehouseID), seqnum)
	cmd := protocol.CreateUPSCommands([]*worldupspb.UGoPickup{goPickup}, nil, 0, false, nil, nil)

	if err := service.ChangeTruckStatus(int(truckID), model.TruckStatus.PICKING); err != nil {
		log.Println(err)
		return err
	}
	// fmt.Printf("Sending UGoPickup command: TruckID=%d, WarehouseID=%d, Seqnum=%d\n", truckID, warehouseID, seqnum) //debugging

	s.addMsg(seqnum, "UGoPickup", goPickup)
	return s.Send(cmd)
}

// Sends a UGoDelivery command to world that tells the truck to delivery the package
func (s *Sender) SendWorldDeliveryRequest(packageID string) error {
	seqnum := util.GenerateSeqNum()
	pack, err := dao.GetPackagesByPackageID(packageID)
	if err != nil {
		return err
	}
	delivery := protocol.CreateDeliveryLocation(int64(*pack.TruckID), int32(pack.Destination.X), int32(pack.Destination.Y))
	goDeliver := protocol.MakeDelivery(int32(*pack.TruckID), seqnum, []*worldupspb.UDeliveryLocation{delivery})
	cmd := protocol.CreateUPSCommands(nil, []*worldupspb.UGoDeliver{goDeliver}, 0, false, nil, nil)

	s.addMsg(seqnum, "UGoDeliver", goDeliver)

	return s.Send(cmd)
}

// Sends a Truck Query to World of truckID
func (s *Sender) SendWorldTruckQuery(truckID model.TruckID) error {

	seqnum := util.GenerateSeqNum()
	query := protocol.MakeTruckQuery(int32(truckID), seqnum)
	cmd := protocol.CreateUPSCommands(nil, nil, 0, false, []*worldupspb.UQuery{query}, nil)
	s.addMsg(seqnum, "UQuery", query)
	return s.Send(cmd)
}

// TODO
func (s *Sender) SendWorldAck(ack int64) error {
	cmd := protocol.CreateUPSCommands(nil, nil, 0, false, nil, []int64{ack})
	//s.recvWindow.Record(ack) <-- NOT HERE!! record Ack in reciever -> !
	//														 This is because we call Record there to determine if we need to do the operation at all.
	return s.Send(cmd)
}

func (s *Sender) addMsg(seqnum int64, msgType string, msg interface{}) {
	s.sendWindow.Add(seqnum, msgType, msg)
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

	log.Printf("Sent %d bytes header + %d bytes data\n", n, len(data))
	return nil
}
