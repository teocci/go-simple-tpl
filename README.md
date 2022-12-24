# go-simple-tpl

[![GoDoc Widget]][GoDoc]  [![GoReportCard Widget]][GoReportCard]

`go-simple-tpl` is a lightweight, minimalist and idiomatic template library based on
golang [html/template][1] for building Go web application.

## Contents

- [Install](#install)
- [Features](#features)
- [Usage](#usage)
    - [Overview](#overview)
    - [Config](#config)
    - [Include syntax](#include-syntax)
    - [Render name](#render-name)
    - [Custom template functions](#custom-template-functions)
- [Examples](#examples)
    - [Basic example](#basic-example)
    - [Advance example](#advance-example)
    - [Multiple example](#multiple-example)
    - [Gin example](#gin-example)
    - [more examples](#more-examples)

## Install

```bash
go get github.com/teocci/go-simple-tpl
```

## Features

* **Lightweight** - use golang html/template syntax.
* **Easy** - easy use for your web application.
* **Fast** - Support configure cache template.
* **Include syntax** - Support include file.
* **Master layout** - Support configure master layout file.
* **Extension** - Support configure template file extension.
* **Easy** - Support configure templates directory.
* **Auto reload** - Support dynamic reload template(disable cache mode).
* **Multiple Engine** - Support multiple templates for frontend and backend.
* **No external dependencies** - plain ol' Go html/template.
* **Gin** - Support gin framework

## Usage

### Overview

Project structure:

```text
|-- app/views/
  |--- index.html
  |--- page.html
  |-- /layouts/
    |--- footer.html
    |--- master.html
```

Use default instance:

```go
// write http.ResponseWriter
//"index" -> index.html
view.Render(writer, http.StatusOK, "index", view.M{})
```

Use new instance with config:

```go
gv := view.New(view.Config{
  Root:      "views",
  Extension: ".tpl",
  Master:    "layouts/master",
  Partials:  []string{"partials/ad"},
  Funcs: template.FuncMap {
    "sub": func (a, b int) int {
        return a - b
    },
    "copy": func () string {
        return time.Now().Format("2006")
    },
  },
  DisableCache: true,
  Delimiters:    Delimiters{Left: "{{", Right: "}}"},
})

//Set new instance
view.Use(gv)

//write http.ResponseWriter
view.Render(writer, http.StatusOK, "index", view.M{})
```

Use multiple instance with config:

```go
    //============== Frontend ============== //
gvFrontend := view.New(view.Config{
  Root:      "views/frontend",
  Extension: ".tpl",
  Master:    "layouts/master",
  Partials:  []string{"partials/ad"},
  Funcs: template.FuncMap{
    "sub": func (a, b int) int {
        return a - b
    },
    "copy": func () string {
        return time.Now().Format("2006")
    },
  },
  DisableCache: true,
  Delimiters:       Delimiters{Left: "{{", Right: "}}"},
})

//write http.ResponseWriter
gvFrontend.Render(writer, http.StatusOK, "index", view.M{})

//============== Backend ============== //
gvBackend := view.New(view.Config{
  Root:      "views/backend",
  Extension: ".tpl",
  Master:    "layouts/master",
  Partials:  []string{"partials/ad"},
  Funcs: template.FuncMap{
    "sub": func (a, b int) int {
      return a - b
    },
    "copy": func () string {
      return time.Now().Format("2006")
    },
  },
  DisableCache: true,
  Delimiters:   Delimiters{Left: "{{", Right: "}}"},
})

//write http.ResponseWriter
gvBackend.Render(writer, http.StatusOK, "index", view.M{})
```

### Config

```go
view.Config{
  Root:      "views", //template root path
  Extension: ".tpl",  //file extension
  Master:    "layouts/master", //master layout file
  Partials:  []string{"partials/head"}, //partial files
  Funcs: template.FuncMap{
    "sub": func (a, b int) int {
    return a - b
  },
  // more funcs
  },
  DisableCache: false, //if disable cache, auto reload template file for debug.
  Delimiters: Delimiters{Left: "{{", Right: "}}"},
}
```

### Include syntax

```twig
//template file
{{include "layouts/footer"}}
```

### Render name:

Render name use `index` without `.html` extension, that will render with master layout.

- **"index"** - Render with master layout.
- **"index.html"** - Not render with master layout.

```
Notice: `.html` is default template extension, you can change with config
```

Render with master

```go
//use name without extension `.html`
view.Render(w, http.StatusOK, "index", view.M{})
```

The `w` is instance of  `http.ResponseWriter`

Render only file(not use master layout)

```go
//use full name with extension `.html`
view.Render(w, http.StatusOK, "page.html", view.M{})
```

### Custom template functions

We have two type of functions `global functions`, and `temporary functions`.

`Global functions` are set within the `config`.

```go
view.Config{
Funcs: template.FuncMap{
"reverse": e.Reverse,
},
}
```

```go
//template file
{{ reverse "route-name" }}
```

`Temporary functions` are set inside the handler.

```go
http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
  err := view.Render(w, http.StatusOK, "index", view.M{
   "reverse": e.Reverse,
  })
  if err != nil {
   fmt.Fprintf(w, "Render index error: %v!", err)
  }
})
```

```go
//template file
{{ call $.reverse "route-name" }}
```

## Examples

See [_examples/](https://github.com/teocci/go-simple-tpl/blob/master/_examples/) for a variety of examples.

### Basic example

```go
package main

import (
	"fmt"
	"net/http"
	
    "github.com/teocci/go-simple-tpl/scr/view"
)

func main() {

	//render index use `index` without `.html` extension, that will render with master layout.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err :=
		view.Render(w, http.StatusOK, "index",
		view.M{
			"title": "Index title!",
			"add": func(a int, b int) int {
				return a + b
			},
		})
		if err != nil {
			fmt.Fprintf(w, "Render index error: %v!", err)
		}

	})

	//render page use `page.tpl` with '.html' will only file template without master layout.
	http.HandleFunc("/page", func(w http.ResponseWriter, r *http.Request) {
		err :=
		view.Render(w, http.StatusOK, "page.html",
		view.M{"title": "Page file title!!"})
		if err != nil {
			fmt.Fprintf(w, "Render page.html error: %v!", err)
		}
	})

	fmt.Println("Listening and serving HTTP on :9090")
	http.ListenAndServe(":9090", nil)
}
```

Project structure:

```go
|-- app/views/
  |--- index.html
  |--- page.html
  |-- layouts/
    |--- footer.html
    |--- master.html


See in "examples/basic" folder
```

[Basic example](https://github.com/teocci/go-simple-tpl/tree/master/_examples/basic)


### Advance example
```go

package main

import (
  "fmt"
  "html/template"
  "net/http"
  "time"
  
  
  "github.com/teocci/go-simple-tpl/scr/view"
)

func main() {
  gv := view.New(view.Config{
    Root:      "views",
    Extension: ".tpl",
    Master:    "layouts/master",
    Partials:  []string{"partials/ad"},
    Funcs: template.FuncMap{
        "sub": func(a, b int) int {
            return a - b
        },
        "copy": func() string {
            return time.Now().Format("2006")
        },
    },
    DisableCache: true,
  })
  
  //Set new instance
  view.Use(gv)
  
  //render index use `index` without `.html` extension, that will render with master layout.
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    err := view.Render(w, http.StatusOK, "index", view.M{
      "title": "Index title!",
      "add": func(a int, b int) int {
          return a + b
      },
    })
    if err != nil {
      fmt.Fprintf(w, "Render index error: %v!", err)
    }
  })
  
  //render page use `page.tpl` with '.html' will only file template without master layout.
  http.HandleFunc("/page", func(w http.ResponseWriter, r *http.Request) {
      err := view.Render(w, http.StatusOK, "page.tpl", view.M{"title": "Page file title!!"})
      if err != nil {
          fmt.Fprintf(w, "Render page.html error: %v!", err)
      }
  })
  
  fmt.Println("Listening and serving HTTP on :9090")
  http.ListenAndServe(":9090", nil)
}

```

Project structure:
```text
|-- app/views/
  |--- index.tpl          
  |--- page.tpl
  |-- layouts/
    |--- footer.tpl
    |--- head.tpl
  |--- master.tpl
    |-- partials/
    |--- ad.tpl
    

See in "examples/advance" folder
```

[Advance example](https://github.com/teocci/go-simple-tpl/tree/master/_examples/advance)

### Multiple example

```go

package main

import (
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
    "github.com/teocci/go-simple-tpl/src/view"
)

func main() {
	router := gin.Default()

	//new template engine
	router.HTMLRender = view.New(view.Config{
		Root:      "views/fontend",
		Extension: ".html",
		Master:    "layouts/master",
		Partials:  []string{"partials/ad"},
		Funcs: template.FuncMap{
			"copy": func() string {
				return time.Now().Format("2006")
			},
		},
		DisableCache: true,
	})

	router.GET("/", func(ctx *gin.Context) {
		// `HTML()` is a helper func to deal with multiple TemplateEngine's.
		// It detects the suitable TemplateEngine for each path automatically.
		view.HTML(ctx, http.StatusOK, "index", gin.H{
			"title": "Fontend title!",
		})
	})

	//=========== Backend ===========//

	//new middleware
	mw := view.NewMiddleware(view.Config{
		Root:      "views/backend",
		Extension: ".html",
		Master:    "layouts/master",
		Partials:  []string{},
		Funcs: template.FuncMap{
			"copy": func() string {
				return time.Now().Format("2006")
			},
		},
		DisableCache: true,
	})

	// You should use helper func `Middleware()` to set the supplied
	// TemplateEngine and make `HTML()` work validly.
	backendGroup := router.Group("/admin", mw)

	backendGroup.GET("/", func(ctx *gin.Context) {
		// With the middleware, `HTML()` can detect the valid TemplateEngine.
		view.HTML(ctx, http.StatusOK, "index", gin.H{
			"title": "Backend title!",
		})
	})

	router.Run(":9090")
}


```

Project structure:

```text
|-- app/views/
  |-- fontend/
    |--- index.html
    |-- layouts/
      |--- footer.html
      |--- head.html
      |--- master.html
    |-- partials/
        |--- ad.html
  |-- backend/
    |--- index.html
    |-- layouts/
      |--- footer.html
      |--- head.html
      |--- master.html

See in "examples/multiple" folder
```

[Multiple example](https://github.com/teocci/go-simple-tpl/tree/master/_examples/multiple)

### Gin example

```bash
go get github.com/teocci/go-simple-tpl/supports/ginview
```

```go

package main

import (
	"net/http"
	
    "github.com/gin-gonic/gin"
    "github.com/teocci/go-simple-tpl/src/ginview"
)

func main() {
	router := gin.Default()

	//new template engine
	router.HTMLRender = ginview.Default("")

	router.GET("/", func(ctx *gin.Context) {
		//render with master
		ctx.HTML(http.StatusOK, "index", gin.H{
			"title": "Index title!",
			"add": func(a int, b int) int {
				return a + b
			},
		})
	})

	router.GET("/page", func(ctx *gin.Context) {
		//render only file, must full name with extension
		ctx.HTML(http.StatusOK, "page.html", gin.H{"title": "Page file title!!"})
	})

	router.Run(":9090")
}

```

Project structure:

```text
|-- app/views/
  |--- index.html
  |--- page.html
  |-- layouts/
    |--- footer.html
    |--- master.html

See in "examples/basic" folder
```

[Gin example](https://github.com/teocci/go-simple-tpl/tree/master/_examples/gin)



### More examples

See [_examples/](https://github.com/teocci/go-simple-tpl/blob/master/_examples/) for a variety of examples.


### Todo

[ ] Add Partials support directory or glob
[ ] Add functions support.

[1]: https://golang.org/pkg/html/template/
[2]: https://github.com/teocci/go-simple-tpl/tree/master/src/ginview

[GoDoc]: https://godoc.org/github.com/teocci/go-simple-tpl
[GoDoc Widget]: https://godoc.org/github.com/teocci/go-simple-tpl?status.svg
[GoReportCard]: https://goreportcard.com/report/github.com/teocci/go-simple-tpl
[GoReportCard Widget]: https://goreportcard.com/badge/github.com/teocci/go-simple-tpl
[GoCover]: https://goreportcard.com/report/github.com/teocci/go-simple-tpl
[GoCover Widget]: https://goreportcard.com/badge/github.com/teocci/go-simple-tpl