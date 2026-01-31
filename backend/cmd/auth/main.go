package main

import (
	"log"
	"os"

	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/database"
)

func main() {
	dsn := os.Getenv("DATABASE_DSN")

	//Call the function from our internal package
	db, err := database.Connect(dsn)
	if err != nil {
		log.Fatal("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("✅ Auth Service started successfully!")

	// Keep service alive
	select {}
}