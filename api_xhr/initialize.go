package api_xhr

import (
	"errors"
)

//------------------------------------------------------------
// Globals
//------------------------------------------------------------

var _onPanic func(string, string, ...interface{}) // msg, details, kvals

//------------------------------------------------------------
// Initialization
//------------------------------------------------------------

// Initializes XHR package.
// OnPanic handler must be provided.
func Initialize(onPanic func(string, string, ...interface{})) (err error) {

	if onPanic == nil {
		err = errors.New("XHR handler needs OnPanic handler")
		return
	}

	_onPanic = onPanic
	return
}
