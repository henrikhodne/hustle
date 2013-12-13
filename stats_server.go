package hustle

import (
	"log"
	"net/http"

	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
)

// StatsServerMain is the whole shebang for the stats, mannn
func StatsServerMain(addr string) {
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(martini.Logger())

	m.Get(`/timeline/:id`, handleStatsJSONP)

	log.Printf("hustle-server Stats HTTP listening at %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, m))
}

func handleStatsJSONP() string {
	return ""
}
