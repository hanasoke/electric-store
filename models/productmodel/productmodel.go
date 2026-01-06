package productmodel

import (
	"database/sql"
	"electric-store/config"
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
