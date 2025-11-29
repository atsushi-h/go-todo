package seed

import (
	"context"

	"go-todo/internal/model"

	"github.com/uptrace/bun"
)

type TodoSeeder struct {
	db *bun.DB
}

func NewTodoSeeder(db *bun.DB) *TodoSeeder {
	return &TodoSeeder{db: db}
}

func (s *TodoSeeder) Name() string {
	return "todos"
}

func (s *TodoSeeder) Seed(ctx context.Context) error {
	todos := []model.Todo{
		{Title: "Goの基礎を学ぶ", Description: "A Tour of Goを完了する", Completed: true},
		{Title: "REST APIを構築", Description: "Todo APIを実装する", Completed: true},
		{Title: "テストを書く", Description: "ユニットテストとE2Eテストを追加", Completed: false},
		{Title: "認証機能を追加", Description: "JWT認証を実装する", Completed: false},
		{Title: "デプロイする", Description: "AWS/GCPにデプロイ", Completed: false},
	}

	_, err := s.db.NewInsert().
		Model(&todos).
		On("CONFLICT DO NOTHING"). // 重複時はスキップ
		Exec(ctx)

	return err
}
