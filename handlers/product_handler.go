package handlers

import (
	"electric-store/config"
	"electric-store/models"
	"html/template"
	"net/http"
	"strconv"
)

type ProductHandler struct {
	model *models.ProductModel
}

func NewProductHandler() (*ProductHandler, error) {
	db, err := config.DBConnection()
	if err != nil {
		return nil, err
	}

	return &ProductHandler{
		model: &models.ProductModel{DB: db},
	}, nil
}

func (ph *ProductHandler) Index(w http.ResponseWriter, r *http.Request) {
	products, err := ph.model.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, products)
}

func (ph *ProductHandler) CreateForm(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/create.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func (ph *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	stock, _ := strconv.Atoi(r.FormValue("stock"))

	product := &models.Product{
		Name:        r.FormValue("name"),
		Category:    r.FormValue("category"),
		Price:       price,
		Stock:       stock,
		Description: r.FormValue("description"),
	}

	err = ph.model.Create(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (ph *ProductHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := ph.model.GetByID(id)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("templates/edit.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, product)
}

func (ph *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(r.FormValue("id"))
	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	stock, _ := strconv.Atoi(r.FormValue("stock"))

	product := &models.Product{
		ID:          id,
		Name:        r.FormValue("name"),
		Category:    r.FormValue("category"),
		Price:       price,
		Stock:       stock,
		Description: r.FormValue("description"),
	}

	err = ph.model.Update(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (ph *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	err = ph.model.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
