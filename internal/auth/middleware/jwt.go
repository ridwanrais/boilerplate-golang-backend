package middleware

import (
	"errors"
	"fmt"
	"strings"

	"backend-golang/internal/auth/usecase"

	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type ProtectedRequest struct {
	Authorization string    `header:"Authorization" doc:"Bearer token"`
	UserID        uuid.UUID `json:"-"`
}

func (r *ProtectedRequest) Resolve(ctx huma.Context) []error {
	if r.Authorization == "" {
		return []error{errors.New("missing authorization header")}
	}

	parts := strings.Split(r.Authorization, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return []error{errors.New("invalid authorization header format")}
	}

	tokenString := parts[1]
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return usecase.JWTSecret, nil
	})

	if err != nil || !token.Valid {
		return []error{errors.New("invalid token")}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return []error{errors.New("invalid token claims")}
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return []error{errors.New("token missing subject")}
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		return []error{errors.New("invalid user ID in token")}
	}

	r.UserID = userID
	return nil
}
