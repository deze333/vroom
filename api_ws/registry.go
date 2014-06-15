package api_ws

import (
    "fmt"
    "bytes"
    "net/http"
    "github.com/deze333/vroom/auth"
)

//------------------------------------------------------------
// Registry of active WebSocket connections
//------------------------------------------------------------

/*
var _openConns_Authd = map[int64]*Conn{}
var _openConns_Public = map[int64]*Conn{}
var _publicConnId int64 = 0
*/

var _connsPublic = map[*http.Request]*Conn{}
var _connsAuthd = map[int64]map[*http.Request]*Conn{}

//------------------------------------------------------------
// Functions
//------------------------------------------------------------

// Registers open connection based on session ID.
func RegisterConn(ws *Conn) {

    if ws.isAuthd {

        // Authenticated connection
        if id := auth.GetAuthdId(ws.r); id != -1 {
            ws.id = id

            if conns, ok := _connsAuthd[id]; ok {
                conns[ws.r] = ws
            } else {
                _connsAuthd[id] = map[*http.Request]*Conn{ws.r:ws}
            }

            // Debug
            fmt.Println(DumpConnsAuthd("REGISTERED CONNS: AUTHD"))

        } else {

            // Can't find ID for authenticated connection
            _onPanic(
                fmt.Sprintf("Cannot find ID for authd connection"),
                fmt.Sprintf("%v", ws.r), "", "")
        }

    } else {

        // Public connection
        _connsPublic[ws.r] = ws

        // Debug
        fmt.Println(DumpConnsPublic("REGISTERED CONNS: PUBLIC"))
    }
}

// Deregisteres connection by removing it from the registry.
func DeregisterConn(ws *Conn) {
    if ws.isAuthd {
        delete(_connsAuthd, ws.id)
    } else {
        delete(_connsPublic, ws.r)
    }
}

// Closes authenticated connection by ID.
func CloseAuthdConn(id int64) {

    if conns, ok := _connsAuthd[id]; ok {

        // Multiple connections may share same authentication ID
        for _, ws := range conns {
            fmt.Printf("__X Closed WebSocket ID = %v, conn = %p\n", id, ws.conn)
            ws.conn.Close()
        }

        DeregisterConn(&Conn{isAuthd: true, id: id})
    }
}

func DumpConnsPublic(header string) string {
    var buf bytes.Buffer

    buf.WriteString("\n")
    buf.WriteString("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")
    buf.WriteString(header)
    buf.WriteString("\n")
    buf.WriteString("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")

    i := 0
    for r, ws := range _connsPublic {
        buf.WriteString(fmt.Sprintf("\t%v: r = %p, wsconn = %p\n", i, r, ws.conn))
        i++
    }

    return buf.String()
}

func DumpConnsAuthd(header string) string {
    var buf bytes.Buffer

    buf.WriteString("\n")
    buf.WriteString("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")
    buf.WriteString(header)
    buf.WriteString("\n")
    buf.WriteString("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")

    for i, conns := range _connsAuthd {
        buf.WriteString(fmt.Sprintf("\t%v:\n", i))
        j := 0
        for r, ws := range conns {
            buf.WriteString(fmt.Sprintf("\t\t%v: r = %p, wsconn = %p\n", j, r, ws.conn))
        }
    }

    return buf.String()
}
