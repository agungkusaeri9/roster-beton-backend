package main

import (
	"go-arch/internal/delivery/http/handler"
	"go-arch/internal/infrastructure/pgsql"
	pgsqlRepo "go-arch/internal/infrastructure/pgsql"
	"go-arch/internal/repository"
	"go-arch/internal/usecase"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 🧩 Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("❌ Error loading .env file")
	}
	log.Println("✅ Env loaded successfully")

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

	// 🧩 Initialize usecase
	jwtSecret := "supersecret" // TODO: ambil dari env
	authUsecase := usecase.NewAuthUseCase(userRepo, jwtSecret)

	// 🧩 Initialize Gin router
	r := gin.Default()

	// 🧩 Register handler (controller)
	handler.NewAuthHandler(r, authUsecase)

	// 🧩 (Optional) log semua route terdaftar
	for _, route := range r.Routes() {
		log.Printf("%s %s\n", route.Method, route.Path)
	}

	// 🧩 Start server
	log.Println("🚀 Server started on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
