package hustle

import (
	"fmt"
	"log"

	"code.google.com/p/go.net/websocket"
)

type wsSubscription struct {
	ws    *websocket.Conn
	h     *hub
	msg   *wsMessage
	qChan chan bool

	socketID string
}

func newWsSubscription(ws *websocket.Conn, h *hub, msg *wsMessage) *wsSubscription {
	if ws == nil {
		log.Panic("ws cannot be nil")
	}

	if h == nil {
		log.Panic("hub cannot be nil")
	}

	return &wsSubscription{
		ws:  ws,
		h:   h,
		msg: msg,

		qChan:    make(chan bool),
		socketID: fmt.Sprintf("%s", ws.RemoteAddr().String()),
	}
}

func (wsSub *wsSubscription) Subscribe() string {
	err := wsSub.sendPayload(wsSub.msg.Channel,
		"pusher_internal:subscription_succeeded", nil)
	if err != nil {
		log.Printf("error subscribing: %v\n", err)
		return wsSub.socketID
	}

	go wsSub.subscribeForever()

	return wsSub.socketID
}

func (wsSub *wsSubscription) subscribeForever() {
	subChan, subQuitChan := wsSub.h.Subscribe(wsSub.msg.Channel, wsSub.socketID)
	if subChan == nil {
		return
	}

	for {
		select {
		case <-wsSub.qChan:
			subQuitChan <- true
			return
		case m := <-subChan:
			wsSub.sendPayload(m.Channel, "no_idea_what_this_should_be", m.Data)
		}
	}
}

func (wsSub *wsSubscription) sendPayload(channelID, eventName string, payload interface{}) error {
	return websocket.JSON.Send(wsSub.ws, &eventPayload{
		Event:    eventName,
		Channel:  channelID,
		Data:     payload,
		SocketID: wsSub.socketID,
	})
}
