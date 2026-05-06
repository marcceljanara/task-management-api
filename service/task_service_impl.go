package service

import (
	"context"
	"database/sql"
	"errors"
	"marcceljanara/task-management-api/exception"
	"marcceljanara/task-management-api/helper"
	"marcceljanara/task-management-api/model/domain"
	"marcceljanara/task-management-api/model/web"
	"marcceljanara/task-management-api/repository"
	"strings"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type TaskServiceImpl struct {
	TaskRepository repository.TaskRepository
	DB             *sql.DB
	validate       *validator.Validate
}

func NewTaskService(taskRepository repository.TaskRepository, DB *sql.DB, validate *validator.Validate) TaskService {
	return &TaskServiceImpl{
		TaskRepository: taskRepository,
		DB:             DB,
		validate:       validate,
	}
}

func (service *TaskServiceImpl) Save(ctx context.Context, request web.TaskCreateRequest, userId string) (web.TaskResponse, error) {
	userId = strings.TrimSpace(userId)
	request.Title = strings.TrimSpace(request.Title)

	if userId == "" {
		return web.TaskResponse{}, exception.New(exception.ErrUnauthorized, "invalid user session")
	}

	if request.Priority == "" {
		request.Priority = "medium"
	}
	if request.Status == "" {
		request.Status = "pending"
	}

	err := service.validate.Struct(request)
	if err != nil {
		return web.TaskResponse{}, exception.Wrap(exception.ErrValidation, "request body not valid", err)
	}

	taskID, err := uuid.NewRandom()
	if err != nil {
		return web.TaskResponse{}, exception.Wrap(exception.ErrInternal, "failed to generate task id", err)
	}

	task := domain.Task{
		Id:          taskID.String(),
		UserId:      userId,
		Title:       request.Title,
		Description: request.Description,
		Status:      request.Status,
		Priority:    request.Priority,
		DueDate:     request.DueDate,
	}

	err = service.TaskRepository.Save(ctx, service.DB, task)
	if err != nil {
		return web.TaskResponse{}, exception.Wrap(exception.ErrInternal, "failed to create task", err)
	}

	createdTask, err := service.TaskRepository.FindById(ctx, service.DB, task.Id, userId)
	if err != nil {
		return web.TaskResponse{}, exception.Wrap(exception.ErrInternal, "failed to find created task", err)
	}

	return helper.ToTaskResponse(createdTask), nil
}

func (service *TaskServiceImpl) FindAll(ctx context.Context, request web.TaskQueryFindAllRequest, userId string) (web.TasksResponse, error) {
	userId = strings.TrimSpace(userId)
	request.Status = strings.TrimSpace(request.Status)

	if userId == "" {
		return web.TasksResponse{}, exception.New(exception.ErrUnauthorized, "invalid user session")
	}

	if request.Page < 1 {
		request.Page = 1
	}
	if request.Limit < 1 {
		request.Limit = 10
	}
	if request.Limit > 100 {
		request.Limit = 100
	}
	if request.Status != "" {
		err := service.validate.Var(request.Status, "oneof=pending in_progress completed")
		if err != nil {
			return web.TasksResponse{}, exception.Wrap(exception.ErrValidation, "invalid task status filter", err)
		}
	}

	filter := domain.Task{
		UserId: userId,
		Title:  strings.TrimSpace(request.Title),
		Status: request.Status,
	}
	offset := (request.Page - 1) * request.Limit

	tasks, totalRows, err := service.TaskRepository.FindAll(ctx, service.DB, filter, request.Limit, offset)
	if err != nil {
		return web.TasksResponse{}, exception.Wrap(exception.ErrInternal, "failed to find tasks", err)
	}

	totalPages := 0
	if totalRows > 0 {
		totalPages = (totalRows + request.Limit - 1) / request.Limit
	}

	return web.TasksResponse{
		Data: helper.ToTaskFindAllResponses(tasks),
		Pagination: web.Pagination{
			Page:       request.Page,
			Limit:      request.Limit,
			TotalRows:  totalRows,
			TotalPages: totalPages,
			HasNext:    request.Page < totalPages,
			HasPrev:    request.Page > 1,
		},
	}, nil
}

func (service *TaskServiceImpl) FindById(ctx context.Context, taskId string, userId string) (web.TaskResponse, error) {
	userId = strings.TrimSpace(userId)
	taskId = strings.TrimSpace(taskId)

	if userId == "" {
		return web.TaskResponse{}, exception.New(exception.ErrUnauthorized, "invalid user session")
	}
	if taskId == "" {
		return web.TaskResponse{}, exception.New(exception.ErrValidation, "task id is required")
	}

	task, err := service.TaskRepository.FindById(ctx, service.DB, taskId, userId)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return web.TaskResponse{}, err
		}
		return web.TaskResponse{}, exception.Wrap(exception.ErrInternal, "failed to find task", err)
	}

	return helper.ToTaskResponse(task), nil
}

func (service *TaskServiceImpl) Update(ctx context.Context, request web.TaskUpdateRequest, userId string) error {
	userId = strings.TrimSpace(userId)
	request.Id = strings.TrimSpace(request.Id)
	request.Title = strings.TrimSpace(request.Title)

	if userId == "" {
		return exception.New(exception.ErrUnauthorized, "invalid user session")
	}

	err := service.validate.Struct(request)
	if err != nil {
		return exception.Wrap(exception.ErrValidation, "request body not valid", err)
	}

	task := domain.Task{
		Id:          request.Id,
		UserId:      userId,
		Title:       request.Title,
		Description: request.Description,
		Status:      request.Status,
		Priority:    request.Priority,
		DueDate:     request.DueDate,
	}

	err = service.TaskRepository.Update(ctx, service.DB, task)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return err
		}
		return exception.Wrap(exception.ErrInternal, "failed to update task", err)
	}

	return nil
}

func (service *TaskServiceImpl) Delete(ctx context.Context, taskId string, userId string) error {
	userId = strings.TrimSpace(userId)
	taskId = strings.TrimSpace(taskId)

	if userId == "" {
		return exception.New(exception.ErrUnauthorized, "invalid user session")
	}
	if taskId == "" {
		return exception.New(exception.ErrValidation, "task id is required")
	}

	err := service.TaskRepository.Delete(ctx, service.DB, taskId, userId)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return err
		}
		return exception.Wrap(exception.ErrInternal, "failed to delete task", err)
	}

	return nil
}
