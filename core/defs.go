package core

import (
	"errors"
	"time"
)

type (
	Direction uint8

	Operation struct {
		ID        uint64
		Time      time.Time
		Direction Direction
		Amount    uint64
	}

	Wallet struct {
		ID      uint64
		Balance int64
	}
)

var (
	ErrArg      = errors.New("incorrect argument")
	ErrNoMoney  = errors.New("not enough money")
	ErrInternal = errors.New("internal error")
)

const (
	DirAny Direction = iota
	DirOut
	DirIn
)
