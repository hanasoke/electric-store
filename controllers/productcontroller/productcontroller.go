package productcontroller

import (
	"net/http"
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
