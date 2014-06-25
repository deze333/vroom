package util

import (
    "net/http"
)

//------------------------------------------------------------
// HTTP Request data
//------------------------------------------------------------

// Returns most probable real IP address from which request came from.
func GetRequestIP(r *http.Request) string {

    var ip string
    if val, ok := r.Header["X-Real-Ip"]; ok && len(val) >= 1 {
        ip = val[0]
    } else 
    if val, ok := r.Header["X-Forwarded-For"]; ok && len(val) >= 1 {
        ip = val[0]
    } else {
        ip = r.RemoteAddr
    }

    return ip
}
