package hustle

import (
	"log"
)

type wsSubscription struct {
	h        *hub
	msg      *wsMessage
	socketID string
	qChan    chan bool
}

func newWsSubscription(socketID string, h *hub, msg *wsMessage) *wsSubscription {
	if h == nil {
		log.Panic("hub cannot be nil")
	}

	if msg == nil {
		log.Panic("msg cannot be nil")
	}

	return &wsSubscription{
		h:        h,
		msg:      msg,
		socketID: socketID,

		qChan: make(chan bool),
	}
}

func (wsSub *wsSubscription) Subscribe(outMsgChan chan *wsMessage) string {
	go wsSub.subscribeForever(outMsgChan)

	return wsSub.socketID
}

func (wsSub *wsSubscription) subscribeForever(outMsgChan chan *wsMessage) {
	subChan, _ := wsSub.h.Subscribe(wsSub.msg.Channel, wsSub.socketID)
	if subChan == nil {
		log.Println("cannot use a nil subscription channel")
		return
	}

	for {
		select {
		case msg := <-subChan:
			log.Printf("subscription %s delegating message to outMsgChan\n", wsSub.socketID)
			outMsgChan <- msg
			log.Printf("subscription %s successfully delegated %#v\n", wsSub.socketID, msg)
		}
	}
}
