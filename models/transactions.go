package models

import "time"

type Transaction struct {
	ID          uint64 `gorm:"primary_key"`
	BlockNumber uint64
	Hash        string
	To          string
	Value       string
	Nonce       uint64
	Data        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Transaction) TableName() string {
	return "transactions"
}
