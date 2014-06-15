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
    r         *http.Request
	conn      *websocket.Conn

    router    *Router

	id        int64
	isAuthd   bool
	isOpen    bool
	chanClose chan int
	chanSend  chan *Broadcast
}
