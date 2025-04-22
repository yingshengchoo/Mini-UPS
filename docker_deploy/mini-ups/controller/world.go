package controller

import (
	"fmt"
	"mini-ups/protocol/worldupspb"
	"mini-ups/util"
)

// keeps listening for world responses
// called in main as goroutine to keep listening in the background.
func ParseWorldResponse() {
	for {
		resp := &worldupspb.UResponses{}
		if err := util.RecvMsg(util.UPSConn, resp); err != nil {
			fmt.Println("Error receiving world response:", err)
			continue
		}

		HandleCompletions(resp.GetCompletions())
		HandleDeliveries(resp.GetDelivered())
		HandleAcks(resp.GetAcks())
		HandleTruckStatus(resp.GetTruckstatus())
		HandleErrors(resp.GetError())

		if finished := resp.GetFinished(); finished {
			fmt.Println("World server disconnected")
			// 這裡要做什麼? 直接結束program嗎?
		}
	}
}

func HandleCompletions(completions []*worldupspb.UFinished) {
}

func HandleDeliveries(delivered []*worldupspb.UDeliveryMade) {
}

func HandleAcks(acks []int64) {
}

func HandleTruckStatus(statusList []*worldupspb.UTruck) {
}

func HandleErrors(errors []*worldupspb.UErr) {
}
