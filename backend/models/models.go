package models

import (
	"time"

	"gorm.io/gorm"
)

// for each different stock
type Investment struct {
	gorm.Model
	Symbol          string    `gorm:"index"`
	LastTransaction time.Time `gorm:"type:date"`
	TotalShares     float64
	TotalInvested   float64
	Version         string `gorm:"index"`
}

// keep track of all transactions to track daily spend
type Transaction struct {
	gorm.Model
	Date              time.Time `gorm:"index;type:date"`
	TransactionsCount int
	TransactionAmount float64
	Version           string `gorm:"index"`
}

type Historical struct {
	Price  float64   `gorm:"column:price"`
	Date   time.Time `gorm:"column:date"`
	Symbol string    `gorm:"column:symbol"`
}

func (Historical) TableName() string {
	return "historical"
}

type MarketCap struct {
	Symbol    string  `grom:"column:symbol"`
	MarketCap float64 `gorm:column:market_cap"`
}

func (MarketCap) TableName() string {
	return "market_cap"
}

type PricePoint struct {
	Date  time.Time
	Price float64
}
