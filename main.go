package main

import (
	"electric-store/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	// Setup logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("🚀 Starting Electric Store Application...")

	// Initialize product handler
	productHandler, err := handlers.NewProductHandler()
	if err != nil {
		log.Printf("❌ Failed to initialize product handler: %v", err)
		log.Println("💡 Please check:")
		log.Println("   1. Is XAMPP MySQL running?")
		log.Println("   2. Is the database 'electric_store' created?")
		log.Println("   3. Are the tables created?")
		os.Exit(1)
	}

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes dengan logging
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("🌐 %s %s", r.Method, r.URL.Path)
		productHandler.Index(w, r)
	})

	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("🌐 %s %s", r.Method, r.URL.Path)
		if r.Method == "GET" {
			productHandler.CreateForm(w, r)
		} else if r.Method == "POST" {
			productHandler.Create(w, r)
		}
	})

	http.HandleFunc("/edit", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("🌐 %s %s", r.Method, r.URL.Path)
		productHandler.EditForm(w, r)
	})

	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("🌐 %s %s", r.Method, r.URL.Path)
		productHandler.Update(w, r)
	})

	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("🌐 %s %s", r.Method, r.URL.Path)
		productHandler.Delete(w, r)
	})

	log.Println("✅ Server started successfully on http://localhost:8080")
	log.Println("📊 Access the application at: http://localhost:8080")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("❌ Server failed to start: %v", err)
		os.Exit(1)
	}
}
