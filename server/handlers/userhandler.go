package handlers

import (
	"passport-v4/model"
	"passport-v4/utils"
	"passport-v4/utils/resp"

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

	successUrl := c.Query("success")
	if successUrl != "" {
		// 302 会导致Cookie丢失
		c.HTML(200, "redirect.html", gin.H{
			"href": successUrl,
		})
	} else {
		resp.Json(c, nil)
	}
}

func GetProfile(c *gin.Context) {
	ss := sessions.Default(c)
	user, ok := ss.Get(utils.SessionKeyUser).(model.User)
	if !ok {
		logrus.Errorf("User not logined!")
		resp.Json(c, gin.H{
			"logined": false,
		})
		return
	}
	logrus.Infof("getting User: %s %s", user.Name, user.ZjuId)
	resp.Json(c, gin.H{
		"logined": true,
		"User":    user,
	})
}
