package categorymodel

import (
	"database/sql"
	"electric-store/config"
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
