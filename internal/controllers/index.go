package controllers

import (
	"html/template"
	"net/http"
)

func ShowIndex(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFiles("internal/views/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
