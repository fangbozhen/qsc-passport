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

	root.GET("/zju/login", handler.ZJU_LoginRequest)
	root.GET("/zju/login_success", handler.ZJU_OauthCodeReturn)
	root.POST("/qsc/login", handler.QSC_Login)
	root.GET("/logout", handler.Logout)
	root.GET("/profile", handler.GetProfile)
}
