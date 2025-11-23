package router

import (
	"encoding/json"
	"net/http"

	"go-todo/internal/handler"
)

// ルーターにルートを設定
func SetupRoutes(r *Router, todoHandler *handler.TodoHandler) {
	// ミドルウェアを設定
	r.Use(RecoveryMiddleware)
	r.Use(LoggingMiddleware)
	r.Use(CORSMiddleware)

	// 共通ルート定義
	r.GET("/", handleHome)
	r.GET("/health", handleHealth)

	// Todo関連のルート（todo.route.goに分離）
	SetupTodoRoutes(r, todoHandler)

	// 404カスタムハンドラー
	r.SetNotFound(handleNotFound)
}

// ホームページハンドラー
func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Todo API",
		"version": "0.0.0",
	})
}

// ヘルスチェックハンドラー
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
