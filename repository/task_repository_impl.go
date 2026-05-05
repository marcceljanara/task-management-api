package repository

import (
	"context"
	"database/sql"
	"errors"
	"marcceljanara/task-management-api/exception"
	"marcceljanara/task-management-api/model/domain"
)

type TaskRepositoryImpl struct {
}

func (repository *TaskRepositoryImpl) Save(ctx context.Context, db *sql.DB, task domain.Task) error {
	query := "INSERT INTO tasks(id, user_id, title, description, status, priority, due_date) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := db.ExecContext(ctx, query, task.Id, task.UserId, task.Title, task.Description, task.Status, task.Priority, task.DueDate)
	if err != nil {
		return err
	}
	return nil
}

func (repository *TaskRepositoryImpl) FindAll(ctx context.Context, db *sql.DB, userId string, limit int, offset int) ([]domain.Task, int, error) {
	panic("not implemented") // TODO: Implement
}

func (repository *TaskRepositoryImpl) FindById(ctx context.Context, db *sql.DB, taskId string, userId string) (domain.Task, error) {
	query := "SELECT id, title, description, status, priority, due_date, created_at, updated_at FROM task WHERE id = $1 AND user_id = $2"
	task := domain.Task{}
	err := db.QueryRowContext(ctx, query, taskId, userId).Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.Priority, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return task, exception.New(exception.ErrNotFound, "task not found")
		}
		return task, err
	}
	return task, nil
}

func (repository *TaskRepositoryImpl) Update(ctx context.Context, db *sql.DB, task domain.Task) error {
	query := "UPDATE tasks SET title = $1, description = $2, status = $3, priority = $4, due_date = $5, updated_at = NOW() WHERE id = $6 AND user_id = $7"
	_, err := db.ExecContext(ctx, query, task.Title, task.Description, task.Status, task.Priority, task.DueDate, task.Id, task.UserId)
	if err != nil {
		return err
	}
	return nil
}

func (repository *TaskRepositoryImpl) Delete(ctx context.Context, db *sql.DB, taskId string, userId string) error {
	query := "DELETE FROM tasks WHERE id = $1 AND user_id = $2"
	_, err := db.ExecContext(ctx, query, taskId, userId)
	if err != nil {
		return err
	}
	return nil
}
