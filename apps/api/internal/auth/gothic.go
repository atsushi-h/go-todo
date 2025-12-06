package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

func InitGothic(sm *SessionManager) {
	gothic.Store = sm.Store()

	// gothicはクエリパラメータからproviderを取得する
	gothic.GetProviderName = func(r *http.Request) (string, error) {
		provider := r.URL.Query().Get("provider")
		if provider != "" {
			return provider, nil
		}
		return "", nil
	}
}

// Echoのパスパラメータをクエリパラメータに設定する
// gothicがproviderを取得できるようにするためのヘルパー
func SetProviderToRequest(c echo.Context) {
	provider := c.Param("provider")
	q := c.Request().URL.Query()
	q.Set("provider", provider)
	c.Request().URL.RawQuery = q.Encode()
}
