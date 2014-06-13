package api_ws

import (
    "fmt"
    "bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
    "github.com/deze333/vroom/errors"
)

//------------------------------------------------------------
// Implementation
//------------------------------------------------------------

// Marshals successful data response.
func NewResponse(id int, op string, data interface{}) []byte {

    res := &Res{
        Id:   id,
        Op:   op,
        Data: data,
    }

    return marshal(res) //postEncode(marshal(res))
}

// Marshals error condition.
func NewResponse_Err(id int, op string, err *errors.ResError) []byte {

    res := &Res{
        Id:  id,
        Op:  op,
        Err: err,
    }

    return marshal(res) //postEncode(marshal(res))
}

// Marshals response to JSON.
func marshal(res *Res) []byte {

    jsonb, err := json.Marshal(res)
    if err == nil {
        return jsonb
    }

    // Report panic: err, url, params, stack
    _onPanic(
        fmt.Sprintf("Error marshalling WS response, error: %v", err),
        fmt.Sprintf("WebSocket JSON encoding"),
        fmt.Sprint(res),
        "Stack not needed")

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
    return bytes.Replace(res, []byte{'%'}, []byte{'%','%'}, -1)
}

// Responds on websocket connection.
func Respond(conn *websocket.Conn, res []byte) {

    err := conn.WriteMessage(websocket.TextMessage, res)
    if err == nil {
        return
    }

    // Report panic: err, url, params, stack
    _onPanic(
        fmt.Sprintf("WebSocket failed to write response, error: %v", err),
        fmt.Sprintf("%v ---> %v", conn.LocalAddr(), conn.RemoteAddr()),
        string(res),
        "Stack not needed")
}

