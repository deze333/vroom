package api_xhr

import (
	"compress/gzip"
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

	params := map[string]interface{}{}

	// Support different content ecnodings
	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		params, err = decodeAsGzipReq(r)

	default:
		params, err = decodeAsUnencodedReq(r)
	}

	req = &reqres.Req{
		Params:        params,
		HttpReq:       r,
		HttpResWriter: w,
	}
	return
}

func decodeAsUnencodedReq(r *http.Request) (params map[string]interface{}, err error) {

	// Decode
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)

	// Empty request is not an error
	if err == io.EOF {
		err = nil
	}

	return
}

// Parses JSON payload in GZIP compressed request.
func decodeAsGzipReq(r *http.Request) (params map[string]interface{}, err error) {

	var reader *gzip.Reader
	reader, err = gzip.NewReader(r.Body)
	if err != nil {
		return
	}

	// Decode
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&params)

	// Empty request is not an error
	if err == io.EOF {
		err = nil
	}

	return
}
