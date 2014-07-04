package api_http

import (
	"net/http"

	"github.com/deze333/vroom/api_ws"
	"github.com/deze333/vroom/reqres"
)

//------------------------------------------------------------
// WS Request handler generator
//------------------------------------------------------------

// Creates new Websocket handler out of handler function.
func makeHandler_WS(ctx *Ctx, router *reqres.WebSocket_Router, needsAuth bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Cookie header
		//w.Header().Add("Vary", "Cookie")

		// Client needs to be authenticated ?
		if !ctx.Presets.IsDebug && needsAuth && !isAuthPassed(w, r, ctx) {
			api_ws.Handle_NotAuthd(w, r)
			return
		}

		// Call handler
		api_ws.Handle(w, r, router, needsAuth)
	}
}
