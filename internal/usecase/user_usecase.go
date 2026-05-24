package usecase

import (
	"context"

	"backend-golang/internal/domain"

	"github.com/google/uuid"
)

type userUsecase struct {
	repo domain.UserRepository
}

// NewUserUsecase creates a new instance of UserUsecase.
func NewUserUsecase(repo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{repo: repo}
}

func (u *userUsecase) CreateUser(ctx context.Context, name, email string) (*domain.User, error) {
	user := &domain.User{
		Name:  name,
		Email: email,
	}
	return u.repo.Create(ctx, user)
}

func (u *userUsecase) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *userUsecase) ListUsers(ctx context.Context) ([]*domain.User, error) {
	return u.repo.List(ctx)
}

func (u *userUsecase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}
