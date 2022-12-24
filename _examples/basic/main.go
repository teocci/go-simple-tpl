// Package webserver
// Created by RTT.
// Author: teocci@yandex.com on 2022-Apr-26

package main

import (
	"fmt"
	"net/http"

	"github.com/teocci/go-simple-tpl/src/view"
)

func main() {
	// Render index use `index` without `.html` extension, that will render with master layout.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := view.Render(w, http.StatusOK, "index", view.M{
			"title": "Index title!",
			"add": func(a int, b int) int {
				return a + b
			},
		})
		if err != nil {
			_, _ = fmt.Fprintf(w, "Render index error: %v!", err)
		}

	})

	// Render page use `page.html` with '.html' will only file template without master layout.
	http.HandleFunc("/page", func(w http.ResponseWriter, r *http.Request) {
		err := view.Render(w, http.StatusOK, "page.html", view.M{"title": "Page file title!!"})
		if err != nil {
			_, _ = fmt.Fprintf(w, "Render page.html error: %v!", err)
		}
	})

	fmt.Println("Listening and serving HTTP on :9090")
	_ = http.ListenAndServe(":9090", nil)
}
