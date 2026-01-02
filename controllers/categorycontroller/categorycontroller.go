package categorycontroller

import (
	"electric-store/controllers"
	"electric-store/entities"
	"electric-store/models/categorymodel"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	// Ambil data kategori dari model
	categories := categorymodel.GetAll()

	// Prepare template data
	data := struct {
		Title      string
		ActivePage string
		Categories []entities.Category // Tambahkan ini
	}{
		Title:      "Categories",
		ActivePage: "categories",
		Categories: categories,
	}

	// Render template
	controllers.RenderTemplate(w, "categories", data)
}
