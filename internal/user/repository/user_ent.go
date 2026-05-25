package repository

import (
	"context"
	"errors"

	"backend-golang/ent"
	entuser "backend-golang/ent/user"
	"backend-golang/internal/user"

	"github.com/google/uuid"
)

type entRepository struct {
	client *ent.Client
}

func NewEntRepository(client *ent.Client) user.Repository {
	return &entRepository{
		client: client,
	}
}

func (r *entRepository) Create(ctx context.Context, u *user.User) (*user.User, error) {
	entUser, err := r.client.User.
		Create().
		SetName(u.Name).
		SetEmail(u.Email).
		SetPassword(u.Password).
		Save(ctx)

	if err != nil {
		return nil, err
	}

	return toDomainUser(entUser), nil
}

func (r *entRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	entUser, err := r.client.User.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return toDomainUser(entUser), nil
}

func (r *entRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	entUser, err := r.client.User.Query().Where(entuser.Email(email)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return toDomainUser(entUser), nil
}

func (r *entRepository) List(ctx context.Context) ([]*user.User, error) {
	entUsers, err := r.client.User.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	var users []*user.User
	for _, eu := range entUsers {
		users = append(users, toDomainUser(eu))
	}

	return users, nil
}

func (r *entRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.client.User.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.New("user not found")
		}
		return err
	}
	return nil
}

func toDomainUser(u *ent.User) *user.User {
	if u == nil {
		return nil
	}
	return &user.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
	}
}
