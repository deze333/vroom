package api_http

import (
	"net/http"

	"github.com/deze333/vroom/reqres"
)

//------------------------------------------------------------
// Model
//------------------------------------------------------------

type RouteHandler struct {
	pattern string
	handler func(http.ResponseWriter, *http.Request)
}

type Presets struct {
	// Placeholder for extras
}

type Auth struct {
	CookieName    string
	CookieStoreId string
	CookiePath    string
	CookieDomain  string
	CookieMaxAge  string
}

//------------------------------------------------------------
// Directories
//------------------------------------------------------------

type Dirs struct {

	// Directory where the app resides
	// clients will be notified on change
	// by broadcast "core/broadcast/app/updated"
	AppWatchNotify []string

	// Directory to store version track file
	VersionFileDir string
}

//------------------------------------------------------------
// Route Handlers
//------------------------------------------------------------

type H func(http.ResponseWriter, *http.Request)
type RH map[string]func(http.ResponseWriter, *http.Request)

//type H_XHR func(map[string]interface{}) (interface{}, error)
//type RH_XHR map[string]func(map[string]interface{}) (interface{}, error)

type Handlers struct {
	NotFound H // 404 handler
	NotAuthd H // User not authenticated handler
	AppErr   H // Application error handler
	Public   RH
	Authd    RH
}

type Handlers_XHR struct {
	Public reqres.XHR_Router
	Authd  reqres.XHR_Router
}

type Handlers_WS struct {
	Public []reqres.WebSocket_Router
	Authd  []reqres.WebSocket_Router
}

//------------------------------------------------------------
// Context
//------------------------------------------------------------

type Ctx struct {
	Presets       Presets
	Auth          Auth
	Dirs          Dirs
	OnPanic       func(string, string, ...interface{}) // msg, details, kvals
	Handlers_FILE Handlers
	Handlers_HTML Handlers
	Handlers_XHR  Handlers_XHR
	Handlers_WS   Handlers_WS
}
