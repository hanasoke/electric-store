package productcontroller

import (
	"electric-store/controllers"
	"electric-store/entities"
	"electric-store/models/productmodel"
	"net/http"
	"strconv"
)

func Index(w http.ResponseWriter, r *http.Request) {
	// Get all products from database
	products, err := productmodel.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get all categories for dropdown
	categories, err := productmodel.GetAllCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// prepare template data
	data := struct {
		Title      string
		ActivePage string
		Products   []entities.Product
		Categories []entities.Category
	}{
		Title:      "Products",
		ActivePage: "products",
		Products:   products,
		Categories: categories,
	}

	controllers.RenderTemplate(w, "products", data)
}

func Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get form values
	name := r.FormValue("name")
	categoryIDStr := r.FormValue("category_id")
	priceStr := r.FormValue("price")
	stockStr := r.FormValue("stock")
	description := r.FormValue("description")

	// Validate required fields
	if name == "" {
		http.Error(w, "Product name is required", http.StatusBadRequest)
		return
	}

	// Convert string to appropriate types
	categoryID, _ := strconv.ParseUint(categoryIDStr, 10, 64)
	price, err := strconv.ParseInt(priceStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid price", http.StatusBadRequest)
		return
	}

	stock, err := strconv.ParseInt(stockStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid stock", http.StatusBadRequest)
		return
	}

	// Create product entity
	product := &entities.Product{
		Name:        name,
		Category:    entities.Category{Id: uint(categoryID)},
		Price:       price,
		Stock:       stock,
		Description: description,
	}

	// Save to database
	err = productmodel.Create(product)
	if err != nil {
		if err == productmodel.ErrDuplicateProduct {
			http.Error(w, "Product already exists", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to products page
	http.Redirect(w, r, "/products", http.StatusSeeOther)
}
