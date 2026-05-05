package repository

import (
	"context"
	"database/sql"
	"marcceljanara/task-management-api/model/domain"
)

type TaskRepositoryImpl struct {
}

func (repository *TaskRepositoryImpl) Save(ctx context.Context, db *sql.DB, task domain.Task) error {
	panic("not implemented") // TODO: Implement
}

func (repository *TaskRepositoryImpl) FindAll(ctx context.Context, db *sql.DB, userId string, limit int, offset int) ([]domain.Task, int, error) {
	panic("not implemented") // TODO: Implement
}

func (repository *TaskRepositoryImpl) FindById(ctx context.Context, db *sql.DB, taskId string, userId string) (domain.Task, error) {
	panic("not implemented") // TODO: Implement
}

func (repository *TaskRepositoryImpl) Update(ctx context.Context, db *sql.DB, task domain.Task, userId string) error {
	panic("not implemented") // TODO: Implement
}

func (repository *TaskRepositoryImpl) Delete(ctx context.Context, db *sql.DB, taskId string, userId string) error {
	panic("not implemented") // TODO: Implement
}
