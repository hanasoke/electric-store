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

	// Tambahkan nomor urut (jika diperlukan)
	type CategoryWithNo struct {
		No       int
		Category entities.Category
	}

	var categoriesWithNo []CategoryWithNo
	for i, category := range categories {
		categoriesWithNo = append(categoriesWithNo, CategoryWithNo{
			No:       i + 1,
			Category: category,
		})
	}

	// Prepare template data
	data := struct {
		Title      string
		ActivePage string
		Categories []CategoryWithNo
	}{
		Title:      "Categories",
		ActivePage: "categories",
		Categories: categoriesWithNo,
	}

	controllers.RenderTemplate(w, "categories", data)
}

func Add(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

	}

	if r.Method == http.MethodPost {

	}
}
