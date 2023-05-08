package models

import (
	"time"
)

type Block struct {
	ID           uint64 `gorm:"primary_key"`
	Number       uint64
	Hash         string
	Time         uint64
	ParentHash   string
	Transactions []*Transaction `gorm:"foreignkey:BlockNumber"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (Block) TableName() string {
	return "blocks"
}
