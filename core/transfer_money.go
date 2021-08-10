package core

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (c *Core) TransferMoney(
	ctx context.Context,
	from, to, amount uint64,
) error {
	db, err := c.db(ctx)
	if err != nil {
		return err
	}

	err = transaction(db, func(tx *gorm.DB) error {
		var walletFrom wallet
		err := tx.Clauses(clause.Locking{Strength: "update"}).
			First(&walletFrom, from).Error
		if err != nil {
			return newErrInternal(err)
		}
		if walletFrom.Balance < amount {
			return fmt.Errorf("%w (id:%d)", ErrNoMoney, from)
		}
		err = tx.Model(&wallet{ID: to}).
			Update("balance", gorm.Expr("balance + ?", amount)).Error
		if err != nil {
			return newErrInternal(err)
		}
		err = tx.Model(&wallet{ID: from}).
			Update("balance", gorm.Expr("balance - ?", amount)).Error
		if err != nil {
			return newErrInternal(err)
		}
		err = c.insertOperation(tx, from, DirOut, amount)
		if err != nil {
			return wrapErrInternal(err)
		}
		return c.insertOperation(tx, to, DirIn, amount)
	})
	return err
}
