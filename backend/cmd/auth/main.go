package main

import (
	"log"
	"net/http"
	"os"

	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/database"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/handler"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/repository"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/service"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Digital Contract Platform API
// @version         1.0
// @description     This is the Auth Service for the Digital Contract Platform.
// @host            localhost:8080
// @BasePath        /
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

	//1. Initiatize Repositories
	repo := repository.NewPostgresUserRepository(db)

	//2. Initiatize Services( Inject Repo)
	svc := service.NewAuthService(repo)

	//3. Initiatize Handlers( Inject Services)
	h := handler.NewAuthHandler(svc)

	//4. Setup Routes
	http.HandleFunc("/auth/register", h.Register)

	//Swagger UI Route
	http.Handle("/swagger/", httpSwagger.WrapHandler)

	log.Println("✅ Auth Service started successfully!")
	log.Println("📖 Swagger Docs available at http://localhost:8080/swagger/index.html")

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