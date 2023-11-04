package models

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type connection struct {
	host     string 
	dbname   string 
	user     string
	password string
	port     int  
}

func NewDBConnection() *connection {
	return &connection{
		host:     "localhost",
		dbname:   "todolist",
		user:     "postgres",
		password: "postgres",
		port:     5432,
	}
}

func ConnectDatabase() {
	postgresConnection := NewDBConnection()
	postgresDsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		postgresConnection.host,
		postgresConnection.user,
		postgresConnection.password,
		postgresConnection.dbname,
		postgresConnection.port)
	database, err := gorm.Open(postgres.Open(postgresDsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	err = database.AutoMigrate(&User{}, &Task{})
	if err != nil {
		return
	}

	DB = database
}
