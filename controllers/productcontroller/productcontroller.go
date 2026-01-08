package productcontroller

import (
	"electric-store/controllers"
	"electric-store/entities"
	"electric-store/models/productmodel"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	// Get all products from database
	products, err := productmodel.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// prepare template data
	data := struct {
		Title      string
		ActivePage string
		Products   []entities.Product
	}{
		Title:      "Products",
		ActivePage: "products",
		Products:   products,
	}

	controllers.RenderTemplate(w, "products", data)
}
