package vnetcontroller

type Controller struct {
	receiver   *Receiver
	sender     *Sender
	recvWindow *RecvWindow
	snedWindow *SendWindow
}

func NewController() *Controller {
	return &Controller{
		receiver:   &Receiver{},
		sender:     &Sender{},
		recvWindow: &RecvWindow{},
		snedWindow: &SendWindow{},
	}
}

// func (c *Controller) SendRequestToGoPickUp(,xxx){
// 	sender.SendRequestToGoPickUp(xxx)
// 	window.addSeq()
// }
