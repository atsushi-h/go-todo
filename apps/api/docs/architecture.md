# アーキテクチャドキュメント

このドキュメントでは、Go APIのアーキテクチャ設計思想と構造を説明します。

## 目次

- [アーキテクチャ設計哲学](#アーキテクチャ設計哲学)
- [システム全体像](#システム全体像)
- [ディレクトリ構成と責務](#ディレクトリ構成と責務)
- [レイヤーアーキテクチャ詳細](#レイヤーアーキテクチャ詳細)
- [コード生成戦略](#コード生成戦略)
- [重要パターンの解説](#重要パターンの解説)
- [リクエスト/レスポンスのデータフロー](#リクエストレスポンスのデータフロー)
- [認証・認可の仕組み](#認証認可の仕組み)
- [新機能追加ガイド](#新機能追加ガイド)
- [テスト戦略](#テスト戦略)
- [なぜこのアーキテクチャなのか](#なぜこのアーキテクチャなのか)

## アーキテクチャ設計哲学

このアプリケーションは、以下の4つの原則に基づいて設計されています：

### 1. OpenAPI First Development

API仕様を**CUE言語**で定義し、そこからOpenAPI YAML、さらにGoコードを自動生成します。

**利点**:
- API仕様が常にコードと同期
- フロントエンドとの契約が明確
- 型安全なハンドラーの自動生成

### 2. Type Safety at Every Layer

2つの層で型安全性を確保しています：

1. **HTTP層**: OpenAPI → oapi-codegen による型安全なリクエスト/レスポンス
2. **データベース層**: SQL → sqlc による型安全なクエリ実行

コンパイル時にほとんどのエラーを検出できます。

### 3. Clear Separation of Concerns

各レイヤーは明確な責務を持ち、**依存性逆転の原則（Dependency Inversion Principle）** を適用しています：

```
HTTP Request
  ↓
Router (ルーティング & ミドルウェア)
  ↓
Handler (リクエスト検証 & レスポンス構築)
  ↓
Service (ビジネスロジック)
  ↓
Repository Interface (抽象化) ← 依存性の逆転
  ↑ 実装
sqlc (具象実装)
  ↓
DB (PostgreSQL)
```

**依存性逆転のポイント**:
- Service層は **Repository Interface** に依存（具象型に依存しない）
- sqlc生成コードが Repository Interface を実装
- これにより、Service層がDB実装の詳細から独立

### 4. Code Generation for Consistency

手動実装によるミスを減らし、一貫性を保つため、可能な限りコード生成を活用します：

- **OpenAPI → Go**: oapi-codegen
- **SQL → Go**: sqlc
- **Interface → Mock**: mockery

## システム全体像

### レイヤー構造

```
┌─────────────────────────────────────────┐
│         HTTP Request (Echo)             │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│  Router Layer (internal/router)         │
│  - CORS, Logger, Recover                │
│  - 認証ミドルウェア                       │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│  Handler Layer (internal/handler)       │
│  - リクエスト検証                         │
│  - 認証チェック (Context経由)            │
│  - エラーハンドリング                     │
│  - レスポンス構築                         │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│  Mapper Layer (internal/mapper)         │
│  - DB型 ⇄ API型の変換                   │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│  Service Layer (internal/service)       │
│  - ビジネスロジック                       │
│  - トランザクション管理                   │
│  - バッチ処理                            │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│  Repository Layer (interface)           │
│  - データアクセス抽象化                   │
│  - テスタビリティ向上                     │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│  sqlc Layer (db/sqlc)                   │
│  - 型安全なSQL実行                       │
│  - 自動生成コード                         │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│         PostgreSQL                      │
└─────────────────────────────────────────┘
```

### 依存関係の方向性

**依存性逆転の原則（Dependency Inversion Principle）を適用**:

- Handler は Service に依存
- Service は **Repository Interface** に依存（具象型ではない）
- sqlc生成コード（`*sqlc.Queries`）が Repository Interface を**実装**
- DB層の変更がService層に影響しない

**重要**: Service層は抽象（インターフェース）に依存し、具象実装（sqlc）には依存しません。これにより：
- **テスタビリティ**: モックに差し替え可能
- **柔軟性**: DB実装を変更してもService層は影響を受けない
- **保守性**: 各層が独立して進化できる

## ディレクトリ構成と責務

```
apps/api/
├── cmd/
│   ├── server/main.go          # サーバーエントリーポイント
│   └── seed/main.go            # シードデータ投入ツール
├── internal/
│   ├── gen/                    # 【生成】OpenAPI生成コード
│   │   └── api.gen.go          # StrictServerInterface, 型定義
│   ├── handler/                # 【手動】HTTPハンドラー
│   │   ├── api.handler.go      # StrictServerInterface実装
│   │   ├── todo.handler.go     # Todo CRUD ロジック
│   │   └── auth.handler.go     # 認証ハンドラー（OAuth）
│   ├── service/                # 【手動】ビジネスロジック
│   │   ├── todo.service.go     # Todoビジネスロジック
│   │   ├── user.service.go     # Userビジネスロジック
│   │   ├── todo.repository.go  # TodoRepository interface
│   │   ├── user.repository.go  # UserRepository interface
│   │   └── mocks/              # 【生成】mockery生成モック
│   ├── mapper/                 # 【手動】DTO変換
│   │   ├── todo.go             # sqlc.Todo → gen.Todo
│   │   └── user.go             # sqlc.User → gen.User
│   ├── router/                 # 【手動】ルーティング設定
│   │   ├── routes.go           # メインルーター & ミドルウェア
│   │   ├── middleware.go       # カスタムミドルウェア
│   │   ├── cors.go             # CORS設定
│   │   └── auth.route.go       # 認証ルート（OAuth）
│   ├── auth/                   # 【手動】認証機能
│   │   ├── session.go          # Redis セッション管理
│   │   ├── middleware.go       # 認証ミドルウェア
│   │   ├── gothic.go           # OAuth設定（Goth）
│   │   ├── provider.go         # OAuth プロバイダー設定
│   │   └── context.go          # Context ヘルパー
│   ├── config/                 # 【手動】設定管理
│   │   └── config.go           # envconfig ベース設定
│   ├── database/               # 【手動】DB接続管理
│   │   ├── database.go         # 接続プール管理
│   │   └── transaction.go      # トランザクション ヘルパー
│   └── seed/                   # 【手動】シードデータ
│       ├── seed.go             # シード実行ロジック
│       └── todo.seed.go        # Todo シードデータ
├── db/
│   ├── schema.sql              # 【手動】スキーマ定義（真実の源）
│   ├── query/                  # 【手動】SQLクエリ定義
│   │   ├── todo.sql            # Todo関連クエリ
│   │   └── user.sql            # User関連クエリ
│   ├── sqlc/                   # 【生成】sqlc生成コード
│   │   ├── db.go
│   │   ├── models.go
│   │   ├── todo.sql.go
│   │   └── user.sql.go
│   └── migrations/             # 【生成+手動】Atlas マイグレーション
│       ├── YYYYMMDDHHMMSS_*.sql
│       └── atlas.sum           # チェックサム（自動生成）
├── openapi/
│   ├── cue/
│   │   └── api.cue             # 【手動】OpenAPI定義（CUE）
│   ├── oapi-codegen.yaml       # 【手動】コード生成設定
│   └── openapi.yaml            # 【生成】OpenAPI YAML
├── docs/                       # ドキュメント
│   ├── architecture.md
│   └── development.md
├── atlas.hcl                   # Atlas設定
├── sqlc.yaml                   # sqlc設定
└── go.mod / go.sum
```

### 重要な区分

**生成コード（編集禁止）**:
- `internal/gen/api.gen.go`
- `db/sqlc/*.go`
- `openapi/openapi.yaml`
- `internal/service/mocks/*.go`

**手動実装**:
- `internal/handler/*.go`
- `internal/service/*.service.go`
- `internal/mapper/*.go`
- `db/schema.sql`, `db/query/*.sql`
- `openapi/cue/api.cue`

## レイヤーアーキテクチャ詳細

### 4.1 Router層（ルーティング & ミドルウェア）

**責務**:
- HTTPルーティング
- グローバルミドルウェアの適用
- 認証ミドルウェアの設定

**実装例** ([internal/router/routes.go:16-36](internal/router/routes.go)):

```go
func SetupRoutes(e *echo.Echo, apiHandler *handler.APIHandler, authHandler *handler.AuthHandler, sm *auth.SessionManager, frontendConfig config.FrontendConfig) {
    // グローバルミドルウェア
    e.Use(middleware.CORSWithConfig(CORSConfig(frontendConfig)))
    e.Use(middleware.Recover())
    e.Use(middleware.Logger())

    // 認証ミドルウェアをstrictmiddlewareとしてラップ
    authMiddleware := createAuthMiddleware(sm)

    // StrictハンドラーをEchoハンドラーにラップ（認証ミドルウェア付き）
    strictHandler := gen.NewStrictHandler(apiHandler, []gen.StrictMiddlewareFunc{authMiddleware})

    // 生成されたルート登録関数を使用
    gen.RegisterHandlers(e, strictHandler)

    // 認証関連のルート（手動で設定）
    SetupAuthRoutes(e, authHandler, sm)

    // カスタムエラーハンドラー
    e.HTTPErrorHandler = customHTTPErrorHandler
}
```

**ポイント**:
- `gen.NewStrictHandler` で認証ミドルウェアを統合
- `gen.RegisterHandlers` で生成されたルートを自動登録
- OAuth認証ルートは手動で設定（複雑なフローのため）

### 4.2 Handler層（HTTPリクエスト処理）

**責務**:
- リクエストパラメータの検証
- ユーザー認証チェック
- サービス層の呼び出し
- エラーハンドリング
- レスポンスオブジェクトの構築

**実装例** ([internal/handler/todo.handler.go:24-37](internal/handler/todo.handler.go)):

```go
func (h *TodoHandler) ListTodos(ctx context.Context, request gen.ListTodosRequestObject) (gen.ListTodosResponseObject, error) {
    // Context から認証済みユーザーIDを取得
    userID, ok := auth.GetUserIDFromContext(ctx)
    if !ok {
        return gen.ListTodos401JSONResponse{Message: "Unauthorized"}, nil
    }

    // サービス層を呼び出し
    todos, err := h.service.GetAllTodos(ctx, userID)
    if err != nil {
        return gen.ListTodos500JSONResponse{Message: "Internal server error"}, nil
    }

    // Mapper を使用してレスポンスを構築
    return gen.ListTodos200JSONResponse(mapper.TodosToResponse(todos)), nil
}
```

**特徴**:
- **型安全なリクエスト/レスポンス**: `gen.ListTodosRequestObject`, `gen.ListTodosResponseObject`
- **明示的なステータスコード**: `gen.ListTodos200JSONResponse`, `gen.ListTodos401JSONResponse`
- **Context経由の認証情報**: ミドルウェアで設定されたユーザーIDを取得

### 4.3 Service層（ビジネスロジック）

**責務**:
- ビジネスロジックの実装
- トランザクション管理
- バッチ処理
- エラー変換（`pgx.ErrNoRows` → `ErrTodoNotFound`）

**実装例** ([internal/service/todo.service.go:92-139](internal/service/todo.service.go)):

```go
func (s *TodoService) BatchCompleteTodos(ctx context.Context, userID int64, ids []int64) (*BatchCompleteResult, error) {
    result := &BatchCompleteResult{
        Succeeded: []sqlc.Todo{},
        Failed:    []BatchFailedItem{},
    }

    // 存在チェック
    existingTodos, err := s.repo.GetTodosByIDs(ctx, sqlc.GetTodosByIDsParams{
        Ids:    ids,
        UserID: userID,
    })
    if err != nil {
        return nil, err
    }

    // 存在するIDのマップを作成
    existingIDMap := make(map[int64]bool)
    for _, todo := range existingTodos {
        existingIDMap[todo.ID] = true
    }

    // 存在しないIDを失敗として記録
    var validIDs []int64
    for _, id := range ids {
        if existingIDMap[id] {
            validIDs = append(validIDs, id)
        } else {
            result.Failed = append(result.Failed, BatchFailedItem{
                ID:    id,
                Error: "Todo not found",
            })
        }
    }

    // 有効なIDがある場合のみバッチ更新
    if len(validIDs) > 0 {
        completedTodos, err := s.repo.BatchCompleteTodos(ctx, sqlc.BatchCompleteTodosParams{
            Ids:    validIDs,
            UserID: userID,
        })
        if err != nil {
            return nil, err
        }
        result.Succeeded = completedTodos
    }

    return result, nil
}
```

**特徴**:
- **複雑なビジネスロジック**: バッチ処理の成功/失敗を分離
- **Repository インターフェース経由**: テスト可能な設計
- **ドメインエラーの定義**: `ErrTodoNotFound` など

### 4.4 Repository層（データアクセス抽象化）

**責務**:
- データアクセスの抽象化
- **依存性逆転の原則**を適用（Service層が具象実装に依存しない）
- テスタビリティの向上（モック化可能）

**実装例** ([internal/service/todo.repository.go](internal/service/todo.repository.go)):

```go
type TodoRepository interface {
    ListTodosByUser(ctx context.Context, userID int64) ([]sqlc.Todo, error)
    GetTodoByID(ctx context.Context, arg sqlc.GetTodoByIDParams) (sqlc.Todo, error)
    CreateTodo(ctx context.Context, arg sqlc.CreateTodoParams) (sqlc.Todo, error)
    UpdateTodo(ctx context.Context, arg sqlc.UpdateTodoParams) (sqlc.Todo, error)
    DeleteTodo(ctx context.Context, arg sqlc.DeleteTodoParams) error
    GetTodosByIDs(ctx context.Context, arg sqlc.GetTodosByIDsParams) ([]sqlc.Todo, error)
    BatchCompleteTodos(ctx context.Context, arg sqlc.BatchCompleteTodosParams) ([]sqlc.Todo, error)
    BatchDeleteTodos(ctx context.Context, arg sqlc.BatchDeleteTodosParams) error
}
```

**依存性逆転の実現**:

通常、上位層が下位層に依存しますが、このアーキテクチャでは**インターフェースを介して依存関係を逆転**させています。

```
従来の依存関係:
Service (上位) → sqlc (下位)    ← Service が具象実装に直接依存

依存性逆転後:
Service (上位) → Repository Interface (抽象)
                         ↑
                    sqlc (下位) が実装    ← 下位層が上位層の定義に従う
```

**巧妙な設計**:

`sqlc.Querier` インターフェースは、sqlcが自動生成するすべてのクエリメソッドを含んでいます。

```go
// sqlc が自動生成
type Querier interface {
    ListTodosByUser(ctx context.Context, userID int64) ([]Todo, error)
    GetTodoByID(ctx context.Context, arg GetTodoByIDParams) (Todo, error)
    // ... その他のメソッド
}
```

私たちの `TodoRepository` は、この `Querier` と**同じシグネチャ**を持つように設計されています。

**結果**: `sqlc.Querier` は自動的に `TodoRepository` インターフェースを満たします。追加のラッパーコードは不要です。

```go
// main.go で直接使用可能
queries := sqlc.New(pool)  // *sqlc.Queries は sqlc.Querier を実装
todoService := service.NewTodoService(queries)  // TodoRepository として渡せる
```

**依存性逆転のメリット**:
1. **Service層がDB実装の詳細を知らない** - sqlcを別のORMに変更してもService層は影響を受けない
2. **テストが容易** - Repository をモックに差し替え可能
3. **コンパイル時の型チェック** - インターフェースを満たさないとコンパイルエラー

### 4.5 Mapper層（DTO変換）

**責務**:
- DB型（`sqlc.Todo`）と API型（`gen.Todo`）の変換
- 各レイヤーの独立性を保証

**実装例** ([internal/mapper/todo.go:9-19](internal/mapper/todo.go)):

```go
func TodoToResponse(t *sqlc.Todo) gen.Todo {
    return gen.Todo{
        Id:          t.ID,
        Title:       t.Title,
        Description: t.Description,
        Completed:   t.Completed,
        UserId:      t.UserID,
        CreatedAt:   t.CreatedAt,
        UpdatedAt:   t.UpdatedAt,
    }
}

func TodosToResponse(todos []sqlc.Todo) []gen.Todo {
    result := make([]gen.Todo, len(todos))
    for i := range todos {
        result[i] = TodoToResponse(&todos[i])
    }
    return result
}
```

**なぜ必要か**:

1. **DB スキーマの変更がAPI仕様に影響しない**
2. **API仕様の変更がDB設計に影響しない**
3. **各層が独立して進化できる**

例: DBカラム名が `deleted_at` でも、APIレスポンスには含めないといった制御が可能。

### 4.6 Generated層（自動生成コード）

**internal/gen/api.gen.go** には以下が含まれます：

1. **型定義**:
   ```go
   type Todo struct {
       Id          int64       `json:"id"`
       Title       string      `json:"title"`
       Description *string     `json:"description,omitempty"`
       Completed   bool        `json:"completed"`
       UserId      int64       `json:"userId"`
       CreatedAt   time.Time   `json:"createdAt"`
       UpdatedAt   time.Time   `json:"updatedAt"`
   }
   ```

2. **StrictServerInterface**:
   ```go
   type StrictServerInterface interface {
       GetInfo(ctx context.Context, request GetInfoRequestObject) (GetInfoResponseObject, error)
       ListTodos(ctx context.Context, request ListTodosRequestObject) (ListTodosResponseObject, error)
       CreateTodo(ctx context.Context, request CreateTodoRequestObject) (CreateTodoResponseObject, error)
       // ... その他のメソッド
   }
   ```

3. **Request/Response型**:
   ```go
   type ListTodosRequestObject struct {}

   type ListTodosResponseObject interface {
       VisitListTodosResponse(w http.ResponseWriter) error
   }

   type ListTodos200JSONResponse []Todo
   type ListTodos401JSONResponse ErrorResponse
   type ListTodos500JSONResponse ErrorResponse
   ```

4. **ルート登録関数**:
   ```go
   func RegisterHandlers(router EchoRouter, si ServerInterface)
   ```

## コード生成戦略

### 5.1 OpenAPI生成フロー（CUE → YAML → Go）

```
┌─────────────────────────┐
│ openapi/cue/api.cue     │  ← CUE言語でAPI定義（手動編集）
│ - 型安全なスキーマ定義    │
│ - バリデーション         │
└───────────┬─────────────┘
            │ make openapi-gen (cue export)
            ▼
┌─────────────────────────┐
│ openapi/openapi.yaml    │  ← OpenAPI 3.0 YAML（自動生成）
│ - 標準的なAPI仕様        │
└───────────┬─────────────┘
            │ make api-gen (oapi-codegen)
            ▼
┌─────────────────────────┐
│ internal/gen/api.gen.go │  ← Go コード（自動生成）
│ - StrictServerInterface  │
│ - 型定義                 │
│ - ルーティング           │
└─────────────────────────┘
```

**コマンド**:

```bash
# 一括生成
make generate

# 段階的に実行
make openapi-gen  # CUE → YAML
make api-gen      # YAML → Go
```

**設定ファイル**: [openapi/oapi-codegen.yaml](openapi/oapi-codegen.yaml)

```yaml
package: gen
output: internal/gen/api.gen.go
generate:
  strict-server: true   # StrictServer形式を生成
  embedded-spec: true
```

### 5.2 Database生成フロー（SQL → Go）

```
┌─────────────────────────┐
│ db/schema.sql           │  ← スキーマ定義（真実の源）
└───────────┬─────────────┘
            │ make migrate-diff
            ▼
┌─────────────────────────┐
│ db/migrations/*.sql     │  ← マイグレーションファイル（自動生成）
│ atlas.sum               │  ← チェックサム（自動生成）
└───────────┬─────────────┘
            │ make migrate-apply
            ▼
       PostgreSQL
```

```
┌─────────────────────────┐
│ db/query/todo.sql       │  ← SQLクエリ定義（手動編集）
│ db/query/user.sql       │
└───────────┬─────────────┘
            │ make sqlc-generate
            ▼
┌─────────────────────────┐
│ db/sqlc/models.go       │  ← 型安全なGoコード（自動生成）
│ db/sqlc/todo.sql.go     │
│ db/sqlc/user.sql.go     │
└─────────────────────────┘
```

**コマンド**:

```bash
# スキーマ変更時
make migrate-diff NAME=add_priority_column
make migrate-apply
make sqlc-generate

# または一括で
make schema-update NAME=add_priority_column
```

### 5.3 Mock生成フロー

```
┌─────────────────────────────────┐
│ internal/service/               │
│   todo.repository.go            │  ← Repository interface
│   user.repository.go            │
└───────────┬─────────────────────┘
            │ make mock
            ▼
┌─────────────────────────────────┐
│ internal/service/mocks/         │  ← モックコード（自動生成）
│   TodoRepository.go             │
│   UserRepository.go             │
└─────────────────────────────────┘
```

**設定ファイル**: `.mockery.yaml`

```yaml
with-expecter: true
dir: "internal/service/mocks"
packages:
  go-todo/internal/service:
    interfaces:
      TodoRepository:
      UserRepository:
```

## 重要パターンの解説

### 6.1 StrictServer パターン

**通常のハンドラー（strict-server: false）**:

```go
func (h *Handler) GetTodo(c echo.Context, id int) error {
    // echo.Context を直接操作
    todo, err := h.service.GetByID(id)
    if err != nil {
        return c.JSON(404, map[string]string{"error": "not found"})
    }
    return c.JSON(200, todo)
}
```

**問題点**:
- 型安全性が低い（`map[string]string` は任意の構造）
- ステータスコードのtypoはランタイムまで検出されない
- OpenAPI仕様との不一致がコンパイル時に検出されない

**StrictServerパターン（strict-server: true）**:

```go
func (h *Handler) GetTodo(ctx context.Context, request gen.GetTodoRequestObject) (gen.GetTodoResponseObject, error) {
    // 型付きリクエスト/レスポンス
    todo, err := h.service.GetByID(request.Id)
    if err != nil {
        return gen.GetTodo404JSONResponse{Message: "not found"}, nil
    }
    return gen.GetTodo200JSONResponse(*todo), nil
}
```

**利点**:
- **コンパイル時型チェック**: 間違ったレスポンス型を返すとコンパイルエラー
- **OpenAPI仕様との一貫性**: 定義されたステータスコードの型が自動生成される
- **テストが書きやすい**: `echo.Context` への依存がない

### 6.2 Repository パターン（依存性逆転の原則）

**依存性逆転の原則（Dependency Inversion Principle）の実装例**:

このアーキテクチャでは、Repository パターンを使って依存性逆転を実現しています。

**従来のアプローチ（依存性逆転なし）**:

```go
// Service が具象実装に直接依存
type TodoService struct {
    queries *sqlc.Queries  // ← 具象型への直接依存
}

// 問題点:
// - sqlc の実装を変更すると Service も変更が必要
// - テスト時にモックに差し替えられない
// - DB実装の詳細が Service 層に漏れる
```

**このアーキテクチャ（依存性逆転あり）**:

```go
// 1. Service 層がインターフェースを定義（上位層が抽象を定義）
// internal/service/todo.repository.go
type TodoRepository interface {
    ListTodosByUser(ctx context.Context, userID int64) ([]sqlc.Todo, error)
    GetTodoByID(ctx context.Context, arg sqlc.GetTodoByIDParams) (sqlc.Todo, error)
    CreateTodo(ctx context.Context, arg sqlc.CreateTodoParams) (sqlc.Todo, error)
    // ... その他のメソッド
}

// 2. Service がインターフェースに依存
type TodoService struct {
    repo TodoRepository  // ← 抽象（インターフェース）に依存
}

// 3. sqlc が自動生成する Querier がインターフェースを実装（下位層が上位層に従う）
// db/sqlc/querier.go (sqlc自動生成)
type Querier interface {
    ListTodosByUser(ctx context.Context, userID int64) ([]Todo, error)
    GetTodoByID(ctx context.Context, arg GetTodoByIDParams) (Todo, error)
    CreateTodo(ctx context.Context, arg CreateTodoParams) (Todo, error)
    // ... その他のメソッド
}

// 4. *sqlc.Queries が Querier を実装 → TodoRepository を満たす
```

**依存関係の方向**:

```
従来:
Service → sqlc (具象) ← 上位が下位に依存

依存性逆転後:
Service → TodoRepository (抽象)
              ↑
         sqlc (具象) が実装  ← 下位が上位の定義に従う
```

**実際の使用例**:

```go
// cmd/server/main.go
queries := sqlc.New(pool)  // *sqlc.Queries は Querier を実装
todoService := service.NewTodoService(queries)  // TodoRepository として渡せる！
```

追加のラッパーコード不要で、自動的にインターフェースを満たします。

**テスト時の利点**:

```go
// テストコード - モックに差し替え可能
mockRepo := mocks.NewTodoRepository(t)
mockRepo.EXPECT().ListTodosByUser(mock.Anything, int64(1)).Return([]sqlc.Todo{...}, nil)
service := service.NewTodoService(mockRepo)
```

**依存性逆転の利点**:
1. **柔軟性**: sqlcを別のORMに変更しても、Service層は変更不要
2. **テスタビリティ**: モックに簡単に差し替え可能
3. **保守性**: DB層とService層が独立して進化できる
4. **型安全性**: インターフェースを満たさないとコンパイルエラー

### 6.3 Mapper パターン

**なぜ必要か**:

DB型とAPI型を直接結合すると、以下の問題が発生します：

1. **DBスキーマ変更がAPIに影響**: カラム名変更 → API互換性破壊
2. **API変更がDBに影響**: 新しいフィールド追加 → DB設計の制約
3. **不要な情報の漏洩**: `deleted_at` などの内部フィールドがAPIに露出

**Mapperによる分離**:

```go
// DB型（sqlc生成）
type Todo struct {
    ID          int64
    UserID      int64
    Title       string
    Description *string
    Completed   bool
    DeletedAt   *time.Time  // ← APIには含めない
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// API型（oapi-codegen生成）
type Todo struct {
    Id          int64      `json:"id"`
    UserId      int64      `json:"userId"`
    Title       string     `json:"title"`
    Description *string    `json:"description,omitempty"`
    Completed   bool       `json:"completed"`
    CreatedAt   time.Time  `json:"createdAt"`
    UpdatedAt   time.Time  `json:"updatedAt"`
}

// Mapper
func TodoToResponse(t *sqlc.Todo) gen.Todo {
    return gen.Todo{
        Id:          t.ID,
        UserId:      t.UserID,
        Title:       t.Title,
        Description: t.Description,
        Completed:   t.Completed,
        CreatedAt:   t.CreatedAt,
        UpdatedAt:   t.UpdatedAt,
        // DeletedAt は含めない
    }
}
```

### 6.4 トランザクション管理

**ヘルパー関数**: [internal/database/transaction.go](internal/database/transaction.go)

```go
func WithTransaction(ctx context.Context, pool *pgxpool.Pool, fn func(tx pgx.Tx) error) error {
    tx, err := pool.Begin(ctx)
    if err != nil {
        return err
    }

    defer func() {
        if p := recover(); p != nil {
            tx.Rollback(ctx)
            panic(p)
        } else if err != nil {
            tx.Rollback(ctx)
        } else {
            err = tx.Commit(ctx)
        }
    }()

    err = fn(tx)
    return err
}
```

**使用例**:

```go
func (s *UserService) CreateUserWithProfile(ctx context.Context, email, name string) error {
    return database.WithTransaction(ctx, s.pool, func(tx pgx.Tx) error {
        queries := s.queries.WithTx(tx)

        user, err := queries.CreateUser(ctx, CreateUserParams{...})
        if err != nil {
            return err
        }

        _, err = queries.CreateProfile(ctx, CreateProfileParams{
            UserID: user.ID,
            ...
        })
        return err
    })
}
```

## リクエスト/レスポンスのデータフロー

具体例: **Todo作成のリクエスト**

```
1. HTTP POST /api/todos
   Body: {"title": "Buy milk", "description": "At supermarket"}

2. ↓ Router層
   - CORS, Logger, Recover ミドルウェア適用
   - 認証ミドルウェア実行
     - セッションからユーザーID取得
     - Context に userID=123 を設定

3. ↓ Handler層 (todo.handler.go:62-82)
   - CreateTodo(ctx, gen.CreateTodoRequestObject)
   - Context から userID=123 を取得
   - リクエストボディの検証（title が空でないか）
   - サービス層呼び出し:
     service.CreateTodo(ctx, 123, "Buy milk", "At supermarket")

4. ↓ Service層 (todo.service.go:56-66)
   - ビジネスロジック実行
   - Repository呼び出し:
     repo.CreateTodo(ctx, CreateTodoParams{
       UserID: 123,
       Title: "Buy milk",
       Description: "At supermarket",
     })

5. ↓ Repository層 (sqlc.Queries)
   - SQL実行: INSERT INTO todos (user_id, title, description, ...) VALUES (...)
   - 結果: sqlc.Todo{ID: 456, UserID: 123, Title: "Buy milk", ...}

6. ↑ Service層
   - sqlc.Todo を返却

7. ↑ Handler層
   - Mapper で変換: mapper.TodoToResponse(&todo)
   - レスポンス構築:
     gen.CreateTodo201JSONResponse{
       Id: 456,
       UserId: 123,
       Title: "Buy milk",
       Description: "At supermarket",
       Completed: false,
       CreatedAt: 2025-01-01T10:00:00Z,
       UpdatedAt: 2025-01-01T10:00:00Z,
     }

8. ↑ Router層
   - StrictHandler が自動的に JSON にシリアライズ
   - HTTP 201 Created
   - Body: {"id":456,"userId":123,"title":"Buy milk",...}
```

## 認証・認可の仕組み

### OAuth認証フロー

```
1. フロントエンド: GET /auth/google
   ↓
2. Goth: Googleの認証ページへリダイレクト
   ↓
3. ユーザー: Googleでログイン
   ↓
4. Google: GET /auth/google/callback?code=xxx
   ↓
5. Goth: codeを使ってアクセストークン取得
   ↓
6. AuthHandler: ユーザー情報取得（email, name, avatar）
   ↓
7. UserService: ユーザー作成 or 更新
   ↓
8. SessionManager: Redis にセッション保存（userID: 123）
   ↓
9. AuthHandler: フロントエンドへリダイレクト + Cookie設定
```

### セッション管理

**Redis ベースのセッション**:

- **Store**: `redisstore` (gorilla/sessions)
- **Key**: セッションID（Cookie に保存）
- **Value**: `{"userID": 123}`
- **有効期限**: 7日間

**実装** ([internal/auth/session.go](internal/auth/session.go)):

```go
type SessionManager struct {
    store *redisstore.RedisStore
}

func (sm *SessionManager) SetUserID(w http.ResponseWriter, r *http.Request, userID int64) error {
    session, _ := sm.store.Get(r, "session")
    session.Values["userID"] = userID
    return session.Save(r, w)
}

func (sm *SessionManager) GetUserID(r *http.Request) (int64, error) {
    session, _ := sm.store.Get(r, "session")
    userID, ok := session.Values["userID"].(int64)
    if !ok {
        return 0, errors.New("user not authenticated")
    }
    return userID, nil
}
```

### ミドルウェア実装（StrictMiddleware）

**実装** ([internal/router/routes.go:39-60](internal/router/routes.go)):

```go
func createAuthMiddleware(sm *auth.SessionManager) gen.StrictMiddlewareFunc {
    return func(f gen.StrictHandlerFunc, operationID string) gen.StrictHandlerFunc {
        return func(ctx echo.Context, request interface{}) (interface{}, error) {
            // 認証不要なエンドポイントをスキップ
            if operationID == "GetInfo" || operationID == "GetHealth" {
                return f(ctx, request)
            }

            // セッションからユーザーIDを取得
            userID, err := sm.GetUserID(ctx.Request())
            if err != nil {
                return nil, echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
            }

            // コンテキストにユーザーIDを設定
            newCtx := auth.WithUserID(ctx.Request().Context(), userID)
            ctx.SetRequest(ctx.Request().WithContext(newCtx))

            return f(ctx, request)
        }
    }
}
```

**ポイント**:
- `operationID` ベースで認証除外を制御
- Context にユーザーIDを設定（各ハンドラーで利用可能）

### Context経由でのユーザーID伝播

**設定** ([internal/auth/context.go](internal/auth/context.go)):

```go
type contextKey string

const userIDKey contextKey = "userID"

func WithUserID(ctx context.Context, userID int64) context.Context {
    return context.WithValue(ctx, userIDKey, userID)
}

func GetUserIDFromContext(ctx context.Context) (int64, bool) {
    userID, ok := ctx.Value(userIDKey).(int64)
    return userID, ok
}
```

**利用** ([internal/handler/todo.handler.go:26-29](internal/handler/todo.handler.go)):

```go
func (h *TodoHandler) ListTodos(ctx context.Context, request gen.ListTodosRequestObject) (gen.ListTodosResponseObject, error) {
    userID, ok := auth.GetUserIDFromContext(ctx)
    if !ok {
        return gen.ListTodos401JSONResponse{Message: "Unauthorized"}, nil
    }
    // ...
}
```

## 新機能追加ガイド

### 新しいエンドポイントの追加

**手順概要**:

1. **CUE定義** ([openapi/cue/api.cue](openapi/cue/api.cue))
2. **コード生成** (`make generate`)
3. **ハンドラー実装** ([internal/handler/](internal/handler/))
4. **サービス実装** ([internal/service/](internal/service/))
5. **SQLクエリ追加**（必要な場合、[db/query/](db/query/)）
6. **Mapper実装**（必要な場合、[internal/mapper/](internal/mapper/))
7. **認証設定** ([internal/router/routes.go](internal/router/routes.go))

詳細は [development.md](development.md) を参照してください。

### テーブル/カラムの追加

**手順概要**:

1. **スキーマ編集** ([db/schema.sql](db/schema.sql))
2. **マイグレーション生成** (`make migrate-diff NAME=xxx`)
3. **マイグレーション適用** (`make migrate-apply`)
4. **SQLクエリ追加** ([db/query/](db/query/))
5. **sqlc再生成** (`make sqlc-generate`)
6. **サービス/ハンドラー更新**

詳細は [development.md](development.md) を参照してください。

## テスト戦略

### ユニットテスト（モック使用）

**Service層のテスト例** ([internal/service/todo.service_test.go](internal/service/todo.service_test.go)):

```go
func TestTodoService_GetTodoByID(t *testing.T) {
    mockRepo := mocks.NewTodoRepository(t)
    service := NewTodoService(mockRepo)

    // モックの期待値設定
    mockRepo.EXPECT().
        GetTodoByID(mock.Anything, sqlc.GetTodoByIDParams{ID: 1, UserID: 123}).
        Return(sqlc.Todo{ID: 1, Title: "Test"}, nil)

    // テスト実行
    todo, err := service.GetTodoByID(context.Background(), 1, 123)

    // アサーション
    assert.NoError(t, err)
    assert.Equal(t, int64(1), todo.ID)
    assert.Equal(t, "Test", todo.Title)
}
```

### インテグレーションテスト（実際のDB使用）

**UserService統合テスト例** ([internal/service/user.service_integration_test.go](internal/service/user.service_integration_test.go)):

```go
func TestUserService_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    // テスト用DB接続
    pool := setupTestDB(t)
    defer pool.Close()

    queries := sqlc.New(pool)
    service := NewUserService(queries, pool)

    // 実際のDB操作をテスト
    user, err := service.CreateOrUpdateUser(context.Background(), UserParams{
        Email: "test@example.com",
        Name:  "Test User",
    })

    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
    assert.Equal(t, "test@example.com", user.Email)
}
```

### テーブル駆動テスト

**パターン**:

```go
func TestTodoService_CreateTodo(t *testing.T) {
    tests := []struct {
        name        string
        title       string
        description *string
        wantErr     bool
    }{
        {
            name:        "valid todo",
            title:       "Test",
            description: ptrString("Description"),
            wantErr:     false,
        },
        {
            name:        "empty title",
            title:       "",
            description: nil,
            wantErr:     true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // テストケース実行
        })
    }
}
```

## なぜこのアーキテクチャなのか？

### 1. 二重の型安全性（OpenAPI + sqlc）

**HTTP層**:
- CUE → OpenAPI → oapi-codegen
- リクエスト/レスポンスが型付けされる
- 間違ったステータスコードやフィールド名 → コンパイルエラー

**データベース層**:
- SQL → sqlc
- クエリパラメータと結果が型付けされる
- 存在しないカラム名やSQL構文エラー → 生成時エラー

**結果**: ランタイムエラーの大幅削減

### 2. コンパイル時エラー検出

**通常のアプローチ**:

```go
// ランタイムまでエラーが検出されない
return c.JSON(200, map[string]interface{}{
    "mesage": "Success",  // typo: "message" が正しい
})
```

**このアーキテクチャ**:

```go
// コンパイル時にエラー検出
return gen.CreateTodo201JSONResponse{
    Mesage: "Success",  // ← コンパイルエラー: フィールド "Mesage" は存在しない
}
```

### 3. 各レイヤーの独立性とテスタビリティ

**Repository インターフェース**:
- Service層はRepositoryに依存（具象型ではない）
- テスト時はモックに差し替え可能
- DB実装の変更がService層に影響しない

**Mapper**:
- DB型とAPI型が分離
- スキーマ変更がAPIに影響しない
- API変更がDBに影響しない

### 4. コード生成による一貫性と保守性

**手動実装の問題**:
- ルーティング定義とハンドラーの不一致
- リクエスト/レスポンス構造のドキュメントとの齟齬
- 人的ミスによるバグ

**コード生成のメリット**:
- OpenAPI仕様が真実の源
- 自動的にコードとドキュメントが同期
- 変更が一箇所（CUE定義）で済む

### 5. 開発者体験の向上

**新機能追加時**:

1. CUE定義を追加
2. `make generate`
3. コンパイラが必要な実装を教えてくれる
   ```
   undefined: ListItems
   ```
4. ハンドラーを実装
5. 型チェックが通れば、基本的に動作する

**リファクタリング時**:

- 型が変更されると、コンパイラが影響箇所を全て教えてくれる
- 「この変更で何が壊れるか」が明確

---

## 参考リンク

- **OpenAPI**: [OpenAPI Specification](https://spec.openapis.org/oas/v3.0.3)
- **CUE**: [CUE Language](https://cuelang.org/)
- **oapi-codegen**: [oapi-codegen GitHub](https://github.com/oapi-codegen/oapi-codegen)
- **sqlc**: [sqlc Documentation](https://docs.sqlc.dev/)
- **Atlas**: [Atlas Documentation](https://atlasgo.io/docs)
- **Echo**: [Echo Framework](https://echo.labstack.com/)
- **pgx**: [pgx Documentation](https://github.com/jackc/pgx)
- **Goth**: [Goth GitHub](https://github.com/markbates/goth)
- **mockery**: [mockery GitHub](https://github.com/vektra/mockery)

---

より詳細な開発手順については、[development.md](development.md) を参照してください。
