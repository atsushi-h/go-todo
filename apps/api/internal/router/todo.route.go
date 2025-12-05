package router

import (
	"go-todo/internal/auth"
	"go-todo/internal/handler"
)

// SetupTodoRoutes はTodo関連のルートを設定
func SetupTodoRoutes(r *Router, todoHandler *handler.TodoHandler, sm *auth.SessionManager) {
	requireAuth := auth.RequireAuth(sm) // 認証ミドルウェア

	// Todoルートの登録
	r.GET("/todos", requireAuth(todoHandler.ListTodos))
	r.POST("/todos", requireAuth(todoHandler.CreateTodo))
	r.GET("/todos/{id}", requireAuth(todoHandler.GetTodo))
	r.PUT("/todos/{id}", requireAuth(todoHandler.UpdateTodo))
	r.DELETE("/todos/{id}", requireAuth(todoHandler.DeleteTodo))
}
