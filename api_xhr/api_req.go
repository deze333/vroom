package api_xhr

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
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

	fmt.Println("=== DECODE REQUEST AS GZIP ===")
	fmt.Println(r)
	fmt.Println("=== Request body reader:", r.Body)

	var reader *gzip.Reader
	reader, err = gzip.NewReader(r.Body)
	if err != nil {
		fmt.Println("=== New READER for GZIP REQUEST FAILED:", err)
		return
	}

	// Decode
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&params)

	if err != nil && err != io.EOF {
		fmt.Println("=== Error decoding:", err)
	} else {
		fmt.Println("=== [OK] decoded body:", params)
	}

	// Empty request is not an error
	if err == io.EOF {
		err = nil
	}

	return
}
