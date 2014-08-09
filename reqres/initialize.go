package reqres

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

// Initializes package.
// OnPanic handler must be provided.
func Initialize(onPanic func(string, string, ...interface{})) (err error) {

	if onPanic == nil {
		err = errors.New("ReqRes package needs OnPanic handler")
		return
	}

	_onPanic = onPanic
	return
}
