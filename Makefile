include .env

# =============================================================================
# ローカル開発用
# =============================================================================

# go library install
## 複数のライブラリを指定する場合は、name="xxx yyy" のように""で囲んで実行すること
go-add-library:
	docker exec -it ${BACKEND_CONTAINER_NAME} sh -c "go get ${name}"

## 依存関係の整理
go-mod-tidy:
	docker exec -i ${BACKEND_CONTAINER_NAME} sh -c "go mod tidy"

## テスト
test:
	docker exec -i ${BACKEND_CONTAINER_NAME} sh -c "go test -v ./..."

lint:
	docker exec -i ${BACKEND_CONTAINER_NAME} sh -c "staticcheck ./..."

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

# =============================================================================
# データベースマイグレーション関連
# =============================================================================

# デフォルトの環境を設定
ATLAS_ENV ?= local

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

# atlas_dev データベースを作成
create-atlas-dev-db:
	docker exec -i go_todo_db psql -U user -d postgres -c "CREATE DATABASE atlas_dev;"

# マイグレーションファイルのハッシュを再計算
migrate-hash:
	docker exec -i $(BACKEND_CONTAINER_NAME) \
		atlas migrate hash \
		--config file://atlas.hcl

# マイグレーションをロールバック（1つ前に戻す）
migrate-down:
	docker exec -i $(BACKEND_CONTAINER_NAME) \
		atlas migrate down \
		--env $(ATLAS_ENV) \
		--config file://atlas.hcl

# =============================================================================
# sqlc関連
# =============================================================================

# sqlcコード生成
sqlc-generate:
	docker exec -i $(BACKEND_CONTAINER_NAME) sqlc generate

# =============================================================================
# シード
# =============================================================================

# シードを実行
seed:
	docker exec -i $(BACKEND_CONTAINER_NAME) go run cmd/seed/main.go

# データをクリアしてシード（fresh）
seed-fresh:
	docker exec -i $(BACKEND_CONTAINER_NAME) go run cmd/seed/main.go -fresh

# =============================================================================
# OpenAPI生成（CUE → YAML → Go）
# =============================================================================

# CUEファイルのフォーマット
cue-fmt:
	docker exec -i ${BACKEND_CONTAINER_NAME} sh -c "cd /app/openapi/cue && cue fmt api.cue"

# CUEファイルの検証
cue-vet:
	docker exec -i ${BACKEND_CONTAINER_NAME} sh -c "cd /app/openapi/cue && cue vet api.cue"

# CUEからOpenAPI YAMLを生成
openapi-gen:
	docker exec -i ${BACKEND_CONTAINER_NAME} sh -c "cd /app/openapi/cue && cue export api.cue -f -o ../openapi.yaml --out yaml"

# OpenAPI YAMLからGoコードを生成
api-gen:
	docker exec -i ${BACKEND_CONTAINER_NAME} sh -c "cd /app/openapi && oapi-codegen --config oapi-codegen.yaml openapi.yaml"

# =============================================================================
# 一括コマンド
# =============================================================================

# 全コード生成（OpenAPI + sqlc）
generate: cue-fmt cue-vet openapi-gen api-gen sqlc-generate
	@echo "✅ All code generated"

# スキーマ変更時の一括処理
schema-update:
	@if [ -z "$(NAME)" ]; then \
		echo "❌ Error: NAME is required. Usage: make schema-update NAME=migration_name"; \
		exit 1; \
	fi
	$(MAKE) migrate-diff NAME=$(NAME)
	$(MAKE) sqlc-generate
	@echo "✅ Schema updated"
