package seed

import (
	"context"

	"go-todo/db/sqlc"
)

type TodoSeeder struct {
	queries *sqlc.Queries
}

func NewTodoSeeder(queries *sqlc.Queries) *TodoSeeder {
	return &TodoSeeder{queries: queries}
}

func (s *TodoSeeder) Name() string {
	return "todos"
}

func (s *TodoSeeder) Seed(ctx context.Context) error {
	// シードデータにはユーザーIDが必要なため、
	// 実際の運用では先にユーザーを作成するか、
	// テスト用ユーザーIDを指定する必要がある
	//
	// 現在はスキップ（OAuth認証後にユーザーが作成されるため）
	return nil
}
