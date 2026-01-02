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

func GetAll() []entities.Category {
	rows, err := config.DB.Query(`SELECT id, name, created_at, updated_at FROM categories ORDER BY id DESC`)
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
			&category.UpdatedAt,
		); err != nil {
			panic(err)
		}
		categories = append(categories, category)
	}

	return categories
}
