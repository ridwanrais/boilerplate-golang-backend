package usecase

import (
	"context"

	"backend-golang/internal/user"

	"github.com/google/uuid"
)

type usecaseImpl struct {
	repo user.Repository
}

func NewUsecase(repo user.Repository) user.Usecase {
	return &usecaseImpl{repo: repo}
}

func (u *usecaseImpl) CreateUser(ctx context.Context, name, email string) (*user.User, error) {
	usr := &user.User{
		Name:  name,
		Email: email,
	}
	return u.repo.Create(ctx, usr)
}

func (u *usecaseImpl) GetUser(ctx context.Context, id uuid.UUID) (*user.User, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *usecaseImpl) ListUsers(ctx context.Context) ([]*user.User, error) {
	return u.repo.List(ctx)
}

func (u *usecaseImpl) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}
