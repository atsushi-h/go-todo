package main

import (
	"fmt"
	"log"
	"net/http"

	"go-todo/handler"
	"go-todo/repository"
	"go-todo/router"
	"go-todo/service"
)

func main() {
	// リポジトリの初期化
	todoRepo := repository.NewTodoRepository()

	// サービスの初期化
	todoService := service.NewTodoService(todoRepo)

	// ハンドラーの初期化
	todoHandler := handler.NewTodoHandler(todoService)

	// ルーターの初期化
	r := router.NewRouter()
	router.SetupRoutes(r, todoHandler)

	// サーバー起動
	fmt.Println("Server starting on :4000...")
	if err := http.ListenAndServe(":4000", r); err != nil {
		log.Fatal(err)
	}
}
