package chttp

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

//Params :
type Params httprouter.Params

// ByName : returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (my Params) ByName(name string) string {
	hp := httprouter.Params(my)
	return hp.ByName(name)
}

type ResponseWriter interface {
	W() http.ResponseWriter
	R() *render.Render

	Header() http.Header
	Write(buf []byte) (int, error)
}
type customResponseWriter struct {
	w http.ResponseWriter
	r *render.Render
}

func (my *customResponseWriter) W() http.ResponseWriter { return my.w }
func (my *customResponseWriter) R() *render.Render      { return my.r }
func (my *customResponseWriter) Write(buf []byte) (int, error) {
	return my.W().Write(buf)
}
func (my *customResponseWriter) Header() http.Header {
	return my.W().Header()
}

type RouterHandle func(w ResponseWriter, req *http.Request, ps Params)

func (my RouterHandle) getHandle(r *render.Render) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		my(&customResponseWriter{w, r}, req, Params(ps))
	}
}

func (my *Classic) ResponsWriter(w http.ResponseWriter) ResponseWriter {
	return &customResponseWriter{
		w: w,
		r: my.renderer,
	}
}

//HandlerFunc :
type HandlerFunc func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)

func (my HandlerFunc) getHandlerFunc() negroni.HandlerFunc {
	return negroni.HandlerFunc(my)
}

type HANDLER interface {
	Method() string
	Path() string
	Handle() RouterHandle
}

type HANDLERLIST []HANDLER

func (my *HANDLERLIST) Append(hs ...HANDLER) {
	*my = append(*my, hs...)
}
func (my *HANDLERLIST) Add(method, path string, h RouterHandle) {
	*my = append(*my,
		MakeHandler(method, path, h),
	)
}

type cHandler struct {
	method string
	path   string
	handle RouterHandle
}

func (my cHandler) Method() string       { return my.method }
func (my cHandler) Path() string         { return my.path }
func (my cHandler) Handle() RouterHandle { return my.handle }

func MakeHandler(method, path string, h RouterHandle) HANDLER {
	return cHandler{
		method: method,
		path:   path,
		handle: h,
	}
}

func newRouterRenderer(serveFilePath string) (*httprouter.Router, *render.Render) {
	applyServeFilePath := ServeFilesPath
	if serveFilePath != "" {
		applyServeFilePath = serveFilePath
	}

	renderPath := NowPath() + "\\" + RenderRootPath
	//Renderer = render.New(render.Options{
	renderer := render.New(render.Options{
		Directory:  renderPath,
		Extensions: []string{".tmpl", ".html"},
	})

	router := httprouter.New()
	//KKJJSS
	staticDir := NowPath() + "\\" + applyServeFilePath

	LogPurple("router.filepath :", staticDir)
	LogPurple("renderer.filepath :", renderPath)

	router.ServeFiles(fmt.Sprintf("/%v/*filepath", applyServeFilePath), http.Dir(staticDir))
	return router, renderer
}
