package api_ws

import (
	"errors"
)

//------------------------------------------------------------
// Globals
//------------------------------------------------------------

var _verGetter func() string
var _onPanic func(string, string, ...interface{}) // msg, details, kvals

//------------------------------------------------------------
// Initialization
//------------------------------------------------------------

// Initializes websocket package.
// OnPanic handler must be provided.
func Initialize(verGetter func() string, onPanic func(string, string, ...interface{})) (err error) {

	if verGetter == nil {
		err = errors.New("Websocket handler needs Version Getter handler")
		return
	}

	if onPanic == nil {
		err = errors.New("Websocket handler needs OnPanic handler")
		return
	}

	_verGetter = verGetter
	_onPanic = onPanic
	return
}
