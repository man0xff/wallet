package core

import (
	"context"
	"fmt"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func (c *Core) initDB(options *Options) {
	c.dbDSN = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		options.DBUser,
		options.DBPassword,
		options.DBAddress,
		options.DBName,
	)
	c.db = c.dbClosure()
}

func (c *Core) connectToDB() (*gorm.DB, error) {
	logger := logger.New(c.log, logger.Config{LogLevel: logger.Info})
	db, err := gorm.Open(mysql.Open(c.dbDSN), &gorm.Config{
		Logger:                 logger,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: database connecting: %v", ErrInternal, err)
	}
	return db, nil
}

// позволяет стартовать приложению даже если соединения с базой пока нет;
// коннект к базе данных не сохраняем, чтобы всегда получать коннект через
// замыкание
func (c *Core) dbClosure() func(context.Context) (*gorm.DB, error) {
	var (
		db    *gorm.DB
		mutex sync.Mutex
	)
	return func(ctx context.Context) (*gorm.DB, error) {
		var err error

		mutex.Lock()
		defer mutex.Unlock()

		if db == nil {
			db, err = c.connectToDB()
			if err != nil {
				return nil, err
			}
		}
		return db.WithContext(ctx), nil
	}
}

func (c *Core) insertOperation(
	db *gorm.DB,
	id uint64,
	dir Direction,
	amount uint64,
) error {
	op := &operation{
		WalletID:  id,
		Direction: uint8(dir),
		Amount:    amount,
	}
	err := db.Omit("Time").Create(op).Error
	return wrapErrInternal(err)
}

func transaction(db *gorm.DB, fn func(*gorm.DB) error) error {
	db.Exec("begin")
	if db.Error != nil {
		return wrapErrInternal(db.Error)
	}
	if err := fn(db); err != nil {
		db.Exec("rollback")
		return err
	}
	if err := db.Exec("commit").Error; err != nil {
		db.Exec("rollback")
		return wrapErrInternal(err)
	}
	return nil
}
