package vnetcontroller

import (
	"mini-ups/util"
	"net"
)

type Controller struct {
	receiver   *Receiver
	Sender     *Sender
	recvWindow *RecvWindow
	sendWindow *SendWindow
}


func NewController(conn net.Conn) *Controller {
	recvWindow := NewRecvWindow()
	sendWindow := NewSendWindow()
	sender := NewSender(conn, sendWindow)
	receiver := NewReceiver(recvWindow, sendWindow)

	return &Controller{
		receiver:   receiver,
		Sender:     sender,
		recvWindow: recvWindow,
		sendWindow: sendWindow,
	}
}

func (c *Controller) Start() {
	go c.receiver.ListenForWorldResponses(util.UPSConn)
}
