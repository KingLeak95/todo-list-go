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

	// Users routes
	r.POST("/createUser", models.CreateUser)
	r.GET("/allUsers", models.GetAllUsers)
	r.DELETE("/deleteUser/:id", models.DeleteUser)

	// Tasks routes
	r.POST("/tasks", models.CreateTask)
	r.DELETE("/tasks/:id", models.DeleteTask)
	r.PUT("/tasks/:id/complete", models.CompleteTask)

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

func TestUserLifecycle(t *testing.T) {
	r := testRouter(t)

	// Create user
	createPayload := map[string]interface{}{
		"name":  "John Doe",
		"email": "john@example.com",
	}
	w := doJSONRequest(t, r, http.MethodPost, "/createUser", createPayload)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body=%s", w.Code, w.Body.String())
	}

	// List users
	w = doJSONRequest(t, r, http.MethodGet, "/allUsers", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}

	// Delete user id 1 (first created)
	w = doJSONRequest(t, r, http.MethodDelete, "/deleteUser/1", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on delete, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestTaskLifecycle(t *testing.T) {
	r := testRouter(t)

	// Create a user to own tasks
	userPayload := map[string]interface{}{
		"name":  "Alice",
		"email": "alice@example.com",
	}
	w := doJSONRequest(t, r, http.MethodPost, "/createUser", userPayload)
	if w.Code != http.StatusCreated {
		t.Fatalf("create user expected 201, got %d, body=%s", w.Code, w.Body.String())
	}

	// Create task for user 1
	taskPayload := map[string]interface{}{
		"task":   "Buy groceries",
		"userId": 1,
	}
	w = doJSONRequest(t, r, http.MethodPost, "/tasks", taskPayload)
	if w.Code != http.StatusCreated {
		t.Fatalf("create task expected 201, got %d, body=%s", w.Code, w.Body.String())
	}

	// Complete task id 1
	w = doJSONRequest(t, r, http.MethodPut, "/tasks/1/complete", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("complete task expected 200, got %d, body=%s", w.Code, w.Body.String())
	}

	// Delete task id 1
	w = doJSONRequest(t, r, http.MethodDelete, "/tasks/1", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("delete task expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
}
