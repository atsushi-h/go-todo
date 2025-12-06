package router

import (
	"go-todo/internal/auth"
	"go-todo/internal/handler"

	"github.com/labstack/echo/v4"
)

// 認証関連のルートを設定
func SetupAuthRoutes(e *echo.Echo, authHandler *handler.AuthHandler, sm *auth.SessionManager) {
	// 認証関連のルート（認証不要）
	authGroup := e.Group("/auth")
	authGroup.GET("/:provider", authHandler.BeginAuth)
	authGroup.GET("/:provider/callback", authHandler.Callback)

	// ログアウト
	e.POST("/logout", authHandler.Logout)

	// ユーザー情報取得（認証必要）
	e.GET("/me", authHandler.Me, auth.RequireAuth(sm))
}
