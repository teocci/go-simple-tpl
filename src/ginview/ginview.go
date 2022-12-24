// Package ginview
// Created by Teocci.
// Author: teocci@yandex.com on 2022-12월-24
package ginview

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"

	"github.com/teocci/go-simple-tpl/src/view"
)

const templateEngineKey = "ginview-engine"

// ViewEngine view engine for gin
type ViewEngine struct {
	*view.Engine
}

// ViewRender view render implement gin interface
type ViewRender struct {
	Engine *ViewEngine
	Name   string
	Data   interface{}
}

// New creates a view engine for gin
func New(config view.Config) *ViewEngine {
	return Wrap(view.New(config))
}

// Wrap wraps a view engine for view.ViewEngine
func Wrap(engine *view.Engine) *ViewEngine {
	return &ViewEngine{
		Engine: engine,
	}
}

// Default new default engine
func Default(path string) *ViewEngine {
	config := view.DefaultConfig
	path = filepath.Join(path, config.Root)
	config.Root = path
	return New(config)
}

// Instance implement gin interface
func (e *ViewEngine) Instance(name string, data interface{}) render.Render {
	return ViewRender{
		Engine: e,
		Name:   name,
		Data:   data,
	}
}

// HTML render html
func (e *ViewEngine) HTML(ctx *gin.Context, code int, name string, data interface{}) {
	instance := e.Instance(name, data)
	ctx.Render(code, instance)
}

// Render (YAML) marshals the given interface object and writes data with custom ContentType.
func (v ViewRender) Render(w http.ResponseWriter) error {
	return v.Engine.RenderWriter(w, v.Name, v.Data)
}

// WriteContentType write html content type
func (v ViewRender) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = view.HTMLContentType
	}
}

// NewMiddleware gin middleware for func `gintemplate.HTML()`
func NewMiddleware(config view.Config) gin.HandlerFunc {
	return Middleware(New(config))
}

// Middleware gin middleware wrapper
func Middleware(e *ViewEngine) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(templateEngineKey, e)
	}
}

// HTML is  a render for template
// You should use helper func `Middleware()` to set the supplied
// TemplateEngine and make `HTML()` work validly.
func HTML(ctx *gin.Context, code int, name string, data interface{}) {
	if val, ok := ctx.Get(templateEngineKey); ok {
		if e, ok := val.(*ViewEngine); ok {
			e.HTML(ctx, code, name, data)
			return
		}
	}
	ctx.HTML(code, name, data)
}
