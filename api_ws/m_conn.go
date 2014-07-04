package api_ws

import (
	"net/http"

	"github.com/deze333/vroom/reqres"
	"github.com/gorilla/websocket"
)

//------------------------------------------------------------
// Websocket connection
//------------------------------------------------------------

type Conn struct {
	w    http.ResponseWriter
	r    *http.Request
	conn *websocket.Conn

	router *reqres.WebSocket_Router

	id      string
	isAuthd bool
	isOpen  bool

	chanIn              chan []byte
	chanOut             chan []byte
	chanProcClose       chan int
	chanProcWriterClose chan int
}
