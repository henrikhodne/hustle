package hustle

import (
	"fmt"
	"log"

	"code.google.com/p/go.net/websocket"
)

const (
	channelBufSize = 100
)

var (
	maxID = int(0)
)

type wsClient struct {
	id   int
	ws   *websocket.Conn
	h    *hub
	srv  *wsServer
	subs map[string]string

	doneChan chan bool
}

type wsErrMsg struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func newClient(ws *websocket.Conn, h *hub, srv *wsServer) *wsClient {
	if ws == nil {
		log.Panicln("ws cannot be nil")
	}

	if srv == nil {
		log.Panicln("server cannot be nil")
	}

	maxID++

	return &wsClient{
		id:   maxID,
		ws:   ws,
		h:    h,
		srv:  srv,
		subs: make(map[string]string),

		doneChan: make(chan bool),
	}
}

func (c *wsClient) Listen() {
	// FIXME: do stuff here, mkay?
	<-c.doneChan
}

func (c *wsClient) pusherSubscribe(msg *wsMessage) {
	channelID := msg.Data.Channel
	if _, ok := c.subs[channelID]; ok {
		c.sendError(-1,
			fmt.Sprintf("Existing subscription to channel %s", channelID))
	}

	c.subs[channelID] = newWsSubscription(c.ws, c.h, msg).Subscribe()
}

func (c *wsClient) pusherUnsubscribe(msg *wsMessage) {
	var (
		subsID string
		ok     bool
	)

	channelID := msg.Data.Channel
	if subsID, ok = c.subs[channelID]; ok {
		delete(c.subs, channelID)
	}

	c.h.Unsubscribe(subsID)
}

func (c *wsClient) sendError(code int, message string) {
	websocket.JSON.Send(c.ws, &wsErrMsg{code, message})
}
