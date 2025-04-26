package queue

import (
	"log"
	"mini-ups/model"
	"mini-ups/service"
	"mini-ups/vnetcontroller"
	"sync"

	"gorm.io/datatypes"
)

type PickupQueue struct {
	queue         chan *PickupReq
	priorityQueue chan *PickupReq
	consuming     bool
	mu            sync.Mutex
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
var PkQueue = NewPickupQueue(100)

func NewPickupQueue(size uint) *PickupQueue {
	return &PickupQueue{
		queue:         make(chan *PickupReq, size),
		priorityQueue: make(chan *PickupReq, size),
		// availableTrucks: truckNum,
	}
}

// add new req
func (q *PickupQueue) AddRequest(req *PickupReq) {
	// log.Println("here")
	q.queue <- req
	q.mu.Lock()
	if !q.consuming {
		q.consuming = true
		go q.tryConsume()
	}
	q.mu.Unlock()
}

func (q *PickupQueue) PrioritizePackage(packageID string) {
	newQueue := make(chan *PickupReq, cap(q.queue))

	for {
		select {
		case req := <-q.queue:
			if req.PackageID == packageID {
				q.priorityQueue <- req // Move it to priority queue
			} else {
				newQueue <- req
			}
		default:
			q.queue = newQueue
			// q.mu.Lock()
			// if !q.consuming {
			// 	q.consuming = true
			// 	go q.tryConsume()
			// }
			// q.mu.Unlock()
			return
		}
	}
}

// // try to consume message
// func (q *PickupQueue) tryConsume() {
// 	// log.Println("here2")
// 	for {
// 		log.Printf("queue size:%d", len(q.queue))
// 		if len(q.queue) == 0 {
// 			break
// 		}
// 		truck, err := service.GetFirstIdleTruck()
// 		if err != nil {
// 			log.Println(err)
// 			break
// 		}
// 		q.Pickup(truck.ID)
// 	}
// }

func (q *PickupQueue) tryConsume() {
	defer func() {
		q.mu.Lock()
		q.consuming = false
		q.mu.Unlock()
	}()
	for {
		log.Printf("priority size:%d, queue size:%d", len(q.priorityQueue), len(q.queue))

		if len(q.priorityQueue) == 0 && len(q.queue) == 0 {
			break
		}

		truck, err := service.GetFirstIdleTruck()
		if err != nil {
			log.Println(err)
			break
		}

		select {
		case req := <-q.priorityQueue:
			q.Pickup(req, truck.ID)
		case req := <-q.queue:
			q.Pickup(req, truck.ID)
		default:
			return // No available package to process
		}
	}
}

// do pick up
func (q *PickupQueue) Pickup(req *PickupReq, truckID model.TruckID) {
	err := service.LinkTruckToPackage(string(req.PackageID), uint(truckID))
	if err != nil {
		log.Fatalln(err)
		return
	}

	err = VnetCtrl.Sender.SendWorldRequestToGoPickUp(truckID, req.WarehouseID)
	if err != nil {
		log.Println("Error sending GoPickUp command:", err)
	}
}
