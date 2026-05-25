package main

import (
	"fmt"
	"log"
	"os"

	"backend-golang/ent"
	authhttp "backend-golang/internal/auth/delivery/http"
	authuc "backend-golang/internal/auth/usecase"
	userhttp "backend-golang/internal/user/delivery/http"
	userrepo "backend-golang/internal/user/repository"
	useruc "backend-golang/internal/user/usecase"

	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
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

	// 2. Setup Standard Library Router
	router := http.NewServeMux()

	// 3. Setup Huma API on top of standard ServeMux
	config := huma.DefaultConfig("Backend API Boilerplate", "1.0.0")
	
	// Add security scheme to API config
	config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"bearerAuth": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		},
	}
	
	api := humago.New(router, config)

	// 4. Wire dependencies (Feature Modules)
	userRepo := userrepo.NewEntRepository(client)
	
	authUC := authuc.NewUsecase(userRepo)
	myUserUC := useruc.NewMyUsecase(userRepo)
	publicUserUC := useruc.NewPublicUsecase(userRepo)

	// 5. Register Routes
	authhttp.RegisterRoutes(api, authUC)
	userhttp.RegisterMyUserRoutes(api, myUserUC)
	userhttp.RegisterPublicUserRoutes(api, publicUserUC)

	// 6. Start Server
	log.Println("Server is running on http://localhost:8080")
	log.Println("Swagger UI is available at http://localhost:8080/docs")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
