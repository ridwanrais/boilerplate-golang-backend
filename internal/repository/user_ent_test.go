package repository_test

import (
	"context"
	"testing"
	"time"

	"backend-golang/ent"
	"backend-golang/internal/domain"
	"backend-golang/internal/repository"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestUserEntRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	// 1. Start PostgreSQL container using testcontainers-go (runs smoothly on OrbStack!)
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

	// Clean up the container after the test suite finishes
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	// 2. Get the dynamically assigned connection string
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err)
	}

	// 3. Connect Ent client to the temporary test database
	client, err := ent.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()

	// 4. Run automatic schema migration (perfect for tests!)
	if err := client.Schema.Create(ctx); err != nil {
		t.Fatalf("failed creating schema resources: %v", err)
	}

	// 5. Initialize the repository with our test DB client
	repo := repository.NewUserEntRepository(client)

	// 6. Run actual tests against the real database!
	t.Run("Create and Get User", func(t *testing.T) {
		dUser := &domain.User{
			Name:  "OrbStack Test User",
			Email: "orbstack@test.com",
		}

		// Test Create
		createdUser, err := repo.Create(ctx, dUser)
		if err != nil {
			t.Fatalf("failed to create user: %v", err)
		}
		if createdUser.Name != dUser.Name {
			t.Errorf("expected name %s, got %s", dUser.Name, createdUser.Name)
		}
		if createdUser.ID.String() == "" {
			t.Error("expected non-empty ID")
		}

		// Test GetByID
		fetchedUser, err := repo.GetByID(ctx, createdUser.ID)
		if err != nil {
			t.Fatalf("failed to get user: %v", err)
		}
		if fetchedUser.Email != dUser.Email {
			t.Errorf("expected email %s, got %s", dUser.Email, fetchedUser.Email)
		}
	})
}
