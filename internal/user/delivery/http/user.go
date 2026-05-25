package http

import (
	"context"
	"net/http"

	"backend-golang/internal/user"

	"github.com/danielgtaylor/huma/v2"
)

type ListUsersResponse struct {
	Body struct {
		Data []*user.User `json:"data"`
	}
}

func RegisterPublicUserRoutes(api huma.API, uc user.PublicUsecase) {
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
