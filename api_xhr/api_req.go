package api_xhr

import (
    "io"
	"encoding/json"
    "net/http"
)

//------------------------------------------------------------
// XHR Request
//------------------------------------------------------------

// Parses JSON request into a map. 
func ParseReq(r *http.Request) (req map[string]interface{}, err error) {

    // Decode
    decoder := json.NewDecoder(r.Body)
    req = map[string]interface{}{}
    err = decoder.Decode(&req)
    if err == nil {
        return
    }

    // Empty request is not an error
    if err == io.EOF {
        err = nil
    } 

    return
}

