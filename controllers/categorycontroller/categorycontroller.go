package categorycontroller

import (
	"electric-store/entities"
	"net/http"
	"time"
)

func Index(w http.ResponseWriter, r *http.Request) {

}

func Store(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.FormValue("name")

	// Validasi null/empty
	if name == "" {
		http.SetCookie(w, &http.Cookie{
			Name:     "error",
			Value:    "Category name cannot be empty",
			Path:     "/",
			MaxAge:   5,
			HttpOnly: true,
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "category_name",
			Value:    name,
			Path:     "/",
			MaxAge:   5,
			HttpOnly: true,
		})
		http.Redirect(w, r, "/categories", http.StatusSeeOther)
		return
	}

	// BUat entity category
	category := entities.Category{
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
