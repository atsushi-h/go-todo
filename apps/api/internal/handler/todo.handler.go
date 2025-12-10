package handler

import (
	"context"

	"go-todo/internal/auth"
	"go-todo/internal/gen"
	"go-todo/internal/mapper"
	"go-todo/internal/service"
)

// TodoのHTTPハンドラー（StrictServerInterface実装）
type TodoHandler struct {
	service *service.TodoService
}

// 新しいTodoHandlerを作成
func NewTodoHandler(service *service.TodoService) *TodoHandler {
	return &TodoHandler{
		service: service,
	}
}

// ListTodos - 全Todoを取得
func (h *TodoHandler) ListTodos(ctx context.Context, request gen.ListTodosRequestObject) (gen.ListTodosResponseObject, error) {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return gen.ListTodos401JSONResponse{Message: "Unauthorized"}, nil
	}

	todos, err := h.service.GetAllTodos(ctx, userID)
	if err != nil {
		return gen.ListTodos500JSONResponse{Message: "Internal server error"}, nil
	}

	return gen.ListTodos200JSONResponse(mapper.TodosToResponse(todos)), nil
}

// GetTodo - IDでTodoを取得
func (h *TodoHandler) GetTodo(ctx context.Context, request gen.GetTodoRequestObject) (gen.GetTodoResponseObject, error) {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return gen.GetTodo401JSONResponse{Message: "Unauthorized"}, nil
	}

	if request.Id < 0 {
		return gen.GetTodo400JSONResponse{Message: "Invalid ID"}, nil
	}

	todo, err := h.service.GetTodoByID(ctx, int64(request.Id), userID)
	if err != nil {
		if err == service.ErrTodoNotFound {
			return gen.GetTodo404JSONResponse{Message: "Todo not found"}, nil
		}
		return gen.GetTodo500JSONResponse{Message: "Internal server error"}, nil
	}

	return gen.GetTodo200JSONResponse(mapper.TodoToResponse(todo)), nil
}

// CreateTodo - 新しいTodoを作成
func (h *TodoHandler) CreateTodo(ctx context.Context, request gen.CreateTodoRequestObject) (gen.CreateTodoResponseObject, error) {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return gen.CreateTodo401JSONResponse{Message: "Unauthorized"}, nil
	}

	if request.Body == nil {
		return gen.CreateTodo400JSONResponse{Message: "Invalid request body"}, nil
	}

	if request.Body.Title == "" {
		return gen.CreateTodo400JSONResponse{Message: "Title is required"}, nil
	}

	todo, err := h.service.CreateTodo(ctx, userID, request.Body.Title, request.Body.Description)
	if err != nil {
		return gen.CreateTodo500JSONResponse{Message: "Internal server error"}, nil
	}

	return gen.CreateTodo201JSONResponse(mapper.TodoToResponse(todo)), nil
}

// UpdateTodo - Todoを更新
func (h *TodoHandler) UpdateTodo(ctx context.Context, request gen.UpdateTodoRequestObject) (gen.UpdateTodoResponseObject, error) {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return gen.UpdateTodo401JSONResponse{Message: "Unauthorized"}, nil
	}

	if request.Id < 0 {
		return gen.UpdateTodo400JSONResponse{Message: "Invalid ID"}, nil
	}

	if request.Body == nil {
		return gen.UpdateTodo400JSONResponse{Message: "Invalid request body"}, nil
	}

	todo, err := h.service.UpdateTodo(
		ctx,
		int64(request.Id),
		userID,
		request.Body.Title,
		request.Body.Description,
		request.Body.Completed,
	)
	if err != nil {
		if err == service.ErrTodoNotFound {
			return gen.UpdateTodo404JSONResponse{Message: "Todo not found"}, nil
		}
		return gen.UpdateTodo500JSONResponse{Message: "Internal server error"}, nil
	}

	return gen.UpdateTodo200JSONResponse(mapper.TodoToResponse(todo)), nil
}

// DeleteTodo - Todoを削除
func (h *TodoHandler) DeleteTodo(ctx context.Context, request gen.DeleteTodoRequestObject) (gen.DeleteTodoResponseObject, error) {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return gen.DeleteTodo401JSONResponse{Message: "Unauthorized"}, nil
	}

	if request.Id < 0 {
		return gen.DeleteTodo400JSONResponse{Message: "Invalid ID"}, nil
	}

	if err := h.service.DeleteTodo(ctx, int64(request.Id), userID); err != nil {
		if err == service.ErrTodoNotFound {
			return gen.DeleteTodo404JSONResponse{Message: "Todo not found"}, nil
		}
		return gen.DeleteTodo500JSONResponse{Message: "Internal server error"}, nil
	}

	return gen.DeleteTodo204Response{}, nil
}

// BatchCompleteTodos - Todoを一括完了
func (h *TodoHandler) BatchCompleteTodos(ctx context.Context, request gen.BatchCompleteTodosRequestObject) (gen.BatchCompleteTodosResponseObject, error) {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return gen.BatchCompleteTodos401JSONResponse{Message: "Unauthorized"}, nil
	}

	if request.Body == nil || len(request.Body.Ids) == 0 {
		return gen.BatchCompleteTodos400JSONResponse{Message: "IDs are required"}, nil
	}

	if len(request.Body.Ids) > 100 {
		return gen.BatchCompleteTodos400JSONResponse{Message: "Too many IDs (max 100)"}, nil
	}

	result, err := h.service.BatchCompleteTodos(ctx, userID, request.Body.Ids)
	if err != nil {
		return gen.BatchCompleteTodos500JSONResponse{Message: "Internal server error"}, nil
	}

	return gen.BatchCompleteTodos200JSONResponse{
		Succeeded: mapper.TodosToResponse(result.Succeeded),
		Failed:    mapper.BatchFailedItemsToResponse(result.Failed),
	}, nil
}

// BatchDeleteTodos - Todoを一括削除
func (h *TodoHandler) BatchDeleteTodos(ctx context.Context, request gen.BatchDeleteTodosRequestObject) (gen.BatchDeleteTodosResponseObject, error) {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return gen.BatchDeleteTodos401JSONResponse{Message: "Unauthorized"}, nil
	}

	if request.Body == nil || len(request.Body.Ids) == 0 {
		return gen.BatchDeleteTodos400JSONResponse{Message: "IDs are required"}, nil
	}

	if len(request.Body.Ids) > 100 {
		return gen.BatchDeleteTodos400JSONResponse{Message: "Too many IDs (max 100)"}, nil
	}

	result, err := h.service.BatchDeleteTodos(ctx, userID, request.Body.Ids)
	if err != nil {
		return gen.BatchDeleteTodos500JSONResponse{Message: "Internal server error"}, nil
	}

	return gen.BatchDeleteTodos200JSONResponse{
		Succeeded: result.Succeeded,
		Failed:    mapper.BatchFailedItemsToResponse(result.Failed),
	}, nil
}
