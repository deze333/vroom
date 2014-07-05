package reqres

import (
	"net/http"

	"github.com/deze333/vroom/auth"
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
