package midware

import (
	. "passport-v4/global"
	"passport-v4/service"
	"passport-v4/util/resp"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func CheckToken(c *gin.Context) {
	id, ok := sessions.Default(c).Get(KEY_ID).(string)
	if !ok {
		logrus.Infof("[CheckToken midware] TOKEN not in DB, rejected")
		resp.ERR(c, "TOKEN not in Server")
		c.Abort()
		return
	}
	_, ok = service.GetUserById(id)
	if !ok {
		logrus.Infof("[CheckToken midware] User not in DB, rejected")
		resp.ERR(c, "Cannot read user data")
		c.Abort()
		return
	}
	c.Next()
}
