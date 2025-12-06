package router

import (
	"go-todo/internal/auth"
	"go-todo/internal/handler"

	"github.com/labstack/echo/v4"
)

// Todo関連のルートを設定
func SetupTodoRoutes(e *echo.Echo, todoHandler *handler.TodoHandler, sm *auth.SessionManager) {
	// ルートグループを使用（認証が必要なルート）
	todosGroup := e.Group("/todos")
	todosGroup.Use(auth.RequireAuth(sm)) // グループにミドルウェアを適用

	// パラメータは :id 形式（Echo標準）
	todosGroup.GET("", todoHandler.ListTodos)
	todosGroup.POST("", todoHandler.CreateTodo)
	todosGroup.GET("/:id", todoHandler.GetTodo)
	todosGroup.PUT("/:id", todoHandler.UpdateTodo)
	todosGroup.DELETE("/:id", todoHandler.DeleteTodo)
}
