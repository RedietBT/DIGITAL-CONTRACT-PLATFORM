package database

import (
	"database/sql"
	"fmt"
	"time"

	_"github.com/lib/pq"
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