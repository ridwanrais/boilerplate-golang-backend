package http

import (
	"context"
	"net/http"

	"backend-golang/internal/user"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type GetUserRequest struct {
	ID uuid.UUID `path:"id" doc:"User ID"`
}

type UserResponse struct {
	Body struct {
		Data *user.User `json:"data"`
	}
}

type ListUsersResponse struct {
	Body struct {
		Data []*user.User `json:"data"`
	}
}

func RegisterPublicUserRoutes(api huma.API, uc user.PublicUsecase) {
	huma.Register(api, huma.Operation{
		OperationID: "get-user",
		Method:      http.MethodGet,
		Path:        "/users/{id}",
		Summary:     "Get user by ID",
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
}
