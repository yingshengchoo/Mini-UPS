package vnetcontroller

import (
	"mini-ups/util"
	"sync"
)

//Recev Window objecct checks and records all seqnum number received from the world simulation
// In the case that a pacakge is lost, we can ensure that no operation is done twice.

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
