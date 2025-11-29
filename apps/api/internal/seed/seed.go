package seed

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

type Seeder interface {
	Seed(ctx context.Context) error
	Name() string
}

// 全シーダーを実行
func RunAll(ctx context.Context, db *bun.DB) error {
	seeders := []Seeder{
		NewTodoSeeder(db),
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
