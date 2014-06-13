package api_ws

import (
    "errors"
)

//------------------------------------------------------------
// Globals
//------------------------------------------------------------

var _onPanic func(string, string, string, string) // err, url, params, stack

//------------------------------------------------------------
// Initialization
//------------------------------------------------------------

// Initializes websocket package.
// OnPanic handler must be provided.
func Initialize(onPanic func(string, string, string, string)) (err error) {
    
    if onPanic == nil {
        err = errors.New("Websocket handler needs OnPanic handler")
        return
    }

    _onPanic = onPanic
    return
}

