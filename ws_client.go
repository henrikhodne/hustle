package hustle

import (
	"fmt"
	"log"
	"strings"

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

	inMsgChan  chan *wsMessage
	outMsgChan chan *wsMessage
	doneChan   chan bool

	socketID string
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

		inMsgChan:  make(chan *wsMessage),
		outMsgChan: make(chan *wsMessage),
		doneChan:   make(chan bool),

		socketID: sha1Sum(fmt.Sprintf("%s-%s-%d", ws.RemoteAddr().String(),
			ws.LocalAddr().String(), maxID)),
	}
}

func (c *wsClient) SocketID() string {
	return c.socketID
}

func (c *wsClient) Listen() {
	log.Printf("client %d listening\n", c.id)
	go c.listenIncoming()
	c.sendPayload("", "pusher:connection_established", &eventPayload{
		SocketID: c.socketID,
	})
	c.listenOutgoing()
}

func (c *wsClient) listenIncoming() {
	log.Printf("client %d listening for incoming messages\n", c.id)
	for {
		msg := newWsMessage()
		websocket.JSON.Receive(c.ws, msg)
		log.Printf("client %d received incoming message %#v\n", c.id, msg)
		if strings.HasPrefix(msg.Event, "client-") {
			c.publishClientEvent(msg)
		} else {
			switch msg.Event {
			case "pusher:ping":
				c.pusherPing(msg)
			case "pusher:pong":
				c.pusherPong(msg)
			case "pusher:subscribe":
				c.pusherSubscribe(msg)
			case "pusher:unsubscribe":
				c.pusherUnsubscribe(msg)
			}
		}
	}
}

func (c *wsClient) listenOutgoing() {
	log.Printf("client %d listening for outgoing messages\n", c.id)
	for {
		select {
		case msg := <-c.outMsgChan:
			log.Printf("client %d received outgoing message\n", c.id)
			websocket.JSON.Send(c.ws, msg)
		}
	}
}

func (c *wsClient) pusherPing(msg *wsMessage) {
	c.sendPayload("", "pusher:pong", nil)
}

func (c *wsClient) pusherPong(msg *wsMessage) {}

func (c *wsClient) pusherSubscribe(msg *wsMessage) {
	log.Printf("adding subscription via %#v\n", msg)

	channelID := msg.Channel
	if channelID == "" {
		switch msg.Data.(type) {
		case map[string]interface{}:
			if value, ok := msg.Data.(map[string]interface{})["channel"]; ok {
				channelID = value.(string)
			}
		}
	}

	if channelID == "" {
		c.sendError(-1, fmt.Sprintf("no channel id present"))
		return
	}

	msg.Channel = channelID

	if _, ok := c.subs[channelID]; ok {
		c.sendError(-1,
			fmt.Sprintf("Existing subscription to channel %s", channelID))
		return
	}

	err := c.sendPayload(msg.Channel,
		"pusher_internal:subscription_succeeded", nil)
	if err != nil {
		log.Printf("error subscribing: %v\n", err)
		return
	}

	subID := newWsSubscription(c.socketID, c.h, msg).Subscribe(c.outMsgChan)

	if subID == "" {
		c.sendError(-1,
			fmt.Sprintf("failed to add subscription to channel %s", channelID))
		return
	}

	c.subs[channelID] = subID
	log.Printf("client %d subscribed to %s with subscription ID %s\n",
		c.id, channelID, c.subs[channelID])
}

func (c *wsClient) pusherUnsubscribe(msg *wsMessage) {
	var (
		subsID string
		ok     bool
	)

	channelID := msg.Channel
	if subsID, ok = c.subs[channelID]; ok {
		delete(c.subs, channelID)
	}

	c.h.Unsubscribe(subsID)
}

func (c *wsClient) sendError(code int, message string) {
	websocket.JSON.Send(c.ws, &wsErrMsg{code, message})
}

func (c *wsClient) sendPayload(channel, event string, payload interface{}) error {
	return websocket.JSON.Send(c.ws, &eventPayload{
		Event:   event,
		Data:    payload,
		Channel: channel,
	})
}

func (c *wsClient) publishClientEvent(msg *wsMessage) {
	log.Printf("publishing client event %#v\n", msg)
	response, err := c.h.PublishEvent(&eventPayload{
		Event:    msg.Event,
		Channel:  msg.Channel,
		SocketID: msg.SocketID,
		Data:     msg.Data,
	})
	if err != nil {
		c.srv.Err(err)
	}
	log.Printf("publish response: %#v\n", response)
}
