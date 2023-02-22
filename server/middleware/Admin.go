package middleware

import (
	"passport-v4/model"
	"passport-v4/utils"
	"passport-v4/utils/resp"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func AdminCheck(c *gin.Context) {
	ss := sessions.Default(c)
	if ss.Get(utils.SessionKeyUser) == nil {
		resp.Err(c, resp.AuthFailedError, "未登录！")
		c.Abort()
		return
	}
	user := ss.Get(utils.SessionKeyUser).(model.User)
	if user.LoginType != model.LoginQSC || user.QscUser == nil {
		log.Warn("LoginType Error")
		resp.Err(c, resp.AuthFailedError, "not qscer")
		c.Abort()
		return
	}
	if user.QscUser.Position != model.PosManager && user.QscUser.Position != model.PosMaster {
		resp.Err(c, resp.AuthFailedError, "不是管理员")
		c.Abort()
		return
	}
	log.Infof("Admin has logined: %s", user.QscUser.QscId)
	c.Next()
}
