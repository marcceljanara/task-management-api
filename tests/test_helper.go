package tests

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"marcceljanara/task-management-api/controller"
	"marcceljanara/task-management-api/middleware"
	"marcceljanara/task-management-api/model/domain"
	"marcceljanara/task-management-api/repository"
	"marcceljanara/task-management-api/service"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-playground/validator"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/require"
)

const testJWTSecret = "integration-test-secret"

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	loadTestEnv()

	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	db, err := sql.Open("pgx", databaseURL)
	require.NoError(t, err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	require.NoError(t, db.PingContext(context.Background()))
	return db
}

func loadTestEnv() {
	_ = godotenv.Load()

	if os.Getenv("TEST_DATABASE_URL") != "" {
		return
	}

	workingDirectory, err := os.Getwd()
	if err != nil {
		return
	}

	for {
		envPath := filepath.Join(workingDirectory, ".env")
		env, err := godotenv.Read(envPath)
		if err == nil {
			if value := env["TEST_DATABASE_URL"]; value != "" {
				_ = os.Setenv("TEST_DATABASE_URL", value)
			}
			return
		}

		parentDirectory := filepath.Dir(workingDirectory)
		if parentDirectory == workingDirectory {
			return
		}
		workingDirectory = parentDirectory
	}
}

func setupTestRouter(db *sql.DB) (http.Handler, service.JWTService) {
	validate := validator.New()
	jwtService := service.NewJWTService(testJWTSecret, "task-management-jwt")

	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository, jwtService, db, validate)
	userController := controller.NewUserController(userService)

	taskRepository := repository.NewTaskRepository()
	taskService := service.NewTaskService(taskRepository, db, validate)
	taskController := controller.NewTaskController(taskService)

	router := httprouter.New()
	router.POST("/api/v1/register", userController.Register)
	router.POST("/api/v1/login", userController.Login)
	router.POST("/api/v1/tasks", middleware.ValidateJWT(jwtService, taskController.CreateTask))
	router.GET("/api/v1/tasks", middleware.ValidateJWT(jwtService, taskController.GetAllTasks))
	router.GET("/api/v1/tasks/:taskId", middleware.ValidateJWT(jwtService, taskController.GetTaskById))
	router.PUT("/api/v1/tasks/:taskId", middleware.ValidateJWT(jwtService, taskController.UpdateTask))
	router.DELETE("/api/v1/tasks/:taskId", middleware.ValidateJWT(jwtService, taskController.DeleteTask))

	return router, jwtService
}

func truncateTables(t *testing.T, db *sql.DB) {
	t.Helper()

	_, err := db.ExecContext(context.Background(), "TRUNCATE TABLE tasks, users RESTART IDENTITY CASCADE")
	require.NoError(t, err)
}

func seedUser(t *testing.T, db *sql.DB, id string, email string) {
	t.Helper()

	query := "INSERT INTO users(id, name, email, password) VALUES ($1, $2, $3, $4)"
	_, err := db.ExecContext(context.Background(), query, id, "Integration Test User", email, "hashed-password")
	require.NoError(t, err)
}

func seedTask(t *testing.T, db *sql.DB, task domain.Task) domain.Task {
	t.Helper()

	if task.CreatedAt.IsZero() {
		task.CreatedAt = time.Now().UTC()
	}
	if task.UpdatedAt.IsZero() {
		task.UpdatedAt = task.CreatedAt
	}

	task.DueDate = task.DueDate.UTC()
	task.CreatedAt = task.CreatedAt.UTC()
	task.UpdatedAt = task.UpdatedAt.UTC()

	query := `
		INSERT INTO tasks(id, user_id, title, description, status, priority, due_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := db.ExecContext(
		context.Background(),
		query,
		task.Id,
		task.UserId,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.DueDate,
		task.CreatedAt,
		task.UpdatedAt,
	)
	require.NoError(t, err)

	return task
}

func authHeader(t *testing.T, jwtService service.JWTService, userID string) string {
	t.Helper()

	token, err := jwtService.GenerateToken(userID)
	require.NoError(t, err)

	return "Bearer " + token
}

func decodeJSONResponse(t *testing.T, response *http.Response) map[string]any {
	t.Helper()
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	var responseBody map[string]any
	require.NoError(t, json.Unmarshal(body, &responseBody))

	return responseBody
}
