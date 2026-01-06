package productcontroller

import (
	"electric-store/controllers"
	"electric-store/entities"
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Ambil kategori untuk dropdown
	categories, err := productmodel.GetCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
	categoryID := ""
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
		categoryID = cookie.Value
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
		CategoryID      string
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
		CategoryID:      categoryID,
		Price:           price,
		Stock:           stock,
		Description:     description,
		ValidationError: validationError,
		EditId:          editId,
	}

	controllers.RenderTemplate(w, "products", data)
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
	price, err := strconv.ParseInt(priceStr, 10, 64)
	if err != nil || price <= 0 {
		setErrorCookies(w, "Price must be a positive integer", name, categoryIDStr, priceStr, stockStr, description)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Parse stock
	stock, err := strconv.ParseInt(stockStr, 10, 64)
	if err != nil || stock < 0 {
		setErrorCookies(w, "Stock must be a non-negative integer", name, categoryIDStr, priceStr, stockStr, description)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Buat entity product
	product := entities.Product{
		Name:        name,
		Category:    entities.Category{Id: uint(categoryID)},
		Price:       price,
		Stock:       stock,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Coba create product
	err = productmodel.Create(product)
	if err != nil {
		var errorMessage string
		if err == productmodel.ErrDuplicateProduct {
			errorMessage = "Product " + name + " already exists"
		} else {
			errorMessage = "Failed to create product: " + err.Error()
		}

		setErrorCookies(w, errorMessage, name, categoryIDStr, priceStr, stockStr, description)
		http.Redirect(w, r, "/", http.StatusSeeOther)
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

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.ParseForm()

	idStr := r.FormValue("id")
	name := r.FormValue("name")
	categoryIDStr := r.FormValue("category_id")
	priceStr := r.FormValue("price")
	stockStr := r.FormValue("stock")
	description := r.FormValue("description")

	// Convert ID ke integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		setErrorCookies(w, "Invalid product ID", name, categoryIDStr, priceStr, stockStr, description)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Validasi input
	if name == "" {
		setErrorCookiesWithEdit(w, "Product name cannot be empty", name, categoryIDStr, priceStr, stockStr, description, idStr)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Parse category ID
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil || categoryID == 0 {
		setErrorCookiesWithEdit(w, "Please select a valid category", name, categoryIDStr, priceStr, stockStr, description, idStr)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Parse price
	price, err := strconv.ParseInt(priceStr, 10, 64)
	if err != nil || price <= 0 {
		setErrorCookiesWithEdit(w, "Price must be a positive integer", name, categoryIDStr, priceStr, stockStr, description, idStr)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Parse stock
	stock, err := strconv.ParseInt(stockStr, 10, 64)
	if err != nil || stock < 0 {
		setErrorCookiesWithEdit(w, "Stock must be a non-negative integer", name, categoryIDStr, priceStr, stockStr, description, idStr)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Update product
	product := entities.Product{
		Name:        name,
		Category:    entities.Category{Id: uint(categoryID)},
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
		} else {
			errorMessage = "Failed to update product: " + err.Error()
		}

		setErrorCookiesWithEdit(w, errorMessage, name, categoryIDStr, priceStr, stockStr, description, idStr)
		http.Redirect(w, r, "/", http.StatusSeeOther)
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

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.ParseForm()

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
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
		http.Redirect(w, r, "/", http.StatusSeeOther)
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

	http.Redirect(w, r, "/", http.StatusSeeOther)
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
