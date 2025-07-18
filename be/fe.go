package main

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// renderIndex renders the main page using the Go HTML template.
func renderIndex(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("fe", "index.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// serveStatic serves static files (CSS, JS) from the fe/ directory.
func serveStatic() http.Handler {
	return http.StripPrefix("/static/", http.FileServer(http.Dir("fe")))
}
