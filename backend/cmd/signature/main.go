package signature

import (
	"log"
	"os"
	"strings"

	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/broker"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/database"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/handler"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/middleware"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/repository"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/service"
	pkgBroker "github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/pkg/broker"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Digital Contract Platform - Signature Service
// @version         1.0
// @description Type your JWT token
// @description     Service for managing digital signatures.
// @host            localhost:8084
// @BasePath        /
// @securityDefinitions.apiKey AuthKey
// @type                       apiKey
// @in                         header
// @name                       Authorization
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
	repo := repository.NewSignatureRepository(db)

	// Create the publisher
	publisher, err := broker.NewContractPublisher(conn)
	if err != nil {
		log.Fatalf("❌ Failed to create publisher: %v", err)
	}
	
	svc := service.NewSignatureService(repo, publisher)
	validate := validator.New()
	validate.RegisterValidation("no_scripts", func(fl validator.FieldLevel) bool {
    return !strings.Contains(strings.ToLower(fl.Field().String()), "<script>")
    })
	h := handler.NewSignatureHandler(svc, validate)

	// Intialixe & Start RabbitMQ Consumer
	// We pass 'repo' because it implements the ContractEventHandler interface
	sigConsumer, err := broker.NewSignatureConsumer(conn, repo)
	if err != nil {
		log.Fatalf("❌ Failed to create signature consumer: %v", err)
	}

	// Running consumer in background so it doesn't block the API
	go sigConsumer.Listen()

	// Gin Engine Setup
	r := gin.Default()

    r.Use(middleware.CORSMiddleware()) // Must be FIRST

	// Swagger & Health Check (Public)
	// Replace your existing r.GET("/swagger/*any", ...) with this:
	// Replace the complex WrapHandler with this in all 3 main.go files
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Signature Service is running"})
	})

	// Protected Contract Routes
	// We use the custom ContractAuthMiddleware to protect these routes
	signatureRoutes := r.Group("/signatures")
	signatureRoutes.Use(middleware.SignatureMiddleware(jwtSecret))
	{
		signatureRoutes.POST("", h.SignContract)
		signatureRoutes.GET("/:id", h.GetSignature)
		signatureRoutes.DELETE("/:id", h.RevokeSignature)
	}

	log.Println("🚀 Signature Service listening on :8084")
	//This rplaces 'select{}' and keeps the app runing
	if err := r.Run(":8084"); err != nil{
		log.Fatalf("❌ Server failed: %v", err)
	}
}