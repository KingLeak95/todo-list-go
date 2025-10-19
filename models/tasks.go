package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Task      string `gorm:"column:name;type:varchar(255);not null"`
	Completed bool
	UserID    int `gorm:"foreignKey:user_id"`
}

type NewTask struct {
	Task   string `json:"task"`
	UserID int    `json:"userId"`
}

func CreateTask(c *gin.Context) {
	var input NewTask
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Format"})
		return
	}
	if input.Task == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task is required"})
		return
	}
	task := Task{Task: input.Task, Completed: false, UserID: input.UserID}
	if err := DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create task"})
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
	task.Completed = true
	if err := DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not complete task"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": task})
}
