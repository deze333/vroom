package auth

import (
    "fmt"
	"net/http"
    "github.com/gorilla/sessions"
    "github.com/deze333/vroom/util"
)

//------------------------------------------------------------
// Authentication functions
//------------------------------------------------------------

// Checks if this user has been authenticated.
func IsAuthd(r *http.Request) bool {

    sess, err := _cookieStore.Get(r, _cookieSessName)
    if err != nil {
        return false
    }
    
    //fmt.Println("\t??? IS LOGGED IN:", sess.Options, sess.Values)

    if _, ok := sess.Values["_auth"]; ok {
        return true
    }

    return false
}

// Authenticates new user by creating a new session.
func Auth(w http.ResponseWriter, r *http.Request, vals map[string]string) (err error) {

    // Get session which may be a new one
    // New session is auto added to the store
    sess, err := _cookieStore.Get(r, _cookieSessName)
    if err != nil {
        return
    }
    
    // Set session values
    sess.Values["_auth"] = util.NewID()
    sess.Values["_ip"] = r.RemoteAddr

    // Add user values
    for k, v := range vals {
        sess.Values[k] = v
    }

    sess.Options = &sessions.Options{
        MaxAge: _cookieMaxAge,
    }

    if _cookiePath != "" {
        sess.Options.Path = _cookiePath
    }

    if _cookieDomain != "" {
        sess.Options.Domain = _cookieDomain
    }

    err = sess.Save(r, w)

    //fmt.Println(">>> LOGIN:", sess.Options, sess.Values)
    return
}

// Retrieves authentication ID from the session.
// Returns -1 if not found.
func GetAuthdId(r *http.Request) int64 {

    sess, _ := _cookieStore.Get(r, _cookieSessName)

    if val, ok := sess.Values["_auth"]; ok {
        if valInt, ok := val.(int64); ok {
            return valInt
        }
    }

    return -1
}

// De-authenticates user.
func DeAuth(w http.ResponseWriter, r *http.Request) (err error) {

    // Get session which may be a new one
    sess, err := _cookieStore.Get(r, _cookieSessName)
    if err != nil {
        return
    }

    //fmt.Println()
    //fmt.Println()
    //fmt.Println("xxx LOGOUT:", sess.Values)

    var authId int64
    if val, ok := sess.Values["_auth"]; ok {
        if authId, ok = val.(int64); !ok {
            authId = -1
        }
    }

    // Remember user, only deauth them
    //delete(sess.Values, "_auth")

    // Forget user completely
    // Set cookie to expire right away
    sess.Values = map[interface{}]interface{}{}
    sess.Options.MaxAge = -1

    err = sess.Save(r, w)

    // Broadcast event
    broadcastDeAuth(authId)
    return
}

// Retrieve existing session.
func GetSessionValues(r *http.Request) (info map[string]string, err error) {

    // Get session which may be a new one
    sess, err := _cookieStore.Get(r, _cookieSessName)
    if err != nil {
        return
    }

    //fmt.Println()
    //fmt.Println("--->", r.RequestURI, ", Session :", sess.Values)

    // Build user info
    info = map[string]string{}
    
    for k, val := range sess.Values {
        key := fmt.Sprint(k)
        if len(key) > 0 && key[0] == '_' {
            continue
        }
        info[key] = fmt.Sprint(val)
    }
    return
}

// XXX Experimental, not used. 
// Adds item to an array of items in session.
func Values_Array_ItemPush(w http.ResponseWriter, r *http.Request, key string, item interface{}) {

    sess, _ := _cookieStore.Get(r, _cookieSessName)
    if val, ok := sess.Values[key]; ok {
        if arr, ok := val.([]interface{}); ok {
            sess.Values[key] = append(arr, item)
        } else {
            sess.Values[key] = []interface{}{item}
        }
    } else {
        sess.Values[key] = []interface{}{item}
    }
    sess.Save(r, w)
}

