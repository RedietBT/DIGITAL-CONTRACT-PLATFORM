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
	_ "github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/cmd/contract/docs"
	pkgBroker "github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/pkg/broker"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Digital Contract Platform - Contract Service
// @version         1.0
// @description     Service for managing digital contracts.
// @host            localhost:8081
// @BasePath        /

// @securityDefinitions.bearerAuth BearerAuth
// @in header
// @name Authorization
// @description Type 'Bearer ' followed by your JWT token.
func main() {
	dsn := os.Getenv("DATABASE_DSN")
	jwtSecret := os.Getenv("JWT_SECRET")
	if dsn == ""{
		log.Fatal("❌ DATABASE_DSN is not set in environment variables")
	}
	if jwtSecret == ""{
		log.Fatal("❌ JWT_SECRET is not set in environment variables")
	}
	//Call the function from our internal package
	db, err := database.Connect(dsn)
	if err != nil{
		log.Fatal("❌ FAILED: Could not connect to DB: %v", err)
	}

	// Correctly cleanup
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	//Calling Migrate to ensure the User table is created
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
	repo := repository.NewContractRepository(db)
	svc := service.NewContractService(repo)
	validate := validator.New()
	h := handler.NewContractHandler(svc, validate)

	// Intialixe & Start RabbitMQ Consumer
	// We pass 'repo' because it implements the ContractEventHandler interface
	contractConsumer, err := broker.NewContractConsumer(conn, repo)
	if err != nil {
		log.Fatalf("❌ Failed to create contract consumer: %v", err)
	}

	// Start listining in a background goroutine
	err = contractConsumer.Start()
	if err != nil {
		log.Fatalf("❌ Failed to start contract consumer: %v", err)
	}

	// Gin Engine Setup
	r := gin.Default()

	// Add CORS Middleware
	r.Use(middleware.CORSMiddleware())

	// Swagger & Health Check (Public)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Contract Service is running"})
	})

	// Protected Contract Routes
	// We use the custom ContractAuthMiddleware to protect these routes
	contractRoutes := r.Group("/contracts")
	contractRoutes.Use(middleware.ContractMiddleware(jwtSecret))
	{
		contractRoutes.POST("/", h.CreateContract)
		contractRoutes.GET("/:id", h.GetContract)
		contractRoutes.GET("/", h.ListContracts)
		contractRoutes.PUT("/:id", h.UpdateContract)
		contractRoutes.DELETE("/:id", h.DeleteContract)
	}

	log.Println("🚀 Contract Service listening on :8081")
	//This rplaces 'select{}' and keeps the app runing
	if err := r.Run(":8081"); err != nil{
		log.Fatalf("❌ Server failed: %v", err)
	}
}