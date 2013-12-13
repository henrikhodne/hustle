package hustle

import (
	"log"
	"net/http"

	"code.google.com/p/go.net/websocket"
	//"github.com/codegangsta/martini"
)

// WSServerMain is the whole shebang for Web Sockets
func WSServerMain(addr string) {
	//m := martini.Classic()
	//m.Use(martini.Logger())

	//m.Get(`/app/:app_id`, websocket.Handler(serveWs))

	http.Handle(`/`, websocket.Handler(serveWs))

	log.Printf("hustle-server WS listening at %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func serveWs(ws *websocket.Conn) {
	for {
		return
	}
}
