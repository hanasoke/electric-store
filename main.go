package main

import (
	"electric-store/handlers"
	"log"
	"net/http"
)

func main() {
	productHandler, err := handlers.NewProductHandler()
	if err != nil {
		log.Fatal("Failed to initialize product handler:", err)
	}

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
	http.HandleFunc("/", productHandler.Index)
	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			productHandler.CreateForm(w, r)
		} else if r.Method == "POST" {
			productHandler.Create(w, r)
		}
	})
	http.HandleFunc("/edit", productHandler.EditForm)
	http.HandleFunc("/update", productHandler.Update)
	http.HandleFunc("/delete", productHandler.Delete)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
