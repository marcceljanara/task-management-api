package helper

import (
	"marcceljanara/task-management-api/model/domain"
	"marcceljanara/task-management-api/model/web"
)

func ToUserResponse(user domain.User) web.UserResponse {
	return web.UserResponse{
		Id: user.Id,
		Name: user.Name,
	}
}