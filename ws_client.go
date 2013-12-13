package hustle

import (
	//"bytes"
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
	id  int
	ws  *websocket.Conn
	srv *wsServer
	ch  chan *wsMessage

	doneChan chan bool
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
	log.Printf("client %d listening for writes\n", c.id)

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
	log.Printf("client %d listening for reads\n", c.id)
	for {
		select {
		case <-c.doneChan:
			log.Printf("client %d received `done` in listenRead\n", c.id)
			c.srv.Del(c)
			c.doneChan <- true
			return
		default:
			log.Printf("client %d reading from ws\n", c.id)
			var msg wsMessage
			err := websocket.JSON.Receive(c.ws, &msg)
			if err == io.EOF {
				log.Printf("client %d hit EOF, sending `done`\n", c.id)
				c.doneChan <- true
			} else if err != nil {
				c.srv.Err(err)
			} else {
				c.srv.SendAll(&msg)
			}
		}
	}
}

//req := c.ws.Request()
//if req != nil {
//log.Printf("serveWs request: %#v\n", req)
//log.Printf("serveWs request URL: %#v\n", req.URL)
//}
//var outbuf bytes.Buffer
//out := io.MultiWriter(c.ws, &outbuf)
//io.Copy(out, c.ws)
//log.Printf("serveWs received: %#v\n", string(outbuf.Bytes()))
