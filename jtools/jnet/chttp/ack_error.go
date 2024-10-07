package chttp

import (
	"jtools/dbg"
	"sort"
)

const (
	//ErrorServerUnderMaintenance : fixed
	ErrorServerUnderMaintenance = 33
)

type errorHelper map[int]string

func (my errorHelper) toSlice() []cError {
	sl := []cError{}
	for k, v := range my {
		sl = append(sl, cError{k, v})
	}
	sort.Slice(sl, func(i, j int) bool { return sl[i].code < sl[j].code })
	return sl
}

// ERROR :
type ERROR interface {
	Code() int
	Desc() string
	String() string
}

type cError struct {
	code int
	desc string
}

func (my cError) Code() int      { return my.code }
func (my cError) Desc() string   { return my.desc }
func (my cError) String() string { return dbg.Cat("[", my.code, "]", my.desc) }

var (
	errorView              = false //SetErrorView(true)
	ErrorNone              = cError{0, "None"}
	serverUnderMaintenance = cError{ErrorServerUnderMaintenance, "ServerUnderMaintenance"} //33 (서버정검)
	errHelper              = errorHelper{
		ErrorNone.code:              ErrorNone.desc,
		serverUnderMaintenance.code: serverUnderMaintenance.desc,
	}
)

// SetErrorView : errorView = true
func SetErrorView() {
	errorView = true
}

// Error :
func Error(code int, desc string) ERROR {
	errHelper[code] = desc
	return cError{code, desc}
}

// ErrorStrings :
func ErrorStrings() string {
	msg := "===== error_code & error_message =====\n"
	array := errHelper.toSlice()
	for _, v := range array {
		msg += Cat(v.code, "  :  ", v.desc, "\n")
	}
	return msg
}
