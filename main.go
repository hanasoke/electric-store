package main

import (
	"electric-store/config"
	"electric-store/controllers"
	"electric-store/controllers/categorycontroller"
	"electric-store/controllers/productcontroller"
	"log"
	"net/http"
)

func main() {
	config.ConnectDB()

	// Initialize templates
	controllers.InitTemplates()

	// Serve status files (if any)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", productcontroller.Index)
	http.HandleFunc("/products", productcontroller.Index)
	http.HandleFunc("/products/update", productcontroller.Update)
	http.HandleFunc("/products/delete", productcontroller.Delete)

	http.HandleFunc("/categories", categorycontroller.Index)
	http.HandleFunc("/categories/update", categorycontroller.Update)
	http.HandleFunc("/categories/delete", categorycontroller.Delete)

	log.Println("Server running on port 7050")
	http.ListenAndServe(":7050", nil)
}
