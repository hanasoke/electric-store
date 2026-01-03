package categorycontroller

import (
	"electric-store/controllers"
	"electric-store/entities"
	"electric-store/models/categorymodel"
	"net/http"
	"time"
)

func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		Store(w, r)
		return
	}

	// Ambil data kategori dari model
	categories := categorymodel.GetAll()

	// Tambahkan nomor urut
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

	// Ambil pesan flash jika ada
	successMessage := ""
	errorMessage := ""
	categoryName := ""
	validationError := ""

	// Cek cookie untuk pesan
	if cookie, err := r.Cookie("success"); err == nil {
		successMessage = cookie.Value
		http.SetCookie(w, &http.Cookie{
			Name:   "success",
			MaxAge: -1,
		})
	}

	if cookie, err := r.Cookie("error"); err == nil {
		errorMessage = cookie.Value
		http.SetCookie(w, &http.Cookie{
			Name:   "error",
			MaxAge: -1,
		})
	}

	if cookie, err := r.Cookie("category_name"); err == nil {
		errorMessage = cookie.Value
		http.SetCookie(w, &http.Cookie{
			Name:   "category_name",
			MaxAge: -1,
		})
	}

	if cookie, err := r.Cookie("validation_error"); err == nil {
		errorMessage = cookie.Value
		http.SetCookie(w, &http.Cookie{
			Name:   "validation_error",
			MaxAge: -1,
		})
	}

	// Prepare template data
	data := struct {
		Title           string
		ActivePage      string
		Categories      []CategoryWithNo
		successMessage  string
		ErrorMessage    string
		CategoryName    string
		validationError string
	}{
		Title:           "Categories",
		ActivePage:      "categories",
		Categories:      categoriesWithNo,
		successMessage:  successMessage,
		ErrorMessage:    errorMessage,
		CategoryName:    categoryName,
		validationError: validationError,
	}

	controllers.RenderTemplate(w, "categories", data)
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

	// Coba create category
	err := categorymodel.Create(category)
	if err != nil {
		var errorMessage string
		if err == categorymodel.ErrDuplicateCategory {
			errorMessage = "Category '" + name + "' already exists"
		} else {
			errorMessage = "Failed to create category: " + err.Error()
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "error",
			Value:    errorMessage,
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

	// Set success message
	http.SetCookie(w, &http.Cookie{
		Name:     "success",
		Value:    "Category '" + name + "' successfully created",
		Path:     "/",
		MaxAge:   5,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}
