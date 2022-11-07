package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string    `gorm:"unique" json:"email"`
	PasswordHash string    `json:"password"`
	Status       string    `gorm:"default:'active'" json:"status"`
	DeletedTime  time.Time `gorm:"default:null" json:"deletedtime"`
}
