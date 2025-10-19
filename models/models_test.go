package models

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite: %v", err)
	}
	// Ensure SQLite enforces foreign keys for cascade behavior
	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		t.Fatalf("failed to enable foreign keys: %v", err)
	}
	if err := db.AutoMigrate(&User{}, &Task{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}
	return db
}

func TestUserUniqueEmail(t *testing.T) {
	db := setupTestDB(t)

	u1 := User{Name: "Alice", Email: "alice@example.com"}
	if err := db.Create(&u1).Error; err != nil {
		t.Fatalf("failed to create first user: %v", err)
	}
	u2 := User{Name: "Alice 2", Email: "alice@example.com"}
	if err := db.Create(&u2).Error; err == nil {
		t.Fatalf("expected error on duplicate email, got nil")
	}
}

func TestTaskForeignKeyToUser(t *testing.T) {
	db := setupTestDB(t)

	u := User{Name: "Bob", Email: "bob@example.com"}
	if err := db.Create(&u).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	task := Task{Task: "Do something", UserID: int(u.ID)}
	if err := db.Create(&task).Error; err != nil {
		t.Fatalf("failed to create task: %v", err)
	}

	var fetched Task
	if err := db.First(&fetched, task.ID).Error; err != nil {
		t.Fatalf("failed to fetch task: %v", err)
	}
	if fetched.UserID != int(u.ID) {
		t.Fatalf("task user id mismatch: got %d want %d", fetched.UserID, u.ID)
	}
}

func TestCascadeDeleteUserDeletesTasks(t *testing.T) {
	db := setupTestDB(t)

	u := User{Name: "Charlie", Email: "charlie@example.com"}
	if err := db.Create(&u).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	t1 := Task{Task: "Task 1", UserID: int(u.ID)}
	t2 := Task{Task: "Task 2", UserID: int(u.ID)}
	if err := db.Create(&t1).Error; err != nil {
		t.Fatalf("failed to create task1: %v", err)
	}
	if err := db.Create(&t2).Error; err != nil {
		t.Fatalf("failed to create task2: %v", err)
	}

	// Use hard delete to trigger FK cascade at DB level
	if err := db.Unscoped().Delete(&u).Error; err != nil {
		t.Fatalf("failed to delete user: %v", err)
	}

	var count int64
	if err := db.Model(&Task{}).Where("user_id = ?", u.ID).Count(&count).Error; err != nil {
		t.Fatalf("failed counting tasks: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected 0 tasks after cascade delete, got %d", count)
	}
}
