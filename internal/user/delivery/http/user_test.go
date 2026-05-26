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
	authhttp "backend-golang/internal/auth/delivery/http"
	authuc "backend-golang/internal/auth/usecase"
	userhttp "backend-golang/internal/user/delivery/http"
	userrepo "backend-golang/internal/user/repository"
	useruc "backend-golang/internal/user/usecase"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestUserAndAuthAPI_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping API integration test")
	}

	ctx := context.Background()

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

	// Wire the application
	userRepo := userrepo.NewEntRepository(client)
	
	authUC := authuc.NewUsecase(userRepo)
	myUserUC := useruc.NewMyUsecase(userRepo)
	publicUserUC := useruc.NewPublicUsecase(userRepo)

	router := http.NewServeMux()
	config := huma.DefaultConfig("Test API", "1.0.0")
	api := humago.New(router, config)

	authhttp.RegisterRoutes(api, authUC)
	userhttp.RegisterMyUserRoutes(api, myUserUC)
	userhttp.RegisterPublicUserRoutes(api, publicUserUC)

	ts := httptest.NewServer(router)
	defer ts.Close()

	var jwtToken string

	t.Run("Register User", func(t *testing.T) {
		reqBody := []byte(`{"name":"Auth User","email":"auth@test.com","password":"securepassword123"}`)
		resp, err := http.Post(ts.URL+"/auth/register", "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("failed to make POST request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected status OK/Created, got %d", resp.StatusCode)
		}
	})

	t.Run("Login User", func(t *testing.T) {
		reqBody := []byte(`{"email":"auth@test.com","password":"securepassword123"}`)
		resp, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("failed to make POST request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status OK, got %d", resp.StatusCode)
		}

		var actualResp struct {
			Token string `json:"token"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&actualResp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		
		if actualResp.Token == "" {
			t.Fatalf("expected token, got empty string")
		}
		jwtToken = actualResp.Token
	})

	var userID uuid.UUID
	
	t.Run("Get My Profile (Protected)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, ts.URL+"/users/me", nil)
		req.Header.Set("Authorization", "Bearer "+jwtToken)
		
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to make GET request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status OK, got %d", resp.StatusCode)
		}

		var fetchResp struct {
			Data struct {
				ID    uuid.UUID `json:"id"`
				Email string    `json:"email"`
			} `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&fetchResp); err != nil {
			t.Fatalf("failed to decode get response: %v", err)
		}
		if fetchResp.Data.Email != "auth@test.com" {
			t.Errorf("expected email auth@test.com, got %s", fetchResp.Data.Email)
		}
		userID = fetchResp.Data.ID
	})

	t.Run("Get My Profile (Unauthorized)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, ts.URL+"/users/me", nil)
		
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to make GET request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized && resp.StatusCode != 422 { // Huma might throw 422 if header is required
			t.Fatalf("expected status 401 or 422, got %d", resp.StatusCode)
		}
	})

	t.Run("Get User Detail (Public)", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/users/" + userID.String())
		if err != nil {
			t.Fatalf("failed to make GET request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status OK, got %d", resp.StatusCode)
		}

		var detailResp struct {
			Data struct {
				ID    uuid.UUID `json:"id"`
				Email string    `json:"email"`
			} `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&detailResp); err != nil {
			t.Fatalf("failed to decode get response: %v", err)
		}
		if detailResp.Data.Email != "auth@test.com" {
			t.Errorf("expected email auth@test.com, got %s", detailResp.Data.Email)
		}
	})

	t.Run("List Users (Public)", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/users")
		if err != nil {
			t.Fatalf("failed to make GET request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status OK, got %d", resp.StatusCode)
		}
	})
}
