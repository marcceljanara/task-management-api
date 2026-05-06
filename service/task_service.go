package service

import (
	"context"
	"marcceljanara/task-management-api/model/web"
)

type TaskService interface {
	Save(ctx context.Context, request web.TaskCreateRequest, userId string) (web.TaskResponse, error)
	FindAll(ctx context.Context, request web.TaskQueryFindAllRequest, userId string) (web.TasksResponse, error)
	FindById(ctx context.Context, taskId string, userId string) (web.TaskResponse, error)
	Update(ctx context.Context, request web.TaskUpdateRequest, userId string) error
	Delete(ctx context.Context, taskId string, userId string) error
}