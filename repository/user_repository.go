package repository

import (
	"context"
	"database/sql"
	"marcceljanara/task-management-api/model/domain"
)

type UserRepository interface {
	InsertUser(ctx context.Context, tx *sql.Tx, user domain.User) error
	FindByEmail(ctx context.Context, tx *sql.Tx, email string) (domain.User, error)
}
