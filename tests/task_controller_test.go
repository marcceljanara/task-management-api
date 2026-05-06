package tests

import (
	"context"
	"marcceljanara/task-management-api/model/domain"
	"marcceljanara/task-management-api/model/web"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testUserID      = "00000000-0000-0000-0000-000000000001"
	otherTestUserID = "00000000-0000-0000-0000-000000000002"
	testTaskID      = "00000000-0000-0000-0000-000000000101"
	otherTestTaskID = "00000000-0000-0000-0000-000000000102"
)

func TestCreateTaskSuccess(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	truncateTables(t, db)

	router, jwtService := setupTestRouter(db)
	seedUser(t, db, testUserID, "task-create@example.com")

	requestBody := strings.NewReader(`{
		"title": "Memasak Mie Goreng",
		"description": "langkah langkah benar dalam membuat mie goreng adalah sebagai berikut....",
		"priority": "high",
		"status": "pending",
		"due_date": "07-05-2026 12:30:00"
	}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/tasks", requestBody)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", authHeader(t, jwtService, testUserID))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	responseBody := decodeJSONResponse(t, response)
	assert.Equal(t, float64(http.StatusOK), responseBody["code"])
	assert.Equal(t, "Berhasil membuat task", responseBody["status"])

	data := responseBody["data"].(map[string]any)
	taskID := data["id"].(string)
	assert.NotEmpty(t, taskID)
	assert.Equal(t, "Memasak Mie Goreng", data["title"])
	assert.Equal(t, "high", data["priority"])
	assert.Equal(t, "pending", data["status"])
	assert.Equal(t, "07-05-2026 12:30:00", data["due_date"])

	var storedDueDate time.Time
	err := db.QueryRowContext(context.Background(), "SELECT due_date FROM tasks WHERE id = $1", taskID).Scan(&storedDueDate)
	require.NoError(t, err)
	assert.Equal(t, "07-05-2026 12:30:00", storedDueDate.UTC().Format(web.TaskDateTimeLayout))
}

func TestCreateTaskValidationFailed(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	truncateTables(t, db)

	router, jwtService := setupTestRouter(db)
	seedUser(t, db, testUserID, "task-create-failed@example.com")

	requestBody := strings.NewReader(`{
		"title": "",
		"description": "invalid task",
		"priority": "high",
		"status": "pending",
		"due_date": "07-05-2026 12:30:00"
	}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/tasks", requestBody)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", authHeader(t, jwtService, testUserID))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	responseBody := decodeJSONResponse(t, response)
	assert.Equal(t, float64(http.StatusBadRequest), responseBody["code"])
}

func TestGetAllTasksSuccess(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	truncateTables(t, db)

	router, jwtService := setupTestRouter(db)
	seedUser(t, db, testUserID, "task-list@example.com")
	seedUser(t, db, otherTestUserID, "task-list-other@example.com")

	firstTask := seedTask(t, db, domain.Task{
		Id:          testTaskID,
		UserId:      testUserID,
		Title:       "Memasak Mie Goreng",
		Description: "task milik user login",
		Status:      "pending",
		Priority:    "high",
		DueDate:     time.Date(2026, time.May, 7, 12, 30, 0, 0, time.UTC),
		CreatedAt:   time.Date(2026, time.May, 1, 10, 0, 0, 0, time.UTC),
	})
	secondTask := seedTask(t, db, domain.Task{
		Id:          otherTestTaskID,
		UserId:      testUserID,
		Title:       "Membaca Buku",
		Description: "task kedua milik user login",
		Status:      "completed",
		Priority:    "medium",
		DueDate:     time.Date(2026, time.May, 8, 12, 30, 0, 0, time.UTC),
		CreatedAt:   time.Date(2026, time.May, 2, 10, 0, 0, 0, time.UTC),
	})
	seedTask(t, db, domain.Task{
		Id:          "00000000-0000-0000-0000-000000000103",
		UserId:      otherTestUserID,
		Title:       "Task User Lain",
		Description: "tidak boleh muncul",
		Status:      "pending",
		Priority:    "low",
		DueDate:     time.Date(2026, time.May, 9, 12, 30, 0, 0, time.UTC),
		CreatedAt:   time.Date(2026, time.May, 3, 10, 0, 0, 0, time.UTC),
	})

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/tasks", nil)
	request.Header.Set("Authorization", authHeader(t, jwtService, testUserID))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	responseBody := decodeJSONResponse(t, response)
	assert.Equal(t, float64(http.StatusOK), responseBody["code"])
	assert.Equal(t, "Berhasil mengambil daftar task", responseBody["status"])

	data := responseBody["data"].([]any)
	require.Len(t, data, 2)

	taskIDs := make([]string, 0, len(data))
	for _, item := range data {
		task := item.(map[string]any)
		taskIDs = append(taskIDs, task["id"].(string))
		assert.Len(t, task, 5)
		assert.Contains(t, task, "id")
		assert.Contains(t, task, "title")
		assert.Contains(t, task, "status")
		assert.Contains(t, task, "priority")
		assert.Contains(t, task, "due_date")
		assert.NotContains(t, task, "description")
		assert.NotContains(t, task, "created_at")
		assert.NotContains(t, task, "updated_at")
	}
	assert.ElementsMatch(t, []string{firstTask.Id, secondTask.Id}, taskIDs)

	pagination := responseBody["pagination"].(map[string]any)
	assert.Equal(t, float64(1), pagination["page"])
	assert.Equal(t, float64(10), pagination["limit"])
	assert.Equal(t, float64(2), pagination["total_rows"])
	assert.Equal(t, float64(1), pagination["total_pages"])
	assert.Equal(t, false, pagination["has_next"])
	assert.Equal(t, false, pagination["has_prev"])
}

