package models

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name  string `gorm:"not null"`
	Email string `gorm:"unique; not null"`
	Tasks []Task `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

type NewUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func CreateUser(c *gin.Context) {
	var input NewUser
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Format"})
		return
	}

	// Validate user name
	if len(input.Name) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Name is invalid: %s", input.Name)})
	}

	// Validate email
	if !strings.Contains(input.Email, "@") {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Email is invalid: %s", input.Email)})
	}

	newUser := User{Name: input.Name, Email: input.Email}

	result := DB.Create(&newUser)

	if result.Error != nil {
		// Check for duplicate email error
		if strings.Contains(result.Error.Error(), "Duplicate") || strings.Contains(result.Error.Error(), "duplicate") {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": newUser})
}

func GetAllUsers(c *gin.Context) {
	var allUsers []User
	if err := DB.Find(&allUsers).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": allUsers})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	var user User
	if err := DB.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve user"})
		}
		return
	}

	if err := DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": id})
}
