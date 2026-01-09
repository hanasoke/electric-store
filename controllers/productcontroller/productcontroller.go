package productcontroller

import (
	"electric-store/controllers"
	"electric-store/entities"
	"electric-store/models/categorymodel"
	"electric-store/models/productmodel"
	"net/http"
	"strconv"
	"time"
)

func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		Store(w, r)
		return
	}

	// Ambil data produk dari model
	products, err := productmodel.GetAll()
	if err != nil {
		// Set error message
		http.SetCookie(w, &http.Cookie{
			Name:     "error",
			Value:    "Failed to load products: " + err.Error(),
			Path:     "/",
			MaxAge:   5,
			HttpOnly: true,
		})
		products = []entities.Product{}
	}

	// Ambil data kategori untuk dropdown
	categories, err := categorymodel.GetAllForSelect()
	if err != nil {
		categories = []entities.Category{}
	}

	// Tambahkan nomor urut
	type ProductWithNo struct {
		No      int
		Product entities.Product
	}

	var productsWithNo []ProductWithNo
	for i, product := range products {
		productsWithNo = append(productsWithNo, ProductWithNo{
			No:      i + 1,
			Product: product,
		})
	}

	// Ambil pesan flash jika ada
	successMessage := ""
	errorMessage := ""
	productName := ""
	categoryId := ""
	price := ""
	stock := ""
	description := ""
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

	if cookie, err := r.Cookie("product_name"); err == nil {
		productName = cookie.Value
		http.SetCookie(w, &http.Cookie{
			Name:   "product_name",
			MaxAge: -1,
		})
	}

	if cookie, err := r.Cookie("category_id"); err == nil {
		categoryId = cookie.Value
		http.SetCookie(w, &http.Cookie{
			Name:   "category_id",
			MaxAge: -1,
		})
	}

	if cookie, err := r.Cookie("price"); err == nil {
		price = cookie.Value
		http.SetCookie(w, &http.Cookie{
			Name:   "price",
			MaxAge: -1,
		})
	}

	if cookie, err := r.Cookie("stock"); err == nil {
		stock = cookie.Value
		http.SetCookie(w, &http.Cookie{
			Name:   "stock",
			MaxAge: -1,
		})
	}

	if cookie, err := r.Cookie("description"); err == nil {
		description = cookie.Value
		http.SetCookie(w, &http.Cookie{
			Name:   "description",
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
		Products        []ProductWithNo
		Categories      []entities.Category
		SuccessMessage  string
		ErrorMessage    string
		ProductName     string
		CategoryId      string
		Price           string
		Stock           string
		Description     string
		ValidationError string
		EditId          string
	}{
		Title:           "Products",
		ActivePage:      "products",
		Products:        productsWithNo,
		Categories:      categories,
		SuccessMessage:  successMessage,
		ErrorMessage:    errorMessage,
		ProductName:     productName,
		CategoryId:      categoryId,
		Price:           price,
		Stock:           stock,
		Description:     description,
		ValidationError: validationError,
		EditId:          editId,
	}

	controllers.RenderTemplate(w, "products", data)
}

func Store(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// Ambil data dari form
	name := r.FormValue("name")
	categoryIdStr := r.FormValue("category_id")
	priceStr := r.FormValue("price")
	stockStr := r.FormValue("stock")
	description := r.FormValue("description")

	// Convert values
	categoryId, _ := strconv.ParseUint(categoryIdStr, 10, 32)
	price, _ := strconv.ParseInt(priceStr, 10, 64)
	stock, _ := strconv.ParseInt(stockStr, 10, 64)

	// Validasi
	if name == "" {
		setProductCookies(w, name, categoryIdStr, priceStr, stockStr, description, "Product name cannot be empty")
		http.Redirect(w, r, "/products", http.StatusSeeOther)
		return
	}

	if categoryId == 0 {
		setProductCookies(w, name, categoryIdStr, priceStr, stockStr, description, "Category must be selected")
		http.Redirect(w, r, "/products", http.StatusSeeOther)
		return
	}

	if price <= 0 {
		setProductCookies(w, name, categoryIdStr, priceStr, stockStr, description, "Price must be greater than 0")
		http.Redirect(w, r, "/products", http.StatusSeeOther)
		return
	}

	if stock < 0 {
		setProductCookies(w, name, categoryIdStr, priceStr, stockStr, description, "Stock cannot be negative")
		http.Redirect(w, r, "/products", http.StatusSeeOther)
		return
	}

	// Buat entity product
	product := entities.Product{
		Name:        name,
		CategoryId:  uint(categoryId),
		Price:       price,
		Stock:       stock,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Coba create product
	err := productmodel.Create(product)
	if err != nil {
		var errorMessage string
		if err == productmodel.ErrDuplicateProduct {
			errorMessage = "Product " + name + " already exists"
		} else if err == productmodel.ErrCategoryNotFound {
			errorMessage = "Selected category not found"
		} else {
			errorMessage = "Failed to create product: " + err.Error()
		}

		setProductCookies(w, name, categoryIdStr, priceStr, stockStr, description, errorMessage)
		http.Redirect(w, r, "/products", http.StatusSeeOther)
		return
	}

	// Set success message
	http.SetCookie(w, &http.Cookie{
		Name:     "success",
		Value:    "Product " + name + " successfully created",
		Path:     "/",
		MaxAge:   5,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/products", http.StatusSeeOther)
}

func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/products", http.StatusSeeOther)
		return
	}

	r.ParseForm()

	// Ambil data dari form
	idStr := r.FormValue("id")
	name := r.FormValue("name")
	categoryIdStr := r.FormValue("category_id")
	priceStr := r.FormValue("price")
	stockStr := r.FormValue("stock")
	description := r.FormValue("description")

	// Convert ID ke integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.SetCookie(w, &http.Cookie{
			Name:     "error",
			Value:    "Invalid product ID",
			Path:     "/",
			MaxAge:   5,
			HttpOnly: true,
		})
		http.Redirect(w, r, "/products", http.StatusSeeOther)
		return
	}

	// Convert values
	categoryId, _ := strconv.ParseUint(categoryIdStr, 10, 32)
	price, _ := strconv.ParseInt(priceStr, 10, 64)
	stock, _ := strconv.ParseInt(stockStr, 10, 64)

	// Validasi
	if name == "" {
		setProductCookies(w, name, categoryIdStr, priceStr, stockStr, description, "Product name cannot be empty")
		http.SetCookie(w, &http.Cookie{
			Name:     "edit_id",
			Value:    idStr,
			Path:     "/",
			MaxAge:   5,
			HttpOnly: true,
		})
		http.Redirect(w, r, "/products", http.StatusSeeOther)
		return
	}

	if categoryId == 0 {
		setProductCookies(w, name, categoryIdStr, priceStr, stockStr, description, "Category must be selected")
		http.SetCookie(w, &http.Cookie{
			Name:     "edit_id",
			Value:    idStr,
			Path:     "/",
			MaxAge:   5,
			HttpOnly: true,
		})
		http.Redirect(w, r, "/products", http.StatusSeeOther)
		return
	}

	if price <= 0 {
		setProductCookies(w, name, categoryIdStr, priceStr, stockStr, description, "Price must be greater than 0")
		http.SetCookie(w, &http.Cookie{
			Name:     "edit_id",
			Value:    idStr,
			Path:     "/",
			MaxAge:   5,
			HttpOnly: true,
		})
		http.Redirect(w, r, "/products", http.StatusSeeOther)
		return
	}

	if stock < 0 {
		setProductCookies(w, name, categoryIdStr, priceStr, stockStr, description, "Stock cannot be negative")
		http.SetCookie(w, &http.Cookie{
			Name:     "edit_id",
			Value:    idStr,
			Path:     "/",
			MaxAge:   5,
			HttpOnly: true,
		})
		http.Redirect(w, r, "/products", http.StatusSeeOther)
		return
	}

	// Update product
	product := entities.Product{
		Name:        name,
		CategoryId:  uint(categoryId),
		Price:       price,
		Stock:       stock,
		Description: description,
		UpdatedAt:   time.Now(),
	}

	err = productmodel.Update(id, product)
	if err != nil {
		var errorMessage string
		if err == productmodel.ErrDuplicateProduct {
			errorMessage = "Product " + name + " already exists"
		} else if err == productmodel.ErrCategoryNotFound {
			errorMessage = "Selected category not found"
		} else {
			errorMessage = "Failed to update product: " + err.Error()
		}

		setProductCookies(w, name, categoryIdStr, priceStr, stockStr, description, errorMessage)
		http.SetCookie(w, &http.Cookie{
			Name:     "edit_id",
			Value:    idStr,
			Path:     "/",
			MaxAge:   5,
			HttpOnly: true,
		})
		http.Redirect(w, r, "/products", http.StatusSeeOther)
		return
	}

	// Set success message
	http.SetCookie(w, &http.Cookie{
		Name:     "success",
		Value:    "Product " + name + " successfully updated",
		Path:     "/",
		MaxAge:   5,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/products", http.StatusSeeOther)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/products", http.StatusSeeOther)
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
			Value:    "Invalid product ID",
			Path:     "/",
			MaxAge:   5,
			HttpOnly: true,
		})
		http.Redirect(w, r, "/products", http.StatusSeeOther)
		return
	}

	// Hapus product
	err = productmodel.Delete(id)
	if err != nil {
		errorMessage := "Failed to delete product: " + err.Error()
		http.SetCookie(w, &http.Cookie{
			Name:     "error",
			Value:    errorMessage,
			Path:     "/",
			MaxAge:   5,
			HttpOnly: true,
		})
		http.Redirect(w, r, "/products", http.StatusSeeOther)
		return
	}

	// Set success message
	http.SetCookie(w, &http.Cookie{
		Name:     "success",
		Value:    "Product successfully deleted",
		Path:     "/",
		MaxAge:   5,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/products", http.StatusSeeOther)
}

// Helper function untuk set cookies
func setProductCookies(w http.ResponseWriter, name, categoryId, price, stock, description, error string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "error",
		Value:    error,
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
		Value:    categoryId,
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
