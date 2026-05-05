package main

import (
	"log"
	"marcceljanara/task-management-api/app"
	"marcceljanara/task-management-api/controller"
	"marcceljanara/task-management-api/repository"
	"marcceljanara/task-management-api/service"
	"net/http"
	"os"

	"github.com/go-playground/validator"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func main() {
	_ = godotenv.Load()
	db, err := app.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	validate := validator.New()

	secret := os.Getenv("JWT_SECRET")
	jwtService := service.NewJWTService(secret, "task-management-jwt")
	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository, jwtService, db, validate)
	userController := controller.NewUserController(userService)

	router := httprouter.New()
	router.POST("/api/v1/register", userController.Register)
	router.POST("/api/v1/login", userController.Login)

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
