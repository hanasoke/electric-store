package categorycontroller

import (
	"electric-store/controllers"
	"electric-store/entities"
	"electric-store/models/categorymodel"
	"net/http"
	"strconv"
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
	editId := ""

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
		categoryName = cookie.Value
		http.SetCookie(w, &http.Cookie{
			Name:   "category_name",
			MaxAge: -1,
		})
	}

	if cookie, err := r.Cookie("validation_error"); err == nil {
		validationError = cookie.Value
		http.SetCookie(w, &http.Cookie{
			Name:   "validation_error",
			MaxAge: -1,
		})
	}

	if cookie, err := r.Cookie("edit_id"); err == nil {
		editId = cookie.Value
		http.SetCookie(w, &http.Cookie{
			Name:   "edit_id",
			MaxAge: -1,
		})
	}

	// Prepare template data
	data := struct {
		Title           string
		ActivePage      string
		Categories      []CategoryWithNo
		SuccessMessage  string
		ErrorMessage    string
		CategoryName    string
		ValidationError string
		EditId          string
	}{
		Title:           "Categories",
		ActivePage:      "categories",
		Categories:      categoriesWithNo,
		SuccessMessage:  successMessage,
		ErrorMessage:    errorMessage,
		CategoryName:    categoryName,
		ValidationError: validationError,
		EditId:          editId,
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

	// Buat entity category
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
			errorMessage = "Category " + name + " already exists"
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
		Value:    "Category " + name + " successfully created",
		Path:     "/",
		MaxAge:   5,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}

func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/categories", http.StatusSeeOther)
		return
	}

	r.ParseForm()

	// Ambil ID dari form
	idStr := r.FormValue("id")
	name := r.FormValue("name")

	// Convert ID ke integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.SetCookie(w, &http.Cookie{
			Name:     "error",
			Value:    "Invalid category ID",
			Path:     "/",
			MaxAge:   5,
			HttpOnly: true,
		})
		http.Redirect(w, r, "/categories", http.StatusSeeOther)
		return
	}

	// Validasi
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
		http.SetCookie(w, &http.Cookie{
			Name:     "edit_id",
			Value:    idStr,
			Path:     "/",
			MaxAge:   5,
			HttpOnly: true,
		})
		http.Redirect(w, r, "/categories", http.StatusSeeOther)
		return
	}

	// Update category
	category := entities.Category{
		Name:      name,
		UpdatedAt: time.Now(),
	}

	err = categorymodel.Update(id, category)
	if err != nil {
		var errorMessage string
		if err == categorymodel.ErrDuplicateCategory {
			errorMessage = "Category " + name + " already exists"
		} else {
			errorMessage = "Failed to update category: " + err.Error()
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
		http.SetCookie(w, &http.Cookie{
			Name:     "edit_id",
			Value:    idStr,
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
		Value:    "Category " + name + " successfully updated",
		Path:     "/",
		MaxAge:   5,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/categories", http.StatusSeeOther)
		return
	}

	r.ParseForm()

	// Ambil ID dari form
	idStr := r.FormValue("id")

	// Convert ID ke integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.SetCookie(w, &http.Cookie{
			Name:     "error",
			Value:    "Invalid category ID",
			Path:     "/",
			MaxAge:   5,
			HttpOnly: true,
		})
		http.Redirect(w, r, "/categories", http.StatusSeeOther)
		return
	}

	// Hapus category
	err = categorymodel.Delete(id)
	if err != nil {
		errorMessage := "Failed to delete category: " + err.Error()
		http.SetCookie(w, &http.Cookie{
			Name:     "error",
			Value:    errorMessage,
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
		Value:    "Category successfully deleted",
		Path:     "/",
		MaxAge:   5,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}
