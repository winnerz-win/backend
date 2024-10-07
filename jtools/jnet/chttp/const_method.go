package chttp

const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	PATCH   = "PATCH"
	HEAD    = "HEAD"
	OPTIONS = "OPTIONS"
)

const (
	DefaultAssetsName    = "assets"
	DefaultTemplatesName = "templates"

	ServeFilesPath = DefaultAssetsName
	RenderRootPath = DefaultTemplatesName
)
