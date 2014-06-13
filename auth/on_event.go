package auth

import ()

//------------------------------------------------------------
// Registry of event listeners
//------------------------------------------------------------

// On user has been de-authenticated
var _deAuthListeners = []func(int64){}

// Add de-auth listener callback.
func AddListener_DeAuth(fn func(int64)) {
    _deAuthListeners = append(_deAuthListeners, fn)
}

// Broadcast de-auth event for specific user.
func broadcastDeAuth(authId int64) {
    for _, fn := range _deAuthListeners {
        fn(authId)
    }
}
