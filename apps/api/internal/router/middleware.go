package router

import (
	"log"
	"net/http"
	"os"
	"time"
)

// リクエストをロギングするミドルウェア
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("[%s] %s %s", time.Now().Format("2006-01-02 15:04:05"), r.Method, r.URL.Path)
		next(w, r)
		log.Printf("Completed in %v", time.Since(start))
	}
}

// CORSヘッダーを設定するミドルウェア
func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allowedOrigin := os.Getenv("FRONTEND_URL")
		if allowedOrigin == "" {
			allowedOrigin = "http://localhost:3000"
		}

		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = allowedOrigin
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// プリフライトリクエストの処理
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// パニックをリカバリーするミドルウェア
func RecoveryMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}

// Content-Typeをapplication/jsonに設定するミドルウェア
func JSONContentTypeMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}
