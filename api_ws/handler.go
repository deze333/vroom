package api_ws

import (
	"fmt"
	"net/http"
	"time"

	"github.com/deze333/vroom/auth"
	"github.com/deze333/vroom/reqres"
	"github.com/gorilla/websocket"
)

//------------------------------------------------------------
// WebSocket Implementation
//------------------------------------------------------------

// Handles WS Not Authd response.
func Handle_NotAuthd(w http.ResponseWriter, r *http.Request) {

	http.Error(w, "Not authorized", http.StatusUnauthorized)
}

// Creates WS connecttion and processes it in forever loop.
func Handle(w http.ResponseWriter, r *http.Request, router *reqres.WebSocket_Router, isAuthd bool) {

	var agentId string

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			r.ParseForm()

			fmt.Printf("\n\n-------------------------------\nCheck WebSocket call: \nURL = %v \nOrigin = %v\nCaller ID = %v\n---------------------------------------\n\n", r.URL, r.Header["Origin"], r.Form["id"])

			// Check that ID is present
			if ids, ok := r.Form["id"]; ok {
				if len(ids) > 0 {
					agentId = ids[0]
				} else {
					return false
				}
			} else {
				return false
			}

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
			"Session", fmt.Sprint(auth.GetSessionValues(r)))
		return
	}

	// Create connection object
	ws := &Conn{
		w:                   w,
		r:                   r,
		conn:                conn,
		router:              router,
		agentId:             agentId,
		isAuthd:             isAuthd,
		chanIn:              make(chan []byte),
		chanOut:             make(chan []byte),
		chanProcClose:       make(chan int),
		chanProcWriterClose: make(chan int),
		isOpen:              true,
	}

	// Register this new connection
	RegisterConn(ws)

	// Retry failed messages, both in & out
	retryFailedMessages(agentId)

	// Close when done
	defer func() {
		DeregisterConn(ws)
		ws.conn.Close()
	}()

	//fmt.Printf("++++ WebSocket %p opened\n", ws.conn)

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

			//fmt.Println("O__________")
			// XXX TEMP TESTING FOR NEWER PROC
			// Choose processing logic depending on user type

			// XXX Force via _chanIn:
			//vals, _ := auth.GetSessionValues(r)
			//if q, ok := vals["qualities"]; ok {
			if true {
				//if strings.Contains(q, "tester") {
				if true {
					_chanIn <- &Message{isAuthd: isAuthd, agentId: agentId, req: msg}
				} else {
					ws.chanIn <- msg
				}
			}
			//fmt.Println("__________X")

		// Binary not supported
		case websocket.BinaryMessage:
			// Do nothing

		case websocket.CloseMessage:
			break

		case websocket.PingMessage:
			// Do nothing

		case websocket.PongMessage:
			// Do nothing

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

		case msg := <-ws.chanIn:
			processMessage(ws, msg)

		case <-ws.chanProcClose:
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

		case res := <-ws.chanOut:
			err := Respond(ws, res)

			// Exit on write error
			if err != nil {
				DeregisterConn(ws)
				ws.conn.Close()
				return
			}

		case <-ws.chanProcWriterClose:
			// Exit
			fmt.Println("---X_W procWriter finished")
			return
		}
	}
}
