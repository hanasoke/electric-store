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

func (pm *ProductModel) GetAll() ([]Product, error) {
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

func (pm *ProductModel) GetByID(id int) (*Product, error) {
	query := `SELECT id, name, category, price, stock, description, created_at FROM products WHERE id = ?`

	row := pm.DB.QueryRow(query, id)

	var p Product
	err := row.Scan(&p.ID, &p.Name, &p.Category, &p.Price, &p.Stock, &p.Description, &p.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (pm *ProductModel) Create(product *Product) error {
	query := `INSERT INTO products (name, category, price, stock, description) VALUES (?, ?, ?, ?, ?)`

	result, err := pm.DB.Exec(query, product.Name, product.Category, product.Price, product.Stock, product.Description)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	product.ID = int(id)
	return nil
}

func (pm *ProductModel) Update(product *Product) error {
	query := `UPDATE products SET name=?, category=?, price=?, stock=?, description=? WHERE id=?`

	_, err := pm.DB.Exec(query, product.Name, product.Category, product.Price, product.Stock, product.Description, product.ID)
	return err
}
