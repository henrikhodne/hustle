package hustle

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
)

// HTTPServerMain is the whole shebang for the HTTP, mannn
func HTTPServerMain(addr string) {
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(martini.Logger())

	m.Get(`/pusher/info`, getPusherInfo)

	m.Post(`/apps/:app_id/events`, createAppEvents)
	m.Get(`/apps/:app_id/channels`, getAppChannels)
	m.Get(`/apps/:app_id/channels/:channel_name`, getAppChannel)
	m.Get(`/apps/:app_id/channels/:channel_name/users`, getAppChannelUsers)

	log.Printf("hustle-server HTTP listening at %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, m))
}

func getPusherInfo() string {
	return fmt.Sprintf(`{
	"hostname": "localhost",
	"websocket": false,
	"origins": ["*:*"],
	"cookie_needed": false,
	"entropy": %v,
	"server_heartbeat_interval": 25000
  }`, rand.Int())
}

func getAppChannels() string {
	return `{"channels": {}}`
}

func getAppChannel() string {
	return `{}`
}

func getAppChannelUsers() string {
	return `{"users": []}`
}

func createAppEvents() string {
	return `{}`
}
