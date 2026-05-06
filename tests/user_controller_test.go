package tests

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

const (
	userControllerTestUserID = "00000000-0000-0000-0000-000000000201"
)

func TestRegisterUserSuccess(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	truncateTables(t, db)

	router, _ := setupTestRouter(db)

	requestBody := strings.NewReader(`{
		"name": "Integration User",
		"email": "register-success@example.com",
		"password": "password123"
	}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/register", requestBody)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	responseBody := decodeJSONResponse(t, response)
	assert.Equal(t, float64(http.StatusOK), responseBody["code"])
	assert.Equal(t, "Berhasil mendaftar akun", responseBody["status"])

	data := responseBody["data"].(map[string]any)
	userID := data["id"].(string)
	assert.NotEmpty(t, userID)
	assert.Equal(t, "Integration User", data["name"])
	assert.NotContains(t, data, "email")
	assert.NotContains(t, data, "password")

	var storedName string
	var storedEmail string
	var storedPassword string
	err := db.QueryRowContext(
		context.Background(),
		"SELECT name, email, password FROM users WHERE id = $1",
		userID,
	).Scan(&storedName, &storedEmail, &storedPassword)
	require.NoError(t, err)

	assert.Equal(t, "Integration User", storedName)
	assert.Equal(t, "register-success@example.com", storedEmail)
	assert.NotEqual(t, "password123", storedPassword)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte("password123")))
}

func TestRegisterUserValidationFailed(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	truncateTables(t, db)

	router, _ := setupTestRouter(db)

	requestBody := strings.NewReader(`{
		"name": "",
		"email": "not-an-email",
		"password": "short"
	}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/register", requestBody)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	responseBody := decodeJSONResponse(t, response)
	assert.Equal(t, float64(http.StatusBadRequest), responseBody["code"])
	assert.Equal(t, "invalid request payload", responseBody["status"])
}

func TestRegisterUserDuplicateEmail(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	truncateTables(t, db)

	router, _ := setupTestRouter(db)
	seedUserWithPassword(t, db, userControllerTestUserID, "Existing User", "duplicate@example.com", "password123")

	requestBody := strings.NewReader(`{
		"name": "New User",
		"email": "duplicate@example.com",
		"password": "password123"
	}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/register", requestBody)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusConflict, response.StatusCode)

	responseBody := decodeJSONResponse(t, response)
	assert.Equal(t, float64(http.StatusConflict), responseBody["code"])
	assert.Equal(t, "email address is already registered", responseBody["status"])
}

func TestRegisterUserInvalidJSON(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	truncateTables(t, db)

	router, _ := setupTestRouter(db)

	requestBody := strings.NewReader(`{"name": "Broken JSON"`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/register", requestBody)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	responseBody := decodeJSONResponse(t, response)
	assert.Equal(t, float64(http.StatusBadRequest), responseBody["code"])
	assert.Equal(t, "invalid JSON request body", responseBody["status"])
}

func TestLoginUserSuccess(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	truncateTables(t, db)

	router, jwtService := setupTestRouter(db)
	seedUserWithPassword(t, db, userControllerTestUserID, "Login User", "login-success@example.com", "password123")

	requestBody := strings.NewReader(`{
		"email": "login-success@example.com",
		"password": "password123"
	}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/login", requestBody)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	responseBody := decodeJSONResponse(t, response)
	assert.Equal(t, float64(http.StatusOK), responseBody["code"])
	assert.Equal(t, "Login berhasil!", responseBody["status"])
	assert.Nil(t, responseBody["data"])

	cookies := response.Cookies()
	require.Len(t, cookies, 1)

	accessTokenCookie := cookies[0]
	assert.Equal(t, "access_token", accessTokenCookie.Name)
	assert.NotEmpty(t, accessTokenCookie.Value)
	assert.Equal(t, "/", accessTokenCookie.Path)
	assert.True(t, accessTokenCookie.HttpOnly)

	userID, err := jwtService.ValidateToken(accessTokenCookie.Value)
	require.NoError(t, err)
	assert.Equal(t, userControllerTestUserID, userID)
}

func TestLoginUserValidationFailed(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	truncateTables(t, db)

	router, _ := setupTestRouter(db)

	requestBody := strings.NewReader(`{
		"email": "not-an-email",
		"password": "short"
	}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/login", requestBody)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	responseBody := decodeJSONResponse(t, response)
	assert.Equal(t, float64(http.StatusBadRequest), responseBody["code"])
	assert.Equal(t, "invalid request payload", responseBody["status"])
}

func TestLoginUserWrongPassword(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	truncateTables(t, db)

	router, _ := setupTestRouter(db)
	seedUserWithPassword(t, db, userControllerTestUserID, "Login User", "wrong-password@example.com", "password123")

	requestBody := strings.NewReader(`{
		"email": "wrong-password@example.com",
		"password": "wrongpass123"
	}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/login", requestBody)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)

	responseBody := decodeJSONResponse(t, response)
	assert.Equal(t, float64(http.StatusUnauthorized), responseBody["code"])
	assert.Equal(t, "email or password is invalid", responseBody["status"])
}

func TestLoginUserEmailNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	truncateTables(t, db)

	router, _ := setupTestRouter(db)

	requestBody := strings.NewReader(`{
		"email": "not-found@example.com",
		"password": "password123"
	}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/login", requestBody)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)

	responseBody := decodeJSONResponse(t, response)
	assert.Equal(t, float64(http.StatusUnauthorized), responseBody["code"])
	assert.Equal(t, "email or password is invalid", responseBody["status"])
}

func TestLoginUserInvalidJSON(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	truncateTables(t, db)

	router, _ := setupTestRouter(db)

	requestBody := strings.NewReader(`{"email": "broken@example.com"`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/login", requestBody)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	responseBody := decodeJSONResponse(t, response)
	assert.Equal(t, float64(http.StatusBadRequest), responseBody["code"])
	assert.Equal(t, "invalid JSON request body", responseBody["status"])
}

func seedUserWithPassword(t *testing.T, db *sql.DB, id string, name string, email string, password string) {
	t.Helper()

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	query := "INSERT INTO users(id, name, email, password) VALUES ($1, $2, $3, $4)"
	_, err = db.ExecContext(context.Background(), query, id, name, email, string(passwordHash))
	require.NoError(t, err)
}
