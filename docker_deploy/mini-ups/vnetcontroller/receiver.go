package vnetcontroller

import (
	"fmt"
	"log"
	"mini-ups/model"
	"mini-ups/protocol/worldupspb"
	"mini-ups/service"
	"net"

	"google.golang.org/protobuf/proto"
)

type Receiver struct {
	recvWindow *RecvWindow
	sendWindow *SendWindow
}

func NewReceiver(rw *RecvWindow, sw *SendWindow) *Receiver {
	return &Receiver{recvWindow: rw, sendWindow: sw}
}

func (r *Receiver) ListenForWorldResponses(conn net.Conn) {
	for {
		resp := &worldupspb.UResponses{}
		if err := r.RecvMsg(conn, resp); err != nil {
			fmt.Println("Error receiving world response:", err)
			continue
		}
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

//Assum all response are correct

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

		//query request based on seqnum and get the packageID
		//resp := r.sendWindow.GetResponse(seqnum) //   <-------- IF we use Goroutine in world handler, watch out here.
		//             What if we try to get response, but ackHandler deletes the response.

		//Also revist logic here. How do we get the package ID associated with
		pack, err := service.GetPackageInfoByTruck(truckID)
		packageID := pack.ID
		if err != nil {
			log.Println("Error:", err)
			return
		}

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
		}
		//udpates truck coordinates
		err = service.ChangeTruckCoord(model.TruckID(int(truckID)), int(x), int(y))
		if err != nil {
			log.Println("Error:", err)
			return
		}

		r.handleSeqnum(seqnum)
	}
}

// TODO
func (r *Receiver) handleDeliveries(delivered []*worldupspb.UDeliveryMade) {
	for _, delivery := range delivered {
		//truckID := int(delivery.GetTruckid())
		packageID := fmt.Sprintf("%d", delivery.GetPackageid())
		seqnum := delivery.GetSeqnum()
		service.ChangePackageStatus(packageID, model.StatusDelivered)
		//I choose to NotifyAmazonDeliveryComplete in HandleCompletion instead of here because it has x,y

		r.handleSeqnum(seqnum)
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
		//update truck info here
		truck, err := service.GetTruckByID(model.TruckID(truckID))
		if err != nil {
			log.Println("Error:", err)
			return
		}
		// truck.status = status  <-- naming from world is different from our ENUM,
		truck.Coord.X = int(x)
		truck.Coord.Y = int(y)

		r.handleSeqnum(seqnum)
	}
}

// TODO
func (r *Receiver) handleErrors(errors []*worldupspb.UErr) {
}

func (r *Receiver) handleSeqnum(seqnum int64) {
	//Check if seqnum has been ack or not. If no, send ack to world.
	if r.recvWindow.Record(seqnum) {
		//sendAckToWorld(seqnum) <-- hmm not sure what to do here.. This logic belongs in sender.. but how do we access it?
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
