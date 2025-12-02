package main

import (
	"fmt"
	"os"

	"go-todo/internal/model"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func main() {
	// ダミーのDB接続（実際には接続しない）
	db := bun.NewDB(nil, pgdialect.New())

	// スキーマSQL生成
	var sql string
	
	// Todoテーブルのスキーマ生成
	sql += generateCreateTableSQL(db, (*model.Todo)(nil))
	
	// Userテーブルのスキーマ生成
	sql += generateCreateTableSQL(db, (*model.User)(nil))
	
	// schema-gen.sqlに出力
	err := os.WriteFile("schema-gen.sql", []byte(sql), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write schema-gen.sql: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("✅ schema-gen.sql generated successfully")
}

func generateCreateTableSQL(db *bun.DB, model interface{}) string {
	query := db.NewCreateTable().
		Model(model).
		IfNotExists()
	
	sql := query.String()
	return string(sql) + ";\n\n"
}
