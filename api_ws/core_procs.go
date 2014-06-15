package api_ws

import (
)

//------------------------------------------------------------
// Core message processors
//------------------------------------------------------------

// Core processors that get matched first
var _coreProcs = map[string]func(*Req) (interface{}, error) {
    "core/pulse": CorePulseProc,
}

// Response to pulse by returning current app deployment info.
func CorePulseProc(req *Req) (data interface{}, err error) {


    data = map[string]interface{}{
        "sent": req.Params["sent"],
        "appv": _verGetter(),
    }

    return
}
