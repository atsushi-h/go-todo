package mapper

import (
	"go-todo/db/sqlc"
	"go-todo/internal/gen"
	"go-todo/internal/service"
)

func TodoToResponse(t *sqlc.Todo) gen.Todo {
	return gen.Todo{
		Id:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Completed:   t.Completed,
		UserId:      t.UserID,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

func TodosToResponse(todos []sqlc.Todo) []gen.Todo {
	result := make([]gen.Todo, len(todos))
	for i := range todos {
		result[i] = TodoToResponse(&todos[i])
	}
	return result
}

func BatchFailedItemsToResponse(items []service.BatchFailedItem) []gen.BatchFailedItem {
	result := make([]gen.BatchFailedItem, len(items))
	for i, item := range items {
		result[i] = gen.BatchFailedItem{
			Id:    item.ID,
			Error: item.Error,
		}
	}
	return result
}
