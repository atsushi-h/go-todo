package router

import (
	"go-todo/internal/auth"
	"go-todo/internal/handler"
)

func SetupAuthRoutes(r *Router, authHandler *handler.AuthHandler, sm *auth.SessionManager) {
	r.GET("/auth/{provider}", authHandler.BeginAuth)
	r.GET("/auth/{provider}/callback", authHandler.Callback)
	r.POST("/logout", authHandler.Logout)
	
	// ミドルウェアにSessionManagerを渡す
	r.GET("/me", auth.RequireAuth(sm)(authHandler.Me))
}
