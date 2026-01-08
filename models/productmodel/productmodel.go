package productmodel

import (
	"database/sql"
	"electric-store/config"
	"electric-store/entities"
	"errors"
)

var ErrDuplicateProduct = errors.New("brand already exists")

func IsProductExists(name string) (bool, error) {
	var id uint
	err := config.DB.QueryRow(
		"SELECT id From products WHERE name = ? LIMIT 1",
		name,
	).Scan(&id)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func IsProductExceptID(name string, id int) (bool, error) {
	var productID int
	err := config.DB.QueryRow(`
		SELECT id FROM products 
		WHERE name = ? AND id != ? 
		LIMIT 1`,
		name, id,
	).Scan(&productID)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func GetAll() ([]entities.Product, error) {
	rows, err := config.DB.Query(`
        SELECT 
            p.id, 
            p.name, 
            p.category_id,
            p.price, 
            p.stock, 
            p.description,
            p.created_at,
            p.updated_at,
            c.id as category_id,
            c.name as category_name
        FROM products p
        LEFT JOIN categories c ON p.category_id = c.id
        ORDER BY p.id DESC`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []entities.Product

	for rows.Next() {
		var product entities.Product
		var categoryID sql.NullInt64
		var categoryName sql.NullString

		if err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Category.Id,
			&product.Price,
			&product.Stock,
			&product.Description,
			&product.CreatedAt,
			&product.UpdatedAt,
			&categoryID,
			&categoryName,
		); err != nil {
			return nil, err
		}

		// Set category data
		if categoryID.Valid && categoryName.Valid {
			product.Category.Id = uint(categoryID.Int64)
			product.Category.Name = categoryName.String
		}

		products = append(products, product)
	}

	return products, nil
}
