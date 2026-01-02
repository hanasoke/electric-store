package categorycontroller

import (
	"electric-store/controllers"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {

	// Prepare template data
	data := struct {
		Title      string
		ActivePage string
	}{
		Title:      "Categories",
		ActivePage: "categories",
	}

	// Render template
	controllers.RenderTemplate(w, "categories", data)
}
