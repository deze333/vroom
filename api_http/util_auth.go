package api_http

import (
	"net/http"
	"github.com/deze333/vroom/auth"
)

//------------------------------------------------------------
// Utils for authorization
//------------------------------------------------------------

// Handles authentication verification.
// Executes corresponding handler on failure and returns false to
// signal stop for further request processing.
func isAuthPassed(w http.ResponseWriter, r *http.Request, ctx *Ctx) bool {

    // Authenticated ?
    if auth.IsAuthd(r) {
        return true
    }

    // Authentication failed:

    // Close persistent connections
    // TODO Close all open authd WS...

    return false
}

