package models

import (
	"database/sql"
	"time"
)

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type ProductModel struct {
	DB *sql.DB
}

func (pm *ProductModel) getAll() ([]Product, error) {
	query := `SELECT id, name, category, price, stock, description, created_at 
              FROM products ORDER BY created_at DESC`

	rows, err := pm.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Category, &p.Price, &p.Stock, &p.Description, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}
