package main

import (
	"electric-store/config"
	"log"
	"net/http"
)

func main() {
	config.ConnectDB()

	log.Println("Server running on port 7050")
	http.ListenAndServe(":7050", nil)
}
