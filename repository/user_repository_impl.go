package repository

import (
	"context"
	"database/sql"
	"errors"
	"marcceljanara/task-management-api/exception"
	"marcceljanara/task-management-api/model/domain"
)

type UserRepositoryImpl struct {
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (repository *UserRepositoryImpl) InsertUser(ctx context.Context, tx *sql.Tx, user domain.User) error {
	query := "INSERT INTO users(id, name, email, password) VALUES ($1, $2, $3, $4)"
	_, err := tx.ExecContext(ctx, query, user.Id, user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (repository *UserRepositoryImpl) FindByEmail(ctx context.Context, tx *sql.Tx, email string) (domain.User, error) {
	query := "SELECT id, email, name, password, created_at FROM users WHERE email = $1"
	user := domain.User{}

	err := tx.QueryRowContext(ctx, query, email).Scan(&user.Id, &user.Email, &user.Name, &user.Password, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, exception.New(exception.ErrNotFound, "email not found")
		}
		return user, err
	}

	return user, nil
}
