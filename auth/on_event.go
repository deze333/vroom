package auth

//------------------------------------------------------------
// Registry of event listeners
//------------------------------------------------------------

// On user has been de-authenticated
var _deAuthListeners = []func(string){}

// Add de-auth listener callback.
func AddListener_DeAuth(fn func(string)) {
	_deAuthListeners = append(_deAuthListeners, fn)
}

// Broadcast de-auth event for specific user.
func broadcastDeAuth(authId string) {
	for _, fn := range _deAuthListeners {
		fn(authId)
	}
}
