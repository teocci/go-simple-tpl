// Package view
// Created by RTT.
// Author: teocci@yandex.com on 2022-12월-24
package view

import "net/http"

var instance *Engine

// Use setting default instance engine
func Use(engine *Engine) {
	instance = engine
}

// Render render view template with default instance
func Render(w http.ResponseWriter, status int, name string, data interface{}) error {
	if instance == nil {
		instance = Default()
		//return fmt.Errorf("instance not yet initialized, please call Init() first before Render()")
	}
	return instance.Render(w, status, name, data)
}
