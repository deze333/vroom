package reqres

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/deze333/vroom/auth"
	"github.com/deze333/vroom/errors"
)

//------------------------------------------------------------
// Request
//------------------------------------------------------------

type Req struct {
	Id     int                    `json:"_id"`
	Op     string                 `json:"op,omitempty"`
	Params map[string]interface{} `json:"params,omitempty"`

	HttpReq       *http.Request
	HttpResWriter http.ResponseWriter
}

//------------------------------------------------------------
// Response
//------------------------------------------------------------

type Res struct {
	Id   int         `json:"_id"`
	Op   string      `json:"op,omitempty"`
	Err  interface{} `json:"_err,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

//------------------------------------------------------------
// Request methods
//------------------------------------------------------------

// Returns session values of underlying HTTP request.
func (r *Req) GetSessionValues() map[string]string {
	vals, err := auth.GetSessionValues(r.HttpReq)
	if err != nil {
		vals = map[string]string{}
	}
	return vals
}

//------------------------------------------------------------
// Response methods
//------------------------------------------------------------

// Marshals response to JSON.
func (r *Res) Marshal(req *Req) []byte {

	jsonb, err := json.Marshal(r)
	if err == nil {
		//return postEncode(jsonb)
		return jsonb
	}

	// Report marshalling error
	if req != nil {
		_onPanic(
			fmt.Sprintf("Error marshalling response: %v", err),
			fmt.Sprintf("%v : %v", req.HttpReq.Host, req.HttpReq.RequestURI),
			"Session", fmt.Sprint(req.GetSessionValues(),
				"Data", r.Data))
	} else {
		_onPanic(
			fmt.Sprintf("Error marshalling response: %v", err),
			"",
			"Data", r.Data)
	}

	// Error marshalling response
	resErr, _ := json.Marshal(
		errors.New_AppErr(err, "Cannot marshal JSON response"))

	return []byte(fmt.Sprintf(
		`{"_id": %v, "op": "%v", "_err": %v}`,
		r.Id, r.Op, string(resErr)))
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
