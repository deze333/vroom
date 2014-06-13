package api_ws

import (
	"fmt"
	"encoding/json"
    "github.com/deze333/vroom/errors"
)

//------------------------------------------------------------
// Websocket Response
//------------------------------------------------------------

type Res struct {
	Id            int          `json:"_id"`
	Op            string       `json:"op,omitempty"`
	Err           interface{}  `json:"_err,omitempty"`
	Data          interface{}  `json:"data,omitempty"`
}

//------------------------------------------------------------
// Methods
//------------------------------------------------------------

func (r Res) Json() (b []byte) {

    b, err := json.Marshal(r)
    if err == nil {
        return
    }

    // Error marshalling response
    resErr := errors.New_AppErr(err, "Cannot marshal JSON response")
    resErrJson, _ := json.Marshal(resErr)
    res := fmt.Sprintf(`{"_id": "%v", "op": "%v", "_err": "%v"}`,
        r.Id, r.Op, string(resErrJson))
    return []byte(res)
}

