package usecase

import (
	"context"
	"errors"
	"time"

	"backend-golang/internal/auth"
	"backend-golang/internal/user"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var JWTSecret = []byte("super-secret-key") // In a real app, this should be in an environment variable

type authUsecase struct {
	userRepo user.Repository
}

func NewUsecase(userRepo user.Repository) auth.Usecase {
	return &authUsecase{userRepo: userRepo}
}

func (u *authUsecase) Register(ctx context.Context, name, email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	usr := &user.User{
		Name:     name,
		Email:    email,
		Password: string(hash),
	}
	_, err = u.userRepo.Create(ctx, usr)
	return err
}

func (u *authUsecase) Login(ctx context.Context, email, password string) (string, error) {
	usr, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Create JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": usr.ID.String(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(JWTSecret)
}
