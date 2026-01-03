package categorycontroller

import (
	"electric-store/controllers"
	"electric-store/entities"
	"electric-store/models/categorymodel"
	"net/http"
	"strings"
	"time"
)

func Index(w http.ResponseWriter, r *http.Request) {
	// Handle POST request untuk tambah kategori
	if r.Method == http.MethodPost {

	}

	// Handle GET request untuk menampilkan data
	displayCategories(w, r, nil, "", "")
}

func handleAddCategory(w http.ResponseWriter, r *http.Request) {
	// Parse form
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.Form.Get("name"))

	// Validasi input
	var errors []string
	var nameValue string

	// Validasi null/empty
	if name == "" {
		errors = append(errors, "Category name is required")
	} else {
		nameValue = name

		// Validasi duplikat
		exists, err := categorymodel.IsCategoryExists(name)
		if err != nil {
			errors = append(errors, "Error checking category existence")
		} else if exists {
			errors = append(errors, "Category already exists")
		}
	}

	// Jika ada error, tampilkan kembali dengan pesan error
	if len(errors) > 0 {
		displayCategories(w, r, errors, nameValue, "is-invalid")
		return
	}

	// Buat kategori baru
	category := entities.Category{
		Name:      name,
		CreatedAt: time.Now(),
	}

	// Simpan ke database
	err = categorymodel.Create(category)
	if err != nil {
		errors = append(errors, "Failed to save category")
		displayCategories(w, r, errors, nameValue, "is-invalid")
		return
	}

	// Redirect untuk menghindari resubmission
	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}

func displayCategories(w http.ResponseWriter, r *http.Request, errors []string, nameValue string, inputClass string) {
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

	// Tentukan modal state
	showModal := false
	if r.Method == http.MethodGet && r.URL.Query().Get("add") == "true" {
		showModal = true
	} else if r.Method == http.MethodPost {
		showModal = true
	}

	// Prepare template data
	data := struct {
		Title      string
		ActivePage string
		Categories []CategoryWithNo
		Errors     []string
		NameValue  string
		InputClass string
		ShowModal  bool
		HasErrors  bool
	}{
		Title:      "Categories",
		ActivePage: "categories",
		Categories: categoriesWithNo,
		Errors:     errors,
		NameValue:  nameValue,
		InputClass: inputClass,
		ShowModal:  showModal,
		HasErrors:  len(errors) > 0,
	}

	controllers.RenderTemplate(w, "categories", data)
}
