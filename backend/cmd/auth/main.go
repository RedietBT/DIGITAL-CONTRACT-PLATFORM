package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/cmd/auth/docs" // Swagger docs
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/broker"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/database"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/handler"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/middleware"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/repository"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/service"
	pkgBroker "github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/pkg/broker"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Digital Contract Platform API (Auth)
// @version         1.0
// @description     Auth Service with RabbitMQ & Gin.
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey ApiKeyAuth
// @in              header
// @name            Authorization
func main() {
	// 1. Env Variables
	dsn := os.Getenv("DATABASE_DSN")
	jwtSecret := os.Getenv("JWT_SECRET")
	if dsn == "" || jwtSecret == "" {
		log.Fatal("❌ Critical Environment Variables (DSN/JWT) are missing")
	}

	// 2. Database Connection & Migration
	db, err := database.Connect(dsn)
	if err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	// 3. RabbitMQ Connection
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@rabbitmq:5672/"
	}
	conn, err := pkgBroker.Connect(rabbitURL)
	if err != nil {
		log.Fatalf("❌ RabbitMQ connection failed: %v", err)
	}
	defer conn.Close()

	authPub, err := broker.NewAuthPublisher(conn)
	if err != nil {
		log.Fatalf("❌ Publisher initialization failed: %v", err)
	}

	// 4. Layers Initialization
	repo := repository.NewPostgresUserRepository(db)
	svc := service.NewAuthService(repo, authPub, jwtSecret)
	h := handler.NewAuthHandler(svc)

	// 5. Gin Setup
	r := gin.Default()

	// 6. Swagger & Health
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/health", func(c *gin.Context) {
		c.String(200, "OK - Auth Service is alive")
	})

	// 7. Route Groups
	
	// PUBLIC
	auth := r.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/forgot-password", h.ForgotPassword)
		auth.POST("/reset-password", h.ResetPassword)
		auth.POST("/refresh", h.Refresh)
		auth.POST("/logout", h.Logout)
	}

	// PROTECTED (Requires Login)
	me := r.Group("/auth/me")
	me.Use(middleware.AuthMiddleware(jwtSecret))
	{
		me.GET("", h.GetProfile)
		me.PUT("/email", h.UpdateEmail)
		me.PUT("/password", h.ChangePassword)
		me.DELETE("", h.DeleteMe)
	}

	// ADMIN ONLY (Requires Login + Admin Role)
	admin := r.Group("/auth/admin")
	admin.Use(middleware.AuthMiddleware(jwtSecret), middleware.RoleMiddleware("admin", svc))
	{
		admin.GET("/users", h.GetAllUsers)
		admin.DELETE("/users", h.DeleteUser)
		admin.PUT("/user-status", h.UpdateUserStatus)
	}

	log.Println("🚀 Auth Service listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("❌ Server failure: %v", err)
	}
}