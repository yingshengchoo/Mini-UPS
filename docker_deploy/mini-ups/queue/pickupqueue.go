package queue

import (
	"log"
	"mini-ups/model"
	"mini-ups/service"
	"mini-ups/vnetcontroller"

	"gorm.io/datatypes"
)

type PickupQueue struct {
	queue chan *PickupReq
	// availableTrucks int
}

// pick up request struct
// seperate it because it will be used somewhere else
type PickupReq struct {
	Action      string         `json:"action"`
	PackageID   string         `json:"package_id" binding:"required"`
	Username    string         `json:"username" binding:"required"`
	Items       datatypes.JSON `json:"items" binding:"required"`
	DestX       int            `json:"destination_x" binding:"required"`
	DestY       int            `json:"destination_y" binding:"required"`
	WarehouseID uint           `json:"warehouse_id" binding:"required"`
	MessageID   string         `json:"message_id" binding:"required"`
}

var VnetCtrl *vnetcontroller.Controller
var PkQueue = NewPickupQueue()

func NewPickupQueue() *PickupQueue {
	return &PickupQueue{
		queue: make(chan *PickupReq),
		// availableTrucks: truckNum,
	}
}

// add new req
func (q *PickupQueue) AddRequest(req *PickupReq) {
	log.Println("here")
	q.queue <- req
	q.tryConsume()
}

// try to consume message
func (q *PickupQueue) tryConsume() {
	log.Println("here2")
	for {
		log.Printf("queue:%d", len(q.queue))
		if len(q.queue) == 0 {
			break
		}
		truck, err := service.GetFirstIdleTruck()
		if err != nil {
			log.Println(err)
			break
		}
		q.Pickup(truck.ID)
	}
}

// do pick up
func (q *PickupQueue) Pickup(truckID model.TruckID) {
	req := <-q.queue

	// send world pickup
	err := service.LinkTruckToPackage(string(req.PackageID), uint(truckID))
	if err != nil {
		log.Fatalln(err)
		return
	}

	err = VnetCtrl.Sender.SendWorldRequestToGoPickUp(truckID, req.WarehouseID) //Do World Command // 用thread maybe?
	if err != nil {
		log.Println("Error sending GoPickUp command:", err)
	}
}
