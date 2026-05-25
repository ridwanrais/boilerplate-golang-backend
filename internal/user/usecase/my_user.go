package usecase

import (
	"context"

	"backend-golang/internal/user"

	"github.com/google/uuid"
)

type myUsecaseImpl struct {
	repo user.Repository
}

func NewMyUsecase(repo user.Repository) user.MyUsecase {
	return &myUsecaseImpl{repo: repo}
}

func (u *myUsecaseImpl) GetProfile(ctx context.Context, id uuid.UUID) (*user.User, error) {
	return u.repo.GetByID(ctx, id)
}
