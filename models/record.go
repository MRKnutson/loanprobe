package models

import (
	"time"

	"gorm.io/gorm"
)

type Record struct {
	gorm.Model
	OperationRefer    int       `json:"operationid"`
	UserRefer         int       `json:"userid"`
	Amount            float64   `json:"amount"`
	UserBalance       float64   `json:"userbalance"`
	DeletedTime       time.Time `gorm:"default:null" json:"deletedtime"`
	OperationResponse string    `json:"operationresponse"`
	Operation         Operation `gorm:"foreignKey:OperationRefer"`
	User              User      `gorm:"foreignKey:UserRefer"`
}
