package api_ws

import (
    "fmt"
    "net/http"
    "sync/atomic"
    "github.com/deze333/vroom/auth"
)

//------------------------------------------------------------
// Registry of active WebSocket connections
//------------------------------------------------------------

var _openConns_Authd = map[int64]*Conn{}
var _openConns_Public = map[int64]*Conn{}
var _publicConnId int64 = 0

//------------------------------------------------------------
// Functions
//------------------------------------------------------------

// Registers connection based on session ID.
func RegisterConn(w http.ResponseWriter, r  *http.Request, ws *Conn) {

    if ws.isAuthd {
        // Authenticated connection
        id := auth.GetAuthdId(r)
        fmt.Println("New Authd conn, id =", id)
        if id == -1 {
            return
        }

        ws.id = id
        _openConns_Authd[id] = ws

    } else {
        // Public connection
        id := atomic.AddInt64(&_publicConnId, 1)
        ws.id = id
        _openConns_Public[id] = ws
    }
}

// Deregisteres connection by removing it from the registry.
func DeregisterConn(ws *Conn) {
    if ws.isAuthd {
        delete(_openConns_Authd, ws.id)
    } else {
        delete(_openConns_Public, ws.id)
    }
}

// Closes authenticated connection by ID.
func CloseAuthdConn(id int64) {

    fmt.Println("\t Close authd WebSockets for session ID =", id)

    if ws, ok := _openConns_Authd[id]; ok {
        fmt.Println("\t\t Closed authd WebSocket for session ID =", id)
        ws.conn.Close()
        DeregisterConn(ws)
    }
}
