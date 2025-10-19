package models

import (
	"net/http"
	"strings"

	"github.com/KingLeak95/todo-list-go/pkg/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"not null" json:"name"`
	Email    string `gorm:"unique; not null" json:"email"`
	Password string `gorm:"not null" json:"-"` // Hidden from JSON
	Role     string `gorm:"default:'user'" json:"role"`
	Tasks    []Task `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"tasks,omitempty"`
}

type NewUser struct {
	Name     string `json:"name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func CreateUser(c *gin.Context) {
	var input NewUser
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Format", "details": err.Error()})
		return
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	newUser := User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
		Role:     "user", // Default role
	}

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

	// Remove password from response
	newUser.Password = ""
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

// Register creates a new user and returns authentication tokens
// @Summary Register a new user
// @Description Create a new user account and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param user body NewUser true "User registration data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var input NewUser
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Format", "details": err.Error()})
		return
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	newUser := User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
		Role:     "user",
	}

	result := DB.Create(&newUser)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "Duplicate") || strings.Contains(result.Error.Error(), "duplicate") {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		}
		return
	}

	// Generate JWT tokens
	jwtManager := auth.NewJWTManager()
	accessToken, refreshToken, err := jwtManager.GenerateTokenPair(newUser.ID, newUser.Email, newUser.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate tokens"})
		return
	}

	// Remove password from response
	newUser.Password = ""
	response := AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         newUser,
	}

	c.JSON(http.StatusCreated, gin.H{"data": response})
}

// Login authenticates a user and returns JWT tokens
// @Summary Login user
// @Description Authenticate user with email and password, return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var input LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Format", "details": err.Error()})
		return
	}

	// Find user by email
	var user User
	if err := DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Check password
	if !auth.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT tokens
	jwtManager := auth.NewJWTManager()
	accessToken, refreshToken, err := jwtManager.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate tokens"})
		return
	}

	// Remove password from response
	user.Password = ""
	response := AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// RefreshToken generates a new access token using a refresh token
func RefreshToken(c *gin.Context) {
	var input RefreshTokenRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Format", "details": err.Error()})
		return
	}

	jwtManager := auth.NewJWTManager()
	newAccessToken, err := jwtManager.RefreshToken(input.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"access_token": newAccessToken}})
}
