package dbg

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
	"txscheduler/brix/tools/database/mongo/tools/cc"
)

func ToJSONString(v interface{}) string {
	return ToJsonString(v)
}
func ToJsonString(v interface{}) string {
	if v == nil {
		return ""
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintln(err)

	} else if !json.Valid(b) {
		return fmt.Sprintln(string(b))

	}

	return fmt.Sprintln(string(b))
}

func Cat(a ...interface{}) string {
	msg := ""
	for _, v := range a {
		msg += fmt.Sprintf("%v", v)
	}
	return msg
}

func Error(a ...interface{}) error {
	msg := []string{}
	for _, v := range a {
		msg = append(msg, fmt.Sprintf("%v", v))
	} //for
	return errors.New(strings.Join(msg, " "))
}

func IsTrue(a interface{}) bool {
	switch v := a.(type) {
	case bool:
		return v

	case []bool:
		if len(v) > 0 {
			return v[0]
		}

	case string:
		item := strings.ToLower(strings.TrimSpace(v))
		return item == "true"
	case []string:
		if len(v) > 0 {
			item := strings.ToLower(strings.TrimSpace(v[0]))
			return item == "true"
		}
	case int, int8, int16, int32, int64:
		return v != 0

	case uint, uint8, uint16, uint32, uint64:
		return v != 0

	case float32, float64:
		return v != 0

	}

	return false
}
func StackError(a ...interface{}) string {
	subject := []interface{}{"[ StackError ]  "}
	subject = append(subject, a...)
	subject = append(subject, "          ")
	stackString := string(debug.Stack())

	cc.RedItalicBG(subject...)
	cc.RedItalic(stackString)

	msg := ""
	for _, s := range subject {
		msg = fmt.Sprintf("%v%v\n", msg, s)
	} //for
	msg = fmt.Sprintf("%v%v\n", msg, stackString)
	return msg

}

func ChangeStruct(src, dst interface{}) error {
	return ParseStruct(src, dst)
}
func ParseStruct(src, dst interface{}) error {
	if src == nil {
		return Error("src is nil")
	}
	if dst == nil {
		return Error("dst is nil")
	}

	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, dst)
	case string:
		return json.Unmarshal([]byte(v), dst)
	}

	b, err := json.Marshal(src)
	if err != nil {
		return Error("[Marshal]", err)
	}

	return json.Unmarshal(b, dst)
}
