package doc

import (
	"fmt"
	"reflect"
)

//Object :
type Object interface {
	Message(v string, color ...Color)
	Comment(v string) Object
	Update(isLocal ...bool)

	URL(v string, params ...pairArray) Object
	URLS(kvs ...string) Object
	WS(v string) Object
	Method(v string) Object
	POST() Object
	GET() Object
	Header(kvlist ...keyValue) Object
	Param(kvlist ...keyValue) Object
	JParam(v interface{}, tag ...string) Object
	//GetParam(pair ...string) Object
	Result(code int, kvlist ...keyValue) Object
	ResultOK(kvlist ...keyValue) Object
	JResultOK(v interface{}, tag ...string) Object
	JAckOK(v interface{}, tag ...string) Object
	JAckNone() Object
	JAckError(err IAckError, tag ...string) Object
	ResultERRR(etype IError, tag ...string) Object
	ResultBadParameter(tag ...string) Object
	ResultBadRequest(tag ...string) Object
	ResultUnauthorized(tag ...string) Object
	ResultServerError(tag ...string) Object
	ResultNotFound(tag ...string) Object
	ResultConflict(tag ...string) Object
	Etc(v interface{}, tag ...string) Object
	ETC(ev etcValue, tag ...string) Object
	ETCVAL(void interface{}, pair ...string) Object
	Tag(tag ...string) Object

	Apply(color ...Color)
	Apply2(size string, w Weight, c Color)
	Set(color ...Color)
	Set2(size string, w Weight, c Color)
	ApplyItems(list ItemList)

	Bytes() []byte
}

//IError :
type IError interface {
	Desc() string
}

//FieldTypeName :
func FieldTypeName(i interface{}) string {
	r := reflect.TypeOf(i)
	return fmt.Sprintf("%v : %v", r.Name(), r.Kind())
}
