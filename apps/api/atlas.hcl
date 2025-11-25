env "local" {
  # Bunモデルから生成したスキーマファイル
  src = "file://schema-gen.sql"
  
  # Dockerネットワーク経由で接続するデータベースURL
  url = "postgres://user:pass@go_todo_db:5432/go_todo_db?sslmode=disable"
  
  # Atlas開発用データベース（差分計算に使用）
  # ローカルに別のDBを作成して使用
  dev = "postgres://user:pass@go_todo_db:5432/atlas_dev?sslmode=disable"
  
  # マイグレーションファイルの格納先
  migration {
    dir = "file://migrations"
  }
  
  # フォーマット設定
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "dev" {
  src = "file://schema-gen.sql"
  url = getenv("DATABASE_URL")
  dev = getenv("ATLAS_DEV_URL")
  
  migration {
    dir = "file://migrations"
  }
}
