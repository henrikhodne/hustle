package hustle

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/codegangsta/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
)

type httpServer struct {
	cfg *Config
	m   *martini.ClassicMartini
	h   *hub
}

type pusherInfo struct {
	Hostname     string   `json:"hostname"`
	Websocket    bool     `json:"websocket"`
	Origins      []string `json:"origins"`
	CookieNeeded bool     `json:"cookie_needed"`
	Entropy      int      `json:"entropy"`
	Heartbeat    int      `json:"server_heartbeat_interval"`
}

// HTTPServerMain is the whole shebang for the HTTP, mannn
func HTTPServerMain(cfg *Config) {
	srv, err := newHTTPServer(cfg)
	if err != nil {
		log.Fatalf("oh well: %v\n", err)
	}
	srv.Listen()
}

func newHTTPServer(cfg *Config) (*httpServer, error) {
	h, err := newHub(cfg.HubAddr)
	if err != nil {
		return nil, err
	}

	return &httpServer{
		cfg: cfg,
		m:   martini.Classic(),
		h:   h,
	}, nil
}

func (srv *httpServer) Listen() {
	srv.setupMiddleware()
	srv.setupRoutes()
	log.Printf("hustle HTTP listening at %s\n", srv.cfg.HTTPAddr)
	log.Fatal(http.ListenAndServe(srv.cfg.HTTPAddr, srv.m))
}

func (srv *httpServer) setupMiddleware() {
	srv.m.Use(render.Renderer(render.Options{
		Directory: "templates",
		Layout:    "layout",
	}))
	srv.m.Use(martini.Logger())
	srv.m.Use(CORSAllowAny())
}

func (srv *httpServer) setupRoutes() {
	srv.m.Get(`/disco`, srv.getDiscoPage)
	srv.m.Get(`/pusher/info`, srv.getPusherInfo)
	srv.m.Post(`/pusher/**`, srv.createUnknownThing)

	srv.m.Post(`/apps/:app_id/events`, binding.Json(Event{}), srv.createAppEvents)
	srv.m.Get(`/apps/:app_id/channels`, srv.getAppChannels)
	srv.m.Get(`/apps/:app_id/channels/:channel_name`, srv.getAppChannel)
	srv.m.Post(`/apps/:app_id/channels/:channel_name/events`, srv.createAppChannelEvents)
	srv.m.Get(`/apps/:app_id/channels/:channel_name/users`, srv.getAppChannelUsers)
}

func (srv *httpServer) getDiscoPage(r render.Render, resp http.ResponseWriter) {
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	r.HTML(200, "disco", srv.cfg)
}

func (srv *httpServer) getPusherInfo(r render.Render, resp http.ResponseWriter) {
	r.JSON(200, &pusherInfo{
		Hostname:     srv.cfg.WSPubHost(),
		Websocket:    false,
		Origins:      []string{"*:*"},
		CookieNeeded: false,
		Entropy:      rand.Int(),
		Heartbeat:    25000,
	})
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
	for _, payload := range evt.Payloads(socketID) {
		_, pubErr := srv.h.PublishEvent(payload)
		if pubErr != nil {
			log.Printf("ERROR: %v\n", pubErr)
			resp.WriteHeader(http.StatusInternalServerError)
			return `{}`
		}
	}
	resp.WriteHeader(http.StatusAccepted)
	return `{}`
}

func (srv *httpServer) createAppChannelEvents(req *http.Request) string {
	return `{}`
}

func (srv *httpServer) createUnknownThing(r render.Render, req *http.Request) {
	r.JSON(200, req)
}
