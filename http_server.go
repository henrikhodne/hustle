package hustle

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
)

type httpServer struct {
	addr string
	m    *martini.ClassicMartini
}

// HTTPServerMain is the whole shebang for the HTTP, mannn
func HTTPServerMain(addr string, ) {
	srv := newHTTPServer(addr)
	srv.Listen()
}

func newHTTPServer(addr string) *httpServer {
	return &httpServer{
		addr: addr,
		m:    martini.Classic(),
	}
}

func (srv *httpServer) Listen() {
	srv.setupMiddleware()
	srv.setupRoutes()
	log.Printf("hustle-server HTTP listening at %s\n", srv.addr)
	log.Fatal(http.ListenAndServe(srv.addr, srv.m))
}

func (srv *httpServer) setupMiddleware() {
	srv.m.Use(render.Renderer())
	srv.m.Use(martini.Logger())
}

func (srv *httpServer) setupRoutes() {
	srv.m.Get(`/pusher/info`, srv.getPusherInfo)
	srv.m.Post(`/pusher/**`, srv.createUnknownThing)

	srv.m.Post(`/apps/:app_id/events`, srv.createAppEvents)
	srv.m.Get(`/apps/:app_id/channels`, srv.getAppChannels)
	srv.m.Get(`/apps/:app_id/channels/:channel_name`, srv.getAppChannel)
	srv.m.Post(`/apps/:app_id/channels/:channel_name/events`, srv.createAppChannelEvents)
	srv.m.Get(`/apps/:app_id/channels/:channel_name/users`, srv.getAppChannelUsers)
}

func (srv *httpServer) getPusherInfo() string {
	return fmt.Sprintf(`{
	"hostname": "localhost",
	"websocket": false,
	"origins": ["*:*"],
	"cookie_needed": false,
	"entropy": %v,
	"server_heartbeat_interval": 25000
  }`, rand.Int())
}

func (srv *httpServer) getAppChannels() string {
	return `{"channels": {}}`
}

func (srv *httpServer) getAppChannel() string {
	return `{}`
}

func (srv *httpServer) getAppChannelUsers() string {
	return `{"users": []}`
}

func (srv *httpServer) createAppEvents(req *http.Request) string {
	dumpRequest(req)
	return `{}`
}

func (srv *httpServer) createAppChannelEvents(req *http.Request) string {
	dumpRequest(req)
	return `{}`
}

func (srv *httpServer) createUnknownThing(r render.Render, req *http.Request) {
	dumpRequest(req)
	r.JSON(200, req)
}

func dumpRequest(req *http.Request) {
	if req != nil {
		log.Printf("request: %#v\n", req)
		var body bytes.Buffer
		io.Copy(&body, req.Body)
		log.Printf("body: %s\n", string(body.Bytes()))
	}
}
