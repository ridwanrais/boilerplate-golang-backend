package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"backend-golang/ent"
	deliveryHTTP "backend-golang/internal/delivery/http"
	"backend-golang/internal/repository"
	"backend-golang/internal/usecase"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestUserAPI_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping API integration test")
	}

	ctx := context.Background()

	// 1. Start PostgreSQL container
	pgContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("test-user"),
		postgres.WithPassword("test-pass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err)
	}

	client, err := ent.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()

	if err := client.Schema.Create(ctx); err != nil {
		t.Fatalf("failed creating schema resources: %v", err)
	}

	// 2. Wire the application exactly like main.go
	userRepo := repository.NewUserEntRepository(client)
	userUC := usecase.NewUserUsecase(userRepo)

	router := http.NewServeMux()
	config := huma.DefaultConfig("Test API", "1.0.0")
	api := humago.New(router, config)

	deliveryHTTP.RegisterUserRoutes(api, userUC)

	// 3. Start test HTTP server
	ts := httptest.NewServer(router)
	defer ts.Close()

	// 4. Test the API Endpoints via actual HTTP requests!
	t.Run("Create and Get User via API", func(t *testing.T) {
		// Test Create
		reqBody := []byte(`{"name":"API Test User","email":"api@test.com"}`)
		resp, err := http.Post(ts.URL+"/users", "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("failed to make POST request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected status OK/Created, got %d", resp.StatusCode)
		}

		var createResp struct {
			Data struct {
				ID    string `json:"id"`
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if createResp.Data.Name != "API Test User" {
			t.Errorf("expected name API Test User, got %s", createResp.Data.Name)
		}

		// Test GetByID
		getResp, err := http.Get(ts.URL + "/users/" + createResp.Data.ID)
		if err != nil {
			t.Fatalf("failed to make GET request: %v", err)
		}
		defer getResp.Body.Close()

		if getResp.StatusCode != http.StatusOK {
			t.Fatalf("expected status OK, got %d", getResp.StatusCode)
		}

		var fetchResp struct {
			Data struct {
				ID    string `json:"id"`
				Email string `json:"email"`
			} `json:"data"`
		}
		if err := json.NewDecoder(getResp.Body).Decode(&fetchResp); err != nil {
			t.Fatalf("failed to decode get response: %v", err)
		}
		if fetchResp.Data.Email != "api@test.com" {
			t.Errorf("expected email api@test.com, got %s", fetchResp.Data.Email)
		}
	})
}
