package chttp

import (
	"fmt"
	"sort"

	"jtools/jmath"
	"txscheduler/brix/tools/dbg"
)

/////////////////////////////////////////////////////////////////////

type errorHelper map[int]string

func (my errorHelper) toSlice() []cError {
	sl := []cError{}
	for k, v := range my {
		sl = append(sl, cError{k, v})
	}
	sort.Slice(sl, func(i, j int) bool { return sl[i].code < sl[j].code })
	return sl
}

type cError struct {
	code int
	desc string
}

// NewError : 이거 대신 Error(code , desc) 로 쓴다!
func NewError(code int) cError {
	return cError{code: code}
}

func (my cError) ToJson() interface{} {
	v := map[string]interface{}{}
	v["code"] = my.code
	v["desc"] = my.desc
	return v
}

func ParseCError(p interface{}) CError {
	v := map[string]interface{}{}
	dbg.ChangeStruct(p, &v)

	er := cError{
		code: jmath.Int(v["code"]),
		desc: dbg.Cat(v["desc"]),
	}
	return er
}

// CError :
type CError interface {
	Code() int
	Int() int
	String() string
	Desc() string
	Help() string

	ToJson() interface{}
}

var (
	errorView              = false //SetErrorView(true)
	ErrorNone              = cError{0, "None"}
	serverUnderMaintenance = cError{ErrorServerUnderMaintenance, "ServerUnderMaintenance"} //33 (서버정검)
	errHelper              = errorHelper{
		ErrorNone.code:              ErrorNone.desc,
		serverUnderMaintenance.code: serverUnderMaintenance.desc,
	}
)

// Error :
func Error(code int, desc string) cError {
	errHelper[code] = desc
	return cError{code, desc}
}

// ErrorStrings :
func ErrorStrings() string { return cError{}.Help() }

const (
	//ErrorServerUnderMaintenance : 서버 정검중
	ErrorServerUnderMaintenance = 33
)

// Int :
func (my cError) Int() int       { return my.code }
func (my cError) Code() int      { return my.code }
func (my cError) String() string { return fmt.Sprint(my.code) }

// Desc : for interface
func (my cError) Desc() string { return my.desc + " ( " + my.String() + " )" }

func (my cError) DescString() string { return my.desc }

// Help :
func (my cError) Help() string {
	msg := "<cc_blue>===== error_code & error_message =====</cc_blue>\n"
	array := errHelper.toSlice()
	for _, v := range array {
		//msg += `{ "code" : ` + v.String() + ` }  ` + v.desc + " \n"
		msg += "<cc_bold>" + v.String() + "</cc_bold>  :  " + v.desc + dbg.ENTER
	}
	return msg
}

// SetErrorView : errorView = true
func SetErrorView(v ...bool) {
	errorView = true
	errorView = dbg.IsTrue2(v...)
}

// ErrorDesc :
func ErrorDesc(e CError, message string) string {
	r := `{ "code":` + e.String() + `, "etc":any } //` + message
	return r
}
