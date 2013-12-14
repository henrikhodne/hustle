package hustle

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/codegangsta/martini-contrib/render"
)

type httpServer struct {
	addr string
	m    *martini.ClassicMartini
	h    *hub
}

// HTTPServerMain is the whole shebang for the HTTP, mannn
func HTTPServerMain(addr string, hubAddr string) {
	srv, err := newHTTPServer(addr, hubAddr)
	if err != nil {
		log.Panicf("oh well: %v\n", err)
	}
	srv.Listen()
}

func newHTTPServer(addr, hubAddr string) (*httpServer, error) {
	h, err := newHub(hubAddr)
	if err != nil {
		return nil, err
	}
	return &httpServer{
		addr: addr,
		m:    martini.Classic(),
		h:    h,
	}, nil
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

	srv.m.Post(`/apps/:app_id/events`, binding.Json(Event{}), srv.createAppEvents)
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

func (srv *httpServer) createAppEvents(evt Event, err binding.Errors, resp http.ResponseWriter, req *http.Request) string {
	if err.Count() > 0 {
		resp.WriteHeader(http.StatusBadRequest)
		return fmt.Sprintf(`{"errors":"%v"}`, err)
	}
	socketID := req.URL.Query().Get("socket_id")
	log.Printf("received event: %#v\n", evt)
	for _, channel := range evt.Channels {
		_, pubErr := srv.h.PublishEvent(channel, evt.Name, socketID, evt.Data)
		if pubErr != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			return `{}`
		}
	}
	resp.WriteHeader(http.StatusAccepted)
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

func (srv *httpServer) publishEvent(channel, name string, data []byte) {
}

func dumpRequest(req *http.Request) {
	if req != nil {
		log.Printf("request: %#v\n", req)
		var body bytes.Buffer
		io.Copy(&body, req.Body)
		log.Printf("body: %s\n", string(body.Bytes()))
	}
}
