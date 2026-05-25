package user

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type Repository interface {
	Create(ctx context.Context, u *User) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	List(ctx context.Context) ([]*User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type MyUsecase interface {
	GetProfile(ctx context.Context, id uuid.UUID) (*User, error)
}

type PublicUsecase interface {
	ListUsers(ctx context.Context) ([]*User, error)
}
