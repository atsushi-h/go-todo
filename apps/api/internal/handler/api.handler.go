package handler

import (
	"context"

	"go-todo/internal/gen"
)

// APIHandler は StrictServerInterface を実装する
// TodoHandlerと他のハンドラーを統合したもの
type APIHandler struct {
	todoHandler *TodoHandler
}

// NewAPIHandler は新しいAPIHandlerを作成
func NewAPIHandler(todoHandler *TodoHandler) *APIHandler {
	return &APIHandler{
		todoHandler: todoHandler,
	}
}

// GetInfo - API情報を取得
func (h *APIHandler) GetInfo(ctx context.Context, request gen.GetInfoRequestObject) (gen.GetInfoResponseObject, error) {
	return gen.GetInfo200JSONResponse{
		Name:    "Todo API",
		Version: "1.0.0",
	}, nil
}

// GetHealth - ヘルスチェック
func (h *APIHandler) GetHealth(ctx context.Context, request gen.GetHealthRequestObject) (gen.GetHealthResponseObject, error) {
	return gen.GetHealth200JSONResponse{
		Status: "ok",
	}, nil
}

// ListTodos - TodoHandlerに委譲
func (h *APIHandler) ListTodos(ctx context.Context, request gen.ListTodosRequestObject) (gen.ListTodosResponseObject, error) {
	return h.todoHandler.ListTodos(ctx, request)
}

// GetTodo - TodoHandlerに委譲
func (h *APIHandler) GetTodo(ctx context.Context, request gen.GetTodoRequestObject) (gen.GetTodoResponseObject, error) {
	return h.todoHandler.GetTodo(ctx, request)
}

// CreateTodo - TodoHandlerに委譲
func (h *APIHandler) CreateTodo(ctx context.Context, request gen.CreateTodoRequestObject) (gen.CreateTodoResponseObject, error) {
	return h.todoHandler.CreateTodo(ctx, request)
}

// UpdateTodo - TodoHandlerに委譲
func (h *APIHandler) UpdateTodo(ctx context.Context, request gen.UpdateTodoRequestObject) (gen.UpdateTodoResponseObject, error) {
	return h.todoHandler.UpdateTodo(ctx, request)
}

// DeleteTodo - TodoHandlerに委譲
func (h *APIHandler) DeleteTodo(ctx context.Context, request gen.DeleteTodoRequestObject) (gen.DeleteTodoResponseObject, error) {
	return h.todoHandler.DeleteTodo(ctx, request)
}

// BatchCompleteTodos - TodoHandlerに委譲
func (h *APIHandler) BatchCompleteTodos(ctx context.Context, request gen.BatchCompleteTodosRequestObject) (gen.BatchCompleteTodosResponseObject, error) {
	return h.todoHandler.BatchCompleteTodos(ctx, request)
}

// BatchDeleteTodos - TodoHandlerに委譲
func (h *APIHandler) BatchDeleteTodos(ctx context.Context, request gen.BatchDeleteTodosRequestObject) (gen.BatchDeleteTodosResponseObject, error) {
	return h.todoHandler.BatchDeleteTodos(ctx, request)
}

// コンパイル時にStrictServerInterfaceを実装していることを確認
var _ gen.StrictServerInterface = (*APIHandler)(nil)
