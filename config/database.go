package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Sesuaikan dengan konfigurasi XAMPP Anda
const (
	username = "root"
	password = "" // Biasanya kosong di XAMPP
	hostname = "127.0.0.1:3306"
	dbname   = "electric_store"
)

func DBConnection() (*sql.DB, error) {
	// Format: username:password@tcp(host:port)/database
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbname)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return nil, err
	}

	// Test koneksi
	err = db.Ping()
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return nil, err
	}

	log.Println("✅ Database connected successfully")
	return db, nil
}
