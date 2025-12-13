package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TxManager interface {
	RunInTx(ctx context.Context, fn func(tx pgx.Tx) error) error
}

type txManager struct {
	pool *pgxpool.Pool
}

func NewTxManager(pool *pgxpool.Pool) TxManager {
	return &txManager{pool: pool}
}

func (m *txManager) RunInTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		tx.Rollback(ctx) // 失敗したらロールバック
		return err
	}

	return tx.Commit(ctx) // 成功したらコミット
}
