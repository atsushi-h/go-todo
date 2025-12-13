package main

import (
	"context"
	"log"

	"go-todo/db/sqlc"
	"go-todo/internal/auth"
	"go-todo/internal/database"
	"go-todo/internal/handler"
	"go-todo/internal/router"
	"go-todo/internal/service"

	"github.com/labstack/echo/v4"
)

func main() {
	ctx := context.Background()

	// DB接続プールの初期化
	pool, err := database.NewPool(ctx)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer pool.Close()

	// DBのヘルスチェック
	if err := database.HealthCheck(ctx, pool); err != nil {
		log.Fatal("Database health check failed:", err)
	}
	log.Println("Database connected successfully.")

	// sqlc Queriesの作成
	queries := sqlc.New(pool)

	// SessionManagerを作成
	sessionManager, err := auth.NewSessionManager()
	if err != nil {
		log.Fatal("Failed to initialize session:", err)
	}
	log.Println("Session store initialized.")

	// Gothicの初期化
	auth.InitProviders()
	auth.InitGothic(sessionManager)

	// サービスの初期化
	todoService := service.NewTodoService(queries)
	userService := service.NewUserService(queries, pool)

	// ハンドラーの初期化
	todoHandler := handler.NewTodoHandler(todoService)
	authHandler := handler.NewAuthHandler(userService, sessionManager)

	// APIHandlerの作成（StrictServerInterface実装）
	apiHandler := handler.NewAPIHandler(todoHandler)

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
