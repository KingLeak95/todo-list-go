package models

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strconv"
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
		host:     getEnv("DB_HOST", "localhost"),
		dbname:   getEnv("DB_NAME", "todolist"),
		user:     getEnv("DB_USER", "postgres"),
		password: getEnv("DB_PASSWORD", "postgres"),
		port:     getEnvInt("DB_PORT", 5432),
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

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if valueStr, ok := os.LookupEnv(key); ok {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return fallback
}
