package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/cmd/auth/docs"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/database"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/handler"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/middleware"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/repository"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/service"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Digital Contract Platform API
// @version         1.0
// @description     This is the Auth Service for the Digital Contract Platform.
// @host            localhost:8080
// @BasePath        /

// @securityDefinitions.apikey ApiKeyAuth
// @in                         header
// @name                       Authorization
// @description                Type "Bearer" followed by a space and JWT token.
func main() {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
        log.Fatal("❌ DATABASE_DSN is not set in environment variables")
    }

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("❌ JWT_SECRET is not set in environment variables")
	}

	//Call the function from our internal package
	db, err := database.Connect(dsn)
	if err != nil {
		log.Fatal("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	//Calling AutoMigrate to ensure the User table is created
	if err := database.RunMigrations(db); err != nil{
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}


	//1. Initiatize Repositories
	repo := repository.NewPostgresUserRepository(db)

	//2. Initiatize Services( Inject Repo)
	svc := service.NewAuthService(repo, jwtSecret)

	//3. Initiatize Handlers( Inject Services)
	h := handler.NewAuthHandler(svc)

	//4. Setup Routes
	// Public routes
	http.HandleFunc("/auth/register", h.Register)
	http.HandleFunc("/auth/login", h.Login)
	http.HandleFunc("/auth/forgot-password", h.ForgotPassword)
	http.HandleFunc("/auth/reset-password", h.ResetPassword)

	// Protected route (Only accessible with a valid token)
	protectedProfile := middleware.AuthMiddleware(jwtSecret)(http.HandlerFunc(h.GetProfile))
	http.Handle("/auth/me", protectedProfile)
	
	// We create the "Admin Only" version of the handler
	// We wrap h.GetAllUsers in RoleMiddleware first
	adminOnly := middleware.RoleMiddleware("admin", svc)(http.HandlerFunc(h.GetAllUsers))
	protectedAdmin := middleware.AuthMiddleware(jwtSecret)(adminOnly)
	http.Handle("/auth/admin/users", protectedAdmin)

	//Swagger UI Route
	http.Handle("/swagger/", httpSwagger.Handler(
    httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // Force the URL
	))

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