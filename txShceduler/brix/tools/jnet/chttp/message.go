package chttp

import (
	"fmt"
	"net/http"
)

//Message :
func Message(message interface{}) JsonType {
	return JsonType{"message": message}
}

//MessageBadParam : StatusBadRequest { "message" : "Bad-parameter" }
func MessageBadParam(msg0 ...string) JsonType {
	if len(msg0) > 0 {
		return JsonType{"message": fmt.Sprintf("Bad-parameter:%v", msg0[0])}
	}
	return JsonType{"message": "Bad-parameter"}
}

//ErrorMessage :
func ErrorMessage(err error) JsonType {
	return Message(err.Error())
}

//ResultBadHeader : StatusBadRequest
func ResultBadHeader(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	ResultJSON(w, StatusBadRequest, err)
	return true
}

/////////////////////////////////////////////////////////////////////////

type chttpResult struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
}

//ChttpResult :
func ChttpResult(success bool, datamap ...interface{}) chttpResult {
	r := chttpResult{
		Success: success,
		Data:    map[string]interface{}{},
	}

	key := ""
	for i, v := range datamap {
		if i%2 == 0 {
			key = v.(string)
		} else {
			r.Data[key] = v
		}
	} //for
	return r
}
