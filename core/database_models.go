package core

import (
	"time"
)

type wallet struct {
	ID      uint64 `gorm:"primaryKey"`
	Name    string
	Balance uint64
}

func (*wallet) TableName() string {
	return "wallets"
}

type operation struct {
	ID        uint64 `gorm:"primaryKey"`
	Time      time.Time
	WalletID  uint64
	Direction uint8
	Amount    uint64
}

func (*operation) TableName() string {
	return "operations"
}
