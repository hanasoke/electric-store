package handlers

import (
	"electric-store/config"
	"electric-store/models"
	"html/template"
	"net/http"
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
