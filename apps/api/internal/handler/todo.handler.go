package handler

import (
	"context"

	"go-todo/internal/auth"
	"go-todo/internal/gen"
	"go-todo/internal/model"
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

	return gen.ListTodos200JSONResponse(convertTodosToGen(todos)), nil
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

	id := uint(request.Id)
	todo, err := h.service.GetTodoByID(ctx, id, userID)
	if err != nil {
		if err == service.ErrTodoNotFound {
			return gen.GetTodo404JSONResponse{Message: "Todo not found"}, nil
		}
		return gen.GetTodo500JSONResponse{Message: "Internal server error"}, nil
	}

	return gen.GetTodo200JSONResponse(convertTodoToGen(todo)), nil
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

	req := model.CreateTodoRequest{
		Title:       request.Body.Title,
		Description: ptrToString(request.Body.Description),
	}

	todo, err := h.service.CreateTodo(ctx, req, userID)
	if err != nil {
		return gen.CreateTodo500JSONResponse{Message: "Internal server error"}, nil
	}

	return gen.CreateTodo201JSONResponse(convertTodoToGen(todo)), nil
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

	id := uint(request.Id)
	req := model.UpdateTodoRequest{
		Title:       request.Body.Title,
		Description: request.Body.Description,
		Completed:   request.Body.Completed,
	}

	todo, err := h.service.UpdateTodo(ctx, id, userID, req)
	if err != nil {
		if err == service.ErrTodoNotFound {
			return gen.UpdateTodo404JSONResponse{Message: "Todo not found"}, nil
		}
		return gen.UpdateTodo500JSONResponse{Message: "Internal server error"}, nil
	}

	return gen.UpdateTodo200JSONResponse(convertTodoToGen(todo)), nil
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

	id := uint(request.Id)

	if err := h.service.DeleteTodo(ctx, id, userID); err != nil {
		if err == service.ErrTodoNotFound {
			return gen.DeleteTodo404JSONResponse{Message: "Todo not found"}, nil
		}
		return gen.DeleteTodo500JSONResponse{Message: "Internal server error"}, nil
	}

	return gen.DeleteTodo204Response{}, nil
}

// model.Todo から gen.Todo への変換
func convertTodoToGen(todo *model.Todo) gen.Todo {
	return gen.Todo{
		Id:          int64(todo.ID),
		Title:       todo.Title,
		Description: stringToPtr(todo.Description),
		Completed:   todo.Completed,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
		UserId:      int64(todo.UserID),
	}
}

// []model.Todo から []gen.Todo への変換
func convertTodosToGen(todos []*model.Todo) []gen.Todo {
	result := make([]gen.Todo, len(todos))
	for i, todo := range todos {
		result[i] = convertTodoToGen(todo)
	}
	return result
}

// ヘルパー: *string → string（nilの場合は空文字）
func ptrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// ヘルパー: string → *string（空文字の場合はnil）
func stringToPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
