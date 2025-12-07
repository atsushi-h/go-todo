package seed

import (
	"context"
	"fmt"

	"go-todo/db/sqlc"
)

type Seeder interface {
	Seed(ctx context.Context) error
	Name() string
}

// 全シーダーを実行
func RunAll(ctx context.Context, queries *sqlc.Queries) error {
	seeders := []Seeder{
		NewTodoSeeder(queries),
	}

	for _, s := range seeders {
		fmt.Printf("🌱 Seeding %s...\n", s.Name())
		if err := s.Seed(ctx); err != nil {
			return fmt.Errorf("failed to seed %s: %w", s.Name(), err)
		}
		fmt.Printf("✅ %s seeded successfully\n", s.Name())
	}

	return nil
}
