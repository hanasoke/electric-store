package handlers

import (
	"electric-store/config"
	"electric-store/models"
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
