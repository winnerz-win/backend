package console

import (
	"strconv"
	"strings"
)

//PsCut :
func PsCut(ps []string, idx int, arg *string) bool {
	if idx >= len(ps) {
		return false
	}
	*arg = ps[idx]
	return true
}
func PsIndex(ps []string, idx int, f func(v string)) {
	if idx >= len(ps) {
		return
	}
	f(ps[idx])
}

//PsDo :
func PsDo(ps []string, tstr string) bool {
	for _, v := range ps {
		if v == tstr {
			return true
		}
	} //for
	return false
}

//PsKeyVal : [key=value]
func PsKeyVal(ps []string, key string, call func(val string)) bool {
	for _, v := range ps {
		if strings.HasPrefix(v, key) {
			ss := strings.Split(v, "=")
			call(ss[1])
			return true
		}
	} //for
	return false
}

//PsKeyValFor : call(key=value)
func PsKeyValFor(ps []string, call func(k, v string)) {
	for _, v := range ps {
		if strings.Contains(v, "=") {
			ss := strings.Split(v, "=")
			call(ss[0], ss[1])
		}
	} //for
}

//PsContains :
func PsContains(ps []string, key string, call func(v string)) {
	for _, v := range ps {
		if strings.Contains(v, key) {
			call(v)
			break
		}
	} //for
}

//Int :
func Int(ps string) int {
	v, _ := strconv.ParseInt(ps, 10, 64)
	return int(v)
}

//Int64 :
func Int64(ps string) int64 {
	v, _ := strconv.ParseInt(ps, 10, 64)
	return v
}