func TestGetAllTasksWithFilter(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	truncateTables(t, db)

	router, jwtService := setupTestRouter(db)
	seedUser(t, db, testUserID, "task-filter@example.com")

	matchedTask := seedTask(t, db, domain.Task{
		Id:          testTaskID,
		UserId:      testUserID,
		Title:       "Memasak Mie Goreng",
		Description: "matched task",
		Status:      "pending",
		Priority:    "high",
		DueDate:     time.Date(2026, time.May, 7, 12, 30, 0, 0, time.UTC),
		CreatedAt:   time.Date(2026, time.May, 1, 10, 0, 0, 0, time.UTC),
	})
	seedTask(t, db, domain.Task{
		Id:          otherTestTaskID,
		UserId:      testUserID,
		Title:       "Memasak Nasi Goreng",
		Description: "status berbeda",
		Status:      "completed",
		Priority:    "high",
		DueDate:     time.Date(2026, time.May, 8, 12, 30, 0, 0, time.UTC),
		CreatedAt:   time.Date(2026, time.May, 2, 10, 0, 0, 0, time.UTC),
	})

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/tasks?status=pending&title=Memasak", nil)
	request.Header.Set("Authorization", authHeader(t, jwtService, testUserID))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	responseBody := decodeJSONResponse(t, response)
	data := responseBody["data"].([]any)
	require.Len(t, data, 1)

	task := data[0].(map[string]any)
	assert.Equal(t, matchedTask.Id, task["id"])
	assert.Equal(t, matchedTask.Title, task["title"])
	assert.Equal(t, matchedTask.Status, task["status"])
}

func TestGetTaskByIdSuccess(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	truncateTables(t, db)

	router, jwtService := setupTestRouter(db)
	seedUser(t, db, testUserID, "task-detail@example.com")
	task := seedTask(t, db, domain.Task{
		Id:          testTaskID,
		UserId:      testUserID,
		Title:       "Memasak Mie Goreng",
		Description: "detail task",
		Status:      "pending",
		Priority:    "high",
		DueDate:     time.Date(2026, time.May, 7, 12, 30, 0, 0, time.UTC),
		CreatedAt:   time.Date(2026, time.May, 1, 10, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2026, time.May, 1, 10, 0, 0, 0, time.UTC),
	})

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/tasks/"+task.Id, nil)
	request.Header.Set("Authorization", authHeader(t, jwtService, testUserID))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	responseBody := decodeJSONResponse(t, response)
	data := responseBody["data"].(map[string]any)
	assert.Equal(t, task.Id, data["id"])
	assert.Equal(t, task.Title, data["title"])
	assert.Equal(t, task.Description, data["description"])
	assert.Equal(t, task.Status, data["status"])
	assert.Equal(t, task.Priority, data["priority"])
	assert.Equal(t, "07-05-2026 12:30:00", data["due_date"])
	assert.Contains(t, data, "created_at")
	assert.Contains(t, data, "updated_at")
}

func TestGetTaskByIdNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	truncateTables(t, db)

	router, jwtService := setupTestRouter(db)
	seedUser(t, db, testUserID, "task-detail-not-found@example.com")
	seedUser(t, db, otherTestUserID, "task-detail-other@example.com")
	seedTask(t, db, domain.Task{
		Id:          testTaskID,
		UserId:      otherTestUserID,
		Title:       "Task User Lain",
		Description: "tidak boleh diakses",
		Status:      "pending",
		Priority:    "high",
		DueDate:     time.Date(2026, time.May, 7, 12, 30, 0, 0, time.UTC),
	})

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/tasks/"+testTaskID, nil)
	request.Header.Set("Authorization", authHeader(t, jwtService, testUserID))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusNotFound, response.StatusCode)

	responseBody := decodeJSONResponse(t, response)
	assert.Equal(t, float64(http.StatusNotFound), responseBody["code"])
}

