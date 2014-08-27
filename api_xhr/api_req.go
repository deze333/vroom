package api_xhr

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/deze333/vroom/reqres"
)

//------------------------------------------------------------
// XHR Request
//------------------------------------------------------------

// Parses JSON request into package request.
func ParseReq(w http.ResponseWriter, r *http.Request) (req *reqres.Req, err error) {

	// Decode
	decoder := json.NewDecoder(r.Body)
	params := map[string]interface{}{}
	err = decoder.Decode(&params)

	// Empty request is not an error
	if err == io.EOF {
		err = nil
	}

	req = &reqres.Req{
		Params:        params,
		HttpReq:       r,
		HttpResWriter: w,
	}
	return
}
