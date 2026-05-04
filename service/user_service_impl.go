package service

import (
	"context"
	"database/sql"
	"errors"
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
	JWTService JWTService
	DB             *sql.DB
	Validate       *validator.Validate
}

func NewUserService(userRepository repository.UserRepository, jwtService JWTService, DB *sql.DB, validate *validator.Validate) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
		JWTService: jwtService,
		DB: DB,
		Validate: validate,
	}
}

func (service *UserServiceImpl) Register(ctx context.Context, request web.UserCreateRequest) web.UserResponse {
	err := service.Validate.Struct(request)
	if err != nil {
		panic(err)
	}
	tx, err := service.DB.Begin()
	if err != nil {
		panic(err)
	}
	defer helper.CommitOrRollback(tx)

	passwordByte := []byte(request.Password)
	passwordHashByte, err := bcrypt.GenerateFromPassword(passwordByte, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	uuid, _ := uuid.NewRandom() // uuidv4
	user := domain.User{
		Id: uuid.String(),
		Name: request.Name,
		Email: request.Email,
		Password: string(passwordHashByte),
	}

	result, _ := service.UserRepository.FindByEmail(ctx, tx, request.Email)
	if result.Email != "" {
		panic(errors.New("Email address is already registered"))
	}

	service.UserRepository.InsertUser(ctx, tx, user)

	return helper.ToUserResponse(user)

}

func (service *UserServiceImpl) Login(ctx context.Context, request web.UserLoginRequest) string {
	err := service.Validate.Struct(request)
	if err != nil {
		panic(err)
	}

	tx, err := service.DB.Begin()
	if err != nil {
		panic(err)
	}

	defer helper.CommitOrRollback(tx)

	user, err := service.UserRepository.FindByEmail(ctx, tx, request.Email)
	if err != nil {
		panic(err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		panic(errors.New("Password salah goblok"))
	}

	tokenString, err := service.JWTService.GenerateToken(user.Id)
	if err != nil {
		panic(errors.New("Gagal generate JWT"))
	}

	return tokenString

}
