package api_xhr

import (
    "errors"
)

//------------------------------------------------------------
// Globals
//------------------------------------------------------------

var _onPanic func(string, string, string, string, string) // err, url, params, session, stack

//------------------------------------------------------------
// Initialization
//------------------------------------------------------------

// Initializes XHR package.
// OnPanic handler must be provided.
func Initialize(onPanic func(string, string, string, string, string)) (err error) {
    
    if onPanic == nil {
        err = errors.New("XHR handler needs OnPanic handler")
        return
    }

    _onPanic = onPanic
    return
}

