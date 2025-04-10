package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"github.com/atsushi-h/go-todo/pkg/db"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("❌ .env の読み込みに失敗しました")
	}

	_, err := db.Connect()
	if err != nil {
		log.Fatalf("❌ DB接続失敗: %v", err)
	}

	r := chi.NewRouter()

	// /api/health エンドポイント
	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	fmt.Println("🚀 サーバー起動中 :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
