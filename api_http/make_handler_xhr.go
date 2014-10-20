package api_http

import (
	"net/http"

	"github.com/deze333/vroom/api_xhr"
	"github.com/deze333/vroom/reqres"
)

//------------------------------------------------------------
// XHR Request handler generator
//------------------------------------------------------------

// Creates new XHR handler out of handler function.
func makeHandler_XHR(ctx *Ctx, fn func(*reqres.Req) (interface{}, error), needsAuth bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// If HEAD do nothing
		if r.Method == "HEAD" {
			return
		}

		// Cookie header
		w.Header().Add("Vary", "Cookie")
		w.Header().Set("Content-Type", "application/json")

		// Client needs to be authenticated ?
		if needsAuth && !isAuthPassed(w, r, ctx) {
			api_xhr.Handle_NotAuthd(w, r)
			return
		}

		// Call handler
		api_xhr.Handle(w, r, fn)
	}
}
