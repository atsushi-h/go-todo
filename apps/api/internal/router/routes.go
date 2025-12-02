package router

import (
	"encoding/json"
	"net/http"

	"go-todo/internal/handler"
)

// ルーターにルートを設定
func SetupRoutes(r *Router, todoHandler *handler.TodoHandler) {
	// グローバルミドルウェアを設定（ルートマッチング前に適用）
	r.UseGlobal(CORSMiddleware)

	// ミドルウェアを設定（ルートマッチング後に適用）
	r.Use(RecoveryMiddleware)
	r.Use(LoggingMiddleware)

	// 共通ルート定義
	r.GET("/", handleHome)
	r.GET("/health", handleHealth)

	// Todo関連のルート（todo.route.goに分離）
	SetupTodoRoutes(r, todoHandler)

	// 404カスタムハンドラー
	r.SetNotFound(handleNotFound)
}

// handleHome godoc
// @Summary API information
// @Description Get API information including name and version
// @Tags general
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router / [get]
func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
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
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

// 404エラーハンドラー
func handleNotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "Resource not found",
		"path":  r.URL.Path,
	})
}
