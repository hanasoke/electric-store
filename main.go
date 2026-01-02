package main

import (
	"electric-store/config"
	"electric-store/controllers/productcontroller"
	"log"
	"net/http"
)

func main() {
	config.ConnectDB()

	http.HandleFunc("/", productcontroller.Index)

	log.Println("Server running on port 7050")
	http.ListenAndServe(":7050", nil)
}
