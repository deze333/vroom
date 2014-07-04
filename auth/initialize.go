package auth

import (
	"errors"
	"strconv"

	"github.com/gorilla/sessions"
)

//------------------------------------------------------------
// Authentication
//------------------------------------------------------------

var _cookieSessName string
var _cookieStoreId string
var _cookieStore *sessions.CookieStore
var _cookiePath string
var _cookieDomain string
var _cookieMaxAge int

var _onPanic func(string, string, ...interface{}) // msg, details, kvals

//------------------------------------------------------------
// Initialization
//------------------------------------------------------------

// Initializes authentication mechanism.
func Initialize(cookieSessName, cookieStoreId, cookiePath, cookieDomain, cookieMaxAge string, onPanic func(string, string, ...interface{})) (err error) {

	if cookieSessName == "" {
		err = errors.New("Auth cannot be configured with empty cookie session name")
		return
	}

	if cookieStoreId == "" {
		err = errors.New("Auth cannot be configured with empty cookie store id")
		return
	}

	_cookieMaxAge, err = strconv.Atoi(cookieMaxAge)
	if err != nil {
		err = errors.New("Auth cannot be configured with cookie maxAge not integer")
		return
	}

	_cookieSessName = cookieSessName
	_cookieStoreId = cookieStoreId
	_cookieStore = sessions.NewCookieStore([]byte(cookieStoreId))
	_cookiePath = cookiePath
	_cookieDomain = cookieDomain

	// Panic handler
	if onPanic != nil {
		_onPanic = onPanic
	} else {
		_onPanic = onPanicDisabled
	}
	return
}

// Dummy handler
func onPanicDisabled(err, details string, kvals ...interface{}) {
	// Do nothing
}
