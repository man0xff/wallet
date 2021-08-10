package core

import (
	"context"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type Core struct {
	db    func(context.Context) (*gorm.DB, error)
	log   *log.Logger
	dbDSN string
}

type Options struct {
	DBUser     string
	DBPassword string
	DBAddress  string
	DBName     string

	Log *log.Logger
}

func New(options *Options) *Core {
	c := &Core{
		log: options.Log,
	}
	c.initDB(options)
	return c
}

func wrapErrInternal(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%w: %v", ErrInternal, err)
}

func newErrInternal(err error) error {
	return fmt.Errorf("%w: %v", ErrInternal, err)
}
