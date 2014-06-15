package api_ws

import (
    "fmt"
)

//------------------------------------------------------------
// WebSocket broadcast
//------------------------------------------------------------

// Returns public broadcaster.
func GetBroadcaster_Public() (func (string, interface{}) (err error)) {
    return broadcastPublic
}

// Returns authenticated broadcaster.
func GetBroadcaster_Authd() (func (string, interface{}) (err error)) {
    return broadcastAuthd
}

// Broadcasts to open public websocket connections.
func broadcastPublic(op string, data interface{}) (err error) {

    fmt.Println(DumpConnsPublic("BROADCAST PUBLIC"))

    for _, ws := range _connsPublic {
        if ws.isOpen {
            ws.chanSend <- &Broadcast{op, data}
        }
    }

    return
}

// Broadcasts to open authenticated websocket connections.
func broadcastAuthd(op string, data interface{}) (err error) {

    fmt.Println(DumpConnsPublic("BROADCAST AUTHD"))

    for _, conns := range _connsAuthd {
        for _, ws := range conns {
            if ws.isOpen {
                ws.chanSend <- &Broadcast{op, data}
            }
        }
    }

    return
}
