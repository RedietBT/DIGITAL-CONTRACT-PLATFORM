package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

//Connect initalizes the database connection pool
func Connect(dsn string) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for i := 0; i < 5; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			err = db.Ping()
		}

		if err == nil {
			return db, nil
		}
		fmt.Printf("⏳ Connecting to DB (attempt %d/5)...\n", i+1)
		time.Sleep(2 * time.Second)
	}

	return  nil, err
}

//RunMigrations handles the "Content" -Synchronizing your SQL tables
func RunMigrations(db *sql.DB) error {
	//1. Create a migration driver from our existing DB connection
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil{
		return  fmt.Errorf("Could not create migration driver: %v", err)
	}

	//2. Point to your migration files (ensure this path is correct in Docker)
	// We use "file://internal/migrations" relative to the project root
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/migrations",
		"postgres", driver)

	if err != nil{
		return fmt.Errorf("migration init failed: %v", err)
	}

	//3. Apply the migrations(Up)
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange{
		return  fmt.Errorf("migration up failed: %v", err)
	}

	log.Println("✅ Database migrations synchronized successfully!")
	return nil
}