package xparam

import (
    "fmt"
    "errors"
    "time"
    "net/http"
	"labix.org/v2/mgo/bson"
)

//------------------------------------------------------------
// Xparam 
//------------------------------------------------------------

func New(vals map[string]interface{}) XP {
    return vals
}

//------------------------------------------------------------
// Xparam access methods
//------------------------------------------------------------

// Gets http.Request.
func (xp XP) Get_HttpRequest() (req *http.Request) {

    if val, ok := xp["_httpReq"]; ok && val != nil {
        if req, ok = val.(*http.Request); ok {
            return
        }
    }
    return 
}

// Gets http.ResponseWriter.
func (xp XP) Get_HttpResponseWriter() (rw http.ResponseWriter) {

    if val, ok := xp["_httpResWriter"]; ok && val != nil {
        if rw, ok = val.(http.ResponseWriter); ok {
            return
        }
    }
    return 
}

// Gets session values.
func (xp XP) Get_SessionValues() (vals map[string]string) {

    if val, ok := xp["_session"]; ok && val != nil {
        if vals, ok = val.(map[string]string); ok {
            return
        }
    }
    return 
}

// Gets parameter as array of xparams.
func (xp XP) As_XP(key string) (data XP) {

    if val, ok := xp[key]; ok && val != nil {
        if axp, ok := val.(map[string]interface{}); ok {
            data = axp
        }
    }
    return
}

// Gets parameter as array of xparams.
func (xp XP) As_ArrayXP(key string) (data []XP) {

    if val, ok := xp[key]; ok && val != nil {
        if arr, ok := val.([]interface{}); ok {
            data = []XP{}
            for _, obj := range arr {
                if axp, ok := obj.(map[string]interface{}); ok {
                    data = append(data, axp)
                }
            }
        }
    }
    return
}

// Gets key value as bson.ObjectId.
func (xp XP) MustBe_ObjectId(key string) (id bson.ObjectId, err error) {

    if val, ok := xp[key]; ok {
        valstr := fmt.Sprint(val)
        if bson.IsObjectIdHex(valstr) { 
            id = bson.ObjectIdHex(valstr)
        } else {
            err = errors.New("Id parameter is not ObjectId: " + valstr)
        }
    } else {
        err = errors.New("Missing parameter: id")
    }
    return
}

// Gets key value as if it should be bson.ObjectId.
func (xp XP) As_ObjectId(key string) (id *bson.ObjectId) {

    if val, ok := xp[key]; ok {

        // Is it already ObjectId ?
        if oid, ok := val.(bson.ObjectId); ok {
            id = &oid
            return
        }

        // Treat as string
        valstr := fmt.Sprint(val)
        if bson.IsObjectIdHex(valstr) { 
            bsonId := bson.ObjectIdHex(valstr)
            id = &bsonId
        }
    }
    return
}

// Gets parameter as bool.
func (xp XP) As_Bool(key string) (b bool) {

    if val, ok := xp[key]; ok && val != nil {
        if b, ok = val.(bool); ok {
            return
        } else {
            str := fmt.Sprint(val)
            if str == "true" {
                return true
            } else if str == "false" {
                return false
            } 
            return false
        }
    }
    return 
}

// Gets parameter as string.
func (xp XP) As_String(key string) (str string) {

    if val, ok := xp[key]; ok && val != nil {
        str = fmt.Sprint(val)
    }
    return 
}

// Sets parameter as string.
func (xp XP) To_String(to *map[string]string, key string) {
    if val, ok := xp[key]; ok && val != nil {
        str := fmt.Sprint(val)
        if str != "" {
            (*to)[key] = str
        }
    }
}

// Gets parameter as string array.
func (xp XP) As_ArrayString(key string) (data []string) {

    if val, ok := xp[key]; ok && val != nil {
        if data, ok = val.([]string); ok {
            return
        }
    }
    return 
}

// Gets parameter as array of string maps.
func (xp XP) As_ArrayMapString(key string, fields ...string) (data []map[string]string) {

    if val, ok := xp[key]; ok && val != nil {
        if arr, ok := val.([]interface{}); ok {
            data = []map[string]string{}

            for _, obj := range arr {
                if inmap, ok := obj.(map[string]interface{}); ok {
                    outmap := map[string]string{}
                    for k, v := range inmap {
                        isPresent := false
                        for _, field := range fields {
                            if k == field {
                                isPresent = true
                                break
                            }
                        }
                        if isPresent && v != nil {
                            outmap[k] = fmt.Sprint(v)
                        }
                    }
                    if len(outmap) != 0 {
                        data = append(data, outmap)
                    }
                }
            }
            return
        }
    }
    return 
}

// Gets parameter as int array.
func (xp XP) As_ArrayInt(key string) (data []int) {

    if val, ok := xp[key]; ok && val != nil {
        if data, ok = val.([]int); ok {
            return
        }
    }
    return 
}

// Gets parameter as time. Call t.IsZero() to check for error.
func (xp XP) As_Time(key string) (t time.Time) {

    if val, ok := xp[key]; ok && val != nil {
        t, _ = time.Parse(time.RFC3339, fmt.Sprint(val))
    }
    return 
}

// Gets parameter as time.Duration array.
func (xp XP) As_ArrayDuration(key string) (data []time.Duration) {

    if val, ok := xp[key]; ok && val != nil {
        if data, ok = val.([]time.Duration); ok {
            return
        }
    }
    return 
}

// Gets parameter as map of string to sting.
func (xp XP) As_MapString(key string) (data map[string]string) {

    if val, ok := xp[key]; ok && val != nil {
        if amap, ok := val.(map[string]interface{}); ok {
            data = map[string]string{}
            for k, v := range amap {
                if v != nil {
                    data[k] = fmt.Sprint(v)
                }
            }
            return
        }
    }
    return 
}

