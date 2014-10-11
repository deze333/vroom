package api_ws

import (
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

//------------------------------------------------------------
// WebSocket messages async processor
//------------------------------------------------------------

const (
	CHAN_MESSAGE_MAX           = 1024
	CHAN_MESSAGE_BROADCAST_MAX = 100
	CHAN_MESSAGE_FAILED_MAX    = 100
	LOST_MESSAGE_MAX_LIFE      = 30 * time.Second
)

var (
	_chanIn          = make(chan *Message, CHAN_MESSAGE_MAX)
	_chanOut         = make(chan *Message, CHAN_MESSAGE_MAX)
	_chanBroadcast   = make(chan *Message, CHAN_MESSAGE_BROADCAST_MAX)
	_chanFailed      = make(chan *Message, CHAN_MESSAGE_FAILED_MAX)
	_failed          = map[string][]*Message{}
	_procFailedMutex = &sync.Mutex{}
)

// Initializes async processors.
func init() {
	go procChanIn()
	go procChanOut()
	go procChanBroadcast()
	go procChanFailed()
	go procFailedCleanup()
}

// Processes incoming channel messages. One message at a time.
func procChanIn() {
	var conn *Conn
	for msg := range _chanIn {
		msg.isProcessed = false
		fmt.Printf("chan IN ---> %v = %v\n", msg.agentId, string(msg.req))

		// Attempt to get active connection
		if conn = GetConn(msg.isAuthd, msg.agentId); conn == nil {
			_chanFailed <- msg
			continue
		}

		// Process message
		// Must be done via goroutine to prevent channel block
		// due to processing delay
		go processAndRespond(conn, msg, _chanOut)
	}
}

// Processes outgoing channel messages. One message at a time.
func procChanOut() {
	var conn *Conn
	for msg := range _chanOut {
		fmt.Printf("chan OUT <--- %v\n", string(msg.res))

		// Attempt to get active connection
		if conn = GetConn(msg.isAuthd, msg.agentId); conn == nil {
			_chanFailed <- msg
			continue
		}

		// Write response
		err := conn.conn.WriteMessage(websocket.TextMessage, msg.res)
		if err != nil {
			// Close connection on write error
			_chanFailed <- msg
			DeregisterConn(conn)
			conn.conn.Close()
			continue
		}
	}
}

// Sends broadcast messages to all active connections.
func procChanBroadcast() {
	for msg := range _chanBroadcast {
		fmt.Printf("chan BROADCAST <--- %v\n", string(msg.res))

		// Apply broadcaster function to each connection
		if msg.isAuthd {
			applyToAuthd(
				func(conn *Conn) {
					broadcaster(conn, msg)
				})
		} else {
			applyToPublic(
				func(conn *Conn) {
					broadcaster(conn, msg)
				})
		}
	}
}

// Broadcasts message to single connection.
// Closes connection on error.
func broadcaster(conn *Conn, msg *Message) {
	err := conn.conn.WriteMessage(websocket.TextMessage, msg.res)
	if err != nil {
		// Close connection on write error
		_chanFailed <- msg
		DeregisterConn(conn)
		conn.conn.Close()
	}
}

// Listens on failed channel and adds messages to failed.
func procChanFailed() {
	for msg := range _chanFailed {
		_procFailedMutex.Lock()

		t := time.Now()
		msg.failTime = &t
		if agentFails, ok := _failed[msg.agentId]; ok {
			_failed[msg.agentId] = append(agentFails, msg)
		} else {
			_failed[msg.agentId] = []*Message{msg}
		}

		_procFailedMutex.Unlock()
	}
}

// Periodically cleans up expired failed messages.
func procFailedCleanup() {

	ticker := time.NewTicker(15 * time.Second)

	for {
		<-ticker.C
		_procFailedMutex.Lock()

		now := time.Now()

		for agentId, agentFails := range _failed {

			if anyExpiredFails(agentFails, now) {
				unexp := getUnexpiredFails(agentFails, now)
				if len(unexp) > 0 {
					_failed[agentId] = unexp
				} else {
					delete(_failed, agentId)
				}
			}
		}

		_procFailedMutex.Unlock()
	}
}

// Has at least one message expired?
func anyExpiredFails(msgs []*Message, t time.Time) bool {
	for _, msg := range msgs {
		if msg.failTime == nil {
			return true // safe guard, nil not allowed
		}
		if msg.failTime.Add(LOST_MESSAGE_MAX_LIFE).After(t) {
			return true
		}
	}
	return false
}

// Filters out expired messages.
func getUnexpiredFails(msgs []*Message, t time.Time) (unexp []*Message) {
	for _, msg := range msgs {
		if msg.failTime == nil {
			continue // safe guard, nil not allowed
		}
		if msg.failTime.Add(LOST_MESSAGE_MAX_LIFE).After(t) {
			continue
		}
		unexp = append(unexp, msg)
	}
	return
}

// Pushes specific agent's failed messages back to channel for processing.
func retryFailedMessages(agentId string) {

	_procFailedMutex.Lock()

	if agentFails, ok := _failed[agentId]; ok {
		delete(_failed, agentId)
		for _, msg := range agentFails {
			msg.failTime = nil
			if msg.isProcessed {
				_chanOut <- msg
			} else {
				_chanIn <- msg
			}
		}
	}

	_procFailedMutex.Unlock()
}
