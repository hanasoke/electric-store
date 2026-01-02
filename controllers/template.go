package controllers

import (
	"log"
	"net/http"
	"text/template"
)

var templates map[string]*template.Template

func InitTemplates() {
	templates = make(map[string]*template.Template)

	// Defire template patterns
	templatePatterns := map[string][]string{
		"products": {
			"views/templates/base.html",
			"views/products/index.html",
		},
		"categories": {
			"views/templates/base.html",
			"views/categories/index.html",
		},
	}

	// Parse all templates
	for name, files := range templatePatterns {
		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			log.Fatalf("Failed to parse template %s: %v", name, err)
		}
		templates[name] = tmpl
	}
}

func RenderTemplate(w http.ResponseWriter, name string, data interface{}) {
	tmpl, ok := templates[name]
	if !ok {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	err := tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
