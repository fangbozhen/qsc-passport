package server

import (
	"github.com/gin-gonic/gin"
)

func configRoutes(e *gin.Engine) {
	// //根节点
	// root := e.Group("/", midware.Response)

	// //通用接口，不需要token
	// common := root.Group("/", midware.RemoveUndefinedQuery)
	// common.GET("/login", handler.UserLogin)

	// //用户接口，需要token
	// user := root.Group("/user", midware.CheckToken)
	// user.GET("/get", handler.UserGet)
	// user.GET("/get_profile", handler.UserGetProfile)
	// user.GET("/set_profile", handler.UserSetProfile)
}
