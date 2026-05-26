package usecase

import (
	"context"

	"backend-golang/internal/user"

	"github.com/google/uuid"
)

type publicUsecaseImpl struct {
	repo user.Repository
}

func NewPublicUsecase(repo user.Repository) user.PublicUsecase {
	return &publicUsecaseImpl{repo: repo}
}

func (u *publicUsecaseImpl) GetUser(ctx context.Context, id uuid.UUID) (*user.User, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *publicUsecaseImpl) ListUsers(ctx context.Context) ([]*user.User, error) {
	return u.repo.List(ctx)
}
