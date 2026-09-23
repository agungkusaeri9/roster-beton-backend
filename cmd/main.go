package main

import (
	"go-arch/internal/delivery/http/handler"
	"go-arch/internal/infrastructure/pgsql"
	pgsqlRepo "go-arch/internal/infrastructure/pgsql"
	"go-arch/internal/repository"
	"go-arch/internal/usecase"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func findEnvFile() string {
	// List of possible .env file locations
	possiblePaths := []string{
		".env",           // Current directory
		"../.env",       // Parent directory (if running from cmd/)
		"../../.env",    // Two levels up
	}
	
	// Try each path
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	
	// Default fallback
	return ".env"
}

func main() {
	// 🧩 Load environment variables
	envFile := findEnvFile()
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("⚠️  Warning: Error loading .env file from %s: %v", envFile, err)
		log.Println("💡 Trying to continue with system environment variables...")
	} else {
		log.Printf("✅ Env loaded successfully from %s", envFile)
	}

	// 🧩 Initialize PostgreSQL connection
	pgdb, err := pgsql.Init()
	if err != nil {
		log.Fatalf("❌ Failed to init PostgreSQL: %v", err)
	}
	if pgdb == nil {
		log.Fatal("❌ PostgreSQL DB object is nil — check Init()")
	}
	log.Println("✅ PostgreSQL connected successfully")

	// 🧩 Initialize repository
	var userRepo repository.UserRepository = pgsqlRepo.NewUserRepoPg(pgdb)
	var productRepo repository.ProductRepository = pgsqlRepo.NewProductRepoPg(pgdb)
	var categoryRepo repository.CategoryRepository = pgsqlRepo.NewCategoryRepoPg(pgdb)
	var galleryRepo repository.GalleryRepository = pgsqlRepo.NewGalleryRepoPg(pgdb)
	var configRepo repository.ConfigRepository = pgsqlRepo.NewConfigRepoPg(pgdb)

	// 🧩 Get JWT Secret from environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("❌ JWT_SECRET environment variable is required")
	}

	// 🧩 Get Token Expiry from environment (default: 24h)
	tokenExpiryStr := os.Getenv("TOKEN_EXPIRY")
	if tokenExpiryStr == "" {
		tokenExpiryStr = "24h" // Default to 24 hours
		log.Println("⚠️  TOKEN_EXPIRY not set, using default: 24h")
	}

	tokenExpiry, err := time.ParseDuration(tokenExpiryStr)
	if err != nil {
		log.Fatalf("❌ Invalid TOKEN_EXPIRY format: %v. Use format like '24h', '1h', '30m', etc.", err)
	}

	log.Printf("✅ JWT Secret loaded, Token expiry: %s", tokenExpiry)

	// 🧩 Initialize usecase
	authUsecase := usecase.NewAuthUseCase(userRepo, jwtSecret, tokenExpiry)
	productUsecase := usecase.NewProductUsecase(productRepo)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)
	galleryUsecase := usecase.NewGalleryUsecase(galleryRepo)
	configUsecase := usecase.NewConfigUsecase(configRepo)

	// 🧩 Initialize Gin router
	r := gin.Default()

	// 🧩 CORS Middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 🧩 Serve static files for uploads (products, gallery, etc.)
	r.Static("/uploads", "./uploads")

	// 🧩 API Route Group (/api)
	api := r.Group("/api")

	// 🧩 Register handlers (controllers)
	handler.NewAuthHandler(api, authUsecase)
	handler.NewProductHandler(api, productUsecase)
	handler.NewCategoryHandler(api, categoryUsecase)
	handler.NewGalleryHandler(api, galleryUsecase)
	handler.NewConfigHandler(api, configUsecase)

	// 🧩 (Optional) log semua route terdaftar
	for _, route := range r.Routes() {
		log.Printf("%s %s\n", route.Method, route.Path)
	}

	// 🧩 Get port from environment (default: 8080)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 🧩 Start server
	log.Printf("🚀 Server started on http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
