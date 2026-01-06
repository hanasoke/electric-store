package productcontroller

import (
	"electric-store/models/productmodel"
	"net/http"
	"strconv"
)

func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		Store(w, r)
		return
	}
}

func Store(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.ParseForm()

	name := r.FormValue("name")
	categoryIDStr := r.FormValue("category_id")
	priceStr := r.FormValue("price")
	stockStr := r.FormValue("stock")
	description := r.FormValue("description")

	// Validasi input
	if name == "" {
		setErrorCookies(w, "Product name cannot be empty", name, categoryIDStr, priceStr, stockStr, description)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Parse category ID
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil || categoryID == 0 {
		setErrorCookies(w, "Please select a valid category", name, categoryIDStr, priceStr, stockStr, description)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Parse price
	price, err := strconv.ParseInt(priceStr, 10, 32)
	if err != nil || price == 0 {
		setErrorCookies(w, "Price must be a positive integer", name, categoryIDStr, priceStr, stockStr, description)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Parse stock
	stock, err := strconv.ParseInt(priceStr, 10, 32)
	if err != nil || stock < 0 {
		setErrorCookies(w, "Price must be a positive integer", name, categoryIDStr, priceStr, stockStr, description)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Coba create product
	err = productmodel.Create(product)

}

// Helper function untuk set error cookies
func setErrorCookies(w http.ResponseWriter, errorMessage, name, categoryID, price, stock, description string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "error",
		Value:    errorMessage,
		Path:     "/",
		MaxAge:   5,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "product_name",
		Value:    name,
		Path:     "/",
		MaxAge:   5,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "category_id",
		Value:    categoryID,
		Path:     "/",
		MaxAge:   5,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "price",
		Value:    price,
		Path:     "/",
		MaxAge:   5,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "stock",
		Value:    stock,
		Path:     "/",
		MaxAge:   5,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "description",
		Value:    description,
		Path:     "/",
		MaxAge:   5,
		HttpOnly: true,
	})
}

func setErrorCookiesWithEdit(w http.ResponseWriter, errorMessage, name, categoryID, price, stock, description, editId string) {
	setErrorCookies(w, errorMessage, name, categoryID, price, stock, description)
	http.SetCookie(w, &http.Cookie{
		Name:     "edit_id",
		Value:    editId,
		Path:     "/",
		MaxAge:   5,
		HttpOnly: true,
	})
}
