package main

import (
	"fmt"
	"log"
	"os"

	"backend-golang/ent"
	deliveryHTTP "backend-golang/internal/delivery/http"
	"backend-golang/internal/repository"
	"backend-golang/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load .env file if it exists
	_ = godotenv.Load()

	// 1. Initialize Postgres Database & Ent Client
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"))
	
	client, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()

	// 2. Setup Gin Router
	router := gin.Default()

	// 3. Setup Huma API on top of Gin
	config := huma.DefaultConfig("Backend API Boilerplate", "1.0.0")
	api := humagin.New(router, config)

	// 4. Wire dependencies (Clean Architecture)
	userRepo := repository.NewUserEntRepository(client)
	userUC := usecase.NewUserUsecase(userRepo)

	// 5. Register Routes
	deliveryHTTP.RegisterUserRoutes(api, userUC)

	// 6. Start Server
	log.Println("Server is running on http://localhost:8080")
	log.Println("Swagger UI is available at http://localhost:8080/docs")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
