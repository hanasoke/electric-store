package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func testDB() {
	dsn := "root:@tcp(127.0.0.1:3306)/electric_store"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("❌ Error opening DB: %v", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Printf("❌ Error pinging DB: %v", err)
		log.Println("💡 Please check:")
		log.Println("   - Is MySQL running?")
		log.Println("   - Is the database 'electric_store' created?")
		log.Println("   - Check with: sudo /opt/lampp/lampp status")
		return
	}

	log.Println("✅ Database connection successful!")

	// Test query
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		log.Printf("❌ Error querying tables: %v", err)
		return
	}
	defer rows.Close()

	log.Println("📋 Tables in database:")
	for rows.Next() {
		var table string
		rows.Scan(&table)
		log.Printf("   - %s", table)
	}
}

// func main() {
// 	testDB()
// }
