package router

import (
	"net/http"

	"go-todo/internal/auth"
	"go-todo/internal/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Echoインスタンスにルートを設定
func SetupRoutes(e *echo.Echo, todoHandler *handler.TodoHandler, authHandler *handler.AuthHandler, sm *auth.SessionManager) {
	// グローバルミドルウェア
	e.Use(middleware.CORSWithConfig(CORSConfig()))
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// 共通ルート定義
	e.GET("/", handleHome)
	e.GET("/health", handleHealth)

	// Todo関連のルート
	SetupTodoRoutes(e, todoHandler, sm)

	// 認証関連のルート
	SetupAuthRoutes(e, authHandler, sm)

	// カスタムエラーハンドラー
	e.HTTPErrorHandler = customHTTPErrorHandler
}

// handleHome godoc
// @Summary API information
// @Description Get API information including name and version
// @Tags general
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router / [get]
func handleHome(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Todo API",
		"version": "0.0.0",
	})
}

// handleHealth godoc
// @Summary Health check
// @Description Check if the API is running
// @Tags general
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func handleHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}

// カスタムエラーハンドラー
func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := "Internal Server Error"

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		if m, ok := he.Message.(string); ok {
			message = m
		}
	}

	if code == http.StatusNotFound {
		c.JSON(code, map[string]string{
			"error": "Resource not found",
			"path":  c.Request().URL.Path,
		})
		return
	}

	c.JSON(code, map[string]string{"error": message})
}
