package main

import (
	"fmt"
	"log"
	"net/http"

	"go-todo/internal/auth"
	"go-todo/internal/database"
	"go-todo/internal/handler"
	"go-todo/internal/repository"
	"go-todo/internal/router"
	"go-todo/internal/service"
)

// @title Todo API
// @version 0.0.0
// @description A simple Todo API built with Go and PostgreSQL
// @host localhost:4000
// @BasePath /
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

	// SessionManagerを作成（値として受け取る）
	sessionManager, err := auth.NewSessionManager()
	if err != nil {
		log.Fatal("Failed to initialize session:", err)
	}
	log.Println("Session store initialized.")

	// Gothicに渡す
	auth.InitProviders()
	auth.InitGothic(sessionManager)

	// リポジトリの初期化
	todoRepo := repository.NewTodoRepository(db)
	userRepo := repository.NewUserRepository(db)

	// サービスの初期化
	todoService := service.NewTodoService(todoRepo)
	userService := service.NewUserService(userRepo)

	// ハンドラーの初期化
	todoHandler := handler.NewTodoHandler(todoService)
	authHandler := handler.NewAuthHandler(userService, sessionManager)

	// ルーターの初期化
	r := router.NewRouter()
	router.SetupRoutes(r, todoHandler)
	router.SetupAuthRoutes(r, authHandler, sessionManager)

	// サーバー起動
	fmt.Println("Server starting on :4000...")
	if err := http.ListenAndServe(":4000", r); err != nil {
		log.Fatal(err)
	}
}
