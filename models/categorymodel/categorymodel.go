package categorymodel

import (
	"database/sql"
	"electric-store/config"
	"electric-store/entities"
	"errors"
)

var ErrDuplicateCategory = errors.New("category already exists")

func IsCategoryExists(name string) (bool, error) {
	var id uint
	err := config.DB.QueryRow(
		"SELECT id FROM categories WHERE name = ? LIMIT 1",
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

func IsCategoryExceptID(name string, id int) (bool, error) {
	var Id int
	err := config.DB.QueryRow(`
		SELECT id FROM categories 
		WHERE name = ? AND id != ?
		LIMIT 1`,
		name, id,
	).Scan(&Id)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func GetAll() []entities.Category {
	rows, err := config.DB.Query(`SELECT id, name, created_at FROM categories ORDER BY id DESC`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var categories []entities.Category

	for rows.Next() {
		var category entities.Category
		if err := rows.Scan(
			&category.Id,
			&category.Name,
			&category.CreatedAt,
		); err != nil {
			panic(err)
		}
		categories = append(categories, category)
	}

	return categories
}

func Create(category entities.Category) error {
	// Validasi input
	if category.Name == "" {
		return errors.New("category name cannot be empty")
	}

	// Cek duplikat
	exists, err := IsCategoryExists(category.Name)
	if err != nil {
		return err
	}

	if exists {
		return ErrDuplicateCategory
	}

	// Insert ke database
	_, err = config.DB.Exec(`
		INSERT INTO categories (
			name, 
			created_at
		) VALUES (?, ?)`,
		category.Name,
		category.CreatedAt,
	)

	return err
}

func Detail(id int) entities.Category {
	row := config.DB.QueryRow(`SELECT id, name FROM categories WHERE id = ?`, id)

	var category entities.Category
	if err := row.Scan(&category.Id, &category.Name); err != nil {
		panic(err.Error())
	}

	return category
}

func Update(id int, category entities.Category) error {
	exists, err := IsCategoryExceptID(category.Name, id)
	if err != nil {
		return err
	}

	if exists {
		return ErrDuplicateCategory
	}

	_, err = config.DB.Exec(`
		UPDATE categories 
		SET name = ?, updated_at = ?
		WHERE id = ?`,
		category.Name,
		category.UpdatedAt,
		id,
	)

	return err
}
