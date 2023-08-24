package models

import (
  "gorm.io/gorm"
)

type User struct {
  gorm.Model
  name string  `gorm:"type:varchar(255);not null"`
  email string `gorm:"uniqueIndex;not null"`
  tasks []Task
} 
