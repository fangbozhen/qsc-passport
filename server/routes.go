package server

import (
	"passport-v4/handler"
	"passport-v4/server/midware"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.String(200, "pong!")
}

func configRoutes(e *gin.Engine) {
	e.LoadHTMLFiles("handler/redirect.html")

	root := e.Group("/", midware.Response)
	// internal := root.Group("/", IPWhiteList())
	root.GET("/ping", Ping)

	root.GET("/zju/login", handler.ZjuLoginRequest)
	root.GET("/zju/login_success", handler.ZjuOauthCodeReturn)
	root.POST("/qsc/login", handler.QscLogin)
	// internal.POST("/qsc/set_password", handler.SetPassword)
	root.GET("/logout", handler.Logout)
	root.GET("/profile", handler.GetProfile)
}
