package api_ws

import (
)

//------------------------------------------------------------
// WebSocket broadcast
//------------------------------------------------------------

// Returns public broadcaster.
func GetBroadcaster_Public() func(string, interface{}) (err error) {
	return broadcastPublic
}

// Returns authenticated broadcaster.
func GetBroadcaster_Authd() func(string, interface{}) (err error) {
	return broadcastAuthd
}

// Broadcasts to open public websocket connections.
func broadcastPublic(op string, data interface{}) (err error) {

	//fmt.Println(DumpConnsPublic("BROADCAST PUBLIC"))

	/* OLD STYLE:
	for _, ws := range _connsPublic {
		if ws.isOpen {
			ws.chanOut <- NewResponse_Broadcast(0, op, data)
		}
	}
	*/

	_chanBroadcast <- &Message{
		isAuthd: false,
		res:     NewResponse_Broadcast(0, op, data),
	}
	return
}

// Broadcasts to open authenticated websocket connections.
func broadcastAuthd(op string, data interface{}) (err error) {

	//fmt.Println(DumpConnsAuthd("BROADCAST AUTHD"))

	/* OLD STYLE:
	for _, conns := range _connsAuthd {
		for _, ws := range conns {
			if ws.isOpen {
				ws.chanOut <- NewResponse_Broadcast(0, op, data)
			}
		}
	}
	*/

	// XXX via channel
	_chanBroadcast <- &Message{
		isAuthd: true,
		res:     NewResponse_Broadcast(0, op, data),
	}

	return
}
