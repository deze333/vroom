package api_http

import (
    "fmt"
    "strings"
	"net/http"
    "github.com/deze333/vroom/auth"
    "github.com/deze333/vroom/util"
)

//------------------------------------------------------------
// HTTP Request handler generator
//------------------------------------------------------------

// Creates new HTML handler out of handler function.
func makeHandler_HTML(ctx *Ctx, fn H, needsAuth bool) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        // If HEAD do nothing
        if r.Method == "HEAD" {
            return
        }

        // Cookie header
        w.Header().Add("Vary", "Cookie")

        // Client needs to be authenticated ?
        if ! ctx.Presets.IsDebug && needsAuth && ! isAuthPassed(w, r, ctx) {
            ctx.Handlers_HTML.NotAuthd(w, r)
            return
        }

        // Catch panic
        defer func() {
            if err := recover(); err != nil {
                handlePanic_HTML(w, r, ctx, err)
            }
        }()

        // Call handler
        fn(w, r)
    }
}

// Handles handler's panic.
func handlePanic_HTML(w http.ResponseWriter, r *http.Request, ctx *Ctx, err interface{}) {

    // Ignore broken pipe (client hang up)
    if strings.Contains(fmt.Sprint(err), "broken pipe") {
        return
    }

    // Show provided panic page
    ctx.Handlers_HTML.AppErr(w, r)

    // Report panic: err, url, params, stack
    ctx.OnPanic(
        fmt.Sprintf("Error processing HTML page: %v", err),
        fmt.Sprintf("%v : %v", r.Host, r.RequestURI),
        fmt.Sprint(r),
        fmt.Sprint(auth.GetSessionValues(r)),
        util.Stack())
}

