package main

import (
	"log"
	"net/http"
	"os"

	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/database"
)

func main() {
	dsn := os.Getenv("DATABASE_DSN")

	if dsn == "" {
        log.Fatal("❌ DATABASE_DSN is not set in environment variables")
    }

	//Call the function from our internal package
	db, err := database.Connect(dsn)
	if err != nil {
		log.Fatal("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("✅ Auth Service started successfully!")

	//Define a simple health check route
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK - Auth Service is alive"))
	})

	log.Println("🚀 Auth Service listening on :8080")

	//This rplaces 'select{}' and keeps the app runing
	if err := http.ListenAndServe(":8080", nil); err != nil{
		log.Fatalf("❌ Server failed: %v", err)
	}
}