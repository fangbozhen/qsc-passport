package midware

import (
	. "passport-v4/global"
	"passport-v4/util/resp"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func CheckToken(c *gin.Context) {

	ck, err := c.Request.Cookie(KEY_TOKEN)
	if err != nil {
		logrus.Infof("[CheckToken midware] TOKEN not in cookie, rejected")
		resp.ERR(c, "TOKEN not in cookie, login first")
		c.Abort()
		return
	}
	client_token := ck.Value

	db_token, ok := sessions.Default(c).Get(KEY_TOKEN).(string)
	if !ok {
		logrus.Infof("[CheckToken midware] TOKEN not in DB, rejected")
		resp.ERR(c, "TOKEN not in Server")
		c.Abort()
		return
	}

	if client_token != db_token {
		logrus.Infof("[CheckToken midware] TOKEN incorrect, rejected")
		resp.ERR(c, "TOKEN incorrect")
		c.Abort()
		return
	}

	c.Next()
}
