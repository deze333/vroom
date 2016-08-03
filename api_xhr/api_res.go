package api_xhr

import (
	"github.com/deze333/vroom/errors"
	"github.com/deze333/vroom/reqres"
)

//------------------------------------------------------------
// XHR Response
//------------------------------------------------------------

// Creates new JSON response.
func NewResponse(req *reqres.Req, data interface{}) []byte {

	res := &reqres.Res{
		Data: data,
	}

	switch req.ResponseMarshalingMode {
	case "DataOnly":
		return res.MarshalDataOnly(req)
	default:
		return res.Marshal(req)
	}
}

// Creates new JSON error response, prefab types.
func NewResponse_Err(req *reqres.Req, err *errors.ResError) []byte {

	res := &reqres.Res{
		Err: err,
	}

	switch req.ResponseMarshalingMode {
	case "DataOnly":
		return res.MarshalDataOnly(req)
	default:
		return res.Marshal(req)
	}
}
