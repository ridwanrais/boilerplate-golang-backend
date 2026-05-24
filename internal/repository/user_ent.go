package repository

import (
	"context"
	"errors"

	"backend-golang/ent"
	"backend-golang/internal/domain"

	"github.com/google/uuid"
)

type userEntRepository struct {
	client *ent.Client
}

// NewUserEntRepository creates a new instance of UserRepository using Ent.
func NewUserEntRepository(client *ent.Client) domain.UserRepository {
	return &userEntRepository{
		client: client,
	}
}

func (r *userEntRepository) Create(ctx context.Context, dUser *domain.User) (*domain.User, error) {
	u, err := r.client.User.
		Create().
		SetName(dUser.Name).
		SetEmail(dUser.Email).
		Save(ctx)

	if err != nil {
		return nil, err
	}

	return toDomainUser(u), nil
}

func (r *userEntRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	u, err := r.client.User.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return toDomainUser(u), nil
}

func (r *userEntRepository) List(ctx context.Context) ([]*domain.User, error) {
	users, err := r.client.User.
		Query().
		All(ctx)

	if err != nil {
		return nil, err
	}

	var dUsers []*domain.User
	for _, u := range users {
		dUsers = append(dUsers, toDomainUser(u))
	}

	return dUsers, nil
}

func (r *userEntRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.client.User.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.New("user not found")
		}
		return err
	}
	return nil
}

// Helper function to map Ent User to Domain User.
func toDomainUser(u *ent.User) *domain.User {
	if u == nil {
		return nil
	}
	return &domain.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}
