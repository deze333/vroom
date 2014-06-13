package api_ws

import (
    "net/http"
    "github.com/deze333/vroom/auth"
)

//------------------------------------------------------------
// Websocket Request
//------------------------------------------------------------

type Req struct {
	Id            int                     `json:"_id"`
	Op            string                  `json:"op,omitempty"`
	Params        map[string]interface{}  `json:"params,omitempty"`
	Data          interface{}             `json:"data,omitempty"`

    httpReq       *http.Request
}

// Returns underlying HTTP request.
func (r *Req) GetHttpRequest() *http.Request {
    return r.httpReq
}

// Returns session values of underlying HTTP request.
func (r *Req) GetSessionValues() map[string]string {
    vals, err := auth.GetSessionValues(r.httpReq)
    if err != nil {
        vals = map[string]string{}
    }
    return vals
}

