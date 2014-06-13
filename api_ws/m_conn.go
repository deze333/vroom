package api_ws

import (
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
    router    *Router

	id        int64
	isAuthd   bool
	isOpen    bool
	conn      *websocket.Conn
	chanClose chan int
	chanSend  chan *Broadcast
}
