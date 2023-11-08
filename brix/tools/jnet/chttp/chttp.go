package chttp

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/urfave/negroni"

	"github.com/julienschmidt/httprouter"
)

const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	PATCH   = "PATCH"
	HEAD    = "HEAD"
	OPTIONS = "OPTIONS"
)

//Context :
type Context struct {
	Method string //GET,POST,...
	Path   string
	Handle RouterHandle
}

//String :
func (my Context) String() string {
	msg := ""
	if my.Method == GET {
		msg = fmt.Sprintf("[ %v  ] ", my.Method)
	} else {
		msg = fmt.Sprintf("[ %v ] ", my.Method)
	}
	msg = fmt.Sprintf("%v%v", msg, my.Path)

	if strings.Contains(my.Path, ":") {
		msg = fmt.Sprintf("%v  ( \":tag\" - multi args. )", msg)
	}

	isWebSocket := false
	if my.Method == GET {
		if strings.HasPrefix(my.Path, "ws") || strings.HasPrefix(my.Path, "wss") {
			isWebSocket = true
		}
	}

	if isWebSocket {
		msg += " --- WebSocket"
	}

	return msg
}

//NewContext :
func NewContext(method string, api string, handle RouterHandle) *Context {
	return &Context{
		Method: method,
		Path:   api,
		Handle: handle,
	}
}

/*PContext :
type Context struct {
	Method string //GET,POST,...
	Path   string
	Handle RouterHandle
}
func(w http.ResponseWriter, req *http.Request, ps chttp.Params)
*/
type PContext *Context

//MContext :
func MContext(method, path string, handle RouterHandle) PContext {
	return &Context{
		Method: method,
		Path:   path,
		Handle: handle,
	}
}

// PContexts : chttp.MContext(chttp.POST , "uri" , nil)
type PContexts []PContext

//Append :
func (my *PContexts) Append(method, path string, handle RouterHandle) {
	(*my) = append((*my), &Context{
		method,
		path,
		handle,
	})
}

//Add :
func (my *PContexts) Add(ctx PContext) {
	(*my) = append((*my), ctx)
}

//AppendList :
func (my *PContexts) AppendList(a PContexts) {
	(*my) = append((*my), a...)
}

//Params :
type Params httprouter.Params

// ByName : returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (my Params) ByName(name string) string {
	hp := httprouter.Params(my)
	return hp.ByName(name)
}

//RouterHandle : func(w http.ResponseWriter, req *http.Request, ps chttp.Params)
type RouterHandle func(w http.ResponseWriter, req *http.Request, ps Params)

func (my RouterHandle) getHandle() httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		my(w, req, Params(ps))
	}
}

//HandlerFunc :
type HandlerFunc negroni.HandlerFunc

func (my HandlerFunc) getHandlerFunc() negroni.HandlerFunc {
	return negroni.HandlerFunc(my)
}
