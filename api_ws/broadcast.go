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

    if l := len(_openConns_Public); l > 0 {
        ids := []int64{}
        for id, _ := range _openConns_Authd {
            ids = append(ids, id)
        }
        fmt.Printf("\n<--- BR (Pub %d: %v) %s = %v\n", l, ids, op, data)
    } else {
        fmt.Printf("\nx--- BR (No Pub)\n")
    }
    for _, ws := range _openConns_Public {
        if ws.isOpen {
            ws.chanSend <- &Broadcast{op, data}
        }
    }

    return
}

// Broadcasts to open authenticated websocket connections.
func broadcastAuthd(op string, data interface{}) (err error) {

    if l := len(_openConns_Authd); l > 0 {
        ids := []int64{}
        for id, _ := range _openConns_Authd {
            ids = append(ids, id)
        }
        fmt.Printf("\n<--- BR (Authd %d: %v) %s = %v\n", l, ids, op, data)
    } else {
        fmt.Printf("\nx--- BR (No Authd)\n")
    }

    for _, ws := range _openConns_Authd {
        if ws.isOpen {
            ws.chanSend <- &Broadcast{op, data}
        }
    }

    return
}
