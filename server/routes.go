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

	//根节点
	root := e.Group("/", midware.Response)
	root.GET("/ping", Ping)

	root.GET("/zju/login", handler.ZjuLoginRequest)
	root.GET("/zju/login_success", handler.ZjuOauthCodeReturn)
	root.POST("/qsc/login", handler.QscLogin)
	root.GET("/logout", handler.Logout)
	root.GET("/profile", handler.GetProfile)
}
