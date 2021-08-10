package core

import (
	"context"

	"gorm.io/gorm"
)

func (c *Core) DepositMoney(
	ctx context.Context, id uint64, amount uint64,
) error {
	db, err := c.db(ctx)
	if err != nil {
		return err
	}

	err = transaction(db, func(tx *gorm.DB) error {
		err := tx.Model(&wallet{ID: id}).
			Update("balance", gorm.Expr("balance + ?", amount)).Error
		if err != nil {
			return wrapErrInternal(err)
		}
		return c.insertOperation(tx, id, DirIn, amount)
	})
	return err
}
