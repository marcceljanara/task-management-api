package service

import (
	"context"
	"database/sql"
	"errors"
	"marcceljanara/task-management-api/exception"
	"marcceljanara/task-management-api/helper"
	"marcceljanara/task-management-api/model/domain"
	"marcceljanara/task-management-api/model/web"
	"marcceljanara/task-management-api/repository"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	UserRepository repository.UserRepository
	JWTService     JWTService
	DB             *sql.DB
	Validate       *validator.Validate
}

func NewUserService(userRepository repository.UserRepository, jwtService JWTService, DB *sql.DB, validate *validator.Validate) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
		JWTService:     jwtService,
		DB:             DB,
		Validate:       validate,
	}
}

func (service *UserServiceImpl) Register(ctx context.Context, request web.UserCreateRequest) (response web.UserResponse, err error) {
	err = service.Validate.Struct(request)
	if err != nil {
		return response, exception.Wrap(exception.ErrValidation, "invalid request payload", err)
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return response, exception.Wrap(exception.ErrInternal, "failed to start transaction", err)
	}
	defer func() {
		beforeCommitErr := err
		helper.CommitOrRollback(tx, &err)
		if beforeCommitErr == nil && err != nil {
			err = exception.Wrap(exception.ErrInternal, "failed to commit transaction", err)
		}
	}()

	passwordHashByte, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return response, exception.Wrap(exception.ErrInternal, "failed to hash password", err)
	}

	userID, err := uuid.NewRandom()
	if err != nil {
		return response, exception.Wrap(exception.ErrInternal, "failed to generate user id", err)
	}

	user := domain.User{
		Id:       userID.String(),
		Name:     request.Name,
		Email:    request.Email,
		Password: string(passwordHashByte),
	}

	existingUser, err := service.UserRepository.FindByEmail(ctx, tx, request.Email)
	if err != nil && !errors.Is(err, exception.ErrNotFound) {
		return response, exception.Wrap(exception.ErrInternal, "failed to check email availability", err)
	}
	if existingUser.Email != "" {
		return response, exception.New(exception.ErrConflict, "email address is already registered")
	}

	err = service.UserRepository.InsertUser(ctx, tx, user)
	if err != nil {
		return response, exception.Wrap(exception.ErrInternal, "failed to create user", err)
	}

	return helper.ToUserResponse(user), nil
}

func (service *UserServiceImpl) Login(ctx context.Context, request web.UserLoginRequest) (tokenString string, err error) {
	err = service.Validate.Struct(request)
	if err != nil {
		return tokenString, exception.Wrap(exception.ErrValidation, "invalid request payload", err)
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return tokenString, exception.Wrap(exception.ErrInternal, "failed to start transaction", err)
	}
	defer func() {
		beforeCommitErr := err
		helper.CommitOrRollback(tx, &err)
		if beforeCommitErr == nil && err != nil {
			err = exception.Wrap(exception.ErrInternal, "failed to commit transaction", err)
		}
	}()

	user, err := service.UserRepository.FindByEmail(ctx, tx, request.Email)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return tokenString, exception.New(exception.ErrUnauthorized, "email or password is invalid")
		}
		return tokenString, exception.Wrap(exception.ErrInternal, "failed to find user", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return tokenString, exception.New(exception.ErrUnauthorized, "email or password is invalid")
	}

	tokenString, err = service.JWTService.GenerateToken(user.Id)
	if err != nil {
		return tokenString, exception.Wrap(exception.ErrInternal, "failed to generate access token", err)
	}

	return tokenString, nil
}
