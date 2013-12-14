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

	inMsgChan  chan *wsMessage
	outMsgChan chan *wsMessage
	doneChan   chan bool
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
	log.Printf("client %d listening\n", c.id)
	go c.channelizeIncomingMessages()
	go c.channelizeOutgoingMessages()
	go c.listenIncoming()
	c.sendPayload("", "pusher:connection_established", &eventPayload{
		SocketID: fmt.Sprintf("%v", c.id),
	})
	c.listenOutgoing()
}

func (c *wsClient) listenIncoming() {
	log.Printf("client %d listening for incoming messages\n", c.id)
	for {
		select {
		case <-c.doneChan:
			return
		case msg := <-c.inMsgChan:
			switch msg.Event {
			case "pusher_ping":
				c.pusherPing(msg)
			case "pusher_pong":
				c.pusherPong(msg)
			case "pusher_subscribe":
				c.pusherSubscribe(msg)
			case "pusher_unsubscribe":
				c.pusherUnsubscribe(msg)
			}
		}
	}
}

func (c *wsClient) channelizeIncomingMessages() {
	log.Printf("client %d setting up incoming message channel\n", c.id)
	for {
		select {
		case <-c.doneChan:
			return
		}

		msg := newWsMessage()
		log.Printf("client %d waiting to receive from %v\n", c.id, c.ws)
		websocket.JSON.Receive(c.ws, msg)
		log.Printf("client %d received message %#v\n", c.id, msg)
		c.inMsgChan <- msg
	}
}

func (c *wsClient) listenOutgoing() {
	log.Printf("client %d listening for outgoing messages\n", c.id)
	for {
		select {
		case <-c.doneChan:
			return
		}
	}
}

func (c *wsClient) channelizeOutgoingMessages() {
	log.Printf("client %d setting up outgoing message channel\n", c.id)
	for {
		select {
		case <-c.doneChan:
			return
		case msg := <-c.outMsgChan:
			websocket.JSON.Send(c.ws, msg)
		}
	}
}

func (c *wsClient) pusherPing(msg *wsMessage) {
	c.sendPayload("", "pusher:pong", nil)
}

func (c *wsClient) pusherPong(msg *wsMessage) {}

func (c *wsClient) pusherSubscribe(msg *wsMessage) {
	channelID := msg.Data.Channel
	if _, ok := c.subs[channelID]; ok {
		c.sendError(-1,
			fmt.Sprintf("Existing subscription to channel %s", channelID))
	}

	c.subs[channelID] = newWsSubscription(c.ws, c.h, msg).Subscribe()
	log.Printf("client %d subscribed to %s with subscription ID %s\n",
		c.id, channelID, c.subs[channelID])
	c.sendPayload(channelID, "pusher_internal:subscription_succeeded", nil)
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

func (c *wsClient) sendPayload(channel, event string, payload interface{}) {
	websocket.JSON.Send(c.ws, &eventPayload{
		Event:   event,
		Data:    payload,
		Channel: channel,
	})
}
