// Package webserver
// Created by RTT.
// Author: teocci@yandex.com on 2022-Apr-26

package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/teocci/go-simple-tpl/src/view"
)

func main() {
	gv := view.New(view.Config{
		Root:      "_examples/advance/views",
		Extension: ".tpl",
		Master:    "layouts/master",
		Partials:  []string{"partials/ad"},
		Funcs: template.FuncMap{
			"safeHTML": func(v string) template.HTML {
				return template.HTML(v)
			},
			"copy": func() string {
				return time.Now().Format("2006")
			},
		},
		DisableCache: true,
	})

	//Set new instance
	view.Use(gv)

	rawContent := `This is <b>HTML</b> content! Posted on <time datetime="2019-05-16 01:02:03">May 16</time> by teocci.`

	//render index use `index` without `.tpl` extension, that will render with master layout.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := view.Render(w, http.StatusOK, "index", view.M{
			"Title":       "Index title!",
			"HtmlContent": template.HTML(rawContent),
			"RawContent":  rawContent,
			"tempConvertHTML": func(v string) template.HTML {
				return template.HTML(v)
			},
		})
		if err != nil {
			_, _ = fmt.Fprintf(w, "Render index error: %v!", err)
		}

	})

	//render page use `page.tpl` with '.tpl' will only file template without master layout.
	http.HandleFunc("/page", func(w http.ResponseWriter, r *http.Request) {
		err := view.Render(w, http.StatusOK, "page.tpl", view.M{
			"Title": "Page file title!!",
		})
		if err != nil {
			_, _ = fmt.Fprintf(w, "Render page.html error: %v!", err)
		}
	})

	fmt.Println("Listening and serving HTTP on :9090")
	_ = http.ListenAndServe(":9090", nil)
}
