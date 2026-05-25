package usecase

import (
	"context"

	"backend-golang/internal/user"
)

type publicUsecaseImpl struct {
	repo user.Repository
}

func NewPublicUsecase(repo user.Repository) user.PublicUsecase {
	return &publicUsecaseImpl{repo: repo}
}

func (u *publicUsecaseImpl) ListUsers(ctx context.Context) ([]*user.User, error) {
	return u.repo.List(ctx)
}
