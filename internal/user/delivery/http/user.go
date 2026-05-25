package http

import (
	"context"
	"net/http"

	"backend-golang/internal/user"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Body struct {
		Name  string `json:"name" doc:"User's full name" example:"John Doe" maxLength:"100"`
		Email string `json:"email" doc:"User's email address" example:"john.doe@example.com" format:"email"`
	}
}

type UserResponse struct {
	Body struct {
		Data *user.User `json:"data"`
	}
}

type GetUserRequest struct {
	ID uuid.UUID `path:"id" doc:"User ID"`
}

type ListUsersResponse struct {
	Body struct {
		Data []*user.User `json:"data"`
	}
}

type DeleteUserRequest struct {
	ID uuid.UUID `path:"id" doc:"User ID"`
}

type DeleteUserResponse struct {
	Status int `doc:"204 No Content"`
}

func RegisterRoutes(api huma.API, uc user.Usecase) {
	huma.Register(api, huma.Operation{
		OperationID: "create-user",
		Method:      http.MethodPost,
		Path:        "/users",
		Summary:     "Create a new user",
		Description: "Creates a new user with the given name and email.",
		Tags:        []string{"Users"},
	}, func(ctx context.Context, input *CreateUserRequest) (*UserResponse, error) {
		usr, err := uc.CreateUser(ctx, input.Body.Name, input.Body.Email)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to create user", err)
		}
		resp := &UserResponse{}
		resp.Body.Data = usr
		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-user",
		Method:      http.MethodGet,
		Path:        "/users/{id}",
		Summary:     "Get a user by ID",
		Tags:        []string{"Users"},
	}, func(ctx context.Context, input *GetUserRequest) (*UserResponse, error) {
		usr, err := uc.GetUser(ctx, input.ID)
		if err != nil {
			return nil, huma.Error404NotFound("user not found", err)
		}
		resp := &UserResponse{}
		resp.Body.Data = usr
		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "list-users",
		Method:      http.MethodGet,
		Path:        "/users",
		Summary:     "List all users",
		Tags:        []string{"Users"},
	}, func(ctx context.Context, input *struct{}) (*ListUsersResponse, error) {
		users, err := uc.ListUsers(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to fetch users", err)
		}
		resp := &ListUsersResponse{}
		if users == nil {
			resp.Body.Data = make([]*user.User, 0)
		} else {
			resp.Body.Data = users
		}
		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-user",
		Method:      http.MethodDelete,
		Path:        "/users/{id}",
		Summary:     "Delete a user by ID",
		Tags:        []string{"Users"},
	}, func(ctx context.Context, input *DeleteUserRequest) (*DeleteUserResponse, error) {
		err := uc.DeleteUser(ctx, input.ID)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to delete user", err)
		}
		resp := &DeleteUserResponse{}
		resp.Status = http.StatusNoContent
		return resp, nil
	})
}
