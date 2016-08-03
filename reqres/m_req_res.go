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

	// "" (empty string) - default
	// "DataOnly"
	ResponseMarshalingMode string

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
func (r *Req) GetSessionValues() map[string]interface{} {
	vals, err := auth.GetSessionValues(r.HttpReq)
	if err != nil {
		vals = map[string]interface{}{}
	}
	return vals
}

// Returns specific session value of underlying HTTP request.
// Returns empty string if value doesn't exist.
func (r *Req) GetSessionValue(key string) interface{} {
	vals, err := auth.GetSessionValues(r.HttpReq)
	if err != nil {
		return nil
	}

	if val, ok := vals[key]; ok {
		return val
	}

	return nil
}

//------------------------------------------------------------
// Response methods
//------------------------------------------------------------

// Marshals response to JSON.
// Includes whole response object.
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

// Marshals response data field to JSON.
// Use for external integrations that don't need our response data structure.
func (r *Res) MarshalDataOnly(req *Req) []byte {

	fmt.Println("Marshal data only =====", r.Data)

	if r.Data == nil {
		return []byte{}
	}

	// Check if that's an empty map to prevent
	// return that will say "null"
	if data, ok := r.Data.(map[string]interface{}); ok {
		if len(data) == 0 {
			return []byte{}
		}
	}

	// Marshal data only
	jsonb, err := json.Marshal(r.Data)
	if err == nil {
		//return postEncode(jsonb)
		fmt.Println("Marshal data only =====", jsonb)
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

/*
func marshalDataOnly(r *Res) (jsonb{

	data, isCorrectType := r.data.(map[string]interface{})
	if !isCorrectType {
		return false
	}
}
*/

// Makes some post encoding adjustements to achieve correct JSON.
func postEncode(res []byte) []byte {

	// XXX Perhaps I need to read the manual...
	// Fix of strange behaviour when writer expects second %
	// after first and otherwise says (MISSING), that breaks JSON parser.
	// SOLUTION:
	// Convert % into %%.
	return bytes.Replace(res, []byte{'%'}, []byte{'%', '%'}, -1)
}
