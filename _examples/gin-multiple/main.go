/*
 * Copyright 2018 Foolin.  All rights reserved.
 *
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 *
 */

package main

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"time"

	ginview "github.com/teocci/go-simple-tpl/src/ginview"
	"github.com/teocci/go-simple-tpl/src/view"
)

func main() {
	router := gin.Default()

	//new template engine
	router.HTMLRender = ginview.New(view.Config{
		Root:      "_examples/gin-multiple/views/frontend",
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
		ginview.HTML(ctx, http.StatusOK, "index", gin.H{
			"title": "Frontend title!",
		})
	})

	//=========== Backend ===========//

	//new middleware
	mw := ginview.NewMiddleware(view.Config{
		Root:      "_examples/gin-multiple/views/backend",
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
		ginview.HTML(ctx, http.StatusOK, "index", gin.H{
			"title":     "Backend title!",
			"page":      "backend",
			"user":      "kogas",
			"stationId": 4,
		})
	})

	router.Run(":9090")
}
