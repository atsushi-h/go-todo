package main

import (
	"fmt"
	"log"
	"net/http"

	"go-todo/internal/database"
	"go-todo/internal/handler"
	"go-todo/internal/repository"
	"go-todo/internal/router"
	"go-todo/internal/service"
)

func main() {
	// DBの初期化
	db, err := database.Init()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close(db)
	// DBのヘルスチェック
	if err := database.HealthCheck(db); err != nil {
		log.Fatal("Database health check failed:", err)
	}
	log.Println("Database connected successfully.")

	// リポジトリの初期化
	todoRepo := repository.NewTodoRepository(db)

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
