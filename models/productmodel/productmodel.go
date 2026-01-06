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

func GetById(id int) (entities.Product, error) {
	var product entities.Product
	var categoryID sql.NullInt64
	var categoryName sql.NullString

	err := config.DB.QueryRow(`
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
        WHERE p.id = ?`, id,
	).Scan(
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
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return product, errors.New("product not found")
		}
		return product, err
	}

	// Set category data
	if categoryID.Valid && categoryName.Valid {
		product.Category.Id = uint(categoryID.Int64)
		product.Category.Name = categoryName.String
	}

	return product, nil
}

func Create(product entities.Product) error {
	// Validasi input
	if product.Name == "" {
		return errors.New("product name cannot be empty")
	}

	if product.Price <= 0 {
		return errors.New("price must be greater than 0")
	}

	if product.Stock < 0 {
		return errors.New("stock cannot be negative")
	}

	// Cek duplikat
	exists, err := IsProductExists(product.Name)
	if err != nil {
		return err
	}

	if exists {
		return ErrDuplicateProduct
	}

	// Insert ke database
	_, err = config.DB.Exec(`
        INSERT INTO products (
            name, 
            category_id,
            price,
            stock,
            description,
            created_at
        ) VALUES (?, ?, ?, ?, ?, ?)`,
		product.Name,
		product.Category.Id,
		product.Price,
		product.Stock,
		product.Description,
		product.CreatedAt,
	)

	return err
}

func Update(id int, product entities.Product) error {
	// Validasi input
	if product.Name == "" {
		return errors.New("product name cannot be empty")
	}

	if product.Price <= 0 {
		return errors.New("price must be greater than 0")
	}

	if product.Stock < 0 {
		return errors.New("stock cannot be negative")
	}

	exists, err := IsProductExceptID(product.Name, id)
	if err != nil {
		return err
	}

	if exists {
		return ErrDuplicateProduct
	}

	_, err = config.DB.Exec(`
        UPDATE products 
        SET name = ?, 
            category_id = ?,
            price = ?,
            stock = ?,
            description = ?,
            updated_at = ?
        WHERE id = ?`,
		product.Name,
		product.Category.Id,
		product.Price,
		product.Stock,
		product.Description,
		product.UpdatedAt,
		id,
	)

	return err
}

func Delete(id int) error {
	_, err := config.DB.Exec(`
		DELETE FROM products
		WHERE id = ?`, id)

	return err
}

// GetCategories untuk dropdown
func GetCategories() ([]entities.Category, error) {
	rows, err := config.DB.Query(`
		SELECT id, name
		FROM categories 
		ORDER BY name`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []entities.Category

	for rows.Next() {
		var category entities.Category
		if err := rows.Scan(
			&category.Id,
			&category.Name,
		); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}
