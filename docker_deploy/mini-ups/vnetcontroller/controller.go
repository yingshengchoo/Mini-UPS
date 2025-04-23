package vnetcontroller

import "net"

type Controller struct {
	receiver   *Receiver
	sender     *Sender
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
		sender:     sender,
		recvWindow: recvWindow,
		sendWindow: sendWindow,
	}
}

func (c *Controller) Start() {
	go c.receiver.ListenForWorldResponses(c.sender.conn)
}
