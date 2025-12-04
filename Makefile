include .env

empty:
	echo "empty"

# 開発環境のdocker compose コマンド
dcb-dev:
	docker compose build
dcu-dev:
	docker compose up -d
dcd-dev:
	docker compose down

# コンテナ環境へsshログイン
backend-ssh:
	docker exec -it ${BACKEND_CONTAINER_NAME} sh

# データベースマイグレーション関連
# デフォルトの環境を設定
ATLAS_ENV ?= local

# スキーマSQLを生成（Goモデルから）
generate-schema:
	docker exec -i $(BACKEND_CONTAINER_NAME) go run cmd/schema/main.go

# マイグレーションの差分ファイルを生成
# 使用例: make migrate-diff NAME=create_todos
migrate-diff:
	@if [ -z "$(NAME)" ]; then \
		echo "❌ Error: NAME is required. Usage: make migrate-diff NAME=migration_name"; \
		exit 1; \
	fi
	docker exec -i $(BACKEND_CONTAINER_NAME) \
		atlas migrate diff $(NAME) \
		--env $(ATLAS_ENV) \
		--config file://atlas.hcl

# マイグレーションを適用
migrate-apply:
	docker exec -i $(BACKEND_CONTAINER_NAME) \
		atlas migrate apply \
		--env $(ATLAS_ENV) \
		--config file://atlas.hcl

# マイグレーション状態を確認
migrate-status:
	docker exec -i $(BACKEND_CONTAINER_NAME) \
		atlas migrate status \
		--env $(ATLAS_ENV) \
		--config file://atlas.hcl

# マイグレーション履歴を検証
migrate-validate:
	docker exec -i $(BACKEND_CONTAINER_NAME) \
		atlas migrate validate \
		--env $(ATLAS_ENV) \
		--config file://atlas.hcl

# 初回セットアップ（スキーマ生成 + 初回マイグレーション作成）
migrate-init:
	@echo "📋 Generating schema-gen.sql..."
	@$(MAKE) schema-generate
	@echo "📝 Creating initial migration..."
	@$(MAKE) migrate-diff NAME=init
	@echo "✅ Migration initialized. Run 'make migrate-apply' to apply."

# atlas_dev データベースを作成
create-atlas-dev-db:
	docker exec -i go_todo_db psql -U user -d postgres -c "CREATE DATABASE atlas_dev;"

# マイグレーションをロールバック（1つ前に戻す）
migrate-down:
	docker exec -i $(BACKEND_CONTAINER_NAME) \
		atlas migrate down \
		--env $(ATLAS_ENV) \
		--config file://atlas.hcl

# シードを実行
seed:
	docker exec -i $(BACKEND_CONTAINER_NAME) go run cmd/seed/main.go

# データをクリアしてシード（fresh）
seed-fresh:
	docker exec -i $(BACKEND_CONTAINER_NAME) go run cmd/seed/main.go -fresh

# ローカル開発用
# go library install
## 複数のライブラリを指定する場合は、name="xxx yyy" のように""で囲んで実行すること
go-add-library:
	docker exec -it ${BACKEND_CONTAINER_NAME} sh -c "go get ${name}"
## テスト
test:
	docker exec -i ${BACKEND_CONTAINER_NAME} sh -c "go test -v ./..."
lint:
	docker exec -i ${BACKEND_CONTAINER_NAME} sh -c "staticcheck ./..."

## OpenAPI YAML生成
openapi:
	docker exec -i ${BACKEND_CONTAINER_NAME} sh -c "cd /app && swag init -g cmd/server/main.go -o openapi --outputTypes yaml"
	mv /Users/atsushi-h/workspace/go-todo/apps/api/openapi/swagger.yaml /Users/atsushi-h/workspace/go-todo/apps/api/openapi/openapi.yaml