package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

var Endpoint = oauth2.Endpoint{
	AuthURL:   "https://discord.com/api/oauth2/authorize",
	TokenURL:  "https://discord.com/api/oauth2/token",
	AuthStyle: oauth2.AuthStyleInParams,
}

var conf = &oauth2.Config{
	Endpoint:     Endpoint,
	Scopes:       []string{},
	RedirectURL:  "http://localhost:3000/auth/callback",
	ClientID:     "id",
	ClientSecret: "secret",
}

func RegisterOauth2(router *gin.RouterGroup) {
	discord := router.Group("/discord")
	discord.GET("/redirect", HandlerOAuth2Redirect)
	discord.GET("/callback", HandlerOAuth2Callback)
}

func HandlerOAuth2Redirect(ctx *gin.Context) {
	state := "testing"
	ctx.Redirect(http.StatusTemporaryRedirect, conf.AuthCodeURL(state))
}

func HandlerOAuth2Callback(ctx *gin.Context) {
}
