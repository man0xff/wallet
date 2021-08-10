// +build test_env

package core_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	// "github.com/davecgh/go-spew/spew"

	"github.com/stretchr/testify/require"

	"example/core"
)

func newCore() *core.Core {
	return core.New(&core.Options{
		DBUser:     "app",
		DBPassword: "gfhjkm-app",
		DBAddress:  "database",
		DBName:     "example",
		Log:        log.New(os.Stdout, "", 0),
	})
}

func TestAll(t *testing.T) {
	ctx := context.Background()
	c := newCore()
	name := fmt.Sprintf("name%d", time.Now().UnixNano())
	id, err := c.AddWallet(ctx, name)
	require.Nil(t, err)
	require.NotZero(t, id)

	err = c.DepositMoney(ctx, id, 100)
	require.Nil(t, err)

	id2, err := c.AddWallet(ctx, name+"_")
	require.Nil(t, err)

	err = c.TransferMoney(ctx, id, id2, 100)
	require.Nil(t, err)

	ops, err := c.OperationsHistory(
		ctx, id,
		time.Now().Add(-2*time.Second),
		time.Now().Add(2*time.Second),
		core.DirAny)
	require.Nil(t, err)
	require.Equal(t, 2, len(ops))
}
