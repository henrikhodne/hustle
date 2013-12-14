package hustle

import (
	"fmt"
	"io"
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
	ch   chan *wsMessage
	subs map[string]string

	doneChan chan bool
}

type wsErrMsg struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func newClient(ws *websocket.Conn, srv *wsServer) *wsClient {
	if ws == nil {
		log.Panicln("ws cannot be nil")
	}

	if srv == nil {
		log.Panicln("server cannot be nil")
	}

	maxID++

	return &wsClient{
		id:  maxID,
		ws:  ws,
		srv: srv,
		ch:  make(chan *wsMessage, channelBufSize),

		doneChan: make(chan bool),
	}
}

func (c *wsClient) Conn() *websocket.Conn {
	return c.ws
}

func (c *wsClient) Write(msg *wsMessage) {
	select {
	case c.ch <- msg:
	default:
		c.srv.Del(c)
		err := fmt.Errorf("client %d is disconnected", c.id)
		c.srv.Err(err)
	}
}

func (c *wsClient) Listen() {
	go c.listenWrite()
	c.listenRead()
}

func (c *wsClient) listenWrite() {
	log.Printf("client %d listening for outgoing messages\n", c.id)

	for {
		select {
		case msg := <-c.ch:
			log.Printf("client %d received send: %v\n", c.id, msg)
			websocket.JSON.Send(c.ws, msg)
		case <-c.doneChan:
			log.Printf("client %d received `done` in listenWrite\n", c.id)
			c.srv.Del(c)
			c.doneChan <- true
			return
		}
	}
}

func (c *wsClient) listenRead() {
	log.Printf("client %d listening for incoming messages\n", c.id)
	for {
		select {
		case <-c.doneChan:
			log.Printf("client %d received `done` in listenRead\n", c.id)
			c.srv.Del(c)
			c.doneChan <- true
			return
		default:
			log.Printf("client %d reading from ws\n", c.id)
			msg := newWsMessage()
			err := websocket.JSON.Receive(c.ws, msg)
			if err == io.EOF {
				log.Printf("client %d hit EOF, sending `done`\n", c.id)
				c.doneChan <- true
			} else if err != nil {
				c.srv.Err(err)
			} else {
				c.srv.SendAll(msg)
			}
		}
	}
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
