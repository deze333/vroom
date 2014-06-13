package api_http

import (
    "strings"
	"net/http"
)

//------------------------------------------------------------
// Utils for requests
//------------------------------------------------------------

// Guesses request type from the header.
func requestType(r *http.Request) string {
    if isRequest_HTML(r) {
        return "HTML"
    }
    if isRequest_XHR(r) {
        return "XHR"
    }
    if isRequest_WS(r) {
        return "WS"
    }
    if isRequest_IMG(r) {
        return "IMG"
    }
    return "TEXT"
}

// Is this HTML request.
func isRequest_HTML(r *http.Request) bool {
    for _, param := range r.Header["Accept"] {
        if strings.Contains(param, "text/html") {
            return true
        }
    }
    return false
}

// Is this XHR request.
func isRequest_XHR(r *http.Request) bool {
    for _, param := range r.Header["Accept"] {
        if strings.Contains(param, "application/json") {
            return true
        }
    }
    return false
}

// Is this Websocket request.
func isRequest_WS(r *http.Request) bool {
    for _, param := range r.Header["Connection"] {
        if strings.Contains(param, "Upgrade") {
            return true
        }
    }
    return false
}

// Is this IMAGE request.
func isRequest_IMG(r *http.Request) bool {
    for _, param := range r.Header["Accept"] {
        if strings.Contains(param, "image/webp") {
            return true
        }
    }
    return false
}

