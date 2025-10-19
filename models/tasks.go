package models

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TaskPriority represents the priority level of a task
type TaskPriority string

const (
	PriorityLow    TaskPriority = "low"
	PriorityMedium TaskPriority = "medium"
	PriorityHigh   TaskPriority = "high"
)

// TaskStatus represents the status of a task
type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusCompleted TaskStatus = "completed"
	StatusCancelled TaskStatus = "cancelled"
)

type Task struct {
	gorm.Model
	Task        string       `gorm:"column:name;type:varchar(255);not null" json:"task"`
	Description string       `gorm:"type:text" json:"description"`
	Priority    TaskPriority `gorm:"type:varchar(20);default:'medium'" json:"priority"`
	Status      TaskStatus   `gorm:"type:varchar(20);default:'pending'" json:"status"`
	DueDate     *time.Time   `json:"dueDate,omitempty"`
	Category    string       `gorm:"type:varchar(100)" json:"category"`
	Completed   bool         `json:"completed"` // Deprecated: use Status instead
	UserID      int          `json:"userId"`
	User        User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type NewTask struct {
	Task        string       `json:"task" binding:"required"`
	Description string       `json:"description"`
	Priority    TaskPriority `json:"priority"`
	Category    string       `json:"category"`
	DueDate     *time.Time   `json:"dueDate,omitempty"`
	UserID      int          `json:"userId" binding:"required"`
}

type UpdateTaskRequest struct {
	Task        *string       `json:"task,omitempty"`
	Description *string       `json:"description,omitempty"`
	Priority    *TaskPriority `json:"priority,omitempty"`
	Category    *string       `json:"category,omitempty"`
	DueDate     *time.Time    `json:"dueDate,omitempty"`
	Status      *TaskStatus   `json:"status,omitempty"`
}

type TaskQuery struct {
	UserID    *int          `form:"userId"`
	Priority  *TaskPriority `form:"priority"`
	Status    *TaskStatus   `form:"status"`
	Category  *string       `form:"category"`
	Page      int           `form:"page,default=1"`
	Limit     int           `form:"limit,default=10"`
	SortBy    string        `form:"sortBy,default=created_at"`
	SortOrder string        `form:"sortOrder,default=desc"`
	Search    string        `form:"search"`
}

// CreateTask creates a new task
// @Summary Create a new task
// @Description Create a new task with optional priority, category, and due date
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body NewTask true "Task data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tasks [post]
func CreateTask(c *gin.Context) {
	var input NewTask
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Format", "details": err.Error()})
		return
	}

	// Set default priority if not provided
	if input.Priority == "" {
		input.Priority = PriorityMedium
	}

	task := Task{
		Task:        input.Task,
		Description: input.Description,
		Priority:    input.Priority,
		Category:    input.Category,
		DueDate:     input.DueDate,
		Status:      StatusPending,
		Completed:   false,
		UserID:      input.UserID,
	}

	if err := DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create task", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": task})
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	var task Task
	if err := DB.Where("id = ?", id).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve task"})
		}
		return
	}
	if err := DB.Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete task"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": id})
}

func CompleteTask(c *gin.Context) {
	id := c.Param("id")
	var task Task
	if err := DB.Where("id = ?", id).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve task"})
		}
		return
	}
	task.Status = StatusCompleted
	task.Completed = true
	if err := DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not complete task"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": task})
}

