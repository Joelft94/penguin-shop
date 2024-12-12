package main

import (
	"log"
	"net/http"
	"os"

	"penguin-store/database"
	"penguin-store/handlers"
)

func main() {
	// Connect to MongoDB
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	database.ConnectDB(mongoURI)

	// Static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
	http.HandleFunc("/", handlers.HandleProducts) // Changed to show products on home page
	http.HandleFunc("/products", handlers.HandleProducts)
	http.HandleFunc("/order", handlers.HandleOrder)
	http.HandleFunc("/order-success", handlers.HandleOrderSuccess)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
