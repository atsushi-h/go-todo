package database

import (
	"context"
	"errors"
	"testing"

	"go-todo/internal/config"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTxManager_RunInTx_Success(t *testing.T) {
	t.Skip("Integration test - requires database")

	ctx := context.Background()
	cfg := config.DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		Database: "go_todo_db",
		User:     "user",
		Password: "pass",
	}
	pool, err := NewPool(ctx, cfg)
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
	cfg := config.DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		Database: "go_todo_db",
		User:     "user",
		Password: "pass",
	}
	pool, err := NewPool(ctx, cfg)
	require.NoError(t, err)
	defer pool.Close()

	tm := NewTxManager(pool)
	testErr := errors.New("test error")

	err = tm.RunInTx(ctx, func(tx pgx.Tx) error {
		return testErr
	})

	assert.ErrorIs(t, err, testErr)
}
