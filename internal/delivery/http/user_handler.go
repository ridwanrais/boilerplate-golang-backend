package http

import (
	"context"
	"net/http"

	"backend-golang/internal/domain"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

// CreateUserRequest represents the input for creating a user.
type CreateUserRequest struct {
	Body struct {
		Name  string `json:"name" doc:"User's full name" example:"John Doe" maxLength:"100"`
		Email string `json:"email" doc:"User's email address" example:"john.doe@example.com" format:"email"`
	}
}

type UserResponse struct {
	Body struct {
		Data *domain.User `json:"data"`
	}
}

type GetUserRequest struct {
	ID uuid.UUID `path:"id" doc:"User ID"`
}

type ListUsersResponse struct {
	Body struct {
		Data []*domain.User `json:"data"`
	}
}

type DeleteUserRequest struct {
	ID uuid.UUID `path:"id" doc:"User ID"`
}

type DeleteUserResponse struct {
	Status int `doc:"204 No Content"`
}

func RegisterUserRoutes(api huma.API, userUC domain.UserUsecase) {
	// POST /users
	huma.Register(api, huma.Operation{
		OperationID: "create-user",
		Method:      http.MethodPost,
		Path:        "/users",
		Summary:     "Create a new user",
		Description: "Creates a new user with the given name and email.",
		Tags:        []string{"Users"},
	}, func(ctx context.Context, input *CreateUserRequest) (*UserResponse, error) {
		user, err := userUC.CreateUser(ctx, input.Body.Name, input.Body.Email)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to create user", err)
		}
		resp := &UserResponse{}
		resp.Body.Data = user
		return resp, nil
	})

	// GET /users/{id}
	huma.Register(api, huma.Operation{
		OperationID: "get-user",
		Method:      http.MethodGet,
		Path:        "/users/{id}",
		Summary:     "Get a user by ID",
		Tags:        []string{"Users"},
	}, func(ctx context.Context, input *GetUserRequest) (*UserResponse, error) {
		user, err := userUC.GetUser(ctx, input.ID)
		if err != nil {
			return nil, huma.Error404NotFound("user not found", err)
		}
		resp := &UserResponse{}
		resp.Body.Data = user
		return resp, nil
	})

	// GET /users
	huma.Register(api, huma.Operation{
		OperationID: "list-users",
		Method:      http.MethodGet,
		Path:        "/users",
		Summary:     "List all users",
		Tags:        []string{"Users"},
	}, func(ctx context.Context, input *struct{}) (*ListUsersResponse, error) {
		users, err := userUC.ListUsers(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to fetch users", err)
		}
		resp := &ListUsersResponse{}
		if users == nil {
			resp.Body.Data = make([]*domain.User, 0)
		} else {
			resp.Body.Data = users
		}
		return resp, nil
	})

	// DELETE /users/{id}
	huma.Register(api, huma.Operation{
		OperationID: "delete-user",
		Method:      http.MethodDelete,
		Path:        "/users/{id}",
		Summary:     "Delete a user by ID",
		Tags:        []string{"Users"},
	}, func(ctx context.Context, input *DeleteUserRequest) (*DeleteUserResponse, error) {
		err := userUC.DeleteUser(ctx, input.ID)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to delete user", err)
		}
		resp := &DeleteUserResponse{}
		resp.Status = http.StatusNoContent
		return resp, nil
	})
}
