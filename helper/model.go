package helper

import (
	"marcceljanara/task-management-api/model/domain"
	"marcceljanara/task-management-api/model/web"
	"time"
)

const dateTimeLayout = time.RFC3339

func ToUserResponse(user domain.User) web.UserResponse {
	return web.UserResponse{
		Id:   user.Id,
		Name: user.Name,
	}
}

func ToTaskResponse(task domain.Task) web.TaskResponse {
	dueDate := ""
	if !task.DueDate.IsZero() {
		dueDate = task.DueDate.Format(dateTimeLayout)
	}

	return web.TaskResponse{
		Id:          task.Id,
		Title:       task.Title,
		Description: task.Description,
		Priority:    task.Priority,
		Status:      task.Status,
		DueDate:     dueDate,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

func ToTaskResponses(tasks []domain.Task) []web.TaskResponse {
	taskResponses := make([]web.TaskResponse, 0, len(tasks))
	for _, task := range tasks {
		taskResponses = append(taskResponses, ToTaskResponse(task))
	}
	return taskResponses
}
