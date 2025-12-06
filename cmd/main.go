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

	// 🧩 Initialize Gin router
	r := gin.Default()

	// 🧩 Register handler (controller)
	handler.NewAuthHandler(r, authUsecase)

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
