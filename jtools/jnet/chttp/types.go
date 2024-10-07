package chttp

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"jtools/dbg"
)

// JsonType :
type JsonType map[string]interface{}

func (my JsonType) String() string { return dbg.ToJsonString(my) }

// InjectStruct :
func (my *JsonType) InjectStruct(p interface{}) JsonType {
	b, _ := json.Marshal(p)
	json.Unmarshal(b, my)
	return *my
}

// ParseStruct :
func (my JsonType) ParseStruct(p interface{}) JsonType {
	my.InjectStruct(p)
	return my
}

// Bytes : json.Marshal
func (my JsonType) Bytes() []byte {
	buffer, _ := json.Marshal(my)
	return buffer
}

// RequestBody : http-POST : application/x-www-form-urlencoded
func (my JsonType) RequestBody() string {
	// b, _ := json.Marshal(my)
	// fdata := map[string]interface{}{}
	// json.Unmarshal(b, &fdata)

	u := url.Values{}
	for k, v := range my {
		u.Set(k, fmt.Sprint(v))
	}
	return u.Encode()
}

// GetString :
func (my JsonType) GetString(key string) string {
	if v, isDo := my[key].(string); isDo == true {
		return v
	}
	return ""
}

// GetBool :
func (my JsonType) GetBool(key string) bool {
	if v, isDo := my[key].(bool); isDo == true {
		return v
	}
	return false
}

// GetFloat64 :
func (my JsonType) GetFloat64(key string) float64 {
	if v, do := my[key]; do == false {
		return 0
	} else {
		var reval float64
		switch v.(type) {
		case int:
			reval = float64(v.(int))
		case uint:
			reval = float64(v.(uint))
		case int8:
			reval = float64(v.(int8))
		case uint8:
			reval = float64(v.(uint8))
		case int16:
			reval = float64(v.(int16))
		case uint16:
			reval = float64(v.(uint16))
		case int32:
			reval = float64(v.(int32))
		case uint32:
			reval = float64(v.(uint32))
		case int64:
			reval = float64(v.(int64))
		case uint64:
			reval = float64(v.(uint64))
		case float32:
			reval = float64(v.(float32))
		case float64:
			reval = float64(v.(float64))
		case string:
			var err error
			reval, err = strconv.ParseFloat(v.(string), 64)
			if err != nil {
				fmt.Println("[JsonType]string_to_float64_fail :", v)
				reval = 0
			}
		default:
			fmt.Println("[JsonType]fail_to_float64 :", v)
			reval = 0
		}
		return reval
	}
}

// GetInt64 :
func (my JsonType) GetInt64(key string) int64 {
	if v, do := my[key]; do == false {
		return 0
	} else {
		var reval int64
		switch v.(type) {
		case int:
			reval = int64(v.(int))
		case uint:
			reval = int64(v.(uint))
		case int8:
			reval = int64(v.(int8))
		case uint8:
			reval = int64(v.(uint8))
		case int16:
			reval = int64(v.(int16))
		case uint16:
			reval = int64(v.(uint16))
		case int32:
			reval = int64(v.(int32))
		case uint32:
			reval = int64(v.(uint32))
		case int64:
			reval = int64(v.(int64))
		case uint64:
			reval = int64(v.(uint64))
		case float32:
			reval = int64(v.(float32))
		case float64:
			reval = int64(v.(float64))
		case string:
			var err error
			reval, err = strconv.ParseInt(v.(string), 10, 64)
			if err != nil {
				fmt.Println("[JsonType]string_to_int64_fail :", v)
				reval = 0
			}
		default:
			fmt.Println("[JsonType]fail_to_int64 :", v)
			reval = 0
		}
		return reval
	}
}

//String:
// func (j JsonType) String() string {
// 	var message string
// 	for key, value := range j {
// 		message = fmt.Sprintf("%s%s\n", message, fmt.Sprintf("%s : %#v", key, value))
// 	} //for
// 	return message
// }

// AndroidType : 안드로이드에서는 http.StatusOK 가 아니면 Exception으로 타기때문에
// 심각한 에러가 아니면 200으로 보낸후 클라쪽에서 예외 처리 하도록 유도 하게끔 하자.
func AndroidType() JsonType {
	return nil
}

// JSONStatusOk :
func JSONStatusOk() JsonType {
	ok := JsonType{
		"UTC":    time.Now().UTC().Unix(),
		"Result": "Success",
	}
	return ok
}

// JSONStatusOkParam :
func JSONStatusOkParam(pm JsonType) JsonType {
	ok := JsonType{
		"UTC":    time.Now().UTC().Unix(),
		"Result": "Success",
	}

	for k, v := range pm {
		ok[k] = v
	}

	return ok
}

// // BsonValue : 만약 키값이 없으면 interface.(nil) 을 반환 함..
// func BsonValue(bs interface{}, key string) interface{} {
// 	bm := bs.(bson.M)
// 	return bm[key]
// }

// JsTimeUnix : javascript Date format!
func JsTimeUnix() int64 {
	return time.Now().Unix() / 1000000 //101625400
}

// YMDSFormat : 2012/12/12 20:30:05
func YMDSFormat(nsec int64) string {

	t := time.Unix(0, nsec).UTC()

	str := t.String()
	ss := strings.Split(str, " ")
	ymd := strings.Split(ss[0], "-")
	hmc := strings.Split(ss[1], ".")

	var result string
	result = fmt.Sprintf("%v/%v/%v", ymd[0], ymd[1], ymd[2])
	result = fmt.Sprintf("%v %v", result, hmc[0])

	return result
}
