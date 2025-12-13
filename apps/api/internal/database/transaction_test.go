package database

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTxManager_RunInTx_Success(t *testing.T) {
	t.Skip("Integration test - requires database")

	ctx := context.Background()
	pool, err := NewPool(ctx)
	require.NoError(t, err)
	defer pool.Close()

	tm := NewTxManager(pool)

	var capturedTx pgx.Tx
	err = tm.RunInTx(ctx, func(tx pgx.Tx) error {
		capturedTx = tx
		return nil
	})

	assert.NoError(t, err)
	assert.NotNil(t, capturedTx)
}

func TestTxManager_RunInTx_Rollback(t *testing.T) {
	t.Skip("Integration test - requires database")

	ctx := context.Background()
	pool, err := NewPool(ctx)
	require.NoError(t, err)
	defer pool.Close()

	tm := NewTxManager(pool)
	testErr := errors.New("test error")

	err = tm.RunInTx(ctx, func(tx pgx.Tx) error {
		return testErr
	})

	assert.ErrorIs(t, err, testErr)
}
