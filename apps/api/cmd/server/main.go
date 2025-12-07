package main

import (
	"log"

	"go-todo/internal/auth"
	"go-todo/internal/database"
	"go-todo/internal/handler"
	"go-todo/internal/repository"
	"go-todo/internal/router"
	"go-todo/internal/service"

	"github.com/labstack/echo/v4"
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

	// SessionManagerを作成
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
	apiHandler := handler.NewAPIHandler(todoService)
	authHandler := handler.NewAuthHandler(userService, sessionManager)

	// Echoインスタンスを作成
	e := echo.New()

	// ルートを設定
	router.SetupRoutes(e, apiHandler, authHandler, sessionManager)

	// サーバー起動
	log.Println("Server starting on :4000...")
	if err := e.Start(":4000"); err != nil {
		log.Fatal(err)
	}
}
