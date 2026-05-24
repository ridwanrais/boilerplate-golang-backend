package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// User represents the domain model for a user.
type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// UserRepository defines the interface for data access.
type UserRepository interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	List(ctx context.Context) ([]*User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// UserUsecase defines the interface for business logic.
type UserUsecase interface {
	CreateUser(ctx context.Context, name, email string) (*User, error)
	GetUser(ctx context.Context, id uuid.UUID) (*User, error)
	ListUsers(ctx context.Context) ([]*User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
