package controller

import (
	"fmt"
	"log"
	"mini-ups/model"
	"mini-ups/protocol/worldupspb"
	"mini-ups/service"
	"mini-ups/util"
)

// keeps listening for world responses
// called in main as goroutine to keep listening in the background.
func ParseWorldResponse() {
	for {
		resp := &worldupspb.UResponses{}
		if err := util.RecvMsg(util.UPSConn, resp); err != nil {
			log.Fatal("Error receiving world response:", err)
			continue
		}

		HandleCompletions(resp.GetCompletions())
		HandleDeliveries(resp.GetDelivered())
		HandleAcks(resp.GetAcks()) //For each command we give, World returns an ack corresponding to the sequence number
		// E.g. send two trucks to go pick up (seqnum 1, 2) and delivery (seqnum 3), the response returns 1, 2, 3
		HandleTruckStatus(resp.GetTruckstatus())
		HandleErrors(resp.GetError())

		if finished := resp.GetFinished(); finished {
			log.Println("World server disconnected")
			// 這裡要做什麼? 直接結束program嗎?
		}
	}
}

//Assum all response are correct

// TODO
// Completion relates to the status of the truck
// Completion: arrive_warehouse or idle(completed delivery)
func HandleCompletions(completions []*worldupspb.UFinished) {
	for _, completion := range completions { //可能有很多個
		truckID := completion.GetTruckid()
		// x := completion.GetX()
		// y := completion.GetY()
		status := completion.GetStatus()
		//seqnum := completion.GetSeqnum()

		//query request based on seqnum and get the packageID
		//packageID :=

		//update package and truck status
		if status == "ARRIVE WAREHOUSE" {
			//query request based on seqnum and get the warehouseID
			//warehouseID :=

			service.ChangeTruckStatus(int(truckID), model.TruckStatus.ARRIVED)
			//service.ChangePackageStatus(packageID, model.StatusPickupComplete)
			//service.NotifyAmazonTruckArrived(truckID, warehouseID)
		} else if status == "IDLE" {
			service.ChangeTruckStatus(int(truckID), model.TruckStatus.IDLE)
			//service.NotifyAmazonDeliveryComplete(packageID, truckID, x, y)
		}
		//UPDATE TRUCK COORDINATE
	}
}

// TODO
func HandleDeliveries(delivered []*worldupspb.UDeliveryMade) {
	for _, delivery := range delivered {
		//truckID := int(delivery.GetTruckid())
		packageID := fmt.Sprintf("%d", delivery.GetPackageid())
		//seqnum := delivery.GetSeqnum()
		service.ChangePackageStatus(packageID, model.StatusDelivered)
		//I choose to NotifyAmazonDeliveryComplete in HandleCompletion instead of here because it has x,y
	}
}

// TODO
// Update our seqnum + response datastructure
func HandleAcks(acks []int64) {
	// for _, ack := range acks {

	// }
}

// TODO
func HandleTruckStatus(statusList []*worldupspb.UTruck) {
	for _, truck := range statusList {
		truckID := truck.GetTruckid()
		//status := truck.GetStatus()
		x := truck.GetX()
		y := truck.GetY()
		//seqnum := truck.GetSeqnum()
		//update truck info here
		truck, err := service.GetTruckByID(model.TruckID(truckID))
		if err != nil {

		}
		// truck.status = status  <-- naming from world is different from our ENUM,
		truck.Coord.X = int(x)
		truck.Coord.Y = int(y)

	}
}

// TODO
func HandleErrors(errors []*worldupspb.UErr) {
}
