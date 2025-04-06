package models

import (
	"time"

	"gorm.io/gorm"
)

type Investments struct {
	gorm.Model
	Symbol          string    `gorm:"index"`
	LastTransaction time.Time `gorm:"type:date"`
	CurrentAmount   float64
	CurrentProfit   float64
	Version         string `gorm:"index"`
}

type Transactions struct {
	gorm.Model
	Date              time.Time `gorm:"index;type:date"`
	TransactionsCount int
	TransactionAmount float64
	Version           string `gorm:"index"`
}
