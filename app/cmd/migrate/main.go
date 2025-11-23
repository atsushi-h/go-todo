package main

import (
    "flag"
    "log"
    
    "go-todo/database"
)

func main() {
    // コマンドラインフラグの定義
    action := flag.String("action", "migrate", "Action to perform: migrate, reset")
    flag.Parse()
    
    // DBの初期化
    db, err := database.Init()
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer database.Close(db)
    
    // ヘルスチェック
    if err := database.HealthCheck(db); err != nil {
        log.Fatal("Database health check failed:", err)
    }
    
    // アクション実行
    switch *action {
    case "migrate":
        if err := database.Migrate(db); err != nil {
            log.Fatal("Failed to migrate database:", err)
        }
        log.Println("Migration completed successfully.")
        
    case "reset":
        if err := database.ResetDatabase(db); err != nil {
            log.Fatal("Failed to reset database:", err)
        }
        log.Println("Database reset completed successfully.")
        
    default:
        log.Fatalf("Unknown action: %s", *action)
    }
}
