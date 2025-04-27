package vnetcontroller

import (
	"log"
	"mini-ups/model"
	"mini-ups/service"
	"sync"

	"gorm.io/datatypes"
)

type PickupQueue struct {
	queue            chan *PickupReq
	prioritizedQueue chan *PickupReq
	waitingMap       map[string]*PickupReq
	mu               sync.Mutex
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

var VnetCtrl *Controller
var PkQueue = NewPickupQueue(100)

func NewPickupQueue(size uint) *PickupQueue {
	return &PickupQueue{
		queue:            make(chan *PickupReq, size),
		prioritizedQueue: make(chan *PickupReq, size),
		waitingMap:       make(map[string]*PickupReq),
	}
}

// add new req to regular queue
func (q *PickupQueue) AddRequest(req *PickupReq) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.waitingMap[req.PackageID] = req
	q.queue <- req
	q.tryConsume()
}

// add to prioritized queue
// use a set to record this package, avoiding repeating pickup
func (q *PickupQueue) PrioritizePackage(packageID string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// check if valid
	if _, exists := q.waitingMap[packageID]; !exists {
		return
	}

	// write to db
	service.PrioritizePackage(packageID)

	// line up
	q.prioritizedQueue <- q.waitingMap[packageID]
}

// try consume with lock
func (q *PickupQueue) TryConsumeWithLock() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.tryConsume()
}

// try to consume message
func (q *PickupQueue) tryConsume() {

	for {
		log.Printf("priority size:%d, queue size:%d", len(q.prioritizedQueue), len(q.queue))

		// check unprocessed req
		if len(q.prioritizedQueue) == 0 && len(q.queue) == 0 {
			break
		}

		// get available truck
		truck, err := service.GetFirstIdleTruck()
		if err != nil {
			log.Println(err)
			break
		}

		// assign by priority
		if len(q.prioritizedQueue) > 0 {
			req := <-q.prioritizedQueue
			if _, exists := q.waitingMap[req.PackageID]; exists {
				delete(q.waitingMap, req.PackageID)
				q.pickup(req, truck.ID)
			}
			continue
		} else if len(q.queue) > 0 {
			req := <-q.queue
			if _, exists := q.waitingMap[req.PackageID]; exists {
				delete(q.waitingMap, req.PackageID)
				q.pickup(req, truck.ID)
			}
			continue
		}
	}
}

// do pick up
func (q *PickupQueue) pickup(req *PickupReq, truckID model.TruckID) {
	log.Println("pick up:", req.PackageID, "truck:", truckID)
	// link package to truck
	err := service.LinkTruckToPackage(string(req.PackageID), uint(truckID))
	if err != nil {
		log.Fatalln(err)
		return
	}

	// send world pickup command
	err = VnetCtrl.Sender.SendWorldRequestToGoPickUp(truckID, req.WarehouseID)
	if err != nil {
		log.Println("Error sending GoPickUp command:", err)
	}
}

// helper function to print out the lengths of the queues for debugging.
func (q *PickupQueue) PrintLengths() {
	log.Print(len(q.queue))
	log.Print(len(q.prioritizedQueue))
}
