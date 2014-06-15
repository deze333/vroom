package api_ws

import (
    "fmt"
    "net/http"
    "time"
	"github.com/gorilla/websocket"
    "github.com/deze333/vroom/errors"
    "github.com/deze333/vroom/util"
)

//------------------------------------------------------------
// WebSocket Implementation
//------------------------------------------------------------

// Handles WS Not Authd response.
func Handle_NotAuthd(w http.ResponseWriter, r *http.Request) {

    http.Error(w, "Not authorized", http.StatusUnauthorized)
}

// Creates WS connecttion and processes it in forever loop.
func Handle(w http.ResponseWriter, r *http.Request, router *Router, isAuthd bool) {

    var upgrader = websocket.Upgrader{
        ReadBufferSize:  1024,
        WriteBufferSize: 1024,
        CheckOrigin: func(r *http.Request) bool {
            fmt.Printf("\n\n-------------------------------\nCheck WebSocket call: \nURL = %v \nOrigin = %v\n---------------------------------------\n\n", r.URL, r.Header["Origin"])
            return true
        },
    } 

    // Open websocket connection
    conn, err := upgrader.Upgrade(w, r, nil)

    // Process error
    if _, ok := err.(websocket.HandshakeError); ok {
        http.Error(w, "Not a websocket handshake", http.StatusBadRequest)
        return

    } else if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Zero time means unlimited connection
    err = conn.SetReadDeadline(time.Time{})
    if err != nil {
        _onPanic(
            fmt.Sprintf("Error setting WS deadline: %v", err),
            fmt.Sprintf("%v : %v", r.Host, r.URL),
            "","")
        return
    }

    // Create connection object
    ws := &Conn{
        isAuthd: isAuthd,
        router: router,
        conn: conn,
        chanClose: make(chan int),
        chanSend:  make(chan *Broadcast),
        isOpen:  true,
    }

    // Register authenticated connections
    RegisterConn(w, r, ws)

    // Close when done
    defer func() {
        DeregisterConn(ws)
        ws.conn.Close()
    }()

    fmt.Printf("++++ WebSocket %p opened\n", ws.conn)

    // Close connection channel

    // Start incoming/outgoing listeners
    go procIncoming(w, r, ws)
    go procOutgoing(w, r, ws)

    // Wait for close message
    <-ws.chanClose
    ws.chanSend <- CloseBroadcastMessage
    ws.isOpen = false

    fmt.Printf("---X WebSocket %p closed\n", ws.conn)
}

// Processor for incoming (passive) messages.
func procIncoming(w http.ResponseWriter, r *http.Request, ws *Conn) {

    var msgType int
    var msg []byte
    var err error

    // Forever loop
    for {

        // Read message
        msgType, msg, err = ws.conn.ReadMessage()
        if err != nil {
            fmt.Printf("---Z WebSocket %p error: %v\n", ws.conn, err)
            break 
        }

        switch msgType {

        // Text message
        case websocket.TextMessage:
            processMessage(w, r, msg, ws)

        // Binary not supported, bye
        case websocket.BinaryMessage:
            fmt.Printf("---> WebSocket %p binary msg\n", ws.conn)

        case websocket.CloseMessage:
            fmt.Printf("---> WebSocket %p CLOSE msg\n", ws.conn)
            break

        case websocket.PingMessage:
            fmt.Printf("---> WebSocket %p ping msg\n", ws.conn)

        case websocket.PongMessage:
            fmt.Printf("---> WebSocket %p pong msg\n", ws.conn)

        default:
            fmt.Printf("---> WebSocket %p other msg: %v\n", ws.conn, msgType)
        }
    }

    // Signal connection closed
    ws.chanClose <- 1
}

// Processor for outgoing (active) messages.
func procOutgoing(w http.ResponseWriter, r *http.Request, ws *Conn) {

    var br *Broadcast

    // Forever loop
    for {

        // Wait on data channel
        br = <-ws.chanSend

        // Close message arrived ?
        if br.IsCloseMessage() {
            break
        }

        res := NewResponse(-1, br.Op, br.Data)
        Respond(ws.conn, res)
    }
}

// Processes incoming messages and invokes matching data processor.
func processMessage(w http.ResponseWriter, r *http.Request, msg []byte, ws *Conn) {

    // Parse request
    var req *Req
    var err error
    if req, err = ParseReq(msg); err != nil {
        res := NewResponse_Err(req.Id, req.Op, 
            errors.New_AppErr(err, "Cannot unmarshal request"))
        Respond(ws.conn, res)
        return
    }

    // Validate request
    if req.Id == 0 {
        res := NewResponse_Err(req.Id, req.Op, 
            errors.New_AppErr(err, "Request must have _reqId"))
        Respond(ws.conn, res)
        return
    }

    // HTTP request
    req.httpReq = r

    // Catch panic
    defer func() {
        if err := recover(); err != nil {
            stack := util.Stack()
            res := NewResponse_Err(req.Id, req.Op,
                errors.New_AppErr(fmt.Errorf("%v", err),
                "Application error, support notified"))
            Respond(ws.conn, res)

            // Report panic: err, url, params, stack
            _onPanic(
                fmt.Sprintf("Error processing WS request: %v", err),
                fmt.Sprintf("%v : %v : %v", r.Host, ws.router.URL, req.Op),
                fmt.Sprint(req.Params),
                stack)
        }
    }()

    // Find proc for given op
    var proc func(*Req) (interface{}, error)
    var ok bool

    // First look in core procs in this package
    // Then look in app supplied procs
    if proc, ok = _coreProcs[req.Op]; !ok {
        if proc, ok = ws.router.Procs[req.Op]; !ok {
            res := NewResponse_Err(req.Id, req.Op, 
                errors.New_NotFound(req.Op, "No matching op processor found"))
            Respond(ws.conn, res)
            return
        }
    }

    // Call proc
    var data interface{}
    data, err = proc(req)

    // Respond
    var res []byte
    if err == nil {
        res = NewResponse(req.Id, req.Op, data)
    } else {
        // Error can be either:
        // Request error: prepended with "ERR:" to be shown to user
        // Application error: all programming logic error
        res = NewResponse_Err(req.Id, req.Op, errors.New(err))
    }
    Respond(ws.conn, res)
}

