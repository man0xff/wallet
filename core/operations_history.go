package core

import (
	"context"
	"time"
)

func (c *Core) OperationsHistory(
	ctx context.Context, id uint64, from, to time.Time, dir Direction,
) ([]Operation, error) {
	db, err := c.db(ctx)
	if err != nil {
		return nil, err
	}

	var operations []operation
	err = db.Where("wallet_id = ? and time >= ? and time <= ?", id, from, to).
		Find(&operations).Error
	if err != nil {
		return nil, newErrInternal(err)
	}

	out := make([]Operation, 0, len(operations))
	for _, op := range operations {
		curDir := Direction(op.Direction)
		if dir == DirAny || curDir == dir {
			out = append(out, Operation{
				ID:        op.ID,
				Time:      op.Time,
				Direction: curDir,
				Amount:    op.Amount,
			})
		}
	}
	return out, nil
}
