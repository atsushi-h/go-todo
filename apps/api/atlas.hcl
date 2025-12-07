env "local" {
  # スキーマファイル
  src = "file://db/schema.sql"

  # Dockerネットワーク経由で接続するデータベースURL
  url = "postgres://user:pass@go_todo_db:5432/go_todo_db?sslmode=disable"

  # Atlas開発用データベース（差分計算に使用）
  dev = "postgres://user:pass@go_todo_db:5432/atlas_dev?sslmode=disable"

  # マイグレーションファイルの格納先
  migration {
    dir = "file://db/migrations"
  }

  # フォーマット設定
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "dev" {
  src = "file://db/schema.sql"
  url = getenv("DATABASE_URL")
  dev = getenv("ATLAS_DEV_URL")

  migration {
    dir = "file://db/migrations"
  }
}
