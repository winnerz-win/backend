package dbg

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"

	"jtools/jmath"
)

func ToJsonString(v interface{}) string {
	return ToJSONString(v)
}

// ToJSONString :
func ToJSONString(v interface{}, istag ...bool) string {
	if v == nil {
		return ""
	}

	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintln("dbg.ViewJSON :", err)

	} else if json.Valid(b) == false {
		return fmt.Sprintln("dbg.ViewJSON is not Valid :", string(b))

	}

	tag := false
	if len(istag) > 0 && istag[0] == true {
		tag = true
	}

	tcode := "---------------------------------------------------"
	result := ""
	if tag {
		result += fmt.Sprintln(tcode)
	}
	result += fmt.Sprintln(string(b))
	if tag {
		result += fmt.Sprintln(tcode)
	}
	return result
}

// ToJSONTag : omitempty , -
func ToJSONTag(v interface{}, tag string) string {
	if v == nil {
		return ""
	}

	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintln("dbg.ToJSONTag :", err)

	} else if json.Valid(b) == false {
		return fmt.Sprintln("dbg.ToJSONTag is not Valid :", string(b))

	}

	tcode := "< " + tag + " >"
	result := fmt.Sprintln(tcode)
	result += fmt.Sprintln(string(b))

	return result
}

// ViewJSON :
func ViewJSON(v interface{}, istag ...bool) {
	fmt.Println(ToJSONString(v, istag...))
}

// ViewJSONTag :
func ViewJSONTag(v interface{}, tag string) {
	fmt.Println(ToJSONTag(v, tag))
}

// TrimToLower : ToLower & TrimSpace
func TrimToLower(str string) string {
	return strings.ToLower(strings.TrimSpace(str))
}

func TrimToLowers(a ...*string) {
	for _, v := range a {
		*v = strings.ToLower(strings.TrimSpace(*v))
	}
}

// TrimToUpper : ToUpper & TrimSpace
func TrimToUpper(str string) string {
	return strings.ToUpper(strings.TrimSpace(str))
}

func TrimToUppers(a ...*string) {
	for _, v := range a {
		*v = strings.ToUpper(strings.TrimSpace(*v))
	}
}

// Trim :
func Trim(a ...*string) {
	for _, v := range a {
		*v = strings.TrimSpace(*v)
	}
}
func Trims(a *[]string) {
	for i, _ := range *a {
		(*a)[i] = strings.TrimSpace((*a)[i])
	}
}

func Stack() string {
	sl := strings.Split(string(debug.Stack()), "\n")
	if len(sl) >= 5 {
		/*
			[ 0 ] goroutine 6 [running]:
			[ 1 ] runtime/debug.Stack()
			[ 2 ] 	C:/Program Files/Go/src/runtime/debug/stack.go:24 +0x7a
			[ 3 ] jtools/dbg._print_stack()
			[ 4 ] 	d:/work/go/src/brix_pkg/jtools/dbg/error.go:32 +0x2e
		*/
		sl = sl[5:]
	}
	// for i, v := range sl {
	// 	Println("[", i, "]", v)
	// }

	return strings.Join(sl, "\n")
}

// BoolsOne :
func BoolsOne(vals ...bool) bool {
	return len(vals) > 0 && vals[0]
}

// IsTrue2 : BoolsOne
func IsTrue2(vals ...bool) bool {
	return BoolsOne(vals...)
}

func IsTrue(p interface{}) bool {
	switch v := p.(type) {
	case bool:
		return v
	case []bool:
		if len(v) > 0 {
			return v[0]
		}
	case string:
		return TrimToLower(v) == "true"
	case []interface{}:
		if len(v) > 0 {
			return IsTrue(v[0])
		}
	}
	return false
}

// IsNum :
func IsNum(val string) bool {
	return jmath.IsNum(val)
}

// DotString :
func DotString(hex string, size ...int) string {
	ds := 8
	if len(size) > 0 && size[0] > 0 && size[0] < len(hex)-1 {
		ds = size[0]
	}
	if len(hex) > ds {
		clone := hex[:ds] + "~"
		return clone
	}
	return hex
}

// Void : ToString param ....
func Void(a interface{}) string {
	return fmt.Sprintf("%v", a)
}

// D : 단순 더미 주석용 (do-space)
func D(a ...interface{}) string {
	msg := ""
	for _, v := range a {
		msg += fmt.Sprintf("%v ", v)
	}
	return msg
}

// Cat : no-space
func Cat(a ...interface{}) string {
	msg := ""
	for _, v := range a {
		msg += fmt.Sprintf("%v", v)
	}
	return msg
}
func JsonString(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
func JsonBuffer(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

// Coma : a,b,c,d,...
func Coma(a ...interface{}) string {
	if len(a) == 0 {
		return ""
	}
	sl := []string{}
	for _, v := range a {
		sl = append(sl, fmt.Sprint(v))
	}
	return strings.Join(sl, ",")
}

// Squard :
func Squard(v int, count int) int {
	if count > 0 {
		d := v
		for count > 0 {
			v = v * d
			count--
		}
	}
	return v
}

// Int64 :
func Int64(val string) int64 {
	i, _ := strconv.ParseInt(val, 10, 64)
	return i
}

// Int :
func Int(val string) int {
	return int(Int64(val))
}

// Int32 :
func Int32(val string) int32 {
	return int32(Int64(val))
}
