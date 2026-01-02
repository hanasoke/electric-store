package productcontroller

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
		Title:      "Products",
		ActivePage: "products",
	}

	controllers.RenderTemplate(w, "products", data)
}
