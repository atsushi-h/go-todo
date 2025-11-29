package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"go-todo/internal/database"
	"go-todo/internal/model"
	"go-todo/internal/seed"
)

func main() {
	// フラグ定義
	fresh := flag.Bool("fresh", false, "テーブルをクリアしてからシード")
	flag.Parse()

	// DB接続
	db, err := database.Init()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close(db)

	ctx := context.Background()

	// freshフラグがある場合はデータを削除
	if *fresh {
		fmt.Println("🗑️  Clearing existing data...")
		if _, err := db.NewTruncateTable().
			Model((*model.Todo)(nil)).
			Cascade().
			Exec(ctx); err != nil {
			// テーブルが存在しない場合は無視
			fmt.Printf("⚠️  Warning: %v\n", err)
		}
	}

	// シード実行
	fmt.Println("🌱 Starting database seeding...")
	if err := seed.RunAll(ctx, db); err != nil {
		log.Fatal("Seeding failed:", err)
	}

	fmt.Println("🎉 Database seeding completed!")
}
