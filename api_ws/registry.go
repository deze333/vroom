package api_ws

import (
	"bytes"
	"fmt"

	"github.com/deze333/vroom/auth"
)

//------------------------------------------------------------
// Registry of active WebSocket connections
//------------------------------------------------------------

var _connsPublic = map[string]*Conn{}           // {agentId: conn}
var _connsAuthd = map[string]map[string]*Conn{} // {authId: {agentId: conn}}

//------------------------------------------------------------
// Functions
//------------------------------------------------------------

// Registers open connection based on session ID.
func RegisterConn(ws *Conn) {

	if ws.isAuthd {

		// Authenticated connection
		if id := auth.GetAuthdId(ws.r); id != "" {

			// Use authentication ID
			ws.authId = id

			// Try to use "email" from session values
			if vals, err := auth.GetSessionValues(ws.r); err == nil {
				ws.authEmail = fmt.Sprint(vals["email"])
			}

			// Add authd connection
			if conns, ok := _connsAuthd[ws.authId]; ok {
				conns[ws.agentId] = ws
			} else {
				_connsAuthd[ws.authId] = map[string]*Conn{ws.agentId: ws}
			}

			// Debug
			//fmt.Println(DumpConnsAuthd("REGISTERED CONNS: AUTHD"))

		} else {

			// Can't find ID for authenticated connection
			_onPanic(
				fmt.Sprintf("Cannot find ID for authd connection"),
				fmt.Sprintf("%v", ws.r))
		}

	} else {

		// Public connection
		_connsPublic[ws.agentId] = ws

		// Debug
		//fmt.Println(DumpConnsPublic("REGISTERED CONNS: PUBLIC"))
	}
}

// Deregisteres connection by removing it from the registry.
func DeregisterConn(ws *Conn) {
	if ws.isAuthd {
		// Authd connections
		if conns, ok := _connsAuthd[ws.authId]; ok {
			delete(conns, ws.agentId)
		}

		// Debug
		//fmt.Println(DumpConnsAuthd("REGISTERED CONNS: AUTHD"))

	} else {
		// Public connections
		delete(_connsPublic, ws.agentId)

		// Debug
		//fmt.Println(DumpConnsPublic("REGISTERED CONNS: PUBLIC"))
	}
}

// Finds connection that corresponds to agentId.
func GetConn(isAuthd bool, agentId string) (ws *Conn) {
	if isAuthd {
		// Authd connections
		for _, conns := range _connsAuthd {
			if ws, ok := conns[agentId]; ok {
				return ws
			}
		}
	} else {
		// Public connections
		ws = _connsPublic[agentId]
	}
	return
}

// Applies function fn to all public connections.
func applyToPublic(fn func(*Conn)) {
	for _, conn := range _connsPublic {
		fn(conn)
	}
}

// Applies function fn to all authd connections.
func applyToAuthd(fn func(*Conn)) {
	for _, conns := range _connsAuthd {
		for _, conn := range conns {
			fn(conn)
		}
	}
}

// Closes authenticated connection by ID.
func CloseAuthdConn(id string) {

	if conns, ok := _connsAuthd[id]; ok {

		// Multiple connections may share same authentication ID
		for _, ws := range conns {
			//fmt.Printf("__X Closed WebSocket ID = %v, conn = %p\n", id, ws.conn)
			ws.conn.Close()
		}

		DeregisterConn(&Conn{isAuthd: true, authId: id})
	}
}

//------------------------------------------------------------
// Monitoring functions
//------------------------------------------------------------

func ConnsInfoAuthd() (infos []map[string]interface{}) {

	infos = []map[string]interface{}{}

	// For each registered request connection
	for authId, conns := range _connsAuthd {

		// Add its opened WebSocket sessions
		for agentId, ws := range conns {

			info := map[string]interface{}{}
			info["_authId"] = authId
			info["_httpReqId"] = fmt.Sprintf("%v", agentId)

			sessVals, _ := auth.GetSessionValues(ws.r)
			for k, v := range sessVals {
				info["sess/"+k] = v
			}

			info["ws/authd"] = fmt.Sprint(ws.isAuthd)
			info["ws/open"] = fmt.Sprint(ws.isOpen)

			infos = append(infos, info)
		}
	}

	return
}

//------------------------------------------------------------
// Debug functions
//------------------------------------------------------------

func DumpConnsPublic(header string) string {
	var buf bytes.Buffer

	buf.WriteString("\n")
	buf.WriteString("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")
	buf.WriteString(header)
	buf.WriteString("\n")
	buf.WriteString("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")

	i := 0
	for agentId, ws := range _connsPublic {
		buf.WriteString(fmt.Sprintf("\t%v: agentId = %v, wsconn = %p\n", i, agentId, ws.conn))
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

	for authId, conns := range _connsAuthd {
		buf.WriteString(fmt.Sprintf("\t%v:\n", authId))
		j := 0
		for agentId, ws := range conns {
			buf.WriteString(fmt.Sprintf("\t\t%v: %v, agentId = %v, wsconn = %p\n", j, ws.authEmail, agentId, ws.conn))
			j++
		}
	}

	return buf.String()
}
