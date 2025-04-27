package vnetcontroller

import (
	"fmt"
	"io"
	"log"
	"mini-ups/dao"
	"mini-ups/model"
	"mini-ups/protocol/worldupspb"
	"mini-ups/service"

	"net"

	"google.golang.org/protobuf/proto"
)

type Receiver struct {
	recvWindow *RecvWindow
	sendWindow *SendWindow
	sender     *Sender
}

func NewReceiver(rw *RecvWindow, sw *SendWindow, s *Sender) *Receiver {
	return &Receiver{recvWindow: rw, sendWindow: sw, sender: s}
}

func (r *Receiver) ListenForWorldResponses(conn net.Conn) {
	for {
		resp := &worldupspb.UResponses{}
		if err := r.RecvMsg(conn, resp); err != nil {
			log.Print(resp)
			fmt.Println("Error receiving world response:", err)
			if err == io.EOF {
				// Gracefully handle EOF, e.g., connection closed
				return
			}
			continue
		}
		log.Print(resp)
		r.handleCompletions(resp.GetCompletions())
		r.handleDeliveries(resp.GetDelivered())
		r.handleAcks(resp.GetAcks())
		r.handleTruckStatus(resp.GetTruckstatus())
		r.handleErrors(resp.GetError())

		if resp.GetFinished() {
			fmt.Println("World server disconnected")
			return
		}
	}
}

func (r *Receiver) handleAcks(acks []int64) {
	for _, ack := range acks {
		r.sendWindow.Ack(ack)
	}
}

//Assume all response are correct

// TODO
// Completion relates to the status of the truck
// Completion: arrive_warehouse or idle(completed delivery)
func (r *Receiver) handleCompletions(completions []*worldupspb.UFinished) {
	for _, completion := range completions { //可能有很多個
		truckID := completion.GetTruckid()
		x := completion.GetX()
		y := completion.GetY()
		status := completion.GetStatus()
		seqnum := completion.GetSeqnum()

		//If Seqnum has already been processed skip to next one
		if !r.handleSeqnum(seqnum) {
			continue
		}

		pack, err := dao.GetPickingPackageInfoByTruckID(truckID)
		packageID := pack.ID
		if err != nil {
			log.Println("Error:", err)
			continue
		}

		log.Println("do this once:", seqnum)
		//update package and truck status
		if status == "ARRIVE WAREHOUSE" {
			//query request based on seqnum and get the warehouseID
			warehouseID := pack.WarehouseID

			service.ChangeTruckStatus(int(truckID), model.TruckStatus.ARRIVED)
			service.ChangePackageStatus(string(packageID), model.StatusPickupComplete)
			service.NotifyAmazonTruckArrived(int(truckID), int(warehouseID))
		} else if status == "IDLE" {
			service.ChangeTruckStatus(int(truckID), model.TruckStatus.IDLE)
			service.NotifyAmazonDeliveryComplete(string(packageID), int(truckID), int(x), int(y))
			PkQueue.TryConsumeWithLock()
		}
		//udpates truck coordinates
		err = service.ChangeTruckCoord(model.TruckID(int(truckID)), int(x), int(y))
		if err != nil {
			log.Println("Error:", err)
			continue
		}

	}
}

// TODO
func (r *Receiver) handleDeliveries(delivered []*worldupspb.UDeliveryMade) {
	for _, delivery := range delivered {
		//truckID := int(delivery.GetTruckid())
		packageID := fmt.Sprintf("%d", delivery.GetPackageid())
		seqnum := delivery.GetSeqnum()

		// If Seqnum has already been processed return.
		if !r.handleSeqnum(seqnum) {
			continue
		}

		service.ChangePackageStatus(packageID, model.StatusDelivered)
		//I choose to NotifyAmazonDeliveryComplete in HandleCompletion instead of here because it has x,y

	}
}

// TODO
func (r *Receiver) handleTruckStatus(statusList []*worldupspb.UTruck) {
	for _, truck := range statusList {
		truckID := truck.GetTruckid()
		//status := truck.GetStatus()
		x := truck.GetX()
		y := truck.GetY()
		seqnum := truck.GetSeqnum()

		//If Seqnum has already been processed return.
		if !r.handleSeqnum(seqnum) {
			continue
		}

		//update truck info here
		truck, err := service.GetTruckByID(model.TruckID(truckID))
		if err != nil {
			log.Println("Error:", err)
			continue
		}
		// truck.status = status  <-- naming from world is different from our ENUM,  <--- HERE!!! TODO
		truck.Coord.X = int(x)
		truck.Coord.Y = int(y)
	}
}

// TODO
func (r *Receiver) handleErrors(errors []*worldupspb.UErr) {
}

// returns true if the seqnum has not been handled, False if it has already been handled
func (r *Receiver) handleSeqnum(seqnum int64) bool {
	//Check if seqnum has been ack or not.
	//SendAck regardless (for each seqnum from World we send it back in case it drops the pacakge)

	err := r.sender.SendWorldAck(seqnum)
	if err != nil {
		log.Println("Error:", err)
	}

	if r.recvWindow.Record(seqnum) {
		return true
	} else {
		return false
	}
}

func (r *Receiver) RecvMsg(conn net.Conn, msg proto.Message) error {
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
