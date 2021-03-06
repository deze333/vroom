package api_xhr

import (
	"fmt"
	"net/http"

	"github.com/deze333/vroom/errors"
	"github.com/deze333/vroom/reqres"
	"github.com/deze333/vroom/util"
)

//------------------------------------------------------------
// XHR Implementation
//------------------------------------------------------------

// Handles XHR Not Authd response.
func Handle_NotAuthd(w http.ResponseWriter, r *http.Request) {

	// Error 401, Unathorized
	w.WriteHeader(http.StatusUnauthorized)

	// Dummy wrapper for consistency
	req := &reqres.Req{
		HttpReq:       r,
		HttpResWriter: w,
	}

	res := NewResponse_Err(req, errors.New_NotAuthd())
	w.Write(res)
}

// Handles XHR request.
func Handle(w http.ResponseWriter, r *http.Request, fn func(*reqres.Req) (interface{}, error)) {

	w.Header().Set("Content-Type", "application/json")

	var req *reqres.Req

	// Check if there is a file in a request.
	// If so - do not parse request as JSON.
	_, _, err := r.FormFile("file")
	if err != nil {
		// Parse request data
		req, err = ParseReq(w, r)

		// On parse error
		if err != nil {
			res := NewResponse_Err(req, errors.New_AppErr(err,
				"Error parsing request data as JSON"))
			w.Write(res)
			return
		}
	} else {
		// Dummy wrapper for consistency.
		req = &reqres.Req{
			HttpReq:       r,
			HttpResWriter: w,
		}
	}

	// Catch panic
	defer func() {
		if err := recover(); err != nil {
			res := NewResponse_Err(req, errors.New_AppErrWithStack(fmt.Errorf("%v", err),
				"Application error, support notified",
				util.Stack()))
			w.Write(res)

			// Report panic: err, url, params, stack
			_onPanic(
				fmt.Sprintf("Error processing XHR request: %v", err),
				fmt.Sprintf("%v : %v", req.HttpReq.Host, req.HttpReq.RequestURI),
				"Session", fmt.Sprint(req.GetSessionValues()),
				"Params", fmt.Sprint(req.Params),
				"Stack", util.Stack())
		}
	}()

	// Call processor
	data, err := fn(req)

	// On processor API error
	if err != nil {
		// Error can be either:
		// Request error: prepended with "ERR:" to be shown to user
		// Application error: all programming logic error
		res := NewResponse_Err(req, errors.New(err))
		w.Write(res)
		return
	}

	// Successful response
	res := NewResponse(req, data)
	w.Write(res)
}
