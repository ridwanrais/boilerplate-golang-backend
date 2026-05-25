package http

import (
	"context"
	"net/http"

	"backend-golang/internal/auth/middleware"
	"backend-golang/internal/user"

	"github.com/danielgtaylor/huma/v2"
)

type MyUserRequest struct {
	middleware.ProtectedRequest
}

type MyUserResponse struct {
	Body struct {
		Data *user.User `json:"data"`
	}
}

func RegisterMyUserRoutes(api huma.API, uc user.MyUsecase) {
	huma.Register(api, huma.Operation{
		OperationID: "get-my-profile",
		Method:      http.MethodGet,
		Path:        "/users/me",
		Summary:     "Get current user profile",
		Tags:        []string{"Users"},
	}, func(ctx context.Context, input *MyUserRequest) (*MyUserResponse, error) {
		usr, err := uc.GetProfile(ctx, input.UserID)
		if err != nil {
			return nil, huma.Error404NotFound("user not found", err)
		}
		resp := &MyUserResponse{}
		resp.Body.Data = usr
		return resp, nil
	})
}
