package database

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect now returns *gorm.DB
func Connect(dsn string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	for i := 0; i < 5; i++ {
		db, err = gorm.Open(gormPostgres.Open(dsn), &gorm.Config{})
		if err == nil {
			// Verify connection
			sqlDB, _ := db.DB()
			if err = sqlDB.Ping(); err == nil {
				return db, nil
			}
		}

		fmt.Printf("⏳ Connecting to DB (attempt %d/5)...\n", i+1)
		time.Sleep(2 * time.Second)
	}
	return nil, err
}

// RunMigrations now accepts *gorm.DB
func RunMigrations(gormDB *gorm.DB) error {
	sqlDB, err := gormDB.DB()
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://internal/migrations", "postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("✅ Migrations complete")
	return nil
}