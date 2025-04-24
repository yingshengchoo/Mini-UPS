package vnetcontroller

import (
	"mini-ups/util"
	"sync"
)

//After careful thought, no point merge values in the list since we never need to actually iterate through the list
//IF world doesn't recieve ack, they resend the response --> we check if ack exist: Yes then send ack, send ack + processe response

type RecvWindow struct {
	mu    sync.Mutex
	acked *util.Int64Set
}

func NewRecvWindow() *RecvWindow {
	return &RecvWindow{
		acked: util.NewInt64Set(),
	}
}

// Record returns true only if seqnum is new
func (rw *RecvWindow) Record(seqnum int64) bool {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	return rw.acked.Add(seqnum)
}
