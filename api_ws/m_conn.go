package api_ws

import (
	"net/http"
	"sync"
	"time"

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
	m    sync.Mutex

	router *reqres.WebSocket_Router

	agentId   string
	authId    string
	authEmail string
	isAuthd   bool
	isOpen    bool

	chanIn              chan []byte
	chanOut             chan []byte
	chanProcClose       chan int
	chanProcWriterClose chan int
}

type Message struct {
	isAuthd     bool
	agentId     string
	failTime    *time.Time
	isProcessed bool
	req         []byte
	res         []byte
}
