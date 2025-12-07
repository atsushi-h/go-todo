package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"go-todo/db/sqlc"
	"go-todo/internal/database"
	"go-todo/internal/seed"
)

func main() {
	// フラグ定義
	fresh := flag.Bool("fresh", false, "テーブルをクリアしてからシード")
	flag.Parse()

	ctx := context.Background()

	// DB接続
	pool, err := database.NewPool(ctx)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer pool.Close()

	// freshフラグがある場合はデータを削除
	if *fresh {
		fmt.Println("🗑️  Clearing existing data...")
		_, err := pool.Exec(ctx, "TRUNCATE TABLE todos CASCADE")
		if err != nil {
			// テーブルが存在しない場合は無視
			fmt.Printf("⚠️  Warning: %v\n", err)
		}
	}

	// sqlc Queriesの作成
	queries := sqlc.New(pool)

	// シード実行
	fmt.Println("🌱 Starting database seeding...")
	if err := seed.RunAll(ctx, queries); err != nil {
		log.Fatal("Seeding failed:", err)
	}

	fmt.Println("🎉 Database seeding completed!")
}
