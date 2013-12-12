package hustle

import (
	"fmt"
	"net/http"
	"os"

	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
)

// Serve is the whole shebang
func Serve(addr string) {
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(martini.Logger())

	m.Post(`/apps/:app_id/events`, createAppEvents)

	fmt.Fprintf(os.Stderr, "hustle-server listening at %s\n", addr)
	http.ListenAndServe(addr, m)
}

func createAppEvents(r render.Render) {
}
