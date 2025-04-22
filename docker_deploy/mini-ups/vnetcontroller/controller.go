package vnetcontroller

type Controller struct {
	receiver *Receiver
	sender   *Sender
	window   *Window
}

func NewController() *Controller {
	return &Controller{
		receiver: &Receiver{},
		sender:   &Sender{},
		window:   &Window{},
	}
}

// func (c *Controller) SendRequestToGoPickUp(,xxx){
// 	sender.SendRequestToGoPickUp(xxx)
// 	window.addSeq()
// }
