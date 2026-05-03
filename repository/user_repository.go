package repository

import (
	"context"
	"database/sql"
	"marcceljanara/task-management-api/model/domain"
)

type UserRepository interface {
	InsertUser(ctx context.Context, tx *sql.DB, user domain.User)
	FindByEmail(ctx context.Context, tx *sql.DB, email string) (domain.User, error)
}