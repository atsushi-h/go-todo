# 開発ガイド

このドキュメントでは、Go APIの日常的な開発タスクと実践的な手順を説明します。

## 目次

- [環境セットアップ](#環境セットアップ)
- [日常的な開発フロー](#日常的な開発フロー)
- [新機能の追加](#新機能の追加)
- [テスト](#テスト)
- [コード品質](#コード品質)
- [トラブルシューティング](#トラブルシューティング)
- [参考情報](#参考情報)

## 環境セットアップ

### 1. Docker環境の起動

```bash
make dcu-dev
```

以下のコンテナが起動します：
- `go_todo_server` - Go APIサーバー（Air でホットリロード）
- `go_todo_db` - PostgreSQL データベース
- `go_todo_redis` - Redis（セッションストア）

### 2. 環境変数の設定

`.env` ファイルが存在することを確認してください：

```bash
# Database
POSTGRES_HOST=go_todo_db
POSTGRES_PORT=5432
POSTGRES_DB=go_todo_db
POSTGRES_USER=user
POSTGRES_PASSWORD=password

# Redis
REDIS_HOST=go_todo_redis
REDIS_PORT=6379

# OAuth
GOOGLE_CLIENT_ID=your-client-id
GOOGLE_CLIENT_SECRET=your-client-secret
OAUTH_CALLBACK_URL=http://localhost:8080/auth/google/callback

# Server
SERVER_PORT=8080

# Frontend
FRONTEND_URL=http://localhost:3000

# Cookie
COOKIE_SECURE=false
```

### 3. atlas_dev データベースの作成

Atlasの差分計算に必要な作業用データベースを作成します（詳細は「[7.7 atlas_dev データベースについて](#77-atlas_dev-データベースについて)」を参照）：

```bash
make create-atlas-dev-db
```

または手動で：

```bash
docker exec -it go_todo_db psql -U user -d go_todo_db -c "CREATE DATABASE atlas_dev;"
```

### 4. 初期マイグレーションの適用

```bash
make migrate-apply
```

### 5. シードデータの投入

```bash
make seed
```

または、既存データを削除してから投入する場合：

```bash
make seed-fresh
```

## 日常的な開発フロー

### 開発サーバーの起動/停止

**起動**:

```bash
make dcu-dev
```

**停止**:

```bash
make dcd-dev
```

**再ビルド（Dockerイメージ）**:

```bash
make dcb-dev
```

### ホットリロード（Air）

`go_todo_server` コンテナでは **Air** が動作しており、Go ファイルの変更を検知して自動的にサーバーを再起動します。

編集するだけで変更が即座に反映されます。

### ログの確認

**サーバーログ**:

```bash
docker logs -f go_todo_server
```

**データベースログ**:

```bash
docker logs -f go_todo_db
```

### コンテナへのアクセス

**サーバーコンテナ**:

```bash
make backend-ssh
```

または：

```bash
docker exec -it go_todo_server sh
```

**データベースコンテナ（psql）**:

```bash
docker exec -it go_todo_db psql -U user -d go_todo_db
```

### 依存関係の追加

**Goパッケージを追加**:

```bash
# コンテナ内で
go get github.com/some/package

# または Makefile 経由
make go-add-library PKG=github.com/some/package
```

**依存関係の整理**:

```bash
make go-mod-tidy
```

## 新機能の追加

このセクションでは、新しいAPIエンドポイントの追加、データベーススキーマの変更など、典型的な開発タスクの実践的な手順を説明します。

### 3.1 新しいAPIエンドポイントの追加

完全なワークフローを、具体例を用いて説明します。

**例: Todoの優先度別一覧取得エンドポイントを追加する**

---

#### Step 1: CUE定義を編集

[openapi/cue/api.cue](openapi/cue/api.cue) に新しいエンドポイントを追加します。

```cue
paths: {
    // 既存のエンドポイント...

    "/todos/by-priority": {
        get: {
            operationId: "ListTodosByPriority"
            summary: "List todos ordered by priority"
            parameters: [
                {
                    name: "priority"
                    in: "query"
                    required: true
                    schema: {
                        type: "string"
                        enum: ["high", "medium", "low"]
                    }
                }
            ]
            responses: {
                "200": {
                    description: "Success"
                    content: "application/json": schema: {
                        type: "array"
                        items: {$ref: "#/components/schemas/Todo"}
                    }
                }
                "400": {
                    description: "Bad Request"
                    content: "application/json": schema: {$ref: "#/components/schemas/ErrorResponse"}
                }
                "401": {
                    description: "Unauthorized"
                    content: "application/json": schema: {$ref: "#/components/schemas/ErrorResponse"}
                }
                "500": {
                    description: "Internal Server Error"
                    content: "application/json": schema: {$ref: "#/components/schemas/ErrorResponse"}
                }
            }
        }
    }
}
```

**重要**:
- `operationId` は **PascalCase** で記述（Goのメソッド名になる）
- レスポンスコードごとに型が生成される

---

#### Step 2: コード生成

```bash
make generate
```

これにより以下が生成されます：
- `openapi/openapi.yaml` - OpenAPI YAML
- `internal/gen/api.gen.go` - Go コード（StrictServerInterface に `ListTodosByPriority` メソッドが追加される）

---

#### Step 3: ハンドラーの実装

[internal/handler/todo.handler.go](internal/handler/todo.handler.go) に新しいメソッドを追加します。

```go
// 必要なインポートを追加
import (
    "context"
    "go-todo/internal/auth"
    "go-todo/internal/gen"
    "go-todo/internal/mapper"
    "go-todo/internal/service"
)

// ListTodosByPriority - 優先度別にTodoを取得
func (h *TodoHandler) ListTodosByPriority(ctx context.Context, request gen.ListTodosByPriorityRequestObject) (gen.ListTodosByPriorityResponseObject, error) {
    userID, ok := auth.GetUserIDFromContext(ctx)
    if !ok {
        return gen.ListTodosByPriority401JSONResponse{Message: "Unauthorized"}, nil
    }

    // クエリパラメータの検証
    if request.Params.Priority == nil {
        return gen.ListTodosByPriority400JSONResponse{Message: "Priority is required"}, nil
    }

    todos, err := h.service.GetTodosByPriority(ctx, userID, *request.Params.Priority)
    if err != nil {
        return gen.ListTodosByPriority500JSONResponse{Message: "Internal server error"}, nil
    }

    return gen.ListTodosByPriority200JSONResponse(mapper.TodosToResponse(todos)), nil
}
```

---

#### Step 4: APIHandlerに追加

[internal/handler/api.handler.go](internal/handler/api.handler.go) でTodoHandlerに委譲します。

```go
func (h *APIHandler) ListTodosByPriority(ctx context.Context, request gen.ListTodosByPriorityRequestObject) (gen.ListTodosByPriorityResponseObject, error) {
    return h.todoHandler.ListTodosByPriority(ctx, request)
}
```

---

#### Step 5: サービス層の実装

[internal/service/todo.service.go](internal/service/todo.service.go) にビジネスロジックを追加します。

```go
func (s *TodoService) GetTodosByPriority(ctx context.Context, userID int64, priority string) ([]sqlc.Todo, error) {
    return s.repo.ListTodosByPriority(ctx, sqlc.ListTodosByPriorityParams{
        UserID:   userID,
        Priority: priority,
    })
}
```

---

#### Step 6: SQLクエリの追加（必要な場合）

優先度カラムが存在する場合、[db/query/todo.sql](db/query/todo.sql) に新しいクエリを追加します。

```sql
-- name: ListTodosByPriority :many
SELECT * FROM todos
WHERE user_id = $1 AND priority = $2 AND deleted_at IS NULL
ORDER BY created_at DESC;
```

sqlcコードを再生成：

```bash
make sqlc-generate
```

---

#### Step 7: Repositoryインターフェースの更新

[internal/service/todo.repository.go](internal/service/todo.repository.go) に新しいメソッドを追加します。

```go
type TodoRepository interface {
    // 既存のメソッド...
    ListTodosByPriority(ctx context.Context, arg sqlc.ListTodosByPriorityParams) ([]sqlc.Todo, error)
}
```

---

#### Step 8: 認証設定（認証不要な場合）

もしこのエンドポイントが認証不要な場合、[internal/router/routes.go](internal/router/routes.go) の `createAuthMiddleware` 関数で除外設定を追加します。

```go
func createAuthMiddleware(sm *auth.SessionManager) gen.StrictMiddlewareFunc {
    return func(f gen.StrictHandlerFunc, operationID string) gen.StrictHandlerFunc {
        return func(ctx echo.Context, request interface{}) (interface{}, error) {
            // 認証不要なエンドポイントをスキップ
            if operationID == "GetInfo" || operationID == "GetHealth" {
                return f(ctx, request)
            }

            // 認証処理...
        }
    }
}
```

今回の例では認証が必要なので、この手順はスキップします。

---

#### Step 9: テストの作成

[internal/service/todo.service_test.go](internal/service/todo.service_test.go) にテストを追加します。

```go
func TestTodoService_GetTodosByPriority(t *testing.T) {
    mockRepo := mocks.NewTodoRepository(t)
    service := NewTodoService(mockRepo)

    mockRepo.EXPECT().
        ListTodosByPriority(mock.Anything, sqlc.ListTodosByPriorityParams{
            UserID:   123,
            Priority: "high",
        }).
        Return([]sqlc.Todo{
            {ID: 1, Title: "Urgent task", Priority: "high"},
        }, nil)

    todos, err := service.GetTodosByPriority(context.Background(), 123, "high")

    assert.NoError(t, err)
    assert.Len(t, todos, 1)
    assert.Equal(t, "high", todos[0].Priority)
}
```

---

### 3.2 データベーススキーマの変更

完全なワークフローを、具体例を用いて説明します。

**例: Todoテーブルに `priority` カラムを追加する**

---

#### Step 1: スキーマ定義を編集

[db/schema.sql](db/schema.sql) を編集します。

```sql
CREATE TABLE todos (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    priority INTEGER NOT NULL DEFAULT 0,  -- 新規追加
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_todos_user_id ON todos(user_id);
CREATE INDEX idx_todos_deleted_at ON todos(deleted_at);
CREATE INDEX idx_todos_priority ON todos(priority);  -- 新規追加
```

---

#### Step 2: マイグレーションファイルを生成

```bash
make migrate-diff NAME=add_priority_to_todos
```

`apps/api/db/migrations/YYYYMMDDHHMMSS_add_priority_to_todos.sql` が生成されます。

---

#### Step 3: マイグレーションファイルをレビュー

生成されたSQLファイルを確認します：

```sql
-- Add column "priority" to table: "todos"
ALTER TABLE "todos" ADD COLUMN "priority" integer NOT NULL DEFAULT 0;
-- Create index "idx_todos_priority" to table: "todos"
CREATE INDEX "idx_todos_priority" ON "todos" ("priority");
```

必要に応じて手動で編集できます。

---

#### Step 4: マイグレーション状態を確認

```bash
make migrate-status
```

出力例：

```
Migration Status: PENDING
  Current Version: 20250125100000
  Next Version:    20250125120000_add_priority_to_todos
  Total Pending:   1
```

---

#### Step 5: マイグレーションを適用

```bash
make migrate-apply
```

成功すると：

```
Migrating to version 20250125120000 (1 migration)
  -> 20250125120000_add_priority_to_todos.sql .... ok (15ms)
  -------------------------------
  -> 0s
```

---

#### Step 6: SQLクエリの追加

[db/query/todo.sql](db/query/todo.sql) に新しいカラムを使うクエリを追加します。

```sql
-- name: ListTodosByPriority :many
SELECT * FROM todos
WHERE user_id = $1 AND deleted_at IS NULL
ORDER BY priority DESC, created_at DESC;

-- name: UpdateTodoPriority :one
UPDATE todos
SET priority = $2, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;
```

---

#### Step 7: sqlcコードを再生成

```bash
make sqlc-generate
```

これにより `db/sqlc/` 配下のコードが更新されます。

---

#### Step 8: サービス/ハンドラーの更新

新しく生成された型とメソッドを使用してビジネスロジックを実装します。

[internal/service/todo.service.go](internal/service/todo.service.go):

```go
func (s *TodoService) UpdateTodoPriority(ctx context.Context, id int64, priority int32) (*sqlc.Todo, error) {
    todo, err := s.repo.UpdateTodoPriority(ctx, sqlc.UpdateTodoPriorityParams{
        ID:       id,
        Priority: priority,
    })
    if err != nil {
        return nil, err
    }
    return &todo, nil
}
```

---

#### Step 9: Gitへコミット

マイグレーションファイルと `atlas.sum` を必ずコミットします：

```bash
git add apps/api/db/
git commit -m "feat: Add priority column to todos table"
```

**重要**: `atlas.sum` はマイグレーションファイルの整合性を保証するため、必ずコミットしてください。

---

### 一括処理コマンド

マイグレーション生成とsqlc生成を一括で実行：

```bash
make schema-update NAME=add_priority_to_todos
```

これは以下を順番に実行します：
1. `make migrate-diff`
2. `make sqlc-generate`

---

### 3.3 認証が不要なエンドポイントの追加

デフォルトでは、すべてのAPIエンドポイントに認証が必要です。

認証を不要にするには、[internal/router/routes.go](internal/router/routes.go) の `createAuthMiddleware` 関数で `operationID` を除外リストに追加します。

```go
func createAuthMiddleware(sm *auth.SessionManager) gen.StrictMiddlewareFunc {
    return func(f gen.StrictHandlerFunc, operationID string) gen.StrictHandlerFunc {
        return func(ctx echo.Context, request interface{}) (interface{}, error) {
            // 認証不要なエンドポイントをスキップ
            if operationID == "GetInfo" || operationID == "GetHealth" || operationID == "YourNewEndpoint" {
                return f(ctx, request)
            }

            // 認証処理...
        }
    }
}
```

**例**:
- `GetInfo` - `/api/info` (API情報取得)
- `GetHealth` - `/api/health` (ヘルスチェック)

## テスト

### テストの実行

**全テスト実行**:

```bash
make test
```

**統合テストをスキップ**:

```bash
go test -short ./...
```

**特定のパッケージのみ**:

```bash
go test ./internal/service
```

**カバレッジ付き**:

```bash
go test -cover ./...
```

### モック生成

Repositoryインターフェースを変更した場合、モックを再生成します：

```bash
make mock
```

これにより `internal/service/mocks/` 配下のモックが更新されます。

### テストの書き方パターン

#### サービス層の単体テスト（mockery使用）

[internal/service/todo.service_test.go](internal/service/todo.service_test.go):

```go
func TestTodoService_GetTodoByID(t *testing.T) {
    // モックRepositoryを作成
    mockRepo := mocks.NewTodoRepository(t)
    service := NewTodoService(mockRepo)

    // モックの期待値設定
    mockRepo.EXPECT().
        GetTodoByID(mock.Anything, sqlc.GetTodoByIDParams{
            ID:     1,
            UserID: 123,
        }).
        Return(sqlc.Todo{
            ID:     1,
            Title:  "Test Todo",
            UserID: 123,
        }, nil)

    // テスト実行
    todo, err := service.GetTodoByID(context.Background(), 1, 123)

    // アサーション
    assert.NoError(t, err)
    assert.Equal(t, int64(1), todo.ID)
    assert.Equal(t, "Test Todo", todo.Title)
}
```

#### 統合テスト（実際のDB使用）

[internal/service/user.service_integration_test.go](internal/service/user.service_integration_test.go):

```go
func TestUserService_CreateOrUpdateUser_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    // テスト用DB接続
    cfg, _ := config.Load()
    pool, _ := database.NewPool(context.Background(), cfg.Database)
    defer pool.Close()

    queries := sqlc.New(pool)
    service := service.NewUserService(queries, pool)

    // 実際のDB操作をテスト
    user, err := service.CreateOrUpdateUser(context.Background(), service.UserParams{
        Email:    "test@example.com",
        Name:     "Test User",
        Provider: "google",
    })

    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
    assert.Equal(t, "test@example.com", user.Email)

    // クリーンアップ
    // ...
}
```

#### テーブル駆動テスト

```go
func TestTodoService_CreateTodo(t *testing.T) {
    tests := []struct {
        name        string
        title       string
        description *string
        setupMock   func(*mocks.TodoRepository)
        wantErr     bool
    }{
        {
            name:        "valid todo",
            title:       "Test",
            description: ptrString("Description"),
            setupMock: func(m *mocks.TodoRepository) {
                m.EXPECT().CreateTodo(mock.Anything, mock.Anything).
                    Return(sqlc.Todo{ID: 1}, nil)
            },
            wantErr: false,
        },
        {
            name:        "repository error",
            title:       "Test",
            description: nil,
            setupMock: func(m *mocks.TodoRepository) {
                m.EXPECT().CreateTodo(mock.Anything, mock.Anything).
                    Return(sqlc.Todo{}, errors.New("db error"))
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := mocks.NewTodoRepository(t)
            tt.setupMock(mockRepo)
            service := NewTodoService(mockRepo)

            _, err := service.CreateTodo(context.Background(), 123, tt.title, tt.description)

            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}

// ヘルパー関数
func ptrString(s string) *string {
    return &s
}
```

## コード品質

### Linterの実行

```bash
make lint
```

内部で `staticcheck` を実行します。

### コードフォーマット

**CUEファイル**:

```bash
make cue-fmt
```

**Goファイル**:

Airが自動的に `gofmt` を実行するため、手動実行は不要です。

手動で実行する場合：

```bash
go fmt ./...
```

### Pre-pushフック（lefthook）

このプロジェクトでは **lefthook** を使用して、`git push` 前に自動的にテストとLintを実行します。

**動作**:

```bash
git push
```

以下が自動実行されます：
1. `make lint` - staticcheck
2. `make test` - 全テスト

**フックをスキップ（非推奨）**:

```bash
git push --no-verify
```

**lefthook設定**: [lefthook.yml](../../lefthook.yml)

```yaml
pre-push:
  commands:
    lint:
      run: make lint
    test:
      run: make test
```

## トラブルシューティング

### 6.1 Docker関連

#### コンテナが起動しない

**確認**:

```bash
docker ps -a | grep go_todo
```

**ログを確認**:

```bash
docker logs go_todo_server
docker logs go_todo_db
```

**再起動**:

```bash
make dcd-dev
make dcu-dev
```

#### ポートが使用中

**エラー例**:

```
Error starting userland proxy: listen tcp4 0.0.0.0:8080: bind: address already in use
```

**解決方法**:

使用中のプロセスを確認：

```bash
lsof -i :8080
```

プロセスを終了するか、`.env` で別のポートを使用：

```env
SERVER_PORT=8081
```

#### ボリュームの問題

データベースが初期化されない場合、ボリュームを削除して再作成：

```bash
docker-compose -f docker/docker-compose.dev.yml down -v
make dcu-dev
```

### 6.2 マイグレーション関連

#### エラー: "atlas_dev database does not exist"

**原因**: Atlas作業用の一時データベースが作成されていない

**解決方法**:

```bash
make create-atlas-dev-db
```

または手動で:

```bash
docker exec -it go_todo_db psql -U user -d go_todo_db -c "CREATE DATABASE atlas_dev;"
```

#### エラー: "migration checksum mismatch"

**原因**: マイグレーションファイルを直接編集したが、`atlas.sum` が更新されていない

**解決方法**:

```bash
make migrate-hash
```

これにより `atlas.sum` が再計算されます。

#### エラー: "connection refused"

**原因**:
- DBコンテナが起動していない
- Atlas CLIの実行場所とDB接続先のホスト名が一致していない

**解決方法**:

```bash
# DBコンテナの状態確認
docker ps | grep go_todo_db

# コンテナが起動していない場合
make dcu-dev
```

**接続先の確認**:
- Dockerコンテナ内から実行: `go_todo_db` を使用
- ホストマシンから実行: `localhost` を使用

#### マイグレーションが中途半端に失敗した場合

```bash
# 現在の状態を確認
make migrate-status

# 手動でSQLを実行してデータを修正
docker exec -it go_todo_db psql -U user -d go_todo_db

# atlas_migrationsテーブルを確認
SELECT * FROM atlas_schema_revisions;

# 必要に応じてバージョンをロールバック（慎重に！）
DELETE FROM atlas_schema_revisions WHERE version = '20250125120000';
```

### 6.3 コード生成関連

#### CUE検証エラー

**エラー例**:

```
cue: field not allowed: operationid
```

**原因**: CUE構文エラー（例: `operationid` → `operationId`）

**解決方法**:

CUE定義を修正してから再実行：

```bash
make cue-vet  # 検証のみ
make openapi-gen  # 生成
```

#### oapi-codegen エラー

**エラー例**:

```
error generating code: invalid operation ID
```

**原因**: OpenAPI仕様が不正（operationIdの重複など）

**解決方法**:

`openapi/openapi.yaml` を確認し、CUE定義を修正：

```bash
make openapi-gen
```

#### sqlc生成エラー

**エラー例**:

```
ERROR: column "priority" does not exist
```

**原因**: `db/query/*.sql` で存在しないカラムを参照している

**解決方法**:

1. `db/schema.sql` にカラムが定義されているか確認
2. マイグレーションが適用されているか確認：
   ```bash
   make migrate-status
   make migrate-apply
   ```
3. sqlc再生成：
   ```bash
   make sqlc-generate
   ```

### 6.4 テスト関連

#### モックが古い

**エラー例**:

```
undefined: ListTodosByPriority
```

**原因**: Repositoryインターフェースを変更したが、モックが再生成されていない

**解決方法**:

```bash
make mock
```

#### DB接続エラー（統合テスト）

**エラー例**:

```
failed to connect to database
```

**原因**: テスト用DBが起動していない

**解決方法**:

統合テストをスキップ：

```bash
go test -short ./...
```

またはDB環境を起動：

```bash
make dcu-dev
```

### 6.5 認証関連

#### セッションが取得できない

**症状**: APIリクエストで常に `401 Unauthorized`

**確認事項**:

1. **Redisが起動しているか**:
   ```bash
   docker ps | grep go_todo_redis
   ```

2. **Cookieが送信されているか**:
   ブラウザのDevToolsでCookieを確認

3. **CORS設定が正しいか**:
   [internal/router/cors.go](internal/router/cors.go) の `AllowOrigins` を確認

**解決方法**:

Redis再起動：

```bash
docker restart go_todo_redis
```

#### OAuth設定エラー

**エラー例**:

```
OAuth provider not found
```

**原因**: 環境変数が設定されていない

**解決方法**:

`.env` ファイルを確認：

```env
GOOGLE_CLIENT_ID=your-client-id
GOOGLE_CLIENT_SECRET=your-client-secret
OAUTH_CALLBACK_URL=http://localhost:8080/auth/google/callback
```

サーバー再起動：

```bash
make dcd-dev
make dcu-dev
```

## 参考情報

### 7.1 ディレクトリ構造

```
apps/api/
├── cmd/
│   ├── server/main.go          # サーバーエントリーポイント
│   └── seed/main.go            # シードデータ投入ツール
├── internal/
│   ├── gen/                    # 【生成】OpenAPI生成コード
│   │   └── api.gen.go
│   ├── handler/                # 【手動】HTTPハンドラー
│   │   ├── api.handler.go      # StrictServerInterface実装
│   │   ├── todo.handler.go     # Todo CRUD ロジック
│   │   └── auth.handler.go     # 認証ハンドラー（OAuth）
│   ├── service/                # 【手動】ビジネスロジック
│   │   ├── todo.service.go
│   │   ├── user.service.go
│   │   ├── todo.repository.go  # TodoRepository interface
│   │   ├── user.repository.go  # UserRepository interface
│   │   └── mocks/              # 【生成】mockery生成モック
│   ├── mapper/                 # 【手動】DTO変換
│   │   ├── todo.go
│   │   └── user.go
│   ├── router/                 # 【手動】ルーティング設定
│   │   ├── routes.go
│   │   ├── middleware.go
│   │   ├── cors.go
│   │   └── auth.route.go
│   ├── auth/                   # 【手動】認証機能
│   │   ├── session.go          # Redis セッション管理
│   │   ├── middleware.go
│   │   ├── gothic.go           # OAuth設定（Goth）
│   │   ├── provider.go
│   │   └── context.go
│   ├── config/                 # 【手動】設定管理
│   │   └── config.go
│   ├── database/               # 【手動】DB接続管理
│   │   ├── database.go
│   │   └── transaction.go
│   └── seed/                   # 【手動】シードデータ
│       ├── seed.go
│       └── todo.seed.go
├── db/
│   ├── schema.sql              # 【手動】スキーマ定義（真実の源）
│   ├── query/                  # 【手動】SQLクエリ定義
│   │   ├── todo.sql
│   │   └── user.sql
│   ├── sqlc/                   # 【生成】sqlc生成コード
│   │   ├── db.go
│   │   ├── models.go
│   │   ├── querier.go
│   │   ├── todo.sql.go
│   │   └── user.sql.go
│   └── migrations/             # 【生成+手動】Atlas マイグレーション
│       ├── YYYYMMDDHHMMSS_*.sql
│       └── atlas.sum
├── openapi/
│   ├── cue/
│   │   └── api.cue             # 【手動】OpenAPI定義（CUE）
│   ├── oapi-codegen.yaml       # 【手動】コード生成設定
│   └── openapi.yaml            # 【生成】OpenAPI YAML
├── docs/
│   ├── architecture.md
│   └── development.md
├── atlas.hcl                   # Atlas設定
├── sqlc.yaml                   # sqlc設定
└── go.mod / go.sum
```

### 7.2 アーキテクチャパターン概要

**レイヤー構成**:

```
HTTP Request
  ↓
Router (ルーティング & ミドルウェア)
  ↓
Handler (リクエスト検証 & レスポンス構築)
  ↓
Mapper (DTO変換)
  ↓
Service (ビジネスロジック)
  ↓
Repository Interface (抽象化) ← 依存性逆転の原則を適用
  ↑ 実装
sqlc (具象実装)
  ↓
PostgreSQL
```

**依存性逆転の原則（Dependency Inversion Principle）**:
- Service層は Repository Interface（抽象）に依存
- sqlc生成コード（具象実装）が Repository Interface を実装
- これにより、Service層がDB実装の詳細から独立

詳細は [architecture.md](architecture.md) を参照してください。

### 7.3 使用技術スタック

**フレームワーク**:
- **Echo v4** - HTTPフレームワーク

**OpenAPI**:
- **CUE** - 型安全なスキーマ定義言語
- **oapi-codegen** - OpenAPI → Go コード生成

**データベース**:
- **PostgreSQL** - リレーショナルデータベース
- **pgx/v5** - Goの高性能PostgreSQLドライバー
- **sqlc** - SQL → 型安全なGoコード生成
- **Atlas** - スキーママイグレーションツール

**認証**:
- **Goth** - マルチプロバイダーOAuth認証ライブラリ
- **gorilla/sessions** - セッション管理
- **redisstore** - Redisベースのセッションストア

**設定管理**:
- **envconfig** - 環境変数ベースの設定管理

**テスト**:
- **testify** - アサーション & モック
- **mockery** - インターフェースからモック自動生成

**開発環境**:
- **Docker Compose** - コンテナオーケストレーション
- **Air** - Goのホットリロードツール

### 7.4 コード生成フロー図

#### OpenAPI生成フロー

```
┌─────────────────────────┐
│ openapi/cue/api.cue     │  ← 手動編集
└───────────┬─────────────┘
            │ make openapi-gen
            ▼
┌─────────────────────────┐
│ openapi/openapi.yaml    │  ← 自動生成
└───────────┬─────────────┘
            │ make api-gen
            ▼
┌─────────────────────────┐
│ internal/gen/api.gen.go │  ← 自動生成
└─────────────────────────┘
```

#### データベース生成フロー

```
┌─────────────────────────┐
│ db/schema.sql           │  ← 手動編集
└───────────┬─────────────┘
            │ make migrate-diff
            ▼
┌─────────────────────────┐
│ db/migrations/*.sql     │  ← 自動生成
│ atlas.sum               │
└───────────┬─────────────┘
            │ make migrate-apply
            ▼
       PostgreSQL
```

```
┌─────────────────────────┐
│ db/query/*.sql          │  ← 手動編集
└───────────┬─────────────┘
            │ make sqlc-generate
            ▼
┌─────────────────────────┐
│ db/sqlc/*.go            │  ← 自動生成
└─────────────────────────┘
```

#### Mock生成フロー

```
┌─────────────────────────────────┐
│ internal/service/               │  ← 手動定義
│   todo.repository.go            │
│   user.repository.go            │
└───────────┬─────────────────────┘
            │ make mock
            ▼
┌─────────────────────────────────┐
│ internal/service/mocks/         │  ← 自動生成
│   TodoRepository.go             │
│   UserRepository.go             │
└─────────────────────────────────┘
```

### 7.5 コマンドリファレンス

| カテゴリ | コマンド | 説明 |
|---------|---------|------|
| **Docker** | `make dcu-dev` | 開発環境起動（docker compose up） |
| | `make dcd-dev` | 開発環境停止（docker compose down） |
| | `make dcb-dev` | イメージビルド |
| | `make backend-ssh` | サーバーコンテナにアクセス |
| **コード生成** | `make generate` | 全コード生成（OpenAPI + sqlc） |
| | `make openapi-gen` | CUE → OpenAPI YAML |
| | `make api-gen` | OpenAPI YAML → Go |
| | `make sqlc-generate` | SQL → Go（sqlc） |
| | `make mock` | モック生成（mockery） |
| **マイグレーション** | `make create-atlas-dev-db` | atlas_dev データベース作成 |
| | `make migrate-diff NAME=xxx` | マイグレーションファイル生成 |
| | `make migrate-apply` | マイグレーション適用 |
| | `make migrate-status` | マイグレーション状態確認 |
| | `make migrate-hash` | atlas.sum 再計算 |
| | `make schema-update NAME=xxx` | マイグレーション生成 + sqlc生成 |
| **テスト** | `make test` | 全テスト実行 |
| | `go test -short ./...` | 統合テストをスキップ |
| **Lint** | `make lint` | staticcheck 実行 |
| | `make cue-fmt` | CUEファイルフォーマット |
| | `make cue-vet` | CUE定義の検証のみ実行 |
| **ログ確認** | `docker logs -f go_todo_server` | サーバーログをリアルタイム表示 |
| | `docker logs -f go_todo_db` | データベースログをリアルタイム表示 |
| **シード** | `make seed` | シードデータ投入 |
| | `make seed-fresh` | データ削除 + シード投入 |
| **依存関係** | `make go-mod-tidy` | go mod tidy 実行 |
| | `make go-add-library PKG=xxx` | Goパッケージ追加 |

### 7.6 環境変数リファレンス

| 変数名 | 説明 | デフォルト/例 |
|-------|------|-------------|
| **Database** | | |
| `POSTGRES_HOST` | PostgreSQLホスト | `go_todo_db` |
| `POSTGRES_PORT` | PostgreSQLポート | `5432` |
| `POSTGRES_DB` | データベース名 | `go_todo_db` |
| `POSTGRES_USER` | データベースユーザー | `user` |
| `POSTGRES_PASSWORD` | データベースパスワード | `password` |
| **Redis** | | |
| `REDIS_HOST` | Redisホスト | `go_todo_redis` |
| `REDIS_PORT` | Redisポート | `6379` |
| **OAuth** | | |
| `GOOGLE_CLIENT_ID` | Google OAuth クライアントID | - |
| `GOOGLE_CLIENT_SECRET` | Google OAuth クライアントシークレット | - |
| `OAUTH_CALLBACK_URL` | OAuthコールバックURL | `http://localhost:8080/auth/google/callback` |
| **Server** | | |
| `SERVER_PORT` | サーバーポート | `8080` |
| **Frontend** | | |
| `FRONTEND_URL` | フロントエンドURL（CORS用） | `http://localhost:3000` |
| **Cookie** | | |
| `COOKIE_SECURE` | Cookie Secure フラグ | `false`（本番: `true`） |

### 7.7 atlas_dev データベースについて

**概要**:

`atlas_dev` は、Atlasが**差分計算のために使用する一時的な作業用データベース**です。

**役割**:

1. `atlas_dev` に現在のスキーマ（`db/schema.sql`）を適用
2. 実際のDB（`go_todo_db`）と `atlas_dev` を比較
3. 差分をマイグレーションファイルとして生成

**特徴**:

- データは保存されず、毎回クリーンな状態で使用される
- 本番DBとは完全に分離されているため安全
- マイグレーション生成時のみ使用される

### 7.8 atlas.sum について

**概要**:

`atlas.sum` は、マイグレーションファイルの**整合性を保証**するためのハッシュファイルです。

**内容**:

各マイグレーションファイルのSHA256ハッシュ値を記録：

```
h1:abc123... 20250125100000_init.sql
h1:def456... 20250125120000_add_priority.sql
```

**重要性**:

- **必ずGitで管理する**（コミット必須）
- マイグレーションファイルを手動編集した場合は `make migrate-hash` で再計算
- チームメンバー間でマイグレーションの整合性を保証

### 7.9 セキュリティのベストプラクティス

#### 環境変数の管理

1. **`.env` ファイルの取り扱い**
   - `.env` ファイルは `.gitignore` に含めること（既に設定済み）
   - 本番環境では環境変数を直接設定（`.env` ファイルを使用しない）
   - シークレット情報は絶対にコミットしない

2. **本番環境での設定**
   - `COOKIE_SECURE=true` を必ず設定
   - `POSTGRES_PASSWORD` などのシークレットは環境変数管理サービス（AWS Secrets Manager等）を使用

#### Cookie設定

本番環境では以下を確認：
- `COOKIE_SECURE=true` - HTTPS接続でのみCookieを送信
- `HttpOnly` 属性 - JavaScriptからのアクセスを防止（gorilla/sessionsで自動設定）
- `SameSite` 属性 - CSRF攻撃対策（設定で調整可能）

#### SQLインジェクション対策

- **sqlcの使用**: パラメータ化クエリが自動的に適用されるため安全
- **生のSQL実行は避ける**: 必ず sqlc を介してクエリを実行
- **ユーザー入力の検証**: Handler層で適切にバリデーション

#### OAuth認証

- Google OAuth クライアントシークレットは厳重に管理
- コールバックURLは本番環境のドメインを正確に設定
- 開発環境と本番環境で異なるOAuthアプリケーションを使用推奨

### 7.10 ベストプラクティス

#### スキーマ管理

1. **db/schema.sql を真実の源として管理**
   - スキーマ変更は必ず `db/schema.sql` から行う
   - マイグレーションファイルは差分として生成される

2. **小さく頻繁にマイグレーション**
   - 1つのマイグレーションで1つの変更
   - 複雑な変更は複数のマイグレーションに分割

3. **本番適用前にdev環境でテスト**
   - 必ず開発環境でテストしてから本番適用
   - ロールバック手順も事前に確認

4. **破壊的変更に注意**
   - カラム削除やテーブル削除は慎重に
   - 可能であれば段階的に実施:
     1. 新カラム追加 → アプリ移行
     2. 旧カラム削除

5. **atlas.sum は必ずコミット**
   - マイグレーションファイルと同時にコミット
   - CIでチェックサム検証を行う

#### コード生成

1. **生成コードは編集しない**
   - `internal/gen/api.gen.go` は直接編集禁止
   - `db/sqlc/*.go` は直接編集禁止
   - 変更が必要な場合は元の定義（CUE, SQL）を編集

2. **operationIdの命名規則**
   - PascalCaseで記述（Goのメソッド名になる）
   - 例: `ListTodos`, `GetTodoById`

3. **生成ファイルのコミット**
   - 生成されたファイルもGitにコミット
   - チームメンバー全員が同じコードベースを共有

### 7.11 関連リンク

- **Atlas**: [Atlas Documentation](https://atlasgo.io/docs)
- **sqlc**: [sqlc Documentation](https://docs.sqlc.dev/)
- **oapi-codegen**: [oapi-codegen GitHub](https://github.com/oapi-codegen/oapi-codegen)
- **CUE**: [CUE Language](https://cuelang.org/)
- **Echo**: [Echo Framework](https://echo.labstack.com/)
- **pgx**: [pgx Documentation](https://github.com/jackc/pgx)
- **Goth**: [Goth GitHub](https://github.com/markbates/goth)
- **mockery**: [mockery GitHub](https://github.com/vektra/mockery)
- **Air**: [Air GitHub](https://github.com/cosmtrek/air)

---

より詳細なアーキテクチャ情報については、[architecture.md](architecture.md) を参照してください。
