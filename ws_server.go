package hustle

import (
	"log"
	"net/http"

	"code.google.com/p/go.net/websocket"
	//"github.com/codegangsta/martini"
)

type wsServer struct {
	errChan     chan error
	doneChan    chan bool
	messages    []*wsMessage
	clients     map[int]*wsClient
	addChan     chan *wsClient
	delChan     chan *wsClient
	sendAllChan chan *wsMessage
}

// WSServerMain is the whole shebang for Web Sockets
func WSServerMain(addr string) {
	//m := martini.Classic()
	//m.Use(martini.Logger())

	//m.Get(`/app/:app_id`, websocket.Handler(serveWs))

	srv := newWsServer()
	srv.Listen(addr)
}

func newWsServer() *wsServer {
	srv := &wsServer{
		errChan:     make(chan error),
		doneChan:    make(chan bool),
		messages:    []*wsMessage{},
		clients:     make(map[int]*wsClient),
		addChan:     make(chan *wsClient),
		delChan:     make(chan *wsClient),
		sendAllChan: make(chan *wsMessage),
	}
	return srv
}

func (srv *wsServer) Add(c *wsClient) {
	srv.addChan <- c
}

func (srv *wsServer) Del(c *wsClient) {
	srv.delChan <- c
}

func (srv *wsServer) SendAll(msg *wsMessage) {
	srv.sendAllChan <- msg
}

func (srv *wsServer) Done() {
	srv.doneChan <- true
}

func (srv *wsServer) Err(err error) {
	srv.errChan <- err
}

func (srv *wsServer) sendPastMessages(c *wsClient) {
	for _, msg := range srv.messages {
		c.Write(msg)
	}
}

func (srv *wsServer) sendAll(msg *wsMessage) {
	for _, c := range srv.clients {
		c.Write(msg)
	}
}

func (srv *wsServer) Listen(addr string) {
	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				srv.errChan <- err
			}
		}()

		client := &wsClient{ws: ws, srv: srv}
		srv.Add(client)
		client.Listen()
	}

	http.Handle(`/`, websocket.Handler(onConnected))

	go func() {
		log.Printf("hustle-server WS listening at %s\n", addr)
		log.Fatal(http.ListenAndServe(addr, nil))
	}()

	for {
		select {
		case c := <-srv.addChan:
			srv.clients[c.id] = c
			log.Printf("Added client %d\n", c.id)
			log.Printf("%d clients connected\n", len(srv.clients))
			srv.sendPastMessages(c)
		case c := <-srv.delChan:
			delete(srv.clients, c.id)
			log.Printf("Deleted client %d\n", c.id)
		case msg := <-srv.sendAllChan:
			log.Println("Send all:", msg)
			srv.messages = append(srv.messages, msg)
			srv.sendAll(msg)
		case err := <-srv.errChan:
			log.Println("Error: ", err.Error())
		case <-srv.doneChan:
			return
		}
	}
}
