package handler

import (
	"passport-v4/database"
	. "passport-v4/global"
	"passport-v4/model"
	"passport-v4/util/resp"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Logout(c *gin.Context) {

	ss := sessions.Default(c)
	ss.Clear()
	ss.Save()

	resp.Json(c, nil)
}

func GetProfile(c *gin.Context) {
	ss := sessions.Default(c)
	user, ok := ss.Get(SS_KEY_USER).(model.User)
	if !ok {
		resp.Json(c, gin.H{
			"logined": false,
		})
		return
	}
	if user.LoginType == model.LT_ZJU {
		resp.Json(c, gin.H{
			"logined": true,
			"user":    user,
		})
		return
	}
	qscuser, err := database.FindByName(user)
	if err != nil {
		logrus.Errorf("err: %s", err.Error())
		resp.Err(c, resp.E_DATABASE_ERROR, "数据库查找失败")
		return
	}
	resp.Json(c, gin.H{
		"logined": true,
		"user":    qscuser,
	})
}
