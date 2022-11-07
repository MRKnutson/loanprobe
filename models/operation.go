package models

import (
	"gorm.io/gorm"

	"time"
)

type Operation struct {
	gorm.Model
	Type        string    `json:"type"`
	Cost        float64   `json:"cost"`
	DeletedTime time.Time `gorm:"default:null" json:"deletedtime"`
}
