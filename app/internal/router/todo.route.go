package router

import "go-todo/internal/handler"

// SetupTodoRoutes はTodo関連のルートを設定
func SetupTodoRoutes(r *Router, todoHandler *handler.TodoHandler) {
	r.GET("/todos", todoHandler.ListTodos)
	r.POST("/todos", todoHandler.CreateTodo)
	r.GET("/todos/{id}", todoHandler.GetTodo)
	r.PUT("/todos/{id}", todoHandler.UpdateTodo)
	r.DELETE("/todos/{id}", todoHandler.DeleteTodo)
}
