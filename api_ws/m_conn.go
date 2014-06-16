package api_ws

import (
    "net/http"
	"github.com/gorilla/websocket"
)

//------------------------------------------------------------
// Websocket connection
//------------------------------------------------------------

type Router struct {
	URL    string
	Procs  map[string]func(*Req) (interface{}, error)
}

type Conn struct {
    w         http.ResponseWriter
    r         *http.Request
	conn      *websocket.Conn

    router    *Router

	id        int64
	isAuthd   bool
	isOpen    bool

    chanIn    chan []byte
    chanOut   chan []byte
	chanProcClose chan int
	chanProcWriterClose chan int
}
