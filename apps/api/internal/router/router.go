package router

import (
	"context"
	"net/http"
	"regexp"

	"go-todo/internal/util"
)

// ルート定義
type Route struct {
	Method  string
	Pattern *regexp.Regexp
	Handler http.HandlerFunc
	Params  []string // パラメータ名のリスト
}

// カスタムルーター
type Router struct {
	routes            []*Route
	middlewares       []Middleware
	globalMiddlewares []Middleware // ルートマッチング前に適用されるミドルウェア
	notFound          http.HandlerFunc
}

// ミドルウェア関数の型
type Middleware func(http.HandlerFunc) http.HandlerFunc

// 新しいルーターを作成
func NewRouter() *Router {
	return &Router{
		routes: []*Route{},
		notFound: func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Not Found", http.StatusNotFound)
		},
	}
}

// ルートを登録
func (rt *Router) Handle(method, pattern string, handler http.HandlerFunc) {
	// パターンを正規表現に変換
	// /todos/{id} -> ^/todos/([^/]+)$
	params := []string{}
	regexPattern := regexp.MustCompile(`\{([^}]+)\}`)

	// パラメータ名を抽出
	matches := regexPattern.FindAllStringSubmatch(pattern, -1)
	for _, match := range matches {
		params = append(params, match[1])
	}

	// {param} を正規表現パターンに置き換え
	pattern = regexPattern.ReplaceAllString(pattern, `([^/]+)`)
	pattern = "^" + pattern + "$"

	rt.routes = append(rt.routes, &Route{
		Method:  method,
		Pattern: regexp.MustCompile(pattern),
		Handler: handler,
		Params:  params,
	})
}

// GETメソッドのルートを登録
func (rt *Router) GET(pattern string, handler http.HandlerFunc) {
	rt.Handle(http.MethodGet, pattern, handler)
}

// POSTメソッドのルートを登録
func (rt *Router) POST(pattern string, handler http.HandlerFunc) {
	rt.Handle(http.MethodPost, pattern, handler)
}

// PUTメソッドのルートを登録
func (rt *Router) PUT(pattern string, handler http.HandlerFunc) {
	rt.Handle(http.MethodPut, pattern, handler)
}

// DELETEメソッドのルートを登録
func (rt *Router) DELETE(pattern string, handler http.HandlerFunc) {
	rt.Handle(http.MethodDelete, pattern, handler)
}

// ミドルウェアを追加（ルートマッチング後に適用）
func (rt *Router) Use(middleware Middleware) {
	rt.middlewares = append(rt.middlewares, middleware)
}

// グローバルミドルウェアを追加（ルートマッチング前に適用、CORS等に使用）
func (rt *Router) UseGlobal(middleware Middleware) {
	rt.globalMiddlewares = append(rt.globalMiddlewares, middleware)
}

// 404ハンドラーを設定
func (rt *Router) SetNotFound(handler http.HandlerFunc) {
	rt.notFound = handler
}

// http.Handlerインターフェースを実装
func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// グローバルミドルウェアでラップした内部ハンドラーを作成
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// ルートをマッチング
		for _, route := range rt.routes {
			// メソッドチェック
			if route.Method != r.Method {
				continue
			}

			// パスマッチング
			matches := route.Pattern.FindStringSubmatch(path)
			if matches == nil {
				continue
			}

			// パラメータをコンテキストに格納
			if len(route.Params) > 0 {
				ctx := r.Context()
				for i, param := range route.Params {
					ctx = context.WithValue(ctx, util.ContextKey(param), matches[i+1])
				}
				r = r.WithContext(ctx)
			}

			// ルートミドルウェアを適用
			h := route.Handler
			for i := len(rt.middlewares) - 1; i >= 0; i-- {
				h = rt.middlewares[i](h)
			}

			h(w, r)
			return
		}

		// マッチするルートが見つからない
		rt.notFound(w, r)
	})

	// グローバルミドルウェアを適用
	for i := len(rt.globalMiddlewares) - 1; i >= 0; i-- {
		handler = rt.globalMiddlewares[i](handler)
	}

	handler(w, r)
}

// パスパラメータを取得（util.GetParamのエイリアス）
func Param(r *http.Request, name string) string {
	return util.GetParam(r, name)
}
