package api_ws

import (
	"fmt"

	"github.com/deze333/vroom/errors"
	"github.com/deze333/vroom/reqres"
	"github.com/deze333/vroom/util"
)

//------------------------------------------------------------
// WebSocket Message processor
//------------------------------------------------------------

// Processes incoming messages and invokes matching data processor.
func processMessage(ws *Conn, msg []byte) {

	r := ws.r

	// Parse request
	var req *reqres.Req
	var err error
	if req, err = ParseReq(msg); err != nil {
		res := NewResponse_Err(req,
			errors.New_AppErr(err, "Cannot unmarshal request"))
		ws.chanOut <- res
		return
	}

	/* NOTE: Relax this requirement
	   // Validate request
	   if req.Id == 0 {
	       err = fmt.Errorf("Request ID is missing")
	       res := NewResponse_Err(req,
	           errors.New_AppErr(err, "Request must have _reqId"))
	       ws.chanOut <- res
	       return
	   }
	*/

	// HTTP request
	req.HttpReq = r

	// Catch panic
	defer func() {
		if err := recover(); err != nil {
			stack := util.Stack()
			res := NewResponse_Err(req,
				errors.New_AppErr(fmt.Errorf("%v", err),
					"Application error, support notified"))
			ws.chanOut <- res

			// Report panic: err, url, params, stack
			_onPanic(
				fmt.Sprintf("Error processing WS request: %v", err),
				fmt.Sprintf("%v : %v : %v", r.Host, ws.router.URL, req.Op),
				"Params", fmt.Sprint(req.Params),
				"Session", fmt.Sprint(req.GetSessionValues()),
				"Stack", stack)
		}
	}()

	// Find proc for given op
	var proc func(*reqres.Req) (interface{}, error)
	var ok bool

	// First look in core procs in this package
	// Then look in app supplied procs
	if proc, ok = _coreProcs[req.Op]; !ok {
		if proc, ok = ws.router.Procs[req.Op]; !ok {
			res := NewResponse_Err(req,
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
		res = NewResponse(req, data)
	} else {
		// Error can be either:
		// Request error: prepended with "ERR:" to be shown to user
		// Application error: all programming logic error
		res = NewResponse_Err(req, errors.New(err))
	}
	ws.chanOut <- res
}

// Processes incoming connection messages and writes resposnse to chanOut.
func processAndRespond(conn *Conn, msg *Message, chanOut chan *Message) {
	processConnectionMessage(conn, msg)
	msg.isProcessed = true
	chanOut <- msg
}

// Processes incoming messages and invokes matching data processor.
func processConnectionMessage(conn *Conn, msg *Message) {

	r := conn.r

	// Parse request
	var req *reqres.Req
	var err error
	if req, err = ParseReq(msg.req); err != nil {
		msg.res = NewResponse_Err(req,
			errors.New_AppErr(err, "Cannot unmarshal request"))
		return
	}

	// Add originating HTTP request
	req.HttpReq = r

	// Catch if panics
	defer func() {
		if err := recover(); err != nil {
			stack := util.Stack()
			msg.res = NewResponse_Err(req,
				errors.New_AppErr(fmt.Errorf("%v", err),
					"Application error, support notified"))

			// Report panic: err, url, params, stack
			_onPanic(
				fmt.Sprintf("Error processing WebSocket request: %v", err),
				fmt.Sprintf("%v : %v : %v", r.Host, conn.router.URL, req.Op),
				"Params", fmt.Sprint(req.Params),
				"Session", fmt.Sprint(req.GetSessionValues()),
				"Stack", stack)
		}
	}()

	// Proc for given op
	var proc func(*reqres.Req) (interface{}, error)
	var ok bool

	// First look in core route procs in this package
	// Then look in app supplied procs
	if proc, ok = _coreProcs[req.Op]; !ok {
		if proc, ok = conn.router.Procs[req.Op]; !ok {
			msg.res = NewResponse_Err(req,
				errors.New_NotFound(req.Op, "No matching op processor found"))
			return
		}
	}

	// Call proc
	var data interface{}
	data, err = proc(req)

	// Respond
	if err == nil {
		msg.res = NewResponse(req, data)
	} else {
		// Error can be either:
		// Request error: prepended with "ERR:" to be shown to user
		// Application error: all programming logic error
		msg.res = NewResponse_Err(req, errors.New(err))
	}
}
