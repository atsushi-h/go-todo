# OpenAPI コード生成フロー

このドキュメントでは、CUE言語からOpenAPI仕様を生成し、oapi-codegenでGoコードを生成する手順を説明します。

## 目次

- [概要](#概要)
- [ディレクトリ構造](#ディレクトリ構造)
- [生成フロー](#生成フロー)
- [コマンド](#コマンド)
- [ハンドラーの実装](#ハンドラーの実装)
- [StrictServer形式について](#strictserver形式について)
- [新しいエンドポイントの追加手順](#新しいエンドポイントの追加手順)
- [注意事項](#注意事項)

## 概要

```
CUE定義 → OpenAPI YAML → Go コード
```

- **CUE**: OpenAPIスキーマを型安全に記述
- **oapi-codegen**: OpenAPI仕様からGoのサーバーコードを生成
- **StrictServer**: 型付きリクエスト/レスポンスによる安全なハンドラー実装

## ディレクトリ構造

```
apps/api/
├── openapi/
│   ├── cue/
│   │   └── api.cue          # OpenAPI定義（CUE）
│   ├── oapi-codegen.yaml    # oapi-codegen設定
│   └── openapi.yaml         # 生成されたOpenAPI仕様
├── internal/
│   ├── gen/
│   │   └── api.gen.go       # 生成されたGoコード
│   └── handler/
│       ├── api.handler.go   # StrictServerInterface実装
│       ├── todo.handler.go  # Todo関連ハンドラー
│       └── auth.handler.go  # 認証ハンドラー（手動実装）
```

## 生成フロー

```
1. CUE定義を編集 (openapi/cue/api.cue)
   ↓
2. OpenAPI YAMLを生成 (make openapi-gen)
   ↓
3. Goコードを生成 (make api-gen)
   ↓
4. ハンドラーを実装/更新
   ↓
5. 変更をGitにコミット
```

## コマンド

### OpenAPI YAML生成

```bash
make openapi-gen
```

CUE定義から `openapi/openapi.yaml` を生成します。

### Goコード生成

```bash
make api-gen
```

OpenAPI仕様から `internal/gen/api.gen.go` を生成します。

### 両方を一括実行

```bash
make generate
```

`openapi-gen` と `api-gen` を順番に実行します。

## ハンドラーの実装

### 生成されるコードと手動実装の分離

| ファイル | 生成/手動 | 説明 |
|---------|----------|------|
| `internal/gen/api.gen.go` | 生成 | インターフェース、型定義、ルーティング |
| `internal/handler/api.handler.go` | 手動 | StrictServerInterface実装 |
| `internal/handler/todo.handler.go` | 手動 | Todo CRUD ロジック |
| `internal/handler/auth.handler.go` | 手動 | 認証ロジック（OAuth, セッション管理） |

### APIHandler構造

```go
type APIHandler struct {
    todoHandler *TodoHandler
}

// StrictServerInterfaceを実装
func (h *APIHandler) GetInfo(ctx context.Context, request gen.GetInfoRequestObject) (gen.GetInfoResponseObject, error) {
    return gen.GetInfo200JSONResponse{
        Name:    "Todo API",
        Version: "1.0.0",
    }, nil
}

// Todo関連はTodoHandlerに委譲
func (h *APIHandler) ListTodos(ctx context.Context, request gen.ListTodosRequestObject) (gen.ListTodosResponseObject, error) {
    return h.todoHandler.ListTodos(ctx, request)
}
```

### 認証ハンドラーについて

`auth.handler.go` は**手動で実装**しています。理由:

1. **OAuth認証フロー**: Gothic/Gothを使用したOAuthプロバイダー連携が必要
2. **セッション管理**: gorilla/sessionsを使用した独自のセッション管理
3. **リダイレクト処理**: 認証後のフロントエンドへのリダイレクト

これらはOpenAPIで定義しにくい複雑な処理のため、生成コードに含めず手動で実装しています。

```go
// auth.handler.go - 手動実装
func (h *AuthHandler) BeginAuth(c echo.Context) error { ... }
func (h *AuthHandler) Callback(c echo.Context) error { ... }
func (h *AuthHandler) Logout(c echo.Context) error { ... }
func (h *AuthHandler) Me(c echo.Context) error { ... }
```

ルーティングも `auth.route.go` で手動設定:

```go
authGroup := e.Group("/auth")
authGroup.GET("/:provider", authHandler.BeginAuth)
authGroup.GET("/:provider/callback", authHandler.Callback)
e.POST("/logout", authHandler.Logout)
e.GET("/me", authHandler.Me, auth.RequireAuth(sm))
```

## StrictServer形式について

oapi-codegenの `strict-server: true` オプションにより、型安全なハンドラーが生成されます。

### 通常形式との比較

**通常形式** (`strict-server: false`):
```go
func (h *Handler) GetTodo(ctx echo.Context, id int) error {
    todo, err := h.service.GetByID(id)
    if err != nil {
        return ctx.JSON(404, map[string]string{"error": "not found"})
    }
    return ctx.JSON(200, todo)
}
```

**StrictServer形式** (`strict-server: true`):
```go
func (h *Handler) GetTodo(ctx context.Context, request gen.GetTodoRequestObject) (gen.GetTodoResponseObject, error) {
    todo, err := h.service.GetByID(request.Id)
    if err != nil {
        return gen.GetTodo404JSONResponse{Message: "not found"}, nil
    }
    return gen.GetTodo200JSONResponse(*todo), nil
}
```

### メリット

1. **型安全性**: リクエストパラメータとレスポンスが型付けされる
2. **コンパイル時エラー**: 間違ったレスポンス型を返すとコンパイルエラー
3. **OpenAPIとの一貫性**: 定義した全レスポンスコードに対応する型が生成される

### 生成される型の例

```go
// リクエストオブジェクト
type GetTodoRequestObject struct {
    Id int `json:"id"`
}

// レスポンスインターフェース
type GetTodoResponseObject interface {
    VisitGetTodoResponse(w http.ResponseWriter) error
}

// 各ステータスコードの具象型
type GetTodo200JSONResponse Todo
type GetTodo404JSONResponse ErrorResponse
type GetTodo500JSONResponse ErrorResponse
```

## 新しいエンドポイントの追加手順

### 1. CUE定義を編集

```cue
// openapi/cue/api.cue
paths: {
    "/items": {
        get: {
            operationId: "ListItems"
            summary: "List all items"
            responses: {
                "200": {
                    description: "Success"
                    content: "application/json": schema: type: "array", items: {$ref: "#/components/schemas/Item"}
                }
            }
        }
    }
}

components: schemas: {
    Item: {
        type: "object"
        required: ["id", "name"]
        properties: {
            id:   {type: "integer", format: "int64"}
            name: {type: "string"}
        }
    }
}
```

### 2. コード生成

```bash
make generate
```

### 3. ハンドラー実装

```go
// internal/handler/item.handler.go
func (h *ItemHandler) ListItems(ctx context.Context, request gen.ListItemsRequestObject) (gen.ListItemsResponseObject, error) {
    items, err := h.service.GetAll(ctx)
    if err != nil {
        return gen.ListItems500JSONResponse{Message: err.Error()}, nil
    }
    // 型変換して返却
    result := make(gen.ListItems200JSONResponse, len(items))
    for i, item := range items {
        result[i] = gen.Item{Id: item.ID, Name: item.Name}
    }
    return result, nil
}
```

### 4. APIHandlerに追加

```go
// internal/handler/api.handler.go
type APIHandler struct {
    todoHandler *TodoHandler
    itemHandler *ItemHandler  // 追加
}

func (h *APIHandler) ListItems(ctx context.Context, request gen.ListItemsRequestObject) (gen.ListItemsResponseObject, error) {
    return h.itemHandler.ListItems(ctx, request)
}
```

## 注意事項

### 生成ファイルの編集禁止

`internal/gen/api.gen.go` は自動生成ファイルです。直接編集しないでください。
CUE定義を変更して再生成してください。

### operationIdの命名規則

operationIdはPascalCaseで記述してください。これがGoのメソッド名になります。

```cue
// Good
operationId: "ListTodos"
operationId: "GetTodoById"

// Bad
operationId: "list_todos"
operationId: "getTodoById"
```

### 認証ミドルウェアとoperationID

認証不要なエンドポイントは `routes.go` の `createAuthMiddleware` で除外設定が必要です:

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

### 生成ファイルのコミット

生成されたファイルはGitにコミットしてください:

```bash
git add apps/api/openapi/openapi.yaml
git add apps/api/internal/gen/api.gen.go
git commit -m "feat: Add new endpoint"
```

## 参考リンク

- [CUE Language](https://cuelang.org/)
- [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen)
- [OpenAPI Specification](https://spec.openapis.org/oas/v3.0.3)
- [Echo Framework](https://echo.labstack.com/)
