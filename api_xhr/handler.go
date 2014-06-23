package api_xhr

import (
    "fmt"
    "net/http"
    "github.com/deze333/vroom/errors"
    "github.com/deze333/vroom/auth"
    "github.com/deze333/vroom/util"
)

//------------------------------------------------------------
// XHR Implementation
//------------------------------------------------------------

// Handles XHR Not Authd response.
func Handle_NotAuthd(w http.ResponseWriter, r *http.Request) {

    // Error 401, Unathorized
    w.WriteHeader(http.StatusUnauthorized)

    res := NewResponse_Err(errors.New_NotAuthd())
    w.Write(res)
}

// Handles XHR request.
func Handle(w http.ResponseWriter, r *http.Request, fn func(map[string]interface{})(interface{}, error)) {

    w.Header().Set("Content-Type", "application/json")

    // Parse request data
    req, err := ParseReq(r)

    // On parse error
    if err != nil {
        res := NewResponse_Err(errors.New_AppErr(err,
            "Error parsing request data as JSON"))
        w.Write(res)
        return
    }

    // Catch panic
    defer func() {
        if err := recover(); err != nil {
            res := NewResponse_Err(errors.New_AppErr(fmt.Errorf("%v", err),
                "Application error, support notified"))
            w.Write(res)

            // Report panic: err, url, params, stack
            _onPanic(
                fmt.Sprintf("Error processing XHR request: %v", err),
                fmt.Sprintf("%v : %v", r.Host, r.RequestURI),
                fmt.Sprint(auth.GetSessionValues(r)),
                fmt.Sprint(req),
                util.Stack())
        }
    }()

    // Add useful request data
    req["_httpReq"] = r
    req["_httpResWriter"] = w
    req["_session"], _ = auth.GetSessionValues(r)

    // Call processor
    data, err := fn(req)

    // On processor API error
    if err != nil {
        // Error can be either:
        // Request error: prepended with "ERR:" to be shown to user
        // Application error: all programming logic error
        res := NewResponse_Err(errors.New(err))
        w.Write(res)
        return
    }

    // Successful response
    res := NewResponse(data)
    w.Write(res)
}
