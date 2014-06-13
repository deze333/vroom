package api_xhr

import (
    "fmt"
    "bytes"
	"encoding/json"
    "github.com/deze333/vroom/errors"
)

//------------------------------------------------------------
// XHR Response
//------------------------------------------------------------

// Creates new JSON response.
func NewResponse(v interface{}) []byte {

    res, err := json.Marshal(v)
    if err != nil {
        res = NewResponse_Err(errors.New_AppErr(err, "Cannot marshal JSON response"))
    }

    return res //postEncode(res)
}

// Creates new JSON error response, prefab types.
func NewResponse_Err(err *errors.ResError) []byte {

    res, errr := json.Marshal(map[string]interface{}{
        "_err": err,
    })

    // Just in case even this failed
    if errr != nil {
        res = []byte(fmt.Sprintf(
            `{"_err": {"code": "APP_ERR", "err": "%v", "msg": "%v, original error: %v"}}`, 
            errr, "Error marshalling error to JSON", err))
    }

    return res //postEncode(res)
}

// Makes some post encoding adjustements to achieve correct JSON.
func postEncode(res []byte) []byte {

    // XXX Perhaps I need to read the manual...
    // Fix of strange behaviour when writer expects second % 
    // after first and otherwise says (MISSING), that breaks JSON parser.
    // SOLUTION:
    // Convert % into %%.
    return bytes.Replace(res, []byte{'%'}, []byte{'%','%'}, -1)
}

