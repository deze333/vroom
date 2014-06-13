package auth

import ()

//------------------------------------------------------------
// Model
//------------------------------------------------------------

type Member struct {
	RemoteIP string
	Name     string
	Initials string
	Email    string
}
