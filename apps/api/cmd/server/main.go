package main

import (
	"context"
	"log"

	"go-todo/db/sqlc"
	"go-todo/internal/auth"
	"go-todo/internal/config"
	"go-todo/internal/database"
	"go-todo/internal/handler"
	"go-todo/internal/router"
	"go-todo/internal/service"

	"github.com/labstack/echo/v4"
)

func main() {
	ctx := context.Background()

	// 環境変数から設定を読み込み
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// DB接続プールの初期化
	pool, err := database.NewPool(ctx, cfg.Database)
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
	sessionManager, err := auth.NewSessionManager(cfg.Redis, cfg.Cookie)
	if err != nil {
		log.Fatal("Failed to initialize session:", err)
	}
	log.Println("Session store initialized.")

	// Gothicの初期化
	auth.InitProviders(cfg.OAuth)
	auth.InitGothic(sessionManager)

	// サービスの初期化
	todoService := service.NewTodoService(queries)
	userService := service.NewUserService(queries, pool)

	// ハンドラーの初期化
	todoHandler := handler.NewTodoHandler(todoService)
	authHandler := handler.NewAuthHandler(userService, sessionManager, cfg.Frontend)

	// APIHandlerの作成（StrictServerInterface実装）
	apiHandler := handler.NewAPIHandler(todoHandler)

	// Echoインスタンスを作成
	e := echo.New()

	// ルートを設定
	router.SetupRoutes(e, apiHandler, authHandler, sessionManager, cfg.Frontend)

	// サーバー起動
	log.Printf("Server starting on %s...", cfg.Server.Address())
	if err := e.Start(cfg.Server.Address()); err != nil {
		log.Fatal(err)
	}
}
