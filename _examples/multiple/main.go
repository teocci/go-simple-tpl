// Package view
// Created by RTT.
// Author: teocci@yandex.com on 2022-12월-24
package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/teocci/go-simple-tpl/src/view"
)

func main() {
	gvFront := view.New(view.Config{
		Root:      "_examples/multiple/views/frontend",
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := gvFront.Render(w, http.StatusOK, "index", view.M{
			"title": "Frontend title!",
		})
		if err != nil {
			_, _ = fmt.Fprintf(w, "Render index error: %v!", err)
		}
	})

	//=========== Backend ===========//
	gvBackend := view.New(view.Config{
		Root:      "./_examples/multiple/views/backend",
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

	http.HandleFunc("/admin/", func(w http.ResponseWriter, r *http.Request) {
		err := gvBackend.Render(w, http.StatusOK, "index", view.M{
			"title": "Backend title!",
		})
		if err != nil {
			_, _ = fmt.Fprintf(w, "Render index error: %v!", err)
		}
	})

	fmt.Printf("Server start on :9090")
	_ = http.ListenAndServe(":9090", nil)
}
