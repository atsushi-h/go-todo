package util

import "net/http"

// コンテキストキーの型
type ContextKey string

// パスパラメータを取得
func GetParam(r *http.Request, name string) string {
	value := r.Context().Value(ContextKey(name))
	if value == nil {
		return ""
	}
	return value.(string)
}
