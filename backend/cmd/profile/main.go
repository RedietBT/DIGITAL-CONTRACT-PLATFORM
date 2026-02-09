package main

import (
	"log"
	"net/http"
	"os"

	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/broker"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/database"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/repository"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/service"
	pkgBroker "github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/pkg/broker"
)

// @title           Digital Contract Platform API
// @version         1.0
// @description     This is the Profile Service for the Digital Contract Platform.
// @host            localhost:8081
// @BasePath        /

func main() {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == ""{
		log.Fatal("❌ DATABASE_DSN is not set in environment variables")
	}

	// Call the funvtion from our internal package
	db, err := database.Connect(dsn)
	if err != nil{
		log.Fatal("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	//Calling AutoMigrate to ensure the User table is created
	if err := database.RunMigrations(db); err != nil{
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == ""{
		rabbitURL = "amqp://guest:guest@rabbitmq:5672/"
	}

	// RabbitMQ Connection (Using our shared pkg/broker)
	conn, err := pkgBroker.Connect(rabbitURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	// Intialize Layers (Dependency Injection)
	repo := repository.NewProfileRepository(db)
	profileSvc := service.NewProfileService(repo)

	// Intialixe & Start RabbitMQ Consumer
	// We pass 'profilerSvc' because it implements the UserEventHandler interface
	profileConsumer, err := broker.NewProfileConsumer(conn, profileSvc)
	if err != nil {
		log.Fatalf("❌ Failed to create profile consumer: %v", err)
	}
	
	// Start listining in a background goroutine
	profileConsumer.Start()
	log.Println("❌ Failed to create consumer: %v", err)

	// 6. Start HTTP Server (Example placeholder)
	// This is where you Gin or Echo handlers would go 
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Profile Service is Healthy"))
	})

	log.Println("🌐 Profile Service API listening on :8082")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}