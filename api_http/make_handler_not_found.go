package api_http

import (
	"net/http"
)

//------------------------------------------------------------
// Default Not Found handler
//------------------------------------------------------------

// Creates default "not found" handler that tries to 
// guess the request format and answers correspondingly.
func makeHandler_NotFound(ctx *Ctx, handlerHTML H) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        // Error 404, Not Found
        w.WriteHeader(http.StatusNotFound)

        // If HEAD do nothing
        if r.Method == "HEAD" {
            return
        }

        // Catch panic
        defer func() {
            if err := recover(); err != nil {
                handlePanic_HTML(w, r, ctx, err)
            }
        }()

        // Respond according to request type
        switch requestType(r) {

        case "HTML":
            handlerHTML(w, r)

        case "XHR":
            // Do nothing

        case "WS":
            // Do nothing

        default:
            handlerHTML(w, r)
        }
    }
}