// GetAllTasks retrieves all tasks with filtering, pagination, and sorting
// @Summary Get all tasks
// @Description Retrieve all tasks with optional filtering, pagination, and sorting
// @Tags tasks
// @Accept json
// @Produce json
// @Param userId query int false "Filter by user ID"
// @Param priority query string false "Filter by priority" Enums(low,medium,high)
// @Param status query string false "Filter by status" Enums(pending,completed,cancelled)
// @Param category query string false "Filter by category"
// @Param search query string false "Search in task name and description"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param sortBy query string false "Sort field" default(created_at)
// @Param sortOrder query string false "Sort order" Enums(asc,desc) default(desc)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tasks [get]
func GetAllTasks(c *gin.Context) {
	var query TaskQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters", "details": err.Error()})
		return
	}

	// Set defaults
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Limit <= 0 {
		query.Limit = 10
	}
	if query.Limit > 100 {
		query.Limit = 100 // Max limit
	}

	var tasks []Task
	queryBuilder := DB.Model(&Task{})

	// Apply filters
	if query.UserID != nil {
		queryBuilder = queryBuilder.Where("user_id = ?", *query.UserID)
	}
	if query.Priority != nil {
		queryBuilder = queryBuilder.Where("priority = ?", *query.Priority)
	}
	if query.Status != nil {
		queryBuilder = queryBuilder.Where("status = ?", *query.Status)
	}
	if query.Category != nil && *query.Category != "" {
		queryBuilder = queryBuilder.Where("category = ?", *query.Category)
	}
	if query.Search != "" {
		queryBuilder = queryBuilder.Where("task ILIKE ? OR description ILIKE ?", "%"+query.Search+"%", "%"+query.Search+"%")
	}

	// Apply sorting
	orderBy := query.SortBy
	if orderBy == "" {
		orderBy = "created_at"
	}
	if query.SortOrder == "asc" {
		orderBy += " ASC"
	} else {
		orderBy += " DESC"
	}
	queryBuilder = queryBuilder.Order(orderBy)

	// Get total count
	var total int64
	queryBuilder.Count(&total)

	// Apply pagination
	offset := (query.Page - 1) * query.Limit
	if err := queryBuilder.Offset(offset).Limit(query.Limit).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve tasks", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": tasks,
		"pagination": gin.H{
			"page":       query.Page,
			"limit":      query.Limit,
			"total":      total,
			"totalPages": (total + int64(query.Limit) - 1) / int64(query.Limit),
		},
	})
}

// GetTasksByUser retrieves tasks for a specific user
func GetTasksByUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var query TaskQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters", "details": err.Error()})
		return
	}

	// Override UserID with the path parameter
	query.UserID = &userID

	// Set defaults
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Limit <= 0 {
		query.Limit = 10
	}

	var tasks []Task
	queryBuilder := DB.Model(&Task{}).Where("user_id = ?", userID)

	// Apply filters
	if query.Priority != nil {
		queryBuilder = queryBuilder.Where("priority = ?", *query.Priority)
	}
	if query.Status != nil {
		queryBuilder = queryBuilder.Where("status = ?", *query.Status)
	}
	if query.Category != nil && *query.Category != "" {
		queryBuilder = queryBuilder.Where("category = ?", *query.Category)
	}
	if query.Search != "" {
		queryBuilder = queryBuilder.Where("task ILIKE ? OR description ILIKE ?", "%"+query.Search+"%", "%"+query.Search+"%")
	}

	// Apply sorting
	orderBy := query.SortBy
	if orderBy == "" {
		orderBy = "created_at"
	}
	if query.SortOrder == "asc" {
		orderBy += " ASC"
	} else {
		orderBy += " DESC"
	}
	queryBuilder = queryBuilder.Order(orderBy)

	// Get total count
	var total int64
	queryBuilder.Count(&total)

	// Apply pagination
	offset := (query.Page - 1) * query.Limit
	if err := queryBuilder.Offset(offset).Limit(query.Limit).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve tasks", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": tasks,
		"pagination": gin.H{
			"page":       query.Page,
			"limit":      query.Limit,
			"total":      total,
			"totalPages": (total + int64(query.Limit) - 1) / int64(query.Limit),
		},
	})
}

// UpdateTask updates an existing task
func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var input UpdateTaskRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Format", "details": err.Error()})
		return
	}

	var task Task
	if err := DB.Where("id = ?", id).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve task"})
		}
		return
	}

	// Update fields if provided
	if input.Task != nil {
		task.Task = *input.Task
	}
	if input.Description != nil {
		task.Description = *input.Description
	}
	if input.Priority != nil {
		task.Priority = *input.Priority
	}
	if input.Category != nil {
		task.Category = *input.Category
	}
	if input.DueDate != nil {
		task.DueDate = input.DueDate
	}
	if input.Status != nil {
		task.Status = *input.Status
		// Update completed field for backward compatibility
		task.Completed = (*input.Status == StatusCompleted)
	}

	if err := DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update task", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": task})
}

// GetTaskByID retrieves a single task by ID
func GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	var task Task
	if err := DB.Where("id = ?", id).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve task"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": task})
}
