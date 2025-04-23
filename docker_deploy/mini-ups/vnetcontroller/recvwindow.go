package vnetcontroller

import (
	"sort"
	"sync"
)

// merge inaterval for ack to world
// 1 2 3
// [[1, 3]]

// 5 6 7
// [[1, 3], [5, 7]]

// 4
// [[1, 7]]

type RecvWindow struct {
	mu      sync.Mutex
	acked   map[int64]bool
	ackList []int64
}

func NewRecvWindow() *RecvWindow {
	return &RecvWindow{
		acked: make(map[int64]bool),
	}
}

// Only ACK new seqnums, skip if already seen
func (rw *RecvWindow) Record(seqnum int64) bool {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	if rw.acked[seqnum] {
		return false
	}
	rw.acked[seqnum] = true
	rw.ackList = append(rw.ackList, seqnum)
	return true
}

// Merge seqnums into intervals like [[1,3], [5,7]]
func (rw *RecvWindow) GetAckRanges() [][2]int64 {
	rw.mu.Lock()
	defer rw.mu.Unlock()

	sort.Slice(rw.ackList, func(i, j int) bool {
		return rw.ackList[i] < rw.ackList[j]
	})

	var ranges [][2]int64
	var start, end int64
	for i, num := range rw.ackList {
		if i == 0 {
			start, end = num, num
			continue
		}
		if num == end+1 {
			end = num
		} else {
			ranges = append(ranges, [2]int64{start, end})
			start, end = num, num
		}
	}
	if len(rw.ackList) > 0 {
		ranges = append(ranges, [2]int64{start, end})
	}
	rw.ackList = nil // reset after sending
	return ranges
}
