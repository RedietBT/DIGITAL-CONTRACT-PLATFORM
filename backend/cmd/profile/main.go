package main

import (
	"log"
	"os"

	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/broker"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/database"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/handler"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/middleware"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/repository"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/service"
	pkgBroker "github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/pkg/broker"
	_ "github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/cmd/profile/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Digital Contract Platform - Profile Service
// @version         1.0
// @description     Profile Service for the Digital Contract Platform.
// @host            localhost:8082
// @BasePath        /
// @securityDefinitions.bearerAuth  BearerAuth
// @in              header
// @name            Authorization
func main() {
	dsn := os.Getenv("DATABASE_DSN")
	jwtSecret := os.Getenv("JWT_SECRET")
	if dsn == ""{
		log.Fatal("❌ DATABASE_DSN is not set in environment variables")
	}
	if jwtSecret == ""{
		log.Fatal("❌ JWT_SECRET is not set in environment variables")
	}

	// Call the funvtion from our internal package
	gormDB, err := database.Connect(dsn)
	if err != nil{
		log.Fatal("❌ Failed to connect to database: %v", err)
	}
	sqlDB, _ := gormDB.DB()
	defer sqlDB.Close()

	//Calling Migrate to ensure the User table is created
	if err := database.RunMigrations(gormDB); err != nil{
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
	repo := repository.NewProfileRepository(sqlDB)
	svc := service.NewProfileService(repo)
	h := handler.NewProfileHandler(svc)

	// Intialixe & Start RabbitMQ Consumer
	// We pass 'profilerSvc' because it implements the UserEventHandler interface
	profileConsumer, err := broker.NewProfileConsumer(conn, svc)
	if err != nil {
    log.Fatalf("❌ Failed to create profile consumer: %v", err)
    }
	
	// Start listining in a background goroutine
	err = profileConsumer.Start() 
	if err != nil {
    log.Fatalf("❌ Failed to start profile consumer: %v", err)
    }
	
	// Gin Engine Setup
	r := gin.Default()

	// Add CORS Middleware
	r.Use(middleware.CORSMiddleware())

	// Swagger & Health Check (Public)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Profile Service is running"})
	})

	// Protected Profile Routes
	// We use the custom ProfileAuthMiddleware to protect these routes
	profileGroup := r.Group("/profile")
	profileGroup.Use(middleware.ProfileAuthMiddleware(jwtSecret))
	{
		profileGroup.GET("/me", h.GetProfile)
		profileGroup.PUT("/me", h.UpdateProfile)
	}

	log.Println("🚀 Profile Service API listening on :8082")
    if err := r.Run(":8082"); err != nil {
        log.Fatalf("❌ Server failed: %v", err)
    }
}