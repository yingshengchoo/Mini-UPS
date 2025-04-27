package vnetcontroller

import (
	"container/list"
	"log"
	"mini-ups/protocol/worldupspb"
	"sync"
	"time"
)

//The sendWindow object maintains the messages sent to the world simulation.
// In the case the a package is lost from us to world, we can resent it.
// Once it is acknlowdged, it is removed.

type SendWindow struct {
	mu       sync.Mutex
	seqMap   map[int64]*list.Element // seqnum → node
	respList *list.List              // linked list of responses
}

// ResponseNode representa a node in our linkedlist.
type ResponseNode struct {
	SeqNum    int64
	MsgType   string //e.g. "UGoPickup", "UQuery", "UTruck"
	Msg       interface{}
	TimeAdded time.Time
}

// creates a new send window
func NewSendWindow() *SendWindow {
	return &SendWindow{
		seqMap:   make(map[int64]*list.Element),
		respList: list.New(),
	}
}

// Adds the message to our objecct.
func (sw *SendWindow) Add(seqnum int64, msgType string, msg interface{}) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	node := &ResponseNode{
		SeqNum:    seqnum,
		MsgType:   msgType,
		Msg:       msg,
		TimeAdded: time.Now(),
	}
	elem := sw.respList.PushBack(node)
	sw.seqMap[seqnum] = elem
}

// Finds the message corresponding to the seqnum.
func (sw *SendWindow) GetResponse(seqnum int64) (string, interface{}) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	if elem, ok := sw.seqMap[seqnum]; ok {
		node := elem.Value.(*ResponseNode)
		return node.MsgType, node.Msg
	}
	return "", nil // or an error if you prefer
}

// removes the message from our datastructure when world has confirmed receiving it.
func (sw *SendWindow) Ack(seqnum int64) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	if elem, ok := sw.seqMap[seqnum]; ok {
		sw.respList.Remove(elem)
		delete(sw.seqMap, seqnum)
	}
}

// Resends stale responses to world.
func (sw *SendWindow) ResendStaleResponses(pickups []*worldupspb.UGoPickup, queries []*worldupspb.UQuery, deliveries []*worldupspb.UGoDeliver, seqnums []int64, resendFunc func(seqnum int64, msgType string, msg interface{})) {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()
	for e := sw.respList.Front(); e != nil; {
		node := e.Value.(*ResponseNode)
		next := e.Next()

		if now.Sub(node.TimeAdded) > 5*time.Second {
			// Resend the message
			log.Printf("resend: %d, %s", node.SeqNum, node.MsgType)

			resendFunc(node.SeqNum, node.MsgType, node.Msg)

			// Remove the node
			// sw.respList.Remove(e)
			// delete(sw.seqMap, node.SeqNum)
		}

		e = next
	}
}
