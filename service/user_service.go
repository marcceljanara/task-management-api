package service

import (
	"context"
	"marcceljanara/task-management-api/model/web"
)

type UserService interface {
	Register(ctx context.Context, request web.UserCreateRequest) (web.UserResponse, error)
	Login(ctx context.Context, request web.UserLoginRequest) (string, error)
}
