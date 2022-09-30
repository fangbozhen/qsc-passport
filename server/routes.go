package server

import (
	"net/http"
	"passport-v4/handler"
	"passport-v4/server/midware"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.String(200, "pong!")
}

func IPWhiteList(whitelist map[string]bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !whitelist[c.ClientIP()] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":  http.StatusForbidden,
				"message": "Access Denied",
			})
			return
		}
	}
}

func configRoutes(e *gin.Engine) {
	e.LoadHTMLFiles("handler/redirect.html")
	whitelist := make(map[string]bool)
	whitelist["127.0.0.1"] = true

	e.Use(IPWhiteList(whitelist))
	//根节点
	root := e.Group("/", midware.Response)
	root.GET("/ping", Ping)

	root.GET("/zju/login", handler.ZjuLoginRequest)
	root.GET("/zju/login_success", handler.ZjuOauthCodeReturn)
	root.POST("/qsc/login", handler.QscLogin)
	root.GET("/qsc/set_password", handler.SetPassword)
	root.GET("/logout", handler.Logout)
	root.GET("/profile", handler.GetProfile)
}
