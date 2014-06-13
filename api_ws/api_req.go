package api_ws

import (
    "bytes"
    "io"
	"encoding/json"
)

//------------------------------------------------------------
// Request
//------------------------------------------------------------

// Parses JSON request into a map. May return ApiErr if parse failed.
func ParseReq(raw []byte) (req *Req, err error) {

    // Decode
    decoder := json.NewDecoder(bytes.NewBuffer(raw))
    req = &Req{}
    err = decoder.Decode(&req)

    // Success ?
    if err == nil {
        return
    }

    // Empty is not an error
    if err == io.EOF {
        err = nil
        return
    } 

    return
}

