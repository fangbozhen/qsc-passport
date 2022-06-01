package handler

import (
	. "passport-v4/global"
	"passport-v4/model"
	"passport-v4/util/resp"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {

	ss := sessions.Default(c)
	ss.Clear()
	ss.Save()

	resp.JSON(c, nil)
}

func GetProfile(c *gin.Context) {
	ss := sessions.Default(c)
	user, ok := ss.Get(SS_KEY_USER).(model.User)
	if !ok {
		resp.JSON(c, gin.H{
			"logined": false,
		})
	}
	resp.JSON(c, gin.H{
		"logined": true,
		"user":    user,
	})
}
