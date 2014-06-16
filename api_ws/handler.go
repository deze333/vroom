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
        w: w,
        r: r,
        conn: conn,
        router: router,
        isAuthd: isAuthd,
        chanIn:  make(chan []byte),
        chanOut:  make(chan []byte),
        chanProcClose: make(chan int),
        chanProcWriterClose: make(chan int),
        isOpen:  true,
    }

    // Register authenticated connections
    RegisterConn(ws)

    // Close when done
    defer func() {
        DeregisterConn(ws)
        ws.conn.Close()
    }()

    fmt.Printf("++++ WebSocket %p opened\n", ws.conn)

    // Start async processors
    go proc(ws)
    go procWriter(ws)

    var msgType int
    var msg []byte

    // Forever loop listening for incoming messages
    for {

        // Read message
        msgType, msg, err = ws.conn.ReadMessage()
        if err != nil {
            break 
        }

        switch msgType {

        // Text message
        case websocket.TextMessage:
            ws.chanIn <- msg

        // Binary not supported, bye
        case websocket.BinaryMessage:

        case websocket.CloseMessage:
            break

        case websocket.PingMessage:

        case websocket.PongMessage:

        default:
        }
    }

    // Signal procs to finish
    ws.chanProcClose <- 1
    ws.chanProcWriterClose <- 1
    fmt.Printf("---X WebSocket %p closed\n", ws.conn)
}

// Processes incoming channel data.
func proc(ws *Conn) {
    for {
        select {

        case msg := <- ws.chanIn:
            processMessage(ws, msg)

        case <- ws.chanProcClose:
            // Exit
            fmt.Println("---X_R proc finished")
            return
        }
    }
}

// Processes outgoing channel responses.
func procWriter(ws *Conn) {
    for {
        select {

        case res := <- ws.chanOut:
            Respond(ws.conn, res)

        case <- ws.chanProcWriterClose:
            // Exit
            fmt.Println("---X_W procWriter finished")
            return
        }
    }
}

// Processes incoming messages and invokes matching data processor.
func processMessage(ws *Conn, msg []byte) {

    r := ws.r

    // Parse request
    var req *Req
    var err error
    if req, err = ParseReq(msg); err != nil {
        res := NewResponse_Err(req.Id, req.Op, 
            errors.New_AppErr(err, "Cannot unmarshal request"))
        ws.chanOut <- res
        return
    }

    /* NOTE: Relax this requirement
    // Validate request
    if req.Id == 0 {
        err = fmt.Errorf("Request ID is missing")
        res := NewResponse_Err(req.Id, req.Op, 
            errors.New_AppErr(err, "Request must have _reqId"))
        ws.chanOut <- res
        return
    }
    */

    // HTTP request
    req.httpReq = r

    // Catch panic
    defer func() {
        if err := recover(); err != nil {
            stack := util.Stack()
            res := NewResponse_Err(req.Id, req.Op,
                errors.New_AppErr(fmt.Errorf("%v", err),
                "Application error, support notified"))
            ws.chanOut <- res

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
            ws.chanOut <- res
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
    ws.chanOut <- res
}

