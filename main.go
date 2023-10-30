package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/KingLeak95/todo-list-go/models"
)

func createUser(c *gin.Context) {

}

func main() {
	r := gin.Default()

	// Enable Middle Ware from Gin Package
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Connect Database
	models.ConnectDatabase()

	// Index for Testing
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})

	// Users
	r.POST("/createUser", models.CreateUser)
	r.GET("/allUsers", models.GetAllUsers)
	r.DELETE("/deleteUser/:id", models.DeleteUser)

	r.Run()
}
