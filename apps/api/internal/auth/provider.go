package auth

import (
	"go-todo/internal/config"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

func InitProviders(oauthConfig config.OAuthConfig) {
	goth.UseProviders(
		google.New(
			oauthConfig.GoogleClientID,
			oauthConfig.GoogleClientSecret,
			oauthConfig.CallbackURL,
			"email", "profile",
		),
	)
}
