package auth

import (
	"net/http"

	"go-todo/internal/util"

	"github.com/markbates/goth/gothic"
)

func InitGothic(sm *SessionManager) {
	gothic.Store = sm.Store()

	gothic.GetProviderName = func(r *http.Request) (string, error) {
		provider := util.GetParam(r, "provider")
		if provider != "" {
			return provider, nil
		}
		return r.URL.Query().Get("provider"), nil
	}
}
