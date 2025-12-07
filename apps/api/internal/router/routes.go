package router

import (
	"net/http"

	"go-todo/internal/auth"
	"go-todo/internal/gen"
	"go-todo/internal/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Echoインスタンスにルートを設定
func SetupRoutes(e *echo.Echo, apiHandler *handler.APIHandler, authHandler *handler.AuthHandler, sm *auth.SessionManager) {
	// グローバルミドルウェア
	e.Use(middleware.CORSWithConfig(CORSConfig()))
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// 認証ミドルウェアをstrictmiddlewareとしてラップ
	authMiddleware := createAuthMiddleware(sm)

	// StrictハンドラーをEchoハンドラーにラップ（認証ミドルウェア付き）
	strictHandler := gen.NewStrictHandler(apiHandler, []gen.StrictMiddlewareFunc{authMiddleware})

	// 生成されたルート登録関数を使用
	gen.RegisterHandlers(e, strictHandler)

	// 認証関連のルート（手動で設定）
	SetupAuthRoutes(e, authHandler, sm)

	// カスタムエラーハンドラー
	e.HTTPErrorHandler = customHTTPErrorHandler
}

// 認証ミドルウェアをStrictMiddlewareFuncに変換
func createAuthMiddleware(sm *auth.SessionManager) gen.StrictMiddlewareFunc {
	return func(f gen.StrictHandlerFunc, operationID string) gen.StrictHandlerFunc {
		return func(ctx echo.Context, request interface{}) (interface{}, error) {
			// 認証不要なエンドポイントをスキップ
			if operationID == "GetInfo" || operationID == "GetHealth" {
				return f(ctx, request)
			}

			// セッションからユーザーIDを取得
			userID, err := sm.GetUserID(ctx.Request())
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}

			// コンテキストにユーザーIDを設定
			newCtx := auth.WithUserID(ctx.Request().Context(), userID)
			ctx.SetRequest(ctx.Request().WithContext(newCtx))

			return f(ctx, request)
		}
	}
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
