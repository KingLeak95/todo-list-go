package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KingLeak95/todo-list-go/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// testRouter creates a Gin engine with routes wired to the handlers, using the provided DB
func testRouter(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	// Initialize in-memory SQLite DB
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite: %v", err)
	}

	// Auto-migrate schemas
	if err := db.AutoMigrate(&models.User{}, &models.Task{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}

	// Override global DB used by handlers
	models.DB = db

	r := gin.New()
	r.Use(gin.Recovery())

	// Public authentication routes
	r.POST("/auth/register", models.Register)
	r.POST("/auth/login", models.Login)
	r.POST("/auth/refresh", models.RefreshToken)

	// Legacy user creation (deprecated, use /auth/register)
	r.POST("/createUser", models.CreateUser)

	// Protected routes - require authentication
	protected := r.Group("/")
	// Note: In tests, we'll bypass auth middleware for simplicity
	{
		// Users
		protected.GET("/allUsers", models.GetAllUsers)
		protected.DELETE("/deleteUser/:id", models.DeleteUser)

		// Tasks
		protected.GET("/tasks", models.GetAllTasks)
		protected.POST("/tasks", models.CreateTask)
		protected.GET("/tasks/:id", models.GetTaskByID)
		protected.PUT("/tasks/:id", models.UpdateTask)
		protected.DELETE("/tasks/:id", models.DeleteTask)
		protected.PUT("/tasks/:id/complete", models.CompleteTask)
		protected.GET("/users/:id/tasks", models.GetTasksByUser)
	}

	return r
}

func doJSONRequest(t *testing.T, r http.Handler, method, path string, body interface{}) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("failed to encode body: %v", err)
		}
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func doJSONRequestWithHeaders(t *testing.T, r http.Handler, method, path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("failed to encode body: %v", err)
		}
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestUserLifecycle(t *testing.T) {
	r := testRouter(t)

	// Register user
	registerPayload := map[string]interface{}{
		"name":     "John Doe",
		"email":    "john@example.com",
		"password": "password123",
	}
	w := doJSONRequest(t, r, http.MethodPost, "/auth/register", registerPayload)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body=%s", w.Code, w.Body.String())
	}

	// Login user
	loginPayload := map[string]interface{}{
		"email":    "john@example.com",
		"password": "password123",
	}
	w = doJSONRequest(t, r, http.MethodPost, "/auth/login", loginPayload)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}

	// Extract token from response
	var loginResponse map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &loginResponse); err != nil {
		t.Fatalf("failed to parse login response: %v", err)
	}

	data := loginResponse["data"].(map[string]interface{})
	accessToken := data["access_token"].(string)

	// List users with authentication
	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
	}
	w = doJSONRequestWithHeaders(t, r, http.MethodGet, "/allUsers", nil, headers)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}

	// Delete user id 1 (first created)
	w = doJSONRequestWithHeaders(t, r, http.MethodDelete, "/deleteUser/1", nil, headers)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on delete, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestTaskLifecycle(t *testing.T) {
	r := testRouter(t)

	// Register a user
	userPayload := map[string]interface{}{
		"name":     "Alice",
		"email":    "alice@example.com",
		"password": "password123",
	}
	w := doJSONRequest(t, r, http.MethodPost, "/auth/register", userPayload)
	if w.Code != http.StatusCreated {
		t.Fatalf("register user expected 201, got %d, body=%s", w.Code, w.Body.String())
	}

	// Extract token from response
	var registerResponse map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &registerResponse); err != nil {
		t.Fatalf("failed to parse register response: %v", err)
	}

	data := registerResponse["data"].(map[string]interface{})
	accessToken := data["access_token"].(string)

	// Create task for user 1
	taskPayload := map[string]interface{}{
		"task":   "Buy groceries",
		"userId": 1,
	}
	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
	}
	w = doJSONRequestWithHeaders(t, r, http.MethodPost, "/tasks", taskPayload, headers)
	if w.Code != http.StatusCreated {
		t.Fatalf("create task expected 201, got %d, body=%s", w.Code, w.Body.String())
	}

	// Complete task id 1
	w = doJSONRequestWithHeaders(t, r, http.MethodPut, "/tasks/1/complete", nil, headers)
	if w.Code != http.StatusOK {
		t.Fatalf("complete task expected 200, got %d, body=%s", w.Code, w.Body.String())
	}

	// Delete task id 1
	w = doJSONRequestWithHeaders(t, r, http.MethodDelete, "/tasks/1", nil, headers)
	if w.Code != http.StatusOK {
		t.Fatalf("delete task expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
}
