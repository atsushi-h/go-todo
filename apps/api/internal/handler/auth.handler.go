package handler

import (
	"log"
	"net/http"

	"go-todo/internal/auth"
	"go-todo/internal/config"
	"go-todo/internal/mapper"
	"go-todo/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

type AuthHandler struct {
	userService    *service.UserService
	sessionManager *auth.SessionManager
	frontendURL    string
}

func NewAuthHandler(userService *service.UserService, sm *auth.SessionManager, frontendConfig config.FrontendConfig) *AuthHandler {
	return &AuthHandler{
		userService:    userService,
		sessionManager: sm,
		frontendURL:    frontendConfig.URL,
	}
}

func (h *AuthHandler) BeginAuth(c echo.Context) error {
	// providerをリクエストに設定
	auth.SetProviderToRequest(c)

	gothic.BeginAuthHandler(c.Response(), c.Request())
	return nil
}

func (h *AuthHandler) Callback(c echo.Context) error {
	// providerをリクエストに設定
	auth.SetProviderToRequest(c)

	gothUser, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		log.Printf("Authentication failed: %v", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "Authentication failed")
	}

	user, err := h.userService.FindOrCreateFromOAuth(c.Request().Context(), gothUser)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
	}

	// SessionManagerを使用
	if err := h.sessionManager.SetUserID(c.Response(), c.Request(), user.ID); err != nil {
		log.Printf("Failed to save session: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save session")
	}

	// フロントエンドURLにリダイレクト
	return c.Redirect(http.StatusTemporaryRedirect, h.frontendURL)
}

func (h *AuthHandler) Logout(c echo.Context) error {
	if err := h.sessionManager.Clear(c.Response(), c.Request()); err != nil {
		log.Printf("Failed to clear session: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to clear session")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Logged out"})
}

func (h *AuthHandler) Me(c echo.Context) error {
	userID, ok := auth.GetUserIDFromContext(c.Request().Context())
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	user, err := h.userService.GetByID(c.Request().Context(), userID)
	if err != nil {
		log.Printf("User not found (id=%d): %v", userID, err)
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	return c.JSON(http.StatusOK, mapper.UserToResponse(user))
}

func (h *AuthHandler) DeleteUserAccount(c echo.Context) error {
	userID, ok := auth.GetUserIDFromContext(c.Request().Context())
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	if err := h.userService.DeleteAccount(c.Request().Context(), userID); err != nil {
		if err == service.ErrUserNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		log.Printf("Failed to delete user account (id=%d): %v", userID, err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete account")
	}

	// セッションをクリア
	if err := h.sessionManager.Clear(c.Response(), c.Request()); err != nil {
		log.Printf("Warning: Failed to clear session after account deletion (id=%d): %v", userID, err)
	}

	return c.NoContent(http.StatusNoContent)
}
