package chttp

import (
	"fmt"
	"net/http"

	"txscheduler/brix/tools/jpath"

	"github.com/urfave/negroni"

	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
)

const (
	//RenderRootPath :
	RenderRootPath = "templates"
	//ServeFilesPath :
	ServeFilesPath = "assert"
)

var (
	//Renderer :
	Renderer = render.New(render.Options{
		Directory:  RenderRootPath,
		Extensions: []string{".tmpl", ".html"},
	})

	router = httprouter.New()

	ApiCount int
)

//AssertDir :
func AssertDir() string {
	return jpath.NowPath() + "\\" + ServeFilesPath
}

//TempleteDir :
func TempleteDir() string {
	return jpath.NowPath() + "\\" + RenderRootPath
}

// GetRouter : static router
func GetRouter() *httprouter.Router {

	Renderer = render.New(render.Options{
		Directory:  jpath.NowPath() + "\\" + RenderRootPath,
		Extensions: []string{".tmpl", ".html"},
	})

	staticDir := jpath.NowPath() + "\\" + ServeFilesPath
	router.ServeFiles(fmt.Sprintf("/%v/*filepath", ServeFilesPath), http.Dir(staticDir))

	return router
}

// SetRouting :
func SetRouting(n *negroni.Negroni, router *httprouter.Router) {
	n.UseHandler(router)
}

// SetHandlerFunc :
func SetHandlerFunc(n *negroni.Negroni, f negroni.HandlerFunc) {
	n.Use(f)
}

// SetContextHandles :
func SetContextHandles(contexts []PContext, countSkip ...bool) {

	if len(countSkip) == 0 {
		ApiCount += len(contexts)
	}

	for _, ctx := range contexts {
		router.Handle(ctx.Method, ctx.Path, ctx.Handle.getHandle())
	}
}
