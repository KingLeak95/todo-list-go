package models

import (
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Task      string `gorm:"column:name;type:varchar(255);not null"`
	Completed bool
	UserID    int `gorm:"foreignKey:user_id"`
}
