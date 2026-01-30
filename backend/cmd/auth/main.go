package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq" // Make sure to run 'go get github.com/lib/pq' in /backend
)

func main() {
	// 1. Get the connection string from Environment Variables
	// Note: In Docker, the host is "db", not "localhost"
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://Rediet:postgres@db:5432/digital_contract_db?sslmode=disable"
	}

	var db *sql.DB
	var err error

	// 2. Retry logic: Wait for DB to be fully ready
	for i := 0; i < 5; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
		}

		if err == nil {
			fmt.Println("✅ SUCCESS: Auth Service connected to Postgres!")
			break
		}

		fmt.Printf("⏳ DB not ready yet (attempt %d/5)... waiting\n", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("❌ FAILED: Could not connect to DB after retries:", err)
	}

	// Keep service alive
	select {}
}