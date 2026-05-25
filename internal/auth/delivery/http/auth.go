package http

import (
	"context"
	"net/http"

	"backend-golang/internal/auth"

	"github.com/danielgtaylor/huma/v2"
)

type RegisterRequest struct {
	Body struct {
		Name     string `json:"name" doc:"User's full name" example:"John Doe"`
		Email    string `json:"email" doc:"User's email" format:"email"`
		Password string `json:"password" doc:"User's password" minLength:"6"`
	}
}

type RegisterResponse struct {
	Status int `doc:"201 Created"`
}

type LoginRequest struct {
	Body struct {
		Email    string `json:"email" doc:"User's email" format:"email"`
		Password string `json:"password" doc:"User's password"`
	}
}

type LoginResponse struct {
	Body struct {
		Token string `json:"token"`
	}
}

func RegisterRoutes(api huma.API, uc auth.Usecase) {
	huma.Register(api, huma.Operation{
		OperationID: "register-user",
		Method:      http.MethodPost,
		Path:        "/auth/register",
		Summary:     "Register a new user",
		Tags:        []string{"Auth"},
	}, func(ctx context.Context, input *RegisterRequest) (*RegisterResponse, error) {
		err := uc.Register(ctx, input.Body.Name, input.Body.Email, input.Body.Password)
		if err != nil {
			return nil, huma.Error400BadRequest("registration failed", err)
		}
		resp := &RegisterResponse{}
		resp.Status = http.StatusCreated
		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "login-user",
		Method:      http.MethodPost,
		Path:        "/auth/login",
		Summary:     "Login user",
		Tags:        []string{"Auth"},
	}, func(ctx context.Context, input *LoginRequest) (*LoginResponse, error) {
		token, err := uc.Login(ctx, input.Body.Email, input.Body.Password)
		if err != nil {
			return nil, huma.Error401Unauthorized("invalid credentials", err)
		}
		resp := &LoginResponse{}
		resp.Body.Token = token
		return resp, nil
	})
}
