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

func (repository *TaskRepositoryImpl) FindAll(ctx context.Context, db *sql.DB, task domain.Task, limit int, offset int) ([]domain.Task, int, error) {
	var tasks []domain.Task
	var totalRows int
	
	query := `SELECT title, status, priority, due_date 
	FROM tasks 
	WHERE user_id = $1,
	status = $2, 
	LOWER(title) LIKE '$3%'
	ORDER BY created_at DESC
	LIMIT $4 OFFSET $5
	`
	rows, err := db.QueryContext(ctx, query, task.UserId, task.Status, task.Title, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	
	defer rows.Close()

	for rows.Next() {
		var task domain.Task
		rows.Scan(&task.Title, &task.Status, &task.Priority, &task.DueDate)
		tasks = append(tasks, task)
	}

	countQuery := "SELECT COUNT(*) FROM tasks WHERE user_id = $1"
	err = db.QueryRowContext(ctx, countQuery, task.UserId).Scan(&totalRows)
	if err != nil {
		return nil, 0, err
	}

	return tasks, totalRows, nil
	 
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
