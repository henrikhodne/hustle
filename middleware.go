package hustle

import (
	"net/http"

	"github.com/codegangsta/martini"
)

// CORSAllowAny sets Access-Control-Allow-Origin: *
func CORSAllowAny() martini.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
}
