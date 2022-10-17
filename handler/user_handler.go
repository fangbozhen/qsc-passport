package handler

import (
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
	ss.Options(sessions.Options{
		MaxAge: -1,
	})
	ss.Save()

	success_url := c.Query("success")
	if success_url != "" {
		// 302 会导致Cookie丢失
		c.HTML(200, "redirect.html", gin.H{
			"href": success_url,
		})
	} else {
		resp.Json(c, nil)
	}
}

func GetProfile(c *gin.Context) {
	ss := sessions.Default(c)
	user, ok := ss.Get(SS_KEY_USER).(model.User)
	if !ok {
		logrus.Errorf("user not logined!")
		resp.Json(c, gin.H{
			"logined": false,
		})
		return
	}
	logrus.Infof("getting user: %s %s", user.Name, user.ZjuId)
	resp.Json(c, gin.H{
		"logined": true,
		"user":    user,
	})
}
