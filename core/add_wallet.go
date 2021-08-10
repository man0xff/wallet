package core

import (
	"context"
	"fmt"
)

func (c *Core) AddWallet(ctx context.Context, name string) (uint64, error) {
	if name == "" {
		return 0, fmt.Errorf("%w: no wallet name", ErrArg)
	}

	db, err := c.db(ctx)
	if err != nil {
		return 0, err
	}

	w := &wallet{Name: name}
	if err := db.Create(w).Error; err != nil {
		return 0, err
	}
	return w.ID, err
}
