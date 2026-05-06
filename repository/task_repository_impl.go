package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"marcceljanara/task-management-api/exception"
	"marcceljanara/task-management-api/model/domain"
	"strings"
	"time"
)

type TaskRepositoryImpl struct {
}

func NewTaskRepository() TaskRepository {
	return &TaskRepositoryImpl{}
}

func (repository *TaskRepositoryImpl) Save(ctx context.Context, db *sql.DB, task domain.Task) error {
	query := "INSERT INTO tasks(id, user_id, title, description, status, priority, due_date) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := db.ExecContext(ctx, query, task.Id, task.UserId, task.Title, task.Description, task.Status, task.Priority, task.DueDate.UTC())
	if err != nil {
		return err
	}
	return nil
}

func (repository *TaskRepositoryImpl) FindAll(ctx context.Context, db *sql.DB, task domain.Task, limit int, offset int) ([]domain.Task, int, error) {
	var tasks []domain.Task
	var totalRows int

	conditions := []string{"user_id = $1"}
	args := []any{task.UserId}
	argIndex := 2

	if task.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, task.Status)
		argIndex++
	}
	if task.Title != "" {
		conditions = append(conditions, fmt.Sprintf("title ILIKE $%d", argIndex))
		args = append(args, task.Title+"%")
		argIndex++
	}

	whereClause := strings.Join(conditions, " AND ")
	query := fmt.Sprintf(`
		SELECT id, user_id, title, COALESCE(description, ''), status, priority, due_date, created_at, COALESCE(updated_at, created_at)
		FROM tasks
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	queryArgs := append(args, limit, offset)
	rows, err := db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var task domain.Task
		err = rows.Scan(&task.Id, &task.UserId, &task.Title, &task.Description, &task.Status, &task.Priority, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM tasks WHERE %s", whereClause)
	err = db.QueryRowContext(ctx, countQuery, args...).Scan(&totalRows)
	if err != nil {
		return nil, 0, err
	}

	return tasks, totalRows, nil
}

func (repository *TaskRepositoryImpl) FindById(ctx context.Context, db *sql.DB, taskId string, userId string) (domain.Task, error) {
	query := `
		SELECT id, user_id, title, COALESCE(description, ''), status, priority, due_date, created_at, COALESCE(updated_at, created_at)
		FROM tasks
		WHERE id = $1 AND user_id = $2
	`
	task := domain.Task{}
	err := db.QueryRowContext(ctx, query, taskId, userId).Scan(&task.Id, &task.UserId, &task.Title, &task.Description, &task.Status, &task.Priority, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return task, exception.New(exception.ErrNotFound, "task not found")
		}
		return task, err
	}
	return task, nil
}

func (repository *TaskRepositoryImpl) Update(ctx context.Context, db *sql.DB, task domain.Task) error {
	query := "UPDATE tasks SET title = $1, description = $2, status = $3, priority = $4, due_date = $5, updated_at = $6 WHERE id = $7 AND user_id = $8"
	result, err := db.ExecContext(ctx, query, task.Title, task.Description, task.Status, task.Priority, task.DueDate.UTC(), time.Now().UTC(), task.Id, task.UserId)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return exception.New(exception.ErrNotFound, "task not found")
	}
	return nil
}

func (repository *TaskRepositoryImpl) Delete(ctx context.Context, db *sql.DB, taskId string, userId string) error {
	query := "DELETE FROM tasks WHERE id = $1 AND user_id = $2"
	result, err := db.ExecContext(ctx, query, taskId, userId)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return exception.New(exception.ErrNotFound, "task not found")
	}
	return nil
}
