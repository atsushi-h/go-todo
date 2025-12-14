package router

import (
	"go-todo/internal/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// EchoのCORSミドルウェア設定
func CORSConfig(frontendConfig config.FrontendConfig) middleware.CORSConfig {
	return middleware.CORSConfig{
		AllowOrigins:     []string{frontendConfig.URL},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders:     []string{echo.HeaderContentType, echo.HeaderAuthorization},
		AllowCredentials: true,
	}
}
