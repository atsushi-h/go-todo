package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func Init() (*bun.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&timezone=Asia/Tokyo",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	// PostgreSQL connector
	connector := pgdriver.NewConnector(pgdriver.WithDSN(dsn))
	sqldb := sql.OpenDB(connector)

	// Bun DB with PostgreSQL dialect
	db := bun.NewDB(sqldb, pgdialect.New())

	// Add query hook for debugging (optional, can be removed in production)
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	return db, nil
}

// データベース接続の確認
func HealthCheck(db *bun.DB) error {
	sqlDB := db.DB
	return sqlDB.Ping()
}

// DBを閉じる
func Close(db *bun.DB) error {
	return db.Close()
}