func TestUpdateTaskSuccess(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	truncateTables(t, db)

	router, jwtService := setupTestRouter(db)
	seedUser(t, db, testUserID, "task-update@example.com")
	task := seedTask(t, db, domain.Task{
		Id:          testTaskID,
		UserId:      testUserID,
		Title:       "Memasak Mie Goreng",
		Description: "sebelum update",
		Status:      "pending",
		Priority:    "high",
		DueDate:     time.Date(2026, time.May, 7, 12, 30, 0, 0, time.UTC),
		CreatedAt:   time.Date(2026, time.May, 1, 10, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2026, time.May, 1, 10, 0, 0, 0, time.UTC),
	})

	requestBody := strings.NewReader(`{
		"title": "Memasak Mie Goreng Spesial",
		"description": "setelah update",
		"priority": "medium",
		"status": "completed",
		"due_date": "08-05-2026 09:15:00"
	}`)
	request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/v1/tasks/"+task.Id, requestBody)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", authHeader(t, jwtService, testUserID))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	responseBody := decodeJSONResponse(t, response)
	data := responseBody["data"].(map[string]any)
	assert.Equal(t, task.Id, data["id"])
	assert.Equal(t, "Memasak Mie Goreng Spesial", data["title"])
	assert.Equal(t, "setelah update", data["description"])
	assert.Equal(t, "medium", data["priority"])
	assert.Equal(t, "completed", data["status"])
	assert.Equal(t, "08-05-2026 09:15:00", data["due_date"])

	var storedTitle string
	var storedDueDate time.Time
	var storedUpdatedAt time.Time
	err := db.QueryRowContext(context.Background(), "SELECT title, due_date, updated_at FROM tasks WHERE id = $1", task.Id).Scan(&storedTitle, &storedDueDate, &storedUpdatedAt)
	require.NoError(t, err)

	assert.Equal(t, "Memasak Mie Goreng Spesial", storedTitle)
	assert.Equal(t, "08-05-2026 09:15:00", storedDueDate.UTC().Format(web.TaskDateTimeLayout))
	assert.WithinDuration(t, time.Now().UTC(), storedUpdatedAt.UTC(), 5*time.Second)
}

func TestUpdateTaskValidationFailed(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	truncateTables(t, db)

	router, jwtService := setupTestRouter(db)
	seedUser(t, db, testUserID, "task-update-failed@example.com")
	task := seedTask(t, db, domain.Task{
		Id:          testTaskID,
		UserId:      testUserID,
		Title:       "Memasak Mie Goreng",
		Description: "sebelum update",
		Status:      "pending",
		Priority:    "high",
		DueDate:     time.Date(2026, time.May, 7, 12, 30, 0, 0, time.UTC),
	})

	requestBody := strings.NewReader(`{
		"title": "Memasak Mie Goreng Spesial",
		"description": "setelah update",
		"priority": "medium",
		"status": "done",
		"due_date": "08-05-2026 09:15:00"
	}`)
	request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/v1/tasks/"+task.Id, requestBody)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", authHeader(t, jwtService, testUserID))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	responseBody := decodeJSONResponse(t, response)
	assert.Equal(t, float64(http.StatusBadRequest), responseBody["code"])
}

func TestDeleteTaskSuccess(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	truncateTables(t, db)

	router, jwtService := setupTestRouter(db)
	seedUser(t, db, testUserID, "task-delete@example.com")
	task := seedTask(t, db, domain.Task{
		Id:          testTaskID,
		UserId:      testUserID,
		Title:       "Memasak Mie Goreng",
		Description: "akan dihapus",
		Status:      "pending",
		Priority:    "high",
		DueDate:     time.Date(2026, time.May, 7, 12, 30, 0, 0, time.UTC),
	})

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/v1/tasks/"+task.Id, nil)
	request.Header.Set("Authorization", authHeader(t, jwtService, testUserID))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	responseBody := decodeJSONResponse(t, response)
	assert.Equal(t, float64(http.StatusOK), responseBody["code"])

	var totalRows int
	err := db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM tasks WHERE id = $1", task.Id).Scan(&totalRows)
	require.NoError(t, err)
	assert.Equal(t, 0, totalRows)
}

func TestDeleteTaskNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	truncateTables(t, db)

	router, jwtService := setupTestRouter(db)
	seedUser(t, db, testUserID, "task-delete-not-found@example.com")

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/v1/tasks/"+testTaskID, nil)
	request.Header.Set("Authorization", authHeader(t, jwtService, testUserID))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusNotFound, response.StatusCode)

	responseBody := decodeJSONResponse(t, response)
	assert.Equal(t, float64(http.StatusNotFound), responseBody["code"])
}

func TestTaskUnauthorized(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	truncateTables(t, db)

	router, _ := setupTestRouter(db)
	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/tasks", nil)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)

	responseBody := decodeJSONResponse(t, response)
	assert.Equal(t, float64(http.StatusUnauthorized), responseBody["code"])
}
