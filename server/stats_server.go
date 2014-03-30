package hustle

import (
	"log"
	"net/http"

	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
)

// StatsServerMain is the whole shebang for the stats, mannn
func StatsServerMain(cfg *Config) {
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(martini.Logger())
	m.Use(CORSAllowAny())

	m.Get(`/timeline/:id`, handleStatsJSONP)

	log.Printf("hustle-server Stats HTTP listening at %s\n", cfg.StatsAddr)
	log.Fatal(http.ListenAndServe(cfg.StatsAddr, m))
}

func handleStatsJSONP(r render.Render) {
	r.JSON(200, map[string]string{})
}
