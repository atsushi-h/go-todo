package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"go-todo/internal/auth"
	"go-todo/internal/service"

	"github.com/markbates/goth/gothic"
)

type AuthHandler struct {
	userService    *service.UserService
	sessionManager *auth.SessionManager
}

func NewAuthHandler(userService *service.UserService, sm *auth.SessionManager) *AuthHandler {
	return &AuthHandler{
		userService:    userService,
		sessionManager: sm,
	}
}

func (h *AuthHandler) BeginAuth(w http.ResponseWriter, r *http.Request) {
	gothic.BeginAuthHandler(w, r)
}

func (h *AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	gothUser, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		log.Printf("Authentication failed: %v", err)
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	user, err := h.userService.FindOrCreateFromOAuth(r.Context(), gothUser)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// SessionManagerを使用
	if err := h.sessionManager.SetUserID(w, r, user.ID); err != nil {
		log.Printf("Failed to save session: %v", err)
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	// リダイレクト先のフロントエンドURLを環境変数から取得
	frontendURL := os.Getenv("FRONTEND_URL")
    if frontendURL == "" {
		log.Printf("Server configuration error: FRONTEND_URL is not set")
        http.Error(w, "Server configuration error", http.StatusInternalServerError)
        return
    }

	http.Redirect(w, r, frontendURL, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.sessionManager.Clear(w, r); err != nil {
		log.Printf("Failed to clear session: %v", err)
		http.Error(w, "Failed to clear session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out"})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.userService.GetByID(r.Context(), userID)
	if err != nil {
		log.Printf("User not found (id=%d): %v", userID, err) 
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
