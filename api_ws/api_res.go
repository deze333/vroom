package api_ws

import (
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

	return res.Marshal(req)
}

// Marshals broadcast response.
func NewResponse_Broadcast(id int, op string, data interface{}) []byte {

	res := &reqres.Res{
		Id:   id,
		Op:   op,
		Data: data,
	}

	return res.Marshal(nil)
}

// Marshals error condition.
func NewResponse_Err(req *reqres.Req, err *errors.ResError) []byte {

	res := &reqres.Res{
		Id:  req.Id,
		Op:  req.Op,
		Err: err,
	}

	return res.Marshal(req)
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
