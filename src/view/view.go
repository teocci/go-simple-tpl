// Package view
// Created by RTT.
// Author: teocci@yandex.com on 2022-12월-24
package view

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// HTMLContentType representation header is used to indicate the original media type of the resource
var HTMLContentType = []string{"text/html; charset=utf-8"}

// DefaultConfig default config
var DefaultConfig = Config{
	Root:         "views",
	Extension:    ".html",
	Master:       "layouts/master",
	Partials:     []string{},
	Funcs:        make(template.FuncMap),
	DisableCache: false,
	Delimiters:   Delimiters{Left: "{{", Right: "}}"},
}

// Engine view template engine
type Engine struct {
	config      Config
	tplMap      map[string]*template.Template
	tplMutex    sync.RWMutex
	fileHandler FileHandler
}

// Config configuration options
type Config struct {
	Root         string           // view root
	Extension    string           // template extension
	Master       string           // template master
	Partials     []string         // template partial, such as head, foot
	Funcs        template.FuncMap // template functions
	DisableCache bool             // disable cache, debug mode
	Delimiters   Delimiters       // delimiters
}

// M map interface for data
type M map[string]interface{}

// Delimiters delimiters for template
type Delimiters struct {
	Left  string
	Right string
}

// FileHandler file handler interface
type FileHandler func(config Config, tplFile string) (content string, err error)

// New creates a template engine
func New(config Config) *Engine {
	return &Engine{
		config:      config,
		tplMap:      make(map[string]*template.Template),
		tplMutex:    sync.RWMutex{},
		fileHandler: DefaultFileHandler(),
	}
}

// Default new default template engine
func Default() *Engine {
	return New(DefaultConfig)
}

// Render render template with http.ResponseWriter
func (e *Engine) Render(w http.ResponseWriter, statusCode int, name string, data interface{}) error {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = HTMLContentType
	}
	w.WriteHeader(statusCode)
	return e.executeRender(w, name, data)
}

// RenderWriter render template with io.Writer
func (e *Engine) RenderWriter(w io.Writer, name string, data interface{}) error {
	return e.executeRender(w, name, data)
}

func (e *Engine) executeRender(out io.Writer, name string, data interface{}) error {
	useMaster := true
	if filepath.Ext(name) == e.config.Extension {
		useMaster = false
		name = strings.TrimSuffix(name, e.config.Extension)
	}

	return e.executeTemplate(out, name, data, useMaster)
}

func (e *Engine) executeTemplate(out io.Writer, name string, data interface{}, useMaster bool) error {
	var tpl *template.Template
	var err error
	var ok bool

	functions := make(template.FuncMap, 0)
	functions["include"] = func(layout string) (template.HTML, error) {
		buf := new(bytes.Buffer)
		err := e.executeTemplate(buf, layout, data, false)
		return template.HTML(buf.String()), err
	}
	// Get the plugin collection
	for k, v := range e.config.Funcs {
		functions[k] = v
	}

	e.tplMutex.RLock()
	tpl, ok = e.tplMap[name]
	e.tplMutex.RUnlock()

	exeName := name
	if useMaster && e.config.Master != "" {
		exeName = e.config.Master
	}

	if !ok || e.config.DisableCache {
		tpls := make([]string, 0)
		if useMaster {
			// render()
			if e.config.Master != "" {
				tpls = append(tpls, e.config.Master)
			}
		}
		tpls = append(tpls, name)
		tpls = append(tpls, e.config.Partials...)

		// Loop through each template and test the full path
		tpl = template.New(name).Funcs(functions).Delims(e.config.Delimiters.Left, e.config.Delimiters.Right)
		for _, v := range tpls {
			var data string
			data, err = e.fileHandler(e.config, v)
			if err != nil {
				return err
			}

			var tmpl *template.Template
			if v == name {
				tmpl = tpl
			} else {
				tmpl = tpl.New(v)
			}
			_, err = tmpl.Parse(data)
			if err != nil {
				return fmt.Errorf("ViewEngine render parser name:%v, error: %v", v, err)
			}
		}
		e.tplMutex.Lock()
		e.tplMap[name] = tpl
		e.tplMutex.Unlock()
	}

	// Display the content to the screen
	err = tpl.Funcs(functions).ExecuteTemplate(out, exeName, data)
	if err != nil {
		return fmt.Errorf("ViewEngine execute template error: %v", err)
	}

	return nil
}

// SetFileHandler set file handler
func (e *Engine) SetFileHandler(handle FileHandler) {
	if handle == nil {
		panic("FileHandler can't set nil!")
	}
	e.fileHandler = handle
}

// DefaultFileHandler new default file handler
func DefaultFileHandler() FileHandler {
	return func(config Config, tplFile string) (content string, err error) {
		// Get the absolute path of the root template
		path, err := filepath.Abs(config.Root + string(os.PathSeparator) + tplFile + config.Extension)
		if err != nil {
			return "", fmt.Errorf("ViewEngine path:%v error: %v", path, err)
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("ViewEngine render read name:%v, path:%v, error: %v", tplFile, path, err)
		}

		return string(data), nil
	}
}
