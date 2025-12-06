package auth

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type contextKey string

const UserIDKey contextKey = "user_id"

// 認証が必要なルートに適用するEchoミドルウェア
func RequireAuth(sm *SessionManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, err := sm.GetUserID(c.Request())
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}

			// 標準contextにUserIDを設定
			ctx := context.WithValue(c.Request().Context(), UserIDKey, userID)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

// コンテキストからユーザーIDを取得
func GetUserIDFromContext(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(UserIDKey).(uint)
	return userID, ok
}
