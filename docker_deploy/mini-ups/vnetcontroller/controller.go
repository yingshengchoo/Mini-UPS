package vnetcontroller

import (
	"log"
	"mini-ups/protocol"
	"mini-ups/protocol/worldupspb"
	"net"
	"time"
)

type Controller struct {
	receiver   *Receiver
	Sender     *Sender
	recvWindow *RecvWindow
	sendWindow *SendWindow
	conn       net.Conn
}


func NewController(conn net.Conn) *Controller {
	recvW := NewRecvWindow()
	sendW := NewSendWindow()
	sender := NewSender(recvW, sendW, conn)
	return &Controller{
		receiver:   NewReceiver(recvW, sendW, sender),
		Sender:     sender,
		recvWindow: recvW,
		sendWindow: sendW,
		conn:       conn,
	}
}

func (c *Controller) Start() {
	go c.receiver.ListenForWorldResponses(c.conn)

	//Goroutine for resent requests
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			var pickups []*worldupspb.UGoPickup
			var queries []*worldupspb.UQuery
			var deliveries []*worldupspb.UGoDeliver // only if needed
			var seqnums []int64

			c.sendWindow.ResendStaleResponses(func(seqnum int64, msgType string, msg interface{}) {
				seqnums = append(seqnums, seqnum)

				switch msgType {
				case "UQuery":
					if q, ok := msg.(*worldupspb.UQuery); ok {
						queries = append(queries, q)
					}
				case "UGoPickup":
					if p, ok := msg.(*worldupspb.UGoPickup); ok {
						pickups = append(pickups, p)
					}
				case "UGoDeliver":
					if d, ok := msg.(*worldupspb.UGoDeliver); ok {
						deliveries = append(deliveries, d)
					}
				default:
					log.Printf("Unknown message type: %s\n", msgType)
				}
			})

			// If any messages were collected, batch resend
			if len(pickups) > 0 || len(queries) > 0 {
				cmd := protocol.CreateUPSCommands(pickups, deliveries, 0, false, queries, nil) // tested: sending empty list doesnt affect functionality with world.
				if err := c.Sender.Send(cmd); err != nil {
					log.Println("Resend batch failed:", err)
				}
			}
		}
	}()
}
