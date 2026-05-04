package repository

import (
	"context"
	"database/sql"
	"errors"
	"marcceljanara/task-management-api/model/domain"
)

type UserRepositoryImpl struct {
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (repository *UserRepositoryImpl) InsertUser(ctx context.Context, tx *sql.Tx, user domain.User) {
	sql := "INSERT INTO users(id ,name, email, password) VALUES ($1, $2, $3, $4)"
	_, err := tx.ExecContext(ctx, sql, user.Id, user.Name, user.Email, user.Password)
	if err != nil {
		panic(err)
	}
}

func (repository *UserRepositoryImpl) FindByEmail(ctx context.Context, tx *sql.Tx, email string) (domain.User, error) {
	sql := "select id, email, name, password, created_at from users WHERE email = $1"
	rows, err := tx.QueryContext(ctx, sql, email)
	if err != nil {
		panic(err)
	}
	user := domain.User{}
	defer rows.Close()
	if rows.Next() {
		err := rows.Scan(&user.Id, &user.Email, &user.Name, &user.Password, &user.CreatedAt)
		if err != nil {
			panic(err)
		}
		return user, nil
	} else {
		return user, errors.New("email not found")
	}
}
