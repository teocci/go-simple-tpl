// Package view
// Created by RTT.
// Author: teocci@yandex.com on 2022-12월-24
package view

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"testing"
	"time"
)

var cases = []struct {
	Name string
	Data M
	Out  string
}{
	{
		Name: "echo.tpl",
		Data: M{"name": "GoView"},
		Out:  "$GoView",
	},
	{
		Name: "include",
		Data: M{"name": "GoView"},
		Out:  "<v>IncGoView</v>",
	},
	{
		Name: "index",
		Data: M{},
		Out:  "<v>Index</v>",
	},
	{
		Name: "sum",
		Data: M{
			"sum": func(a int, b int) int {
				return a + b
			},
			"a": 1,
			"b": 2,
		},
		Out: "<v>3</v>",
	},
}

func TestViewEngine_RenderWriter(t *testing.T) {
	gv := New(Config{
		Root:      "../../_examples/test",
		Extension: ".tpl",
		Master:    "layouts/master",
		Partials:  []string{},
		Funcs: template.FuncMap{
			"echo": func(v string) string {
				return "$" + v
			},
		},
		DisableCache: true,
	})

	for _, v := range cases {
		buff := new(bytes.Buffer)
		err := gv.RenderWriter(buff, v.Name, v.Data)
		if err != nil {
			t.Errorf("name: %v, data: %v, error: %v", v.Name, v.Data, err)
			continue
		}
		val := string(buff.Bytes())
		if val != v.Out {
			t.Errorf("actual: %v, expect: %v", val, v.Out)
		}
	}
}

func ExampleDefault() {
	//Project structure:
	//|-- app/views/
	//   |--- index.html
	//   |--- page.html
	//   |-- layouts/
	//	   |--- footer.html
	//	   |--- master.html

	//render index use `index` without `.html` extension, that will render with master layout.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := Render(w, http.StatusOK, "index", M{
			"title": "Index title!",
			"add": func(a int, b int) int {
				return a + b
			},
		})
		if err != nil {
			_, _ = fmt.Fprintf(w, "Render index error: %v!", err)
		}

	})

	//render page use `page.html` with '.html' will only file template without master layout.
	http.HandleFunc("/page", func(w http.ResponseWriter, r *http.Request) {
		err := Render(w, http.StatusOK, "page.html", M{"title": "Page file title!!"})
		if err != nil {
			_, _ = fmt.Fprintf(w, "Render page.html error: %v!", err)
		}
	})

	fmt.Println("Listening and serving HTTP on :9090")
	_ = http.ListenAndServe(":9090", nil)
}

func ExampleNew() {
	//Project structure:
	//
	//|-- app/views/
	//    |--- index.tpl
	//    |--- page.tpl
	//    |-- layouts/
	//        |--- footer.tpl
	//        |--- head.tpl
	//        |--- master.tpl
	//    |-- partials/
	//        |--- ad.tpl

	gv := New(Config{
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

	// Set new instance
	Use(gv)

	// Render index use `index` without `.tpl` extension, that will render with master layout.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := Render(w, http.StatusOK, "index", M{
			"title": "Index title!",
			"add": func(a int, b int) int {
				return a + b
			},
		})
		if err != nil {
			_, _ = fmt.Fprintf(w, "Render index error: %v!", err)
		}

	})

	// Render page use `page.tpl` with '.tpl' will only file template without master layout.
	http.HandleFunc("/page", func(w http.ResponseWriter, r *http.Request) {
		err := Render(w, http.StatusOK, "page.tpl", M{"title": "Page file title!!"})
		if err != nil {
			_, _ = fmt.Fprintf(w, "Render page.html error: %v!", err)
		}
	})

	fmt.Println("Listening and serving HTTP on :9090")
	_ = http.ListenAndServe(":9090", nil)
}
