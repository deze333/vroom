package api_ws

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/deze333/vroom/errors"
	"github.com/deze333/vroom/reqres"
	"github.com/gorilla/websocket"
)

//------------------------------------------------------------
// Implementation
//------------------------------------------------------------

// Marshals successful data response.
func NewResponse(req *reqres.Req, data interface{}) []byte {

	res := &reqres.Res{
		Id:   req.Id,
		Op:   req.Op,
		Data: data,
	}

	return marshal(req, res) //postEncode(marshal(res))
}

// Marshals broadcast response.
func NewResponse_Broadcast(id int, op string, data interface{}) []byte {

	res := &reqres.Res{
		Id:   id,
		Op:   op,
		Data: data,
	}

	return marshal(nil, res) //postEncode(marshal(res))
}

// Marshals error condition.
func NewResponse_Err(req *reqres.Req, err *errors.ResError) []byte {

	res := &reqres.Res{
		Id:  req.Id,
		Op:  req.Op,
		Err: err,
	}

	return marshal(req, res) //postEncode(marshal(res))
}

// Marshals response to JSON.
func marshal(req *reqres.Req, res *reqres.Res) []byte {

	jsonb, err := json.Marshal(res)
	if err == nil {
		return jsonb
	}

	// Report marshalling error
	if req != nil {
		_onPanic(
			fmt.Sprintf("Error marshalling WS response: %v", err),
			fmt.Sprintf("%v : %v", req.HttpReq.Host, req.HttpReq.RequestURI),
			"Session", fmt.Sprint(req.GetSessionValues(),
				"Data", res.Data))
	} else {
		_onPanic(
			fmt.Sprintf("Error marshalling WS broadcast: %v", err),
			"",
			"Data", res.Data)
	}

	// Error marshalling response
	resErr, _ := json.Marshal(
		errors.New_AppErr(err, "Cannot marshal JSON response"))

	return []byte(fmt.Sprintf(
		`{"_id": %v, "op": "%v", "_err": %v}`,
		res.Id, res.Op, string(resErr)))
}

// Makes some post encoding adjustements to achieve correct JSON.
func postEncode(res []byte) []byte {

	// XXX Perhaps I need to read the manual...
	// Fix of strange behaviour when writer expects second %
	// after first and otherwise says (MISSING), that breaks JSON parser.
	// SOLUTION:
	// Convert % into %%.
	return bytes.Replace(res, []byte{'%'}, []byte{'%', '%'}, -1)
}

// Responds on websocket connection.
func Respond(conn *Conn, res []byte) (err error) {

	wsConn := conn.conn
	err = wsConn.WriteMessage(websocket.TextMessage, res)
	if err == nil {
		return
	}

	// TODO Include reporting in admin
	/*
	   // Get session details
	   sess, _ := auth.GetSessionValues(conn.r)

	   // Report panic: err, url, params, session, stack
	   _onPanic(
	       fmt.Sprintf("WebSocket failed to write response, error: %v", err),
	       fmt.Sprintf("%v #%v @ %v", sess["initials"], sess["_auth"], sess["_ip"]),
	       "Response", string(res),
	       "Session", fmt.Sprint(sess))
	*/

	return
}
