package repository

import (
	"context"
	"database/sql"
	"marcceljanara/task-management-api/model/domain"
)

type TaskRepository interface {
	Save(ctx context.Context, db *sql.DB, task domain.Task) error
	FindAll(ctx context.Context, db *sql.DB, userId string, limit int, offset int) ([]domain.Task, int, error)
	FindById(ctx context.Context, db *sql.DB, taskId string, userId string) (domain.Task, error)
	Update(ctx context.Context, db *sql.DB, task domain.Task) error
	Delete(ctx context.Context, db *sql.DB, taskId string, userId string) error
}