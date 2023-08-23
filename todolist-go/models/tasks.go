package models

import (
  "gorm.io/gorm"
)

type Task struct {
  gorm.Model
  task string 
  user_id uint 
} 
