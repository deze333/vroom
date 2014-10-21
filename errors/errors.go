package errors

import (
	"strings"
)

//------------------------------------------------------------
// Response error
//------------------------------------------------------------

const (
	// Prefix for errors
	ERR_PREFIX_BAD_DATA = "(BAD_DATA)" // user supplied bad data

	// Values passed to browser
	ERR_NOT_FOUND = "NOT_FOUND" // requested URL not found
	ERR_NOT_AUTHD = "NOT_AUTHD" // user not authenticated for this request
	ERR_APP       = "APP_ERR"   // error on app side (can't open file, etc.)
	ERR_REQ       = "REQ_ERR"   // erroneous request (wrong data, etc.)
)

// Response error
type ResError struct {
	Code  string `json:"code,omitempty"`
	Err   string `json:"err,omitempty"`
	Msg   string `json:"msg,omitempty"`
	Stack string `json:"stack,omitempty"`
}

// Satisfies error interface.
func (e *ResError) Error() string {
	return e.Code + ": " + e.Msg + ": " + e.Err
}

// Parses error and creates new ResError.
func New(err error) *ResError {

	s := err.Error()
	if strings.HasPrefix(s, ERR_PREFIX_BAD_DATA) {
		// Remove prefix
		return &ResError{ERR_REQ, "Request error", s[len(ERR_PREFIX_BAD_DATA):], ""}
	} else {
		return &ResError{ERR_APP, s, "Application error", ""}
	}
}

// New appliation error.
func New_AppErr(err error, msg string) *ResError {
	return &ResError{ERR_APP, err.Error(), msg, ""}
}

// New appliation error with stack.
func New_AppErrWithStack(err error, msg, stack string) *ResError {
	return &ResError{ERR_APP, err.Error(), msg, stack}
}

// URL not found.
func New_NotFound(url, msg string) *ResError {
	return &ResError{ERR_NOT_FOUND, url, msg, ""}
}

// User not authenticated.
func New_NotAuthd() *ResError {
	return &ResError{ERR_NOT_AUTHD, "User not authenticated", "Please login to proceed", ""}
}
