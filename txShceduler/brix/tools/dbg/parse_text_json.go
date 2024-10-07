package dbg

import (
	"reflect"
	"strconv"
	"strings"
)

type jTextMap struct {
	data map[string]string
}

func (my jTextMap) String(key string) string {
	if v, do := my.data[key]; do {
		return v
	}
	return ""
}
func (my jTextMap) Int64(key string) int64 {
	if v, do := my.data[key]; do {
		val, _ := strconv.ParseInt(v, 10, 64)
		return val
	}
	return 0
}
func (my jTextMap) Int(key string) int {
	return int(my.Int64(key))
}
func (my jTextMap) Bool(key string) bool {
	if v, do := my.data[key]; do {
		return IsTrue(v)
	}
	return false
}

func (my jTextMap) Float64(key string) float64 {
	if v, do := my.data[key]; do {
		val, _ := strconv.ParseFloat(v, 64)
		return val
	}
	return 0
}
func (my jTextMap) Float32(key string) float32 {
	return float32(my.Float64(key))
}

//JTextParsing : [obsoluted] instead JTextParseMap
func JTextParsing(text string, keys []string) jTextMap {
	jmap := jTextMap{
		data: map[string]string{},
	}
	//key , val , bool
	isKey := func(line string) (string, string, bool) {
		for _, key := range keys {
			if strings.Contains(line, key) {
				div := strings.Split(line, ":")
				div[1] = strings.ReplaceAll(div[1], "\"", "")
				div[1] = strings.ReplaceAll(div[1], ",", "")
				if key == "_id" {
					objIDs := strings.Split(div[1], "(")
					div[1] = strings.ReplaceAll(objIDs[1], ")", "")

					// div[1] = strings.ToLower(div[1])
					// div[1] = strings.ReplaceAll(div[1], "objectidhex(", "") //ObjectIdHex
					// div[1] = strings.ReplaceAll(div[1], "objectid(", "")    //ObjectId
					// div[1] = strings.ReplaceAll(div[1], ")", "")
				}
				div[1] = strings.TrimSpace(div[1])

				return key, div[1], true
			}
		} //for
		return "", "", false
	}
	ss := strings.Split(text, "\n")
	for _, line := range ss {
		key, val, do := isKey(line)
		if do == false {
			continue
		}
		jmap.data[key] = val
	} //for

	return jmap
}

type textMap struct {
	jTextMap
}

//JTextParseMap : [string]FieldPointer
func JTextParseMap(text string, keyValue map[string]interface{}) jTextMap {
	keys := []string{}
	for key, _ := range keyValue {
		keys = append(keys, key)
	} //for
	jmap := JTextParsing(text, keys)

	for key, _ := range jmap.data {
		if keyValue[key] == nil {
			continue
		}
		fieldPtr := reflect.ValueOf(keyValue[key])
		elem := fieldPtr.Elem()
		switch elem.Kind() {
		case reflect.Int:
			elem.SetInt(jmap.Int64(key))
		case reflect.Int64:
			elem.SetInt(jmap.Int64(key))
		case reflect.Int32:
			elem.SetInt(jmap.Int64(key))
		case reflect.Int16:
			elem.SetInt(jmap.Int64(key))
		case reflect.Int8:
			elem.SetInt(jmap.Int64(key))

		case reflect.String:
			elem.SetString(jmap.String(key))

		case reflect.Float64:
			elem.SetFloat(jmap.Float64(key))
		case reflect.Float32:
			elem.SetFloat(jmap.Float64(key))

		case reflect.Bool:
			elem.SetBool(jmap.Bool(key))

		case reflect.Struct:
			{
				Red("struct?")
			}
		case reflect.Ptr: //Struct??
			{
				Red("ppp")
			}
		} //switch
	} //for

	return jmap
}
