package reqres

//------------------------------------------------------------
// Router for XHR
//------------------------------------------------------------

type XHR_Router map[string]func(*Req) (interface{}, error)

//------------------------------------------------------------
// Router for WebSocket
//------------------------------------------------------------

type WebSocket_Router struct {
	URL   string
	Procs map[string]func(*Req) (interface{}, error)
}
