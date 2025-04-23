package vnetcontroller

import (
	"container/list"
	"sync"
)

type SendWindow struct {
	mu       sync.Mutex
	seqMap   map[int64]*ResponseNode // seqnum → node
	respList *list.List              // linked list of responses
}

type ResponseNode struct {
	SeqNum int64
	Msg    interface{}
}

func NewSendWindow() *SendWindow {
	return &SendWindow{
		seqMap:   make(map[int64]*ResponseNode),
		respList: list.New(),
	}
}

func (sw *SendWindow) Add(seqnum int64, msg interface{}) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	node := &ResponseNode{SeqNum: seqnum, Msg: msg}
	sw.respList.PushBack(node) //don't need to maintain any sort of order since we have the map
	sw.seqMap[seqnum] = node
}

func (sw *SendWindow) GetResponse(seqnum int64) interface{} {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return sw.seqMap[seqnum].Msg
}

func (sw *SendWindow) Ack(seqnum int64) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	if _ /*node*/, ok := sw.seqMap[seqnum]; ok {
		for e := sw.respList.Front(); e != nil; e = e.Next() {
			if e.Value.(*ResponseNode).SeqNum == seqnum {
				sw.respList.Remove(e)
				break
			}
		}
		delete(sw.seqMap, seqnum)
	}
}
